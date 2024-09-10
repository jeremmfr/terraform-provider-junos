package providerfwk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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
	_ resource.Resource                   = &rstp{}
	_ resource.ResourceWithConfigure      = &rstp{}
	_ resource.ResourceWithValidateConfig = &rstp{}
	_ resource.ResourceWithImportState    = &rstp{}
)

type rstp struct {
	client *junos.Client
}

func newRstpResource() resource.Resource {
	return &rstp{}
}

func (rsc *rstp) typeName() string {
	return providerName + "_rstp"
}

func (rsc *rstp) junosName() string {
	return "protocols rstp"
}

func (rsc *rstp) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *rstp) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *rstp) Configure(
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

func (rsc *rstp) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for rstp protocol if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"backup_bridge_priority": schema.StringAttribute{
				Optional:    true,
				Description: "Priority of the bridge.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d\d?k$`),
						"must be a number with increments of 4k - 4k,8k,..60k",
					),
				},
			},
			"bpdu_block_on_edge": schema.BoolAttribute{
				Optional:    true,
				Description: "Block BPDU on all interfaces configured as edge (BPDU Protect).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"bpdu_destination_mac_address_provider_bridge_group": schema.BoolAttribute{
				Optional:    true,
				Description: "Destination MAC address in the spanning tree BPDUs is 802.1ad provider bridge group address.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"bridge_priority": schema.StringAttribute{
				Optional:    true,
				Description: "Priority of the bridge.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(0|\d\d?k)$`),
						"must be a number with increments of 4k - 0,4k,8k,..60k",
					),
				},
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable STP.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"extended_system_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Extended system identifier.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4095),
				},
			},
			"force_version_stp": schema.BoolAttribute{
				Optional:    true,
				Description: "Force protocol version STP.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"forward_delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Time spent in listening or learning state (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(4, 30),
				},
			},
			"hello_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Time interval between configuration BPDUs (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 10),
				},
			},
			"max_age": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum age of received protocol bpdu (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(6, 40),
				},
			},
			"priority_hold_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Hold time before switching to primary priority when core domain becomes up (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 255),
				},
			},
			"system_identifier": schema.StringAttribute{
				Optional:    true,
				Description: "System identifier to represent this node.",
				Validators: []validator.String{
					tfvalidator.StringMACAddress().WithMac48ColonHexa(),
				},
			},
			"vpls_flush_on_topology_change": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable VPLS MAC flush on root protected CE interface receiving topology change.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"system_id": schema.SetNestedBlock{
				Description: "For each ID, System ID to IP mapping.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "System ID.",
							Validators: []validator.String{
								tfvalidator.StringMACAddress().WithMac48ColonHexa(),
							},
						},
						"ip_address": schema.StringAttribute{
							Optional:    true,
							Description: "Peer ID (IP Address).",
							Validators: []validator.String{
								tfvalidator.StringCIDR(),
							},
						},
					},
				},
			},
		},
	}
}

//nolint:lll
type rstpData struct {
	ID                                           types.String        `tfsdk:"id"`
	RoutingInstance                              types.String        `tfsdk:"routing_instance"`
	BackupBridgePriority                         types.String        `tfsdk:"backup_bridge_priority"`
	BpduBlockOnEdge                              types.Bool          `tfsdk:"bpdu_block_on_edge"`
	BpduDestinationMacAddressProviderBridgeGroup types.Bool          `tfsdk:"bpdu_destination_mac_address_provider_bridge_group"`
	BridgePriority                               types.String        `tfsdk:"bridge_priority"`
	Disable                                      types.Bool          `tfsdk:"disable"`
	ExtendedSystemID                             types.Int64         `tfsdk:"extended_system_id"`
	ForceVersionStp                              types.Bool          `tfsdk:"force_version_stp"`
	ForwardDelay                                 types.Int64         `tfsdk:"forward_delay"`
	HelloTime                                    types.Int64         `tfsdk:"hello_time"`
	MaxAge                                       types.Int64         `tfsdk:"max_age"`
	PriorityHoldTime                             types.Int64         `tfsdk:"priority_hold_time"`
	SystemIdentifier                             types.String        `tfsdk:"system_identifier"`
	VplsFlushOnTopologyChange                    types.Bool          `tfsdk:"vpls_flush_on_topology_change"`
	SystemID                                     []rstpBlockSystemID `tfsdk:"system_id"`
}

