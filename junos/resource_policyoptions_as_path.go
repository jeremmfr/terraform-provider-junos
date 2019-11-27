package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type asPathOptions struct {
	dynamicDb bool
	name      string
	path      string
}

func resourcePolicyoptionsAsPath() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyoptionsAsPathCreate,
		Read:   resourcePolicyoptionsAsPathRead,
		Update: resourcePolicyoptionsAsPathUpdate,
		Delete: resourcePolicyoptionsAsPathDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePolicyoptionsAsPathImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dynamic_db": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePolicyoptionsAsPathCreate(d *schema.ResourceData, m interface{}) error {
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
	policyoptsAsPathExists, err := checkPolicyoptionsAsPathExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if policyoptsAsPathExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("policy-options as-path %v already exists", d.Get("name").(string))
	}

	err = setPolicyoptionsAsPath(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_policyoptions_as_path", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	policyoptsAsPathExists, err = checkPolicyoptionsAsPathExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if policyoptsAsPathExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("policy-options as-path %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourcePolicyoptionsAsPathRead(d, m)
}
func resourcePolicyoptionsAsPathRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	asPathOptions, err := readPolicyoptionsAsPath(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if asPathOptions.name == "" {
		d.SetId("")
	} else {
		fillPolicyoptionsAsPathData(d, asPathOptions)
	}
	return nil
}
func resourcePolicyoptionsAsPathUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delPolicyoptionsAsPath(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setPolicyoptionsAsPath(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_policyoptions_as_path", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourcePolicyoptionsAsPathRead(d, m)
}
func resourcePolicyoptionsAsPathDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delPolicyoptionsAsPath(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_policyoptions_as_path", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourcePolicyoptionsAsPathImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	policyoptsAsPathExists, err := checkPolicyoptionsAsPathExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !policyoptsAsPathExists {
		return nil, fmt.Errorf("don't find policy-options as-path with id '%v' (id must be <name>)", d.Id())
	}
	asPathOptions, err := readPolicyoptionsAsPath(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillPolicyoptionsAsPathData(d, asPathOptions)

	result[0] = d
	return result, nil
}

func checkPolicyoptionsAsPathExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	asPathConfig, err := sess.command("show configuration policy-options as-path "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if asPathConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setPolicyoptionsAsPath(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	if d.Get("path").(string) != "" {
		configSet = append(configSet, "set policy-options as-path "+d.Get("name").(string)+
			" \""+d.Get("path").(string)+"\"\n")
	}
	if d.Get("dynamic_db").(bool) {
		configSet = append(configSet, "set policy-options as-path "+d.Get("name").(string)+
			" dynamic-db\n")
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readPolicyoptionsAsPath(asPath string, m interface{}, jnprSess *NetconfObject) (asPathOptions, error) {
	sess := m.(*Session)
	var confRead asPathOptions

	asPathConfig, err := sess.command("show configuration "+
		"policy-options as-path "+asPath+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if asPathConfig != emptyWord {
		confRead.name = asPath
		for _, item := range strings.Split(asPathConfig, "\n") {
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
			default:
				confRead.path = strings.Trim(itemTrim, "\"")
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}

func delPolicyoptionsAsPath(asPath string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete policy-options as-path "+asPath+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func fillPolicyoptionsAsPathData(d *schema.ResourceData, asPathOptions asPathOptions) {
	tfErr := d.Set("name", asPathOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("path", asPathOptions.path)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("dynamic_db", asPathOptions.dynamicDb)
	if tfErr != nil {
		panic(tfErr)
	}
}
