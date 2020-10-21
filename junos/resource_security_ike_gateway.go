package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		CreateContext: resourceIkeGatewayCreate,
		ReadContext:   resourceIkeGatewayRead,
		UpdateContext: resourceIkeGatewayUpdate,
		DeleteContext: resourceIkeGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIkeGatewayImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}),
			},
			"address": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"local_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
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
							ValidateFunc: validation.IntBetween(10, 60),
						},
						"threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},
						"send_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"always-send", "optimized", "probe-idle-tunnel"}, false),
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
							ValidateFunc: validation.StringInSlice([]string{
								"distinguished-name", "hostname", "inet", "inet6", "user-at-hostname"}, false),
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
							ValidateFunc: validation.StringInSlice([]string{
								"distinguished-name", "hostname", "inet", "inet6", "user-at-hostname"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"v1-only", "v2-only"}, false),
			},
		},
	}
}

func resourceIkeGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security ike gateway not compatible with Junos device %s",
			jnprSess.Platform[0].Model))
	}
	sess.configLock(jnprSess)
	ikeGatewayExists, err := checkIkeGatewayExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if ikeGatewayExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security ike gateway %v already exists", d.Get("name").(string)))
	}
	if err := setIkeGateway(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := sess.commitConf("create resource junos_security_ike_gateway", jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	ikeGatewayExists, err = checkIkeGatewayExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return diag.FromErr(err)
	}
	if ikeGatewayExists {
		d.SetId(d.Get("name").(string))
	} else {
		return diag.FromErr(fmt.Errorf("security ike gateway %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))
	}

	return resourceIkeGatewayRead(ctx, d, m)
}
func resourceIkeGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	ikeGatewayOptions, err := readIkeGateway(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ikeGatewayOptions.name == "" {
		d.SetId("")
	} else {
		fillIkeGatewayData(d, ikeGatewayOptions)
	}

	return nil
}
func resourceIkeGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delIkeGateway(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setIkeGateway(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := sess.commitConf("update resource junos_security_ike_gateway", jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	d.Partial(false)

	return resourceIkeGatewayRead(ctx, d, m)
}
func resourceIkeGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delIkeGateway(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := sess.commitConf("delete resource junos_security_ike_gateway", jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
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
		return nil, fmt.Errorf("don't find security ike gateway with id '%v' (id must be <name>)", d.Id())
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
		_, errs := validation.IsIPAddress(v, "address")
		if len(errs) > 0 {
			return errs[0]
		}
		configSet = append(configSet, setPrefix+" address "+v.(string))
	}
	if d.Get("local_address").(string) != "" {
		configSet = append(configSet, setPrefix+" local-address "+d.Get("local_address").(string))
	}
	if d.Get("policy").(string) != "" {
		configSet = append(configSet, setPrefix+" ike-policy "+d.Get("policy").(string))
	}
	if d.Get("external_interface").(string) != "" {
		configSet = append(configSet, setPrefix+" external-interface "+d.Get("external_interface").(string))
	}
	if d.Get("general_ike_id").(bool) {
		configSet = append(configSet, setPrefix+" general-ikeid")
	}
	if d.Get("no_nat_traversal").(bool) {
		configSet = append(configSet, setPrefix+" no-nat-traversal")
	}
	if len(d.Get("dead_peer_detection").([]interface{})) != 0 {
		configSet = append(configSet, setPrefix+" dead-peer-detection")
		for _, v := range d.Get("dead_peer_detection").([]interface{}) {
			if v != nil {
				deadPeerOptions := v.(map[string]interface{})
				if deadPeerOptions["interval"].(int) != 0 {
					configSet = append(configSet, setPrefix+" dead-peer-detection interval "+
						strconv.Itoa(deadPeerOptions["interval"].(int)))
				}
				if deadPeerOptions["threshold"].(int) != 0 {
					configSet = append(configSet, setPrefix+" dead-peer-detection threshold "+
						strconv.Itoa(deadPeerOptions["threshold"].(int)))
				}
				if deadPeerOptions["send_mode"].(string) != "" {
					configSet = append(configSet, setPrefix+" dead-peer-detection "+
						deadPeerOptions["send_mode"].(string))
				}
			}
		}
	}
	for _, v := range d.Get("local_identity").([]interface{}) {
		localIdentity := v.(map[string]interface{})
		if localIdentity["type"].(string) == "distinguished-name" {
			if localIdentity["value"].(string) != "" {
				return fmt.Errorf("no value for option distinguished-name in security ike gateway config")
			}
		} else {
			if localIdentity["value"].(string) == "" {
				return fmt.Errorf("missing value for option local-identity %s in security ike gateway config",
					localIdentity["type"].(string))
			}
		}
		configSet = append(configSet, setPrefix+" local-identity "+
			localIdentity["type"].(string)+" "+localIdentity["value"].(string))
	}
	for _, v := range d.Get("remote_identity").([]interface{}) {
		remoteIdentity := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+" remote-identity "+
			remoteIdentity["type"].(string)+" "+remoteIdentity["value"].(string))
	}
	if d.Get("version").(string) != "" {
		configSet = append(configSet, setPrefix+" version "+d.Get("version").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
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
					"send_mode": "",
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
				case strings.HasSuffix(itemTrim, " always-send"):
					deadPeerOptions["send_mode"] = "always-send"
				case strings.HasSuffix(itemTrim, " optimized"):
					deadPeerOptions["send_mode"] = "optimized"
				case strings.HasSuffix(itemTrim, " probe-idle-tunnel"):
					deadPeerOptions["send_mode"] = "probe-idle-tunnel"
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
	configSet = append(configSet, "delete security ike gateway "+d.Get("name").(string))
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillIkeGatewayData(d *schema.ResourceData, ikeGatewayOptions ikeGatewayOptions) {
	if tfErr := d.Set("name", ikeGatewayOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", ikeGatewayOptions.address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_address", ikeGatewayOptions.localAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("policy", ikeGatewayOptions.policy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("external_interface", ikeGatewayOptions.externalInterface); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("general_ike_id", ikeGatewayOptions.generalIkeid); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_nat_traversal", ikeGatewayOptions.noNatTraversal); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dead_peer_detection", ikeGatewayOptions.deadPeerDetection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_identity", ikeGatewayOptions.localIdentity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("remote_identity", ikeGatewayOptions.remoteIdentity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version", ikeGatewayOptions.version); tfErr != nil {
		panic(tfErr)
	}
}
