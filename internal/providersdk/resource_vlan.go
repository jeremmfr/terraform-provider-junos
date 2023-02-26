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

type vlanOptions struct {
	isolatedVlan        int
	serviceID           int
	vlanID              int
	name                string
	description         string
	forwardFilterInput  string
	forwardFilterOutput string
	forwardFloodInput   string
	l3Interface         string
	privateVlan         string
	communityVlans      []int
	vlanIDList          []string
	vxlan               []map[string]interface{}
}

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVlanCreate,
		ReadWithoutTimeout:   resourceVlanRead,
		UpdateWithoutTimeout: resourceVlanUpdate,
		DeleteWithoutTimeout: resourceVlanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVlanImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"community_vlans": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forward_filter_input": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"forward_filter_output": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"forward_flood_input": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"isolated_vlan": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"l3_interface": {
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
			"private_vlan": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"community", "isolated"}, false),
			},
			"service_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"vlan_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(1, 4094),
				ConflictsWith: []string{"vlan_id_list"},
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

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setVlan(d, junSess); err != nil {
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
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	vlanExists, err := checkVlansExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if vlanExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("vlan %v already exists", d.Get("name").(string)))...)
	}

	if err := setVlan(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_vlan")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	vlanExists, err = checkVlansExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if vlanExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns,
			diag.FromErr(fmt.Errorf("vlan %v not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceVlanReadWJunSess(d, junSess)...)
}

func resourceVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceVlanReadWJunSess(d, junSess)
}

func resourceVlanReadWJunSess(d *schema.ResourceData, junSess *junos.Session) diag.Diagnostics {
	junos.MutexLock()
	vlanOptions, err := readVlan(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if vlanOptions.name == "" {
		d.SetId("")
	} else {
		fillVlanData(d, vlanOptions)
	}

	return nil
}

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if d.HasChange("vxlan") {
			oldVxlan, _ := d.GetChange("vxlan")
			if err := delVlan(d.Get("name").(string), oldVxlan.([]interface{}), junSess); err != nil {
				return diag.FromErr(err)
			}
		} else if err := delVlan(d.Get("name").(string), d.Get("vxlan").([]interface{}), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setVlan(d, junSess); err != nil {
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
		if err := delVlan(d.Get("name").(string), oldVxlan.([]interface{}), junSess); err != nil {
			appendDiagWarns(&diagWarns, junSess.ConfigClear())

			return append(diagWarns, diag.FromErr(err)...)
		}
	} else if err := delVlan(d.Get("name").(string), d.Get("vxlan").([]interface{}), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setVlan(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_vlan")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVlanReadWJunSess(d, junSess)...)
}

func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delVlan(d.Get("name").(string), d.Get("vxlan").([]interface{}), junSess); err != nil {
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
	if err := delVlan(d.Get("name").(string), d.Get("vxlan").([]interface{}), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_vlan")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceVlanImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	vlanExists, err := checkVlansExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !vlanExists {
		return nil, fmt.Errorf("don't find vlan with id '%v' (id must be <name>)", d.Id())
	}
	vlanOptions, err := readVlan(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillVlanData(d, vlanOptions)

	result[0] = d

	return result, nil
}

func checkVlansExists(vlan string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "vlans " + vlan + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setVlan(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set vlans " + d.Get("name").(string) + " "
	for _, v := range d.Get("community_vlans").(*schema.Set).List() {
		configSet = append(configSet, setPrefix+"community-vlans "+strconv.Itoa(v.(int)))
	}
	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+"description \""+d.Get("description").(string)+"\"")
	}
	if d.Get("forward_filter_input").(string) != "" {
		configSet = append(configSet, setPrefix+
			"forwarding-options filter input "+d.Get("forward_filter_input").(string))
	}
	if d.Get("forward_filter_output").(string) != "" {
		configSet = append(configSet, setPrefix+
			"forwarding-options filter output "+d.Get("forward_filter_output").(string))
	}
	if d.Get("forward_flood_input").(string) != "" {
		configSet = append(configSet, setPrefix+
			"forwarding-options flood input "+d.Get("forward_flood_input").(string))
	}
	if d.Get("isolated_vlan").(int) != 0 {
		configSet = append(configSet, setPrefix+"isolated-vlan "+strconv.Itoa(d.Get("isolated_vlan").(int)))
	}
	if d.Get("l3_interface").(string) != "" {
		configSet = append(configSet, setPrefix+"l3-interface "+d.Get("l3_interface").(string))
	}
	if d.Get("private_vlan").(string) != "" {
		configSet = append(configSet, setPrefix+"private-vlan "+d.Get("private_vlan").(string))
	}
	if d.Get("service_id").(int) != 0 {
		configSet = append(configSet, setPrefix+"service-id "+strconv.Itoa(d.Get("service_id").(int)))
	}
	if d.Get("vlan_id").(int) != 0 {
		configSet = append(configSet, setPrefix+"vlan-id "+strconv.Itoa(d.Get("vlan_id").(int)))
	}
	for _, v := range sortSetOfString(d.Get("vlan_id_list").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"vlan-id-list "+v)
	}
	for _, v := range d.Get("vxlan").([]interface{}) {
		vxlan := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"vxlan vni "+strconv.Itoa(vxlan["vni"].(int)))

		if vxlan["vni_extend_evpn"].(bool) {
			configSet = append(configSet, "set protocols evpn extended-vni-list "+strconv.Itoa(vxlan["vni"].(int)))
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

func readVlan(vlan string, junSess *junos.Session,
) (confRead vlanOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "vlans " + vlan + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = vlan
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
				commVlan, err := strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
				confRead.communityVlans = append(confRead.communityVlans, commVlan)
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "forwarding-options filter input "):
				confRead.forwardFilterInput = itemTrim
			case balt.CutPrefixInString(&itemTrim, "forwarding-options filter output "):
				confRead.forwardFilterOutput = itemTrim
			case balt.CutPrefixInString(&itemTrim, "forwarding-options flood input "):
				confRead.forwardFloodInput = itemTrim
			case balt.CutPrefixInString(&itemTrim, "isolated-vlan "):
				confRead.isolatedVlan, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "l3-interface "):
				confRead.l3Interface = itemTrim
			case balt.CutPrefixInString(&itemTrim, "private-vlan "):
				confRead.privateVlan = itemTrim
			case balt.CutPrefixInString(&itemTrim, "service-id "):
				confRead.serviceID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "vlan-id "):
				confRead.vlanID, err = strconv.Atoi(itemTrim)
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "vlan-id-list "):
				confRead.vlanIDList = append(confRead.vlanIDList, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "vxlan "):
				if len(confRead.vxlan) == 0 {
					confRead.vxlan = append(confRead.vxlan, map[string]interface{}{
						"vni":                          -1,
						"vni_extend_evpn":              false,
						"encapsulate_inner_vlan":       false,
						"ingress_node_replication":     false,
						"multicast_group":              "",
						"ovsdb_managed":                false,
						"unreachable_vtep_aging_timer": 0,
					})
				}
				vxlan := confRead.vxlan[0]
				switch {
				case balt.CutPrefixInString(&itemTrim, "vni "):
					vxlan["vni"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
					showConfigEvpn, err := junSess.Command(junos.CmdShowConfig + "protocols evpn" + junos.PipeDisplaySetRelative)
					if err != nil {
						return confRead, err
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
							}
						}
					}
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
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			}
		}
	}

	return confRead, nil
}

func delVlan(vlan string, vxlan []interface{}, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete vlans "+vlan)
	for _, v := range vxlan {
		vxlanParams := v.(map[string]interface{})
		if vxlanParams["vni_extend_evpn"].(bool) {
			configSet = append(configSet, "delete protocols evpn extended-vni-list "+strconv.Itoa(vxlanParams["vni"].(int)))
		}
	}

	return junSess.ConfigSet(configSet)
}

func fillVlanData(d *schema.ResourceData, vlanOptions vlanOptions) {
	if tfErr := d.Set("name", vlanOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("community_vlans", vlanOptions.communityVlans); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", vlanOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forward_filter_input", vlanOptions.forwardFilterInput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forward_filter_output", vlanOptions.forwardFilterOutput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forward_flood_input", vlanOptions.forwardFloodInput); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("l3_interface", vlanOptions.l3Interface); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("isolated_vlan", vlanOptions.isolatedVlan); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("private_vlan", vlanOptions.privateVlan); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("service_id", vlanOptions.serviceID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_id", vlanOptions.vlanID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vlan_id_list", vlanOptions.vlanIDList); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("vxlan", vlanOptions.vxlan); tfErr != nil {
		panic(tfErr)
	}
}
