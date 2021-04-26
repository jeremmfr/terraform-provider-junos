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
	jdecode "github.com/jeremmfr/junosdecode"
)

type servicesOptions struct {
	appIdent             []map[string]interface{}
	securityIntelligence []map[string]interface{}
}

func resourceServices() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicesCreate,
		ReadContext:   resourceServicesRead,
		UpdateContext: resourceServicesUpdate,
		DeleteContext: resourceServicesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServicesImport,
		},
		Schema: map[string]*schema.Schema{
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
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^[a-zA-Z0-9]{32}$`),
								"Auth token must be consisted of 32 alphanumeric characters"),
							ConflictsWith: []string{"security_intelligence.0.authentication_tls_profile"},
						},
						"authentication_tls_profile": {
							Type:          schema.TypeString,
							Optional:      true,
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
						},
						"url_parameter": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceServicesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setServices(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("services")

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := setServices(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_services", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("services")

	return append(diagWarns, resourceServicesReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceServicesReadWJnprSess(d, m, jnprSess)
}

func resourceServicesReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	servicesOptions, err := readServices(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillServices(d, servicesOptions)

	return nil
}

func resourceServicesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delServices(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServices(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_services", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceServicesImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	servicesOptions, err := readServices(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillServices(d, servicesOptions)
	d.SetId("services")
	result[0] = d

	return result, nil
}

func setServices(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	// setPrefix := "set services "
	configSet := make([]string, 0)

	for _, v := range d.Get("application_identification").([]interface{}) {
		configSetApplicationIdentification, err := setServicesApplicationIdentification(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetApplicationIdentification...)
	}
	for _, v := range d.Get("security_intelligence").([]interface{}) {
		configSetSecurityIntel, err := setServicesSecurityIntell(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetSecurityIntel...)
	}

	return sess.configSet(configSet, jnprSess)
}

func setServicesApplicationIdentification(appID interface{}) ([]string, error) {
	setPrefix := "set services application-identification "
	configSet := make([]string, 0)
	if appID != nil {
		appIDM := appID.(map[string]interface{})
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
	} else {
		return configSet, fmt.Errorf("application_identification block is empty")
	}

	return configSet, nil
}

func setServicesSecurityIntell(secuIntel interface{}) ([]string, error) {
	setPrefix := "set services security-intelligence "
	configSet := make([]string, 0)
	if secuIntel != nil {
		secuIntelM := secuIntel.(map[string]interface{})
		if v := secuIntelM["authentication_token"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication auth-token "+v)
		}
		if v := secuIntelM["authentication_tls_profile"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication tls-profile \""+v+"\"")
		}
		for _, v := range secuIntelM["category_disable"].(*schema.Set).List() {
			if v.(string) == "all" {
				configSet = append(configSet, setPrefix+"category all disable")
			} else {
				configSet = append(configSet, setPrefix+"category category-name "+v.(string)+" disable")
			}
		}
		for _, v := range secuIntelM["default_policy"].([]interface{}) {
			defPolicy := v.(map[string]interface{})
			configSet = append(configSet, setPrefix+"default-policy "+
				defPolicy["category_name"].(string)+" "+defPolicy["profile_name"].(string))
		}
		if v := secuIntelM["proxy_profile"].(string); v != "" {
			configSet = append(configSet, setPrefix+"proxy-profile \""+v+"\"")
		}
		if v := secuIntelM["url"].(string); v != "" {
			configSet = append(configSet, setPrefix+"url \""+v+"\"")
		}
		if v := secuIntelM["url_parameter"].(string); v != "" {
			configSet = append(configSet, setPrefix+"url-parameter \""+v+"\"")
		}
	} else {
		return configSet, fmt.Errorf("security_intelligence block is empty")
	}

	return configSet, nil
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
		"security-intelligence authentication",
		"security-intelligence category",
		"security-intelligence default-policy",
		"security-intelligence proxy-profile",
		"security-intelligence url",
		"security-intelligence url-parameter",
	}
}

func delServices(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, listLinesServicesApplicationIdentification()...)
	listLinesToDelete = append(listLinesToDelete, listLinesServicesSecurityIntel()...)
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete services "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return sess.configSet(configSet, jnprSess)
}

func readServices(m interface{}, jnprSess *NetconfObject) (servicesOptions, error) {
	sess := m.(*Session)
	var confRead servicesOptions

	servicesConfig, err := sess.command("show configuration services"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if servicesConfig != emptyWord {
		for _, item := range strings.Split(servicesConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case checkStringHasPrefixInList(itemTrim, listLinesServicesApplicationIdentification()):
				if err := readServicesApplicationIdentification(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case checkStringHasPrefixInList(itemTrim, listLinesServicesSecurityIntel()):
				if err := readServicesSecurityIntel(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
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
		confRead.securityIntelligence[0]["url_parameter"], err =
			jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim, "url-parameter "), "\""))
		if err != nil {
			return fmt.Errorf("failed to decode url-parameter : %w", err)
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
		confRead.appIdent[0]["application_system_cache_timeout"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "application-system-cache-timeout "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
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
		switch {
		case strings.HasPrefix(itemTrim, "download automatic interval "):
			var err error
			confRead.appIdent[0]["download"].([]map[string]interface{})[0]["automatic_interval"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "download automatic interval "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "download automatic start-time "):
			confRead.appIdent[0]["download"].([]map[string]interface{})[0]["automatic_start_time"] =
				strings.TrimPrefix(itemTrim, "download automatic start-time ")
		case itemTrim == "download ignore-server-validation":
			confRead.appIdent[0]["download"].([]map[string]interface{})[0]["ignore_server_validation"] =
				true
		case strings.HasPrefix(itemTrim, "download proxy-profile "):
			confRead.appIdent[0]["download"].([]map[string]interface{})[0]["proxy_profile"] =
				strings.Trim(strings.TrimPrefix(itemTrim, "download proxy-profile "), "\"")
		case strings.HasPrefix(itemTrim, "download url "):
			confRead.appIdent[0]["download"].([]map[string]interface{})[0]["url"] =
				strings.Trim(strings.TrimPrefix(itemTrim, "download url "), "\"")
		}
	case strings.HasPrefix(itemTrim, "enable-performance-mode"):
		if len(confRead.appIdent[0]["enable_performance_mode"].([]map[string]interface{})) == 0 {
			confRead.appIdent[0]["enable_performance_mode"] = append(
				confRead.appIdent[0]["enable_performance_mode"].([]map[string]interface{}), map[string]interface{}{
					"max_packet_threshold": 0,
				})
		}
		if strings.HasPrefix(itemTrim, "enable-performance-mode max-packet-threshold ") {
			var err error
			confRead.appIdent[0]["enable_performance_mode"].([]map[string]interface{})[0]["max_packet_threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "enable-performance-mode max-packet-threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "global-offload-byte-limit "):
		var err error
		confRead.appIdent[0]["global_offload_byte_limit"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "global-offload-byte-limit "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "imap-cache-size "):
		var err error
		confRead.appIdent[0]["imap_cache_size"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "imap-cache-size "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "imap-cache-timeout "):
		var err error
		confRead.appIdent[0]["imap_cache_timeout"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "imap-cache-timeout "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "inspection-limit tcp"):
		if len(confRead.appIdent[0]["inspection_limit_tcp"].([]map[string]interface{})) == 0 {
			confRead.appIdent[0]["inspection_limit_tcp"] = append(
				confRead.appIdent[0]["inspection_limit_tcp"].([]map[string]interface{}), map[string]interface{}{
					"byte_limit":   -1,
					"packet_limit": -1,
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "inspection-limit tcp byte-limit "):
			var err error
			confRead.appIdent[0]["inspection_limit_tcp"].([]map[string]interface{})[0]["byte_limit"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "inspection-limit tcp byte-limit "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "inspection-limit tcp packet-limit "):
			var err error
			confRead.appIdent[0]["inspection_limit_tcp"].([]map[string]interface{})[0]["packet_limit"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "inspection-limit tcp packet-limit "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
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
		switch {
		case strings.HasPrefix(itemTrim, "inspection-limit udp byte-limit "):
			var err error
			confRead.appIdent[0]["inspection_limit_udp"].([]map[string]interface{})[0]["byte_limit"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "inspection-limit udp byte-limit "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "inspection-limit udp packet-limit "):
			var err error
			confRead.appIdent[0]["inspection_limit_udp"].([]map[string]interface{})[0]["packet_limit"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "inspection-limit udp packet-limit "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "max-memory "):
		var err error
		confRead.appIdent[0]["max_memory"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "max-memory "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "max-transactions "):
		var err error
		confRead.appIdent[0]["max_transactions"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "max-transactions "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case itemTrim == "micro-apps":
		confRead.appIdent[0]["micro_apps"] = true
	case strings.HasPrefix(itemTrim, "statistics interval "):
		var err error
		confRead.appIdent[0]["statistics_interval"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "statistics interval "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func fillServices(d *schema.ResourceData, servicesOptions servicesOptions) {
	if tfErr := d.Set("application_identification", servicesOptions.appIdent); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("security_intelligence", servicesOptions.securityIntelligence); tfErr != nil {
		panic(tfErr)
	}
}
