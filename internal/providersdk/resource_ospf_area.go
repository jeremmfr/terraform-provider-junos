package providersdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	jdecode "github.com/jeremmfr/junosdecode"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type ospfAreaOptions struct {
	noContextIdentifierAdvertisement bool
	areaID                           string
	routingInstance                  string
	realm                            string
	version                          string
	contextIdentifier                []string
	interAreaPrefixExport            []string
	interAreaPrefixImport            []string
	networkSummaryExport             []string
	networkSummaryImport             []string
	areaRange                        []map[string]interface{}
	interFace                        []map[string]interface{}
	nssa                             []map[string]interface{}
	stub                             []map[string]interface{}
	virtualLink                      []map[string]interface{}
}

func resourceOspfArea() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceOspfAreaCreate,
		ReadWithoutTimeout:   resourceOspfAreaRead,
		UpdateWithoutTimeout: resourceOspfAreaUpdate,
		DeleteWithoutTimeout: resourceOspfAreaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOspfAreaImport,
		},
		Schema: map[string]*schema.Schema{
			"area_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^(\d|\.)+$`), "should be usually in the IP format"),
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          junos.DefaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"realm": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4-unicast", "ipv4-multicast", "ipv6-multicast"}, false),
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "v2",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"v2", "v3"}, false),
			},
			"interface": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"authentication_simple_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"authentication_md5": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_id": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"key": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile(`^\d{4}\-\d\d?\-\d\d?\.\d{2}:\d{2}:\d{2}$`),
											"must be in the format 'YYYY-MM-DD.HH:MM:SS'"),
									},
								},
							},
						},
						"bandwidth_based_metrics": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bandwidth": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(\d)+(m|k|g)?$`),
											`must be a bandwidth ^(\d)+(m|k|g)?$`),
									},
									"metric": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
								},
							},
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
									"full_neighbors_only": {
										Type:     schema.TypeBool,
										Optional: true,
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
									"no_adaptation": {
										Type:     schema.TypeBool,
										Optional: true,
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
						"dead_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"demand_circuit": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"dynamic_neighbors": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"flood_reduction": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"hello_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255),
						},
						"interface_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipsec_sa": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv4_adjacency_segment_protected_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"dynamic", "index", "label"}, false),
						},
						"ipv4_adjacency_segment_protected_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv4_adjacency_segment_unprotected_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"dynamic", "index", "label"}, false),
						},
						"ipv4_adjacency_segment_unprotected_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"link_protection": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"metric": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"mtu": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(128, 65535),
						},
						"neighbor": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPAddress,
									},
									"eligible": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"no_advertise_adjacency_segment": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"no_eligible_backup": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"no_eligible_remote_backup": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"no_interface_state_traps": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"no_neighbor_down_notification": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"node_link_protection": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"passive": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"passive_traffic_engineering_remote_node_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"passive_traffic_engineering_remote_node_router_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"poll_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 255),
						},
						"retransmit_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"secondary": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"strict_bfd": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"te_metric": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"transit_delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
					},
				},
			},
			"area_range": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
						"exact": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"override_metric": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(1, 16777215),
						},
						"restrict": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"context_identifier": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`), "must be in the IP format"),
				},
				ConflictsWith: []string{"no_context_identifier_advertisement"},
			},
			"inter_area_prefix_export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"inter_area_prefix_import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"network_summary_export": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"network_summary_import": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"no_context_identifier_advertisement": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"context_identifier"},
			},
			"nssa": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"stub"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"area_range": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"range": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsCIDRNetwork(0, 128),
									},
									"exact": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"override_metric": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntBetween(1, 16777215),
									},
									"restrict": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"default_lsa": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_metric": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 16777215),
									},
									"metric_type": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 2),
									},
									"type_7": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"no_summaries": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"nssa.0.summaries"},
						},
						"summaries": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"nssa.0.no_summaries"},
						},
					},
				},
			},
			"stub": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"nssa"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_metric": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 16777215),
						},
						"no_summaries": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"stub.0.summaries"},
						},
						"summaries": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"stub.0.no_summaries"},
						},
					},
				},
			},
			"virtual_link": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"neighbor_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPv4Address,
						},
						"transit_area": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`), "must be in the IP format"),
						},
						"dead_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"demand_circuit": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"flood_reduction": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"hello_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(1, 255),
						},
						"ipsec_sa": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"mtu": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(128, 65535),
						},
						"retransmit_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"transit_delay": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
					},
				},
			},
		},
	}
}

func resourceOspfAreaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setOspfArea(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if realm := d.Get("realm").(string); realm != "" {
			d.SetId(d.Get("area_id").(string) + junos.IDSeparator + d.Get("version").(string) +
				junos.IDSeparator + realm + junos.IDSeparator + d.Get("routing_instance").(string))
		} else {
			d.SetId(d.Get("area_id").(string) + junos.IDSeparator + d.Get("version").(string) +
				junos.IDSeparator + d.Get("routing_instance").(string))
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
	ospfAreaExists, err := checkOspfAreaExists(
		d.Get("area_id").(string),
		d.Get("version").(string),
		d.Get("realm").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ospfAreaExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		if realm := d.Get("realm").(string); realm != "" {
			return append(diagWarns, diag.FromErr(fmt.Errorf("ospf %v realm %v area %v already exists in routing instance %v",
				d.Get("version").(string), realm, d.Get("area_id").(string), d.Get("routing_instance").(string)))...)
		}

		return append(diagWarns, diag.FromErr(fmt.Errorf("ospf %v area %v already exists in routing instance %v",
			d.Get("version").(string), d.Get("area_id").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setOspfArea(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_ospf_area")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	ospfAreaExists, err = checkOspfAreaExists(
		d.Get("area_id").(string),
		d.Get("version").(string),
		d.Get("realm").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ospfAreaExists {
		if realm := d.Get("realm").(string); realm != "" {
			d.SetId(d.Get("area_id").(string) + junos.IDSeparator + d.Get("version").(string) +
				junos.IDSeparator + realm + junos.IDSeparator + d.Get("routing_instance").(string))
		} else {
			d.SetId(d.Get("area_id").(string) + junos.IDSeparator + d.Get("version").(string) +
				junos.IDSeparator + d.Get("routing_instance").(string))
		}
	} else {
		if realm := d.Get("realm").(string); realm != "" {
			return append(diagWarns,
				diag.FromErr(fmt.Errorf("ospf %v realm %v area %v in routing instance %v "+
					"not exists after commit => check your config",
					d.Get("version").(string), realm, d.Get("area_id").(string), d.Get("routing_instance").(string)))...)
		}

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("ospf %v area %v in routing instance %v not exists after commit => check your config",
				d.Get("version").(string), d.Get("area_id").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceOspfAreaReadWJunSess(d, junSess)...)
}

func resourceOspfAreaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceOspfAreaReadWJunSess(d, junSess)
}

func resourceOspfAreaReadWJunSess(d *schema.ResourceData, junSess *junos.Session) diag.Diagnostics {
	mutex.Lock()
	ospfAreaOptions, err := readOspfArea(
		d.Get("area_id").(string),
		d.Get("version").(string),
		d.Get("realm").(string),
		d.Get("routing_instance").(string),
		junSess,
	)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if ospfAreaOptions.areaID == "" {
		d.SetId("")
	} else {
		fillOspfAreaData(d, ospfAreaOptions)
	}

	return nil
}

func resourceOspfAreaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delOspfArea(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setOspfArea(d, junSess); err != nil {
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
	if err := delOspfArea(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setOspfArea(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_ospf_area")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceOspfAreaReadWJunSess(d, junSess)...)
}

func resourceOspfAreaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delOspfArea(d, junSess); err != nil {
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
	if err := delOspfArea(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_ospf_area")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceOspfAreaImport(ctx context.Context, d *schema.ResourceData, m interface{},
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
	if len(idSplit) == 3 {
		ospfAreaExists, err := checkOspfAreaExists(idSplit[0], idSplit[1], "", idSplit[2], junSess)
		if err != nil {
			return nil, err
		}
		if !ospfAreaExists {
			return nil,
				fmt.Errorf(
					"don't find ospf area with id '%v' (id must be "+
						"<aread_id>"+junos.IDSeparator+"<version>"+junos.IDSeparator+"<routing_instance> or "+
						"<aread_id>"+junos.IDSeparator+"<version>"+junos.IDSeparator+"<realm>"+junos.IDSeparator+"<routing_instance>)",
					d.Id(),
				)
		}
		ospfAreaOptions, err := readOspfArea(idSplit[0], idSplit[1], "", idSplit[2], junSess)
		if err != nil {
			return nil, err
		}
		fillOspfAreaData(d, ospfAreaOptions)
		result[0] = d

		return result, nil
	}
	ospfAreaExists, err := checkOspfAreaExists(idSplit[0], idSplit[1], idSplit[2], idSplit[3], junSess)
	if err != nil {
		return nil, err
	}
	if !ospfAreaExists {
		return nil,
			fmt.Errorf(
				"don't find ospf area with id '%v' (id must be "+
					"<aread_id>"+junos.IDSeparator+"<version>"+junos.IDSeparator+"<routing_instance> or "+
					"<aread_id>"+junos.IDSeparator+"<version>"+junos.IDSeparator+"<realm>"+junos.IDSeparator+"<routing_instance>)",
				d.Id(),
			)
	}
	ospfAreaOptions, err := readOspfArea(idSplit[0], idSplit[1], idSplit[2], idSplit[3], junSess)
	if err != nil {
		return nil, err
	}
	fillOspfAreaData(d, ospfAreaOptions)
	if ospfAreaOptions.realm == "" {
		d.SetId(idSplit[0] + junos.IDSeparator + idSplit[1] + junos.IDSeparator + idSplit[3])
	}
	result[0] = d

	return result, nil
}

func checkOspfAreaExists(idArea, version, realm, routingInstance string, junSess *junos.Session,
) (_ bool, err error) {
	var showConfig string
	ospfVersion := junos.OspfV2
	if version == "v3" {
		ospfVersion = junos.OspfV3
	} else if realm != "" {
		return false, fmt.Errorf("realm can't set if version != v3")
	}
	switch {
	case routingInstance == junos.DefaultW && realm == "":
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols " + ospfVersion + " area " + idArea + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
	case routingInstance == junos.DefaultW && realm != "":
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols " + ospfVersion + " realm " + realm + " area " + idArea + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
	case realm != "":
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols " + ospfVersion + " realm " + realm + " area " + idArea + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
	default:
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols " + ospfVersion + " area " + idArea + junos.PipeDisplaySet)
		if err != nil {
			return false, err
		}
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setOspfArea(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if d.Get("routing_instance").(string) != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	ospfVersion := junos.OspfV2
	if d.Get("version").(string) == "v3" {
		ospfVersion = junos.OspfV3
	} else if d.Get("realm").(string) != "" {
		return fmt.Errorf("realm can't set if version != v3")
	}
	setPrefix += "protocols " + ospfVersion + " "
	if realm := d.Get("realm").(string); realm != "" {
		setPrefix += "realm " + realm + " "
	}
	setPrefix += "area " + d.Get("area_id").(string) + " "

	interfaceNameList := make([]string, 0)
	for _, v := range d.Get("interface").([]interface{}) {
		ospfInterface := v.(map[string]interface{})
		if bchk.InSlice(ospfInterface["name"].(string), interfaceNameList) {
			return fmt.Errorf("multiple blocks interface with the same name %s", ospfInterface["name"].(string))
		}
		interfaceNameList = append(interfaceNameList, ospfInterface["name"].(string))
		setPrefixInterface := setPrefix + "interface " + ospfInterface["name"].(string) + " "
		configSetInterface, err := setOspfAreaInterface(setPrefixInterface, ospfInterface)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetInterface...)
	}
	areaRangeList := make([]string, 0)
	for _, areaRangeBlock := range d.Get("area_range").(*schema.Set).List() {
		areaRange := areaRangeBlock.(map[string]interface{})
		if bchk.InSlice(areaRange["range"].(string), areaRangeList) {
			return fmt.Errorf("multiple blocks area_range with the same range %s", areaRange["range"].(string))
		}
		areaRangeList = append(areaRangeList, areaRange["range"].(string))
		configSet = append(configSet, setPrefix+"area-range "+areaRange["range"].(string))
		if areaRange["exact"].(bool) {
			configSet = append(configSet,
				setPrefix+"area-range "+areaRange["range"].(string)+" exact",
			)
		}
		if v := areaRange["override_metric"].(int); v != 0 {
			configSet = append(configSet,
				setPrefix+"area-range "+areaRange["range"].(string)+" override-metric "+strconv.Itoa(v),
			)
		}
		if areaRange["restrict"].(bool) {
			configSet = append(configSet,
				setPrefix+"area-range "+areaRange["range"].(string)+" restrict",
			)
		}
	}
	for _, contextIdentifier := range sortSetOfString(d.Get("context_identifier").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"context-identifier "+contextIdentifier)
	}
	for _, interAreaPrefExport := range d.Get("inter_area_prefix_export").([]interface{}) {
		configSet = append(configSet, setPrefix+"inter-area-prefix-export "+interAreaPrefExport.(string))
	}
	for _, interAreaPrefImport := range d.Get("inter_area_prefix_import").([]interface{}) {
		configSet = append(configSet, setPrefix+"inter-area-prefix-import "+interAreaPrefImport.(string))
	}
	for _, networkSummaryExport := range d.Get("network_summary_export").([]interface{}) {
		configSet = append(configSet, setPrefix+"network-summary-export "+networkSummaryExport.(string))
	}
	for _, networkSummaryImport := range d.Get("network_summary_import").([]interface{}) {
		configSet = append(configSet, setPrefix+"network-summary-import "+networkSummaryImport.(string))
	}
	if d.Get("no_context_identifier_advertisement").(bool) {
		configSet = append(configSet, setPrefix+"no-context-identifier-advertisement")
	}
	for _, nssaBlock := range d.Get("nssa").([]interface{}) {
		setPrefixNssa := setPrefix + "nssa "
		configSet = append(configSet, setPrefixNssa)
		if nssaBlock != nil {
			nssa := nssaBlock.(map[string]interface{})
			nssaAreaRangeList := make([]string, 0)
			for _, areaRangeBlock := range nssa["area_range"].(*schema.Set).List() {
				areaRange := areaRangeBlock.(map[string]interface{})
				if bchk.InSlice(areaRange["range"].(string), nssaAreaRangeList) {
					return fmt.Errorf("multiple blocks area_range with the same range %s in nssa block", areaRange["range"].(string))
				}
				nssaAreaRangeList = append(nssaAreaRangeList, areaRange["range"].(string))
				configSet = append(configSet, setPrefixNssa+"area-range "+areaRange["range"].(string))
				if areaRange["exact"].(bool) {
					configSet = append(configSet,
						setPrefixNssa+"area-range "+areaRange["range"].(string)+" exact",
					)
				}
				if v := areaRange["override_metric"].(int); v != 0 {
					configSet = append(configSet,
						setPrefixNssa+"area-range "+areaRange["range"].(string)+" override-metric "+strconv.Itoa(v),
					)
				}
				if areaRange["restrict"].(bool) {
					configSet = append(configSet,
						setPrefixNssa+"area-range "+areaRange["range"].(string)+" restrict",
					)
				}
			}
			for _, defaultLsaBlock := range nssa["default_lsa"].([]interface{}) {
				configSet = append(configSet, setPrefixNssa+"default-lsa")
				if defaultLsaBlock != nil {
					defaultLsa := defaultLsaBlock.(map[string]interface{})
					if v := defaultLsa["default_metric"].(int); v != 0 {
						configSet = append(configSet, setPrefixNssa+"default-lsa default-metric "+strconv.Itoa(v))
					}
					if v := defaultLsa["metric_type"].(int); v != 0 {
						configSet = append(configSet, setPrefixNssa+"default-lsa metric-type "+strconv.Itoa(v))
					}
					if defaultLsa["type_7"].(bool) {
						configSet = append(configSet, setPrefixNssa+"default-lsa type-7")
					}
				}
			}
			if nssa["no_summaries"].(bool) {
				configSet = append(configSet, setPrefixNssa+"no-summaries")
			}
			if nssa["summaries"].(bool) {
				configSet = append(configSet, setPrefixNssa+"summaries")
			}
		}
	}
	for _, stubBlock := range d.Get("stub").([]interface{}) {
		configSet = append(configSet, setPrefix+"stub")
		if stubBlock != nil {
			stub := stubBlock.(map[string]interface{})
			if v := stub["default_metric"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"stub default-metric "+strconv.Itoa(v))
			}
			if stub["no_summaries"].(bool) {
				configSet = append(configSet, setPrefix+"stub no-summaries")
			}
			if stub["summaries"].(bool) {
				configSet = append(configSet, setPrefix+"stub summaries")
			}
		}
	}
	virtualLinkList := make([]string, 0)
	for _, virtualLinkBlock := range d.Get("virtual_link").(*schema.Set).List() {
		virtualLink := virtualLinkBlock.(map[string]interface{})
		if bchk.InSlice(
			virtualLink["neighbor_id"].(string)+junos.IDSeparator+virtualLink["transit_area"].(string),
			virtualLinkList,
		) {
			return fmt.Errorf(
				"multiple blocks virtual_link with the same neighbor_id '%s' and transit_area '%s'",
				virtualLink["neighbor_id"].(string), virtualLink["transit_area"].(string))
		}
		virtualLinkList = append(
			virtualLinkList, virtualLink["neighbor_id"].(string)+junos.IDSeparator+virtualLink["transit_area"].(string),
		)
		setPrefixVirtualLink := setPrefix + "virtual-link " +
			"neighbor-id " + virtualLink["neighbor_id"].(string) +
			" transit-area " + virtualLink["transit_area"].(string) + " "
		configSet = append(configSet, setPrefixVirtualLink)
		if v := virtualLink["dead_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixVirtualLink+"dead-interval "+strconv.Itoa(v))
		}
		if virtualLink["demand_circuit"].(bool) {
			configSet = append(configSet, setPrefixVirtualLink+"demand-circuit")
		}
		if virtualLink["disable"].(bool) {
			configSet = append(configSet, setPrefixVirtualLink+"disable")
		}
		if virtualLink["flood_reduction"].(bool) {
			configSet = append(configSet, setPrefixVirtualLink+"flood-reduction")
		}
		if v := virtualLink["hello_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixVirtualLink+"hello-interval "+strconv.Itoa(v))
		}
		if v := virtualLink["ipsec_sa"].(string); v != "" {
			configSet = append(configSet, setPrefixVirtualLink+"ipsec-sa \""+v+"\"")
		}
		if v := virtualLink["mtu"].(int); v != 0 {
			configSet = append(configSet, setPrefixVirtualLink+"mtu "+strconv.Itoa(v))
		}
		if v := virtualLink["retransmit_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixVirtualLink+"retransmit-interval "+strconv.Itoa(v))
		}
		if v := virtualLink["transit_delay"].(int); v != 0 {
			configSet = append(configSet, setPrefixVirtualLink+"transit-delay "+strconv.Itoa(v))
		}
	}

	return junSess.ConfigSet(configSet)
}

func setOspfAreaInterface(setPrefix string, ospfInterface map[string]interface{}) ([]string, error) {
	configSet := make([]string, 0)

	configSet = append(configSet, setPrefix)
	if v := ospfInterface["authentication_simple_password"].(string); v != "" {
		if len(ospfInterface["authentication_md5"].([]interface{})) != 0 {
			return configSet, fmt.Errorf("conflict between 'authentication_simple_password' and 'authentication_md5'"+
				" in interface '%s'", ospfInterface["name"].(string))
		}
		configSet = append(configSet, setPrefix+"authentication simple-password \""+v+"\"")
	}
	authenticationMD5List := make([]int, 0)
	for _, mAuthMD5 := range ospfInterface["authentication_md5"].([]interface{}) {
		authMD5 := mAuthMD5.(map[string]interface{})
		if bchk.InSlice(authMD5["key_id"].(int), authenticationMD5List) {
			return configSet, fmt.Errorf("multiple blocks authentication_md5 with the same key_id %d in interface with name %s",
				authMD5["key_id"].(int), ospfInterface["name"].(string))
		}
		authenticationMD5List = append(authenticationMD5List, authMD5["key_id"].(int))
		configSet = append(configSet, setPrefix+"authentication md5 "+
			strconv.Itoa(authMD5["key_id"].(int))+" key \""+authMD5["key"].(string)+"\"")
		if v := authMD5["start_time"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication md5 "+
				strconv.Itoa(authMD5["key_id"].(int))+" start-time "+v)
		}
	}
	bandwidthBasedMetricsList := make([]string, 0)
	for _, mBandBaseMet := range ospfInterface["bandwidth_based_metrics"].(*schema.Set).List() {
		bandwidthBaseMetrics := mBandBaseMet.(map[string]interface{})
		if bchk.InSlice(bandwidthBaseMetrics["bandwidth"].(string), bandwidthBasedMetricsList) {
			return configSet, fmt.Errorf("multiple blocks bandwidth_based_metrics "+
				"with the same bandwidth %s in interface with name %s",
				bandwidthBaseMetrics["bandwidth"].(string), ospfInterface["name"].(string))
		}
		bandwidthBasedMetricsList = append(bandwidthBasedMetricsList, bandwidthBaseMetrics["bandwidth"].(string))
		configSet = append(configSet, setPrefix+"bandwidth-based-metrics bandwidth "+
			bandwidthBaseMetrics["bandwidth"].(string)+" metric "+strconv.Itoa(bandwidthBaseMetrics["metric"].(int)))
	}
	for _, mBFDLivDet := range ospfInterface["bfd_liveness_detection"].([]interface{}) {
		if mBFDLivDet == nil {
			return configSet, fmt.Errorf("bfd_liveness_detection block is empty in interface %s", ospfInterface["name"].(string))
		}
		setPrefixBfd := setPrefix + "bfd-liveness-detection "
		bfdLiveDetect := mBFDLivDet.(map[string]interface{})
		if v := bfdLiveDetect["authentication_algorithm"].(string); v != "" {
			configSet = append(configSet, setPrefixBfd+"authentication algorithm "+v)
		}
		if v := bfdLiveDetect["authentication_key_chain"].(string); v != "" {
			configSet = append(configSet, setPrefixBfd+"authentication key-chain \""+v+"\"")
		}
		if bfdLiveDetect["authentication_loose_check"].(bool) {
			configSet = append(configSet, setPrefixBfd+"authentication loose-check")
		}
		if v := bfdLiveDetect["detection_time_threshold"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"detection-time threshold "+strconv.Itoa(v))
		}
		if bfdLiveDetect["full_neighbors_only"].(bool) {
			configSet = append(configSet, setPrefixBfd+"full-neighbors-only")
		}
		if v := bfdLiveDetect["holddown_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"holddown-interval "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["minimum_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"minimum-interval "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["minimum_receive_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"minimum-receive-interval "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["multiplier"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+"multiplier "+strconv.Itoa(v))
		}
		if bfdLiveDetect["no_adaptation"].(bool) {
			configSet = append(configSet, setPrefixBfd+"no-adaptation")
		}
		if v := bfdLiveDetect["transmit_interval_minimum_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+
				"transmit-interval minimum-interval "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["transmit_interval_threshold"].(int); v != 0 {
			configSet = append(configSet, setPrefixBfd+
				"transmit-interval threshold "+strconv.Itoa(v))
		}
		if v := bfdLiveDetect["version"].(string); v != "" {
			configSet = append(configSet, setPrefixBfd+"version "+v)
		}
		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefixBfd) {
			return configSet, fmt.Errorf("bfd_liveness_detection block is empty in interface %s", ospfInterface["name"].(string))
		}
	}
	if v := ospfInterface["dead_interval"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"dead-interval "+strconv.Itoa(v))
	}
	if ospfInterface["demand_circuit"].(bool) {
		configSet = append(configSet, setPrefix+"demand-circuit")
	}
	if ospfInterface["disable"].(bool) {
		configSet = append(configSet, setPrefix+"disable")
	}
	if ospfInterface["dynamic_neighbors"].(bool) {
		configSet = append(configSet, setPrefix+"dynamic-neighbors")
	}
	if ospfInterface["flood_reduction"].(bool) {
		configSet = append(configSet, setPrefix+"flood-reduction")
	}
	if v := ospfInterface["hello_interval"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"hello-interval "+strconv.Itoa(v))
	}
	if v := ospfInterface["interface_type"].(string); v != "" {
		configSet = append(configSet, setPrefix+"interface-type "+v)
	}
	if v := ospfInterface["ipsec_sa"].(string); v != "" {
		configSet = append(configSet, setPrefix+"ipsec-sa \""+v+"\"")
	}
	if t := ospfInterface["ipv4_adjacency_segment_protected_type"].(string); t != "" {
		if v := ospfInterface["ipv4_adjacency_segment_protected_value"].(string); v != "" {
			configSet = append(configSet, setPrefix+"ipv4-adjacency-segment protected "+t+" "+v)
		} else {
			configSet = append(configSet, setPrefix+"ipv4-adjacency-segment protected "+t)
		}
	} else if ospfInterface["ipv4_adjacency_segment_protected_value"].(string) != "" {
		return configSet, fmt.Errorf("ipv4_adjacency_segment_protected_type need to be set with " +
			"ipv4_adjacency_segment_protected_value")
	}
	if t := ospfInterface["ipv4_adjacency_segment_unprotected_type"].(string); t != "" {
		if v := ospfInterface["ipv4_adjacency_segment_unprotected_value"].(string); v != "" {
			configSet = append(configSet, setPrefix+"ipv4-adjacency-segment unprotected "+t+" "+v)
		} else {
			configSet = append(configSet, setPrefix+"ipv4-adjacency-segment unprotected "+t)
		}
	} else if ospfInterface["ipv4_adjacency_segment_unprotected_value"].(string) != "" {
		return configSet, fmt.Errorf("ipv4_adjacency_segment_unprotected_type need to be set with " +
			"ipv4_adjacency_segment_unprotected_value")
	}
	if ospfInterface["link_protection"].(bool) {
		configSet = append(configSet, setPrefix+"link-protection")
	}
	if v := ospfInterface["metric"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"metric "+strconv.Itoa(v))
	}
	if v := ospfInterface["mtu"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"mtu "+strconv.Itoa(v))
	}
	neighborList := make([]string, 0)
	for _, mNeighbor := range ospfInterface["neighbor"].(*schema.Set).List() {
		neighbor := mNeighbor.(map[string]interface{})
		if bchk.InSlice(neighbor["address"].(string), neighborList) {
			return configSet, fmt.Errorf("multiple blocks neighbor with the same address %s in interface with name %s",
				neighbor["address"].(string), ospfInterface["name"].(string))
		}
		neighborList = append(neighborList, neighbor["address"].(string))
		configSet = append(configSet, setPrefix+"neighbor "+neighbor["address"].(string))
		if neighbor["eligible"].(bool) {
			configSet = append(configSet, setPrefix+"neighbor "+neighbor["address"].(string)+" eligible")
		}
	}
	if ospfInterface["no_advertise_adjacency_segment"].(bool) {
		configSet = append(configSet, setPrefix+"no-advertise-adjacency-segment")
	}
	if ospfInterface["no_eligible_backup"].(bool) {
		configSet = append(configSet, setPrefix+"no-eligible-backup")
	}
	if ospfInterface["no_eligible_remote_backup"].(bool) {
		configSet = append(configSet, setPrefix+"no-eligible-remote-backup")
	}
	if ospfInterface["no_interface_state_traps"].(bool) {
		configSet = append(configSet, setPrefix+"no-interface-state-traps")
	}
	if ospfInterface["no_neighbor_down_notification"].(bool) {
		configSet = append(configSet, setPrefix+"no-neighbor-down-notification")
	}
	if ospfInterface["node_link_protection"].(bool) {
		configSet = append(configSet, setPrefix+"node-link-protection")
	}
	if ospfInterface["passive"].(bool) {
		configSet = append(configSet, setPrefix+"passive")
		if v := ospfInterface["passive_traffic_engineering_remote_node_id"].(string); v != "" {
			configSet = append(configSet, setPrefix+"passive traffic-engineering remote-node-id "+v)
		}
		if v := ospfInterface["passive_traffic_engineering_remote_node_router_id"].(string); v != "" {
			configSet = append(configSet, setPrefix+"passive traffic-engineering remote-node-router-id "+v)
		}
	} else if ospfInterface["passive_traffic_engineering_remote_node_id"].(string) != "" ||
		ospfInterface["passive_traffic_engineering_remote_node_router_id"].(string) != "" {
		return configSet, fmt.Errorf("passive need to be true with " +
			"passive_traffic_engineering_remote_node_id and passive_traffic_engineering_remote_node_router_id")
	}
	if v := ospfInterface["poll_interval"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"poll-interval "+strconv.Itoa(v))
	}
	if v := ospfInterface["priority"].(int); v != -1 {
		configSet = append(configSet, setPrefix+"priority "+strconv.Itoa(v))
	}
	if v := ospfInterface["retransmit_interval"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"retransmit-interval "+strconv.Itoa(v))
	}
	if ospfInterface["secondary"].(bool) {
		configSet = append(configSet, setPrefix+"secondary")
	}
	if ospfInterface["strict_bfd"].(bool) {
		configSet = append(configSet, setPrefix+"strict-bfd")
	}
	if v := ospfInterface["te_metric"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"te-metric "+strconv.Itoa(v))
	}
	if v := ospfInterface["transit_delay"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"transit-delay "+strconv.Itoa(v))
	}

	return configSet, nil
}

