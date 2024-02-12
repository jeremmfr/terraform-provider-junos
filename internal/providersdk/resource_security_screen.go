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

const (
	userDefinedOptionTypeRegex = `^([1-9]|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])( to ([1-9]|[1-9]\d|1\d\d|2[0-4]\d|25[0-5]))?$`
	userDefinedHeaderTypeRegex = `^(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])( to ([1-9]|[1-9]\d|1\d\d|2[0-4]\d|25[0-5]))?$`
)

type screenOptions struct {
	alarmWithoutDrop bool
	name             string
	description      string
	icmp             []map[string]interface{}
	ip               []map[string]interface{}
	limitSession     []map[string]interface{}
	tcp              []map[string]interface{}
	udp              []map[string]interface{}
}

func resourceSecurityScreen() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityScreenCreate,
		ReadWithoutTimeout:   resourceSecurityScreenRead,
		UpdateWithoutTimeout: resourceSecurityScreenUpdate,
		DeleteWithoutTimeout: resourceSecurityScreenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityScreenImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"alarm_without_drop": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"icmp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flood": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 1000000),
									},
								},
							},
						},
						"fragment": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"icmpv6_malformed": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"large": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ping_death": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sweep": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1000, 1000000),
									},
								},
							},
						},
					},
				},
			},
			"ip": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bad_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"block_frag": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ipv6_extension_header": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ah_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"esp_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"hip_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"destination_header": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ilnp_nonce_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"home_address_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"line_identification_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"tunnel_encapsulation_limit_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"user_defined_option_type": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														ValidateFunc: validation.StringMatch(regexp.MustCompile(userDefinedOptionTypeRegex),
															"doesn't match '(1..255)' or '(1..255) to (1..255)'"),
													},
												},
											},
										},
									},
									"fragment_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"hop_by_hop_header": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"calipso_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"rpl_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"smf_dpd_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"jumbo_payload_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"quick_start_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"router_alert_option": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"user_defined_option_type": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
														ValidateFunc: validation.StringMatch(regexp.MustCompile(userDefinedOptionTypeRegex),
															"doesn't match '(1..255)' or '(1..255) to (1..255)'"),
													},
												},
											},
										},
									},
									"mobility_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_next_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"routing_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"shim6_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"user_defined_header_type": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringMatch(regexp.MustCompile(userDefinedHeaderTypeRegex),
												"doesn't match '(0..255)' or '(0..255) to (0..255)'"),
										},
									},
								},
							},
						},
						"ipv6_extension_header_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 32),
						},
						"ipv6_malformed_header": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"loose_source_route_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"record_route_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"security_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"source_route_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"spoofing": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"stream_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"strict_source_route_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tear_drop": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"timestamp_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tunnel": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bad_inner_header": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"gre": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"gre_4in4": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"gre_4in6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"gre_6in4": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"gre_6in6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"ip_in_udp_teredo": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ipip": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipip_4in4": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ipip_4in6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ipip_6in4": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ipip_6in6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ipip_6over4": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ipip_6to4relay": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"dslite": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"isatap": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"unknown_protocol": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"limit_session": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_ip_based": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 2000000),
						},
						"source_ip_based": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 2000000),
						},
					},
				},
			},
			"tcp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fin_no_ack": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"land": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"no_flag": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"port_scan": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1000, 1000000),
									},
								},
							},
						},
						"syn_ack_ack_proxy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 250000),
									},
								},
							},
						},
						"syn_fin": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"syn_flood": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alarm_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"attack_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 500000),
									},
									"destination_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(4, 500000),
									},
									"source_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(4, 500000),
									},
									"timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 50),
									},
									"whitelist": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:             schema.TypeString,
													Required:         true,
													ValidateDiagFunc: validateNameObjectJunos([]string{}, 32, formatDefault),
												},
												"destination_address": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type:             schema.TypeString,
														ValidateDiagFunc: validateCIDRNetworkFunc(),
													},
												},
												"source_address": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type:             schema.TypeString,
														ValidateDiagFunc: validateCIDRNetworkFunc(),
													},
												},
											},
										},
									},
								},
							},
						},
						"syn_frag": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sweep": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1000, 1000000),
									},
								},
							},
						},
						"winnuke": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"udp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flood": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 1000000),
									},
									"whitelist": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"port_scan": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1000, 1000000),
									},
								},
							},
						},
						"sweep": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1000, 1000000),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSecurityScreenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityScreen(d, junSess); err != nil {
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
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security screen not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	securityScreenExists, err := checkSecurityScreenExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityScreenExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf("security screen %v already exists", d.Get("name").(string)))...)
	}

	if err := setSecurityScreen(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "create resource junos_security_screen")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityScreenExists, err = checkSecurityScreenExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityScreenExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security screen %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityScreenReadWJunSess(d, junSess)...)
}

func resourceSecurityScreenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityScreenReadWJunSess(d, junSess)
}

func resourceSecurityScreenReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	screenOptions, err := readSecurityScreen(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if screenOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityScreenData(d, screenOptions)
	}

	return nil
}

func resourceSecurityScreenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityScreen(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityScreen(d, junSess); err != nil {
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
	if err := delSecurityScreen(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityScreen(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "update resource junos_security_screen")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityScreenReadWJunSess(d, junSess)...)
}

func resourceSecurityScreenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityScreen(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityScreen(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf(ctx, "delete resource junos_security_screen")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityScreenImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	securityScreenExists, err := checkSecurityScreenExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityScreenExists {
		return nil, fmt.Errorf("don't find screen with id '%v' (id must be <name>)", d.Id())
	}
	screenOptions, err := readSecurityScreen(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityScreenData(d, screenOptions)

	result[0] = d

	return result, nil
}

func checkSecurityScreenExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security screen ids-option \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityScreen(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security screen ids-option \"" + d.Get("name").(string) + "\" "

	if d.Get("alarm_without_drop").(bool) {
		configSet = append(configSet, setPrefix+"alarm-without-drop")
	}
	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+"description \""+d.Get("description").(string)+"\"")
	}
	for _, v := range d.Get("icmp").([]interface{}) {
		if v == nil {
			return errors.New("icmp block is empty")
		}
		icmp := v.(map[string]interface{})
		configSet = append(configSet, setSecurityScreenIcmp(icmp, setPrefix)...)
	}
	for _, v := range d.Get("ip").([]interface{}) {
		ipSet, err := setSecurityScreenIP(v.(map[string]interface{}), setPrefix)
		if err != nil {
			return err
		}
		configSet = append(configSet, ipSet...)
	}
	for _, v := range d.Get("limit_session").([]interface{}) {
		if v == nil {
			return errors.New("limit_session block is empty")
		}
		limitSession := v.(map[string]interface{})
		if limitSession["destination_ip_based"].(int) != 0 {
			configSet = append(configSet, setPrefix+"limit-session destination-ip-based "+
				strconv.Itoa(limitSession["destination_ip_based"].(int)))
		}
		if limitSession["source_ip_based"].(int) != 0 {
			configSet = append(configSet, setPrefix+"limit-session source-ip-based "+
				strconv.Itoa(limitSession["source_ip_based"].(int)))
		}
	}
	for _, v := range d.Get("tcp").([]interface{}) {
		if v == nil {
			return errors.New("tcp block is empty")
		}
		tcp := v.(map[string]interface{})
		configSetTCP, err := setSecurityScreenTCP(tcp, setPrefix)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetTCP...)
	}
	for _, v := range d.Get("udp").([]interface{}) {
		if v == nil {
			return errors.New("udp block is empty")
		}
		udp := v.(map[string]interface{})
		configSet = append(configSet, setSecurityScreenUDP(udp, setPrefix)...)
	}

	return junSess.ConfigSet(configSet)
}

func setSecurityScreenIcmp(icmp map[string]interface{}, setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "icmp "
	for _, v := range icmp["flood"].([]interface{}) {
		configSet = append(configSet, setPrefix+"flood")
		if v != nil {
			icmpFlood := v.(map[string]interface{})
			if icmpFlood["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"flood threshold "+
					strconv.Itoa(icmpFlood["threshold"].(int)))
			}
		}
	}
	if icmp["fragment"].(bool) {
		configSet = append(configSet, setPrefix+"fragment")
	}
	if icmp["icmpv6_malformed"].(bool) {
		configSet = append(configSet, setPrefix+"icmpv6-malformed")
	}
	if icmp["large"].(bool) {
		configSet = append(configSet, setPrefix+"large")
	}
	if icmp["ping_death"].(bool) {
		configSet = append(configSet, setPrefix+"ping-death")
	}
	for _, v2 := range icmp["sweep"].([]interface{}) {
		configSet = append(configSet, setPrefix+"ip-sweep")
		if v2 != nil {
			icmpSweep := v2.(map[string]interface{})
			if icmpSweep["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"ip-sweep threshold "+
					strconv.Itoa(icmpSweep["threshold"].(int)))
			}
		}
	}

	return configSet
}

