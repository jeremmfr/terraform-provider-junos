package providerfwk

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                   = &iccpPeer{}
	_ resource.ResourceWithConfigure      = &iccpPeer{}
	_ resource.ResourceWithValidateConfig = &iccpPeer{}
	_ resource.ResourceWithImportState    = &iccpPeer{}
)

type iccpPeer struct {
	client *junos.Client
}

func newIccpPeerResource() resource.Resource {
	return &iccpPeer{}
}

func (rsc *iccpPeer) typeName() string {
	return providerName + "_iccp_peer"
}

func (rsc *iccpPeer) junosName() string {
	return "iccp peer"
}

func (rsc *iccpPeer) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *iccpPeer) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *iccpPeer) Configure(
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

func (rsc *iccpPeer) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<ip_address>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_address": schema.StringAttribute{
				Required:    true,
				Description: "IP address for this peer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
			"redundancy_group_id_list": schema.SetAttribute{
				ElementType: types.Int64Type,
				Required:    true,
				Description: "List of redundancy groups this peer is part of.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueInt64sAre(
						int64validator.Between(1, 4294967294),
					),
				},
			},
			"authentication_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "MD5 authentication key.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 126),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"local_ip_addr": schema.StringAttribute{
				Optional:    true,
				Description: "Local IP address to use for this peer alone.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
			"session_establishment_hold_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Time within which connection must succeed with this peer (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(45, 600),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"backup_liveness_detection": schema.SingleNestedBlock{
				Description: "Backup liveness detection.",
				Attributes: map[string]schema.Attribute{
					"backup_peer_ip": schema.StringAttribute{
						Optional:    true,
						Description: "Backup liveness detection peer's IP address.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress().IPv4Only(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"liveness_detection": schema.SingleNestedBlock{
				Description: "Bidirectional Forwarding Detection options for the peer.",
				Attributes: map[string]schema.Attribute{
					"detection_time_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "High detection-time triggering a trap (milliseconds).",
						Validators: []validator.Int64{
							int64validator.Between(0, 4294967295),
						},
					},
					"minimum_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Minimum transmit and receive interval (milliseconds).",
						Validators: []validator.Int64{
							int64validator.Between(1, 255000),
						},
					},
					"minimum_receive_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Minimum receive interval (milliseconds).",
						Validators: []validator.Int64{
							int64validator.Between(1, 255000),
						},
					},
					"multiplier": schema.Int64Attribute{
						Optional:    true,
						Description: "Detection time multiplier (1..255).",
						Validators: []validator.Int64{
							int64validator.Between(1, 255),
						},
					},
					"no_adaptation": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable adaptation.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"transmit_interval_minimum_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Minimum transmit interval (milliseconds).",
						Validators: []validator.Int64{
							int64validator.Between(1, 255000),
						},
					},
					"transmit_interval_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "High transmit interval triggering a trap (milliseconds).",
						Validators: []validator.Int64{
							int64validator.Between(0, 4294967295),
						},
					},
					"version": schema.StringAttribute{
						Optional:    true,
						Description: "BFD protocol version number.",
						Validators: []validator.String{
							stringvalidator.OneOf("0", "1", "automatic"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
		},
	}
}

type iccpPeerData struct {
	ID                           types.String                          `tfsdk:"id"`
	IPAddress                    types.String                          `tfsdk:"ip_address"`
	RedundancyGroupIDList        []types.Int64                         `tfsdk:"redundancy_group_id_list"`
	AuthenticationKey            types.String                          `tfsdk:"authentication_key"`
	LocalIPAddr                  types.String                          `tfsdk:"local_ip_addr"`
	SessionEstablishmentHoldTime types.Int64                           `tfsdk:"session_establishment_hold_time"`
	BackupLivenessDetection      *iccpPeerBlockBackupLivenessDetection `tfsdk:"backup_liveness_detection"`
	LivenessDetection            *iccpPeerBlockLivenessDetection       `tfsdk:"liveness_detection"`
}

