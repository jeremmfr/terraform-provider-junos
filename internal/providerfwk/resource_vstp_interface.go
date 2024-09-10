package providerfwk

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                   = &vstpInterface{}
	_ resource.ResourceWithConfigure      = &vstpInterface{}
	_ resource.ResourceWithValidateConfig = &vstpInterface{}
	_ resource.ResourceWithImportState    = &vstpInterface{}
)

type vstpInterface struct {
	client *junos.Client
}

func newVstpInterfaceResource() resource.Resource {
	return &vstpInterface{}
}

func (rsc *vstpInterface) typeName() string {
	return providerName + "_vstp_interface"
}

func (rsc *vstpInterface) junosName() string {
	return "vstp interface"
}

func (rsc *vstpInterface) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *vstpInterface) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *vstpInterface) Configure(
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

func (rsc *vstpInterface) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<name>" + junos.IDSeparator + junos.IDSeparator + "<routing_instance>`, " +
					"`<name>" + junos.IDSeparator + "v_<vlan>" + junos.IDSeparator + "<routing_instance>` or " +
					"`<name>" + junos.IDSeparator + "vg_<vlan_group>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Interface name or `all`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.StringDotExclusion(),
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
			"vlan": schema.StringAttribute{
				Optional:    true,
				Description: "Configure interface in VSTP vlan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9]|all)$`),
						"must be a VLAN id (1..4094) or all"),
				},
			},
			"vlan_group": schema.StringAttribute{
				Optional:    true,
				Description: "Configure interface in VSTP vlan-group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"access_trunk": schema.BoolAttribute{
				Optional:    true,
				Description: "Send/Receive untagged RSTP BPDUs on this interface.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"bpdu_timeout_action_alarm": schema.BoolAttribute{
				Optional:    true,
				Description: "Generate an alarm on BPDU expiry (Loop Protect).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"bpdu_timeout_action_block": schema.BoolAttribute{
				Optional:    true,
				Description: "Block the interface on BPDU expiry (Loop Protect).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"cost": schema.Int64Attribute{
				Optional:    true,
				Description: "Cost of the interface.",
				Validators: []validator.Int64{
					int64validator.Between(1, 200000000),
				},
			},
			"edge": schema.BoolAttribute{
				Optional:    true,
				Description: "Port is an edge port.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"mode": schema.StringAttribute{
				Optional:    true,
				Description: "Interface mode (P2P or shared).",
				Validators: []validator.String{
					stringvalidator.OneOf("point-to-point", "shared"),
				},
			},
			"no_root_port": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not allow the interface to become root (Root Protect).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"priority": schema.Int64Attribute{
				Optional:    true,
				Description: "Interface priority (in increments of 16).",
				Validators: []validator.Int64{
					int64validator.Between(0, 240),
				},
			},
		},
	}
}

type vstpInterfaceData struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	RoutingInstance        types.String `tfsdk:"routing_instance"`
	Vlan                   types.String `tfsdk:"vlan"`
	VlanGroup              types.String `tfsdk:"vlan_group"`
	AccessTrunk            types.Bool   `tfsdk:"access_trunk"`
	BpduTimeoutActionAlarm types.Bool   `tfsdk:"bpdu_timeout_action_alarm"`
	BpduTimeoutActionBlock types.Bool   `tfsdk:"bpdu_timeout_action_block"`
	Cost                   types.Int64  `tfsdk:"cost"`
	Edge                   types.Bool   `tfsdk:"edge"`
	Mode                   types.String `tfsdk:"mode"`
	NoRootPort             types.Bool   `tfsdk:"no_root_port"`
	Priority               types.Int64  `tfsdk:"priority"`
}

func (rsc *vstpInterface) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config vstpInterfaceData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Vlan.IsNull() && !config.Vlan.IsUnknown() &&
		!config.VlanGroup.IsNull() && !config.VlanGroup.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vlan"),
			tfdiag.ConflictConfigErrSummary,
			"vlan and vlan_group cannot be configured together",
		)
	}
	if !config.Priority.IsNull() && !config.Priority.IsUnknown() {
		if config.Priority.ValueInt64()%16 != 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("priority"),
				"Bad Value Error",
				"priority must be a multiple of 16",
			)
		}
	}
}

