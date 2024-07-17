package providerfwk

import (
	"context"
	"fmt"
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
	_ resource.Resource                   = &vstp{}
	_ resource.ResourceWithConfigure      = &vstp{}
	_ resource.ResourceWithValidateConfig = &vstp{}
	_ resource.ResourceWithImportState    = &vstp{}
)

type vstp struct {
	client *junos.Client
}

func newVstpResource() resource.Resource {
	return &vstp{}
}

func (rsc *vstp) typeName() string {
	return providerName + "_vstp"
}

func (rsc *vstp) junosName() string {
	return "protocols vstp"
}

func (rsc *vstp) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *vstp) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *vstp) Configure(
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

func (rsc *vstp) Schema(
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
				Description: "Routing instance for vstp protocol if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},

			"bpdu_block_on_edge": schema.BoolAttribute{
				Optional:    true,
				Description: "Block BPDU on all interfaces configured as edge (BPDU Protect).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable STP.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"force_version_stp": schema.BoolAttribute{
				Optional:    true,
				Description: "Force protocol version STP.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"priority_hold_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Hold time before switching to primary priority when core domain becomes up (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 255),
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

type vstpData struct {
	ID                        types.String        `tfsdk:"id"`
	RoutingInstance           types.String        `tfsdk:"routing_instance"`
	BpduBlockOnEdge           types.Bool          `tfsdk:"bpdu_block_on_edge"`
	Disable                   types.Bool          `tfsdk:"disable"`
	ForceVersionStp           types.Bool          `tfsdk:"force_version_stp"`
	PriorityHoldTime          types.Int64         `tfsdk:"priority_hold_time"`
	VplsFlushOnTopologyChange types.Bool          `tfsdk:"vpls_flush_on_topology_change"`
	SystemID                  []vstpBlockSystemID `tfsdk:"system_id"`
}

type vstpConfig struct {
	ID                        types.String `tfsdk:"id"`
	RoutingInstance           types.String `tfsdk:"routing_instance"`
	BpduBlockOnEdge           types.Bool   `tfsdk:"bpdu_block_on_edge"`
	Disable                   types.Bool   `tfsdk:"disable"`
	ForceVersionStp           types.Bool   `tfsdk:"force_version_stp"`
	PriorityHoldTime          types.Int64  `tfsdk:"priority_hold_time"`
	VplsFlushOnTopologyChange types.Bool   `tfsdk:"vpls_flush_on_topology_change"`
	SystemID                  types.Set    `tfsdk:"system_id"`
}

type vstpBlockSystemID struct {
	ID        types.String `tfsdk:"id"`
	IPAddress types.String `tfsdk:"ip_address"`
}

func (rsc *vstp) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config vstpConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.SystemID.IsNull() && !config.SystemID.IsUnknown() {
		var systemID []vstpBlockSystemID
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

func (rsc *vstp) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan vstpData
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

func (rsc *vstp) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data vstpData
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

func (rsc *vstp) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state vstpData
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

func (rsc *vstp) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state vstpData
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

func (rsc *vstp) ImportState(
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

	var data vstpData
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

func (rscData *vstpData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(v)
	} else {
		rscData.ID = types.StringValue(junos.DefaultW)
	}
}

func (rscData *vstpData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *vstpData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols vstp "

	if rscData.BpduBlockOnEdge.ValueBool() {
		configSet = append(configSet, setPrefix+"bpdu-block-on-edge")
	}
	if rscData.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if rscData.ForceVersionStp.ValueBool() {
		configSet = append(configSet, setPrefix+"force-version stp")
	}

	if !rscData.PriorityHoldTime.IsNull() {
		configSet = append(configSet, setPrefix+"priority-hold-time "+
			utils.ConvI64toa(rscData.PriorityHoldTime.ValueInt64()))
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

func (rscData *vstpData) read(
	_ context.Context, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols vstp" + junos.PipeDisplaySetRelative)
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
			case itemTrim == "bpdu-block-on-edge":
				rscData.BpduBlockOnEdge = types.BoolValue(true)
			case itemTrim == junos.DisableW:
				rscData.Disable = types.BoolValue(true)
			case itemTrim == "force-version stp":
				rscData.ForceVersionStp = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "priority-hold-time "):
				rscData.PriorityHoldTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "vpls-flush-on-topology-change":
				rscData.VplsFlushOnTopologyChange = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "system-id "):
				itemTrimFields := strings.Split(itemTrim, " ")
				switch len(itemTrimFields) { // <id> (ip-address <ip_address>)?
				case 1:
					rscData.SystemID = append(rscData.SystemID, vstpBlockSystemID{
						ID: types.StringValue(itemTrimFields[0]),
					})
				case 3:
					rscData.SystemID = append(rscData.SystemID, vstpBlockSystemID{
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

func (rscData *vstpData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "protocols vstp "

	configSet := []string{
		delPrefix + "bpdu-block-on-edge",
		delPrefix + "disable",
		delPrefix + "force-version",
		delPrefix + "priority-hold-time",
		delPrefix + "vpls-flush-on-topology-change",
	}
	for _, block := range rscData.SystemID {
		configSet = append(configSet, delPrefix+"system-id "+block.ID.ValueString())
	}

	return junSess.ConfigSet(configSet)
}
