package junos

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type interfaceOptions struct {
	vlanTagging       bool
	inet              bool
	inet6             bool
	trunk             bool
	vlanNative        int
	aeMinLink         int
	inetMtu           int
	inet6Mtu          int
	inetFilterInput   string
	inetFilterOutput  string
	inet6FilterInput  string
	inet6FilterOutput string
	description       string
	v8023ad           string
	aeLacp            string
	aeLinkSpeed       string
	securityZones     string
	routingInstances  string
	vlanMembers       []string
	inetAddress       []map[string]interface{}
	inet6Address      []map[string]interface{}
}

func resourceInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceInterfaceCreate,
		Read:   resourceInterfaceRead,
		Update: resourceInterfaceUpdate,
		Delete: resourceInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceInterfaceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if strings.Count(value, ".") > 1 {
						errors = append(errors, fmt.Errorf(
							"%q in %q cannot have more of 1 dot", value, k))
					}
					return
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_tagging": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"inet": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"inet6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"inet_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPMaskFunc(),
						},
						"vrrp_group": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"identifier": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateIntRange(1, 255),
									},
									"virtual_address": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"accept_data": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"advertise_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 255),
									},
									"advertisements_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 15),
									},
									"authentication_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"authentication_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(string)
											if value != "md5" && value != "simple" {
												errors = append(errors, fmt.Errorf(
													"%q for %q is not' md5' or 'simple'", value, k))
											}
											return
										},
									},
									"no_accept_data": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_preempt": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"preempt": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"priority": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 255),
									},
									"track_interface": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interface": {
													Type:     schema.TypeString,
													Required: true,
												},
												"priority_cost": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validateIntRange(1, 254),
												},
											},
										},
									},
									"track_route": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"route": {
													Type:     schema.TypeString,
													Required: true,
												},
												"routing_instance": {
													Type:     schema.TypeString,
													Required: true,
												},
												"priority_cost": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validateIntRange(1, 254),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"inet6_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIPMaskFunc(),
						},
						"vrrp_group": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"identifier": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateIntRange(1, 255),
									},
									"virtual_address": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"virtual_link_local_address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateIPFunc(),
									},
									"accept_data": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"advertise_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 255),
									},
									"authentication_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"authentication_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(string)
											if !stringInSlice(value, []string{"md5", "simple"}) {
												errors = append(errors, fmt.Errorf(
													"%q for %q is not' md5' or 'simple'", value, k))
											}
											return
										},
									},
									"no_accept_data": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_preempt": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"preempt": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"priority": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 255),
									},
									"track_interface": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interface": {
													Type:     schema.TypeString,
													Required: true,
												},
												"priority_cost": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validateIntRange(1, 254),
												},
											},
										},
									},
									"track_route": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"route": {
													Type:     schema.TypeString,
													Required: true,
												},
												"routing_instance": {
													Type:     schema.TypeString,
													Required: true,
												},
												"priority_cost": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validateIntRange(1, 254),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"inet_mtu": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 500 || value > 9192 {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not a valid mtu (500-9192)", value, k))
					}
					return
				},
			},
			"inet6_mtu": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 500 || value > 9192 {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not a valid mtu (500-9192)", value, k))
					}
					return
				},
			},
			"inet_filter_input": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"inet_filter_output": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"inet6_filter_input": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"inet6_filter_output": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"ether802_3ad": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !strings.HasPrefix(value, "ae") {
						errors = append(errors, fmt.Errorf(
							"%q in %q isn't an ae interface", value, k))
					}
					return
				},
			},
			"trunk": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vlan_members": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vlan_native": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 4094 {
						errors = append(errors, fmt.Errorf(
							"%q in %q is not in default vlan id (1-4094)", value, k))
					}
					return
				},
			},
			"ae_lacp": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "active" && value != "passive" {
						errors = append(errors, fmt.Errorf(
							"%q is not active or passive", k))
					}
					return
				},
			},
			"ae_link_speed": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					validSpeed := []string{"100m", "1g", "8g", "10g", "40g", "50g", "80g", "100g"}
					if !stringInSlice(value, validSpeed) {
						errors = append(errors, fmt.Errorf(
							"%q in %q is not valid speed", value, k))
					}
					return
				},
			},
			"ae_minimum_links": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"security_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"routing_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
		},
	}
}

func resourceInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	intExists, err := checkInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	if intExists {
		err = checkInterfaceNC(d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)
			return err
		}
		err = delInterfaceElement("apply-groups interface-NC", d, m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)
			return err
		}

	}
	if d.Get("security_zone").(string) != "" {
		zonesExists, err := checkSecurityZonesExists(d.Get("security_zone").(string), m, jnprSess)
		if err != nil {
			return err
		}
		if !zonesExists {
			return fmt.Errorf("security zones %v doesn't exist", d.Get("security_zone").(string))
		}
	}
	if d.Get("routing_instance").(string) != "" {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			return err
		}
		if !instanceExists {
			return fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string))
		}
	}
	err = setInterface(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	intExists, err = checkInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if intExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("interface %v not exists after commit => check your config", d.Get("name").(string))
	}

	return resourceInterfaceRead(d, m)
}
func resourceInterfaceRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	intExists, err := checkInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		mutex.Unlock()
		return err
	}
	if !intExists {
		d.SetId("")
	}
	interfaceOpt, err := readInterface(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	fillInterfaceData(d, interfaceOpt)

	return nil
}
func resourceInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	defer sess.closeSession(jnprSess)
	if err != nil {
		return err
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delInterfaceOpts(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if d.HasChange("ether802_3ad") {
		oAE, nAE := d.GetChange("ether802_3ad")
		if oAE.(string) != "" {
			var newAE string
			if nAE.(string) == "" {
				newAE = "ae-1"
			} else {
				newAE = nAE.(string)
			}
			aggregatedCount, err := aggregatedCountSearchMax(newAE, oAE.(string), m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)
				return err
			}
			if aggregatedCount == "0" {
				err = sess.configSet([]string{"delete chassis aggregated-devices ethernet device-count"}, jnprSess)
				if err != nil {
					sess.configClear(jnprSess)
					return err
				}
			} else {
				err = sess.configSet([]string{"set chassis aggregated-devices ethernet device-count " +
					aggregatedCount + "\n"}, jnprSess)
				if err != nil {
					sess.configClear(jnprSess)
					return err
				}
			}
		}
	}
	if d.HasChange("security_zone") {
		oSecurityZone, nSecurityZone := d.GetChange("security_zone")
		if nSecurityZone.(string) != "" {
			zonesExists, err := checkSecurityZonesExists(nSecurityZone.(string), m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)
				return err
			}
			if !zonesExists {
				sess.configClear(jnprSess)
				return fmt.Errorf("security zones %v doesn't exist", nSecurityZone.(string))
			}
		}
		if oSecurityZone.(string) != "" {
			err = delZoneInterface(oSecurityZone.(string), d, m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)
				return err
			}
		}
	}
	if d.HasChange("routing_instance") {
		oRoutingInstance, nRoutingInstance := d.GetChange("routing_instance")
		if nRoutingInstance.(string) != "" {
			instanceExists, err := checkRoutingInstanceExists(nRoutingInstance.(string), m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)
				return err
			}
			if !instanceExists {
				sess.configClear(jnprSess)
				return fmt.Errorf("routing instance %v doesn't exist", nRoutingInstance.(string))
			}
		}
		if oRoutingInstance.(string) != "" {
			err = delRoutingInstanceInterface(oRoutingInstance.(string), d, m, jnprSess)
			if err != nil {
				sess.configClear(jnprSess)
				return err
			}
		}
	}
	err = setInterface(d, m, jnprSess)
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
	return resourceInterfaceRead(d, m)
}
func resourceInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	defer sess.closeSession(jnprSess)
	if err != nil {
		return err
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delInterface(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	intExists, err := checkInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if intExists {
		err = addInterfaceNC(d.Get("name").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)
			return err
		}
		err = sess.commitConf(jnprSess)
		if err != nil {
			sess.configClear(jnprSess)
			return err
		}
	}

	return nil
}
func resourceInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	intExists, err := checkInterfaceExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !intExists {
		return nil, fmt.Errorf("don't find interface with id '%v' (id must be <name>)", d.Id())
	}
	interfaceOpt, err := readInterface(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	tfErr := d.Set("name", d.Id())
	if tfErr != nil {
		panic(tfErr)
	}
	fillInterfaceData(d, interfaceOpt)

	result[0] = d
	return result, nil

}

func checkInterfaceNC(interFace string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intConfigLines := make([]string, 0)
	intConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return err
	}
	// remove unused lines
	for _, item := range strings.Split(intConfig, "\n") {
		// show parameters root on interface exclude unit parameters (except ethernet-switching)
		if !strings.Contains(interFace, ".") && strings.Contains(item, "unit") &&
			!strings.Contains(item, "ethernet-switching") {
			continue
		}
		if strings.Contains(item, "<configuration-output>") {
			continue
		}
		if strings.Contains(item, "</configuration-output>") {
			break
		}
		intConfigLines = append(intConfigLines, item)
	}
	intConfig = strings.Join(intConfigLines, "\n")
	if sess.junosLogFile != "" {
		logFile(fmt.Sprintf("[intConfig] '%s'", intConfig), sess.junosLogFile)
	}
	if strings.Contains(intConfig, "set apply-groups interface-NC") ||
		strings.Contains(intConfig, "set disable\nset description NC") ||
		intConfig == emptyWord ||
		strings.Count(intConfig, "") <= 2 ||
		intConfig == "\nset \n" {
		return nil
	}
	return fmt.Errorf("interface %s already configured", interFace)
}
func addInterfaceNC(interFace string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := make([]string, 0, 2)
	var setName string
	var err error
	if strings.Contains(interFace, ".") {
		intCut = strings.Split(interFace, ".")
	} else {
		intCut = append(intCut, interFace)
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dot", interFace)
	}
	if intCut[0] == "st0" {
		err = sess.configSet([]string{"set interfaces " + setName + " disable description NC\n"}, jnprSess)
	} else {
		err = sess.configSet([]string{"set interfaces " + setName + " apply-groups interface-NC\n"}, jnprSess)
	}
	if err != nil {
		return err
	}
	return nil
}

