package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ipsecVpnOptions struct {
	name             string
	establishTunnels string
	bindInterface    string
	dfBit            string
	ike              []map[string]interface{}
	vpnMonitor       []map[string]interface{}
}

func resourceIpsecVpn() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpsecVpnCreate,
		ReadContext:   resourceIpsecVpnRead,
		UpdateContext: resourceIpsecVpnUpdate,
		DeleteContext: resourceIpsecVpnDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpsecVpnImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"establish_tunnels": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !stringInSlice(value, []string{"immediately", "on-traffic"}) {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not 'immediately' or 'on-traffic'", value, k))
					}

					return
				},
			},
			"bind_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bind_interface_auto": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"df_bit": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !stringInSlice(value, []string{"clear", "copy", "set"}) {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not 'clear', 'copy' or 'set'", value, k))
					}

					return
				},
			},
			"ike": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway": {
							Type:     schema.TypeString,
							Required: true,
						},
						"policy": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
						"identity_local": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateNetworkFunc(),
						},
						"identity_remote": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateNetworkFunc(),
						},
						"identity_service": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
					},
				},
			},
			"vpn_monitor": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_interface": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_interface_auto": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"destination_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateIPFunc(),
						},
						"optimized": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceIpsecVpnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security ipsec vpn not compatible with Junos device %s", jnprSess.Platform[0].Model))
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return diag.FromErr(err)
	}
	ipsecVpnExists, err := checkIpsecVpnExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if ipsecVpnExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security ipsec vpn %v already exists", d.Get("name").(string)))
	}
	if d.Get("bind_interface_auto").(bool) {
		newSt0, err := searchInterfaceSt0ToCreate(m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return diag.FromErr(fmt.Errorf("error for find new bind interface: %q", err))
		}
		tfErr := d.Set("bind_interface", newSt0)
		if tfErr != nil {
			panic(tfErr)
		}
	}
	err = setIpsecVpn(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	err = sess.commitConf("create resource junos_security_ipsec_vpn", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	ipsecVpnExists, err = checkIpsecVpnExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return diag.FromErr(err)
	}
	if ipsecVpnExists {
		d.SetId(d.Get("name").(string))
	} else {
		return diag.FromErr(fmt.Errorf("security ipsec vpn %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))
	}

	return resourceIpsecVpnRead(ctx, d, m)
}
func resourceIpsecVpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	ipsecVpnOptions, err := readIpsecVpn(d.Get("name").(string), m, jnprSess)
	// copy state vpn_monitor.0.source_interface_auto to struct
	if len(ipsecVpnOptions.vpnMonitor) > 0 {
		for _, v := range d.Get("vpn_monitor").([]interface{}) {
			if v != nil {
				stateMonitor := v.(map[string]interface{})
				vpnMonitor := ipsecVpnOptions.vpnMonitor[0]
				vpnMonitor["source_interface_auto"] = stateMonitor["source_interface_auto"].(bool)
				ipsecVpnOptions.vpnMonitor = []map[string]interface{}{vpnMonitor}
			}
		}
	}
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ipsecVpnOptions.name == "" {
		d.SetId("")
	} else {
		fillIpsecVpnData(d, ipsecVpnOptions)
	}

	return nil
}
func resourceIpsecVpnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return diag.FromErr(err)
	}
	err = delIpsecVpnConf(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	err = setIpsecVpn(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	err = sess.commitConf("update resource junos_security_ipsec_vpn", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	d.Partial(false)

	return resourceIpsecVpnRead(ctx, d, m)
}
func resourceIpsecVpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return diag.FromErr(err)
	}
	err = delIpsecVpn(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	err = sess.commitConf("delete resource junos_security_ipsec_vpn", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	return nil
}
func resourceIpsecVpnImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ipsecVpnExists, err := checkIpsecVpnExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ipsecVpnExists {
		return nil, fmt.Errorf("don't find security ipsec vpn with id '%v' (id must be <name>)", d.Id())
	}
	ipsecVpnOptions, err := readIpsecVpn(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillIpsecVpnData(d, ipsecVpnOptions)
	result[0] = d

	return result, nil
}

func checkIpsecVpnExists(ipsecVpn string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	ipsecVpnConfig, err := sess.command("show configuration security ipsec vpn "+ipsecVpn+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if ipsecVpnConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setIpsecVpn(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	if d.Get("bind_interface_auto").(bool) {
		configSet = append(configSet, "set interfaces "+d.Get("bind_interface").(string))
	}
	setPrefix := "set security ipsec vpn " + d.Get("name").(string)
	if d.Get("establish_tunnels").(string) != "" {
		configSet = append(configSet, setPrefix+" establish-tunnels "+d.Get("establish_tunnels").(string))
	}
	if d.Get("bind_interface").(string) != "" {
		configSet = append(configSet, setPrefix+" bind-interface "+d.Get("bind_interface").(string))
	}
	if d.Get("df_bit").(string) != "" {
		configSet = append(configSet, setPrefix+" df-bit "+d.Get("df_bit").(string))
	}
	for _, v := range d.Get("ike").([]interface{}) {
		ike := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+" ike gateway "+ike["gateway"].(string))
		configSet = append(configSet, setPrefix+" ike ipsec-policy "+ike["policy"].(string))
		if ike["identity_local"].(string) != "" {
			configSet = append(configSet, setPrefix+" ike proxy-identity local "+ike["identity_local"].(string))
		}
		if ike["identity_remote"].(string) != "" {
			configSet = append(configSet, setPrefix+" ike proxy-identity remote "+ike["identity_remote"].(string))
		}
		if ike["identity_service"].(string) != "" {
			configSet = append(configSet, setPrefix+" ike proxy-identity service "+ike["identity_service"].(string))
		}
	}
	for _, v := range d.Get("vpn_monitor").([]interface{}) {
		monitor := v.(map[string]interface{})
		configSet = append(configSet, "set security ipsec vpn "+d.Get("name").(string)+" vpn-monitor")
		if monitor["source_interface"].(string) != "" {
			configSet = append(configSet, setPrefix+" vpn-monitor source-interface "+
				monitor["source_interface"].(string))
		}
		if monitor["source_interface_auto"].(bool) {
			configSet = append(configSet, setPrefix+" vpn-monitor source-interface "+
				d.Get("bind_interface").(string))
		}
		if monitor["destination_ip"].(string) != "" {
			configSet = append(configSet, setPrefix+" vpn-monitor destination-ip "+
				monitor["destination_ip"].(string))
		}
		if monitor["optimized"].(bool) {
			configSet = append(configSet, setPrefix+" vpn-monitor optimized")
		}
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func readIpsecVpn(ipsecVpn string, m interface{}, jnprSess *NetconfObject) (ipsecVpnOptions, error) {
	sess := m.(*Session)
	var confRead ipsecVpnOptions

	ipsecVpnConfig, err := sess.command("show configuration"+
		" security ipsec vpn "+ipsecVpn+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if ipsecVpnConfig != emptyWord {
		confRead.name = ipsecVpn
		for _, item := range strings.Split(ipsecVpnConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "bind-interface "):
				confRead.bindInterface = strings.TrimPrefix(itemTrim, "bind-interface ")
			case strings.HasPrefix(itemTrim, "establish-tunnels "):
				confRead.establishTunnels = strings.TrimPrefix(itemTrim, "establish-tunnels ")
			case strings.HasPrefix(itemTrim, "df-bit "):
				confRead.dfBit = strings.TrimPrefix(itemTrim, "df-bit ")
			case strings.HasPrefix(itemTrim, "ike "):
				ikeOptions := map[string]interface{}{
					"gateway":          "",
					"policy":           "",
					"identity_local":   "",
					"identity_remote":  "",
					"identity_service": "",
				}
				if len(confRead.ike) > 0 {
					for k, v := range confRead.ike[0] {
						ikeOptions[k] = v
					}
				}
				switch {
				case strings.HasPrefix(itemTrim, "ike gateway "):
					ikeOptions["gateway"] = strings.TrimPrefix(itemTrim, "ike gateway ")
				case strings.HasPrefix(itemTrim, "ike ipsec-policy "):
					ikeOptions["policy"] = strings.TrimPrefix(itemTrim, "ike ipsec-policy ")
				case strings.HasPrefix(itemTrim, "ike proxy-identity local "):
					ikeOptions["identity_local"] = strings.TrimPrefix(itemTrim, "ike proxy-identity local ")
				case strings.HasPrefix(itemTrim, "ike proxy-identity remote "):
					ikeOptions["identity_remote"] = strings.TrimPrefix(itemTrim, "ike proxy-identity remote ")
				case strings.HasPrefix(itemTrim, "ike proxy-identity service "):
					ikeOptions["identity_service"] = strings.TrimPrefix(itemTrim, "ike proxy-identity service ")
				}
				// override (maxItem = 1)
				confRead.ike = []map[string]interface{}{ikeOptions}
			case strings.HasPrefix(itemTrim, "vpn-monitor "):
				monitorOptions := map[string]interface{}{
					"source_interface": "",
					"destination_ip":   "",
					"optimized":        false,
				}
				if len(confRead.vpnMonitor) > 0 {
					for k, v := range confRead.vpnMonitor[0] {
						monitorOptions[k] = v
					}
				}
				switch {
				case strings.HasPrefix(itemTrim, "vpn-monitor source-interface "):
					monitorOptions["source_interface"] = strings.TrimPrefix(itemTrim, "vpn-monitor source-interface ")
				case strings.HasPrefix(itemTrim, "vpn-monitor destination-ip "):
					monitorOptions["destination_ip"] = strings.TrimPrefix(itemTrim, "vpn-monitor destination-ip ")
				case strings.HasPrefix(itemTrim, "vpn-monitor optimized"):
					monitorOptions["optimized"] = true
				}
				// override (maxItem = 1)
				confRead.vpnMonitor = []map[string]interface{}{monitorOptions}
			}
		}
	} else {
		confRead.name = ""

		return confRead, nil
	}

	return confRead, nil
}
func delIpsecVpnConf(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec vpn "+d.Get("name").(string))
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func delIpsecVpn(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec vpn "+d.Get("name").(string))
	if d.Get("bind_interface_auto").(bool) {
		empty := checkInterfaceNC(d.Get("bind_interface").(string), m, jnprSess)
		if empty == nil {
			configSet = append(configSet, "delete interfaces "+d.Get("bind_interface").(string))
		}
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}

func fillIpsecVpnData(d *schema.ResourceData, ipsecVpnOptions ipsecVpnOptions) {
	tfErr := d.Set("name", ipsecVpnOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("establish_tunnels", ipsecVpnOptions.establishTunnels)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("bind_interface", ipsecVpnOptions.bindInterface)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("df_bit", ipsecVpnOptions.dfBit)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("ike", ipsecVpnOptions.ike)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("vpn_monitor", ipsecVpnOptions.vpnMonitor)
	if tfErr != nil {
		panic(tfErr)
	}
}

func searchInterfaceSt0ToCreate(m interface{}, jnprSess *NetconfObject) (string, error) {
	sess := m.(*Session)
	st0, err := sess.command("show interfaces st0 terse", jnprSess)
	if err != nil {
		return "", err
	}
	st0Line := strings.Split(st0, "\n")
	st0int := make([]string, 0)
	for _, line := range st0Line {
		if strings.HasPrefix(line, "st0.") {
			lineSplit := strings.Split(line, " ")
			st0int = append(st0int, lineSplit[0])
		}
	}
	for i := 0; i <= 1073741823; i++ {
		if !stringInSlice("st0."+strconv.Itoa(i), st0int) {
			return "st0." + strconv.Itoa(i), nil
		}
	}

	return "", fmt.Errorf("error for find st0 unit to create")
}
