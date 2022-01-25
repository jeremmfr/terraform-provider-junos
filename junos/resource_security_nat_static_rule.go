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

type natStaticRuleOptions struct {
	destinationPort        int
	destinationPortTo      int
	name                   string
	destinationAddress     string
	destinationAddressName string
	ruleSet                string
	sourceAddress          []string
	sourceAddressName      []string
	sourcePort             []string
	then                   []map[string]interface{}
}

func resourceSecurityNatStaticRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityNatStaticRuleCreate,
		ReadContext:   resourceSecurityNatStaticRuleRead,
		UpdateContext: resourceSecurityNatStaticRuleUpdate,
		DeleteContext: resourceSecurityNatStaticRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityNatStaticRuleImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"rule_set": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"destination_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"destination_address", "destination_address_name"},
				ValidateFunc: validation.IsCIDRNetwork(0, 128),
			},
			"destination_address_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"destination_address", "destination_address_name"},
			},
			"destination_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"destination_port_to": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"destination_port"},
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"source_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_address_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_port": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"then": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{inetWord, prefixWord, prefixNameWord}, false),
						},
						"mapped_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"mapped_port_to": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"then.0.mapped_port"},
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"routing_instance": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
					},
				},
			},
		},
	}
}

func resourceSecurityNatStaticRuleCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityNatStaticRule(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("rule_set").(string) + idSeparator + d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security nat static rule in rule-set not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	natStaticExists, err := checkSecurityNatStaticExists(d.Get("rule_set").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !natStaticExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security nat static rule-set %v doesn't exist", d.Get("rule_set").(string)))...)
	}
	natStaticRuleExists, err := checkSecurityNatStaticRuleExists(
		d.Get("rule_set").(string), d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if natStaticRuleExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security nat static rule %v already exists in rule-set %s",
			d.Get("name").(string), d.Get("rule_set").(string)))...)
	}

	if err := setSecurityNatStaticRule(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_nat_static_rule", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	natStaticRuleExists, err = checkSecurityNatStaticRuleExists(
		d.Get("rule_set").(string), d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if natStaticRuleExists {
		d.SetId(d.Get("rule_set").(string) + idSeparator + d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security nat statuc rule %v not exists in rule-set %s after commit "+
				"=> check your config", d.Get("name").(string), d.Get("rule_set").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatStaticRuleReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityNatStaticRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityNatStaticRuleReadWJnprSess(d, m, jnprSess)
}

func resourceSecurityNatStaticRuleReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	natStaticRuleOptions, err := readSecurityNatStaticRule(d.Get("rule_set").(string), d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if natStaticRuleOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatStaticRuleData(d, natStaticRuleOptions)
	}

	return nil
}

func resourceSecurityNatStaticRuleUpdate(ctx context.Context,
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
	if err := delSecurityNatStaticRule(d.Get("rule_set").(string), d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityNatStaticRule(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_nat_static_rule", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatStaticRuleReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityNatStaticRuleDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityNatStaticRule(d.Get("rule_set").(string), d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_nat_static_rule", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatStaticRuleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	natStaticRuleExists, err := checkSecurityNatStaticRuleExists(idList[0], idList[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !natStaticRuleExists {
		return nil, fmt.Errorf(
			"don't find static nat rule with id '%v' (id must be <rule_set>"+idSeparator+"<name>)", d.Id())
	}
	natStaticRuleOptions, err := readSecurityNatStaticRule(idList[0], idList[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatStaticRuleData(d, natStaticRuleOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatStaticRuleExists(ruleSet, name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration"+
		" security nat static rule-set "+ruleSet+" rule "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSecurityNatStaticRule(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	regexpSourcePort := regexp.MustCompile(`^\d+( to \d+)?$`)

	setPrefix := "set security nat static rule-set " + d.Get("rule_set").(string) +
		" rule " + d.Get("name").(string) + " "
	configSet = append(configSet, setPrefix)
	if v := d.Get("destination_address").(string); v != "" {
		configSet = append(configSet, setPrefix+"match destination-address "+v)
	}
	if v := d.Get("destination_address_name").(string); v != "" {
		configSet = append(configSet, setPrefix+"match destination-address-name \""+v+"\"")
	}
	if v := d.Get("destination_port").(int); v != 0 {
		configSet = append(configSet, setPrefix+"match destination-port "+strconv.Itoa(v))
		if vv := d.Get("destination_port_to").(int); vv != 0 {
			configSet = append(configSet, setPrefix+"match destination-port to "+strconv.Itoa(vv))
		}
	} else if d.Get("destination_port_to").(int) != 0 {
		return fmt.Errorf("destination_port need to be not 0 with destination_port_to")
	}
	for _, v := range sortSetOfString(d.Get("source_address").(*schema.Set).List()) {
		if err := validateCIDRNetwork(v); err != nil {
			return err
		}
		configSet = append(configSet, setPrefix+"match source-address "+v)
	}
	for _, v := range sortSetOfString(d.Get("source_address_name").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"match source-address-name \""+v+"\"")
	}
	for _, v := range sortSetOfString(d.Get("source_port").(*schema.Set).List()) {
		if !regexpSourcePort.MatchString(v) {
			return fmt.Errorf("source_port need to have format `x` or `x to y`")
		}
		configSet = append(configSet, setPrefix+"match source-port "+v)
	}
	for _, v := range d.Get("then").([]interface{}) {
		then := v.(map[string]interface{})
		if then["type"].(string) == inetWord {
			if then["routing_instance"].(string) == "" {
				return fmt.Errorf("missing routing_instance with type = inet")
			}
			if then["prefix"].(string) != "" ||
				then["mapped_port"].(int) != 0 ||
				then["mapped_port_to"].(int) != 0 {
				return fmt.Errorf("only routing_instance need to be set with type = inet")
			}
			configSet = append(configSet, setPrefix+"then static-nat inet routing-instance "+
				then["routing_instance"].(string))
		}
		if then["type"].(string) == prefixWord || then["type"].(string) == prefixNameWord {
			setPrefixRuleThenStaticNat := setPrefix + "then static-nat "
			if then["type"].(string) == prefixWord {
				setPrefixRuleThenStaticNat += "prefix "
				if then["prefix"].(string) == "" {
					return fmt.Errorf("missing prefix with type = prefix")
				}
				if err := validateCIDRNetwork(then["prefix"].(string)); err != nil {
					return err
				}
			}
			if then["type"].(string) == prefixNameWord {
				setPrefixRuleThenStaticNat += "prefix-name "
				if then["prefix"].(string) == "" {
					return fmt.Errorf("missing prefix with type = prefix-name")
				}
			}
			configSet = append(configSet, setPrefixRuleThenStaticNat+"\""+then["prefix"].(string)+"\"")
			if vv := then["mapped_port"].(int); vv != 0 {
				configSet = append(configSet, setPrefixRuleThenStaticNat+"mapped-port "+strconv.Itoa(vv))
				if vvTo := then["mapped_port_to"].(int); vvTo != 0 {
					configSet = append(configSet, setPrefixRuleThenStaticNat+"mapped-port to "+strconv.Itoa(vvTo))
				}
			} else if then["mapped_port_to"].(int) != 0 {
				return fmt.Errorf("mapped_port need to be not 0 with mapped_port_to")
			}
			if vv := then["routing_instance"].(string); vv != "" {
				configSet = append(configSet, setPrefixRuleThenStaticNat+"routing-instance "+vv)
			}
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readSecurityNatStaticRule(ruleSet, name string,
	m interface{}, jnprSess *NetconfObject) (natStaticRuleOptions, error) {
	sess := m.(*Session)
	var confRead natStaticRuleOptions

	showConfig, err := sess.command("show configuration"+
		" security nat static rule-set "+ruleSet+" rule "+name+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = name
		confRead.ruleSet = ruleSet
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "match destination-address "):
				confRead.destinationAddress = strings.TrimPrefix(itemTrim, "match destination-address ")
			case strings.HasPrefix(itemTrim, "match destination-address-name "):
				confRead.destinationAddressName = strings.Trim(strings.TrimPrefix(
					itemTrim, "match destination-address-name "), "\"")
			case strings.HasPrefix(itemTrim, "match destination-port to "):
				var err error
				confRead.destinationPortTo, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "match destination-port to "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "match destination-port "):
				var err error
				confRead.destinationPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim,
					"match destination-port "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "match source-address "):
				confRead.sourceAddress = append(confRead.sourceAddress, strings.TrimPrefix(itemTrim, "match source-address "))
			case strings.HasPrefix(itemTrim, "match source-address-name "):
				confRead.sourceAddressName = append(confRead.sourceAddressName,
					strings.Trim(strings.TrimPrefix(itemTrim, "match source-address-name "), "\""))
			case strings.HasPrefix(itemTrim, "match source-port "):
				confRead.sourcePort = append(confRead.sourcePort, strings.TrimPrefix(itemTrim, "match source-port "))
			case strings.HasPrefix(itemTrim, "then static-nat "):
				itemThen := strings.TrimPrefix(itemTrim, "then static-nat ")
				if len(confRead.then) == 0 {
					confRead.then = append(confRead.then, map[string]interface{}{
						"type":             "",
						"mapped_port":      0,
						"mapped_port_to":   0,
						"prefix":           "",
						"routing_instance": "",
					})
				}
				ruleThenOptions := confRead.then[0]
				switch {
				case strings.HasPrefix(itemThen, "prefix ") || strings.HasPrefix(itemThen, "prefix-name "):
					if strings.HasPrefix(itemThen, "prefix ") {
						ruleThenOptions["type"] = prefixWord
						itemThen = strings.TrimPrefix(itemThen, "prefix ")
					}
					if strings.HasPrefix(itemThen, "prefix-name ") {
						ruleThenOptions["type"] = prefixNameWord
						itemThen = strings.TrimPrefix(itemThen, "prefix-name ")
					}
					switch {
					case strings.HasPrefix(itemThen, "routing-instance "):
						ruleThenOptions["routing_instance"] = strings.TrimPrefix(itemThen, "routing-instance ")
					case strings.HasPrefix(itemThen, "mapped-port to "):
						var err error
						ruleThenOptions["mapped_port_to"], err = strconv.Atoi(strings.TrimPrefix(itemThen, "mapped-port to "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					case strings.HasPrefix(itemThen, "mapped-port "):
						var err error
						ruleThenOptions["mapped_port"], err = strconv.Atoi(strings.TrimPrefix(itemThen, "mapped-port "))
						if err != nil {
							return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
						}
					default:
						ruleThenOptions["prefix"] = strings.Trim(itemThen, "\"")
					}
				case strings.HasPrefix(itemThen, "inet "):
					ruleThenOptions["type"] = inetWord
					ruleThenOptions["routing_instance"] = strings.TrimPrefix(itemThen, "inet routing-instance ")
				}
			}
		}
	}

	return confRead, nil
}

func delSecurityNatStaticRule(ruleSet, name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete security nat static rule-set " + ruleSet + " rule " + name}

	return sess.configSet(configSet, jnprSess)
}

func fillSecurityNatStaticRuleData(d *schema.ResourceData, natStaticRuleOptions natStaticRuleOptions) {
	if tfErr := d.Set("name", natStaticRuleOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("rule_set", natStaticRuleOptions.ruleSet); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("destination_address", natStaticRuleOptions.destinationAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("destination_address_name", natStaticRuleOptions.destinationAddressName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("destination_port", natStaticRuleOptions.destinationPort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("destination_port_to", natStaticRuleOptions.destinationPortTo); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("source_address", natStaticRuleOptions.sourceAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("source_address_name", natStaticRuleOptions.sourceAddressName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("source_port", natStaticRuleOptions.sourcePort); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("then", natStaticRuleOptions.then); tfErr != nil {
		panic(tfErr)
	}
}
