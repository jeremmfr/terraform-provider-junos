package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *interfaceLogical) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"st0_also_on_destroy": schema.BoolAttribute{
						Optional: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"disable": schema.BoolAttribute{
						Optional: true,
					},
					"routing_instance": schema.StringAttribute{
						Optional: true,
					},
					"security_inbound_protocols": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"security_inbound_services": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"security_zone": schema.StringAttribute{
						Optional: true,
					},
					"vlan_id": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"vlan_no_compute": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"family_inet": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"filter_input": schema.StringAttribute{
									Optional: true,
								},
								"filter_output": schema.StringAttribute{
									Optional: true,
								},
								"mtu": schema.Int64Attribute{
									Optional: true,
								},
								"sampling_input": schema.BoolAttribute{
									Optional: true,
								},
								"sampling_output": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"address": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"cidr_ip": schema.StringAttribute{
												Required: true,
											},
											"preferred": schema.BoolAttribute{
												Optional: true,
											},
											"primary": schema.BoolAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"vrrp_group": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"identifier": schema.Int64Attribute{
															Required: true,
														},
														"virtual_address": schema.ListAttribute{
															ElementType: types.StringType,
															Required:    true,
														},
														"accept_data": schema.BoolAttribute{
															Optional: true,
														},
														"no_accept_data": schema.BoolAttribute{
															Optional: true,
														},
														"advertise_interval": schema.Int64Attribute{
															Optional: true,
														},
														"advertisements_threshold": schema.Int64Attribute{
															Optional: true,
														},
														"authentication_key": schema.StringAttribute{
															Optional:  true,
															Sensitive: true,
														},
														"authentication_type": schema.StringAttribute{
															Optional: true,
														},
														"preempt": schema.BoolAttribute{
															Optional: true,
														},
														"no_preempt": schema.BoolAttribute{
															Optional: true,
														},
														"priority": schema.Int64Attribute{
															Optional: true,
														},
													},
													Blocks: map[string]schema.Block{
														"track_interface": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"interface": schema.StringAttribute{
																		Required: true,
																	},
																	"priority_cost": schema.Int64Attribute{
																		Required: true,
																	},
																},
															},
														},
														"track_route": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"route": schema.StringAttribute{
																		Required: true,
																	},
																	"routing_instance": schema.StringAttribute{
																		Required: true,
																	},
																	"priority_cost": schema.Int64Attribute{
																		Required: true,
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
								"dhcp": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"srx_old_option_name": schema.BoolAttribute{
												Optional: true,
											},
											"client_identifier_ascii": schema.StringAttribute{
												Optional: true,
											},
											"client_identifier_hexadecimal": schema.StringAttribute{
												Optional: true,
											},
											"client_identifier_prefix_hostname": schema.BoolAttribute{
												Optional: true,
											},
											"client_identifier_prefix_routing_instance_name": schema.BoolAttribute{
												Optional: true,
											},
											"client_identifier_use_interface_description": schema.StringAttribute{
												Optional: true,
											},
											"client_identifier_userid_ascii": schema.StringAttribute{
												Optional: true,
											},
											"client_identifier_userid_hexadecimal": schema.StringAttribute{
												Optional: true,
											},
											"force_discover": schema.BoolAttribute{
												Optional: true,
											},
											"lease_time": schema.Int64Attribute{
												Optional: true,
											},
											"lease_time_infinite": schema.BoolAttribute{
												Optional: true,
											},
											"metric": schema.Int64Attribute{
												Optional: true,
											},
											"no_dns_install": schema.BoolAttribute{
												Optional: true,
											},
											"options_no_hostname": schema.BoolAttribute{
												Optional: true,
											},
											"retransmission_attempt": schema.Int64Attribute{
												Optional: true,
											},
											"retransmission_interval": schema.Int64Attribute{
												Optional: true,
											},
											"server_address": schema.StringAttribute{
												Optional: true,
											},
											"update_server": schema.BoolAttribute{
												Optional: true,
											},
											"vendor_id": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"rpf_check": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"fail_filter": schema.StringAttribute{
												Optional: true,
											},
											"mode_loose": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"family_inet6": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"dad_disable": schema.BoolAttribute{
									Optional: true,
								},
								"filter_input": schema.StringAttribute{
									Optional: true,
								},
								"filter_output": schema.StringAttribute{
									Optional: true,
								},
								"mtu": schema.Int64Attribute{
									Optional: true,
								},
								"sampling_input": schema.BoolAttribute{
									Optional: true,
								},
								"sampling_output": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"address": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"cidr_ip": schema.StringAttribute{
												Required: true,
											},
											"preferred": schema.BoolAttribute{
												Optional: true,
											},
											"primary": schema.BoolAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"vrrp_group": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"identifier": schema.Int64Attribute{
															Required: true,
														},
														"virtual_address": schema.ListAttribute{
															ElementType: types.StringType,
															Required:    true,
														},
														"virtual_link_local_address": schema.StringAttribute{
															Required: true,
														},
														"accept_data": schema.BoolAttribute{
															Optional: true,
														},
														"no_accept_data": schema.BoolAttribute{
															Optional: true,
														},
														"advertise_interval": schema.Int64Attribute{
															Optional: true,
														},
														"advertisements_threshold": schema.Int64Attribute{
															Optional: true,
														},
														"preempt": schema.BoolAttribute{
															Optional: true,
														},
														"no_preempt": schema.BoolAttribute{
															Optional: true,
														},
														"priority": schema.Int64Attribute{
															Optional: true,
														},
													},
													Blocks: map[string]schema.Block{
														"track_interface": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"interface": schema.StringAttribute{
																		Required: true,
																	},
																	"priority_cost": schema.Int64Attribute{
																		Required: true,
																	},
																},
															},
														},
														"track_route": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"route": schema.StringAttribute{
																		Required: true,
																	},
																	"routing_instance": schema.StringAttribute{
																		Required: true,
																	},
																	"priority_cost": schema.Int64Attribute{
																		Required: true,
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
								"dhcpv6_client": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"client_identifier_duid_type": schema.StringAttribute{
												Required: true,
											},
											"client_type": schema.StringAttribute{
												Required: true,
											},
											"client_ia_type_na": schema.BoolAttribute{
												Optional: true,
											},
											"client_ia_type_pd": schema.BoolAttribute{
												Optional: true,
											},
											"no_dns_install": schema.BoolAttribute{
												Optional: true,
											},
											"prefix_delegating_preferred_prefix_length": schema.Int64Attribute{
												Optional: true,
											},
											"prefix_delegating_sub_prefix_length": schema.Int64Attribute{
												Optional: true,
											},
											"rapid_commit": schema.BoolAttribute{
												Optional: true,
											},
											"req_option": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"retransmission_attempt": schema.Int64Attribute{
												Optional: true,
											},
											"update_router_advertisement_interface": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"update_server": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"rpf_check": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"fail_filter": schema.StringAttribute{
												Optional: true,
											},
											"mode_loose": schema.BoolAttribute{
												Optional: true,
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
								"destination": schema.StringAttribute{
									Required: true,
								},
								"source": schema.StringAttribute{
									Required: true,
								},
								"allow_fragmentation": schema.BoolAttribute{
									Optional: true,
								},
								"do_not_fragment": schema.BoolAttribute{
									Optional: true,
								},
								"flow_label": schema.Int64Attribute{
									Optional: true,
								},
								"path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"no_path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"routing_instance_destination": schema.StringAttribute{
									Optional: true,
								},
								"traffic_class": schema.Int64Attribute{
									Optional: true,
								},
								"ttl": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeInterfaceLogicalV0toV1,
		},
	}
}