type iccpPeerConfig struct {
	ID                           types.String                          `tfsdk:"id"`
	IPAddress                    types.String                          `tfsdk:"ip_address"`
	RedundancyGroupIDList        types.Set                             `tfsdk:"redundancy_group_id_list"`
	AuthenticationKey            types.String                          `tfsdk:"authentication_key"`
	LocalIPAddr                  types.String                          `tfsdk:"local_ip_addr"`
	SessionEstablishmentHoldTime types.Int64                           `tfsdk:"session_establishment_hold_time"`
	BackupLivenessDetection      *iccpPeerBlockBackupLivenessDetection `tfsdk:"backup_liveness_detection"`
	LivenessDetection            *iccpPeerBlockLivenessDetection       `tfsdk:"liveness_detection"`
}

type iccpPeerBlockBackupLivenessDetection struct {
	BackupPeerIP types.String `tfsdk:"backup_peer_ip"`
}

type iccpPeerBlockLivenessDetection struct {
	DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
	MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
	MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
	Multiplier                      types.Int64  `tfsdk:"multiplier"`
	NoAdaptation                    types.Bool   `tfsdk:"no_adaptation"`
	TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
	TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
	Version                         types.String `tfsdk:"version"`
}

func (rsc *iccpPeer) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config iccpPeerConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.LivenessDetection == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("liveness_detection").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"liveness_detection block must be specified",
		)
	} else if config.LivenessDetection.MinimumInterval.IsNull() {
		if config.LivenessDetection.MinimumReceiveInterval.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("liveness_detection").AtName("minimum_receive_interval"),
				tfdiag.MissingConfigErrSummary,
				"minimum_receive_interval must be specified if minimum_interval is not specified"+
					" in liveness_detection block",
			)
		}
		if config.LivenessDetection.TransmitIntervalMinimumInterval.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("liveness_detection").AtName("transmit_interval_minimum_interval"),
				tfdiag.MissingConfigErrSummary,
				"transmit_interval_minimum_interval must be specified if minimum_interval is not specified"+
					" in liveness_detection block",
			)
		}
	}
}

func (rsc *iccpPeer) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan iccpPeerData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.IPAddress.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("ip_address"),
			"Empty IP Address",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "ip_address"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			peerExists, err := checkIccpPeerExists(fnCtx, plan.IPAddress.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if peerExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.IPAddress),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			peerExists, err := checkIccpPeerExists(fnCtx, plan.IPAddress.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !peerExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.IPAddress),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *iccpPeer) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data iccpPeerData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.IPAddress.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *iccpPeer) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state iccpPeerData
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

func (rsc *iccpPeer) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state iccpPeerData
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

func (rsc *iccpPeer) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data iccpPeerData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "ip_address"),
	)
}

func checkIccpPeerExists(
	_ context.Context, ipAddress string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols iccp peer " + ipAddress + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *iccpPeerData) fillID() {
	rscData.ID = types.StringValue(rscData.IPAddress.ValueString())
}

