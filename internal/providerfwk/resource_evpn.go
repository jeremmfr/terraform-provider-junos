package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &evpn{}
	_ resource.ResourceWithConfigure      = &evpn{}
	_ resource.ResourceWithValidateConfig = &evpn{}
	_ resource.ResourceWithImportState    = &evpn{}
	_ resource.ResourceWithUpgradeState   = &evpn{}
)

type evpn struct {
	client *junos.Client
}

func newEvpnResource() resource.Resource {
	return &evpn{}
}

func (rsc *evpn) typeName() string {
	return providerName + "_evpn"
}

func (rsc *evpn) junosName() string {
	return "protocols evpn"
}

func (rsc *evpn) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *evpn) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *evpn) Configure(
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

func (rsc *evpn) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"encapsulation": schema.StringAttribute{
				Required:    true,
				Description: "Encapsulation type for EVPN.",
				Validators: []validator.String{
					stringvalidator.OneOf("mpls", "vxlan"),
				},
			},
			"default_gateway": schema.StringAttribute{
				Optional:    true,
				Description: "Default gateway mode.",
				Validators: []validator.String{
					stringvalidator.OneOf("advertise", "do-not-advertise", "no-gateway-community"),
				},
			},
			"multicast_mode": schema.StringAttribute{
				Optional:    true,
				Description: "Multicast mode for EVPN.",
				Validators: []validator.String{
					stringvalidator.OneOf("ingress-replication"),
				},
			},
			"no_core_isolation": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable EVPN Core isolation.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"routing_instance_evpn": schema.BoolAttribute{
				Optional:    true,
				Description: "Configure routing instance is an evpn instance-type.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"duplicate_mac_detection": schema.SingleNestedBlock{
				Description: "Duplicate MAC detection settings",
				Attributes: map[string]schema.Attribute{
					"auto_recovery_time": schema.Int64Attribute{
						Optional:    true,
						Description: "Automatically unblock duplicate MACs after a time delay.",
						Validators: []validator.Int64{
							int64validator.Between(1, 360),
						},
					},
					"detection_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of moves to trigger duplicate MAC detection.",
						Validators: []validator.Int64{
							int64validator.Between(2, 20),
						},
					},
					"detection_window": schema.Int64Attribute{
						Optional:    true,
						Description: "Time window for detection of duplicate MACs.",
						Validators: []validator.Int64{
							int64validator.Between(5, 600),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"switch_or_ri_options": schema.SingleNestedBlock{
				Description: "Declare `switch-options` or `routing-instance` configuration.",
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
					tfplanmodifier.BlockSetUnsetRequireReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"route_distinguisher": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Route distinguisher for this instance.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^(\d|\.)+L?:\d+$`),
								"must have valid route distinguisher. Use format 'x:y'"),
						},
					},
					"vrf_export": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Export policy for VRF instance RIBs.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.ValueStringsAre(
								stringvalidator.LengthBetween(1, 250),
								tfvalidator.StringDoubleQuoteExclusion(),
							),
						},
					},
					"vrf_import": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Import policy for VRF instance RIBs.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.ValueStringsAre(
								stringvalidator.LengthBetween(1, 250),
								tfvalidator.StringDoubleQuoteExclusion(),
							),
						},
					},
					"vrf_target": schema.StringAttribute{
						Optional:    true,
						Description: "VRF target community configuration.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^target:(\d|\.)+L?:\d+$`),
								"must have valid target. Use format 'target:x:y'"),
						},
					},
					"vrf_target_auto": schema.BoolAttribute{
						Optional:    true,
						Description: "Auto derive import and export target community from BGP AS & L2.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"vrf_target_export": schema.StringAttribute{
						Optional:    true,
						Description: "Target community to use when marking routes on export.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^target:(\d|\.)+L?:\d+$`),
								"must have valid target. Use format 'target:x:y'"),
						},
					},
					"vrf_target_import": schema.StringAttribute{
						Optional:    true,
						Description: "Target community to use when filtering on import.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^target:(\d|\.)+L?:\d+$`),
								"must have valid target. Use format 'target:x:y'"),
						},
					},
				},
			},
		},
	}
}

