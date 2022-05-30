package junos

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	jdecode "github.com/jeremmfr/junosdecode"
)

type servicesOptions struct {
	advAntiMalware       []map[string]interface{}
	appIdent             []map[string]interface{}
	securityIntelligence []map[string]interface{}
	userIdentification   []map[string]interface{}
}

func resourceServices() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceServicesCreate,
		ReadWithoutTimeout:   resourceServicesRead,
		UpdateWithoutTimeout: resourceServicesUpdate,
		DeleteWithoutTimeout: resourceServicesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServicesImport,
		},
		Schema: map[string]*schema.Schema{
			"clean_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"advanced_anti_malware": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth_tls_profile": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"proxy_profile": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_address": {
										Type:          schema.TypeString,
										Optional:      true,
										ValidateFunc:  validation.IsIPAddress,
										ConflictsWith: []string{"advanced_anti_malware.0.connection.0.source_interface"},
									},
									"source_interface": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"advanced_anti_malware.0.connection.0.source_address"},
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"default_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blacklist_notification_log": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"default_notification_log": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"fallback_options_action": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"block", "permit"}, false),
									},
									"fallback_options_notification_log": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"http_action": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"block", "permit"}, false),
										RequiredWith: []string{"advanced_anti_malware.0.default_policy.0.http_inspection_profile"},
									},
									"http_client_notify_file": {
										Type:     schema.TypeString,
										Optional: true,
										ConflictsWith: []string{
											"advanced_anti_malware.0.default_policy.0.http_client_notify_message",
											"advanced_anti_malware.0.default_policy.0.http_client_notify_redirect_url",
										},
										RequiredWith: []string{
											"advanced_anti_malware.0.default_policy.0.http_action",
											"advanced_anti_malware.0.default_policy.0.http_inspection_profile",
										},
									},
									"http_client_notify_message": {
										Type:     schema.TypeString,
										Optional: true,
										ConflictsWith: []string{
											"advanced_anti_malware.0.default_policy.0.http_client_notify_file",
											"advanced_anti_malware.0.default_policy.0.http_client_notify_redirect_url",
										},
										RequiredWith: []string{
											"advanced_anti_malware.0.default_policy.0.http_action",
											"advanced_anti_malware.0.default_policy.0.http_inspection_profile",
										},
									},
									"http_client_notify_redirect_url": {
										Type:     schema.TypeString,
										Optional: true,
										ConflictsWith: []string{
											"advanced_anti_malware.0.default_policy.0.http_client_notify_file",
											"advanced_anti_malware.0.default_policy.0.http_client_notify_message",
										},
										RequiredWith: []string{
											"advanced_anti_malware.0.default_policy.0.http_action",
											"advanced_anti_malware.0.default_policy.0.http_inspection_profile",
										},
									},
									"http_file_verdict_unknown": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"block", "permit"}, false),
										RequiredWith: []string{
											"advanced_anti_malware.0.default_policy.0.http_action",
											"advanced_anti_malware.0.default_policy.0.http_inspection_profile",
										},
									},
									"http_inspection_profile": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{"advanced_anti_malware.0.default_policy.0.http_action"},
									},
									"http_notification_log": {
										Type:     schema.TypeBool,
										Optional: true,
										RequiredWith: []string{
											"advanced_anti_malware.0.default_policy.0.http_action",
											"advanced_anti_malware.0.default_policy.0.http_inspection_profile",
										},
									},
									"imap_inspection_profile": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"imap_notification_log": {
										Type:         schema.TypeBool,
										Optional:     true,
										RequiredWith: []string{"advanced_anti_malware.0.default_policy.0.imap_inspection_profile"},
									},
									"smtp_inspection_profile": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"smtp_notification_log": {
										Type:         schema.TypeBool,
										Optional:     true,
										RequiredWith: []string{"advanced_anti_malware.0.default_policy.0.smtp_inspection_profile"},
									},
									"verdict_threshold": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"recommended", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
										}, false),
									},
									"whitelist_notification_log": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"application_identification": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_system_cache": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: []string{"application_identification.0.no_application_system_cache"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"no_miscellaneous_services": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"security_services": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"no_application_system_cache": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"application_identification.0.application_system_cache"},
						},

						"application_system_cache_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 1000000),
						},
						"download": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"automatic_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(6, 720),
									},
									"automatic_start_time": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(
											`^([0-9]{4}-)?(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]).(2[0-3]|[01][0-9]):[0-5][0-9](:[0-5][0-9])?$`),
											"Invalid date; format is MM-DD.hh:mm / YYYY-MM-DD.hh:mm:ss"),
									},
									"ignore_server_validation": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"proxy_profile": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"enable_performance_mode": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_packet_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
								},
							},
						},
						"global_offload_byte_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
						"imap_cache_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 512000),
						},
						"imap_cache_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 86400),
						},
						"inspection_limit_tcp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"byte_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"packet_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
								},
							},
						},
						"inspection_limit_udp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"byte_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"packet_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
								},
							},
						},
						"max_memory": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 200000),
						},
						"max_transactions": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 25),
						},
						"micro_apps": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"statistics_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 1440),
						},
					},
				},
			},
			"security_intelligence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_token": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^[a-zA-Z0-9]{32}$`),
								"Auth token must be consisted of 32 alphanumeric characters"),
							ConflictsWith: []string{"security_intelligence.0.authentication_tls_profile"},
						},
						"authentication_tls_profile": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"security_intelligence.0.authentication_token"},
						},
						"category_disable": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"default_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"profile_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"proxy_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"url_parameter": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"user_identification": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_access": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"user_identification.0.identity_management"},
							MaxItems:      1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth_entry_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 1440),
									},
									"filter_exclude": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"filter_include": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"firewall_auth_forced_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(10, 1440),
									},
									"invalid_auth_entry_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 1440),
									},
									"no_on_demand_probe": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"wmi_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(3, 120),
									},
								},
							},
						},
						"device_info_auth_source": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"active-directory", "network-access-controller"}, false),
						},
						"identity_management": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"user_identification.0.ad_access"},
							MaxItems:      1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connection": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"primary_address": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.IsIPAddress,
												},
												"primary_client_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"primary_client_secret": {
													Type:      schema.TypeString,
													Required:  true,
													Sensitive: true,
												},
												"connect_method": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
												},
												"port": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(1, 65535),
												},
												"primary_ca_certificate": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"query_api": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"secondary_address": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.IsIPAddress,
												},
												"secondary_ca_certificate": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"secondary_client_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"secondary_client_secret": {
													Type:      schema.TypeString,
													Optional:  true,
													Sensitive: true,
												},
												"token_api": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"authentication_entry_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 1440),
									},
									"batch_query_items_per_batch": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(100, 1000),
									},
									"batch_query_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 60),
									},
									"filter_domain": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"filter_exclude_ip_address_book": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{"user_identification.0.identity_management.0.filter_exclude_ip_address_set"},
									},
									"filter_exclude_ip_address_set": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{"user_identification.0.identity_management.0.filter_exclude_ip_address_book"},
									},
									"filter_include_ip_address_book": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{"user_identification.0.identity_management.0.filter_include_ip_address_set"},
									},
									"filter_include_ip_address_set": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{"user_identification.0.identity_management.0.filter_include_ip_address_book"},
									},
									"invalid_authentication_entry_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 1440),
									},
									"ip_query_disable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ip_query_delay_time": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceServicesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setServices(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("services")

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := setServices(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_services", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("services")

	return append(diagWarns, resourceServicesReadWJunSess(d, clt, junSess)...)
}

func resourceServicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceServicesReadWJunSess(d, clt, junSess)
}

func resourceServicesReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession) diag.Diagnostics {
	mutex.Lock()
	servicesOptions, err := readServices(clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillServices(d, servicesOptions)

	return nil
}

func resourceServicesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delServices(d, false, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setServices(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delServices(d, false, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServices(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_services", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesReadWJunSess(d, clt, junSess)...)
}

func resourceServicesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("clean_on_destroy").(bool) {
		clt := m.(*Client)
		if clt.fakeDeleteAlso {
			if err := delServices(d, true, clt, nil); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
		junSess, err := clt.startNewSession(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		defer clt.closeSession(junSess)
		if err := clt.configLock(ctx, junSess); err != nil {
			return diag.FromErr(err)
		}
		var diagWarns diag.Diagnostics
		if err := delServices(d, true, clt, junSess); err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		warns, err := clt.commitConf("delete resource junos_services", junSess)
		appendDiagWarns(&diagWarns, warns)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
	}

	return nil
}

func resourceServicesImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	servicesOptions, err := readServices(clt, junSess)
	if err != nil {
		return nil, err
	}
	fillServices(d, servicesOptions)
	d.SetId("services")
	result[0] = d

	return result, nil
}

func setServices(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)
	for _, v := range d.Get("advanced_anti_malware").([]interface{}) {
		configSetAdvAntiMalware, err := setServicesAdvancedAntiMalware(d, v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetAdvAntiMalware...)
	}
	for _, v := range d.Get("application_identification").([]interface{}) {
		configSetApplicationIdentification, err := setServicesApplicationIdentification(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetApplicationIdentification...)
	}
	for _, v := range d.Get("security_intelligence").([]interface{}) {
		configSetSecurityIntelligence, err := setServicesSecurityIntell(d, v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetSecurityIntelligence...)
	}
	for _, v := range d.Get("user_identification").([]interface{}) {
		configSetUserIdent, err := setServicesUserIdentification(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetUserIdent...)
	}

	return clt.configSet(configSet, junSess)
}

func setServicesAdvancedAntiMalware(d *schema.ResourceData, advAntiMalware interface{}) ([]string, error) {
	setPrefix := "set services advanced-anti-malware "
	configSet := make([]string, 0)
	if advAntiMalware != nil {
		advAntiMalwareM := advAntiMalware.(map[string]interface{})
		for _, v := range advAntiMalwareM["connection"].([]interface{}) {
			setPrefixConn := setPrefix + "connection "
			configSet = append(configSet, setPrefixConn)
			if v != nil {
				connection := v.(map[string]interface{})
				if d.HasChange("advanced_anti_malware.0.connection.0.auth_tls_profile") &&
					connection["auth_tls_profile"].(string) != "" {
					configSet = append(configSet, setPrefixConn+"authentication tls-profile \""+
						connection["auth_tls_profile"].(string)+"\"")
				}
				if v2 := connection["proxy_profile"].(string); v2 != "" {
					configSet = append(configSet, setPrefixConn+"proxy-profile \""+v2+"\"")
				}
				if v2 := connection["source_address"].(string); v2 != "" {
					configSet = append(configSet, setPrefixConn+"source-address "+v2)
				}
				if v2 := connection["source_interface"].(string); v2 != "" {
					configSet = append(configSet, setPrefixConn+"source-interface "+v2)
				}
				if d.HasChange("advanced_anti_malware.0.connection.0.url") &&
					connection["url"].(string) != "" {
					configSet = append(configSet, setPrefixConn+"url \""+connection["url"].(string)+"\"")
				}
			}
		}
		for _, v := range advAntiMalwareM["default_policy"].([]interface{}) {
			setPrefixDefPolicy := setPrefix + "default-policy "
			configSet = append(configSet, setPrefixDefPolicy)
			if v != nil {
				defPolicy := v.(map[string]interface{})
				if defPolicy["blacklist_notification_log"].(bool) {
					configSet = append(configSet, setPrefixDefPolicy+"blacklist-notification log")
				}
				if defPolicy["default_notification_log"].(bool) {
					configSet = append(configSet, setPrefixDefPolicy+"default-notification log")
				}
				if v2 := defPolicy["fallback_options_action"].(string); v2 != "" {
					configSet = append(configSet, setPrefixDefPolicy+"fallback-options action "+v2)
				}
				if defPolicy["fallback_options_notification_log"].(bool) {
					configSet = append(configSet, setPrefixDefPolicy+"fallback-options notification log")
				}
				if v2 := defPolicy["http_action"].(string); v2 != "" {
					configSet = append(configSet, setPrefixDefPolicy+"http action "+v2)
				}
				if v := defPolicy["http_client_notify_file"].(string); v != "" {
					configSet = append(configSet, setPrefixDefPolicy+"http client-notify file \""+v+"\"")
				}
				if v := defPolicy["http_client_notify_message"].(string); v != "" {
					configSet = append(configSet, setPrefixDefPolicy+"http client-notify message \""+v+"\"")
				}
				if v := defPolicy["http_client_notify_redirect_url"].(string); v != "" {
					configSet = append(configSet, setPrefixDefPolicy+"http client-notify redirect-url \""+v+"\"")
				}
				if v := defPolicy["http_file_verdict_unknown"].(string); v != "" {
					configSet = append(configSet, setPrefixDefPolicy+"http file-verdict-unknown "+v)
				}
				if v2 := defPolicy["http_inspection_profile"].(string); v2 != "" {
					configSet = append(configSet, setPrefixDefPolicy+"http inspection-profile \""+v2+"\"")
				}
				if defPolicy["http_notification_log"].(bool) {
					configSet = append(configSet, setPrefixDefPolicy+"http notification log")
				}
				if v2 := defPolicy["imap_inspection_profile"].(string); v2 != "" {
					configSet = append(configSet, setPrefixDefPolicy+"imap inspection-profile \""+v2+"\"")
				}
				if defPolicy["imap_notification_log"].(bool) {
					configSet = append(configSet, setPrefixDefPolicy+"imap notification log")
				}
				if v2 := defPolicy["smtp_inspection_profile"].(string); v2 != "" {
					configSet = append(configSet, setPrefixDefPolicy+"smtp inspection-profile \""+v2+"\"")
				}
				if defPolicy["smtp_notification_log"].(bool) {
					configSet = append(configSet, setPrefixDefPolicy+"smtp notification log")
				}
				if v2 := defPolicy["verdict_threshold"].(string); v2 != "" {
					configSet = append(configSet, setPrefixDefPolicy+"verdict-threshold "+v2)
				}
				if defPolicy["whitelist_notification_log"].(bool) {
					configSet = append(configSet, setPrefixDefPolicy+"whitelist-notification log")
				}
			}
		}
	} else {
		return configSet, fmt.Errorf("advanced_anti_malware block is empty")
	}

	return configSet, nil
}

func setServicesApplicationIdentification(appID interface{}) ([]string, error) {
	setPrefix := "set services application-identification "
	configSet := make([]string, 0)
	appIDM := appID.(map[string]interface{})
	configSet = append(configSet, setPrefix)
	for _, v := range appIDM["application_system_cache"].([]interface{}) {
		configSet = append(configSet, setPrefix+"application-system-cache")
		if v != nil {
			appSysCache := v.(map[string]interface{})
			if appSysCache["no_miscellaneous_services"].(bool) {
				configSet = append(configSet, setPrefix+"application-system-cache no-miscellaneous-services")
			}
			if appSysCache["security_services"].(bool) {
				configSet = append(configSet, setPrefix+"application-system-cache security-services")
			}
		}
	}
	if appIDM["no_application_system_cache"].(bool) {
		configSet = append(configSet, setPrefix+"no-application-system-cache")
	}
	if v := appIDM["application_system_cache_timeout"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"application-system-cache-timeout "+strconv.Itoa(v))
	}
	for _, v := range appIDM["download"].([]interface{}) {
		if v != nil {
			download := v.(map[string]interface{})
			if v2 := download["automatic_interval"].(int); v2 != 0 {
				configSet = append(configSet, setPrefix+"download automatic interval "+strconv.Itoa(v2))
			}
			if v2 := download["automatic_start_time"].(string); v2 != "" {
				configSet = append(configSet, setPrefix+"download automatic start-time "+v2)
			}
			if download["ignore_server_validation"].(bool) {
				configSet = append(configSet, setPrefix+"download ignore-server-validation")
			}
			if v2 := download["proxy_profile"].(string); v2 != "" {
				configSet = append(configSet, setPrefix+"download proxy-profile \""+v2+"\"")
			}
			if v2 := download["url"].(string); v2 != "" {
				configSet = append(configSet, setPrefix+"download url \""+v2+"\"")
			}
		} else {
			return configSet, fmt.Errorf("application_identification.0.download block is empty")
		}
	}
	for _, v := range appIDM["enable_performance_mode"].([]interface{}) {
		configSet = append(configSet, setPrefix+"enable-performance-mode")
		if v != nil {
			enPerfMode := v.(map[string]interface{})
			if v := enPerfMode["max_packet_threshold"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"enable-performance-mode max-packet-threshold "+strconv.Itoa(v))
			}
		}
	}
	if v := appIDM["global_offload_byte_limit"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"global-offload-byte-limit "+strconv.Itoa(v))
	}
	if v := appIDM["imap_cache_size"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"imap-cache-size "+strconv.Itoa(v))
	}
	if v := appIDM["imap_cache_timeout"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"imap-cache-timeout "+strconv.Itoa(v))
	}
	for _, v := range appIDM["inspection_limit_tcp"].([]interface{}) {
		configSet = append(configSet, setPrefix+"inspection-limit tcp")
		if v != nil {
			inspLimitTCP := v.(map[string]interface{})
			if v := inspLimitTCP["byte_limit"].(int); v != -1 {
				configSet = append(configSet, setPrefix+"inspection-limit tcp byte-limit "+strconv.Itoa(v))
			}
			if v := inspLimitTCP["packet_limit"].(int); v != -1 {
				configSet = append(configSet, setPrefix+"inspection-limit tcp packet-limit "+strconv.Itoa(v))
			}
		}
	}
	for _, v := range appIDM["inspection_limit_udp"].([]interface{}) {
		configSet = append(configSet, setPrefix+"inspection-limit udp")
		if v != nil {
			inspLimitUDP := v.(map[string]interface{})
			if v := inspLimitUDP["byte_limit"].(int); v != -1 {
				configSet = append(configSet, setPrefix+"inspection-limit udp byte-limit "+strconv.Itoa(v))
			}
			if v := inspLimitUDP["packet_limit"].(int); v != -1 {
				configSet = append(configSet, setPrefix+"inspection-limit udp packet-limit "+strconv.Itoa(v))
			}
		}
	}
	if v := appIDM["max_memory"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"max-memory "+strconv.Itoa(v))
	}
	if v := appIDM["max_transactions"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"max-transactions "+strconv.Itoa(v))
	}
	if appIDM["micro_apps"].(bool) {
		configSet = append(configSet, setPrefix+"micro-apps")
	}
	if v := appIDM["statistics_interval"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"statistics interval "+strconv.Itoa(v))
	}

	return configSet, nil
}

func setServicesSecurityIntell(d *schema.ResourceData, secuIntel interface{}) ([]string, error) {
	setPrefix := "set services security-intelligence "
	configSet := make([]string, 0)
	if secuIntel != nil {
		secuIntelM := secuIntel.(map[string]interface{})
		if d.HasChange("security_intelligence.0.authentication_token") &&
			secuIntelM["authentication_token"].(string) != "" {
			configSet = append(configSet, "delete services security-intelligence authentication")
			configSet = append(configSet,
				setPrefix+"authentication auth-token "+secuIntelM["authentication_token"].(string))
		}
		if d.HasChange("security_intelligence.0.authentication_tls_profile") &&
			secuIntelM["authentication_tls_profile"].(string) != "" {
			configSet = append(configSet, "delete services security-intelligence authentication")
			configSet = append(configSet,
				setPrefix+"authentication tls-profile \""+secuIntelM["authentication_tls_profile"].(string)+"\"")
		}
		for _, v := range sortSetOfString(secuIntelM["category_disable"].(*schema.Set).List()) {
			if v == "all" {
				configSet = append(configSet, setPrefix+"category all disable")
			} else {
				configSet = append(configSet, setPrefix+"category category-name "+v+" disable")
			}
		}
		defaultPolicyCatNameList := make([]string, 0)
		for _, v := range secuIntelM["default_policy"].([]interface{}) {
			defPolicy := v.(map[string]interface{})
			if bchk.StringInSlice(defPolicy["category_name"].(string), defaultPolicyCatNameList) {
				return configSet, fmt.Errorf("multiple blocks default_policy with the same category_name %s",
					defPolicy["category_name"].(string))
			}
			defaultPolicyCatNameList = append(defaultPolicyCatNameList, defPolicy["category_name"].(string))
			configSet = append(configSet, setPrefix+"default-policy "+
				defPolicy["category_name"].(string)+" "+defPolicy["profile_name"].(string))
		}
		if v := secuIntelM["proxy_profile"].(string); v != "" {
			configSet = append(configSet, setPrefix+"proxy-profile \""+v+"\"")
		}
		if d.HasChange("security_intelligence.0.url") &&
			secuIntelM["url"].(string) != "" {
			configSet = append(configSet, setPrefix+"url \""+secuIntelM["url"].(string)+"\"")
		}
		if v := secuIntelM["url_parameter"].(string); v != "" {
			configSet = append(configSet, setPrefix+"url-parameter \""+v+"\"")
		}
	}

	return configSet, nil
}

func setServicesUserIdentification(userIdentification interface{}) ([]string, error) {
	setPrefix := "set services user-identification "
	configSet := make([]string, 0)
	if userIdentification != nil {
		userIdent := userIdentification.(map[string]interface{})
		if len(userIdent["ad_access"].([]interface{})) == 0 {
			configSet = append(configSet, "delete services user-identification active-directory-access")
		}
		for _, v := range userIdent["ad_access"].([]interface{}) {
			adAccess := v.(map[string]interface{})
			configSet = append(configSet, setPrefix+"active-directory-access")
			if v2 := adAccess["auth_entry_timeout"].(int); v2 != -1 {
				configSet = append(configSet, setPrefix+"active-directory-access authentication-entry-timeout "+
					strconv.Itoa(v2))
			}
			for _, v2 := range sortSetOfString(adAccess["filter_exclude"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"active-directory-access filter exclude "+v2)
			}
			for _, v2 := range sortSetOfString(adAccess["filter_include"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"active-directory-access filter include "+v2)
			}
			if v2 := adAccess["firewall_auth_forced_timeout"].(int); v2 != 0 {
				configSet = append(configSet, setPrefix+"active-directory-access firewall-authentication-forced-timeout "+
					strconv.Itoa(v2))
			}
			if v2 := adAccess["invalid_auth_entry_timeout"].(int); v2 != -1 {
				configSet = append(configSet, setPrefix+"active-directory-access invalid-authentication-entry-timeout "+
					strconv.Itoa(v2))
			}
			if adAccess["no_on_demand_probe"].(bool) {
				configSet = append(configSet, setPrefix+"active-directory-access no-on-demand-probe")
			}
			if v2 := adAccess["wmi_timeout"].(int); v2 != 0 {
				configSet = append(configSet, setPrefix+"active-directory-access wmi-timeout "+strconv.Itoa(v2))
			}
		}
		if v := userIdent["device_info_auth_source"].(string); v != "" {
			configSet = append(configSet, setPrefix+"device-information authentication-source "+v)
		}
		for _, v := range userIdent["identity_management"].([]interface{}) {
			setPrefixIdentMgmt := setPrefix + "identity-management "
			identMgmt := v.(map[string]interface{})
			for _, v2 := range identMgmt["connection"].([]interface{}) {
				connection := v2.(map[string]interface{})
				setPrefixIMConn := setPrefixIdentMgmt + "connection "
				configSet = append(configSet, setPrefixIMConn+"primary address "+
					connection["primary_address"].(string))
				configSet = append(configSet, setPrefixIMConn+"primary client-id \""+
					connection["primary_client_id"].(string)+"\"")
				configSet = append(configSet, setPrefixIMConn+"primary client-secret \""+
					connection["primary_client_secret"].(string)+"\"")
				if v3 := connection["connect_method"].(string); v3 != "" {
					configSet = append(configSet, setPrefixIMConn+"connect-method "+v3)
				}
				if v3 := connection["port"].(int); v3 != 0 {
					configSet = append(configSet, setPrefixIMConn+"port "+strconv.Itoa(v3))
				}
				if v3 := connection["primary_ca_certificate"].(string); v3 != "" {
					configSet = append(configSet, setPrefixIMConn+"primary ca-certificate \""+v3+"\"")
				}
				if v3 := connection["query_api"].(string); v3 != "" {
					configSet = append(configSet, setPrefixIMConn+"query-api \""+v3+"\"")
				}
				if v3 := connection["secondary_address"].(string); v3 != "" {
					configSet = append(configSet, setPrefixIMConn+"secondary address "+v3)
				}
				if v3 := connection["secondary_ca_certificate"].(string); v3 != "" {
					configSet = append(configSet, setPrefixIMConn+"secondary ca-certificate \""+v3+"\"")
				}
				if v3 := connection["secondary_client_id"].(string); v3 != "" {
					configSet = append(configSet, setPrefixIMConn+"secondary client-id \""+v3+"\"")
				}
				if v3 := connection["secondary_client_secret"].(string); v3 != "" {
					configSet = append(configSet, setPrefixIMConn+"secondary client-secret \""+v3+"\"")
				}
				if v3 := connection["token_api"].(string); v3 != "" {
					configSet = append(configSet, setPrefixIMConn+"token-api \""+v3+"\"")
				}
			}
			if v2 := identMgmt["authentication_entry_timeout"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixIdentMgmt+"authentication-entry-timeout "+strconv.Itoa(v2))
			}
			if v2 := identMgmt["batch_query_items_per_batch"].(int); v2 != 0 {
				configSet = append(configSet, setPrefixIdentMgmt+"batch-query items-per-batch "+strconv.Itoa(v2))
			}
			if v2 := identMgmt["batch_query_interval"].(int); v2 != 0 {
				configSet = append(configSet, setPrefixIdentMgmt+"batch-query query-interval "+strconv.Itoa(v2))
			}
			for _, v2 := range sortSetOfString(identMgmt["filter_domain"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefixIdentMgmt+"filter domain "+v2)
			}
			if v2 := identMgmt["filter_exclude_ip_address_book"].(string); v2 != "" {
				configSet = append(configSet, setPrefixIdentMgmt+"filter exclude-ip address-book \""+v2+"\"")
			}
			if v2 := identMgmt["filter_exclude_ip_address_set"].(string); v2 != "" {
				configSet = append(configSet, setPrefixIdentMgmt+"filter exclude-ip address-set \""+v2+"\"")
			}
			if v2 := identMgmt["filter_include_ip_address_book"].(string); v2 != "" {
				configSet = append(configSet, setPrefixIdentMgmt+"filter include-ip address-book \""+v2+"\"")
			}
			if v2 := identMgmt["filter_include_ip_address_set"].(string); v2 != "" {
				configSet = append(configSet, setPrefixIdentMgmt+"filter include-ip address-set \""+v2+"\"")
			}
			if v2 := identMgmt["invalid_authentication_entry_timeout"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixIdentMgmt+"invalid-authentication-entry-timeout "+strconv.Itoa(v2))
			}
			if identMgmt["ip_query_disable"].(bool) {
				configSet = append(configSet, setPrefixIdentMgmt+"ip-query no-ip-query")
			}
			if v2 := identMgmt["ip_query_delay_time"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixIdentMgmt+"ip-query query-delay-time "+strconv.Itoa(v2))
			}
		}
	} else {
		return configSet, fmt.Errorf("user_identification block is empty")
	}

	return configSet, nil
}

func listLinesServicesAdvancedAntiMalware() []string {
	return []string{
		"advanced-anti-malware connection proxy-profile",
		"advanced-anti-malware connection source-address",
		"advanced-anti-malware connection source-interface",
		"advanced-anti-malware default-policy",
	}
}

func listLinesServicesApplicationIdentification() []string {
	return []string{
		"application-identification application-system-cache",
		"application-identification no-application-system-cache",
		"application-identification application-system-cache-timeout",
		"application-identification download",
		"application-identification global-offload-byte-limit",
		"application-identification enable-performance-mode",
		"application-identification imap-cache-size",
		"application-identification imap-cache-timeout",
		"application-identification inspection-limit tcp",
		"application-identification inspection-limit udp",
		"application-identification max-memory",
		"application-identification max-transactions",
		"application-identification micro-apps",
		"application-identification statistics interval",
	}
}

func listLinesServicesSecurityIntel() []string {
	return []string{
		"security-intelligence category",
		"security-intelligence default-policy",
		"security-intelligence proxy-profile",
		"security-intelligence url-parameter",
	}
}

func listLinesServicesUserIdentification() []string {
	r := []string{
		"user-identification device-information authentication-source",
		"user-identification identity-management",
	}
	r = append(r, listLinesServicesUserIdentificationAdAccess()...)

	return r
}

func listLinesServicesUserIdentificationAdAccess() []string {
	return []string{
		"user-identification active-directory-access authentication-entry-timeout",
		"user-identification active-directory-access filter",
		"user-identification active-directory-access firewall-authentication-forced-timeout",
		"user-identification active-directory-access invalid-authentication-entry-timeout",
		"user-identification active-directory-access no-on-demand-probe",
		"user-identification active-directory-access wmi-timeout",
	}
}

func delServices(d *schema.ResourceData, cleanAll bool, clt *Client, junSess *junosSession) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, listLinesServicesAdvancedAntiMalware()...)
	listLinesToDelete = append(listLinesToDelete, listLinesServicesApplicationIdentification()...)
	listLinesToDelete = append(listLinesToDelete, listLinesServicesSecurityIntel()...)
	listLinesToDelete = append(listLinesToDelete, listLinesServicesUserIdentification()...)

	configSet := make([]string, 0)
	delPrefix := "delete services "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}
	if len(d.Get("advanced_anti_malware").([]interface{})) == 0 || cleanAll {
		configSet = append(configSet, delPrefix+"advanced-anti-malware connection")
	} else {
		advAntiMalware := d.Get("advanced_anti_malware").([]interface{})[0].(map[string]interface{})
		if len(advAntiMalware["connection"].([]interface{})) == 0 {
			configSet = append(configSet, delPrefix+"advanced-anti-malware connection")
		}
	}
	if len(d.Get("application_identification").([]interface{})) == 0 || cleanAll {
		configSet = append(configSet, delPrefix+"application-identification")
	}
	if len(d.Get("security_intelligence").([]interface{})) == 0 || cleanAll {
		configSet = append(configSet, delPrefix+"security-intelligence authentication")
		configSet = append(configSet, delPrefix+"security-intelligence url")
	}
	if len(d.Get("user_identification").([]interface{})) == 0 || cleanAll {
		configSet = append(configSet, delPrefix+"user-identification active-directory-access")
	}

	return clt.configSet(configSet, junSess)
}

func readServices(clt *Client, junSess *junosSession) (servicesOptions, error) {
	var confRead servicesOptions

	showConfig, err := clt.command(cmdShowConfig+"services"+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case bchk.StringHasOneOfPrefixes(itemTrim, listLinesServicesAdvancedAntiMalware()) ||
				strings.HasPrefix(itemTrim, "advanced-anti-malware connection"):
				readServicesAdvancedAntiMalware(&confRead, itemTrim)
			case strings.HasPrefix(itemTrim, "application-identification"):
				if err := readServicesApplicationIdentification(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, listLinesServicesSecurityIntel()) ||
				strings.HasPrefix(itemTrim, "security-intelligence authentication ") ||
				strings.HasPrefix(itemTrim, "security-intelligence url "):
				if err := readServicesSecurityIntel(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, listLinesServicesUserIdentification()) ||
				strings.HasPrefix(itemTrim, "user-identification active-directory-access"):
				if err := readServicesUserIdentification(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}

func readServicesAdvancedAntiMalware(confRead *servicesOptions, itemTrimAdvAntiMalware string) {
	itemTrim := strings.TrimPrefix(itemTrimAdvAntiMalware, "advanced-anti-malware ")
	if len(confRead.advAntiMalware) == 0 {
		confRead.advAntiMalware = append(confRead.advAntiMalware, map[string]interface{}{
			"connection":     make([]map[string]interface{}, 0),
			"default_policy": make([]map[string]interface{}, 0),
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "connection"):
		if len(confRead.advAntiMalware[0]["connection"].([]map[string]interface{})) == 0 {
			confRead.advAntiMalware[0]["connection"] = append(
				confRead.advAntiMalware[0]["connection"].([]map[string]interface{}),
				map[string]interface{}{
					"auth_tls_profile": "",
					"proxy_profile":    "",
					"source_address":   "",
					"source_interface": "",
					"url":              "",
				})
		}
		connection := confRead.advAntiMalware[0]["connection"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "connection authentication tls-profile "):
			connection["auth_tls_profile"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "connection authentication tls-profile "), "\"")
		case strings.HasPrefix(itemTrim, "connection proxy-profile "):
			connection["proxy_profile"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "connection proxy-profile "), "\"")
		case strings.HasPrefix(itemTrim, "connection source-address "):
			connection["source_address"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "connection source-address "), "\"")
		case strings.HasPrefix(itemTrim, "connection source-interface "):
			connection["source_interface"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "connection source-interface "), "\"")
		case strings.HasPrefix(itemTrim, "connection url "):
			connection["url"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "connection url "), "\"")
		}
	case strings.HasPrefix(itemTrim, "default-policy"):
		if len(confRead.advAntiMalware[0]["default_policy"].([]map[string]interface{})) == 0 {
			confRead.advAntiMalware[0]["default_policy"] = append(
				confRead.advAntiMalware[0]["default_policy"].([]map[string]interface{}),
				map[string]interface{}{
					"blacklist_notification_log":        false,
					"default_notification_log":          false,
					"fallback_options_action":           "",
					"fallback_options_notification_log": false,
					"http_action":                       "",
					"http_client_notify_file":           "",
					"http_client_notify_message":        "",
					"http_client_notify_redirect_url":   "",
					"http_file_verdict_unknown":         "",
					"http_inspection_profile":           "",
					"http_notification_log":             false,
					"imap_inspection_profile":           "",
					"imap_notification_log":             false,
					"smtp_inspection_profile":           "",
					"smtp_notification_log":             false,
					"verdict_threshold":                 "",
					"whitelist_notification_log":        false,
				})
		}
		defaultPolicy := confRead.advAntiMalware[0]["default_policy"].([]map[string]interface{})[0]
		switch {
		case itemTrim == "default-policy blacklist-notification log":
			defaultPolicy["blacklist_notification_log"] = true
		case itemTrim == "default-policy default-notification log":
			defaultPolicy["default_notification_log"] = true
		case strings.HasPrefix(itemTrim, "default-policy fallback-options action "):
			defaultPolicy["fallback_options_action"] = strings.TrimPrefix(itemTrim, "default-policy fallback-options action ")
		case itemTrim == "default-policy fallback-options notification log":
			defaultPolicy["fallback_options_notification_log"] = true
		case strings.HasPrefix(itemTrim, "default-policy http action "):
			defaultPolicy["http_action"] = strings.TrimPrefix(itemTrim, "default-policy http action ")
		case strings.HasPrefix(itemTrim, "default-policy http client-notify file "):
			defaultPolicy["http_client_notify_file"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "default-policy http client-notify file "), "\"")
		case strings.HasPrefix(itemTrim, "default-policy http client-notify message "):
			defaultPolicy["http_client_notify_message"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "default-policy http client-notify message "), "\"")
		case strings.HasPrefix(itemTrim, "default-policy http client-notify redirect-url "):
			defaultPolicy["http_client_notify_redirect_url"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "default-policy http client-notify redirect-url "), "\"")
		case strings.HasPrefix(itemTrim, "default-policy http file-verdict-unknown "):
			defaultPolicy["http_file_verdict_unknown"] = strings.TrimPrefix(
				itemTrim, "default-policy http file-verdict-unknown ")
		case strings.HasPrefix(itemTrim, "default-policy http inspection-profile "):
			defaultPolicy["http_inspection_profile"] = strings.TrimPrefix(strings.TrimPrefix(
				itemTrim, "default-policy http inspection-profile "), "\"")
		case itemTrim == "default-policy http notification log":
			defaultPolicy["http_notification_log"] = true
		case strings.HasPrefix(itemTrim, "default-policy imap inspection-profile "):
			defaultPolicy["imap_inspection_profile"] = strings.TrimPrefix(strings.TrimPrefix(
				itemTrim, "default-policy imap inspection-profile "), "\"")
		case itemTrim == "default-policy imap notification log":
			defaultPolicy["imap_notification_log"] = true
		case strings.HasPrefix(itemTrim, "default-policy smtp inspection-profile "):
			defaultPolicy["smtp_inspection_profile"] = strings.TrimPrefix(strings.TrimPrefix(
				itemTrim, "default-policy smtp inspection-profile "), "\"")
		case itemTrim == "default-policy smtp notification log":
			defaultPolicy["smtp_notification_log"] = true
		case strings.HasPrefix(itemTrim, "default-policy verdict-threshold "):
			defaultPolicy["verdict_threshold"] = strings.TrimPrefix(itemTrim, "default-policy verdict-threshold ")
		case itemTrim == "default-policy whitelist-notification log":
			defaultPolicy["whitelist_notification_log"] = true
		}
	}
}

func readServicesSecurityIntel(confRead *servicesOptions, itemTrimSecurityIntel string) error {
	itemTrim := strings.TrimPrefix(itemTrimSecurityIntel, "security-intelligence ")
	if len(confRead.securityIntelligence) == 0 {
		confRead.securityIntelligence = append(confRead.securityIntelligence, map[string]interface{}{
			"authentication_token":       "",
			"authentication_tls_profile": "",
			"category_disable":           make([]string, 0),
			"default_policy":             make([]map[string]interface{}, 0),
			"proxy_profile":              "",
			"url":                        "",
			"url_parameter":              "",
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "authentication auth-token "):
		confRead.securityIntelligence[0]["authentication_token"] = strings.TrimPrefix(itemTrim, "authentication auth-token ")
	case strings.HasPrefix(itemTrim, "authentication tls-profile "):
		confRead.securityIntelligence[0]["authentication_tls_profile"] = strings.Trim(strings.TrimPrefix(itemTrim,
			"authentication tls-profile "), "\"")
	case strings.HasPrefix(itemTrim, "category "):
		if itemTrim == "category all disable" {
			confRead.securityIntelligence[0]["category_disable"] = append(
				confRead.securityIntelligence[0]["category_disable"].([]string), "all")
		} else {
			confRead.securityIntelligence[0]["category_disable"] = append(
				confRead.securityIntelligence[0]["category_disable"].([]string),
				strings.TrimSuffix(strings.TrimPrefix(itemTrim, "category category-name "), " disable"))
		}
	case strings.HasPrefix(itemTrim, "default-policy "):
		if itemTrimSplit := strings.Split(itemTrim, " "); len(itemTrimSplit) == 3 {
			confRead.securityIntelligence[0]["default_policy"] = append(
				confRead.securityIntelligence[0]["default_policy"].([]map[string]interface{}), map[string]interface{}{
					"category_name": itemTrimSplit[1],
					"profile_name":  itemTrimSplit[2],
				})
		}
	case strings.HasPrefix(itemTrim, "proxy-profile "):
		confRead.securityIntelligence[0]["proxy_profile"] = strings.Trim(strings.TrimPrefix(itemTrim,
			"proxy-profile "), "\"")
	case strings.HasPrefix(itemTrim, "url "):
		confRead.securityIntelligence[0]["url"] = strings.Trim(strings.TrimPrefix(itemTrim,
			"url "), "\"")
	case strings.HasPrefix(itemTrim, "url-parameter "):
		var err error
		confRead.securityIntelligence[0]["url_parameter"], err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
			itemTrim, "url-parameter "), "\""))
		if err != nil {
			return fmt.Errorf("failed to decode url-parameter: %w", err)
		}
	}

	return nil
}

func readServicesApplicationIdentification(confRead *servicesOptions, itemTrimAppID string) error {
	itemTrim := strings.TrimPrefix(itemTrimAppID, "application-identification ")
	if len(confRead.appIdent) == 0 {
		confRead.appIdent = append(confRead.appIdent, map[string]interface{}{
			"application_system_cache":         make([]map[string]interface{}, 0),
			"no_application_system_cache":      false,
			"application_system_cache_timeout": -1,
			"download":                         make([]map[string]interface{}, 0),
			"enable_performance_mode":          make([]map[string]interface{}, 0),
			"global_offload_byte_limit":        -1,
			"imap_cache_size":                  0,
			"imap_cache_timeout":               0,
			"inspection_limit_tcp":             make([]map[string]interface{}, 0),
			"inspection_limit_udp":             make([]map[string]interface{}, 0),
			"max_memory":                       0,
			"max_transactions":                 -1,
			"micro_apps":                       false,
			"statistics_interval":              0,
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "application-system-cache-timeout "):
		var err error
		confRead.appIdent[0]["application_system_cache_timeout"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "application-system-cache-timeout "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "application-system-cache"):
		if len(confRead.appIdent[0]["application_system_cache"].([]map[string]interface{})) == 0 {
			confRead.appIdent[0]["application_system_cache"] = append(
				confRead.appIdent[0]["application_system_cache"].([]map[string]interface{}),
				map[string]interface{}{
					"no_miscellaneous_services": false,
					"security_services":         false,
				})
		}
		applicationSystemCache := confRead.appIdent[0]["application_system_cache"].([]map[string]interface{})[0]
		switch {
		case itemTrim == "application-system-cache no-miscellaneous-services":
			applicationSystemCache["no_miscellaneous_services"] = true
		case itemTrim == "application-system-cache security-services":
			applicationSystemCache["security_services"] = true
		}
	case itemTrim == "no-application-system-cache":
		confRead.appIdent[0]["no_application_system_cache"] = true
	case strings.HasPrefix(itemTrim, "download "):
		if len(confRead.appIdent[0]["download"].([]map[string]interface{})) == 0 {
			confRead.appIdent[0]["download"] = append(
				confRead.appIdent[0]["download"].([]map[string]interface{}), map[string]interface{}{
					"automatic_interval":       0,
					"automatic_start_time":     "",
					"ignore_server_validation": false,
					"proxy_profile":            "",
					"url":                      "",
				})
		}
		download := confRead.appIdent[0]["download"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "download automatic interval "):
			var err error
			download["automatic_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "download automatic interval "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "download automatic start-time "):
			download["automatic_start_time"] = strings.TrimPrefix(itemTrim, "download automatic start-time ")
		case itemTrim == "download ignore-server-validation":
			download["ignore_server_validation"] = true
		case strings.HasPrefix(itemTrim, "download proxy-profile "):
			download["proxy_profile"] = strings.Trim(strings.TrimPrefix(itemTrim, "download proxy-profile "), "\"")
		case strings.HasPrefix(itemTrim, "download url "):
			download["url"] = strings.Trim(strings.TrimPrefix(itemTrim, "download url "), "\"")
		}
	case strings.HasPrefix(itemTrim, "enable-performance-mode"):
		if len(confRead.appIdent[0]["enable_performance_mode"].([]map[string]interface{})) == 0 {
			confRead.appIdent[0]["enable_performance_mode"] = append(
				confRead.appIdent[0]["enable_performance_mode"].([]map[string]interface{}), map[string]interface{}{
					"max_packet_threshold": 0,
				})
		}
		enablePerfMode := confRead.appIdent[0]["enable_performance_mode"].([]map[string]interface{})[0]
		if strings.HasPrefix(itemTrim, "enable-performance-mode max-packet-threshold ") {
			var err error
			enablePerfMode["max_packet_threshold"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "enable-performance-mode max-packet-threshold "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "global-offload-byte-limit "):
		var err error
		confRead.appIdent[0]["global_offload_byte_limit"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "global-offload-byte-limit "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "imap-cache-size "):
		var err error
		confRead.appIdent[0]["imap_cache_size"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "imap-cache-size "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "imap-cache-timeout "):
		var err error
		confRead.appIdent[0]["imap_cache_timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "imap-cache-timeout "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "inspection-limit tcp"):
		if len(confRead.appIdent[0]["inspection_limit_tcp"].([]map[string]interface{})) == 0 {
			confRead.appIdent[0]["inspection_limit_tcp"] = append(
				confRead.appIdent[0]["inspection_limit_tcp"].([]map[string]interface{}), map[string]interface{}{
					"byte_limit":   -1,
					"packet_limit": -1,
				})
		}
		inspLimitTCP := confRead.appIdent[0]["inspection_limit_tcp"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "inspection-limit tcp byte-limit "):
			var err error
			inspLimitTCP["byte_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "inspection-limit tcp byte-limit "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "inspection-limit tcp packet-limit "):
			var err error
			inspLimitTCP["packet_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "inspection-limit tcp packet-limit "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "inspection-limit udp"):
		if len(confRead.appIdent[0]["inspection_limit_udp"].([]map[string]interface{})) == 0 {
			confRead.appIdent[0]["inspection_limit_udp"] = append(
				confRead.appIdent[0]["inspection_limit_udp"].([]map[string]interface{}), map[string]interface{}{
					"byte_limit":   -1,
					"packet_limit": -1,
				})
		}
		inspLimitUDP := confRead.appIdent[0]["inspection_limit_udp"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "inspection-limit udp byte-limit "):
			var err error
			inspLimitUDP["byte_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "inspection-limit udp byte-limit "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "inspection-limit udp packet-limit "):
			var err error
			inspLimitUDP["packet_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "inspection-limit udp packet-limit "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "max-memory "):
		var err error
		confRead.appIdent[0]["max_memory"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-memory "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "max-transactions "):
		var err error
		confRead.appIdent[0]["max_transactions"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-transactions "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "micro-apps":
		confRead.appIdent[0]["micro_apps"] = true
	case strings.HasPrefix(itemTrim, "statistics interval "):
		var err error
		confRead.appIdent[0]["statistics_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "statistics interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	}

	return nil
}

func readServicesUserIdentification(confRead *servicesOptions, itemTrimUserIdentification string) error {
	itemTrim := strings.TrimPrefix(itemTrimUserIdentification, "user-identification ")
	if len(confRead.userIdentification) == 0 {
		confRead.userIdentification = append(confRead.userIdentification, map[string]interface{}{
			"ad_access":               make([]map[string]interface{}, 0),
			"device_info_auth_source": "",
			"identity_management":     make([]map[string]interface{}, 0),
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "active-directory-access"):
		if len(confRead.userIdentification[0]["ad_access"].([]map[string]interface{})) == 0 {
			confRead.userIdentification[0]["ad_access"] = append(
				confRead.userIdentification[0]["ad_access"].([]map[string]interface{}),
				map[string]interface{}{
					"auth_entry_timeout":           -1,
					"filter_exclude":               make([]string, 0),
					"filter_include":               make([]string, 0),
					"firewall_auth_forced_timeout": 0,
					"invalid_auth_entry_timeout":   -1,
					"no_on_demand_probe":           false,
					"wmi_timeout":                  0,
				})
		}
		adAccess := confRead.userIdentification[0]["ad_access"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "active-directory-access authentication-entry-timeout "):
			var err error
			adAccess["auth_entry_timeout"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "active-directory-access authentication-entry-timeout "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "active-directory-access filter exclude "):
			adAccess["filter_exclude"] = append(adAccess["filter_exclude"].([]string), strings.TrimPrefix(
				itemTrim, "active-directory-access filter exclude "))
		case strings.HasPrefix(itemTrim, "active-directory-access filter include "):
			adAccess["filter_include"] = append(adAccess["filter_include"].([]string), strings.TrimPrefix(
				itemTrim, "active-directory-access filter include "))
		case strings.HasPrefix(itemTrim, "active-directory-access firewall-authentication-forced-timeout "):
			var err error
			adAccess["firewall_auth_forced_timeout"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "active-directory-access firewall-authentication-forced-timeout "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "active-directory-access invalid-authentication-entry-timeout "):
			var err error
			adAccess["invalid_auth_entry_timeout"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "active-directory-access invalid-authentication-entry-timeout "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case itemTrim == "active-directory-access no-on-demand-probe":
			adAccess["no_on_demand_probe"] = true
		case strings.HasPrefix(itemTrim, "active-directory-access wmi-timeout "):
			var err error
			adAccess["wmi_timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "active-directory-access wmi-timeout "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "device-information authentication-source "):
		confRead.userIdentification[0]["device_info_auth_source"] = strings.TrimPrefix(
			itemTrim, "device-information authentication-source ")
	case strings.HasPrefix(itemTrim, "identity-management "):
		if len(confRead.userIdentification[0]["identity_management"].([]map[string]interface{})) == 0 {
			confRead.userIdentification[0]["identity_management"] = append(
				confRead.userIdentification[0]["identity_management"].([]map[string]interface{}),
				map[string]interface{}{
					"connection":                           make([]map[string]interface{}, 0),
					"authentication_entry_timeout":         -1,
					"batch_query_items_per_batch":          0,
					"batch_query_interval":                 0,
					"filter_domain":                        make([]string, 0),
					"filter_exclude_ip_address_book":       "",
					"filter_exclude_ip_address_set":        "",
					"filter_include_ip_address_book":       "",
					"filter_include_ip_address_set":        "",
					"invalid_authentication_entry_timeout": -1,
					"ip_query_disable":                     false,
					"ip_query_delay_time":                  -1,
				})
		}
		userIdentIdentityMgmt := confRead.userIdentification[0]["identity_management"].([]map[string]interface{})[0]
		itemTrimIdentMgmt := strings.TrimPrefix(itemTrim, "identity-management ")
		switch {
		case strings.HasPrefix(itemTrimIdentMgmt, "authentication-entry-timeout "):
			var err error
			userIdentIdentityMgmt["authentication_entry_timeout"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimIdentMgmt, "authentication-entry-timeout "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrimIdentMgmt, "batch-query items-per-batch "):
			var err error
			userIdentIdentityMgmt["batch_query_items_per_batch"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimIdentMgmt, "batch-query items-per-batch "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrimIdentMgmt, "batch-query query-interval "):
			var err error
			userIdentIdentityMgmt["batch_query_interval"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrimIdentMgmt, "batch-query query-interval "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrimIdentMgmt, "connection "):
			if len(userIdentIdentityMgmt["connection"].([]map[string]interface{})) == 0 {
				userIdentIdentityMgmt["connection"] = append(
					userIdentIdentityMgmt["connection"].([]map[string]interface{}), map[string]interface{}{
						"primary_address":          "",
						"primary_client_id":        "",
						"primary_client_secret":    "",
						"connect_method":           "",
						"port":                     0,
						"primary_ca_certificate":   "",
						"query_api":                "",
						"secondary_address":        "",
						"secondary_ca_certificate": "",
						"secondary_client_id":      "",
						"secondary_client_secret":  "",
						"token_api":                "",
					})
			}
			userIdentIdentityMgmtConnect := userIdentIdentityMgmt["connection"].([]map[string]interface{})[0]
			switch {
			case strings.HasPrefix(itemTrimIdentMgmt, "connection primary address "):
				userIdentIdentityMgmtConnect["primary_address"] = strings.TrimPrefix(
					itemTrimIdentMgmt, "connection primary address ")
			case strings.HasPrefix(itemTrimIdentMgmt, "connection primary client-id "):
				userIdentIdentityMgmtConnect["primary_client_id"] = strings.Trim(strings.TrimPrefix(
					itemTrimIdentMgmt, "connection primary client-id "), "\"")
			case strings.HasPrefix(itemTrimIdentMgmt, "connection primary client-secret "):
				var err error
				userIdentIdentityMgmtConnect["primary_client_secret"], err = jdecode.Decode(
					strings.Trim(strings.TrimPrefix(itemTrimIdentMgmt, "connection primary client-secret "), "\""))
				if err != nil {
					return fmt.Errorf("failed to decode primary client-secret: %w", err)
				}
			case strings.HasPrefix(itemTrimIdentMgmt, "connection connect-method "):
				userIdentIdentityMgmtConnect["connect_method"] = strings.TrimPrefix(itemTrimIdentMgmt, "connection connect-method ")
			case strings.HasPrefix(itemTrimIdentMgmt, "connection port "):
				var err error
				userIdentIdentityMgmtConnect["port"], err = strconv.Atoi(strings.TrimPrefix(itemTrimIdentMgmt, "connection port "))
				if err != nil {
					return fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrimIdentMgmt, "connection primary ca-certificate "):
				userIdentIdentityMgmtConnect["primary_ca_certificate"] = strings.Trim(strings.TrimPrefix(
					itemTrimIdentMgmt, "connection primary ca-certificate "), "\"")
			case strings.HasPrefix(itemTrimIdentMgmt, "connection query-api "):
				userIdentIdentityMgmtConnect["query_api"] = strings.Trim(strings.TrimPrefix(
					itemTrimIdentMgmt, "connection query-api "), "\"")
			case strings.HasPrefix(itemTrimIdentMgmt, "connection secondary address "):
				userIdentIdentityMgmtConnect["secondary_address"] = strings.TrimPrefix(
					itemTrimIdentMgmt, "connection secondary address ")
			case strings.HasPrefix(itemTrimIdentMgmt, "connection secondary ca-certificate "):
				userIdentIdentityMgmtConnect["secondary_ca_certificate"] = strings.Trim(strings.TrimPrefix(
					itemTrimIdentMgmt, "connection secondary ca-certificate "), "\"")
			case strings.HasPrefix(itemTrimIdentMgmt, "connection secondary client-id "):
				userIdentIdentityMgmtConnect["secondary_client_id"] = strings.Trim(strings.TrimPrefix(
					itemTrimIdentMgmt, "connection secondary client-id "), "\"")
			case strings.HasPrefix(itemTrimIdentMgmt, "connection secondary client-secret "):
				var err error
				userIdentIdentityMgmtConnect["secondary_client_secret"], err = jdecode.Decode(
					strings.Trim(strings.TrimPrefix(itemTrimIdentMgmt, "connection secondary client-secret "), "\""))
				if err != nil {
					return fmt.Errorf("failed to decode secondary client-secret: %w", err)
				}
			case strings.HasPrefix(itemTrimIdentMgmt, "connection token-api "):
				userIdentIdentityMgmtConnect["token_api"] = strings.Trim(strings.TrimPrefix(
					itemTrimIdentMgmt, "connection token-api "), "\"")
			}
		case strings.HasPrefix(itemTrimIdentMgmt, "filter domain "):
			userIdentIdentityMgmt["filter_domain"] = append(userIdentIdentityMgmt["filter_domain"].([]string),
				strings.TrimPrefix(itemTrimIdentMgmt, "filter domain "))
		case strings.HasPrefix(itemTrimIdentMgmt, "filter exclude-ip address-book "):
			userIdentIdentityMgmt["filter_exclude_ip_address_book"] = strings.Trim(strings.TrimPrefix(
				itemTrimIdentMgmt, "filter exclude-ip address-book "), "\"")
		case strings.HasPrefix(itemTrimIdentMgmt, "filter exclude-ip address-set "):
			userIdentIdentityMgmt["filter_exclude_ip_address_set"] = strings.Trim(strings.TrimPrefix(
				itemTrimIdentMgmt, "filter exclude-ip address-set "), "\"")
		case strings.HasPrefix(itemTrimIdentMgmt, "filter include-ip address-book "):
			userIdentIdentityMgmt["filter_include_ip_address_book"] = strings.Trim(strings.TrimPrefix(
				itemTrimIdentMgmt, "filter include-ip address-book "), "\"")
		case strings.HasPrefix(itemTrimIdentMgmt, "filter include-ip address-set "):
			userIdentIdentityMgmt["filter_include_ip_address_set"] = strings.Trim(strings.TrimPrefix(
				itemTrimIdentMgmt, "filter include-ip address-set "), "\"")
		case strings.HasPrefix(itemTrimIdentMgmt, "invalid-authentication-entry-timeout "):
			var err error
			userIdentIdentityMgmt["invalid_authentication_entry_timeout"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrimIdentMgmt, "invalid-authentication-entry-timeout "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case itemTrimIdentMgmt == "ip-query no-ip-query":
			userIdentIdentityMgmt["ip_query_disable"] = true
		case strings.HasPrefix(itemTrimIdentMgmt, "ip-query query-delay-time "):
			var err error
			userIdentIdentityMgmt["ip_query_delay_time"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrimIdentMgmt, "ip-query query-delay-time "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	}

	return nil
}

func fillServices(d *schema.ResourceData, servicesOptions servicesOptions) {
	if tfErr := d.Set("advanced_anti_malware", servicesOptions.advAntiMalware); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("application_identification", servicesOptions.appIdent); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_intelligence", servicesOptions.securityIntelligence); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("user_identification", servicesOptions.userIdentification); tfErr != nil {
		panic(tfErr)
	}
}
