package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type natStaticOptions struct {
	name        string
	description string
	from        []map[string]interface{}
	rule        []map[string]interface{}
}

func resourceSecurityNatStatic() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityNatStaticCreate,
		ReadWithoutTimeout:   resourceSecurityNatStaticRead,
		UpdateWithoutTimeout: resourceSecurityNatStaticUpdate,
		DeleteWithoutTimeout: resourceSecurityNatStaticDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityNatStaticImport,
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
			"configure_rules_singly": {
				Type:         schema.TypeBool,
				Optional:     true,
				ExactlyOneOf: []string{"configure_rules_singly", "rule"},
			},
			"rule": {
				Type:         schema.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"rule", "configure_rules_singly"},
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
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateCIDRNetworkFunc(),
							},
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
										ValidateFunc: validation.StringInSlice([]string{"inet", "prefix", "prefix-name"}, false),
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityNatStatic(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security nat static not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityNatStaticExists, err := checkSecurityNatStaticExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatStaticExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat static %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityNatStatic(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_nat_static")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityNatStaticExists, err = checkSecurityNatStaticExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatStaticExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat static %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatStaticReadWJunSess(d, junSess)...)
}

func resourceSecurityNatStaticRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityNatStaticReadWJunSess(d, junSess)
}

func resourceSecurityNatStaticReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	natStaticOptions, err := readSecurityNatStatic(d.Get("name").(string), junSess)
	junos.MutexUnlock()
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
	var diagWarns diag.Diagnostics
	configureRulesSingly := d.Get("configure_rules_singly").(bool)
	if d.HasChange("configure_rules_singly") {
		if o, _ := d.GetChange("configure_rules_singly"); o.(bool) {
			configureRulesSingly = o.(bool)
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary: "Disable configure_rules_singly on resource already created doesn't " +
					"delete rule(s) already configured.",
				Detail:        "So refresh resource after apply to detect rule(s) that need to be deleted",
				AttributePath: cty.Path{cty.GetAttrStep{Name: "configure_rules_singly"}},
			})
		} else {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary: "Enable configure_rules_singly on resource already created doesn't " +
					"delete rules already configured.",
				Detail:        "So import rule(s) in dedicated resource(s) to be able to manage them",
				AttributePath: cty.Path{cty.GetAttrStep{Name: "configure_rules_singly"}},
			})
		}
	}
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if configureRulesSingly {
			if err := delSecurityNatStaticWithoutRules(d.Get("name").(string), junSess); err != nil {
				return append(diagWarns, diag.FromErr(err)...)
			}
		} else {
			if err := delSecurityNatStatic(d.Get("name").(string), junSess); err != nil {
				return append(diagWarns, diag.FromErr(err)...)
			}
		}
		if err := setSecurityNatStatic(d, junSess); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		d.Partial(false)

		return diagWarns
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	if configureRulesSingly {
		if err := delSecurityNatStaticWithoutRules(d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	} else {
		if err := delSecurityNatStatic(d.Get("name").(string), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	}
	if err := setSecurityNatStatic(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_nat_static")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatStaticReadWJunSess(d, junSess)...)
}

func resourceSecurityNatStaticDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityNatStatic(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityNatStatic(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_nat_static")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatStaticImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), junos.IDSeparator)
	securityNatStaticExists, err := checkSecurityNatStaticExists(idList[0], junSess)
	if err != nil {
		return nil, err
	}
	if !securityNatStaticExists {
		return nil, fmt.Errorf("don't find nat static with id '%v' "+
			"(id must be <name> or <name>"+junos.IDSeparator+"no_rules)", idList[0])
	}
	natStaticOptions, err := readSecurityNatStatic(idList[0], junSess)
	if err != nil {
		return nil, err
	}
	if len(idList) > 1 && idList[1] == "no_rules" {
		if tfErr := d.Set("configure_rules_singly", true); tfErr != nil {
			panic(tfErr)
		}
	}
	fillSecurityNatStaticData(d, natStaticOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatStaticExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat static rule-set " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityNatStatic(d *schema.ResourceData, junSess *junos.Session) error {
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
	if !d.Get("configure_rules_singly").(bool) {
		for _, v := range d.Get("rule").([]interface{}) {
			rule := v.(map[string]interface{})
			if bchk.InSlice(rule["name"].(string), ruleNameList) {
				return fmt.Errorf("multiple blocks rule with the same name %s", rule["name"].(string))
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
			for _, thenV := range rule["then"].([]interface{}) {
				then := thenV.(map[string]interface{})
				if then["type"].(string) == junos.InetW {
					if then["prefix"].(string) != "" ||
						then["mapped_port"].(int) != 0 ||
						then["mapped_port_to"].(int) != 0 {
						return fmt.Errorf("only routing_instance can be set in rule %s with type = inet", rule["name"].(string))
					}
					configSet = append(configSet, setPrefixRule+" then static-nat inet")
					if rI := then["routing_instance"].(string); rI != "" {
						configSet = append(configSet, setPrefixRule+" then static-nat inet routing-instance "+rI)
					}
				}
				if then["type"].(string) == "prefix" || then["type"].(string) == "prefix-name" {
					setPrefixRuleThenStaticNat := setPrefixRule + " then static-nat "
					if then["type"].(string) == "prefix" {
						setPrefixRuleThenStaticNat += "prefix "
						if then["prefix"].(string) == "" {
							return fmt.Errorf("missing prefix in rule %s with type = prefix", rule["name"].(string))
						}
						if err := validateCIDRNetwork(then["prefix"].(string)); err != nil {
							return err
						}
					}
					if then["type"].(string) == "prefix-name" {
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
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+" description \""+v+"\"")
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityNatStatic(name string, junSess *junos.Session,
) (confRead natStaticOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat static rule-set " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "from "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <type> <value>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "from", itemTrim)
				}
				if len(confRead.from) == 0 {
					confRead.from = append(confRead.from, map[string]interface{}{
						"type":  itemTrimFields[0],
						"value": make([]string, 0),
					})
				}
				confRead.from[0]["value"] = append(confRead.from[0]["value"].([]string), itemTrimFields[1])
			case balt.CutPrefixInString(&itemTrim, "rule "):
				itemTrimFields := strings.Split(itemTrim, " ")
				ruleOptions := map[string]interface{}{
					"name":                     itemTrimFields[0],
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
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "match destination-address "):
					ruleOptions["destination_address"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "match destination-address-name "):
					ruleOptions["destination_address_name"] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "match destination-port to "):
					ruleOptions["destination_port_to"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "match destination-port "):
					ruleOptions["destination_port"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "match source-address "):
					ruleOptions["source_address"] = append(
						ruleOptions["source_address"].([]string),
						itemTrim,
					)
				case balt.CutPrefixInString(&itemTrim, "match source-address-name "):
					ruleOptions["source_address_name"] = append(
						ruleOptions["source_address_name"].([]string),
						strings.Trim(itemTrim, "\""),
					)
				case balt.CutPrefixInString(&itemTrim, "match source-port "):
					ruleOptions["source_port"] = append(
						ruleOptions["source_port"].([]string),
						itemTrim,
					)
				case balt.CutPrefixInString(&itemTrim, "then static-nat "):
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
					case balt.CutPrefixInString(&itemTrim, "prefix"):
						ruleThenOptions["type"] = "prefix"
						if balt.CutPrefixInString(&itemTrim, "-name") {
							ruleThenOptions["type"] = "prefix-name"
						}
						balt.CutPrefixInString(&itemTrim, " ")
						switch {
						case balt.CutPrefixInString(&itemTrim, "routing-instance "):
							ruleThenOptions["routing_instance"] = itemTrim
						case balt.CutPrefixInString(&itemTrim, "mapped-port to "):
							ruleThenOptions["mapped_port_to"], err = strconv.Atoi(itemTrim)
							if err != nil {
								return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
							}
						case balt.CutPrefixInString(&itemTrim, "mapped-port "):
							ruleThenOptions["mapped_port"], err = strconv.Atoi(itemTrim)
							if err != nil {
								return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
							}
						default:
							ruleThenOptions["prefix"] = strings.Trim(itemTrim, "\"")
						}
					case balt.CutPrefixInString(&itemTrim, junos.InetW):
						ruleThenOptions["type"] = junos.InetW
						if balt.CutPrefixInString(&itemTrim, " routing-instance ") {
							ruleThenOptions["routing_instance"] = itemTrim
						}
					}
				}
				confRead.rule = append(confRead.rule, ruleOptions)
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			}
		}
	}

	return confRead, nil
}

func delSecurityNatStatic(natStatic string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat static rule-set "+natStatic)

	return junSess.ConfigSet(configSet)
}

func delSecurityNatStaticWithoutRules(natStatic string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat static rule-set "+natStatic+" description")
	configSet = append(configSet, "delete security nat static rule-set "+natStatic+" from")

	return junSess.ConfigSet(configSet)
}

func fillSecurityNatStaticData(d *schema.ResourceData, natStaticOptions natStaticOptions) {
	if tfErr := d.Set("name", natStaticOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("from", natStaticOptions.from); tfErr != nil {
		panic(tfErr)
	}
	if !d.Get("configure_rules_singly").(bool) {
		if tfErr := d.Set("rule", natStaticOptions.rule); tfErr != nil {
			panic(tfErr)
		}
	}
	if tfErr := d.Set("description", natStaticOptions.description); tfErr != nil {
		panic(tfErr)
	}
}
