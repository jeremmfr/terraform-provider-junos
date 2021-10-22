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
)

type rpmProbeOptions struct {
	delegateProbes bool
	name           string
	test           []map[string]interface{}
}

func resourceServicesRpmProbe() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServicesRpmProbeCreate,
		ReadContext:   resourceServicesRpmProbeRead,
		UpdateContext: resourceServicesRpmProbeUpdate,
		DeleteContext: resourceServicesRpmProbeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServicesRpmProbeImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"delegate_probes": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"test": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringDoesNotContainAny(" "),
						},
						"data_fill": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(
								`^[0-9a-fA-F]+$`), "must be hexadecimal digits (0-9, a-f, A-F)"),
						},
						"data_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 65400),
						},
						"destination_interface": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(7, 65535),
						},
						"dscp_code_points": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hardware_timestamp": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"history_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 512),
						},
						"inet6_source_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPv6Address,
						},
						"moving_average_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 1024),
						},
						"one_way_hardware_timestamp": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"probe_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 15),
						},
						"probe_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 255),
						},
						"probe_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http-get", "http-metadata-get",
								"icmp-ping", "icmp-ping-timestamp", "icmp6-ping",
								"tcp-ping",
								"udp-ping", "udp-ping-timestamp",
							}, false),
						},
						"routing_instance": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
						},
						"rpm_scale": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tests_count": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"destination_interface": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"destination_subunit_cnt": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"source_address_base": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"source_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"source_step": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"source_inet6_address_base": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv6Address,
									},
									"source_inet6_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"source_inet6_step": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv6Address,
									},
									"target_address_base": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"target_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"target_step": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
									"target_inet6_address_base": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv6Address,
									},
									"target_inet6_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"target_inet6_step": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv6Address,
									},
								},
							},
						},
						"source_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPv4Address,
						},
						"target_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"address", "inet6-address", "inet6-url", "url"}, false),
						},
						"target_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"test_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
						"thresholds": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"egress_time": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"ingress_time": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"jitter_egress": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"jitter_ingress": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"jitter_rtt": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"rtt": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"std_dev_egress": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"std_dev_ingress": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"std_dev_rtt": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 60000000),
									},
									"successive_loss": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 15),
									},
									"total_loss": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 15),
									},
								},
							},
						},
						"traps": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ttl": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 254),
						},
					},
				},
			},
		},
	}
}

