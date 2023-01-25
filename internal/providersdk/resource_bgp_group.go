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

func resourceBgpGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceBgpGroupCreate,
		ReadWithoutTimeout:   resourceBgpGroupRead,
		UpdateWithoutTimeout: resourceBgpGroupUpdate,
		DeleteWithoutTimeout: resourceBgpGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBgpGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          junos.DefaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "external",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"internal", "external"}, false),
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
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"automatic", "multihop", "single-hop"}, false),
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
							ConflictsWith: []string{
								"graceful_restart.0.restart_time",
								"graceful_restart.0.stale_route_time",
							},
						},
						"restart_time": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 600),
							ConflictsWith: []string{"graceful_restart.0.disable"},
						},
						"stale_route_time": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 600),
							ConflictsWith: []string{"graceful_restart.0.disable"},
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

func resourceBgpGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		if err := setBgpGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string) + junos.IDSeparator + d.Get("routing_instance").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if d.Get("routing_instance").(string) != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	bgpGroupxists, err := checkBgpGroupExists(d.Get("name").(string), d.Get("routing_instance").(string), clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if bgpGroupxists {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("bgp group %v already exists in routing-instance %v",
			d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setBgpGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("create resource junos_bgp_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	bgpGroupxists, err = checkBgpGroupExists(d.Get("name").(string), d.Get("routing_instance").(string), clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if bgpGroupxists {
		d.SetId(d.Get("name").(string) + junos.IDSeparator + d.Get("routing_instance").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("bgp group %v not exists in routing-instance %v after commit "+
			"=> check your config", d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceBgpGroupReadWJunSess(d, clt, junSess)...)
}

func resourceBgpGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)

	return resourceBgpGroupReadWJunSess(d, clt, junSess)
}

func resourceBgpGroupReadWJunSess(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) diag.Diagnostics {
	mutex.Lock()
	bgpGroupOptions, err := readBgpGroup(d.Get("name").(string), d.Get("routing_instance").(string), clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if bgpGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillBgpGroupData(d, bgpGroupOptions)
	}

	return nil
}

func resourceBgpGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		if err := delBgpOpts(d, "group", clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setBgpGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delBgpOpts(d, "group", clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setBgpGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("update resource junos_bgp_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceBgpGroupReadWJunSess(d, clt, junSess)...)
}

func resourceBgpGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		if err := delBgpGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.CloseSession(junSess)
	if err := clt.ConfigLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delBgpGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.CommitConf("delete resource junos_bgp_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.ConfigClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceBgpGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.CloseSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), junos.IDSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	bgpGroupxists, err := checkBgpGroupExists(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	if !bgpGroupxists {
		return nil, fmt.Errorf("don't find bgp group with id '%v' "+
			"(id must be <name>"+junos.IDSeparator+"<routing_instance>)", d.Id())
	}
	bgpGroupOptions, err := readBgpGroup(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	fillBgpGroupData(d, bgpGroupOptions)
	result[0] = d

	return result, nil
}

func checkBgpGroupExists(bgpGroup, instance string, clt *junos.Client, junSess *junos.Session) (_ bool, err error) {
	var showConfig string
	if instance == junos.DefaultW {
		showConfig, err = clt.Command(junos.CmdShowConfig+
			"protocols bgp group "+bgpGroup+junos.PipeDisplaySet, junSess)
		if err != nil {
			return false, err
		}
	} else {
		showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+instance+" "+
			"protocols bgp group "+bgpGroup+junos.PipeDisplaySet, junSess)
		if err != nil {
			return false, err
		}
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setBgpGroup(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	setPrefix := junos.SetLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	setPrefix += "protocols bgp group " + d.Get("name").(string) + " "

	if err := clt.ConfigSet([]string{setPrefix + "type " + d.Get("type").(string)}, junSess); err != nil {
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
	if err := setBgpOptsSimple(setPrefix, d, clt, junSess); err != nil {
		return err
	}
	if err := setBgpOptsBfd(setPrefix, d.Get("bfd_liveness_detection").([]interface{}), clt, junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, "evpn", d.Get("family_evpn").([]interface{}), clt, junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, junos.InetW, d.Get("family_inet").([]interface{}), clt, junSess); err != nil {
		return err
	}
	if err := setBgpOptsFamily(setPrefix, junos.Inet6W, d.Get("family_inet6").([]interface{}), clt, junSess); err != nil {
		return err
	}

	return setBgpOptsGrafefulRestart(setPrefix, d.Get("graceful_restart").([]interface{}), clt, junSess)
}

func readBgpGroup(bgpGroup, instance string, clt *junos.Client, junSess *junos.Session,
) (confRead bgpOptions, err error) {
	// default -1
	confRead.localPreference = -1
	confRead.metricOut = -1
	confRead.preference = -1
	var showConfig string
	if instance == junos.DefaultW {
		showConfig, err = clt.Command(junos.CmdShowConfig+
			"protocols bgp group "+bgpGroup+junos.PipeDisplaySetRelative, junSess)
		if err != nil {
			return confRead, err
		}
	} else {
		showConfig, err = clt.Command(junos.CmdShowConfig+junos.RoutingInstancesWS+instance+" "+
			"protocols bgp group "+bgpGroup+junos.PipeDisplaySetRelative, junSess)
		if err != nil {
			return confRead, err
		}
	}
	if showConfig != junos.EmptyW {
		confRead.name = bgpGroup
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

func delBgpGroup(d *schema.ResourceData, clt *junos.Client, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	if d.Get("routing_instance").(string) == junos.DefaultW {
		configSet = append(configSet, "delete protocols bgp group "+d.Get("name").(string))
	} else {
		configSet = append(configSet, junos.DelRoutingInstances+d.Get("routing_instance").(string)+
			" protocols bgp group "+d.Get("name").(string))
	}

	return clt.ConfigSet(configSet, junSess)
}

func fillBgpGroupData(d *schema.ResourceData, bgpGroupOptions bgpOptions) {
	if tfErr := d.Set("name", bgpGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", bgpGroupOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("accept_remote_nexthop", bgpGroupOptions.acceptRemoteNexthop); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advertise_external", bgpGroupOptions.advertiseExternal); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advertise_external_conditional", bgpGroupOptions.advertiseExternalConditional); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advertise_inactive", bgpGroupOptions.advertiseInactive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("advertise_peer_as", bgpGroupOptions.advertisePeerAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("as_override", bgpGroupOptions.asOverride); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_algorithm", bgpGroupOptions.authenticationAlgorithm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_key", bgpGroupOptions.authenticationKey); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_key_chain", bgpGroupOptions.authenticationKeyChain); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bfd_liveness_detection", bgpGroupOptions.bfdLivenessDetection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("cluster", bgpGroupOptions.cluster); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("damping", bgpGroupOptions.damping); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("export", bgpGroupOptions.exportPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_evpn", bgpGroupOptions.familyEvpn); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet", bgpGroupOptions.familyInet); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("family_inet6", bgpGroupOptions.familyInet6); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("graceful_restart", bgpGroupOptions.gracefulRestart); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("hold_time", bgpGroupOptions.holdTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("import", bgpGroupOptions.importPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("keep_all", bgpGroupOptions.keepAll); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("keep_none", bgpGroupOptions.keepNone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_address", bgpGroupOptions.localAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as", bgpGroupOptions.localAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as_alias", bgpGroupOptions.localAsAlias); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as_loops", bgpGroupOptions.localAsLoops); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as_no_prepend_global_as", bgpGroupOptions.localAsNoPrependGlobalAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_as_private", bgpGroupOptions.localAsPrivate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_interface", bgpGroupOptions.localInterface); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("local_preference", bgpGroupOptions.localPreference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("log_updown", bgpGroupOptions.logUpdown); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out", bgpGroupOptions.metricOut); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_igp", bgpGroupOptions.metricOutIgp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_igp_delay_med_update", bgpGroupOptions.metricOutIgpDelayMedUpdate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_igp_offset", bgpGroupOptions.metricOutIgpOffset); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_minimum_igp", bgpGroupOptions.metricOutMinimumIgp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("metric_out_minimum_igp_offset", bgpGroupOptions.metricOutMinimumIgpOffset); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("mtu_discovery", bgpGroupOptions.mtuDiscovery); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("multihop", bgpGroupOptions.multihop); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("bgp_multipath", bgpGroupOptions.bgpMultipath); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_advertise_peer_as", bgpGroupOptions.noAdvertisePeerAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("out_delay", bgpGroupOptions.outDelay); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("passive", bgpGroupOptions.passive); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("peer_as", bgpGroupOptions.peerAs); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("preference", bgpGroupOptions.preference); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("remove_private", bgpGroupOptions.removePrivate); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("type", bgpGroupOptions.bgpType); tfErr != nil {
		panic(tfErr)
	}
}
