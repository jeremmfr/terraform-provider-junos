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

type natSourceOptions struct {
	name        string
	description string
	from        []map[string]interface{}
	to          []map[string]interface{}
	rule        []map[string]interface{}
}

func resourceSecurityNatSource() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityNatSourceCreate,
		ReadWithoutTimeout:   resourceSecurityNatSourceRead,
		UpdateWithoutTimeout: resourceSecurityNatSourceUpdate,
		DeleteWithoutTimeout: resourceSecurityNatSourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityNatSourceImport,
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
			"to": {
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
						"match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_address": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validateCIDRNetworkFunc(),
										},
									},
									"destination_address_name": {
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
									"source_port": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
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
										ValidateFunc: validation.StringInSlice([]string{"interface", "pool", "off"}, false),
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

func resourceSecurityNatSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityNatSource(d, junSess); err != nil {
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
		return diag.FromErr(fmt.Errorf("security nat source not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityNatSourceExists, err := checkSecurityNatSourceExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatSourceExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat source %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityNatSource(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_nat_source")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityNatSourceExists, err = checkSecurityNatSourceExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatSourceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat source %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatSourceReadWJunSess(d, junSess)...)
}

func resourceSecurityNatSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityNatSourceReadWJunSess(d, junSess)
}

func resourceSecurityNatSourceReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	natSourceOptions, err := readSecurityNatSource(d.Get("name").(string), junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if natSourceOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatSourceData(d, natSourceOptions)
	}

	return nil
}

func resourceSecurityNatSourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityNatSource(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityNatSource(d, junSess); err != nil {
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
	if err := delSecurityNatSource(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityNatSource(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_nat_source")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatSourceReadWJunSess(d, junSess)...)
}

func resourceSecurityNatSourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityNatSource(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityNatSource(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_nat_source")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatSourceImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	securityNatSourceExists, err := checkSecurityNatSourceExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityNatSourceExists {
		return nil, fmt.Errorf("don't find nat source with id '%v' (id must be <name>)", d.Id())
	}
	natSourceOptions, err := readSecurityNatSource(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatSourceData(d, natSourceOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatSourceExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat source rule-set " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityNatSource(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	regexpPort := regexp.MustCompile(`^\d+( to \d+)?$`)

	setPrefix := "set security nat source rule-set " + d.Get("name").(string)
	for _, v := range d.Get("from").([]interface{}) {
		from := v.(map[string]interface{})
		for _, value := range sortSetOfString(from["value"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+" from "+from["type"].(string)+" "+value)
		}
	}
	for _, v := range d.Get("to").([]interface{}) {
		to := v.(map[string]interface{})
		for _, value := range sortSetOfString(to["value"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+" to "+to["type"].(string)+" "+value)
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
		for _, matchV := range rule["match"].([]interface{}) {
			if matchV == nil {
				return fmt.Errorf("match block in rule %s need to have an argument", rule["name"].(string))
			}
			match := matchV.(map[string]interface{})
			if len(match["destination_address"].(*schema.Set).List()) == 0 &&
				len(match["destination_address_name"].(*schema.Set).List()) == 0 &&
				len(match["source_address"].(*schema.Set).List()) == 0 &&
				len(match["source_address_name"].(*schema.Set).List()) == 0 {
				return fmt.Errorf("one of destination_address, destination_address_name, " +
					"source_address or source_address_name arguments must be set")
			}
			for _, vv := range sortSetOfString(match["application"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefixRule+" match application \""+vv+"\"")
			}
			for _, address := range sortSetOfString(match["destination_address"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefixRule+" match destination-address "+address)
			}
			for _, vv := range sortSetOfString(match["destination_address_name"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefixRule+" match destination-address-name \""+vv+"\"")
			}
			for _, vv := range sortSetOfString(match["destination_port"].(*schema.Set).List()) {
				if !regexpPort.MatchString(vv) {
					return fmt.Errorf("destination_port need to have format `x` or `x to y` in rule %s", rule["name"].(string))
				}
				configSet = append(configSet, setPrefixRule+" match destination-port "+vv)
			}
			for _, proto := range sortSetOfString(match["protocol"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefixRule+" match protocol "+proto)
			}
			for _, address := range sortSetOfString(match["source_address"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefixRule+" match source-address "+address)
			}
			for _, vv := range sortSetOfString(match["source_address_name"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefixRule+" match source-address-name \""+vv+"\"")
			}
			for _, vv := range sortSetOfString(match["source_port"].(*schema.Set).List()) {
				if !regexpPort.MatchString(vv) {
					return fmt.Errorf("source_port need to have format `x` or `x to y` in rule %s", rule["name"].(string))
				}
				configSet = append(configSet, setPrefixRule+" match source-port "+vv)
			}
		}
		for _, thenV := range rule["then"].([]interface{}) {
			then := thenV.(map[string]interface{})
			if then["type"].(string) == "interface" {
				configSet = append(configSet, setPrefixRule+" then source-nat interface")
			}
			if then["type"].(string) == "off" {
				configSet = append(configSet, setPrefixRule+" then source-nat off")
			}
			if then["type"].(string) == "pool" {
				if then["pool"].(string) == "" {
					return fmt.Errorf("missing pool for source-nat pool for rule %v in %v",
						rule["name"].(string), d.Get("name").(string))
				}
				configSet = append(configSet, setPrefixRule+" then source-nat pool "+then["pool"].(string))
			}
		}
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+" description \""+v+"\"")
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityNatSource(name string, junSess *junos.Session,
) (confRead natSourceOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat source rule-set " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "to "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <type> <value>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "to", itemTrim)
				}
				if len(confRead.to) == 0 {
					confRead.to = append(confRead.to, map[string]interface{}{
						"type":  itemTrimFields[0],
						"value": make([]string, 0),
					})
				}
				confRead.to[0]["value"] = append(confRead.to[0]["value"].([]string), itemTrimFields[1])
			case balt.CutPrefixInString(&itemTrim, "rule "):
				itemTrimFields := strings.Split(itemTrim, " ")
				ruleOptions := map[string]interface{}{
					"name":  itemTrimFields[0],
					"match": make([]map[string]interface{}, 0),
					"then":  make([]map[string]interface{}, 0),
				}
				confRead.rule = copyAndRemoveItemMapList("name", ruleOptions, confRead.rule)
				switch {
				case balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" match "):
					if len(ruleOptions["match"].([]map[string]interface{})) == 0 {
						ruleOptions["match"] = append(ruleOptions["match"].([]map[string]interface{}),
							map[string]interface{}{
								"application":              make([]string, 0),
								"destination_address":      make([]string, 0),
								"destination_address_name": make([]string, 0),
								"destination_port":         make([]string, 0),
								"protocol":                 make([]string, 0),
								"source_address":           make([]string, 0),
								"source_address_name":      make([]string, 0),
								"source_port":              make([]string, 0),
							})
					}
					ruleMatchOptions := ruleOptions["match"].([]map[string]interface{})[0]
					switch {
					case balt.CutPrefixInString(&itemTrim, "application "):
						ruleMatchOptions["application"] = append(
							ruleMatchOptions["application"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-address "):
						ruleMatchOptions["destination_address"] = append(
							ruleMatchOptions["destination_address"].([]string),
							itemTrim,
						)
					case balt.CutPrefixInString(&itemTrim, "destination-address-name "):
						ruleMatchOptions["destination_address_name"] = append(
							ruleMatchOptions["destination_address_name"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-port "):
						ruleMatchOptions["destination_port"] = append(
							ruleMatchOptions["destination_port"].([]string),
							itemTrim,
						)
					case balt.CutPrefixInString(&itemTrim, "protocol "):
						ruleMatchOptions["protocol"] = append(
							ruleMatchOptions["protocol"].([]string),
							itemTrim,
						)
					case balt.CutPrefixInString(&itemTrim, "source-address "):
						ruleMatchOptions["source_address"] = append(
							ruleMatchOptions["source_address"].([]string),
							itemTrim,
						)
					case balt.CutPrefixInString(&itemTrim, "source-address-name "):
						ruleMatchOptions["source_address_name"] = append(
							ruleMatchOptions["source_address_name"].([]string),
							strings.Trim(itemTrim, "\""),
						)
					case balt.CutPrefixInString(&itemTrim, "source-port "):
						ruleMatchOptions["source_port"] = append(
							ruleMatchOptions["source_port"].([]string),
							itemTrim,
						)
					}
				case balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" then source-nat "):
					if len(ruleOptions["then"].([]map[string]interface{})) == 0 {
						ruleOptions["then"] = append(ruleOptions["then"].([]map[string]interface{}),
							map[string]interface{}{
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

func delSecurityNatSource(natSource string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat source rule-set "+natSource)

	return junSess.ConfigSet(configSet)
}

func fillSecurityNatSourceData(d *schema.ResourceData, natSourceOptions natSourceOptions) {
	if tfErr := d.Set("name", natSourceOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("from", natSourceOptions.from); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("to", natSourceOptions.to); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("rule", natSourceOptions.rule); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", natSourceOptions.description); tfErr != nil {
		panic(tfErr)
	}
}
