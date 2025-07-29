package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *services) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"advanced_anti_malware": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"connection": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"auth_tls_profile": schema.StringAttribute{
												Optional: true,
												Computed: true,
											},
											"proxy_profile": schema.StringAttribute{
												Optional: true,
											},
											"source_address": schema.StringAttribute{
												Optional: true,
											},
											"source_interface": schema.StringAttribute{
												Optional: true,
											},
											"url": schema.StringAttribute{
												Optional: true,
												Computed: true,
											},
										},
									},
								},
								"default_policy": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"blacklist_notification_log": schema.BoolAttribute{
												Optional: true,
											},
											"default_notification_log": schema.BoolAttribute{
												Optional: true,
											},
											"fallback_options_action": schema.StringAttribute{
												Optional: true,
											},
											"fallback_options_notification_log": schema.BoolAttribute{
												Optional: true,
											},
											"http_action": schema.StringAttribute{
												Optional: true,
											},
											"http_client_notify_file": schema.StringAttribute{
												Optional: true,
											},
											"http_client_notify_message": schema.StringAttribute{
												Optional: true,
											},
											"http_client_notify_redirect_url": schema.StringAttribute{
												Optional: true,
											},
											"http_file_verdict_unknown": schema.StringAttribute{
												Optional: true,
											},
											"http_inspection_profile": schema.StringAttribute{
												Optional: true,
											},
											"http_notification_log": schema.BoolAttribute{
												Optional: true,
											},
											"imap_inspection_profile": schema.StringAttribute{
												Optional: true,
											},
											"imap_notification_log": schema.BoolAttribute{
												Optional: true,
											},
											"smtp_inspection_profile": schema.StringAttribute{
												Optional: true,
											},
											"smtp_notification_log": schema.BoolAttribute{
												Optional: true,
											},
											"verdict_threshold": schema.StringAttribute{
												Optional: true,
											},
											"whitelist_notification_log": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"application_identification": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"application_system_cache_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"global_offload_byte_limit": schema.Int64Attribute{
									Optional: true,
								},
								"imap_cache_size": schema.Int64Attribute{
									Optional: true,
								},
								"imap_cache_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"max_memory": schema.Int64Attribute{
									Optional: true,
								},
								"max_transactions": schema.Int64Attribute{
									Optional: true,
								},
								"micro_apps": schema.BoolAttribute{
									Optional: true,
								},
								"no_application_system_cache": schema.BoolAttribute{
									Optional: true,
								},
								"statistics_interval": schema.Int64Attribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"application_system_cache": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"no_miscellaneous_services": schema.BoolAttribute{
												Optional: true,
											},
											"security_services": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"download": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"automatic_interval": schema.Int64Attribute{
												Optional: true,
											},
											"automatic_start_time": schema.StringAttribute{
												Optional: true,
											},
											"ignore_server_validation": schema.BoolAttribute{
												Optional: true,
											},
											"proxy_profile": schema.StringAttribute{
												Optional: true,
											},
											"url": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"enable_performance_mode": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"max_packet_threshold": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"inspection_limit_tcp": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"byte_limit": schema.Int64Attribute{
												Optional: true,
											},
											"packet_limit": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"inspection_limit_udp": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"byte_limit": schema.Int64Attribute{
												Optional: true,
											},
											"packet_limit": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"security_intelligence": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"authentication_tls_profile": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"authentication_token": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"category_disable": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"proxy_profile": schema.StringAttribute{
									Optional: true,
								},
								"url": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"url_parameter": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
							},
							Blocks: map[string]schema.Block{
								"default_policy": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"category_name": schema.StringAttribute{
												Required: true,
											},
											"profile_name": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
							},
						},
					},
					"user_identification": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"device_info_auth_source": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"ad_access": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"auth_entry_timeout": schema.Int64Attribute{
												Optional: true,
											},
											"filter_exclude": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"filter_include": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"firewall_auth_forced_timeout": schema.Int64Attribute{
												Optional: true,
											},
											"invalid_auth_entry_timeout": schema.Int64Attribute{
												Optional: true,
											},
											"no_on_demand_probe": schema.BoolAttribute{
												Optional: true,
											},
											"wmi_timeout": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"identity_management": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"authentication_entry_timeout": schema.Int64Attribute{
												Optional: true,
											},
											"batch_query_interval": schema.Int64Attribute{
												Optional: true,
											},
											"batch_query_items_per_batch": schema.Int64Attribute{
												Optional: true,
											},
											"filter_domain": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"filter_exclude_ip_address_book": schema.StringAttribute{
												Optional: true,
											},
											"filter_exclude_ip_address_set": schema.StringAttribute{
												Optional: true,
											},
											"filter_include_ip_address_book": schema.StringAttribute{
												Optional: true,
											},
											"filter_include_ip_address_set": schema.StringAttribute{
												Optional: true,
											},
											"invalid_authentication_entry_timeout": schema.Int64Attribute{
												Optional: true,
											},
											"ip_query_disable": schema.BoolAttribute{
												Optional: true,
											},
											"ip_query_delay_time": schema.Int64Attribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"connection": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"primary_address": schema.StringAttribute{
															Required: true,
														},
														"primary_client_id": schema.StringAttribute{
															Required: true,
														},
														"primary_client_secret": schema.StringAttribute{
															Required:  true,
															Sensitive: true,
														},
														"connect_method": schema.StringAttribute{
															Optional: true,
														},
														"port": schema.Int64Attribute{
															Optional: true,
														},
														"primary_ca_certificate": schema.StringAttribute{
															Optional: true,
														},
														"query_api": schema.StringAttribute{
															Optional: true,
														},
														"secondary_address": schema.StringAttribute{
															Optional: true,
														},
														"secondary_ca_certificate": schema.StringAttribute{
															Optional: true,
														},
														"secondary_client_id": schema.StringAttribute{
															Optional: true,
														},
														"secondary_client_secret": schema.StringAttribute{
															Optional:  true,
															Sensitive: true,
														},
														"token_api": schema.StringAttribute{
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
				},
			},
			StateUpgrader: upgradeServicesV0toV1,
		},
	}
}