func checkInterfaceExists(interFace string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	rpcIntName := "<get-interface-information><interface-name>" + interFace +
		"</interface-name></get-interface-information>"
	reply, err := sess.commandXML(rpcIntName, jnprSess)
	if err != nil {
		return false, err
	}
	if strings.Contains(reply, interFace+" not found") {
		return false, nil
	}
	return true, nil
}
func setInterface(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	var setName string
	intCut := make([]string, 0, 2)
	configSet := make([]string, 0)
	if strings.Contains(d.Get("name").(string), ".") {
		intCut = strings.Split(d.Get("name").(string), ".")
	} else {
		intCut = append(intCut, d.Get("name").(string))
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dot", d.Get("name").(string))
	}
	err := checkResourceInterfaceConfigAndName(len(intCut), d)
	if err != nil {
		return err
	}
	if d.Get("description").(string) != "" {
		configSet = append(configSet, "set interfaces "+setName+" description \""+d.Get("description").(string)+"\"\n")
	}
	if d.Get("vlan_tagging").(bool) {
		configSet = append(configSet, "set interfaces "+setName+" vlan-tagging\n")
	}
	if len(intCut) == 2 && intCut[0] != "st0" && intCut[1] != "0" {
		configSet = append(configSet, "set interfaces "+setName+" vlan-id "+intCut[1]+"\n")
	}
	if d.Get("inet").(bool) {
		configSet = append(configSet, "set interfaces "+setName+" family inet\n")
	}
	if d.Get("inet6").(bool) {
		configSet = append(configSet, "set interfaces "+setName+" family inet6\n")
	}
	for _, address := range d.Get("inet_address").([]interface{}) {
		configSet, err = setFamilyAddress(address, intCut, configSet, setName, inetWord)
		if err != nil {
			return err
		}
	}
	for _, address := range d.Get("inet6_address").([]interface{}) {
		configSet, err = setFamilyAddress(address, intCut, configSet, setName, inet6Word)
		if err != nil {
			return err
		}
	}
	if d.Get("inet_mtu").(int) > 0 {
		configSet = append(configSet, "set interfaces "+setName+" family inet mtu "+
			strconv.Itoa(d.Get("inet_mtu").(int))+"\n")
	}
	if d.Get("inet6_mtu").(int) > 0 {
		configSet = append(configSet, "set interfaces "+setName+" family inet6 mtu "+
			strconv.Itoa(d.Get("inet6_mtu").(int))+"\n")
	}
	if d.Get("inet_filter_input").(string) != "" {
		configSet = append(configSet, "set interfaces "+setName+" family inet filter input "+
			d.Get("inet_filter_input").(string)+"\n")
	}
	if d.Get("inet_filter_output").(string) != "" {
		configSet = append(configSet, "set interfaces "+setName+" family inet filter output "+
			d.Get("inet_filter_output").(string)+"\n")
	}
	if d.Get("inet6_filter_input").(string) != "" {
		configSet = append(configSet, "set interfaces "+setName+" family inet6 filter input "+
			d.Get("inet6_filter_input").(string)+"\n")
	}
	if d.Get("inet6_filter_output").(string) != "" {
		configSet = append(configSet, "set interfaces "+setName+" family inet6 filter output "+
			d.Get("inet6_filter_output").(string)+"\n")
	}
	if d.Get("ether802_3ad").(string) != "" {
		configSet = append(configSet, "set interfaces "+setName+" ether-options 802.3ad "+
			d.Get("ether802_3ad").(string)+"\n")
		configSet = append(configSet, "set interfaces "+setName+" gigether-options 802.3ad "+
			d.Get("ether802_3ad").(string)+"\n")
		oldAE := "ae-1"
		if d.HasChange("ether802_3ad") {
			oldAEtf, _ := d.GetChange("ether802_3ad")
			if oldAEtf.(string) != "" {
				oldAE = oldAEtf.(string)
			}
		}
		aggregatedCount, err := aggregatedCountSearchMax(d.Get("ether802_3ad").(string), oldAE, m, jnprSess)
		if err != nil {
			return err
		}
		configSet = append(configSet, "set chassis aggregated-devices ethernet device-count "+aggregatedCount+"\n")
	}
	if d.Get("trunk").(bool) {
		configSet = append(configSet, "set interfaces "+setName+" unit 0 family ethernet-switching interface-mode trunk\n")
	}
	if len(d.Get("vlan_members").([]interface{})) > 0 {
		for _, v := range d.Get("vlan_members").([]interface{}) {
			configSet = append(configSet, "set interfaces "+setName+
				" unit 0 family ethernet-switching vlan members "+v.(string)+"\n")
		}
	}
	if d.Get("vlan_native").(int) != 0 {
		configSet = append(configSet, "set interfaces  native-vlan-id "+strconv.Itoa(d.Get("vlan_native").(int))+"\n")
	}
	if d.Get("ae_lacp").(string) != "" {
		if !strings.Contains(intCut[0], "ae") {
			return fmt.Errorf("ae_lacp invalid for this interface")
		}
		configSet = append(configSet, "set interfaces "+setName+
			" aggregated-ether-options lacp "+d.Get("ae_lacp").(string)+"\n")
	}
	if d.Get("ae_link_speed").(string) != "" {
		if !strings.Contains(intCut[0], "ae") {
			return fmt.Errorf("ae_link_speed invalid for this interface")
		}
		configSet = append(configSet, "set interfaces "+setName+
			" aggregated-ether-options link-speed "+d.Get("ae_link_speed").(string)+"\n")
	}
	if d.Get("ae_minimum_links").(int) > 0 {
		if !strings.Contains(intCut[0], "ae") {
			return fmt.Errorf("ae_minimum_links invalid for this interface")
		}
		configSet = append(configSet, "set interfaces "+setName+
			" aggregated-ether-options minimum-links "+strconv.Itoa(d.Get("ae_minimum_links").(int))+"\n")
	}
	if d.Get("security_zone").(string) != "" {
		configSet = append(configSet, "set security zones security-zone "+
			d.Get("security_zone").(string)+" interfaces "+d.Get("name").(string)+"\n")
	}
	if d.Get("routing_instance").(string) != "" {
		configSet = append(configSet, "set routing-instances "+d.Get("routing_instance").(string)+
			" interface "+d.Get("name").(string)+"\n")
	}

	err = sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func readInterface(interFace string, m interface{}, jnprSess *NetconfObject) (interfaceOptions, error) {
	sess := m.(*Session)
	var confRead interfaceOptions

	intConfig, err := sess.command("show configuration interfaces "+interFace+" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	inetAddress := make([]map[string]interface{}, 0)
	inet6Address := make([]map[string]interface{}, 0)

	if intConfig != emptyWord {
		for _, item := range strings.Split(intConfig, "\n") {
			if !strings.Contains(interFace, ".") && strings.Contains(item, " unit ") &&
				!strings.Contains(item, "ethernet-switching") {
				continue
			}
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, "set ")
			switch {
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")

			case strings.HasPrefix(itemTrim, "vlan-tagging"):
				confRead.vlanTagging = true
			case strings.HasPrefix(itemTrim, "family inet6"):
				confRead.inet6 = true
				switch {
				case strings.HasPrefix(itemTrim, "family inet6 address "):
					inet6Address, err = fillFamilyInetAddress(itemTrim, inet6Address, inet6Word)
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "family inet6 mtu"):
					confRead.inetMtu, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "family inet6 mtu "))
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "family inet6 filter input "):
					confRead.inet6FilterInput = strings.TrimPrefix(itemTrim, "family inet6 filter input ")
				case strings.HasPrefix(itemTrim, "family inet6 filter output "):
					confRead.inet6FilterOutput = strings.TrimPrefix(itemTrim, "family inet6 filter output ")
				}
			case strings.HasPrefix(itemTrim, "family inet"):
				confRead.inet = true
				switch {
				case strings.HasPrefix(itemTrim, "family inet address "):
					inetAddress, err = fillFamilyInetAddress(itemTrim, inetAddress, inetWord)
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "family inet mtu "):
					confRead.inetMtu, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "family inet mtu "))
					if err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrim, "family inet filter input "):
					confRead.inetFilterInput = strings.TrimPrefix(itemTrim, "family inet filter input ")
				case strings.HasPrefix(itemTrim, "family inet filter output "):
					confRead.inetFilterOutput = strings.TrimPrefix(itemTrim, "family inet filter output ")
				}
			case strings.HasPrefix(itemTrim, "ether-options 802.3ad "):
				confRead.v8023ad = strings.TrimPrefix(itemTrim, "ether-options 802.3ad ")
			case strings.HasPrefix(itemTrim, "gigether-options 802.3ad "):
				confRead.v8023ad = strings.TrimPrefix(itemTrim, "gigether-options 802.3ad ")
			case strings.HasPrefix(itemTrim, "unit 0 family ethernet-switching interface-mode trunk"):
				confRead.trunk = true
			case strings.HasPrefix(itemTrim, "unit 0 family ethernet-switching vlan members"):
				confRead.vlanMembers = append(confRead.vlanMembers, strings.TrimPrefix(itemTrim,
					"unit 0 family ethernet-switching vlan members "))
			case strings.HasPrefix(itemTrim, "native-vlan-id"):
				confRead.vlanNative, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "native-vlan-id "))
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "aggregated-ether-options lacp "):
				confRead.aeLacp = strings.TrimPrefix(itemTrim, "aggregated-ether-options lacp ")
			case strings.HasPrefix(itemTrim, "aggregated-ether-options link-speed "):
				confRead.aeLinkSpeed = strings.TrimPrefix(itemTrim, "aggregated-ether-options link-speed ")
			case strings.HasPrefix(itemTrim, "aggregated-ether-options minimum-links "):
				confRead.aeMinLink, err = strconv.Atoi(strings.TrimPrefix(itemTrim,
					"aggregated-ether-options minimum-links "))
				if err != nil {
					return confRead, err
				}
			default:
				continue
			}
		}
		confRead.inetAddress = inetAddress
		confRead.inet6Address = inet6Address
	}
	zonesConfig, err := sess.command("show configuration security zones | display set", jnprSess)
	if err != nil {
		return confRead, err
	}
	regexpInts := regexp.MustCompile(`.*interfaces ` + interFace + `$`)
	for _, item := range strings.Split(zonesConfig, "\n") {
		intMatch := regexpInts.MatchString(item)
		if intMatch {
			confRead.securityZones = strings.TrimPrefix(strings.TrimSuffix(item, " interfaces "+interFace),
				"set security zones security-zone ")
			break
		}
	}
	routingConfig, err := sess.command("show configuration routing-instances | display set", jnprSess)
	if err != nil {
		return confRead, err
	}
	regexpInt := regexp.MustCompile(`.*interface ` + interFace + `$`)
	for _, item := range strings.Split(routingConfig, "\n") {
		intMatch := regexpInt.MatchString(item)
		if intMatch {
			confRead.routingInstances = strings.TrimPrefix(strings.TrimSuffix(item, " interface "+interFace),
				"set routing-instances ")
			break
		}
	}
	return confRead, nil
}
func delInterface(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := make([]string, 0, 2)
	var setName string
	if strings.Contains(d.Get("name").(string), ".") {
		intCut = strings.Split(d.Get("name").(string), ".")
	} else {
		intCut = append(intCut, d.Get("name").(string))
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dot", d.Get("name").(string))
	}
	err := sess.configSet([]string{"delete interfaces " + setName + "\n"}, jnprSess)
	if err != nil {
		return err
	}
	if strings.Contains(d.Get("name").(string), "st0.") {
		// interface totally delete by resource_security_ipsec_vpn when bind_interface_auto
		// else there is an interface st0.x empty
		err := sess.configSet([]string{"set interfaces " + setName + "\n"}, jnprSess)
		if err != nil {
			return err
		}
	}
	if d.Get("ether802_3ad").(string) != "" {
		aggregatedCount, err := aggregatedCountSearchMax("ae-1", d.Get("ether802_3ad").(string), m, jnprSess)
		if err != nil {
			return err
		}
		if aggregatedCount == "0" {
			err = sess.configSet([]string{"delete chassis aggregated-devices ethernet device-count"}, jnprSess)
			if err != nil {
				return err
			}
		} else {
			err = sess.configSet([]string{"set chassis aggregated-devices ethernet device-count " +
				aggregatedCount + "\n"}, jnprSess)
			if err != nil {
				return err
			}
		}
	}
	if d.Get("security_zone").(string) != "" {
		err = delZoneInterface(d.Get("security_zone").(string), d, m, jnprSess)
		if err != nil {
			return err
		}
	}
	if d.Get("routing_instance").(string) != "" {
		err = delRoutingInstanceInterface(d.Get("routing_instance").(string), d, m, jnprSess)
		if err != nil {
			return err
		}
	}

	return nil
}

