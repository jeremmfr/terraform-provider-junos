package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type natSourcePoolOptions struct {
	portNoTranslation     bool
	portOverloadingFactor int
	name                  string
	routingInstance       string
	portRange             string
	address               []string
}

func resourceSecurityNatSourcePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityNatSourcePoolCreate,
		Read:   resourceSecurityNatSourcePoolRead,
		Update: resourceSecurityNatSourcePoolUpdate,
		Delete: resourceSecurityNatSourcePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityNatSourcePoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"address": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"routing_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"port_no_translation": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"port_overloading_factor", "port_range"},
			},
			"port_overloading_factor": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validateIntRange(2, 32),
				ConflictsWith: []string{"port_no_translation", "port_range"},
			},
			"port_range": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"port_overloading_factor", "port_no_translation"},
				ValidateFunc:  validateSourcePoolPortRange(),
			},
		},
	}
}

func resourceSecurityNatSourcePoolCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security nat source pool not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	securityNatSourcePoolExists, err := checkSecurityNatSourcePoolExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if securityNatSourcePoolExists {
		sess.configClear(jnprSess)

		return fmt.Errorf("security nat source pool %v already exists", d.Get("name").(string))
	}

	err = setSecurityNatSourcePool(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_security_nat_source_pool", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	securityNatSourcePoolExists, err = checkSecurityNatSourcePoolExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if securityNatSourcePoolExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security nat source pool %v not exists after commit => check your config", d.Get("name").(string))
	}

	return resourceSecurityNatSourcePoolRead(d, m)
}
func resourceSecurityNatSourcePoolRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	natSourcePoolOptions, err := readSecurityNatSourcePool(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if natSourcePoolOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityNatSourcePoolData(d, natSourcePoolOptions)
	}

	return nil
}
func resourceSecurityNatSourcePoolUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delSecurityNatSourcePool(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = setSecurityNatSourcePool(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_security_nat_source_pool", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceSecurityNatSourcePoolRead(d, m)
}
func resourceSecurityNatSourcePoolDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delSecurityNatSourcePool(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_security_nat_source_pool", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
}
func resourceSecurityNatSourcePoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	securityNatSourcePoolExists, err := checkSecurityNatSourcePoolExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityNatSourcePoolExists {
		return nil, fmt.Errorf("don't find nat source pool with id '%v' (id must be <name>)", d.Id())
	}
	natSourcePoolOptions, err := readSecurityNatSourcePool(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityNatSourcePoolData(d, natSourcePoolOptions)

	result[0] = d

	return result, nil
}

func checkSecurityNatSourcePoolExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	natSourcePoolConfig, err := sess.command("show configuration"+
		" security nat source pool "+name+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if natSourcePoolConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSecurityNatSourcePool(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security nat source pool " + d.Get("name").(string)
	for _, v := range d.Get("address").([]interface{}) {
		err := validateIPwithMask(v.(string))
		if err != nil {
			return err
		}
		configSet = append(configSet, setPrefix+" address "+v.(string)+"\n")
	}

	if d.Get("port_no_translation").(bool) {
		configSet = append(configSet, setPrefix+" port no-translation "+"\n")
	}
	if d.Get("port_overloading_factor").(int) != 0 {
		configSet = append(configSet, setPrefix+" port port-overloading-factor "+
			strconv.Itoa(d.Get("port_overloading_factor").(int))+"\n")
	}
	if d.Get("port_range").(string) != "" {
		rangePort := strings.Split(d.Get("port_range").(string), "-")
		configSet = append(configSet, setPrefix+" port range "+rangePort[0]+" to "+rangePort[1]+"\n")
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+d.Get("routing_instance").(string))
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func readSecurityNatSourcePool(natSourcePool string,
	m interface{}, jnprSess *NetconfObject) (natSourcePoolOptions, error) {
	sess := m.(*Session)
	var confRead natSourcePoolOptions

	natSourcePoolConfig, err := sess.command("show configuration"+
		" security nat source pool "+natSourcePool+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if natSourcePoolConfig != emptyWord {
		confRead.name = natSourcePool
		var portRange string
		for _, item := range strings.Split(natSourcePoolConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "address "):
				confRead.address = append(confRead.address, strings.TrimPrefix(itemTrim, "address "))
			case strings.HasPrefix(itemTrim, "port no-translation"):
				confRead.portNoTranslation = true
			case strings.HasPrefix(itemTrim, "port port-overloading-factor"):
				confRead.portOverloadingFactor, err = strconv.Atoi(strings.TrimPrefix(itemTrim,
					"port port-overloading-factor "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "port range to"):
				portRange += "-" + strings.TrimPrefix(itemTrim, "port range to ")
			case strings.HasPrefix(itemTrim, "port range "):
				portRange = strings.TrimPrefix(itemTrim, "port range ")
			case strings.HasPrefix(itemTrim, "routing-instance"):
				confRead.routingInstance = strings.TrimPrefix(itemTrim, "routing-instance ")
			}
		}
		confRead.portRange = portRange
	} else {
		confRead.name = ""

		return confRead, nil
	}

	return confRead, nil
}

func delSecurityNatSourcePool(natSourcePool string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security nat source pool "+natSourcePool+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func fillSecurityNatSourcePoolData(d *schema.ResourceData, natSourcePoolOptions natSourcePoolOptions) {
	tfErr := d.Set("name", natSourcePoolOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("address", natSourcePoolOptions.address)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("port_no_translation", natSourcePoolOptions.portNoTranslation)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("port_overloading_factor", natSourcePoolOptions.portOverloadingFactor)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("port_range", natSourcePoolOptions.portRange)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", natSourcePoolOptions.routingInstance)
	if tfErr != nil {
		panic(tfErr)
	}
}

func validateSourcePoolPortRange() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v := i.(string)
		vSplit := strings.Split(v, "-")
		if len(vSplit) < 2 {
			es = append(es, fmt.Errorf(
				"%q missing range separtor - in %q", k, i))
		}
		low, err := strconv.Atoi(vSplit[0])
		if err != nil {
			es = append(es, err)
		}
		high, err := strconv.Atoi(vSplit[1])
		if err != nil {
			es = append(es, err)
		}
		if low > high {
			es = append(es, fmt.Errorf(
				"%q low in %q bigger than high", k, i))
		}
		if low < 1024 {
			es = append(es, fmt.Errorf(
				"%q low in %q is too small (min 1024)", k, i))
		}
		if high > 63487 {
			es = append(es, fmt.Errorf(
				"%q high in %q is too big (max 63487)", k, i))
		}

		return
	}
}
