package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
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
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"match_source_address": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_destination_address": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_from_zone": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_to_zone": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
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
						"match_application": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_destination_address_excluded": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"match_dynamic_application": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_source_address_excluded": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"match_source_end_user_profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "profile can also be named device identity profile",
						},
						"permit_application_services": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"advanced_anti_malware_policy": {
										Type:     schema.TypeString,
										Optional: true,
									},
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
									"idp_policy": {
										Type:     schema.TypeString,
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
	var diagWarns diag.Diagnostics
	glbPolicy, err := readSecurityGlobalPolicy(m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if len(glbPolicy.policy) != 0 {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("security policies global already set"))...)
	}

	if err := setSecurityGlobalPolicy(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_global_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

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
	if sess.junosFakeUpdateAlso {
		if err := delSecurityGlobalPolicy(m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityGlobalPolicy(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityGlobalPolicy(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityGlobalPolicy(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_global_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityGlobalPolicyReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityGlobalPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delSecurityGlobalPolicy(m, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}

	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityGlobalPolicy(m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_global_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

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
	policyNameList := make([]string, 0)
	for _, v := range d.Get("policy").([]interface{}) {
		policy := v.(map[string]interface{})
		if bchk.StringInSlice(policy["name"].(string), policyNameList) {
			return fmt.Errorf("multiple blocks policy with the same name %s", policy["name"].(string))
		}
		policyNameList = append(policyNameList, policy["name"].(string))
		setPrefixPolicy := setPrefix + policy["name"].(string)
		for _, address := range sortSetOfString(policy["match_source_address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixPolicy+" match source-address "+address)
		}
		for _, address := range sortSetOfString(policy["match_destination_address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixPolicy+" match destination-address "+address)
		}
		for _, v := range sortSetOfString(policy["match_from_zone"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixPolicy+" match from-zone "+v)
		}
		for _, v := range sortSetOfString(policy["match_to_zone"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixPolicy+" match to-zone "+v)
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
		if len(policy["match_application"].(*schema.Set).List()) == 0 &&
			len(policy["match_dynamic_application"].(*schema.Set).List()) == 0 {
			return fmt.Errorf("1 minimum item must be set in 'match_application' or 'match_dynamic_application' "+
				"argument in '%s' policy", policy["name"].(string))
		}
		for _, app := range sortSetOfString(policy["match_application"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixPolicy+" match application "+app)
		}
		if policy["match_destination_address_excluded"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" match destination-address-excluded")
		}
		for _, v := range sortSetOfString(policy["match_dynamic_application"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixPolicy+" match dynamic-application "+v)
		}
		if policy["match_source_address_excluded"].(bool) {
			configSet = append(configSet, setPrefixPolicy+" match source-address-excluded")
		}
		if v := policy["match_source_end_user_profile"].(string); v != "" {
			configSet = append(configSet, setPrefixPolicy+" match source-end-user-profile \""+v+"\"")
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

	showConfig, err := sess.command(cmdShowConfig+"security policies global | display set relative ", jnprSess)
	if err != nil {
		return confRead, err
	}
	policyList := make([]map[string]interface{}, 0)
	if showConfig != emptyWord {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.HasPrefix(itemTrim, "policy ") {
				policyLineCut := strings.Split(itemTrim, " ")
				policy := genMapGlobalPolicyWithName(policyLineCut[1])
				policyList = copyAndRemoveItemMapList("name", policy, policyList)
				itemTrimPolicy := strings.TrimPrefix(itemTrim, "policy "+policyLineCut[1]+" ")
				switch {
				case strings.HasPrefix(itemTrimPolicy, "match source-address "):
					policy["match_source_address"] = append(policy["match_source_address"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match source-address "))
				case strings.HasPrefix(itemTrimPolicy, "match destination-address "):
					policy["match_destination_address"] = append(policy["match_destination_address"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match destination-address "))
				case strings.HasPrefix(itemTrimPolicy, "match application "):
					policy["match_application"] = append(policy["match_application"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match application "))
				case strings.HasPrefix(itemTrimPolicy, "match from-zone "):
					policy["match_from_zone"] = append(policy["match_from_zone"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match from-zone "))
				case strings.HasPrefix(itemTrimPolicy, "match to-zone "):
					policy["match_to_zone"] = append(policy["match_to_zone"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match to-zone "))
				case strings.HasPrefix(itemTrimPolicy, "match destination-address-excluded"):
					policy["match_destination_address_excluded"] = true
				case strings.HasPrefix(itemTrimPolicy, "match dynamic-application "):
					policy["match_dynamic_application"] = append(policy["match_dynamic_application"].([]string),
						strings.TrimPrefix(itemTrimPolicy, "match dynamic-application "))
				case strings.HasPrefix(itemTrimPolicy, "match source-address-excluded"):
					policy["match_source_address_excluded"] = true
				case strings.HasPrefix(itemTrimPolicy, "match source-end-user-profile "):
					policy["match_source_end_user_profile"] = strings.Trim(strings.TrimPrefix(
						itemTrimPolicy, "match source-end-user-profile "), "\"")
				case strings.HasPrefix(itemTrimPolicy, "then "):
					switch {
					case strings.HasSuffix(itemTrimPolicy, permitWord),
						strings.HasSuffix(itemTrimPolicy, "deny"),
						strings.HasSuffix(itemTrimPolicy, "reject"):
						policy["then"] = strings.TrimPrefix(itemTrimPolicy, "then ")
					case itemTrimPolicy == "then count":
						policy["count"] = true
					case itemTrimPolicy == "then log session-init":
						policy["log_init"] = true
					case itemTrimPolicy == "then log session-close":
						policy["log_close"] = true
					case strings.HasPrefix(itemTrimPolicy, "then permit application-services"):
						policy["then"] = permitWord
						if len(policy["permit_application_services"].([]map[string]interface{})) == 0 {
							policy["permit_application_services"] = append(
								policy["permit_application_services"].([]map[string]interface{}),
								map[string]interface{}{
									"advanced_anti_malware_policy":         "",
									"application_firewall_rule_set":        "",
									"application_traffic_control_rule_set": "",
									"gprs_gtp_profile":                     "",
									"gprs_sctp_profile":                    "",
									"idp":                                  false,
									"idp_policy":                           "",
									"redirect_wx":                          false,
									"reverse_redirect_wx":                  false,
									"security_intelligence_policy":         "",
									"ssl_proxy":                            make([]map[string]interface{}, 0, 1),
									"uac_policy":                           make([]map[string]interface{}, 0, 1),
									"utm_policy":                           "",
								})
						}
						readGlobalPolicyPermitApplicationServices(itemTrimPolicy,
							policy["permit_application_services"].([]map[string]interface{})[0])
					}
				}
				policyList = append(policyList, policy)
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
		"match_source_end_user_profile":      "",
		"permit_application_services":        make([]map[string]interface{}, 0),
	}
}

func readGlobalPolicyPermitApplicationServices(itemTrimPolicy string, applicationServices map[string]interface{}) {
	itemTrimPolicyPermitAppSvc := strings.TrimPrefix(itemTrimPolicy, "then permit application-services ")
	switch {
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "advanced-anti-malware-policy "):
		applicationServices["advanced_anti_malware_policy"] = strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc,
			"advanced-anti-malware-policy "), "\"")
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
	case strings.HasPrefix(itemTrimPolicyPermitAppSvc, "idp-policy "):
		applicationServices["idp_policy"] = strings.Trim(strings.TrimPrefix(itemTrimPolicyPermitAppSvc, "idp-policy "), "\"")
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
}

func setGlobalPolicyPermitApplicationServices(setPrefixPolicy string,
	policyPermitApplicationServices map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	setPrefixPolicyPermitAppSvc := setPrefixPolicy + " then permit application-services "
	if v := policyPermitApplicationServices["advanced_anti_malware_policy"].(string); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"advanced-anti-malware-policy \""+v+"\"")
	}
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
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"idp")
	}
	if v := policyPermitApplicationServices["idp_policy"].(string); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"idp-policy \""+v+"\"")
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
			sslProxy := policyPermitApplicationServices["ssl_proxy"].([]interface{})[0].(map[string]interface{})
			if sslProxy["profile_name"].(string) != "" {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+
					"ssl-proxy profile-name \""+sslProxy["profile_name"].(string)+"\"")
			} else {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+"ssl-proxy")
			}
		} else {
			configSet = append(configSet, setPrefixPolicyPermitAppSvc+"ssl-proxy")
		}
	}
	if len(policyPermitApplicationServices["uac_policy"].([]interface{})) > 0 {
		if policyPermitApplicationServices["uac_policy"].([]interface{})[0] != nil {
			uacPolicy := policyPermitApplicationServices["uac_policy"].([]interface{})[0].(map[string]interface{})
			if uacPolicy["captive_portal"].(string) != "" {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+
					"uac-policy captive-portal \""+uacPolicy["captive_portal"].(string)+"\"")
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