func setSecurityScreenIP(ip map[string]interface{}, setPrefix string) ([]string, error) {
	configSet := make([]string, 0)
	setPrefix += "ip "
	if ip["bad_option"].(bool) {
		configSet = append(configSet, setPrefix+"bad-option")
	}
	if ip["block_frag"].(bool) {
		configSet = append(configSet, setPrefix+"block-frag")
	}
	for _, v := range ip["ipv6_extension_header"].([]interface{}) {
		if v == nil {
			return configSet, errors.New("ip.0.ipv6_extension_header block is empty")
		}
		ipIPv6ExtHeader := v.(map[string]interface{})
		if ipIPv6ExtHeader["ah_header"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header AH-header")
		}
		if ipIPv6ExtHeader["esp_header"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header ESP-header")
		}
		if ipIPv6ExtHeader["hip_header"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header HIP-header")
		}
		for _, v2 := range ipIPv6ExtHeader["destination_header"].([]interface{}) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header destination-header")
			if v2 != nil {
				ipIPv6ExtHeaderDestHeader := v2.(map[string]interface{})
				if ipIPv6ExtHeaderDestHeader["ilnp_nonce_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header destination-header ILNP-nonce-option")
				}
				if ipIPv6ExtHeaderDestHeader["home_address_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header destination-header home-address-option")
				}
				if ipIPv6ExtHeaderDestHeader["line_identification_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header destination-header line-identification-option")
				}
				if ipIPv6ExtHeaderDestHeader["tunnel_encapsulation_limit_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header destination-header "+
						"tunnel-encapsulation-limit-option")
				}
				for _, v3 := range ipIPv6ExtHeaderDestHeader["user_defined_option_type"].([]interface{}) {
					configSet = append(configSet,
						setPrefix+"ipv6-extension-header destination-header user-defined-option-type "+v3.(string))
				}
			}
		}
		if ipIPv6ExtHeader["fragment_header"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header fragment-header")
		}
		for _, v2 := range ipIPv6ExtHeader["hop_by_hop_header"].([]interface{}) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header hop-by-hop-header")
			if v2 != nil {
				ipIPv6ExtHeaderHopByHopHeader := v2.(map[string]interface{})
				if ipIPv6ExtHeaderHopByHopHeader["calipso_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header hop-by-hop-header CALIPSO-option")
				}
				if ipIPv6ExtHeaderHopByHopHeader["rpl_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header hop-by-hop-header RPL-option")
				}
				if ipIPv6ExtHeaderHopByHopHeader["smf_dpd_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header hop-by-hop-header SMF-DPD-option")
				}
				if ipIPv6ExtHeaderHopByHopHeader["jumbo_payload_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header hop-by-hop-header jumbo-payload-option")
				}
				if ipIPv6ExtHeaderHopByHopHeader["quick_start_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header hop-by-hop-header quick-start-option")
				}
				if ipIPv6ExtHeaderHopByHopHeader["router_alert_option"].(bool) {
					configSet = append(configSet, setPrefix+"ipv6-extension-header hop-by-hop-header router-alert-option")
				}
				for _, v3 := range ipIPv6ExtHeaderHopByHopHeader["user_defined_option_type"].([]interface{}) {
					configSet = append(configSet,
						setPrefix+"ipv6-extension-header hop-by-hop-header user-defined-option-type "+v3.(string))
				}
			}
		}
		if ipIPv6ExtHeader["mobility_header"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header mobility-header")
		}
		if ipIPv6ExtHeader["no_next_header"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header no-next-header")
		}
		if ipIPv6ExtHeader["routing_header"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header routing-header")
		}
		if ipIPv6ExtHeader["shim6_header"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header shim6-header")
		}
		for _, v2 := range ipIPv6ExtHeader["user_defined_header_type"].([]interface{}) {
			configSet = append(configSet, setPrefix+"ipv6-extension-header user-defined-header-type "+v2.(string))
		}
	}
	if ip["ipv6_extension_header_limit"].(int) != -1 {
		configSet = append(configSet, setPrefix+"ipv6-extension-header-limit "+
			strconv.Itoa(ip["ipv6_extension_header_limit"].(int)))
	}
	if ip["ipv6_malformed_header"].(bool) {
		configSet = append(configSet, setPrefix+"ipv6-malformed-header")
	}
	if ip["loose_source_route_option"].(bool) {
		configSet = append(configSet, setPrefix+"loose-source-route-option")
	}
	if ip["record_route_option"].(bool) {
		configSet = append(configSet, setPrefix+"record-route-option")
	}
	if ip["security_option"].(bool) {
		configSet = append(configSet, setPrefix+"security-option")
	}
	if ip["source_route_option"].(bool) {
		configSet = append(configSet, setPrefix+"source-route-option")
	}
	if ip["spoofing"].(bool) {
		configSet = append(configSet, setPrefix+"spoofing")
	}
	if ip["stream_option"].(bool) {
		configSet = append(configSet, setPrefix+"stream-option")
	}
	if ip["strict_source_route_option"].(bool) {
		configSet = append(configSet, setPrefix+"strict-source-route-option")
	}
	if ip["tear_drop"].(bool) {
		configSet = append(configSet, setPrefix+"tear-drop")
	}
	if ip["timestamp_option"].(bool) {
		configSet = append(configSet, setPrefix+"timestamp-option")
	}
	for _, v := range ip["tunnel"].([]interface{}) {
		if v == nil {
			return configSet, errors.New("ip.0.tunnel block is empty")
		}
		ipTunnel := v.(map[string]interface{})
		if ipTunnel["bad_inner_header"].(bool) {
			configSet = append(configSet, setPrefix+"tunnel bad-inner-header")
		}
		for _, v2 := range ipTunnel["gre"].([]interface{}) {
			if v2 == nil {
				return configSet, errors.New("ip.0.tunnel.0.gre block is empty")
			}
			ipTunnelGre := v2.(map[string]interface{})
			if ipTunnelGre["gre_4in4"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel gre gre-4in4")
			}
			if ipTunnelGre["gre_4in6"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel gre gre-4in6")
			}
			if ipTunnelGre["gre_6in4"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel gre gre-6in4")
			}
			if ipTunnelGre["gre_6in6"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel gre gre-6in6")
			}
		}
		if ipTunnel["ip_in_udp_teredo"].(bool) {
			configSet = append(configSet, setPrefix+"tunnel ip-in-udp teredo")
		}
		for _, v2 := range ipTunnel["ipip"].([]interface{}) {
			if v2 == nil {
				return configSet, errors.New("ip.0.tunnel.0.ipip block is empty")
			}
			ipTunnelIPIP := v2.(map[string]interface{})
			if ipTunnelIPIP["ipip_4in4"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel ipip ipip-4in4")
			}
			if ipTunnelIPIP["ipip_4in6"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel ipip ipip-4in6")
			}
			if ipTunnelIPIP["ipip_6in4"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel ipip ipip-6in4")
			}
			if ipTunnelIPIP["ipip_6in6"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel ipip ipip-6in6")
			}
			if ipTunnelIPIP["ipip_6over4"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel ipip ipip-6over4")
			}
			if ipTunnelIPIP["ipip_6to4relay"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel ipip ipip-6to4relay")
			}
			if ipTunnelIPIP["dslite"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel ipip dslite")
			}
			if ipTunnelIPIP["isatap"].(bool) {
				configSet = append(configSet, setPrefix+"tunnel ipip isatap")
			}
		}
	}
	if ip["unknown_protocol"].(bool) {
		configSet = append(configSet, setPrefix+"unknown-protocol")
	}
	if len(configSet) == 0 {
		return configSet, errors.New("ip block is empty")
	}

	return configSet, nil
}