//nolint:lll
func upgradeServicesV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                  types.String `tfsdk:"id"`
		CleanOnDestroy      types.Bool   `tfsdk:"clean_on_destroy"`
		AdvancedAntiMalware []struct {
			Connection []struct {
				AuthTLSProfile  types.String `tfsdk:"auth_tls_profile"`
				ProxyProfile    types.String `tfsdk:"proxy_profile"`
				SourceAddress   types.String `tfsdk:"source_address"`
				SourceInterface types.String `tfsdk:"source_interface"`
				URL             types.String `tfsdk:"url"`
			} `tfsdk:"connection"`
			DefaultPolicy []struct {
				BlacklistNotificationLog       types.Bool   `tfsdk:"blacklist_notification_log"`
				DefaultNotificationLog         types.Bool   `tfsdk:"default_notification_log"`
				FallbackOptionsAction          types.String `tfsdk:"fallback_options_action"`
				FallbackOptionsNotificationLog types.Bool   `tfsdk:"fallback_options_notification_log"`
				HTTPAction                     types.String `tfsdk:"http_action"`
				HTTPClientNotifyFile           types.String `tfsdk:"http_client_notify_file"`
				HTTPClientNotifyMessage        types.String `tfsdk:"http_client_notify_message"`
				HTTPClientNotifyRedirectURL    types.String `tfsdk:"http_client_notify_redirect_url"`
				HTTPFileVerdictUnknown         types.String `tfsdk:"http_file_verdict_unknown"`
				HTTPInspectionProfile          types.String `tfsdk:"http_inspection_profile"`
				HTTPNotificationLog            types.Bool   `tfsdk:"http_notification_log"`
				IMAPInspectionProfile          types.String `tfsdk:"imap_inspection_profile"`
				IMAPNotificationLog            types.Bool   `tfsdk:"imap_notification_log"`
				SMTPInspectionProfile          types.String `tfsdk:"smtp_inspection_profile"`
				SMTPNotificationLog            types.Bool   `tfsdk:"smtp_notification_log"`
				VerdictThreshold               types.String `tfsdk:"verdict_threshold"`
				WhitelistNotificationLog       types.Bool   `tfsdk:"whitelist_notification_log"`
			} `tfsdk:"default_policy"`
		} `tfsdk:"advanced_anti_malware"`
		ApplicationIdentification []struct {
			ApplicationSystemCacheTimeout types.Int64 `tfsdk:"application_system_cache_timeout"`
			GlobalOffloadByteLimit        types.Int64 `tfsdk:"global_offload_byte_limit"`
			IMAPCacheSize                 types.Int64 `tfsdk:"imap_cache_size"`
			IMAPCacheTimeout              types.Int64 `tfsdk:"imap_cache_timeout"`
			MaxMemory                     types.Int64 `tfsdk:"max_memory"`
			MaxTransactions               types.Int64 `tfsdk:"max_transactions"`
			MicroApps                     types.Bool  `tfsdk:"micro_apps"`
			NoApplicationSystemCache      types.Bool  `tfsdk:"no_application_system_cache"`
			StatisticsInterval            types.Int64 `tfsdk:"statistics_interval"`
			ApplicationSystemCache        []struct {
				NoMiscellaneousServices types.Bool `tfsdk:"no_miscellaneous_services"`
				SecurityServices        types.Bool `tfsdk:"security_services"`
			} `tfsdk:"application_system_cache"`
			Download []struct {
				AutomaticInterval      types.Int64  `tfsdk:"automatic_interval"`
				AutomaticStartTime     types.String `tfsdk:"automatic_start_time"`
				IgnoreServerValidation types.Bool   `tfsdk:"ignore_server_validation"`
				ProxyProfile           types.String `tfsdk:"proxy_profile"`
				URL                    types.String `tfsdk:"url"`
			} `tfsdk:"download"`
			EnablePerformanceMode []struct {
				MaxPacketThreshold types.Int64 `tfsdk:"max_packet_threshold"`
			} `tfsdk:"enable_performance_mode"`
			InspectionLimitTCP []struct {
				ByteLimit   types.Int64 `tfsdk:"byte_limit"`
				PacketLimit types.Int64 `tfsdk:"packet_limit"`
			} `tfsdk:"inspection_limit_tcp"`
			InspectionLimitUDP []struct {
				ByteLimit   types.Int64 `tfsdk:"byte_limit"`
				PacketLimit types.Int64 `tfsdk:"packet_limit"`
			} `tfsdk:"inspection_limit_udp"`
		} `tfsdk:"application_identification"`
		SecurityIntelligence []struct {
			AuthenticationTLSProfile types.String   `tfsdk:"authentication_tls_profile"`
			AuthenticationToken      types.String   `tfsdk:"authentication_token"`
			CategoryDisable          []types.String `tfsdk:"category_disable"`
			ProxyProfile             types.String   `tfsdk:"proxy_profile"`
			URL                      types.String   `tfsdk:"url"`
			URLParameter             types.String   `tfsdk:"url_parameter"`
			DefaultPolicy            []struct {
				CategoryName types.String `tfsdk:"category_name"`
				ProfileName  types.String `tfsdk:"profile_name"`
			} `tfsdk:"default_policy"`
		} `tfsdk:"security_intelligence"`
		UserIdentification []struct {
			DeviceInfoAuthSource types.String `tfsdk:"device_info_auth_source"`
			ADAccess             []struct {
				AuthEntryTimeout          types.Int64    `tfsdk:"auth_entry_timeout"`
				FilterExclude             []types.String `tfsdk:"filter_exclude"`
				FilterInclude             []types.String `tfsdk:"filter_include"`
				FirewallAuthForcedTimeout types.Int64    `tfsdk:"firewall_auth_forced_timeout"`
				InvalidAuthEntryTimeout   types.Int64    `tfsdk:"invalid_auth_entry_timeout"`
				NoOnDemandProbe           types.Bool     `tfsdk:"no_on_demand_probe"`
				WmiTimeout                types.Int64    `tfsdk:"wmi_timeout"`
			} `tfsdk:"ad_access"`
			IdentityManagement []struct {
				AuthenticationEntryTimeout        types.Int64    `tfsdk:"authentication_entry_timeout"`
				BatchQueryInterval                types.Int64    `tfsdk:"batch_query_interval"`
				BatchQueryItemsPerBatch           types.Int64    `tfsdk:"batch_query_items_per_batch"`
				FilterDomain                      []types.String `tfsdk:"filter_domain"`
				FilterExcludeIPAddressBook        types.String   `tfsdk:"filter_exclude_ip_address_book"`
				FilterExcludeIPAddressSet         types.String   `tfsdk:"filter_exclude_ip_address_set"`
				FilterIncludeIPAddressBook        types.String   `tfsdk:"filter_include_ip_address_book"`
				FilterIncludeIPAddressSet         types.String   `tfsdk:"filter_include_ip_address_set"`
				InvalidAuthenticationEntryTimeout types.Int64    `tfsdk:"invalid_authentication_entry_timeout"`
				IPQueryDisable                    types.Bool     `tfsdk:"ip_query_disable"`
				IPQueryDelayTime                  types.Int64    `tfsdk:"ip_query_delay_time"`
				Connection                        []struct {
					PrimaryAddress         types.String `tfsdk:"primary_address"`
					PrimaryClientID        types.String `tfsdk:"primary_client_id"`
					PrimaryClientSecret    types.String `tfsdk:"primary_client_secret"`
					ConnectMethod          types.String `tfsdk:"connect_method"`
					Port                   types.Int64  `tfsdk:"port"`
					PrimaryCACertificate   types.String `tfsdk:"primary_ca_certificate"`
					QueryAPI               types.String `tfsdk:"query_api"`
					SecondaryAddress       types.String `tfsdk:"secondary_address"`
					SecondaryCACertificate types.String `tfsdk:"secondary_ca_certificate"`
					SecondaryClientID      types.String `tfsdk:"secondary_client_id"`
					SecondaryClientSecret  types.String `tfsdk:"secondary_client_secret"`
					TokenAPI               types.String `tfsdk:"token_api"`
				} `tfsdk:"connection"`
			} `tfsdk:"identity_management"`
		} `tfsdk:"user_identification"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 servicesData
	dataV1.ID = dataV0.ID
	dataV1.CleanOnDestroy = dataV0.CleanOnDestroy
	if !dataV1.CleanOnDestroy.IsNull() && !dataV1.CleanOnDestroy.ValueBool() {
		dataV1.CleanOnDestroy = types.BoolNull()
	}
	if len(dataV0.AdvancedAntiMalware) > 0 {
		dataV1.AdvancedAntiMalware = &servicesBlockAdvancedAntiMalware{}
		if len(dataV0.AdvancedAntiMalware[0].Connection) > 0 {
			dataV1.AdvancedAntiMalware.Connection = &servicesBlockAdvancedAntiMalwareBlockConnection{
				AuthTLSProfile:  dataV0.AdvancedAntiMalware[0].Connection[0].AuthTLSProfile,
				ProxyProfile:    dataV0.AdvancedAntiMalware[0].Connection[0].ProxyProfile,
				SourceAddress:   dataV0.AdvancedAntiMalware[0].Connection[0].SourceAddress,
				SourceInterface: dataV0.AdvancedAntiMalware[0].Connection[0].SourceInterface,
				URL:             dataV0.AdvancedAntiMalware[0].Connection[0].URL,
			}
		}
		if len(dataV0.AdvancedAntiMalware[0].DefaultPolicy) > 0 {
			dataV1.AdvancedAntiMalware.DefaultPolicy = &servicesBlockAdvancedAntiMalwareBlockDefaultPolicy{
				BlacklistNotificationLog:       dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].BlacklistNotificationLog,
				DefaultNotificationLog:         dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].DefaultNotificationLog,
				FallbackOptionsAction:          dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].FallbackOptionsAction,
				FallbackOptionsNotificationLog: dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].FallbackOptionsNotificationLog,
				HTTPAction:                     dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].HTTPAction,
				HTTPClientNotifyFile:           dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].HTTPClientNotifyFile,
				HTTPClientNotifyMessage:        dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].HTTPClientNotifyMessage,
				HTTPClientNotifyRedirectURL:    dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].HTTPClientNotifyRedirectURL,
				HTTPFileVerdictUnknown:         dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].HTTPFileVerdictUnknown,
				HTTPInspectionProfile:          dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].HTTPInspectionProfile,
				HTTPNotificationLog:            dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].HTTPNotificationLog,
				IMAPInspectionProfile:          dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].IMAPInspectionProfile,
				IMAPNotificationLog:            dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].IMAPNotificationLog,
				SMTPInspectionProfile:          dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].SMTPInspectionProfile,
				SMTPNotificationLog:            dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].SMTPNotificationLog,
				VerdictThreshold:               dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].VerdictThreshold,
				WhitelistNotificationLog:       dataV0.AdvancedAntiMalware[0].DefaultPolicy[0].WhitelistNotificationLog,
			}
		}
	}
	if len(dataV0.ApplicationIdentification) > 0 {
		dataV1.ApplicationIdentification = &servicesBlockApplicationIdentification{
			ApplicationSystemCacheTimeout: dataV0.ApplicationIdentification[0].ApplicationSystemCacheTimeout,
			GlobalOffloadByteLimit:        dataV0.ApplicationIdentification[0].GlobalOffloadByteLimit,
			IMAPCacheSize:                 dataV0.ApplicationIdentification[0].IMAPCacheSize,
			IMAPCacheTimeout:              dataV0.ApplicationIdentification[0].IMAPCacheTimeout,
			MaxMemory:                     dataV0.ApplicationIdentification[0].MaxMemory,
			MaxTransactions:               dataV0.ApplicationIdentification[0].MaxTransactions,
			MicroApps:                     dataV0.ApplicationIdentification[0].MicroApps,
			NoApplicationSystemCache:      dataV0.ApplicationIdentification[0].NoApplicationSystemCache,
			StatisticsInterval:            dataV0.ApplicationIdentification[0].StatisticsInterval,
		}
		if len(dataV0.ApplicationIdentification[0].ApplicationSystemCache) > 0 {
			dataV1.ApplicationIdentification.ApplicationSystemCache = &servicesBlockApplicationIdentificationBlockApplicationSystemCache{
				NoMiscellaneousServices: dataV0.ApplicationIdentification[0].ApplicationSystemCache[0].NoMiscellaneousServices,
				SecurityServices:        dataV0.ApplicationIdentification[0].ApplicationSystemCache[0].SecurityServices,
			}
		}
		if len(dataV0.ApplicationIdentification[0].Download) > 0 {
			dataV1.ApplicationIdentification.Download = &servicesBlockApplicationIdentificationBlockDownload{
				AutomaticInterval:      dataV0.ApplicationIdentification[0].Download[0].AutomaticInterval,
				AutomaticStartTime:     dataV0.ApplicationIdentification[0].Download[0].AutomaticStartTime,
				IgnoreServerValidation: dataV0.ApplicationIdentification[0].Download[0].IgnoreServerValidation,
				ProxyProfile:           dataV0.ApplicationIdentification[0].Download[0].ProxyProfile,
				URL:                    dataV0.ApplicationIdentification[0].Download[0].URL,
			}
		}
		if len(dataV0.ApplicationIdentification[0].EnablePerformanceMode) > 0 {
			dataV1.ApplicationIdentification.EnablePerformanceMode = &servicesBlockApplicationIdentificationBlockEnablePerformanceMode{
				MaxPacketThreshold: dataV0.ApplicationIdentification[0].EnablePerformanceMode[0].MaxPacketThreshold,
			}
		}
		if len(dataV0.ApplicationIdentification[0].InspectionLimitTCP) > 0 {
			dataV1.ApplicationIdentification.InspectionLimitTCP = &servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP{
				ByteLimit:   dataV0.ApplicationIdentification[0].InspectionLimitTCP[0].ByteLimit,
				PacketLimit: dataV0.ApplicationIdentification[0].InspectionLimitTCP[0].PacketLimit,
			}
		}
		if len(dataV0.ApplicationIdentification[0].InspectionLimitUDP) > 0 {
			dataV1.ApplicationIdentification.InspectionLimitUDP = &servicesBlockApplicationIdentificationBlockInspectionLimitTCPUDP{
				ByteLimit:   dataV0.ApplicationIdentification[0].InspectionLimitUDP[0].ByteLimit,
				PacketLimit: dataV0.ApplicationIdentification[0].InspectionLimitUDP[0].PacketLimit,
			}
		}
	}
	if len(dataV0.SecurityIntelligence) > 0 {
		dataV1.SecurityIntelligence = &servicesBlockSecurityIntelligence{
			AuthenticationTLSProfile: dataV0.SecurityIntelligence[0].AuthenticationTLSProfile,
			AuthenticationToken:      dataV0.SecurityIntelligence[0].AuthenticationToken,
			CategoryDisable:          dataV0.SecurityIntelligence[0].CategoryDisable,
			ProxyProfile:             dataV0.SecurityIntelligence[0].ProxyProfile,
			URL:                      dataV0.SecurityIntelligence[0].URL,
			URLParameter:             dataV0.SecurityIntelligence[0].URLParameter,
		}
		for _, block := range dataV0.SecurityIntelligence[0].DefaultPolicy {
			dataV1.SecurityIntelligence.DefaultPolicy = append(dataV1.SecurityIntelligence.DefaultPolicy,
				servicesBlockSecurityIntelligenceBlockDefaultPolicy{
					CategoryName: block.CategoryName,
					ProfileName:  block.ProfileName,
				},
			)
		}
	}
	if len(dataV0.UserIdentification) > 0 {
		dataV1.UserIdentification = &servicesBlockUserIdentification{
			DeviceInfoAuthSource: dataV0.UserIdentification[0].DeviceInfoAuthSource,
		}
		if len(dataV0.UserIdentification[0].ADAccess) > 0 {
			dataV1.UserIdentification.ADAccess = &servicesBlockUserIdentificationBlockADAccess{
				AuthEntryTimeout:          dataV0.UserIdentification[0].ADAccess[0].AuthEntryTimeout,
				FilterExclude:             dataV0.UserIdentification[0].ADAccess[0].FilterExclude,
				FilterInclude:             dataV0.UserIdentification[0].ADAccess[0].FilterInclude,
				FirewallAuthForcedTimeout: dataV0.UserIdentification[0].ADAccess[0].FirewallAuthForcedTimeout,
				InvalidAuthEntryTimeout:   dataV0.UserIdentification[0].ADAccess[0].InvalidAuthEntryTimeout,
				NoOnDemandProbe:           dataV0.UserIdentification[0].ADAccess[0].NoOnDemandProbe,
				WmiTimeout:                dataV0.UserIdentification[0].ADAccess[0].WmiTimeout,
			}
		}
		if len(dataV0.UserIdentification[0].IdentityManagement) > 0 {
			dataV1.UserIdentification.IdentityManagement = &servicesBlockUserIdentificationBlockIdentityManagement{
				AuthenticationEntryTimeout:        dataV0.UserIdentification[0].IdentityManagement[0].AuthenticationEntryTimeout,
				BatchQueryInterval:                dataV0.UserIdentification[0].IdentityManagement[0].BatchQueryInterval,
				BatchQueryItemsPerBatch:           dataV0.UserIdentification[0].IdentityManagement[0].BatchQueryItemsPerBatch,
				FilterDomain:                      dataV0.UserIdentification[0].IdentityManagement[0].FilterDomain,
				FilterExcludeIPAddressBook:        dataV0.UserIdentification[0].IdentityManagement[0].FilterExcludeIPAddressBook,
				FilterExcludeIPAddressSet:         dataV0.UserIdentification[0].IdentityManagement[0].FilterExcludeIPAddressSet,
				FilterIncludeIPAddressBook:        dataV0.UserIdentification[0].IdentityManagement[0].FilterIncludeIPAddressBook,
				FilterIncludeIPAddressSet:         dataV0.UserIdentification[0].IdentityManagement[0].FilterIncludeIPAddressSet,
				InvalidAuthenticationEntryTimeout: dataV0.UserIdentification[0].IdentityManagement[0].InvalidAuthenticationEntryTimeout,
				IPQueryDisable:                    dataV0.UserIdentification[0].IdentityManagement[0].IPQueryDisable,
				IPQueryDelayTime:                  dataV0.UserIdentification[0].IdentityManagement[0].IPQueryDelayTime,
			}
			if len(dataV0.UserIdentification[0].IdentityManagement[0].Connection) > 0 {
				dataV1.UserIdentification.IdentityManagement.Connection = &servicesBlockUserIdentificationBlockIdentityManagementBlockConnection{
					PrimaryAddress:         dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].PrimaryAddress,
					PrimaryClientID:        dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].PrimaryClientID,
					PrimaryClientSecret:    dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].PrimaryClientSecret,
					ConnectMethod:          dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].ConnectMethod,
					Port:                   dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].Port,
					PrimaryCACertificate:   dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].PrimaryCACertificate,
					QueryAPI:               dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].QueryAPI,
					SecondaryAddress:       dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].SecondaryAddress,
					SecondaryCACertificate: dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].SecondaryCACertificate,
					SecondaryClientID:      dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].SecondaryClientID,
					SecondaryClientSecret:  dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].SecondaryClientSecret,
					TokenAPI:               dataV0.UserIdentification[0].IdentityManagement[0].Connection[0].TokenAPI,
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
