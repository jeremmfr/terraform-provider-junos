package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	jdecode "github.com/jeremmfr/junosdecode"
)

type radiusServerOptions struct {
	port                    int
	accountingPort          int
	dynamicRequestPort      int
	preauthenticationPort   int
	retry                   int
	accountingRetry         int
	timeout                 int
	accoutingTimeout        int
	maxOutstandingRequests  int
	address                 string
	sourceAddress           string
	secret                  string
	preauthenticationSecret string
	routingInstance         string
}

func resourceSystemRadiusServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSystemRadiusServerCreate,
		Read:   resourceSystemRadiusServerRead,
		Update: resourceSystemRadiusServerUpdate,
		Delete: resourceSystemRadiusServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemRadiusServerImport,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIPFunc(),
			},
			"secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"preauthentication_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"source_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIPFunc(),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 65535),
			},
			"accounting_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 65535),
			},
			"dynamic_request_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 65535),
			},
			"preauthentication_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 65535),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 1000),
			},
			"accouting_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validateIntRange(0, 1000),
			},
			"retry": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 100),
			},
			"accounting_retry": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validateIntRange(0, 100),
			},
			"max_outstanding_requests": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validateIntRange(0, 2000),
			},
			"routing_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
		},
	}
}

func resourceSystemRadiusServerCreate(d *schema.ResourceData, m interface{}) error {
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
	radiusServerExists, err := checkSystemRadiusServerExists(d.Get("address").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if radiusServerExists {
		sess.configClear(jnprSess)

		return fmt.Errorf("system radius-server %v already exists", d.Get("address").(string))
	}

	err = setSystemRadiusServer(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_system_radius_server", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	radiusServerExists, err = checkSystemRadiusServerExists(d.Get("address").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if radiusServerExists {
		d.SetId(d.Get("address").(string))
	} else {
		return fmt.Errorf("system radius-server %v not exists after commit => check your config", d.Get("address").(string))
	}

	return resourceSystemRadiusServerRead(d, m)
}
func resourceSystemRadiusServerRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	radiusServerOptions, err := readSystemRadiusServer(d.Get("address").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if radiusServerOptions.address == "" {
		d.SetId("")
	} else {
		fillSystemRadiusServerData(d, radiusServerOptions)
	}

	return nil
}
func resourceSystemRadiusServerUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delSystemRadiusServer(d.Get("address").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = setSystemRadiusServer(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_system_radius_server", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceSystemRadiusServerRead(d, m)
}
func resourceSystemRadiusServerDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delSystemRadiusServer(d.Get("address").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_system_radius_server", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
}
func resourceSystemRadiusServerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)

	radiusServerExists, err := checkSystemRadiusServerExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !radiusServerExists {
		return nil, fmt.Errorf("don't find system radius-server with id '%v' (id must be <address>)", d.Id())
	}
	radiusServerOptions, err := readSystemRadiusServer(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystemRadiusServerData(d, radiusServerOptions)

	result[0] = d

	return result, nil
}

func checkSystemRadiusServerExists(address string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	radiusServerConfig, err := sess.command("show configuration"+
		" system radius-server "+address+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if radiusServerConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSystemRadiusServer(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set system radius-server " + d.Get("address").(string)
	configSet := []string{setPrefix + " secret \"" + d.Get("secret").(string) + "\""}

	if d.Get("preauthentication_secret").(string) != "" {
		configSet = append(configSet, setPrefix+" preauthentication-secret \""+
			d.Get("preauthentication_secret").(string)+"\"")
	}
	if d.Get("source_address").(string) != "" {
		configSet = append(configSet, setPrefix+" source-address "+
			d.Get("source_address").(string))
	}
	if d.Get("port").(int) != 0 {
		configSet = append(configSet, setPrefix+" port "+
			strconv.Itoa(d.Get("port").(int)))
	}
	if d.Get("accounting_port").(int) != 0 {
		configSet = append(configSet, setPrefix+" accounting-port "+
			strconv.Itoa(d.Get("accounting_port").(int)))
	}
	if d.Get("dynamic_request_port").(int) != 0 {
		configSet = append(configSet, setPrefix+" dynamic-request-port "+
			strconv.Itoa(d.Get("dynamic_request_port").(int)))
	}
	if d.Get("preauthentication_port").(int) != 0 {
		configSet = append(configSet, setPrefix+" preauthentication-port "+
			strconv.Itoa(d.Get("preauthentication_port").(int)))
	}
	if d.Get("timeout").(int) != 0 {
		configSet = append(configSet, setPrefix+" timeout "+
			strconv.Itoa(d.Get("timeout").(int)))
	}
	if d.Get("accouting_timeout").(int) != -1 {
		configSet = append(configSet, setPrefix+" accounting-timeout "+
			strconv.Itoa(d.Get("accouting_timeout").(int)))
	}
	if d.Get("retry").(int) != 0 {
		configSet = append(configSet, setPrefix+" retry "+
			strconv.Itoa(d.Get("retry").(int)))
	}
	if d.Get("accounting_retry").(int) != -1 {
		configSet = append(configSet, setPrefix+" accounting-retry "+
			strconv.Itoa(d.Get("accounting_retry").(int)))
	}
	if d.Get("max_outstanding_requests").(int) != -1 {
		configSet = append(configSet, setPrefix+" max-outstanding-requests "+
			strconv.Itoa(d.Get("max_outstanding_requests").(int)))
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, setPrefix+" routing-instance "+
			d.Get("routing_instance").(string))
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func readSystemRadiusServer(server string, m interface{}, jnprSess *NetconfObject) (radiusServerOptions, error) {
	sess := m.(*Session)
	var confRead radiusServerOptions
	confRead.accountingRetry = -1
	confRead.accoutingTimeout = -1
	confRead.maxOutstandingRequests = -1

	radiusServerConfig, err := sess.command("show configuration"+
		" system radius-server "+server+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if radiusServerConfig != emptyWord {
		confRead.address = server
		for _, item := range strings.Split(radiusServerConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "secret "):
				var err error
				confRead.secret, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"secret "), "\""))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "preauthentication-secret "):
				var err error
				confRead.preauthenticationSecret, err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
					"preauthentication-secret "), "\""))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "source-address "):
				confRead.sourceAddress = strings.TrimPrefix(itemTrim, "source-address ")
			case strings.HasPrefix(itemTrim, "port "):
				var err error
				confRead.port, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "port "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "accounting-port "):
				var err error
				confRead.accountingPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accounting-port "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "dynamic-request-port "):
				var err error
				confRead.dynamicRequestPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "dynamic-request-port "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "preauthentication-port "):
				var err error
				confRead.preauthenticationPort, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "preauthentication-port "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "timeout "):
				var err error
				confRead.timeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "timeout "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "accounting-timeout "):
				var err error
				confRead.accoutingTimeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accounting-timeout "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "retry "):
				var err error
				confRead.retry, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "retry "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "accounting-retry "):
				var err error
				confRead.accountingRetry, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "accounting-retry "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "max-outstanding-requests "):
				var err error
				confRead.maxOutstandingRequests, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-outstanding-requests "))
				if err != nil {
					return confRead, err
				}
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

func delSystemRadiusServer(server string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete system radius-server "+server)
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func fillSystemRadiusServerData(d *schema.ResourceData, radiusServerOptions radiusServerOptions) {
	tfErr := d.Set("address", radiusServerOptions.address)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("secret", radiusServerOptions.secret)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("preauthentication_secret", radiusServerOptions.preauthenticationSecret)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("source_address", radiusServerOptions.sourceAddress)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("port", radiusServerOptions.port)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("accounting_port", radiusServerOptions.accountingPort)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("dynamic_request_port", radiusServerOptions.dynamicRequestPort)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("preauthentication_port", radiusServerOptions.preauthenticationPort)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("timeout", radiusServerOptions.timeout)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("accouting_timeout", radiusServerOptions.accoutingTimeout)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("retry", radiusServerOptions.retry)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("accounting_retry", radiusServerOptions.accountingRetry)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("max_outstanding_requests", radiusServerOptions.maxOutstandingRequests)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", radiusServerOptions.routingInstance)
	if tfErr != nil {
		panic(tfErr)
	}
}