func (rscData *iccpPeerData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *iccpPeerData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set protocols iccp peer " + rscData.IPAddress.ValueString() + " "

	for _, v := range rscData.RedundancyGroupIDList {
		configSet = append(configSet, setPrefix+"redundancy-group-id-list "+
			utils.ConvI64toa(v.ValueInt64()))
	}
	if v := rscData.AuthenticationKey.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication-key \""+v+"\"")
	}
	if v := rscData.LocalIPAddr.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"local-ip-addr "+v)
	}
	if !rscData.SessionEstablishmentHoldTime.IsNull() {
		configSet = append(configSet, setPrefix+"session-establishment-hold-time "+
			utils.ConvI64toa(rscData.SessionEstablishmentHoldTime.ValueInt64()))
	}
	if rscData.BackupLivenessDetection != nil {
		configSet = append(configSet, setPrefix+"backup-liveness-detection")

		if v := rscData.BackupLivenessDetection.BackupPeerIP.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"backup-liveness-detection backup-peer-ip "+v)
		}
	}
	if rscData.LivenessDetection != nil {
		configSet = append(configSet, setPrefix+"liveness-detection")

		if !rscData.LivenessDetection.DetectionTimeThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"liveness-detection detection-time threshold "+
				utils.ConvI64toa(rscData.LivenessDetection.DetectionTimeThreshold.ValueInt64()))
		}
		if !rscData.LivenessDetection.MinimumInterval.IsNull() {
			configSet = append(configSet, setPrefix+"liveness-detection minimum-interval "+
				utils.ConvI64toa(rscData.LivenessDetection.MinimumInterval.ValueInt64()))
		}
		if !rscData.LivenessDetection.MinimumReceiveInterval.IsNull() {
			configSet = append(configSet, setPrefix+"liveness-detection minimum-receive-interval "+
				utils.ConvI64toa(rscData.LivenessDetection.MinimumReceiveInterval.ValueInt64()))
		}
		if !rscData.LivenessDetection.Multiplier.IsNull() {
			configSet = append(configSet, setPrefix+"liveness-detection multiplier "+
				utils.ConvI64toa(rscData.LivenessDetection.Multiplier.ValueInt64()))
		}
		if rscData.LivenessDetection.NoAdaptation.ValueBool() {
			configSet = append(configSet, setPrefix+"liveness-detection no-adaptation")
		}
		if !rscData.LivenessDetection.TransmitIntervalMinimumInterval.IsNull() {
			configSet = append(configSet, setPrefix+"liveness-detection transmit-interval minimum-interval "+
				utils.ConvI64toa(rscData.LivenessDetection.TransmitIntervalMinimumInterval.ValueInt64()))
		}
		if !rscData.LivenessDetection.Multiplier.IsNull() {
			configSet = append(configSet, setPrefix+"liveness-detection transmit-interval threshold "+
				utils.ConvI64toa(rscData.LivenessDetection.TransmitIntervalThreshold.ValueInt64()))
		}
		if v := rscData.LivenessDetection.Version.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"liveness-detection version "+v)
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *iccpPeerData) read(
	_ context.Context, ipAddress string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols iccp peer " + ipAddress + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.IPAddress = types.StringValue(ipAddress)
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
			case balt.CutPrefixInString(&itemTrim, "redundancy-group-id-list "):
				value, err := tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
				rscData.RedundancyGroupIDList = append(rscData.RedundancyGroupIDList, value)
			case balt.CutPrefixInString(&itemTrim, "authentication-key "):
				rscData.AuthenticationKey, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "local-ip-addr "):
				rscData.LocalIPAddr = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "session-establishment-hold-time "):
				rscData.SessionEstablishmentHoldTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "backup-liveness-detection"):
				if rscData.BackupLivenessDetection == nil {
					rscData.BackupLivenessDetection = &iccpPeerBlockBackupLivenessDetection{}
				}
				if balt.CutPrefixInString(&itemTrim, " backup-peer-ip ") {
					rscData.BackupLivenessDetection.BackupPeerIP = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "liveness-detection"):
				if rscData.LivenessDetection == nil {
					rscData.LivenessDetection = &iccpPeerBlockLivenessDetection{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " detection-time threshold "):
					rscData.LivenessDetection.DetectionTimeThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " minimum-interval "):
					rscData.LivenessDetection.MinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " minimum-receive-interval "):
					rscData.LivenessDetection.MinimumReceiveInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " multiplier "):
					rscData.LivenessDetection.Multiplier, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " no-adaptation"):
					rscData.LivenessDetection.NoAdaptation = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, " transmit-interval minimum-interval "):
					rscData.LivenessDetection.TransmitIntervalMinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " transmit-interval threshold "):
					rscData.LivenessDetection.TransmitIntervalThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " version "):
					rscData.LivenessDetection.Version = types.StringValue(itemTrim)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *iccpPeerData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete protocols iccp peer " + rscData.IPAddress.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
