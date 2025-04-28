package providerfwk

import (
	"context"
	"errors"
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
	_ resource.Resource                   = &routingOptions{}
	_ resource.ResourceWithConfigure      = &routingOptions{}
	_ resource.ResourceWithValidateConfig = &routingOptions{}
	_ resource.ResourceWithImportState    = &routingOptions{}
	_ resource.ResourceWithUpgradeState   = &routingOptions{}
)

type routingOptions struct {
	client *junos.Client
}

func newRoutingOptionsResource() resource.Resource {
	return &routingOptions{}
}

func (rsc *routingOptions) typeName() string {
	return providerName + "_routing_options"
}

func (rsc *routingOptions) junosName() string {
	return "routing-options"
}

func (rsc *routingOptions) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *routingOptions) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *routingOptions) Configure(
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

func (rsc *routingOptions) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `routing_options`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clean_on_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "Clean supported lines when destroy this resource.",
			},
			"forwarding_table_export_configure_singly": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable management of `forwarding-table export` in this resource.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"instance_export": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Export policy for instance RIBs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"instance_import": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Import policy for instance RIBs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"ipv6_router_id": schema.StringAttribute{
				Optional:    true,
				Description: "IPv6 router identifier.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv6Only(),
				},
			},
			"router_id": schema.StringAttribute{
				Optional:    true,
				Description: "Router identifier.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"autonomous_system": schema.SingleNestedBlock{
				Description: "Declare `autonomous-system` configuration.",
				Attributes: map[string]schema.Attribute{
					"number": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Autonomous system number.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^\d+(\.\d+)?$`),
								"must be in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format"),
						},
					},
					"asdot_notation": schema.BoolAttribute{
						Optional:    true,
						Description: "Use AS-Dot notation to display true 4 byte AS numbers.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"loops": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum number of times this AS can be in an AS path.",
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"forwarding_table": schema.SingleNestedBlock{
				Description: "Declare `forwarding-table` configuration.",
				Attributes: map[string]schema.Attribute{
					"chain_composite_max_label_count": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum labels inside chain composite for the platform.",
						Validators: []validator.Int64{
							int64validator.Between(1, 8),
						},
					},
					"chained_composite_next_hop_ingress": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Next-hop chaining mode -> Ingress LSP nexthop settings.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.NoNullValues(),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							),
						},
					},
					"chained_composite_next_hop_transit": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Next-hop chaining mode -> Transit LSP nexthops settings.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.NoNullValues(),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							),
						},
					},
					"dynamic_list_next_hop": schema.BoolAttribute{
						Optional:    true,
						Description: "Dynamic next-hop mode for EVPN.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ecmp_fast_reroute": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable fast reroute for ECMP next hops.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_ecmp_fast_reroute": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable fast reroute for ECMP next hops.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"export": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Export policy.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.NoNullValues(),
							listvalidator.ValueStringsAre(
								stringvalidator.LengthBetween(1, 250),
								tfvalidator.StringDoubleQuoteExclusion(),
							),
						},
					},
					"indirect_next_hop": schema.BoolAttribute{
						Optional:    true,
						Description: "Install indirect next hops in Packet Forwarding Engine.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_indirect_next_hop": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't install indirect next hops in Packet Forwarding Engine.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"indirect_next_hop_change_acknowledgements": schema.BoolAttribute{
						Optional:    true,
						Description: "Request acknowledgements for Indirect next hop changes.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_indirect_next_hop_change_acknowledgements": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't request acknowledgements for Indirect next hop changes.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"krt_nexthop_ack_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Kernel nexthop ack timeout interval.",
						Validators: []validator.Int64{
							int64validator.Between(1, 400),
						},
					},
					"remnant_holdtime": schema.Int64Attribute{
						Optional:    true,
						Description: "Time to hold inherited routes from FIB.",
						Validators: []validator.Int64{
							int64validator.Between(0, 10000),
						},
					},
					"unicast_reverse_path": schema.StringAttribute{
						Optional:    true,
						Description: "Unicast reverse path (RP) verification.",
						Validators: []validator.String{
							stringvalidator.OneOf("active-paths", "feasible-paths"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"graceful_restart": schema.SingleNestedBlock{
				Description: "Graceful or hitless routing restart options.",
				Attributes: map[string]schema.Attribute{
					"disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable graceful restart.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"restart_duration": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum time for which router is in graceful restart.",
						Validators: []validator.Int64{
							int64validator.Between(120, 10000),
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

//nolint:lll
type routingOptionsData struct {
	ID                                   types.String                         `tfsdk:"id"`
	CleanOnDestroy                       types.Bool                           `tfsdk:"clean_on_destroy"`
	ForwardingTableExportConfigureSingly types.Bool                           `tfsdk:"forwarding_table_export_configure_singly"`
	InstanceExport                       []types.String                       `tfsdk:"instance_export"`
	InstanceImport                       []types.String                       `tfsdk:"instance_import"`
	IPv6RouterID                         types.String                         `tfsdk:"ipv6_router_id"`
	RouterID                             types.String                         `tfsdk:"router_id"`
	AutonomousSystem                     *routingOptionsBlockAutonomousSystem `tfsdk:"autonomous_system"`
	ForwardingTable                      *routingOptionsBlockForwardingTable  `tfsdk:"forwarding_table"`
	GracefulRestart                      *routingOptionsBlockGracefulRestart  `tfsdk:"graceful_restart"`
}

//nolint:lll
type routingOptionsConfig struct {
	ID                                   types.String                              `tfsdk:"id"`
	CleanOnDestroy                       types.Bool                                `tfsdk:"clean_on_destroy"`
	ForwardingTableExportConfigureSingly types.Bool                                `tfsdk:"forwarding_table_export_configure_singly"`
	InstanceExport                       types.List                                `tfsdk:"instance_export"`
	InstanceImport                       types.List                                `tfsdk:"instance_import"`
	IPv6RouterID                         types.String                              `tfsdk:"ipv6_router_id"`
	RouterID                             types.String                              `tfsdk:"router_id"`
	AutonomousSystem                     *routingOptionsBlockAutonomousSystem      `tfsdk:"autonomous_system"`
	ForwardingTable                      *routingOptionsBlockForwardingTableConfig `tfsdk:"forwarding_table"`
	GracefulRestart                      *routingOptionsBlockGracefulRestart       `tfsdk:"graceful_restart"`
}

type routingOptionsBlockAutonomousSystem struct {
	Number        types.String `tfsdk:"number"`
	ASdotNotation types.Bool   `tfsdk:"asdot_notation"`
	Loops         types.Int64  `tfsdk:"loops"`
}

type routingOptionsBlockForwardingTable struct {
	ChainCompositeMaxLabelCount             types.Int64    `tfsdk:"chain_composite_max_label_count"`
	ChainedCompositeNextHopIngress          []types.String `tfsdk:"chained_composite_next_hop_ingress"`
	ChainedCompositeNextHopTransit          []types.String `tfsdk:"chained_composite_next_hop_transit"`
	DynamicListNextHop                      types.Bool     `tfsdk:"dynamic_list_next_hop"`
	EcmpFastReroute                         types.Bool     `tfsdk:"ecmp_fast_reroute"`
	NoEcmpFastReroute                       types.Bool     `tfsdk:"no_ecmp_fast_reroute"`
	Export                                  []types.String `tfsdk:"export"`
	IndirectNextHop                         types.Bool     `tfsdk:"indirect_next_hop"`
	NoIndirectNextHop                       types.Bool     `tfsdk:"no_indirect_next_hop"`
	IndirectNextHopChangeAcknowledgements   types.Bool     `tfsdk:"indirect_next_hop_change_acknowledgements"`
	NoIndirectNextHopChangeAcknowledgements types.Bool     `tfsdk:"no_indirect_next_hop_change_acknowledgements"`
	KrtNexthopAckTimeout                    types.Int64    `tfsdk:"krt_nexthop_ack_timeout"`
	RemnantHoldtime                         types.Int64    `tfsdk:"remnant_holdtime"`
	UnicastReversePath                      types.String   `tfsdk:"unicast_reverse_path"`
}

func (block *routingOptionsBlockForwardingTable) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type routingOptionsBlockForwardingTableConfig struct {
	ChainCompositeMaxLabelCount             types.Int64  `tfsdk:"chain_composite_max_label_count"`
	ChainedCompositeNextHopIngress          types.Set    `tfsdk:"chained_composite_next_hop_ingress"`
	ChainedCompositeNextHopTransit          types.Set    `tfsdk:"chained_composite_next_hop_transit"`
	DynamicListNextHop                      types.Bool   `tfsdk:"dynamic_list_next_hop"`
	EcmpFastReroute                         types.Bool   `tfsdk:"ecmp_fast_reroute"`
	NoEcmpFastReroute                       types.Bool   `tfsdk:"no_ecmp_fast_reroute"`
	Export                                  types.List   `tfsdk:"export"`
	IndirectNextHop                         types.Bool   `tfsdk:"indirect_next_hop"`
	NoIndirectNextHop                       types.Bool   `tfsdk:"no_indirect_next_hop"`
	IndirectNextHopChangeAcknowledgements   types.Bool   `tfsdk:"indirect_next_hop_change_acknowledgements"`
	NoIndirectNextHopChangeAcknowledgements types.Bool   `tfsdk:"no_indirect_next_hop_change_acknowledgements"`
	KrtNexthopAckTimeout                    types.Int64  `tfsdk:"krt_nexthop_ack_timeout"`
	RemnantHoldtime                         types.Int64  `tfsdk:"remnant_holdtime"`
	UnicastReversePath                      types.String `tfsdk:"unicast_reverse_path"`
}

func (block *routingOptionsBlockForwardingTableConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type routingOptionsBlockGracefulRestart struct {
	Disable         types.Bool  `tfsdk:"disable"`
	RestartDuration types.Int64 `tfsdk:"restart_duration"`
}

func (rsc *routingOptions) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config routingOptionsConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.AutonomousSystem != nil {
		if config.AutonomousSystem.Number.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("autonomous_system").AtName("number"),
				tfdiag.MissingConfigErrSummary,
				"number must be specified in autonomous_system block",
			)
		}
	}

	if config.ForwardingTable != nil {
		if config.ForwardingTable.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("forwarding_table"),
				tfdiag.MissingConfigErrSummary,
				"forwarding_table block is empty",
			)
		}

		if !config.ForwardingTable.EcmpFastReroute.IsNull() &&
			!config.ForwardingTable.EcmpFastReroute.IsUnknown() &&
			!config.ForwardingTable.NoEcmpFastReroute.IsNull() &&
			!config.ForwardingTable.NoEcmpFastReroute.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("forwarding_table").AtName("ecmp_fast_reroute"),
				tfdiag.ConflictConfigErrSummary,
				"ecmp_fast_reroute and no_ecmp_fast_reroute cannot be configured together"+
					" in forwarding_table block",
			)
		}
		if !config.ForwardingTable.Export.IsNull() &&
			!config.ForwardingTable.Export.IsUnknown() &&
			!config.ForwardingTableExportConfigureSingly.IsNull() &&
			!config.ForwardingTableExportConfigureSingly.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("forwarding_table").AtName("export"),
				tfdiag.ConflictConfigErrSummary,
				"export in forwarding_table block and forwarding_table_export_configure_singly"+
					" cannot be configured together",
			)
		}
		if !config.ForwardingTable.IndirectNextHop.IsNull() &&
			!config.ForwardingTable.IndirectNextHop.IsUnknown() &&
			!config.ForwardingTable.NoIndirectNextHop.IsNull() &&
			!config.ForwardingTable.NoIndirectNextHop.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("forwarding_table").AtName("indirect_next_hop"),
				tfdiag.ConflictConfigErrSummary,
				"indirect_next_hop and no_indirect_next_hop cannot be configured together"+
					" in forwarding_table block",
			)
		}
		if !config.ForwardingTable.IndirectNextHopChangeAcknowledgements.IsNull() &&
			!config.ForwardingTable.IndirectNextHopChangeAcknowledgements.IsUnknown() &&
			!config.ForwardingTable.NoIndirectNextHopChangeAcknowledgements.IsNull() &&
			!config.ForwardingTable.NoIndirectNextHopChangeAcknowledgements.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("forwarding_table").AtName("indirect_next_hop_change_acknowledgements"),
				tfdiag.ConflictConfigErrSummary,
				"indirect_next_hop_change_acknowledgements and no_indirect_next_hop_change_acknowledgements"+
					" cannot be configured together"+
					" in forwarding_table block",
			)
		}
	}
}

func (rsc *routingOptions) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan routingOptionsData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		nil,
		nil,
		&plan,
		resp,
	)
}

func (rsc *routingOptions) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data routingOptionsData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadWithoutArg = &data
	defaultResourceRead(
		ctx,
		rsc,
		nil,
		&data,
		func() {
			data.CleanOnDestroy = state.CleanOnDestroy
			data.ForwardingTableExportConfigureSingly = state.ForwardingTableExportConfigureSingly
			if data.ForwardingTableExportConfigureSingly.ValueBool() {
				if data.ForwardingTable != nil {
					data.ForwardingTable.Export = nil
				}
			}
		},
		resp,
	)
}

func (rsc *routingOptions) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state routingOptionsData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	forwardingTableExportConfigureSingly := plan.ForwardingTableExportConfigureSingly.ValueBool()
	if !plan.ForwardingTableExportConfigureSingly.Equal(state.ForwardingTableExportConfigureSingly) {
		if state.ForwardingTableExportConfigureSingly.ValueBool() {
			forwardingTableExportConfigureSingly = state.ForwardingTableExportConfigureSingly.ValueBool()
			resp.Diagnostics.AddAttributeWarning(
				path.Root("forwarding_table_export_configure_singly"),
				"Disable forwarding_table_export_configure_singly on resource already created",
				"It's doesn't delete export list already configured. "+
					"So refresh resource after apply to detect export list entries that need to be deleted",
			)
		} else {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("forwarding_table_export_configure_singly"),
				"Enable forwarding_table_export_configure_singly on resource already created",
				"It's doesn't delete export list already configured. "+
					"So add `add_it_to_forwarding_table_export` argument on each `junos_policyoptions_policy_statement` "+
					"resource to be able to manage each element of the export list",
			)
		}
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.delOpts(ctx, forwardingTableExportConfigureSingly, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigUnlockWarnSummary, junSess.ConfigUnlock())...)
	}()

	if err := state.delOpts(ctx, forwardingTableExportConfigureSingly, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "update resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *routingOptions) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state routingOptionsData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.CleanOnDestroy.ValueBool() {
		defaultResourceDelete(
			ctx,
			rsc,
			&state,
			resp,
		)
	}
}

func (rsc *routingOptions) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data routingOptionsData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"",
	)
}

func (rscData *routingOptionsData) fillID() {
	rscData.ID = types.StringValue("routing_options")
}

func (rscData *routingOptionsData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *routingOptionsData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set routing-options "
	configSet := make([]string, 0, 100)

	for _, v := range rscData.InstanceExport {
		configSet = append(configSet, setPrefix+"instance-export \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.InstanceImport {
		configSet = append(configSet, setPrefix+"instance-import \""+v.ValueString()+"\"")
	}
	if v := rscData.IPv6RouterID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ipv6-router-id "+v)
	}
	if v := rscData.RouterID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"router-id "+v)
	}

	if rscData.AutonomousSystem != nil {
		configSet = append(configSet, setPrefix+"autonomous-system "+rscData.AutonomousSystem.Number.ValueString())
		if rscData.AutonomousSystem.ASdotNotation.ValueBool() {
			configSet = append(configSet, setPrefix+"autonomous-system asdot-notation")
		}
		if !rscData.AutonomousSystem.Loops.IsNull() {
			configSet = append(configSet, setPrefix+"autonomous-system loops "+
				utils.ConvI64toa(rscData.AutonomousSystem.Loops.ValueInt64()))
		}
	}
	if rscData.ForwardingTable != nil {
		if rscData.ForwardingTable.isEmpty() {
			return path.Root("forwarding_table"),
				errors.New("forwarding_table block is empty")
		}
		if rscData.ForwardingTableExportConfigureSingly.ValueBool() &&
			len(rscData.ForwardingTable.Export) > 0 {
			return path.Root("forwarding_table").AtName("export"),
				errors.New("export in forwarding_table block and forwarding_table_export_configure_singly" +
					" cannot be configured together")
		}

		configSet = append(configSet, rscData.ForwardingTable.configSet()...)
	}
	if rscData.GracefulRestart != nil {
		configSet = append(configSet, setPrefix+"graceful-restart")

		if rscData.GracefulRestart.Disable.ValueBool() {
			configSet = append(configSet, setPrefix+"graceful-restart disable")
		}
		if !rscData.GracefulRestart.RestartDuration.IsNull() {
			configSet = append(configSet, setPrefix+"graceful-restart restart-duration "+
				utils.ConvI64toa(rscData.GracefulRestart.RestartDuration.ValueInt64()))
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *routingOptionsBlockForwardingTable) configSet() []string {
	setPrefix := "set routing-options forwarding-table "
	configSet := make([]string, 0, 100)

	if !block.ChainCompositeMaxLabelCount.IsNull() {
		configSet = append(configSet, setPrefix+"chain-composite-max-label-count "+
			utils.ConvI64toa(block.ChainCompositeMaxLabelCount.ValueInt64()))
	}
	for _, v := range block.ChainedCompositeNextHopIngress {
		configSet = append(configSet, setPrefix+"chained-composite-next-hop ingress "+v.ValueString())
	}
	for _, v := range block.ChainedCompositeNextHopTransit {
		configSet = append(configSet, setPrefix+"chained-composite-next-hop transit "+v.ValueString())
	}
	if block.DynamicListNextHop.ValueBool() {
		configSet = append(configSet, setPrefix+"dynamic-list-next-hop")
	}
	if block.EcmpFastReroute.ValueBool() {
		configSet = append(configSet, setPrefix+"ecmp-fast-reroute")
	}
	if block.NoEcmpFastReroute.ValueBool() {
		configSet = append(configSet, setPrefix+"no-ecmp-fast-reroute")
	}
	for _, v := range block.Export {
		configSet = append(configSet, setPrefix+"export \""+v.ValueString()+"\"")
	}
	if block.IndirectNextHop.ValueBool() {
		configSet = append(configSet, setPrefix+"indirect-next-hop")
	}
	if block.NoIndirectNextHop.ValueBool() {
		configSet = append(configSet, setPrefix+"no-indirect-next-hop")
	}
	if block.IndirectNextHopChangeAcknowledgements.ValueBool() {
		configSet = append(configSet, setPrefix+"indirect-next-hop-change-acknowledgements")
	}
	if block.NoIndirectNextHopChangeAcknowledgements.ValueBool() {
		configSet = append(configSet, setPrefix+"no-indirect-next-hop-change-acknowledgements")
	}
	if !block.KrtNexthopAckTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"krt-nexthop-ack-timeout "+
			utils.ConvI64toa(block.KrtNexthopAckTimeout.ValueInt64()))
	}
	if !block.RemnantHoldtime.IsNull() {
		configSet = append(configSet, setPrefix+"remnant-holdtime "+
			utils.ConvI64toa(block.RemnantHoldtime.ValueInt64()))
	}
	if v := block.UnicastReversePath.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"unicast-reverse-path "+v)
	}

	return configSet
}

func (rscData *routingOptionsData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"routing-options" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
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
			case balt.CutPrefixInString(&itemTrim, "autonomous-system "):
				if rscData.AutonomousSystem == nil {
					rscData.AutonomousSystem = &routingOptionsBlockAutonomousSystem{}
				}

				switch {
				case balt.CutPrefixInString(&itemTrim, "loops "):
					rscData.AutonomousSystem.Loops, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case itemTrim == "asdot-notation":
					rscData.AutonomousSystem.ASdotNotation = types.BoolValue(true)
				default:
					rscData.AutonomousSystem.Number = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "forwarding-table "):
				if rscData.ForwardingTable == nil {
					rscData.ForwardingTable = &routingOptionsBlockForwardingTable{}
				}

				if err := rscData.ForwardingTable.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "graceful-restart"):
				if rscData.GracefulRestart == nil {
					rscData.GracefulRestart = &routingOptionsBlockGracefulRestart{}
				}

				switch {
				case itemTrim == " disable":
					rscData.GracefulRestart.Disable = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, " restart-duration "):
					rscData.GracefulRestart.RestartDuration, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "instance-export "):
				rscData.InstanceExport = append(rscData.InstanceExport, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "instance-import "):
				rscData.InstanceImport = append(rscData.InstanceImport, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "ipv6-router-id "):
				rscData.IPv6RouterID = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "router-id "):
				rscData.RouterID = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (block *routingOptionsBlockForwardingTable) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "chain-composite-max-label-count "):
		block.ChainCompositeMaxLabelCount, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "chained-composite-next-hop ingress "):
		block.ChainedCompositeNextHopIngress = append(block.ChainedCompositeNextHopIngress,
			types.StringValue(itemTrim),
		)
	case balt.CutPrefixInString(&itemTrim, "chained-composite-next-hop transit "):
		block.ChainedCompositeNextHopTransit = append(block.ChainedCompositeNextHopTransit,
			types.StringValue(itemTrim),
		)
	case itemTrim == "dynamic-list-next-hop":
		block.DynamicListNextHop = types.BoolValue(true)
	case itemTrim == "ecmp-fast-reroute":
		block.EcmpFastReroute = types.BoolValue(true)
	case itemTrim == "no-ecmp-fast-reroute":
		block.NoEcmpFastReroute = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "export "):
		block.Export = append(block.Export,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case itemTrim == "indirect-next-hop":
		block.IndirectNextHop = types.BoolValue(true)
	case itemTrim == "no-indirect-next-hop":
		block.NoIndirectNextHop = types.BoolValue(true)
	case itemTrim == "indirect-next-hop-change-acknowledgements":
		block.IndirectNextHopChangeAcknowledgements = types.BoolValue(true)
	case itemTrim == "no-indirect-next-hop-change-acknowledgements":
		block.NoIndirectNextHopChangeAcknowledgements = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "krt-nexthop-ack-timeout "):
		block.KrtNexthopAckTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "remnant-holdtime "):
		block.RemnantHoldtime, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "unicast-reverse-path "):
		block.UnicastReversePath = types.StringValue(itemTrim)
	}

	return err
}

func (rscData *routingOptionsData) delOpts(
	_ context.Context, forwardingTableExportConfigureSingly bool, junSess *junos.Session,
) error {
	listLinesToDelete := []string{
		"autonomous-system",
		"graceful-restart",
		"instance-export",
		"instance-import",
		"ipv6-router-id",
		"router-id",
	}
	if forwardingTableExportConfigureSingly {
		listLinesToDeleteFwTable := []string{
			"forwarding-table chain-composite-max-label-count",
			"forwarding-table chained-composite-next-hop",
			"forwarding-table dynamic-list-next-hop",
			"forwarding-table ecmp-fast-reroute",
			"forwarding-table no-ecmp-fast-reroute",
			"forwarding-table indirect-next-hop",
			"forwarding-table no-indirect-next-hop",
			"forwarding-table indirect-next-hop-change-acknowledgements",
			"forwarding-table no-indirect-next-hop-change-acknowledgements",
			"forwarding-table krt-nexthop-ack-timeout",
			"forwarding-table remnant-holdtime",
			"forwarding-table unicast-reverse-path",
		}
		listLinesToDelete = append(listLinesToDelete, listLinesToDeleteFwTable...)
	} else {
		listLinesToDelete = append(listLinesToDelete, "forwarding-table")
	}

	configSet := make([]string, len(listLinesToDelete))
	delPrefix := "delete routing-options "
	for i, line := range listLinesToDelete {
		configSet[i] = delPrefix + line
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *routingOptionsData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	return rscData.delOpts(ctx, rscData.ForwardingTableExportConfigureSingly.ValueBool(), junSess)
}
