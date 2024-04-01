package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *ospfArea) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"area_id": schema.StringAttribute{
						Required: true,
					},
					"version": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"realm": schema.StringAttribute{
						Optional: true,
					},
					"routing_instance": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"context_identifier": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"inter_area_prefix_export": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"inter_area_prefix_import": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"network_summary_export": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"network_summary_import": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"no_context_identifier_advertisement": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"interface": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"authentication_simple_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"dead_interval": schema.Int64Attribute{
									Optional: true,
								},
								"demand_circuit": schema.BoolAttribute{
									Optional: true,
								},
								"disable": schema.BoolAttribute{
									Optional: true,
								},
								"dynamic_neighbors": schema.BoolAttribute{
									Optional: true,
								},
								"flood_reduction": schema.BoolAttribute{
									Optional: true,
								},
								"hello_interval": schema.Int64Attribute{
									Optional: true,
								},
								"interface_type": schema.StringAttribute{
									Optional: true,
								},
								"ipsec_sa": schema.StringAttribute{
									Optional: true,
								},
								"ipv4_adjacency_segment_protected_type": schema.StringAttribute{
									Optional: true,
								},
								"ipv4_adjacency_segment_protected_value": schema.StringAttribute{
									Optional: true,
								},
								"ipv4_adjacency_segment_unprotected_type": schema.StringAttribute{
									Optional: true,
								},
								"ipv4_adjacency_segment_unprotected_value": schema.StringAttribute{
									Optional: true,
								},
								"link_protection": schema.BoolAttribute{
									Optional: true,
								},
								"metric": schema.Int64Attribute{
									Optional: true,
								},
								"mtu": schema.Int64Attribute{
									Optional: true,
								},
								"no_advertise_adjacency_segment": schema.BoolAttribute{
									Optional: true,
								},
								"no_eligible_backup": schema.BoolAttribute{
									Optional: true,
								},
								"no_eligible_remote_backup": schema.BoolAttribute{
									Optional: true,
								},
								"no_interface_state_traps": schema.BoolAttribute{
									Optional: true,
								},
								"no_neighbor_down_notification": schema.BoolAttribute{
									Optional: true,
								},
								"node_link_protection": schema.BoolAttribute{
									Optional: true,
								},
								"passive": schema.BoolAttribute{
									Optional: true,
								},
								"passive_traffic_engineering_remote_node_id": schema.StringAttribute{
									Optional: true,
								},
								"passive_traffic_engineering_remote_node_router_id": schema.StringAttribute{
									Optional: true,
								},
								"poll_interval": schema.Int64Attribute{
									Optional: true,
								},
								"priority": schema.Int64Attribute{
									Optional: true,
								},
								"retransmit_interval": schema.Int64Attribute{
									Optional: true,
								},
								"secondary": schema.BoolAttribute{
									Optional: true,
								},
								"strict_bfd": schema.BoolAttribute{
									Optional: true,
								},
								"te_metric": schema.Int64Attribute{
									Optional: true,
								},
								"transit_delay": schema.Int64Attribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"authentication_md5": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"key_id": schema.Int64Attribute{
												Required: true,
											},
											"key": schema.StringAttribute{
												Required:  true,
												Sensitive: true,
											},
											"start_time": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"bandwidth_based_metrics": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"bandwidth": schema.StringAttribute{
												Required: true,
											},
											"metric": schema.Int64Attribute{
												Required: true,
											},
										},
									},
								},
								"bfd_liveness_detection": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
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
											"full_neighbors_only": schema.BoolAttribute{
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
								"neighbor": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"address": schema.StringAttribute{
												Required: true,
											},
											"eligible": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"area_range": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"range": schema.StringAttribute{
									Required: true,
								},
								"exact": schema.BoolAttribute{
									Optional: true,
								},
								"override_metric": schema.Int64Attribute{
									Optional: true,
								},
								"restrict": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"nssa": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"summaries": schema.BoolAttribute{
									Optional: true,
								},
								"no_summaries": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"area_range": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"range": schema.StringAttribute{
												Required: true,
											},
											"exact": schema.BoolAttribute{
												Optional: true,
											},
											"override_metric": schema.Int64Attribute{
												Optional: true,
											},
											"restrict": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"default_lsa": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"default_metric": schema.Int64Attribute{
												Optional: true,
											},
											"metric_type": schema.Int64Attribute{
												Optional: true,
											},
											"type_7": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"stub": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"default_metric": schema.Int64Attribute{
									Optional: true,
								},
								"summaries": schema.BoolAttribute{
									Optional: true,
								},
								"no_summaries": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"virtual_link": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"neighbor_id": schema.StringAttribute{
									Required: true,
								},
								"transit_area": schema.StringAttribute{
									Required: true,
								},
								"dead_interval": schema.Int64Attribute{
									Optional: true,
								},
								"demand_circuit": schema.BoolAttribute{
									Optional: true,
								},
								"disable": schema.BoolAttribute{
									Optional: true,
								},
								"flood_reduction": schema.BoolAttribute{
									Optional: true,
								},
								"hello_interval": schema.Int64Attribute{
									Optional: true,
								},
								"ipsec_sa": schema.StringAttribute{
									Optional: true,
								},
								"mtu": schema.Int64Attribute{
									Optional: true,
								},
								"retransmit_interval": schema.Int64Attribute{
									Optional: true,
								},
								"transit_delay": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeOspfAreaStateV0toV1,
		},
	}
}

func upgradeOspfAreaStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                               types.String   `tfsdk:"id"`
		AreaID                           types.String   `tfsdk:"area_id"`
		Version                          types.String   `tfsdk:"version"`
		Realm                            types.String   `tfsdk:"realm"`
		RoutingInstance                  types.String   `tfsdk:"routing_instance"`
		ContextIdentifier                []types.String `tfsdk:"context_identifier"`
		InterAreaPrefixExport            []types.String `tfsdk:"inter_area_prefix_export"`
		InterAreaPrefixImport            []types.String `tfsdk:"inter_area_prefix_import"`
		NetworkSummaryExport             []types.String `tfsdk:"network_summary_export"`
		NetworkSummaryImport             []types.String `tfsdk:"network_summary_import"`
		NoContextIdentifierAdvertisement types.Bool     `tfsdk:"no_context_identifier_advertisement"`
		Interface                        []struct {
			Name                                        types.String `tfsdk:"name"`
			AuthenticationSimplePassword                types.String `tfsdk:"authentication_simple_password"`
			DeadInterval                                types.Int64  `tfsdk:"dead_interval"`
			DemandCircuit                               types.Bool   `tfsdk:"demand_circuit"`
			Disable                                     types.Bool   `tfsdk:"disable"`
			DynamicNeighbors                            types.Bool   `tfsdk:"dynamic_neighbors"`
			FloodReduction                              types.Bool   `tfsdk:"flood_reduction"`
			HelloInterval                               types.Int64  `tfsdk:"hello_interval"`
			InterfaceType                               types.String `tfsdk:"interface_type"`
			IpsecSA                                     types.String `tfsdk:"ipsec_sa"`
			IPv4AdjacencySegmentProtectedType           types.String `tfsdk:"ipv4_adjacency_segment_protected_type"`
			IPv4AdjacencySegmentProtectedValue          types.String `tfsdk:"ipv4_adjacency_segment_protected_value"`
			IPv4AdjacencySegmentUnprotectedType         types.String `tfsdk:"ipv4_adjacency_segment_unprotected_type"`
			IPv4AdjacencySegmentUnprotectedValue        types.String `tfsdk:"ipv4_adjacency_segment_unprotected_value"`
			LinkProtection                              types.Bool   `tfsdk:"link_protection"`
			Metric                                      types.Int64  `tfsdk:"metric"`
			Mtu                                         types.Int64  `tfsdk:"mtu"`
			NoAdvertiseAdjacencySegment                 types.Bool   `tfsdk:"no_advertise_adjacency_segment"`
			NoEligibleBackup                            types.Bool   `tfsdk:"no_eligible_backup"`
			NoEligibleRemoteBackup                      types.Bool   `tfsdk:"no_eligible_remote_backup"`
			NoInterfaceStateTraps                       types.Bool   `tfsdk:"no_interface_state_traps"`
			NoNeighborDownNotification                  types.Bool   `tfsdk:"no_neighbor_down_notification"`
			NodeLinkProtection                          types.Bool   `tfsdk:"node_link_protection"`
			Passive                                     types.Bool   `tfsdk:"passive"`
			PassiveTrafficEngineeringRemoteNodeID       types.String `tfsdk:"passive_traffic_engineering_remote_node_id"`
			PassiveTrafficEngineeringRemoteNodeRouterID types.String `tfsdk:"passive_traffic_engineering_remote_node_router_id"`
			PollInterval                                types.Int64  `tfsdk:"poll_interval"`
			Priority                                    types.Int64  `tfsdk:"priority"`
			RetransmitInterval                          types.Int64  `tfsdk:"retransmit_interval"`
			Secondary                                   types.Bool   `tfsdk:"secondary"`
			StrictBfd                                   types.Bool   `tfsdk:"strict_bfd"`
			TeMetric                                    types.Int64  `tfsdk:"te_metric"`
			TransitDelay                                types.Int64  `tfsdk:"transit_delay"`
			AuthenticationMD5                           []struct {
				KeyID     types.Int64  `tfsdk:"key_id"`
				Key       types.String `tfsdk:"key"`
				StartTime types.String `tfsdk:"start_time"`
			} `tfsdk:"authentication_md5"`
			BandwidthBasedMetrics []struct {
				Bandwidth types.String `tfsdk:"bandwidth"`
				Metric    types.Int64  `tfsdk:"metric"`
			} `tfsdk:"bandwidth_based_metrics"`
			BfdLivenessDetection []struct {
				AuthenticationAlgorithm         types.String `tfsdk:"authentication_algorithm"`
				AuthenticationKeyChain          types.String `tfsdk:"authentication_key_chain"`
				AuthenticationLooseCheck        types.Bool   `tfsdk:"authentication_loose_check"`
				DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
				FullNeighborsOnly               types.Bool   `tfsdk:"full_neighbors_only"`
				HolddownInterval                types.Int64  `tfsdk:"holddown_interval"`
				MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
				MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
				Multiplier                      types.Int64  `tfsdk:"multiplier"`
				NoAdaptation                    types.Bool   `tfsdk:"no_adaptation"`
				TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
				TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
				Version                         types.String `tfsdk:"version"`
			} `tfsdk:"bfd_liveness_detection"`
			Neighbor []struct {
				Address  types.String `tfsdk:"address"`
				Eligbile types.Bool   `tfsdk:"eligible"`
			} `tfsdk:"neighbor"`
		} `tfsdk:"interface"`
		AreaRange []struct {
			Range          types.String `tfsdk:"range"`
			Exact          types.Bool   `tfsdk:"exact"`
			OverrideMetric types.Int64  `tfsdk:"override_metric"`
			Restrict       types.Bool   `tfsdk:"restrict"`
		} `tfsdk:"area_range"`
		Nssa []struct {
			Summaries   types.Bool `tfsdk:"summaries"`
			NoSummaries types.Bool `tfsdk:"no_summaries"`
			AreaRange   []struct {
				Range          types.String `tfsdk:"range"`
				Exact          types.Bool   `tfsdk:"exact"`
				OverrideMetric types.Int64  `tfsdk:"override_metric"`
				Restrict       types.Bool   `tfsdk:"restrict"`
			} `tfsdk:"area_range"`
			DefaultLsa []struct {
				DefaultMetric types.Int64 `tfsdk:"default_metric"`
				MetricType    types.Int64 `tfsdk:"metric_type"`
				Type7         types.Bool  `tfsdk:"type_7"`
			} `tfsdk:"default_lsa"`
		} `tfsdk:"nssa"`
		Stub []struct {
			DefaultMetric types.Int64 `tfsdk:"default_metric"`
			Summaries     types.Bool  `tfsdk:"summaries"`
			NoSummaries   types.Bool  `tfsdk:"no_summaries"`
		} `tfsdk:"stub"`
		VirtualLink []struct {
			NeighborID         types.String `tfsdk:"neighbor_id"`
			TransitArea        types.String `tfsdk:"transit_area"`
			DeadInterval       types.Int64  `tfsdk:"dead_interval"`
			DemandCircuit      types.Bool   `tfsdk:"demand_circuit"`
			Disable            types.Bool   `tfsdk:"disable"`
			FloodReduction     types.Bool   `tfsdk:"flood_reduction"`
			HelloInterval      types.Int64  `tfsdk:"hello_interval"`
			IpsecSA            types.String `tfsdk:"ipsec_sa"`
			Mtu                types.Int64  `tfsdk:"mtu"`
			RetransmitInterval types.Int64  `tfsdk:"retransmit_interval"`
			TransitDelay       types.Int64  `tfsdk:"transit_delay"`
		} `tfsdk:"virtual_link"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 ospfAreaData
	dataV1.ID = dataV0.ID
	dataV1.AreaID = dataV0.AreaID
	dataV1.Version = dataV0.Version
	dataV1.Realm = dataV0.Realm
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.ContextIdentifier = dataV0.ContextIdentifier
	dataV1.InterAreaPrefixExport = dataV0.InterAreaPrefixExport
	dataV1.InterAreaPrefixImport = dataV0.InterAreaPrefixImport
	dataV1.NetworkSummaryExport = dataV0.NetworkSummaryExport
	dataV1.NetworkSummaryImport = dataV0.NetworkSummaryImport
	dataV1.NoContextIdentifierAdvertisement = dataV0.NoContextIdentifierAdvertisement
	for _, blockV0 := range dataV0.Interface {
		blockV1 := ospfAreaBlockInterface{
			Name:                                  blockV0.Name,
			AuthenticationSimplePassword:          blockV0.AuthenticationSimplePassword,
			DeadInterval:                          blockV0.DeadInterval,
			DemandCircuit:                         blockV0.DemandCircuit,
			Disable:                               blockV0.Disable,
			DynamicNeighbors:                      blockV0.DynamicNeighbors,
			FloodReduction:                        blockV0.FloodReduction,
			HelloInterval:                         blockV0.HelloInterval,
			InterfaceType:                         blockV0.InterfaceType,
			IpsecSA:                               blockV0.IpsecSA,
			IPv4AdjacencySegmentProtectedType:     blockV0.IPv4AdjacencySegmentProtectedType,
			IPv4AdjacencySegmentProtectedValue:    blockV0.IPv4AdjacencySegmentProtectedValue,
			IPv4AdjacencySegmentUnprotectedType:   blockV0.IPv4AdjacencySegmentUnprotectedType,
			IPv4AdjacencySegmentUnprotectedValue:  blockV0.IPv4AdjacencySegmentUnprotectedValue,
			LinkProtection:                        blockV0.LinkProtection,
			Metric:                                blockV0.Metric,
			Mtu:                                   blockV0.Mtu,
			NoAdvertiseAdjacencySegment:           blockV0.NoAdvertiseAdjacencySegment,
			NoEligibleBackup:                      blockV0.NoEligibleBackup,
			NoEligibleRemoteBackup:                blockV0.NoEligibleBackup,
			NoInterfaceStateTraps:                 blockV0.NoInterfaceStateTraps,
			NoNeighborDownNotification:            blockV0.NoNeighborDownNotification,
			NodeLinkProtection:                    blockV0.NodeLinkProtection,
			Passive:                               blockV0.Passive,
			PassiveTrafficEngineeringRemoteNodeID: blockV0.PassiveTrafficEngineeringRemoteNodeID,
			PassiveTrafficEngineeringRemoteNodeRouterID: blockV0.PassiveTrafficEngineeringRemoteNodeRouterID,
			PollInterval:       blockV0.PollInterval,
			Priority:           blockV0.Priority,
			RetransmitInterval: blockV0.RetransmitInterval,
			Secondary:          blockV0.Secondary,
			StrictBfd:          blockV0.StrictBfd,
			TeMetric:           blockV0.TeMetric,
			TransitDelay:       blockV0.TransitDelay,
		}
		for _, subBlockV0 := range blockV0.AuthenticationMD5 {
			subBlockV1 := ospfAreaBlockInterfaceBlockAuthenticationMD5{
				KeyID:     subBlockV0.KeyID,
				Key:       subBlockV0.Key,
				StartTime: subBlockV0.StartTime,
			}
			blockV1.AuthenticationMD5 = append(blockV1.AuthenticationMD5, subBlockV1)
		}
		for _, subBlockV0 := range blockV0.BandwidthBasedMetrics {
			subBlockV1 := ospfAreaBlockInterfaceBlockBandwidthBasedMetrics{
				Bandwidth: subBlockV0.Bandwidth,
				Metric:    subBlockV0.Metric,
			}
			blockV1.BandwidthBasedMetrics = append(blockV1.BandwidthBasedMetrics, subBlockV1)
		}
		if len(blockV0.BfdLivenessDetection) > 0 {
			blockV1.BfdLivenessDetection = &ospfAreaBlockInterfaceBlockBfdLivenessDetection{
				AuthenticationAlgorithm:         blockV0.BfdLivenessDetection[0].AuthenticationAlgorithm,
				AuthenticationKeyChain:          blockV0.BfdLivenessDetection[0].AuthenticationKeyChain,
				AuthenticationLooseCheck:        blockV0.BfdLivenessDetection[0].AuthenticationLooseCheck,
				DetectionTimeThreshold:          blockV0.BfdLivenessDetection[0].DetectionTimeThreshold,
				FullNeighborsOnly:               blockV0.BfdLivenessDetection[0].FullNeighborsOnly,
				HolddownInterval:                blockV0.BfdLivenessDetection[0].HolddownInterval,
				MinimumInterval:                 blockV0.BfdLivenessDetection[0].MinimumInterval,
				MinimumReceiveInterval:          blockV0.BfdLivenessDetection[0].MinimumReceiveInterval,
				Multiplier:                      blockV0.BfdLivenessDetection[0].Multiplier,
				NoAdaptation:                    blockV0.BfdLivenessDetection[0].NoAdaptation,
				TransmitIntervalMinimumInterval: blockV0.BfdLivenessDetection[0].TransmitIntervalMinimumInterval,
				TransmitIntervalThreshold:       blockV0.BfdLivenessDetection[0].TransmitIntervalThreshold,
				Version:                         blockV0.BfdLivenessDetection[0].Version,
			}
		}
		for _, subBlockV0 := range blockV0.Neighbor {
			subBlockV1 := ospfAreaBlockInterfaceBlockNeighbor{
				Address:  subBlockV0.Address,
				Eligbile: subBlockV0.Eligbile,
			}
			blockV1.Neighbor = append(blockV1.Neighbor, subBlockV1)
		}
		dataV1.Interface = append(dataV1.Interface, blockV1)
	}
	for _, blockV0 := range dataV0.AreaRange {
		blockV1 := ospfAreaBlockAreaRange{
			Range:          blockV0.Range,
			Exact:          blockV0.Exact,
			OverrideMetric: blockV0.OverrideMetric,
			Restrict:       blockV0.Restrict,
		}
		dataV1.AreaRange = append(dataV1.AreaRange, blockV1)
	}
	if len(dataV0.Nssa) > 0 {
		dataV1.Nssa = &ospfAreaBlockNssa{
			Summaries:   dataV0.Nssa[0].Summaries,
			NoSummaries: dataV0.Nssa[0].NoSummaries,
		}
		for _, subBlockV0 := range dataV0.Nssa[0].AreaRange {
			subBlockV1 := ospfAreaBlockAreaRange{
				Range:          subBlockV0.Range,
				Exact:          subBlockV0.Exact,
				OverrideMetric: subBlockV0.OverrideMetric,
				Restrict:       subBlockV0.Restrict,
			}
			dataV1.Nssa.AreaRange = append(dataV1.Nssa.AreaRange, subBlockV1)
		}
		if len(dataV0.Nssa[0].DefaultLsa) > 0 {
			dataV1.Nssa.DefaultLsa = &ospfAreaBlockNssaBlockDefaultLsa{
				DefaultMetric: dataV0.Nssa[0].DefaultLsa[0].DefaultMetric,
				MetricType:    dataV0.Nssa[0].DefaultLsa[0].MetricType,
				Type7:         dataV0.Nssa[0].DefaultLsa[0].Type7,
			}
		}
	}
	if len(dataV0.Stub) > 0 {
		dataV1.Stub = &ospfAreaBlockStub{
			DefaultMetric: dataV0.Stub[0].DefaultMetric,
			Summaries:     dataV0.Stub[0].Summaries,
			NoSummaries:   dataV0.Stub[0].NoSummaries,
		}
	}
	for _, blockV0 := range dataV0.VirtualLink {
		blockV1 := ospfAreaBlockVirtualLink{
			NeighborID:         blockV0.NeighborID,
			TransitArea:        blockV0.TransitArea,
			DeadInterval:       blockV0.DeadInterval,
			DemandCircuit:      blockV0.DemandCircuit,
			Disable:            blockV0.Disable,
			FloodReduction:     blockV0.FloodReduction,
			HelloInterval:      blockV0.HelloInterval,
			IpsecSA:            blockV0.IpsecSA,
			Mtu:                blockV0.Mtu,
			RetransmitInterval: blockV0.RetransmitInterval,
			TransitDelay:       blockV0.TransitDelay,
		}
		dataV1.VirtualLink = append(dataV1.VirtualLink, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
