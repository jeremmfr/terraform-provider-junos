package provider

import (
	"context"
	"fmt"
	"strconv"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &chassisCluster{}
	_ resource.ResourceWithConfigure      = &chassisCluster{}
	_ resource.ResourceWithValidateConfig = &chassisCluster{}
	_ resource.ResourceWithImportState    = &chassisCluster{}
	_ resource.ResourceWithUpgradeState   = &chassisCluster{}
)

type chassisCluster struct {
	client *junos.Client
}

func newChassisClusterResource() resource.Resource {
	return &chassisCluster{}
}

func (rsc *chassisCluster) typeName() string {
	return providerName + "_chassis_cluster"
}

func (rsc *chassisCluster) junosName() string {
	return "chassis cluster"
}

func (rsc *chassisCluster) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *chassisCluster) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *chassisCluster) Configure(
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

func (rsc *chassisCluster) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version: 1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block" +
			" and configure fab0 and fab1 interfaces.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with value " +
					"`cluster`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reth_count": schema.Int64Attribute{
				Required:    true,
				Description: "Number of redundant ethernet interfaces.",
				Validators: []validator.Int64{
					int64validator.Between(1, 128),
				},
			},
			"config_sync_no_secondary_bootup_auto": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable auto configuration synchronize on secondary bootup.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"control_link_recovery": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable automatic control link recovery.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"heartbeat_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval between successive heartbeats (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(1000, 2000),
				},
			},
			"heartbeat_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of consecutive missed heartbeats to indicate device failure.",
				Validators: []validator.Int64{
					int64validator.Between(3, 8),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"fab0": schema.SingleNestedBlock{
				Description: "Declare `interfaces fab0` configuration.",
				Attributes: map[string]schema.Attribute{
					"member_interfaces": schema.ListAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "Member interfaces for the fabric interface.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.NoNullValues(),
							listvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
								tfvalidator.StringDotExclusion(),
							),
						},
					},
					"description": schema.StringAttribute{
						Optional:    true,
						Description: "Text description of interface.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 900),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"fab1": schema.SingleNestedBlock{
				Description: "Declare `interfaces fab1` configuration.",
				Attributes: map[string]schema.Attribute{
					"member_interfaces": schema.ListAttribute{
						ElementType: types.StringType,
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Member interfaces for the fabric interface.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.NoNullValues(),
							listvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
								tfvalidator.StringDotExclusion(),
							),
						},
					},
					"description": schema.StringAttribute{
						Optional:    true,
						Description: "Text description of interface.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 900),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"redundancy_group": schema.ListNestedBlock{
				Description: "For each redundancy-group to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"node0_priority": schema.Int64Attribute{
							Required:    true,
							Description: "Priority of the node0 in the redundancy-group.",
							Validators: []validator.Int64{
								int64validator.Between(1, 254),
							},
						},
						"node1_priority": schema.Int64Attribute{
							Required:    true,
							Description: "Priority of the node1 in the redundancy-group.",
							Validators: []validator.Int64{
								int64validator.Between(1, 254),
							},
						},
						"gratuitous_arp_count": schema.Int64Attribute{
							Optional:    true,
							Description: "Number of gratuitous ARPs to send on an active interface after failover.",
							Validators: []validator.Int64{
								int64validator.Between(1, 16),
							},
						},
						"hold_down_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "RG failover interval (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(0, 1800),
							},
						},
						"preempt": schema.BoolAttribute{
							Optional:    true,
							Description: "Allow preemption of primaryship based on priority.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"preempt_delay": schema.Int64Attribute{
							Optional:    true,
							Description: "Time to wait before taking over mastership (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 21600),
							},
						},
						"preempt_limit": schema.Int64Attribute{
							Optional:    true,
							Description: "Max number of preemptive failovers allowed.",
							Validators: []validator.Int64{
								int64validator.Between(1, 50),
							},
						},
						"preempt_period": schema.Int64Attribute{
							Optional:    true,
							Description: "Time period during which the limit is applied (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 1400),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"interface_monitor": schema.ListNestedBlock{
							Description: "For each monitoring interface to declare.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required:    true,
										Description: "Name of the interface to monitor.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
											tfvalidator.StringDotExclusion(),
										},
									},
									"weight": schema.Int64Attribute{
										Required:    true,
										Description: "Weight assigned to this interface that influences failover.",
										Validators: []validator.Int64{
											int64validator.Between(0, 255),
										},
									},
								},
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(128),
				},
			},
			"control_ports": schema.SetNestedBlock{
				Description: "For each combination of block arguments," +
					" enable the specific control port to use as a control link for the chassis cluster.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"fpc": schema.Int64Attribute{
							Required:    true,
							Description: "Flexible PIC Concentrator (FPC) slot number.",
							Validators: []validator.Int64{
								int64validator.Between(0, 23),
							},
						},
						"port": schema.Int64Attribute{
							Required:    true,
							Description: "Port number on which to configure the control port.",
							Validators: []validator.Int64{
								int64validator.Between(0, 1),
							},
						},
					},
				},
			},
		},
	}
}