type evpnData struct {
	ID                    types.String                    `tfsdk:"id"`
	RoutingInstance       types.String                    `tfsdk:"routing_instance"`
	Encapsulation         types.String                    `tfsdk:"encapsulation"`
	DefaultGateway        types.String                    `tfsdk:"default_gateway"`
	MulticastMode         types.String                    `tfsdk:"multicast_mode"`
	NoCoreIsolation       types.Bool                      `tfsdk:"no_core_isolation"`
	RoutingInstanceEvpn   types.Bool                      `tfsdk:"routing_instance_evpn"`
	DuplicateMacDetection *evpnBlockDuplicateMACDetection `tfsdk:"duplicate_mac_detection"`
	SwitchOrRIOptions     *evpnBlockSwitchOrRIOptions     `tfsdk:"switch_or_ri_options"`
}

type evpnConfig struct {
	ID                    types.String                      `tfsdk:"id"`
	RoutingInstance       types.String                      `tfsdk:"routing_instance"`
	Encapsulation         types.String                      `tfsdk:"encapsulation"`
	DefaultGateway        types.String                      `tfsdk:"default_gateway"`
	MulticastMode         types.String                      `tfsdk:"multicast_mode"`
	NoCoreIsolation       types.Bool                        `tfsdk:"no_core_isolation"`
	RoutingInstanceEvpn   types.Bool                        `tfsdk:"routing_instance_evpn"`
	DuplicateMacDetection *evpnBlockDuplicateMACDetection   `tfsdk:"duplicate_mac_detection"`
	SwitchOrRIOptions     *evpnBlockSwitchOrRIOptionsConfig `tfsdk:"switch_or_ri_options"`
}

type evpnBlockDuplicateMACDetection struct {
	AutoRecoveryTime   types.Int64 `tfsdk:"auto_recovery_time"`
	DetectionThreshold types.Int64 `tfsdk:"detection_threshold"`
	DetectionWindow    types.Int64 `tfsdk:"detection_window"`
}

func (block *evpnBlockDuplicateMACDetection) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type evpnBlockSwitchOrRIOptions struct {
	RouteDistinguisher types.String   `tfsdk:"route_distinguisher"`
	VRFExport          []types.String `tfsdk:"vrf_export"`
	VRFImport          []types.String `tfsdk:"vrf_import"`
	VRFTarget          types.String   `tfsdk:"vrf_target"`
	VRFTargetAuto      types.Bool     `tfsdk:"vrf_target_auto"`
	VRFTargetExport    types.String   `tfsdk:"vrf_target_export"`
	VRFTargetImport    types.String   `tfsdk:"vrf_target_import"`
}

type evpnBlockSwitchOrRIOptionsConfig struct {
	RouteDistinguisher types.String `tfsdk:"route_distinguisher"`
	VRFExport          types.List   `tfsdk:"vrf_export"`
	VRFImport          types.List   `tfsdk:"vrf_import"`
	VRFTarget          types.String `tfsdk:"vrf_target"`
	VRFTargetAuto      types.Bool   `tfsdk:"vrf_target_auto"`
	VRFTargetExport    types.String `tfsdk:"vrf_target_export"`
	VRFTargetImport    types.String `tfsdk:"vrf_target_import"`
}

