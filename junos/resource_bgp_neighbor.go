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
		CreateWithoutTimeout: resourceBgpNeighborCreate,
		ReadWithoutTimeout:   resourceBgpNeighborRead,
		UpdateWithoutTimeout: resourceBgpNeighborUpdate,
		DeleteWithoutTimeout: resourceBgpNeighborDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBgpNeighborImport,
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
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"group": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
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
				Sensitive:     true,
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
			"bgp_multipath": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"multipath"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_protection": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"multiple_as": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"cluster": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"damping": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"family_evpn": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nlri_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "signaling",
							ValidateFunc: validation.StringInSlice([]string{"signaling"}, false),
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
			"family_inet": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nlri_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"any", "flow", "labeled-unicast", "unicast", "multicast"}, false),
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
							ValidateFunc: validation.StringInSlice(
								[]string{"any", "flow", "labeled-unicast", "unicast", "multicast"}, false),
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
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"keep_all": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"keep_none"},
			},
			"keep_none": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"keep_all"},
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
				ConflictsWith: []string{
					"metric_out_igp",
					"metric_out_igp_offset",
					"metric_out_igp_delay_med_update",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset",
				},
			},
			"metric_out_igp": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"metric_out",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset",
				},
			},
			"metric_out_igp_delay_med_update": {
				Type:     schema.TypeBool,
				Optional: true,
				ConflictsWith: []string{
					"metric_out",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset",
				},
			},
			"metric_out_igp_offset": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(-2147483648, 2147483647),
				ConflictsWith: []string{
					"metric_out",
					"metric_out_minimum_igp",
					"metric_out_minimum_igp_offset",
				},
			},
			"metric_out_minimum_igp": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"metric_out",
					"metric_out_igp",
					"metric_out_igp_offset",
					"metric_out_igp_delay_med_update",
				},
			},
			"metric_out_minimum_igp_offset": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(-2147483648, 2147483647),
				ConflictsWith: []string{
					"metric_out",
					"metric_out_igp",
					"metric_out_igp_offset",
					"metric_out_igp_delay_med_update",
				},
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
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"bgp_multipath"},
				Deprecated:    "use bgp_multipath instead",
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setBgpNeighbor(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("ip").(string) +
			idSeparator + d.Get("routing_instance").(string) +
			idSeparator + d.Get("group").(string))

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	bgpGroupExists, err := checkBgpGroupExists(d.Get("group").(string), d.Get("routing_instance").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !bgpGroupExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("bgp group %v doesn't exist", d.Get("group").(string)))...)
	}
	bgpNeighborxists, err := checkBgpNeighborExists(
		d.Get("ip").(string),
		d.Get("routing_instance").(string),
		d.Get("group").(string),
		clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if bgpNeighborxists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("bgp neighbor %v already exists in group %v (routing-instance %v)",
			d.Get("ip").(string), d.Get("group").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setBgpNeighbor(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_bgp_neighbor", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	bgpNeighborxists, err = checkBgpNeighborExists(
		d.Get("ip").(string),
		d.Get("routing_instance").(string),
		d.Get("group").(string),
		clt, junSess)
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

	return append(diagWarns, resourceBgpNeighborReadWJunSess(d, clt, junSess)...)
}

func resourceBgpNeighborRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceBgpNeighborReadWJunSess(d, clt, junSess)
}

func resourceBgpNeighborReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession) diag.Diagnostics {
	mutex.Lock()
	bgpNeighborOptions, err := readBgpNeighbor(
		d.Get("ip").(string),
		d.Get("routing_instance").(string),
		d.Get("group").(string),
		clt, junSess)
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
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delBgpOpts(d, "neighbor", clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setBgpNeighbor(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delBgpOpts(d, "neighbor", clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setBgpNeighbor(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_bgp_neighbor", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceBgpNeighborReadWJunSess(d, clt, junSess)...)
}

func resourceBgpNeighborDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delBgpNeighbor(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delBgpNeighbor(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_bgp_neighbor", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceBgpNeighborImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	bgpNeighborxists, err := checkBgpNeighborExists(idSplit[0], idSplit[1], idSplit[2], clt, junSess)
	if err != nil {
		return nil, err
	}
	if !bgpNeighborxists {
		return nil, fmt.Errorf("don't find bgp neighbor with id '%v' "+
			"(id must be <ip>"+idSeparator+"<routing_instance>"+idSeparator+"<group>)", d.Id())
	}
	bgpNeighborOptions, err := readBgpNeighbor(idSplit[0], idSplit[1], idSplit[2], clt, junSess)
	if err != nil {
		return nil, err
	}
	fillBgpNeighborData(d, bgpNeighborOptions)
	result[0] = d

	return result, nil
}

func checkBgpNeighborExists(ip, instance, group string, clt *Client, junSess *junosSession) (bool, error) {
	var showConfig string
	var err error
	if instance == defaultW {
		showConfig, err = clt.command(cmdShowConfig+
			"protocols bgp group "+group+" neighbor "+ip+pipeDisplaySet, junSess)
		if err != nil {
			return false, err
		}
	} else {
		showConfig, err = clt.command(cmdShowConfig+routingInstancesWS+instance+" "+
			"protocols bgp group "+group+" neighbor "+ip+pipeDisplaySet, junSess)
		if err != nil {
			return false, err
		}
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setBgpNeighbor(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	setPrefix += "protocols bgp group " + d.Get("group").(string) + " neighbor " + d.Get("ip").(string) + " "

	if err := setBgpOptsSimple(setPrefix, d, clt, junSess); err != nil {
		return err
	}
	if err := setBgpOptsBfd(setPrefix, d.Get("bfd_liveness_detection").([]interface{}), clt, junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, "evpn", d.Get("family_evpn").([]interface{}), clt, junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, inetW, d.Get("family_inet").([]interface{}), clt, junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, inet6W, d.Get("family_inet6").([]interface{}), clt, junSess); err != nil {
		return err
	}

	return setBgpOptsGrafefulRestart(setPrefix, d.Get("graceful_restart").([]interface{}), clt, junSess)
}

func readBgpNeighbor(ip, instance, group string, clt *Client, junSess *junosSession) (bgpOptions, error) {
	var confRead bgpOptions
	var showConfig string
	var err error
	// default -1
	confRead.localPreference = -1
	confRead.metricOut = -1
	confRead.preference = -1

	if instance == defaultW {
		showConfig, err = clt.command(cmdShowConfig+
			"protocols bgp group "+group+" neighbor "+ip+pipeDisplaySetRelative, junSess)
		if err != nil {
			return confRead, err
		}
	} else {
		showConfig, err = clt.command(cmdShowConfig+routingInstancesWS+instance+" "+
			"protocols bgp group "+group+" neighbor "+ip+pipeDisplaySetRelative, junSess)
		if err != nil {
			return confRead, err
		}
	}
	if showConfig != emptyW {
		confRead.ip = ip
		confRead.routingInstance = instance
		confRead.name = group
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "family evpn "):
				confRead.familyEvpn, err = readBgpOptsFamily(itemTrim, "evpn", confRead.familyEvpn)
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "family inet "):
				confRead.familyInet, err = readBgpOptsFamily(itemTrim, inetW, confRead.familyInet)
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "family inet6 "):
				confRead.familyInet6, err = readBgpOptsFamily(itemTrim, inet6W, confRead.familyInet6)
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "bfd-liveness-detection "):
				if len(confRead.bfdLivenessDetection) == 0 {
					confRead.bfdLivenessDetection = append(confRead.bfdLivenessDetection,
						map[string]interface{}{
							"authentication_algorithm":           "",
							"authentication_key_chain":           "",
							"authentication_loose_check":         false,
							"detection_time_threshold":           0,
							"holddown_interval":                  0,
							"minimum_interval":                   0,
							"minimum_receive_interval":           0,
							"multiplier":                         0,
							"session_mode":                       "",
							"transmit_interval_minimum_interval": 0,
							"transmit_interval_threshold":        0,
							"version":                            "",
						})
				}
				if err := readBgpOptsBfd(itemTrim, confRead.bfdLivenessDetection[0]); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "graceful-restart "):
				if len(confRead.gracefulRestart) == 0 {
					confRead.gracefulRestart = append(confRead.gracefulRestart, map[string]interface{}{
						"disable":          false,
						"restart_time":     0,
						"stale_route_time": 0,
					})
				}
				if err := readBgpOptsGracefulRestart(itemTrim, confRead.gracefulRestart[0]); err != nil {
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

func delBgpNeighbor(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0, 1)
	if d.Get("routing_instance").(string) == defaultW {
		configSet = append(configSet, "delete protocols bgp"+
			" group "+d.Get("group").(string)+
			" neighbor "+d.Get("ip").(string))
	} else {
		configSet = append(configSet, delRoutingInstances+d.Get("routing_instance").(string)+
			" protocols bgp group "+d.Get("group").(string)+
			" neighbor "+d.Get("ip").(string))
	}

	return clt.configSet(configSet, junSess)
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
	if tfErr := d.Set("cluster", bgpNeighborOptions.cluster); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("damping", bgpNeighborOptions.damping); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("export", bgpNeighborOptions.exportPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_evpn", bgpNeighborOptions.familyEvpn); tfErr != nil {
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
	if tfErr := d.Set("keep_all", bgpNeighborOptions.keepAll); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("keep_none", bgpNeighborOptions.keepNone); tfErr != nil {
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
	if _, ok := d.GetOk("multipath"); ok {
		if tfErr := d.Set("multipath", bgpNeighborOptions.multipath); tfErr != nil {
			panic(tfErr)
		}
	} else {
		if tfErr := d.Set("bgp_multipath", bgpNeighborOptions.bgpMultipath); tfErr != nil {
			panic(tfErr)
		}
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
