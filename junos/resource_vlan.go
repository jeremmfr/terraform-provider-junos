package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type vlanOptions struct {
	vlanID              int
	serviceID           int
	isolatedVlan        int
	name                string
	description         string
	forwardFilterInput  string
	forwardFilterOutput string
	forwardFloodInput   string
	privateVlan         string
	l3Interface         string
	communityVlans      []int
	vlanIDList          []string
	vxlan               []map[string]interface{}
}

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceVlanCreate,
		Read:   resourceVlanRead,
		Update: resourceVlanUpdate,
		Delete: resourceVlanDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVlanImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validateIntRange(1, 4094),
				ConflictsWith: []string{"vlan_id_list"},
			},
			"vlan_id_list": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"vlan_id"},
			},
			"service_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 65535),
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
			"forward_filter_input": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"forward_filter_output": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"forward_flood_input": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"private_vlan": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "community" && value != "isolated" {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not 'community' or 'isolated'", value, k))
					}
					return
				},
			},
			"community_vlans": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"isolated_vlan": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 65535),
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
							ValidateFunc: validateIntRange(0, 16777214),
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
							ValidateFunc: validateIPFunc(),
						},
						"ovsdb_managed": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"unreachable_vtep_aging_timer": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(300, 1800),
						},
					},
				},
			},
		},
	}
}

