package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *systemServicesDhcpLocalserverGroup) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"version": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"access_profile": schema.StringAttribute{
						Optional: true,
					},
					"authentication_password": schema.StringAttribute{
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
					"liveness_detection_failure_action": schema.StringAttribute{
						Optional: true,
					},
					"reauthenticate_lease_renewal": schema.BoolAttribute{
						Optional: true,
					},
					"reauthenticate_remote_id_mismatch": schema.BoolAttribute{
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
					"service_profile": schema.StringAttribute{
						Optional: true,
					},
					"short_cycle_protection_lockout_max_time": schema.Int64Attribute{
						Optional: true,
					},
					"short_cycle_protection_lockout_min_time": schema.Int64Attribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
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
					"interface": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"access_profile": schema.StringAttribute{
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
								"exclude": schema.BoolAttribute{
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
								"trace": schema.BoolAttribute{
									Optional: true,
								},
								"upto": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"overrides_v4": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"allow_no_end_option": schema.BoolAttribute{
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
											"delay_offer_delay_time": schema.Int64Attribute{
												Optional: true,
											},
											"delete_binding_on_renegotiation": schema.BoolAttribute{
												Optional: true,
											},
											"dual_stack": schema.StringAttribute{
												Optional: true,
											},
											"include_option_82_forcerenew": schema.BoolAttribute{
												Optional: true,
											},
											"include_option_82_nak": schema.BoolAttribute{
												Optional: true,
											},
											"interface_client_limit": schema.Int64Attribute{
												Optional: true,
											},
											"process_inform": schema.BoolAttribute{
												Optional: true,
											},
											"process_inform_pool": schema.StringAttribute{
												Optional: true,
											},
											"protocol_attributes": schema.StringAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"delay_offer_based_on": schema.SetNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"option": schema.StringAttribute{
															Required: true,
														},
														"compare": schema.StringAttribute{
															Required: true,
														},
														"value_type": schema.StringAttribute{
															Required: true,
														},
														"value": schema.StringAttribute{
															Required: true,
														},
													},
												},
											},
										},
									},
								},
								"overrides_v6": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"always_add_option_dns_server": schema.BoolAttribute{
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
											"delay_advertise_delay_time": schema.Int64Attribute{
												Optional: true,
											},
											"delegated_pool": schema.StringAttribute{
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
											"multi_address_embedded_option_response": schema.BoolAttribute{
												Optional: true,
											},
											"process_inform": schema.BoolAttribute{
												Optional: true,
											},
											"process_inform_pool": schema.StringAttribute{
												Optional: true,
											},
											"protocol_attributes": schema.StringAttribute{
												Optional: true,
											},
											"rapid_commit": schema.BoolAttribute{
												Optional: true,
											},
											"top_level_status_code": schema.BoolAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"delay_advertise_based_on": schema.SetNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"option": schema.StringAttribute{
															Required: true,
														},
														"compare": schema.StringAttribute{
															Required: true,
														},
														"value_type": schema.StringAttribute{
															Required: true,
														},
														"value": schema.StringAttribute{
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
					"lease_time_validation": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"lease_time_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"violation_action": schema.StringAttribute{
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
								"asymmetric_lease_time": schema.Int64Attribute{
									Optional: true,
								},
								"bootp_support": schema.BoolAttribute{
									Optional: true,
								},
								"client_discover_match": schema.StringAttribute{
									Optional: true,
								},
								"delay_offer_delay_time": schema.Int64Attribute{
									Optional: true,
								},
								"delete_binding_on_renegotiation": schema.BoolAttribute{
									Optional: true,
								},
								"dual_stack": schema.StringAttribute{
									Optional: true,
								},
								"include_option_82_forcerenew": schema.BoolAttribute{
									Optional: true,
								},
								"include_option_82_nak": schema.BoolAttribute{
									Optional: true,
								},
								"interface_client_limit": schema.Int64Attribute{
									Optional: true,
								},
								"process_inform": schema.BoolAttribute{
									Optional: true,
								},
								"process_inform_pool": schema.StringAttribute{
									Optional: true,
								},
								"protocol_attributes": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"delay_offer_based_on": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"option": schema.StringAttribute{
												Required: true,
											},
											"compare": schema.StringAttribute{
												Required: true,
											},
											"value_type": schema.StringAttribute{
												Required: true,
											},
											"value": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
							},
						},
					},
					"overrides_v6": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"always_add_option_dns_server": schema.BoolAttribute{
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
								"delay_advertise_delay_time": schema.Int64Attribute{
									Optional: true,
								},
								"delegated_pool": schema.StringAttribute{
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
								"multi_address_embedded_option_response": schema.BoolAttribute{
									Optional: true,
								},
								"process_inform": schema.BoolAttribute{
									Optional: true,
								},
								"process_inform_pool": schema.StringAttribute{
									Optional: true,
								},
								"protocol_attributes": schema.StringAttribute{
									Optional: true,
								},
								"rapid_commit": schema.BoolAttribute{
									Optional: true,
								},
								"top_level_status_code": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"delay_advertise_based_on": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"option": schema.StringAttribute{
												Required: true,
											},
											"compare": schema.StringAttribute{
												Required: true,
											},
											"value_type": schema.StringAttribute{
												Required: true,
											},
											"value": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
							},
						},
					},
					"reconfigure": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"attempts": schema.Int64Attribute{
									Optional: true,
								},
								"clear_on_abort": schema.BoolAttribute{
									Optional: true,
								},
								"support_option_pd_exclude": schema.BoolAttribute{
									Optional: true,
								},
								"timeout": schema.Int64Attribute{
									Optional: true,
								},
								"token": schema.StringAttribute{
									Optional: true,
								},
								"trigger_radius_disconnect": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSystemServicesDhcpLocalserverGroupV0toV1,
		},
	}
}

