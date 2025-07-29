package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityScreen) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"alarm_without_drop": schema.BoolAttribute{
						Optional: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"icmp": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"fragment": schema.BoolAttribute{
									Optional: true,
								},
								"icmpv6_malformed": schema.BoolAttribute{
									Optional: true,
								},
								"large": schema.BoolAttribute{
									Optional: true,
								},
								"ping_death": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"flood": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threshold": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"sweep": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threshold": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"ip": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"bad_option": schema.BoolAttribute{
									Optional: true,
								},
								"block_frag": schema.BoolAttribute{
									Optional: true,
								},
								"ipv6_extension_header_limit": schema.Int64Attribute{
									Optional: true,
								},
								"ipv6_malformed_header": schema.BoolAttribute{
									Optional: true,
								},
								"loose_source_route_option": schema.BoolAttribute{
									Optional: true,
								},
								"record_route_option": schema.BoolAttribute{
									Optional: true,
								},
								"security_option": schema.BoolAttribute{
									Optional: true,
								},
								"source_route_option": schema.BoolAttribute{
									Optional: true,
								},
								"spoofing": schema.BoolAttribute{
									Optional: true,
								},
								"stream_option": schema.BoolAttribute{
									Optional: true,
								},
								"strict_source_route_option": schema.BoolAttribute{
									Optional: true,
								},
								"tear_drop": schema.BoolAttribute{
									Optional: true,
								},
								"timestamp_option": schema.BoolAttribute{
									Optional: true,
								},
								"unknown_protocol": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"ipv6_extension_header": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"ah_header": schema.BoolAttribute{
												Optional: true,
											},
											"esp_header": schema.BoolAttribute{
												Optional: true,
											},
											"hip_header": schema.BoolAttribute{
												Optional: true,
											},
											"fragment_header": schema.BoolAttribute{
												Optional: true,
											},
											"mobility_header": schema.BoolAttribute{
												Optional: true,
											},
											"no_next_header": schema.BoolAttribute{
												Optional: true,
											},
											"routing_header": schema.BoolAttribute{
												Optional: true,
											},
											"shim6_header": schema.BoolAttribute{
												Optional: true,
											},
											"user_defined_header_type": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
										},
										Blocks: map[string]schema.Block{
											"destination_header": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"home_address_option": schema.BoolAttribute{
															Optional: true,
														},
														"ilnp_nonce_option": schema.BoolAttribute{
															Optional: true,
														},
														"line_identification_option": schema.BoolAttribute{
															Optional: true,
														},
														"tunnel_encapsulation_limit_option": schema.BoolAttribute{
															Optional: true,
														},
														"user_defined_option_type": schema.ListAttribute{
															ElementType: types.StringType,
															Optional:    true,
														},
													},
												},
											},
											"hop_by_hop_header": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"calipso_option": schema.BoolAttribute{
															Optional: true,
														},
														"jumbo_payload_option": schema.BoolAttribute{
															Optional: true,
														},
														"quick_start_option": schema.BoolAttribute{
															Optional: true,
														},
														"router_alert_option": schema.BoolAttribute{
															Optional: true,
														},
														"rpl_option": schema.BoolAttribute{
															Optional: true,
														},
														"smf_dpd_option": schema.BoolAttribute{
															Optional: true,
														},
														"user_defined_option_type": schema.ListAttribute{
															ElementType: types.StringType,
															Optional:    true,
														},
													},
												},
											},
										},
									},
								},
								"tunnel": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"bad_inner_header": schema.BoolAttribute{
												Optional: true,
											},
											"ip_in_udp_teredo": schema.BoolAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"gre": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"gre_4in4": schema.BoolAttribute{
															Optional: true,
														},
														"gre_4in6": schema.BoolAttribute{
															Optional: true,
														},
														"gre_6in4": schema.BoolAttribute{
															Optional: true,
														},
														"gre_6in6": schema.BoolAttribute{
															Optional: true,
														},
													},
												},
											},
											"ipip": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"dslite": schema.BoolAttribute{
															Optional: true,
														},
														"ipip_4in4": schema.BoolAttribute{
															Optional: true,
														},
														"ipip_4in6": schema.BoolAttribute{
															Optional: true,
														},
														"ipip_6in4": schema.BoolAttribute{
															Optional: true,
														},
														"ipip_6in6": schema.BoolAttribute{
															Optional: true,
														},
														"ipip_6over4": schema.BoolAttribute{
															Optional: true,
														},
														"ipip_6to4relay": schema.BoolAttribute{
															Optional: true,
														},
														"isatap": schema.BoolAttribute{
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
					},
					"limit_session": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"destination_ip_based": schema.Int64Attribute{
									Optional: true,
								},
								"source_ip_based": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"tcp": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"fin_no_ack": schema.BoolAttribute{
									Optional: true,
								},
								"land": schema.BoolAttribute{
									Optional: true,
								},
								"no_flag": schema.BoolAttribute{
									Optional: true,
								},
								"syn_fin": schema.BoolAttribute{
									Optional: true,
								},
								"syn_frag": schema.BoolAttribute{
									Optional: true,
								},
								"winnuke": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"port_scan": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threshold": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"sweep": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threshold": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"syn_ack_ack_proxy": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threshold": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"syn_flood": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"alarm_threshold": schema.Int64Attribute{
												Optional: true,
											},
											"attack_threshold": schema.Int64Attribute{
												Optional: true,
											},
											"destination_threshold": schema.Int64Attribute{
												Optional: true,
											},
											"source_threshold": schema.Int64Attribute{
												Optional: true,
											},
											"timeout": schema.Int64Attribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"whitelist": schema.SetNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Required: true,
														},
														"destination_address": schema.SetAttribute{
															ElementType: types.StringType,
															Optional:    true,
														},
														"source_address": schema.SetAttribute{
															ElementType: types.StringType,
															Optional:    true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"udp": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"flood": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threshold": schema.Int64Attribute{
												Optional: true,
											},
											"whitelist": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
										},
									},
								},
								"port_scan": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threshold": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"sweep": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threshold": schema.Int64Attribute{
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
			StateUpgrader: upgradeSecurityScreenV0toV1,
		},
	}
}