func setSecurityScreenTCP(tcp map[string]interface{}, setPrefix string) ([]string, error) {
	configSet := make([]string, 0)
	setPrefix += "tcp "
	if tcp["fin_no_ack"].(bool) {
		configSet = append(configSet, setPrefix+"fin-no-ack")
	}
	if tcp["land"].(bool) {
		configSet = append(configSet, setPrefix+"land")
	}
	if tcp["no_flag"].(bool) {
		configSet = append(configSet, setPrefix+"tcp-no-flag")
	}
	for _, v := range tcp["port_scan"].([]interface{}) {
		configSet = append(configSet, setPrefix+"port-scan")
		if v != nil {
			tcpPortScan := v.(map[string]interface{})
			if tcpPortScan["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"port-scan threshold "+
					strconv.Itoa(tcpPortScan["threshold"].(int)))
			}
		}
	}
	for _, v := range tcp["sweep"].([]interface{}) {
		configSet = append(configSet, setPrefix+"tcp-sweep")
		if v != nil {
			tcpSweep := v.(map[string]interface{})
			if tcpSweep["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"tcp-sweep threshold "+
					strconv.Itoa(tcpSweep["threshold"].(int)))
			}
		}
	}
	for _, v := range tcp["syn_ack_ack_proxy"].([]interface{}) {
		configSet = append(configSet, setPrefix+"syn-ack-ack-proxy")
		if v != nil {
			tcpSAAP := v.(map[string]interface{})
			if tcpSAAP["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"syn-ack-ack-proxy threshold "+
					strconv.Itoa(tcpSAAP["threshold"].(int)))
			}
		}
	}
	if tcp["syn_fin"].(bool) {
		configSet = append(configSet, setPrefix+"syn-fin")
	}
	for _, v := range tcp["syn_flood"].([]interface{}) {
		configSet = append(configSet, setPrefix+"syn-flood")
		if v != nil {
			tcpSynFlood := v.(map[string]interface{})
			if tcpSynFlood["alarm_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"syn-flood alarm-threshold "+
					strconv.Itoa(tcpSynFlood["alarm_threshold"].(int)))
			}
			if tcpSynFlood["attack_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"syn-flood attack-threshold "+
					strconv.Itoa(tcpSynFlood["attack_threshold"].(int)))
			}
			if tcpSynFlood["destination_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"syn-flood destination-threshold "+
					strconv.Itoa(tcpSynFlood["destination_threshold"].(int)))
			}
			if tcpSynFlood["source_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"syn-flood source-threshold "+
					strconv.Itoa(tcpSynFlood["source_threshold"].(int)))
			}
			if tcpSynFlood["timeout"].(int) != 0 {
				configSet = append(configSet, setPrefix+"syn-flood timeout "+
					strconv.Itoa(tcpSynFlood["timeout"].(int)))
			}
			whitelistNameList := make([]string, 0)
			for _, v2 := range tcpSynFlood["whitelist"].(*schema.Set).List() {
				whitelist := v2.(map[string]interface{})
				if len(whitelist["source_address"].(*schema.Set).List()) == 0 &&
					len(whitelist["destination_address"].(*schema.Set).List()) == 0 {
					return configSet, fmt.Errorf("white-list %s need to have a source or destination address set",
						whitelist["name"].(string))
				}
				if slices.Contains(whitelistNameList, whitelist["name"].(string)) {
					return configSet, fmt.Errorf("multiple blocks whitelist with the same name %s", whitelist["name"].(string))
				}
				whitelistNameList = append(whitelistNameList, whitelist["name"].(string))
				for _, destination := range sortSetOfString(whitelist["destination_address"].(*schema.Set).List()) {
					configSet = append(configSet, setPrefix+"syn-flood white-list "+whitelist["name"].(string)+
						" destination-address "+destination)
				}
				for _, source := range sortSetOfString(whitelist["source_address"].(*schema.Set).List()) {
					configSet = append(configSet, setPrefix+"syn-flood white-list "+whitelist["name"].(string)+
						" source-address "+source)
				}
			}
		}
	}
	if tcp["syn_frag"].(bool) {
		configSet = append(configSet, setPrefix+"syn-frag")
	}
	if tcp["winnuke"].(bool) {
		configSet = append(configSet, setPrefix+"winnuke")
	}

	return configSet, nil
}

