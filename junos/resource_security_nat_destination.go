package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type natDestinationOptions struct {
	name string
	from []map[string]interface{}
	rule []map[string]interface{}
}

func resourceSecurityNatDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityNatDestinationCreate,
		Read:   resourceSecurityNatDestinationRead,
		Update: resourceSecurityNatDestinationUpdate,
		Delete: resourceSecurityNatDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityNatDestinationImport,
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
						"destination_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNetworkFunc(),
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
											if !stringInSlice(value, []string{"off", "pool"}) {
												errors = append(errors, fmt.Errorf(
													"%q for %q is not 'off' or 'pool'", value, k))
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

func resourceSecurityNatDestinationCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security nat destination not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	securityNatDestinationExists, err := checkSecurityNatDestinationExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if securityNatDestinationExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("security nat destination %v already exists", d.Get("name").(string))
	}

	err = setSecurityNatDestination(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_security_nat_destination", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	securityNatDestinationExists, err = checkSecurityNatDestinationExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if securityNatDestinationExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security nat destination %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceSecurityNatDestinationRead(d, m)
}
func resourceSecurityNatDestinationRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	natDestinationOptions, err := readSecurityNatDestination(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if natDestinationOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatDestinationData(d, natDestinationOptions)
	}
	return nil
}
func resourceSecurityNatDestinationUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delSecurityNatDestination(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setSecurityNatDestination(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_security_nat_destination", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceSecurityNatDestinationRead(d, m)
}
func resourceSecurityNatDestinationDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delSecurityNatDestination(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_security_nat_destination", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourceSecurityNatDestinationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	securityNatDestinationExists, err := checkSecurityNatDestinationExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityNatDestinationExists {
		return nil, fmt.Errorf("don't find nat destination with id '%v' (id must be <name>)", d.Id())
	}
	natDestinationOptions, err := readSecurityNatDestination(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatDestinationData(d, natDestinationOptions)

	result[0] = d
	return result, nil
}

func checkSecurityNatDestinationExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	natDestinationConfig, err := sess.command("show configuration"+
		" security nat destination rule-set "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if natDestinationConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setSecurityNatDestination(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security nat destination rule-set " + d.Get("name").(string)
	for _, v := range d.Get("from").([]interface{}) {
		from := v.(map[string]interface{})
		for _, value := range from["value"].([]interface{}) {
			configSet = append(configSet, setPrefix+" from "+from["type"].(string)+" "+value.(string)+"\n")
		}
	}
	for _, v := range d.Get("rule").([]interface{}) {
		rule := v.(map[string]interface{})
		setPrefixRule := setPrefix + " rule " + rule["name"].(string)
		configSet = append(configSet, setPrefixRule+
			" match destination-address "+rule["destination_address"].(string)+"\n")
		for _, thenV := range rule[thenWord].([]interface{}) {
			then := thenV.(map[string]interface{})
			if then["type"].(string) == "off" {
				configSet = append(configSet, setPrefixRule+" then destination-nat off\n")
			}
			if then["type"].(string) == "pool" {
				if then["pool"].(string) == "" {
					return fmt.Errorf("missing pool for destination-nat pool for rule %v in %v",
						then["name"].(string), d.Get("name").(string))
				}
				configSet = append(configSet, setPrefixRule+" then destination-nat pool "+then["pool"].(string)+"\n")
			}
		}
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readSecurityNatDestination(natDestination string,
	m interface{}, jnprSess *NetconfObject) (natDestinationOptions, error) {
	sess := m.(*Session)
	var confRead natDestinationOptions

	natDestinationConfig, err := sess.command("show configuration"+
		" security nat destination rule-set "+natDestination+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if natDestinationConfig != emptyWord {
		confRead.name = natDestination
		for _, item := range strings.Split(natDestinationConfig, "\n") {
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
			case strings.HasPrefix(itemTrim, "rule "):
				ruleConfig := strings.Split(strings.TrimPrefix(itemTrim, "rule "), " ")

				ruleOptions := map[string]interface{}{
					"name":                ruleConfig[0],
					"destination_address": "",
					thenWord:              make([]map[string]interface{}, 0),
				}
				ruleOptions, confRead.rule = copyAndRemoveItemMapList("name", false, ruleOptions, confRead.rule)
				switch {
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" match destination-address "):
					ruleOptions["destination_address"] = strings.TrimPrefix(itemTrim,
						"rule "+ruleConfig[0]+" match destination-address ")
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" then destination-nat "):
					itemTrimThen := strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" then destination-nat ")
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

func delSecurityNatDestination(natDestination string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat destination rule-set "+natDestination+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func fillSecurityNatDestinationData(d *schema.ResourceData, natDestinationOptions natDestinationOptions) {
	tfErr := d.Set("name", natDestinationOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("from", natDestinationOptions.from)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("rule", natDestinationOptions.rule)
	if tfErr != nil {
		panic(tfErr)
	}
}
