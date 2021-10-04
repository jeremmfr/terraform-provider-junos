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
)

type securityIntellProfileOptions struct {
	name            string
	category        string
	description     string
	defaultRuleThen []map[string]interface{}
	rule            []map[string]interface{}
}

func resourceServicesSecurityIntellProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicesSecurityIntellProfileCreate,
		ReadContext:   resourceServicesSecurityIntellProfileRead,
		UpdateContext: resourceServicesSecurityIntellProfileUpdate,
		DeleteContext: resourceServicesSecurityIntellProfileDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServicesSecurityIntellProfileImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
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
								Schema: map[string]*schema.Schema{
									"threat_level": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"feed_name": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"then_action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^(permit|recommended|block (drop|close( http (file|message|redirect-url) .+)?))$`),
								"must have valid action (permit|recommended|block...)"),
						},
						"then_log": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"default_rule_then": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^(permit|recommended|block (drop|close( http (file|message|redirect-url) .+)?))$`),
								"must have valid action (permit|recommended|block...)"),
						},
						"log": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"default_rule_then.0.no_log"},
						},
						"no_log": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"default_rule_then.0.log"},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceServicesSecurityIntellProfileCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setServicesSecurityIntellProfile(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	securityIntellProfileExists, err := checkServicesSecurityIntellProfileExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityIntellProfileExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services security-intelligence profile %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesSecurityIntellProfile(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_services_security_intelligence_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityIntellProfileExists, err = checkServicesSecurityIntellProfileExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityIntellProfileExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services security-intelligence profile %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesSecurityIntellProfileReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesSecurityIntellProfileRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceServicesSecurityIntellProfileReadWJnprSess(d, m, jnprSess)
}

func resourceServicesSecurityIntellProfileReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	securityIntellProfileOptions, err := readServicesSecurityIntellProfile(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if securityIntellProfileOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesSecurityIntellProfileData(d, securityIntellProfileOptions)
	}

	return nil
}

func resourceServicesSecurityIntellProfileUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delServicesSecurityIntellProfile(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesSecurityIntellProfile(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_services_security_intelligence_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesSecurityIntellProfileReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesSecurityIntellProfileDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delServicesSecurityIntellProfile(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_services_security_intelligence_profile", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesSecurityIntellProfileImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityIntellProfileExists, err := checkServicesSecurityIntellProfileExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityIntellProfileExists {
		return nil, fmt.Errorf("don't find services security-intelligence profile with id '%v' (id must be <name>)", d.Id())
	}
	securityIntellProfileOptions, err := readServicesSecurityIntellProfile(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillServicesSecurityIntellProfileData(d, securityIntellProfileOptions)

	result[0] = d

	return result, nil
}

func checkServicesSecurityIntellProfileExists(profile string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration"+
		" services security-intelligence profile \""+profile+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setServicesSecurityIntellProfile(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set services security-intelligence profile \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix+"category "+d.Get("category").(string))
	ruleNameList := make([]string, 0)
	for _, v := range d.Get("rule").([]interface{}) {
		rule := v.(map[string]interface{})
		if stringInSlice(rule["name"].(string), ruleNameList) {
			return fmt.Errorf("multiple rule blocks with the same name")
		}
		ruleNameList = append(ruleNameList, rule["name"].(string))
		setPrefixRule := setPrefix + "rule \"" + rule["name"].(string) + "\" "
		for _, v2 := range rule["match"].([]interface{}) {
			match := v2.(map[string]interface{})
			for _, v3 := range match["threat_level"].([]interface{}) {
				configSet = append(configSet, setPrefixRule+"match threat-level "+strconv.Itoa(v3.(int)))
			}
			for _, v3 := range match["feed_name"].([]interface{}) {
				configSet = append(configSet, setPrefixRule+"match feed-name "+v3.(string))
			}
		}
		configSet = append(configSet, setPrefixRule+"then action "+rule["then_action"].(string))
		if rule["then_log"].(bool) {
			configSet = append(configSet, setPrefixRule+"then log")
		}
	}
	for _, v := range d.Get("default_rule_then").([]interface{}) {
		rule := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"default-rule then action "+rule["action"].(string))
		if rule["log"].(bool) {
			configSet = append(configSet, setPrefix+"default-rule then log")
		}
		if rule["no_log"].(bool) {
			configSet = append(configSet, setPrefix+"default-rule then no-log")
		}
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	return sess.configSet(configSet, jnprSess)
}

func readServicesSecurityIntellProfile(profile string, m interface{}, jnprSess *NetconfObject) (
	securityIntellProfileOptions, error) {
	sess := m.(*Session)
	var confRead securityIntellProfileOptions

	showConfig, err := sess.command("show configuration"+
		" services security-intelligence profile \""+profile+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = profile
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "category "):
				confRead.category = strings.TrimPrefix(itemTrim, "category ")
			case strings.HasPrefix(itemTrim, "rule "):
				ruleLineCut := strings.Split(itemTrim, " ")
				rule := map[string]interface{}{
					"name":        strings.Trim(ruleLineCut[1], "\""),
					"match":       make([]map[string]interface{}, 0),
					"then_action": "",
					"then_log":    false,
				}
				confRead.rule = copyAndRemoveItemMapList("name", rule, confRead.rule)
				itemTrimRule := strings.TrimPrefix(itemTrim, "rule "+ruleLineCut[1]+" ")
				if err := readServicesSecurityIntellProfileRule(itemTrimRule, rule); err != nil {
					return confRead, err
				}
				confRead.rule = append(confRead.rule, rule)
			case strings.HasPrefix(itemTrim, "default-rule then "):
				if len(confRead.defaultRuleThen) == 0 {
					confRead.defaultRuleThen = append(confRead.defaultRuleThen, map[string]interface{}{
						"action": "",
						"log":    false,
						"no_log": false,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "default-rule then action "):
					confRead.defaultRuleThen[0]["action"] = strings.TrimPrefix(itemTrim, "default-rule then action ")
				case itemTrim == "default-rule then log":
					confRead.defaultRuleThen[0]["log"] = true
				case itemTrim == "default-rule then no-log":
					confRead.defaultRuleThen[0]["no_log"] = true
				}
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			}
		}
	}

	return confRead, nil
}

func readServicesSecurityIntellProfileRule(itemTrimPolicyRule string, ruleMap map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrimPolicyRule, "match "):
		if len(ruleMap["match"].([]map[string]interface{})) == 0 {
			ruleMap["match"] = append(ruleMap["match"].([]map[string]interface{}), map[string]interface{}{
				"threat_level": make([]int, 0),
				"feed_name":    make([]string, 0),
			})
		}
		switch {
		case strings.HasPrefix(itemTrimPolicyRule, "match threat-level "):
			threatLevel, err :=
				strconv.Atoi(strings.TrimPrefix(itemTrimPolicyRule, "match threat-level "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrimPolicyRule, err)
			}
			ruleMap["match"].([]map[string]interface{})[0]["threat_level"] = append(
				ruleMap["match"].([]map[string]interface{})[0]["threat_level"].([]int), threatLevel)
		case strings.HasPrefix(itemTrimPolicyRule, "match feed-name "):
			ruleMap["match"].([]map[string]interface{})[0]["feed_name"] = append(
				ruleMap["match"].([]map[string]interface{})[0]["feed_name"].([]string),
				strings.TrimPrefix(itemTrimPolicyRule, "match feed-name "))
		}
	case strings.HasPrefix(itemTrimPolicyRule, "then action "):
		ruleMap["then_action"] = strings.TrimPrefix(itemTrimPolicyRule, "then action ")
	case itemTrimPolicyRule == "then log":
		ruleMap["then_log"] = true
	}

	return nil
}

func delServicesSecurityIntellProfile(profile string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete services security-intelligence profile \"" + profile + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillServicesSecurityIntellProfileData(
	d *schema.ResourceData, securityIntellProfileOptions securityIntellProfileOptions) {
	if tfErr := d.Set("name", securityIntellProfileOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("category", securityIntellProfileOptions.category); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("rule", securityIntellProfileOptions.rule); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("default_rule_then", securityIntellProfileOptions.defaultRuleThen); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", securityIntellProfileOptions.description); tfErr != nil {
		panic(tfErr)
	}
}