func setSecurityScreenUDP(udp map[string]interface{}, setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "udp "
	for _, v := range udp["flood"].([]interface{}) {
		configSet = append(configSet, setPrefix+"flood")
		if v != nil {
			udpFlood := v.(map[string]interface{})
			if udpFlood["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"flood threshold "+
					strconv.Itoa(udpFlood["threshold"].(int)))
			}
			for _, whitelist := range sortSetOfString(udpFlood["whitelist"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"flood white-list "+whitelist)
			}
		}
	}
	for _, v := range udp["port_scan"].([]interface{}) {
		configSet = append(configSet, setPrefix+"port-scan")
		if v != nil {
			udpPortScan := v.(map[string]interface{})
			if udpPortScan["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"port-scan threshold "+
					strconv.Itoa(udpPortScan["threshold"].(int)))
			}
		}
	}
	for _, v := range udp["sweep"].([]interface{}) {
		configSet = append(configSet, setPrefix+"udp-sweep")
		if v != nil {
			udpSweep := v.(map[string]interface{})
			if udpSweep["threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"udp-sweep threshold "+
					strconv.Itoa(udpSweep["threshold"].(int)))
			}
		}
	}

	return configSet
}

func readSecurityScreen(name string, junSess *junos.Session,
) (confRead screenOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security screen ids-option \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "alarm-without-drop":
				confRead.alarmWithoutDrop = true
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "icmp "):
				if err := confRead.readSecurityScreenIcmp(itemTrim); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "ip"):
				if err := confRead.readSecurityScreenIP(itemTrim); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "limit-session "):
				if len(confRead.limitSession) == 0 {
					confRead.limitSession = append(confRead.limitSession, map[string]interface{}{
						"destination_ip_based": 0,
						"source_ip_based":      0,
					})
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "destination-ip-based "):
					confRead.limitSession[0]["destination_ip_based"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case balt.CutPrefixInString(&itemTrim, "source-ip-based "):
					confRead.limitSession[0]["source_ip_based"], err = strconv.Atoi(itemTrim)
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case balt.CutPrefixInString(&itemTrim, "tcp"):
				if err := confRead.readSecurityScreenTCP(itemTrim); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "udp"):
				if err := confRead.readSecurityScreenUDP(itemTrim); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}

