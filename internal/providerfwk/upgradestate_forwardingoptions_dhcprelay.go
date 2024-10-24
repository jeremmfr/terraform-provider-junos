package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *forwardingoptionsDhcprelay) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"routing_instance": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"version": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"access_profile": schema.StringAttribute{
						Optional: true,
					},
					"active_server_group": schema.StringAttribute{
						Optional: true,
					},
					"active_server_group_allow_server_change": schema.BoolAttribute{
						Optional: true,
					},
					"arp_inspection": schema.BoolAttribute{
						Optional: true,
					},
					"authentication_password": schema.StringAttribute{
						Optional: true,
					},
					"client_response_ttl": schema.Int64Attribute{
						Optional: true,
					},
					"duplicate_clients_in_subnet": schema.StringAttribute{
						Optional: true,
					},
					"duplicate_clients_incoming_interface": schema.BoolAttribute{
						Optional: true,
					},
					"dynamic_profile": schema.StringAttribute{
						Optional: true,
					},
					"dynamic_profile_aggregate_clients": schema.BoolAttribute{
						Optional: true,
					},
					"dynamic_profile_aggregate_clients_action": schema.StringAttribute{
						Optional: true,
					},
					"dynamic_profile_use_primary": schema.StringAttribute{
						Optional: true,
					},
					"exclude_relay_agent_identifier": schema.BoolAttribute{
						Optional: true,
					},
					"forward_only": schema.BoolAttribute{
						Optional: true,
					},
					"forward_only_replies": schema.BoolAttribute{
						Optional: true,
					},
					"forward_only_routing_instance": schema.StringAttribute{
						Optional: true,
					},
					"forward_snooped_clients": schema.StringAttribute{
						Optional: true,
					},
					"liveness_detection_failure_action": schema.StringAttribute{
						Optional: true,
					},
					"maximum_hop_count": schema.Int64Attribute{
						Optional: true,
					},
					"minimum_wait_time": schema.Int64Attribute{
						Optional: true,
					},
					"no_snoop": schema.BoolAttribute{
						Optional: true,
					},
					"persistent_storage_automatic": schema.BoolAttribute{
						Optional: true,
					},
					"relay_agent_option_79": schema.BoolAttribute{
						Optional: true,
					},
					"remote_id_mismatch_disconnect": schema.BoolAttribute{
						Optional: true,
					},
					"route_suppression_access": schema.BoolAttribute{
						Optional: true,
					},
					"route_suppression_access_internal": schema.BoolAttribute{
						Optional: true,
					},
					"route_suppression_destination": schema.BoolAttribute{
						Optional: true,
					},
					"server_match_default_action": schema.StringAttribute{
						Optional: true,
					},
					"server_response_time": schema.Int64Attribute{
						Optional: true,
					},
					"service_profile": schema.StringAttribute{
						Optional: true,
					},
					"short_cycle_protection_lockout_max_time": schema.Int64Attribute{
						Optional: true,
					},
					"short_cycle_protection_lockout_min_time": schema.Int64Attribute{
						Optional: true,
					},
					"source_ip_change": schema.BoolAttribute{
						Optional: true,
					},
					"vendor_specific_information_host_name": schema.BoolAttribute{
						Optional: true,
					},
					"vendor_specific_information_location": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"active_leasequery": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"idle_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"peer_address": schema.StringAttribute{
									Optional: true,
								},
								"timeout": schema.Int64Attribute{
									Optional: true,
								},
								"topology_discover": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"authentication_username_include": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"circuit_type": schema.BoolAttribute{
									Optional: true,
								},
								"client_id": schema.BoolAttribute{
									Optional: true,
								},
								"client_id_exclude_headers": schema.BoolAttribute{
									Optional: true,
								},
								"client_id_use_automatic_ascii_hex_encoding": schema.BoolAttribute{
									Optional: true,
								},
								"delimiter": schema.StringAttribute{
									Optional: true,
								},
								"domain_name": schema.StringAttribute{
									Optional: true,
								},
								"interface_description": schema.StringAttribute{
									Optional: true,
								},
								"interface_name": schema.BoolAttribute{
									Optional: true,
								},
								"mac_address": schema.BoolAttribute{
									Optional: true,
								},
								"option_60": schema.BoolAttribute{
									Optional: true,
								},
								"option_82": schema.BoolAttribute{
									Optional: true,
								},
								"option_82_circuit_id": schema.BoolAttribute{
									Optional: true,
								},
								"option_82_remote_id": schema.BoolAttribute{
									Optional: true,
								},
								"relay_agent_interface_id": schema.BoolAttribute{
									Optional: true,
								},
								"relay_agent_remote_id": schema.BoolAttribute{
									Optional: true,
								},
								"relay_agent_subscriber_id": schema.BoolAttribute{
									Optional: true,
								},
								"routing_instance_name": schema.BoolAttribute{
									Optional: true,
								},
								"user_prefix": schema.StringAttribute{
									Optional: true,
								},
								"vlan_tags": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"bulk_leasequery": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"attempts": schema.Int64Attribute{
									Optional: true,
								},
								"timeout": schema.Int64Attribute{
									Optional: true,
								},
								"trigger_automatic": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"lease_time_validation": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"lease_time_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"violation_action_drop": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"leasequery": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"attempts": schema.Int64Attribute{
									Optional: true,
								},
								"timeout": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"liveness_detection_method_bfd": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
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
								"no_adaptation": schema.BoolAttribute{
									Optional: true,
								},
								"session_mode": schema.StringAttribute{
									Optional: true,
								},
								"transmit_interval_minimum": schema.Int64Attribute{
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
					"liveness_detection_method_layer2": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"max_consecutive_retries": schema.Int64Attribute{
									Optional: true,
								},
								"transmit_interval": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"overrides_v4": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"allow_no_end_option": schema.BoolAttribute{
									Optional: true,
								},
								"allow_snooped_clients": schema.BoolAttribute{
									Optional: true,
								},
								"no_allow_snooped_clients": schema.BoolAttribute{
									Optional: true,
								},
								"always_write_giaddr": schema.BoolAttribute{
									Optional: true,
								},
								"always_write_option_82": schema.BoolAttribute{
									Optional: true,
								},
								"asymmetric_lease_time": schema.Int64Attribute{
									Optional: true,
								},
								"bootp_support": schema.BoolAttribute{
									Optional: true,
								},
								"client_discover_match": schema.StringAttribute{
									Optional: true,
								},
								"delay_authentication": schema.BoolAttribute{
									Optional: true,
								},
								"delete_binding_on_renegotiation": schema.BoolAttribute{
									Optional: true,
								},
								"disable_relay": schema.BoolAttribute{
									Optional: true,
								},
								"dual_stack": schema.StringAttribute{
									Optional: true,
								},
								"interface_client_limit": schema.Int64Attribute{
									Optional: true,
								},
								"layer2_unicast_replies": schema.BoolAttribute{
									Optional: true,
								},
								"no_bind_on_request": schema.BoolAttribute{
									Optional: true,
								},
								"no_unicast_replies": schema.BoolAttribute{
									Optional: true,
								},
								"proxy_mode": schema.BoolAttribute{
									Optional: true,
								},
								"relay_source": schema.StringAttribute{
									Optional: true,
								},
								"replace_ip_source_with_giaddr": schema.BoolAttribute{
									Optional: true,
								},
								"send_release_on_delete": schema.BoolAttribute{
									Optional: true,
								},
								"trust_option_82": schema.BoolAttribute{
									Optional: true,
								},
								"user_defined_option_82": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"overrides_v6": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"allow_snooped_clients": schema.BoolAttribute{
									Optional: true,
								},
								"no_allow_snooped_clients": schema.BoolAttribute{
									Optional: true,
								},
								"always_process_option_request_option": schema.BoolAttribute{
									Optional: true,
								},
								"asymmetric_lease_time": schema.Int64Attribute{
									Optional: true,
								},
								"asymmetric_prefix_lease_time": schema.Int64Attribute{
									Optional: true,
								},
								"client_negotiation_match_incoming_interface": schema.BoolAttribute{
									Optional: true,
								},
								"delay_authentication": schema.BoolAttribute{
									Optional: true,
								},
								"delete_binding_on_renegotiation": schema.BoolAttribute{
									Optional: true,
								},
								"dual_stack": schema.StringAttribute{
									Optional: true,
								},
								"interface_client_limit": schema.Int64Attribute{
									Optional: true,
								},
								"no_bind_on_request": schema.BoolAttribute{
									Optional: true,
								},
								"relay_source": schema.StringAttribute{
									Optional: true,
								},
								"send_release_on_delete": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"relay_agent_interface_id": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"keep_incoming_id": schema.BoolAttribute{
									Optional: true,
								},
								"keep_incoming_id_strict": schema.BoolAttribute{
									Optional: true,
								},
								"include_irb_and_l2": schema.BoolAttribute{
									Optional: true,
								},
								"no_vlan_interface_name": schema.BoolAttribute{
									Optional: true,
								},
								"prefix_host_name": schema.BoolAttribute{
									Optional: true,
								},
								"prefix_routing_instance_name": schema.BoolAttribute{
									Optional: true,
								},
								"use_interface_description": schema.StringAttribute{
									Optional: true,
								},
								"use_option_82": schema.BoolAttribute{
									Optional: true,
								},
								"use_option_82_strict": schema.BoolAttribute{
									Optional: true,
								},
								"use_vlan_id": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"relay_agent_remote_id": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"keep_incoming_id": schema.BoolAttribute{
									Optional: true,
								},
								"include_irb_and_l2": schema.BoolAttribute{
									Optional: true,
								},
								"no_vlan_interface_name": schema.BoolAttribute{
									Optional: true,
								},
								"prefix_host_name": schema.BoolAttribute{
									Optional: true,
								},
								"prefix_routing_instance_name": schema.BoolAttribute{
									Optional: true,
								},
								"use_interface_description": schema.StringAttribute{
									Optional: true,
								},
								"use_option_82": schema.BoolAttribute{
									Optional: true,
								},
								"use_option_82_strict": schema.BoolAttribute{
									Optional: true,
								},
								"use_vlan_id": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"relay_option": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"option_order": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
							},
							Blocks: map[string]schema.Block{
								"option_15": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"compare": schema.StringAttribute{
												Required: true,
											},
											"value_type": schema.StringAttribute{
												Required: true,
											},
											"value": schema.StringAttribute{
												Required: true,
											},
											"action": schema.StringAttribute{
												Required: true,
											},
											"group": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"option_15_default_action": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"group": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"option_16": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"compare": schema.StringAttribute{
												Required: true,
											},
											"value_type": schema.StringAttribute{
												Required: true,
											},
											"value": schema.StringAttribute{
												Required: true,
											},
											"action": schema.StringAttribute{
												Required: true,
											},
											"group": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"option_16_default_action": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"group": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"option_60": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"compare": schema.StringAttribute{
												Required: true,
											},
											"value_type": schema.StringAttribute{
												Required: true,
											},
											"value": schema.StringAttribute{
												Required: true,
											},
											"action": schema.StringAttribute{
												Required: true,
											},
											"group": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"option_60_default_action": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"group": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"option_77": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"compare": schema.StringAttribute{
												Required: true,
											},
											"value_type": schema.StringAttribute{
												Required: true,
											},
											"value": schema.StringAttribute{
												Required: true,
											},
											"action": schema.StringAttribute{
												Required: true,
											},
											"group": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"option_77_default_action": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"group": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"relay_option_82": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"exclude_relay_agent_identifier": schema.BoolAttribute{
									Optional: true,
								},
								"link_selection": schema.BoolAttribute{
									Optional: true,
								},
								"server_id_override": schema.BoolAttribute{
									Optional: true,
								},
								"vendor_specific_host_name": schema.BoolAttribute{
									Optional: true,
								},
								"vendor_specific_location": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"circuit_id": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"include_irb_and_l2": schema.BoolAttribute{
												Optional: true,
											},
											"keep_incoming_circuit_id": schema.BoolAttribute{
												Optional: true,
											},
											"no_vlan_interface_name": schema.BoolAttribute{
												Optional: true,
											},
											"prefix_host_name": schema.BoolAttribute{
												Optional: true,
											},
											"prefix_routing_instance_name": schema.BoolAttribute{
												Optional: true,
											},
											"use_interface_description": schema.StringAttribute{
												Optional: true,
											},
											"use_vlan_id": schema.BoolAttribute{
												Optional: true,
											},
											"user_defined": schema.BoolAttribute{
												Optional: true,
											},
											"vlan_id_only": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"remote_id": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"hostname_only": schema.BoolAttribute{
												Optional: true,
											},
											"include_irb_and_l2": schema.BoolAttribute{
												Optional: true,
											},
											"keep_incoming_remote_id": schema.BoolAttribute{
												Optional: true,
											},
											"no_vlan_interface_name": schema.BoolAttribute{
												Optional: true,
											},
											"prefix_host_name": schema.BoolAttribute{
												Optional: true,
											},
											"prefix_routing_instance_name": schema.BoolAttribute{
												Optional: true,
											},
											"use_interface_description": schema.StringAttribute{
												Optional: true,
											},
											"use_string": schema.StringAttribute{
												Optional: true,
											},
											"use_vlan_id": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"server_match_address": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"address": schema.StringAttribute{
									Required: true,
								},
								"action": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"server_match_duid": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"compare": schema.StringAttribute{
									Required: true,
								},
								"value_type": schema.StringAttribute{
									Required: true,
								},
								"value": schema.StringAttribute{
									Required: true,
								},
								"action": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeForwardingoptionsDhcprelayV0toV1,
		},
	}
}