func (rsc *vstpInterface) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan vstpInterfaceData
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
	if !plan.Vlan.IsNull() && !plan.VlanGroup.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vlan"),
			tfdiag.ConflictConfigErrSummary,
			"vlan and vlan_group cannot be configured together",
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
			if v := plan.Vlan.ValueString(); v != "" {
				vlanExists, err := checkVstpVlanExists(fnCtx, v, plan.RoutingInstance.ValueString(), junSess)
				if err != nil {
					resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

					return false
				}
				if !vlanExists {
					if vRI := plan.RoutingInstance.ValueString(); vRI != "" && vRI != junos.DefaultW {
						resp.Diagnostics.AddAttributeError(
							path.Root("vlan"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("vstp vlan %q in routing-instance %q doesn't exist", v, vRI),
						)
					} else {
						resp.Diagnostics.AddAttributeError(
							path.Root("vlan"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("vstp vlan %q doesn't exist", v),
						)
					}

					return false
				}
			}
			if v := plan.VlanGroup.ValueString(); v != "" {
				groupExists, err := checkVstpVlanGroupExists(fnCtx, v, plan.RoutingInstance.ValueString(), junSess)
				if err != nil {
					resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

					return false
				}
				if !groupExists {
					if vRI := plan.RoutingInstance.ValueString(); vRI != "" && vRI != junos.DefaultW {
						resp.Diagnostics.AddAttributeError(
							path.Root("vlan_group"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("vstp vlan-group group %q in routing-instance %q doesn't exist", v, vRI),
						)
					} else {
						resp.Diagnostics.AddAttributeError(
							path.Root("vlan_group"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("vstp vlan-group group %q doesn't exist", v),
						)
					}

					return false
				}
			}
			interfaceExists, err := checkVstpInterfaceExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				plan.Vlan.ValueString(),
				plan.VlanGroup.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if interfaceExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					switch {
					case plan.Vlan.ValueString() != "":
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(rsc.junosName()+" %q already exists in routing-instance %q in vlan %q",
								plan.Name.ValueString(), v, plan.Vlan.ValueString(),
							),
						)
					case plan.VlanGroup.ValueString() != "":
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(rsc.junosName()+" %q already exists in routing-instance %q in vlan-group group %q",
								plan.Name.ValueString(), v, plan.VlanGroup.ValueString(),
							),
						)
					default:
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							defaultResourceAlreadyExistsInRoutingInstanceMessage(rsc, plan.Name, v),
						)
					}
				} else {
					switch {
					case plan.Vlan.ValueString() != "":
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(rsc.junosName()+" %q already exists in vlan %q",
								plan.Name.ValueString(), plan.Vlan.ValueString(),
							),
						)
					case plan.VlanGroup.ValueString() != "":
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(rsc.junosName()+" %q already exists in vlan-group group %q",
								plan.Name.ValueString(), plan.VlanGroup.ValueString(),
							),
						)
					default:
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							defaultResourceAlreadyExistsMessage(rsc, plan.Name),
						)
					}
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			interfaceExists, err := checkVstpInterfaceExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				plan.Vlan.ValueString(),
				plan.VlanGroup.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !interfaceExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					switch {
					case plan.Vlan.ValueString() != "":
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							fmt.Sprintf(rsc.junosName()+" %q does not exists in routing-instance %q in vlan %q after commit "+
								"=> check your config",
								plan.Name.ValueString(), v, plan.Vlan.ValueString()),
						)
					case plan.VlanGroup.ValueString() != "":
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							fmt.Sprintf(rsc.junosName()+" %q does not exists in routing-instance %q in vlan-group group %q after commit "+
								"=> check your config",
								plan.Name.ValueString(), v, plan.VlanGroup.ValueString()),
						)
					default:
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(rsc, plan.Name, v),
						)
					}
				} else {
					switch {
					case plan.Vlan.ValueString() != "":
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							fmt.Sprintf(rsc.junosName()+" %q does not exists in vlan %q after commit "+
								"=> check your config",
								plan.Name.ValueString(), plan.Vlan.ValueString()),
						)
					case plan.VlanGroup.ValueString() != "":
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							fmt.Sprintf(rsc.junosName()+" %q does not exists in vlan-group group %q after commit "+
								"=> check your config",
								plan.Name.ValueString(), plan.VlanGroup.ValueString()),
						)
					default:
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
						)
					}
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *vstpInterface) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data vstpInterfaceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom4String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.RoutingInstance.ValueString(),
			state.Vlan.ValueString(),
			state.VlanGroup.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *vstpInterface) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state vstpInterfaceData
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

func (rsc *vstpInterface) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state vstpInterfaceData
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

