package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type instanceOptions struct {
	name         string
	instanceType string
	as           string
}

func resourceRoutingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoutingInstanceCreate,
		Read:   resourceRoutingInstanceRead,
		Update: resourceRoutingInstanceUpdate,
		Delete: resourceRoutingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRoutingInstanceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "virtual-router",
			},
			"as": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRoutingInstanceCreate(d *schema.ResourceData, m interface{}) error {
	if d.Get("name").(string) == "default" {
		return fmt.Errorf("name default isn't valid")
	}
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	routingInstanceExists, err := checkRoutingInstanceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if routingInstanceExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("routing-instance %v already exists", d.Get("name").(string))
	}
	err = setRoutingInstance(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	routingInstanceExists, err = checkRoutingInstanceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if routingInstanceExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("routing-instance %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceRoutingInstanceRead(d, m)
}
func resourceRoutingInstanceRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	instanceOptions, err := readRoutingInstance(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if instanceOptions.name == "" {
		d.SetId("")
	} else {
		fillRoutingInstanceData(d, instanceOptions)
	}
	return nil
}
func resourceRoutingInstanceUpdate(d *schema.ResourceData, m interface{}) error {
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

	err = delRoutingInstanceOpts(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setRoutingInstance(d, m, jnprSess)
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
	return resourceRoutingInstanceRead(d, m)
}
func resourceRoutingInstanceDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delRoutingInstance(d, m, jnprSess)
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
func resourceRoutingInstanceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	routingInstanceExists, err := checkRoutingInstanceExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !routingInstanceExists {
		return nil, fmt.Errorf("don't find routing instance with id '%v' (id must be <name>)", d.Id())
	}
	instanceOptions, err := readRoutingInstance(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillRoutingInstanceData(d, instanceOptions)
	result[0] = d
	return result, nil
}

func checkRoutingInstanceExists(instance string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	routingInstanceConfig, err := sess.command("show configuration routing-instances "+instance+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if routingInstanceConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setRoutingInstance(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set routing-instances " + d.Get("name").(string) + " "
	if d.Get("type").(string) != "" {
		configSet = append(configSet, setPrefix+"instance-type "+d.Get("type").(string)+"\n")
	} else {
		configSet = append(configSet, setPrefix+"\n")
	}
	if d.Get("as").(string) != "" {
		configSet = append(configSet, setPrefix+
			"routing-options autonomous-system "+d.Get("as").(string)+"\n")
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readRoutingInstance(instance string, m interface{}, jnprSess *NetconfObject) (instanceOptions, error) {
	sess := m.(*Session)
	var confRead instanceOptions

	instanceConfig, err := sess.command("show configuration"+
		" routing-instances "+instance+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if instanceConfig != emptyWord {
		confRead.name = instance
		for _, item := range strings.Split(instanceConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, "set ")
			switch {
			case strings.HasPrefix(itemTrim, "instance-type "):
				confRead.instanceType = strings.TrimPrefix(itemTrim, "instance-type ")
			case strings.HasPrefix(itemTrim, "routing-options autonomous-system "):
				confRead.as = strings.TrimPrefix(itemTrim, "routing-options autonomous-system ")
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}
func delRoutingInstanceOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := "delete routing-instances " + d.Get("name").(string) + " "
	configSet = append(configSet,
		setPrefix+"instance-type\n",
		setPrefix+"routing-options autonomous-system\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func delRoutingInstance(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-instances "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillRoutingInstanceData(d *schema.ResourceData, instanceOptions instanceOptions) {
	tfErr := d.Set("name", instanceOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("type", instanceOptions.instanceType)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("as", instanceOptions.as)
	if tfErr != nil {
		panic(tfErr)
	}
}
