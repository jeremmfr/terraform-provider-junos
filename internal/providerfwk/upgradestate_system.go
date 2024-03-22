package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *system) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"authentication_order": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"auto_snapshot": schema.BoolAttribute{
						Optional: true,
					},
					"default_address_selection": schema.BoolAttribute{
						Optional: true,
					},
					"domain_name": schema.StringAttribute{
						Optional: true,
					},
					"host_name": schema.StringAttribute{
						Optional: true,
					},
					"max_configuration_rollbacks": schema.Int64Attribute{
						Optional: true,
					},
					"max_configurations_on_flash": schema.Int64Attribute{
						Optional: true,
					},
					"name_server": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"no_multicast_echo": schema.BoolAttribute{
						Optional: true,
					},
					"no_ping_record_route": schema.BoolAttribute{
						Optional: true,
					},
					"no_ping_time_stamp": schema.BoolAttribute{
						Optional: true,
					},
					"no_redirects": schema.BoolAttribute{
						Optional: true,
					},
					"no_redirects_ipv6": schema.BoolAttribute{
						Optional: true,
					},
					"radius_options_attributes_nas_ipaddress": schema.StringAttribute{
						Optional: true,
					},
					"radius_options_enhanced_accounting": schema.BoolAttribute{
						Optional: true,
					},
					"radius_options_password_protocol_mschapv2": schema.BoolAttribute{
						Optional: true,
					},
					"time_zone": schema.StringAttribute{
						Optional: true,
					},
					"tracing_dest_override_syslog_host": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"archival_configuration": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"transfer_interval": schema.Int64Attribute{
									Optional: true,
								},
								"transfer_on_commit": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"archive_site": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"url": schema.StringAttribute{
												Required: true,
											},
											"password": schema.StringAttribute{
												Optional:  true,
												Sensitive: true,
											},
										},
									},
								},
							},
						},
					},
					"inet6_backup_router": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"address": schema.StringAttribute{
									Required: true,
								},
								"destination": schema.SetAttribute{
									ElementType: types.StringType,
									Required:    true,
								},
							},
						},
					},
					"internet_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"gre_path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"no_gre_path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"ipip_path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"no_ipip_path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"ipv6_duplicate_addr_detection_transmits": schema.Int64Attribute{
									Optional: true,
								},
								"ipv6_path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"no_ipv6_path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"ipv6_path_mtu_discovery_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"ipv6_reject_zero_hop_limit": schema.BoolAttribute{
									Optional: true,
								},
								"no_ipv6_reject_zero_hop_limit": schema.BoolAttribute{
									Optional: true,
								},
								"no_tcp_reset": schema.StringAttribute{
									Optional: true,
								},
								"no_tcp_rfc1323": schema.BoolAttribute{
									Optional: true,
								},
								"no_tcp_rfc1323_paws": schema.BoolAttribute{
									Optional: true,
								},
								"path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"no_path_mtu_discovery": schema.BoolAttribute{
									Optional: true,
								},
								"source_port_upper_limit": schema.Int64Attribute{
									Optional: true,
								},
								"source_quench": schema.BoolAttribute{
									Optional: true,
								},
								"no_source_quench": schema.BoolAttribute{
									Optional: true,
								},
								"tcp_drop_synfin_set": schema.BoolAttribute{
									Optional: true,
								},
								"tcp_mss": schema.Int64Attribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"icmpv4_rate_limit": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"bucket_size": schema.Int64Attribute{
												Optional: true,
											},
											"packet_rate": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"icmpv6_rate_limit": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"bucket_size": schema.Int64Attribute{
												Optional: true,
											},
											"packet_rate": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"license": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"autoupdate": schema.BoolAttribute{
									Optional: true,
								},
								"autoupdate_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"autoupdate_url": schema.StringAttribute{
									Optional: true,
								},
								"renew_before_expiration": schema.Int64Attribute{
									Optional: true,
								},
								"renew_interval": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"login": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"announcement": schema.StringAttribute{
									Optional: true,
								},
								"deny_sources_address": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"idle_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"message": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"password": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"change_type": schema.StringAttribute{
												Optional: true,
											},
											"format": schema.StringAttribute{
												Optional: true,
											},
											"maximum_length": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_changes": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_character_changes": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_length": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_lower_cases": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_numerics": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_punctuations": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_reuse": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_upper_cases": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"retry_options": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"backoff_factor": schema.Int64Attribute{
												Optional: true,
											},
											"backoff_threshold": schema.Int64Attribute{
												Optional: true,
											},
											"lockout_period": schema.Int64Attribute{
												Optional: true,
											},
											"maximum_time": schema.Int64Attribute{
												Optional: true,
											},
											"minimum_time": schema.Int64Attribute{
												Optional: true,
											},
											"tries_before_disconnect": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"ntp": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"boot_server": schema.StringAttribute{
									Optional: true,
								},
								"broadcast_client": schema.BoolAttribute{
									Optional: true,
								},
								"interval_range": schema.Int64Attribute{
									Optional: true,
								},
								"multicast_client": schema.BoolAttribute{
									Optional: true,
								},
								"multicast_client_address": schema.StringAttribute{
									Optional: true,
								},
								"threshold_action": schema.StringAttribute{
									Optional: true,
								},
								"threshold_value": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"ports": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"auxiliary_authentication_order": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"auxiliary_disable": schema.BoolAttribute{
									Optional: true,
								},
								"auxiliary_insecure": schema.BoolAttribute{
									Optional: true,
								},
								"auxiliary_logout_on_disconnect": schema.BoolAttribute{
									Optional: true,
								},
								"auxiliary_type": schema.StringAttribute{
									Optional: true,
								},
								"console_authentication_order": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"console_disable": schema.BoolAttribute{
									Optional: true,
								},
								"console_insecure": schema.BoolAttribute{
									Optional: true,
								},
								"console_logout_on_disconnect": schema.BoolAttribute{
									Optional: true,
								},
								"console_type": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"services": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"netconf_ssh": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"client_alive_count_max": schema.Int64Attribute{
												Optional: true,
											},
											"client_alive_interval": schema.Int64Attribute{
												Optional: true,
											},
											"connection_limit": schema.Int64Attribute{
												Optional: true,
											},
											"rate_limit": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"netconf_traceoptions": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"file_name": schema.StringAttribute{
												Optional: true,
											},
											"file_files": schema.Int64Attribute{
												Optional: true,
											},
											"file_match": schema.StringAttribute{
												Optional: true,
											},
											"file_size": schema.Int64Attribute{
												Optional: true,
											},
											"file_world_readable": schema.BoolAttribute{
												Optional: true,
											},
											"file_no_world_readable": schema.BoolAttribute{
												Optional: true,
											},
											"flag": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"no_remote_trace": schema.BoolAttribute{
												Optional: true,
											},
											"on_demand": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"ssh": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"authentication_order": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"ciphers": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"client_alive_count_max": schema.Int64Attribute{
												Optional: true,
											},
											"client_alive_interval": schema.Int64Attribute{
												Optional: true,
											},
											"connection_limit": schema.Int64Attribute{
												Optional: true,
											},
											"fingerprint_hash": schema.StringAttribute{
												Optional: true,
											},
											"hostkey_algorithm": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"key_exchange": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"log_key_changes": schema.BoolAttribute{
												Optional: true,
											},
											"macs": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"max_pre_authentication_packets": schema.Int64Attribute{
												Optional: true,
											},
											"max_sessions_per_connection": schema.Int64Attribute{
												Optional: true,
											},
											"no_passwords": schema.BoolAttribute{
												Optional: true,
											},
											"no_public_keys": schema.BoolAttribute{
												Optional: true,
											},
											"port": schema.Int64Attribute{
												Optional: true,
											},
											"protocol_version": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"rate_limit": schema.Int64Attribute{
												Optional: true,
											},
											"root_login": schema.StringAttribute{
												Optional: true,
											},
											"tcp_forwarding": schema.BoolAttribute{
												Optional: true,
											},
											"no_tcp_forwarding": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"web_management_http": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"interface": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"port": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"web_management_https": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"interface": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"local_certificate": schema.StringAttribute{
												Optional: true,
											},
											"pki_local_certificate": schema.StringAttribute{
												Optional: true,
											},
											"port": schema.Int64Attribute{
												Optional: true,
											},
											"system_generated_certificate": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"syslog": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"log_rotate_frequency": schema.Int64Attribute{
									Optional: true,
								},
								"source_address": schema.StringAttribute{
									Optional: true,
								},
								"time_format_millisecond": schema.BoolAttribute{
									Optional: true,
								},
								"time_format_year": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"archive": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"binary_data": schema.BoolAttribute{
												Optional: true,
											},
											"no_binary_data": schema.BoolAttribute{
												Optional: true,
											},
											"files": schema.Int64Attribute{
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
								"console": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"any_severity": schema.StringAttribute{
												Optional: true,
											},
											"authorization_severity": schema.StringAttribute{
												Optional: true,
											},
											"changelog_severity": schema.StringAttribute{
												Optional: true,
											},
											"conflictlog_severity": schema.StringAttribute{
												Optional: true,
											},
											"daemon_severity": schema.StringAttribute{
												Optional: true,
											},
											"dfc_severity": schema.StringAttribute{
												Optional: true,
											},
											"external_severity": schema.StringAttribute{
												Optional: true,
											},
											"firewall_severity": schema.StringAttribute{
												Optional: true,
											},
											"ftp_severity": schema.StringAttribute{
												Optional: true,
											},
											"interactivecommands_severity": schema.StringAttribute{
												Optional: true,
											},
											"kernel_severity": schema.StringAttribute{
												Optional: true,
											},
											"ntp_severity": schema.StringAttribute{
												Optional: true,
											},
											"pfe_severity": schema.StringAttribute{
												Optional: true,
											},
											"security_severity": schema.StringAttribute{
												Optional: true,
											},
											"user_severity": schema.StringAttribute{
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
			StateUpgrader: upgradeSystemV0toV1,
		},
	}
}

