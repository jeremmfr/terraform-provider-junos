package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
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
		CreateWithoutTimeout: resourceServicesSecurityIntellProfileCreate,
		ReadWithoutTimeout:   resourceServicesSecurityIntellProfileRead,
		UpdateWithoutTimeout: resourceServicesSecurityIntellProfileUpdate,
		DeleteWithoutTimeout: resourceServicesSecurityIntellProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServicesSecurityIntellProfileImport,
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

func resourceServicesSecurityIntellProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setServicesSecurityIntellProfile(d, junSess); err != nil {
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
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityIntellProfileExists, err := checkServicesSecurityIntellProfileExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityIntellProfileExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services security-intelligence profile %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesSecurityIntellProfile(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_services_security_intelligence_profile")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityIntellProfileExists, err = checkServicesSecurityIntellProfileExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityIntellProfileExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services security-intelligence profile %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesSecurityIntellProfileReadWJunSess(d, junSess)...)
}

func resourceServicesSecurityIntellProfileRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceServicesSecurityIntellProfileReadWJunSess(d, junSess)
}

func resourceServicesSecurityIntellProfileReadWJunSess(
	d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	securityIntellProfileOptions, err := readServicesSecurityIntellProfile(d.Get("name").(string), junSess)
	junos.MutexUnlock()
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

func resourceServicesSecurityIntellProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delServicesSecurityIntellProfile(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesSecurityIntellProfile(d, junSess); err != nil {
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
	if err := delServicesSecurityIntellProfile(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesSecurityIntellProfile(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_services_security_intelligence_profile")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesSecurityIntellProfileReadWJunSess(d, junSess)...)
}

func resourceServicesSecurityIntellProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delServicesSecurityIntellProfile(d.Get("name").(string), junSess); err != nil {
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
	if err := delServicesSecurityIntellProfile(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_services_security_intelligence_profile")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesSecurityIntellProfileImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	securityIntellProfileExists, err := checkServicesSecurityIntellProfileExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityIntellProfileExists {
		return nil, fmt.Errorf("don't find services security-intelligence profile with id '%v' (id must be <name>)", d.Id())
	}
	securityIntellProfileOptions, err := readServicesSecurityIntellProfile(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillServicesSecurityIntellProfileData(d, securityIntellProfileOptions)

	result[0] = d

	return result, nil
}

func checkServicesSecurityIntellProfileExists(profile string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services security-intelligence profile \"" + profile + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setServicesSecurityIntellProfile(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set services security-intelligence profile \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix+"category "+d.Get("category").(string))
	ruleNameList := make([]string, 0)
	for _, v := range d.Get("rule").([]interface{}) {
		rule := v.(map[string]interface{})
		if slices.Contains(ruleNameList, rule["name"].(string)) {
			return fmt.Errorf("multiple blocks rule with the same name %s", rule["name"].(string))
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

	return junSess.ConfigSet(configSet)
}

func readServicesSecurityIntellProfile(profile string, junSess *junos.Session,
) (confRead securityIntellProfileOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services security-intelligence profile \"" + profile + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = profile
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "category "):
				confRead.category = itemTrim
			case balt.CutPrefixInString(&itemTrim, "rule "):
				itemTrimFields := strings.Split(itemTrim, " ")
				rule := map[string]interface{}{
					"name":        strings.Trim(itemTrimFields[0], "\""),
					"match":       make([]map[string]interface{}, 0),
					"then_action": "",
					"then_log":    false,
				}
				confRead.rule = copyAndRemoveItemMapList("name", rule, confRead.rule)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				if err := readServicesSecurityIntellProfileRule(itemTrim, rule); err != nil {
					return confRead, err
				}
				confRead.rule = append(confRead.rule, rule)
			case balt.CutPrefixInString(&itemTrim, "default-rule then "):
				if len(confRead.defaultRuleThen) == 0 {
					confRead.defaultRuleThen = append(confRead.defaultRuleThen, map[string]interface{}{
						"action": "",
						"log":    false,
						"no_log": false,
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "action "):
					confRead.defaultRuleThen[0]["action"] = itemTrim
				case itemTrim == "log":
					confRead.defaultRuleThen[0]["log"] = true
				case itemTrim == "no-log":
					confRead.defaultRuleThen[0]["no_log"] = true
				}
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			}
		}
	}

	return confRead, nil
}

func readServicesSecurityIntellProfileRule(itemTrim string, ruleMap map[string]interface{}) error {
	switch {
	case balt.CutPrefixInString(&itemTrim, "match "):
		if len(ruleMap["match"].([]map[string]interface{})) == 0 {
			ruleMap["match"] = append(ruleMap["match"].([]map[string]interface{}), map[string]interface{}{
				"threat_level": make([]int, 0),
				"feed_name":    make([]string, 0),
			})
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "threat-level "):
			threatLevel, err := strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
			ruleMap["match"].([]map[string]interface{})[0]["threat_level"] = append(
				ruleMap["match"].([]map[string]interface{})[0]["threat_level"].([]int),
				threatLevel,
			)
		case balt.CutPrefixInString(&itemTrim, "feed-name "):
			ruleMap["match"].([]map[string]interface{})[0]["feed_name"] = append(
				ruleMap["match"].([]map[string]interface{})[0]["feed_name"].([]string),
				itemTrim,
			)
		}
	case balt.CutPrefixInString(&itemTrim, "then action "):
		ruleMap["then_action"] = itemTrim
	case itemTrim == "then log":
		ruleMap["then_log"] = true
	}

	return nil
}

func delServicesSecurityIntellProfile(profile string, junSess *junos.Session) error {
	configSet := []string{"delete services security-intelligence profile \"" + profile + "\""}

	return junSess.ConfigSet(configSet)
}

func fillServicesSecurityIntellProfileData(
	d *schema.ResourceData, securityIntellProfileOptions securityIntellProfileOptions,
) {
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