func upgradeInterfaceLogicalV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                       types.String   `tfsdk:"id"`
		Name                     types.String   `tfsdk:"name"`
		St0AlsoOnDestroy         types.Bool     `tfsdk:"st0_also_on_destroy"`
		Disable                  types.Bool     `tfsdk:"disable"`
		Description              types.String   `tfsdk:"description"`
		RoutingInstance          types.String   `tfsdk:"routing_instance"`
		SecurityInboundProtocols []types.String `tfsdk:"security_inbound_protocols"`
		SecurityInboundServices  []types.String `tfsdk:"security_inbound_services"`
		SecurityZone             types.String   `tfsdk:"security_zone"`
		VlanID                   types.Int64    `tfsdk:"vlan_id"`
		VlanNoCompute            types.Bool     `tfsdk:"vlan_no_compute"`
		FamilyInet               []struct {
			FilterInput    types.String `tfsdk:"filter_input"`
			FilterOutput   types.String `tfsdk:"filter_output"`
			Mtu            types.Int64  `tfsdk:"mtu"`
			SamplingInput  types.Bool   `tfsdk:"sampling_input"`
			SamplingOutput types.Bool   `tfsdk:"sampling_output"`
			Address        []struct {
				CidrIP    types.String `tfsdk:"cidr_ip"`
				Preferred types.Bool   `tfsdk:"preferred"`
				Primary   types.Bool   `tfsdk:"primary"`
				VRRPGroup []struct {
					Identifier              types.Int64    `tfsdk:"identifier"`
					VirtualAddress          []types.String `tfsdk:"virtual_address"`
					AcceptData              types.Bool     `tfsdk:"accept_data"`
					NoAcceptData            types.Bool     `tfsdk:"no_accept_data"`
					AdvertiseInterval       types.Int64    `tfsdk:"advertise_interval"`
					AdvertisementsThreshold types.Int64    `tfsdk:"advertisements_threshold"`
					AuthenticationKey       types.String   `tfsdk:"authentication_key"`
					AuthenticationType      types.String   `tfsdk:"authentication_type"`
					Preempt                 types.Bool     `tfsdk:"preempt"`
					NoPreempt               types.Bool     `tfsdk:"no_preempt"`
					Priority                types.Int64    `tfsdk:"priority"`
					TrackInterface          []struct {
						Interface    types.String `tfsdk:"interface"`
						PriorityCost types.Int64  `tfsdk:"priority_cost"`
					} `tfsdk:"track_interface"`
					TrackRoute []struct {
						Route           types.String `tfsdk:"route"`
						RoutingInstance types.String `tfsdk:"routing_instance"`
						PriorityCost    types.Int64  `tfsdk:"priority_cost"`
					} `tfsdk:"track_route"`
				} `tfsdk:"vrrp_group"`
			} `tfsdk:"address"`
			DHCP []struct {
				SrxOldOptionName                          types.Bool   `tfsdk:"srx_old_option_name"`
				ClientIdentifierASCII                     types.String `tfsdk:"client_identifier_ascii"`
				ClientIdentifierHexadecimal               types.String `tfsdk:"client_identifier_hexadecimal"`
				ClientIdentifierPrefixHostname            types.Bool   `tfsdk:"client_identifier_prefix_hostname"`
				ClientIdentifierPrefixRoutingInstanceName types.Bool   `tfsdk:"client_identifier_prefix_routing_instance_name"`
				ClientIdentifierUseInterfaceDescription   types.String `tfsdk:"client_identifier_use_interface_description"`
				ClientIdentifierUseridASCII               types.String `tfsdk:"client_identifier_userid_ascii"`
				ClientIdentifierUseridHexadecimal         types.String `tfsdk:"client_identifier_userid_hexadecimal"`
				ForceDiscover                             types.Bool   `tfsdk:"force_discover"`
				LeaseTime                                 types.Int64  `tfsdk:"lease_time"`
				LeaseTimeInfinite                         types.Bool   `tfsdk:"lease_time_infinite"`
				Metric                                    types.Int64  `tfsdk:"metric"`
				NoDNSInstall                              types.Bool   `tfsdk:"no_dns_install"`
				OptionsNoHostname                         types.Bool   `tfsdk:"options_no_hostname"`
				RetransmissionAttempt                     types.Int64  `tfsdk:"retransmission_attempt"`
				RetransmissionInterval                    types.Int64  `tfsdk:"retransmission_interval"`
				ServerAddress                             types.String `tfsdk:"server_address"`
				UpdateServer                              types.Bool   `tfsdk:"update_server"`
				VendorID                                  types.String `tfsdk:"vendor_id"`
			} `tfsdk:"dhcp"`
			RPFCheck []struct {
				FailFilter types.String `tfsdk:"fail_filter"`
				ModeLoose  types.Bool   `tfsdk:"mode_loose"`
			} `tfsdk:"rpf_check"`
		} `tfsdk:"family_inet"`
		FamilyInet6 []struct {
			DadDisable     types.Bool   `tfsdk:"dad_disable"`
			FilterInput    types.String `tfsdk:"filter_input"`
			FilterOutput   types.String `tfsdk:"filter_output"`
			Mtu            types.Int64  `tfsdk:"mtu"`
			SamplingInput  types.Bool   `tfsdk:"sampling_input"`
			SamplingOutput types.Bool   `tfsdk:"sampling_output"`
			Address        []struct {
				CidrIP    types.String `tfsdk:"cidr_ip"`
				Preferred types.Bool   `tfsdk:"preferred"`
				Primary   types.Bool   `tfsdk:"primary"`
				VRRPGroup []struct {
					Identifier              types.Int64    `tfsdk:"identifier"`
					VirtualAddress          []types.String `tfsdk:"virtual_address"`
					VirutalLinkLocalAddress types.String   `tfsdk:"virtual_link_local_address"`
					AcceptData              types.Bool     `tfsdk:"accept_data"`
					NoAcceptData            types.Bool     `tfsdk:"no_accept_data"`
					AdvertiseInterval       types.Int64    `tfsdk:"advertise_interval"`
					AdvertisementsThreshold types.Int64    `tfsdk:"advertisements_threshold"`
					Preempt                 types.Bool     `tfsdk:"preempt"`
					NoPreempt               types.Bool     `tfsdk:"no_preempt"`
					Priority                types.Int64    `tfsdk:"priority"`
					TrackInterface          []struct {
						Interface    types.String `tfsdk:"interface"`
						PriorityCost types.Int64  `tfsdk:"priority_cost"`
					} `tfsdk:"track_interface"`
					TrackRoute []struct {
						Route           types.String `tfsdk:"route"`
						RoutingInstance types.String `tfsdk:"routing_instance"`
						PriorityCost    types.Int64  `tfsdk:"priority_cost"`
					} `tfsdk:"track_route"`
				} `tfsdk:"vrrp_group"`
			} `tfsdk:"address"`
			DHCPv6Client []struct {
				ClientIdentifierDuidType              types.String   `tfsdk:"client_identifier_duid_type"`
				ClientType                            types.String   `tfsdk:"client_type"`
				ClientIATypeNA                        types.Bool     `tfsdk:"client_ia_type_na"`
				ClientIATypePD                        types.Bool     `tfsdk:"client_ia_type_pd"`
				NoDNSInstall                          types.Bool     `tfsdk:"no_dns_install"`
				PrefixDelegatingPreferredPrefixLength types.Int64    `tfsdk:"prefix_delegating_preferred_prefix_length"`
				PrefixDelegatingSubPrefixLength       types.Int64    `tfsdk:"prefix_delegating_sub_prefix_length"`
				RapidCommit                           types.Bool     `tfsdk:"rapid_commit"`
				ReqOption                             []types.String `tfsdk:"req_option"`
				RetransmissionAttempt                 types.Int64    `tfsdk:"retransmission_attempt"`
				UpdateRouterAdvertisementInterface    []types.String `tfsdk:"update_router_advertisement_interface"`
				UpdateServer                          types.Bool     `tfsdk:"update_server"`
			} `tfsdk:"dhcpv6_client"`
			RPFCheck []struct {
				FailFilter types.String `tfsdk:"fail_filter"`
				ModeLoose  types.Bool   `tfsdk:"mode_loose"`
			} `tfsdk:"rpf_check"`
		} `tfsdk:"family_inet6"`
		Tunnel []struct {
			Destination                types.String `tfsdk:"destination"`
			Source                     types.String `tfsdk:"source"`
			AllowFragmentation         types.Bool   `tfsdk:"allow_fragmentation"`
			DoNotFragment              types.Bool   `tfsdk:"do_not_fragment"`
			FlowLabel                  types.Int64  `tfsdk:"flow_label"`
			PathMtuDiscovery           types.Bool   `tfsdk:"path_mtu_discovery"`
			NoPathMtuDiscovery         types.Bool   `tfsdk:"no_path_mtu_discovery"`
			RoutingInstanceDestination types.String `tfsdk:"routing_instance_destination"`
			TrafficClass               types.Int64  `tfsdk:"traffic_class"`
			TTL                        types.Int64  `tfsdk:"ttl"`
		} `tfsdk:"tunnel"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 interfaceLogicalData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.St0AlsoOnDestroy = dataV0.St0AlsoOnDestroy
	if !dataV1.St0AlsoOnDestroy.IsNull() && !dataV1.St0AlsoOnDestroy.ValueBool() {
		dataV1.St0AlsoOnDestroy = types.BoolNull()
	}
	dataV1.Description = dataV0.Description
	dataV1.Disable = dataV0.Disable
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.SecurityInboundProtocols = dataV0.SecurityInboundProtocols
	dataV1.SecurityInboundServices = dataV0.SecurityInboundServices
	dataV1.SecurityZone = dataV0.SecurityZone
	dataV1.VlanID = dataV0.VlanID
	dataV1.VlanNoCompute = dataV0.VlanNoCompute
	if !dataV1.VlanNoCompute.IsNull() && !dataV1.VlanNoCompute.ValueBool() {
		dataV1.VlanNoCompute = types.BoolNull()
	}
	if len(dataV0.FamilyInet) > 0 {
		dataV1.FamilyInet = &interfaceLogicalBlockFamilyInet{
			SamplingInput:  dataV0.FamilyInet[0].SamplingInput,
			SamplingOutput: dataV0.FamilyInet[0].SamplingOutput,
			FilterInput:    dataV0.FamilyInet[0].FilterInput,
			FilterOutput:   dataV0.FamilyInet[0].FilterOutput,
			Mtu:            dataV0.FamilyInet[0].Mtu,
		}
		for _, blockV0 := range dataV0.FamilyInet[0].Address {
			blockV1 := interfaceLogicalBlockFamilyInetBlockAddress{
				Preferred: blockV0.Preferred,
				Primary:   blockV0.Primary,
				CidrIP:    blockV0.CidrIP,
			}
			for _, subBlockV0 := range blockV0.VRRPGroup {
				subBlockV1 := interfaceLogicalBlockFamilyInetBlockAddressBlockVRRPGroup{
					AcceptData:              subBlockV0.AcceptData,
					NoAcceptData:            subBlockV0.NoAcceptData,
					Preempt:                 subBlockV0.Preempt,
					NoPreempt:               subBlockV0.NoPreempt,
					Identifier:              subBlockV0.Identifier,
					VirtualAddress:          subBlockV0.VirtualAddress,
					AdvertiseInterval:       subBlockV0.AdvertiseInterval,
					AdvertisementsThreshold: subBlockV0.AdvertisementsThreshold,
					AuthenticationKey:       subBlockV0.AuthenticationKey,
					AuthenticationType:      subBlockV0.AuthenticationType,
					Priority:                subBlockV0.Priority,
				}
				for _, subSubBlockV0 := range subBlockV0.TrackInterface {
					subBlockV1.TrackInterface = append(subBlockV1.TrackInterface,
						interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface{
							Interface:    subSubBlockV0.Interface,
							PriorityCost: subSubBlockV0.PriorityCost,
						},
					)
				}
				for _, subSubBlockV0 := range subBlockV0.TrackRoute {
					subBlockV1.TrackRoute = append(subBlockV1.TrackRoute,
						interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute{
							Route:           subSubBlockV0.Route,
							RoutingInstance: subSubBlockV0.RoutingInstance,
							PriorityCost:    subSubBlockV0.PriorityCost,
						},
					)
				}
				blockV1.VRRPGroup = append(blockV1.VRRPGroup, subBlockV1)
			}
			dataV1.FamilyInet.Address = append(dataV1.FamilyInet.Address, blockV1)
		}
		if len(dataV0.FamilyInet[0].DHCP) > 0 {
			dataV1.FamilyInet.DHCP = &interfaceLogicalBlockFamilyInetBlockDhcp{
				SrxOldOptionName:                          dataV0.FamilyInet[0].DHCP[0].SrxOldOptionName,
				ClientIdentifierPrefixHostname:            dataV0.FamilyInet[0].DHCP[0].ClientIdentifierPrefixHostname,
				ClientIdentifierPrefixRoutingInstanceName: dataV0.FamilyInet[0].DHCP[0].ClientIdentifierPrefixRoutingInstanceName,
				ForceDiscover:                             dataV0.FamilyInet[0].DHCP[0].ForceDiscover,
				LeaseTimeInfinite:                         dataV0.FamilyInet[0].DHCP[0].LeaseTimeInfinite,
				NoDNSInstall:                              dataV0.FamilyInet[0].DHCP[0].NoDNSInstall,
				OptionsNoHostname:                         dataV0.FamilyInet[0].DHCP[0].OptionsNoHostname,
				UpdateServer:                              dataV0.FamilyInet[0].DHCP[0].UpdateServer,
				ClientIdentifierASCII:                     dataV0.FamilyInet[0].DHCP[0].ClientIdentifierASCII,
				ClientIdentifierHexadecimal:               dataV0.FamilyInet[0].DHCP[0].ClientIdentifierHexadecimal,
				ClientIdentifierUseInterfaceDescription:   dataV0.FamilyInet[0].DHCP[0].ClientIdentifierUseInterfaceDescription,
				ClientIdentifierUseridASCII:               dataV0.FamilyInet[0].DHCP[0].ClientIdentifierUseridASCII,
				ClientIdentifierUseridHexadecimal:         dataV0.FamilyInet[0].DHCP[0].ClientIdentifierUseridHexadecimal,
				LeaseTime:                                 dataV0.FamilyInet[0].DHCP[0].LeaseTime,
				Metric:                                    dataV0.FamilyInet[0].DHCP[0].Metric,
				RetransmissionAttempt:                     dataV0.FamilyInet[0].DHCP[0].RetransmissionAttempt,
				RetransmissionInterval:                    dataV0.FamilyInet[0].DHCP[0].RetransmissionInterval,
				ServerAddress:                             dataV0.FamilyInet[0].DHCP[0].ServerAddress,
				VendorID:                                  dataV0.FamilyInet[0].DHCP[0].VendorID,
			}
		}
		if len(dataV0.FamilyInet[0].RPFCheck) > 0 {
			dataV1.FamilyInet.RPFCheck = &interfaceLogicalBlockFamilyBlockRPFCheck{
				FailFilter: dataV0.FamilyInet[0].RPFCheck[0].FailFilter,
				ModeLoose:  dataV0.FamilyInet[0].RPFCheck[0].ModeLoose,
			}
		}
	}
	if len(dataV0.FamilyInet6) > 0 {
		dataV1.FamilyInet6 = &interfaceLogicalBlockFamilyInet6{
			DadDisable:     dataV0.FamilyInet6[0].DadDisable,
			SamplingInput:  dataV0.FamilyInet6[0].SamplingInput,
			SamplingOutput: dataV0.FamilyInet6[0].SamplingOutput,
			FilterInput:    dataV0.FamilyInet6[0].FilterInput,
			FilterOutput:   dataV0.FamilyInet6[0].FilterOutput,
			Mtu:            dataV0.FamilyInet6[0].Mtu,
		}
		for _, blockV0 := range dataV0.FamilyInet6[0].Address {
			blockV1 := interfaceLogicalBlockFamilyInet6BlockAddress{
				Preferred: blockV0.Preferred,
				Primary:   blockV0.Primary,
				CidrIP:    blockV0.CidrIP,
			}
			for _, subBlockV0 := range blockV0.VRRPGroup {
				subBlockV1 := interfaceLogicalBlockFamilyInet6BlockAddressBlockVRRPGroup{
					AcceptData:              subBlockV0.AcceptData,
					NoAcceptData:            subBlockV0.NoAcceptData,
					Preempt:                 subBlockV0.Preempt,
					NoPreempt:               subBlockV0.NoPreempt,
					Identifier:              subBlockV0.Identifier,
					VirtualAddress:          subBlockV0.VirtualAddress,
					VirutalLinkLocalAddress: subBlockV0.VirutalLinkLocalAddress,
					AdvertiseInterval:       subBlockV0.AdvertiseInterval,
					AdvertisementsThreshold: subBlockV0.AdvertisementsThreshold,
					Priority:                subBlockV0.Priority,
				}
				for _, subSubBlockV0 := range subBlockV0.TrackInterface {
					subBlockV1.TrackInterface = append(subBlockV1.TrackInterface,
						interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface{
							Interface:    subSubBlockV0.Interface,
							PriorityCost: subSubBlockV0.PriorityCost,
						},
					)
				}
				for _, subSubBlockV0 := range subBlockV0.TrackRoute {
					subBlockV1.TrackRoute = append(subBlockV1.TrackRoute,
						interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute{
							Route:           subSubBlockV0.Route,
							RoutingInstance: subSubBlockV0.RoutingInstance,
							PriorityCost:    subSubBlockV0.PriorityCost,
						},
					)
				}
				blockV1.VRRPGroup = append(blockV1.VRRPGroup, subBlockV1)
			}
			dataV1.FamilyInet6.Address = append(dataV1.FamilyInet6.Address, blockV1)
		}
		if len(dataV0.FamilyInet6[0].DHCPv6Client) > 0 {
			dataV1.FamilyInet6.DHCPv6Client = &interfaceLogicalBlockFamilyInet6BlockDhcpV6Client{
				ClientIATypeNA:                        dataV0.FamilyInet6[0].DHCPv6Client[0].ClientIATypeNA,
				ClientIATypePD:                        dataV0.FamilyInet6[0].DHCPv6Client[0].ClientIATypePD,
				NoDNSInstall:                          dataV0.FamilyInet6[0].DHCPv6Client[0].NoDNSInstall,
				RapidCommit:                           dataV0.FamilyInet6[0].DHCPv6Client[0].RapidCommit,
				ClientIdentifierDuidType:              dataV0.FamilyInet6[0].DHCPv6Client[0].ClientIdentifierDuidType,
				ClientType:                            dataV0.FamilyInet6[0].DHCPv6Client[0].ClientType,
				PrefixDelegatingPreferredPrefixLength: dataV0.FamilyInet6[0].DHCPv6Client[0].PrefixDelegatingPreferredPrefixLength,
				PrefixDelegatingSubPrefixLength:       dataV0.FamilyInet6[0].DHCPv6Client[0].PrefixDelegatingSubPrefixLength,
				ReqOption:                             dataV0.FamilyInet6[0].DHCPv6Client[0].ReqOption,
				RetransmissionAttempt:                 dataV0.FamilyInet6[0].DHCPv6Client[0].RetransmissionAttempt,
				UpdateRouterAdvertisementInterface:    dataV0.FamilyInet6[0].DHCPv6Client[0].UpdateRouterAdvertisementInterface,
				UpdateServer:                          dataV0.FamilyInet6[0].DHCPv6Client[0].UpdateServer,
			}
		}
		if len(dataV0.FamilyInet6[0].RPFCheck) > 0 {
			dataV1.FamilyInet6.RPFCheck = &interfaceLogicalBlockFamilyBlockRPFCheck{
				FailFilter: dataV0.FamilyInet6[0].RPFCheck[0].FailFilter,
				ModeLoose:  dataV0.FamilyInet6[0].RPFCheck[0].ModeLoose,
			}
		}
	}
	if len(dataV0.Tunnel) > 0 {
		dataV1.Tunnel = &interfaceLogicalBlockTunnel{
			AllowFragmentation:         dataV0.Tunnel[0].AllowFragmentation,
			DoNotFragment:              dataV0.Tunnel[0].DoNotFragment,
			PathMtuDiscovery:           dataV0.Tunnel[0].PathMtuDiscovery,
			NoPathMtuDiscovery:         dataV0.Tunnel[0].NoPathMtuDiscovery,
			Destination:                dataV0.Tunnel[0].Destination,
			Source:                     dataV0.Tunnel[0].Source,
			FlowLabel:                  dataV0.Tunnel[0].FlowLabel,
			RoutingInstanceDestination: dataV0.Tunnel[0].RoutingInstanceDestination,
			TrafficClass:               dataV0.Tunnel[0].TrafficClass,
			TTL:                        dataV0.Tunnel[0].TTL,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