func delInterfaceElement(element string, d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := make([]string, 0, 2)
	var setName string
	configSet := make([]string, 0, 1)
	if strings.Contains(d.Get("name").(string), ".") {
		intCut = strings.Split(d.Get("name").(string), ".")
	} else {
		intCut = append(intCut, d.Get("name").(string))
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dot", d.Get("name").(string))
	}
	configSet = append(configSet, "delete interfaces "+setName+" "+element+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func delInterfaceOpts(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	intCut := make([]string, 0, 2)
	var setName string
	configSet := make([]string, 0, 1)
	if strings.Contains(d.Get("name").(string), ".") {
		intCut = strings.Split(d.Get("name").(string), ".")
	} else {
		intCut = append(intCut, d.Get("name").(string))
	}
	switch len(intCut) {
	case 2:
		setName = intCut[0] + " unit " + intCut[1]
	case 1:
		setName = intCut[0]
	default:
		return fmt.Errorf("the name %s contains too dot", d.Get("name").(string))
	}
	delPrefix := "delete interfaces " + setName + " "
	configSet = append(configSet,
		delPrefix+"vlan-tagging\n",
		delPrefix+"family inet\n",
		delPrefix+"family inet6\n",
		delPrefix+"ether-options 802.3ad\n",
		delPrefix+"gigether-options 802.3ad\n",
		delPrefix+"unit 0 family ethernet-switching interface-mode\n",
		delPrefix+"unit 0 family ethernet-switching vlan members\n",
		delPrefix+"native-vlan-id\n",
		delPrefix+"aggregated-ether-options\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func delZoneInterface(zone string, d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security zones security-zone "+zone+" interfaces "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}
func delRoutingInstanceInterface(instance string, d *schema.ResourceData,
	m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete routing-instances "+instance+" interface "+d.Get("name").(string)+"\n")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillInterfaceData(d *schema.ResourceData, interfaceOpt interfaceOptions) {
	tfErr := d.Set("description", interfaceOpt.description)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("vlan_tagging", interfaceOpt.vlanTagging)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet", interfaceOpt.inet)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet6", interfaceOpt.inet6)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet_address", interfaceOpt.inetAddress)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet6_address", interfaceOpt.inet6Address)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet_mtu", interfaceOpt.inetMtu)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet6_mtu", interfaceOpt.inet6Mtu)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet_filter_input", interfaceOpt.inetFilterInput)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet_filter_output", interfaceOpt.inetFilterOutput)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet6_filter_input", interfaceOpt.inet6FilterInput)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("inet6_filter_output", interfaceOpt.inet6FilterOutput)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("ether802_3ad", interfaceOpt.v8023ad)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("trunk", interfaceOpt.trunk)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("vlan_members", interfaceOpt.vlanMembers)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("vlan_native", interfaceOpt.vlanNative)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("ae_lacp", interfaceOpt.aeLacp)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("ae_link_speed", interfaceOpt.aeLinkSpeed)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("ae_minimum_links", interfaceOpt.aeMinLink)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("security_zone", interfaceOpt.securityZones)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", interfaceOpt.routingInstances)
	if tfErr != nil {
		panic(tfErr)
	}
}
func fillFamilyInetAddress(item string, inetAddress []map[string]interface{},
	family string) ([]map[string]interface{}, error) {
	var addressConfig []string
	var itemTrim string
	switch family {
	case inetWord:
		addressConfig = strings.Split(strings.TrimPrefix(item, "family inet address "), " ")
		itemTrim = strings.TrimPrefix(item, "family inet address "+addressConfig[0]+" ")
	case inet6Word:
		addressConfig = strings.Split(strings.TrimPrefix(item, "family inet6 address "), " ")
		itemTrim = strings.TrimPrefix(item, "family inet6 address "+addressConfig[0]+" ")
	}

	m := genFamilyInetAddress(addressConfig[0])
	m, inetAddress = copyAndRemoveItemMapList("address", false, m, inetAddress)

	if strings.HasPrefix(itemTrim, "vrrp-group ") || strings.HasPrefix(itemTrim, "vrrp-inet6-group ") {
		vrrpGroup := genVRRPGroup(family)
		vrrpID, err := strconv.Atoi(addressConfig[2])
		if err != nil {
			return inetAddress, nil
		}
		itemTrimVrrp := strings.TrimPrefix(itemTrim, "vrrp-group "+strconv.Itoa(vrrpID)+" ")
		if strings.HasPrefix(itemTrim, "vrrp-inet6-group ") {
			itemTrimVrrp = strings.TrimPrefix(itemTrim, "vrrp-inet6-group "+strconv.Itoa(vrrpID)+" ")
		}
		vrrpGroup["identifier"] = vrrpID
		vrrpGroup, m["vrrp_group"] = copyAndRemoveItemMapList("identifier", true, vrrpGroup,
			m["vrrp_group"].([]map[string]interface{}))
		switch {
		case strings.HasPrefix(itemTrimVrrp, "virtual-address "):
			vrrpGroup["virtual_address"] = append(vrrpGroup["virtual_address"].([]string),
				strings.TrimPrefix(itemTrimVrrp, "virtual-address "))
		case strings.HasPrefix(itemTrimVrrp, "virtual-inet6-address "):
			vrrpGroup["virtual_address"] = append(vrrpGroup["virtual_address"].([]string),
				strings.TrimPrefix(itemTrimVrrp, "virtual-inet6-address "))
		case strings.HasPrefix(itemTrimVrrp, "virtual-link-local-address "):
			vrrpGroup["virtual_link_local_address"] = strings.TrimPrefix(itemTrimVrrp,
				"virtual-link-local-address ")
		case strings.HasPrefix(itemTrimVrrp, "accept-data"):
			vrrpGroup["accept_data"] = true
		case strings.HasPrefix(itemTrimVrrp, "advertise-interval "):
			vrrpGroup["advertise_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrimVrrp,
				"advertise-interval "))
			if err != nil {
				return inetAddress, err
			}
		case strings.HasPrefix(itemTrimVrrp, "inet6-advertise-interval "):
			vrrpGroup["advertise_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrimVrrp,
				"inet6-advertise-interval "))
			if err != nil {
				return inetAddress, err
			}
		case strings.HasPrefix(itemTrimVrrp, "advertisements-threshold "):
			vrrpGroup["advertisements_threshold"] = strings.TrimPrefix(itemTrimVrrp,
				"advertisements-threshold ")
		case strings.HasPrefix(itemTrimVrrp, "authentication-key "):
			vrrpGroup["authentication_key"] = strings.TrimPrefix(itemTrimVrrp,
				"authentication-key ")
		case strings.HasPrefix(itemTrimVrrp, "authentication-type "):
			vrrpGroup["authentication_type"] = strings.TrimPrefix(itemTrimVrrp,
				"authentication-type ")
		case strings.HasPrefix(itemTrimVrrp, "no-accept-data"):
			vrrpGroup["no_accept_data"] = true
		case strings.HasPrefix(itemTrimVrrp, "no-preempt"):
			vrrpGroup["no_preempt"] = true
		case strings.HasPrefix(itemTrimVrrp, "preempt"):
			vrrpGroup["preempt"] = true
		case strings.HasPrefix(itemTrimVrrp, "priority"):
			vrrpGroup["priority"] = strings.TrimPrefix(itemTrimVrrp,
				"priority ")
		case strings.HasPrefix(itemTrimVrrp, "track interface "):
			vrrpSlit := strings.Split(itemTrimVrrp, " ")
			cost, err := strconv.Atoi(vrrpSlit[len(vrrpSlit)-1])
			if err != nil {
				return inetAddress, err
			}
			trackInt := map[string]interface{}{
				"interface":     vrrpSlit[3],
				"priority_cost": cost,
			}
			vrrpGroup["track_interface"] = append(vrrpGroup["track_interface"].([]map[string]interface{}), trackInt)
		case strings.HasPrefix(itemTrimVrrp, "track route "):
			vrrpSlit := strings.Split(itemTrimVrrp, " ")
			cost, err := strconv.Atoi(vrrpSlit[len(vrrpSlit)-1])
			if err != nil {
				return inetAddress, err
			}
			trackRoute := map[string]interface{}{
				"route":            vrrpSlit[3],
				"routing_instance": vrrpSlit[5],
				"priority_cost":    cost,
			}
			vrrpGroup["track_route"] = append(vrrpGroup["track_route"].([]map[string]interface{}), trackRoute)
		}
		m["vrrp_group"] = append(m["vrrp_group"].([]map[string]interface{}), vrrpGroup)
	}
	inetAddress = append(inetAddress, m)
	return inetAddress, nil
}
func setFamilyAddress(inetAddress interface{}, intCut []string, configSet []string, setName string,
	family string) ([]string, error) {
	if family != inetWord && family != inet6Word {
		return configSet, fmt.Errorf("setFamilyAddress() unknown family %v", family)
	}
	inetAddressMap := inetAddress.(map[string]interface{})
	configSet = append(configSet, "set interfaces "+setName+" family "+family+
		" address "+inetAddressMap["address"].(string)+"\n")
	for _, vrrpGroup := range inetAddressMap["vrrp_group"].([]interface{}) {
		if intCut[0] == "st0" {
			return configSet, fmt.Errorf("vrrp not available on st0")
		}
		vrrpGroupMap := vrrpGroup.(map[string]interface{})
		if vrrpGroupMap["no_preempt"].(bool) && vrrpGroupMap["preempt"].(bool) {
			return configSet, fmt.Errorf("ConflictsWith no_preempt and preempt")
		}
		if vrrpGroupMap["no_accept_data"].(bool) && vrrpGroupMap["accept_data"].(bool) {
			return configSet, fmt.Errorf("ConflictsWith no_accept_data and accept_data")
		}
		var setNameAddVrrp string
		switch family {
		case inetWord:
			setNameAddVrrp = "set interfaces " + setName + " family inet address " + inetAddressMap["address"].(string) +
				" vrrp-group " + strconv.Itoa(vrrpGroupMap["identifier"].(int))
			for _, ip := range vrrpGroupMap["virtual_address"].([]interface{}) {
				err := validateIP(ip.(string))
				if err != nil {
					return configSet, err
				}
				configSet = append(configSet, setNameAddVrrp+" virtual-address "+ip.(string)+"\n")
			}
			if vrrpGroupMap["advertise_interval"].(int) != 0 {
				configSet = append(configSet, setNameAddVrrp+" advertise-interval "+
					strconv.Itoa(vrrpGroupMap["advertise_interval"].(int))+"\n")
			}
			if vrrpGroupMap["advertisements_threshold"].(int) != 0 {
				configSet = append(configSet, setNameAddVrrp+" advertisements-threshold "+
					strconv.Itoa(vrrpGroupMap["advertisements_threshold"].(int))+"\n")
			}
		case inet6Word:
			setNameAddVrrp = "set interfaces " + setName + " family inet6 address " + inetAddressMap["address"].(string) +
				" vrrp-inet6-group " + strconv.Itoa(vrrpGroupMap["identifier"].(int))
			for _, ip := range vrrpGroupMap["virtual_address"].([]interface{}) {
				err := validateIP(ip.(string))
				if err != nil {
					return configSet, err
				}
				configSet = append(configSet, setNameAddVrrp+" virtual-inet6-address "+ip.(string)+"\n")
			}
			configSet = append(configSet, setNameAddVrrp+" virtual-link-local-address "+
				vrrpGroupMap["virtual_link_local_address"].(string)+"\n")
			if vrrpGroupMap["advertise_interval"].(int) != 0 {
				configSet = append(configSet, setNameAddVrrp+" inet6-advertise-interval "+
					strconv.Itoa(vrrpGroupMap["advertise_interval"].(int))+"\n")
			}
		}
		if vrrpGroupMap["accept_data"].(bool) {
			configSet = append(configSet, setNameAddVrrp+" accept-data"+"\n")
		}
		if vrrpGroupMap["authentication_key"].(string) != "" {
			configSet = append(configSet, setNameAddVrrp+" authentication-key "+
				vrrpGroupMap["authentication_key"].(string)+"\n")
		}
		if vrrpGroupMap["authentication_type"].(string) != "" {
			configSet = append(configSet, setNameAddVrrp+" authentication-type "+
				vrrpGroupMap["authentication_type"].(string)+"\n")
		}
		if vrrpGroupMap["no_accept_data"].(bool) {
			configSet = append(configSet, setNameAddVrrp+" no-accept-data"+"\n")
		}
		if vrrpGroupMap["no_preempt"].(bool) {
			configSet = append(configSet, setNameAddVrrp+" no-preempt"+"\n")
		}
		if vrrpGroupMap["preempt"].(bool) {
			configSet = append(configSet, setNameAddVrrp+" preempt"+"\n")
		}
		if vrrpGroupMap["priority"].(int) != 0 {
			configSet = append(configSet, setNameAddVrrp+" priority "+strconv.Itoa(vrrpGroupMap["priority"].(int))+"\n")
		}
		for _, trackInterface := range vrrpGroupMap["track_interface"].([]interface{}) {
			trackInterfaceMap := trackInterface.(map[string]interface{})
			configSet = append(configSet, setNameAddVrrp+" track interface "+trackInterfaceMap["interface"].(string)+
				" priority-cost "+strconv.Itoa(trackInterfaceMap["priority_cost"].(int))+"\n")
		}
		for _, trackRoute := range vrrpGroupMap["track_route"].([]interface{}) {
			trackRouteMap := trackRoute.(map[string]interface{})
			configSet = append(configSet, setNameAddVrrp+" track route "+trackRouteMap["route"].(string)+
				" routing-instance "+trackRouteMap["routing_instance"].(string)+
				" priority-cost "+strconv.Itoa(trackRouteMap["priority_cost"].(int))+"\n")
		}
	}
	return configSet, nil
}

