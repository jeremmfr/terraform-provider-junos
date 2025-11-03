package provider

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &vstpVlanGroup{}
	_ resource.ResourceWithConfigure      = &vstpVlanGroup{}
	_ resource.ResourceWithValidateConfig = &vstpVlanGroup{}
	_ resource.ResourceWithImportState    = &vstpVlanGroup{}
)

type vstpVlanGroup struct {
	client *junos.Client
}

func newVstpVlanGroupResource() resource.Resource {
	return &vstpVlanGroup{}
}

func (rsc *vstpVlanGroup) typeName() string {
	return providerName + "_vstp_vlan_group"
}

func (rsc *vstpVlanGroup) junosName() string {
	return "vstp vlan-group group"
}

func (rsc *vstpVlanGroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *vstpVlanGroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *vstpVlanGroup) Configure(
	ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedResourceConfigureType(ctx, req, resp)

		return
	}
	rsc.client = client
}

func (rsc *vstpVlanGroup) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			Description: "An identifier for the resource with format " +
				"`<name>" + junos.IDSeparator + "<routing_instance>`.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Required:    true,
			Description: "VLAN group name.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				tfvalidator.StringFormat(tfvalidator.DefaultFormat),
			},
		},
		"routing_instance": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(junos.DefaultW),
			Description: "Routing instance for vstp protocol if not root level.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				tfvalidator.StringFormat(tfvalidator.DefaultFormat),
			},
		},
		"vlan": schema.SetAttribute{
			ElementType: types.StringType,
			Required:    true,
			Description: " VLAN IDs or VLAN ID ranges.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9])`+
							`(-(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9]))?$`),
						"must be a VLAN id (1..4094) or a range of VLAN id (1..4094)-(1..4094)"),
				),
			},
		},
	}
	maps.Copy(attributes, vstpVlanAttrData{}.attributesSchema())

	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes:  attributes,
	}
}

type vstpVlanGroupData struct {
	vstpVlanAttrData

	ID              types.String   `tfsdk:"id"`
	Name            types.String   `tfsdk:"name"`
	RoutingInstance types.String   `tfsdk:"routing_instance"`
	Vlan            []types.String `tfsdk:"vlan"`
}

type vstpVlanGroupConfig struct {
	vstpVlanAttrData

	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
	Vlan            types.Set    `tfsdk:"vlan"`
}

func (rsc *vstpVlanGroup) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config vstpVlanGroupConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.vstpVlanAttrData.validateConfig(ctx, resp)
}

func (rsc *vstpVlanGroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan vstpVlanGroupData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Empty Name",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "name"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
				instanceExists, err := checkRoutingInstanceExists(fnCtx, v, junSess)
				if err != nil {
					resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

					return false
				}
				if !instanceExists {
					resp.Diagnostics.AddAttributeError(
						path.Root("routing_instance"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("routing instance %q doesn't exist", v),
					)

					return false
				}
			}
			groupExists, err := checkVstpVlanGroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if groupExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsInRoutingInstanceMessage(rsc, plan.Name, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsMessage(rsc, plan.Name),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			groupExists, err := checkVstpVlanGroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !groupExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(rsc, plan.Name, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
					)
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *vstpVlanGroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data vstpVlanGroupData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *vstpVlanGroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state vstpVlanGroupData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataDelWithOpts = &state
	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *vstpVlanGroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state vstpVlanGroupData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceDelete(
		ctx,
		rsc,
		&state,
		resp,
	)
}

func (rsc *vstpVlanGroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data vstpVlanGroupData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <name>"+junos.IDSeparator+"<routing_instance>)",
	)
}

func checkVstpVlanGroupExists(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols vstp vlan-group group " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *vstpVlanGroupData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *vstpVlanGroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *vstpVlanGroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols vstp vlan-group group " + rscData.Name.ValueString() + " "

	for _, v := range rscData.Vlan {
		configSet = append(configSet, setPrefix+"vlan "+v.ValueString())
	}
	configSet = append(configSet, rscData.vstpVlanAttrData.configSet(setPrefix)...)

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *vstpVlanGroupData) read(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols vstp vlan-group group " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		if routingInstance == "" {
			rscData.RoutingInstance = types.StringValue(junos.DefaultW)
		} else {
			rscData.RoutingInstance = types.StringValue(routingInstance)
		}
		rscData.fillID()
		for item := range strings.SplitSeq(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "vlan "):
				rscData.Vlan = append(rscData.Vlan, types.StringValue(itemTrim))
			default:
				if err := rscData.vstpVlanAttrData.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *vstpVlanGroupData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "protocols vstp vlan-group group " + rscData.Name.ValueString() + " "

	configSet := append(
		[]string{
			delPrefix + "vlan",
		},
		rscData.vstpVlanAttrData.configOptsToDel(delPrefix)...,
	)

	return junSess.ConfigSet(configSet)
}

func (rscData *vstpVlanGroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "protocols vstp vlan-group group " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