func (rsc *evpn) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config evpnConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.RoutingInstance.IsNull() ||
		(!config.RoutingInstance.IsUnknown() && config.RoutingInstance.ValueString() == junos.DefaultW) {
		if config.SwitchOrRIOptions == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("switch_or_ri_options"),
				tfdiag.MissingConfigErrSummary,
				fmt.Sprintf("switch_or_ri_options must be specified when routing_instance = %q", junos.DefaultW),
			)
		}
		if !config.RoutingInstanceEvpn.IsNull() && !config.RoutingInstanceEvpn.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("routing_instance_evpn"),
				tfdiag.ConflictConfigErrSummary,
				fmt.Sprintf("routing_instance_evpn cannot be configured when routing_instance = %q", junos.DefaultW),
			)
		}
		if !config.DefaultGateway.IsNull() && !config.DefaultGateway.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("default_gateway"),
				tfdiag.ConflictConfigErrSummary,
				fmt.Sprintf("default_gateway cannot be configured when routing_instance = %q", junos.DefaultW),
			)
		}
	}
	if !config.RoutingInstance.IsNull() && !config.RoutingInstance.IsUnknown() &&
		config.RoutingInstance.ValueString() != junos.DefaultW &&
		!config.NoCoreIsolation.IsNull() && !config.NoCoreIsolation.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("no_core_isolation"),
			tfdiag.ConflictConfigErrSummary,
			fmt.Sprintf("no_core_isolation cannot be configured when routing_instance != %q", junos.DefaultW),
		)
	}
	if config.DuplicateMacDetection != nil {
		if config.DuplicateMacDetection.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("duplicate_mac_detection").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"duplicate_mac_detection block is empty",
			)
		}
	}
	if !config.RoutingInstanceEvpn.IsNull() && !config.RoutingInstanceEvpn.IsUnknown() {
		if config.RoutingInstance.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("routing_instance_evpn"),
				tfdiag.MissingConfigErrSummary,
				"routing_instance must be specified with routing_instance_evpn",
			)
		} else if !config.RoutingInstance.IsUnknown() {
			if routingInstance := config.RoutingInstance.ValueString(); routingInstance == junos.DefaultW {
				resp.Diagnostics.AddAttributeError(
					path.Root("routing_instance_evpn"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("routing_instance must be specified with routing_instance_evpn and not with value %q", junos.DefaultW),
				)
			}
		}
		if config.SwitchOrRIOptions == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("routing_instance_evpn"),
				tfdiag.MissingConfigErrSummary,
				"switch_or_ri_options must be specified with routing_instance_evpn",
			)
		}
	}
	if config.SwitchOrRIOptions != nil {
		if config.SwitchOrRIOptions.RouteDistinguisher.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("switch_or_ri_options").AtName("route_distinguisher"),
				tfdiag.MissingConfigErrSummary,
				"route_distinguisher must be specified in switch_or_ri_options block",
			)
		}
	}
}

func (rsc *evpn) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan evpnData
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

func (rsc *evpn) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data evpnData
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

	if state.SwitchOrRIOptions == nil {
		data.SwitchOrRIOptions = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *evpn) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state evpnData
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

func (rsc *evpn) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state evpnData
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

func (rsc *evpn) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	idList := strings.Split(req.ID, junos.IDSeparator)
	if idList[0] != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, idList[0], junSess)
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			resp.Diagnostics.AddError(
				tfdiag.NotFoundErrSummary,
				fmt.Sprintf("routing instance %q doesn't exist", idList[0]),
			)

			return
		}
	}

	var data evpnData
	if err := data.read(ctx, idList[0], junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "routing_instance"),
		)

		return
	}

	if idList[0] != junos.DefaultW && len(idList) == 1 {
		data.SwitchOrRIOptions = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *evpnData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(v)
	} else {
		rscData.ID = types.StringValue(junos.DefaultW)
	}
}

