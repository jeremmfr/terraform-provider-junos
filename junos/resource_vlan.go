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
		CreateContext: resourceVlanCreate,
		ReadContext:   resourceVlanRead,
		UpdateContext: resourceVlanUpdate,
		DeleteContext: resourceVlanDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVlanImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"community_vlans": {
				Type:     schema.TypeList,
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
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"forward_filter_output": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"forward_flood_input": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
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
					if !strings.HasPrefix(value, "irb.") {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not start with 'irb.'", value, k))
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
				Type:          schema.TypeList,
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	vlanExists, err := checkVlansExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if vlanExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("vlan %v already exists", d.Get("name").(string)))
	}

	if err := setVlan(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_vlan", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	vlanExists, err = checkVlansExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if vlanExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns,
			diag.FromErr(fmt.Errorf("vlan %v not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceVlanReadWJnprSess(d, m, jnprSess)...)
}
func resourceVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceVlanReadWJnprSess(d, m, jnprSess)
}
func resourceVlanReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	vlanOptions, err := readVlan(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delVlan(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setVlan(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_vlan", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceVlanReadWJnprSess(d, m, jnprSess)...)
}
func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delVlan(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_vlan", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceVlanImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	vlanExists, err := checkVlansExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !vlanExists {
		return nil, fmt.Errorf("don't find vlan with id '%v' (id must be <name>)", d.Id())
	}
	vlanOptions, err := readVlan(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillVlanData(d, vlanOptions)

	result[0] = d

	return result, nil
}

func checkVlansExists(vlan string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	vlanConfig, err := sess.command("show configuration vlans "+vlan+" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if vlanConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setVlan(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set vlans " + d.Get("name").(string) + " "
	for _, v := range d.Get("community_vlans").([]interface{}) {
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
	for _, v := range d.Get("vlan_id_list").([]interface{}) {
		configSet = append(configSet, setPrefix+"vlan-id-list "+v.(string))
	}
	for _, v := range d.Get("vxlan").([]interface{}) {
		vxlan := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"vxlan vni "+strconv.Itoa(vxlan["vni"].(int)))

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

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readVlan(vlan string, m interface{}, jnprSess *NetconfObject) (vlanOptions, error) {
	sess := m.(*Session)
	var confRead vlanOptions

	vlanConfig, err := sess.command("show configuration"+
		" vlans "+vlan+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if vlanConfig != emptyWord {
		confRead.name = vlan
		for _, item := range strings.Split(vlanConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "community-vlans "):
				commVlan, err := strconv.Atoi(strings.TrimPrefix(itemTrim, "community-vlans "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
				confRead.communityVlans = append(confRead.communityVlans, commVlan)
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "forwarding-options filter input "):
				confRead.forwardFilterInput = strings.TrimPrefix(itemTrim, "forwarding-options filter input ")
			case strings.HasPrefix(itemTrim, "forwarding-options filter output "):
				confRead.forwardFilterOutput = strings.TrimPrefix(itemTrim, "forwarding-options filter output ")
			case strings.HasPrefix(itemTrim, "forwarding-options flood input "):
				confRead.forwardFloodInput = strings.TrimPrefix(itemTrim, "forwarding-options flood input ")
			case strings.HasPrefix(itemTrim, "isolated-vlan "):
				confRead.isolatedVlan, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "isolated-vlan "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "l3-interface "):
				confRead.l3Interface = strings.TrimPrefix(itemTrim, "l3-interface ")
			case strings.HasPrefix(itemTrim, "private-vlan "):
				confRead.privateVlan = strings.TrimPrefix(itemTrim, "private-vlan ")
			case strings.HasPrefix(itemTrim, "service-id "):
				confRead.serviceID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "service-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "vlan-id "):
				confRead.vlanID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "vlan-id "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "vlan-id-list "):
				confRead.vlanIDList = append(confRead.vlanIDList, strings.TrimPrefix(itemTrim, "vlan-id-list "))
			case strings.HasPrefix(itemTrim, "vxlan "):
				vxlan := map[string]interface{}{
					"vni":                          -1,
					"encapsulate_inner_vlan":       false,
					"ingress_node_replication":     false,
					"multicast_group":              "",
					"ovsdb_managed":                false,
					"unreachable_vtep_aging_timer": 0,
				}
				if len(confRead.vxlan) > 0 {
					for k, v := range confRead.vxlan[0] {
						vxlan[k] = v
					}
				}
				switch {
				case strings.HasPrefix(itemTrim, "vxlan vni "):
					vxlan["vni"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "vxlan vni "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case itemTrim == "vxlan encapsulate-inner-vlan":
					vxlan["encapsulate_inner_vlan"] = true
				case itemTrim == "vxlan ingress-node-replication":
					vxlan["ingress_node_replication"] = true
				case strings.HasPrefix(itemTrim, "vxlan multicast-group "):
					vxlan["multicast_group"] = strings.TrimPrefix(itemTrim, "vxlan multicast-group ")
				case itemTrim == "vxlan ovsdb-managed":
					vxlan["ovsdb_managed"] = true
				case strings.HasPrefix(itemTrim, "vxlan unreachable-vtep-aging-timer "):
					vxlan["unreachable_vtep_aging_timer"], err = strconv.Atoi(strings.TrimPrefix(itemTrim,
						"vxlan unreachable-vtep-aging-timer "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				}
				confRead.vxlan = []map[string]interface{}{vxlan}
			}
		}
	}

	return confRead, nil
}

func delVlan(vlan string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete vlans "+vlan)
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
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
