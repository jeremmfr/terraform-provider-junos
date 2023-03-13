package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
				Default:          junos.DefaultW,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setBgpNeighbor(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("ip").(string) +
			junos.IDSeparator + d.Get("routing_instance").(string) +
			junos.IDSeparator + d.Get("group").(string))

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
	bgpGroupExists, err := checkBgpGroupExists(d.Get("group").(string), d.Get("routing_instance").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !bgpGroupExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("bgp group %v doesn't exist", d.Get("group").(string)))...)
	}
	bgpNeighborxists, err := checkBgpNeighborExists(
		d.Get("ip").(string),
		d.Get("routing_instance").(string),
		d.Get("group").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if bgpNeighborxists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("bgp neighbor %v already exists in group %v (routing-instance %v)",
			d.Get("ip").(string), d.Get("group").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setBgpNeighbor(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_bgp_neighbor")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	bgpNeighborxists, err = checkBgpNeighborExists(
		d.Get("ip").(string),
		d.Get("routing_instance").(string),
		d.Get("group").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if bgpNeighborxists {
		d.SetId(d.Get("ip").(string) +
			junos.IDSeparator + d.Get("routing_instance").(string) +
			junos.IDSeparator + d.Get("group").(string))
	} else {
		return append(diagWarns,
			diag.FromErr(fmt.Errorf("bgp neighbor %v not exists in group %v (routing-instance %v) after commit "+
				"=> check your config", d.Get("ip").(string), d.Get("group").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceBgpNeighborReadWJunSess(d, junSess)...)
}

func resourceBgpNeighborRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceBgpNeighborReadWJunSess(d, junSess)
}

func resourceBgpNeighborReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	bgpNeighborOptions, err := readBgpNeighbor(
		d.Get("ip").(string),
		d.Get("routing_instance").(string),
		d.Get("group").(string),
		junSess,
	)
	junos.MutexUnlock()
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
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delBgpOpts(d, "neighbor", junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setBgpNeighbor(d, junSess); err != nil {
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
	if err := delBgpOpts(d, "neighbor", junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setBgpNeighbor(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_bgp_neighbor")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceBgpNeighborReadWJunSess(d, junSess)...)
}

func resourceBgpNeighborDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delBgpNeighbor(d, junSess); err != nil {
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
	if err := delBgpNeighbor(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_bgp_neighbor")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceBgpNeighborImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), junos.IDSeparator)
	if len(idSplit) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	bgpNeighborxists, err := checkBgpNeighborExists(idSplit[0], idSplit[1], idSplit[2], junSess)
	if err != nil {
		return nil, err
	}
	if !bgpNeighborxists {
		return nil, fmt.Errorf("don't find bgp neighbor with id '%v' "+
			"(id must be <ip>"+junos.IDSeparator+"<routing_instance>"+junos.IDSeparator+"<group>)", d.Id())
	}
	bgpNeighborOptions, err := readBgpNeighbor(idSplit[0], idSplit[1], idSplit[2], junSess)
	if err != nil {
		return nil, err
	}
	fillBgpNeighborData(d, bgpNeighborOptions)
	result[0] = d

	return result, nil
}

func checkBgpNeighborExists(ip, instance, group string, junSess *junos.Session) (_ bool, err error) {
	var showConfig string
	if instance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols bgp group " + group + " neighbor " + ip + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
			"protocols bgp group " + group + " neighbor " + ip + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setBgpNeighbor(d *schema.ResourceData, junSess *junos.Session) error {
	setPrefix := junos.SetLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	setPrefix += "protocols bgp group " + d.Get("group").(string) + " neighbor " + d.Get("ip").(string) + " "

	if err := setBgpOptsSimple(setPrefix, d, junSess); err != nil {
		return err
	}
	if err := setBgpOptsBfd(setPrefix, d.Get("bfd_liveness_detection").([]interface{}), junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, "evpn", d.Get("family_evpn").([]interface{}), junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, junos.InetW, d.Get("family_inet").([]interface{}), junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, junos.Inet6W, d.Get("family_inet6").([]interface{}), junSess); err != nil {
		return err
	}

	return setBgpOptsGrafefulRestart(setPrefix, d.Get("graceful_restart").([]interface{}), junSess)
}

func readBgpNeighbor(ip, instance, group string, junSess *junos.Session,
) (confRead bgpOptions, err error) {
	// default -1
	confRead.localPreference = -1
	confRead.metricOut = -1
	confRead.preference = -1
	var showConfig string
	if instance == junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols bgp group " + group + " neighbor " + ip + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + instance + " " +
			"protocols bgp group " + group + " neighbor " + ip + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	}
	if showConfig != junos.EmptyW {
		confRead.ip = ip
		confRead.routingInstance = instance
		confRead.name = group
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "family evpn "):
				confRead.familyEvpn, err = readBgpOptsFamily(itemTrim, confRead.familyEvpn)
				if err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet "):
				confRead.familyInet, err = readBgpOptsFamily(itemTrim, confRead.familyInet)
				if err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet6 "):
				confRead.familyInet6, err = readBgpOptsFamily(itemTrim, confRead.familyInet6)
				if err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
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
			case balt.CutPrefixInString(&itemTrim, "graceful-restart "):
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
				err = confRead.readBgpOptsSimple(itemTrim)
				if err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}

func delBgpNeighbor(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	if d.Get("routing_instance").(string) == junos.DefaultW {
		configSet = append(configSet, "delete protocols bgp"+
			" group "+d.Get("group").(string)+
			" neighbor "+d.Get("ip").(string))
	} else {
		configSet = append(configSet, junos.DelRoutingInstances+d.Get("routing_instance").(string)+
			" protocols bgp group "+d.Get("group").(string)+
			" neighbor "+d.Get("ip").(string))
	}

	return junSess.ConfigSet(configSet)
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
	if tfErr := d.Set("bgp_multipath", bgpNeighborOptions.bgpMultipath); tfErr != nil {
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
