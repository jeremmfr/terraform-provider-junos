package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &snmpCommunity{}
	_ resource.ResourceWithConfigure      = &snmpCommunity{}
	_ resource.ResourceWithValidateConfig = &snmpCommunity{}
	_ resource.ResourceWithImportState    = &snmpCommunity{}
)

type snmpCommunity struct {
	client *junos.Client
}

func newSnmpCommunityResource() resource.Resource {
	return &snmpCommunity{}
}

func (rsc *snmpCommunity) typeName() string {
	return providerName + "_snmp_community"
}

func (rsc *snmpCommunity) junosName() string {
	return "snmp community"
}

func (rsc *snmpCommunity) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *snmpCommunity) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *snmpCommunity) Configure(
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

func (rsc *snmpCommunity) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of snmp community.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"authorization_read_only": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow read-only access.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"authorization_read_write": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow read and write access.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"client_list_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of client list or prefix list.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"clients": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of source address prefix ranges to accept.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringCIDR(),
					),
				},
			},
			"view": schema.StringAttribute{
				Optional:    true,
				Description: "View name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"routing_instance": schema.ListNestedBlock{
				Description: "For each name of routing instance, accept clients.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of routing instance.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
								stringvalidator.NoneOfCaseInsensitive(junos.DefaultW),
							},
						},
						"client_list_name": schema.StringAttribute{
							Optional:    true,
							Description: "The name of client list or prefix list.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 128),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"clients": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of source address prefix ranges to accept.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									tfvalidator.StringCIDR(),
								),
							},
						},
					},
				},
			},
		},
	}
}

type snmpCommunityData struct {
	ID                     types.String                        `tfsdk:"id"`
	Name                   types.String                        `tfsdk:"name"`
	AuthorizationReadOnly  types.Bool                          `tfsdk:"authorization_read_only"`
	AuthorizationReadWrite types.Bool                          `tfsdk:"authorization_read_write"`
	ClientListName         types.String                        `tfsdk:"client_list_name"`
	Clients                []types.String                      `tfsdk:"clients"`
	View                   types.String                        `tfsdk:"view"`
	RoutingInstance        []snmpCommunityBlockRoutingInstance `tfsdk:"routing_instance"`
}

type snmpCommunityConfig struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	AuthorizationReadOnly  types.Bool   `tfsdk:"authorization_read_only"`
	AuthorizationReadWrite types.Bool   `tfsdk:"authorization_read_write"`
	ClientListName         types.String `tfsdk:"client_list_name"`
	Clients                types.Set    `tfsdk:"clients"`
	View                   types.String `tfsdk:"view"`
	RoutingInstance        types.List   `tfsdk:"routing_instance"`
}

type snmpCommunityBlockRoutingInstance struct {
	Name           types.String   `tfsdk:"name"`
	ClientListName types.String   `tfsdk:"client_list_name"`
	Clients        []types.String `tfsdk:"clients"`
}

type snmpCommunityBlockRoutingInstanceConfig struct {
	Name           types.String `tfsdk:"name"`
	ClientListName types.String `tfsdk:"client_list_name"`
	Clients        types.Set    `tfsdk:"clients"`
}

func (rsc *snmpCommunity) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config snmpCommunityConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.AuthorizationReadOnly.IsNull() && !config.AuthorizationReadOnly.IsUnknown() &&
		!config.AuthorizationReadWrite.IsNull() && !config.AuthorizationReadWrite.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authorization_read_only"),
			tfdiag.ConflictConfigErrSummary,
			"authorization_read_only and authorization_read_write cannot be configured together",
		)
	}
	if !config.ClientListName.IsNull() && !config.ClientListName.IsUnknown() &&
		!config.Clients.IsNull() && !config.Clients.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_list_name"),
			tfdiag.ConflictConfigErrSummary,
			"client_list_name and clients cannot be configured together",
		)
	}
	if !config.RoutingInstance.IsNull() && !config.RoutingInstance.IsUnknown() {
		var routingInstance []snmpCommunityBlockRoutingInstanceConfig
		asDiags := config.RoutingInstance.ElementsAs(ctx, &routingInstance, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		routingInstanceName := make(map[string]struct{})
		for i, block := range routingInstance {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := routingInstanceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("routing_instance").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple routing_instance blocks with the same name %q", name),
					)
				}
				routingInstanceName[name] = struct{}{}
			}

			if !block.ClientListName.IsNull() && !block.ClientListName.IsUnknown() &&
				!block.Clients.IsNull() && !block.Clients.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("routing_instance").AtListIndex(i).AtName("client_list_name"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("client_list_name and clients cannot be configured together"+
						" in routing_instance block %q", block.Name.ValueString()),
				)
			}
		}
	}
}

