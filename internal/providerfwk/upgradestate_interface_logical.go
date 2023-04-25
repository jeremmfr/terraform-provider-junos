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
		St0AlsoOnDestroy         types.Bool     `tfsdk:"st0_also_on_destroy"`
		VlanNoCompute            types.Bool     `tfsdk:"vlan_no_compute"`
		Disable                  types.Bool     `tfsdk:"disable"`
		ID                       types.String   `tfsdk:"id"`
		Name                     types.String   `tfsdk:"name"`
		Description              types.String   `tfsdk:"description"`
		RoutingInstance          types.String   `tfsdk:"routing_instance"`
		SecurityInboundProtocols []types.String `tfsdk:"security_inbound_protocols"`
		SecurityInboundServices  []types.String `tfsdk:"security_inbound_services"`
		SecurityZone             types.String   `tfsdk:"security_zone"`
		VlanID                   types.Int64    `tfsdk:"vlan_id"`
		FamilyInet               []struct {
			SamplingInput  types.Bool                                    `tfsdk:"sampling_input"`
			SamplingOutput types.Bool                                    `tfsdk:"sampling_output"`
			FilterInput    types.String                                  `tfsdk:"filter_input"`
			FilterOutput   types.String                                  `tfsdk:"filter_output"`
			Mtu            types.Int64                                   `tfsdk:"mtu"`
			Address        []interfaceLogicalBlockFamilyInetBlockAddress `tfsdk:"address"`
			DHCP           []interfaceLogicalBlockFamilyInetBlockDhcp    `tfsdk:"dhcp"`
			RPFCheck       []interfaceLogicalBlockFamilyBlockRPFCheck    `tfsdk:"rpf_check"`
		} `tfsdk:"family_inet"`
		FamilyInet6 []struct {
			DadDisable     types.Bool                                          `tfsdk:"dad_disable"`
			SamplingInput  types.Bool                                          `tfsdk:"sampling_input"`
			SamplingOutput types.Bool                                          `tfsdk:"sampling_output"`
			FilterInput    types.String                                        `tfsdk:"filter_input"`
			FilterOutput   types.String                                        `tfsdk:"filter_output"`
			Mtu            types.Int64                                         `tfsdk:"mtu"`
			Address        []interfaceLogicalBlockFamilyInet6BlockAddress      `tfsdk:"address"`
			DHCPv6Client   []interfaceLogicalBlockFamilyInet6BlockDhcpV6Client `tfsdk:"dhcpv6_client"`
			RPFCheck       []interfaceLogicalBlockFamilyBlockRPFCheck          `tfsdk:"rpf_check"`
		} `tfsdk:"family_inet6"`
		Tunnel []interfaceLogicalBlockTunnel `tfsdk:"tunnel"`
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
			FilterInput:    dataV0.FamilyInet[0].FilterInput,
			FilterOutput:   dataV0.FamilyInet[0].FilterOutput,
			Mtu:            dataV0.FamilyInet[0].Mtu,
			SamplingInput:  dataV0.FamilyInet[0].SamplingInput,
			SamplingOutput: dataV0.FamilyInet[0].SamplingOutput,
			Address:        dataV0.FamilyInet[0].Address,
		}
		if len(dataV0.FamilyInet[0].DHCP) > 0 {
			dataV1.FamilyInet.DHCP = &dataV0.FamilyInet[0].DHCP[0]
		}
		if len(dataV0.FamilyInet[0].RPFCheck) > 0 {
			dataV1.FamilyInet.RPFCheck = &dataV0.FamilyInet[0].RPFCheck[0]
		}
	}
	if len(dataV0.FamilyInet6) > 0 {
		dataV1.FamilyInet6 = &interfaceLogicalBlockFamilyInet6{
			DadDisable:     dataV0.FamilyInet6[0].DadDisable,
			FilterInput:    dataV0.FamilyInet6[0].FilterInput,
			FilterOutput:   dataV0.FamilyInet6[0].FilterOutput,
			Mtu:            dataV0.FamilyInet6[0].Mtu,
			SamplingInput:  dataV0.FamilyInet6[0].SamplingInput,
			SamplingOutput: dataV0.FamilyInet6[0].SamplingOutput,
			Address:        dataV0.FamilyInet6[0].Address,
		}
		if len(dataV0.FamilyInet6[0].DHCPv6Client) > 0 {
			dataV1.FamilyInet6.DHCPv6Client = &dataV0.FamilyInet6[0].DHCPv6Client[0]
		}
		if len(dataV0.FamilyInet6[0].RPFCheck) > 0 {
			dataV1.FamilyInet6.RPFCheck = &dataV0.FamilyInet6[0].RPFCheck[0]
		}
	}
	if len(dataV0.Tunnel) > 0 {
		dataV1.Tunnel = &dataV0.Tunnel[0]
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
