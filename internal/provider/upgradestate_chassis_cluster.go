package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *chassisCluster) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"reth_count": schema.Int64Attribute{
						Required: true,
					},
					"config_sync_no_secondary_bootup_auto": schema.BoolAttribute{
						Optional: true,
					},
					"control_link_recovery": schema.BoolAttribute{
						Optional: true,
					},
					"heartbeat_interval": schema.Int64Attribute{
						Optional: true,
					},
					"heartbeat_threshold": schema.Int64Attribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"fab0": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"member_interfaces": schema.ListAttribute{
									ElementType: types.StringType,
									Required:    true,
								},
								"description": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"fab1": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"member_interfaces": schema.ListAttribute{
									ElementType: types.StringType,
									Required:    true,
								},
								"description": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"redundancy_group": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"node0_priority": schema.Int64Attribute{
									Required: true,
								},
								"node1_priority": schema.Int64Attribute{
									Required: true,
								},
								"gratuitous_arp_count": schema.Int64Attribute{
									Optional: true,
								},
								"hold_down_interval": schema.Int64Attribute{
									Optional: true,
								},
								"preempt": schema.BoolAttribute{
									Optional: true,
								},
								"preempt_delay": schema.Int64Attribute{
									Optional: true,
								},
								"preempt_limit": schema.Int64Attribute{
									Optional: true,
								},
								"preempt_period": schema.Int64Attribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"interface_monitor": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required: true,
											},
											"weight": schema.Int64Attribute{
												Required: true,
											},
										},
									},
								},
							},
						},
					},
					"control_ports": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"fpc": schema.Int64Attribute{
									Required: true,
								},
								"port": schema.Int64Attribute{
									Required: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeChassisClusterV0toV1,
		},
	}
}

func upgradeChassisClusterV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                              types.String `tfsdk:"id"`
		RethCount                       types.Int64  `tfsdk:"reth_count"`
		ConfigSyncNoSecondaryBootupAuto types.Bool   `tfsdk:"config_sync_no_secondary_bootup_auto"`
		ControlLinkRecovery             types.Bool   `tfsdk:"control_link_recovery"`
		HeartbeatInterval               types.Int64  `tfsdk:"heartbeat_interval"`
		HeartbeatThreshold              types.Int64  `tfsdk:"heartbeat_threshold"`
		Fab0                            []struct {
			MemberInterfaces []types.String `tfsdk:"member_interfaces"`
			Description      types.String   `tfsdk:"description"`
		} `tfsdk:"fab0"`
		Fab1 []struct {
			MemberInterfaces []types.String `tfsdk:"member_interfaces"`
			Description      types.String   `tfsdk:"description"`
		} `tfsdk:"fab1"`
		RedundancyGroup []struct {
			Node0Priority      types.Int64 `tfsdk:"node0_priority"`
			Node1Priority      types.Int64 `tfsdk:"node1_priority"`
			GratuitousArpCount types.Int64 `tfsdk:"gratuitous_arp_count"`
			HoldDownInterval   types.Int64 `tfsdk:"hold_down_interval"`
			Preempt            types.Bool  `tfsdk:"preempt"`
			PreemptDelay       types.Int64 `tfsdk:"preempt_delay"`
			PreemptLimit       types.Int64 `tfsdk:"preempt_limit"`
			PreemptPeriod      types.Int64 `tfsdk:"preempt_period"`
			InterfaceMonitor   []struct {
				Name   types.String `tfsdk:"name"`
				Weight types.Int64  `tfsdk:"weight"`
			} `tfsdk:"interface_monitor"`
		} `tfsdk:"redundancy_group"`
		ControlPorts []struct {
			Fpc  types.Int64 `tfsdk:"fpc"`
			Port types.Int64 `tfsdk:"port"`
		} `tfsdk:"control_ports"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 chassisClusterData
	dataV1.ID = dataV0.ID
	dataV1.RethCount = dataV0.RethCount
	dataV1.ConfigSyncNoSecondaryBootupAuto = dataV0.ConfigSyncNoSecondaryBootupAuto
	dataV1.ControlLinkRecovery = dataV0.ControlLinkRecovery
	dataV1.HeartbeatInterval = dataV0.HeartbeatInterval
	dataV1.HeartbeatThreshold = dataV0.HeartbeatThreshold
	if len(dataV0.Fab0) > 0 {
		dataV1.Fab0 = &chassisClusterBlockFab{
			MemberInterfaces: dataV0.Fab0[0].MemberInterfaces,
			Description:      dataV0.Fab0[0].Description,
		}
	}
	if len(dataV0.Fab1) > 0 {
		dataV1.Fab1 = &chassisClusterBlockFab{
			MemberInterfaces: dataV0.Fab1[0].MemberInterfaces,
			Description:      dataV0.Fab1[0].Description,
		}
	}
	for _, blockV0 := range dataV0.RedundancyGroup {
		blockV1 := chassisClusterBlockRedundancyGroup{
			Node0Priority:      blockV0.Node0Priority,
			Node1Priority:      blockV0.Node1Priority,
			GratuitousArpCount: blockV0.GratuitousArpCount,
			HoldDownInterval:   blockV0.HoldDownInterval,
			Preempt:            blockV0.Preempt,
			PreemptDelay:       blockV0.PreemptDelay,
			PreemptLimit:       blockV0.PreemptLimit,
			PreemptPeriod:      blockV0.PreemptPeriod,
		}
		for _, subBlockV0 := range blockV0.InterfaceMonitor {
			blockV1.InterfaceMonitor = append(blockV1.InterfaceMonitor,
				chassisClusterBlockRedundancyGroupBlockInterfaceMonitor{
					Name:   subBlockV0.Name,
					Weight: subBlockV0.Weight,
				},
			)
		}
		dataV1.RedundancyGroup = append(dataV1.RedundancyGroup, blockV1)
	}
	for _, blockV0 := range dataV0.ControlPorts {
		dataV1.ControlPorts = append(dataV1.ControlPorts, chassisClusterBlockControlPorts{
			Fpc:  blockV0.Fpc,
			Port: blockV0.Port,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
