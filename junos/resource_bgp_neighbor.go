package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Create: resourceBgpNeighborCreate,
		Read:   resourceBgpNeighborRead,
		Update: resourceBgpNeighborUpdate,
		Delete: resourceBgpNeighborDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBgpNeighborImport,
		},
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIPFunc(),
			},
			"routing_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      defaultWord,
				ValidateFunc: validateNameObjectJunos(),
			},
			"group": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"accept_remote_nexthop": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"advertise_external": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"advertise_external_conditional"},
			},
			"advertise_external_conditional": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"advertise_external"},
			},
			"advertise_inactive": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"advertise_peer_as": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"as_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"damping": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"log_updown": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mtu_discovery": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"multihop": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"multipath": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_advertise_peer_as": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"remove_private": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"passive": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hold_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(3, 65535),
			},
			"local_as": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_as_private": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"local_as_alias", "local_as_no_prepend_global_as"},
			},
			"local_as_alias": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"local_as_private", "local_as_no_prepend_global_as"},
			},
			"local_as_no_prepend_global_as": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"local_as_private", "local_as_alias"},
			},
			"local_as_loops": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 10),
			},
			"local_preference": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(0, 4294967295),
				Default:      -1,
			},
			"metric_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(0, 4294967295),
				Default:      -1,
				ConflictsWith: []string{"metric_out_igp",
					"metric_out_igp_offset",
					"metric_out_igp_delay_med_update",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset"},
			},
			"metric_out_igp": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"metric_out",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset"},
			},
			"metric_out_igp_offset": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(-2147483648, 2147483647),
				ConflictsWith: []string{"metric_out",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset"},
			},
			"metric_out_igp_delay_med_update": {
				Type:     schema.TypeBool,
				Optional: true,
				ConflictsWith: []string{"metric_out",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset"},
			},
			"metric_out_minimum_igp": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"metric_out",
					"metric_out_igp",
					"metric_out_igp_offset",
					"metric_out_igp_delay_med_update"},
			},
			"metric_out_minimum_igp_offset": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(-2147483648, 2147483647),
				ConflictsWith: []string{"metric_out",
					"metric_out_igp",
					"metric_out_igp_offset",
					"metric_out_igp_delay_med_update"},
			},
			"out_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 65535),
			},
			"peer_as": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"preference": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(0, 4294967295),
				Default:      -1,
			},
			"authentication_algorithm": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"authentication_key"},
			},
			"authentication_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"authentication_algorithm", "authentication_key_chain"},
			},
			"authentication_key_chain": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"authentication_key"},
			},
			"local_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIPFunc(),
			},
			"local_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bfd_liveness_detection": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_key_chain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authentication_algorithm": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authentication_loose_check": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"detection_time_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 4294967295),
						},
						"transmit_interval_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 4294967295),
						},
						"transmit_interval_minimum_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 255000),
						},
						"holddown_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 255000),
						},
						"minimum_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 255000),
						},
						"minimum_receive_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 255000),
						},
						"multiplier": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 255),
						},
						"session_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"automatic", "multihop", "single-hop"}) {
									errors = append(errors, fmt.Errorf(
										"%q for %q is not 'automatic', 'multihop' or 'single-hop'", value, k))
								}
								return
							},
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"family_inet": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nlri_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"any", "flow", "labeled-unicast", "unicast", "multicast"}) {
									errors = append(errors, fmt.Errorf(
										"%q for %q is not valid nlri type", value, k))
								}
								return
							},
						},
						"accepted_prefix_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"maximum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateIntRange(1, 4294967295),
									},
									"teardown": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 100),
									},
									"teardown_idle_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 2400),
									},
									"teardown_idle_timeout_forever": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"prefix_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"maximum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateIntRange(1, 4294967295),
									},
									"teardown": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 100),
									},
									"teardown_idle_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 2400),
									},
									"teardown_idle_timeout_forever": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"family_inet6": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nlri_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"any", "flow", "labeled-unicast", "unicast", "multicast"}) {
									errors = append(errors, fmt.Errorf(
										"%q for %q is not valid nlri type", value, k))
								}
								return
							},
						},
						"accepted_prefix_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"maximum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateIntRange(1, 4294967295),
									},
									"teardown": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 100),
									},
									"teardown_idle_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 2400),
									},
									"teardown_idle_timeout_forever": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"prefix_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"maximum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateIntRange(1, 4294967295),
									},
									"teardown": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 100),
									},
									"teardown_idle_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntRange(1, 2400),
									},
									"teardown_idle_timeout_forever": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"graceful_restart": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"restart_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 600),
						},
						"stale_route_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntRange(1, 600),
						},
					},
				},
			},
		},
	}
}