type chassisClusterData struct {
	ID                              types.String                         `tfsdk:"id"`
	RethCount                       types.Int64                          `tfsdk:"reth_count"`
	ConfigSyncNoSecondaryBootupAuto types.Bool                           `tfsdk:"config_sync_no_secondary_bootup_auto"`
	ControlLinkRecovery             types.Bool                           `tfsdk:"control_link_recovery"`
	HeartbeatInterval               types.Int64                          `tfsdk:"heartbeat_interval"`
	HeartbeatThreshold              types.Int64                          `tfsdk:"heartbeat_threshold"`
	Fab0                            *chassisClusterBlockFab              `tfsdk:"fab0"`
	Fab1                            *chassisClusterBlockFab              `tfsdk:"fab1"`
	RedundancyGroup                 []chassisClusterBlockRedundancyGroup `tfsdk:"redundancy_group"`
	ControlPorts                    []chassisClusterBlockControlPorts    `tfsdk:"control_ports"`
}

type chassisClusterConfig struct {
	ID                              types.String                  `tfsdk:"id"`
	RethCount                       types.Int64                   `tfsdk:"reth_count"`
	ConfigSyncNoSecondaryBootupAuto types.Bool                    `tfsdk:"config_sync_no_secondary_bootup_auto"`
	ControlLinkRecovery             types.Bool                    `tfsdk:"control_link_recovery"`
	HeartbeatInterval               types.Int64                   `tfsdk:"heartbeat_interval"`
	HeartbeatThreshold              types.Int64                   `tfsdk:"heartbeat_threshold"`
	Fab0                            *chassisClusterBlockFabConfig `tfsdk:"fab0"`
	Fab1                            *chassisClusterBlockFabConfig `tfsdk:"fab1"`
	RedundancyGroup                 types.List                    `tfsdk:"redundancy_group"`
	ControlPorts                    types.Set                     `tfsdk:"control_ports"`
}

type chassisClusterBlockFab struct {
	MemberInterfaces []types.String `tfsdk:"member_interfaces"`
	Description      types.String   `tfsdk:"description"`
}

type chassisClusterBlockFabConfig struct {
	MemberInterfaces types.List   `tfsdk:"member_interfaces"`
	Description      types.String `tfsdk:"description"`
}

type chassisClusterBlockRedundancyGroup struct {
	Node0Priority      types.Int64                                               `tfsdk:"node0_priority"`
	Node1Priority      types.Int64                                               `tfsdk:"node1_priority"`
	GratuitousArpCount types.Int64                                               `tfsdk:"gratuitous_arp_count"`
	HoldDownInterval   types.Int64                                               `tfsdk:"hold_down_interval"`
	Preempt            types.Bool                                                `tfsdk:"preempt"`
	PreemptDelay       types.Int64                                               `tfsdk:"preempt_delay"`
	PreemptLimit       types.Int64                                               `tfsdk:"preempt_limit"`
	PreemptPeriod      types.Int64                                               `tfsdk:"preempt_period"`
	InterfaceMonitor   []chassisClusterBlockRedundancyGroupBlockInterfaceMonitor `tfsdk:"interface_monitor"`
}

type chassisClusterBlockRedundancyGroupConfig struct {
	Node0Priority      types.Int64 `tfsdk:"node0_priority"`
	Node1Priority      types.Int64 `tfsdk:"node1_priority"`
	GratuitousArpCount types.Int64 `tfsdk:"gratuitous_arp_count"`
	HoldDownInterval   types.Int64 `tfsdk:"hold_down_interval"`
	Preempt            types.Bool  `tfsdk:"preempt"`
	PreemptDelay       types.Int64 `tfsdk:"preempt_delay"`
	PreemptLimit       types.Int64 `tfsdk:"preempt_limit"`
	PreemptPeriod      types.Int64 `tfsdk:"preempt_period"`
	InterfaceMonitor   types.List  `tfsdk:"interface_monitor"`
}

type chassisClusterBlockRedundancyGroupBlockInterfaceMonitor struct {
	Name   types.String `tfsdk:"name"   tfdata:"identifier"`
	Weight types.Int64  `tfsdk:"weight"`
}

