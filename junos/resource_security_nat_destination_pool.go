package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type natDestinationPoolOptions struct {
	addressPort     int
	name            string
	address         string
	addressTo       string
	routingInstance string
}

func resourceSecurityNatDestinationPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityNatDestinationPoolCreate,
		Read:   resourceSecurityNatDestinationPoolRead,
		Update: resourceSecurityNatDestinationPoolUpdate,
		Delete: resourceSecurityNatDestinationPoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityNatDestinationPoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIPMaskFunc(),
			},
			"address_to": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateIPMaskFunc(),
				ConflictsWith: []string{"address_port"},
			},
			"address_port": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validateIntRange(1, 65535),
				ConflictsWith: []string{"address_to"},
			},
			"routing_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
		},
	}
}

func resourceSecurityNatDestinationPoolCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security nat destination pool not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	securityNatDestinationPoolExists, err := checkSecurityNatDestinationPoolExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if securityNatDestinationPoolExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("security nat destination pool %v already exists", d.Get("name").(string))
	}

	err = setSecurityNatDestinationPool(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_security_nat_destination_pool", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	securityNatDestinationPoolExists, err = checkSecurityNatDestinationPoolExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if securityNatDestinationPoolExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security nat destination pool %v not exists after commit "+
			"=> check your config", d.Get("name").(string))
	}
	return resourceSecurityNatDestinationPoolRead(d, m)
}
func resourceSecurityNatDestinationPoolRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	natDestinationPoolOptions, err := readSecurityNatDestinationPool(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if natDestinationPoolOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatDestinationPoolData(d, natDestinationPoolOptions)
	}
	return nil
}
func resourceSecurityNatDestinationPoolUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delSecurityNatDestinationPool(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setSecurityNatDestinationPool(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_security_nat_destination_pool", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceSecurityNatDestinationPoolRead(d, m)
}
func resourceSecurityNatDestinationPoolDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delSecurityNatDestinationPool(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_security_nat_destination_pool", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourceSecurityNatDestinationPoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	securityNatDestinationPoolExists, err := checkSecurityNatDestinationPoolExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityNatDestinationPoolExists {
		return nil, fmt.Errorf("don't find nat destination pool with id '%v' (id must be <name>)", d.Id())
	}
	natDestinationPoolOptions, err := readSecurityNatDestinationPool(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatDestinationPoolData(d, natDestinationPoolOptions)

	result[0] = d
	return result, nil
}

func checkSecurityNatDestinationPoolExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	natDestinationPoolConfig, err := sess.command("show configuration"+
		" security nat destination pool "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if natDestinationPoolConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setSecurityNatDestinationPool(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security nat destination pool " + d.Get("name").(string)
	configSet = append(configSet, setPrefix+" address "+d.Get("address").(string)+"\n")
	if d.Get("address_to").(string) != "" {
		configSet = append(configSet, setPrefix+" address to "+d.Get("address_to").(string)+"\n")
	}
	if d.Get("address_port").(int) != 0 {
		configSet = append(configSet, setPrefix+" address port "+strconv.Itoa(d.Get("address_port").(int))+"\n")
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+d.Get("routing_instance").(string)+"\n")
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readSecurityNatDestinationPool(natDestinationPool string,
	m interface{}, jnprSess *NetconfObject) (natDestinationPoolOptions, error) {
	sess := m.(*Session)
	var confRead natDestinationPoolOptions

	natDestinationPoolConfig, err := sess.command("show configuration"+
		" security nat destination pool "+natDestinationPool+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if natDestinationPoolConfig != emptyWord {
		confRead.name = natDestinationPool
		for _, item := range strings.Split(natDestinationPoolConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "address to"):
				confRead.addressTo = strings.TrimPrefix(itemTrim, "address to ")
			case strings.HasPrefix(itemTrim, "address port"):
				confRead.addressPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "address port "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "address "):
				confRead.address = strings.TrimPrefix(itemTrim, "address ")
			case strings.HasPrefix(itemTrim, "routing-instance "):
				confRead.routingInstance = strings.TrimPrefix(itemTrim, "routing-instance ")
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}

func delSecurityNatDestinationPool(natDestinationPool string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat destination pool "+natDestinationPool+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func fillSecurityNatDestinationPoolData(d *schema.ResourceData, natDestinationPoolOptions natDestinationPoolOptions) {
	tfErr := d.Set("name", natDestinationPoolOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("address", natDestinationPoolOptions.address)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("address_to", natDestinationPoolOptions.addressTo)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("address_port", natDestinationPoolOptions.addressPort)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", natDestinationPoolOptions.routingInstance)
	if tfErr != nil {
		panic(tfErr)
	}
}