func resourceBgpNeighborCreate(d *schema.ResourceData, m interface{}) error {
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
	if d.Get("routing_instance").(string) != defaultWord {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)
			return err
		}
		if !instanceExists {
			sess.configClear(jnprSess)
			return fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string))
		}
	}
	bgpGroupExists, err := checkBgpGroupExists(d.Get("group").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if !bgpGroupExists {
		sess.configClear(jnprSess)
		return fmt.Errorf("bgp group %v doesn't exist", d.Get("group").(string))
	}
	bgpNeighborxists, err := checkBgpNeighborExists(d.Get("ip").(string),
		d.Get("routing_instance").(string), d.Get("group").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	if bgpNeighborxists {
		sess.configClear(jnprSess)
		return fmt.Errorf("bgp neighbor %v already exists in group %v (routing-instance %v)",
			d.Get("ip").(string), d.Get("group").(string), d.Get("routing_instance").(string))
	}
	err = setBgpNeighbor(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = sess.commitConf(jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	bgpNeighborxists, err = checkBgpNeighborExists(d.Get("ip").(string),
		d.Get("routing_instance").(string), d.Get("group").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if bgpNeighborxists {
		d.SetId(d.Get("ip").(string) +
			idSeparator + d.Get("routing_instance").(string) +
			idSeparator + d.Get("group").(string))
	} else {
		return fmt.Errorf("bgp neighbor %v not exists in group %v (routing-instance %v) after commit "+
			"=> check your config", d.Get("ip").(string), d.Get("group").(string), d.Get("routing_instance").(string))
	}
	return resourceBgpNeighborRead(d, m)
}
func resourceBgpNeighborRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()
		return err
	}
	defer sess.closeSession(jnprSess)
	bgpNeighborOptions, err := readBgpNeighbor(d.Get("ip").(string),
		d.Get("routing_instance").(string), d.Get("group").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if bgpNeighborOptions.ip == "" {
		d.SetId("")
	} else {
		fillBgpNeighborData(d, bgpNeighborOptions)
	}
	return nil
}
func resourceBgpNeighborUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delBgpOpts(d, "neighbor", m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)
		return err
	}
	err = setBgpNeighbor(d, m, jnprSess)
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
	return resourceBgpNeighborRead(d, m)
}
func resourceBgpNeighborDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delBgpNeighbor(d, m, jnprSess)
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
func resourceBgpNeighborImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	bgpNeighborxists, err := checkBgpNeighborExists(idSplit[0], idSplit[1], idSplit[2], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !bgpNeighborxists {
		return nil, fmt.Errorf("don't find bgp neighbor with id '%v' "+
			"(id must be <ip>"+idSeparator+"<routing_instace>"+idSeparator+"<group>)", d.Id())
	}
	bgpNeighborOptions, err := readBgpNeighbor(idSplit[0], idSplit[1], idSplit[2], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillBgpNeighborData(d, bgpNeighborOptions)
	result[0] = d
	return result, nil
}

func checkBgpNeighborExists(ip, instance, group string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	var bgpNeighborConfig string
	var err error
	if instance == defaultWord {
		bgpNeighborConfig, err = sess.command("show configuration protocols bgp group "+
			group+" neighbor "+ip+" | display set", jnprSess)
		if err != nil {
			return false, err
		}
	} else {
		bgpNeighborConfig, err = sess.command("show configuration routing-instances "+
			instance+" protocols bgp group "+group+" neighbor "+ip+" | display set", jnprSess)
		if err != nil {
			return false, err
		}
	}
	if bgpNeighborConfig == emptyWord {
		return false, nil
	}
	return true, nil
}
func setBgpNeighbor(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	setPrefix := setLineStart
	if d.Get("routing_instance").(string) == defaultWord {
		setPrefix += "protocols bgp group " + d.Get("group").(string) +
			" neighbor " + d.Get("ip").(string) + " "
	} else {
		setPrefix += "routing-instances " + d.Get("routing_instance").(string) +
			" protocols bgp group " + d.Get("group").(string) +
			" neighbor " + d.Get("ip").(string) + " "
	}
	err := setBgpOptsSimple(setPrefix, d, m, jnprSess)
	if err != nil {
		return err
	}
	err = setBgpOptsBfd(setPrefix, d.Get("bfd_liveness_detection").([]interface{}), m, jnprSess)
	if err != nil {
		return err
	}
	err = setBgpOptsFamily(setPrefix, inetWord, d.Get("family_inet").([]interface{}), m, jnprSess)
	if err != nil {
		return err
	}
	err = setBgpOptsFamily(setPrefix, inet6Word, d.Get("family_inet6").([]interface{}), m, jnprSess)
	if err != nil {
		return err
	}
	err = setBgpOptsGrafefulRestart(setPrefix, d.Get("graceful_restart").([]interface{}), m, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func readBgpNeighbor(ip, instance, group string, m interface{}, jnprSess *NetconfObject) (bgpOptions, error) {
	sess := m.(*Session)
	var confRead bgpOptions
	var bgpNeighborConfig string
	var err error
	// default -1
	confRead.localPreference = -1
	confRead.metricOut = -1
	confRead.preference = -1

	if instance == defaultWord {
		bgpNeighborConfig, err = sess.command("show configuration"+
			" protocols bgp group "+group+
			" neighbor "+ip+" | display set relative", jnprSess)
		if err != nil {
			return confRead, err
		}
	} else {
		bgpNeighborConfig, err = sess.command("show configuration"+
			" routing-instances "+instance+
			" protocols bgp group "+group+
			" neighbor "+ip+" | display set relative", jnprSess)
		if err != nil {
			return confRead, err
		}
	}
	if bgpNeighborConfig != emptyWord {
		confRead.ip = ip
		confRead.routingInstance = instance
		confRead.name = group
		for _, item := range strings.Split(bgpNeighborConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "family inet "):
				confRead.familyInet, err = readBgpOptsFamily(itemTrim, inetWord, confRead.familyInet)
				if err != nil {
					return confRead, err
				}

			case strings.HasPrefix(itemTrim, "family inet6 "):
				confRead.familyInet6, err = readBgpOptsFamily(itemTrim, inet6Word, confRead.familyInet6)
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "bfd-liveness-detection "):
				confRead.bfdLivenessDetection, err = readBgpOptsBfd(itemTrim, confRead.bfdLivenessDetection)
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "graceful-restart "):
				confRead.gracefulRestart, err = readBgpOptsGracefulRestart(itemTrim, confRead.gracefulRestart)
				if err != nil {
					return confRead, err
				}
			default:
				confRead, err = readBgpOptsSimple(itemTrim, confRead)
				if err != nil {
					return confRead, err
				}
			}
		}
	} else {
		confRead.ip = ""
		return confRead, nil
	}
	return confRead, nil
}
func delBgpNeighbor(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	if d.Get("routing_instance").(string) == defaultWord {
		configSet = append(configSet, "delete protocols bgp"+
			" group "+d.Get("group").(string)+
			" neighbor "+d.Get("ip").(string)+"\n")
	} else {
		configSet = append(configSet, "delete"+
			" routing-instances "+d.Get("routing_instance").(string)+
			" protocols bgp group "+d.Get("group").(string)+
			" neighbor "+d.Get("ip").(string)+"\n")
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}
	return nil
}

