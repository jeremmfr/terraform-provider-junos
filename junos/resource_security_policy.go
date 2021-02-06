package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type policyOptions struct {
	fromZone string
	toZone   string
	policy   []map[string]interface{}
}

func resourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPolicyCreate,
		ReadContext:   resourceSecurityPolicyRead,
		UpdateContext: resourceSecurityPolicyUpdate,
		DeleteContext: resourceSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"from_zone": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"to_zone": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"policy": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
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
						"permit_tunnel_ipsec_vpn": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSecurityPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security policy not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	securityPolicyExists, err := checkSecurityPolicyExists(d.Get("from_zone").(string), d.Get("to_zone").(string),
		m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if securityPolicyExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security policy from %v to %v already exists",
			d.Get("from_zone").(string), d.Get("to_zone").(string)))
	}

	if err := setSecurityPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityPolicyExists, err = checkSecurityPolicyExists(d.Get("from_zone").(string), d.Get("to_zone").(string),
		m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityPolicyExists {
		d.SetId(d.Get("from_zone").(string) + idSeparator + d.Get("to_zone").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security policy from %v to %v not exists after commit "+
			"=> check your config", d.Get("from_zone").(string), d.Get("to_zone").(string)))...)
	}

	return append(diagWarns, resourceSecurityPolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityPolicyReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityPolicyReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	policyOptions, err := readSecurityPolicy(d.Get("from_zone").(string)+idSeparator+d.Get("to_zone").(string),
		m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if len(policyOptions.policy) == 0 {
		d.SetId("")
	} else {
		fillSecurityPolicyData(d, policyOptions)
	}

	return nil
}
func resourceSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := delSecurityPolicy(d.Get("from_zone").(string), d.Get("to_zone").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	if err := setSecurityPolicy(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityPolicyReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityPolicy(d.Get("from_zone").(string), d.Get("to_zone").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_policy", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSecurityPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), idSeparator)
	if len(idList) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	securityPolicyExists, err := checkSecurityPolicyExists(idList[0], idList[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityPolicyExists {
		return nil, fmt.Errorf("don't find policy with id '%v' (id must be <from_zone>"+idSeparator+"<to_zone>)", d.Id())
	}
	policyOptions, err := readSecurityPolicy(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityPolicyData(d, policyOptions)

	result[0] = d

	return result, nil
}

func checkSecurityPolicyExists(fromZone, toZone string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	policyConfig, err := sess.command("show configuration"+
		" security policies from-zone "+fromZone+" to-zone "+toZone+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if policyConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSecurityPolicy(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security policies" +
		" from-zone " + d.Get("from_zone").(string) +
		" to-zone " + d.Get("to_zone").(string) +
		" policy "
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
		if policy["permit_tunnel_ipsec_vpn"].(string) != "" {
			if policy["then"].(string) != permitWord {
				return fmt.Errorf("conflict policy then %v and policy permit_tunnel_ipsec_vpn",
					policy["then"].(string))
			}
			configSet = append(configSet, setPrefixPolicy+" then permit tunnel ipsec-vpn "+
				policy["permit_tunnel_ipsec_vpn"].(string))
		}
		if len(policy["permit_application_services"].([]interface{})) > 0 {
			if policy["permit_application_services"].([]interface{})[0] == nil {
				return fmt.Errorf("permit_application_services block is empty")
			}
			if policy["then"].(string) != permitWord {
				return fmt.Errorf("conflict policy then %v and policy permit_application_services",
					policy["then"].(string))
			}
			configSetAppSvc, err := setPolicyPermitApplicationServices(setPrefixPolicy,
				policy["permit_application_services"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetAppSvc...)
		}
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSecurityPolicy(idPolicy string, m interface{}, jnprSess *NetconfObject) (policyOptions, error) {
	zone := strings.Split(idPolicy, idSeparator)
	fromZone := zone[0]
	toZone := zone[1]

	sess := m.(*Session)
	var confRead policyOptions

	policyConfig, err := sess.command("show configuration"+
		" security policies from-zone "+fromZone+" to-zone "+toZone+" | display set relative ", jnprSess)
	if err != nil {
		return confRead, err
	}
	policyList := make([]map[string]interface{}, 0)
	if policyConfig != emptyWord {
		confRead.fromZone = fromZone
		confRead.toZone = toZone
		for _, item := range strings.Split(policyConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.Contains(itemTrim, " match ") || strings.Contains(itemTrim, " then ") {
				policyLineCut := strings.Split(itemTrim, " ")
				m := genMapPolicyWithName(policyLineCut[1])
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
					case strings.HasPrefix(itemTrimPolicy, "then permit tunnel ipsec-vpn "):
						m["then"] = permitWord
						m["permit_tunnel_ipsec_vpn"] = strings.TrimPrefix(itemTrimPolicy,
							"then permit tunnel ipsec-vpn ")
					case strings.HasPrefix(itemTrimPolicy, "then permit application-services"):
						m["then"] = permitWord
						m["permit_application_services"] = readPolicyPermitApplicationServices(itemTrimPolicy,
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
func delSecurityPolicy(fromZone string, toZone string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security policies from-zone "+fromZone+" to-zone "+toZone)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillSecurityPolicyData(d *schema.ResourceData, policyOptions policyOptions) {
	if tfErr := d.Set("from_zone", policyOptions.fromZone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("to_zone", policyOptions.toZone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("policy", policyOptions.policy); tfErr != nil {
		panic(tfErr)
	}
}

func genMapPolicyWithName(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":                        name,
		"match_source_address":        make([]string, 0),
		"match_destination_address":   make([]string, 0),
		"match_application":           make([]string, 0),
		"then":                        "",
		"count":                       false,
		"log_init":                    false,
		"log_close":                   false,
		"permit_application_services": make([]map[string]interface{}, 0),
		"permit_tunnel_ipsec_vpn":     "",
	}
}

func readPolicyPermitApplicationServices(itemTrimPolicy string,
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

func setPolicyPermitApplicationServices(setPrefixPolicy string,
	policyPermitApplicationServices map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	setPrefixPolicyPermitAppSvc := setPrefixPolicy + " then permit application-services"
	if policyPermitApplicationServices["application_firewall_rule_set"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" application-firewall rule-set \""+
			policyPermitApplicationServices["application_firewall_rule_set"].(string)+"\"")
	}
	if policyPermitApplicationServices["application_traffic_control_rule_set"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" application-traffic-control rule-set \""+
			policyPermitApplicationServices["application_traffic_control_rule_set"].(string)+"\"")
	}
	if policyPermitApplicationServices["gprs_gtp_profile"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" gprs-gtp-profile \""+policyPermitApplicationServices["gprs_gtp_profile"].(string)+"\"")
	}
	if policyPermitApplicationServices["gprs_sctp_profile"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" gprs-sctp-profile \""+policyPermitApplicationServices["gprs_sctp_profile"].(string)+"\"")
	}
	if policyPermitApplicationServices["idp"].(bool) {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" idp")
	}
	if policyPermitApplicationServices["redirect_wx"].(bool) &&
		policyPermitApplicationServices["reverse_redirect_wx"].(bool) {
		return configSet, fmt.Errorf("conflict redirect_wx and reverse_redirect_wx enabled both")
	}
	if policyPermitApplicationServices["redirect_wx"].(bool) {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" redirect-wx")
	}
	if policyPermitApplicationServices["reverse_redirect_wx"].(bool) {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" reverse-redirect-wx")
	}
	if policyPermitApplicationServices["security_intelligence_policy"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" security-intelligence-policy \""+
			policyPermitApplicationServices["security_intelligence_policy"].(string)+"\"")
	}
	if len(policyPermitApplicationServices["ssl_proxy"].([]interface{})) > 0 {
		if policyPermitApplicationServices["ssl_proxy"].([]interface{})[0] != nil {
			policyPermitApplicationServicesSSLProxy :=
				policyPermitApplicationServices["ssl_proxy"].([]interface{})[0].(map[string]interface{})
			if policyPermitApplicationServicesSSLProxy["profile_name"].(string) != "" {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+
					" ssl-proxy profile-name \""+
					policyPermitApplicationServicesSSLProxy["profile_name"].(string)+"\"")
			} else {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+" ssl-proxy")
			}
		} else {
			configSet = append(configSet, setPrefixPolicyPermitAppSvc+" ssl-proxy")
		}
	}
	if len(policyPermitApplicationServices["uac_policy"].([]interface{})) > 0 {
		if policyPermitApplicationServices["uac_policy"].([]interface{})[0] != nil {
			policyPermitApplicationServicesUACPolicy :=
				policyPermitApplicationServices["uac_policy"].([]interface{})[0].(map[string]interface{})
			if policyPermitApplicationServicesUACPolicy["captive_portal"].(string) != "" {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+
					" uac-policy captive-portal \""+
					policyPermitApplicationServicesUACPolicy["captive_portal"].(string)+"\"")
			} else {
				configSet = append(configSet, setPrefixPolicyPermitAppSvc+" uac-policy")
			}
		} else {
			configSet = append(configSet, setPrefixPolicyPermitAppSvc+" uac-policy")
		}
	}
	if policyPermitApplicationServices["utm_policy"].(string) != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+
			" utm-policy \""+policyPermitApplicationServices["utm_policy"].(string)+"\"")
	}

	return configSet, nil
}