func readOspfArea(idArea, version, realm, routingInstance string, junSess *junos.Session,
) (confRead ospfAreaOptions, err error) {
	var showConfig string
	ospfVersion := junos.OspfV2
	if version == "v3" {
		ospfVersion = junos.OspfV3
	} else if realm != "" {
		return confRead, fmt.Errorf("realm can't set if version != v3")
	}
	switch {
	case routingInstance == junos.DefaultW && realm == "":
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols " + ospfVersion + " area " + idArea + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	case routingInstance == junos.DefaultW && realm != "":
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"protocols " + ospfVersion + " realm " + realm + " area " + idArea + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	case realm != "":
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols " + ospfVersion + " realm " + realm + " area " + idArea + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	default:
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"protocols " + ospfVersion + " area " + idArea + junos.PipeDisplaySetRelative)
		if err != nil {
			return confRead, err
		}
	}

	if showConfig != junos.EmptyW {
		confRead.areaID = idArea
		confRead.realm = realm
		confRead.version = version
		confRead.routingInstance = routingInstance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "interface "):
				itemTrimFields := strings.Split(itemTrim, " ")
				interfaceOptions := map[string]interface{}{
					"name":                                       itemTrimFields[0],
					"authentication_simple_password":             "",
					"authentication_md5":                         make([]map[string]interface{}, 0),
					"bandwidth_based_metrics":                    make([]map[string]interface{}, 0),
					"bfd_liveness_detection":                     make([]map[string]interface{}, 0),
					"dead_interval":                              0,
					"demand_circuit":                             false,
					"disable":                                    false,
					"dynamic_neighbors":                          false,
					"flood_reduction":                            false,
					"hello_interval":                             0,
					"interface_type":                             "",
					"ipsec_sa":                                   "",
					"ipv4_adjacency_segment_protected_type":      "",
					"ipv4_adjacency_segment_protected_value":     "",
					"ipv4_adjacency_segment_unprotected_type":    "",
					"ipv4_adjacency_segment_unprotected_value":   "",
					"link_protection":                            false,
					"metric":                                     0,
					"mtu":                                        0,
					"neighbor":                                   make([]map[string]interface{}, 0),
					"no_advertise_adjacency_segment":             false,
					"no_eligible_backup":                         false,
					"no_eligible_remote_backup":                  false,
					"no_interface_state_traps":                   false,
					"no_neighbor_down_notification":              false,
					"node_link_protection":                       false,
					"passive":                                    false,
					"passive_traffic_engineering_remote_node_id": "",
					"passive_traffic_engineering_remote_node_router_id": "",
					"poll_interval":       0,
					"priority":            -1,
					"retransmit_interval": 0,
					"secondary":           false,
					"strict_bfd":          false,
					"te_metric":           0,
					"transit_delay":       0,
				}
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				confRead.interFace = copyAndRemoveItemMapList("name", interfaceOptions, confRead.interFace)
				if err := readOspfAreaInterface(itemTrim, interfaceOptions); err != nil {
					return confRead, err
				}
				confRead.interFace = append(confRead.interFace, interfaceOptions)
			case balt.CutPrefixInString(&itemTrim, "area-range "):
				itemTrimFields := strings.Split(itemTrim, " ")
				areaRange := map[string]interface{}{
					"range":           itemTrimFields[0],
					"exact":           false,
					"override_metric": 0,
					"restrict":        false,
				}
				confRead.areaRange = copyAndRemoveItemMapList("range", areaRange, confRead.areaRange)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case itemTrim == "exact":
					areaRange["exact"] = true
				case balt.CutPrefixInString(&itemTrim, "override-metric "):
					areaRange["override_metric"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case itemTrim == "restrict":
					areaRange["restrict"] = true
				}
				confRead.areaRange = append(confRead.areaRange, areaRange)
			case balt.CutPrefixInString(&itemTrim, "context-identifier "):
				confRead.contextIdentifier = append(confRead.contextIdentifier, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "inter-area-prefix-export "):
				confRead.interAreaPrefixExport = append(confRead.interAreaPrefixExport, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "inter-area-prefix-import "):
				confRead.interAreaPrefixImport = append(confRead.interAreaPrefixImport, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "network-summary-export "):
				confRead.networkSummaryExport = append(confRead.networkSummaryExport, itemTrim)
			case balt.CutPrefixInString(&itemTrim, "network-summary-import "):
				confRead.networkSummaryImport = append(confRead.networkSummaryImport, itemTrim)
			case itemTrim == "no-context-identifier-advertisement":
				confRead.noContextIdentifierAdvertisement = true
			case balt.CutPrefixInString(&itemTrim, "nssa"):
				if len(confRead.nssa) == 0 {
					confRead.nssa = append(confRead.nssa, map[string]interface{}{
						"area_range":   make([]map[string]interface{}, 0),
						"default_lsa":  make([]map[string]interface{}, 0),
						"no_summaries": false,
						"summaries":    false,
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " area-range "):
					itemTrimFields := strings.Split(itemTrim, " ")
					areaRange := map[string]interface{}{
						"range":           itemTrimFields[0],
						"exact":           false,
						"override_metric": 0,
						"restrict":        false,
					}
					confRead.nssa[0]["area_range"] = copyAndRemoveItemMapList(
						"range",
						areaRange,
						confRead.nssa[0]["area_range"].([]map[string]interface{}),
					)
					balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
					switch {
					case itemTrim == "exact":
						areaRange["exact"] = true
					case balt.CutPrefixInString(&itemTrim, "override-metric "):
						areaRange["override_metric"], err = strconv.Atoi(itemTrim)
						if err != nil {
							return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
						}
					case itemTrim == "restrict":
						areaRange["restrict"] = true
					}
					confRead.nssa[0]["area_range"] = append(confRead.nssa[0]["area_range"].([]map[string]interface{}), areaRange)
				case balt.CutPrefixInString(&itemTrim, " default-lsa"):
					if len(confRead.nssa[0]["default_lsa"].([]map[string]interface{})) == 0 {
						confRead.nssa[0]["default_lsa"] = append(
							confRead.nssa[0]["default_lsa"].([]map[string]interface{}),
							map[string]interface{}{
								"default_metric": 0,
								"metric_type":    0,
								"type_7":         false,
							})
					}
					defaultLsa := confRead.nssa[0]["default_lsa"].([]map[string]interface{})[0]
					switch {
					case balt.CutPrefixInString(&itemTrim, " default-metric "):
						defaultLsa["default_metric"], err = strconv.Atoi(itemTrim)
						if err != nil {
							return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
						}
					case balt.CutPrefixInString(&itemTrim, " metric-type "):
						defaultLsa["metric_type"], err = strconv.Atoi(itemTrim)
						if err != nil {
							return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
						}
					case itemTrim == " type-7":
						defaultLsa["type_7"] = true
					}
				case itemTrim == " no-summaries":
					confRead.nssa[0]["no_summaries"] = true
				case itemTrim == " summaries":
					confRead.nssa[0]["summaries"] = true
				}
			case balt.CutPrefixInString(&itemTrim, "stub"):
				if len(confRead.stub) == 0 {
					confRead.stub = append(confRead.stub, map[string]interface{}{
						"default_metric": 0,
						"no_summaries":   false,
						"summaries":      false,
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " default-metric "):
					confRead.stub[0]["default_metric"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
					}
				case itemTrim == " no-summaries":
					confRead.stub[0]["no_summaries"] = true
				case itemTrim == " summaries":
					confRead.stub[0]["summaries"] = true
				}
			case balt.CutPrefixInString(&itemTrim, "virtual-link "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 4 { // neighbor-id <neighbor_id> transit-area <transit_area>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "virtual-link", itemTrim)
				}
				virtualLink := map[string]interface{}{
					"neighbor_id":         itemTrimFields[1],
					"transit_area":        itemTrimFields[3],
					"dead_interval":       0,
					"demand_circuit":      false,
					"disable":             false,
					"flood_reduction":     false,
					"hello_interval":      0,
					"ipsec_sa":            "",
					"mtu":                 0,
					"retransmit_interval": 0,
					"transit_delay":       0,
				}
				confRead.virtualLink = copyAndRemoveItemMapList2("neighbor_id", "transit_area", virtualLink, confRead.virtualLink)
				balt.CutPrefixInString(&itemTrim, "neighbor-id "+itemTrimFields[1]+" transit-area "+itemTrimFields[3]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "dead-interval "):
					virtualLink["dead_interval"], err = strconv.Atoi(itemTrim)
				case itemTrim == "demand-circuit":
					virtualLink["demand_circuit"] = true
				case itemTrim == "disable":
					virtualLink["disable"] = true
				case itemTrim == "flood-reduction":
					virtualLink["flood_reduction"] = true
				case balt.CutPrefixInString(&itemTrim, "hello-interval "):
					virtualLink["hello_interval"], err = strconv.Atoi(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "ipsec-sa "):
					virtualLink["ipsec_sa"] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "mtu "):
					virtualLink["mtu"], err = strconv.Atoi(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "retransmit-interval "):
					virtualLink["retransmit_interval"], err = strconv.Atoi(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "transit-delay "):
					virtualLink["transit_delay"], err = strconv.Atoi(itemTrim)
				}
				if err != nil {
					return confRead, fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
				}
				confRead.virtualLink = append(confRead.virtualLink, virtualLink)
			}
		}
	}

	return confRead, nil
}

