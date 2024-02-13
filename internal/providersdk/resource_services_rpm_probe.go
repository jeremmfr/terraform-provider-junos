package providersdk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

type rpmProbeOptions struct {
	delegateProbes bool
	name           string
	test           []map[string]interface{}
}

func resourceServicesRpmProbe() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceServicesRpmProbeCreate,
		ReadWithoutTimeout:   resourceServicesRpmProbeRead,
		UpdateWithoutTimeout: resourceServicesRpmProbeUpdate,
		DeleteWithoutTimeout: resourceServicesRpmProbeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServicesRpmProbeImport,
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
							ValidateFunc: validateIsIPv6Address,
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
										ValidateFunc: validateIsIPv6Address,
									},
									"source_inet6_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"source_inet6_step": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateIsIPv6Address,
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
										ValidateFunc: validateIsIPv6Address,
									},
									"target_inet6_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"target_inet6_step": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateIsIPv6Address,
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

func resourceServicesRpmProbeCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setServicesRpmProbe(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

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
	rpmProbeExists, err := checkServicesRpmProbeExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if rpmProbeExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("services rpm probe %v already exists", d.Get("name").(string)))...)
	}

	if err := setServicesRpmProbe(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_services_rpm_probe")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	rpmProbeExists, err = checkServicesRpmProbeExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if rpmProbeExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("services rpm probe %v "+
			"not exists after commit => check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceServicesRpmProbeReadWJunSess(d, junSess)...)
}

func resourceServicesRpmProbeRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceServicesRpmProbeReadWJunSess(d, junSess)
}

func resourceServicesRpmProbeReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	rpmProbeOptions, err := readServicesRpmProbe(d.Get("name").(string), junSess)
	junos.MutexUnlock()
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

func resourceServicesRpmProbeUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delServicesRpmProbe(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setServicesRpmProbe(d, junSess); err != nil {
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
	if err := delServicesRpmProbe(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setServicesRpmProbe(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_services_rpm_probe")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceServicesRpmProbeReadWJunSess(d, junSess)...)
}

func resourceServicesRpmProbeDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delServicesRpmProbe(d.Get("name").(string), junSess); err != nil {
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
	if err := delServicesRpmProbe(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_services_rpm_probe")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceServicesRpmProbeImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	rpmProbeExists, err := checkServicesRpmProbeExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !rpmProbeExists {
		return nil, fmt.Errorf("don't find services rpm probe with id '%v' (id must be <name>)", d.Id())
	}
	rpmProbeOptions, err := readServicesRpmProbe(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillServicesRpmProbeData(d, rpmProbeOptions)

	result[0] = d

	return result, nil
}

func checkServicesRpmProbeExists(probe string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "services rpm probe \"" + probe + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setServicesRpmProbe(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set services rpm probe \"" + d.Get("name").(string) + "\" "
	configSet = append(configSet, setPrefix)
	if d.Get("delegate_probes").(bool) {
		configSet = append(configSet, setPrefix+"delegate-probes")
	}
	testNameList := make([]string, 0)
	for _, t := range d.Get("test").([]interface{}) {
		test := t.(map[string]interface{})
		if slices.Contains(testNameList, test["name"].(string)) {
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
			if v2int, v2cnt := rpmScale["destination_interface"].(string),
				rpmScale["destination_subunit_cnt"].(int); v2int != "" || v2cnt != 0 {
				if v2int == "" || v2cnt == 0 {
					return errors.New("all of `destination_interface,destination_subunit_cnt` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale destination interface "+v2int)
				configSet = append(configSet, setPrefixTest+"rpm-scale destination subunit-cnt "+strconv.Itoa(v2cnt))
			}
			if v2add, v2cnt, v2step := rpmScale["source_address_base"].(string), rpmScale["source_count"].(int),
				rpmScale["source_step"].(string); v2add != "" || v2cnt != 0 || v2step != "" {
				if v2add == "" || v2cnt == 0 || v2step == "" {
					return errors.New("all of `source_address_base,source_count,source_step` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale source address-base "+v2add)
				configSet = append(configSet, setPrefixTest+"rpm-scale source count "+strconv.Itoa(v2cnt))
				configSet = append(configSet, setPrefixTest+"rpm-scale source step "+v2step)
			}
			if v2add, v2cnt, v2step := rpmScale["source_inet6_address_base"].(string), rpmScale["source_inet6_count"].(int),
				rpmScale["source_inet6_step"].(string); v2add != "" || v2cnt != 0 || v2step != "" {
				if v2add == "" || v2cnt == 0 || v2step == "" {
					return errors.New("all of `source_inet6_address_base,source_inet6_count,source_inet6_step` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale source-inet6 address-base "+v2add)
				configSet = append(configSet, setPrefixTest+"rpm-scale source-inet6 count "+strconv.Itoa(v2cnt))
				configSet = append(configSet, setPrefixTest+"rpm-scale source-inet6 step "+v2step)
			}
			if v2add, v2cnt, v2step := rpmScale["target_address_base"].(string), rpmScale["target_count"].(int),
				rpmScale["target_step"].(string); v2add != "" || v2cnt != 0 || v2step != "" {
				if v2add == "" || v2cnt == 0 || v2step == "" {
					return errors.New("all of `target_address_base,target_count,target_step` must be specified")
				}
				configSet = append(configSet, setPrefixTest+"rpm-scale target address-base "+v2add)
				configSet = append(configSet, setPrefixTest+"rpm-scale target count "+strconv.Itoa(v2cnt))
				configSet = append(configSet, setPrefixTest+"rpm-scale target step "+v2step)
			}
			if v2add, v2cnt, v2step := rpmScale["target_inet6_address_base"].(string), rpmScale["target_inet6_count"].(int),
				rpmScale["target_inet6_step"].(string); v2add != "" || v2cnt != 0 || v2step != "" {
				if v2add == "" || v2cnt == 0 || v2step == "" {
					return errors.New("all of `target_inet6_address_base,target_inet6_count,target_inet6_step` must be specified")
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
				return errors.New("all of `target_type,target_value` must be specified")
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

	return junSess.ConfigSet(configSet)
}

func readServicesRpmProbe(probe string, junSess *junos.Session,
) (confRead rpmProbeOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services rpm probe \"" + probe + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = probe
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "delegate-probes":
				confRead.delegateProbes = true
			case balt.CutPrefixInString(&itemTrim, "test "):
				itemTrimFields := strings.Split(itemTrim, " ")
				test := map[string]interface{}{
					"name":                       strings.Trim(itemTrimFields[0], "\""),
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
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				if err := readServicesRpmProbeTest(itemTrim, test); err != nil {
					return confRead, err
				}
				confRead.test = append(confRead.test, test)
			}
		}
	}

	return confRead, nil
}

func readServicesRpmProbeTest(itemTrim string, test map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "target "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 2 { // <type> <value>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "target", itemTrim)
		}
		test["target_type"] = itemTrimFields[0]
		test["target_value"] = strings.Trim(itemTrimFields[1], "\"")
	case balt.CutPrefixInString(&itemTrim, "data-fill "):
		test["data_fill"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "data-size "):
		test["data_size"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "destination-interface "):
		test["destination_interface"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		test["destination_port"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "dscp-code-points "):
		test["dscp_code_points"] = itemTrim
	case itemTrim == "hardware-timestamp":
		test["hardware_timestamp"] = true
	case balt.CutPrefixInString(&itemTrim, "history-size "):
		test["history_size"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "inet6-options source-address "):
		test["inet6_source_address"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "moving-average-size "):
		test["moving_average_size"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "one-way-hardware-timestamp":
		test["one_way_hardware_timestamp"] = true
	case balt.CutPrefixInString(&itemTrim, "probe-count "):
		test["probe_count"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "probe-interval "):
		test["probe_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "probe-type "):
		test["probe_type"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "routing-instance "):
		test["routing_instance"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "rpm-scale "):
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
		switch {
		case balt.CutPrefixInString(&itemTrim, "tests-count "):
			rpmScale["tests_count"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "destination interface "):
			rpmScale["destination_interface"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "destination subunit-cnt "):
			rpmScale["destination_subunit_cnt"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "source address-base "):
			rpmScale["source_address_base"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "source count "):
			rpmScale["source_count"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "source step "):
			rpmScale["source_step"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "source-inet6 address-base "):
			rpmScale["source_inet6_address_base"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "source-inet6 count "):
			rpmScale["source_inet6_count"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "source-inet6 step "):
			rpmScale["source_inet6_step"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "target address-base "):
			rpmScale["target_address_base"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "target count "):
			rpmScale["target_count"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "target step "):
			rpmScale["target_step"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "target-inet6 address-base "):
			rpmScale["target_inet6_address_base"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "target-inet6 count "):
			rpmScale["target_inet6_count"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "target-inet6 step "):
			rpmScale["target_inet6_step"] = itemTrim
		}
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		test["source_address"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "test-interval "):
		test["test_interval"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case balt.CutPrefixInString(&itemTrim, "thresholds"):
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
		case balt.CutPrefixInString(&itemTrim, " egress-time "):
			thresholds["egress_time"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " ingress-time "):
			thresholds["ingress_time"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " jitter-egress "):
			thresholds["jitter_egress"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " jitter-ingress "):
			thresholds["jitter_ingress"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " jitter-rtt "):
			thresholds["jitter_rtt"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " rtt "):
			thresholds["rtt"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " std-dev-egress "):
			thresholds["std_dev_egress"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " std-dev-ingress "):
			thresholds["std_dev_ingress"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " std-dev-rtt "):
			thresholds["std_dev_rtt"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " successive-loss "):
			thresholds["successive_loss"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " total-loss "):
			thresholds["total_loss"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case balt.CutPrefixInString(&itemTrim, "traps "):
		test["traps"] = append(test["traps"].([]string), itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ttl "):
		test["ttl"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	}

	return nil
}

func delServicesRpmProbe(probe string, junSess *junos.Session) error {
	configSet := []string{"delete services rpm probe \"" + probe + "\""}

	return junSess.ConfigSet(configSet)
}

func fillServicesRpmProbeData(d *schema.ResourceData, rpmProbeOptions rpmProbeOptions,
) {
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
