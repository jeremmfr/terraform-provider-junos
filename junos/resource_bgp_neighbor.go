package junos

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBgpNeighborCreate,
		ReadContext:   resourceBgpNeighborRead,
		UpdateContext: resourceBgpNeighborUpdate,
		DeleteContext: resourceBgpNeighborDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBgpNeighborImport,
		},
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultWord,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
			},
			"group": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64),
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
			"bfd_liveness_detection": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_algorithm": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authentication_key_chain": {
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
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"holddown_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255000),
						},
						"minimum_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255000),
						},
						"minimum_receive_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255000),
						},
						"multiplier": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255),
						},
						"session_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"automatic", "multihop", "single-hop"}, false),
						},
						"transmit_interval_minimum_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255000),
						},
						"transmit_interval_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"damping": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"family_inet": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nlri_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"any", "flow", "labeled-unicast", "unicast", "multicast"}, false),
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
										ValidateFunc: validation.IntBetween(1, 4294967295),
									},
									"teardown": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
									"teardown_idle_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 2400),
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
										ValidateFunc: validation.IntBetween(1, 4294967295),
									},
									"teardown": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
									"teardown_idle_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 2400),
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
							ValidateFunc: validation.StringInSlice([]string{
								"any", "flow", "labeled-unicast", "unicast", "multicast"}, false),
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
										ValidateFunc: validation.IntBetween(1, 4294967295),
									},
									"teardown": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
									"teardown_idle_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 2400),
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
										ValidateFunc: validation.IntBetween(1, 4294967295),
									},
									"teardown": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
									"teardown_idle_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 2400),
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
							ValidateFunc: validation.IntBetween(1, 600),
						},
						"stale_route_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 600),
						},
					},
				},
			},
			"hold_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(3, 65535),
			},
			"import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"local_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"local_as": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_as_alias": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"local_as_private", "local_as_no_prepend_global_as"},
			},
			"local_as_loops": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"local_as_no_prepend_global_as": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"local_as_private", "local_as_alias"},
			},
			"local_as_private": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"local_as_alias", "local_as_no_prepend_global_as"},
			},
			"local_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_preference": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4294967295),
				Default:      -1,
			},
			"log_updown": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"metric_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4294967295),
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
			"metric_out_igp_delay_med_update": {
				Type:     schema.TypeBool,
				Optional: true,
				ConflictsWith: []string{"metric_out",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset"},
			},
			"metric_out_igp_offset": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(-2147483648, 2147483647),
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
				ValidateFunc: validation.IntBetween(-2147483648, 2147483647),
				ConflictsWith: []string{"metric_out",
					"metric_out_igp",
					"metric_out_igp_offset",
					"metric_out_igp_delay_med_update"},
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
				ConflictsWith: []string{"multipath_options"},
			},
                        "multipath_options": {
                                Type:     schema.TypeList,
                                Optional: true,
				MaxItems: 1,
				ConflictsWith: []string{"multipath"},
                                Elem: &schema.Resource{
                                        Schema: map[string]*schema.Schema{
                                                "disable": {
                                                        Type:     schema.TypeBool,
                                                        Optional: true,
                                                },
                                                "allow_protection": {
                                                        Type:          schema.TypeBool,
                                                        Optional:      true,
                                                },
                                                "multiple_as": {
                                                        Type:          schema.TypeBool,
                                                        Optional:      true,
                                                },
                                        },
                                },

                        },

			"out_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"passive": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"peer_as": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"preference": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4294967295),
				Default:      -1,
			},
			"remove_private": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceBgpNeighborCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if d.Get("routing_instance").(string) != defaultWord {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), m, jnprSess)
		if err != nil {
			sess.configClear(jnprSess)

			return diag.FromErr(err)
		}
		if !instanceExists {
			sess.configClear(jnprSess)

			return diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))
		}
	}
	bgpGroupExists, err := checkBgpGroupExists(d.Get("group").(string), d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if !bgpGroupExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("bgp group %v doesn't exist", d.Get("group").(string)))
	}
	bgpNeighborxists, err := checkBgpNeighborExists(d.Get("ip").(string),
		d.Get("routing_instance").(string), d.Get("group").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if bgpNeighborxists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("bgp neighbor %v already exists in group %v (routing-instance %v)",
			d.Get("ip").(string), d.Get("group").(string), d.Get("routing_instance").(string)))
	}
	if err := setBgpNeighbor(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_bgp_neighbor", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	bgpNeighborxists, err = checkBgpNeighborExists(d.Get("ip").(string),
		d.Get("routing_instance").(string), d.Get("group").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if bgpNeighborxists {
		d.SetId(d.Get("ip").(string) +
			idSeparator + d.Get("routing_instance").(string) +
			idSeparator + d.Get("group").(string))
	} else {
		return append(diagWarns,
			diag.FromErr(fmt.Errorf("bgp neighbor %v not exists in group %v (routing-instance %v) after commit "+
				"=> check your config", d.Get("ip").(string), d.Get("group").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceBgpNeighborReadWJnprSess(d, m, jnprSess)...)
}
func resourceBgpNeighborRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceBgpNeighborReadWJnprSess(d, m, jnprSess)
}
func resourceBgpNeighborReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	bgpNeighborOptions, err := readBgpNeighbor(d.Get("ip").(string),
		d.Get("routing_instance").(string), d.Get("group").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if bgpNeighborOptions.ip == "" {
		d.SetId("")
	} else {
		fillBgpNeighborData(d, bgpNeighborOptions)
	}

	return nil
}
func resourceBgpNeighborUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delBgpOpts(d, "neighbor", m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setBgpNeighbor(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_bgp_neighbor", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceBgpNeighborReadWJnprSess(d, m, jnprSess)...)
}
func resourceBgpNeighborDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delBgpNeighbor(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_bgp_neighbor", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
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
			"(id must be <ip>"+idSeparator+"<routing_instance>"+idSeparator+"<group>)", d.Id())
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
	if err := setBgpOptsSimple(setPrefix, d, m, jnprSess); err != nil {
		return err
	}
	if err := setBgpOptsBfd(setPrefix, d.Get("bfd_liveness_detection").([]interface{}), m, jnprSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, inetWord, d.Get("family_inet").([]interface{}), m, jnprSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, inet6Word, d.Get("family_inet6").([]interface{}), m, jnprSess); err != nil {
		return err
	}
	if err := setBgpOptsGrafefulRestart(setPrefix, d.Get("graceful_restart").([]interface{}), m, jnprSess); err != nil {
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
                        case strings.HasPrefix(itemTrim, "multipath "):
                                confRead.multipath_options, err = readBgpOptsMultipath(itemTrim, confRead.multipath_options)
                                if err != nil {
                                        return confRead, err
                                }

			default:
				err = readBgpOptsSimple(itemTrim, &confRead)
				if err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}
func delBgpNeighbor(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	if d.Get("routing_instance").(string) == defaultWord {
		configSet = append(configSet, "delete protocols bgp"+
			" group "+d.Get("group").(string)+
			" neighbor "+d.Get("ip").(string))
	} else {
		configSet = append(configSet, deleteWord+
			" routing-instances "+d.Get("routing_instance").(string)+
			" protocols bgp group "+d.Get("group").(string)+
			" neighbor "+d.Get("ip").(string))
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func fillBgpNeighborData(d *schema.ResourceData, bgpNeighborOptions bgpOptions) {
	if tfErr := d.Set("ip", bgpNeighborOptions.ip); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", bgpNeighborOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("group", bgpNeighborOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("accept_remote_nexthop", bgpNeighborOptions.acceptRemoteNexthop); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advertise_external", bgpNeighborOptions.advertiseExternal); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advertise_external_conditional", bgpNeighborOptions.advertiseExternalConditional); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advertise_inactive", bgpNeighborOptions.advertiseInactive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advertise_peer_as", bgpNeighborOptions.advertisePeerAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_override", bgpNeighborOptions.asOverride); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_algorithm", bgpNeighborOptions.authenticationAlgorithm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_key", bgpNeighborOptions.authenticationKey); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_key_chain", bgpNeighborOptions.authenticationKeyChain); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bfd_liveness_detection", bgpNeighborOptions.bfdLivenessDetection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("damping", bgpNeighborOptions.damping); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("export", bgpNeighborOptions.exportPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet", bgpNeighborOptions.familyInet); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet6", bgpNeighborOptions.familyInet6); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("graceful_restart", bgpNeighborOptions.gracefulRestart); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hold_time", bgpNeighborOptions.holdTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("import", bgpNeighborOptions.importPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_address", bgpNeighborOptions.localAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as", bgpNeighborOptions.localAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as_alias", bgpNeighborOptions.localAsAlias); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as_loops", bgpNeighborOptions.localAsLoops); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as_no_prepend_global_as", bgpNeighborOptions.localAsNoPrependGlobalAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as_private", bgpNeighborOptions.localAsPrivate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_interface", bgpNeighborOptions.localInterface); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_preference", bgpNeighborOptions.localPreference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("log_updown", bgpNeighborOptions.logUpdown); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out", bgpNeighborOptions.metricOut); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_igp", bgpNeighborOptions.metricOutIgp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_igp_delay_med_update", bgpNeighborOptions.metricOutIgpDelayMedUpdate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_igp_offset", bgpNeighborOptions.metricOutIgpOffset); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_minimum_igp", bgpNeighborOptions.metricOutMinimumIgp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_minimum_igp_offset", bgpNeighborOptions.metricOutMinimumIgpOffset); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("mtu_discovery", bgpNeighborOptions.mtuDiscovery); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("multihop", bgpNeighborOptions.multihop); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("multipath", bgpNeighborOptions.multipath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("multipath_options", bgpNeighborOptions.multipath_options); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_advertise_peer_as", bgpNeighborOptions.noAdvertisePeerAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("out_delay", bgpNeighborOptions.outDelay); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("passive", bgpNeighborOptions.passive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("peer_as", bgpNeighborOptions.peerAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preference", bgpNeighborOptions.preference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("remove_private", bgpNeighborOptions.removePrivate); tfErr != nil {
		panic(tfErr)
	}
}
