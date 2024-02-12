package providersdk

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type idpPolicyOptions struct {
	name       string
	exemptRule []map[string]interface{}
	ipsRule    []map[string]interface{}
}

func resourceSecurityIdpPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityIdpPolicyCreate,
		ReadWithoutTimeout:   resourceSecurityIdpPolicyRead,
		UpdateWithoutTimeout: resourceSecurityIdpPolicyUpdate,
		DeleteWithoutTimeout: resourceSecurityIdpPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityIdpPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"exempt_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringDoesNotContainAny(" "),
						},
						"match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: schemaSecurityIdpPolicyRuleMatch(true),
							},
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ips_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringDoesNotContainAny(" "),
						},
						"match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: schemaSecurityIdpPolicyRuleMatch(false),
							},
						},
						"then": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"class-of-service",
											"close-client",
											"close-client-and-server",
											"close-server",
											"drop-connection",
											"drop-packet",
											"ignore-connection",
											"mark-diffserv",
											"no-action",
											"recommended",
										}, false),
									},
									"action_cos_forwarding_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"action_dscp_code_point": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 63),
									},
									"ip_action": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ip-block", "ip-close", "ip-notify"}, false),
									},
									"ip_action_log": {
										Type:     schema.TypeBool,
										Optional: true,
										// RequiredWith ip_action
									},
									"ip_action_log_create": {
										Type:     schema.TypeBool,
										Optional: true,
										// RequiredWith ip_action
									},
									"ip_action_refresh_timeout": {
										Type:     schema.TypeBool,
										Optional: true,
										// RequiredWith ip_action
									},
									"ip_action_target": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"destination-address",
											"service",
											"source-address",
											"source-zone",
											"source-zone-address",
											"zone-service",
										}, false),
										// RequiredWith ip_action
									},
									"ip_action_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 64800),
										// RequiredWith ip_action
									},
									"notification": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"notification_log_attacks": {
										Type:     schema.TypeBool,
										Optional: true,
										// RequiredWith notification
									},
									"notification_log_attacks_alert": {
										Type:     schema.TypeBool,
										Optional: true,
										// RequiredWith notification, notification_log_attacks
									},
									"notification_packet_log": {
										Type:     schema.TypeBool,
										Optional: true,
										// RequiredWith notification
									},
									"notification_packet_log_post_attack": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
										// RequiredWith notification, notification_packet_log
									},
									"notification_packet_log_post_attack_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 1800),
										// RequiredWith notification, notification_packet_log
									},
									"notification_packet_log_pre_attack": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 255),
										// RequiredWith notification, notification_packet_log
									},
									"severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"critical", "info", "major", "minor", "warning"}, false),
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"terminal": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func schemaSecurityIdpPolicyRuleMatch(exempt bool) map[string]*schema.Schema {
	r := map[string]*schema.Schema{
		"application": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"custom_attack_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"custom_attack": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"destination_address": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"destination_address_except": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dynamic_attack_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"from_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"predefined_attack_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"predefined_attack": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"source_address": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"source_address_except": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"to_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	if exempt {
		delete(r, "application")
	}

	return r
}

func resourceSecurityIdpPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityIdpPolicy(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security idp policy not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	idpPolicyExists, err := checkSecurityIdpPolicyExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if idpPolicyExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security idp idp-policy %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityIdpPolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_idp_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	idpPolicyExists, err = checkSecurityIdpPolicyExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if idpPolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security idp idp-policy %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityIdpPolicyReadWJunSess(d, junSess)...)
}

func resourceSecurityIdpPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityIdpPolicyReadWJunSess(d, junSess)
}

func resourceSecurityIdpPolicyReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	idpPolicyOptions, err := readSecurityIdpPolicy(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if idpPolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityIdpPolicyData(d, idpPolicyOptions)
	}

	return nil
}

func resourceSecurityIdpPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityIdpPolicy(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityIdpPolicy(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityIdpPolicy(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityIdpPolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_idp_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityIdpPolicyReadWJunSess(d, junSess)...)
}

func resourceSecurityIdpPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityIdpPolicy(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityIdpPolicy(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_idp_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityIdpPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idpPolicyExists, err := checkSecurityIdpPolicyExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !idpPolicyExists {
		return nil, fmt.Errorf("don't find security idp idp-policy with id '%v' (id must be <name>)", d.Id())
	}
	idpPolicyOptions, err := readSecurityIdpPolicy(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityIdpPolicyData(d, idpPolicyOptions)

	result[0] = d

	return result, nil
}

func checkSecurityIdpPolicyExists(policy string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp idp-policy \"" + policy + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityIdpPolicy(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security idp idp-policy \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix)
	exemptRuleNameList := make([]string, 0)
	for _, e := range d.Get("exempt_rule").([]interface{}) {
		eM := e.(map[string]interface{})
		if slices.Contains(exemptRuleNameList, eM["name"].(string)) {
			return fmt.Errorf("multiple blocks exempt_rule with the same name %s", eM["name"].(string))
		}
		exemptRuleNameList = append(exemptRuleNameList, eM["name"].(string))
		sets, err := setSecurityIdpPolicyExemptRule(setPrefix, eM)
		if err != nil {
			return err
		}
		configSet = append(configSet, sets...)
	}
	ipsRuleNameList := make([]string, 0)
	for _, e := range d.Get("ips_rule").([]interface{}) {
		eM := e.(map[string]interface{})
		if slices.Contains(ipsRuleNameList, eM["name"].(string)) {
			return fmt.Errorf("multiple blocks ips_rule with the same name %s", eM["name"].(string))
		}
		ipsRuleNameList = append(ipsRuleNameList, eM["name"].(string))
		sets, err := setSecurityIdpPolicyIpsRule(setPrefix, eM)
		if err != nil {
			return err
		}
		configSet = append(configSet, sets...)
	}

	return junSess.ConfigSet(configSet)
}

func setSecurityIdpPolicyExemptRule(setPrefix string, rule map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	setPrefixExeRule := setPrefix + "rulebase-exempt rule \"" + rule["name"].(string) + "\" "
	for _, em := range rule["match"].([]interface{}) {
		if em == nil {
			return configSet, fmt.Errorf("match block in exempt rule '%s' is empty", rule["name"].(string))
		}
		match := em.(map[string]interface{})
		for _, v := range sortSetOfString(match["custom_attack_group"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match attacks custom-attack-groups \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["custom_attack"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match attacks custom-attacks \""+v+"\"")
		}
		if len(match["destination_address"].(*schema.Set).List()) != 0 &&
			len(match["destination_address_except"].(*schema.Set).List()) != 0 {
			return configSet, fmt.Errorf("destination_address and destination_address_except can't set in same time "+
				"in exempt rule '%s'", rule["name"].(string))
		}
		for _, v := range sortSetOfString(match["destination_address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match destination-address \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["destination_address_except"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match destination-except \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["dynamic_attack_group"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match attacks dynamic-attack-groups \""+v+"\"")
		}
		if v := match["from_zone"].(string); v != "" {
			configSet = append(configSet, setPrefixExeRule+"match from-zone \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["predefined_attack_group"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match attacks predefined-attack-groups \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["predefined_attack"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match attacks predefined-attacks \""+v+"\"")
		}
		if len(match["source_address"].(*schema.Set).List()) != 0 &&
			len(match["source_address_except"].(*schema.Set).List()) != 0 {
			return configSet, fmt.Errorf("source_address and source_address_except can't set in same time "+
				"in exempt rule '%s'", rule["name"].(string))
		}
		for _, v := range sortSetOfString(match["source_address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match source-address \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["source_address_except"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixExeRule+"match source-except \""+v+"\"")
		}
		if v := match["to_zone"].(string); v != "" {
			configSet = append(configSet, setPrefixExeRule+"match to-zone \""+v+"\"")
		}
	}
	if v := rule["description"].(string); v != "" {
		configSet = append(configSet, setPrefixExeRule+"description \""+v+"\"")
	}

	return configSet, nil
}

func setSecurityIdpPolicyIpsRule(setPrefix string, rule map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)
	setPrefixIpsRule := setPrefix + "rulebase-ips rule \"" + rule["name"].(string) + "\" "
	for _, em := range rule["match"].([]interface{}) {
		if em == nil {
			return configSet, fmt.Errorf("match block in ips rule '%s' is empty", rule["name"].(string))
		}
		match := em.(map[string]interface{})
		if v := match["application"].(string); v != "" {
			configSet = append(configSet, setPrefixIpsRule+"match application \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["custom_attack_group"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match attacks custom-attack-groups \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["custom_attack"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match attacks custom-attacks \""+v+"\"")
		}
		if len(match["destination_address"].(*schema.Set).List()) != 0 &&
			len(match["destination_address_except"].(*schema.Set).List()) != 0 {
			return configSet, fmt.Errorf("destination_address and destination_address_except can't set in same time "+
				"in ips rule '%s'", rule["name"].(string))
		}
		for _, v := range sortSetOfString(match["destination_address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match destination-address \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["destination_address_except"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match destination-except \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["dynamic_attack_group"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match attacks dynamic-attack-groups \""+v+"\"")
		}
		if v := match["from_zone"].(string); v != "" {
			configSet = append(configSet, setPrefixIpsRule+"match from-zone \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["predefined_attack_group"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match attacks predefined-attack-groups \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["predefined_attack"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match attacks predefined-attacks \""+v+"\"")
		}
		if len(match["source_address"].(*schema.Set).List()) != 0 &&
			len(match["source_address_except"].(*schema.Set).List()) != 0 {
			return configSet, fmt.Errorf("source_address and source_address_except can't set in same time "+
				"in ips rule '%s'", rule["name"].(string))
		}
		for _, v := range sortSetOfString(match["source_address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match source-address \""+v+"\"")
		}
		for _, v := range sortSetOfString(match["source_address_except"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixIpsRule+"match source-except \""+v+"\"")
		}
		if v := match["to_zone"].(string); v != "" {
			configSet = append(configSet, setPrefixIpsRule+"match to-zone \""+v+"\"")
		}
	}
	for _, et := range rule["then"].([]interface{}) {
		then := et.(map[string]interface{})
		configSet = append(configSet, setPrefixIpsRule+"then action "+then["action"].(string))
		if v := then["action_cos_forwarding_class"].(string); v != "" {
			if then["action"].(string) != "class-of-service" {
				return configSet, fmt.Errorf("action_cos_forwarding_class can't set "+
					"if action is not 'class-of-service' in ips rule '%s'", rule["name"].(string))
			}
			configSet = append(configSet, setPrefixIpsRule+"then action "+then["action"].(string)+" forwarding-class \""+v+"\"")
		}
		if v := then["action_dscp_code_point"].(int); v != -1 {
			if then["action"].(string) != "class-of-service" && then["action"].(string) != "mark-diffserv" {
				return configSet, fmt.Errorf("action_dscp_code_point can't set "+
					"if action is not 'class-of-service' or 'mark-diffserv' in ips rule '%s'", rule["name"].(string))
			}
			switch {
			case then["action"].(string) == "class-of-service":
				configSet = append(configSet,
					setPrefixIpsRule+"then action "+then["action"].(string)+" dscp-code-point "+strconv.Itoa(v))
			case then["action"].(string) == "mark-diffserv":
				configSet = append(configSet, setPrefixIpsRule+"then action "+then["action"].(string)+" "+strconv.Itoa(v))
			}
		} else if then["action"].(string) == "mark-diffserv" {
			return configSet, fmt.Errorf("action_dscp_code_point need to be set "+
				"if action == 'mark-diffserv' in rule '%s'", rule["name"].(string))
		}
		if v := then["ip_action"].(string); v != "" {
			configSet = append(configSet, setPrefixIpsRule+"then ip-action "+v)
			if then["ip_action_log"].(bool) {
				configSet = append(configSet, setPrefixIpsRule+"then ip-action log")
			}
			if then["ip_action_log_create"].(bool) {
				configSet = append(configSet, setPrefixIpsRule+"then ip-action log-create")
			}
			if then["ip_action_refresh_timeout"].(bool) {
				configSet = append(configSet, setPrefixIpsRule+"then ip-action refresh-timeout")
			}
			if v2 := then["ip_action_target"].(string); v2 != "" {
				configSet = append(configSet, setPrefixIpsRule+"then ip-action target "+v2)
			}
			if v2 := then["ip_action_timeout"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixIpsRule+"then ip-action timeout "+strconv.Itoa(v2))
			}
		} else if then["ip_action_log"].(bool) ||
			then["ip_action_log_create"].(bool) ||
			then["ip_action_refresh_timeout"].(bool) ||
			then["ip_action_target"].(string) != "" ||
			then["ip_action_timeout"].(int) != -1 {
			return configSet, fmt.Errorf("ip_action need to be set "+
				"with ip_action_* arguments in rule '%s'", rule["name"].(string))
		}
		if then["notification"].(bool) {
			configSet = append(configSet, setPrefixIpsRule+"then notification")
			if then["notification_log_attacks"].(bool) {
				configSet = append(configSet, setPrefixIpsRule+"then notification log-attacks")
				if then["notification_log_attacks_alert"].(bool) {
					configSet = append(configSet, setPrefixIpsRule+"then notification log-attacks alert")
				}
			} else if then["notification_log_attacks_alert"].(bool) {
				return configSet, fmt.Errorf("notification_log_attacks need to be true "+
					"with notification_log_attacks_alert argument in ips rule '%s'", rule["name"].(string))
			}
			if then["notification_packet_log"].(bool) {
				configSet = append(configSet, setPrefixIpsRule+"then notification packet-log")
				if v := then["notification_packet_log_post_attack"].(int); v != -1 {
					configSet = append(configSet, setPrefixIpsRule+"then notification packet-log post-attack "+strconv.Itoa(v))
				}
				if v := then["notification_packet_log_post_attack_timeout"].(int); v != -1 {
					configSet = append(configSet, setPrefixIpsRule+"then notification packet-log post-attack-timeout "+strconv.Itoa(v))
				}
				if v := then["notification_packet_log_pre_attack"].(int); v != 0 {
					configSet = append(configSet, setPrefixIpsRule+"then notification packet-log pre-attack "+strconv.Itoa(v))
				}
			} else if then["notification_packet_log_post_attack"].(int) != -1 ||
				then["notification_packet_log_post_attack_timeout"].(int) != -1 ||
				then["notification_packet_log_pre_attack"].(int) != 0 {
				return configSet, fmt.Errorf("notification_packet_log need to be true "+
					"with notification_packet_log_* arguments in ips rule '%s'", rule["name"].(string))
			}
		} else if then["notification_log_attacks"].(bool) ||
			then["notification_log_attacks_alert"].(bool) ||
			then["notification_packet_log"].(bool) ||
			then["notification_packet_log_post_attack"].(int) != -1 ||
			then["notification_packet_log_post_attack_timeout"].(int) != -1 ||
			then["notification_packet_log_pre_attack"].(int) != 0 {
			return configSet, fmt.Errorf("notification need to be true "+
				"with notification_* arguments in ips rule '%s'", rule["name"].(string))
		}
		if v := then["severity"].(string); v != "" {
			configSet = append(configSet, setPrefixIpsRule+"then severity "+v)
		}
	}
	if v := rule["description"].(string); v != "" {
		configSet = append(configSet, setPrefixIpsRule+"description \""+v+"\"")
	}
	if rule["terminal"].(bool) {
		configSet = append(configSet, setPrefixIpsRule+"terminal")
	}

	return configSet, nil
}

func readSecurityIdpPolicy(policy string, junSess *junos.Session,
) (confRead idpPolicyOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp idp-policy \"" + policy + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = policy
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "rulebase-exempt rule "):
				itemTrimFields := strings.Split(itemTrim, " ")
				rule := map[string]interface{}{
					"name":        strings.Trim(itemTrimFields[0], "\""),
					"match":       make([]map[string]interface{}, 0),
					"description": "",
				}
				confRead.exemptRule = copyAndRemoveItemMapList("name", rule, confRead.exemptRule)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "match "):
					if len(rule["match"].([]map[string]interface{})) == 0 {
						rule["match"] = append(rule["match"].([]map[string]interface{}), map[string]interface{}{
							"custom_attack_group":        make([]string, 0),
							"custom_attack":              make([]string, 0),
							"destination_address":        make([]string, 0),
							"destination_address_except": make([]string, 0),
							"dynamic_attack_group":       make([]string, 0),
							"from_zone":                  "",
							"predefined_attack_group":    make([]string, 0),
							"predefined_attack":          make([]string, 0),
							"source_address":             make([]string, 0),
							"source_address_except":      make([]string, 0),
							"to_zone":                    "",
						})
					}
					match := rule["match"].([]map[string]interface{})[0]
					switch {
					case balt.CutPrefixInString(&itemTrim, "attacks custom-attack-groups "):
						match["custom_attack_group"] = append(
							match["custom_attack_group"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "attacks custom-attacks "):
						match["custom_attack"] = append(match["custom_attack"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-address "):
						match["destination_address"] = append(
							match["destination_address"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-except "):
						match["destination_address_except"] = append(
							match["destination_address_except"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "attacks dynamic-attack-groups "):
						match["dynamic_attack_group"] = append(
							match["dynamic_attack_group"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "from-zone "):
						match["from_zone"] = strings.Trim(itemTrim, "\"")
					case balt.CutPrefixInString(&itemTrim, "attacks predefined-attack-groups "):
						match["predefined_attack_group"] = append(
							match["predefined_attack_group"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "attacks predefined-attacks "):
						match["predefined_attack"] = append(
							match["predefined_attack"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "source-address "):
						match["source_address"] = append(
							match["source_address"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "source-except "):
						match["source_address_except"] = append(
							match["source_address_except"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "to-zone "):
						match["to_zone"] = strings.Trim(itemTrim, "\"")
					}
				case balt.CutPrefixInString(&itemTrim, "description "):
					rule["description"] = strings.Trim(itemTrim, "\"")
				}
				confRead.exemptRule = append(confRead.exemptRule, rule)
			case balt.CutPrefixInString(&itemTrim, "rulebase-ips rule "):
				itemTrimFields := strings.Split(itemTrim, " ")
				rule := map[string]interface{}{
					"name":        strings.Trim(itemTrimFields[0], "\""),
					"match":       make([]map[string]interface{}, 0),
					"then":        make([]map[string]interface{}, 0),
					"description": "",
					"terminal":    false,
				}
				confRead.ipsRule = copyAndRemoveItemMapList("name", rule, confRead.ipsRule)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "match "):
					if len(rule["match"].([]map[string]interface{})) == 0 {
						rule["match"] = append(rule["match"].([]map[string]interface{}), map[string]interface{}{
							"application":                "",
							"custom_attack_group":        make([]string, 0),
							"custom_attack":              make([]string, 0),
							"destination_address":        make([]string, 0),
							"destination_address_except": make([]string, 0),
							"dynamic_attack_group":       make([]string, 0),
							"from_zone":                  "",
							"predefined_attack_group":    make([]string, 0),
							"predefined_attack":          make([]string, 0),
							"source_address":             make([]string, 0),
							"source_address_except":      make([]string, 0),
							"to_zone":                    "",
						})
					}
					match := rule["match"].([]map[string]interface{})[0]
					switch {
					case balt.CutPrefixInString(&itemTrim, "application "):
						match["application"] = strings.Trim(itemTrim, "\"")
					case balt.CutPrefixInString(&itemTrim, "attacks custom-attack-groups "):
						match["custom_attack_group"] = append(
							match["custom_attack_group"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "attacks custom-attacks "):
						match["custom_attack"] = append(
							match["custom_attack"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-address "):
						match["destination_address"] = append(
							match["destination_address"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-except "):
						match["destination_address_except"] = append(
							match["destination_address_except"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "attacks dynamic-attack-groups "):
						match["dynamic_attack_group"] = append(
							match["dynamic_attack_group"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "from-zone "):
						match["from_zone"] = strings.Trim(itemTrim, "\"")
					case balt.CutPrefixInString(&itemTrim, "attacks predefined-attack-groups "):
						match["predefined_attack_group"] = append(
							match["predefined_attack_group"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "attacks predefined-attacks "):
						match["predefined_attack"] = append(
							match["predefined_attack"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "source-address "):
						match["source_address"] = append(
							match["source_address"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "source-except "):
						match["source_address_except"] = append(
							match["source_address_except"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "to-zone "):
						match["to_zone"] = strings.Trim(itemTrim, "\"")
					}
				case balt.CutPrefixInString(&itemTrim, "then "):
					if len(rule["then"].([]map[string]interface{})) == 0 {
						rule["then"] = append(rule["then"].([]map[string]interface{}), map[string]interface{}{
							"action":                                      "",
							"action_cos_forwarding_class":                 "",
							"action_dscp_code_point":                      -1,
							"ip_action":                                   "",
							"ip_action_log":                               false,
							"ip_action_log_create":                        false,
							"ip_action_refresh_timeout":                   false,
							"ip_action_target":                            "",
							"ip_action_timeout":                           -1,
							"notification":                                false,
							"notification_log_attacks":                    false,
							"notification_log_attacks_alert":              false,
							"notification_packet_log":                     false,
							"notification_packet_log_post_attack":         -1,
							"notification_packet_log_post_attack_timeout": -1,
							"notification_packet_log_pre_attack":          0,
							"severity":                                    "",
						})
					}
					then := rule["then"].([]map[string]interface{})[0]
					switch {
					case balt.CutPrefixInString(&itemTrim, "action "):
						switch {
						case balt.CutPrefixInString(&itemTrim, "class-of-service "):
							then["action"] = "class-of-service"
							switch {
							case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
								then["action_cos_forwarding_class"] = strings.Trim(itemTrim, "\"")
							case balt.CutPrefixInString(&itemTrim, "dscp-code-point "):
								then["action_dscp_code_point"], err = strconv.Atoi(itemTrim)
								if err != nil {
									return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
								}
							}
						case balt.CutPrefixInString(&itemTrim, "mark-diffserv "):
							then["action"] = "mark-diffserv"
							then["action_dscp_code_point"], err = strconv.Atoi(itemTrim)
							if err != nil {
								return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
							}
						default:
							then["action"] = itemTrim
						}
					case balt.CutPrefixInString(&itemTrim, "ip-action "):
						switch {
						case itemTrim == "log":
							then["ip_action_log"] = true
						case itemTrim == "log-create":
							then["ip_action_log_create"] = true
						case itemTrim == "refresh-timeout":
							then["ip_action_refresh_timeout"] = true
						case balt.CutPrefixInString(&itemTrim, "target "):
							then["ip_action_target"] = itemTrim
						case balt.CutPrefixInString(&itemTrim, "timeout "):
							then["ip_action_timeout"], err = strconv.Atoi(itemTrim)
							if err != nil {
								return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
							}
						default:
							then["ip_action"] = itemTrim
						}
					case balt.CutPrefixInString(&itemTrim, "notification"):
						then["notification"] = true
						switch {
						case balt.CutPrefixInString(&itemTrim, " log-attacks"):
							then["notification_log_attacks"] = true
							if itemTrim == " alert" {
								then["notification_log_attacks_alert"] = true
							}
						case balt.CutPrefixInString(&itemTrim, " packet-log"):
							then["notification_packet_log"] = true
							switch {
							case balt.CutPrefixInString(&itemTrim, " post-attack "):
								then["notification_packet_log_post_attack"], err = strconv.Atoi(itemTrim)
							case balt.CutPrefixInString(&itemTrim, " post-attack-timeout "):
								then["notification_packet_log_post_attack_timeout"], err = strconv.Atoi(itemTrim)
							case balt.CutPrefixInString(&itemTrim, " pre-attack "):
								then["notification_packet_log_pre_attack"], err = strconv.Atoi(itemTrim)
							}
							if err != nil {
								return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
							}
						}
					case balt.CutPrefixInString(&itemTrim, "severity "):
						then["severity"] = itemTrim
					}
				case balt.CutPrefixInString(&itemTrim, "description "):
					rule["description"] = strings.Trim(itemTrim, "\"")
				case itemTrim == "terminal":
					rule["terminal"] = true
				}
				confRead.ipsRule = append(confRead.ipsRule, rule)
			}
		}
	}

	return confRead, nil
}

func delSecurityIdpPolicy(policy string, junSess *junos.Session) error {
	configSet := []string{"delete security idp idp-policy \"" + policy + "\""}

	return junSess.ConfigSet(configSet)
}

func fillSecurityIdpPolicyData(d *schema.ResourceData, idpPolicyOptions idpPolicyOptions) {
	if tfErr := d.Set("name", idpPolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("exempt_rule", idpPolicyOptions.exemptRule); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ips_rule", idpPolicyOptions.ipsRule); tfErr != nil {
		panic(tfErr)
	}
}
