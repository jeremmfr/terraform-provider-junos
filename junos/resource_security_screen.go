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
		CreateContext: resourceSecurityScreenCreate,
		ReadContext:   resourceSecurityScreenRead,
		UpdateContext: resourceSecurityScreenUpdate,
		DeleteContext: resourceSecurityScreenDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityScreenImport,
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
													Elem:     &schema.Schema{Type: schema.TypeString},
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
													Elem:     &schema.Schema{Type: schema.TypeString},
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
										Elem:     &schema.Schema{Type: schema.TypeString},
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
													ValidateDiagFunc: validateNameObjectJunos([]string{}, 32),
												},
												"destination_address": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"source_address": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
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
										Type:     schema.TypeList,
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security screen not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)
	securityScreenExists, err := checkSecurityScreenExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if securityScreenExists {
		sess.configClear(jnprSess)

		return diag.FromErr(fmt.Errorf("security screen %v already exists", d.Get("name").(string)))
	}

	if err := setSecurityScreen(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security_screen", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityScreenExists, err = checkSecurityScreenExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityScreenExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security screen %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityScreenReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityScreenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityScreenReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityScreenReadWJnprSess(
	d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	screenOptions, err := readSecurityScreen(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
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
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := delSecurityScreen(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}

	if err := setSecurityScreen(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security_screen", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityScreenReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityScreenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurityScreen(d.Get("name").(string), m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("delete resource junos_security_screen", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}
func resourceSecurityScreenImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityScreenExists, err := checkSecurityScreenExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !securityScreenExists {
		return nil, fmt.Errorf("don't find screen with id '%v' (id must be <name>)", d.Id())
	}
	screenOptions, err := readSecurityScreen(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurityScreenData(d, screenOptions)

	result[0] = d

	return result, nil
}

func checkSecurityScreenExists(name string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	screenConfig, err := sess.command("show configuration"+
		" security screen ids-option \""+name+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if screenConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setSecurityScreen(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
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
			return fmt.Errorf("icmp block is empty")
		}
		icmp := v.(map[string]interface{})
		configSet = append(configSet, setSecurityScreenIcmp(icmp, setPrefix)...)
	}
	for _, v := range d.Get("ip").([]interface{}) {
		if v == nil {
			return fmt.Errorf("ip block is empty")
		}
		ip := v.(map[string]interface{})
		if err := checkSetSecurityScreenIP(ip); err != nil {
			return err
		}
		ipSet, err := setSecurityScreenIP(ip, setPrefix)
		if err != nil {
			return err
		}
		configSet = append(configSet, ipSet...)
	}
	for _, v := range d.Get("limit_session").([]interface{}) {
		if v == nil {
			return fmt.Errorf("limit_session block is empty")
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
			return fmt.Errorf("tcp block is empty")
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
			return fmt.Errorf("udp block is empty")
		}
		udp := v.(map[string]interface{})
		configSet = append(configSet, setSecurityScreenUDP(udp, setPrefix)...)
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
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
func checkSetSecurityScreenIP(ip map[string]interface{}) error {
	if !ip["bad_option"].(bool) &&
		!ip["block_frag"].(bool) &&
		len(ip["ipv6_extension_header"].([]interface{})) == 0 &&
		ip["ipv6_extension_header_limit"].(int) == -1 &&
		!ip["ipv6_malformed_header"].(bool) &&
		!ip["loose_source_route_option"].(bool) &&
		!ip["record_route_option"].(bool) &&
		!ip["security_option"].(bool) &&
		!ip["source_route_option"].(bool) &&
		!ip["spoofing"].(bool) &&
		!ip["stream_option"].(bool) &&
		!ip["strict_source_route_option"].(bool) &&
		!ip["tear_drop"].(bool) &&
		!ip["timestamp_option"].(bool) &&
		len(ip["tunnel"].([]interface{})) == 0 &&
		!ip["unknown_protocol"].(bool) {
		return fmt.Errorf("ip block is empty")
	}

	return nil
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
			return configSet, fmt.Errorf("ip.0.ipv6_extension_header block is empty")
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
					if r := regexp.MustCompile(userDefinedOptionTypeRegex); !r.Match([]byte(v3.(string))) {
						return configSet, fmt.Errorf(
							"user_defined_option_type %v doesn't match '(1..255)' or '(1..255) to (1..255)'", v3.(string))
					}
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
					if r := regexp.MustCompile(userDefinedOptionTypeRegex); !r.Match([]byte(v3.(string))) {
						return configSet, fmt.Errorf(
							"user_defined_option_type %v doesn't match '(1..255)' or '(1..255) to (1..255)'", v3.(string))
					}
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
			if r := regexp.MustCompile(userDefinedHeaderTypeRegex); !r.Match([]byte(v2.(string))) {
				return configSet, fmt.Errorf(
					"user_defined_header_type %v doesn't match '(0..255)' or '(0..255) to (0..255)'", v2.(string))
			}
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
			return configSet, fmt.Errorf("ip.0.tunnel block is empty")
		}
		ipTunnel := v.(map[string]interface{})
		if ipTunnel["bad_inner_header"].(bool) {
			configSet = append(configSet, setPrefix+"tunnel bad-inner-header")
		}
		for _, v2 := range ipTunnel["gre"].([]interface{}) {
			if v2 == nil {
				return configSet, fmt.Errorf("ip.0.tunnel.0.gre block is empty")
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
				return configSet, fmt.Errorf("ip.0.tunnel.0.ipip block is empty")
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
			for _, v2 := range tcpSynFlood["whitelist"].(*schema.Set).List() {
				whitelist := v2.(map[string]interface{})
				if len(whitelist["source_address"].([]interface{})) == 0 &&
					len(whitelist["destination_address"].([]interface{})) == 0 {
					return configSet, fmt.Errorf("white-list %s need to have a source or destination address set",
						whitelist["name"].(string))
				}
				for _, destination := range whitelist["destination_address"].([]interface{}) {
					if err := validateCIDRNetwork(destination.(string)); err != nil {
						return configSet, err
					}
					configSet = append(configSet, setPrefix+"syn-flood white-list "+whitelist["name"].(string)+
						" destination-address "+destination.(string))
				}
				for _, source := range whitelist["source_address"].([]interface{}) {
					if err := validateCIDRNetwork(source.(string)); err != nil {
						return configSet, err
					}
					configSet = append(configSet, setPrefix+"syn-flood white-list "+whitelist["name"].(string)+
						" source-address "+source.(string))
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
			for _, whitelist := range udpFlood["whitelist"].([]interface{}) {
				configSet = append(configSet, setPrefix+"flood white-list "+whitelist.(string))
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
func readSecurityScreen(name string, m interface{}, jnprSess *NetconfObject) (screenOptions, error) {
	sess := m.(*Session)
	var confRead screenOptions

	screenConfig, err := sess.command("show configuration security screen ids-option "+
		"\""+name+"\" | display set relative ", jnprSess)
	if err != nil {
		return confRead, err
	}
	if screenConfig != emptyWord {
		confRead.name = name
		for _, item := range strings.Split(screenConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case itemTrim == "alarm-without-drop":
				confRead.alarmWithoutDrop = true
			case strings.HasPrefix(itemTrim, "description "):
				confRead.description = strings.Trim(strings.TrimPrefix(itemTrim, "description "), "\"")
			case strings.HasPrefix(itemTrim, "icmp "):
				if err := readSecurityScreenIcmp(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "ip"):
				if err := readSecurityScreenIP(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "limit-session "):
				if len(confRead.limitSession) == 0 {
					confRead.limitSession = append(confRead.limitSession, map[string]interface{}{
						"destination_ip_based": 0,
						"source_ip_based":      0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "limit-session destination-ip-based "):
					var err error
					confRead.limitSession[0]["destination_ip_based"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "limit-session destination-ip-based "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "limit-session source-ip-based "):
					var err error
					confRead.limitSession[0]["source_ip_based"], err = strconv.Atoi(
						strings.TrimPrefix(itemTrim, "limit-session source-ip-based "))
					if err != nil {
						return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
					}
				}
			case strings.HasPrefix(itemTrim, "tcp"):
				if err := readSecurityScreenTCP(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "udp"):
				if err := readSecurityScreenUDP(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			}
		}
	}

	return confRead, nil
}
func readSecurityScreenIcmp(confRead *screenOptions, itemTrim string) error {
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
	case strings.HasPrefix(itemTrim, "icmp flood"):
		if len(confRead.icmp[0]["flood"].([]map[string]interface{})) == 0 {
			confRead.icmp[0]["flood"] = append(confRead.icmp[0]["flood"].([]map[string]interface{}), map[string]interface{}{
				"threshold": 0,
			})
		}
		if strings.HasPrefix(itemTrim, "icmp flood threshold ") {
			var err error
			confRead.icmp[0]["flood"].([]map[string]interface{})[0]["threshold"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "icmp flood threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case itemTrim == "icmp fragment":
		confRead.icmp[0]["fragment"] = true
	case itemTrim == "icmp icmpv6-malformed":
		confRead.icmp[0]["icmpv6_malformed"] = true
	case itemTrim == "icmp large":
		confRead.icmp[0]["large"] = true
	case itemTrim == "icmp ping-death":
		confRead.icmp[0]["ping_death"] = true
	case strings.HasPrefix(itemTrim, "icmp ip-sweep"):
		if len(confRead.icmp[0]["sweep"].([]map[string]interface{})) == 0 {
			confRead.icmp[0]["sweep"] = append(confRead.icmp[0]["sweep"].([]map[string]interface{}), map[string]interface{}{
				"threshold": 0,
			})
		}
		if strings.HasPrefix(itemTrim, "icmp ip-sweep threshold ") {
			var err error
			confRead.icmp[0]["sweep"].([]map[string]interface{})[0]["threshold"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "icmp ip-sweep threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	}

	return nil
}
func readSecurityScreenIP(confRead *screenOptions, itemTrim string) error {
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
	case itemTrim == "ip bad-option":
		confRead.ip[0]["bad_option"] = true
	case itemTrim == "ip block-frag":
		confRead.ip[0]["block_frag"] = true
	case strings.HasPrefix(itemTrim, "ip ipv6-extension-header "):
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
		switch {
		case strings.HasPrefix(itemTrim, "ip ipv6-extension-header AH-header"):
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["ah_header"] = true
		case strings.HasPrefix(itemTrim, "ip ipv6-extension-header ESP-header"):
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["esp_header"] = true
		case strings.HasPrefix(itemTrim, "ip ipv6-extension-header HIP-header"):
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hip_header"] = true
		case strings.HasPrefix(itemTrim, "ip ipv6-extension-header destination-header"):
			if len(confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"].([]map[string]interface{})) == 0 { // nolint: lll
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"] = append(
					confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"].([]map[string]interface{}), // nolint: lll
					map[string]interface{}{
						"ilnp_nonce_option":                 false,
						"home_address_option":               false,
						"line_identification_option":        false,
						"tunnel_encapsulation_limit_option": false,
						"user_defined_option_type":          make([]string, 0),
					})
			}
			switch {
			case itemTrim == "ip ipv6-extension-header destination-header ILNP-nonce-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"].([]map[string]interface{})[0]["ilnp_nonce_option"] = // nolint: lll
					true
			case itemTrim == "ip ipv6-extension-header destination-header home-address-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"].([]map[string]interface{})[0]["home_address_option"] = // nolint: lll
					true
			case itemTrim == "ip ipv6-extension-header destination-header line-identification-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"].([]map[string]interface{})[0]["line_identification_option"] = // nolint: lll
					true
			case itemTrim == "ip ipv6-extension-header destination-header tunnel-encapsulation-limit-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"].([]map[string]interface{})[0]["tunnel_encapsulation_limit_option"] = // nolint: lll
					true
			case strings.HasPrefix(itemTrim, "ip ipv6-extension-header destination-header user-defined-option-type "):
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"].([]map[string]interface{})[0]["user_defined_option_type"] = // nolint: lll
					append(
						confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["destination_header"].([]map[string]interface{})[0]["user_defined_option_type"].([]string), // nolint: lll
						strings.TrimPrefix(itemTrim, "ip ipv6-extension-header destination-header user-defined-option-type "))
			}
		case strings.HasPrefix(itemTrim, "ip ipv6-extension-header fragment-header"):
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["fragment_header"] = true
		case strings.HasPrefix(itemTrim, "ip ipv6-extension-header hop-by-hop-header"):
			if len(confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})) == 0 { // nolint: lll
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"] = append(
					confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{}), // nolint: lll
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
			switch {
			case itemTrim == "ip ipv6-extension-header hop-by-hop-header CALIPSO-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})[0]["calipso_option"] = // nolint: lll
					true
			case itemTrim == "ip ipv6-extension-header hop-by-hop-header RPL-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})[0]["rpl_option"] = // nolint: lll
					true
			case itemTrim == "ip ipv6-extension-header hop-by-hop-header SMF-DPD-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})[0]["smf_dpd_option"] = // nolint: lll
					true
			case itemTrim == "ip ipv6-extension-header hop-by-hop-header jumbo-payload-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})[0]["jumbo_payload_option"] = // nolint: lll
					true
			case itemTrim == "ip ipv6-extension-header hop-by-hop-header quick-start-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})[0]["quick_start_option"] = // nolint: lll
					true
			case itemTrim == "ip ipv6-extension-header hop-by-hop-header router-alert-option":
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})[0]["router_alert_option"] = // nolint: lll
					true
			case strings.HasPrefix(itemTrim, "ip ipv6-extension-header hop-by-hop-header user-defined-option-type "):
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})[0]["user_defined_option_type"] = // nolint: lll
					append(
						confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["hop_by_hop_header"].([]map[string]interface{})[0]["user_defined_option_type"].([]string), // nolint: lll
						strings.TrimPrefix(itemTrim, "ip ipv6-extension-header hop-by-hop-header user-defined-option-type "))
			}
		case itemTrim == "ip ipv6-extension-header mobility-header":
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["mobility_header"] = true
		case itemTrim == "ip ipv6-extension-header no-next-header":
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["no_next_header"] = true
		case itemTrim == "ip ipv6-extension-header routing-header":
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["routing_header"] = true
		case itemTrim == "ip ipv6-extension-header shim6-header":
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["shim6_header"] = true
		case strings.HasPrefix(itemTrim, "ip ipv6-extension-header user-defined-header-type "):
			confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["user_defined_header_type"] = append(
				confRead.ip[0]["ipv6_extension_header"].([]map[string]interface{})[0]["user_defined_header_type"].([]string),
				strings.TrimPrefix(itemTrim, "ip ipv6-extension-header user-defined-header-type "))
		}
	case strings.HasPrefix(itemTrim, "ip ipv6-extension-header-limit "):
		var err error
		confRead.ip[0]["ipv6_extension_header_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim,
			"ip ipv6-extension-header-limit "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case itemTrim == "ip ipv6-malformed-header":
		confRead.ip[0]["ipv6_malformed_header"] = true
	case itemTrim == "ip loose-source-route-option":
		confRead.ip[0]["loose_source_route_option"] = true
	case itemTrim == "ip record-route-option":
		confRead.ip[0]["record_route_option"] = true
	case itemTrim == "ip security-option":
		confRead.ip[0]["security_option"] = true
	case itemTrim == "ip source-route-option":
		confRead.ip[0]["source_route_option"] = true
	case itemTrim == "ip spoofing":
		confRead.ip[0]["spoofing"] = true
	case itemTrim == "ip stream-option":
		confRead.ip[0]["stream_option"] = true
	case itemTrim == "ip strict-source-route-option":
		confRead.ip[0]["strict_source_route_option"] = true
	case itemTrim == "ip tear-drop":
		confRead.ip[0]["tear_drop"] = true
	case itemTrim == "ip timestamp-option":
		confRead.ip[0]["timestamp_option"] = true
	case strings.HasPrefix(itemTrim, "ip tunnel "):
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
		case itemTrim == "ip tunnel bad-inner-header":
			confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["bad_inner_header"] = true
		case strings.HasPrefix(itemTrim, "ip tunnel gre "):
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
			case itemTrim == "ip tunnel gre gre-4in4":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})[0]["gre_4in4"] = true
			case itemTrim == "ip tunnel gre gre-4in6":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})[0]["gre_4in6"] = true
			case itemTrim == "ip tunnel gre gre-6in4":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})[0]["gre_6in4"] = true
			case itemTrim == "ip tunnel gre gre-6in6":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["gre"].([]map[string]interface{})[0]["gre_6in6"] = true
			}
		case itemTrim == "ip tunnel ip-in-udp teredo":
			confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ip_in_udp_teredo"] = true
		case strings.HasPrefix(itemTrim, "ip tunnel ipip "):
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
			switch {
			case itemTrim == "ip tunnel ipip ipip-4in4":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]["ipip_4in4"] = true
			case itemTrim == "ip tunnel ipip ipip-4in6":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]["ipip_4in6"] = true
			case itemTrim == "ip tunnel ipip ipip-6in4":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]["ipip_6in4"] = true
			case itemTrim == "ip tunnel ipip ipip-6in6":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]["ipip_6in6"] = true
			case itemTrim == "ip tunnel ipip ipip-6over4":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]["ipip_6over4"] = true
			case itemTrim == "ip tunnel ipip ipip-6to4relay":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]["ipip_6to4relay"] =
					true
			case itemTrim == "ip tunnel ipip dslite":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]["dslite"] = true
			case itemTrim == "ip tunnel ipip isatap":
				confRead.ip[0]["tunnel"].([]map[string]interface{})[0]["ipip"].([]map[string]interface{})[0]["isatap"] = true
			}
		}
	case itemTrim == "ip unknown-protocol":
		confRead.ip[0]["unknown_protocol"] = true
	}

	return nil
}
func readSecurityScreenTCP(confRead *screenOptions, itemTrim string) error {
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
	case itemTrim == "tcp fin-no-ack":
		confRead.tcp[0]["fin_no_ack"] = true
	case itemTrim == "tcp land":
		confRead.tcp[0]["land"] = true
	case itemTrim == "tcp tcp-no-flag":
		confRead.tcp[0]["no_flag"] = true
	case strings.HasPrefix(itemTrim, "tcp port-scan"):
		if len(confRead.tcp[0]["port_scan"].([]map[string]interface{})) == 0 {
			confRead.tcp[0]["port_scan"] = append(confRead.tcp[0]["port_scan"].([]map[string]interface{}),
				map[string]interface{}{
					"threshold": 0,
				})
		}
		if strings.HasPrefix(itemTrim, "tcp port-scan threshold ") {
			var err error
			confRead.tcp[0]["port_scan"].([]map[string]interface{})[0]["threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp port-scan threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "tcp tcp-sweep"):
		if len(confRead.tcp[0]["sweep"].([]map[string]interface{})) == 0 {
			confRead.tcp[0]["sweep"] = append(confRead.tcp[0]["sweep"].([]map[string]interface{}),
				map[string]interface{}{
					"threshold": 0,
				})
		}
		if strings.HasPrefix(itemTrim, "tcp tcp-sweep threshold ") {
			var err error
			confRead.tcp[0]["sweep"].([]map[string]interface{})[0]["threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp tcp-sweep threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "tcp syn-ack-ack-proxy"):
		if len(confRead.tcp[0]["syn_ack_ack_proxy"].([]map[string]interface{})) == 0 {
			confRead.tcp[0]["syn_ack_ack_proxy"] = append(confRead.tcp[0]["syn_ack_ack_proxy"].([]map[string]interface{}),
				map[string]interface{}{
					"threshold": 0,
				})
		}
		if strings.HasPrefix(itemTrim, "tcp syn-ack-ack-proxy threshold ") {
			var err error
			confRead.tcp[0]["syn_ack_ack_proxy"].([]map[string]interface{})[0]["threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp syn-ack-ack-proxy threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case itemTrim == "tcp syn-fin":
		confRead.tcp[0]["syn_fin"] = true
	case strings.HasPrefix(itemTrim, "tcp syn-flood"):
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
		switch {
		case strings.HasPrefix(itemTrim, "tcp syn-flood alarm-threshold "):
			var err error
			confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["alarm_threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp syn-flood alarm-threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "tcp syn-flood attack-threshold "):
			var err error
			confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["attack_threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp syn-flood attack-threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "tcp syn-flood destination-threshold "):
			var err error
			confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["destination_threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp syn-flood destination-threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "tcp syn-flood source-threshold "):
			var err error
			confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["source_threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp syn-flood source-threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "tcp syn-flood timeout "):
			var err error
			confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["timeout"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp syn-flood timeout "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "tcp syn-flood white-list "):
			whiteListLineCut := strings.Split(strings.TrimPrefix(itemTrim, "tcp syn-flood white-list "), " ")
			m := map[string]interface{}{
				"name":                whiteListLineCut[0],
				"destination_address": make([]string, 0),
				"source_address":      make([]string, 0),
			}
			m, confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["whitelist"] = copyAndRemoveItemMapList(
				"name", false, m,
				confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["whitelist"].([]map[string]interface{}))
			itemTrimWhiteList := strings.TrimPrefix(itemTrim, "tcp syn-flood white-list "+whiteListLineCut[0]+" ")
			switch {
			case strings.HasPrefix(itemTrimWhiteList, "destination-address "):
				m["destination_address"] = append(m["destination_address"].([]string),
					strings.TrimPrefix(itemTrimWhiteList, "destination-address "))
			case strings.HasPrefix(itemTrimWhiteList, "source-address "):
				m["source_address"] = append(m["source_address"].([]string),
					strings.TrimPrefix(itemTrimWhiteList, "source-address "))
			}
			confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["whitelist"] = append(
				confRead.tcp[0]["syn_flood"].([]map[string]interface{})[0]["whitelist"].([]map[string]interface{}), m)
		}
	case itemTrim == "tcp syn-frag":
		confRead.tcp[0]["syn_frag"] = true
	case itemTrim == "tcp winnuke":
		confRead.tcp[0]["winnuke"] = true
	}

	return nil
}
func readSecurityScreenUDP(confRead *screenOptions, itemTrim string) error {
	if len(confRead.udp) == 0 {
		confRead.udp = append(confRead.udp, map[string]interface{}{
			"flood":     make([]map[string]interface{}, 0),
			"port_scan": make([]map[string]interface{}, 0),
			"sweep":     make([]map[string]interface{}, 0),
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "udp flood"):
		if len(confRead.udp[0]["flood"].([]map[string]interface{})) == 0 {
			confRead.udp[0]["flood"] = append(confRead.udp[0]["flood"].([]map[string]interface{}), map[string]interface{}{
				"threshold": 0,
				"whitelist": make([]string, 0),
			})
		}
		switch {
		case strings.HasPrefix(itemTrim, "udp flood threshold "):
			var err error
			confRead.udp[0]["flood"].([]map[string]interface{})[0]["threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "udp flood threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "udp flood white-list "):
			confRead.udp[0]["flood"].([]map[string]interface{})[0]["whitelist"] = append(
				confRead.udp[0]["flood"].([]map[string]interface{})[0]["whitelist"].([]string),
				strings.TrimPrefix(itemTrim, "udp flood white-list "))
		}
	case strings.HasPrefix(itemTrim, "udp port-scan"):
		if len(confRead.udp[0]["port_scan"].([]map[string]interface{})) == 0 {
			confRead.udp[0]["port_scan"] = append(
				confRead.udp[0]["port_scan"].([]map[string]interface{}), map[string]interface{}{
					"threshold": 0,
				})
		}
		if strings.HasPrefix(itemTrim, "udp port-scan threshold ") {
			var err error
			confRead.udp[0]["port_scan"].([]map[string]interface{})[0]["threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "udp port-scan threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "udp udp-sweep"):
		if len(confRead.udp[0]["sweep"].([]map[string]interface{})) == 0 {
			confRead.udp[0]["sweep"] = append(
				confRead.udp[0]["sweep"].([]map[string]interface{}), map[string]interface{}{
					"threshold": 0,
				})
		}
		if strings.HasPrefix(itemTrim, "udp udp-sweep threshold ") {
			var err error
			confRead.udp[0]["sweep"].([]map[string]interface{})[0]["threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "udp udp-sweep threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	}

	return nil
}

func delSecurityScreen(name string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security screen ids-option \""+name+"\"")
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
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