type chassisClusterBlockControlPorts struct {
	Fpc  types.Int64 `tfsdk:"fpc"  tfdata:"identifier"`
	Port types.Int64 `tfsdk:"port"`
}

func (rsc *chassisCluster) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config chassisClusterConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Fab0 == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("fab0"),
			tfdiag.MissingConfigErrSummary,
			"fab0 block must be specified",
		)
	}
	if config.Fab1 != nil &&
		config.Fab1.MemberInterfaces.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("fab1").AtName("member_interfaces"),
			tfdiag.MissingConfigErrSummary,
			"member_interfaces must be specified"+
				" in fab1 block",
		)
	}
	if config.RedundancyGroup.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("redundancy_group"),
			tfdiag.MissingConfigErrSummary,
			"redundancy_group block must be specified",
		)
	} else if !config.RedundancyGroup.IsUnknown() {
		var configRedundancyGroup []chassisClusterBlockRedundancyGroupConfig
		asDiags := config.RedundancyGroup.ElementsAs(ctx, &configRedundancyGroup, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		for i, block := range configRedundancyGroup {
			if block.Preempt.IsNull() {
				if !block.PreemptDelay.IsNull() &&
					!block.PreemptDelay.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("redundancy_group").AtListIndex(i).AtName("preempt_delay"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("preempt must be specified with preempt_delay"+
							" in redundancy_group block n°%d", i),
					)
				}
				if !block.PreemptLimit.IsNull() &&
					!block.PreemptLimit.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("redundancy_group").AtListIndex(i).AtName("preempt_limit"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("preempt must be specified with preempt_limit"+
							" in redundancy_group block n°%d", i),
					)
				}
				if !block.PreemptPeriod.IsNull() &&
					!block.PreemptPeriod.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("redundancy_group").AtListIndex(i).AtName("preempt_period"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("preempt must be specified with preempt_period"+
							" in redundancy_group block n°%d", i),
					)
				}
			}
			if !block.InterfaceMonitor.IsNull() &&
				!block.InterfaceMonitor.IsUnknown() {
				var configInterfaceMonitor []chassisClusterBlockRedundancyGroupBlockInterfaceMonitor
				asDiags := block.InterfaceMonitor.ElementsAs(ctx, &configInterfaceMonitor, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				interfaceMonitorName := make(map[string]struct{})
				for ii, subBlock := range configInterfaceMonitor {
					if subBlock.Name.IsUnknown() {
						continue
					}

					name := subBlock.Name.ValueString()
					if _, ok := interfaceMonitorName[name]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("redundancy_group").AtListIndex(i).AtName("interface_monitor").AtListIndex(ii).AtName("name"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple interface_monitor blocks with the same name %q"+
								" in redundancy_group block n°%d", name, i),
						)
					}
					interfaceMonitorName[name] = struct{}{}
				}
			}
		}
	}
	if !config.ControlPorts.IsNull() &&
		!config.ControlPorts.IsUnknown() {
		var configControlPorts []chassisClusterBlockControlPorts
		asDiags := config.ControlPorts.ElementsAs(ctx, &configControlPorts, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		controlPortsFpc := make(map[int64]struct{})
		for _, block := range configControlPorts {
			if block.Fpc.IsUnknown() {
				continue
			}

			fpc := block.Fpc.ValueInt64()
			if _, ok := controlPortsFpc[fpc]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("control_ports"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple control_ports blocks with the same fpc %d", fpc),
				)
			}
			controlPortsFpc[fpc] = struct{}{}
		}
	}
}

func (rsc *chassisCluster) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan chassisClusterData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(_ context.Context, junSess *junos.Session) bool {
			if !junSess.CheckCompatibilityChassisCluster() {
				resp.Diagnostics.AddError(
					tfdiag.CompatibilityErrSummary,
					fmt.Sprintf(rsc.junosName()+" not compatible "+
						"with Junos device %q", junSess.SystemInformation.HardwareModel),
				)

				return false
			}

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *chassisCluster) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data chassisClusterData
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
		nil,
		resp,
	)
}

func (rsc *chassisCluster) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state chassisClusterData
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

func (rsc *chassisCluster) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state chassisClusterData
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

func (rsc *chassisCluster) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data chassisClusterData

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

func (rscData *chassisClusterData) fillID() {
	rscData.ID = types.StringValue("cluster")
}

