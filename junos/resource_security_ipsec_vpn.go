package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type ipsecVpnOptions struct {
	name             string
	establishTunnels string
	bindInterface    string
	dfBit            string
	ike              []map[string]interface{}
	trafficSelector  []map[string]interface{}
	vpnMonitor       []map[string]interface{}
}

func resourceIpsecVpn() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIpsecVpnCreate,
		ReadWithoutTimeout:   resourceIpsecVpnRead,
		UpdateWithoutTimeout: resourceIpsecVpnUpdate,
		DeleteWithoutTimeout: resourceIpsecVpnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpsecVpnImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
			},
			"bind_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bind_interface_auto": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"traffic_selector"},
				Deprecated:    "Use the junos_interface_st0_unit resource instead",
			},
			"df_bit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"clear", "copy", "set"}, false),
			},
			"establish_tunnels": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"immediately", "on-traffic"}, false),
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
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
						},
						"identity_local": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
						"identity_remote": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
						"identity_service": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
					},
				},
			},
			"traffic_selector": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"bind_interface_auto"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
						},
						"local_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
						"remote_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
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
						"destination_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"optimized": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"source_interface": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_interface_auto": {
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
	if sess.junosFakeCreateSetFile != "" && !d.Get("bind_interface_auto").(bool) {
		if err := setIpsecVpn(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security ipsec vpn not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ipsecVpnExists, err := checkIpsecVpnExists(d.Get("name").(string), sess, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecVpnExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec vpn %v already exists", d.Get("name").(string)))...)
	}
	if d.Get("bind_interface_auto").(bool) {
		newSt0, err := searchInterfaceSt0UnitToCreate(sess, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(fmt.Errorf("error to find new bind interface: %w", err))...)
		}
		tfErr := d.Set("bind_interface", newSt0)
		if tfErr != nil {
			panic(tfErr)
		}
	}
	if err := setIpsecVpn(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_security_ipsec_vpn", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ipsecVpnExists, err = checkIpsecVpnExists(d.Get("name").(string), sess, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecVpnExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec vpn %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIpsecVpnReadWJnprSess(d, sess, jnprSess)...)
}

func resourceIpsecVpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceIpsecVpnReadWJnprSess(d, sess, jnprSess)
}

func resourceIpsecVpnReadWJnprSess(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ipsecVpnOptions, err := readIpsecVpn(d.Get("name").(string), sess, jnprSess)
	mutex.Unlock()
	// copy state vpn_monitor.0.source_interface_auto to struct
	if len(ipsecVpnOptions.vpnMonitor) > 0 {
		for _, v := range d.Get("vpn_monitor").([]interface{}) {
			if v != nil {
				stateMonitor := v.(map[string]interface{})
				ipsecVpnOptions.vpnMonitor[0]["source_interface_auto"] = stateMonitor["source_interface_auto"].(bool)
			}
		}
	}
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
	var diagWarns diag.Diagnostics
	if d.HasChanges("bind_interface_auto") {
		diagWarns = append(diagWarns, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       "Modifying bind_interface_auto on resource already created has no effect",
			AttributePath: cty.Path{cty.GetAttrStep{Name: "bind_interface_auto"}},
		})
	}
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delIpsecVpnConf(d, sess, nil); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		if d.HasChanges("bind_interface") && d.Get("bind_interface_auto").(bool) {
			oldInt, _ := d.GetChange("bind_interface")
			if err := sess.configSet([]string{"delete interfaces " + oldInt.(string)}, nil); err != nil {
				return append(diagWarns, diag.FromErr(err)...)
			}
		}
		if err := setIpsecVpn(d, sess, nil); err != nil {
			return append(diagWarns, diag.FromErr(err)...)
		}
		d.Partial(false)

		return diagWarns
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	if err := delIpsecVpnConf(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if d.HasChanges("bind_interface") && d.Get("bind_interface_auto").(bool) {
		oldInt, _ := d.GetChange("bind_interface")
		st0NC, st0Emtpy, _, err := checkInterfaceLogicalNCEmpty(oldInt.(string), sess, jnprSess)
		if err != nil {
			appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if st0NC || st0Emtpy {
			if err := sess.configSet([]string{"delete interfaces " + oldInt.(string)}, jnprSess); err != nil {
				appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

				return append(diagWarns, diag.FromErr(err)...)
			}
		}
	}
	if err := setIpsecVpn(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_security_ipsec_vpn", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIpsecVpnReadWJnprSess(d, sess, jnprSess)...)
}

func resourceIpsecVpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delIpsecVpn(d, sess, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delIpsecVpn(d, sess, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_security_ipsec_vpn", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIpsecVpnImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	ipsecVpnExists, err := checkIpsecVpnExists(d.Id(), sess, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ipsecVpnExists {
		return nil, fmt.Errorf("don't find security ipsec vpn with id '%v' (id must be <name>)", d.Id())
	}
	ipsecVpnOptions, err := readIpsecVpn(d.Id(), sess, jnprSess)
	if err != nil {
		return nil, err
	}
	fillIpsecVpnData(d, ipsecVpnOptions)
	result[0] = d

	return result, nil
}

func checkIpsecVpnExists(ipsecVpn string, sess *Session, jnprSess *NetconfObject) (bool, error) {
	showConfig, err := sess.command(cmdShowConfig+"security ipsec vpn "+ipsecVpn+pipeDisplaySet, jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setIpsecVpn(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
	configSet := make([]string, 0)

	if d.Get("bind_interface").(string) != "" {
		configSet = append(configSet, "set interfaces "+d.Get("bind_interface").(string))
	}
	setPrefix := "set security ipsec vpn " + d.Get("name").(string)
	if d.Get("bind_interface").(string) != "" {
		configSet = append(configSet, setPrefix+" bind-interface "+d.Get("bind_interface").(string))
	}
	if d.Get("df_bit").(string) != "" {
		configSet = append(configSet, setPrefix+" df-bit "+d.Get("df_bit").(string))
	}
	if d.Get("establish_tunnels").(string) != "" {
		configSet = append(configSet, setPrefix+" establish-tunnels "+d.Get("establish_tunnels").(string))
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
	trafficSelectorName := make([]string, 0)
	for _, v := range d.Get("traffic_selector").([]interface{}) {
		tS := v.(map[string]interface{})
		if bchk.StringInSlice(tS["name"].(string), trafficSelectorName) {
			return fmt.Errorf("multiple blocks traffic_selector with the same name %s", tS["name"].(string))
		}
		trafficSelectorName = append(trafficSelectorName, tS["name"].(string))
		configSet = append(configSet, "set security ipsec vpn "+d.Get("name").(string)+" traffic-selector "+
			tS["name"].(string)+" local-ip "+tS["local_ip"].(string))
		configSet = append(configSet, "set security ipsec vpn "+d.Get("name").(string)+" traffic-selector "+
			tS["name"].(string)+" remote-ip "+tS["remote_ip"].(string))
	}
	for _, v := range d.Get("vpn_monitor").([]interface{}) {
		monitor := v.(map[string]interface{})
		configSet = append(configSet, "set security ipsec vpn "+d.Get("name").(string)+" vpn-monitor")
		if monitor["destination_ip"].(string) != "" {
			configSet = append(configSet, setPrefix+" vpn-monitor destination-ip "+
				monitor["destination_ip"].(string))
		}
		if monitor["optimized"].(bool) {
			configSet = append(configSet, setPrefix+" vpn-monitor optimized")
		}
		if monitor["source_interface"].(string) != "" {
			configSet = append(configSet, setPrefix+" vpn-monitor source-interface "+
				monitor["source_interface"].(string))
		}
		if monitor["source_interface_auto"].(bool) {
			configSet = append(configSet, setPrefix+" vpn-monitor source-interface "+
				d.Get("bind_interface").(string))
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readIpsecVpn(ipsecVpn string, sess *Session, jnprSess *NetconfObject) (ipsecVpnOptions, error) {
	var confRead ipsecVpnOptions

	showConfig, err := sess.command(cmdShowConfig+
		"security ipsec vpn "+ipsecVpn+pipeDisplaySetRelative, jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		confRead.name = ipsecVpn
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "bind-interface "):
				confRead.bindInterface = strings.TrimPrefix(itemTrim, "bind-interface ")
			case strings.HasPrefix(itemTrim, "df-bit "):
				confRead.dfBit = strings.TrimPrefix(itemTrim, "df-bit ")
			case strings.HasPrefix(itemTrim, "establish-tunnels "):
				confRead.establishTunnels = strings.TrimPrefix(itemTrim, "establish-tunnels ")
			case strings.HasPrefix(itemTrim, "ike "):
				if len(confRead.ike) == 0 {
					confRead.ike = append(confRead.ike, map[string]interface{}{
						"gateway":          "",
						"policy":           "",
						"identity_local":   "",
						"identity_remote":  "",
						"identity_service": "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "ike gateway "):
					confRead.ike[0]["gateway"] = strings.TrimPrefix(itemTrim, "ike gateway ")
				case strings.HasPrefix(itemTrim, "ike ipsec-policy "):
					confRead.ike[0]["policy"] = strings.TrimPrefix(itemTrim, "ike ipsec-policy ")
				case strings.HasPrefix(itemTrim, "ike proxy-identity local "):
					confRead.ike[0]["identity_local"] = strings.TrimPrefix(itemTrim, "ike proxy-identity local ")
				case strings.HasPrefix(itemTrim, "ike proxy-identity remote "):
					confRead.ike[0]["identity_remote"] = strings.TrimPrefix(itemTrim, "ike proxy-identity remote ")
				case strings.HasPrefix(itemTrim, "ike proxy-identity service "):
					confRead.ike[0]["identity_service"] = strings.TrimPrefix(itemTrim, "ike proxy-identity service ")
				}
			case strings.HasPrefix(itemTrim, "traffic-selector "):
				tsSplit := strings.Split(strings.TrimPrefix(itemTrim, "traffic-selector "), " ")
				tsOptions := map[string]interface{}{
					"name":      tsSplit[0],
					"local_ip":  "",
					"remote_ip": "",
				}
				itemTrimTS := strings.TrimPrefix(itemTrim, "traffic-selector "+tsSplit[0]+" ")
				if len(confRead.trafficSelector) > 0 {
					confRead.trafficSelector = copyAndRemoveItemMapList("name", tsOptions, confRead.trafficSelector)
				}
				switch {
				case strings.HasPrefix(itemTrimTS, "local-ip "):
					tsOptions["local_ip"] = strings.TrimPrefix(itemTrimTS, "local-ip ")
				case strings.HasPrefix(itemTrimTS, "remote-ip "):
					tsOptions["remote_ip"] = strings.TrimPrefix(itemTrimTS, "remote-ip ")
				}
				confRead.trafficSelector = append(confRead.trafficSelector, tsOptions)
			case strings.HasPrefix(itemTrim, "vpn-monitor "):
				if len(confRead.vpnMonitor) == 0 {
					confRead.vpnMonitor = append(confRead.vpnMonitor, map[string]interface{}{
						"destination_ip":   "",
						"optimized":        false,
						"source_interface": "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "vpn-monitor destination-ip "):
					confRead.vpnMonitor[0]["destination_ip"] = strings.TrimPrefix(itemTrim, "vpn-monitor destination-ip ")
				case itemTrim == "vpn-monitor optimized":
					confRead.vpnMonitor[0]["optimized"] = true
				case strings.HasPrefix(itemTrim, "vpn-monitor source-interface "):
					confRead.vpnMonitor[0]["source_interface"] = strings.TrimPrefix(itemTrim, "vpn-monitor source-interface ")
				}
			}
		}
	}

	return confRead, nil
}

func delIpsecVpnConf(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec vpn "+d.Get("name").(string))

	return sess.configSet(configSet, jnprSess)
}

func delIpsecVpn(d *schema.ResourceData, sess *Session, jnprSess *NetconfObject) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec vpn "+d.Get("name").(string))
	if d.Get("bind_interface_auto").(bool) {
		if jnprSess == nil {
			configSet = append(configSet, "delete interfaces "+d.Get("bind_interface").(string))
		} else {
			st0NC, st0Emtpy, _, err := checkInterfaceLogicalNCEmpty(d.Get("bind_interface").(string), sess, jnprSess)
			if err != nil {
				return err
			}
			if st0NC || st0Emtpy {
				configSet = append(configSet, "delete interfaces "+d.Get("bind_interface").(string))
			}
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func fillIpsecVpnData(d *schema.ResourceData, ipsecVpnOptions ipsecVpnOptions) {
	if tfErr := d.Set("name", ipsecVpnOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bind_interface", ipsecVpnOptions.bindInterface); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("df_bit", ipsecVpnOptions.dfBit); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("establish_tunnels", ipsecVpnOptions.establishTunnels); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ike", ipsecVpnOptions.ike); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("traffic_selector", ipsecVpnOptions.trafficSelector); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vpn_monitor", ipsecVpnOptions.vpnMonitor); tfErr != nil {
		panic(tfErr)
	}
}
