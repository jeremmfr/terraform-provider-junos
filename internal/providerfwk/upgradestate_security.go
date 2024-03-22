package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *security) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"clean_on_destroy": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"alg": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"dns_disable": schema.BoolAttribute{
									Optional: true,
								},
								"ftp_disable": schema.BoolAttribute{
									Optional: true,
								},
								"h323_disable": schema.BoolAttribute{
									Optional: true,
								},
								"mgcp_disable": schema.BoolAttribute{
									Optional: true,
								},
								"msrpc_disable": schema.BoolAttribute{
									Optional: true,
								},
								"pptp_disable": schema.BoolAttribute{
									Optional: true,
								},
								"rsh_disable": schema.BoolAttribute{
									Optional: true,
								},
								"rtsp_disable": schema.BoolAttribute{
									Optional: true,
								},
								"sccp_disable": schema.BoolAttribute{
									Optional: true,
								},
								"sip_disable": schema.BoolAttribute{
									Optional: true,
								},
								"sql_disable": schema.BoolAttribute{
									Optional: true,
								},
								"sunrpc_disable": schema.BoolAttribute{
									Optional: true,
								},
								"talk_disable": schema.BoolAttribute{
									Optional: true,
								},
								"tftp_disable": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"flow": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"allow_dns_reply": schema.BoolAttribute{
									Optional: true,
								},
								"allow_embedded_icmp": schema.BoolAttribute{
									Optional: true,
								},
								"allow_reverse_ecmp": schema.BoolAttribute{
									Optional: true,
								},
								"enable_reroute_uniform_link_check_nat": schema.BoolAttribute{
									Optional: true,
								},
								"force_ip_reassembly": schema.BoolAttribute{
									Optional: true,
								},
								"ipsec_performance_acceleration": schema.BoolAttribute{
									Optional: true,
								},
								"mcast_buffer_enhance": schema.BoolAttribute{
									Optional: true,
								},
								"pending_sess_queue_length": schema.StringAttribute{
									Optional: true,
								},
								"preserve_incoming_fragment_size": schema.BoolAttribute{
									Optional: true,
								},
								"route_change_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"syn_flood_protection_mode": schema.StringAttribute{
									Optional: true,
								},
								"sync_icmp_session": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"advanced_options": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"drop_matching_link_local_address": schema.BoolAttribute{
												Optional: true,
											},
											"drop_matching_reserved_ip_address": schema.BoolAttribute{
												Optional: true,
											},
											"reverse_route_packet_mode_vr": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"aging": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"early_ageout": schema.Int64Attribute{
												Optional: true,
											},
											"high_watermark": schema.Int64Attribute{
												Optional: true,
											},
											"low_watermark": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"ethernet_switching": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"block_non_ip_all": schema.BoolAttribute{
												Optional: true,
											},
											"bypass_non_ip_unicast": schema.BoolAttribute{
												Optional: true,
											},
											"bpdu_vlan_flooding": schema.BoolAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"no_packet_flooding": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"no_trace_route": schema.BoolAttribute{
															Optional: true,
														},
													},
												},
											},
										},
									},
								},
								"tcp_mss": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"all_tcp_mss": schema.Int64Attribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"gre_in": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"mss": schema.Int64Attribute{
															Optional: true,
														},
													},
												},
											},
											"gre_out": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"mss": schema.Int64Attribute{
															Optional: true,
														},
													},
												},
											},
											"ipsec_vpn": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"mss": schema.Int64Attribute{
															Optional: true,
														},
													},
												},
											},
										},
									},
								},
								"tcp_session": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"fin_invalidate_session": schema.BoolAttribute{
												Optional: true,
											},
											"maximum_window": schema.StringAttribute{
												Optional: true,
											},
											"no_sequence_check": schema.BoolAttribute{
												Optional: true,
											},
											"no_syn_check": schema.BoolAttribute{
												Optional: true,
											},
											"no_syn_check_in_tunnel": schema.BoolAttribute{
												Optional: true,
											},
											"rst_invalidate_session": schema.BoolAttribute{
												Optional: true,
											},
											"rst_sequence_check": schema.BoolAttribute{
												Optional: true,
											},
											"strict_syn_check": schema.BoolAttribute{
												Optional: true,
											},
											"tcp_initial_timeout": schema.Int64Attribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"time_wait_state": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"apply_to_half_close_state": schema.BoolAttribute{
															Optional: true,
														},
														"session_ageout": schema.BoolAttribute{
															Optional: true,
														},
														"session_timeout": schema.Int64Attribute{
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
					"forwarding_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"inet6_mode": schema.StringAttribute{
									Optional: true,
								},
								"iso_mode_packet_based": schema.BoolAttribute{
									Optional: true,
								},
								"mpls_mode": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"forwarding_process": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"enhanced_services_mode": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"idp_security_package": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"automatic_enable": schema.BoolAttribute{
									Optional: true,
								},
								"automatic_interval": schema.Int64Attribute{
									Optional: true,
								},
								"automatic_start_time": schema.StringAttribute{
									Optional: true,
								},
								"install_ignore_version_check": schema.BoolAttribute{
									Optional: true,
								},
								"proxy_profile": schema.StringAttribute{
									Optional: true,
								},
								"source_address": schema.StringAttribute{
									Optional: true,
								},
								"url": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"idp_sensor_configuration": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"log_cache_size": schema.Int64Attribute{
									Optional: true,
								},
								"security_configuration_protection_mode": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"log_suppression": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"disable": schema.BoolAttribute{
												Optional: true,
											},
											"include_destination_address": schema.BoolAttribute{
												Optional: true,
											},
											"no_include_destination_address": schema.BoolAttribute{
												Optional: true,
											},
											"max_logs_operate": schema.Int64Attribute{
												Optional: true,
											},
											"max_time_report": schema.Int64Attribute{
												Optional: true,
											},
											"start_log": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"packet_log": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"source_address": schema.StringAttribute{
												Required: true,
											},
											"host_address": schema.StringAttribute{
												Optional: true,
											},
											"host_port": schema.Int64Attribute{
												Optional: true,
											},
											"max_sessions": schema.Int64Attribute{
												Optional: true,
											},
											"threshold_logging_interval": schema.Int64Attribute{
												Optional: true,
											},
											"total_memory": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"ike_traceoptions": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"flag": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"no_remote_trace": schema.BoolAttribute{
									Optional: true,
								},
								"rate_limit": schema.Int64Attribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"file": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Optional: true,
											},
											"files": schema.Int64Attribute{
												Optional: true,
											},
											"match": schema.StringAttribute{
												Optional: true,
											},
											"size": schema.Int64Attribute{
												Optional: true,
											},
											"world_readable": schema.BoolAttribute{
												Optional: true,
											},
											"no_world_readable": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"log": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"disable": schema.BoolAttribute{
									Optional: true,
								},
								"event_rate": schema.Int64Attribute{
									Optional: true,
								},
								"facility_override": schema.StringAttribute{
									Optional: true,
								},
								"format": schema.StringAttribute{
									Optional: true,
								},
								"max_database_record": schema.Int64Attribute{
									Optional: true,
								},
								"mode": schema.StringAttribute{
									Optional: true,
								},
								"rate_cap": schema.Int64Attribute{
									Optional: true,
								},
								"report": schema.BoolAttribute{
									Optional: true,
								},
								"source_address": schema.StringAttribute{
									Optional: true,
								},
								"source_interface": schema.StringAttribute{
									Optional: true,
								},
								"utc_timestamp": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"file": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"files": schema.Int64Attribute{
												Optional: true,
											},
											"name": schema.StringAttribute{
												Optional: true,
											},
											"path": schema.StringAttribute{
												Optional: true,
											},
											"size": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"transport": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"protocol": schema.StringAttribute{
												Optional: true,
											},
											"tcp_connections": schema.Int64Attribute{
												Optional: true,
											},
											"tls_profile": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"policies": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"policy_rematch": schema.BoolAttribute{
									Optional: true,
								},
								"policy_rematch_extensive": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"user_identification_auth_source": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ad_auth_priority": schema.Int64Attribute{
									Optional: true,
								},
								"aruba_clearpass_priority": schema.Int64Attribute{
									Optional: true,
								},
								"firewall_auth_priority": schema.Int64Attribute{
									Optional: true,
								},
								"local_auth_priority": schema.Int64Attribute{
									Optional: true,
								},
								"unified_access_control_priority": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"utm": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"feature_profile_web_filtering_type": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"feature_profile_web_filtering_juniper_enhanced_server": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"host": schema.StringAttribute{
												Optional: true,
											},
											"port": schema.Int64Attribute{
												Optional: true,
											},
											"proxy_profile": schema.StringAttribute{
												Optional: true,
											},
											"routing_instance": schema.StringAttribute{
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
			StateUpgrader: upgradeSecurityV0toV1,
		},
	}
}