type rstpConfig struct {
	ID                                           types.String `tfsdk:"id"`
	RoutingInstance                              types.String `tfsdk:"routing_instance"`
	BackupBridgePriority                         types.String `tfsdk:"backup_bridge_priority"`
	BpduBlockOnEdge                              types.Bool   `tfsdk:"bpdu_block_on_edge"`
	BpduDestinationMacAddressProviderBridgeGroup types.Bool   `tfsdk:"bpdu_destination_mac_address_provider_bridge_group"`
	BridgePriority                               types.String `tfsdk:"bridge_priority"`
	Disable                                      types.Bool   `tfsdk:"disable"`
	ExtendedSystemID                             types.Int64  `tfsdk:"extended_system_id"`
	ForceVersionStp                              types.Bool   `tfsdk:"force_version_stp"`
	ForwardDelay                                 types.Int64  `tfsdk:"forward_delay"`
	HelloTime                                    types.Int64  `tfsdk:"hello_time"`
	MaxAge                                       types.Int64  `tfsdk:"max_age"`
	PriorityHoldTime                             types.Int64  `tfsdk:"priority_hold_time"`
	SystemIdentifier                             types.String `tfsdk:"system_identifier"`
	VplsFlushOnTopologyChange                    types.Bool   `tfsdk:"vpls_flush_on_topology_change"`
	SystemID                                     types.Set    `tfsdk:"system_id"`
}

type rstpBlockSystemID struct {
	ID        types.String `tfsdk:"id"`
	IPAddress types.String `tfsdk:"ip_address"`
}

func (rsc *rstp) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config rstpConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.BackupBridgePriority.IsNull() && !config.BackupBridgePriority.IsUnknown() {
		if v, err := strconv.Atoi(strings.TrimSuffix(
			config.BackupBridgePriority.ValueString(), "k",
		)); err == nil {
			if v%4 != 0 {
				resp.Diagnostics.AddAttributeError(
					path.Root("backup_bridge_priority"),
					"Bad Value Error",
					"backup_bridge_priority must be a multiple of 4k",
				)
			}
			if v < 4 || v > 60 {
				resp.Diagnostics.AddAttributeError(
					path.Root("backup_bridge_priority"),
					"Bad Value Error",
					"backup_bridge_priority must be between 4k and 60k",
				)
			}
			if !config.BridgePriority.IsNull() && !config.BridgePriority.IsUnknown() {
				if bridgePriority, err := strconv.Atoi(strings.TrimSuffix(
					config.BridgePriority.ValueString(), "k",
				)); err == nil {
					if v <= bridgePriority {
						resp.Diagnostics.AddAttributeError(
							path.Root("backup_bridge_priority"),
							"Bad Value Error",
							"backup_bridge_priority must be worse (higher value) than bridge_priority",
						)
					}
				}
			}
		}
	}
	if !config.BridgePriority.IsNull() && !config.BridgePriority.IsUnknown() {
		if v, err := strconv.Atoi(strings.TrimSuffix(
			config.BridgePriority.ValueString(), "k",
		)); err == nil {
			if v%4 != 0 {
				resp.Diagnostics.AddAttributeError(
					path.Root("bridge_priority"),
					"Bad Value Error",
					"bridge_priority must be a multiple of 4k",
				)
			}
			if v < 0 || v > 60 {
				resp.Diagnostics.AddAttributeError(
					path.Root("bridge_priority"),
					"Bad Value Error",
					"bridge_priority must be between 0 and 60k",
				)
			}
		}
	}
	if !config.SystemID.IsNull() && !config.SystemID.IsUnknown() {
		var systemID []rstpBlockSystemID
		asDiags := config.SystemID.ElementsAs(ctx, &systemID, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		systemIDID := make(map[string]struct{})
		for _, block := range systemID {
			if block.ID.IsUnknown() {
				continue
			}
			id := block.ID.ValueString()
			if _, ok := systemIDID[id]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("system_id"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple system_id blocks with the same id %q", id),
				)
			}
			systemIDID[id] = struct{}{}
		}
	}
}

func (rsc *rstp) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan rstpData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *rstp) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data rstpData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	if v := state.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, v, junSess)
		if err != nil {
			junos.MutexUnlock()
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			junos.MutexUnlock()
			resp.State.RemoveResource(ctx)

			return
		}
	}

	err = data.read(ctx, state.RoutingInstance.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	if data.nullID() {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *rstp) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state rstpData
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

func (rsc *rstp) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state rstpData
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

func (rsc *rstp) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	if req.ID != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, req.ID, junSess)
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			resp.Diagnostics.AddError(
				tfdiag.NotFoundErrSummary,
				fmt.Sprintf("routing instance %q doesn't exist", req.ID),
			)

			return
		}
	}

	var data rstpData
	if err := data.read(ctx, req.ID, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be <routing_instance>)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *rstpData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(v)
	} else {
		rscData.ID = types.StringValue(junos.DefaultW)
	}
}