func (confRead *screenOptions) readSecurityScreenIcmp(itemTrim string) (err error) {
	if len(confRead.icmp) == 0 {
		confRead.icmp = append(confRead.icmp, map[string]interface{}{
			"flood":            make([]map[string]interface{}, 0),
			"fragment":         false,
			"icmpv6_malformed": false,
			"sweep":            make([]map[string]interface{}, 0),
			"large":            false,
			"ping_death":       false,
		})
	}
	switch {
	case balt.CutPrefixInString(&itemTrim, "flood"):
		if len(confRead.icmp[0]["flood"].([]map[string]interface{})) == 0 {
			confRead.icmp[0]["flood"] = append(confRead.icmp[0]["flood"].([]map[string]interface{}), map[string]interface{}{
				"threshold": 0,
			})
		}
		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			confRead.icmp[0]["flood"].([]map[string]interface{})[0]["threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case itemTrim == "fragment":
		confRead.icmp[0]["fragment"] = true
	case itemTrim == "icmpv6-malformed":
		confRead.icmp[0]["icmpv6_malformed"] = true
	case itemTrim == "large":
		confRead.icmp[0]["large"] = true
	case itemTrim == "ping-death":
		confRead.icmp[0]["ping_death"] = true
	case balt.CutPrefixInString(&itemTrim, "ip-sweep"):
		if len(confRead.icmp[0]["sweep"].([]map[string]interface{})) == 0 {
			confRead.icmp[0]["sweep"] = append(confRead.icmp[0]["sweep"].([]map[string]interface{}), map[string]interface{}{
				"threshold": 0,
			})
		}
		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			confRead.icmp[0]["sweep"].([]map[string]interface{})[0]["threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	}

	return nil
}

func (confRead *screenOptions) readSecurityScreenIP(itemTrim string) (err error) {
	if len(confRead.ip) == 0 {
		confRead.ip = append(confRead.ip, map[string]interface{}{
			"bad_option":                  false,
			"block_frag":                  false,
			"ipv6_extension_header":       make([]map[string]interface{}, 0),
			"ipv6_extension_header_limit": -1,
			"ipv6_malformed_header":       false,
			"loose_source_route_option":   false,
			"record_route_option":         false,
			"security_option":             false,
			"source_route_option":         false,
			"spoofing":                    false,
			"stream_option":               false,
			"strict_source_route_option":  false,
			"tear_drop":                   false,
			"timestamp_option":            false,
			"tunnel":                      make([]map[string]interface{}, 0),
			"unknown_protocol":            false,
		})
	}
	switch {
	case itemTrim == " bad-option":
		confRead.ip[0]["bad_option"] = true
	case itemTrim == " block-frag":
		confRead.ip[0]["block_frag"] = true
	case balt.CutPrefixInString(&itemTrim, " ipv6-extension-header "):
		if len(confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})) == 0 {
			confRead.ip[0]["ipv6_extension_header"] = append(
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{}), map[string]interface{}{
					"ah_header":                false,
					"esp_header":               false,
					"hip_header":               false,
					"destination_header":       make([]map[string]interface{}, 0),
					"fragment_header":          false,
					"hop_by_hop_header":        make([]map[string]interface{}, 0),
					"mobility_header":          false,
					"no_next_header":           false,
					"routing_header":           false,
					"shim6_header":             false,
					"user_defined_header_type": make([]string, 0),
				})
		}
		ipIPv6ExtensionHeader := confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]
		switch {
		case itemTrim == "AH-header":
			ipIPv6ExtensionHeader["ah_header"] = true
		case itemTrim == "ESP-header":
			ipIPv6ExtensionHeader["esp_header"] = true
		case itemTrim == "HIP-header":
			ipIPv6ExtensionHeader["hip_header"] = true
		case balt.CutPrefixInString(&itemTrim, "destination-header"):
			if len(ipIPv6ExtensionHeader["destination_header"].([]map[string]interface{})) == 0 {
				ipIPv6ExtensionHeader["destination_header"] = append(
					ipIPv6ExtensionHeader["destination_header"].([]map[string]interface{}),
					map[string]interface{}{
						"ilnp_nonce_option":                 false,
						"home_address_option":               false,
						"line_identification_option":        false,
						"tunnel_encapsulation_limit_option": false,
						"user_defined_option_type":          make([]string, 0),
					})
			}
			ipIPv6ExtensionHeaderDstHeader := ipIPv6ExtensionHeader["destination_header"].([]map[string]interface{})[0]
			switch {
			case itemTrim == " ILNP-nonce-option":
				ipIPv6ExtensionHeaderDstHeader["ilnp_nonce_option"] = true
			case itemTrim == " home-address-option":
				ipIPv6ExtensionHeaderDstHeader["home_address_option"] = true
			case itemTrim == " line-identification-option":
				ipIPv6ExtensionHeaderDstHeader["line_identification_option"] = true
			case itemTrim == " tunnel-encapsulation-limit-option":
				ipIPv6ExtensionHeaderDstHeader["tunnel_encapsulation_limit_option"] = true
			case balt.CutPrefixInString(&itemTrim, " user-defined-option-type "):
				ipIPv6ExtensionHeaderDstHeader["user_defined_option_type"] = append(
					ipIPv6ExtensionHeaderDstHeader["user_defined_option_type"].([]string),
					itemTrim,
				)
			}
		case itemTrim == "fragment-header":
			ipIPv6ExtensionHeader["fragment_header"] = true
		case balt.CutPrefixInString(&itemTrim, "hop-by-hop-header"):
			if len(ipIPv6ExtensionHeader["hop_by_hop_header"].([]map[string]interface{})) == 0 {
				ipIPv6ExtensionHeader["hop_by_hop_header"] = append(
					ipIPv6ExtensionHeader["hop_by_hop_header"].([]map[string]interface{}),
					map[string]interface{}{
						"calipso_option":           false,
						"rpl_option":               false,
						"smf_dpd_option":           false,
						"jumbo_payload_option":     false,
						"quick_start_option":       false,
						"router_alert_option":      false,
						"user_defined_option_type": make([]string, 0),
					})
			}
			ipIPv6ExtensionHeaderHopByHopHeader := ipIPv6ExtensionHeader["hop_by_hop_header"].([]map[string]interface{})[0]
			switch {
			case itemTrim == " CALIPSO-option":
				ipIPv6ExtensionHeaderHopByHopHeader["calipso_option"] = true
			case itemTrim == " RPL-option":
				ipIPv6ExtensionHeaderHopByHopHeader["rpl_option"] = true
			case itemTrim == " SMF-DPD-option":
				ipIPv6ExtensionHeaderHopByHopHeader["smf_dpd_option"] = true
			case itemTrim == " jumbo-payload-option":
				ipIPv6ExtensionHeaderHopByHopHeader["jumbo_payload_option"] = true
			case itemTrim == " quick-start-option":
				ipIPv6ExtensionHeaderHopByHopHeader["quick_start_option"] = true
			case itemTrim == " router-alert-option":
				ipIPv6ExtensionHeaderHopByHopHeader["router_alert_option"] = true
			case balt.CutPrefixInString(&itemTrim, " user-defined-option-type "):
				ipIPv6ExtensionHeaderHopByHopHeader["user_defined_option_type"] = append(
					ipIPv6ExtensionHeaderHopByHopHeader["user_defined_option_type"].([]string),
					itemTrim,
				)
			}
		case itemTrim == "mobility-header":
			ipIPv6ExtensionHeader["mobility_header"] = true
		case itemTrim == "no-next-header":
			ipIPv6ExtensionHeader["no_next_header"] = true
		case itemTrim == "routing-header":
			ipIPv6ExtensionHeader["routing_header"] = true
		case itemTrim == "shim6-header":
			ipIPv6ExtensionHeader["shim6_header"] = true
		case balt.CutPrefixInString(&itemTrim, "user-defined-header-type "):
			ipIPv6ExtensionHeader["user_defined_header_type"] = append(
				ipIPv6ExtensionHeader["user_defined_header_type"].([]string),
				itemTrim,
			)
		}
	case balt.CutPrefixInString(&itemTrim, " ipv6-extension-header-limit "):
		confRead.ip[0]["ipv6_extension_header_limit"], err = strconv.Atoi(itemTrim)
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == " ipv6-malformed-header":
		confRead.ip[0]["ipv6_malformed_header"] = true
	case itemTrim == " loose-source-route-option":
		confRead.ip[0]["loose_source_route_option"] = true
	case itemTrim == " record-route-option":
		confRead.ip[0]["record_route_option"] = true
	case itemTrim == " security-option":
		confRead.ip[0]["security_option"] = true
	case itemTrim == " source-route-option":
		confRead.ip[0]["source_route_option"] = true
	case itemTrim == " spoofing":
		confRead.ip[0]["spoofing"] = true
	case itemTrim == " stream-option":
		confRead.ip[0]["stream_option"] = true
	case itemTrim == " strict-source-route-option":
		confRead.ip[0]["strict_source_route_option"] = true
	case itemTrim == " tear-drop":
		confRead.ip[0]["tear_drop"] = true
	case itemTrim == " timestamp-option":
		confRead.ip[0]["timestamp_option"] = true
	case balt.CutPrefixInString(&itemTrim, " tunnel "):
		if len(confRead.ip[0]["tunnel"].([]map[string]interface{})) == 0 {
			confRead.ip[0]["tunnel"] = append(
				confRead.ip[0]["tunnel"].([]map[string]interface{}), map[string]interface{}{
					"bad_inner_header": false,
					"gre":              make([]map[string]interface{}, 0),
					"ip_in_udp_teredo": false,
					"ipip":             make([]map[string]interface{}, 0),
				})
		}
		switch {
		case itemTrim == "bad-inner-header":
			confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["bad_inner_header"] = true
		case balt.CutPrefixInString(&itemTrim, "gre "):
			if len(confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})) == 0 {
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"] = append(
					confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{}),
					map[string]interface{}{
						"gre_4in4": false,
						"gre_4in6": false,
						"gre_6in4": false,
						"gre_6in6": false,
					})
			}
			switch {
			case itemTrim == "gre-4in4":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})[0]["gre_4in4"] = true
			case itemTrim == "gre-4in6":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})[0]["gre_4in6"] = true
			case itemTrim == "gre-6in4":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})[0]["gre_6in4"] = true
			case itemTrim == "gre-6in6":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})[0]["gre_6in6"] = true
			}
		case itemTrim == "ip-in-udp teredo":
			confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ip_in_udp_teredo"] = true
		case balt.CutPrefixInString(&itemTrim, "ipip "):
			if len(confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})) == 0 {
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"] = append(
					confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{}),
					map[string]interface{}{
						"ipip_4in4":      false,
						"ipip_4in6":      false,
						"ipip_6in4":      false,
						"ipip_6in6":      false,
						"ipip_6over4":    false,
						"ipip_6to4relay": false,
						"dslite":         false,
						"isatap":         false,
					})
			}
			ipTunnelIPIP := confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]
			switch {
			case itemTrim == "ipip-4in4":
				ipTunnelIPIP["ipip_4in4"] = true
			case itemTrim == "ipip-4in6":
				ipTunnelIPIP["ipip_4in6"] = true
			case itemTrim == "ipip-6in4":
				ipTunnelIPIP["ipip_6in4"] = true
			case itemTrim == "ipip-6in6":
				ipTunnelIPIP["ipip_6in6"] = true
			case itemTrim == "ipip-6over4":
				ipTunnelIPIP["ipip_6over4"] = true
			case itemTrim == "ipip-6to4relay":
				ipTunnelIPIP["ipip_6to4relay"] = true
			case itemTrim == "dslite":
				ipTunnelIPIP["dslite"] = true
			case itemTrim == "isatap":
				ipTunnelIPIP["isatap"] = true
			}
		}
	case itemTrim == " unknown-protocol":
		confRead.ip[0]["unknown_protocol"] = true
	}

	return nil
}

