package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type natStaticOptions struct {
	name string
	from []map[string]interface{}
	rule []map[string]interface{}
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
				ValidateDiagFunc: validateNameObjectJunos([]string{}),
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
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}),
						},
						"destination_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
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
										ValidateFunc: validation.StringInSlice([]string{inetWord, prefixWord}, false),
									},
									"routing_instance": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validateNameObjectJunos([]string{}),
									},
									"prefix": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsCIDRNetwork(0, 128),
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

func resourceSecurityNatStaticCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security nat static not compatible with Junos device %s", jnprSess.Platform[0].Model))
	}
	sess.configLock(jnprSess)
	securityNatStaticExists, err := checkSecurityNatStaticExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if securityNatStaticExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security nat static %v already exists", d.Get("name").(string)))
	}

	err = setSecurityNatStatic(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	err = sess.commitConf("create resource junos_security_nat_static", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	securityNatStaticExists, err = checkSecurityNatStaticExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return diag.FromErr(err)
	}
	if securityNatStaticExists {
		d.SetId(d.Get("name").(string))
	} else {
		return diag.FromErr(fmt.Errorf("security nat static %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))
	}

	return resourceSecurityNatStaticRead(ctx, d, m)
}
func resourceSecurityNatStaticRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
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
	err = delSecurityNatStatic(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	err = setSecurityNatStatic(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	err = sess.commitConf("update resource junos_security_nat_static", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	d.Partial(false)

	return resourceSecurityNatStaticRead(ctx, d, m)
}
func resourceSecurityNatStaticDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	err = delSecurityNatStatic(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	err = sess.commitConf("delete resource junos_security_nat_static", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	return nil
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

	setPrefix := "set security nat static rule-set " + d.Get("name").(string)
	for _, v := range d.Get("from").([]interface{}) {
		from := v.(map[string]interface{})
		for _, value := range from["value"].([]interface{}) {
			configSet = append(configSet, setPrefix+" from "+from["type"].(string)+" "+value.(string))
		}
	}
	for _, v := range d.Get("rule").([]interface{}) {
		rule := v.(map[string]interface{})
		setPrefixRule := setPrefix + " rule " + rule["name"].(string)
		configSet = append(configSet, setPrefixRule+" match destination-address "+
			rule["destination_address"].(string))
		for _, thenV := range rule[thenWord].([]interface{}) {
			then := thenV.(map[string]interface{})
			if then["type"].(string) == inetWord {
				if then["routing_instance"].(string) == "" {
					return fmt.Errorf("missing routing_instance for static-nat inet for rule %v in %v ",
						rule["name"].(string), d.Get("name").(string))
				}
				configSet = append(configSet, setPrefixRule+" then static-nat inet routing-instance "+
					then["routing_instance"].(string))
			}
			if then["type"].(string) == prefixWord {
				if then[prefixWord].(string) == "" {
					return fmt.Errorf("missing prefix for static-nat prefix for rule %v in %v",
						rule["name"].(string), d.Get("name").(string))
				}
				configSet = append(configSet, setPrefixRule+" then static-nat prefix "+then[prefixWord].(string))
				if then["routing_instance"].(string) != "" {
					configSet = append(configSet, setPrefixRule+" then static-nat prefix routing-instance "+
						then["routing_instance"].(string))
				}
			}
		}
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
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
				case strings.HasPrefix(itemTrim, "rule "+ruleConfig[0]+" then static-nat "):
					itemThen := strings.TrimPrefix(itemTrim, "rule "+ruleConfig[0]+" then static-nat ")
					ruleThenOptions := map[string]interface{}{
						"type":             "",
						"routing_instance": "",
						prefixWord:         "",
					}
					if len(ruleOptions[thenWord].([]map[string]interface{})) > 0 {
						for k, v := range ruleOptions[thenWord].([]map[string]interface{})[0] {
							ruleThenOptions[k] = v
						}
					}
					switch {
					case strings.HasPrefix(itemThen, "prefix "):
						ruleThenOptions["type"] = prefixWord
						if strings.HasPrefix(itemThen, "prefix routing-instance ") {
							ruleThenOptions["routing_instance"] = strings.TrimPrefix(itemThen, "prefix routing-instance ")
						} else {
							ruleThenOptions[prefixWord] = strings.TrimPrefix(itemThen, "prefix ")
						}
					case strings.HasPrefix(itemThen, "inet "):
						ruleThenOptions["type"] = inetWord
						ruleThenOptions["routing_instance"] = strings.TrimPrefix(itemThen, "inet routing-instance ")
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

func delSecurityNatStatic(natStatic string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat static rule-set "+natStatic)
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func fillSecurityNatStaticData(d *schema.ResourceData, natStaticOptions natStaticOptions) {
	tfErr := d.Set("name", natStaticOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("from", natStaticOptions.from)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("rule", natStaticOptions.rule)
	if tfErr != nil {
		panic(tfErr)
	}
}
