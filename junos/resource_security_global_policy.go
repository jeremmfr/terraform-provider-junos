package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type globalPolicyOptions struct {
	policy []map[string]interface{}
}

func resourceSecurityGlobalPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGlobalPolicyCreate,
		ReadContext:   resourceSecurityGlobalPolicyRead,
		UpdateContext: resourceSecurityGlobalPolicyUpdate,
		DeleteContext: resourceSecurityGlobalPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityGlobalPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"policy": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, FormatDefault),
						},
						"match_source_address": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_destination_address": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_application": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_from_zone": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_to_zone": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"then": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      permitWord,
							ValidateFunc: validation.StringInSlice([]string{permitWord, "reject", "deny"}, false),
						},
						"count": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"log_init": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"log_close": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"match_destination_address_excluded": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"match_dynamic_application": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_source_address_excluded": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"permit_application_services": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application_firewall_rule_set": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"application_traffic_control_rule_set": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"gprs_gtp_profile": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"gprs_sctp_profile": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"idp": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"redirect_wx": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"reverse_redirect_wx": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"security_intelligence_policy": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ssl_proxy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"profile_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"uac_policy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"captive_portal": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"utm_policy": {
										Type:     schema.TypeString,
										Optional: true,
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

func resourceSecurityGlobalPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityGlobalPolicy(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("security_global_policy")

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security policies global not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	glbPolicy, err := readSecurityGlobalPolicy(m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if len(glbPolicy.policy) != 0 {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security policies global already set"))
	}

	if err := setSecurityGlobalPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_global_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("security_global_policy")

	return append(diagWarns, resourceSecurityGlobalPolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityGlobalPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityGlobalPolicyReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityGlobalPolicyReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	globalPolicyOptions, err := readSecurityGlobalPolicy(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSecurityGlobalPolicyData(d, globalPolicyOptions)

	return nil
}
func resourceSecurityGlobalPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := delSecurityGlobalPolicy(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	if err := setSecurityGlobalPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_global_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityGlobalPolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityGlobalPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityGlobalPolicy(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_global_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSecurityGlobalPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	globalPolicyOptions, err := readSecurityGlobalPolicy(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityGlobalPolicyData(d, globalPolicyOptions)
	result[0] = d

	return result, nil
}

func setSecurityGlobalPolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security policies global policy "
	for _, v := range d.Get("policy").([]interface{}) {
		policy := v.(map[string]interface{})
		setPrefixPolicy := setPrefix + policy["name"].(string)
		if len(policy["match_source_address"].([]interface{})) != 0 {
			for _, address := range policy["match_source_address"].([]interface{}) {
				configSet = append(configSet, setPrefixPolicy+" match source-address "+address.(string))
			}
		} else {
			configSet = append(configSet, setPrefixPolicy+" match source-address any")
		}
		if len(policy["match_destination_address"].([]interface{})) != 0 {
			for _, address := range policy["match_destination_address"].([]interface{}) {
				configSet = append(configSet, setPrefixPolicy+" match destination-address "+address.(string))
			}
		} else {
			configSet = append(configSet, setPrefixPolicy+" match destination-address any")
		}
		if len(policy["match_application"].([]interface{})) != 0 {
			for _, app := range policy["match_application"].([]interface{}) {
				configSet = append(configSet, setPrefixPolicy+" match application "+app.(string))
			}
		} else {
			configSet = append(configSet, setPrefixPolicy+" match application any")
		}
		for _, v := range policy["match_from_zone"].([]interface{}) {
			configSet = append(configSet, setPrefixPolicy+" match from-zone "+v.(string))
		}
		for _, v := range policy["match_to_zone"].([]interface{}) {
			configSet = append(configSet, setPrefixPolicy+" match to-zone "+v.(string))
		}
		configSet = append(configSet, setPrefixPolicy+" then "+policy["then"].(string))
		if policy["count"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" then count")
		}
		if policy["log_init"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" then log session-init")
		}
		if policy["log_close"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" then log session-close")
		}
		if policy["match_destination_address_excluded"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" match destination-address-excluded")
		}
		for _, v := range policy["match_dynamic_application"].([]interface{}) {
			configSet = append(configSet, setPrefixPolicy+" match dynamic-application "+v.(string))
		}
		if policy["match_source_address_excluded"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" match source-address-excluded")
		}
		if len(policy["permit_application_services"].([]interface{})) > 0 {
			if policy["permit_application_services"].([]interface{})[0] == nil {
				return fmt.Errorf("permit_application_services block is empty")
			}
			if policy["then"].(string) != permitWord {
				return fmt.Errorf("conflict policy then %v and policy permit_application_services",
					policy["then"].(string))
			}
			configSetAppSvc, err := setGlobalPolicyPermitApplicationServices(setPrefixPolicy,
				policy["permit_application_services"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetAppSvc...)
		}
	}

	return sess.configSet(configSet, jnprSess)
}
func readSecurityGlobalPolicy(m interface{}, jnprSess *NetconfObject) (globalPolicyOptions, error) {
	sess := m.(*Session)
	var confRead globalPolicyOptions

	policyConfig, err := sess.command("show configuration security policies global | display set relative ", jnprSess)
	if err != nil {
		return confRead, err
	}
	policyList := make([]map[string]interface{}, 0)
	if policyConfig != emptyWord {
		for _, item := range strings.Split(policyConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.HasPrefix(itemTrim, "policy ") {
				policyLineCut := strings.Split(itemTrim, " ")
				m := genMapGlobalPolicyWithName(policyLineCut[1])
				m, policyList = copyAndRemoveItemMapList("name", false, m, policyList)
				itemTrimPolicy := strings.TrimPrefix(itemTrim, "policy "+policyLineCut[1]+" ")
				switch {
				case strings.HasPrefix(itemTrimPolicy, "match source-address "):
					m["match_source_address"] = append(m["match_source_address"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match source-address "))
				case strings.HasPrefix(itemTrimPolicy, "match destination-address "):
					m["match_destination_address"] = append(m["match_destination_address"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match destination-address "))
				case strings.HasPrefix(itemTrimPolicy, "match application "):
					m["match_application"] = append(m["match_application"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match application "))
				case strings.HasPrefix(itemTrimPolicy, "match from-zone "):
					m["match_from_zone"] = append(m["match_from_zone"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match from-zone "))
				case strings.HasPrefix(itemTrimPolicy, "match to-zone "):
					m["match_to_zone"] = append(m["match_to_zone"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match to-zone "))
				case strings.HasPrefix(itemTrimPolicy, "match destination-address-excluded"):
					m["match_destination_address_excluded"] = true
				case strings.HasPrefix(itemTrimPolicy, "match dynamic-application "):
					m["match_dynamic_application"] = append(m["match_dynamic_application"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match dynamic-application "))
				case strings.HasPrefix(itemTrimPolicy, "match source-address-excluded"):
					m["match_source_address_excluded"] = true
				case strings.HasPrefix(itemTrimPolicy, "then "):
					switch {
					case strings.HasSuffix(itemTrimPolicy, permitWord),
						strings.HasSuffix(itemTrimPolicy, "deny"),
						strings.HasSuffix(itemTrimPolicy, "reject"):
						m["then"] = strings.TrimPrefix(itemTrimPolicy, "then ")
					case itemTrimPolicy == "then count":
						m["count"] = true
					case itemTrimPolicy == "then log session-init":
						m["log_init"] = true
					case itemTrimPolicy == "then log session-close":
						m["log_close"] = true
					case strings.HasPrefix(itemTrimPolicy, "then permit application-services"):
						m["then"] = permitWord
						m["permit_application_services"] = readGlobalPolicyPermitApplicationServices(itemTrimPolicy,
							m["permit_application_services"])
					}
				}
				policyList = append(policyList, m)
			}
		}
	}
	confRead.policy = policyList

	return confRead, nil
}
func delSecurityGlobalPolicy(m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security policies global")

	return sess.configSet(configSet, jnprSess)
}

func fillSecurityGlobalPolicyData(d *schema.ResourceData, globalPolicyOptions globalPolicyOptions) {
	if tfErr := d.Set("policy", globalPolicyOptions.policy); tfErr != nil {
		panic(tfErr)
	}
}

func genMapGlobalPolicyWithName(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":                               name,
		"match_source_address":               make([]string, 0),
		"match_destination_address":          make([]string, 0),
		"match_application":                  make([]string, 0),
		"match_from_zone":                    make([]string, 0),
		"match_to_zone":                      make([]string, 0),
		"then":                               "",
		"count":                              false,
		"log_init":                           false,
		"log_close":                          false,
		"match_destination_address_excluded": false,
		"match_dynamic_application":          make([]string, 0),
		"match_source_address_excluded":      false,
		"permit_application_services":        make([]map[string]interface{}, 0),
	}
}

func readGlobalPolicyPermitApplicationServices(itemTrimPolicy string,
	permitApplicationServices interface{}) []map[string]interface{} {
	applicationServices := map[string]interface{}{
		"application_firewall_rule_set":        "",
		"application_traffic_control_rule_set": "",
		"gprs_gtp_profile":                     "",
		"gprs_sctp_profile":                    "",
		"idp":                                  false,
		"redirect_wx":                          false,
		"reverse_redirect_wx":                  false,
		"security_intelligence_policy":         "",
		"ssl_proxy":                            make([]map[string]interface{}, 0, 1),
		"uac_policy":                           make([]map[string]interface{}, 0, 1),
		"utm_policy":                           "",
	}
	if len(permitApplicationServices.([]map[string]interface{})) > 0 {
		for k, v := range permitApplicationServices.([]map[string]interface{})[0] {
			applicationServices[k] = v
		}
	}
	itemTrimPolicyPermitAppSvc := strings.TrimPrefix(itemTrimPolicy, "then permit application-services ")
	switch {
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "application-firewall rule-set "):
		applicationServices["application_firewall_rule_set"] = strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc,
			"application-firewall rule-set "), "\"")
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "application-traffic-control rule-set "):
		applicationServices["application_traffic_control_rule_set"] = strings.Trim(
			strings.TrimPrefix(itemTrimPolicyPermitAppSvc, "application-traffic-control rule-set "), "\"")
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "gprs-gtp-profile "):
		applicationServices["gprs_gtp_profile"] = strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc,
			"gprs-gtp-profile "), "\"")
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "gprs-sctp-profile "):
		applicationServices["gprs_sctp_profile"] = strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc,
			"gprs-sctp-profile "), "\"")
	case itemTrimPolicyPermitAppSvc == "idp":
		applicationServices["idp"] = true
	case itemTrimPolicyPermitAppSvc == "redirect-wx":
		applicationServices["redirect_wx"] = true
	case itemTrimPolicyPermitAppSvc == "reverse-redirect-wx":
		applicationServices["reverse_redirect_wx"] = true
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "security-intelligence-policy "):
		applicationServices["security_intelligence_policy"] = strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc,
			"security-intelligence-policy "), "\"")
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "ssl-proxy"):
		if strings.HasPrefix(itemTrimPolicyPermitAppSvc, "ssl-proxy profile-name ") {
			applicationServices["ssl_proxy"] = append(applicationServices["ssl_proxy"].([]map[string]interface{}),
				map[string]interface{}{
					"profile_name": strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc, "ssl-proxy profile-name "), "\""),
				})
		} else {
			applicationServices["ssl_proxy"] = append(applicationServices["ssl_proxy"].([]map[string]interface{}),
				map[string]interface{}{
					"profile_name": "",
				})
		}
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "uac-policy"):
		if strings.HasPrefix(itemTrimPolicyPermitAppSvc, "uac-policy captive-portal ") {
			applicationServices["uac_policy"] = append(applicationServices["uac_policy"].([]map[string]interface{}),
				map[string]interface{}{
					"captive_portal": strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc, "uac-policy captive-portal "), "\""),
				})
		} else {
			applicationServices["uac_policy"] = append(applicationServices["uac_policy"].([]map[string]interface{}),
				map[string]interface{}{
					"captive_portal": "",
				})
		}
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "utm-policy "):
		applicationServices["utm_policy"] = strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc, "utm-policy "), "\"")
	}

	// override (maxItem = 1)
	return []map[string]interface{}{applicationServices}
}