func (rsc *snmpCommunity) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan snmpCommunityData
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
			communityExists, err := checkSnmpCommunityExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if communityExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			communityExists, err := checkSnmpCommunityExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !communityExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *snmpCommunity) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data snmpCommunityData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *snmpCommunity) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state snmpCommunityData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *snmpCommunity) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state snmpCommunityData
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

func (rsc *snmpCommunity) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data snmpCommunityData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "name"),
	)
}

func checkSnmpCommunityExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp community \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *snmpCommunityData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *snmpCommunityData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *snmpCommunityData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set snmp community \"" + rscData.Name.ValueString() + "\" "

	if rscData.AuthorizationReadOnly.ValueBool() {
		configSet = append(configSet, setPrefix+"authorization read-only")
	}
	if rscData.AuthorizationReadWrite.ValueBool() {
		configSet = append(configSet, setPrefix+"authorization read-write")
	}
	if v := rscData.ClientListName.ValueString(); v != "" {
		if len(rscData.Clients) > 0 {
			return path.Root("client_list_name"),
				errors.New("client_list_name and clients cannot be configured together")
		}
		configSet = append(configSet, setPrefix+"client-list-name \""+v+"\"")
	}
	for _, v := range rscData.Clients {
		configSet = append(configSet, setPrefix+"clients "+v.ValueString())
	}
	if v := rscData.View.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"view \""+v+"\"")
	}
	routingInstanceName := make(map[string]struct{})
	for i, block := range rscData.RoutingInstance {
		name := block.Name.ValueString()
		if len(block.Clients) > 0 && block.ClientListName.ValueString() != "" {
			return path.Root("routing_instance").AtListIndex(i).AtName("client_list_name"),
				fmt.Errorf("client_list_name and clients cannot be configured together"+
					" in routing_instance block %q", name)
		}
		if _, ok := routingInstanceName[name]; ok {
			return path.Root("routing_instance").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple routing_instance blocks with the same name %q", name)
		}
		routingInstanceName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"routing-instance "+name)
		if v := block.ClientListName.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"routing-instance "+name+" client-list-name \""+v+"\"")
		}
		for _, v := range block.Clients {
			configSet = append(configSet, setPrefix+"routing-instance "+name+" clients "+v.ValueString())
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *snmpCommunityData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp community \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
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
			case itemTrim == "authorization read-only":
				rscData.AuthorizationReadOnly = types.BoolValue(true)
			case itemTrim == "authorization read-write":
				rscData.AuthorizationReadWrite = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "client-list-name "):
				rscData.ClientListName = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "clients "):
				rscData.Clients = append(rscData.Clients, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "view "):
				rscData.View = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var routingInstance snmpCommunityBlockRoutingInstance
				rscData.RoutingInstance, routingInstance = tfdata.ExtractBlockWithTFTypesString(
					rscData.RoutingInstance, "Name", itemTrimFields[0],
				)
				routingInstance.Name = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "client-list-name "):
					routingInstance.ClientListName = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "clients "):
					routingInstance.Clients = append(routingInstance.Clients, types.StringValue(itemTrim))
				}
				rscData.RoutingInstance = append(rscData.RoutingInstance, routingInstance)
			}
		}
	}

	return nil
}

func (rscData *snmpCommunityData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete snmp community \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
