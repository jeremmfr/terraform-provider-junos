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
)

type natStaticOptions struct {
	name        string
	description string
	from        []map[string]interface{}
	rule        []map[string]interface{}
}

func resourceSecurityNatStatic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityNatStaticCreate,
		ReadContext:   resourceSecurityNatStaticRead,
		UpdateContext: resourceSecurityNatStaticUpdate,
		DeleteContext: resourceSecurityNatStaticDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityNatStaticImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"from": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"interface", "routing-instance", "zone"}, false),
						},
						"value": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
						},
						"destination_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
						"destination_address_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"destination_port_to": {
							Type:         schema.TypeInt,
							Optional:     true,
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
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSecurityNatStaticCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityNatStatic(d, m, nil); err != nil {
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
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security nat static not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	securityNatStaticExists, err := checkSecurityNatStaticExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatStaticExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat static %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityNatStatic(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_nat_static", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityNatStaticExists, err = checkSecurityNatStaticExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatStaticExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat static %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatStaticReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityNatStaticRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityNatStaticReadWJnprSess(d, m, jnprSess)
}

func resourceSecurityNatStaticReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	natStaticOptions, err := readSecurityNatStatic(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if natStaticOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatStaticData(d, natStaticOptions)
	}

	return nil
}

func resourceSecurityNatStaticUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityNatStatic(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityNatStatic(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_nat_static", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatStaticReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityNatStaticDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityNatStatic(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_nat_static", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatStaticImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	securityNatStaticExists, err := checkSecurityNatStaticExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityNatStaticExists {
		return nil, fmt.Errorf("don't find nat static with id '%v' (id must be <name>)", d.Id())
	}
	natStaticOptions, err := readSecurityNatStatic(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatStaticData(d, natStaticOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatStaticExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	natStaticConfig, err := sess.command("show configuration"+
		" security nat static rule-set "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if natStaticConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSecurityNatStatic(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	regexpSourcePort := regexp.MustCompile(`^\d+( to \d+)?$`)

	setPrefix := "set security nat static rule-set " + d.Get("name").(string)
	for _, v := range d.Get("from").([]interface{}) {
		from := v.(map[string]interface{})
		for _, value := range sortSetOfString(from["value"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+" from "+from["type"].(string)+" "+value)
		}
	}
	ruleNameList := make([]string, 0)
	for _, v := range d.Get("rule").([]interface{}) {
		rule := v.(map[string]interface{})
		if bchk.StringInSlice(rule["name"].(string), ruleNameList) {
			return fmt.Errorf("multiple rule blocks with the same name")
		}
		ruleNameList = append(ruleNameList, rule["name"].(string))
		setPrefixRule := setPrefix + " rule " + rule["name"].(string)
		if rule["destination_address"].(string) == "" && rule["destination_address_name"].(string) == "" {
			return fmt.Errorf("missing destination_address or destination_address_name in rule %s", rule["name"].(string))
		}
		if rule["destination_address"].(string) != "" && rule["destination_address_name"].(string) != "" {
			return fmt.Errorf("destination_address and destination_address_name must not be set at the same time "+
				"in rule %s", rule["name"].(string))
		}
		if vv := rule["destination_address"].(string); vv != "" {
			configSet = append(configSet, setPrefixRule+" match destination-address "+vv)
		}
		if vv := rule["destination_address_name"].(string); vv != "" {
			configSet = append(configSet, setPrefixRule+" match destination-address-name \""+vv+"\"")
		}
		if vv := rule["destination_port"].(int); vv != 0 {
			configSet = append(configSet, setPrefixRule+" match destination-port "+strconv.Itoa(vv))
			if vvTo := rule["destination_port_to"].(int); vvTo != 0 {
				configSet = append(configSet, setPrefixRule+" match destination-port to "+strconv.Itoa(vvTo))
			}
		} else if rule["destination_port_to"].(int) != 0 {
			return fmt.Errorf("destination_port need to be set with destination_port_to in rule %s", rule["name"].(string))
		}
		for _, vv := range sortSetOfString(rule["source_address"].(*schema.Set).List()) {
			if err := validateCIDRNetwork(vv); err != nil {
				return err
			}
			configSet = append(configSet, setPrefixRule+" match source-address "+vv)
		}
		for _, vv := range sortSetOfString(rule["source_address_name"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixRule+" match source-address-name \""+vv+"\"")
		}
		for _, vv := range sortSetOfString(rule["source_port"].(*schema.Set).List()) {
			if !regexpSourcePort.MatchString(vv) {
				return fmt.Errorf("source_port need to have format `x` or `x to y` in rule %s", rule["name"].(string))
			}
			configSet = append(configSet, setPrefixRule+" match source-port "+vv)
		}
		for _, thenV := range rule[thenWord].([]interface{}) {
			then := thenV.(map[string]interface{})
			if then["type"].(string) == inetWord {
				if then["routing_instance"].(string) == "" {
					return fmt.Errorf("missing routing_instance in rule %s with type = inet", rule["name"].(string))
				}
				if then["prefix"].(string) != "" ||
					then["mapped_port"].(int) != 0 ||
					then["mapped_port_to"].(int) != 0 {
					return fmt.Errorf("only routing_instance need to be set in rule %s with type = inet", rule["name"].(string))
				}
				configSet = append(configSet, setPrefixRule+" then static-nat inet routing-instance "+
					then["routing_instance"].(string))
			}
			if then["type"].(string) == prefixWord || then["type"].(string) == prefixNameWord {
				setPrefixRuleThenStaticNat := setPrefixRule + " then static-nat "
				if then["type"].(string) == prefixWord {
					setPrefixRuleThenStaticNat += "prefix "
					if then["prefix"].(string) == "" {
						return fmt.Errorf("missing prefix in rule %s with type = prefix", rule["name"].(string))
					}
					if err := validateCIDRNetwork(then["prefix"].(string)); err != nil {
						return err
					}
				}
				if then["type"].(string) == prefixNameWord {
					setPrefixRuleThenStaticNat += "prefix-name "
					if then["prefix"].(string) == "" {
						return fmt.Errorf("missing prefix in rule %s with type = prefix-name", rule["name"].(string))
					}
				}
				configSet = append(configSet, setPrefixRuleThenStaticNat+"\""+then["prefix"].(string)+"\"")
				if vv := then["mapped_port"].(int); vv != 0 {
					configSet = append(configSet, setPrefixRuleThenStaticNat+"mapped-port "+strconv.Itoa(vv))
					if vvTo := then["mapped_port_to"].(int); vvTo != 0 {
						configSet = append(configSet, setPrefixRuleThenStaticNat+"mapped-port to "+strconv.Itoa(vvTo))
					}
				} else if then["mapped_port_to"].(int) != 0 {
					return fmt.Errorf("mapped_port need to set with mapped_port_to in rule %s", rule["name"].(string))
				}
				if vv := then["routing_instance"].(string); vv != "" {
					configSet = append(configSet, setPrefixRuleThenStaticNat+"routing-instance "+vv)
				}
			}
		}
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+" description \""+v+"\"")
	}

	return sess.configSet(configSet, jnprSess)
}

func readSecurityNatStatic(natStatic string, m interface{}, jnprSess *NetconfObject) (natStaticOptions, error) {
	sess := m.(*Session)
	var confRead natStaticOptions

	natStaticConfig, err := sess.command("show configuration"+
		" security nat static rule-set "+natStatic+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if natStaticConfig != emptyWord {
		confRead.name = natStatic
		for _, item := range strings.Split(natStaticConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "from "):
				fromWords := strings.Split(strings.TrimPrefix(itemTrim, "from "), " ")
				if len(confRead.from) == 0 {
					confRead.from = append(confRead.from, map[string]interface{}{
						"type":  fromWords[0],
						"value": make([]string, 0),
					})
				}
				confRead.from[0]["value"] = append(confRead.from[0]["value"].([]string), fromWords[1])
			case strings.HasPrefix(itemTrim, "rule "):
				ruleConfig := strings.Split(strings.TrimPrefix(itemTrim, "rule "), " ")
				ruleOptions := map[string]interface{}{
					"name":                     ruleConfig[0],
					"destination_address":      "",
					"destination_address_name": "",
					"destination_port":         0,
					"destination_port_to":      0,
					"source_address":           make([]string, 0),
					"source_address_name":      make([]string, 0),
					"source_port":              make([]string, 0),
					"then":                     make([]map[string]interface{}, 0),
				}
				confRead.rule = copyAndRemoveItemMapList("name", ruleOptions, confRead.rule)
				switch {
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match destination-address "):
					ruleOptions["destination_address"] = strings.TrimPrefix(itemTrim,
						"rule "+ruleConfig[0]+" match destination-address ")
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match destination-address-name "):
					ruleOptions["destination_address_name"] = strings.Trim(strings.TrimPrefix(itemTrim,
						"rule "+ruleConfig[0]+" match destination-address-name "), "\"")
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match destination-port to "):
					var err error
					ruleOptions["destination_port_to"], err = strconv.Atoi(strings.TrimPrefix(itemTrim,
						"rule "+ruleConfig[0]+" match destination-port to "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match destination-port "):
					var err error
					ruleOptions["destination_port"], err = strconv.Atoi(strings.TrimPrefix(itemTrim,
						"rule "+ruleConfig[0]+" match destination-port "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match source-address "):
					ruleOptions["source_address"] = append(ruleOptions["source_address"].([]string),
						strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" match source-address "))
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match source-address-name "):
					ruleOptions["source_address_name"] = append(ruleOptions["source_address_name"].([]string),
						strings.Trim(strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" match source-address-name "), "\""))
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match source-port "):
					ruleOptions["source_port"] = append(ruleOptions["source_port"].([]string),
						strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" match source-port "))
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" then static-nat "):
					itemThen := strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" then static-nat ")
					if len(ruleOptions["then"].([]map[string]interface{})) == 0 {
						ruleOptions["then"] = append(ruleOptions["then"].([]map[string]interface{}),
							map[string]interface{}{
								"type":             "",
								"mapped_port":      0,
								"mapped_port_to":   0,
								"prefix":           "",
								"routing_instance": "",
							})
					}
					ruleThenOptions := ruleOptions["then"].([]map[string]interface{})[0]
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
				confRead.rule = append(confRead.rule, ruleOptions)
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			}
		}
	}

	return confRead, nil
}

func delSecurityNatStatic(natStatic string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat static rule-set "+natStatic)

	return sess.configSet(configSet, jnprSess)
}

func fillSecurityNatStaticData(d *schema.ResourceData, natStaticOptions natStaticOptions) {
	if tfErr := d.Set("name", natStaticOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("from", natStaticOptions.from); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("rule", natStaticOptions.rule); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", natStaticOptions.description); tfErr != nil {
		panic(tfErr)
	}
}