func fillBgpNeighborData(d *schema.ResourceData, bgpNeighborOptions bgpOptions) {
	tfErr := d.Set("ip", bgpNeighborOptions.ip)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", bgpNeighborOptions.routingInstance)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("group", bgpNeighborOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("accept_remote_nexthop", bgpNeighborOptions.acceptRemoteNexthop)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("advertise_external", bgpNeighborOptions.advertiseExternal)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("advertise_external_conditional", bgpNeighborOptions.advertiseExternalConditional)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("advertise_inactive", bgpNeighborOptions.advertiseInactive)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("advertise_peer_as", bgpNeighborOptions.advertisePeerAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("as_override", bgpNeighborOptions.asOverride)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("damping", bgpNeighborOptions.damping)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as_private", bgpNeighborOptions.localAsPrivate)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as_alias", bgpNeighborOptions.localAsAlias)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as_no_prepend_global_as", bgpNeighborOptions.localAsNoPrependGlobalAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("log_updown", bgpNeighborOptions.logUpdown)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_igp", bgpNeighborOptions.metricOutIgp)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_igp_delay_med_update", bgpNeighborOptions.metricOutIgpDelayMedUpdate)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_minimum_igp", bgpNeighborOptions.metricOutMinimumIgp)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("mtu_discovery", bgpNeighborOptions.mtuDiscovery)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("multihop", bgpNeighborOptions.multihop)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("multipath", bgpNeighborOptions.multipath)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("no_advertise_peer_as", bgpNeighborOptions.noAdvertisePeerAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("remove_private", bgpNeighborOptions.removePrivate)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("passive", bgpNeighborOptions.passive)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("hold_time", bgpNeighborOptions.holdTime)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as_loops", bgpNeighborOptions.localAsLoops)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_preference", bgpNeighborOptions.localPreference)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out", bgpNeighborOptions.metricOut)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_igp_offset", bgpNeighborOptions.metricOutIgpOffset)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_minimum_igp_offset", bgpNeighborOptions.metricOutMinimumIgpOffset)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("out_delay", bgpNeighborOptions.outDelay)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("preference", bgpNeighborOptions.preference)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_algorithm", bgpNeighborOptions.authenticationAlgorithm)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_key", bgpNeighborOptions.authenticationKey)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_key_chain", bgpNeighborOptions.authenticationKeyChain)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_address", bgpNeighborOptions.localAddress)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as", bgpNeighborOptions.localAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_interface", bgpNeighborOptions.localInterface)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("peer_as", bgpNeighborOptions.peerAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("export", bgpNeighborOptions.exportPolicy)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("import", bgpNeighborOptions.importPolicy)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("bfd_liveness_detection", bgpNeighborOptions.bfdLivenessDetection)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("family_inet", bgpNeighborOptions.familyInet)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("family_inet6", bgpNeighborOptions.familyInet6)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("graceful_restart", bgpNeighborOptions.gracefulRestart)
	if tfErr != nil {
		panic(tfErr)
	}
}
