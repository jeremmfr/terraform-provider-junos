package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type applicationSetOptions struct {
	name         string
	applications []string
}

func resourceApplicationSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationSetCreate,
		Read:   resourceApplicationSetRead,
		Update: resourceApplicationSetUpdate,
		Delete: resourceApplicationSetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceApplicationSetImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"applications": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceApplicationSetCreate(d *schema.ResourceData, m interface{}) error {
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
	appSetExists, err := checkApplicationSetExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if appSetExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("application-set %v already exists", d.Get("name").(string))
	}
	err = setApplicationSet(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	appSetExists, err = checkApplicationSetExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if appSetExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("application-set %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceApplicationSetRead(d, m)
}
func resourceApplicationSetRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	applicationSetOptions, err := readApplicationSet(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if applicationSetOptions.name == "" {
		d.SetId("")
	} else {
		fillApplicationSetData(d, applicationSetOptions)
	}
	return nil
}
func resourceApplicationSetUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delApplicationSet(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setApplicationSet(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceApplicationSetRead(d, m)

}
func resourceApplicationSetDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delApplicationSet(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourceApplicationSetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	appSetExists, err := checkApplicationSetExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !appSetExists {
		return nil, fmt.Errorf("don't find application-set with id '%v' (id must be <name>)", d.Id())
	}
	applicationSetOptions, err := readApplicationSet(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillApplicationSetData(d, applicationSetOptions)
	result[0] = d
	return result, nil
}

func checkApplicationSetExists(applicationSet string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	applicationSetConfig, err := sess.command("show configuration applications application-set "+
		applicationSet+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if applicationSetConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setApplicationSet(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set applications application-set " + d.Get("name").(string)
	for _, v := range d.Get("applications").([]interface{}) {
		configSet = append(configSet, setPrefix+" application "+v.(string)+"\n")
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readApplicationSet(applicationSet string, m interface{}, jnprSess *NetconfObject) (applicationSetOptions, error) {
	sess := m.(*Session)
	var confRead applicationSetOptions

	applicationSetConfig, err := sess.command("show configuration applications application-set "+
		applicationSet+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if applicationSetConfig != emptyWord {
		confRead.name = applicationSet
		for _, item := range strings.Split(applicationSetConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			if strings.HasPrefix(itemTrim, "application ") {
				confRead.applications = append(confRead.applications, strings.TrimPrefix(itemTrim, "application "))
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}
func delApplicationSet(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete applications application-set "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillApplicationSetData(d *schema.ResourceData, applicationSetOptions applicationSetOptions) {
	tfErr := d.Set("name", applicationSetOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("applications", applicationSetOptions.applications)
	if tfErr != nil {
		panic(tfErr)
	}
}