func (rscData *chassisClusterData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *chassisClusterData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set chassis cluster "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix + "reth-count " +
		utils.ConvI64toa(rscData.RethCount.ValueInt64())

	if rscData.ConfigSyncNoSecondaryBootupAuto.ValueBool() {
		configSet = append(configSet, setPrefix+"configuration-synchronize no-secondary-bootup-auto")
	}
	if rscData.ControlLinkRecovery.ValueBool() {
		configSet = append(configSet, setPrefix+"control-link-recovery")
	}
	if !rscData.HeartbeatInterval.IsNull() {
		configSet = append(configSet, setPrefix+"heartbeat-interval "+
			utils.ConvI64toa(rscData.HeartbeatInterval.ValueInt64()))
	}
	if !rscData.HeartbeatThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"heartbeat-threshold "+
			utils.ConvI64toa(rscData.HeartbeatThreshold.ValueInt64()))
	}

	if rscData.Fab0 != nil {
		configSet = append(configSet, rscData.Fab0.configSet("fab0")...)
	}
	if rscData.Fab1 != nil {
		configSet = append(configSet, rscData.Fab1.configSet("fab1")...)
	}
	for i, v := range rscData.RedundancyGroup {
		blockSet, pathErr, err := v.configSet(setPrefix, i, path.Root("redundancy_group").AtListIndex(i))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	controlPortsFpc := make(map[int64]struct{})
	for _, v := range rscData.ControlPorts {
		fpc := v.Fpc.ValueInt64()
		if _, ok := controlPortsFpc[fpc]; ok {
			return path.Root("control_ports"),
				fmt.Errorf("multiple control_ports blocks with the same fpc %d", fpc)
		}
		controlPortsFpc[fpc] = struct{}{}

		configSet = append(configSet, setPrefix+"control-ports fpc "+utils.ConvI64toa(fpc)+
			" port "+utils.ConvI64toa(v.Port.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *chassisClusterBlockFab) configSet(iface string) []string {
	setPrefix := "set interfaces " + iface + " "

	configSet := make([]string, 1, 100)
	configSet[0] = "delete interfaces " + iface + " disable"

	for _, v := range block.MemberInterfaces {
		configSet = append(configSet, setPrefix+"fabric-options member-interfaces "+v.ValueString())
	}
	if v := block.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	return configSet
}

func (block *chassisClusterBlockRedundancyGroup) configSet(
	setPrefix string, index int, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "redundancy-group " + strconv.Itoa(index) + " "

	configSet := make([]string, 2, 100)
	configSet[0] = setPrefix + "node 0 priority " +
		utils.ConvI64toa(block.Node0Priority.ValueInt64())
	configSet[1] = setPrefix + "node 1 priority " +
		utils.ConvI64toa(block.Node1Priority.ValueInt64())

	if !block.GratuitousArpCount.IsNull() {
		configSet = append(configSet, setPrefix+"gratuitous-arp-count "+
			utils.ConvI64toa(block.GratuitousArpCount.ValueInt64()))
	}
	if !block.HoldDownInterval.IsNull() {
		configSet = append(configSet, setPrefix+"hold-down-interval "+
			utils.ConvI64toa(block.HoldDownInterval.ValueInt64()))
	}
	if block.Preempt.ValueBool() {
		configSet = append(configSet, setPrefix+"preempt")

		if !block.PreemptDelay.IsNull() {
			configSet = append(configSet, setPrefix+"preempt delay "+
				utils.ConvI64toa(block.PreemptDelay.ValueInt64()))
		}
		if !block.PreemptLimit.IsNull() {
			configSet = append(configSet, setPrefix+"preempt limit "+
				utils.ConvI64toa(block.PreemptLimit.ValueInt64()))
		}
		if !block.PreemptPeriod.IsNull() {
			configSet = append(configSet, setPrefix+"preempt period "+
				utils.ConvI64toa(block.PreemptPeriod.ValueInt64()))
		}
	} else {
		if !block.PreemptDelay.IsNull() {
			return configSet,
				pathRoot.AtName("preempt_delay"),
				fmt.Errorf("preempt must be specified with preempt_delay"+
					" in redundancy_group block n°%d", index)
		}
		if !block.PreemptLimit.IsNull() {
			return configSet,
				pathRoot.AtName("preempt_limit"),
				fmt.Errorf("preempt must be specified with preempt_limit"+
					" in redundancy_group block n°%d", index)
		}
		if !block.PreemptPeriod.IsNull() {
			return configSet,
				pathRoot.AtName("preempt_period"),
				fmt.Errorf("preempt must be specified with preempt_period"+
					" in redundancy_group block n°%d", index)
		}
	}

	interfaceMonitorName := make(map[string]struct{})
	for i, v := range block.InterfaceMonitor {
		name := v.Name.ValueString()
		if _, ok := interfaceMonitorName[name]; ok {
			return configSet,
				pathRoot.AtName("interface_monitor").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple interface_monitor blocks with the same name %q"+
					" in redundancy_group block n°%d", name, index)
		}
		interfaceMonitorName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"interface-monitor "+name+
			" weight "+utils.ConvI64toa(v.Weight.ValueInt64()))
	}

	return configSet, path.Empty(), nil
}

func (rscData *chassisClusterData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"chassis cluster" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "redundancy-group "):
				idStr := tfdata.FirstElementOfJunosLine(itemTrim)
				id, err := tfdata.ConvAtoi64Value(idStr)
				if err != nil {
					return err
				}
				idInt := int(id.ValueInt64())

				if len(rscData.RedundancyGroup) < idInt+1 {
					for i := len(rscData.RedundancyGroup); i < idInt+1; i++ {
						rscData.RedundancyGroup = append(rscData.RedundancyGroup, chassisClusterBlockRedundancyGroup{})
					}
				}
				redundancyGroup := &rscData.RedundancyGroup[idInt]
				balt.CutPrefixInString(&itemTrim, idStr+" ")

				if err := redundancyGroup.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "reth-count "):
				rscData.RethCount, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "configuration-synchronize no-secondary-bootup-auto":
				rscData.ConfigSyncNoSecondaryBootupAuto = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "control-ports fpc "):
				fpcStr := tfdata.FirstElementOfJunosLine(itemTrim)
				fpc, err := tfdata.ConvAtoi64Value(fpcStr)
				if err != nil {
					return err
				}
				var controlPorts chassisClusterBlockControlPorts
				rscData.ControlPorts, controlPorts = tfdata.ExtractBlock(rscData.ControlPorts, fpc)
				balt.CutPrefixInString(&itemTrim, fpcStr+" ")

				if balt.CutPrefixInString(&itemTrim, "port ") {
					controlPorts.Port, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
				rscData.ControlPorts = append(rscData.ControlPorts, controlPorts)
			case itemTrim == "control-link-recovery":
				rscData.ControlLinkRecovery = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "heartbeat-interval "):
				rscData.HeartbeatInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "heartbeat-threshold "):
				rscData.HeartbeatThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	showConfigFab0, err := junSess.Command(junos.CmdShowConfig + "interfaces fab0" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfigFab0 != junos.EmptyW {
		if rscData.Fab0 == nil {
			rscData.Fab0 = &chassisClusterBlockFab{}
		}

		for _, item := range strings.Split(showConfigFab0, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			rscData.Fab0.read(itemTrim)
		}
	}

	showConfigFab1, err := junSess.Command(junos.CmdShowConfig + "interfaces fab1" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfigFab1 != junos.EmptyW {
		if rscData.Fab1 == nil {
			rscData.Fab1 = &chassisClusterBlockFab{}
		}

		for _, item := range strings.Split(showConfigFab1, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			rscData.Fab1.read(itemTrim)
		}
	}

	return nil
}

func (block *chassisClusterBlockFab) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "description "):
		block.Description = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "fabric-options member-interfaces "):
		block.MemberInterfaces = append(block.MemberInterfaces, types.StringValue(itemTrim))
	}
}

