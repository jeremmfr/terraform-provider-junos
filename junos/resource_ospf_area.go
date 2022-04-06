package junos

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	jdecode "github.com/jeremmfr/junosdecode"
)

type ospfAreaOptions struct {
	areaID          string
	routingInstance string
	version         string
	interFace       []map[string]interface{}
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
			},
			"routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          defaultW,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
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
		},
	}
}

func resourceOspfAreaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setOspfArea(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("area_id").(string) + idSeparator + d.Get("version").(string) +
			idSeparator + d.Get("routing_instance").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	ospfAreaExists, err := checkOspfAreaExists(d.Get("area_id").(string), d.Get("version").(string),
		d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if ospfAreaExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(fmt.Errorf("ospf %v area %v already exists in routing instance %v",
			d.Get("version").(string), d.Get("area_id").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setOspfArea(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_ospf_area", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	ospfAreaExists, err = checkOspfAreaExists(d.Get("area_id").(string), d.Get("version").(string),
		d.Get("routing_instance").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if ospfAreaExists {
		d.SetId(d.Get("area_id").(string) + idSeparator + d.Get("version").(string) +
			idSeparator + d.Get("routing_instance").(string))
	} else {
		return append(diagWarns,
			diag.FromErr(fmt.Errorf("ospf %v area %v in routing instance %v not exists after commit => check your config",
				d.Get("version").(string), d.Get("area_id").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceOspfAreaReadWJnprSess(d, m, jnprSess)...)
}

func resourceOspfAreaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceOspfAreaReadWJnprSess(d, m, jnprSess)
}

func resourceOspfAreaReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	ospfAreaOptions, err := readOspfArea(d.Get("area_id").(string), d.Get("version").(string),
		d.Get("routing_instance").(string), m, jnprSess)
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
	sess := m.(*Session)
	if sess.junosFakeUpdateAlso {
		if err := delOspfArea(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setOspfArea(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delOspfArea(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setOspfArea(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_ospf_area", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceOspfAreaReadWJnprSess(d, m, jnprSess)...)
}

func resourceOspfAreaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeDeleteAlso {
		if err := delOspfArea(d, m, nil); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if err := sess.configLock(ctx, jnprSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delOspfArea(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_ospf_area", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceOspfAreaImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	ospfAreaExists, err := checkOspfAreaExists(idSplit[0], idSplit[1], idSplit[2], m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !ospfAreaExists {
		return nil, fmt.Errorf("don't find ospf area with id '%v' (id must be "+
			"<aread_id>"+idSeparator+"<version>"+idSeparator+"<routing_instance>)", d.Id())
	}
	ospfAreaOptions, err := readOspfArea(idSplit[0], idSplit[1], idSplit[2], m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillOspfAreaData(d, ospfAreaOptions)
	result[0] = d

	return result, nil
}

func checkOspfAreaExists(idArea, version, routingInstance string, m interface{}, jnprSess *NetconfObject,
) (bool, error) {
	sess := m.(*Session)
	var showConfig string
	var err error
	ospfVersion := ospfV2
	if version == "v3" {
		ospfVersion = ospfV3
	}
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			"protocols "+ospfVersion+" area "+idArea+pipeDisplaySet, jnprSess)
		if err != nil {
			return false, err
		}
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols "+ospfVersion+" area "+idArea+pipeDisplaySet, jnprSess)
		if err != nil {
			return false, err
		}
	}
	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setOspfArea(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)
	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	ospfVersion := ospfV2
	if d.Get("version").(string) == "v3" {
		ospfVersion = ospfV3
	}
	setPrefix += "protocols " + ospfVersion + " area " + d.Get("area_id").(string) + " "

	interfaceNameList := make([]string, 0)
	for _, v := range d.Get("interface").([]interface{}) {
		ospfInterface := v.(map[string]interface{})
		if bchk.StringInSlice(ospfInterface["name"].(string), interfaceNameList) {
			return fmt.Errorf("multiple blocks interface with the same name %s", ospfInterface["name"].(string))
		}
		interfaceNameList = append(interfaceNameList, ospfInterface["name"].(string))
		setPrefixInterface := setPrefix + "interface " + ospfInterface["name"].(string) + " "
		configSet = append(configSet, setPrefixInterface)
		if v := ospfInterface["authentication_simple_password"].(string); v != "" {
			if len(ospfInterface["authentication_md5"].([]interface{})) != 0 {
				return fmt.Errorf("conflict between 'authentication_simple_password' and 'authentication_md5'"+
					" in interface '%s'", ospfInterface["name"].(string))
			}
			configSet = append(configSet, setPrefixInterface+"authentication simple-password \""+v+"\"")
		}
		authenticationMD5List := make([]int, 0)
		for _, mAuthMD5 := range ospfInterface["authentication_md5"].([]interface{}) {
			authMD5 := mAuthMD5.(map[string]interface{})
			if bchk.IntInSlice(authMD5["key_id"].(int), authenticationMD5List) {
				return fmt.Errorf("multiple blocks authentication_md5 with the same key_id %d in interface with name %s",
					authMD5["key_id"].(int), ospfInterface["name"].(string))
			}
			authenticationMD5List = append(authenticationMD5List, authMD5["key_id"].(int))
			configSet = append(configSet, setPrefixInterface+"authentication md5 "+
				strconv.Itoa(authMD5["key_id"].(int))+" key \""+authMD5["key"].(string)+"\"")
			if v := authMD5["start_time"].(string); v != "" {
				configSet = append(configSet, setPrefixInterface+"authentication md5 "+
					strconv.Itoa(authMD5["key_id"].(int))+" start-time "+v)
			}
		}
		bandwidthBasedMetricsList := make([]string, 0)
		for _, mBandBaseMet := range ospfInterface["bandwidth_based_metrics"].(*schema.Set).List() {
			bandwidthBaseMetrics := mBandBaseMet.(map[string]interface{})
			if bchk.StringInSlice(bandwidthBaseMetrics["bandwidth"].(string), bandwidthBasedMetricsList) {
				return fmt.Errorf("multiple blocks bandwidth_based_metrics with the same bandwidth %s in interface with name %s",
					bandwidthBaseMetrics["bandwidth"].(string), ospfInterface["name"].(string))
			}
			bandwidthBasedMetricsList = append(bandwidthBasedMetricsList, bandwidthBaseMetrics["bandwidth"].(string))
			configSet = append(configSet, setPrefixInterface+"bandwidth-based-metrics bandwidth "+
				bandwidthBaseMetrics["bandwidth"].(string)+" metric "+strconv.Itoa(bandwidthBaseMetrics["metric"].(int)))
		}
		for _, mBFDLivDet := range ospfInterface["bfd_liveness_detection"].([]interface{}) {
			if mBFDLivDet == nil {
				return fmt.Errorf("bfd_liveness_detection block is empty in interface %s", ospfInterface["name"].(string))
			}
			setPrefixBfd := setPrefixInterface + "bfd-liveness-detection "
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
				return fmt.Errorf("bfd_liveness_detection block is empty in interface %s", ospfInterface["name"].(string))
			}
		}
		if v := ospfInterface["dead_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixInterface+"dead-interval "+strconv.Itoa(v))
		}
		if ospfInterface["demand_circuit"].(bool) {
			configSet = append(configSet, setPrefixInterface+"demand-circuit")
		}
		if ospfInterface["disable"].(bool) {
			configSet = append(configSet, setPrefixInterface+"disable")
		}
		if ospfInterface["dynamic_neighbors"].(bool) {
			configSet = append(configSet, setPrefixInterface+"dynamic-neighbors")
		}
		if ospfInterface["flood_reduction"].(bool) {
			configSet = append(configSet, setPrefixInterface+"flood-reduction")
		}
		if v := ospfInterface["hello_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixInterface+"hello-interval "+strconv.Itoa(v))
		}
		if v := ospfInterface["interface_type"].(string); v != "" {
			configSet = append(configSet, setPrefixInterface+"interface-type "+v)
		}
		if v := ospfInterface["ipsec_sa"].(string); v != "" {
			configSet = append(configSet, setPrefixInterface+"ipsec-sa \""+v+"\"")
		}
		if t := ospfInterface["ipv4_adjacency_segment_protected_type"].(string); t != "" {
			if v := ospfInterface["ipv4_adjacency_segment_protected_value"].(string); v != "" {
				configSet = append(configSet, setPrefixInterface+"ipv4-adjacency-segment protected "+t+" "+v)
			} else {
				configSet = append(configSet, setPrefixInterface+"ipv4-adjacency-segment protected "+t)
			}
		} else if ospfInterface["ipv4_adjacency_segment_protected_value"].(string) != "" {
			return fmt.Errorf("ipv4_adjacency_segment_protected_type need to be set with " +
				"ipv4_adjacency_segment_protected_value")
		}
		if t := ospfInterface["ipv4_adjacency_segment_unprotected_type"].(string); t != "" {
			if v := ospfInterface["ipv4_adjacency_segment_unprotected_value"].(string); v != "" {
				configSet = append(configSet, setPrefixInterface+"ipv4-adjacency-segment unprotected "+t+" "+v)
			} else {
				configSet = append(configSet, setPrefixInterface+"ipv4-adjacency-segment unprotected "+t)
			}
		} else if ospfInterface["ipv4_adjacency_segment_unprotected_value"].(string) != "" {
			return fmt.Errorf("ipv4_adjacency_segment_unprotected_type need to be set with " +
				"ipv4_adjacency_segment_unprotected_value")
		}
		if ospfInterface["link_protection"].(bool) {
			configSet = append(configSet, setPrefixInterface+"link-protection")
		}
		if v := ospfInterface["metric"].(int); v != 0 {
			configSet = append(configSet, setPrefixInterface+"metric "+strconv.Itoa(v))
		}
		if v := ospfInterface["mtu"].(int); v != 0 {
			configSet = append(configSet, setPrefixInterface+"mtu "+strconv.Itoa(v))
		}
		neighborList := make([]string, 0)
		for _, mNeighbor := range ospfInterface["neighbor"].(*schema.Set).List() {
			neighbor := mNeighbor.(map[string]interface{})
			if bchk.StringInSlice(neighbor["address"].(string), neighborList) {
				return fmt.Errorf("multiple blocks neighbor with the same address %s in interface with name %s",
					neighbor["address"].(string), ospfInterface["name"].(string))
			}
			neighborList = append(neighborList, neighbor["address"].(string))
			configSet = append(configSet, setPrefixInterface+"neighbor "+neighbor["address"].(string))
			if neighbor["eligible"].(bool) {
				configSet = append(configSet, setPrefixInterface+"neighbor "+neighbor["address"].(string)+" eligible")
			}
		}
		if ospfInterface["no_advertise_adjacency_segment"].(bool) {
			configSet = append(configSet, setPrefixInterface+"no-advertise-adjacency-segment")
		}
		if ospfInterface["no_eligible_backup"].(bool) {
			configSet = append(configSet, setPrefixInterface+"no-eligible-backup")
		}
		if ospfInterface["no_eligible_remote_backup"].(bool) {
			configSet = append(configSet, setPrefixInterface+"no-eligible-remote-backup")
		}
		if ospfInterface["no_interface_state_traps"].(bool) {
			configSet = append(configSet, setPrefixInterface+"no-interface-state-traps")
		}
		if ospfInterface["no_neighbor_down_notification"].(bool) {
			configSet = append(configSet, setPrefixInterface+"no-neighbor-down-notification")
		}
		if ospfInterface["node_link_protection"].(bool) {
			configSet = append(configSet, setPrefixInterface+"node-link-protection")
		}
		if ospfInterface["passive"].(bool) {
			configSet = append(configSet, setPrefixInterface+"passive")
			if v := ospfInterface["passive_traffic_engineering_remote_node_id"].(string); v != "" {
				configSet = append(configSet, setPrefixInterface+"passive traffic-engineering remote-node-id "+v)
			}
			if v := ospfInterface["passive_traffic_engineering_remote_node_router_id"].(string); v != "" {
				configSet = append(configSet, setPrefixInterface+"passive traffic-engineering remote-node-router-id "+v)
			}
		} else if ospfInterface["passive_traffic_engineering_remote_node_id"].(string) != "" ||
			ospfInterface["passive_traffic_engineering_remote_node_router_id"].(string) != "" {
			return fmt.Errorf("passive need to be true with " +
				"passive_traffic_engineering_remote_node_id and passive_traffic_engineering_remote_node_router_id")
		}
		if v := ospfInterface["poll_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixInterface+"poll-interval "+strconv.Itoa(v))
		}
		if v := ospfInterface["priority"].(int); v != -1 {
			configSet = append(configSet, setPrefixInterface+"priority "+strconv.Itoa(v))
		}
		if v := ospfInterface["retransmit_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixInterface+"retransmit-interval "+strconv.Itoa(v))
		}
		if ospfInterface["secondary"].(bool) {
			configSet = append(configSet, setPrefixInterface+"secondary")
		}
		if ospfInterface["strict_bfd"].(bool) {
			configSet = append(configSet, setPrefixInterface+"strict-bfd")
		}
		if v := ospfInterface["te_metric"].(int); v != 0 {
			configSet = append(configSet, setPrefixInterface+"te-metric "+strconv.Itoa(v))
		}
		if v := ospfInterface["transit_delay"].(int); v != 0 {
			configSet = append(configSet, setPrefixInterface+"transit-delay "+strconv.Itoa(v))
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readOspfArea(idArea, version, routingInstance string, m interface{}, jnprSess *NetconfObject,
) (ospfAreaOptions, error) {
	sess := m.(*Session)
	var confRead ospfAreaOptions
	var showConfig string
	var err error
	ospfVersion := ospfV2
	if version == "v3" {
		ospfVersion = ospfV3
	}
	if routingInstance == defaultW {
		showConfig, err = sess.command(cmdShowConfig+
			"protocols "+ospfVersion+" area "+idArea+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
	} else {
		showConfig, err = sess.command(cmdShowConfig+routingInstancesWS+routingInstance+" "+
			"protocols "+ospfVersion+" area "+idArea+pipeDisplaySetRelative, jnprSess)
		if err != nil {
			return confRead, err
		}
	}

	if showConfig != emptyW {
		confRead.areaID = idArea
		confRead.version = version
		confRead.routingInstance = routingInstance
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			if strings.HasPrefix(itemTrim, "interface ") {
				itemInterfaceList := strings.Split(strings.TrimPrefix(itemTrim, "interface "), " ")
				interfaceOptions := map[string]interface{}{
					"name":                                       itemInterfaceList[0],
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
				itemTrimInterface := strings.TrimPrefix(itemTrim, "interface "+itemInterfaceList[0]+" ")
				confRead.interFace = copyAndRemoveItemMapList("name", interfaceOptions, confRead.interFace)
				switch {
				case strings.HasPrefix(itemTrimInterface, "authentication simple-password "):
					var err error
					interfaceOptions["authentication_simple_password"], err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
						itemTrimInterface, "authentication simple-password "), "\""))
					if err != nil {
						return confRead, fmt.Errorf("failed to decode authentication simple-password : %w", err)
					}
				case strings.HasPrefix(itemTrimInterface, "authentication md5 "):
					itemTrimInterfaceSplit := strings.Split(strings.TrimPrefix(itemTrimInterface, "authentication md5 "), " ")
					keyID, err := strconv.Atoi(itemTrimInterfaceSplit[0])
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterfaceSplit[0], err)
					}
					authMD5 := map[string]interface{}{
						"key_id":     keyID,
						"key":        "",
						"start_time": "",
					}
					interfaceOptions["authentication_md5"] = copyAndRemoveItemMapList("key_id", authMD5,
						interfaceOptions["authentication_md5"].([]map[string]interface{}))
					itemTrimAuthMD5 := strings.TrimPrefix(itemTrimInterface, "authentication md5 "+itemTrimInterfaceSplit[0]+" ")
					switch {
					case strings.HasPrefix(itemTrimAuthMD5, "key "):
						var err error
						authMD5["key"], err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
							itemTrimAuthMD5, "key "), "\""))
						if err != nil {
							return confRead, fmt.Errorf("failed to decode authentication md5 key : %w", err)
						}
					case strings.HasPrefix(itemTrimAuthMD5, "start-time "):
						authMD5["start_time"] = strings.Split(strings.Trim(strings.TrimPrefix(
							itemTrimAuthMD5, "start-time "), "\""), " ")[0]
					}
					interfaceOptions["authentication_md5"] = append(
						interfaceOptions["authentication_md5"].([]map[string]interface{}), authMD5)
				case strings.HasPrefix(itemTrimInterface, "bandwidth-based-metrics bandwidth "):
					itemTrimInterfaceSplit := strings.Split(strings.TrimPrefix(
						itemTrimInterface, "bandwidth-based-metrics bandwidth "), " ")
					if len(itemTrimInterfaceSplit) < 3 {
						return confRead, fmt.Errorf("can't read values for bandwidth_based_metrics in %s", itemTrimInterface)
					}
					metric, err := strconv.Atoi(itemTrimInterfaceSplit[2])
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterfaceSplit[2], err)
					}
					interfaceOptions["bandwidth_based_metrics"] = append(
						interfaceOptions["bandwidth_based_metrics"].([]map[string]interface{}), map[string]interface{}{
							"bandwidth": itemTrimInterfaceSplit[0],
							"metric":    metric,
						})
				case strings.HasPrefix(itemTrimInterface, "bfd-liveness-detection "):
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
					if err := readOspfAreaInterfaceBfd(strings.TrimPrefix(itemTrimInterface, "bfd-liveness-detection "),
						interfaceOptions["bfd_liveness_detection"].([]map[string]interface{})[0]); err != nil {
						return confRead, err
					}
				case strings.HasPrefix(itemTrimInterface, "dead-interval "):
					interfaceOptions["dead_interval"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "dead-interval "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				case itemTrimInterface == "demand-circuit":
					interfaceOptions["demand_circuit"] = true
				case itemTrimInterface == disableW:
					interfaceOptions["disable"] = true
				case itemTrimInterface == "dynamic-neighbors":
					interfaceOptions["dynamic_neighbors"] = true
				case itemTrimInterface == "flood-reduction":
					interfaceOptions["flood_reduction"] = true
				case strings.HasPrefix(itemTrimInterface, "hello-interval "):
					interfaceOptions["hello_interval"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "hello-interval "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				case strings.HasPrefix(itemTrimInterface, "interface-type "):
					interfaceOptions["interface_type"] = strings.TrimPrefix(itemTrimInterface, "interface-type ")
				case strings.HasPrefix(itemTrimInterface, "ipsec-sa "):
					interfaceOptions["ipsec_sa"] = strings.Trim(strings.TrimPrefix(itemTrimInterface, "ipsec-sa "), "\"")
				case strings.HasPrefix(itemTrimInterface, "ipv4-adjacency-segment protected "):
					itemTrimInterfaceSplit := strings.Split(strings.TrimPrefix(
						itemTrimInterface, "ipv4-adjacency-segment protected "), " ")
					interfaceOptions["ipv4_adjacency_segment_protected_type"] = itemTrimInterfaceSplit[0]
					if len(itemTrimInterfaceSplit) > 1 {
						interfaceOptions["ipv4_adjacency_segment_protected_value"] = itemTrimInterfaceSplit[1]
					}
				case strings.HasPrefix(itemTrimInterface, "ipv4-adjacency-segment unprotected "):
					itemTrimInterfaceSplit := strings.Split(strings.TrimPrefix(
						itemTrimInterface, "ipv4-adjacency-segment unprotected "), " ")
					interfaceOptions["ipv4_adjacency_segment_unprotected_type"] = itemTrimInterfaceSplit[0]
					if len(itemTrimInterfaceSplit) > 1 {
						interfaceOptions["ipv4_adjacency_segment_unprotected_value"] = itemTrimInterfaceSplit[1]
					}
				case itemTrimInterface == "link-protection":
					interfaceOptions["link_protection"] = true
				case strings.HasPrefix(itemTrimInterface, "metric "):
					interfaceOptions["metric"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "metric "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				case strings.HasPrefix(itemTrimInterface, "mtu "):
					interfaceOptions["mtu"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "mtu "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				case strings.HasPrefix(itemTrimInterface, "neighbor "):
					itemTrimInterfaceSplit := strings.Split(strings.TrimPrefix(itemTrimInterface, "neighbor "), " ")
					address := itemTrimInterfaceSplit[0]
					if len(itemTrimInterfaceSplit) > 1 && itemTrimInterfaceSplit[1] == "eligible" {
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
				case itemTrimInterface == "no-advertise-adjacency-segment":
					interfaceOptions["no_advertise_adjacency_segment"] = true
				case itemTrimInterface == "no-eligible-backup":
					interfaceOptions["no_eligible_backup"] = true
				case itemTrimInterface == "no-eligible-remote-backup":
					interfaceOptions["no_eligible_remote_backup"] = true
				case itemTrimInterface == "no-interface-state-traps":
					interfaceOptions["no_interface_state_traps"] = true
				case itemTrimInterface == "no-neighbor-down-notification":
					interfaceOptions["no_neighbor_down_notification"] = true
				case itemTrimInterface == "node-link-protection":
					interfaceOptions["node_link_protection"] = true
				case strings.HasPrefix(itemTrimInterface, "passive"):
					interfaceOptions["passive"] = true
					switch {
					case strings.HasPrefix(itemTrimInterface, "passive traffic-engineering remote-node-id "):
						interfaceOptions["passive_traffic_engineering_remote_node_id"] = strings.TrimPrefix(
							itemTrimInterface, "passive traffic-engineering remote-node-id ")
					case strings.HasPrefix(itemTrimInterface, "passive traffic-engineering remote-node-router-id "):
						interfaceOptions["passive_traffic_engineering_remote_node_router_id"] = strings.TrimPrefix(
							itemTrimInterface, "passive traffic-engineering remote-node-router-id ")
					}
				case strings.HasPrefix(itemTrimInterface, "poll-interval "):
					interfaceOptions["poll_interval"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "poll-interval "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				case strings.HasPrefix(itemTrimInterface, "priority "):
					interfaceOptions["priority"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "priority "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				case strings.HasPrefix(itemTrimInterface, "retransmit-interval "):
					interfaceOptions["retransmit_interval"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "retransmit-interval "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				case itemTrimInterface == "secondary":
					interfaceOptions["secondary"] = true
				case itemTrimInterface == "strict-bfd":
					interfaceOptions["strict_bfd"] = true
				case strings.HasPrefix(itemTrimInterface, "te-metric "):
					interfaceOptions["te_metric"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "te-metric "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				case strings.HasPrefix(itemTrimInterface, "transit-delay "):
					interfaceOptions["transit_delay"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrimInterface, "transit-delay "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrimInterface, err)
					}
				}
				confRead.interFace = append(confRead.interFace, interfaceOptions)
			}
		}
	}

	return confRead, nil
}

func readOspfAreaInterfaceBfd(itemTrim string, bfd map[string]interface{}) error {
	switch {
	case strings.HasPrefix(itemTrim, "authentication algorithm "):
		bfd["authentication_algorithm"] = strings.TrimPrefix(itemTrim, "authentication algorithm ")
	case strings.HasPrefix(itemTrim, "authentication key-chain "):
		bfd["authentication_key_chain"] = strings.Trim(strings.TrimPrefix(itemTrim, "authentication key-chain "), "\"")
	case itemTrim == "authentication loose-check":
		bfd["authentication_loose_check"] = true
	case strings.HasPrefix(itemTrim, "detection-time threshold "):
		var err error
		bfd["detection_time_threshold"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "detection-time threshold "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "full-neighbors-only":
		bfd["full_neighbors_only"] = true
	case strings.HasPrefix(itemTrim, "holddown-interval "):
		var err error
		bfd["holddown_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "holddown-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "minimum-interval "):
		var err error
		bfd["minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "minimum-receive-interval "):
		var err error
		bfd["minimum_receive_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-receive-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "multiplier "):
		var err error
		bfd["multiplier"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "multiplier "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "no-adaptation":
		bfd["no_adaptation"] = true
	case strings.HasPrefix(itemTrim, "transmit-interval minimum-interval "):
		var err error
		bfd["transmit_interval_minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "transmit-interval minimum-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "transmit-interval threshold "):
		var err error
		bfd["transmit_interval_threshold"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "transmit-interval threshold "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "version "):
		bfd["version"] = strings.TrimPrefix(itemTrim, "version ")
	}

	return nil
}

func delOspfArea(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	ospfVersion := ospfV2
	if d.Get("version").(string) == "v3" {
		ospfVersion = ospfV3
	}
	if d.Get("routing_instance").(string) == defaultW {
		configSet = append(configSet, "delete protocols "+ospfVersion+" area "+d.Get("area_id").(string))
	} else {
		configSet = append(configSet, delRoutingInstances+d.Get("routing_instance").(string)+
			" protocols "+ospfVersion+" area "+d.Get("area_id").(string))
	}

	return sess.configSet(configSet, jnprSess)
}

func fillOspfAreaData(d *schema.ResourceData, ospfAreaOptions ospfAreaOptions) {
	if tfErr := d.Set("area_id", ospfAreaOptions.areaID); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface", ospfAreaOptions.interFace); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance", ospfAreaOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version", ospfAreaOptions.version); tfErr != nil {
		panic(tfErr)
	}
}