func (rscData *evpnData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *evpnData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	setSwitchRIPrefix := junos.SetLS
	switch routingInstance := rscData.RoutingInstance.ValueString(); routingInstance {
	case junos.DefaultW, "":
		if rscData.SwitchOrRIOptions == nil {
			return path.Root("switch_or_ri_options"),
				fmt.Errorf("switch_or_ri_options must be specified when routing_instance = %q", junos.DefaultW)
		}
		if !rscData.RoutingInstanceEvpn.IsNull() {
			return path.Root("routing_instance_evpn"),
				fmt.Errorf("routing_instance_evpn cannot be configured when routing_instance = %q", junos.DefaultW)
		}
		if !rscData.DefaultGateway.IsNull() {
			return path.Root("default_gateway"),
				fmt.Errorf("default_gateway cannot be configured when routing_instance = %q", junos.DefaultW)
		}
		setSwitchRIPrefix += "switch-options "
	default:
		setPrefix += junos.RoutingInstancesWS + routingInstance + " "
		setSwitchRIPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	setPrefix += "protocols evpn "

	if rscData.RoutingInstanceEvpn.ValueBool() {
		if rscData.SwitchOrRIOptions == nil {
			return path.Root("switch_or_ri_options"),
				errors.New("switch_or_ri_options must be specified with routing_instance_evpn")
		}
		configSet = append(configSet, setSwitchRIPrefix+"instance-type evpn")
	}
	if v := rscData.Encapsulation.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"encapsulation "+v)
	}
	if v := rscData.DefaultGateway.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"default-gateway "+v)
	}
	if rscData.DuplicateMacDetection != nil {
		if rscData.DuplicateMacDetection.isEmpty() {
			return path.Root("duplicate_mac_detection"),
				errors.New("duplicate_mac_detection block is empty")
		}
		if !rscData.DuplicateMacDetection.AutoRecoveryTime.IsNull() {
			configSet = append(configSet, setPrefix+"duplicate-mac-detection auto-recovery-time "+
				utils.ConvI64toa(rscData.DuplicateMacDetection.AutoRecoveryTime.ValueInt64()),
			)
		}
		if !rscData.DuplicateMacDetection.DetectionThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"duplicate-mac-detection detection-threshold "+
				utils.ConvI64toa(rscData.DuplicateMacDetection.DetectionThreshold.ValueInt64()),
			)
		}
		if !rscData.DuplicateMacDetection.DetectionWindow.IsNull() {
			configSet = append(configSet, setPrefix+"duplicate-mac-detection detection-window "+
				utils.ConvI64toa(rscData.DuplicateMacDetection.DetectionWindow.ValueInt64()),
			)
		}
	}
	if v := rscData.MulticastMode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"multicast-mode "+v)
	}
	if rscData.NoCoreIsolation.ValueBool() {
		configSet = append(configSet, setPrefix+"no-core-isolation")
	}
	if rscData.SwitchOrRIOptions != nil {
		configSet = append(configSet,
			setSwitchRIPrefix+"route-distinguisher "+rscData.SwitchOrRIOptions.RouteDistinguisher.ValueString())
		for _, v := range rscData.SwitchOrRIOptions.VRFExport {
			configSet = append(configSet, setSwitchRIPrefix+"vrf-export \""+v.ValueString()+"\"")
		}
		for _, v := range rscData.SwitchOrRIOptions.VRFImport {
			configSet = append(configSet, setSwitchRIPrefix+"vrf-import \""+v.ValueString()+"\"")
		}
		if v := rscData.SwitchOrRIOptions.VRFTarget.ValueString(); v != "" {
			configSet = append(configSet, setSwitchRIPrefix+"vrf-target "+v)
		}
		if rscData.SwitchOrRIOptions.VRFTargetAuto.ValueBool() {
			configSet = append(configSet, setSwitchRIPrefix+"vrf-target auto")
		}
		if v := rscData.SwitchOrRIOptions.VRFTargetExport.ValueString(); v != "" {
			configSet = append(configSet, setSwitchRIPrefix+"vrf-target export "+v)
		}
		if v := rscData.SwitchOrRIOptions.VRFTargetImport.ValueString(); v != "" {
			configSet = append(configSet, setSwitchRIPrefix+"vrf-target import "+v)
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *evpnData) read(
	_ context.Context, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	showSwitchRI := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
		showSwitchRI += junos.RoutingInstancesWS + routingInstance
	} else {
		showSwitchRI += "switch-options"
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols evpn" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	showConfigSwitchRI, err := junSess.Command(showSwitchRI + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
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
			case balt.CutPrefixInString(&itemTrim, "default-gateway "):
				rscData.DefaultGateway = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "duplicate-mac-detection "):
				if rscData.DuplicateMacDetection == nil {
					rscData.DuplicateMacDetection = &evpnBlockDuplicateMACDetection{}
				}
				var err error
				switch {
				case balt.CutPrefixInString(&itemTrim, "auto-recovery-time "):
					rscData.DuplicateMacDetection.AutoRecoveryTime, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "detection-threshold "):
					rscData.DuplicateMacDetection.DetectionThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "detection-window "):
					rscData.DuplicateMacDetection.DetectionWindow, err = tfdata.ConvAtoi64Value(itemTrim)
				}
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "encapsulation "):
				rscData.Encapsulation = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "multicast-mode "):
				rscData.MulticastMode = types.StringValue(itemTrim)
			case itemTrim == "no-core-isolation":
				rscData.NoCoreIsolation = types.BoolValue(true)
			}
		}
	}
	if showConfigSwitchRI != junos.EmptyW {
		for _, item := range strings.Split(showConfigSwitchRI, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "instance-type evpn":
				rscData.RoutingInstanceEvpn = types.BoolValue(true)
			case strings.HasPrefix(itemTrim, "route-distinguisher "),
				strings.HasPrefix(itemTrim, "vrf-export"),
				strings.HasPrefix(itemTrim, "vrf-import"),
				strings.HasPrefix(itemTrim, "vrf-target"):
				if rscData.SwitchOrRIOptions == nil {
					rscData.SwitchOrRIOptions = &evpnBlockSwitchOrRIOptions{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "route-distinguisher "):
					rscData.SwitchOrRIOptions.RouteDistinguisher = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "vrf-export "):
					rscData.SwitchOrRIOptions.VRFExport = append(rscData.SwitchOrRIOptions.VRFExport,
						types.StringValue(strings.Trim(itemTrim, "\"")),
					)
				case balt.CutPrefixInString(&itemTrim, "vrf-import "):
					rscData.SwitchOrRIOptions.VRFImport = append(rscData.SwitchOrRIOptions.VRFImport,
						types.StringValue(strings.Trim(itemTrim, "\"")),
					)
				case itemTrim == "vrf-target auto":
					rscData.SwitchOrRIOptions.VRFTargetAuto = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "vrf-target export "):
					rscData.SwitchOrRIOptions.VRFTargetExport = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "vrf-target import "):
					rscData.SwitchOrRIOptions.VRFTargetImport = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "vrf-target "):
					rscData.SwitchOrRIOptions.VRFTarget = types.StringValue(itemTrim)
				}
			}
		}
	}

	return nil
}