func readOspfAreaInterface(itemTrim string, interfaceOptions map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication simple-password "):
		interfaceOptions["authentication_simple_password"], err = jdecode.Decode(strings.Trim(itemTrim, "\""))
		if err != nil {
			return fmt.Errorf("failed to decode authentication simple-password: %w", err)
		}
	case balt.CutPrefixInString(&itemTrim, "authentication md5 "):
		itemTrimFields := strings.Split(itemTrim, " ")
		keyID, err := strconv.Atoi(itemTrimFields[0])
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrimFields[0], err)
		}
		authMD5 := map[string]interface{}{
			"key_id":     keyID,
			"key":        "",
			"start_time": "",
		}
		interfaceOptions["authentication_md5"] = copyAndRemoveItemMapList("key_id", authMD5,
			interfaceOptions["authentication_md5"].([]map[string]interface{}))
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		switch {
		case balt.CutPrefixInString(&itemTrim, "key "):
			authMD5["key"], err = jdecode.Decode(strings.Trim(itemTrim, "\""))
			if err != nil {
				return fmt.Errorf("failed to decode authentication md5 key: %w", err)
			}
		case balt.CutPrefixInString(&itemTrim, "start-time "):
			authMD5["start_time"] = strings.Split(strings.Trim(itemTrim, "\""), " ")[0]
		}
		interfaceOptions["authentication_md5"] = append(
			interfaceOptions["authentication_md5"].([]map[string]interface{}), authMD5)
	case balt.CutPrefixInString(&itemTrim, "bandwidth-based-metrics bandwidth "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 3 { // <bandwidth> metric <metric>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "bandwidth-based-metrics bandwidth", itemTrim)
		}
		metric, err := strconv.Atoi(itemTrimFields[2])
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrimFields[2], err)
		}
		interfaceOptions["bandwidth_based_metrics"] = append(
			interfaceOptions["bandwidth_based_metrics"].([]map[string]interface{}), map[string]interface{}{
				"bandwidth": itemTrimFields[0],
				"metric":    metric,
			})
	case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
		if len(interfaceOptions["bfd_liveness_detection"].([]map[string]interface{})) == 0 {
			interfaceOptions["bfd_liveness_detection"] = append(
				interfaceOptions["bfd_liveness_detection"].([]map[string]interface{}), map[string]interface{}{
					"authentication_algorithm":           "",
					"authentication_key_chain":           "",
					"authentication_loose_check":         false,
					"detection_time_threshold":           0,
					"full_neighbors_only":                false,
					"holddown_interval":                  0,
					"minimum_interval":                   0,
					"minimum_receive_interval":           0,
					"multiplier":                         0,
					"no_adaptation":                      false,
					"transmit_interval_minimum_interval": 0,
					"transmit_interval_threshold":        0,
					"version":                            "",
				})
		}
		if err := readOspfAreaInterfaceBfd(itemTrim,
			interfaceOptions["bfd_liveness_detection"].([]map[string]interface{})[0]); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "dead-interval "):
		interfaceOptions["dead_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "demand-circuit":
		interfaceOptions["demand_circuit"] = true
	case itemTrim == junos.DisableW:
		interfaceOptions["disable"] = true
	case itemTrim == "dynamic-neighbors":
		interfaceOptions["dynamic_neighbors"] = true
	case itemTrim == "flood-reduction":
		interfaceOptions["flood_reduction"] = true
	case balt.CutPrefixInString(&itemTrim, "hello-interval "):
		interfaceOptions["hello_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "interface-type "):
		interfaceOptions["interface_type"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "ipsec-sa "):
		interfaceOptions["ipsec_sa"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "ipv4-adjacency-segment protected "):
		itemTrimFields := strings.Split(itemTrim, " ")
		interfaceOptions["ipv4_adjacency_segment_protected_type"] = itemTrimFields[0]
		if len(itemTrimFields) > 1 { // <type> <value>
			interfaceOptions["ipv4_adjacency_segment_protected_value"] = itemTrimFields[1]
		}
	case balt.CutPrefixInString(&itemTrim, "ipv4-adjacency-segment unprotected "):
		itemTrimFields := strings.Split(itemTrim, " ")
		interfaceOptions["ipv4_adjacency_segment_unprotected_type"] = itemTrimFields[0]
		if len(itemTrimFields) > 1 { // <type> <value>
			interfaceOptions["ipv4_adjacency_segment_unprotected_value"] = itemTrimFields[1]
		}
	case itemTrim == "link-protection":
		interfaceOptions["link_protection"] = true
	case balt.CutPrefixInString(&itemTrim, "metric "):
		interfaceOptions["metric"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "mtu "):
		interfaceOptions["mtu"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "neighbor "):
		itemTrimFields := strings.Split(itemTrim, " ") // <address> (eligible)?
		address := itemTrimFields[0]
		if len(itemTrimFields) > 1 && itemTrimFields[1] == "eligible" {
			interfaceOptions["neighbor"] = append(
				interfaceOptions["neighbor"].([]map[string]interface{}), map[string]interface{}{
					"address":  address,
					"eligible": true,
				})
		} else {
			interfaceOptions["neighbor"] = append(
				interfaceOptions["neighbor"].([]map[string]interface{}), map[string]interface{}{
					"address":  address,
					"eligible": false,
				})
		}
	case itemTrim == "no-advertise-adjacency-segment":
		interfaceOptions["no_advertise_adjacency_segment"] = true
	case itemTrim == "no-eligible-backup":
		interfaceOptions["no_eligible_backup"] = true
	case itemTrim == "no-eligible-remote-backup":
		interfaceOptions["no_eligible_remote_backup"] = true
	case itemTrim == "no-interface-state-traps":
		interfaceOptions["no_interface_state_traps"] = true
	case itemTrim == "no-neighbor-down-notification":
		interfaceOptions["no_neighbor_down_notification"] = true
	case itemTrim == "node-link-protection":
		interfaceOptions["node_link_protection"] = true
	case balt.CutPrefixInString(&itemTrim, "passive"):
		interfaceOptions["passive"] = true
		switch {
		case balt.CutPrefixInString(&itemTrim, " traffic-engineering remote-node-id "):
			interfaceOptions["passive_traffic_engineering_remote_node_id"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, " traffic-engineering remote-node-router-id "):
			interfaceOptions["passive_traffic_engineering_remote_node_router_id"] = itemTrim
		}
	case balt.CutPrefixInString(&itemTrim, "poll-interval "):
		interfaceOptions["poll_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "priority "):
		interfaceOptions["priority"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "retransmit-interval "):
		interfaceOptions["retransmit_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "secondary":
		interfaceOptions["secondary"] = true
	case itemTrim == "strict-bfd":
		interfaceOptions["strict_bfd"] = true
	case balt.CutPrefixInString(&itemTrim, "te-metric "):
		interfaceOptions["te_metric"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "transit-delay "):
		interfaceOptions["transit_delay"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	}

	return nil
}

func readOspfAreaInterfaceBfd(itemTrim string, bfd map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication algorithm "):
		bfd["authentication_algorithm"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "authentication key-chain "):
		bfd["authentication_key_chain"] = strings.Trim(itemTrim, "\"")
	case itemTrim == "authentication loose-check":
		bfd["authentication_loose_check"] = true
	case balt.CutPrefixInString(&itemTrim, "detection-time threshold "):
		bfd["detection_time_threshold"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "full-neighbors-only":
		bfd["full_neighbors_only"] = true
	case balt.CutPrefixInString(&itemTrim, "holddown-interval "):
		bfd["holddown_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-interval "):
		bfd["minimum_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-receive-interval "):
		bfd["minimum_receive_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "multiplier "):
		bfd["multiplier"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "no-adaptation":
		bfd["no_adaptation"] = true
	case balt.CutPrefixInString(&itemTrim, "transmit-interval minimum-interval "):
		bfd["transmit_interval_minimum_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "transmit-interval threshold "):
		bfd["transmit_interval_threshold"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(junos.FailedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "version "):
		bfd["version"] = itemTrim
	}

	return nil
}

func delOspfArea(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	ospfVersion := junos.OspfV2
	if d.Get("version").(string) == "v3" {
		ospfVersion = junos.OspfV3
	} else if d.Get("realm").(string) != "" {
		return fmt.Errorf("realm can't set if version != v3")
	}
	switch {
	case d.Get("routing_instance").(string) == junos.DefaultW && d.Get("realm").(string) == "":
		configSet = append(configSet, junos.DeleteW+
			" protocols "+ospfVersion+
			" area "+d.Get("area_id").(string))
	case d.Get("routing_instance").(string) == junos.DefaultW && d.Get("realm").(string) != "":
		configSet = append(configSet, junos.DeleteW+
			" protocols "+ospfVersion+
			" realm "+d.Get("realm").(string)+
			" area "+d.Get("area_id").(string))
	case d.Get("realm").(string) != "":
		configSet = append(configSet, junos.DelRoutingInstances+d.Get("routing_instance").(string)+
			" protocols "+ospfVersion+
			" realm "+d.Get("realm").(string)+
			" area "+d.Get("area_id").(string))
	default:
		configSet = append(configSet, junos.DelRoutingInstances+d.Get("routing_instance").(string)+
			" protocols "+ospfVersion+
			" area "+d.Get("area_id").(string))
	}

	return junSess.ConfigSet(configSet)
}

func fillOspfAreaData(d *schema.ResourceData, ospfAreaOptions ospfAreaOptions) {
	if tfErr := d.Set("area_id", ospfAreaOptions.areaID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", ospfAreaOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("realm", ospfAreaOptions.realm); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version", ospfAreaOptions.version); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface", ospfAreaOptions.interFace); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("area_range", ospfAreaOptions.areaRange); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("context_identifier", ospfAreaOptions.contextIdentifier); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inter_area_prefix_export", ospfAreaOptions.interAreaPrefixExport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inter_area_prefix_import", ospfAreaOptions.interAreaPrefixImport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("network_summary_export", ospfAreaOptions.networkSummaryExport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("network_summary_import", ospfAreaOptions.networkSummaryImport); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"no_context_identifier_advertisement",
		ospfAreaOptions.noContextIdentifierAdvertisement,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("nssa", ospfAreaOptions.nssa); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("stub", ospfAreaOptions.stub); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("virtual_link", ospfAreaOptions.virtualLink); tfErr != nil {
		panic(tfErr)
	}
}