func (rscData *rstpData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *rstpData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols rstp "

	if v := rscData.BackupBridgePriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"backup-bridge-priority "+v)
	}
	if rscData.BpduBlockOnEdge.ValueBool() {
		configSet = append(configSet, setPrefix+"bpdu-block-on-edge")
	}
	if rscData.BpduDestinationMacAddressProviderBridgeGroup.ValueBool() {
		configSet = append(configSet, setPrefix+"bpdu-destination-mac-address provider-bridge-group")
	}
	if v := rscData.BridgePriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"bridge-priority "+v)
	}
	if rscData.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if !rscData.ExtendedSystemID.IsNull() {
		configSet = append(configSet, setPrefix+"extended-system-id "+
			utils.ConvI64toa(rscData.ExtendedSystemID.ValueInt64()))
	}
	if rscData.ForceVersionStp.ValueBool() {
		configSet = append(configSet, setPrefix+"force-version stp")
	}
	if !rscData.ForwardDelay.IsNull() {
		configSet = append(configSet, setPrefix+"forward-delay "+
			utils.ConvI64toa(rscData.ForwardDelay.ValueInt64()))
	}
	if !rscData.HelloTime.IsNull() {
		configSet = append(configSet, setPrefix+"hello-time "+
			utils.ConvI64toa(rscData.HelloTime.ValueInt64()))
	}
	if !rscData.MaxAge.IsNull() {
		configSet = append(configSet, setPrefix+"max-age "+
			utils.ConvI64toa(rscData.MaxAge.ValueInt64()))
	}
	if !rscData.PriorityHoldTime.IsNull() {
		configSet = append(configSet, setPrefix+"priority-hold-time "+
			utils.ConvI64toa(rscData.PriorityHoldTime.ValueInt64()))
	}
	if v := rscData.SystemIdentifier.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"system-identifier "+v)
	}
	if rscData.VplsFlushOnTopologyChange.ValueBool() {
		configSet = append(configSet, setPrefix+"vpls-flush-on-topology-change")
	}
	systemIDID := make(map[string]struct{})
	for _, block := range rscData.SystemID {
		id := block.ID.ValueString()
		if _, ok := systemIDID[id]; ok {
			return path.Root("system_id"),
				fmt.Errorf("multiple system_id blocks with the same id %q", id)
		}
		systemIDID[id] = struct{}{}

		configSet = append(configSet, setPrefix+"system-id "+id)
		if v := block.IPAddress.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"system-id "+id+" ip-address "+v)
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *rstpData) read(
	_ context.Context, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols rstp" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if routingInstance == "" {
		rscData.RoutingInstance = types.StringValue(junos.DefaultW)
	} else {
		rscData.RoutingInstance = types.StringValue(routingInstance)
	}
	rscData.fillID()
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "backup-bridge-priority "):
				rscData.BackupBridgePriority = types.StringValue(itemTrim)
			case itemTrim == "bpdu-block-on-edge":
				rscData.BpduBlockOnEdge = types.BoolValue(true)
			case itemTrim == "bpdu-destination-mac-address provider-bridge-group":
				rscData.BpduDestinationMacAddressProviderBridgeGroup = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "bridge-priority "):
				rscData.BridgePriority = types.StringValue(itemTrim)
			case itemTrim == junos.DisableW:
				rscData.Disable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "extended-system-id "):
				rscData.ExtendedSystemID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "force-version stp":
				rscData.ForceVersionStp = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "forward-delay "):
				rscData.ForwardDelay, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "hello-time "):
				rscData.HelloTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "max-age "):
				rscData.MaxAge, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "priority-hold-time "):
				rscData.PriorityHoldTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "system-identifier "):
				rscData.SystemIdentifier = types.StringValue(itemTrim)
			case itemTrim == "vpls-flush-on-topology-change":
				rscData.VplsFlushOnTopologyChange = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "system-id "):
				itemTrimFields := strings.Split(itemTrim, " ")
				switch len(itemTrimFields) { // <id> (ip-address <ip_address>)?
				case 1:
					rscData.SystemID = append(rscData.SystemID, rstpBlockSystemID{
						ID: types.StringValue(itemTrimFields[0]),
					})
				case 3:
					rscData.SystemID = append(rscData.SystemID, rstpBlockSystemID{
						ID:        types.StringValue(itemTrimFields[0]),
						IPAddress: types.StringValue(itemTrimFields[2]),
					})
				default:
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "system-id", itemTrim)
				}
			}
		}
	}

	return nil
}

func (rscData *rstpData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "protocols rstp "

	configSet := []string{
		delPrefix + "backup-bridge-priority",
		delPrefix + "bpdu-block-on-edge",
		delPrefix + "bpdu-destination-mac-address",
		delPrefix + "bridge-priority",
		delPrefix + "disable",
		delPrefix + "extended-system-id",
		delPrefix + "force-version",
		delPrefix + "forward-delay",
		delPrefix + "hello-time",
		delPrefix + "max-age",
		delPrefix + "priority-hold-time",
		delPrefix + "system-identifier",
		delPrefix + "vpls-flush-on-topology-change",
	}
	for _, block := range rscData.SystemID {
		configSet = append(configSet, delPrefix+"system-id "+block.ID.ValueString())
	}

	return junSess.ConfigSet(configSet)
}