func (rsc *vstpInterface) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	idList := strings.Split(req.ID, junos.IDSeparator)
	var name, routingInstance, vlan, vlanGroup string
	switch len(idList) {
	case 1:
		name = idList[0]
	case 2:
		name = idList[0]
		routingInstance = idList[1]
	default:
		name = idList[0]
		routingInstance = idList[2]
		if balt.CutPrefixInString(&idList[1], "v_") {
			vlan = idList[1]
		} else if balt.CutPrefixInString(&idList[1], "vg_") {
			vlanGroup = idList[1]
		}
	}

	var data vstpInterfaceData
	if err := data.read(ctx, name, routingInstance, vlan, vlanGroup, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	if data.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			fmt.Sprintf("don't find "+rsc.junosName()+" with id %q"+
				" (id must be <name>"+junos.IDSeparator+junos.IDSeparator+"<routing_instance>, "+
				"<name>"+junos.IDSeparator+"v_<vlan>"+junos.IDSeparator+"<routing_instance> or "+
				"<name>"+junos.IDSeparator+"vg_<vlan_group>"+junos.IDSeparator+"<routing_instance>)",
				req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkVstpInterfaceExists(
	_ context.Context, name, routingInstance, vlan, vlanGroup string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showPrefix += "protocols vstp "
	if vlan != "" {
		showPrefix += "vlan " + vlan + " "
	} else if vlanGroup != "" {
		showPrefix += "vlan-group group " + vlanGroup + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"interface " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *vstpInterfaceData) fillID() {
	idPrefix := rscData.Name.ValueString() + junos.IDSeparator
	if rscData.Vlan.ValueString() != "" {
		idPrefix += "v_" + rscData.Vlan.ValueString()
	} else if rscData.VlanGroup.ValueString() != "" {
		idPrefix += "vg_" + rscData.VlanGroup.ValueString()
	}

	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(idPrefix + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(idPrefix + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *vstpInterfaceData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *vstpInterfaceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols vstp "
	if v := rscData.Vlan.ValueString(); v != "" {
		setPrefix += "vlan " + v + " "
	} else if v := rscData.VlanGroup.ValueString(); v != "" {
		setPrefix += "vlan-group group " + v + " "
	}
	setPrefix += "interface " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if rscData.AccessTrunk.ValueBool() {
		configSet = append(configSet, setPrefix+"access-trunk")
	}
	if rscData.BpduTimeoutActionAlarm.ValueBool() {
		configSet = append(configSet, setPrefix+"bpdu-timeout-action alarm")
	}
	if rscData.BpduTimeoutActionBlock.ValueBool() {
		configSet = append(configSet, setPrefix+"bpdu-timeout-action block")
	}
	if !rscData.Cost.IsNull() {
		configSet = append(configSet, setPrefix+"cost "+
			utils.ConvI64toa(rscData.Cost.ValueInt64()))
	}
	if rscData.Edge.ValueBool() {
		configSet = append(configSet, setPrefix+"edge")
	}
	if v := rscData.Mode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"mode "+v)
	}
	if rscData.NoRootPort.ValueBool() {
		configSet = append(configSet, setPrefix+"no-root-port")
	}
	if !rscData.Priority.IsNull() {
		configSet = append(configSet, setPrefix+"priority "+
			utils.ConvI64toa(rscData.Priority.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *vstpInterfaceData) read(
	_ context.Context, name, routingInstance, vlan, vlanGroup string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showPrefix += "protocols vstp "
	if vlan != "" {
		showPrefix += "vlan " + vlan + " "
	} else if vlanGroup != "" {
		showPrefix += "vlan-group group " + vlanGroup + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"interface " + name + junos.PipeDisplaySetRelative)
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
		if vlan != "" {
			rscData.Vlan = types.StringValue(vlan)
		} else if vlanGroup != "" {
			rscData.VlanGroup = types.StringValue(vlanGroup)
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
			case itemTrim == "access-trunk":
				rscData.AccessTrunk = types.BoolValue(true)
			case itemTrim == "bpdu-timeout-action alarm":
				rscData.BpduTimeoutActionAlarm = types.BoolValue(true)
			case itemTrim == "bpdu-timeout-action block":
				rscData.BpduTimeoutActionBlock = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "cost "):
				rscData.Cost, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "edge":
				rscData.Edge = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "mode "):
				rscData.Mode = types.StringValue(itemTrim)
			case itemTrim == "no-root-port":
				rscData.NoRootPort = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "priority "):
				rscData.Priority, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *vstpInterfaceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "protocols vstp "
	if v := rscData.Vlan.ValueString(); v != "" {
		delPrefix += "vlan " + v + " "
	} else if v := rscData.VlanGroup.ValueString(); v != "" {
		delPrefix += "vlan-group group " + v + " "
	}

	configSet := []string{
		delPrefix + "interface " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
