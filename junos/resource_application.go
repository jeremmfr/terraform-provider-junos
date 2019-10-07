package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type applicationOptions struct {
	name            string
	protocol        string
	destinationPort string
	sourcePort      string
}

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceApplicationImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, m interface{}) error {
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
	appExists, err := checkApplicationExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if appExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("application %v already exists", d.Get("name").(string))
	}
	err = setApplication(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	appExists, err = checkApplicationExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if appExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("application %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceApplicationRead(d, m)
}
func resourceApplicationRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	applicationOptions, err := readApplication(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if applicationOptions.name == "" {
		d.SetId("")
	} else {
		fillApplicationData(d, applicationOptions)
	}
	return nil
}
func resourceApplicationUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delApplication(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setApplication(d, m, jnprSess)
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
	return resourceApplicationRead(d, m)
}
func resourceApplicationDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delApplication(d, m, jnprSess)
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
func resourceApplicationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	appExists, err := checkApplicationExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !appExists {
		return nil, fmt.Errorf("don't find application with id '%v' (id must be <name>)", d.Id())
	}
	applicationOptions, err := readApplication(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillApplicationData(d, applicationOptions)
	result[0] = d
	return result, nil
}

func checkApplicationExists(application string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	applicationConfig, err := sess.command("show configuration applications application "+
		application+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if applicationConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setApplication(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set applications application " + d.Get("name").(string)
	if d.Get("protocol").(string) != "" {
		configSet = append(configSet, setPrefix+" protocol "+d.Get("protocol").(string)+"\n")
	}
	if d.Get("destination_port").(string) != "" {
		configSet = append(configSet, setPrefix+" destination-port "+d.Get("destination_port").(string)+"\n")
	}
	if d.Get("source_port").(string) != "" {
		configSet = append(configSet, setPrefix+" source-port "+d.Get("source_port").(string)+"\n")
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readApplication(application string, m interface{}, jnprSess *NetconfObject) (applicationOptions, error) {
	sess := m.(*Session)
	var confRead applicationOptions

	applicationConfig, err := sess.command("show configuration applications application "+
		application+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if applicationConfig != emptyWord {
		confRead.name = application
		for _, item := range strings.Split(applicationConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "protocol "):
				confRead.protocol = strings.TrimPrefix(itemTrim, "protocol ")
			case strings.HasPrefix(itemTrim, "destination-port "):
				confRead.destinationPort = strings.TrimPrefix(itemTrim, "destination-port ")
			case strings.HasPrefix(itemTrim, "source-port "):
				confRead.sourcePort = strings.TrimPrefix(itemTrim, "source-port ")
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}
func delApplication(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete applications application "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillApplicationData(d *schema.ResourceData, applicationOptions applicationOptions) {
	tfErr := d.Set("name", applicationOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("protocol", applicationOptions.protocol)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("destination_port", applicationOptions.destinationPort)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("source_port", applicationOptions.sourcePort)
	if tfErr != nil {
		panic(tfErr)
	}
}
