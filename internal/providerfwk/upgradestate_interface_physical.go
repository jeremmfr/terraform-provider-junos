package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *interfacePhysical) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"no_disable_on_destroy": schema.BoolAttribute{
						Optional: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"disable": schema.BoolAttribute{
						Optional: true,
					},
					"mtu": schema.Int64Attribute{
						Optional: true,
					},
					"trunk": schema.BoolAttribute{
						Optional: true,
					},
					"vlan_members": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"vlan_native": schema.Int64Attribute{
						Optional: true,
					},
					"vlan_tagging": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"esi": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"mode": schema.StringAttribute{
									Required: true,
								},
								"auto_derive_lacp": schema.BoolAttribute{
									Optional: true,
								},
								"df_election_type": schema.StringAttribute{
									Optional: true,
								},
								"identifier": schema.StringAttribute{
									Optional: true,
								},
								"source_bmac": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"ether_opts": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ae_8023ad": schema.StringAttribute{
									Optional: true,
								},
								"auto_negotiation": schema.BoolAttribute{
									Optional: true,
								},
								"no_auto_negotiation": schema.BoolAttribute{
									Optional: true,
								},
								"flow_control": schema.BoolAttribute{
									Optional: true,
								},
								"no_flow_control": schema.BoolAttribute{
									Optional: true,
								},
								"loopback": schema.BoolAttribute{
									Optional: true,
								},
								"no_loopback": schema.BoolAttribute{
									Optional: true,
								},
								"redundant_parent": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"gigether_opts": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ae_8023ad": schema.StringAttribute{
									Optional: true,
								},
								"auto_negotiation": schema.BoolAttribute{
									Optional: true,
								},
								"no_auto_negotiation": schema.BoolAttribute{
									Optional: true,
								},
								"flow_control": schema.BoolAttribute{
									Optional: true,
								},
								"no_flow_control": schema.BoolAttribute{
									Optional: true,
								},
								"loopback": schema.BoolAttribute{
									Optional: true,
								},
								"no_loopback": schema.BoolAttribute{
									Optional: true,
								},
								"redundant_parent": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"parent_ether_opts": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"flow_control": schema.BoolAttribute{
									Optional: true,
								},
								"no_flow_control": schema.BoolAttribute{
									Optional: true,
								},
								"loopback": schema.BoolAttribute{
									Optional: true,
								},
								"no_loopback": schema.BoolAttribute{
									Optional: true,
								},
								"link_speed": schema.StringAttribute{
									Optional: true,
								},
								"minimum_bandwidth": schema.StringAttribute{
									Optional: true,
								},
								"minimum_links": schema.Int64Attribute{
									Optional: true,
								},
								"redundancy_group": schema.Int64Attribute{
									Optional: true,
								},
								"source_address_filter": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"source_filtering": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"bfd_liveness_detection": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"local_address": schema.StringAttribute{
												Required: true,
											},
											"authentication_algorithm": schema.StringAttribute{
												Optional: true,
											},
											"authentication_key_chain": schema.StringAttribute{
												Optional: true,
											},
											"authentication_loose_check": schema.BoolAttribute{
												Optional: true,
											},
											"detection_time_threshold": schema.Int64Attribute{
												Optional: true,
											},
											"holddown_interval": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_interval": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_receive_interval": schema.Int64Attribute{
												Optional: true,
											},
											"multiplier": schema.Int64Attribute{
												Optional: true,
											},
											"neighbor": schema.StringAttribute{
												Optional: true,
											},
											"no_adaptation": schema.BoolAttribute{
												Optional: true,
											},
											"transmit_interval_minimum_interval": schema.Int64Attribute{
												Optional: true,
											},
											"transmit_interval_threshold": schema.Int64Attribute{
												Optional: true,
											},
											"version": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"lacp": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"mode": schema.StringAttribute{
												Required: true,
											},
											"admin_key": schema.Int64Attribute{
												Optional: true,
											},
											"periodic": schema.StringAttribute{
												Optional: true,
											},
											"sync_reset": schema.StringAttribute{
												Optional: true,
											},
											"system_id": schema.StringAttribute{
												Optional: true,
											},
											"system_priority": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeInterfacePhysicalV0toV1,
		},
	}
}

