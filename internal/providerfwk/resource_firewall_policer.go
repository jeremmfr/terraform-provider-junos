package providerfwk

import (
	"context"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                   = &firewallPolicer{}
	_ resource.ResourceWithConfigure      = &firewallPolicer{}
	_ resource.ResourceWithValidateConfig = &firewallPolicer{}
	_ resource.ResourceWithImportState    = &firewallPolicer{}
	_ resource.ResourceWithUpgradeState   = &firewallPolicer{}
)

type firewallPolicer struct {
	client *junos.Client
}

func newFirewallPolicerResource() resource.Resource {
	return &firewallPolicer{}
}

func (rsc *firewallPolicer) typeName() string {
	return providerName + "_firewall_policer"
}

func (rsc *firewallPolicer) junosName() string {
	return "firewall policer"
}

func (rsc *firewallPolicer) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *firewallPolicer) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *firewallPolicer) Configure(
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

func (rsc *firewallPolicer) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
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
				Description: "Policer name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"filter_specific": schema.BoolAttribute{
				Optional:    true,
				Description: "Policer is filter-specific.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"logical_bandwidth_policer": schema.BoolAttribute{
				Optional:    true,
				Description: "Policer uses logical interface bandwidth.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"logical_interface_policer": schema.BoolAttribute{
				Optional:    true,
				Description: "Policer is logical interface policer.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"physical_interface_policer": schema.BoolAttribute{
				Optional:    true,
				Description: "Policer is physical interface policer.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"shared_bandwidth_policer": schema.BoolAttribute{
				Optional:    true,
				Description: "Share policer bandwidth among bundle links.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"if_exceeding": schema.SingleNestedBlock{
				Description: "Define rate limits options.",
				Attributes: map[string]schema.Attribute{
					"burst_size_limit": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Burst size limit in bytes.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^(\d)+(m|k|g)?$`),
								`must be a bandwidth ^(\d)+(m|k|g)?$`),
						},
					},
					"bandwidth_limit": schema.StringAttribute{
						Optional:    true,
						Description: "Bandwidth limit in bits/second.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^(\d)+(m|k|g)?$`),
								`must be a bandwidth ^(\d)+(m|k|g)?$`),
						},
					},
					"bandwidth_percent": schema.Int64Attribute{
						Optional:    true,
						Description: "Bandwidth limit in percentage (1..100 percent).",
						Validators: []validator.Int64{
							int64validator.Between(1, 100),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"if_exceeding_pps": schema.SingleNestedBlock{
				Description: "Define pps limits options.",
				Attributes: map[string]schema.Attribute{
					"packet_burst": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "PPS burst size limit.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^(\d)+(m|k|g)?$`),
								`must be a pps ^(\d)+(m|k|g)?$`),
						},
					},
					"pps_limit": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "PPS limit.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^(\d)+(m|k|g)?$`),
								`must be a pps ^(\d)+(m|k|g)?$`),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"then": schema.SingleNestedBlock{
				Description: "Define action to take if the rate limits are exceeded.",
				Attributes: map[string]schema.Attribute{
					"discard": schema.BoolAttribute{
						Optional:    true,
						Description: "Discard the packet.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"forwarding_class": schema.StringAttribute{
						Optional:    true,
						Description: "Classify packet to forwarding class.",
						Validators: []validator.String{
							tfvalidator.StringFormat(tfvalidator.DefaultFormat),
						},
					},
					"loss_priority": schema.StringAttribute{
						Optional:    true,
						Description: "Packet's loss priority.",
						Validators: []validator.String{
							stringvalidator.OneOf("high", "low", "medium-high", "medium-low"),
						},
					},
					"out_of_profile": schema.BoolAttribute{
						Optional:    true,
						Description: "Discard packets only if both congested and over threshold.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
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

type firewallPolicerData struct {
	ID                       types.String                        `tfsdk:"id"`
	Name                     types.String                        `tfsdk:"name"`
	FilterSpecific           types.Bool                          `tfsdk:"filter_specific"`
	LogicalBandwidthPolicer  types.Bool                          `tfsdk:"logical_bandwidth_policer"`
	LogicalInterfacePolicer  types.Bool                          `tfsdk:"logical_interface_policer"`
	PhysicalInterfacePolicer types.Bool                          `tfsdk:"physical_interface_policer"`
	SharedBandwidthPolicer   types.Bool                          `tfsdk:"shared_bandwidth_policer"`
	IfExceeding              *firewallPolicerBlockIfExceeding    `tfsdk:"if_exceeding"`
	IfExceedingPPS           *firewallPolicerBlockIfExceedingPPS `tfsdk:"if_exceeding_pps"`
	Then                     *firewallPolicerBlockThen           `tfsdk:"then"`
}

type firewallPolicerBlockIfExceeding struct {
	BurstSizeLimit   types.String `tfsdk:"burst_size_limit"`
	BandwidthPercent types.Int64  `tfsdk:"bandwidth_percent"`
	BandwidthLimit   types.String `tfsdk:"bandwidth_limit"`
}

func (block *firewallPolicerBlockIfExceeding) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type firewallPolicerBlockIfExceedingPPS struct {
	PacketBurst types.String `tfsdk:"packet_burst"`
	PPSLimit    types.String `tfsdk:"pps_limit"`
}

func (block *firewallPolicerBlockIfExceedingPPS) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type firewallPolicerBlockThen struct {
	Discard         types.Bool   `tfsdk:"discard"`
	ForwardingClass types.String `tfsdk:"forwarding_class"`
	LossPriority    types.String `tfsdk:"loss_priority"`
	OutOfProfile    types.Bool   `tfsdk:"out_of_profile"`
}

func (block *firewallPolicerBlockThen) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (rsc *firewallPolicer) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config firewallPolicerData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.PhysicalInterfacePolicer.IsNull() && !config.PhysicalInterfacePolicer.IsUnknown() {
		if !config.FilterSpecific.IsNull() && !config.FilterSpecific.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("logical_bandwidth_policer"),
				tfdiag.ConflictConfigErrSummary,
				"filter_specific and physical_interface_policer cannot be configured together",
			)
		}
		if !config.LogicalBandwidthPolicer.IsNull() && !config.LogicalBandwidthPolicer.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("logical_bandwidth_policer"),
				tfdiag.ConflictConfigErrSummary,
				"logical_bandwidth_policer and physical_interface_policer cannot be configured together",
			)
		}
		if !config.LogicalInterfacePolicer.IsNull() && !config.LogicalInterfacePolicer.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("logical_interface_policer"),
				tfdiag.ConflictConfigErrSummary,
				"logical_interface_policer and physical_interface_policer cannot be configured together",
			)
		}
	}

	if config.IfExceeding == nil &&
		config.IfExceedingPPS == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"one of if_exceeding or if_exceeding_pps block must be specified",
		)
	}
	if config.IfExceeding != nil && config.IfExceeding.hasKnownValue() &&
		config.IfExceedingPPS != nil && config.IfExceedingPPS.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("if_exceeding").AtName("*"),
			tfdiag.ConflictConfigErrSummary,
			"only one of if_exceeding or if_exceeding_pps block can be specified",
		)
	}
	if config.IfExceeding != nil {
		if config.IfExceeding.BurstSizeLimit.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("if_exceeding").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"burst_size_limit must be specified in if_exceeding block",
			)
		}
		if !config.IfExceeding.BandwidthLimit.IsNull() && !config.IfExceeding.BandwidthLimit.IsUnknown() &&
			!config.IfExceeding.BandwidthPercent.IsNull() && !config.IfExceeding.BandwidthPercent.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("if_exceeding").AtName("bandwidth_percent"),
				tfdiag.ConflictConfigErrSummary,
				"bandwidth_percent and bandwidth_limit cannot be configured together "+
					"in if_exceeding block",
			)
		}
	}
	if config.IfExceedingPPS != nil {
		if config.IfExceedingPPS.PacketBurst.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("if_exceeding_pps").AtName("packet_burst"),
				tfdiag.MissingConfigErrSummary,
				"packet_burst must be specified in if_exceeding_pps block",
			)
		}
		if config.IfExceedingPPS.PPSLimit.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("if_exceeding_pps").AtName("pps_limit"),
				tfdiag.MissingConfigErrSummary,
				"pps_limit must be specified in if_exceeding_pps block",
			)
		}
	}
	if config.Then != nil {
		if config.Then.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("then").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"then block is empty",
			)
		}
		if !config.Then.Discard.IsNull() && !config.Then.Discard.IsUnknown() {
			if !config.Then.ForwardingClass.IsNull() && !config.Then.ForwardingClass.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("forwarding_class"),
					tfdiag.ConflictConfigErrSummary,
					"discard and forwarding_class cannot be configured together "+
						"in then block",
				)
			}
			if !config.Then.LossPriority.IsNull() && !config.Then.LossPriority.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("loss_priority"),
					tfdiag.ConflictConfigErrSummary,
					"discard and loss_priority cannot be configured together "+
						"in then block",
				)
			}
			if !config.Then.OutOfProfile.IsNull() && !config.Then.OutOfProfile.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("out_of_profile"),
					tfdiag.ConflictConfigErrSummary,
					"discard and out_of_profile cannot be configured together "+
						"in then block",
				)
			}
		}
	} else {
		resp.Diagnostics.AddAttributeError(
			path.Root("then"),
			tfdiag.MissingConfigErrSummary,
			"then block must be specified",
		)
	}
}