func resourceVlanCreate(d *schema.ResourceData, m interface{}) error {
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
	vlanExists, err := checkVlansExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if vlanExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("vlan %v already exists", d.Get("name").(string))
	}

	err = setVlan(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	mutex.Lock()
	vlanExists, err = checkVlansExists(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if vlanExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("vlan%v not exists after commit => check your config", d.Get("name").(string))
	}
	return resourceVlanRead(d, m)
}
func resourceVlanRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	vlanOptions, err := readVlan(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if vlanOptions.name == "" {
		d.SetId("")
	} else {
		fillVlanData(d, vlanOptions)
	}
	return nil
}
func resourceVlanUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delVlan(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setVlan(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	d.Partial(false)
	return resourceVlanRead(d, m)
}
func resourceVlanDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delVlan(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	return nil
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
	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+"description \""+d.Get("description").(string)+"\"\n")
	}
	if d.Get("vland_id").(int) != 0 {
		configSet = append(configSet, setPrefix+"vlan-id "+strconv.Itoa(d.Get("vlan_id").(int))+"\n")
	}
	for _, v := range d.Get("vland_id_list").([]interface{}) {
		configSet = append(configSet, setPrefix+"vlan-id-list "+v.(string)+"\n")
	}
	if d.Get("service_id").(int) != 0 {
		configSet = append(configSet, setPrefix+"service_id "+strconv.Itoa(d.Get("service_id").(int))+"\n")
	}
	if d.Get("l3_interface").(string) != "" {
		configSet = append(configSet, setPrefix+"l3-interface "+d.Get("l3_interface").(string)+"\n")
	}
	if d.Get("forward_filter_input").(string) != "" {
		configSet = append(configSet, setPrefix+"forwarding-options filter input "+d.Get("forward_filter_input").(string)+"\n")
	}
	if d.Get("forward_filter_output").(string) != "" {
		configSet = append(configSet, setPrefix+"forwarding-options filter output "+d.Get("forward_filter_output").(string)+"\n")
	}
	if d.Get("forward_flood_input").(string) != "" {
		configSet = append(configSet, setPrefix+"forwarding-options flood input "+d.Get("forward_flood_input").(string)+"\n")
	}
	if d.Get("private_vlan").(string) != "" {
		configSet = append(configSet, setPrefix+"private-vlan "+d.Get("private_vlan").(string)+"\n")
	}
	for _, v := range d.Get("community_vlans").([]interface{}) {
		configSet = append(configSet, setPrefix+"community-vlans "+strconv.Itoa(v.(int))+"\n")
	}
	if d.Get("isolated_vlan").(int) != 0 {
		configSet = append(configSet, setPrefix+"isolated-vlan "+strconv.Itoa(d.Get("isolated_vlan").(int))+"\n")
	}
	for _, v := range d.Get("vxlan").([]interface{}) {
		vxlan := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"vxlan vni "+strconv.Itoa(vxlan["vni"].(int))+"\n")

		if vxlan["encapsulate_inner_vlan"].(bool) {
			configSet = append(configSet, setPrefix+"vxlan encapsulate-inner-vlan\n")
		}
		if vxlan["ingress_node_replication"].(bool) {
			configSet = append(configSet, setPrefix+"vxlan ingress-node-replication\n")
		}
		if vxlan["multicast_group"].(string) != "" {
			configSet = append(configSet, setPrefix+"vxlan multicast-group "+vxlan["multicast_group"].(string)+"\n")
		}
		if vxlan["ovsdb_managed"].(bool) {
			configSet = append(configSet, setPrefix+"vxlan ovsdb-managed\n")
		}
		if vxlan["unreachable_vtep_aging_timer"].(int) != 0 {
			configSet = append(configSet, setPrefix+"vxlan unreachable-vtep-aging-timer"+strconv.Itoa(vxlan["unreachable_vtep_aging_timer"].(int))+"\n")
		}
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
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
			case strings.HasPrefix(itemTrim, "vlan-id "):
				confRead.vlanID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "vlan-id "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "vlan-id-list "):
				confRead.vlanIDList = append(confRead.vlanIDList, strings.TrimPrefix(itemTrim, "vlan-id-list "))
			case strings.HasPrefix(itemTrim, "service-id "):
				confRead.serviceID, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "service-id "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "l3-interface "):
				confRead.l3Interface = strings.TrimPrefix(itemTrim, "l3-interface ")
			case strings.HasPrefix(itemTrim, "forwarding-options filter input "):
				confRead.forwardFilterInput = strings.TrimPrefix(itemTrim, "forwarding-options filter input ")
			case strings.HasPrefix(itemTrim, "forwarding-options filter output "):
				confRead.forwardFilterOutput = strings.TrimPrefix(itemTrim, "forwarding-options filter output ")
			case strings.HasPrefix(itemTrim, "forwarding-options flood input "):
				confRead.forwardFloodInput = strings.TrimPrefix(itemTrim, "forwarding-options flood input ")
			case strings.HasPrefix(itemTrim, "private-vlan "):
				confRead.privateVlan = strings.TrimPrefix(itemTrim, "private-vlan ")
			case strings.HasPrefix(itemTrim, "community-vlans "):
				commVlan, err := strconv.Atoi(strings.TrimPrefix(itemTrim, "community-vlans "))
				if err != nil {
					return confRead, err
				}
				confRead.communityVlans = append(confRead.communityVlans, commVlan)
			case strings.HasPrefix(itemTrim, "isolated-vlan "):
				confRead.isolatedVlan, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "isolated-vlan "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "vxlan "):
				vxlan := map[string]interface{}{
					"vni":                          -1,
					"encapsulate_inner_vlan":       false,
					"ingress_node_replication":     false,
					"multicast_group":              "",
					"ovsdb_managed":                false,
					"unreachable_vtep_aging_timer": 0,
				}
				switch {
				case strings.HasPrefix(itemTrim, "vxlan vni "):
					vxlan["vni"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "vxlan vni "))
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "vxlan encapsulate-inner-vlan"):
					vxlan["encapsulate_inner_vlan"] = true
				case strings.HasPrefix(itemTrim, "vxlan ingress-node-replication"):
					vxlan["ingress_node_replication"] = true
				case strings.HasPrefix(itemTrim, "vxlan multicast-group "):
					vxlan["multicast_group"] = strings.TrimPrefix(itemTrim, "vxlan multicast-group ")
				case strings.HasPrefix(itemTrim, "vxlan ovsdb-managed"):
					vxlan["ovsdb_managed"] = true
				case strings.HasPrefix(itemTrim, "vxlan unreachable-vtep-aging-timer "):
					vxlan["unreachable_vtep_aging_timer"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "vxlan unreachable-vtep-aging-timer "))
					if err != nil {
						return confRead, err
					}
				}
				confRead.vxlan = []map[string]interface{}{vxlan}
			}
		}
	} else {
		confRead.name = ""
		return confRead, nil
	}
	return confRead, nil
}

func delVlan(vlan string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete vlans "+vlan+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillVlanData(d *schema.ResourceData, vlanOptions vlanOptions) {
	tfErr := d.Set("name", vlanOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("description", vlanOptions.description)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("vlan_id", vlanOptions.vlanID)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("vlan_id_list", vlanOptions.vlanIDList)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("service_id", vlanOptions.serviceID)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("l3_interface", vlanOptions.l3Interface)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("forward_filter_input", vlanOptions.forwardFilterInput)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("forward_filter_output", vlanOptions.forwardFilterOutput)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("forward_flood_input", vlanOptions.forwardFloodInput)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("private_vlan", vlanOptions.privateVlan)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("community_vlans", vlanOptions.communityVlans)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("isolated_vlan", vlanOptions.isolatedVlan)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("vxlan", vlanOptions.vxlan)
	if tfErr != nil {
		panic(tfErr)
	}
}
