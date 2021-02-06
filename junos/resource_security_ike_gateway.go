package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	jdecode "github.com/jeremmfr/junosdecode"
)

type ikeGatewayOptions struct {
	generalIkeid      bool
	noNatTraversal    bool
	name              string
	externalInterface string
	policy            string
	localAddress      string
	version           string
	address           []string
	aaa               []map[string]interface{}
	deadPeerDetection []map[string]interface{}
	dynamicRemote     []map[string]interface{}
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
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
			},
			"external_interface": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				MaxItems:     5,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"address", "dynamic_remote"},
			},
			"dynamic_remote": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  []string{"address", "dynamic_remote"},
				ConflictsWith: []string{"general_ike_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connections_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"distinguished_name": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ConflictsWith: []string{
								"dynamic_remote.0.hostname",
								"dynamic_remote.0.inet",
								"dynamic_remote.0.inet6",
								"dynamic_remote.0.user_at_hostname",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"wildcard": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"hostname": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateAddress(),
							ConflictsWith: []string{
								"dynamic_remote.0.distinguished_name",
								"dynamic_remote.0.inet",
								"dynamic_remote.0.inet6",
								"dynamic_remote.0.user_at_hostname",
							},
						},
						"ike_user_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"shared-ike-id", "group-ike-id"}, false),
						},
						"inet": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
							ConflictsWith: []string{
								"dynamic_remote.0.distinguished_name",
								"dynamic_remote.0.hostname",
								"dynamic_remote.0.inet6",
								"dynamic_remote.0.user_at_hostname",
							},
						},
						"inet6": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
							ConflictsWith: []string{
								"dynamic_remote.0.distinguished_name",
								"dynamic_remote.0.hostname",
								"dynamic_remote.0.inet",
								"dynamic_remote.0.user_at_hostname",
							},
						},
						"reject_duplicate_connection": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"user_at_hostname": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"dynamic_remote.0.distinguished_name",
								"dynamic_remote.0.hostname",
								"dynamic_remote.0.inet",
								"dynamic_remote.0.inet6",
							},
						},
					},
				},
			},
			"aaa": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_profile": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
							ConflictsWith: []string{
								"aaa.0.client_password",
								"aaa.0.client_username",
							},
						},
						"client_password": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"aaa.0.access_profile",
							},
							ValidateFunc: validation.StringLenBetween(1, 128),
						},
						"client_username": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"aaa.0.access_profile",
							},
							ValidateFunc: validation.StringLenBetween(1, 128),
						},
					},
				},
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
						"send_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"always-send", "optimized", "probe-idle-tunnel"}, false),
						},
						"threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},
					},
				},
			},
			"general_ike_id": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"dynamic_remote"},
			},
			"local_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
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
			"no_nat_traversal": {
				Type:     schema.TypeBool,
				Optional: true,
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
			jnprSess.SystemInformation.HardwareModel))
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
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_ike_gateway", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	ikeGatewayExists, err = checkIkeGatewayExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ikeGatewayExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ike gateway %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIkeGatewayReadWJnprSess(d, m, jnprSess)...)
}
func resourceIkeGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceIkeGatewayReadWJnprSess(d, m, jnprSess)
}
func resourceIkeGatewayReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
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
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_ike_gateway", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIkeGatewayReadWJnprSess(d, m, jnprSess)...)
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
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_ike_gateway", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
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
	configSet = append(configSet, setPrefix+" ike-policy "+d.Get("policy").(string))
	configSet = append(configSet, setPrefix+" external-interface "+d.Get("external_interface").(string))
	for _, v := range d.Get("address").([]interface{}) {
		_, errs := validation.IsIPAddress(v, "address")
		if len(errs) > 0 {
			return errs[0]
		}
		configSet = append(configSet, setPrefix+" address "+v.(string))
	}
	for _, v := range d.Get("dynamic_remote").([]interface{}) {
		if v != nil {
			dynamicRemote := v.(map[string]interface{})
			if dynamicRemote["connections_limit"].(int) > 0 {
				configSet = append(configSet, setPrefix+" dynamic connections-limit "+
					strconv.Itoa(dynamicRemote["connections_limit"].(int)))
			}
			for _, v2 := range dynamicRemote["distinguished_name"].([]interface{}) {
				configSet = append(configSet, setPrefix+" dynamic distinguished-name")
				if v2 != nil {
					distinguishedName := v2.(map[string]interface{})
					if distinguishedName["container"].(string) != "" {
						configSet = append(configSet, setPrefix+" dynamic distinguished-name container \""+
							distinguishedName["container"].(string)+"\"")
					}
					if distinguishedName["wildcard"].(string) != "" {
						configSet = append(configSet, setPrefix+" dynamic distinguished-name wildcard \""+
							distinguishedName["container"].(string)+"\"")
					}
				}
			}
			if dynamicRemote["hostname"].(string) != "" {
				configSet = append(configSet, setPrefix+" dynamic hostname "+
					dynamicRemote["hostname"].(string))
			}
			if dynamicRemote["ike_user_type"].(string) != "" {
				configSet = append(configSet, setPrefix+" dynamic ike-user-type "+
					dynamicRemote["ike_user_type"].(string))
			}
			if dynamicRemote["inet"].(string) != "" {
				configSet = append(configSet, setPrefix+" dynamic inet "+
					dynamicRemote["inet"].(string))
			}
			if dynamicRemote["inet6"].(string) != "" {
				configSet = append(configSet, setPrefix+" dynamic inet6 "+
					dynamicRemote["inet6"].(string))
			}
			if dynamicRemote["reject_duplicate_connection"].(bool) {
				configSet = append(configSet, setPrefix+" dynamic reject-duplicate-connection")
			}
			if dynamicRemote["user_at_hostname"].(string) != "" {
				configSet = append(configSet, setPrefix+" dynamic user-at-hostname \""+
					dynamicRemote["user_at_hostname"].(string)+"\"")
			}
		}
	}
	for _, v := range d.Get("aaa").([]interface{}) {
		if v != nil {
			aaa := v.(map[string]interface{})
			if aaa["access_profile"].(string) != "" {
				configSet = append(configSet, setPrefix+" aaa access-profile "+aaa["access_profile"].(string))
			}
			if aaa["client_password"].(string) != "" {
				configSet = append(configSet, setPrefix+" aaa client password "+aaa["client_password"].(string))
			}
			if aaa["client_username"].(string) != "" {
				configSet = append(configSet, setPrefix+" aaa client username "+aaa["client_username"].(string))
			}
		}
	}
	for _, v := range d.Get("dead_peer_detection").([]interface{}) {
		configSet = append(configSet, setPrefix+" dead-peer-detection")
		if v != nil {
			deadPeerOptions := v.(map[string]interface{})
			if deadPeerOptions["interval"].(int) != 0 {
				configSet = append(configSet, setPrefix+" dead-peer-detection interval "+
					strconv.Itoa(deadPeerOptions["interval"].(int)))
			}
			if deadPeerOptions["send_mode"].(string) != "" {
				configSet = append(configSet, setPrefix+" dead-peer-detection "+
					deadPeerOptions["send_mode"].(string))
			}
			if deadPeerOptions["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+" dead-peer-detection threshold "+
					strconv.Itoa(deadPeerOptions["threshold"].(int)))
			}
		}
	}
	if d.Get("general_ike_id").(bool) {
		configSet = append(configSet, setPrefix+" general-ikeid")
	}
	if d.Get("local_address").(string) != "" {
		configSet = append(configSet, setPrefix+" local-address "+d.Get("local_address").(string))
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
	if d.Get("no_nat_traversal").(bool) {
		configSet = append(configSet, setPrefix+" no-nat-traversal")
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
			case strings.HasPrefix(itemTrim, "external-interface "):
				confRead.externalInterface = strings.TrimPrefix(itemTrim, "external-interface ")
			case strings.HasPrefix(itemTrim, "ike-policy "):
				confRead.policy = strings.TrimPrefix(itemTrim, "ike-policy ")
			case strings.HasPrefix(itemTrim, "address "):
				confRead.address = append(confRead.address, strings.TrimPrefix(itemTrim, "address "))
			case strings.HasPrefix(itemTrim, "dynamic "):
				if len(confRead.dynamicRemote) == 0 {
					confRead.dynamicRemote = append(confRead.dynamicRemote, map[string]interface{}{
						"connections_limit":           0,
						"distinguished_name":          make([]map[string]interface{}, 0),
						"hostname":                    "",
						"ike_user_type":               "",
						"inet":                        "",
						"inet6":                       "",
						"reject_duplicate_connection": false,
						"user_at_hostname":            "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "dynamic connections-limit "):
					confRead.dynamicRemote[0]["connections_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim,
						"dynamic connections-limit "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "dynamic distinguished-name"):
					if len(confRead.dynamicRemote[0]["distinguished_name"].([]map[string]interface{})) == 0 {
						confRead.dynamicRemote[0]["distinguished_name"] = append(
							confRead.dynamicRemote[0]["distinguished_name"].([]map[string]interface{}), map[string]interface{}{
								"container": "",
								"wildcard":  "",
							})
					}
					switch {
					case strings.HasPrefix(itemTrim, "dynamic distinguished-name container "):
						confRead.dynamicRemote[0]["distinguished_name"].([]map[string]interface{})[0]["container"] =
							strings.Trim(strings.TrimPrefix(itemTrim, "dynamic distinguished-name container "), "\"")
					case strings.HasPrefix(itemTrim, "dynamic distinguished-name wildcard "):
						confRead.dynamicRemote[0]["distinguished_name"].([]map[string]interface{})[0]["wildcard"] =
							strings.Trim(strings.TrimPrefix(itemTrim, "dynamic distinguished-name wildcard "), "\"")
					}
				case strings.HasPrefix(itemTrim, "dynamic hostname "):
					confRead.dynamicRemote[0]["hostname"] = strings.TrimPrefix(itemTrim, "dynamic hostname ")
				case strings.HasPrefix(itemTrim, "dynamic ike-user-type "):
					confRead.dynamicRemote[0]["ike_user_type"] = strings.TrimPrefix(itemTrim, "dynamic ike-user-type ")
				case strings.HasPrefix(itemTrim, "dynamic inet "):
					confRead.dynamicRemote[0]["inet"] = strings.TrimPrefix(itemTrim, "dynamic inet ")
				case strings.HasPrefix(itemTrim, "dynamic inet6 "):
					confRead.dynamicRemote[0]["inet6"] = strings.TrimPrefix(itemTrim, "dynamic inet6 ")
				case itemTrim == "dynamic reject-duplicate-connection":
					confRead.dynamicRemote[0]["reject_duplicate_connection"] = true
				case strings.HasPrefix(itemTrim, "dynamic user-at-hostname "):
					confRead.dynamicRemote[0]["user_at_hostname"] = strings.Trim(strings.TrimPrefix(
						itemTrim, "dynamic user-at-hostname "), "\"")
				}
			case strings.HasPrefix(itemTrim, "aaa "):
				if len(confRead.aaa) == 0 {
					confRead.aaa = append(confRead.aaa, map[string]interface{}{
						"access_profile":  "",
						"client_password": "",
						"client_username": "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "aaa access-profile "):
					confRead.aaa[0]["access_profile"] = strings.TrimPrefix(itemTrim, "aaa access-profile ")
				case strings.HasPrefix(itemTrim, "aaa client password "):
					confRead.aaa[0]["client_password"], err = jdecode.Decode(strings.Trim(strings.TrimPrefix(itemTrim,
						"aaa client password "), "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode aaa client password : %w", err)
					}
				case strings.HasPrefix(itemTrim, "aaa client username "):
					confRead.aaa[0]["client_username"] = strings.TrimPrefix(itemTrim, "aaa client username ")
				}
			case strings.HasPrefix(itemTrim, "dead-peer-detection"):
				deadPeerOptions := map[string]interface{}{
					"interval":  0,
					"send_mode": "",
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
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasSuffix(itemTrim, " always-send"):
					deadPeerOptions["send_mode"] = "always-send"
				case strings.HasSuffix(itemTrim, " optimized"):
					deadPeerOptions["send_mode"] = "optimized"
				case strings.HasSuffix(itemTrim, " probe-idle-tunnel"):
					deadPeerOptions["send_mode"] = "probe-idle-tunnel"
				case strings.HasPrefix(itemTrim, "dead-peer-detection threshold "):
					deadPeerOptions["threshold"], err = strconv.Atoi(strings.TrimPrefix(itemTrim,
						"dead-peer-detection threshold "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				}
				// override (maxItem = 1)
				confRead.deadPeerDetection = []map[string]interface{}{deadPeerOptions}
			case itemTrim == "general-ikeid":
				confRead.generalIkeid = true
			case strings.HasPrefix(itemTrim, "local-address "):
				confRead.localAddress = strings.TrimPrefix(itemTrim, "local-address ")
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
			case itemTrim == "no-nat-traversal":
				confRead.noNatTraversal = true
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
	if tfErr := d.Set("external_interface", ikeGatewayOptions.externalInterface); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("policy", ikeGatewayOptions.policy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address", ikeGatewayOptions.address); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_remote", ikeGatewayOptions.dynamicRemote); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("aaa", ikeGatewayOptions.aaa); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dead_peer_detection", ikeGatewayOptions.deadPeerDetection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("general_ike_id", ikeGatewayOptions.generalIkeid); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_address", ikeGatewayOptions.localAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_identity", ikeGatewayOptions.localIdentity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_nat_traversal", ikeGatewayOptions.noNatTraversal); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("remote_identity", ikeGatewayOptions.remoteIdentity); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version", ikeGatewayOptions.version); tfErr != nil {
		panic(tfErr)
	}
}
