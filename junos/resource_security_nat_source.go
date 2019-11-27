package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type natSourceOptions struct {
	name string
	from []map[string]interface{}
	to   []map[string]interface{}
	rule []map[string]interface{}
}

func resourceSecurityNatSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityNatSourceCreate,
		Read:   resourceSecurityNatSourceRead,
		Update: resourceSecurityNatSourceUpdate,
		Delete: resourceSecurityNatSourceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityNatSourceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"from": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"interface", "routing-instance", "zone"}) {
									errors = append(errors, fmt.Errorf(
										"%q for %q is not 'interface', 'routing-instance' or 'zone'", value, k))
								}
								return
							},
						},
						"value": {
							Type:     schema.TypeList,
							Required: true,
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
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"interface", "routing-instance", "zone"}) {
									errors = append(errors, fmt.Errorf(
										"%q for %q is not 'interface', 'routing-instance' or 'zone'", value, k))
								}
								return
							},
						},
						"value": {
							Type:     schema.TypeList,
							Required: true,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
						"match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_address": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
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
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(string)
											if !stringInSlice(value, []string{"interface", "pool", "off"}) {
												errors = append(errors, fmt.Errorf(
													"%q for %q is not 'interface', 'pool' or 'off'", value, k))
											}
											return
										},
									},
									"pool": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateNameObjectJunos(),
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

func resourceSecurityNatSourceCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security nat source not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	securityNatSourceExists, err := checkSecurityNatSourceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if securityNatSourceExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("security nat source %v already exists", d.Get("name").(string))
	}

	err = setSecurityNatSource(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_security_nat_source", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	securityNatSourceExists, err = checkSecurityNatSourceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if securityNatSourceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security nat source %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceSecurityNatSourceRead(d, m)
}
func resourceSecurityNatSourceRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	natSourceOptions, err := readSecurityNatSource(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if natSourceOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatSourceData(d, natSourceOptions)
	}
	return nil
}
func resourceSecurityNatSourceUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delSecurityNatSource(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setSecurityNatSource(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_security_nat_source", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceSecurityNatSourceRead(d, m)
}
func resourceSecurityNatSourceDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delSecurityNatSource(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_security_nat_source", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
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
		for _, value := range from["value"].([]interface{}) {
			configSet = append(configSet, setPrefix+" from "+from["type"].(string)+" "+value.(string)+"\n")
		}
	}
	for _, v := range d.Get("to").([]interface{}) {
		to := v.(map[string]interface{})
		for _, value := range to["value"].([]interface{}) {
			configSet = append(configSet, setPrefix+" to "+to["type"].(string)+" "+value.(string)+"\n")
		}
	}
	for _, v := range d.Get("rule").([]interface{}) {
		rule := v.(map[string]interface{})
		setPrefixRule := setPrefix + " rule " + rule["name"].(string)
		for _, matchV := range rule[matchWord].([]interface{}) {
			match := matchV.(map[string]interface{})
			for _, address := range match["source_address"].([]interface{}) {
				err := validateNetwork(address.(string))
				if err != nil {
					return err
				}
				configSet = append(configSet, setPrefixRule+" match source-address "+address.(string)+"\n")
			}
			for _, address := range match["destination_address"].([]interface{}) {
				err := validateNetwork(address.(string))
				if err != nil {
					return err
				}
				configSet = append(configSet, setPrefixRule+" match destination-address "+address.(string)+"\n")
			}
			for _, proto := range match["protocol"].([]interface{}) {
				configSet = append(configSet, setPrefixRule+" match protocol "+proto.(string)+"\n")
			}
		}
		for _, thenV := range rule[thenWord].([]interface{}) {
			then := thenV.(map[string]interface{})
			if then["type"].(string) == "interface" {
				configSet = append(configSet, setPrefixRule+" then source-nat interface\n")
			}
			if then["type"].(string) == "off" {
				configSet = append(configSet, setPrefixRule+" then source-nat off\n")
			}
			if then["type"].(string) == "pool" {
				if then["pool"].(string) == "" {
					return fmt.Errorf("missing pool for source-nat pool for rule %v in %v",
						rule["name"].(string), d.Get("name").(string))
				}
				configSet = append(configSet, setPrefixRule+" then source-nat pool "+then["pool"].(string)+"\n")
			}
		}
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
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
				fromOptions := map[string]interface{}{
					"type":  "",
					"value": []string{},
				}
				if len(confRead.from) > 0 {
					for k, v := range confRead.from[0] {
						fromOptions[k] = v
					}
				}
				fromWords := strings.Split(strings.TrimPrefix(itemTrim, "from "), " ")
				fromOptions["type"] = fromWords[0]
				fromOptions["value"] = append(fromOptions["value"].([]string), fromWords[1])
				confRead.from = []map[string]interface{}{fromOptions}
			case strings.HasPrefix(itemTrim, "to "):
				toOptions := map[string]interface{}{
					"type":  "",
					"value": []string{},
				}
				if len(confRead.to) > 0 {
					for k, v := range confRead.to[0] {
						toOptions[k] = v
					}
				}
				toWords := strings.Split(strings.TrimPrefix(itemTrim, "to "), " ")
				toOptions["type"] = toWords[0]
				toOptions["value"] = append(toOptions["value"].([]string), toWords[1])
				confRead.to = []map[string]interface{}{toOptions}
			case strings.HasPrefix(itemTrim, "rule "):
				ruleConfig := strings.Split(strings.TrimPrefix(itemTrim, "rule "), " ")

				ruleOptions := map[string]interface{}{
					"name":    ruleConfig[0],
					matchWord: make([]map[string]interface{}, 0),
					thenWord:  make([]map[string]interface{}, 0),
				}
				ruleOptions, confRead.rule = copyAndRemoveItemMapList("name", false, ruleOptions, confRead.rule)
				switch {
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match "):
					itemTrimMatch := strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" match ")
					ruleMatchOptions := map[string]interface{}{
						"source_address":      []string{},
						"destination_address": []string{},
						"protocol":            []string{},
					}
					if len(ruleOptions[matchWord].([]map[string]interface{})) > 0 {
						for k, v := range ruleOptions[matchWord].([]map[string]interface{})[0] {
							ruleMatchOptions[k] = v
						}
					}
					switch {
					case strings.HasPrefix(itemTrimMatch, "source-address "):
						ruleMatchOptions["source_address"] = append(ruleMatchOptions["source_address"].([]string),
							strings.TrimPrefix(itemTrimMatch, "source-address "))
					case strings.HasPrefix(itemTrimMatch, "destination-address "):
						ruleMatchOptions["destination_address"] = append(ruleMatchOptions["destination_address"].([]string),
							strings.TrimPrefix(itemTrimMatch, "destination-address "))
					case strings.HasPrefix(itemTrimMatch, "protocol "):
						ruleMatchOptions["protocol"] = append(ruleMatchOptions["protocol"].([]string),
							strings.TrimPrefix(itemTrimMatch, "protocol "))
					}
					// override (maxItem = 1)
					ruleOptions[matchWord] = []map[string]interface{}{ruleMatchOptions}
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" then source-nat "):
					itemTrimThen := strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" then source-nat ")
					ruleThenOptions := map[string]interface{}{
						"type": "",
						"pool": "",
					}
					if len(ruleOptions[thenWord].([]map[string]interface{})) > 0 {
						for k, v := range ruleOptions[thenWord].([]map[string]interface{})[0] {
							ruleThenOptions[k] = v
						}
					}
					if strings.HasPrefix(itemTrimThen, "pool ") {
						thenSplit := strings.Split(itemTrimThen, " ")
						ruleThenOptions["type"] = thenSplit[0]
						ruleThenOptions["pool"] = thenSplit[1]
					} else {
						ruleThenOptions["type"] = itemTrimThen
					}
					// override (maxItem = 1)
					ruleOptions[thenWord] = []map[string]interface{}{ruleThenOptions}
				}
				confRead.rule = append(confRead.rule, ruleOptions)
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}

func delSecurityNatSource(natSource string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat source rule-set "+natSource+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func fillSecurityNatSourceData(d *schema.ResourceData, natSourceOptions natSourceOptions) {
	tfErr := d.Set("name", natSourceOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("from", natSourceOptions.from)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("to", natSourceOptions.to)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("rule", natSourceOptions.rule)
	if tfErr != nil {
		panic(tfErr)
	}
}