func aggregatedCountSearchMax(newAE string, oldAE string, m interface{}, jnprSess *NetconfObject) (string, error) {
	sess := m.(*Session)
	newAENum := strings.TrimPrefix(newAE, "ae")
	newAENumInt, err := strconv.Atoi(newAENum)
	if err != nil {
		return "", err
	}
	intShowInt, err := sess.command("show interfaces terse", jnprSess)
	if err != nil {
		return "", err
	}

	intShowIntLines := strings.Split(intShowInt, "\n")
	intShowAE := make([]string, 0)
	regexpAE := regexp.MustCompile(`ae\d*\s`)
	for _, line := range intShowIntLines {
		aematch := regexpAE.MatchString(line)
		if aematch {
			wordsLine := strings.Fields(line)
			if wordsLine[0] != oldAE {
				if (len(intShowAE) > 0 && intShowAE[len(intShowAE)-1] != wordsLine[0]) || len(intShowAE) == 0 {
					intShowAE = append(intShowAE, wordsLine[0])
				}
			}
		}
	}
	showConf, err := sess.command("show configuration interfaces | display set", jnprSess)
	if err != nil {
		return "", err
	}
	if strings.Count(showConf, " "+oldAE+"\n") > 1 {
		intShowAE = append(intShowAE, oldAE)
	}
	if len(intShowAE) > 0 {
		lastAeInt, err := strconv.Atoi(strings.TrimPrefix(intShowAE[len(intShowAE)-1], "ae"))
		if err != nil {
			return "", err
		}
		if lastAeInt > newAENumInt {
			return strconv.Itoa(lastAeInt + 1), nil
		}
	}

	return strconv.Itoa(newAENumInt + 1), nil
}
func genFamilyInetAddress(address string) map[string]interface{} {
	return map[string]interface{}{
		"address":    address,
		"vrrp_group": make([]map[string]interface{}, 0),
	}
}
func genVRRPGroup(family string) map[string]interface{} {
	m := map[string]interface{}{
		"identifier":          0,
		"virtual_address":     make([]string, 0),
		"accept_data":         false,
		"advertise_interval":  0,
		"authentication_key":  "",
		"authentication_type": "",
		"no_accept_data":      false,
		"no_preempt":          false,
		"preempt":             false,
		"priority":            0,
		"track_interface":     make([]map[string]interface{}, 0),
		"track_route":         make([]map[string]interface{}, 0),
	}
	if family == inetWord {
		m["advertisements_threshold"] = 0
	}
	if family == inet6Word {
		m["virtual_link_local_address"] = ""
	}
	return m
}
func checkResourceInterfaceConfigAndName(length int, d *schema.ResourceData) error {
	if length == 1 {
		if d.Get("inet").(bool) {
			return fmt.Errorf("inet invalid for this interface")
		}
		if d.Get("inet6").(bool) {
			return fmt.Errorf("inet6 invalid for this interface")
		}
		if len(d.Get("inet_address").([]interface{})) > 0 {
			return fmt.Errorf("inet address invalid for this interface")
		}
		if len(d.Get("inet6_address").([]interface{})) > 0 {
			return fmt.Errorf("inet6 address invalid for this interface")
		}
		if d.Get("inet_mtu").(int) > 0 {
			return fmt.Errorf("inet_mtu invalid for this interface")
		}
		if d.Get("inet6_mtu").(int) > 0 {
			return fmt.Errorf("inet6_mtu invalid for this interface")
		}
		if d.Get("security_zone").(string) != "" {
			return fmt.Errorf("security_zone invalid for this interface")
		}
		if d.Get("routing_instance").(string) != "" {
			return fmt.Errorf("routing_instance invalid for this interface")
		}
	}
	if length == 2 {
		if d.Get("vlan_tagging").(bool) {
			return fmt.Errorf("vlan tagging invalid for this interface")
		}
		if d.Get("ether802_3ad").(string) != "" {
			return fmt.Errorf("ether802_3ad invalid for this interface")
		}
		if d.Get("trunk").(bool) {
			return fmt.Errorf("trunk invalid for this interface (remove .0)")
		}
		if len(d.Get("vlan_members").([]interface{})) > 0 {
			return fmt.Errorf("vlan_members invalid for this interface (remove .0)")
		}
		if d.Get("vlan_native").(int) != 0 {
			return fmt.Errorf("vlan_members invalid for this interface (remove .0)")
		}
		if d.Get("ae_lacp").(string) != "" {
			return fmt.Errorf("ae_lacp invalid for this interface")
		}
		if d.Get("ae_link_speed").(string) != "" {
			return fmt.Errorf("ae_link_speed invalid for this interface")
		}
		if d.Get("ae_minimum_links").(int) > 0 {
			return fmt.Errorf("ae_minimum_links invalid for this interface")
		}
	}
	return nil
}