func (rscData *evpnData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	delSwitchRIPrefix := junos.DeleteLS
	switch routingInstance := rscData.RoutingInstance.ValueString(); routingInstance {
	case junos.DefaultW, "":
		delSwitchRIPrefix += "switch-options "
	default:
		delPrefix += junos.RoutingInstancesWS + routingInstance + " "
		delSwitchRIPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	delPrefix += "protocols evpn "

	configSet := []string{
		delPrefix + "default-gateway",
		delPrefix + "duplicate-mac-detection",
		delPrefix + "encapsulation",
		delPrefix + "multicast-mode",
		delPrefix + "no-core-isolation",
	}

	if rscData.RoutingInstanceEvpn.ValueBool() {
		configSet = append(configSet,
			delSwitchRIPrefix+"instance-type")
	}
	if rscData.SwitchOrRIOptions != nil {
		configSet = append(configSet,
			delSwitchRIPrefix+"route-distinguisher",
			delSwitchRIPrefix+"vrf-export",
			delSwitchRIPrefix+"vrf-import",
			delSwitchRIPrefix+"vrf-target",
		)
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *evpnData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	delSwitchRIPrefix := junos.DeleteLS
	switch routingInstance := rscData.RoutingInstance.ValueString(); routingInstance {
	case junos.DefaultW, "":
		delSwitchRIPrefix += "switch-options "
	default:
		delPrefix += junos.RoutingInstancesWS + routingInstance + " "
		delSwitchRIPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}

	configSet := []string{
		delPrefix + "protocols evpn",
	}

	if rscData.RoutingInstanceEvpn.ValueBool() {
		configSet = append(configSet,
			delSwitchRIPrefix+"instance-type")
	}
	if rscData.SwitchOrRIOptions != nil {
		configSet = append(configSet,
			delSwitchRIPrefix+"route-distinguisher",
			delSwitchRIPrefix+"vrf-export",
			delSwitchRIPrefix+"vrf-import",
			delSwitchRIPrefix+"vrf-target",
		)
	}

	return junSess.ConfigSet(configSet)
}