func (block *chassisClusterBlockRedundancyGroup) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "node 0 priority "):
		block.Node0Priority, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "node 1 priority "):
		block.Node1Priority, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "gratuitous-arp-count "):
		block.GratuitousArpCount, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "hold-down-interval "):
		block.HoldDownInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "interface-monitor "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		block.InterfaceMonitor = tfdata.AppendPotentialNewBlock(block.InterfaceMonitor, types.StringValue(name))
		interfaceMonitor := &block.InterfaceMonitor[len(block.InterfaceMonitor)-1]
		balt.CutPrefixInString(&itemTrim, name+" ")

		if balt.CutPrefixInString(&itemTrim, "weight ") {
			interfaceMonitor.Weight, err = tfdata.ConvAtoi64Value(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "preempt"):
		block.Preempt = types.BoolValue(true)

		if balt.CutPrefixInString(&itemTrim, " ") {
			switch {
			case balt.CutPrefixInString(&itemTrim, "delay "):
				block.PreemptDelay, err = tfdata.ConvAtoi64Value(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "limit "):
				block.PreemptLimit, err = tfdata.ConvAtoi64Value(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "period "):
				block.PreemptPeriod, err = tfdata.ConvAtoi64Value(itemTrim)
			}
		}
	}

	return err
}

func (rscData *chassisClusterData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete chassis cluster",
		"delete interfaces fab0",
		"delete interfaces fab1",
	}

	return junSess.ConfigSet(configSet)
}
