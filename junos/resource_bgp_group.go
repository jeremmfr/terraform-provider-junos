package junos

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBgpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBgpGroupCreate,
		Read:   resourceBgpGroupRead,
		Update: resourceBgpGroupUpdate,
		Delete: resourceBgpGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBgpGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNameObjectJunos(),
			},
			"routing_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      defaultWord,
				ValidateFunc: validateNameObjectJunos(),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "external",
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !stringInSlice(value, []string{"internal", "external"}) {
						errors = append(errors, fmt.Errorf(
							"%q for %q is not 'internal' or 'external'", value, k))
					}

					return
				},
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
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"no_advertise_peer_as"},
			},
			"no_advertise_peer_as": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"advertise_peer_as"},
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
							ConflictsWith: []string{"graceful_restart.0.restart_time",
								"graceful_restart.0.stale_route_time"},
						},
						"restart_time": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validateIntRange(1, 600),
							ConflictsWith: []string{"graceful_restart.0.disable"},
						},
						"stale_route_time": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validateIntRange(1, 600),
							ConflictsWith: []string{"graceful_restart.0.disable"},
						},
					},
				},
			},
		},
	}
}

func resourceBgpGroupCreate(d *schema.ResourceData, m interface{}) error {
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
	bgpGroupxists, err := checkBgpGroupExists(d.Get("name").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if bgpGroupxists {
		sess.configClear(jnprSess)

		return fmt.Errorf("bgp group %v already exists in routing-instance %v",
			d.Get("name").(string), d.Get("routing_instance").(string))
	}
	err = setBgpGroup(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_bgp_group", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	bgpGroupxists, err = checkBgpGroupExists(d.Get("name").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		return err
	}
	if bgpGroupxists {
		d.SetId(d.Get("name").(string) + idSeparator + d.Get("routing_instance").(string))
	} else {
		return fmt.Errorf("bgp group %v not exists in routing-instance %v after commit "+
			"=> check your config", d.Get("name").(string), d.Get("routing_instance").(string))
	}

	return resourceBgpGroupRead(d, m)
}
func resourceBgpGroupRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	bgpGroupOptions, err := readBgpGroup(d.Get("name").(string), d.Get("routing_instance").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if bgpGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillBgpGroupData(d, bgpGroupOptions)
	}

	return nil
}
func resourceBgpGroupUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delBgpOpts(d, "group", m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = setBgpGroup(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_bgp_group", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceBgpGroupRead(d, m)
}
func resourceBgpGroupDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delBgpGroup(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_bgp_group", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
}
func resourceBgpGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	bgpGroupxists, err := checkBgpGroupExists(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !bgpGroupxists {
		return nil, fmt.Errorf("don't find bgp group with id '%v' "+
			"(id must be <name>"+idSeparator+"<routing_instance>)", d.Id())
	}
	bgpGroupOptions, err := readBgpGroup(idSplit[0], idSplit[1], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillBgpGroupData(d, bgpGroupOptions)
	result[0] = d

	return result, nil
}

func checkBgpGroupExists(bgpGroup, instance string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	var bgpGroupConfig string
	var err error
	if instance == defaultWord {
		bgpGroupConfig, err = sess.command("show configuration protocols bgp group "+
			bgpGroup+" | display set", jnprSess)
		if err != nil {
			return false, err
		}
	} else {
		bgpGroupConfig, err = sess.command("show configuration routing-instances "+
			instance+" protocols bgp group "+bgpGroup+" | display set", jnprSess)
		if err != nil {
			return false, err
		}
	}
	if bgpGroupConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setBgpGroup(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	setPrefix := setLineStart
	if d.Get("routing_instance").(string) == defaultWord {
		setPrefix += "protocols bgp group " + d.Get("name").(string) + " "
	} else {
		setPrefix += "routing-instances " + d.Get("routing_instance").(string) +
			" protocols bgp group " + d.Get("name").(string) + " "
	}
	sess := m.(*Session)
	err := sess.configSet([]string{setPrefix + "type " + d.Get("type").(string) + "\n"}, jnprSess)
	if err != nil {
		return err
	}
	if d.Get("type").(string) == "external" {
		if d.Get("advertise_external").(bool) {
			return fmt.Errorf("conflict between type=external and advertise_external")
		}
		if d.Get("accept_remote_nexthop").(bool) && d.Get("multihop").(bool) {
			return fmt.Errorf("conflict between type=external and accept_remote_nexthop + multihop")
		}
	}
	err = setBgpOptsSimple(setPrefix, d, m, jnprSess)
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
func readBgpGroup(bgpGroup, instance string, m interface{}, jnprSess *NetconfObject) (bgpOptions, error) {
	sess := m.(*Session)
	var confRead bgpOptions
	var bgpGroupConfig string
	var err error
	// default -1
	confRead.localPreference = -1
	confRead.metricOut = -1
	confRead.preference = -1

	if instance == defaultWord {
		bgpGroupConfig, err = sess.command("show configuration protocols bgp group "+
			bgpGroup+" | display set relative", jnprSess)
		if err != nil {
			return confRead, err
		}
	} else {
		bgpGroupConfig, err = sess.command("show configuration routing-instances "+
			instance+" protocols bgp group "+bgpGroup+" | display set relative", jnprSess)
		if err != nil {
			return confRead, err
		}
	}
	if bgpGroupConfig != emptyWord {
		confRead.name = bgpGroup
		confRead.routingInstance = instance
		for _, item := range strings.Split(bgpGroupConfig, "\n") {
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
		confRead.name = ""

		return confRead, nil
	}

	return confRead, nil
}
func delBgpGroup(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	if d.Get("routing_instance").(string) == defaultWord {
		configSet = append(configSet, "delete protocols bgp group "+d.Get("name").(string)+"\n")
	} else {
		configSet = append(configSet, "delete routing-instances "+d.Get("routing_instance").(string)+
			" protocols bgp group "+d.Get("name").(string)+"\n")
	}
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}

func fillBgpGroupData(d *schema.ResourceData, bgpGroupOptions bgpOptions) {
	tfErr := d.Set("name", bgpGroupOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("routing_instance", bgpGroupOptions.routingInstance)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("accept_remote_nexthop", bgpGroupOptions.acceptRemoteNexthop)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("advertise_external", bgpGroupOptions.advertiseExternal)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("advertise_external_conditional", bgpGroupOptions.advertiseExternalConditional)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("advertise_inactive", bgpGroupOptions.advertiseInactive)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("advertise_peer_as", bgpGroupOptions.advertisePeerAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("as_override", bgpGroupOptions.asOverride)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("damping", bgpGroupOptions.damping)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as_private", bgpGroupOptions.localAsPrivate)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as_alias", bgpGroupOptions.localAsAlias)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as_no_prepend_global_as", bgpGroupOptions.localAsNoPrependGlobalAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("log_updown", bgpGroupOptions.logUpdown)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_igp", bgpGroupOptions.metricOutIgp)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_igp_delay_med_update", bgpGroupOptions.metricOutIgpDelayMedUpdate)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_minimum_igp", bgpGroupOptions.metricOutMinimumIgp)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("mtu_discovery", bgpGroupOptions.mtuDiscovery)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("multihop", bgpGroupOptions.multihop)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("multipath", bgpGroupOptions.multipath)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("no_advertise_peer_as", bgpGroupOptions.noAdvertisePeerAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("remove_private", bgpGroupOptions.removePrivate)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("passive", bgpGroupOptions.passive)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("hold_time", bgpGroupOptions.holdTime)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as_loops", bgpGroupOptions.localAsLoops)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_preference", bgpGroupOptions.localPreference)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out", bgpGroupOptions.metricOut)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_igp_offset", bgpGroupOptions.metricOutIgpOffset)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("metric_out_minimum_igp_offset", bgpGroupOptions.metricOutMinimumIgpOffset)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("out_delay", bgpGroupOptions.outDelay)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("preference", bgpGroupOptions.preference)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_algorithm", bgpGroupOptions.authenticationAlgorithm)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_key", bgpGroupOptions.authenticationKey)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("authentication_key_chain", bgpGroupOptions.authenticationKeyChain)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("type", bgpGroupOptions.bgpType)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_address", bgpGroupOptions.localAddress)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_as", bgpGroupOptions.localAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("local_interface", bgpGroupOptions.localInterface)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("peer_as", bgpGroupOptions.peerAs)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("export", bgpGroupOptions.exportPolicy)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("import", bgpGroupOptions.importPolicy)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("bfd_liveness_detection", bgpGroupOptions.bfdLivenessDetection)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("family_inet", bgpGroupOptions.familyInet)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("family_inet6", bgpGroupOptions.familyInet6)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("graceful_restart", bgpGroupOptions.gracefulRestart)
	if tfErr != nil {
		panic(tfErr)
	}
}