func (confRead *screenOptions) readSecurityScreenTCP(itemTrim string) (err error) {
	if len(confRead.tcp) == 0 {
		confRead.tcp = append(confRead.tcp, map[string]interface{}{
			"fin_no_ack":        false,
			"land":              false,
			"no_flag":           false,
			"port_scan":         make([]map[string]interface{}, 0),
			"syn_ack_ack_proxy": make([]map[string]interface{}, 0),
			"syn_fin":           false,
			"syn_flood":         make([]map[string]interface{}, 0),
			"syn_frag":          false,
			"sweep":             make([]map[string]interface{}, 0),
			"winnuke":           false,
		})
	}
	switch {
	case itemTrim == " fin-no-ack":
		confRead.tcp[0]["fin_no_ack"] = true
	case itemTrim == " land":
		confRead.tcp[0]["land"] = true
	case itemTrim == " tcp-no-flag":
		confRead.tcp[0]["no_flag"] = true
	case balt.CutPrefixInString(&itemTrim, " port-scan"):
		if len(confRead.tcp[0]["port_scan"].([]map[string]interface{})) == 0 {
			confRead.tcp[0]["port_scan"] = append(confRead.tcp[0]["port_scan"].([]map[string]interface{}),
				map[string]interface{}{
					"threshold": 0,
				})
		}
		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			confRead.tcp[0]["port_scan"].([]map[string]interface{})[0]["threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case balt.CutPrefixInString(&itemTrim, " tcp-sweep"):
		if len(confRead.tcp[0]["sweep"].([]map[string]interface{})) == 0 {
			confRead.tcp[0]["sweep"] = append(confRead.tcp[0]["sweep"].([]map[string]interface{}),
				map[string]interface{}{
					"threshold": 0,
				})
		}
		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			confRead.tcp[0]["sweep"].([]map[string]interface{})[0]["threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case balt.CutPrefixInString(&itemTrim, " syn-ack-ack-proxy"):
		if len(confRead.tcp[0]["syn_ack_ack_proxy"].([]map[string]interface{})) == 0 {
			confRead.tcp[0]["syn_ack_ack_proxy"] = append(confRead.tcp[0]["syn_ack_ack_proxy"].([]map[string]interface{}),
				map[string]interface{}{
					"threshold": 0,
				})
		}
		synAckAckProxy := confRead.tcp[0]["syn_ack_ack_proxy"].([]map[string]interface{})[0]
		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			synAckAckProxy["threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case itemTrim == " syn-fin":
		confRead.tcp[0]["syn_fin"] = true
	case balt.CutPrefixInString(&itemTrim, " syn-flood"):
		if len(confRead.tcp[0]["syn_flood"].([]map[string]interface{})) == 0 {
			confRead.tcp[0]["syn_flood"] = append(
				confRead.tcp[0]["syn_flood"].([]map[string]interface{}), map[string]interface{}{
					"alarm_threshold":       0,
					"attack_threshold":      0,
					"destination_threshold": 0,
					"source_threshold":      0,
					"timeout":               0,
					"whitelist":             make([]map[string]interface{}, 0),
				})
		}
		synFlood := confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]
		switch {
		case balt.CutPrefixInString(&itemTrim, " alarm-threshold "):
			synFlood["alarm_threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " attack-threshold "):
			synFlood["attack_threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " destination-threshold "):
			synFlood["destination_threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " source-threshold "):
			synFlood["source_threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " timeout "):
			synFlood["timeout"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " white-list "):
			itemTrimFields := strings.Split(itemTrim, " ")
			wList := map[string]interface{}{
				"name":                itemTrimFields[0],
				"destination_address": make([]string, 0),
				"source_address":      make([]string, 0),
			}
			synFlood["whitelist"] = copyAndRemoveItemMapList(
				"name", wList, synFlood["whitelist"].([]map[string]interface{}))
			balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
			switch {
			case balt.CutPrefixInString(&itemTrim, "destination-address "):
				wList["destination_address"] = append(
					wList["destination_address"].([]string),
					itemTrim,
				)
			case balt.CutPrefixInString(&itemTrim, "source-address "):
				wList["source_address"] = append(
					wList["source_address"].([]string),
					itemTrim,
				)
			}
			synFlood["whitelist"] = append(synFlood["whitelist"].([]map[string]interface{}), wList)
		}
	case itemTrim == " syn-frag":
		confRead.tcp[0]["syn_frag"] = true
	case itemTrim == " winnuke":
		confRead.tcp[0]["winnuke"] = true
	}

	return nil
}

func (confRead *screenOptions) readSecurityScreenUDP(itemTrim string) (err error) {
	if len(confRead.udp) == 0 {
		confRead.udp = append(confRead.udp, map[string]interface{}{
			"flood":     make([]map[string]interface{}, 0),
			"port_scan": make([]map[string]interface{}, 0),
			"sweep":     make([]map[string]interface{}, 0),
		})
	}
	switch {
	case balt.CutPrefixInString(&itemTrim, " flood"):
		if len(confRead.udp[0]["flood"].([]map[string]interface{})) == 0 {
			confRead.udp[0]["flood"] = append(confRead.udp[0]["flood"].([]map[string]interface{}), map[string]interface{}{
				"threshold": 0,
				"whitelist": make([]string, 0),
			})
		}
		flood := confRead.udp[0]["flood"].([]map[string]interface{})[0]
		switch {
		case balt.CutPrefixInString(&itemTrim, " threshold "):
			flood["threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, " white-list "):
			flood["whitelist"] = append(flood["whitelist"].([]string), itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, " port-scan"):
		if len(confRead.udp[0]["port_scan"].([]map[string]interface{})) == 0 {
			confRead.udp[0]["port_scan"] = append(
				confRead.udp[0]["port_scan"].([]map[string]interface{}), map[string]interface{}{
					"threshold": 0,
				})
		}
		portScan := confRead.udp[0]["port_scan"].([]map[string]interface{})[0]
		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			portScan["threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case balt.CutPrefixInString(&itemTrim, " udp-sweep"):
		if len(confRead.udp[0]["sweep"].([]map[string]interface{})) == 0 {
			confRead.udp[0]["sweep"] = append(
				confRead.udp[0]["sweep"].([]map[string]interface{}), map[string]interface{}{
					"threshold": 0,
				})
		}
		sweep := confRead.udp[0]["sweep"].([]map[string]interface{})[0]
		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			sweep["threshold"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	}

	return nil
}

func delSecurityScreen(name string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security screen ids-option \""+name+"\"")

	return junSess.ConfigSet(configSet)
}

func fillSecurityScreenData(d *schema.ResourceData, screenOptions screenOptions) {
	if tfErr := d.Set("name", screenOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("alarm_without_drop", screenOptions.alarmWithoutDrop); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", screenOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("icmp", screenOptions.icmp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ip", screenOptions.ip); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("limit_session", screenOptions.limitSession); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tcp", screenOptions.tcp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("udp", screenOptions.udp); tfErr != nil {
		panic(tfErr)
	}
}