func upgradeSystemV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                                    types.String   `tfsdk:"id"`
		AuthenticationOrder                   []types.String `tfsdk:"authentication_order"`
		AutoSnapshot                          types.Bool     `tfsdk:"auto_snapshot"`
		DefaultAddressSelection               types.Bool     `tfsdk:"default_address_selection"`
		DomainName                            types.String   `tfsdk:"domain_name"`
		HostName                              types.String   `tfsdk:"host_name"`
		MaxConfigurationRollbacks             types.Int64    `tfsdk:"max_configuration_rollbacks"`
		MaxConfigurationsOnFlash              types.Int64    `tfsdk:"max_configurations_on_flash"`
		NameServer                            []types.String `tfsdk:"name_server"`
		NoMulticastEcho                       types.Bool     `tfsdk:"no_multicast_echo"`
		NoPingRecordRoute                     types.Bool     `tfsdk:"no_ping_record_route"`
		NoPingTimestamp                       types.Bool     `tfsdk:"no_ping_time_stamp"`
		NoRedirects                           types.Bool     `tfsdk:"no_redirects"`
		NoRedirectsIPv6                       types.Bool     `tfsdk:"no_redirects_ipv6"`
		RadiusOptionsAttributesNasIpaddress   types.String   `tfsdk:"radius_options_attributes_nas_ipaddress"`
		RadiusOptionsEnhancedAccounting       types.Bool     `tfsdk:"radius_options_enhanced_accounting"`
		RadiusOptionsPasswordProtocolMschapv2 types.Bool     `tfsdk:"radius_options_password_protocol_mschapv2"`
		TimeZone                              types.String   `tfsdk:"time_zone"`
		TracingDestOverrideSyslogHost         types.String   `tfsdk:"tracing_dest_override_syslog_host"`
		ArchivalConfiguration                 []struct {
			TransferInterval types.Int64 `tfsdk:"transfer_interval"`
			TransferOnCommit types.Bool  `tfsdk:"transfer_on_commit"`
			ArchiveSite      []struct {
				URL      types.String `tfsdk:"url"`
				Password types.String `tfsdk:"password"`
			} `tfsdk:"archive_site"`
		} `tfsdk:"archival_configuration"`
		Inet6BackupRouter []struct {
			Address     types.String   `tfsdk:"address"`
			Destination []types.String `tfsdk:"destination"`
		} `tfsdk:"inet6_backup_router"`
		InternetOptions []struct {
			GrePathMtuDiscovery                 types.Bool   `tfsdk:"gre_path_mtu_discovery"`
			NoGrePathMtuDiscovery               types.Bool   `tfsdk:"no_gre_path_mtu_discovery"`
			IpipPathMtuDiscovery                types.Bool   `tfsdk:"ipip_path_mtu_discovery"`
			NoIpipPathMtuDiscovery              types.Bool   `tfsdk:"no_ipip_path_mtu_discovery"`
			IPv6DuplicateAddrDetectionTransmits types.Int64  `tfsdk:"ipv6_duplicate_addr_detection_transmits"`
			IPv6PathMtuDiscovery                types.Bool   `tfsdk:"ipv6_path_mtu_discovery"`
			NoIPv6PathMtuDiscovery              types.Bool   `tfsdk:"no_ipv6_path_mtu_discovery"`
			IPv6PathMtuDiscoveryTimeout         types.Int64  `tfsdk:"ipv6_path_mtu_discovery_timeout"`
			IPv6RejectZeroHopLimit              types.Bool   `tfsdk:"ipv6_reject_zero_hop_limit"`
			NoIPv6RejectZeroHopLimit            types.Bool   `tfsdk:"no_ipv6_reject_zero_hop_limit"`
			NoTCPReset                          types.String `tfsdk:"no_tcp_reset"`
			NoTCPRFC1323                        types.Bool   `tfsdk:"no_tcp_rfc1323"`
			NoTCPRFC1323Paws                    types.Bool   `tfsdk:"no_tcp_rfc1323_paws"`
			PathMtuDiscovery                    types.Bool   `tfsdk:"path_mtu_discovery"`
			NoPathMtuDiscovery                  types.Bool   `tfsdk:"no_path_mtu_discovery"`
			SourcePortUpperLimit                types.Int64  `tfsdk:"source_port_upper_limit"`
			SourceQuench                        types.Bool   `tfsdk:"source_quench"`
			NoSourceQuench                      types.Bool   `tfsdk:"no_source_quench"`
			TCPDropSynfinSet                    types.Bool   `tfsdk:"tcp_drop_synfin_set"`
			TCPMss                              types.Int64  `tfsdk:"tcp_mss"`
			IcmpV4RateLimit                     []struct {
				BucketSize types.Int64 `tfsdk:"bucket_size"`
				PacketRate types.Int64 `tfsdk:"packet_rate"`
			} `tfsdk:"icmpv4_rate_limit"`
			IcmpV6RateLimit []struct {
				BucketSize types.Int64 `tfsdk:"bucket_size"`
				PacketRate types.Int64 `tfsdk:"packet_rate"`
			} `tfsdk:"icmpv6_rate_limit"`
		} `tfsdk:"internet_options"`
		License []struct {
			Autoupdate            types.Bool   `tfsdk:"autoupdate"`
			AutoupdatePassword    types.String `tfsdk:"autoupdate_password"`
			AutoupdateURL         types.String `tfsdk:"autoupdate_url"`
			RenewBeforeExpiration types.Int64  `tfsdk:"renew_before_expiration"`
			RenewInterval         types.Int64  `tfsdk:"renew_interval"`
		} `tfsdk:"license"`
		Login []struct {
			Announcement       types.String   `tfsdk:"announcement"`
			DenySourcesAddress []types.String `tfsdk:"deny_sources_address"`
			IdleTimeout        types.Int64    `tfsdk:"idle_timeout"`
			Message            types.String   `tfsdk:"message"`
			Password           []struct {
				ChangeType              types.String `tfsdk:"change_type"`
				Format                  types.String `tfsdk:"format"`
				MaximumLength           types.Int64  `tfsdk:"maximum_length"`
				MinimumChanges          types.Int64  `tfsdk:"minimum_changes"`
				MinimumCharacterChanges types.Int64  `tfsdk:"minimum_character_changes"`
				MinimumLength           types.Int64  `tfsdk:"minimum_length"`
				MinimumLowerCases       types.Int64  `tfsdk:"minimum_lower_cases"`
				MinimumNumerics         types.Int64  `tfsdk:"minimum_numerics"`
				MinimumPunctuations     types.Int64  `tfsdk:"minimum_punctuations"`
				MinimumReuse            types.Int64  `tfsdk:"minimum_reuse"`
				MinimumUpperCases       types.Int64  `tfsdk:"minimum_upper_cases"`
			} `tfsdk:"password"`
			RetryOptions []struct {
				BackoffFactor         types.Int64 `tfsdk:"backoff_factor"`
				BackoffThreshold      types.Int64 `tfsdk:"backoff_threshold"`
				LockoutPeriod         types.Int64 `tfsdk:"lockout_period"`
				MaximumTime           types.Int64 `tfsdk:"maximum_time"`
				MinimumTime           types.Int64 `tfsdk:"minimum_time"`
				TriesBeforeDisconnect types.Int64 `tfsdk:"tries_before_disconnect"`
			} `tfsdk:"retry_options"`
		} `tfsdk:"login"`
		Ntp []struct {
			BootServer             types.String `tfsdk:"boot_server"`
			BroadcastClient        types.Bool   `tfsdk:"broadcast_client"`
			IntervalRange          types.Int64  `tfsdk:"interval_range"`
			MulticastClient        types.Bool   `tfsdk:"multicast_client"`
			MulticastClientAddress types.String `tfsdk:"multicast_client_address"`
			ThresholdAction        types.String `tfsdk:"threshold_action"`
			ThresholdValue         types.Int64  `tfsdk:"threshold_value"`
		} `tfsdk:"ntp"`
		Ports []struct {
			AuxiliaryAuthenticationOrder []types.String `tfsdk:"auxiliary_authentication_order"`
			AuxiliaryDisable             types.Bool     `tfsdk:"auxiliary_disable"`
			AuxiliaryInsecure            types.Bool     `tfsdk:"auxiliary_insecure"`
			AuxiliaryLogoutOnDisconnect  types.Bool     `tfsdk:"auxiliary_logout_on_disconnect"`
			AuxiliaryType                types.String   `tfsdk:"auxiliary_type"`
			ConsoleAuthenticationOrder   []types.String `tfsdk:"console_authentication_order"`
			ConsoleDisable               types.Bool     `tfsdk:"console_disable"`
			ConsoleInsecure              types.Bool     `tfsdk:"console_insecure"`
			ConsoleLogoutOnDisconnect    types.Bool     `tfsdk:"console_logout_on_disconnect"`
			ConsoleType                  types.String   `tfsdk:"console_type"`
		} `tfsdk:"ports"`
		Services []struct {
			NetconfSSH []struct {
				ClientAliveCountMax types.Int64 `tfsdk:"client_alive_count_max"`
				ClientAliveInterval types.Int64 `tfsdk:"client_alive_interval"`
				ConnectionLimit     types.Int64 `tfsdk:"connection_limit"`
				RateLimit           types.Int64 `tfsdk:"rate_limit"`
			} `tfsdk:"netconf_ssh"`
			NetconfTraceoptions []struct {
				FileName            types.String   `tfsdk:"file_name"`
				FileFiles           types.Int64    `tfsdk:"file_files"`
				FileMatch           types.String   `tfsdk:"file_match"`
				FileSize            types.Int64    `tfsdk:"file_size"`
				FileWorldReadable   types.Bool     `tfsdk:"file_world_readable"`
				FileNoWorldReadable types.Bool     `tfsdk:"file_no_world_readable"`
				Flag                []types.String `tfsdk:"flag"`
				NoRemoteTrace       types.Bool     `tfsdk:"no_remote_trace"`
				OnDemand            types.Bool     `tfsdk:"on_demand"`
			} `tfsdk:"netconf_traceoptions"`
			SSH []struct {
				AuthenticationOrder         []types.String `tfsdk:"authentication_order"`
				Ciphers                     []types.String `tfsdk:"ciphers"`
				ClientAliveCountMax         types.Int64    `tfsdk:"client_alive_count_max"`
				ClientAliveInterval         types.Int64    `tfsdk:"client_alive_interval"`
				ConnectionLimit             types.Int64    `tfsdk:"connection_limit"`
				FingerprintHash             types.String   `tfsdk:"fingerprint_hash"`
				HostkeyAlgorithm            []types.String `tfsdk:"hostkey_algorithm"`
				KeyExchange                 []types.String `tfsdk:"key_exchange"`
				LogKeyChanges               types.Bool     `tfsdk:"log_key_changes"`
				Macs                        []types.String `tfsdk:"macs"`
				MaxPreAuthenticationPackets types.Int64    `tfsdk:"max_pre_authentication_packets"`
				MaxSessionsPerConnection    types.Int64    `tfsdk:"max_sessions_per_connection"`
				NoPasswords                 types.Bool     `tfsdk:"no_passwords"`
				NoPublicKeys                types.Bool     `tfsdk:"no_public_keys"`
				Port                        types.Int64    `tfsdk:"port"`
				ProtocolVersion             []types.String `tfsdk:"protocol_version"`
				RateLimit                   types.Int64    `tfsdk:"rate_limit"`
				RootLogin                   types.String   `tfsdk:"root_login"`
				TCPForwarding               types.Bool     `tfsdk:"tcp_forwarding"`
				NoTCPForwarding             types.Bool     `tfsdk:"no_tcp_forwarding"`
			} `tfsdk:"ssh"`
			WebManagementHTTP []struct {
				Interface []types.String `tfsdk:"interface"`
				Port      types.Int64    `tfsdk:"port"`
			} `tfsdk:"web_management_http"`
			WebManagementHTTPS []struct {
				Interface                  []types.String `tfsdk:"interface"`
				LocalCertificate           types.String   `tfsdk:"local_certificate"`
				PkiLocalCertificate        types.String   `tfsdk:"pki_local_certificate"`
				Port                       types.Int64    `tfsdk:"port"`
				SystemGeneratedCertificate types.Bool     `tfsdk:"system_generated_certificate"`
			} `tfsdk:"web_management_https"`
		} `tfsdk:"services"`
		Syslog []struct {
			LogRotateFrequency    types.Int64  `tfsdk:"log_rotate_frequency"`
			SourceAddress         types.String `tfsdk:"source_address"`
			TimeFormatMillisecond types.Bool   `tfsdk:"time_format_millisecond"`
			TimeFormatYear        types.Bool   `tfsdk:"time_format_year"`
			Archive               []struct {
				BinaryData      types.Bool  `tfsdk:"binary_data"`
				NoBinaryData    types.Bool  `tfsdk:"no_binary_data"`
				Files           types.Int64 `tfsdk:"files"`
				Size            types.Int64 `tfsdk:"size"`
				WorldReadable   types.Bool  `tfsdk:"world_readable"`
				NoWorldReadable types.Bool  `tfsdk:"no_world_readable"`
			} `tfsdk:"archive"`
			Console []struct {
				AnySeverity                 types.String `tfsdk:"any_severity"`
				AuthorizationSeverity       types.String `tfsdk:"authorization_severity"`
				ChangelogSeverity           types.String `tfsdk:"changelog_severity"`
				ConflictlogSeverity         types.String `tfsdk:"conflictlog_severity"`
				DaemonSeverity              types.String `tfsdk:"daemon_severity"`
				DfcSeverity                 types.String `tfsdk:"dfc_severity"`
				ExternalSeverity            types.String `tfsdk:"external_severity"`
				FirewallSeverity            types.String `tfsdk:"firewall_severity"`
				FtpSeverity                 types.String `tfsdk:"ftp_severity"`
				InteractivecommandsSeverity types.String `tfsdk:"interactivecommands_severity"`
				KernelSeverity              types.String `tfsdk:"kernel_severity"`
				NtpSeverity                 types.String `tfsdk:"ntp_severity"`
				PfeSeverity                 types.String `tfsdk:"pfe_severity"`
				SecuritySeverity            types.String `tfsdk:"security_severity"`
				UserSeverity                types.String `tfsdk:"user_severity"`
			} `tfsdk:"console"`
		} `tfsdk:"syslog"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 systemData
	dataV1.ID = dataV0.ID
	dataV1.AutoSnapshot = dataV0.AutoSnapshot
	dataV1.DefaultAddressSelection = dataV0.DefaultAddressSelection
	dataV1.NoMulticastEcho = dataV0.NoMulticastEcho
	dataV1.NoPingRecordRoute = dataV0.NoPingRecordRoute
	dataV1.NoPingTimestamp = dataV0.NoPingTimestamp
	dataV1.NoRedirects = dataV0.NoRedirects
	dataV1.NoRedirectsIPv6 = dataV0.NoRedirectsIPv6
	dataV1.RadiusOptionsEnhancedAccounting = dataV0.RadiusOptionsEnhancedAccounting
	dataV1.RadiusOptionsPasswordProtocolMschapv2 = dataV0.RadiusOptionsPasswordProtocolMschapv2
	dataV1.AuthenticationOrder = dataV0.AuthenticationOrder
	dataV1.DomainName = dataV0.DomainName
	dataV1.HostName = dataV0.HostName
	dataV1.MaxConfigurationRollbacks = dataV0.MaxConfigurationRollbacks
	dataV1.MaxConfigurationsOnFlash = dataV0.MaxConfigurationsOnFlash
	dataV1.NameServer = dataV0.NameServer
	dataV1.RadiusOptionsAttributesNasIpaddress = dataV0.RadiusOptionsAttributesNasIpaddress
	dataV1.TimeZone = dataV0.TimeZone
	dataV1.TracingDestOverrideSyslogHost = dataV0.TracingDestOverrideSyslogHost
	if len(dataV0.ArchivalConfiguration) > 0 {
		dataV1.ArchivalConfiguration = &systemBlockArchivalConfiguration{
			TransferOnCommit: dataV0.ArchivalConfiguration[0].TransferOnCommit,
			TransferInterval: dataV0.ArchivalConfiguration[0].TransferInterval,
		}
		for _, block := range dataV0.ArchivalConfiguration[0].ArchiveSite {
			dataV1.ArchivalConfiguration.ArchiveSite = append(dataV1.ArchivalConfiguration.ArchiveSite,
				systemBlockArchivalConfigurationBlockArchiveSite{
					URL:      block.URL,
					Password: block.Password,
				},
			)
		}
	}
	if len(dataV0.Inet6BackupRouter) > 0 {
		dataV1.Inet6BackupRouter = &systemBlockInet6BackupRouter{
			Address:     dataV0.Inet6BackupRouter[0].Address,
			Destination: dataV0.Inet6BackupRouter[0].Destination,
		}
	}
	if len(dataV0.InternetOptions) > 0 {
		dataV1.InternetOptions = &systemBlockInternetOptions{
			GrePathMtuDiscovery:                 dataV0.InternetOptions[0].GrePathMtuDiscovery,
			NoGrePathMtuDiscovery:               dataV0.InternetOptions[0].NoGrePathMtuDiscovery,
			IpipPathMtuDiscovery:                dataV0.InternetOptions[0].IpipPathMtuDiscovery,
			NoIpipPathMtuDiscovery:              dataV0.InternetOptions[0].NoIpipPathMtuDiscovery,
			IPv6PathMtuDiscovery:                dataV0.InternetOptions[0].IPv6PathMtuDiscovery,
			NoIPv6PathMtuDiscovery:              dataV0.InternetOptions[0].NoIPv6PathMtuDiscovery,
			IPv6RejectZeroHopLimit:              dataV0.InternetOptions[0].IPv6RejectZeroHopLimit,
			NoIPv6RejectZeroHopLimit:            dataV0.InternetOptions[0].NoIPv6RejectZeroHopLimit,
			NoTCPRFC1323:                        dataV0.InternetOptions[0].NoTCPRFC1323,
			NoTCPRFC1323Paws:                    dataV0.InternetOptions[0].NoTCPRFC1323Paws,
			PathMtuDiscovery:                    dataV0.InternetOptions[0].PathMtuDiscovery,
			NoPathMtuDiscovery:                  dataV0.InternetOptions[0].NoPathMtuDiscovery,
			SourceQuench:                        dataV0.InternetOptions[0].SourceQuench,
			NoSourceQuench:                      dataV0.InternetOptions[0].NoSourceQuench,
			TCPDropSynfinSet:                    dataV0.InternetOptions[0].TCPDropSynfinSet,
			IPv6DuplicateAddrDetectionTransmits: dataV0.InternetOptions[0].IPv6DuplicateAddrDetectionTransmits,
			IPv6PathMtuDiscoveryTimeout:         dataV0.InternetOptions[0].IPv6PathMtuDiscoveryTimeout,
			NoTCPReset:                          dataV0.InternetOptions[0].NoTCPReset,
			SourcePortUpperLimit:                dataV0.InternetOptions[0].SourcePortUpperLimit,
			TCPMss:                              dataV0.InternetOptions[0].TCPMss,
		}
		if len(dataV0.InternetOptions[0].IcmpV4RateLimit) > 0 {
			dataV1.InternetOptions.IcmpV4RateLimit = &systemBlockInternetOptionsBlockIcmpRateLimit{
				BucketSize: dataV0.InternetOptions[0].IcmpV4RateLimit[0].BucketSize,
				PacketRate: dataV0.InternetOptions[0].IcmpV4RateLimit[0].PacketRate,
			}
		}
		if len(dataV0.InternetOptions[0].IcmpV6RateLimit) > 0 {
			dataV1.InternetOptions.IcmpV6RateLimit = &systemBlockInternetOptionsBlockIcmpRateLimit{
				BucketSize: dataV0.InternetOptions[0].IcmpV6RateLimit[0].BucketSize,
				PacketRate: dataV0.InternetOptions[0].IcmpV6RateLimit[0].PacketRate,
			}
		}
	}
	if len(dataV0.License) > 0 {
		dataV1.License = &systemBlockLicense{
			Autoupdate:            dataV0.License[0].Autoupdate,
			AutoupdatePassword:    dataV0.License[0].AutoupdatePassword,
			AutoupdateURL:         dataV0.License[0].AutoupdateURL,
			RenewBeforeExpiration: dataV0.License[0].RenewBeforeExpiration,
			RenewInterval:         dataV0.License[0].RenewInterval,
		}
	}
	if len(dataV0.Login) > 0 {
		dataV1.Login = &systemBlockLogin{
			Announcement:       dataV0.Login[0].Announcement,
			DenySourcesAddress: dataV0.Login[0].DenySourcesAddress,
			IdleTimeout:        dataV0.Login[0].IdleTimeout,
			Message:            dataV0.Login[0].Message,
		}
		if len(dataV0.Login[0].Password) > 0 {
			dataV1.Login.Password = &systemBlockLoginBlockPassword{
				ChangeType:              dataV0.Login[0].Password[0].ChangeType,
				Format:                  dataV0.Login[0].Password[0].Format,
				MaximumLength:           dataV0.Login[0].Password[0].MaximumLength,
				MinimumChanges:          dataV0.Login[0].Password[0].MinimumChanges,
				MinimumCharacterChanges: dataV0.Login[0].Password[0].MinimumCharacterChanges,
				MinimumLength:           dataV0.Login[0].Password[0].MinimumLength,
				MinimumLowerCases:       dataV0.Login[0].Password[0].MinimumLowerCases,
				MinimumNumerics:         dataV0.Login[0].Password[0].MinimumNumerics,
				MinimumPunctuations:     dataV0.Login[0].Password[0].MinimumPunctuations,
				MinimumReuse:            dataV0.Login[0].Password[0].MinimumReuse,
				MinimumUpperCases:       dataV0.Login[0].Password[0].MinimumUpperCases,
			}
		}
		if len(dataV0.Login[0].RetryOptions) > 0 {
			dataV1.Login.RetryOptions = &systemBlockLoginBlockRetryOptions{
				BackoffFactor:         dataV0.Login[0].RetryOptions[0].BackoffFactor,
				BackoffThreshold:      dataV0.Login[0].RetryOptions[0].BackoffThreshold,
				LockoutPeriod:         dataV0.Login[0].RetryOptions[0].LockoutPeriod,
				MaximumTime:           dataV0.Login[0].RetryOptions[0].MaximumTime,
				MinimumTime:           dataV0.Login[0].RetryOptions[0].MinimumTime,
				TriesBeforeDisconnect: dataV0.Login[0].RetryOptions[0].TriesBeforeDisconnect,
			}
		}
	}
	if len(dataV0.Ntp) > 0 {
		dataV1.Ntp = &systemBlockNtp{
			BroadcastClient:        dataV0.Ntp[0].BroadcastClient,
			MulticastClient:        dataV0.Ntp[0].MulticastClient,
			BootServer:             dataV0.Ntp[0].BootServer,
			IntervalRange:          dataV0.Ntp[0].IntervalRange,
			MulticastClientAddress: dataV0.Ntp[0].MulticastClientAddress,
			ThresholdAction:        dataV0.Ntp[0].ThresholdAction,
			ThresholdValue:         dataV0.Ntp[0].ThresholdValue,
		}
	}
	if len(dataV0.Ports) > 0 {
		dataV1.Ports = &systemBlockPorts{
			AuxiliaryDisable:             dataV0.Ports[0].AuxiliaryDisable,
			AuxiliaryInsecure:            dataV0.Ports[0].AuxiliaryInsecure,
			AuxiliaryLogoutOnDisconnect:  dataV0.Ports[0].AuxiliaryLogoutOnDisconnect,
			ConsoleDisable:               dataV0.Ports[0].ConsoleDisable,
			ConsoleInsecure:              dataV0.Ports[0].ConsoleInsecure,
			ConsoleLogoutOnDisconnect:    dataV0.Ports[0].ConsoleLogoutOnDisconnect,
			AuxiliaryAuthenticationOrder: dataV0.Ports[0].AuxiliaryAuthenticationOrder,
			AuxiliaryType:                dataV0.Ports[0].AuxiliaryType,
			ConsoleAuthenticationOrder:   dataV0.Ports[0].ConsoleAuthenticationOrder,
			ConsoleType:                  dataV0.Ports[0].ConsoleType,
		}
	}
	if len(dataV0.Services) > 0 {
		dataV1.Services = &systemBlockServices{}
		if len(dataV0.Services[0].NetconfSSH) > 0 {
			dataV1.Services.NetconfSSH = &systemBlockServicesBlockNetconfSSH{
				ClientAliveCountMax: dataV0.Services[0].NetconfSSH[0].ClientAliveCountMax,
				ClientAliveInterval: dataV0.Services[0].NetconfSSH[0].ClientAliveInterval,
				ConnectionLimit:     dataV0.Services[0].NetconfSSH[0].ConnectionLimit,
				RateLimit:           dataV0.Services[0].NetconfSSH[0].RateLimit,
			}
		}
		if len(dataV0.Services[0].NetconfTraceoptions) > 0 {
			dataV1.Services.NetconfTraceoptions = &systemBlockServicesBlockNetconfTraceoptions{
				FileWorldReadable:   dataV0.Services[0].NetconfTraceoptions[0].FileWorldReadable,
				FileNoWorldReadable: dataV0.Services[0].NetconfTraceoptions[0].FileNoWorldReadable,
				NoRemoteTrace:       dataV0.Services[0].NetconfTraceoptions[0].NoRemoteTrace,
				OnDemand:            dataV0.Services[0].NetconfTraceoptions[0].OnDemand,
				FileName:            dataV0.Services[0].NetconfTraceoptions[0].FileName,
				FileFiles:           dataV0.Services[0].NetconfTraceoptions[0].FileFiles,
				FileMatch:           dataV0.Services[0].NetconfTraceoptions[0].FileMatch,
				FileSize:            dataV0.Services[0].NetconfTraceoptions[0].FileSize,
				Flag:                dataV0.Services[0].NetconfTraceoptions[0].Flag,
			}
		}
		if len(dataV0.Services[0].SSH) > 0 {
			dataV1.Services.SSH = &systemBlockServicesBlockSSH{
				LogKeyChanges:               dataV0.Services[0].SSH[0].LogKeyChanges,
				NoPasswords:                 dataV0.Services[0].SSH[0].NoPasswords,
				NoPublicKeys:                dataV0.Services[0].SSH[0].NoPublicKeys,
				TCPForwarding:               dataV0.Services[0].SSH[0].TCPForwarding,
				NoTCPForwarding:             dataV0.Services[0].SSH[0].NoTCPForwarding,
				AuthenticationOrder:         dataV0.Services[0].SSH[0].AuthenticationOrder,
				Ciphers:                     dataV0.Services[0].SSH[0].Ciphers,
				ClientAliveCountMax:         dataV0.Services[0].SSH[0].ClientAliveCountMax,
				ClientAliveInterval:         dataV0.Services[0].SSH[0].ClientAliveInterval,
				ConnectionLimit:             dataV0.Services[0].SSH[0].ConnectionLimit,
				FingerprintHash:             dataV0.Services[0].SSH[0].FingerprintHash,
				HostkeyAlgorithm:            dataV0.Services[0].SSH[0].HostkeyAlgorithm,
				KeyExchange:                 dataV0.Services[0].SSH[0].KeyExchange,
				Macs:                        dataV0.Services[0].SSH[0].Macs,
				MaxPreAuthenticationPackets: dataV0.Services[0].SSH[0].MaxPreAuthenticationPackets,
				MaxSessionsPerConnection:    dataV0.Services[0].SSH[0].MaxSessionsPerConnection,
				Port:                        dataV0.Services[0].SSH[0].Port,
				ProtocolVersion:             dataV0.Services[0].SSH[0].ProtocolVersion,
				RateLimit:                   dataV0.Services[0].SSH[0].RateLimit,
				RootLogin:                   dataV0.Services[0].SSH[0].RootLogin,
			}
		}
		if len(dataV0.Services[0].WebManagementHTTP) > 0 {
			dataV1.Services.WebManagementHTTP = &systemBlockServicesBlockWebManagementHTTP{
				Interface: dataV0.Services[0].WebManagementHTTP[0].Interface,
				Port:      dataV0.Services[0].WebManagementHTTP[0].Port,
			}
		}
		if len(dataV0.Services[0].WebManagementHTTPS) > 0 {
			dataV1.Services.WebManagementHTTPS = &systemBlockServicesBlockWebManagementHTTPS{
				SystemGeneratedCertificate: dataV0.Services[0].WebManagementHTTPS[0].SystemGeneratedCertificate,
				Interface:                  dataV0.Services[0].WebManagementHTTPS[0].Interface,
				LocalCertificate:           dataV0.Services[0].WebManagementHTTPS[0].LocalCertificate,
				PkiLocalCertificate:        dataV0.Services[0].WebManagementHTTPS[0].PkiLocalCertificate,
				Port:                       dataV0.Services[0].WebManagementHTTPS[0].Port,
			}
		}
	}
	if len(dataV0.Syslog) > 0 {
		dataV1.Syslog = &systemBlockSyslog{
			TimeFormatMillisecond: dataV0.Syslog[0].TimeFormatMillisecond,
			TimeFormatYear:        dataV0.Syslog[0].TimeFormatYear,
			LogRotateFrequency:    dataV0.Syslog[0].LogRotateFrequency,
			SourceAddress:         dataV0.Syslog[0].SourceAddress,
		}
		if len(dataV0.Syslog[0].Archive) > 0 {
			dataV1.Syslog.Archive = &systemBlockSyslogBlockArchive{
				BinaryData:      dataV0.Syslog[0].Archive[0].BinaryData,
				NoBinaryData:    dataV0.Syslog[0].Archive[0].NoBinaryData,
				WorldReadable:   dataV0.Syslog[0].Archive[0].WorldReadable,
				NoWorldReadable: dataV0.Syslog[0].Archive[0].NoWorldReadable,
				Files:           dataV0.Syslog[0].Archive[0].Files,
				Size:            dataV0.Syslog[0].Archive[0].Size,
			}
		}
		if len(dataV0.Syslog[0].Console) > 0 {
			dataV1.Syslog.Console = &systemBlockSyslogBlockConsole{
				AnySeverity:                 dataV0.Syslog[0].Console[0].AnySeverity,
				AuthorizationSeverity:       dataV0.Syslog[0].Console[0].AuthorizationSeverity,
				ChangelogSeverity:           dataV0.Syslog[0].Console[0].ChangelogSeverity,
				ConflictlogSeverity:         dataV0.Syslog[0].Console[0].ConflictlogSeverity,
				DaemonSeverity:              dataV0.Syslog[0].Console[0].DaemonSeverity,
				DfcSeverity:                 dataV0.Syslog[0].Console[0].DfcSeverity,
				ExternalSeverity:            dataV0.Syslog[0].Console[0].ExternalSeverity,
				FirewallSeverity:            dataV0.Syslog[0].Console[0].FirewallSeverity,
				FtpSeverity:                 dataV0.Syslog[0].Console[0].FtpSeverity,
				InteractivecommandsSeverity: dataV0.Syslog[0].Console[0].InteractivecommandsSeverity,
				KernelSeverity:              dataV0.Syslog[0].Console[0].KernelSeverity,
				NtpSeverity:                 dataV0.Syslog[0].Console[0].NtpSeverity,
				PfeSeverity:                 dataV0.Syslog[0].Console[0].PfeSeverity,
				SecuritySeverity:            dataV0.Syslog[0].Console[0].SecuritySeverity,
				UserSeverity:                dataV0.Syslog[0].Console[0].UserSeverity,
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
