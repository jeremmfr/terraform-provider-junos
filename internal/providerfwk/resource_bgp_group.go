package providerfwk

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

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
	_ resource.Resource                   = &bgpGroup{}
	_ resource.ResourceWithConfigure      = &bgpGroup{}
	_ resource.ResourceWithModifyPlan     = &bgpGroup{}
	_ resource.ResourceWithValidateConfig = &bgpGroup{}
	_ resource.ResourceWithImportState    = &bgpGroup{}
	_ resource.ResourceWithUpgradeState   = &bgpGroup{}
)

type bgpGroup struct {
	client *junos.Client
}

func newBgpGroupResource() resource.Resource {
	return &bgpGroup{}
}

func (rsc *bgpGroup) typeName() string {
	return providerName + "_bgp_group"
}

func (rsc *bgpGroup) junosName() string {
	return "bgp group"
}

func (rsc *bgpGroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *bgpGroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *bgpGroup) Configure(
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

func (rsc *bgpGroup) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "An identifier for the resource with format `<name>" + junos.IDSeparator + "<routing_instance>`.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Required:    true,
			Description: "Name of group.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 250),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"routing_instance": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(junos.DefaultW),
			Description: "Routing instance for bgp protocol if not root level.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				tfvalidator.StringFormat(tfvalidator.DefaultFormat),
			},
		},
		"type": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("external"),
			Description: "Type of peer group.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("internal", "external"),
			},
		},
	}
	maps.Copy(attributes, bgpAttrData{}.attributesSchema())

	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes:  attributes,
		Blocks:      bgpAttrData{}.blocksSchema(),
	}
}

type bgpGroupData struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
	Type            types.String `tfsdk:"type"`
	bgpAttrData
}

type bgpGroupConfig struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
	Type            types.String `tfsdk:"type"`
	bgpAttrConfig
}

func (rsc *bgpGroup) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config bgpGroupConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.bgpAttrConfig.validateConfig(ctx, resp)
}

func (rsc *bgpGroup) ModifyPlan(
	ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var config, plan bgpGroupConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.bgpAttrConfig.modifyPlan(ctx, &plan.bgpAttrConfig)

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (rsc *bgpGroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan bgpGroupData
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

	if plan.AdvertiseExternal.IsUnknown() {
		plan.AdvertiseExternal = types.BoolNull()
		if plan.AdvertiseExternalConditional.ValueBool() {
			plan.AdvertiseExternal = types.BoolValue(true)
		}
	}
	if plan.MetricOutIgp.IsUnknown() {
		plan.MetricOutIgp = types.BoolNull()
		if plan.MetricOutIgpDelayMedUpdate.ValueBool() {
			plan.MetricOutIgp = types.BoolValue(true)
		}
		if !plan.MetricOutIgpOffset.IsNull() {
			plan.MetricOutIgp = types.BoolValue(true)
		}
	}
	if plan.MetricOutMinimumIgp.IsUnknown() {
		plan.MetricOutMinimumIgp = types.BoolNull()
		if !plan.MetricOutMinimumIgpOffset.IsNull() {
			plan.MetricOutMinimumIgp = types.BoolValue(true)
		}
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
			bgpGroupExists, err := checkBgpGroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if bgpGroupExists {
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
			bgpGroupExists, err := checkBgpGroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !bgpGroupExists {
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

func (rsc *bgpGroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data bgpGroupData
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

func (rsc *bgpGroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state bgpGroupData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.AdvertiseExternal.IsUnknown() {
		plan.AdvertiseExternal = types.BoolNull()
		if plan.AdvertiseExternalConditional.ValueBool() {
			plan.AdvertiseExternal = types.BoolValue(true)
		}
	}
	if plan.MetricOutIgp.IsUnknown() {
		plan.MetricOutIgp = types.BoolNull()
		if plan.MetricOutIgpDelayMedUpdate.ValueBool() {
			plan.MetricOutIgp = types.BoolValue(true)
		}
		if !plan.MetricOutIgpOffset.IsNull() {
			plan.MetricOutIgp = types.BoolValue(true)
		}
	}
	if plan.MetricOutMinimumIgp.IsUnknown() {
		plan.MetricOutMinimumIgp = types.BoolNull()
		if !plan.MetricOutMinimumIgpOffset.IsNull() {
			plan.MetricOutMinimumIgp = types.BoolValue(true)
		}
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

func (rsc *bgpGroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state bgpGroupData
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

func (rsc *bgpGroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data bgpGroupData

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

func checkBgpGroupExists(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols bgp group \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *bgpGroupData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *bgpGroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *bgpGroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols bgp group \"" + rscData.Name.ValueString() + "\" "

	configSet := []string{
		setPrefix + "type " + rscData.Type.ValueString(),
	}

	dataConfigSet, errPath, err := rscData.bgpAttrData.configSet(setPrefix)
	if err != nil {
		return errPath, err
	}
	configSet = append(configSet, dataConfigSet...)

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *bgpGroupData) read(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols bgp group \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "type "):
				rscData.Type = types.StringValue(itemTrim)
			default:
				if err := rscData.bgpAttrData.read(itemTrim, junSess); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *bgpGroupData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "protocols bgp group \"" + rscData.Name.ValueString() + "\" "

	return junSess.ConfigSet(rscData.bgpAttrData.configOptsToDel(delPrefix))
}

func (rscData *bgpGroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "protocols bgp group \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
