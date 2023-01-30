package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
		CreateWithoutTimeout: resourceSecurityNatStaticRuleCreate,
		ReadWithoutTimeout:   resourceSecurityNatStaticRuleRead,
		UpdateWithoutTimeout: resourceSecurityNatStaticRuleUpdate,
		DeleteWithoutTimeout: resourceSecurityNatStaticRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityNatStaticRuleImport,
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

func resourceSecurityNatStaticRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityNatStaticRule(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("rule_set").(string) + junos.IDSeparator + d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security nat static rule in rule-set not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	natStaticExists, err := checkSecurityNatStaticExists(d.Get("rule_set").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !natStaticExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security nat static rule-set %v doesn't exist", d.Get("rule_set").(string)))...)
	}
	natStaticRuleExists, err := checkSecurityNatStaticRuleExists(
		d.Get("rule_set").(string),
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if natStaticRuleExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security nat static rule %v already exists in rule-set %s",
			d.Get("name").(string), d.Get("rule_set").(string)))...)
	}

	if err := setSecurityNatStaticRule(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_nat_static_rule")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	natStaticRuleExists, err = checkSecurityNatStaticRuleExists(
		d.Get("rule_set").(string),
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if natStaticRuleExists {
		d.SetId(d.Get("rule_set").(string) + junos.IDSeparator + d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security nat statuc rule %v not exists in rule-set %s after commit "+
				"=> check your config", d.Get("name").(string), d.Get("rule_set").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatStaticRuleReadWJunSess(d, junSess)...)
}

func resourceSecurityNatStaticRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityNatStaticRuleReadWJunSess(d, junSess)
}

func resourceSecurityNatStaticRuleReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	natStaticRuleOptions, err := readSecurityNatStaticRule(
		d.Get("rule_set").(string),
		d.Get("name").(string),
		junSess,
	)
	junos.MutexUnlock()
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

func resourceSecurityNatStaticRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityNatStaticRule(d.Get("rule_set").(string), d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityNatStaticRule(d, junSess); err != nil {
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
	if err := delSecurityNatStaticRule(d.Get("rule_set").(string), d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityNatStaticRule(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_nat_static_rule")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatStaticRuleReadWJunSess(d, junSess)...)
}

func resourceSecurityNatStaticRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityNatStaticRule(d.Get("rule_set").(string), d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityNatStaticRule(d.Get("rule_set").(string), d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_nat_static_rule")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatStaticRuleImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), junos.IDSeparator)
	if len(idList) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	natStaticRuleExists, err := checkSecurityNatStaticRuleExists(idList[0], idList[1], junSess)
	if err != nil {
		return nil, err
	}
	if !natStaticRuleExists {
		return nil, fmt.Errorf(
			"don't find static nat rule with id '%v' (id must be <rule_set>"+junos.IDSeparator+"<name>)", d.Id())
	}
	natStaticRuleOptions, err := readSecurityNatStaticRule(idList[0], idList[1], junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatStaticRuleData(d, natStaticRuleOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatStaticRuleExists(ruleSet, name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat static rule-set " + ruleSet + " rule " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityNatStaticRule(d *schema.ResourceData, junSess *junos.Session) error {
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
		if then["type"].(string) == junos.InetW {
			if then["prefix"].(string) != "" ||
				then["mapped_port"].(int) != 0 ||
				then["mapped_port_to"].(int) != 0 {
				return fmt.Errorf("only routing_instance can be set with type = inet")
			}
			configSet = append(configSet, setPrefix+"then static-nat inet")
			if rI := then["routing_instance"].(string); rI != "" {
				configSet = append(configSet, setPrefix+"then static-nat inet routing-instance "+rI)
			}
		}
		if then["type"].(string) == "prefix" || then["type"].(string) == "prefix-name" {
			setPrefixRuleThenStaticNat := setPrefix + "then static-nat "
			if then["type"].(string) == "prefix" {
				setPrefixRuleThenStaticNat += "prefix "
				if then["prefix"].(string) == "" {
					return fmt.Errorf("missing prefix with type = prefix")
				}
				if err := validateCIDRNetwork(then["prefix"].(string)); err != nil {
					return err
				}
			}
			if then["type"].(string) == "prefix-name" {
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

	return junSess.ConfigSet(configSet)
}

func readSecurityNatStaticRule(ruleSet, name string, junSess *junos.Session,
) (confRead natStaticRuleOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat static rule-set " + ruleSet + " rule " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		confRead.ruleSet = ruleSet
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "match destination-address "):
				confRead.destinationAddress = itemTrim
			case balt.CutPrefixInString(&itemTrim, "match destination-address-name "):
				confRead.destinationAddressName = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "match destination-port to "):
				confRead.destinationPortTo, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "match destination-port "):
				confRead.destinationPort, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "match source-address "):
				confRead.sourceAddress = append(confRead.sourceAddress, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "match source-address-name "):
				confRead.sourceAddressName = append(confRead.sourceAddressName, strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "match source-port "):
				confRead.sourcePort = append(confRead.sourcePort, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "then static-nat "):
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
							return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
						}
					case balt.CutPrefixInString(&itemTrim, "mapped-port "):
						ruleThenOptions["mapped_port"], err = strconv.Atoi(itemTrim)
						if err != nil {
							return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
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
		}
	}

	return confRead, nil
}

func delSecurityNatStaticRule(ruleSet, name string, junSess *junos.Session) error {
	configSet := []string{"delete security nat static rule-set " + ruleSet + " rule " + name}

	return junSess.ConfigSet(configSet)
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