//nolint:lll
func upgradeSecurityV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID             types.String `tfsdk:"id"`
		CleanOnDestroy types.Bool   `tfsdk:"clean_on_destroy"`
		Alg            []struct {
			DNSDisable    types.Bool `tfsdk:"dns_disable"`
			FtpDisable    types.Bool `tfsdk:"ftp_disable"`
			H323Disable   types.Bool `tfsdk:"h323_disable"`
			MgcpDisable   types.Bool `tfsdk:"mgcp_disable"`
			MsrpcDisable  types.Bool `tfsdk:"msrpc_disable"`
			PptpDisable   types.Bool `tfsdk:"pptp_disable"`
			RshDisable    types.Bool `tfsdk:"rsh_disable"`
			RtspDisable   types.Bool `tfsdk:"rtsp_disable"`
			SccpDisable   types.Bool `tfsdk:"sccp_disable"`
			SIPDisable    types.Bool `tfsdk:"sip_disable"`
			SQLDisable    types.Bool `tfsdk:"sql_disable"`
			SunrpcDisable types.Bool `tfsdk:"sunrpc_disable"`
			TalkDisable   types.Bool `tfsdk:"talk_disable"`
			TftpDisable   types.Bool `tfsdk:"tftp_disable"`
		} `tfsdk:"alg"`
		Flow []struct {
			AllowDNSReply                    types.Bool   `tfsdk:"allow_dns_reply"`
			AllowEmbeddedIcmp                types.Bool   `tfsdk:"allow_embedded_icmp"`
			AllowReverseEcmp                 types.Bool   `tfsdk:"allow_reverse_ecmp"`
			EnableRerouteUniformLinkCheckNat types.Bool   `tfsdk:"enable_reroute_uniform_link_check_nat"`
			ForceIPReassembly                types.Bool   `tfsdk:"force_ip_reassembly"`
			IpsecPerformanceAcceleration     types.Bool   `tfsdk:"ipsec_performance_acceleration"`
			McastBufferEnhance               types.Bool   `tfsdk:"mcast_buffer_enhance"`
			PendingSessQueueLength           types.String `tfsdk:"pending_sess_queue_length"`
			PreserveIncomingFragmentSize     types.Bool   `tfsdk:"preserve_incoming_fragment_size"`
			RouteChangeTimeout               types.Int64  `tfsdk:"route_change_timeout"`
			SynFloodProtectionMode           types.String `tfsdk:"syn_flood_protection_mode"`
			SyncIcmpSession                  types.Bool   `tfsdk:"sync_icmp_session"`
			AdvancedOptions                  []struct {
				DropMatchingLinkLocalAddress  types.Bool `tfsdk:"drop_matching_link_local_address"`
				DropMatchingReservedIPAddress types.Bool `tfsdk:"drop_matching_reserved_ip_address"`
				ReverseRoutePacketModeVR      types.Bool `tfsdk:"reverse_route_packet_mode_vr"`
			} `tfsdk:"advanced_options"`
			Aging []struct {
				EarlyAgeout   types.Int64 `tfsdk:"early_ageout"`
				HighWatermark types.Int64 `tfsdk:"high_watermark"`
				LowWatermark  types.Int64 `tfsdk:"low_watermark"`
			} `tfsdk:"aging"`
			EthernetSwitching []struct {
				BlockNonIPAll      types.Bool `tfsdk:"block_non_ip_all"`
				BypassNonIPUnicast types.Bool `tfsdk:"bypass_non_ip_unicast"`
				BpduVlanFlooding   types.Bool `tfsdk:"bpdu_vlan_flooding"`
				NoPacketFlooding   []struct {
					NoTraceRoute types.Bool `tfsdk:"no_trace_route"`
				} `tfsdk:"no_packet_flooding"`
			} `tfsdk:"ethernet_switching"`
			TCPMss []struct {
				AllTCPMss types.Int64 `tfsdk:"all_tcp_mss"`
				GreIn     []struct {
					Mss types.Int64 `tfsdk:"mss"`
				} `tfsdk:"gre_in"`
				GreOut []struct {
					Mss types.Int64 `tfsdk:"mss"`
				} `tfsdk:"gre_out"`
				IpsecVpn []struct {
					Mss types.Int64 `tfsdk:"mss"`
				} `tfsdk:"ipsec_vpn"`
			} `tfsdk:"tcp_mss"`
			TCPSession []struct {
				FinInvalidateSession types.Bool   `tfsdk:"fin_invalidate_session"`
				MaximumWindow        types.String `tfsdk:"maximum_window"`
				NoSequenceCheck      types.Bool   `tfsdk:"no_sequence_check"`
				NoSynCheck           types.Bool   `tfsdk:"no_syn_check"`
				NoSynCheckInTunnel   types.Bool   `tfsdk:"no_syn_check_in_tunnel"`
				RstInvalidateSession types.Bool   `tfsdk:"rst_invalidate_session"`
				RstSequenceCheck     types.Bool   `tfsdk:"rst_sequence_check"`
				StrictSynCheck       types.Bool   `tfsdk:"strict_syn_check"`
				TCPInitialTimeout    types.Int64  `tfsdk:"tcp_initial_timeout"`
				TimeWaitState        []struct {
					ApplyToHalfCloseState types.Bool  `tfsdk:"apply_to_half_close_state"`
					SessionAgeout         types.Bool  `tfsdk:"session_ageout"`
					SessionTimeout        types.Int64 `tfsdk:"session_timeout"`
				} `tfsdk:"time_wait_state"`
			} `tfsdk:"tcp_session"`
		} `tfsdk:"flow"`
		ForwardingOptions []struct {
			Inet6Mode          types.String `tfsdk:"inet6_mode"`
			IsoModePacketBased types.Bool   `tfsdk:"iso_mode_packet_based"`
			MplsMode           types.String `tfsdk:"mpls_mode"`
		} `tfsdk:"forwarding_options"`
		ForwardingProcess []struct {
			EnhancedServicesMode types.Bool `tfsdk:"enhanced_services_mode"`
		} `tfsdk:"forwarding_process"`
		IdpSecurityPackage []struct {
			AutomaticEnable           types.Bool   `tfsdk:"automatic_enable"`
			AutomaticInterval         types.Int64  `tfsdk:"automatic_interval"`
			AutomaticStartTime        types.String `tfsdk:"automatic_start_time"`
			InstallIgnoreVersionCheck types.Bool   `tfsdk:"install_ignore_version_check"`
			ProxyProfile              types.String `tfsdk:"proxy_profile"`
			SourceAddress             types.String `tfsdk:"source_address"`
			URL                       types.String `tfsdk:"url"`
		} `tfsdk:"idp_security_package"`
		IdpSensorConfiguration []struct {
			LogCacheSize                        types.Int64  `tfsdk:"log_cache_size"`
			SecurityConfigurationProtectionMode types.String `tfsdk:"security_configuration_protection_mode"`
			LogSuppression                      []struct {
				Disable                     types.Bool  `tfsdk:"disable"`
				IncludeDestinationAddress   types.Bool  `tfsdk:"include_destination_address"`
				NoIncludeDestinationAddress types.Bool  `tfsdk:"no_include_destination_address"`
				MaxLogsOperate              types.Int64 `tfsdk:"max_logs_operate"`
				MaxTimeReport               types.Int64 `tfsdk:"max_time_report"`
				StartLog                    types.Int64 `tfsdk:"start_log"`
			} `tfsdk:"log_suppression"`
			PacketLog []struct {
				SourceAddress            types.String `tfsdk:"source_address"`
				HostAddress              types.String `tfsdk:"host_address"`
				HostPort                 types.Int64  `tfsdk:"host_port"`
				MaxSessions              types.Int64  `tfsdk:"max_sessions"`
				ThresholdLoggingInterval types.Int64  `tfsdk:"threshold_logging_interval"`
				TotalMemory              types.Int64  `tfsdk:"total_memory"`
			} `tfsdk:"packet_log"`
		} `tfsdk:"idp_sensor_configuration"`
		IkeTraceoptions []struct {
			Flag          []types.String `tfsdk:"flag"`
			NoRemoteTrace types.Bool     `tfsdk:"no_remote_trace"`
			RateLimit     types.Int64    `tfsdk:"rate_limit"`
			File          []struct {
				Name            types.String `tfsdk:"name"`
				Files           types.Int64  `tfsdk:"files"`
				Match           types.String `tfsdk:"match"`
				Size            types.Int64  `tfsdk:"size"`
				WorldReadable   types.Bool   `tfsdk:"world_readable"`
				NoWorldReadable types.Bool   `tfsdk:"no_world_readable"`
			} `tfsdk:"file"`
		} `tfsdk:"ike_traceoptions"`
		Log []struct {
			Disable           types.Bool   `tfsdk:"disable"`
			EventRate         types.Int64  `tfsdk:"event_rate"`
			FacilityOverride  types.String `tfsdk:"facility_override"`
			Format            types.String `tfsdk:"format"`
			MaxDatabaseRecord types.Int64  `tfsdk:"max_database_record"`
			Mode              types.String `tfsdk:"mode"`
			RateCap           types.Int64  `tfsdk:"rate_cap"`
			Report            types.Bool   `tfsdk:"report"`
			SourceAddress     types.String `tfsdk:"source_address"`
			SourceInterface   types.String `tfsdk:"source_interface"`
			UtcTimestamp      types.Bool   `tfsdk:"utc_timestamp"`
			File              []struct {
				Files types.Int64  `tfsdk:"files"`
				Name  types.String `tfsdk:"name"`
				Path  types.String `tfsdk:"path"`
				Size  types.Int64  `tfsdk:"size"`
			} `tfsdk:"file"`
			Transport []struct {
				Protocol       types.String `tfsdk:"protocol"`
				TCPConnections types.Int64  `tfsdk:"tcp_connections"`
				TLSProfile     types.String `tfsdk:"tls_profile"`
			} `tfsdk:"transport"`
		} `tfsdk:"log"`
		Policies []struct {
			PolicyRematch          types.Bool `tfsdk:"policy_rematch"`
			PolicyRematchExtensive types.Bool `tfsdk:"policy_rematch_extensive"`
		} `tfsdk:"policies"`
		UserIdentificationAuthSource []struct {
			ADAuthPriority               types.Int64 `tfsdk:"ad_auth_priority"`
			ArubaClearpassPriority       types.Int64 `tfsdk:"aruba_clearpass_priority"`
			FirewallAuthPriority         types.Int64 `tfsdk:"firewall_auth_priority"`
			LocalAuthPriority            types.Int64 `tfsdk:"local_auth_priority"`
			UnifiedAccessControlPriority types.Int64 `tfsdk:"unified_access_control_priority"`
		} `tfsdk:"user_identification_auth_source"`
		Utm []struct {
			FeatureProfileWebFilteringType                  types.String `tfsdk:"feature_profile_web_filtering_type"`
			FeatureProfileWebFilteringJuniperEnhancedServer []struct {
				Host            types.String `tfsdk:"host"`
				Port            types.Int64  `tfsdk:"port"`
				ProxyProfile    types.String `tfsdk:"proxy_profile"`
				RoutingInstance types.String `tfsdk:"routing_instance"`
			} `tfsdk:"feature_profile_web_filtering_juniper_enhanced_server"`
		} `tfsdk:"utm"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityData
	dataV1.ID = dataV0.ID
	dataV1.CleanOnDestroy = dataV0.CleanOnDestroy
	if !dataV1.CleanOnDestroy.IsNull() && !dataV1.CleanOnDestroy.ValueBool() {
		dataV1.CleanOnDestroy = types.BoolNull()
	}
	if len(dataV0.Alg) > 0 {
		dataV1.Alg = &securityBlockAlg{
			DNSDisable:    dataV0.Alg[0].DNSDisable,
			FtpDisable:    dataV0.Alg[0].FtpDisable,
			H323Disable:   dataV0.Alg[0].H323Disable,
			MgcpDisable:   dataV0.Alg[0].MgcpDisable,
			MsrpcDisable:  dataV0.Alg[0].MsrpcDisable,
			PptpDisable:   dataV0.Alg[0].PptpDisable,
			RshDisable:    dataV0.Alg[0].RshDisable,
			RtspDisable:   dataV0.Alg[0].RtspDisable,
			SccpDisable:   dataV0.Alg[0].SccpDisable,
			SIPDisable:    dataV0.Alg[0].SIPDisable,
			SQLDisable:    dataV0.Alg[0].SQLDisable,
			SunrpcDisable: dataV0.Alg[0].SunrpcDisable,
			TalkDisable:   dataV0.Alg[0].TalkDisable,
			TftpDisable:   dataV0.Alg[0].TftpDisable,
		}
	}
	if len(dataV0.Flow) > 0 {
		dataV1.Flow = &securityBlockFlow{
			AllowDNSReply:                    dataV0.Flow[0].AllowDNSReply,
			AllowEmbeddedIcmp:                dataV0.Flow[0].AllowEmbeddedIcmp,
			AllowReverseEcmp:                 dataV0.Flow[0].AllowReverseEcmp,
			EnableRerouteUniformLinkCheckNat: dataV0.Flow[0].EnableRerouteUniformLinkCheckNat,
			ForceIPReassembly:                dataV0.Flow[0].ForceIPReassembly,
			IpsecPerformanceAcceleration:     dataV0.Flow[0].IpsecPerformanceAcceleration,
			McastBufferEnhance:               dataV0.Flow[0].McastBufferEnhance,
			PreserveIncomingFragmentSize:     dataV0.Flow[0].PreserveIncomingFragmentSize,
			SyncIcmpSession:                  dataV0.Flow[0].SyncIcmpSession,
			PendingSessQueueLength:           dataV0.Flow[0].PendingSessQueueLength,
			RouteChangeTimeout:               dataV0.Flow[0].RouteChangeTimeout,
			SynFloodProtectionMode:           dataV0.Flow[0].SynFloodProtectionMode,
		}
		if len(dataV0.Flow[0].AdvancedOptions) > 0 {
			dataV1.Flow.AdvancedOptions = &securityBlockFlowBlockAdvancedOptions{
				DropMatchingLinkLocalAddress:  dataV0.Flow[0].AdvancedOptions[0].DropMatchingLinkLocalAddress,
				DropMatchingReservedIPAddress: dataV0.Flow[0].AdvancedOptions[0].DropMatchingReservedIPAddress,
				ReverseRoutePacketModeVR:      dataV0.Flow[0].AdvancedOptions[0].ReverseRoutePacketModeVR,
			}
		}
		if len(dataV0.Flow[0].Aging) > 0 {
			dataV1.Flow.Aging = &securityBlockFlowBlockAging{
				EarlyAgeout:   dataV0.Flow[0].Aging[0].EarlyAgeout,
				HighWatermark: dataV0.Flow[0].Aging[0].HighWatermark,
				LowWatermark:  dataV0.Flow[0].Aging[0].LowWatermark,
			}
		}
		if len(dataV0.Flow[0].EthernetSwitching) > 0 {
			dataV1.Flow.EthernetSwitching = &securityBlockFlowBlockEthernetSwitching{
				BlockNonIPAll:      dataV0.Flow[0].EthernetSwitching[0].BlockNonIPAll,
				BypassNonIPUnicast: dataV0.Flow[0].EthernetSwitching[0].BypassNonIPUnicast,
				BpduVlanFlooding:   dataV0.Flow[0].EthernetSwitching[0].BpduVlanFlooding,
			}
			if len(dataV0.Flow[0].EthernetSwitching[0].NoPacketFlooding) > 0 {
				dataV1.Flow.EthernetSwitching.NoPacketFlooding = &dataV0.Flow[0].EthernetSwitching[0].NoPacketFlooding[0]
			}
		}
		if len(dataV0.Flow[0].TCPMss) > 0 {
			dataV1.Flow.TCPMss = &securityBlockFlowBlockTCPMss{
				AllTCPMss: dataV0.Flow[0].TCPMss[0].AllTCPMss,
			}
			if len(dataV0.Flow[0].TCPMss[0].GreIn) > 0 {
				dataV1.Flow.TCPMss.GreIn = &dataV0.Flow[0].TCPMss[0].GreIn[0]
			}
			if len(dataV0.Flow[0].TCPMss[0].GreOut) > 0 {
				dataV1.Flow.TCPMss.GreOut = &dataV0.Flow[0].TCPMss[0].GreOut[0]
			}
			if len(dataV0.Flow[0].TCPMss[0].IpsecVpn) > 0 {
				dataV1.Flow.TCPMss.IpsecVpn = &dataV0.Flow[0].TCPMss[0].IpsecVpn[0]
			}
		}
		if len(dataV0.Flow[0].TCPSession) > 0 {
			dataV1.Flow.TCPSession = &securityBlockFlowBlockTCPSession{
				FinInvalidateSession: dataV0.Flow[0].TCPSession[0].FinInvalidateSession,
				NoSequenceCheck:      dataV0.Flow[0].TCPSession[0].NoSequenceCheck,
				NoSynCheck:           dataV0.Flow[0].TCPSession[0].NoSynCheck,
				NoSynCheckInTunnel:   dataV0.Flow[0].TCPSession[0].NoSynCheckInTunnel,
				RstInvalidateSession: dataV0.Flow[0].TCPSession[0].RstInvalidateSession,
				RstSequenceCheck:     dataV0.Flow[0].TCPSession[0].RstSequenceCheck,
				StrictSynCheck:       dataV0.Flow[0].TCPSession[0].StrictSynCheck,
				MaximumWindow:        dataV0.Flow[0].TCPSession[0].MaximumWindow,
				TCPInitialTimeout:    dataV0.Flow[0].TCPSession[0].TCPInitialTimeout,
			}
			if len(dataV0.Flow[0].TCPSession[0].TimeWaitState) > 0 {
				dataV1.Flow.TCPSession.TimeWaitState = &securityBlockFlowBlockTCPSessionBlockTimeWaitState{
					ApplyToHalfCloseState: dataV0.Flow[0].TCPSession[0].TimeWaitState[0].ApplyToHalfCloseState,
					SessionAgeout:         dataV0.Flow[0].TCPSession[0].TimeWaitState[0].SessionAgeout,
					SessionTimeout:        dataV0.Flow[0].TCPSession[0].TimeWaitState[0].SessionTimeout,
				}
			}
		}
	}
	if len(dataV0.ForwardingOptions) > 0 {
		dataV1.ForwardingOptions = &securityBlockForwardingOptions{
			Inet6Mode:          dataV0.ForwardingOptions[0].Inet6Mode,
			IsoModePacketBased: dataV0.ForwardingOptions[0].IsoModePacketBased,
			MplsMode:           dataV0.ForwardingOptions[0].MplsMode,
		}
	}
	if len(dataV0.ForwardingProcess) > 0 {
		dataV1.ForwardingProcess = &securityBlockForwardingProcess{
			EnhancedServicesMode: dataV0.ForwardingProcess[0].EnhancedServicesMode,
		}
	}
	if len(dataV0.IdpSecurityPackage) > 0 {
		dataV1.IdpSecurityPackage = &securityBlockIdpSecurityPackage{
			AutomaticEnable:           dataV0.IdpSecurityPackage[0].AutomaticEnable,
			InstallIgnoreVersionCheck: dataV0.IdpSecurityPackage[0].InstallIgnoreVersionCheck,
			AutomaticInterval:         dataV0.IdpSecurityPackage[0].AutomaticInterval,
			AutomaticStartTime:        dataV0.IdpSecurityPackage[0].AutomaticStartTime,
			ProxyProfile:              dataV0.IdpSecurityPackage[0].ProxyProfile,
			SourceAddress:             dataV0.IdpSecurityPackage[0].SourceAddress,
			URL:                       dataV0.IdpSecurityPackage[0].URL,
		}
	}
	if len(dataV0.IdpSensorConfiguration) > 0 {
		dataV1.IdpSensorConfiguration = &securityBlockIdpSensorConfiguration{
			LogCacheSize:                        dataV0.IdpSensorConfiguration[0].LogCacheSize,
			SecurityConfigurationProtectionMode: dataV0.IdpSensorConfiguration[0].SecurityConfigurationProtectionMode,
		}
		if len(dataV0.IdpSensorConfiguration[0].LogSuppression) > 0 {
			dataV1.IdpSensorConfiguration.LogSuppression = &securityBlockIdpSensorConfigurationBlockLogSuppression{
				Disable:                     dataV0.IdpSensorConfiguration[0].LogSuppression[0].Disable,
				IncludeDestinationAddress:   dataV0.IdpSensorConfiguration[0].LogSuppression[0].IncludeDestinationAddress,
				NoIncludeDestinationAddress: dataV0.IdpSensorConfiguration[0].LogSuppression[0].NoIncludeDestinationAddress,
				MaxLogsOperate:              dataV0.IdpSensorConfiguration[0].LogSuppression[0].MaxLogsOperate,
				MaxTimeReport:               dataV0.IdpSensorConfiguration[0].LogSuppression[0].MaxTimeReport,
				StartLog:                    dataV0.IdpSensorConfiguration[0].LogSuppression[0].StartLog,
			}
		}
		if len(dataV0.IdpSensorConfiguration[0].PacketLog) > 0 {
			dataV1.IdpSensorConfiguration.PacketLog = &securityBlockIdpSensorConfigurationBlockPacketLog{
				SourceAddress:            dataV0.IdpSensorConfiguration[0].PacketLog[0].SourceAddress,
				HostAddress:              dataV0.IdpSensorConfiguration[0].PacketLog[0].HostAddress,
				HostPort:                 dataV0.IdpSensorConfiguration[0].PacketLog[0].HostPort,
				MaxSessions:              dataV0.IdpSensorConfiguration[0].PacketLog[0].MaxSessions,
				ThresholdLoggingInterval: dataV0.IdpSensorConfiguration[0].PacketLog[0].ThresholdLoggingInterval,
				TotalMemory:              dataV0.IdpSensorConfiguration[0].PacketLog[0].TotalMemory,
			}
		}
	}
	if len(dataV0.IkeTraceoptions) > 0 {
		dataV1.IkeTraceoptions = &securityBlockIkeTraceoptions{
			Flag:          dataV0.IkeTraceoptions[0].Flag,
			NoRemoteTrace: dataV0.IkeTraceoptions[0].NoRemoteTrace,
			RateLimit:     dataV0.IkeTraceoptions[0].RateLimit,
		}
		if len(dataV0.IkeTraceoptions[0].File) > 0 {
			dataV1.IkeTraceoptions.File = &securityBlockIkeTraceoptionsBlockFile{
				NoWorldReadable: dataV0.IkeTraceoptions[0].File[0].NoWorldReadable,
				WorldReadable:   dataV0.IkeTraceoptions[0].File[0].WorldReadable,
				Name:            dataV0.IkeTraceoptions[0].File[0].Name,
				Files:           dataV0.IkeTraceoptions[0].File[0].Files,
				Match:           dataV0.IkeTraceoptions[0].File[0].Match,
				Size:            dataV0.IkeTraceoptions[0].File[0].Size,
			}
		}
	}
	if len(dataV0.Log) > 0 {
		dataV1.Log = &securityBlockLog{
			Disable:           dataV0.Log[0].Disable,
			Report:            dataV0.Log[0].Report,
			UtcTimestamp:      dataV0.Log[0].UtcTimestamp,
			EventRate:         dataV0.Log[0].EventRate,
			FacilityOverride:  dataV0.Log[0].FacilityOverride,
			Format:            dataV0.Log[0].Format,
			MaxDatabaseRecord: dataV0.Log[0].MaxDatabaseRecord,
			Mode:              dataV0.Log[0].Mode,
			RateCap:           dataV0.Log[0].RateCap,
			SourceAddress:     dataV0.Log[0].SourceAddress,
			SourceInterface:   dataV0.Log[0].SourceInterface,
		}
		if len(dataV0.Log[0].File) > 0 {
			dataV1.Log.File = &securityBlockLogBlockFile{
				Files: dataV0.Log[0].File[0].Files,
				Name:  dataV0.Log[0].File[0].Name,
				Path:  dataV0.Log[0].File[0].Path,
				Size:  dataV0.Log[0].File[0].Size,
			}
		}
		if len(dataV0.Log[0].Transport) > 0 {
			dataV1.Log.Transport = &securityBlockLogBlockTransport{
				Protocol:       dataV0.Log[0].Transport[0].Protocol,
				TCPConnections: dataV0.Log[0].Transport[0].TCPConnections,
				TLSProfile:     dataV0.Log[0].Transport[0].TLSProfile,
			}
		}
	}
	if len(dataV0.Policies) > 0 {
		dataV1.Policies = &securityBlockPolicies{
			PolicyRematch:          dataV0.Policies[0].PolicyRematch,
			PolicyRematchExtensive: dataV0.Policies[0].PolicyRematchExtensive,
		}
	}
	if len(dataV0.UserIdentificationAuthSource) > 0 {
		dataV1.UserIdentificationAuthSource = &securityBlockUserIdentificationAuthSource{
			ADAuthPriority:               dataV0.UserIdentificationAuthSource[0].ADAuthPriority,
			ArubaClearpassPriority:       dataV0.UserIdentificationAuthSource[0].ArubaClearpassPriority,
			FirewallAuthPriority:         dataV0.UserIdentificationAuthSource[0].FirewallAuthPriority,
			LocalAuthPriority:            dataV0.UserIdentificationAuthSource[0].LocalAuthPriority,
			UnifiedAccessControlPriority: dataV0.UserIdentificationAuthSource[0].UnifiedAccessControlPriority,
		}
	}
	if len(dataV0.Utm) > 0 {
		dataV1.Utm = &securityBlockUtm{
			FeatureProfileWebFilteringType: dataV0.Utm[0].FeatureProfileWebFilteringType,
		}
		if len(dataV0.Utm[0].FeatureProfileWebFilteringJuniperEnhancedServer) > 0 {
			dataV1.Utm.FeatureProfileWebFilteringJuniperEnhancedServer = &securityBlockUtmBlockFeatureProfileWebFilteringJuniperEnhancedServer{
				Host:            dataV0.Utm[0].FeatureProfileWebFilteringJuniperEnhancedServer[0].Host,
				Port:            dataV0.Utm[0].FeatureProfileWebFilteringJuniperEnhancedServer[0].Port,
				ProxyProfile:    dataV0.Utm[0].FeatureProfileWebFilteringJuniperEnhancedServer[0].ProxyProfile,
				RoutingInstance: dataV0.Utm[0].FeatureProfileWebFilteringJuniperEnhancedServer[0].RoutingInstance,
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
