package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type asPathGroupOptions struct {
	dynamicDb bool
	name      string
	asPath    []map[string]interface{}
}

func resourcePolicyoptionsAsPathGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyoptionsAsPathGroupCreate,
		Read:   resourcePolicyoptionsAsPathGroupRead,
		Update: resourcePolicyoptionsAsPathGroupUpdate,
		Delete: resourcePolicyoptionsAsPathGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePolicyoptionsAsPathGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"as_path": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dynamic_db": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePolicyoptionsAsPathGroupCreate(d *schema.ResourceData, m interface{}) error {
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
	policyoptsAsPathGroupExists, err := checkPolicyoptionsAsPathGroupExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if policyoptsAsPathGroupExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("policy-options as-path-group %v already exists", d.Get("name").(string))
	}

	err = setPolicyoptionsAsPathGroup(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_policyoptions_as_path_group", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	policyoptsAsPathGroupExists, err = checkPolicyoptionsAsPathGroupExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if policyoptsAsPathGroupExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("policy-options as-path-group %v not exists after commit "+
			"=> check your config", d.Get("name").(string))
	}
	return resourcePolicyoptionsAsPathGroupRead(d, m)
}
func resourcePolicyoptionsAsPathGroupRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	asPathGroupOptions, err := readPolicyoptionsAsPathGroup(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if asPathGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsAsPathGroupData(d, asPathGroupOptions)
	}
	return nil
}
func resourcePolicyoptionsAsPathGroupUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delPolicyoptionsAsPathGroup(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setPolicyoptionsAsPathGroup(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_policyoptions_as_path_group", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourcePolicyoptionsAsPathGroupRead(d, m)
}
func resourcePolicyoptionsAsPathGroupDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delPolicyoptionsAsPathGroup(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_policyoptions_as_path_group", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourcePolicyoptionsAsPathGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsAsPathGroupExists, err := checkPolicyoptionsAsPathGroupExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsAsPathGroupExists {
		return nil, fmt.Errorf("don't find policy-options as-path-group with id '%v' (id must be <name>)", d.Id())
	}
	asPathGroupOptions, err := readPolicyoptionsAsPathGroup(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsAsPathGroupData(d, asPathGroupOptions)

	result[0] = d
	return result, nil
}

func checkPolicyoptionsAsPathGroupExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	asPathGroupConfig, err := sess.command("show configuration "+
		"policy-options as-path-group "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if asPathGroupConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setPolicyoptionsAsPathGroup(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set policy-options as-path-group " + d.Get("name").(string)
	for _, v := range d.Get("as_path").([]interface{}) {
		asPath := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+
			" as-path "+asPath["name"].(string)+
			" \""+asPath["path"].(string)+"\"\n")
	}
	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, setPrefix+" dynamic-db\n")
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readPolicyoptionsAsPathGroup(asPathGroup string,
	m interface{}, jnprSess *NetconfObject) (asPathGroupOptions, error) {
	sess := m.(*Session)
	var confRead asPathGroupOptions

	asPathGroupConfig, err := sess.command("show configuration"+
		" policy-options as-path-group "+asPathGroup+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if asPathGroupConfig != emptyWord {
		confRead.name = asPathGroup
		for _, item := range strings.Split(asPathGroupConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "dynamic-db"):
				confRead.dynamicDb = true
			case strings.HasPrefix(itemTrim, "as-path "):
				asPath := map[string]interface{}{
					"name": "",
					"path": "",
				}
				itemSplit := strings.Split(strings.TrimPrefix(itemTrim, "as-path "), " ")
				asPath["name"] = itemSplit[0]
				asPath["path"] = strings.Trim(strings.TrimPrefix(itemTrim,
					"as-path "+asPath["name"].(string)+" "), "\"")
				confRead.asPath = append(confRead.asPath, asPath)
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}

func delPolicyoptionsAsPathGroup(asPathGroup string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options as-path-group "+asPathGroup+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func fillPolicyoptionsAsPathGroupData(d *schema.ResourceData, asPathGroupOptions asPathGroupOptions) {
	tfErr := d.Set("name", asPathGroupOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("as_path", asPathGroupOptions.asPath)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("dynamic_db", asPathGroupOptions.dynamicDb)
	if tfErr != nil {
		panic(tfErr)
	}
}
