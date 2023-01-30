package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
				Type:     schema.TypeList,
				Optional: true,
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setIpsecVpn(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security ipsec vpn not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ipsecVpnExists, err := checkIpsecVpnExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecVpnExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec vpn %v already exists", d.Get("name").(string)))...)
	}
	if err := setIpsecVpn(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_ipsec_vpn")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ipsecVpnExists, err = checkIpsecVpnExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ipsecVpnExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security ipsec vpn %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceIpsecVpnReadWJunSess(d, junSess)...)
}

func resourceIpsecVpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceIpsecVpnReadWJunSess(d, junSess)
}

func resourceIpsecVpnReadWJunSess(d *schema.ResourceData, junSess *junos.Session) diag.Diagnostics {
	mutex.Lock()
	ipsecVpnOptions, err := readIpsecVpn(d.Get("name").(string), junSess)
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

	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIpsecVpnConf(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setIpsecVpn(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delIpsecVpnConf(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setIpsecVpn(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_ipsec_vpn")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceIpsecVpnReadWJunSess(d, junSess)...)
}

func resourceIpsecVpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delIpsecVpn(d, junSess); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delIpsecVpn(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_ipsec_vpn")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceIpsecVpnImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	ipsecVpnExists, err := checkIpsecVpnExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !ipsecVpnExists {
		return nil, fmt.Errorf("don't find security ipsec vpn with id '%v' (id must be <name>)", d.Id())
	}
	ipsecVpnOptions, err := readIpsecVpn(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillIpsecVpnData(d, ipsecVpnOptions)
	result[0] = d

	return result, nil
}

func checkIpsecVpnExists(ipsecVpn string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "security ipsec vpn " + ipsecVpn + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setIpsecVpn(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

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
		if bchk.InSlice(tS["name"].(string), trafficSelectorName) {
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

	return junSess.ConfigSet(configSet)
}

func readIpsecVpn(ipsecVpn string, junSess *junos.Session,
) (confRead ipsecVpnOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec vpn " + ipsecVpn + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = ipsecVpn
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "bind-interface "):
				confRead.bindInterface = itemTrim
			case balt.CutPrefixInString(&itemTrim, "df-bit "):
				confRead.dfBit = itemTrim
			case balt.CutPrefixInString(&itemTrim, "establish-tunnels "):
				confRead.establishTunnels = itemTrim
			case balt.CutPrefixInString(&itemTrim, "ike "):
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
				case balt.CutPrefixInString(&itemTrim, "gateway "):
					confRead.ike[0]["gateway"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "ipsec-policy "):
					confRead.ike[0]["policy"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "proxy-identity local "):
					confRead.ike[0]["identity_local"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "proxy-identity remote "):
					confRead.ike[0]["identity_remote"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "proxy-identity service "):
					confRead.ike[0]["identity_service"] = itemTrim
				}
			case balt.CutPrefixInString(&itemTrim, "traffic-selector "):
				itemTrimFields := strings.Split(itemTrim, " ")
				tsOptions := map[string]interface{}{
					"name":      itemTrimFields[0],
					"local_ip":  "",
					"remote_ip": "",
				}
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				if len(confRead.trafficSelector) > 0 {
					confRead.trafficSelector = copyAndRemoveItemMapList("name", tsOptions, confRead.trafficSelector)
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "local-ip "):
					tsOptions["local_ip"] = itemTrim
				case balt.CutPrefixInString(&itemTrim, "remote-ip "):
					tsOptions["remote_ip"] = itemTrim
				}
				confRead.trafficSelector = append(confRead.trafficSelector, tsOptions)
			case balt.CutPrefixInString(&itemTrim, "vpn-monitor "):
				if len(confRead.vpnMonitor) == 0 {
					confRead.vpnMonitor = append(confRead.vpnMonitor, map[string]interface{}{
						"destination_ip":   "",
						"optimized":        false,
						"source_interface": "",
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "destination-ip "):
					confRead.vpnMonitor[0]["destination_ip"] = itemTrim
				case itemTrim == "optimized":
					confRead.vpnMonitor[0]["optimized"] = true
				case balt.CutPrefixInString(&itemTrim, "source-interface "):
					confRead.vpnMonitor[0]["source_interface"] = itemTrim
				}
			}
		}
	}

	return confRead, nil
}

func delIpsecVpnConf(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec vpn "+d.Get("name").(string))

	return junSess.ConfigSet(configSet)
}

func delIpsecVpn(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security ipsec vpn "+d.Get("name").(string))

	return junSess.ConfigSet(configSet)
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
