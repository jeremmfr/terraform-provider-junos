package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *accessAddressAssignmentPool) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"routing_instance": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"active_drain": schema.BoolAttribute{
						Optional: true,
					},
					"hold_down": schema.BoolAttribute{
						Optional: true,
					},
					"link": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"family": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required: true,
								},
								"network": schema.StringAttribute{
									Required: true,
								},
								"excluded_address": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"xauth_attributes_primary_dns": schema.StringAttribute{
									Optional: true,
								},
								"xauth_attributes_primary_wins": schema.StringAttribute{
									Optional: true,
								},
								"xauth_attributes_secondary_dns": schema.StringAttribute{
									Optional: true,
								},
								"xauth_attributes_secondary_wins": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"dhcp_attributes": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"boot_file": schema.StringAttribute{
												Optional: true,
											},
											"boot_server": schema.StringAttribute{
												Optional: true,
											},
											"dns_server": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"domain_name": schema.StringAttribute{
												Optional: true,
											},
											"exclude_prefix_len": schema.Int64Attribute{
												Optional: true,
											},
											"grace_period": schema.Int64Attribute{
												Optional: true,
											},
											"maximum_lease_time": schema.Int64Attribute{
												Optional: true,
											},
											"maximum_lease_time_infinite": schema.BoolAttribute{
												Optional: true,
											},
											"name_server": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"netbios_node_type": schema.StringAttribute{
												Optional: true,
											},
											"next_server": schema.StringAttribute{
												Optional: true,
											},
											"option": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"preferred_lifetime": schema.Int64Attribute{
												Optional: true,
											},
											"preferred_lifetime_infinite": schema.BoolAttribute{
												Optional: true,
											},
											"propagate_ppp_settings": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"propagate_settings": schema.StringAttribute{
												Optional: true,
											},
											"router": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"server_identifier": schema.StringAttribute{
												Optional: true,
											},
											"sip_server_inet_address": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"sip_server_inet_domain_name": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"sip_server_inet6_address": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"sip_server_inet6_domain_name": schema.StringAttribute{
												Optional: true,
											},
											"t1_percentage": schema.Int64Attribute{
												Optional: true,
											},
											"t1_renewal_time": schema.Int64Attribute{
												Optional: true,
											},
											"t2_percentage": schema.Int64Attribute{
												Optional: true,
											},
											"t2_rebinding_time": schema.Int64Attribute{
												Optional: true,
											},
											"tftp_server": schema.StringAttribute{
												Optional: true,
											},
											"valid_lifetime": schema.Int64Attribute{
												Optional: true,
											},
											"valid_lifetime_infinite": schema.BoolAttribute{
												Optional: true,
											},
											"wins_server": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
										},
										Blocks: map[string]schema.Block{
											"option_match_82_circuit_id": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"value": schema.StringAttribute{
															Required: true,
														},
														"range": schema.StringAttribute{
															Required: true,
														},
													},
												},
											},
											"option_match_82_remote_id": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"value": schema.StringAttribute{
															Required: true,
														},
														"range": schema.StringAttribute{
															Required: true,
														},
													},
												},
											},
										},
									},
								},
								"excluded_range": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required: true,
											},
											"low": schema.StringAttribute{
												Required: true,
											},
											"high": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
								"host": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required: true,
											},
											"hardware_address": schema.StringAttribute{
												Required: true,
											},
											"ip_address": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
								"inet_range": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required: true,
											},
											"low": schema.StringAttribute{
												Required: true,
											},
											"high": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
								"inet6_range": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required: true,
											},
											"low": schema.StringAttribute{
												Optional: true,
											},
											"high": schema.StringAttribute{
												Optional: true,
											},
											"prefix_length": schema.Int64Attribute{
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
			StateUpgrader: upgradeAccessAddressAssignmentPoolStateV0toV1,
		},
	}
}

func upgradeAccessAddressAssignmentPoolStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID              types.String `tfsdk:"id"`
		Name            types.String `tfsdk:"name"`
		RoutingInstance types.String `tfsdk:"routing_instance"`
		ActiveDrain     types.Bool   `tfsdk:"active_drain"`
		HoldDown        types.Bool   `tfsdk:"hold_down"`
		Link            types.String `tfsdk:"link"`
		Family          []struct {
			Type                         types.String   `tfsdk:"type"`
			Network                      types.String   `tfsdk:"network"`
			ExcludedAddress              []types.String `tfsdk:"excluded_address"`
			XauthAttributesPrimaryDNS    types.String   `tfsdk:"xauth_attributes_primary_dns"`
			XauthAttributesPrimaryWins   types.String   `tfsdk:"xauth_attributes_primary_wins"`
			XauthAttributesSecondaryDNS  types.String   `tfsdk:"xauth_attributes_secondary_dns"`
			XauthAttributesSecondaryWins types.String   `tfsdk:"xauth_attributes_secondary_wins"`
			DhcpAttributes               []struct {
				BootFile                  types.String   `tfsdk:"boot_file"`
				BootServer                types.String   `tfsdk:"boot_server"`
				DNSServer                 []types.String `tfsdk:"dns_server"`
				DomainName                types.String   `tfsdk:"domain_name"`
				ExcludePrefixLen          types.Int64    `tfsdk:"exclude_prefix_len"`
				GracePeriod               types.Int64    `tfsdk:"grace_period"`
				MaximumLeaseTime          types.Int64    `tfsdk:"maximum_lease_time"`
				MaximumLeaseTimeInfinite  types.Bool     `tfsdk:"maximum_lease_time_infinite"`
				NameServer                []types.String `tfsdk:"name_server"`
				NetbiosNodeType           types.String   `tfsdk:"netbios_node_type"`
				NextServer                types.String   `tfsdk:"next_server"`
				Option                    []types.String `tfsdk:"option"`
				PreferredLifetime         types.Int64    `tfsdk:"preferred_lifetime"`
				PreferredLifetimeInfinite types.Bool     `tfsdk:"preferred_lifetime_infinite"`
				PropagatePppSettings      []types.String `tfsdk:"propagate_ppp_settings"`
				PropagateSettings         types.String   `tfsdk:"propagate_settings"`
				Router                    []types.String `tfsdk:"router"`
				ServerIdentifier          types.String   `tfsdk:"server_identifier"`
				SIPServerInetAddress      []types.String `tfsdk:"sip_server_inet_address"`
				SIPServerInetDomainName   []types.String `tfsdk:"sip_server_inet_domain_name"`
				SIPServerInet6Address     []types.String `tfsdk:"sip_server_inet6_address"`
				SIPServerInet6DomainName  types.String   `tfsdk:"sip_server_inet6_domain_name"`
				T1Percentage              types.Int64    `tfsdk:"t1_percentage"`
				T1RenewalTime             types.Int64    `tfsdk:"t1_renewal_time"`
				T2Percentage              types.Int64    `tfsdk:"t2_percentage"`
				T2RebindingTime           types.Int64    `tfsdk:"t2_rebinding_time"`
				TftpServer                types.String   `tfsdk:"tftp_server"`
				ValidLifetime             types.Int64    `tfsdk:"valid_lifetime"`
				ValidLifetimeInfinite     types.Bool     `tfsdk:"valid_lifetime_infinite"`
				WinsServer                []types.String `tfsdk:"wins_server"`
				OptionMatch82CircuitID    []struct {
					Value types.String `tfsdk:"value"`
					Range types.String `tfsdk:"range"`
				} `tfsdk:"option_match_82_circuit_id"`
				OptionMatch82RemoteID []struct {
					Value types.String `tfsdk:"value"`
					Range types.String `tfsdk:"range"`
				} `tfsdk:"option_match_82_remote_id"`
			} `tfsdk:"dhcp_attributes"`
			ExcludedRange []struct {
				Name types.String `tfsdk:"name"`
				Low  types.String `tfsdk:"low"`
				High types.String `tfsdk:"high"`
			} `tfsdk:"excluded_range"`
			Host []struct {
				Name            types.String `tfsdk:"name"`
				HardwareAddress types.String `tfsdk:"hardware_address"`
				IPAddress       types.String `tfsdk:"ip_address"`
			} `tfsdk:"host"`
			InetRange []struct {
				Name types.String `tfsdk:"name"`
				Low  types.String `tfsdk:"low"`
				High types.String `tfsdk:"high"`
			} `tfsdk:"inet_range"`
			Inet6Range []struct {
				Name         types.String `tfsdk:"name"`
				Low          types.String `tfsdk:"low"`
				High         types.String `tfsdk:"high"`
				PrefixLength types.Int64  `tfsdk:"prefix_length"`
			} `tfsdk:"inet6_range"`
		} `tfsdk:"family"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 accessAddressAssignmentPoolData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.ActiveDrain = dataV0.ActiveDrain
	dataV1.HoldDown = dataV0.HoldDown
	dataV1.Link = dataV0.Link
	if len(dataV0.Family) > 0 {
		dataV1.Family = &accessAddressAssignmentPoolBlockFamily{
			Type:                         dataV0.Family[0].Type,
			Network:                      dataV0.Family[0].Network,
			ExcludedAddress:              dataV0.Family[0].ExcludedAddress,
			XauthAttributesPrimaryDNS:    dataV0.Family[0].XauthAttributesPrimaryDNS,
			XauthAttributesPrimaryWins:   dataV0.Family[0].XauthAttributesPrimaryWins,
			XauthAttributesSecondaryDNS:  dataV0.Family[0].XauthAttributesSecondaryDNS,
			XauthAttributesSecondaryWins: dataV0.Family[0].XauthAttributesSecondaryWins,
		}
		if len(dataV0.Family[0].DhcpAttributes) > 0 {
			dataV1.Family.DhcpAttributes = &accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributes{
				BootFile:                  dataV0.Family[0].DhcpAttributes[0].BootFile,
				BootServer:                dataV0.Family[0].DhcpAttributes[0].BootServer,
				DNSServer:                 dataV0.Family[0].DhcpAttributes[0].DNSServer,
				DomainName:                dataV0.Family[0].DhcpAttributes[0].DomainName,
				ExcludePrefixLen:          dataV0.Family[0].DhcpAttributes[0].ExcludePrefixLen,
				GracePeriod:               dataV0.Family[0].DhcpAttributes[0].GracePeriod,
				MaximumLeaseTime:          dataV0.Family[0].DhcpAttributes[0].MaximumLeaseTime,
				MaximumLeaseTimeInfinite:  dataV0.Family[0].DhcpAttributes[0].MaximumLeaseTimeInfinite,
				NameServer:                dataV0.Family[0].DhcpAttributes[0].NameServer,
				NetbiosNodeType:           dataV0.Family[0].DhcpAttributes[0].NetbiosNodeType,
				NextServer:                dataV0.Family[0].DhcpAttributes[0].NextServer,
				Option:                    dataV0.Family[0].DhcpAttributes[0].Option,
				PreferredLifetime:         dataV0.Family[0].DhcpAttributes[0].PreferredLifetime,
				PreferredLifetimeInfinite: dataV0.Family[0].DhcpAttributes[0].PreferredLifetimeInfinite,
				PropagatePppSettings:      dataV0.Family[0].DhcpAttributes[0].PropagatePppSettings,
				PropagateSettings:         dataV0.Family[0].DhcpAttributes[0].PropagateSettings,
				Router:                    dataV0.Family[0].DhcpAttributes[0].Router,
				ServerIdentifier:          dataV0.Family[0].DhcpAttributes[0].ServerIdentifier,
				SIPServerInetAddress:      dataV0.Family[0].DhcpAttributes[0].SIPServerInetAddress,
				SIPServerInetDomainName:   dataV0.Family[0].DhcpAttributes[0].SIPServerInetDomainName,
				SIPServerInet6Address:     dataV0.Family[0].DhcpAttributes[0].SIPServerInet6Address,
				SIPServerInet6DomainName:  dataV0.Family[0].DhcpAttributes[0].SIPServerInet6DomainName,
				T1Percentage:              dataV0.Family[0].DhcpAttributes[0].T1Percentage,
				T1RenewalTime:             dataV0.Family[0].DhcpAttributes[0].T1RenewalTime,
				T2Percentage:              dataV0.Family[0].DhcpAttributes[0].T2Percentage,
				T2RebindingTime:           dataV0.Family[0].DhcpAttributes[0].T2RebindingTime,
				TftpServer:                dataV0.Family[0].DhcpAttributes[0].TftpServer,
				ValidLifetime:             dataV0.Family[0].DhcpAttributes[0].ValidLifetime,
				ValidLifetimeInfinite:     dataV0.Family[0].DhcpAttributes[0].ValidLifetimeInfinite,
				WinsServer:                dataV0.Family[0].DhcpAttributes[0].WinsServer,
			}
			for _, blockV0 := range dataV0.Family[0].DhcpAttributes[0].OptionMatch82CircuitID {
				dataV1.Family.DhcpAttributes.OptionMatch82CircuitID = append(dataV1.Family.DhcpAttributes.OptionMatch82CircuitID,
					accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82{
						Value: blockV0.Value,
						Range: blockV0.Range,
					},
				)
			}
			for _, blockV0 := range dataV0.Family[0].DhcpAttributes[0].OptionMatch82RemoteID {
				dataV1.Family.DhcpAttributes.OptionMatch82RemoteID = append(dataV1.Family.DhcpAttributes.OptionMatch82RemoteID,
					accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82{
						Value: blockV0.Value,
						Range: blockV0.Range,
					},
				)
			}
		}
		for _, blockV0 := range dataV0.Family[0].ExcludedRange {
			dataV1.Family.ExcludedRange = append(dataV1.Family.ExcludedRange,
				accessAddressAssignmentPoolBlockFamilyBlockExcludedRange{
					Name: blockV0.Name,
					Low:  blockV0.Low,
					High: blockV0.High,
				},
			)
		}
		for _, blockV0 := range dataV0.Family[0].Host {
			dataV1.Family.Host = append(dataV1.Family.Host,
				accessAddressAssignmentPoolBlockFamilyBlockHost{
					Name:            blockV0.Name,
					HardwareAddress: blockV0.HardwareAddress,
					IPAddress:       blockV0.IPAddress,
				},
			)
		}
		for _, blockV0 := range dataV0.Family[0].InetRange {
			dataV1.Family.InetRange = append(dataV1.Family.InetRange,
				accessAddressAssignmentPoolBlockFamilyBlockInetRange{
					Name: blockV0.Name,
					Low:  blockV0.Low,
					High: blockV0.High,
				},
			)
		}
		for _, blockV0 := range dataV0.Family[0].Inet6Range {
			dataV1.Family.Inet6Range = append(dataV1.Family.Inet6Range,
				accessAddressAssignmentPoolBlockFamilyBlockInet6Range{
					Name:         blockV0.Name,
					Low:          blockV0.Low,
					High:         blockV0.High,
					PrefixLength: blockV0.PrefixLength,
				},
			)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
