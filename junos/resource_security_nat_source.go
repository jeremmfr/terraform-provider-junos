package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type natSourceOptions struct {
	name string
	from []map[string]interface{}
	to   []map[string]interface{}
	rule []map[string]interface{}
}

func resourceSecurityNatSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityNatSourceCreate,
		ReadContext:   resourceSecurityNatSourceRead,
		UpdateContext: resourceSecurityNatSourceUpdate,
		DeleteContext: resourceSecurityNatSourceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityNatSourceImport,
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
									"destination_address": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"protocol": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"source_address": {
										Type:     schema.TypeList,
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
		},
	}
}

func resourceSecurityNatSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setSecurityNatSource(d, m, nil); err != nil {
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
		return diag.FromErr(fmt.Errorf("security nat source not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	securityNatSourceExists, err := checkSecurityNatSourceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatSourceExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat source %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityNatSource(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_nat_source", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityNatSourceExists, err = checkSecurityNatSourceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityNatSourceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security nat source %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityNatSourceReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityNatSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityNatSourceReadWJnprSess(d, m, jnprSess)
}

func resourceSecurityNatSourceReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	natSourceOptions, err := readSecurityNatSource(d.Get("name").(string), m, jnprSess)
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityNatSource(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityNatSource(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_nat_source", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityNatSourceReadWJnprSess(d, m, jnprSess)...)
}

func resourceSecurityNatSourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delSecurityNatSource(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_nat_source", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityNatSourceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	securityNatSourceExists, err := checkSecurityNatSourceExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityNatSourceExists {
		return nil, fmt.Errorf("don't find nat source with id '%v' (id must be <name>)", d.Id())
	}
	natSourceOptions, err := readSecurityNatSource(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatSourceData(d, natSourceOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatSourceExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	natSourceConfig, err := sess.command("show configuration"+
		" security nat source rule-set "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if natSourceConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setSecurityNatSource(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

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
	for _, v := range d.Get("rule").([]interface{}) {
		rule := v.(map[string]interface{})
		setPrefixRule := setPrefix + " rule " + rule["name"].(string)
		for _, matchV := range rule[matchWord].([]interface{}) {
			match := matchV.(map[string]interface{})
			for _, address := range match["destination_address"].([]interface{}) {
				err := validateCIDRNetwork(address.(string))
				if err != nil {
					return err
				}
				configSet = append(configSet, setPrefixRule+" match destination-address "+address.(string))
			}
			for _, proto := range match["protocol"].([]interface{}) {
				configSet = append(configSet, setPrefixRule+" match protocol "+proto.(string))
			}
			for _, address := range match["source_address"].([]interface{}) {
				err := validateCIDRNetwork(address.(string))
				if err != nil {
					return err
				}
				configSet = append(configSet, setPrefixRule+" match source-address "+address.(string))
			}
		}
		for _, thenV := range rule[thenWord].([]interface{}) {
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

	return sess.configSet(configSet, jnprSess)
}

func readSecurityNatSource(natSource string, m interface{}, jnprSess *NetconfObject) (natSourceOptions, error) {
	sess := m.(*Session)
	var confRead natSourceOptions

	natSourceConfig, err := sess.command("show configuration"+
		" security nat source rule-set "+natSource+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if natSourceConfig != emptyWord {
		confRead.name = natSource
		for _, item := range strings.Split(natSourceConfig, "\n") {
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
			case strings.HasPrefix(itemTrim, "to "):
				toWords := strings.Split(strings.TrimPrefix(itemTrim, "to "), " ")
				if len(confRead.to) == 0 {
					confRead.to = append(confRead.to, map[string]interface{}{
						"type":  toWords[0],
						"value": make([]string, 0),
					})
				}
				confRead.to[0]["value"] = append(confRead.to[0]["value"].([]string), toWords[1])
			case strings.HasPrefix(itemTrim, "rule "):
				ruleConfig := strings.Split(strings.TrimPrefix(itemTrim, "rule "), " ")
				ruleOptions := map[string]interface{}{
					"name":  ruleConfig[0],
					"match": make([]map[string]interface{}, 0),
					"then":  make([]map[string]interface{}, 0),
				}
				confRead.rule = copyAndRemoveItemMapList("name", ruleOptions, confRead.rule)
				switch {
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match "):
					itemTrimMatch := strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" match ")
					if len(ruleOptions["match"].([]map[string]interface{})) == 0 {
						ruleOptions["match"] = append(ruleOptions["match"].([]map[string]interface{}),
							map[string]interface{}{
								"destination_address": make([]string, 0),
								"protocol":            make([]string, 0),
								"source_address":      make([]string, 0),
							})
					}
					ruleMatchOptions := ruleOptions["match"].([]map[string]interface{})[0]
					switch {
					case strings.HasPrefix(itemTrimMatch, "destination-address "):
						ruleMatchOptions["destination_address"] = append(ruleMatchOptions["destination_address"].([]string),
							strings.TrimPrefix(itemTrimMatch, "destination-address "))
					case strings.HasPrefix(itemTrimMatch, "protocol "):
						ruleMatchOptions["protocol"] = append(ruleMatchOptions["protocol"].([]string),
							strings.TrimPrefix(itemTrimMatch, "protocol "))
					case strings.HasPrefix(itemTrimMatch, "source-address "):
						ruleMatchOptions["source_address"] = append(ruleMatchOptions["source_address"].([]string),
							strings.TrimPrefix(itemTrimMatch, "source-address "))
					}
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" then source-nat "):
					itemTrimThen := strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" then source-nat ")
					if len(ruleOptions["then"].([]map[string]interface{})) == 0 {
						ruleOptions["then"] = append(ruleOptions["then"].([]map[string]interface{}),
							map[string]interface{}{
								"type": "",
								"pool": "",
							})
					}
					ruleThenOptions := ruleOptions["then"].([]map[string]interface{})[0]
					if strings.HasPrefix(itemTrimThen, "pool ") {
						thenSplit := strings.Split(itemTrimThen, " ")
						ruleThenOptions["type"] = thenSplit[0]
						ruleThenOptions["pool"] = thenSplit[1]
					} else {
						ruleThenOptions["type"] = itemTrimThen
					}
				}
				confRead.rule = append(confRead.rule, ruleOptions)
			}
		}
	}

	return confRead, nil
}

func delSecurityNatSource(natSource string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat source rule-set "+natSource)

	return sess.configSet(configSet, jnprSess)
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
}