func (rsc *firewallPolicer) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan firewallPolicerData
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
			policerExists, err := checkFirewallPolicerExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if policerExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			policerExists, err := checkFirewallPolicerExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !policerExists {
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

func (rsc *firewallPolicer) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data firewallPolicerData
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

func (rsc *firewallPolicer) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state firewallPolicerData
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

func (rsc *firewallPolicer) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state firewallPolicerData
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

func (rsc *firewallPolicer) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data firewallPolicerData

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

func checkFirewallPolicerExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"firewall policer \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *firewallPolicerData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *firewallPolicerData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *firewallPolicerData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set firewall policer \"" + rscData.Name.ValueString() + "\" "

	if rscData.FilterSpecific.ValueBool() {
		configSet = append(configSet, setPrefix+"filter-specific")
	}
	if rscData.LogicalBandwidthPolicer.ValueBool() {
		configSet = append(configSet, setPrefix+"logical-bandwidth-policer")
	}
	if rscData.LogicalInterfacePolicer.ValueBool() {
		configSet = append(configSet, setPrefix+"logical-interface-policer")
	}
	if rscData.PhysicalInterfacePolicer.ValueBool() {
		configSet = append(configSet, setPrefix+"physical-interface-policer")
	}
	if rscData.SharedBandwidthPolicer.ValueBool() {
		configSet = append(configSet, setPrefix+"shared-bandwidth-policer")
	}
	if rscData.IfExceeding != nil {
		configSet = append(configSet, setPrefix+"if-exceeding burst-size-limit "+
			rscData.IfExceeding.BurstSizeLimit.ValueString())
		if !rscData.IfExceeding.BandwidthPercent.IsNull() {
			configSet = append(configSet, setPrefix+"if-exceeding bandwidth-percent "+
				utils.ConvI64toa(rscData.IfExceeding.BandwidthPercent.ValueInt64()))
		}
		if v := rscData.IfExceeding.BandwidthLimit.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"if-exceeding bandwidth-limit "+v)
		}
	}
	if rscData.IfExceedingPPS != nil {
		configSet = append(configSet, setPrefix+"if-exceeding-pps packet-burst "+
			rscData.IfExceedingPPS.PacketBurst.ValueString())
		configSet = append(configSet, setPrefix+"if-exceeding-pps pps-limit "+
			rscData.IfExceedingPPS.PPSLimit.ValueString())
	}
	if rscData.Then != nil {
		if rscData.Then.Discard.ValueBool() {
			configSet = append(configSet, setPrefix+"then discard")
		}
		if v := rscData.Then.ForwardingClass.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"then forwarding-class "+v)
		}
		if v := rscData.Then.LossPriority.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"then loss-priority "+v)
		}
		if rscData.Then.OutOfProfile.ValueBool() {
			configSet = append(configSet, setPrefix+"then out-of-profile")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *firewallPolicerData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"firewall policer \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case itemTrim == "filter-specific":
				rscData.FilterSpecific = types.BoolValue(true)
			case itemTrim == "logical-bandwidth-policer":
				rscData.LogicalBandwidthPolicer = types.BoolValue(true)
			case itemTrim == "logical-interface-policer":
				rscData.LogicalInterfacePolicer = types.BoolValue(true)
			case itemTrim == "physical-interface-policer":
				rscData.PhysicalInterfacePolicer = types.BoolValue(true)
			case itemTrim == "shared-bandwidth-policer":
				rscData.SharedBandwidthPolicer = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "if-exceeding "):
				if rscData.IfExceeding == nil {
					rscData.IfExceeding = &firewallPolicerBlockIfExceeding{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "burst-size-limit "):
					rscData.IfExceeding.BurstSizeLimit = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "bandwidth-percent "):
					rscData.IfExceeding.BandwidthPercent, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "bandwidth-limit "):
					rscData.IfExceeding.BandwidthLimit = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "if-exceeding-pps "):
				if rscData.IfExceedingPPS == nil {
					rscData.IfExceedingPPS = &firewallPolicerBlockIfExceedingPPS{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "packet-burst "):
					rscData.IfExceedingPPS.PacketBurst = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "pps-limit "):
					rscData.IfExceedingPPS.PPSLimit = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "then "):
				if rscData.Then == nil {
					rscData.Then = &firewallPolicerBlockThen{}
				}
				switch {
				case itemTrim == "discard":
					rscData.Then.Discard = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
					rscData.Then.ForwardingClass = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "loss-priority "):
					rscData.Then.LossPriority = types.StringValue(itemTrim)
				case itemTrim == "out-of-profile":
					rscData.Then.OutOfProfile = types.BoolValue(true)
				}
			}
		}
	}

	return nil
}

func (rscData *firewallPolicerData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete firewall policer \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
