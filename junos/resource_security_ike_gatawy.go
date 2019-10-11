package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type ikeGatewayOptions struct {
	generalIkeid      bool
	noNatTraversal    bool
	name              string
	externalInterface string
	policy            string
	version           string
	localAddress      string
	address           []string
	deadPeerDetection []map[string]interface{}
	localIdentity     []map[string]interface{}
	remoteIdentity    []map[string]interface{}
}

func resourceIkeGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceIkeGatewayCreate,
		Read:   resourceIkeGatewayRead,
		Update: resourceIkeGatewayUpdate,
		Delete: resourceIkeGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIkeGatewayImport,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"local_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIPFunc(),
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_interface": {
				Type:     schema.TypeString,
				Required: true,
			},
			"general_ike_id": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_nat_traversal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dead_peer_detection": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(10, 60),
						},
						"threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 5),
						},
					},
				},
			},
			"local_identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"distinguished-name", "hostname", "inet", "inet6", "user-at-hostname"}) {
									errors = append(errors, fmt.Errorf(
										"%q for %q is not valid", value, k))
								}
								return
							},
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"remote_identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"distinguished-name", "hostname", "inet", "inet6", "user-at-hostname"}) {
									errors = append(errors, fmt.Errorf(
										"%q for %q is not valid", value, k))
								}
								return
							},
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !stringInSlice(value, []string{"v1-only", "v2-only"}) {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not 'v1-only' or 'v2-only'", value, k))
					}
					return
				},
			},
		},
	}
}

func resourceIkeGatewayCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security ike gateway not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	ikeGatewayExists, err := checkIkeGatewayExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if ikeGatewayExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("ike gateway %v already exists", d.Get("name").(string))
	}
	err = setIkeGateway(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("create resource junos_ike_gateway", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	ikeGatewayExists, err = checkIkeGatewayExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if ikeGatewayExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("ike gateway %v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceIkeGatewayRead(d, m)
}
func resourceIkeGatewayRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	ikeGatewayOptions, err := readIkeGateway(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if ikeGatewayOptions.name == "" {
		d.SetId("")
	} else {
		fillIkeGatewayData(d, ikeGatewayOptions)
	}
	return nil
}
func resourceIkeGatewayUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delIkeGateway(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setIkeGateway(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("update resource junos_ike_gateway", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceIkeGatewayRead(d, m)
}
func resourceIkeGatewayDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delIkeGateway(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf("delete resource junos_ike_gateway", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
}
func resourceIkeGatewayImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ikeGatewayExists, err := checkIkeGatewayExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ikeGatewayExists {
		return nil, fmt.Errorf("don't find ike gateway with id '%v' (id must be <name>)", d.Id())
	}
	ikeGatewayOptions, err := readIkeGateway(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillIkeGatewayData(d, ikeGatewayOptions)
	result[0] = d
	return result, nil
}

func checkIkeGatewayExists(ikeGateway string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	ikeGatewayConfig, err := sess.command("show configuration security ike gateway "+ikeGateway+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if ikeGatewayConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setIkeGateway(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security ike gateway " + d.Get("name").(string)
	for _, v := range d.Get("address").([]interface{}) {
		err := validateIP(v.(string))
		if err != nil {
			return err
		}
		configSet = append(configSet, setPrefix+" address "+v.(string)+"\n")
	}
	if d.Get("local_address").(string) != "" {
		configSet = append(configSet, setPrefix+" local-address "+d.Get("local_address").(string)+"\n")
	}
	if d.Get("policy").(string) != "" {
		configSet = append(configSet, setPrefix+" ike-policy "+d.Get("policy").(string)+"\n")
	}
	if d.Get("external_interface").(string) != "" {
		configSet = append(configSet, setPrefix+" external-interface "+d.Get("external_interface").(string)+"\n")
	}
	if d.Get("general_ike_id").(bool) {
		configSet = append(configSet, setPrefix+" general-ikeid\n")
	}
	if d.Get("no_nat_traversal").(bool) {
		configSet = append(configSet, setPrefix+" no-nat-traversal\n")
	}
	if len(d.Get("dead_peer_detection").([]interface{})) != 0 {
		configSet = append(configSet, setPrefix+" dead-peer-detection\n")
		for _, v := range d.Get("dead_peer_detection").([]interface{}) {
			deadPeerOptions := v.(map[string]interface{})
			if deadPeerOptions["interval"].(int) != 0 {
				configSet = append(configSet, setPrefix+" dead-peer-detection interval "+
					strconv.Itoa(deadPeerOptions["interval"].(int))+"\n")
			}
			if deadPeerOptions["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+" dead-peer-detection threshold "+
					strconv.Itoa(deadPeerOptions["threshold"].(int))+"\n")
			}
		}
	}
	for _, v := range d.Get("local_identity").([]interface{}) {
		localIdentity := v.(map[string]interface{})
		if localIdentity["type"].(string) == "distinguished-name" {
			if localIdentity["value"].(string) != "" {
				return fmt.Errorf("no value for option distinguished-name in ike gateway confifg")
			}
		} else {
			if localIdentity["value"].(string) == "" {
				return fmt.Errorf("missing value for option local-identity %s in ike gateway confifg",
					localIdentity["type"].(string))
			}
		}
		configSet = append(configSet, setPrefix+" local-identity "+
			localIdentity["type"].(string)+" "+localIdentity["value"].(string)+"\n")
	}
	for _, v := range d.Get("remote_identity").([]interface{}) {
		remoteIdentity := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+" remote-identity "+
			remoteIdentity["type"].(string)+" "+remoteIdentity["value"].(string)+"\n")
	}
	if d.Get("version").(string) != "" {
		configSet = append(configSet, setPrefix+" version "+d.Get("version").(string)+"\n")
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readIkeGateway(ikeGateway string, m interface{}, jnprSess *NetconfObject) (ikeGatewayOptions, error) {
	sess := m.(*Session)
	var confRead ikeGatewayOptions

	ikeGatewayConfig, err := sess.command("show configuration"+
		" security ike gateway "+ikeGateway+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if ikeGatewayConfig != emptyWord {
		confRead.name = ikeGateway
		for _, item := range strings.Split(ikeGatewayConfig, "\n") {
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
			case strings.HasPrefix(itemTrim, "local-address "):
				confRead.localAddress = strings.TrimPrefix(itemTrim, "local-address ")
			case strings.HasPrefix(itemTrim, "ike-policy "):
				confRead.policy = strings.TrimPrefix(itemTrim, "ike-policy ")
			case strings.HasPrefix(itemTrim, "external-interface "):
				confRead.externalInterface = strings.TrimPrefix(itemTrim, "external-interface ")
			case strings.HasPrefix(itemTrim, "general-ikeid"):
				confRead.generalIkeid = true
			case strings.HasPrefix(itemTrim, "no-nat-traversal"):
				confRead.noNatTraversal = true
			case strings.HasPrefix(itemTrim, "dead-peer-detection"):
				deadPeerOptions := map[string]interface{}{
					"interval":  0,
					"threshold": 0,
				}
				if len(confRead.deadPeerDetection) > 0 {
					for k, v := range confRead.deadPeerDetection[0] {
						deadPeerOptions[k] = v
					}
				}
				switch {
				case strings.HasPrefix(itemTrim, "dead-peer-detection interval "):
					deadPeerOptions["interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim,
						"dead-peer-detection interval "))
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "dead-peer-detection threshold "):
					deadPeerOptions["threshold"], err = strconv.Atoi(strings.TrimPrefix(itemTrim,
						"dead-peer-detection threshold "))
					if err != nil {
						return confRead, err
					}
				}
				// override (maxItem = 1)
				confRead.deadPeerDetection = []map[string]interface{}{deadPeerOptions}
			case strings.HasPrefix(itemTrim, "local-identity "):
				localIdentityOptions := map[string]interface{}{
					"type":  "",
					"value": "",
				}
				readLocalIdentity := strings.TrimPrefix(itemTrim, "local-identity ")
				readLocalIdentityList := strings.Split(readLocalIdentity, " ")
				localIdentityOptions["type"] = readLocalIdentityList[0]
				if len(readLocalIdentityList) > 1 {
					localIdentityOptions["value"] = readLocalIdentityList[1]
				}
				// override (maxItem = 1)
				confRead.localIdentity = []map[string]interface{}{localIdentityOptions}
			case strings.HasPrefix(itemTrim, "remote-identity "):
				remoteIdentityOptions := map[string]interface{}{
					"type":  "",
					"value": "",
				}
				readRemoteIdentity := strings.TrimPrefix(itemTrim, "remote-identity ")
				readRemoteIdentityList := strings.Split(readRemoteIdentity, " ")
				remoteIdentityOptions["type"] = readRemoteIdentityList[0]
				remoteIdentityOptions["value"] = readRemoteIdentityList[1]
				// override (maxItem = 1)
				confRead.remoteIdentity = []map[string]interface{}{remoteIdentityOptions}
			case strings.HasPrefix(itemTrim, "version "):
				confRead.version = strings.TrimPrefix(itemTrim, "version ")
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}
func delIkeGateway(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ike gateway "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillIkeGatewayData(d *schema.ResourceData, ikeGatewayOptions ikeGatewayOptions) {
	tfErr := d.Set("name", ikeGatewayOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("address", ikeGatewayOptions.address)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_address", ikeGatewayOptions.localAddress)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("policy", ikeGatewayOptions.policy)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("external_interface", ikeGatewayOptions.externalInterface)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("general_ike_id", ikeGatewayOptions.generalIkeid)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("no_nat_traversal", ikeGatewayOptions.noNatTraversal)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("dead_peer_detection", ikeGatewayOptions.deadPeerDetection)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_identity", ikeGatewayOptions.localIdentity)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("remote_identity", ikeGatewayOptions.remoteIdentity)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("version", ikeGatewayOptions.version)
	if tfErr != nil {
		panic(tfErr)
	}
}