func upgradeInterfacePhysicalV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                 types.String   `tfsdk:"id"`
		Name               types.String   `tfsdk:"name"`
		NoDisableOnDestroy types.Bool     `tfsdk:"no_disable_on_destroy"`
		Description        types.String   `tfsdk:"description"`
		Disable            types.Bool     `tfsdk:"disable"`
		Mtu                types.Int64    `tfsdk:"mtu"`
		Trunk              types.Bool     `tfsdk:"trunk"`
		VlanMembers        []types.String `tfsdk:"vlan_members"`
		VlanNative         types.Int64    `tfsdk:"vlan_native"`
		VlanTagging        types.Bool     `tfsdk:"vlan_tagging"`
		ESI                []struct {
			Mode           types.String `tfsdk:"mode"`
			AutoDeriveLACP types.Bool   `tfsdk:"auto_derive_lacp"`
			DFElectionType types.String `tfsdk:"df_election_type"`
			Identifier     types.String `tfsdk:"identifier"`
			SourceBMAC     types.String `tfsdk:"source_bmac"`
		} `tfsdk:"esi"`
		EtherOpts []struct {
			Ae8023ad          types.String `tfsdk:"ae_8023ad"`
			AutoNegotiation   types.Bool   `tfsdk:"auto_negotiation"`
			NoAutoNegotiation types.Bool   `tfsdk:"no_auto_negotiation"`
			FlowControl       types.Bool   `tfsdk:"flow_control"`
			NoFlowControl     types.Bool   `tfsdk:"no_flow_control"`
			Loopback          types.Bool   `tfsdk:"loopback"`
			NoLoopback        types.Bool   `tfsdk:"no_loopback"`
			RedundantParent   types.String `tfsdk:"redundant_parent"`
		} `tfsdk:"ether_opts"`
		GigetherOpts []struct {
			Ae8023ad          types.String `tfsdk:"ae_8023ad"`
			AutoNegotiation   types.Bool   `tfsdk:"auto_negotiation"`
			NoAutoNegotiation types.Bool   `tfsdk:"no_auto_negotiation"`
			FlowControl       types.Bool   `tfsdk:"flow_control"`
			NoFlowControl     types.Bool   `tfsdk:"no_flow_control"`
			Loopback          types.Bool   `tfsdk:"loopback"`
			NoLoopback        types.Bool   `tfsdk:"no_loopback"`
			RedundantParent   types.String `tfsdk:"redundant_parent"`
		} `tfsdk:"gigether_opts"`
		ParentEtherOpts []struct {
			FlowControl          types.Bool     `tfsdk:"flow_control"`
			NoFlowControl        types.Bool     `tfsdk:"no_flow_control"`
			Loopback             types.Bool     `tfsdk:"loopback"`
			NoLoopback           types.Bool     `tfsdk:"no_loopback"`
			LinkSpeed            types.String   `tfsdk:"link_speed"`
			MinimumBandwidth     types.String   `tfsdk:"minimum_bandwidth"`
			MinimumLinks         types.Int64    `tfsdk:"minimum_links"`
			RedundancyGroup      types.Int64    `tfsdk:"redundancy_group"`
			SourceAddressFilter  []types.String `tfsdk:"source_address_filter"`
			SourceFiltering      types.Bool     `tfsdk:"source_filtering"`
			BFDLivenessDetection []struct {
				LocalAddress                    types.String `tfsdk:"local_address"`
				AuthenticationAlgorithm         types.String `tfsdk:"authentication_algorithm"`
				AuthenticationLooseCheck        types.Bool   `tfsdk:"authentication_loose_check"`
				AuthenticationKeyChain          types.String `tfsdk:"authentication_key_chain"`
				DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
				HolddownInterval                types.Int64  `tfsdk:"holddown_interval"`
				MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
				MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
				Multiplier                      types.Int64  `tfsdk:"multiplier"`
				Neighbor                        types.String `tfsdk:"neighbor"`
				NoAdaptation                    types.Bool   `tfsdk:"no_adaptation"`
				TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
				TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
				Version                         types.String `tfsdk:"version"`
			} `tfsdk:"bfd_liveness_detection"`
			LACP []struct {
				Mode           types.String `tfsdk:"mode"`
				AdminKey       types.Int64  `tfsdk:"admin_key"`
				Periodic       types.String `tfsdk:"periodic"`
				SyncReset      types.String `tfsdk:"sync_reset"`
				SystemID       types.String `tfsdk:"system_id"`
				SystemPriority types.Int64  `tfsdk:"system_priority"`
			} `tfsdk:"lacp"`
		} `tfsdk:"parent_ether_opts"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 interfacePhysicalData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.NoDisableOnDestroy = dataV0.NoDisableOnDestroy
	if !dataV1.NoDisableOnDestroy.IsNull() && !dataV1.NoDisableOnDestroy.ValueBool() {
		dataV1.NoDisableOnDestroy = types.BoolNull()
	}
	dataV1.Description = dataV0.Description
	dataV1.Disable = dataV0.Disable
	dataV1.Mtu = dataV0.Mtu
	dataV1.Trunk = dataV0.Trunk
	dataV1.VlanMembers = dataV0.VlanMembers
	dataV1.VlanNative = dataV0.VlanNative
	dataV1.VlanTagging = dataV0.VlanTagging
	if len(dataV0.ESI) > 0 {
		dataV1.ESI = &interfacePhysicalBlockESI{
			AutoDeriveLACP: dataV0.ESI[0].AutoDeriveLACP,
			Mode:           dataV0.ESI[0].Mode,
			DFElectionType: dataV0.ESI[0].DFElectionType,
			Identifier:     dataV0.ESI[0].Identifier,
			SourceBMAC:     dataV0.ESI[0].SourceBMAC,
		}
	}
	if len(dataV0.EtherOpts) > 0 {
		dataV1.EtherOpts = &interfacePhysicalBlockEtherOpts{
			AutoNegotiation:   dataV0.EtherOpts[0].AutoNegotiation,
			NoAutoNegotiation: dataV0.EtherOpts[0].NoAutoNegotiation,
			FlowControl:       dataV0.EtherOpts[0].FlowControl,
			NoFlowControl:     dataV0.EtherOpts[0].NoFlowControl,
			Loopback:          dataV0.EtherOpts[0].Loopback,
			NoLoopback:        dataV0.EtherOpts[0].NoLoopback,
			Ae8023ad:          dataV0.EtherOpts[0].Ae8023ad,
			RedundantParent:   dataV0.EtherOpts[0].RedundantParent,
		}
	}
	if len(dataV0.GigetherOpts) > 0 {
		dataV1.GigetherOpts = &interfacePhysicalBlockEtherOpts{
			AutoNegotiation:   dataV0.GigetherOpts[0].AutoNegotiation,
			NoAutoNegotiation: dataV0.GigetherOpts[0].NoAutoNegotiation,
			FlowControl:       dataV0.GigetherOpts[0].FlowControl,
			NoFlowControl:     dataV0.GigetherOpts[0].NoFlowControl,
			Loopback:          dataV0.GigetherOpts[0].Loopback,
			NoLoopback:        dataV0.GigetherOpts[0].NoLoopback,
			Ae8023ad:          dataV0.GigetherOpts[0].Ae8023ad,
			RedundantParent:   dataV0.GigetherOpts[0].RedundantParent,
		}
	}
	if len(dataV0.ParentEtherOpts) > 0 {
		dataV1.ParentEtherOpts = &interfacePhysicalBlockParentEtherOpts{
			FlowControl:         dataV0.ParentEtherOpts[0].FlowControl,
			NoFlowControl:       dataV0.ParentEtherOpts[0].NoFlowControl,
			Loopback:            dataV0.ParentEtherOpts[0].Loopback,
			NoLoopback:          dataV0.ParentEtherOpts[0].NoLoopback,
			LinkSpeed:           dataV0.ParentEtherOpts[0].LinkSpeed,
			MinimumBandwidth:    dataV0.ParentEtherOpts[0].MinimumBandwidth,
			MinimumLinks:        dataV0.ParentEtherOpts[0].MinimumLinks,
			RedundancyGroup:     dataV0.ParentEtherOpts[0].RedundancyGroup,
			SourceAddressFilter: dataV0.ParentEtherOpts[0].SourceAddressFilter,
			SourceFiltering:     dataV0.ParentEtherOpts[0].SourceFiltering,
		}
		if len(dataV0.ParentEtherOpts[0].BFDLivenessDetection) > 0 {
			dataV1.ParentEtherOpts.BFDLivenessDetection = &interfacePhysicalBlockParentEtherOptsBlockBFDLivenessDetection{
				AuthenticationLooseCheck:        dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].AuthenticationLooseCheck,
				NoAdaptation:                    dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].NoAdaptation,
				LocalAddress:                    dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].LocalAddress,
				AuthenticationAlgorithm:         dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].AuthenticationAlgorithm,
				AuthenticationKeyChain:          dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].AuthenticationKeyChain,
				DetectionTimeThreshold:          dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].DetectionTimeThreshold,
				HolddownInterval:                dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].HolddownInterval,
				MinimumInterval:                 dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].MinimumInterval,
				MinimumReceiveInterval:          dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].MinimumReceiveInterval,
				Multiplier:                      dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].Multiplier,
				Neighbor:                        dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].Neighbor,
				TransmitIntervalMinimumInterval: dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].TransmitIntervalMinimumInterval,
				TransmitIntervalThreshold:       dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].TransmitIntervalThreshold,
				Version:                         dataV0.ParentEtherOpts[0].BFDLivenessDetection[0].Version,
			}
		}
		if len(dataV0.ParentEtherOpts[0].LACP) > 0 {
			dataV1.ParentEtherOpts.LACP = &interfacePhysicalBlockParentEtherOptsBlockLACP{
				Mode:           dataV0.ParentEtherOpts[0].LACP[0].Mode,
				AdminKey:       dataV0.ParentEtherOpts[0].LACP[0].AdminKey,
				Periodic:       dataV0.ParentEtherOpts[0].LACP[0].Periodic,
				SyncReset:      dataV0.ParentEtherOpts[0].LACP[0].SyncReset,
				SystemID:       dataV0.ParentEtherOpts[0].LACP[0].SystemID,
				SystemPriority: dataV0.ParentEtherOpts[0].LACP[0].SystemPriority,
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
