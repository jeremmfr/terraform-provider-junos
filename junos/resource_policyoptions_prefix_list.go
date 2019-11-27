package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type prefixListOptions struct {
	name   string
	prefix []string
}

func resourcePolicyoptionsPrefixList() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyoptionsPrefixListCreate,
		Read:   resourcePolicyoptionsPrefixListRead,
		Update: resourcePolicyoptionsPrefixListUpdate,
		Delete: resourcePolicyoptionsPrefixListDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePolicyoptionsPrefixListImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"prefix": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePolicyoptionsPrefixListCreate(d *schema.ResourceData, m interface{}) error {
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
	policyoptsPrefixListExists, err := checkPolicyoptionsPrefixListExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if policyoptsPrefixListExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("policy-options prefix-list %v already exists", d.Get("name").(string))
	}

	err = setPolicyoptionsPrefixList(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_policyoptions_prefix_list", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	policyoptsPrefixListExists, err = checkPolicyoptionsPrefixListExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if policyoptsPrefixListExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("policy-options prefix-list %v not exists after commit "+
			"=> check your config", d.Get("name").(string))
	}
	return resourcePolicyoptionsPrefixListRead(d, m)
}
func resourcePolicyoptionsPrefixListRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	prefixListOptions, err := readPolicyoptionsPrefixList(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if prefixListOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsPrefixListData(d, prefixListOptions)
	}
	return nil
}
func resourcePolicyoptionsPrefixListUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delPolicyoptionsPrefixList(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setPolicyoptionsPrefixList(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_policyoptions_prefix_list", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourcePolicyoptionsPrefixListRead(d, m)
}
func resourcePolicyoptionsPrefixListDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delPolicyoptionsPrefixList(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_policyoptions_prefix_list", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourcePolicyoptionsPrefixListImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsPrefixListExists, err := checkPolicyoptionsPrefixListExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsPrefixListExists {
		return nil, fmt.Errorf("don't find policy-options prefix-list with id '%v' (id must be <name>)", d.Id())
	}
	prefixListOptions, err := readPolicyoptionsPrefixList(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsPrefixListData(d, prefixListOptions)

	result[0] = d
	return result, nil
}

func checkPolicyoptionsPrefixListExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	prefixListConfig, err := sess.command("show configuration policy-options prefix-list "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if prefixListConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setPolicyoptionsPrefixList(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	for _, v := range d.Get("prefix").([]interface{}) {
		err := validateNetwork(v.(string))
		if err != nil {
			return err
		}
		configSet = append(configSet, "set policy-options prefix-list "+d.Get("name").(string)+" "+v.(string)+"\n")
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readPolicyoptionsPrefixList(prefixList string, m interface{}, jnprSess *NetconfObject) (prefixListOptions, error) {
	sess := m.(*Session)
	var confRead prefixListOptions

	prefixListConfig, err := sess.command("show configuration policy-options prefix-list "+
		prefixList+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if prefixListConfig != emptyWord {
		confRead.name = prefixList
		for _, item := range strings.Split(prefixListConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			if strings.Contains(item, "/") {
				confRead.prefix = append(confRead.prefix, strings.TrimPrefix(item, setLineStart))
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}

func delPolicyoptionsPrefixList(prefixList string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options prefix-list "+prefixList+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func fillPolicyoptionsPrefixListData(d *schema.ResourceData, prefixListOptions prefixListOptions) {
	tfErr := d.Set("name", prefixListOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("prefix", prefixListOptions.prefix)
	if tfErr != nil {
		panic(tfErr)
	}
}
