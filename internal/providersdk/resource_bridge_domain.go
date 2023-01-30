package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type bridgeDomainOptions struct {
	domainTypeBridge bool
	domainID         int
	isolatedVlan     int
	serviceID        int
	vlanID           int
	description      string
	name             string
	routingInstance  string
	routingInterface string
	communityVlans   []string
	vlanIDList       []string
	vxlan            []map[string]interface{}
}

func resourceBridgeDomain() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceBridgeDomainCreate,
		ReadWithoutTimeout:   resourceBridgeDomainRead,
		UpdateWithoutTimeout: resourceBridgeDomainUpdate,
		DeleteWithoutTimeout: resourceBridgeDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBridgeDomainImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          junos.DefaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"community_vlans": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 15),
			},
			"domain_type_bridge": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"isolated_vlan": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"routing_interface": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !bchk.StringHasOneOfPrefixes(value, []string{"irb.", "vlan."}) {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not start with 'irb.' or 'vlan.'", value, k))
					}

					return
				},
			},
			"service_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"vlan_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 4094),
			},
			"vlan_id_list": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"vlan_id"},
			},
			"vxlan": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vni": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 16777214),
						},
						"vni_extend_evpn": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"decapsulate_accept_inner_vlan": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"encapsulate_inner_vlan": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ingress_node_replication": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"multicast_group": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"ovsdb_managed": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"unreachable_vtep_aging_timer": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(300, 1800),
						},
					},
				},
			},
		},
	}
}

func resourceBridgeDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setBridgeDomain(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string) + junos.IDSeparator + d.Get("routing_instance").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilityRouter() {
		return diag.FromErr(fmt.Errorf("bridge domain "+
			"not compatible with Junos device %s", junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if d.Get("routing_instance").(string) != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	bridgeDomainExists, err := checkBridgeDomainExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if bridgeDomainExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("bridge domain %v already exists in routing_instance %s",
			d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setBridgeDomain(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_bridge_domain")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	bridgeDomainExists, err = checkBridgeDomainExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if bridgeDomainExists {
		d.SetId(d.Get("name").(string) + junos.IDSeparator + d.Get("routing_instance").(string))
	} else {
		return append(diagWarns,
			diag.FromErr(fmt.Errorf("bridge domain %v not exists in routing_instance %v after commit "+
				"=> check your config", d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceBridgeDomainReadWJunSess(d, junSess)...)
}

func resourceBridgeDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceBridgeDomainReadWJunSess(d, junSess)
}

func resourceBridgeDomainReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	bridgeDomainOptions, err := readBridgeDomain(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if bridgeDomainOptions.name == "" {
		d.SetId("")
	} else {
		fillBridgeDomainData(d, bridgeDomainOptions)
	}

	return nil
}

func resourceBridgeDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if d.HasChange("vxlan") {
			oldVxlan, _ := d.GetChange("vxlan")
			if err := delBridgeDomainOpts(
				d.Get("name").(string),
				d.Get("routing_instance").(string),
				oldVxlan.([]interface{}),
				junSess,
			); err != nil {
				return diag.FromErr(err)
			}
		} else if err := delBridgeDomainOpts(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			d.Get("vxlan").([]interface{}),
			junSess,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setBridgeDomain(d, junSess); err != nil {
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
	if d.HasChange("vxlan") {
		oldVxlan, _ := d.GetChange("vxlan")
		if err := delBridgeDomainOpts(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			oldVxlan.([]interface{}),
			junSess,
		); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	} else if err := delBridgeDomainOpts(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("vxlan").([]interface{}),
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setBridgeDomain(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_bridge_domain")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceBridgeDomainReadWJunSess(d, junSess)...)
}

func resourceBridgeDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delBridgeDomain(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			d.Get("vxlan").([]interface{}),
			junSess,
		); err != nil {
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
	if err := delBridgeDomain(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("vxlan").([]interface{}),
		junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_bridge_domain")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceBridgeDomainImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), junos.IDSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	bridgeDomainExists, err := checkBridgeDomainExists(idSplit[0], idSplit[1], junSess)
	if err != nil {
		return nil, err
	}
	if !bridgeDomainExists {
		return nil, fmt.Errorf("don't find bridge domain with id '%v' (id must be "+
			"<name>"+junos.IDSeparator+"<routing_instance>)", d.Id())
	}
	bridgeDomainOptions, err := readBridgeDomain(idSplit[0], idSplit[1], junSess)
	if err != nil {
		return nil, err
	}
	fillBridgeDomainData(d, bridgeDomainOptions)

	result[0] = d

	return result, nil
}

func checkBridgeDomainExists(name, instance string, junSess *junos.Session) (_ bool, err error) {
	var showConfig string
	if instance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"bridge-domains \"" + name + "\"" + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
			"bridge-domains \"" + name + "\"" + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
	}

	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setBridgeDomain(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := junos.SetLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	setPrefix += "bridge-domains \"" + d.Get("name").(string) + "\" "

	for _, v := range sortSetOfString(d.Get("community_vlans").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"community-vlans "+v)
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := d.Get("domain_id").(int); v != 0 {
		configSet = append(configSet, setPrefix+"domain-id "+strconv.Itoa(v))
	}
	if d.Get("domain_type_bridge").(bool) {
		configSet = append(configSet, setPrefix+"domain-type bridge")
	}
	if v := d.Get("isolated_vlan").(int); v != 0 {
		configSet = append(configSet, setPrefix+"isolated-vlan "+strconv.Itoa(v))
	}
	if v := d.Get("routing_interface").(string); v != "" {
		configSet = append(configSet, setPrefix+"routing-interface "+v)
	}
	if v := d.Get("service_id").(int); v != 0 {
		configSet = append(configSet, setPrefix+"service-id "+strconv.Itoa(v))
	}
	if v := d.Get("vlan_id").(int); v != 0 {
		configSet = append(configSet, setPrefix+"vlan-id "+strconv.Itoa(v))
	}
	for _, v := range sortSetOfString(d.Get("vlan_id_list").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"vlan-id-list "+v)
	}
	for _, v := range d.Get("vxlan").([]interface{}) {
		vxlan := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"vxlan vni "+strconv.Itoa(vxlan["vni"].(int)))

		if vxlan["vni_extend_evpn"].(bool) {
			if d.Get("routing_instance").(string) == junos.DefaultW {
				configSet = append(configSet, "set protocols evpn extended-vni-list "+strconv.Itoa(vxlan["vni"].(int)))
			} else {
				configSet = append(configSet, junos.SetRoutingInstances+d.Get("routing_instance").(string)+
					" protocols evpn extended-vni-list "+strconv.Itoa(vxlan["vni"].(int)))
			}
		}
		if vxlan["decapsulate_accept_inner_vlan"].(bool) {
			configSet = append(configSet, setPrefix+"vxlan decapsulate-accept-inner-vlan")
		}
		if vxlan["encapsulate_inner_vlan"].(bool) {
			configSet = append(configSet, setPrefix+"vxlan encapsulate-inner-vlan")
		}
		if vxlan["ingress_node_replication"].(bool) {
			configSet = append(configSet, setPrefix+"vxlan ingress-node-replication")
		}
		if vxlan["multicast_group"].(string) != "" {
			configSet = append(configSet, setPrefix+"vxlan multicast-group "+vxlan["multicast_group"].(string))
		}
		if vxlan["ovsdb_managed"].(bool) {
			configSet = append(configSet, setPrefix+"vxlan ovsdb-managed")
		}
		if vxlan["unreachable_vtep_aging_timer"].(int) != 0 {
			configSet = append(configSet, setPrefix+
				"vxlan unreachable-vtep-aging-timer "+strconv.Itoa(vxlan["unreachable_vtep_aging_timer"].(int)))
		}
	}

	return junSess.ConfigSet(configSet)
}

func readBridgeDomain(name, instance string, junSess *junos.Session,
) (confRead bridgeDomainOptions, err error) {
	var showConfig string
	if instance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"bridge-domains \"" + name + "\"" + junos.PipeDisplaySetRelative)
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
			"bridge-domains \"" + name + "\"" + junos.PipeDisplaySetRelative)
	}
	if err != nil {
		return confRead, err
	}

	if showConfig != junos.EmptyW {
		confRead.name = name
		confRead.routingInstance = instance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "community-vlans "):
				confRead.communityVlans = append(confRead.communityVlans, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "domain-id "):
				confRead.domainID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "domain-type bridge":
				confRead.domainTypeBridge = true
			case balt.CutPrefixInString(&itemTrim, "isolated-vlan "):
				confRead.isolatedVlan, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "routing-interface "):
				confRead.routingInterface = itemTrim
			case balt.CutPrefixInString(&itemTrim, "service-id "):
				confRead.serviceID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "vlan-id "):
				confRead.vlanID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "vlan-id-list "):
				confRead.vlanIDList = append(confRead.vlanIDList, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "vxlan "):
				if len(confRead.vxlan) == 0 {
					confRead.vxlan = append(confRead.vxlan, map[string]interface{}{
						"vni":                           -1,
						"vni_extend_evpn":               false,
						"decapsulate_accept_inner_vlan": false,
						"encapsulate_inner_vlan":        false,
						"ingress_node_replication":      false,
						"multicast_group":               "",
						"ovsdb_managed":                 false,
						"unreachable_vtep_aging_timer":  0,
					})
				}
				vxlan := confRead.vxlan[0]
				switch {
				case balt.CutPrefixInString(&itemTrim, "vni "):
					vxlan["vni"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
					var showConfigEvpn string
					if confRead.routingInstance == junos.DefaultW {
						showConfigEvpn, err = junSess.Command(junos.CmdShowConfig + "protocols evpn" + junos.PipeDisplaySetRelative)
						if err != nil {
							return confRead, err
						}
					} else {
						showConfigEvpn, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
							"protocols evpn" + junos.PipeDisplaySetRelative)
						if err != nil {
							return confRead, err
						}
					}
					if showConfigEvpn != junos.EmptyW {
						for _, itemEvpn := range strings.Split(showConfigEvpn, "\n") {
							if strings.Contains(itemEvpn, junos.XMLStartTagConfigOut) {
								continue
							}
							if strings.Contains(itemEvpn, junos.XMLEndTagConfigOut) {
								break
							}
							if strings.HasPrefix(itemEvpn, junos.SetLS+"extended-vni-list "+strconv.Itoa(vxlan["vni"].(int))) {
								vxlan["vni_extend_evpn"] = true

								break
							}
						}
					}
				case itemTrim == "decapsulate-accept-inner-vlan":
					vxlan["decapsulate_accept_inner_vlan"] = true
				case itemTrim == "encapsulate-inner-vlan":
					vxlan["encapsulate_inner_vlan"] = true
				case itemTrim == "ingress-node-replication":
					vxlan["ingress_node_replication"] = true
				case balt.CutPrefixInString(&itemTrim, "multicast-group "):
					vxlan["multicast_group"] = itemTrim
				case itemTrim == "ovsdb-managed":
					vxlan["ovsdb_managed"] = true
				case balt.CutPrefixInString(&itemTrim, "unreachable-vtep-aging-timer "):
					vxlan["unreachable_vtep_aging_timer"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				}
			}
		}
	}

	return confRead, nil
}

