package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ntpServerOptions struct {
	prefer          bool
	key             int
	version         int
	address         string
	routingInstance string
}

func resourceSystemNtpServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSystemNtpServerCreate,
		Read:   resourceSystemNtpServerRead,
		Update: resourceSystemNtpServerUpdate,
		Delete: resourceSystemNtpServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemNtpServerImport,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIPFunc(),
			},
			"key": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 65534),
			},
			"prefer": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"routing_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 4),
			},
		},
	}
}

func resourceSystemNtpServerCreate(d *schema.ResourceData, m interface{}) error {
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
	ntpServerExists, err := checkSystemNtpServerExists(d.Get("address").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if ntpServerExists {
		sess.configClear(jnprSess)

		return fmt.Errorf("system ntp server %v already exists", d.Get("address").(string))
	}

	err = setSystemNtpServer(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_system_ntp_server", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	ntpServerExists, err = checkSystemNtpServerExists(d.Get("address").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if ntpServerExists {
		d.SetId(d.Get("address").(string))
	} else {
		return fmt.Errorf("system ntp server %v not exists after commit => check your config", d.Get("address").(string))
	}

	return resourceSystemNtpServerRead(d, m)
}
func resourceSystemNtpServerRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	ntpServerOptions, err := readSystemNtpServer(d.Get("address").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if ntpServerOptions.address == "" {
		d.SetId("")
	} else {
		fillSystemNtpServerData(d, ntpServerOptions)
	}

	return nil
}
func resourceSystemNtpServerUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delSystemNtpServer(d.Get("address").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = setSystemNtpServer(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_system_ntp_server", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceSystemNtpServerRead(d, m)
}
func resourceSystemNtpServerDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delSystemNtpServer(d.Get("address").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_system_ntp_server", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
}
func resourceSystemNtpServerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	ntpServerExists, err := checkSystemNtpServerExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ntpServerExists {
		return nil, fmt.Errorf("don't find system ntp server with id '%v' (id must be <address>)", d.Id())
	}
	ntpServerOptions, err := readSystemNtpServer(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystemNtpServerData(d, ntpServerOptions)

	result[0] = d

	return result, nil
}

func checkSystemNtpServerExists(address string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	ntpServerConfig, err := sess.command("show configuration"+
		" system ntp server "+address+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if ntpServerConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSystemNtpServer(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set system ntp server " + d.Get("address").(string)
	configSet := []string{setPrefix}

	if d.Get("key").(int) != 0 {
		configSet = append(configSet, setPrefix+" key "+strconv.Itoa(d.Get("key").(int)))
	}
	if d.Get("version").(int) != 0 {
		configSet = append(configSet, setPrefix+" version "+strconv.Itoa(d.Get("version").(int)))
	}
	if d.Get("prefer").(bool) {
		configSet = append(configSet, setPrefix+" prefer")
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
func readSystemNtpServer(address string, m interface{}, jnprSess *NetconfObject) (ntpServerOptions, error) {
	sess := m.(*Session)
	var confRead ntpServerOptions

	ntpServerConfig, err := sess.command("show configuration"+
		" system ntp server "+address+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if ntpServerConfig != emptyWord {
		confRead.address = address
		for _, item := range strings.Split(ntpServerConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "key "):
				var err error
				confRead.key, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "key "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "version "):
				var err error
				confRead.version, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "version "))
				if err != nil {
					return confRead, err
				}
			case strings.HasSuffix(itemTrim, "prefer"):
				confRead.prefer = true
			case strings.HasPrefix(itemTrim, "routing-instance "):
				confRead.routingInstance = strings.TrimPrefix(itemTrim, "routing-instance ")
			}
		}
	} else {
		confRead.address = ""

		return confRead, nil
	}

	return confRead, nil
}

func delSystemNtpServer(address string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system ntp server "+address)
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func fillSystemNtpServerData(d *schema.ResourceData, ntpServerOptions ntpServerOptions) {
	tfErr := d.Set("address", ntpServerOptions.address)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("key", ntpServerOptions.key)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("version", ntpServerOptions.version)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("prefer", ntpServerOptions.prefer)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", ntpServerOptions.routingInstance)
	if tfErr != nil {
		panic(tfErr)
	}
}