func upgradeForwardingoptionsDhcprelayV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	//nolint:lll
	type modelV0 struct {
		ID                                   types.String `tfsdk:"id"`
		RoutingInstance                      types.String `tfsdk:"routing_instance"`
		Version                              types.String `tfsdk:"version"`
		AccessProfile                        types.String `tfsdk:"access_profile"`
		ActiveServerGroup                    types.String `tfsdk:"active_server_group"`
		ActiveServerGroupAllowServerChange   types.Bool   `tfsdk:"active_server_group_allow_server_change"`
		ARPInspection                        types.Bool   `tfsdk:"arp_inspection"`
		AuthenticationPassword               types.String `tfsdk:"authentication_password"`
		ClientResponseTTL                    types.Int64  `tfsdk:"client_response_ttl"`
		DuplicateClientsInSubnet             types.String `tfsdk:"duplicate_clients_in_subnet"`
		DuplicateClientsIncomingInterface    types.Bool   `tfsdk:"duplicate_clients_incoming_interface"`
		DynamicProfile                       types.String `tfsdk:"dynamic_profile"`
		DynamicProfileAggregateClients       types.Bool   `tfsdk:"dynamic_profile_aggregate_clients"`
		DynamicProfileAggregateClientsAction types.String `tfsdk:"dynamic_profile_aggregate_clients_action"`
		DynamicProfileUsePrimary             types.String `tfsdk:"dynamic_profile_use_primary"`
		ExcludeRelayAgentIdentifier          types.Bool   `tfsdk:"exclude_relay_agent_identifier"`
		ForwardOnly                          types.Bool   `tfsdk:"forward_only"`
		ForwardOnlyReplies                   types.Bool   `tfsdk:"forward_only_replies"`
		ForwardOnlyRoutingInstance           types.String `tfsdk:"forward_only_routing_instance"`
		ForwardSnoopedClients                types.String `tfsdk:"forward_snooped_clients"`
		LivenessDetectionFailureAction       types.String `tfsdk:"liveness_detection_failure_action"`
		MaximumHopCount                      types.Int64  `tfsdk:"maximum_hop_count"`
		MinimumWaitTime                      types.Int64  `tfsdk:"minimum_wait_time"`
		NoSnoop                              types.Bool   `tfsdk:"no_snoop"`
		PersistentStorageAutomatic           types.Bool   `tfsdk:"persistent_storage_automatic"`
		RelayAgentOption79                   types.Bool   `tfsdk:"relay_agent_option_79"`
		RemoteIDMismatchDisconnect           types.Bool   `tfsdk:"remote_id_mismatch_disconnect"`
		RouteSuppressionAccess               types.Bool   `tfsdk:"route_suppression_access"`
		RouteSuppressionAccessInternal       types.Bool   `tfsdk:"route_suppression_access_internal"`
		RouteSuppressionDestination          types.Bool   `tfsdk:"route_suppression_destination"`
		ServerMatchDefaultAction             types.String `tfsdk:"server_match_default_action"`
		ServerResponseTime                   types.Int64  `tfsdk:"server_response_time"`
		ServiceProfile                       types.String `tfsdk:"service_profile"`
		ShortCycleProtectionLockoutMaxTime   types.Int64  `tfsdk:"short_cycle_protection_lockout_max_time"`
		ShortCycleProtectionLockoutMinTime   types.Int64  `tfsdk:"short_cycle_protection_lockout_min_time"`
		SourceIPChange                       types.Bool   `tfsdk:"source_ip_change"`
		VendorSpecificInformationHostName    types.Bool   `tfsdk:"vendor_specific_information_host_name"`
		VendorSpecificInformationLocation    types.Bool   `tfsdk:"vendor_specific_information_location"`
		ActiveLeasequery                     []struct {
			IdleTimeout      types.Int64  `tfsdk:"idle_timeout"`
			PeerAddress      types.String `tfsdk:"peer_address"`
			Timeout          types.Int64  `tfsdk:"timeout"`
			TopologyDiscover types.Bool   `tfsdk:"topology_discover"`
		} `tfsdk:"active_leasequery"`
		AuthenticationUsernameInclude []struct {
			CircuitType                          types.Bool   `tfsdk:"circuit_type"`
			ClientID                             types.Bool   `tfsdk:"client_id"`
			ClientIDExcludeHeaders               types.Bool   `tfsdk:"client_id_exclude_headers"`
			ClientIDUseAutomaticASCIIHexEncoding types.Bool   `tfsdk:"client_id_use_automatic_ascii_hex_encoding"`
			Delimiter                            types.String `tfsdk:"delimiter"`
			DomainName                           types.String `tfsdk:"domain_name"`
			InterfaceDescription                 types.String `tfsdk:"interface_description"`
			InterfaceName                        types.Bool   `tfsdk:"interface_name"`
			MACAddress                           types.Bool   `tfsdk:"mac_address"`
			Option60                             types.Bool   `tfsdk:"option_60"`
			Option82                             types.Bool   `tfsdk:"option_82"`
			Option82CircuitID                    types.Bool   `tfsdk:"option_82_circuit_id"`
			Option82RemoteID                     types.Bool   `tfsdk:"option_82_remote_id"`
			RelayAgentInterfaceID                types.Bool   `tfsdk:"relay_agent_interface_id"`
			RelayAgentRemoteID                   types.Bool   `tfsdk:"relay_agent_remote_id"`
			RelayAgentSubscriberID               types.Bool   `tfsdk:"relay_agent_subscriber_id"`
			RoutingInstanceName                  types.Bool   `tfsdk:"routing_instance_name"`
			UserPrefix                           types.String `tfsdk:"user_prefix"`
			VlanTags                             types.Bool   `tfsdk:"vlan_tags"`
		} `tfsdk:"authentication_username_include"`
		BulkLeasequery []struct {
			Attempts         types.Int64 `tfsdk:"attempts"`
			Timeout          types.Int64 `tfsdk:"timeout"`
			TriggerAutomatic types.Bool  `tfsdk:"trigger_automatic"`
		} `tfsdk:"bulk_leasequery"`
		LeaseTimeValidation []struct {
			LeaseTimeThreshold  types.Int64 `tfsdk:"lease_time_threshold"`
			ViolationActionDrop types.Bool  `tfsdk:"violation_action_drop"`
		} `tfsdk:"lease_time_validation"`
		Leasequery []struct {
			Attempts types.Int64 `tfsdk:"attempts"`
			Timeout  types.Int64 `tfsdk:"timeout"`
		} `tfsdk:"leasequery"`
		LivenessDetectionMethodBfd []struct {
			DetectionTimeThreshold    types.Int64  `tfsdk:"detection_time_threshold"`
			HolddownInterval          types.Int64  `tfsdk:"holddown_interval"`
			MinimumInterval           types.Int64  `tfsdk:"minimum_interval"`
			MinimumReceiveInterval    types.Int64  `tfsdk:"minimum_receive_interval"`
			Multiplier                types.Int64  `tfsdk:"multiplier"`
			NoAdaptation              types.Bool   `tfsdk:"no_adaptation"`
			SessionMode               types.String `tfsdk:"session_mode"`
			TransmitIntervalMinimum   types.Int64  `tfsdk:"transmit_interval_minimum"`
			TransmitIntervalThreshold types.Int64  `tfsdk:"transmit_interval_threshold"`
			Version                   types.String `tfsdk:"version"`
		} `tfsdk:"liveness_detection_method_bfd"`
		LivenessDetectionMethodLayer2 []struct {
			MaxConsecutiveRetries types.Int64 `tfsdk:"max_consecutive_retries"`
			TransmitInterval      types.Int64 `tfsdk:"transmit_interval"`
		} `tfsdk:"liveness_detection_method_layer2"`
		OverridesV4 []struct {
			AllowNoEndOption             types.Bool   `tfsdk:"allow_no_end_option"`
			AllowSnoopedClients          types.Bool   `tfsdk:"allow_snooped_clients"`
			NoAllowSnoopedClients        types.Bool   `tfsdk:"no_allow_snooped_clients"`
			AlwaysWriteGiaddr            types.Bool   `tfsdk:"always_write_giaddr"`
			AlwaysWriteOption82          types.Bool   `tfsdk:"always_write_option_82"`
			AsymmetricLeaseTime          types.Int64  `tfsdk:"asymmetric_lease_time"`
			BootpSupport                 types.Bool   `tfsdk:"bootp_support"`
			ClientDiscoverMatch          types.String `tfsdk:"client_discover_match"`
			DelayAuthentication          types.Bool   `tfsdk:"delay_authentication"`
			DeleteBindingOnRenegotiation types.Bool   `tfsdk:"delete_binding_on_renegotiation"`
			DisableRelay                 types.Bool   `tfsdk:"disable_relay"`
			DualStack                    types.String `tfsdk:"dual_stack"`
			InterfaceClientLimit         types.Int64  `tfsdk:"interface_client_limit"`
			Layer2UnicastReplies         types.Bool   `tfsdk:"layer2_unicast_replies"`
			NoBindOnRequest              types.Bool   `tfsdk:"no_bind_on_request"`
			NoUnicastReplies             types.Bool   `tfsdk:"no_unicast_replies"`
			ProxyMode                    types.Bool   `tfsdk:"proxy_mode"`
			RelaySource                  types.String `tfsdk:"relay_source"`
			ReplaceIPSourceWithGiaddr    types.Bool   `tfsdk:"replace_ip_source_with_giaddr"`
			SendReleaseOnDelete          types.Bool   `tfsdk:"send_release_on_delete"`
			TrustOption82                types.Bool   `tfsdk:"trust_option_82"`
			UserDefinedOption82          types.String `tfsdk:"user_defined_option_82"`
		} `tfsdk:"overrides_v4"`
		OverridesV6 []struct {
			AllowSnoopedClients                     types.Bool   `tfsdk:"allow_snooped_clients"`
			NoAllowSnoopedClients                   types.Bool   `tfsdk:"no_allow_snooped_clients"`
			AlwaysProcessOptionRequestOption        types.Bool   `tfsdk:"always_process_option_request_option"`
			AsymmetricLeaseTime                     types.Int64  `tfsdk:"asymmetric_lease_time"`
			AsymmetricPrefixLeaseTime               types.Int64  `tfsdk:"asymmetric_prefix_lease_time"`
			ClientNegotiationMatchIncomingInterface types.Bool   `tfsdk:"client_negotiation_match_incoming_interface"`
			DelayAuthentication                     types.Bool   `tfsdk:"delay_authentication"`
			DeleteBindingOnRenegotiation            types.Bool   `tfsdk:"delete_binding_on_renegotiation"`
			DualStack                               types.String `tfsdk:"dual_stack"`
			InterfaceClientLimit                    types.Int64  `tfsdk:"interface_client_limit"`
			NoBindOnRequest                         types.Bool   `tfsdk:"no_bind_on_request"`
			RelaySource                             types.String `tfsdk:"relay_source"`
			SendReleaseOnDelete                     types.Bool   `tfsdk:"send_release_on_delete"`
		} `tfsdk:"overrides_v6"`
		RelayAgentInterfaceID []struct {
			KeepIncomingID            types.Bool   `tfsdk:"keep_incoming_id"`
			KeepIncomingIDStrict      types.Bool   `tfsdk:"keep_incoming_id_strict"`
			IncludeIrbAndL2           types.Bool   `tfsdk:"include_irb_and_l2"`
			NoVlanInterfaceName       types.Bool   `tfsdk:"no_vlan_interface_name"`
			PrefixHostName            types.Bool   `tfsdk:"prefix_host_name"`
			PrefixRoutingInstanceName types.Bool   `tfsdk:"prefix_routing_instance_name"`
			UseInterfaceDescription   types.String `tfsdk:"use_interface_description"`
			UseOption82               types.Bool   `tfsdk:"use_option_82"`
			UseOption82Strict         types.Bool   `tfsdk:"use_option_82_strict"`
			UseVlanID                 types.Bool   `tfsdk:"use_vlan_id"`
		} `tfsdk:"relay_agent_interface_id"`
		RelayAgentRemoteID []struct {
			KeepIncomingID            types.Bool   `tfsdk:"keep_incoming_id"`
			IncludeIrbAndL2           types.Bool   `tfsdk:"include_irb_and_l2"`
			NoVlanInterfaceName       types.Bool   `tfsdk:"no_vlan_interface_name"`
			PrefixHostName            types.Bool   `tfsdk:"prefix_host_name"`
			PrefixRoutingInstanceName types.Bool   `tfsdk:"prefix_routing_instance_name"`
			UseInterfaceDescription   types.String `tfsdk:"use_interface_description"`
			UseOption82               types.Bool   `tfsdk:"use_option_82"`
			UseOption82Strict         types.Bool   `tfsdk:"use_option_82_strict"`
			UseVlanID                 types.Bool   `tfsdk:"use_vlan_id"`
		} `tfsdk:"relay_agent_remote_id"`
		RelayOption []struct {
			OptionOrder           []types.String                                                        `tfsdk:"option_order"`
			Option15              []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN              `tfsdk:"option_15"`
			Option15DefaultAction []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_15_default_action"`
			Option16              []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN              `tfsdk:"option_16"`
			Option16DefaultAction []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_16_default_action"`
			Option60              []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN              `tfsdk:"option_60"`
			Option60DefaultAction []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_60_default_action"`
			Option77              []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionN              `tfsdk:"option_77"`
			Option77DefaultAction []forwardingoptionsDhcprelayBlockRelayOptionBlockOptionNDefaultAction `tfsdk:"option_77_default_action"`
		} `tfsdk:"relay_option"`
		RelayOption82 []struct {
			ExcludeRelayAgentIdentifier types.Bool `tfsdk:"exclude_relay_agent_identifier"`
			LinkSelection               types.Bool `tfsdk:"link_selection"`
			ServerIDOverride            types.Bool `tfsdk:"server_id_override"`
			VendorSpecificHostName      types.Bool `tfsdk:"vendor_specific_host_name"`
			VendorSpecificLocation      types.Bool `tfsdk:"vendor_specific_location"`
			CircuitID                   []struct {
				IncludeIrbAndL2           types.Bool   `tfsdk:"include_irb_and_l2"`
				KeepIncomingCircuitID     types.Bool   `tfsdk:"keep_incoming_circuit_id"`
				NoVlanInterfaceName       types.Bool   `tfsdk:"no_vlan_interface_name"`
				PrefixHostName            types.Bool   `tfsdk:"prefix_host_name"`
				PrefixRoutingInstanceName types.Bool   `tfsdk:"prefix_routing_instance_name"`
				UseInterfaceDescription   types.String `tfsdk:"use_interface_description"`
				UseVlanID                 types.Bool   `tfsdk:"use_vlan_id"`
				UserDefined               types.Bool   `tfsdk:"user_defined"`
				VlanIDOnly                types.Bool   `tfsdk:"vlan_id_only"`
			} `tfsdk:"circuit_id"`
			RemoteID []struct {
				HostnameOnly              types.Bool   `tfsdk:"hostname_only"`
				IncludeIrbAndL2           types.Bool   `tfsdk:"include_irb_and_l2"`
				KeepIncomingRemoteID      types.Bool   `tfsdk:"keep_incoming_remote_id"`
				NoVlanInterfaceName       types.Bool   `tfsdk:"no_vlan_interface_name"`
				PrefixHostName            types.Bool   `tfsdk:"prefix_host_name"`
				PrefixRoutingInstanceName types.Bool   `tfsdk:"prefix_routing_instance_name"`
				UseInterfaceDescription   types.String `tfsdk:"use_interface_description"`
				UseString                 types.String `tfsdk:"use_string"`
				UseVlanID                 types.Bool   `tfsdk:"use_vlan_id"`
			} `tfsdk:"remote_id"`
		} `tfsdk:"relay_option_82"`
		ServerMatchAddress []struct {
			Address types.String `tfsdk:"address"`
			Action  types.String `tfsdk:"action"`
		} `tfsdk:"server_match_address"`
		ServerMatchDuid []struct {
			Compare   types.String `tfsdk:"compare"`
			ValueType types.String `tfsdk:"value_type"`
			Value     types.String `tfsdk:"value"`
			Action    types.String `tfsdk:"action"`
		} `tfsdk:"server_match_duid"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 forwardingoptionsDhcprelayData
	dataV1.ID = dataV0.ID
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.Version = dataV0.Version
	dataV1.AccessProfile = dataV0.AccessProfile
	dataV1.ActiveServerGroup = dataV0.ActiveServerGroup
	dataV1.ActiveServerGroupAllowServerChange = dataV0.ActiveServerGroupAllowServerChange
	dataV1.ARPInspection = dataV0.ARPInspection
	dataV1.AuthenticationPassword = dataV0.AuthenticationPassword
	dataV1.ClientResponseTTL = dataV0.ClientResponseTTL
	dataV1.DuplicateClientsInSubnet = dataV0.DuplicateClientsInSubnet
	dataV1.DuplicateClientsIncomingInterface = dataV0.DuplicateClientsIncomingInterface
	dataV1.DynamicProfile = dataV0.DynamicProfile
	dataV1.DynamicProfileAggregateClients = dataV0.DynamicProfileAggregateClients
	dataV1.DynamicProfileAggregateClientsAction = dataV0.DynamicProfileAggregateClientsAction
	dataV1.DynamicProfileUsePrimary = dataV0.DynamicProfileUsePrimary
	dataV1.ExcludeRelayAgentIdentifier = dataV0.ExcludeRelayAgentIdentifier
	dataV1.ForwardOnly = dataV0.ForwardOnly
	dataV1.ForwardOnlyReplies = dataV0.ForwardOnlyReplies
	dataV1.ForwardOnlyRoutingInstance = dataV0.ForwardOnlyRoutingInstance
	dataV1.ForwardSnoopedClients = dataV0.ForwardSnoopedClients
	dataV1.LivenessDetectionFailureAction = dataV0.LivenessDetectionFailureAction
	dataV1.MaximumHopCount = dataV0.MaximumHopCount
	dataV1.MinimumWaitTime = dataV0.MinimumWaitTime
	dataV1.NoSnoop = dataV0.NoSnoop
	dataV1.PersistentStorageAutomatic = dataV0.PersistentStorageAutomatic
	dataV1.RelayAgentOption79 = dataV0.RelayAgentOption79
	dataV1.RemoteIDMismatchDisconnect = dataV0.RemoteIDMismatchDisconnect
	dataV1.RouteSuppressionAccess = dataV0.RouteSuppressionAccess
	dataV1.RouteSuppressionAccessInternal = dataV0.RouteSuppressionAccessInternal
	dataV1.RouteSuppressionDestination = dataV0.RouteSuppressionDestination
	dataV1.ServerMatchDefaultAction = dataV0.ServerMatchDefaultAction
	dataV1.ServerResponseTime = dataV0.ServerResponseTime
	dataV1.ServiceProfile = dataV0.ServiceProfile
	dataV1.ShortCycleProtectionLockoutMaxTime = dataV0.ShortCycleProtectionLockoutMaxTime
	dataV1.ShortCycleProtectionLockoutMinTime = dataV0.ShortCycleProtectionLockoutMinTime
	dataV1.SourceIPChange = dataV0.SourceIPChange
	dataV1.VendorSpecificInformationHostName = dataV0.VendorSpecificInformationHostName
	dataV1.VendorSpecificInformationLocation = dataV0.VendorSpecificInformationLocation
	if len(dataV0.ActiveLeasequery) > 0 {
		dataV1.ActiveLeasequery = &forwardingoptionsDhcprelayBlockActiveLeasequery{
			IdleTimeout:      dataV0.ActiveLeasequery[0].IdleTimeout,
			PeerAddress:      dataV0.ActiveLeasequery[0].PeerAddress,
			Timeout:          dataV0.ActiveLeasequery[0].Timeout,
			TopologyDiscover: dataV0.ActiveLeasequery[0].TopologyDiscover,
		}
	}
	if len(dataV0.AuthenticationUsernameInclude) > 0 {
		dataV1.AuthenticationUsernameInclude = &forwardingoptionsDhcprelayBlockAuthenticationUsernameInclude{
			CircuitType:                          dataV0.AuthenticationUsernameInclude[0].CircuitType,
			ClientID:                             dataV0.AuthenticationUsernameInclude[0].ClientID,
			ClientIDExcludeHeaders:               dataV0.AuthenticationUsernameInclude[0].ClientIDExcludeHeaders,
			ClientIDUseAutomaticASCIIHexEncoding: dataV0.AuthenticationUsernameInclude[0].ClientIDUseAutomaticASCIIHexEncoding,
			Delimiter:                            dataV0.AuthenticationUsernameInclude[0].Delimiter,
			DomainName:                           dataV0.AuthenticationUsernameInclude[0].DomainName,
			InterfaceDescription:                 dataV0.AuthenticationUsernameInclude[0].InterfaceDescription,
			InterfaceName:                        dataV0.AuthenticationUsernameInclude[0].InterfaceName,
			MACAddress:                           dataV0.AuthenticationUsernameInclude[0].MACAddress,
			Option60:                             dataV0.AuthenticationUsernameInclude[0].Option60,
			Option82:                             dataV0.AuthenticationUsernameInclude[0].Option82,
			Option82CircuitID:                    dataV0.AuthenticationUsernameInclude[0].Option82CircuitID,
			Option82RemoteID:                     dataV0.AuthenticationUsernameInclude[0].Option82RemoteID,
			RelayAgentInterfaceID:                dataV0.AuthenticationUsernameInclude[0].RelayAgentInterfaceID,
			RelayAgentRemoteID:                   dataV0.AuthenticationUsernameInclude[0].RelayAgentRemoteID,
			RelayAgentSubscriberID:               dataV0.AuthenticationUsernameInclude[0].RelayAgentSubscriberID,
			RoutingInstanceName:                  dataV0.AuthenticationUsernameInclude[0].RoutingInstanceName,
			UserPrefix:                           dataV0.AuthenticationUsernameInclude[0].UserPrefix,
			VlanTags:                             dataV0.AuthenticationUsernameInclude[0].VlanTags,
		}
	}
	if len(dataV0.BulkLeasequery) > 0 {
		dataV1.BulkLeasequery = &forwardingoptionsDhcprelayBlockBulkLeasequery{
			Attempts:         dataV0.BulkLeasequery[0].Attempts,
			Timeout:          dataV0.BulkLeasequery[0].Timeout,
			TriggerAutomatic: dataV0.BulkLeasequery[0].TriggerAutomatic,
		}
	}
	if len(dataV0.LeaseTimeValidation) > 0 {
		dataV1.LeaseTimeValidation = &forwardingoptionsDhcprelayBlockLeaseTimeValidation{
			LeaseTimeThreshold:  dataV0.LeaseTimeValidation[0].LeaseTimeThreshold,
			ViolationActionDrop: dataV0.LeaseTimeValidation[0].ViolationActionDrop,
		}
	}
	if len(dataV0.Leasequery) > 0 {
		dataV1.Leasequery = &forwardingoptionsDhcprelayBlockLeasequery{
			Attempts: dataV0.Leasequery[0].Attempts,
			Timeout:  dataV0.Leasequery[0].Timeout,
		}
	}
	if len(dataV0.LivenessDetectionMethodBfd) > 0 {
		dataV1.LivenessDetectionMethodBfd = &forwardingoptionsDhcprelayBlockLivenessDetectionMethodBfd{
			DetectionTimeThreshold:    dataV0.LivenessDetectionMethodBfd[0].DetectionTimeThreshold,
			HolddownInterval:          dataV0.LivenessDetectionMethodBfd[0].HolddownInterval,
			MinimumInterval:           dataV0.LivenessDetectionMethodBfd[0].MinimumInterval,
			MinimumReceiveInterval:    dataV0.LivenessDetectionMethodBfd[0].MinimumReceiveInterval,
			Multiplier:                dataV0.LivenessDetectionMethodBfd[0].Multiplier,
			NoAdaptation:              dataV0.LivenessDetectionMethodBfd[0].NoAdaptation,
			SessionMode:               dataV0.LivenessDetectionMethodBfd[0].SessionMode,
			TransmitIntervalMinimum:   dataV0.LivenessDetectionMethodBfd[0].TransmitIntervalMinimum,
			TransmitIntervalThreshold: dataV0.LivenessDetectionMethodBfd[0].TransmitIntervalThreshold,
			Version:                   dataV0.LivenessDetectionMethodBfd[0].Version,
		}
	}
	if len(dataV0.LivenessDetectionMethodLayer2) > 0 {
		dataV1.LivenessDetectionMethodLayer2 = &forwardingoptionsDhcprelayBlockLivenessDetectionMethodLayer2{
			MaxConsecutiveRetries: dataV0.LivenessDetectionMethodLayer2[0].MaxConsecutiveRetries,
			TransmitInterval:      dataV0.LivenessDetectionMethodLayer2[0].TransmitInterval,
		}
	}
	if len(dataV0.OverridesV4) > 0 {
		dataV1.OverridesV4 = &forwardingoptionsDhcprelayBlockOverridesV4{
			AllowNoEndOption:             dataV0.OverridesV4[0].AllowNoEndOption,
			AllowSnoopedClients:          dataV0.OverridesV4[0].AllowSnoopedClients,
			NoAllowSnoopedClients:        dataV0.OverridesV4[0].NoAllowSnoopedClients,
			AlwaysWriteGiaddr:            dataV0.OverridesV4[0].AlwaysWriteGiaddr,
			AlwaysWriteOption82:          dataV0.OverridesV4[0].AlwaysWriteOption82,
			AsymmetricLeaseTime:          dataV0.OverridesV4[0].AsymmetricLeaseTime,
			BootpSupport:                 dataV0.OverridesV4[0].BootpSupport,
			ClientDiscoverMatch:          dataV0.OverridesV4[0].ClientDiscoverMatch,
			DelayAuthentication:          dataV0.OverridesV4[0].DelayAuthentication,
			DeleteBindingOnRenegotiation: dataV0.OverridesV4[0].DeleteBindingOnRenegotiation,
			DisableRelay:                 dataV0.OverridesV4[0].DisableRelay,
			DualStack:                    dataV0.OverridesV4[0].DualStack,
			InterfaceClientLimit:         dataV0.OverridesV4[0].InterfaceClientLimit,
			Layer2UnicastReplies:         dataV0.OverridesV4[0].Layer2UnicastReplies,
			NoBindOnRequest:              dataV0.OverridesV4[0].NoBindOnRequest,
			NoUnicastReplies:             dataV0.OverridesV4[0].NoUnicastReplies,
			ProxyMode:                    dataV0.OverridesV4[0].ProxyMode,
			RelaySource:                  dataV0.OverridesV4[0].RelaySource,
			ReplaceIPSourceWithGiaddr:    dataV0.OverridesV4[0].ReplaceIPSourceWithGiaddr,
			SendReleaseOnDelete:          dataV0.OverridesV4[0].SendReleaseOnDelete,
			TrustOption82:                dataV0.OverridesV4[0].TrustOption82,
			UserDefinedOption82:          dataV0.OverridesV4[0].UserDefinedOption82,
		}
	}
	if len(dataV0.OverridesV6) > 0 {
		dataV1.OverridesV6 = &forwardingoptionsDhcprelayBlockOverridesV6{
			AllowSnoopedClients:                     dataV0.OverridesV6[0].AllowSnoopedClients,
			NoAllowSnoopedClients:                   dataV0.OverridesV6[0].NoAllowSnoopedClients,
			AlwaysProcessOptionRequestOption:        dataV0.OverridesV6[0].AlwaysProcessOptionRequestOption,
			AsymmetricLeaseTime:                     dataV0.OverridesV6[0].AsymmetricLeaseTime,
			AsymmetricPrefixLeaseTime:               dataV0.OverridesV6[0].AsymmetricPrefixLeaseTime,
			ClientNegotiationMatchIncomingInterface: dataV0.OverridesV6[0].ClientNegotiationMatchIncomingInterface,
			DelayAuthentication:                     dataV0.OverridesV6[0].DelayAuthentication,
			DeleteBindingOnRenegotiation:            dataV0.OverridesV6[0].DeleteBindingOnRenegotiation,
			DualStack:                               dataV0.OverridesV6[0].DualStack,
			InterfaceClientLimit:                    dataV0.OverridesV6[0].InterfaceClientLimit,
			NoBindOnRequest:                         dataV0.OverridesV6[0].NoBindOnRequest,
			RelaySource:                             dataV0.OverridesV6[0].RelaySource,
			SendReleaseOnDelete:                     dataV0.OverridesV6[0].SendReleaseOnDelete,
		}
	}
	if len(dataV0.RelayAgentInterfaceID) > 0 {
		dataV1.RelayAgentInterfaceID = &forwardingoptionsDhcprelayBlockRelayAgentInterfaceID{
			KeepIncomingID:       dataV0.RelayAgentInterfaceID[0].KeepIncomingID,
			KeepIncomingIDStrict: dataV0.RelayAgentInterfaceID[0].KeepIncomingIDStrict,
		}
		dataV1.RelayAgentInterfaceID.IncludeIrbAndL2 = dataV0.RelayAgentInterfaceID[0].IncludeIrbAndL2
		dataV1.RelayAgentInterfaceID.NoVlanInterfaceName = dataV0.RelayAgentInterfaceID[0].NoVlanInterfaceName
		dataV1.RelayAgentInterfaceID.PrefixHostName = dataV0.RelayAgentInterfaceID[0].PrefixHostName
		dataV1.RelayAgentInterfaceID.PrefixRoutingInstanceName = dataV0.RelayAgentInterfaceID[0].PrefixRoutingInstanceName
		dataV1.RelayAgentInterfaceID.UseInterfaceDescription = dataV0.RelayAgentInterfaceID[0].UseInterfaceDescription
		dataV1.RelayAgentInterfaceID.UseOption82 = dataV0.RelayAgentInterfaceID[0].UseOption82
		dataV1.RelayAgentInterfaceID.UseOption82Strict = dataV0.RelayAgentInterfaceID[0].UseOption82Strict
		dataV1.RelayAgentInterfaceID.UseVlanID = dataV0.RelayAgentInterfaceID[0].UseVlanID
	}
	if len(dataV0.RelayAgentRemoteID) > 0 {
		dataV1.RelayAgentRemoteID = &forwardingoptionsDhcprelayBlockRelayAgentRemoteID{
			KeepIncomingID: dataV0.RelayAgentRemoteID[0].KeepIncomingID,
		}
		dataV1.RelayAgentRemoteID.IncludeIrbAndL2 = dataV0.RelayAgentRemoteID[0].IncludeIrbAndL2
		dataV1.RelayAgentRemoteID.NoVlanInterfaceName = dataV0.RelayAgentRemoteID[0].NoVlanInterfaceName
		dataV1.RelayAgentRemoteID.PrefixHostName = dataV0.RelayAgentRemoteID[0].PrefixHostName
		dataV1.RelayAgentRemoteID.PrefixRoutingInstanceName = dataV0.RelayAgentRemoteID[0].PrefixRoutingInstanceName
		dataV1.RelayAgentRemoteID.UseInterfaceDescription = dataV0.RelayAgentRemoteID[0].UseInterfaceDescription
		dataV1.RelayAgentRemoteID.UseOption82 = dataV0.RelayAgentRemoteID[0].UseOption82
		dataV1.RelayAgentRemoteID.UseOption82Strict = dataV0.RelayAgentRemoteID[0].UseOption82Strict
		dataV1.RelayAgentRemoteID.UseVlanID = dataV0.RelayAgentRemoteID[0].UseVlanID
	}
	if len(dataV0.RelayOption) > 0 {
		dataV1.RelayOption = &forwardingoptionsDhcprelayBlockRelayOption{
			OptionOrder: dataV0.RelayOption[0].OptionOrder,
			Option15:    dataV0.RelayOption[0].Option15,
			Option16:    dataV0.RelayOption[0].Option16,
			Option60:    dataV0.RelayOption[0].Option60,
			Option77:    dataV0.RelayOption[0].Option77,
		}
		if len(dataV0.RelayOption[0].Option15DefaultAction) > 0 {
			dataV1.RelayOption.Option15DefaultAction = &dataV0.RelayOption[0].Option15DefaultAction[0]
		}
		if len(dataV0.RelayOption[0].Option16DefaultAction) > 0 {
			dataV1.RelayOption.Option16DefaultAction = &dataV0.RelayOption[0].Option16DefaultAction[0]
		}
		if len(dataV0.RelayOption[0].Option60DefaultAction) > 0 {
			dataV1.RelayOption.Option60DefaultAction = &dataV0.RelayOption[0].Option60DefaultAction[0]
		}
		if len(dataV0.RelayOption[0].Option77DefaultAction) > 0 {
			dataV1.RelayOption.Option77DefaultAction = &dataV0.RelayOption[0].Option77DefaultAction[0]
		}
	}
	if len(dataV0.RelayOption82) > 0 {
		dataV1.RelayOption82 = &forwardingoptionsDhcprelayBlockRelayOption82{
			ExcludeRelayAgentIdentifier: dataV0.RelayOption82[0].ExcludeRelayAgentIdentifier,
			LinkSelection:               dataV0.RelayOption82[0].LinkSelection,
			ServerIDOverride:            dataV0.RelayOption82[0].ServerIDOverride,
			VendorSpecificHostName:      dataV0.RelayOption82[0].VendorSpecificHostName,
			VendorSpecificLocation:      dataV0.RelayOption82[0].VendorSpecificLocation,
		}
		if len(dataV0.RelayOption82[0].CircuitID) > 0 {
			dataV1.RelayOption82.CircuitID = &forwardingoptionsDhcprelayBlockRelayOption82BlockCircuitID{
				IncludeIrbAndL2:           dataV0.RelayOption82[0].CircuitID[0].IncludeIrbAndL2,
				KeepIncomingCircuitID:     dataV0.RelayOption82[0].CircuitID[0].KeepIncomingCircuitID,
				NoVlanInterfaceName:       dataV0.RelayOption82[0].CircuitID[0].NoVlanInterfaceName,
				PrefixHostName:            dataV0.RelayOption82[0].CircuitID[0].PrefixHostName,
				PrefixRoutingInstanceName: dataV0.RelayOption82[0].CircuitID[0].PrefixRoutingInstanceName,
				UseInterfaceDescription:   dataV0.RelayOption82[0].CircuitID[0].UseInterfaceDescription,
				UseVlanID:                 dataV0.RelayOption82[0].CircuitID[0].UseVlanID,
				UserDefined:               dataV0.RelayOption82[0].CircuitID[0].UserDefined,
				VlanIDOnly:                dataV0.RelayOption82[0].CircuitID[0].VlanIDOnly,
			}
		}
		if len(dataV0.RelayOption82[0].RemoteID) > 0 {
			dataV1.RelayOption82.RemoteID = &forwardingoptionsDhcprelayBlockRelayOption82BlockRemoteID{
				HostnameOnly:              dataV0.RelayOption82[0].RemoteID[0].HostnameOnly,
				IncludeIrbAndL2:           dataV0.RelayOption82[0].RemoteID[0].IncludeIrbAndL2,
				KeepIncomingRemoteID:      dataV0.RelayOption82[0].RemoteID[0].KeepIncomingRemoteID,
				NoVlanInterfaceName:       dataV0.RelayOption82[0].RemoteID[0].NoVlanInterfaceName,
				PrefixHostName:            dataV0.RelayOption82[0].RemoteID[0].PrefixHostName,
				PrefixRoutingInstanceName: dataV0.RelayOption82[0].RemoteID[0].PrefixRoutingInstanceName,
				UseInterfaceDescription:   dataV0.RelayOption82[0].RemoteID[0].UseInterfaceDescription,
				UseString:                 dataV0.RelayOption82[0].RemoteID[0].UseString,
				UseVlanID:                 dataV0.RelayOption82[0].RemoteID[0].UseVlanID,
			}
		}
	}
	for _, block := range dataV0.ServerMatchAddress {
		dataV1.ServerMatchAddress = append(dataV1.ServerMatchAddress,
			forwardingoptionsDhcprelayBlockServerMatchAddress{
				Address: block.Address,
				Action:  block.Action,
			})
	}
	for _, block := range dataV0.ServerMatchDuid {
		dataV1.ServerMatchDuid = append(dataV1.ServerMatchDuid,
			forwardingoptionsDhcprelayBlockServerMatchDuid{
				Compare:   block.Compare,
				ValueType: block.ValueType,
				Value:     block.Value,
				Action:    block.Action,
			})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