func delBridgeDomainOpts(name, instance string, vxlan []interface{}, junSess *junos.Session) error {
	configSet := make([]string, 0)
	delPrefix := junos.DeleteLS
	if instance != junos.DefaultW {
		delPrefix = junos.DelRoutingInstances + instance + " "
	}
	delPrefix += "bridge-domains \"" + name + "\" "

	configSet = append(configSet,
		delPrefix+"community-vlans",
		delPrefix+"description",
		delPrefix+"domain-id",
		delPrefix+"domain-type",
		delPrefix+"isolated-vlan",
		delPrefix+"routing-interface",
		delPrefix+"service-id",
		delPrefix+"vlan-id",
		delPrefix+"vlan-id-list",
		delPrefix+"vxlan",
	)
	for _, v := range vxlan {
		vxlanParams := v.(map[string]interface{})
		if vxlanParams["vni_extend_evpn"].(bool) {
			if instance == junos.DefaultW {
				configSet = append(configSet, junos.DeleteLS+
					"protocols evpn extended-vni-list "+strconv.Itoa(vxlanParams["vni"].(int)))
			} else {
				configSet = append(configSet, junos.DelRoutingInstances+instance+" "+
					"protocols evpn extended-vni-list "+strconv.Itoa(vxlanParams["vni"].(int)))
			}
		}
	}

	return junSess.ConfigSet(configSet)
}

func delBridgeDomain(name, instance string, vxlan []interface{}, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	if instance == junos.DefaultW {
		configSet = append(configSet, "delete bridge-domains \""+name+"\"")
	} else {
		configSet = append(configSet, junos.DelRoutingInstances+instance+" bridge-domains \""+name+"\"")
	}
	for _, v := range vxlan {
		vxlanParams := v.(map[string]interface{})
		if vxlanParams["vni_extend_evpn"].(bool) {
			if instance == junos.DefaultW {
				configSet = append(configSet, junos.DeleteLS+
					"protocols evpn extended-vni-list "+strconv.Itoa(vxlanParams["vni"].(int)))
			} else {
				configSet = append(configSet, junos.DelRoutingInstances+instance+" "+
					"protocols evpn extended-vni-list "+strconv.Itoa(vxlanParams["vni"].(int)))
			}
		}
	}

	return junSess.ConfigSet(configSet)
}

func fillBridgeDomainData(d *schema.ResourceData, bridgeDomainOptions bridgeDomainOptions) {
	if tfErr := d.Set("name", bridgeDomainOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", bridgeDomainOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("community_vlans", bridgeDomainOptions.communityVlans); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", bridgeDomainOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("domain_id", bridgeDomainOptions.domainID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("domain_type_bridge", bridgeDomainOptions.domainTypeBridge); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("isolated_vlan", bridgeDomainOptions.isolatedVlan); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_interface", bridgeDomainOptions.routingInterface); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("service_id", bridgeDomainOptions.serviceID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_id", bridgeDomainOptions.vlanID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_id_list", bridgeDomainOptions.vlanIDList); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vxlan", bridgeDomainOptions.vxlan); tfErr != nil {
		panic(tfErr)
	}
}