//nolint:lll
func upgradeSecurityScreenV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID               types.String `tfsdk:"id"`
		Name             types.String `tfsdk:"name"`
		AlarmWithoutDrop types.Bool   `tfsdk:"alarm_without_drop"`
		Description      types.String `tfsdk:"description"`
		Icmp             []struct {
			Fragment        types.Bool `tfsdk:"fragment"`
			Icmpv6Malformed types.Bool `tfsdk:"icmpv6_malformed"`
			Large           types.Bool `tfsdk:"large"`
			PingDeath       types.Bool `tfsdk:"ping_death"`
			Flood           []struct {
				Threshold types.Int64 `tfsdk:"threshold"`
			} `tfsdk:"flood"`
			Sweep []struct {
				Threshold types.Int64 `tfsdk:"threshold"`
			} `tfsdk:"sweep"`
		} `tfsdk:"icmp"`
		IP []struct {
			BadOption                types.Bool  `tfsdk:"bad_option"`
			BlockFrag                types.Bool  `tfsdk:"block_frag"`
			IPv6ExtensionHeaderLimit types.Int64 `tfsdk:"ipv6_extension_header_limit"`
			IPv6MalformedHeader      types.Bool  `tfsdk:"ipv6_malformed_header"`
			LooseSourceRouteOption   types.Bool  `tfsdk:"loose_source_route_option"`
			RecordRouteOption        types.Bool  `tfsdk:"record_route_option"`
			SecurityOption           types.Bool  `tfsdk:"security_option"`
			SourceRouteOption        types.Bool  `tfsdk:"source_route_option"`
			Spoofing                 types.Bool  `tfsdk:"spoofing"`
			StreamOption             types.Bool  `tfsdk:"stream_option"`
			StrictSourceRouteOption  types.Bool  `tfsdk:"strict_source_route_option"`
			TearDrop                 types.Bool  `tfsdk:"tear_drop"`
			TimestampOption          types.Bool  `tfsdk:"timestamp_option"`
			UnknownProtocol          types.Bool  `tfsdk:"unknown_protocol"`
			IPv6ExtensionHeader      []struct {
				AhHeader              types.Bool     `tfsdk:"ah_header"`
				EspHeader             types.Bool     `tfsdk:"esp_header"`
				HipHeader             types.Bool     `tfsdk:"hip_header"`
				FragmentHeader        types.Bool     `tfsdk:"fragment_header"`
				MobilityHeader        types.Bool     `tfsdk:"mobility_header"`
				NoNextHeader          types.Bool     `tfsdk:"no_next_header"`
				RoutingHeader         types.Bool     `tfsdk:"routing_header"`
				Shim6Header           types.Bool     `tfsdk:"shim6_header"`
				UserDefinedHeaderType []types.String `tfsdk:"user_defined_header_type"`
				DestinationHeader     []struct {
					HomeAddressOption              types.Bool     `tfsdk:"home_address_option"`
					IlnpNonceOption                types.Bool     `tfsdk:"ilnp_nonce_option"`
					LineIdentificationOption       types.Bool     `tfsdk:"line_identification_option"`
					TunnelEncapsulationLimitOption types.Bool     `tfsdk:"tunnel_encapsulation_limit_option"`
					UserDefinedOptionType          []types.String `tfsdk:"user_defined_option_type"`
				} `tfsdk:"destination_header"`
				HopByHopHeader []struct {
					CalipsoOption         types.Bool     `tfsdk:"calipso_option"`
					JumboPayloadOption    types.Bool     `tfsdk:"jumbo_payload_option"`
					QuickStartOption      types.Bool     `tfsdk:"quick_start_option"`
					RouterAlertOption     types.Bool     `tfsdk:"router_alert_option"`
					RplOption             types.Bool     `tfsdk:"rpl_option"`
					SmfDpdOption          types.Bool     `tfsdk:"smf_dpd_option"`
					UserDefinedOptionType []types.String `tfsdk:"user_defined_option_type"`
				} `tfsdk:"hop_by_hop_header"`
			} `tfsdk:"ipv6_extension_header"`
			Tunnel []struct {
				BadInnerHeader types.Bool `tfsdk:"bad_inner_header"`
				IPInUDPTeredo  types.Bool `tfsdk:"ip_in_udp_teredo"`
				Gre            []struct {
					Gre4in4 types.Bool `tfsdk:"gre_4in4"`
					Gre4in6 types.Bool `tfsdk:"gre_4in6"`
					Gre6in4 types.Bool `tfsdk:"gre_6in4"`
					Gre6in6 types.Bool `tfsdk:"gre_6in6"`
				} `tfsdk:"gre"`
				Ipip []struct {
					Dslite        types.Bool `tfsdk:"dslite"`
					Ipip4in4      types.Bool `tfsdk:"ipip_4in4"`
					Ipip4in6      types.Bool `tfsdk:"ipip_4in6"`
					Ipip6in4      types.Bool `tfsdk:"ipip_6in4"`
					Ipip6in6      types.Bool `tfsdk:"ipip_6in6"`
					Ipip6over4    types.Bool `tfsdk:"ipip_6over4"`
					Ipip6to4relay types.Bool `tfsdk:"ipip_6to4relay"`
					Isatap        types.Bool `tfsdk:"isatap"`
				} `tfsdk:"ipip"`
			} `tfsdk:"tunnel"`
		} `tfsdk:"ip"`
		LimitSession []struct {
			DestinationIPBased types.Int64 `tfsdk:"destination_ip_based"`
			SourceIPBased      types.Int64 `tfsdk:"source_ip_based"`
		} `tfsdk:"limit_session"`
		TCP []struct {
			FinNoAck types.Bool `tfsdk:"fin_no_ack"`
			Land     types.Bool `tfsdk:"land"`
			NoFlag   types.Bool `tfsdk:"no_flag"`
			SynFin   types.Bool `tfsdk:"syn_fin"`
			SynFrag  types.Bool `tfsdk:"syn_frag"`
			Winnuke  types.Bool `tfsdk:"winnuke"`
			PortScan []struct {
				Threshold types.Int64 `tfsdk:"threshold"`
			} `tfsdk:"port_scan"`
			Sweep []struct {
				Threshold types.Int64 `tfsdk:"threshold"`
			} `tfsdk:"sweep"`
			SynAckAckProxy []struct {
				Threshold types.Int64 `tfsdk:"threshold"`
			} `tfsdk:"syn_ack_ack_proxy"`
			SynFlood []struct {
				AlarmThreshold       types.Int64 `tfsdk:"alarm_threshold"`
				AttackThreshold      types.Int64 `tfsdk:"attack_threshold"`
				DestinationThreshold types.Int64 `tfsdk:"destination_threshold"`
				SourceThreshold      types.Int64 `tfsdk:"source_threshold"`
				Timeout              types.Int64 `tfsdk:"timeout"`
				Whitelist            []struct {
					Name               types.String   `tfsdk:"name"`
					DestinationAddress []types.String `tfsdk:"destination_address"`
					SourceAddress      []types.String `tfsdk:"source_address"`
				} `tfsdk:"whitelist"`
			} `tfsdk:"syn_flood"`
		} `tfsdk:"tcp"`
		UDP []struct {
			Flood []struct {
				Threshold types.Int64    `tfsdk:"threshold"`
				Whitelist []types.String `tfsdk:"whitelist"`
			} `tfsdk:"flood"`
			PortScan []struct {
				Threshold types.Int64 `tfsdk:"threshold"`
			} `tfsdk:"port_scan"`
			Sweep []struct {
				Threshold types.Int64 `tfsdk:"threshold"`
			} `tfsdk:"sweep"`
		} `tfsdk:"udp"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityScreenData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.AlarmWithoutDrop = dataV0.AlarmWithoutDrop
	dataV1.Description = dataV0.Description
	if len(dataV0.Icmp) > 0 {
		dataV1.Icmp = &securityScreenBlockIcmp{
			Fragment:        dataV0.Icmp[0].Fragment,
			Icmpv6Malformed: dataV0.Icmp[0].Icmpv6Malformed,
			Large:           dataV0.Icmp[0].Large,
			PingDeath:       dataV0.Icmp[0].PingDeath,
		}
		if len(dataV0.Icmp[0].Flood) > 0 {
			dataV1.Icmp.Flood = &securityScreenBlockWithThreshold{
				Threshold: dataV0.Icmp[0].Flood[0].Threshold,
			}
		}
		if len(dataV0.Icmp[0].Sweep) > 0 {
			dataV1.Icmp.Sweep = &securityScreenBlockWithThreshold{
				Threshold: dataV0.Icmp[0].Sweep[0].Threshold,
			}
		}
	}
	if len(dataV0.IP) > 0 {
		dataV1.IP = &securityScreenBlockIP{
			BadOption:                dataV0.IP[0].BadOption,
			BlockFrag:                dataV0.IP[0].BlockFrag,
			IPv6ExtensionHeaderLimit: dataV0.IP[0].IPv6ExtensionHeaderLimit,
			IPv6MalformedHeader:      dataV0.IP[0].IPv6MalformedHeader,
			LooseSourceRouteOption:   dataV0.IP[0].LooseSourceRouteOption,
			RecordRouteOption:        dataV0.IP[0].RecordRouteOption,
			SecurityOption:           dataV0.IP[0].SecurityOption,
			SourceRouteOption:        dataV0.IP[0].SourceRouteOption,
			Spoofing:                 dataV0.IP[0].Spoofing,
			StreamOption:             dataV0.IP[0].StreamOption,
			StrictSourceRouteOption:  dataV0.IP[0].StrictSourceRouteOption,
			TearDrop:                 dataV0.IP[0].TearDrop,
			TimestampOption:          dataV0.IP[0].TimestampOption,
			UnknownProtocol:          dataV0.IP[0].UnknownProtocol,
		}
		if len(dataV0.IP[0].IPv6ExtensionHeader) > 0 {
			dataV1.IP.IPv6ExtensionHeader = &securityScreenBlockIPBlockIPv6ExtensionHeader{
				AhHeader:              dataV0.IP[0].IPv6ExtensionHeader[0].AhHeader,
				EspHeader:             dataV0.IP[0].IPv6ExtensionHeader[0].EspHeader,
				HipHeader:             dataV0.IP[0].IPv6ExtensionHeader[0].HipHeader,
				FragmentHeader:        dataV0.IP[0].IPv6ExtensionHeader[0].FragmentHeader,
				MobilityHeader:        dataV0.IP[0].IPv6ExtensionHeader[0].MobilityHeader,
				NoNextHeader:          dataV0.IP[0].IPv6ExtensionHeader[0].NoNextHeader,
				RoutingHeader:         dataV0.IP[0].IPv6ExtensionHeader[0].RoutingHeader,
				Shim6Header:           dataV0.IP[0].IPv6ExtensionHeader[0].Shim6Header,
				UserDefinedHeaderType: dataV0.IP[0].IPv6ExtensionHeader[0].UserDefinedHeaderType,
			}
			if len(dataV0.IP[0].IPv6ExtensionHeader[0].DestinationHeader) > 0 {
				dataV1.IP.IPv6ExtensionHeader.DestinationHeader = &securityScreenBlockIPBlockIPv6ExtensionHeaderBlockDestinationHeader{
					HomeAddressOption:              dataV0.IP[0].IPv6ExtensionHeader[0].DestinationHeader[0].HomeAddressOption,
					IlnpNonceOption:                dataV0.IP[0].IPv6ExtensionHeader[0].DestinationHeader[0].IlnpNonceOption,
					LineIdentificationOption:       dataV0.IP[0].IPv6ExtensionHeader[0].DestinationHeader[0].LineIdentificationOption,
					TunnelEncapsulationLimitOption: dataV0.IP[0].IPv6ExtensionHeader[0].DestinationHeader[0].TunnelEncapsulationLimitOption,
					UserDefinedOptionType:          dataV0.IP[0].IPv6ExtensionHeader[0].DestinationHeader[0].UserDefinedOptionType,
				}
			}
			if len(dataV0.IP[0].IPv6ExtensionHeader[0].HopByHopHeader) > 0 {
				dataV1.IP.IPv6ExtensionHeader.HopByHopHeader = &securityScreenBlockIPBlockIPv6ExtensionHeaderBlockHopByHopHeader{
					CalipsoOption:         dataV0.IP[0].IPv6ExtensionHeader[0].HopByHopHeader[0].CalipsoOption,
					JumboPayloadOption:    dataV0.IP[0].IPv6ExtensionHeader[0].HopByHopHeader[0].JumboPayloadOption,
					QuickStartOption:      dataV0.IP[0].IPv6ExtensionHeader[0].HopByHopHeader[0].QuickStartOption,
					RouterAlertOption:     dataV0.IP[0].IPv6ExtensionHeader[0].HopByHopHeader[0].RouterAlertOption,
					RplOption:             dataV0.IP[0].IPv6ExtensionHeader[0].HopByHopHeader[0].RplOption,
					SmfDpdOption:          dataV0.IP[0].IPv6ExtensionHeader[0].HopByHopHeader[0].SmfDpdOption,
					UserDefinedOptionType: dataV0.IP[0].IPv6ExtensionHeader[0].HopByHopHeader[0].UserDefinedOptionType,
				}
			}
		}
		if len(dataV0.IP[0].Tunnel) > 0 {
			dataV1.IP.Tunnel = &securityScreenBlockIPBlockTunnel{
				BadInnerHeader: dataV0.IP[0].Tunnel[0].BadInnerHeader,
				IPInUDPTeredo:  dataV0.IP[0].Tunnel[0].IPInUDPTeredo,
			}
			if len(dataV0.IP[0].Tunnel[0].Gre) > 0 {
				dataV1.IP.Tunnel.Gre = &securityScreenBlockIPBlockTunnelBlockGre{
					Gre4in4: dataV0.IP[0].Tunnel[0].Gre[0].Gre4in4,
					Gre4in6: dataV0.IP[0].Tunnel[0].Gre[0].Gre4in6,
					Gre6in4: dataV0.IP[0].Tunnel[0].Gre[0].Gre6in4,
					Gre6in6: dataV0.IP[0].Tunnel[0].Gre[0].Gre6in6,
				}
			}
			if len(dataV0.IP[0].Tunnel[0].Ipip) > 0 {
				dataV1.IP.Tunnel.Ipip = &securityScreenBlockIPBlockTunnelBlockIpip{
					Dslite:        dataV0.IP[0].Tunnel[0].Ipip[0].Dslite,
					Ipip4in4:      dataV0.IP[0].Tunnel[0].Ipip[0].Ipip4in4,
					Ipip4in6:      dataV0.IP[0].Tunnel[0].Ipip[0].Ipip4in6,
					Ipip6in4:      dataV0.IP[0].Tunnel[0].Ipip[0].Ipip6in4,
					Ipip6in6:      dataV0.IP[0].Tunnel[0].Ipip[0].Ipip6in6,
					Ipip6over4:    dataV0.IP[0].Tunnel[0].Ipip[0].Ipip6over4,
					Ipip6to4relay: dataV0.IP[0].Tunnel[0].Ipip[0].Ipip6to4relay,
					Isatap:        dataV0.IP[0].Tunnel[0].Ipip[0].Isatap,
				}
			}
		}
	}
	if len(dataV0.LimitSession) > 0 {
		dataV1.LimitSession = &securityScreenBlockLimitSession{
			DestinationIPBased: dataV0.LimitSession[0].DestinationIPBased,
			SourceIPBased:      dataV0.LimitSession[0].SourceIPBased,
		}
	}
	if len(dataV0.TCP) > 0 {
		dataV1.TCP = &securityScreenBlockTCP{
			FinNoAck: dataV0.TCP[0].FinNoAck,
			Land:     dataV0.TCP[0].Land,
			NoFlag:   dataV0.TCP[0].NoFlag,
			SynFin:   dataV0.TCP[0].SynFin,
			SynFrag:  dataV0.TCP[0].SynFrag,
			Winnuke:  dataV0.TCP[0].Winnuke,
		}
		if len(dataV0.TCP[0].PortScan) > 0 {
			dataV1.TCP.PortScan = &securityScreenBlockWithThreshold{
				Threshold: dataV0.TCP[0].PortScan[0].Threshold,
			}
		}
		if len(dataV0.TCP[0].Sweep) > 0 {
			dataV1.TCP.Sweep = &securityScreenBlockWithThreshold{
				Threshold: dataV0.TCP[0].Sweep[0].Threshold,
			}
		}
		if len(dataV0.TCP[0].SynAckAckProxy) > 0 {
			dataV1.TCP.SynAckAckProxy = &securityScreenBlockWithThreshold{
				Threshold: dataV0.TCP[0].SynAckAckProxy[0].Threshold,
			}
		}
		if len(dataV0.TCP[0].SynFlood) > 0 {
			dataV1.TCP.SynFlood = &securityScreenBlockTCPBlockSynFlood{
				AlarmThreshold:       dataV0.TCP[0].SynFlood[0].AlarmThreshold,
				AttackThreshold:      dataV0.TCP[0].SynFlood[0].AttackThreshold,
				DestinationThreshold: dataV0.TCP[0].SynFlood[0].DestinationThreshold,
				SourceThreshold:      dataV0.TCP[0].SynFlood[0].SourceThreshold,
				Timeout:              dataV0.TCP[0].SynFlood[0].Timeout,
			}
			for _, blockV0 := range dataV0.TCP[0].SynFlood[0].Whitelist {
				dataV1.TCP.SynFlood.Whitelist = append(dataV1.TCP.SynFlood.Whitelist,
					securityScreenBlockTCPBlockSynFloodBlockWhitelist{
						Name:               blockV0.Name,
						DestinationAddress: blockV0.DestinationAddress,
						SourceAddress:      blockV0.SourceAddress,
					})
			}
		}
	}
	if len(dataV0.UDP) > 0 {
		dataV1.UDP = &securityScreenBlockUDP{}
		if len(dataV0.UDP[0].Flood) > 0 {
			dataV1.UDP.Flood = &securityScreenBlockUDPBlockFlood{
				Threshold: dataV0.UDP[0].Flood[0].Threshold,
				Whitelist: dataV0.UDP[0].Flood[0].Whitelist,
			}
		}
		if len(dataV0.UDP[0].PortScan) > 0 {
			dataV1.UDP.PortScan = &securityScreenBlockWithThreshold{
				Threshold: dataV0.UDP[0].PortScan[0].Threshold,
			}
		}
		if len(dataV0.UDP[0].Sweep) > 0 {
			dataV1.UDP.Sweep = &securityScreenBlockWithThreshold{
				Threshold: dataV0.UDP[0].Sweep[0].Threshold,
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