func resourceServicesRpmProbeCreate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	if sess.junosFakeCreateSetFile != "" {
		if err := setServicesRpmProbe(d, m, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	rpmProbeExists, err := checkServicesRpmProbeExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if rpmProbeExists {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services rpm probe %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesRpmProbe(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("create resource junos_services_rpm_probe", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	rpmProbeExists, err = checkServicesRpmProbeExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if rpmProbeExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services rpm probe %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesRpmProbeReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesRpmProbeRead(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceServicesRpmProbeReadWJnprSess(d, m, jnprSess)
}

func resourceServicesRpmProbeReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	rpmProbeOptions, err := readServicesRpmProbe(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if rpmProbeOptions.name == "" {
		d.SetId("")
	} else {
		fillServicesRpmProbeData(d, rpmProbeOptions)
	}

	return nil
}

func resourceServicesRpmProbeUpdate(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delServicesRpmProbe(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesRpmProbe(d, m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("update resource junos_services_rpm_probe", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesRpmProbeReadWJnprSess(d, m, jnprSess)...)
}

func resourceServicesRpmProbeDelete(ctx context.Context,
	d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	var diagWarns diag.Diagnostics
	if err := delServicesRpmProbe(d.Get("name").(string), m, jnprSess); err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := sess.commitConf("delete resource junos_services_rpm_probe", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, sess.configClear(jnprSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesRpmProbeImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	rpmProbeExists, err := checkServicesRpmProbeExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !rpmProbeExists {
		return nil, fmt.Errorf("don't find services rpm probe with id '%v' (id must be <name>)", d.Id())
	}
	rpmProbeOptions, err := readServicesRpmProbe(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillServicesRpmProbeData(d, rpmProbeOptions)

	result[0] = d

	return result, nil
}

func checkServicesRpmProbeExists(probe string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	showConfig, err := sess.command("show configuration services rpm probe \""+probe+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if showConfig == emptyWord {
		return false, nil
	}

	return true, nil
}

func setServicesRpmProbe(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set services rpm probe \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix)
	if d.Get("delegate_probes").(bool) {
		configSet = append(configSet, setPrefix+"delegate-probes")
	}
	testNameList := make([]string, 0)
	for _, t := range d.Get("test").([]interface{}) {
		test := t.(map[string]interface{})
		if bchk.StringInSlice(test["name"].(string), testNameList) {
			return fmt.Errorf("multiple blocks test with the same name %s", test["name"].(string))
		}
		testNameList = append(testNameList, test["name"].(string))
		setPrefixTest := setPrefix + "test \"" + test["name"].(string) + "\" "
		configSet = append(configSet, setPrefixTest)
		if v := test["data_fill"].(string); v != "" {
			configSet = append(configSet, setPrefixTest+"data-fill "+v)
		}
		if v := test["data_size"].(int); v != -1 {
			configSet = append(configSet, setPrefixTest+"data-size "+strconv.Itoa(v))
		}
		if v := test["destination_interface"].(string); v != "" {
			configSet = append(configSet, setPrefixTest+"destination-interface "+v)
		}
		if v := test["destination_port"].(int); v != 0 {
			configSet = append(configSet, setPrefixTest+"destination-port "+strconv.Itoa(v))
		}
		if v := test["dscp_code_points"].(string); v != "" {
			configSet = append(configSet, setPrefixTest+"dscp-code-points "+v)
		}
		if test["hardware_timestamp"].(bool) {
			configSet = append(configSet, setPrefixTest+"hardware-timestamp")
		}
		if v := test["history_size"].(int); v != -1 {
			configSet = append(configSet, setPrefixTest+"history-size "+strconv.Itoa(v))
		}
		if v := test["inet6_source_address"].(string); v != "" {
			configSet = append(configSet, setPrefixTest+"inet6-options source-address "+v)
		}
		if v := test["moving_average_size"].(int); v != -1 {
			configSet = append(configSet, setPrefixTest+"moving-average-size "+strconv.Itoa(v))
		}
		if test["one_way_hardware_timestamp"].(bool) {
			configSet = append(configSet, setPrefixTest+"one-way-hardware-timestamp")
		}
		if v := test["probe_count"].(int); v != 0 {
			configSet = append(configSet, setPrefixTest+"probe-count "+strconv.Itoa(v))
		}
		if v := test["probe_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixTest+"probe-interval "+strconv.Itoa(v))
		}
		if v := test["probe_type"].(string); v != "" {
			configSet = append(configSet, setPrefixTest+"probe-type "+v)
		}
		if v := test["routing_instance"].(string); v != "" {
			configSet = append(configSet, setPrefixTest+"routing-instance \""+v+"\"")
		}
		for _, v := range test["rpm_scale"].([]interface{}) {
			rpmScale := v.(map[string]interface{})
			configSet = append(configSet,
				setPrefixTest+"rpm-scale tests-count "+strconv.Itoa(rpmScale["tests_count"].(int)))
			if v2int, v2cnt :=
				rpmScale["destination_interface"].(string),
				rpmScale["destination_subunit_cnt"].(int); v2int != "" || v2cnt != 0 {
				if v2int == "" || v2cnt == 0 {
					return fmt.Errorf("all of `destination_interface,destination_subunit_cnt` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale destination interface "+v2int)
				configSet = append(configSet, setPrefixTest+"rpm-scale destination subunit-cnt "+strconv.Itoa(v2cnt))
			}
			if v2add, v2cnt, v2step :=
				rpmScale["source_address_base"].(string),
				rpmScale["source_count"].(int),
				rpmScale["source_step"].(string); v2add != "" || v2cnt != 0 || v2step != "" {
				if v2add == "" || v2cnt == 0 || v2step == "" {
					return fmt.Errorf("all of `source_address_base,source_count,source_step` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale source address-base "+v2add)
				configSet = append(configSet, setPrefixTest+"rpm-scale source count "+strconv.Itoa(v2cnt))
				configSet = append(configSet, setPrefixTest+"rpm-scale source step "+v2step)
			}
			if v2add, v2cnt, v2step :=
				rpmScale["source_inet6_address_base"].(string),
				rpmScale["source_inet6_count"].(int),
				rpmScale["source_inet6_step"].(string); v2add != "" || v2cnt != 0 || v2step != "" {
				if v2add == "" || v2cnt == 0 || v2step == "" {
					return fmt.Errorf("all of `source_inet6_address_base,source_inet6_count,source_inet6_step` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale source-inet6 address-base "+v2add)
				configSet = append(configSet, setPrefixTest+"rpm-scale source-inet6 count "+strconv.Itoa(v2cnt))
				configSet = append(configSet, setPrefixTest+"rpm-scale source-inet6 step "+v2step)
			}
			if v2add, v2cnt, v2step :=
				rpmScale["target_address_base"].(string),
				rpmScale["target_count"].(int),
				rpmScale["target_step"].(string); v2add != "" || v2cnt != 0 || v2step != "" {
				if v2add == "" || v2cnt == 0 || v2step == "" {
					return fmt.Errorf("all of `target_address_base,target_count,target_step` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale target address-base "+v2add)
				configSet = append(configSet, setPrefixTest+"rpm-scale target count "+strconv.Itoa(v2cnt))
				configSet = append(configSet, setPrefixTest+"rpm-scale target step "+v2step)
			}
			if v2add, v2cnt, v2step :=
				rpmScale["target_inet6_address_base"].(string),
				rpmScale["target_inet6_count"].(int),
				rpmScale["target_inet6_step"].(string); v2add != "" || v2cnt != 0 || v2step != "" {
				if v2add == "" || v2cnt == 0 || v2step == "" {
					return fmt.Errorf("all of `target_inet6_address_base,target_inet6_count,target_inet6_step` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale target-inet6 address-base "+v2add)
				configSet = append(configSet, setPrefixTest+"rpm-scale target-inet6 count "+strconv.Itoa(v2cnt))
				configSet = append(configSet, setPrefixTest+"rpm-scale target-inet6 step "+v2step)
			}
		}
		if v := test["source_address"].(string); v != "" {
			configSet = append(configSet, setPrefixTest+"source-address "+v)
		}
		if v, v2 := test["target_type"].(string), test["target_value"].(string); v != "" || v2 != "" {
			if v == "" || v2 == "" {
				return fmt.Errorf("all of `target_type,target_value` must be specified")
			}
			configSet = append(configSet, setPrefixTest+"target "+v+" \""+v2+"\"")
		}
		if v := test["test_interval"].(int); v != -1 {
			configSet = append(configSet, setPrefixTest+"test-interval "+strconv.Itoa(v))
		}
		for _, v := range test["thresholds"].([]interface{}) {
			thresholds := v.(map[string]interface{})
			configSet = append(configSet, setPrefixTest+"thresholds")
			if v2 := thresholds["egress_time"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds egress-time "+strconv.Itoa(v2))
			}
			if v2 := thresholds["ingress_time"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds ingress-time "+strconv.Itoa(v2))
			}
			if v2 := thresholds["jitter_egress"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds jitter-egress "+strconv.Itoa(v2))
			}
			if v2 := thresholds["jitter_ingress"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds jitter-ingress "+strconv.Itoa(v2))
			}
			if v2 := thresholds["jitter_rtt"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds jitter-rtt "+strconv.Itoa(v2))
			}
			if v2 := thresholds["rtt"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds rtt "+strconv.Itoa(v2))
			}
			if v2 := thresholds["std_dev_egress"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds std-dev-egress "+strconv.Itoa(v2))
			}
			if v2 := thresholds["std_dev_ingress"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds std-dev-ingress "+strconv.Itoa(v2))
			}
			if v2 := thresholds["std_dev_rtt"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds std-dev-rtt "+strconv.Itoa(v2))
			}
			if v2 := thresholds["successive_loss"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds successive-loss "+strconv.Itoa(v2))
			}
			if v2 := thresholds["total_loss"].(int); v2 != -1 {
				configSet = append(configSet, setPrefixTest+"thresholds total-loss "+strconv.Itoa(v2))
			}
		}
		for _, v := range sortSetOfString(test["traps"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixTest+"traps "+v)
		}
		if v := test["ttl"].(int); v != 0 {
			configSet = append(configSet, setPrefixTest+"ttl "+strconv.Itoa(v))
		}
	}

	return sess.configSet(configSet, jnprSess)
}

func readServicesRpmProbe(probe string, m interface{}, jnprSess *NetconfObject) (
	rpmProbeOptions, error) {
	sess := m.(*Session)
	var confRead rpmProbeOptions

	showConfig, err := sess.command("show configuration services rpm probe \""+probe+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyWord {
		confRead.name = probe
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case itemTrim == "delegate-probes":
				confRead.delegateProbes = true
			case strings.HasPrefix(itemTrim, "test "):
				lineCut := strings.Split(itemTrim, " ")
				test := map[string]interface{}{
					"name":                       strings.Trim(lineCut[1], "\""),
					"target_type":                "",
					"target_value":               "",
					"data_fill":                  "",
					"data_size":                  -1,
					"destination_interface":      "",
					"destination_port":           0,
					"dscp_code_points":           "",
					"hardware_timestamp":         false,
					"history_size":               -1,
					"inet6_source_address":       "",
					"moving_average_size":        -1,
					"one_way_hardware_timestamp": false,
					"probe_count":                0,
					"probe_interval":             0,
					"probe_type":                 "",
					"routing_instance":           "",
					"rpm_scale":                  make([]map[string]interface{}, 0),
					"source_address":             "",
					"test_interval":              -1,
					"thresholds":                 make([]map[string]interface{}, 0),
					"traps":                      make([]string, 0),
					"ttl":                        0,
				}
				confRead.test = copyAndRemoveItemMapList("name", test, confRead.test)
				itemTrimTest := strings.TrimPrefix(itemTrim, "test "+lineCut[1]+" ")
				if err := readServicesRpmProbeTest(itemTrimTest, test); err != nil {
					return confRead, err
				}
				confRead.test = append(confRead.test, test)
			}
		}
	}

	return confRead, nil
}

func readServicesRpmProbeTest(itemTrim string, test map[string]interface{}) error {
	var err error
	switch {
	case strings.HasPrefix(itemTrim, "target "):
		itemTrimSplit := strings.Split(itemTrim, " ")
		if len(itemTrimSplit) != 3 {
			return fmt.Errorf("can't read words in line for target : %s", itemTrim)
		}
		test["target_type"] = itemTrimSplit[1]
		test["target_value"] = strings.Trim(itemTrimSplit[2], "\"")
	case strings.HasPrefix(itemTrim, "data-fill "):
		test["data_fill"] = strings.TrimPrefix(itemTrim, "data-fill ")
	case strings.HasPrefix(itemTrim, "data-size "):
		test["data_size"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "data-size "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "destination-interface "):
		test["destination_interface"] = strings.TrimPrefix(itemTrim, "destination-interface ")
	case strings.HasPrefix(itemTrim, "destination-port "):
		test["destination_port"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "destination-port "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "dscp-code-points "):
		test["dscp_code_points"] = strings.TrimPrefix(itemTrim, "dscp-code-points ")
	case itemTrim == "hardware-timestamp":
		test["hardware_timestamp"] = true
	case strings.HasPrefix(itemTrim, "history-size "):
		test["history_size"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "history-size "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "inet6-options source-address "):
		test["inet6_source_address"] = strings.TrimPrefix(itemTrim, "inet6-options source-address ")
	case strings.HasPrefix(itemTrim, "moving-average-size "):
		test["moving_average_size"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "moving-average-size "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case itemTrim == "one-way-hardware-timestamp":
		test["one_way_hardware_timestamp"] = true
	case strings.HasPrefix(itemTrim, "probe-count "):
		test["probe_count"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "probe-count "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "probe-interval "):
		test["probe_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "probe-interval "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "probe-type "):
		test["probe_type"] = strings.TrimPrefix(itemTrim, "probe-type ")
	case strings.HasPrefix(itemTrim, "routing-instance "):
		test["routing_instance"] = strings.TrimPrefix(itemTrim, "routing-instance ")
	case strings.HasPrefix(itemTrim, "rpm-scale "):
		if len(test["rpm_scale"].([]map[string]interface{})) == 0 {
			test["rpm_scale"] = append(test["rpm_scale"].([]map[string]interface{}), map[string]interface{}{
				"tests_count":               0,
				"destination_interface":     "",
				"destination_subunit_cnt":   0,
				"source_address_base":       "",
				"source_count":              0,
				"source_step":               "",
				"source_inet6_address_base": "",
				"source_inet6_count":        0,
				"source_inet6_step":         "",
				"target_address_base":       "",
				"target_count":              0,
				"target_step":               "",
				"target_inet6_address_base": "",
				"target_inet6_count":        0,
				"target_inet6_step":         "",
			})
		}
		rpmScale := test["rpm_scale"].([]map[string]interface{})[0]
		itemTrimRpmScale := strings.TrimPrefix(itemTrim, "rpm-scale ")
		switch {
		case strings.HasPrefix(itemTrimRpmScale, "tests-count "):
			rpmScale["tests_count"], err = strconv.Atoi(strings.TrimPrefix(itemTrimRpmScale, "tests-count "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimRpmScale, "destination interface "):
			rpmScale["destination_interface"] = strings.TrimPrefix(itemTrimRpmScale, "destination interface ")
		case strings.HasPrefix(itemTrimRpmScale, "destination subunit-cnt "):
			rpmScale["destination_subunit_cnt"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrimRpmScale, "destination subunit-cnt "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimRpmScale, "source address-base "):
			rpmScale["source_address_base"] = strings.TrimPrefix(itemTrimRpmScale, "source address-base ")
		case strings.HasPrefix(itemTrimRpmScale, "source count "):
			rpmScale["source_count"], err = strconv.Atoi(strings.TrimPrefix(itemTrimRpmScale, "source count "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimRpmScale, "source step "):
			rpmScale["source_step"] = strings.TrimPrefix(itemTrimRpmScale, "source step ")
		case strings.HasPrefix(itemTrimRpmScale, "source-inet6 address-base "):
			rpmScale["source_inet6_address_base"] = strings.TrimPrefix(itemTrimRpmScale, "source-inet6 address-base ")
		case strings.HasPrefix(itemTrimRpmScale, "source-inet6 count "):
			rpmScale["source_inet6_count"], err = strconv.Atoi(strings.TrimPrefix(itemTrimRpmScale, "source-inet6 count "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimRpmScale, "source-inet6 step "):
			rpmScale["source_inet6_step"] = strings.TrimPrefix(itemTrimRpmScale, "source-inet6 step ")
		case strings.HasPrefix(itemTrimRpmScale, "target address-base "):
			rpmScale["target_address_base"] = strings.TrimPrefix(itemTrimRpmScale, "target address-base ")
		case strings.HasPrefix(itemTrimRpmScale, "target count "):
			rpmScale["target_count"], err = strconv.Atoi(strings.TrimPrefix(itemTrimRpmScale, "target count "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimRpmScale, "target step "):
			rpmScale["target_step"] = strings.TrimPrefix(itemTrimRpmScale, "target step ")
		case strings.HasPrefix(itemTrimRpmScale, "target-inet6 address-base "):
			rpmScale["target_inet6_address_base"] = strings.TrimPrefix(itemTrimRpmScale, "target-inet6 address-base ")
		case strings.HasPrefix(itemTrimRpmScale, "target-inet6 count "):
			rpmScale["target_inet6_count"], err = strconv.Atoi(strings.TrimPrefix(itemTrimRpmScale, "target-inet6 count "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrimRpmScale, "target-inet6 step "):
			rpmScale["target_inet6_step"] = strings.TrimPrefix(itemTrimRpmScale, "target-inet6 step ")
		}
	case strings.HasPrefix(itemTrim, "source-address "):
		test["source_address"] = strings.TrimPrefix(itemTrim, "source-address ")
	case strings.HasPrefix(itemTrim, "test-interval "):
		test["test_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "test-interval "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "thresholds"):
		if len(test["thresholds"].([]map[string]interface{})) == 0 {
			test["thresholds"] = append(test["thresholds"].([]map[string]interface{}), map[string]interface{}{
				"egress_time":     -1,
				"ingress_time":    -1,
				"jitter_egress":   -1,
				"jitter_ingress":  -1,
				"jitter_rtt":      -1,
				"rtt":             -1,
				"std_dev_egress":  -1,
				"std_dev_ingress": -1,
				"std_dev_rtt":     -1,
				"successive_loss": -1,
				"total_loss":      -1,
			})
		}
		thresholds := test["thresholds"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "thresholds egress-time "):
			thresholds["egress_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds egress-time "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds ingress-time "):
			thresholds["ingress_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds ingress-time "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds jitter-egress "):
			thresholds["jitter_egress"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds jitter-egress "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds jitter-ingress "):
			thresholds["jitter_ingress"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds jitter-ingress "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds jitter-rtt "):
			thresholds["jitter_rtt"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds jitter-rtt "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds rtt "):
			thresholds["rtt"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds rtt "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds std-dev-egress "):
			thresholds["std_dev_egress"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds std-dev-egress "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds std-dev-ingress "):
			thresholds["std_dev_ingress"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds std-dev-ingress "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds std-dev-rtt "):
			thresholds["std_dev_rtt"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds std-dev-rtt "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds successive-loss "):
			thresholds["successive_loss"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds successive-loss "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "thresholds total-loss "):
			thresholds["total_loss"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "thresholds total-loss "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "traps "):
		test["traps"] = append(test["traps"].([]string), strings.TrimPrefix(itemTrim, "traps "))
	case strings.HasPrefix(itemTrim, "ttl "):
		test["ttl"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "ttl "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func delServicesRpmProbe(probe string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := []string{"delete services rpm probe \"" + probe + "\""}

	return sess.configSet(configSet, jnprSess)
}

func fillServicesRpmProbeData(
	d *schema.ResourceData, rpmProbeOptions rpmProbeOptions) {
	if tfErr := d.Set("name", rpmProbeOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("delegate_probes", rpmProbeOptions.delegateProbes); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("test", rpmProbeOptions.test); tfErr != nil {
		panic(tfErr)
	}
}