func setGlobalPolicyPermitApplicationServices(setPrefixPolicy string,
	policyPermitApplicationServices map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	setPrefixPolicyPermitAppSvc := setPrefixPolicy + " then permit application-services "
	if policyPermitApplicationServices["application_firewall_rule_set"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"application-firewall rule-set \""+
			policyPermitApplicationServices["application_firewall_rule_set"].(string)+"\"")
	}
	if policyPermitApplicationServices["application_traffic_control_rule_set"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"application-traffic-control rule-set \""+
			policyPermitApplicationServices["application_traffic_control_rule_set"].(string)+"\"")
	}
	if policyPermitApplicationServices["gprs_gtp_profile"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"gprs-gtp-profile \""+policyPermitApplicationServices["gprs_gtp_profile"].(string)+"\"")
	}
	if policyPermitApplicationServices["gprs_sctp_profile"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"gprs-sctp-profile \""+policyPermitApplicationServices["gprs_sctp_profile"].(string)+"\"")
	}
	if policyPermitApplicationServices["idp"].(bool) {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"idp")
	}
	if policyPermitApplicationServices["redirect_wx"].(bool) &&
		policyPermitApplicationServices["reverse_redirect_wx"].(bool) {
		return configSet, fmt.Errorf("conflict redirect_wx and reverse_redirect_wx enabled both")
	}
	if policyPermitApplicationServices["redirect_wx"].(bool) {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"redirect-wx")
	}
	if policyPermitApplicationServices["reverse_redirect_wx"].(bool) {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"reverse-redirect-wx")
	}
	if policyPermitApplicationServices["security_intelligence_policy"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"security-intelligence-policy \""+
			policyPermitApplicationServices["security_intelligence_policy"].(string)+"\"")
	}
	if len(policyPermitApplicationServices["ssl_proxy"].([]interface{})) > 0 {
		if policyPermitApplicationServices["ssl_proxy"].([]interface{})[0] != nil {
			policyPermitApplicationServicesSSLProxy :=
				policyPermitApplicationServices["ssl_proxy"].([]interface{})[0].(map[string]interface{})
			if policyPermitApplicationServicesSSLProxy["profile_name"].(string) != "" {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+
					"ssl-proxy profile-name \""+
					policyPermitApplicationServicesSSLProxy["profile_name"].(string)+"\"")
			} else {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+"ssl-proxy")
			}
		} else {
			configSet = append(configSet, setPrefixPolicyPermitAppSvc+"ssl-proxy")
		}
	}
	if len(policyPermitApplicationServices["uac_policy"].([]interface{})) > 0 {
		if policyPermitApplicationServices["uac_policy"].([]interface{})[0] != nil {
			policyPermitApplicationServicesUACPolicy :=
				policyPermitApplicationServices["uac_policy"].([]interface{})[0].(map[string]interface{})
			if policyPermitApplicationServicesUACPolicy["captive_portal"].(string) != "" {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+
					"uac-policy captive-portal \""+
					policyPermitApplicationServicesUACPolicy["captive_portal"].(string)+"\"")
			} else {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+"uac-policy")
			}
		} else {
			configSet = append(configSet, setPrefixPolicyPermitAppSvc+"uac-policy")
		}
	}
	if policyPermitApplicationServices["utm_policy"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			"utm-policy \""+policyPermitApplicationServices["utm_policy"].(string)+"\"")
	}

	return configSet, nil
}