func upgradeSystemServicesDhcpLocalserverGroupV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	//nolint:lll
	type modelV0 struct {
		ID                                   types.String `tfsdk:"id"`
		Name                                 types.String `tfsdk:"name"`
		RoutingInstance                      types.String `tfsdk:"routing_instance"`
		Version                              types.String `tfsdk:"version"`
		AccessProfile                        types.String `tfsdk:"access_profile"`
		AuthenticationPassword               types.String `tfsdk:"authentication_password"`
		DynamicProfile                       types.String `tfsdk:"dynamic_profile"`
		DynamicProfileAggregateClients       types.Bool   `tfsdk:"dynamic_profile_aggregate_clients"`
		DynamicProfileAggregateClientsAction types.String `tfsdk:"dynamic_profile_aggregate_clients_action"`
		DynamicProfileUsePrimary             types.String `tfsdk:"dynamic_profile_use_primary"`
		LivenessDetectionFailureAction       types.String `tfsdk:"liveness_detection_failure_action"`
		ReauthenticateLeaseRenewal           types.Bool   `tfsdk:"reauthenticate_lease_renewal"`
		ReauthenticateRemoteIDMismatch       types.Bool   `tfsdk:"reauthenticate_remote_id_mismatch"`
		RemoteIDMismatchDisconnect           types.Bool   `tfsdk:"remote_id_mismatch_disconnect"`
		RouteSuppressionAccess               types.Bool   `tfsdk:"route_suppression_access"`
		RouteSuppressionAccessInternal       types.Bool   `tfsdk:"route_suppression_access_internal"`
		RouteSuppressionDestination          types.Bool   `tfsdk:"route_suppression_destination"`
		ServiceProfile                       types.String `tfsdk:"service_profile"`
		ShortCycleProtectionLockoutMaxTime   types.Int64  `tfsdk:"short_cycle_protection_lockout_max_time"`
		ShortCycleProtectionLockoutMinTime   types.Int64  `tfsdk:"short_cycle_protection_lockout_min_time"`
		AuthenticationUsernameInclude        []struct {
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
		Interface []struct {
			Name                                 types.String `tfsdk:"name"`
			AccessProfile                        types.String `tfsdk:"access_profile"`
			DynamicProfile                       types.String `tfsdk:"dynamic_profile"`
			DynamicProfileAggregateClients       types.Bool   `tfsdk:"dynamic_profile_aggregate_clients"`
			DynamicProfileAggregateClientsAction types.String `tfsdk:"dynamic_profile_aggregate_clients_action"`
			DynamicProfileUsePrimary             types.String `tfsdk:"dynamic_profile_use_primary"`
			Exclude                              types.Bool   `tfsdk:"exclude"`
			ServiceProfile                       types.String `tfsdk:"service_profile"`
			ShortCycleProtectionLockoutMaxTime   types.Int64  `tfsdk:"short_cycle_protection_lockout_max_time"`
			ShortCycleProtectionLockoutMinTime   types.Int64  `tfsdk:"short_cycle_protection_lockout_min_time"`
			Trace                                types.Bool   `tfsdk:"trace"`
			Upto                                 types.String `tfsdk:"upto"`
			OverridesV4                          []struct {
				AllowNoEndOption             types.Bool                                                          `tfsdk:"allow_no_end_option"`
				AsymmetricLeaseTime          types.Int64                                                         `tfsdk:"asymmetric_lease_time"`
				BootpSupport                 types.Bool                                                          `tfsdk:"bootp_support"`
				ClientDiscoverMatch          types.String                                                        `tfsdk:"client_discover_match"`
				DelayOfferDelayTime          types.Int64                                                         `tfsdk:"delay_offer_delay_time"`
				DeleteBindingOnRenegotiation types.Bool                                                          `tfsdk:"delete_binding_on_renegotiation"`
				DualStack                    types.String                                                        `tfsdk:"dual_stack"`
				IncludeOption82Forcerenew    types.Bool                                                          `tfsdk:"include_option_82_forcerenew"`
				IncludeOption82Nak           types.Bool                                                          `tfsdk:"include_option_82_nak"`
				InterfaceClientLimit         types.Int64                                                         `tfsdk:"interface_client_limit"`
				ProcessInform                types.Bool                                                          `tfsdk:"process_inform"`
				ProcessInformPool            types.String                                                        `tfsdk:"process_inform_pool"`
				ProtocolAttributes           types.String                                                        `tfsdk:"protocol_attributes"`
				DelayOfferBasedOn            []systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn `tfsdk:"delay_offer_based_on"`
			} `tfsdk:"overrides_v4"`
			OverridesV6 []struct {
				AlwaysAddOptionDNSServer                types.Bool                                                          `tfsdk:"always_add_option_dns_server"`
				AlwaysProcessOptionRequestOption        types.Bool                                                          `tfsdk:"always_process_option_request_option"`
				AsymmetricLeaseTime                     types.Int64                                                         `tfsdk:"asymmetric_lease_time"`
				AsymmetricPrefixLeaseTime               types.Int64                                                         `tfsdk:"asymmetric_prefix_lease_time"`
				ClientNegotiationMatchIncomingInterface types.Bool                                                          `tfsdk:"client_negotiation_match_incoming_interface"`
				DelayAdvertiseDelayTime                 types.Int64                                                         `tfsdk:"delay_advertise_delay_time"`
				DelegatedPool                           types.String                                                        `tfsdk:"delegated_pool"`
				DeleteBindingOnRenegotiation            types.Bool                                                          `tfsdk:"delete_binding_on_renegotiation"`
				DualStack                               types.String                                                        `tfsdk:"dual_stack"`
				InterfaceClientLimit                    types.Int64                                                         `tfsdk:"interface_client_limit"`
				MultiAddressEmbeddedOptionResponse      types.Bool                                                          `tfsdk:"multi_address_embedded_option_response"`
				ProcessInform                           types.Bool                                                          `tfsdk:"process_inform"`
				ProcessInformPool                       types.String                                                        `tfsdk:"process_inform_pool"`
				ProtocolAttributes                      types.String                                                        `tfsdk:"protocol_attributes"`
				RapidCommit                             types.Bool                                                          `tfsdk:"rapid_commit"`
				TopLevelStatusCode                      types.Bool                                                          `tfsdk:"top_level_status_code"`
				DelayAdvertiseBasedOn                   []systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn `tfsdk:"delay_advertise_based_on"`
			} `tfsdk:"overrides_v6"`
		} `tfsdk:"interface"`
		LeaseTimeValidation []struct {
			LeaseTimeThreshold types.Int64  `tfsdk:"lease_time_threshold"`
			ViolationAction    types.String `tfsdk:"violation_action"`
		} `tfsdk:"lease_time_validation"`
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
			AllowNoEndOption             types.Bool                                                          `tfsdk:"allow_no_end_option"`
			AsymmetricLeaseTime          types.Int64                                                         `tfsdk:"asymmetric_lease_time"`
			BootpSupport                 types.Bool                                                          `tfsdk:"bootp_support"`
			ClientDiscoverMatch          types.String                                                        `tfsdk:"client_discover_match"`
			DelayOfferDelayTime          types.Int64                                                         `tfsdk:"delay_offer_delay_time"`
			DeleteBindingOnRenegotiation types.Bool                                                          `tfsdk:"delete_binding_on_renegotiation"`
			DualStack                    types.String                                                        `tfsdk:"dual_stack"`
			IncludeOption82Forcerenew    types.Bool                                                          `tfsdk:"include_option_82_forcerenew"`
			IncludeOption82Nak           types.Bool                                                          `tfsdk:"include_option_82_nak"`
			InterfaceClientLimit         types.Int64                                                         `tfsdk:"interface_client_limit"`
			ProcessInform                types.Bool                                                          `tfsdk:"process_inform"`
			ProcessInformPool            types.String                                                        `tfsdk:"process_inform_pool"`
			ProtocolAttributes           types.String                                                        `tfsdk:"protocol_attributes"`
			DelayOfferBasedOn            []systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn `tfsdk:"delay_offer_based_on"`
		} `tfsdk:"overrides_v4"`
		OverridesV6 []struct {
			AlwaysAddOptionDNSServer                types.Bool                                                          `tfsdk:"always_add_option_dns_server"`
			AlwaysProcessOptionRequestOption        types.Bool                                                          `tfsdk:"always_process_option_request_option"`
			AsymmetricLeaseTime                     types.Int64                                                         `tfsdk:"asymmetric_lease_time"`
			AsymmetricPrefixLeaseTime               types.Int64                                                         `tfsdk:"asymmetric_prefix_lease_time"`
			ClientNegotiationMatchIncomingInterface types.Bool                                                          `tfsdk:"client_negotiation_match_incoming_interface"`
			DelayAdvertiseDelayTime                 types.Int64                                                         `tfsdk:"delay_advertise_delay_time"`
			DelegatedPool                           types.String                                                        `tfsdk:"delegated_pool"`
			DeleteBindingOnRenegotiation            types.Bool                                                          `tfsdk:"delete_binding_on_renegotiation"`
			DualStack                               types.String                                                        `tfsdk:"dual_stack"`
			InterfaceClientLimit                    types.Int64                                                         `tfsdk:"interface_client_limit"`
			MultiAddressEmbeddedOptionResponse      types.Bool                                                          `tfsdk:"multi_address_embedded_option_response"`
			ProcessInform                           types.Bool                                                          `tfsdk:"process_inform"`
			ProcessInformPool                       types.String                                                        `tfsdk:"process_inform_pool"`
			ProtocolAttributes                      types.String                                                        `tfsdk:"protocol_attributes"`
			RapidCommit                             types.Bool                                                          `tfsdk:"rapid_commit"`
			TopLevelStatusCode                      types.Bool                                                          `tfsdk:"top_level_status_code"`
			DelayAdvertiseBasedOn                   []systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn `tfsdk:"delay_advertise_based_on"`
		} `tfsdk:"overrides_v6"`
		Reconfigure []struct {
			Attempts                types.Int64  `tfsdk:"attempts"`
			ClearOnAbort            types.Bool   `tfsdk:"clear_on_abort"`
			SupportOptionPdExclude  types.Bool   `tfsdk:"support_option_pd_exclude"`
			Timeout                 types.Int64  `tfsdk:"timeout"`
			Token                   types.String `tfsdk:"token"`
			TriggerRadiusDisconnect types.Bool   `tfsdk:"trigger_radius_disconnect"`
		} `tfsdk:"reconfigure"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 systemServicesDhcpLocalserverGroupData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.Version = dataV0.Version
	dataV1.AccessProfile = dataV0.AccessProfile
	dataV1.AuthenticationPassword = dataV0.AuthenticationPassword
	dataV1.DynamicProfile = dataV0.DynamicProfile
	dataV1.DynamicProfileAggregateClients = dataV0.DynamicProfileAggregateClients
	dataV1.DynamicProfileAggregateClientsAction = dataV0.DynamicProfileAggregateClientsAction
	dataV1.DynamicProfileUsePrimary = dataV0.DynamicProfileUsePrimary
	dataV1.LivenessDetectionFailureAction = dataV0.LivenessDetectionFailureAction
	dataV1.ReauthenticateLeaseRenewal = dataV0.ReauthenticateLeaseRenewal
	dataV1.ReauthenticateRemoteIDMismatch = dataV0.ReauthenticateRemoteIDMismatch
	dataV1.RemoteIDMismatchDisconnect = dataV0.RemoteIDMismatchDisconnect
	dataV1.RouteSuppressionAccess = dataV0.RouteSuppressionAccess
	dataV1.RouteSuppressionAccessInternal = dataV0.RouteSuppressionAccessInternal
	dataV1.RouteSuppressionDestination = dataV0.RouteSuppressionDestination
	dataV1.ServiceProfile = dataV0.ServiceProfile
	dataV1.ShortCycleProtectionLockoutMaxTime = dataV0.ShortCycleProtectionLockoutMaxTime
	dataV1.ShortCycleProtectionLockoutMinTime = dataV0.ShortCycleProtectionLockoutMinTime
	if len(dataV0.AuthenticationUsernameInclude) > 0 {
		dataV1.AuthenticationUsernameInclude = &dhcpBlockAuthenticationUsernameInclude{
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
	for _, blockV0 := range dataV0.Interface {
		blockV1 := systemServicesDhcpLocalserverGroupBlockInterface{
			Name:                                 blockV0.Name,
			AccessProfile:                        blockV0.AccessProfile,
			DynamicProfile:                       blockV0.DynamicProfile,
			DynamicProfileAggregateClients:       blockV0.DynamicProfileAggregateClients,
			DynamicProfileAggregateClientsAction: blockV0.DynamicProfileAggregateClientsAction,
			DynamicProfileUsePrimary:             blockV0.DynamicProfileUsePrimary,
			Exclude:                              blockV0.Exclude,
			ServiceProfile:                       blockV0.ServiceProfile,
			ShortCycleProtectionLockoutMaxTime:   blockV0.ShortCycleProtectionLockoutMaxTime,
			ShortCycleProtectionLockoutMinTime:   blockV0.ShortCycleProtectionLockoutMinTime,
			Trace:                                blockV0.Trace,
			Upto:                                 blockV0.Upto,
		}
		if len(blockV0.OverridesV4) > 0 {
			blockV1.OverridesV4 = &systemServicesDhcpLocalserverGroupBlockOverridesV4{
				AllowNoEndOption:             blockV0.OverridesV4[0].AllowNoEndOption,
				AsymmetricLeaseTime:          blockV0.OverridesV4[0].AsymmetricLeaseTime,
				BootpSupport:                 blockV0.OverridesV4[0].BootpSupport,
				ClientDiscoverMatch:          blockV0.OverridesV4[0].ClientDiscoverMatch,
				DelayOfferDelayTime:          blockV0.OverridesV4[0].DelayOfferDelayTime,
				DeleteBindingOnRenegotiation: blockV0.OverridesV4[0].DeleteBindingOnRenegotiation,
				DualStack:                    blockV0.OverridesV4[0].DualStack,
				IncludeOption82Forcerenew:    blockV0.OverridesV4[0].IncludeOption82Forcerenew,
				IncludeOption82Nak:           blockV0.OverridesV4[0].IncludeOption82Nak,
				InterfaceClientLimit:         blockV0.OverridesV4[0].InterfaceClientLimit,
				ProcessInform:                blockV0.OverridesV4[0].ProcessInform,
				ProcessInformPool:            blockV0.OverridesV4[0].ProcessInformPool,
				ProtocolAttributes:           blockV0.OverridesV4[0].ProtocolAttributes,
				DelayOfferBasedOn:            blockV0.OverridesV4[0].DelayOfferBasedOn,
			}
		}
		if len(blockV0.OverridesV6) > 0 {
			blockV1.OverridesV6 = &systemServicesDhcpLocalserverGroupBlockOverridesV6{
				AlwaysAddOptionDNSServer:                blockV0.OverridesV6[0].AlwaysAddOptionDNSServer,
				AlwaysProcessOptionRequestOption:        blockV0.OverridesV6[0].AlwaysProcessOptionRequestOption,
				AsymmetricLeaseTime:                     blockV0.OverridesV6[0].AsymmetricLeaseTime,
				AsymmetricPrefixLeaseTime:               blockV0.OverridesV6[0].AsymmetricPrefixLeaseTime,
				ClientNegotiationMatchIncomingInterface: blockV0.OverridesV6[0].ClientNegotiationMatchIncomingInterface,
				DelayAdvertiseDelayTime:                 blockV0.OverridesV6[0].DelayAdvertiseDelayTime,
				DelegatedPool:                           blockV0.OverridesV6[0].DelegatedPool,
				DeleteBindingOnRenegotiation:            blockV0.OverridesV6[0].DeleteBindingOnRenegotiation,
				DualStack:                               blockV0.OverridesV6[0].DualStack,
				InterfaceClientLimit:                    blockV0.OverridesV6[0].InterfaceClientLimit,
				MultiAddressEmbeddedOptionResponse:      blockV0.OverridesV6[0].MultiAddressEmbeddedOptionResponse,
				ProcessInform:                           blockV0.OverridesV6[0].ProcessInform,
				ProcessInformPool:                       blockV0.OverridesV6[0].ProcessInformPool,
				ProtocolAttributes:                      blockV0.OverridesV6[0].ProtocolAttributes,
				RapidCommit:                             blockV0.OverridesV6[0].RapidCommit,
				TopLevelStatusCode:                      blockV0.OverridesV6[0].TopLevelStatusCode,
				DelayAdvertiseBasedOn:                   blockV0.OverridesV6[0].DelayAdvertiseBasedOn,
			}
		}
		dataV1.Interface = append(dataV1.Interface, blockV1)
	}
	if len(dataV0.LeaseTimeValidation) > 0 {
		dataV1.LeaseTimeValidation = &systemServicesDhcpLocalserverGroupBlockLeaseTimeValidation{
			LeaseTimeThreshold: dataV0.LeaseTimeValidation[0].LeaseTimeThreshold,
			ViolationAction:    dataV0.LeaseTimeValidation[0].ViolationAction,
		}
	}
	if len(dataV0.LivenessDetectionMethodBfd) > 0 {
		dataV1.LivenessDetectionMethodBfd = &dhcpBlockLivenessDetectionMethodBfd{
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
		dataV1.LivenessDetectionMethodLayer2 = &dhcpBlockLivenessDetectionMethodLayer2{
			MaxConsecutiveRetries: dataV0.LivenessDetectionMethodLayer2[0].MaxConsecutiveRetries,
			TransmitInterval:      dataV0.LivenessDetectionMethodLayer2[0].TransmitInterval,
		}
	}
	if len(dataV0.OverridesV4) > 0 {
		dataV1.OverridesV4 = &systemServicesDhcpLocalserverGroupBlockOverridesV4{
			AllowNoEndOption:             dataV0.OverridesV4[0].AllowNoEndOption,
			AsymmetricLeaseTime:          dataV0.OverridesV4[0].AsymmetricLeaseTime,
			BootpSupport:                 dataV0.OverridesV4[0].BootpSupport,
			ClientDiscoverMatch:          dataV0.OverridesV4[0].ClientDiscoverMatch,
			DelayOfferDelayTime:          dataV0.OverridesV4[0].DelayOfferDelayTime,
			DeleteBindingOnRenegotiation: dataV0.OverridesV4[0].DeleteBindingOnRenegotiation,
			DualStack:                    dataV0.OverridesV4[0].DualStack,
			IncludeOption82Forcerenew:    dataV0.OverridesV4[0].IncludeOption82Forcerenew,
			IncludeOption82Nak:           dataV0.OverridesV4[0].IncludeOption82Nak,
			InterfaceClientLimit:         dataV0.OverridesV4[0].InterfaceClientLimit,
			ProcessInform:                dataV0.OverridesV4[0].ProcessInform,
			ProcessInformPool:            dataV0.OverridesV4[0].ProcessInformPool,
			ProtocolAttributes:           dataV0.OverridesV4[0].ProtocolAttributes,
			DelayOfferBasedOn:            dataV0.OverridesV4[0].DelayOfferBasedOn,
		}
	}
	if len(dataV0.OverridesV6) > 0 {
		dataV1.OverridesV6 = &systemServicesDhcpLocalserverGroupBlockOverridesV6{
			AlwaysAddOptionDNSServer:                dataV0.OverridesV6[0].AlwaysAddOptionDNSServer,
			AlwaysProcessOptionRequestOption:        dataV0.OverridesV6[0].AlwaysProcessOptionRequestOption,
			AsymmetricLeaseTime:                     dataV0.OverridesV6[0].AsymmetricLeaseTime,
			AsymmetricPrefixLeaseTime:               dataV0.OverridesV6[0].AsymmetricPrefixLeaseTime,
			ClientNegotiationMatchIncomingInterface: dataV0.OverridesV6[0].ClientNegotiationMatchIncomingInterface,
			DelayAdvertiseDelayTime:                 dataV0.OverridesV6[0].DelayAdvertiseDelayTime,
			DelegatedPool:                           dataV0.OverridesV6[0].DelegatedPool,
			DeleteBindingOnRenegotiation:            dataV0.OverridesV6[0].DeleteBindingOnRenegotiation,
			DualStack:                               dataV0.OverridesV6[0].DualStack,
			InterfaceClientLimit:                    dataV0.OverridesV6[0].InterfaceClientLimit,
			MultiAddressEmbeddedOptionResponse:      dataV0.OverridesV6[0].MultiAddressEmbeddedOptionResponse,
			ProcessInform:                           dataV0.OverridesV6[0].ProcessInform,
			ProcessInformPool:                       dataV0.OverridesV6[0].ProcessInformPool,
			ProtocolAttributes:                      dataV0.OverridesV6[0].ProtocolAttributes,
			RapidCommit:                             dataV0.OverridesV6[0].RapidCommit,
			TopLevelStatusCode:                      dataV0.OverridesV6[0].TopLevelStatusCode,
			DelayAdvertiseBasedOn:                   dataV0.OverridesV6[0].DelayAdvertiseBasedOn,
		}
	}
	if len(dataV0.Reconfigure) > 0 {
		dataV1.Reconfigure = &systemServicesDhcpLocalserverGroupBlockReconfigure{
			Attempts:                dataV0.Reconfigure[0].Attempts,
			ClearOnAbort:            dataV0.Reconfigure[0].ClearOnAbort,
			SupportOptionPdExclude:  dataV0.Reconfigure[0].SupportOptionPdExclude,
			Timeout:                 dataV0.Reconfigure[0].Timeout,
			Token:                   dataV0.Reconfigure[0].Token,
			TriggerRadiusDisconnect: dataV0.Reconfigure[0].TriggerRadiusDisconnect,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
