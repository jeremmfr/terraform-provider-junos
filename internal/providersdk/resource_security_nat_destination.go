package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type natDestinationOptions struct {
	name        string
	description string
	from        []map[string]interface{}
	rule        []map[string]interface{}
}

func resourceSecurityNatDestination() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityNatDestinationCreate,
		ReadWithoutTimeout:   resourceSecurityNatDestinationRead,
		UpdateWithoutTimeout: resourceSecurityNatDestinationUpdate,
		DeleteWithoutTimeout: resourceSecurityNatDestinationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityNatDestinationImport,
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
						"application": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"destination_port": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"protocol": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
						"then": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"off", "pool"}, false),
									},
									"pool": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
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

func resourceSecurityNatDestinationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityNatDestination(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security nat destination not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityNatDestinationExists, err := checkSecurityNatDestinationExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatDestinationExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security nat destination %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityNatDestination(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_nat_destination")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityNatDestinationExists, err = checkSecurityNatDestinationExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatDestinationExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat destination %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatDestinationReadWJunSess(d, junSess)...)
}

func resourceSecurityNatDestinationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityNatDestinationReadWJunSess(d, junSess)
}

func resourceSecurityNatDestinationReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	natDestinationOptions, err := readSecurityNatDestination(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if natDestinationOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatDestinationData(d, natDestinationOptions)
	}

	return nil
}

func resourceSecurityNatDestinationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityNatDestination(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityNatDestination(d, junSess); err != nil {
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
	if err := delSecurityNatDestination(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityNatDestination(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_nat_destination")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatDestinationReadWJunSess(d, junSess)...)
}

func resourceSecurityNatDestinationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityNatDestination(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityNatDestination(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_nat_destination")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatDestinationImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	securityNatDestinationExists, err := checkSecurityNatDestinationExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityNatDestinationExists {
		return nil, fmt.Errorf("don't find nat destination with id '%v' (id must be <name>)", d.Id())
	}
	natDestinationOptions, err := readSecurityNatDestination(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatDestinationData(d, natDestinationOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatDestinationExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat destination rule-set " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityNatDestination(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	regexpDestPort := regexp.MustCompile(`^\d+( to \d+)?$`)

	setPrefix := "set security nat destination rule-set " + d.Get("name").(string)
	for _, v := range d.Get("from").([]interface{}) {
		from := v.(map[string]interface{})
		for _, value := range sortSetOfString(from["value"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+" from "+from["type"].(string)+" "+value)
		}
	}
	ruleNameList := make([]string, 0)
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
		for _, vv := range sortSetOfString(rule["application"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixRule+" match application \""+vv+"\"")
		}
		for _, vv := range sortSetOfString(rule["destination_port"].(*schema.Set).List()) {
			if !regexpDestPort.MatchString(vv) {
				return fmt.Errorf("destination_port need to have format `x` or `x to y` in rule %s", rule["name"].(string))
			}
			configSet = append(configSet, setPrefixRule+" match destination-port "+vv)
		}
		for _, vv := range sortSetOfString(rule["protocol"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixRule+" match protocol "+vv)
		}
		for _, vv := range sortSetOfString(rule["source_address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixRule+" match source-address "+vv)
		}
		for _, vv := range sortSetOfString(rule["source_address_name"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixRule+" match source-address-name \""+vv+"\"")
		}
		for _, thenV := range rule["then"].([]interface{}) {
			then := thenV.(map[string]interface{})
			if then["type"].(string) == "off" {
				configSet = append(configSet, setPrefixRule+" then destination-nat off")
			}
			if then["type"].(string) == "pool" {
				if then["pool"].(string) == "" {
					return fmt.Errorf("missing pool for destination-nat pool for rule %v in %v",
						then["name"].(string), d.Get("name").(string))
				}
				configSet = append(configSet, setPrefixRule+" then destination-nat pool "+then["pool"].(string))
			}
		}
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+" description \""+v+"\"")
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityNatDestination(name string, junSess *junos.Session,
) (confRead natDestinationOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat destination rule-set " + name + junos.PipeDisplaySetRelative)
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
					"application":              make([]string, 0),
					"destination_port":         make([]string, 0),
					"protocol":                 make([]string, 0),
					"source_address":           make([]string, 0),
					"source_address_name":      make([]string, 0),
					"then":                     make([]map[string]interface{}, 0),
				}
				confRead.rule = copyAndRemoveItemMapList("name", ruleOptions, confRead.rule)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "match destination-address "):
					ruleOptions["destination_address"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "match destination-address-name "):
					ruleOptions["destination_address_name"] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "match application "):
					ruleOptions["application"] = append(
						ruleOptions["application"].([]string),
						strings.Trim(itemTrim, "\""),
					)
				case balt.CutPrefixInString(&itemTrim, "match destination-port "):
					ruleOptions["destination_port"] = append(
						ruleOptions["destination_port"].([]string),
						itemTrim,
					)
				case balt.CutPrefixInString(&itemTrim, "match protocol "):
					ruleOptions["protocol"] = append(
						ruleOptions["protocol"].([]string),
						itemTrim,
					)
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
				case balt.CutPrefixInString(&itemTrim, "then destination-nat "):
					if len(ruleOptions["then"].([]map[string]interface{})) == 0 {
						ruleOptions["then"] = append(ruleOptions["then"].([]map[string]interface{}), map[string]interface{}{
							"type": "",
							"pool": "",
						})
					}
					ruleThenOptions := ruleOptions["then"].([]map[string]interface{})[0]
					if balt.CutPrefixInString(&itemTrim, "pool ") {
						ruleThenOptions["type"] = "pool"
						ruleThenOptions["pool"] = itemTrim
					} else {
						ruleThenOptions["type"] = itemTrim
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

func delSecurityNatDestination(natDestination string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat destination rule-set "+natDestination)

	return junSess.ConfigSet(configSet)
}

func fillSecurityNatDestinationData(d *schema.ResourceData, natDestinationOptions natDestinationOptions) {
	if tfErr := d.Set("name", natDestinationOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("from", natDestinationOptions.from); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("rule", natDestinationOptions.rule); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", natDestinationOptions.description); tfErr != nil {
		panic(tfErr)
	}
}
