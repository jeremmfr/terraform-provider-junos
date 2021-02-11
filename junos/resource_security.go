package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type securityOptions struct {
	alg             []map[string]interface{}
	flow            []map[string]interface{}
	forwardingOpts  []map[string]interface{}
	log             []map[string]interface{}
	ikeTraceoptions []map[string]interface{}
	utm             []map[string]interface{}
}

func resourceSecurity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityCreate,
		ReadContext:   resourceSecurityRead,
		UpdateContext: resourceSecurityUpdate,
		DeleteContext: resourceSecurityDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityImport,
		},
		Schema: map[string]*schema.Schema{
			"alg": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ftp_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"h323_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"mgcp_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"msrpc_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"pptp_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rtsp_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sccp_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sip_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sunrpc_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"talk_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tftp_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"flow": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"drop_matching_reserved_ip_address": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"drop_matching_link_local_address": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"reverse_route_packet_mode_vr": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"aging": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"early_ageout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"high_watermark": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 100),
									},
									"low_watermark": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 100),
									},
								},
							},
						},
						"allow_dns_reply": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_embedded_icmp": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_reverse_ecmp": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_reroute_uniform_link_check_nat": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ethernet_switching": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_non_ip_all": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"flow.0.ethernet_switching.0.bypass_non_ip_unicast"},
									},
									"bypass_non_ip_unicast": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"flow.0.ethernet_switching.0.block_non_ip_all"},
									},
									"bpdu_vlan_flooding": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_packet_flooding": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"no_trace_route": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"force_ip_reassembly": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ipsec_performance_acceleration": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"mcast_buffer_enhance": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"pending_sess_queue_length": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"high", "moderate", "normal"}, false),
						},
						"preserve_incoming_fragment_size": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"route_change_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(6, 1800),
						},
						"syn_flood_protection_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"syn-cookie", "syn-proxy"}, false),
						},
						"sync_icmp_session": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tcp_mss": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"all_tcp_mss": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(64, 65535),
									},
									"gre_in": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mss": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(64, 65535),
												},
											},
										},
									},
									"gre_out": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mss": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(64, 65535),
												},
											},
										},
									},
									"ipsec_vpn": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mss": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(64, 65535),
												},
											},
										},
									},
								},
							},
						},
						"tcp_session": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fin_invalidate_session": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"maximum_window": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"64K", "128K", "256K", "512K", "1M"}, false),
									},
									"no_sequence_check": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_syn_check": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"flow.0.tcp_session.0.strict_syn_check"},
									},
									"no_syn_check_in_tunnel": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"flow.0.tcp_session.0.strict_syn_check"},
									},
									"rst_invalidate_session": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"rst_sequence_check": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"strict_syn_check": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"flow.0.tcp_session.0.no_syn_check", "flow.0.tcp_session.0.no_syn_check_in_tunnel"},
									},
									"tcp_initial_timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(4, 300),
									},
									"time_wait_state": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"apply_to_half_close_state": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"session_ageout": {
													Type:          schema.TypeBool,
													Optional:      true,
													ConflictsWith: []string{"flow.0.tcp_session.0.time_wait_state.0.session_timeout"},
												},
												"session_timeout": {
													Type:          schema.TypeInt,
													Optional:      true,
													ValidateFunc:  validation.IntBetween(2, 600),
													ConflictsWith: []string{"flow.0.tcp_session.0.time_wait_state.0.session_ageout"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"forwarding_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inet6_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"drop", "flow-based", "packet-based"}, false),
						},
						"iso_mode_packet_based": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"mpls_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"flow-based", "packet-based"}, false),
						},
					},
				},
			},
			"ike_traceoptions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"files": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(2, 1000),
									},
									"match": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"size": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(10240, 1073741824),
									},
									"no_world_readable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"world_readable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"flag": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"no_remote_trace": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rate_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
					},
				},
			},
			"log": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"event_rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 1500),
						},
						"facility_override": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(listOfSyslogFacility(), false),
						},
						"file": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"files": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(2, 10),
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"size": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 10),
									},
								},
							},
						},
						"format": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"binary", "sd-syslog", "syslog"}, false),
						},
						"max_database_record": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 1000000),
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"event", "stream"}, false),
						},
						"rate_cap": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 5000),
						},
						"report": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"source_address": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  validation.IsIPAddress,
							ConflictsWith: []string{"log.0.source_interface"},
						},
						"source_interface": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"log.0.source_address"},
						},
						"transport": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"tcp", "tls", "udp"}, false),
									},
									"tcp_connections": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 5),
									},
									"tls_profile": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"utc_timestamp": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"utm": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"feature_profile_web_filtering_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"juniper-enhanced", "juniper-local", "web-filtering-none", "websense-redirect",
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceSecurityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return diag.FromErr(fmt.Errorf("security not compatible with Junos device %s",
			jnprSess.SystemInformation.HardwareModel))
	}
	sess.configLock(jnprSess)

	if err := setSecurity(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_security", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("security")

	return append(diagWarns, resourceSecurityReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSecurityReadWJnprSess(d, m, jnprSess)
}
func resourceSecurityReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	securityOptions, err := readSecurity(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSecurity(d, securityOptions)

	return nil
}
func resourceSecurityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSecurity(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSecurity(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_security", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityReadWJnprSess(d, m, jnprSess)...)
}
func resourceSecurityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
func resourceSecurityImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	securityOptions, err := readSecurity(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSecurity(d, securityOptions)
	d.SetId("security")
	result[0] = d

	return result, nil
}

func setSecurity(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set security "
	configSet := make([]string, 0)

	for _, v := range d.Get("alg").([]interface{}) {
		configSetAlg, err := setSecurityAlg(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetAlg...)
	}
	for _, v := range d.Get("flow").([]interface{}) {
		configSetFlow, err := setSecurityFlow(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetFlow...)
	}
	for _, v := range d.Get("forwarding_options").([]interface{}) {
		configSetForwOpts, err := setSecurityForwOpts(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetForwOpts...)
	}
	for _, ikeTrace := range d.Get("ike_traceoptions").([]interface{}) {
		configSetIkeTrace, err := setSecurityIkeTraceOpts(ikeTrace)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetIkeTrace...)
	}
	for _, v := range d.Get("log").([]interface{}) {
		configSetLog, err := setSecurityLog(v)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetLog...)
	}
	for _, v := range d.Get("utm").([]interface{}) {
		if v != nil {
			utm := v.(map[string]interface{})
			if utm["feature_profile_web_filtering_type"].(string) != "" {
				configSet = append(configSet, setPrefix+"utm feature-profile web-filtering type "+
					utm["feature_profile_web_filtering_type"].(string))
			}
		} else {
			return fmt.Errorf("utm block is empty")
		}
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func setSecurityAlg(alg interface{}) ([]string, error) {
	setPrefix := "set security alg "
	configSet := make([]string, 0)
	if alg != nil {
		algM := alg.(map[string]interface{})
		if algM["dns_disable"].(bool) {
			configSet = append(configSet, setPrefix+"dns disable")
		}
		if algM["ftp_disable"].(bool) {
			configSet = append(configSet, setPrefix+"ftp disable")
		}
		if algM["h323_disable"].(bool) {
			configSet = append(configSet, setPrefix+"h323 disable")
		}
		if algM["mgcp_disable"].(bool) {
			configSet = append(configSet, setPrefix+"mgcp disable")
		}
		if algM["msrpc_disable"].(bool) {
			configSet = append(configSet, setPrefix+"msrpc disable")
		}
		if algM["pptp_disable"].(bool) {
			configSet = append(configSet, setPrefix+"pptp disable")
		}
		if algM["rtsp_disable"].(bool) {
			configSet = append(configSet, setPrefix+"rtsp disable")
		}
		if algM["sccp_disable"].(bool) {
			configSet = append(configSet, setPrefix+"sccp disable")
		}
		if algM["sip_disable"].(bool) {
			configSet = append(configSet, setPrefix+"sip disable")
		}
		if algM["sunrpc_disable"].(bool) {
			configSet = append(configSet, setPrefix+"sunrpc disable")
		}
		if algM["talk_disable"].(bool) {
			configSet = append(configSet, setPrefix+"talk disable")
		}
		if algM["tftp_disable"].(bool) {
			configSet = append(configSet, setPrefix+"tftp disable")
		}
	} else {
		return configSet, fmt.Errorf("alg block is empty")
	}

	return configSet, nil
}
func setSecurityFlow(flow interface{}) ([]string, error) { // nolint: gocognit
	setPrefix := "set security flow "
	configSet := make([]string, 0)
	if flow != nil {
		flowM := flow.(map[string]interface{})
		for _, v := range flowM["advanced_options"].([]interface{}) {
			if v != nil {
				advOpt := v.(map[string]interface{})
				if advOpt["drop_matching_reserved_ip_address"].(bool) {
					configSet = append(configSet, setPrefix+"advanced-options drop-matching-reserved-ip-address")
				}
				if advOpt["drop_matching_link_local_address"].(bool) {
					configSet = append(configSet, setPrefix+"advanced-options drop-matching-link-local-address")
				}
				if advOpt["reverse_route_packet_mode_vr"].(bool) {
					configSet = append(configSet, setPrefix+"advanced-options reverse-route-packet-mode-vr")
				}
			} else {
				return configSet, fmt.Errorf("flow advanced_options block is empty")
			}
		}
		for _, v := range flowM["aging"].([]interface{}) {
			if v != nil {
				aging := v.(map[string]interface{})
				if aging["early_ageout"].(int) != 0 {
					configSet = append(configSet, setPrefix+"aging early-ageout "+
						strconv.Itoa(aging["early_ageout"].(int)))
				}
				if aging["high_watermark"].(int) != -1 {
					configSet = append(configSet, setPrefix+"aging high-watermark "+
						strconv.Itoa(aging["high_watermark"].(int)))
				}
				if aging["low_watermark"].(int) != -1 {
					configSet = append(configSet, setPrefix+"aging low-watermark "+
						strconv.Itoa(aging["low_watermark"].(int)))
				}
			} else {
				return configSet, fmt.Errorf("flow aging block is empty")
			}
		}
		if flowM["allow_dns_reply"].(bool) {
			configSet = append(configSet, setPrefix+"allow-dns-reply")
		}
		if flowM["allow_embedded_icmp"].(bool) {
			configSet = append(configSet, setPrefix+"allow-embedded-icmp")
		}
		if flowM["allow_reverse_ecmp"].(bool) {
			configSet = append(configSet, setPrefix+"allow-reverse-ecmp")
		}
		if flowM["enable_reroute_uniform_link_check_nat"].(bool) {
			configSet = append(configSet, setPrefix+"enable-reroute-uniform-link-check nat")
		}
		for _, v := range flowM["ethernet_switching"].([]interface{}) {
			if v != nil {
				ethSwitch := v.(map[string]interface{})
				if ethSwitch["block_non_ip_all"].(bool) {
					configSet = append(configSet, setPrefix+"ethernet-switching block-non-ip-all")
				}
				if ethSwitch["bypass_non_ip_unicast"].(bool) {
					configSet = append(configSet, setPrefix+"ethernet-switching bypass-non-ip-unicast")
				}
				if ethSwitch["bpdu_vlan_flooding"].(bool) {
					configSet = append(configSet, setPrefix+"ethernet-switching bpdu-vlan-flooding")
				}
				for _, v2 := range ethSwitch["no_packet_flooding"].([]interface{}) {
					configSet = append(configSet, setPrefix+"ethernet-switching no-packet-flooding")
					if v2 != nil {
						noPktFlo := v2.(map[string]interface{})
						if noPktFlo["no_trace_route"].(bool) {
							configSet = append(configSet, setPrefix+"ethernet-switching no-packet-flooding no-trace-route")
						}
					}
				}
			} else {
				return configSet, fmt.Errorf("flow ethernet_switching block is empty")
			}
		}
		if flowM["force_ip_reassembly"].(bool) {
			configSet = append(configSet, setPrefix+"force-ip-reassembly")
		}
		if flowM["ipsec_performance_acceleration"].(bool) {
			configSet = append(configSet, setPrefix+"ipsec-performance-acceleration")
		}
		if flowM["mcast_buffer_enhance"].(bool) {
			configSet = append(configSet, setPrefix+"mcast-buffer-enhance")
		}
		if flowM["pending_sess_queue_length"].(string) != "" {
			configSet = append(configSet, setPrefix+"pending-sess-queue-length "+
				flowM["pending_sess_queue_length"].(string))
		}
		if flowM["preserve_incoming_fragment_size"].(bool) {
			configSet = append(configSet, setPrefix+"preserve-incoming-fragment-size")
		}
		if flowM["route_change_timeout"].(int) != 0 {
			configSet = append(configSet, setPrefix+"route-change-timeout "+
				strconv.Itoa(flowM["route_change_timeout"].(int)))
		}
		if flowM["syn_flood_protection_mode"].(string) != "" {
			configSet = append(configSet, setPrefix+"syn-flood-protection-mode "+
				flowM["syn_flood_protection_mode"].(string))
		}
		if flowM["sync_icmp_session"].(bool) {
			configSet = append(configSet, setPrefix+"sync-icmp-session")
		}
		for _, v := range flowM["tcp_mss"].([]interface{}) {
			if v != nil {
				tcpMss := v.(map[string]interface{})
				if tcpMss["all_tcp_mss"].(int) != 0 {
					configSet = append(configSet, setPrefix+"tcp-mss all-tcp mss "+
						strconv.Itoa(tcpMss["all_tcp_mss"].(int)))
				}
				for _, v2 := range tcpMss["gre_in"].([]interface{}) {
					configSet = append(configSet, setPrefix+"tcp-mss gre-in")
					if v2 != nil {
						greIn := v2.(map[string]interface{})
						if greIn["mss"].(int) != 0 {
							configSet = append(configSet, setPrefix+"tcp-mss gre-in mss "+
								strconv.Itoa(greIn["mss"].(int)))
						}
					}
				}
				for _, v2 := range tcpMss["gre_out"].([]interface{}) {
					configSet = append(configSet, setPrefix+"tcp-mss gre-out")
					if v2 != nil {
						greOut := v2.(map[string]interface{})
						if greOut["mss"].(int) != 0 {
							configSet = append(configSet, setPrefix+"tcp-mss gre-out mss "+
								strconv.Itoa(greOut["mss"].(int)))
						}
					}
				}
				for _, v2 := range tcpMss["ipsec_vpn"].([]interface{}) {
					configSet = append(configSet, setPrefix+"tcp-mss ipsec-vpn")
					if v2 != nil {
						ipsecVpn := v2.(map[string]interface{})
						if ipsecVpn["mss"].(int) != 0 {
							configSet = append(configSet, setPrefix+"tcp-mss ipsec-vpn mss "+
								strconv.Itoa(ipsecVpn["mss"].(int)))
						}
					}
				}
			} else {
				return configSet, fmt.Errorf("flow tcp_mss block is empty")
			}
		}
		for _, v := range flowM["tcp_session"].([]interface{}) {
			if v != nil {
				tcpSess := v.(map[string]interface{})
				if tcpSess["fin_invalidate_session"].(bool) {
					configSet = append(configSet, setPrefix+"tcp-session fin-invalidate-session")
				}
				if tcpSess["maximum_window"].(string) != "" {
					configSet = append(configSet, setPrefix+"tcp-session maximum-window "+
						tcpSess["maximum_window"].(string))
				}
				if tcpSess["no_sequence_check"].(bool) {
					configSet = append(configSet, setPrefix+"tcp-session no-sequence-check")
				}
				if tcpSess["no_syn_check"].(bool) {
					configSet = append(configSet, setPrefix+"tcp-session no-syn-check")
				}
				if tcpSess["no_syn_check_in_tunnel"].(bool) {
					configSet = append(configSet, setPrefix+"tcp-session no-syn-check-in-tunnel")
				}
				if tcpSess["rst_invalidate_session"].(bool) {
					configSet = append(configSet, setPrefix+"tcp-session rst-invalidate-session")
				}
				if tcpSess["rst_sequence_check"].(bool) {
					configSet = append(configSet, setPrefix+"tcp-session rst-sequence-check")
				}
				if tcpSess["strict_syn_check"].(bool) {
					configSet = append(configSet, setPrefix+"tcp-session strict-syn-check")
				}
				if tcpSess["tcp_initial_timeout"].(int) != 0 {
					configSet = append(configSet, setPrefix+"tcp-session tcp-initial-timeout "+
						strconv.Itoa(tcpSess["tcp_initial_timeout"].(int)))
				}
				for _, v2 := range tcpSess["time_wait_state"].([]interface{}) {
					configSet = append(configSet, setPrefix+"tcp-session time-wait-state")
					if v2 != nil {
						timWaiSta := v2.(map[string]interface{})
						if timWaiSta["apply_to_half_close_state"].(bool) {
							configSet = append(configSet, setPrefix+"tcp-session time-wait-state apply-to-half-close-state")
						}
						if timWaiSta["session_ageout"].(bool) {
							configSet = append(configSet, setPrefix+"tcp-session time-wait-state session-ageout")
						}
						if timWaiSta["session_timeout"].(int) != 0 {
							configSet = append(configSet, setPrefix+"tcp-session time-wait-state session-timeout "+
								strconv.Itoa(timWaiSta["session_timeout"].(int)))
						}
					}
				}
			} else {
				return configSet, fmt.Errorf("flow tcp_session block is empty")
			}
		}
	} else {
		return configSet, fmt.Errorf("flow block is empty")
	}

	return configSet, nil
}
func setSecurityForwOpts(forwOpts interface{}) ([]string, error) {
	setPrefix := "set security forwarding-options "
	configSet := make([]string, 0)
	if forwOpts != nil {
		forwOptsM := forwOpts.(map[string]interface{})
		if forwOptsM["inet6_mode"].(string) != "" {
			configSet = append(configSet, setPrefix+"family inet6 mode "+forwOptsM["inet6_mode"].(string))
		}
		if forwOptsM["iso_mode_packet_based"].(bool) {
			configSet = append(configSet, setPrefix+"family iso mode packet-based")
		}
		if forwOptsM["mpls_mode"].(string) != "" {
			configSet = append(configSet, setPrefix+"family mpls mode "+forwOptsM["mpls_mode"].(string))
		}
	} else {
		return configSet, fmt.Errorf("forwarding_options block is empty")
	}

	return configSet, nil
}
func setSecurityIkeTraceOpts(ikeTrace interface{}) ([]string, error) {
	setPrefix := "set security ike traceoptions "
	configSet := make([]string, 0)
	if ikeTrace != nil {
		ikeTraceM := ikeTrace.(map[string]interface{})
		for _, v := range ikeTraceM["file"].([]interface{}) {
			if v != nil {
				ikeTraceFile := v.(map[string]interface{})
				if ikeTraceFile["name"].(string) != "" {
					configSet = append(configSet, setPrefix+"file "+
						ikeTraceFile["name"].(string))
				}
				if ikeTraceFile["files"].(int) > 0 {
					configSet = append(configSet, setPrefix+"file files "+
						strconv.Itoa(ikeTraceFile["files"].(int)))
				}
				if ikeTraceFile["match"].(string) != "" {
					configSet = append(configSet, setPrefix+"file match \""+
						ikeTraceFile["match"].(string)+"\"")
				}
				if ikeTraceFile["size"].(int) > 0 {
					configSet = append(configSet, setPrefix+"file size "+
						strconv.Itoa(ikeTraceFile["size"].(int)))
				}
				if ikeTraceFile["world_readable"].(bool) && ikeTraceFile["no_world_readable"].(bool) {
					return configSet, fmt.Errorf("conflict between 'world_readable' and 'no_world_readable' for ike_traceoptions file")
				}
				if ikeTraceFile["world_readable"].(bool) {
					configSet = append(configSet, setPrefix+"file world-readable")
				}
				if ikeTraceFile["no_world_readable"].(bool) {
					configSet = append(configSet, setPrefix+"file no-world-readable")
				}
			} else {
				return configSet, fmt.Errorf("ike_traceoptions file block is empty")
			}
		}
		for _, ikeTraceFlag := range ikeTraceM["flag"].([]interface{}) {
			configSet = append(configSet, setPrefix+"flag "+ikeTraceFlag.(string))
		}
		if ikeTraceM["no_remote_trace"].(bool) {
			configSet = append(configSet, setPrefix+"no-remote-trace")
		}
		if ikeTraceM["rate_limit"].(int) > -1 {
			configSet = append(configSet, setPrefix+"rate-limit "+
				strconv.Itoa(ikeTraceM["rate_limit"].(int)))
		}
	} else {
		return configSet, fmt.Errorf("ike_traceoptions block is empty")
	}

	return configSet, nil
}
func setSecurityLog(log interface{}) ([]string, error) {
	setPrefix := "set security log "
	configSet := make([]string, 0)
	if log != nil {
		logM := log.(map[string]interface{})
		if logM["disable"].(bool) {
			configSet = append(configSet, setPrefix+"disable")
		}
		if logM["event_rate"].(int) != -1 {
			configSet = append(configSet, setPrefix+"event-rate "+strconv.Itoa(logM["event_rate"].(int)))
		}
		if logM["facility_override"].(string) != "" {
			configSet = append(configSet, setPrefix+"facility-override "+logM["facility_override"].(string))
		}
		for _, v := range logM["file"].([]interface{}) {
			if v != nil {
				file := v.(map[string]interface{})
				if file["files"].(int) != 0 {
					configSet = append(configSet, setPrefix+"file files "+strconv.Itoa(file["files"].(int)))
				}
				if file["name"].(string) != "" {
					configSet = append(configSet, setPrefix+"file name "+file["name"].(string))
				}
				if file["path"].(string) != "" {
					configSet = append(configSet, setPrefix+"file path "+file["path"].(string))
				}
				if file["size"].(int) != 0 {
					configSet = append(configSet, setPrefix+"file size "+strconv.Itoa(file["size"].(int)))
				}
			} else {
				return configSet, fmt.Errorf("log file block is empty")
			}
		}
		if logM["format"].(string) != "" {
			configSet = append(configSet, setPrefix+"format "+logM["format"].(string))
		}
		if logM["max_database_record"].(int) != -1 {
			configSet = append(configSet, setPrefix+"max-database-record "+strconv.Itoa(logM["max_database_record"].(int)))
		}
		if logM["mode"].(string) != "" {
			configSet = append(configSet, setPrefix+"mode "+logM["mode"].(string))
		}
		if logM["rate_cap"].(int) != -1 {
			configSet = append(configSet, setPrefix+"rate-cap "+strconv.Itoa(logM["rate_cap"].(int)))
		}
		if logM["report"].(bool) {
			configSet = append(configSet, setPrefix+"report")
		}
		if logM["source_address"].(string) != "" {
			configSet = append(configSet, setPrefix+"source-address "+logM["source_address"].(string))
		}
		if logM["source_interface"].(string) != "" {
			configSet = append(configSet, setPrefix+"source-interface "+logM["source_interface"].(string))
		}
		for _, v := range logM["transport"].([]interface{}) {
			configSet = append(configSet, setPrefix+"transport")
			if v != nil {
				trans := v.(map[string]interface{})
				if trans["protocol"].(string) != "" {
					configSet = append(configSet, setPrefix+"transport protocol "+trans["protocol"].(string))
				}
				if trans["tcp_connections"].(int) != 0 {
					configSet = append(configSet, setPrefix+"transport tcp-connections "+strconv.Itoa(trans["tcp_connections"].(int)))
				}
				if trans["tls_profile"].(string) != "" {
					configSet = append(configSet, setPrefix+"transport tls-profile "+trans["tls_profile"].(string))
				}
			}
		}
		if logM["utc_timestamp"].(bool) {
			configSet = append(configSet, setPrefix+"utc-timestamp")
		}
	} else {
		return configSet, fmt.Errorf("log block is empty")
	}

	return configSet, nil
}

func listLinesSecurityAlg() []string {
	return []string{
		"alg dns disable",
		"alg ftp disable",
		"alg h323 disable",
		"alg mgcp disable",
		"alg msrpc disable",
		"alg pptp disable",
		"alg rtsp disable",
		"alg sccp disable",
		"alg sip disable",
		"alg sunrpc disable",
		"alg talk disable",
		"alg tftp disable",
	}
}
func listLinesSecurityFlow() []string {
	return []string{
		"flow advanced-options",
		"flow aging",
		"flow allow-dns-reply",
		"flow allow-embedded-icmp",
		"flow allow-reverse-ecmp",
		"flow enable-reroute-uniform-link-check",
		"flow ethernet-switching",
		"flow force-ip-reassembly",
		"flow ipsec-performance-acceleration",
		"flow mcast-buffer-enhance",
		"flow pending-sess-queue-length",
		"flow preserve-incoming-fragment-size",
		"flow route-change-timeout",
		"flow syn-flood-protection-mode",
		"flow sync-icmp-session",
		"flow tcp-mss",
		"flow tcp-session",
	}
}
func listLinesSecurityForwardingOptions() []string {
	return []string{
		"forwarding-options family mpls mode",
		"forwarding-options family inet6 mode",
		"forwarding-options family iso mode",
	}
}
func listLinesSecurityLog() []string {
	return []string{
		"log disable",
		"log event-rate",
		"log facility-override",
		"log file",
		"log format",
		"log max-database-record",
		"log mode",
		"log rate-cap",
		"log report",
		"log source-address",
		"log source-interface",
		"log transport",
		"log utc-timestamp",
	}
}
func listLinesSecurityUtm() []string {
	return []string{
		"utm feature-profile web-filtering type",
	}
}

func delSecurity(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := []string{
		"ike traceoptions",
	}
	listLinesToDelete = append(listLinesToDelete, listLinesSecurityAlg()...)
	listLinesToDelete = append(listLinesToDelete, listLinesSecurityFlow()...)
	listLinesToDelete = append(listLinesToDelete, listLinesSecurityForwardingOptions()...)
	listLinesToDelete = append(listLinesToDelete, listLinesSecurityLog()...)
	listLinesToDelete = append(listLinesToDelete, listLinesSecurityUtm()...)
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete security "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSecurity(m interface{}, jnprSess *NetconfObject) (securityOptions, error) {
	sess := m.(*Session)
	var confRead securityOptions

	securityConfig, err := sess.command("show configuration security"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if securityConfig != emptyWord {
		for _, item := range strings.Split(securityConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case checkStringHasPrefixInList(itemTrim, listLinesSecurityAlg()):
				readSecurityAlg(&confRead, itemTrim)
			case checkStringHasPrefixInList(itemTrim, listLinesSecurityFlow()):
				err := readSecurityFlow(&confRead, itemTrim)
				if err != nil {
					return confRead, err
				}
			case checkStringHasPrefixInList(itemTrim, listLinesSecurityForwardingOptions()):
				readSecurityForwardingOpts(&confRead, itemTrim)

			case checkStringHasPrefixInList(itemTrim, listLinesSecurityLog()):
				err := readSecurityLog(&confRead, itemTrim)
				if err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "ike traceoptions"):
				err := readSecurityIkeTraceOptions(&confRead, itemTrim)
				if err != nil {
					return confRead, err
				}
			case checkStringHasPrefixInList(itemTrim, listLinesSecurityUtm()):
				if len(confRead.utm) == 0 {
					confRead.utm = append(confRead.utm, map[string]interface{}{
						"feature_profile_web_filtering_type": "",
					})
				}
				if strings.HasPrefix(itemTrim, "utm feature-profile web-filtering type ") {
					confRead.utm[0]["feature_profile_web_filtering_type"] = strings.TrimPrefix(itemTrim,
						"utm feature-profile web-filtering type ")
				}
			}
		}
	}

	return confRead, nil
}

func readSecurityAlg(confRead *securityOptions, itemTrimAlg string) {
	itemTrim := strings.TrimPrefix(itemTrimAlg, "alg ")
	if len(confRead.alg) == 0 {
		confRead.alg = append(confRead.alg, map[string]interface{}{
			"dns_disable":    false,
			"ftp_disable":    false,
			"h323_disable":   false,
			"mgcp_disable":   false,
			"msrpc_disable":  false,
			"pptp_disable":   false,
			"rtsp_disable":   false,
			"sccp_disable":   false,
			"sip_disable":    false,
			"sunrpc_disable": false,
			"talk_disable":   false,
			"tftp_disable":   false,
		})
	}
	if itemTrim == "dns disable" {
		confRead.alg[0]["dns_disable"] = true
	}
	if itemTrim == "ftp disable" {
		confRead.alg[0]["ftp_disable"] = true
	}
	if itemTrim == "h323 disable" {
		confRead.alg[0]["h323_disable"] = true
	}
	if itemTrim == "mgcp disable" {
		confRead.alg[0]["mgcp_disable"] = true
	}
	if itemTrim == "msrpc disable" {
		confRead.alg[0]["msrpc_disable"] = true
	}
	if itemTrim == "pptp disable" {
		confRead.alg[0]["pptp_disable"] = true
	}
	if itemTrim == "rtsp disable" {
		confRead.alg[0]["rtsp_disable"] = true
	}
	if itemTrim == "sccp disable" {
		confRead.alg[0]["sccp_disable"] = true
	}
	if itemTrim == "sip disable" {
		confRead.alg[0]["sip_disable"] = true
	}
	if itemTrim == "sunrpc disable" {
		confRead.alg[0]["sunrpc_disable"] = true
	}
	if itemTrim == "talk disable" {
		confRead.alg[0]["talk_disable"] = true
	}
	if itemTrim == "tftp disable" {
		confRead.alg[0]["tftp_disable"] = true
	}
}
func readSecurityFlow(confRead *securityOptions, itemTrimFlow string) error {
	itemTrim := strings.TrimPrefix(itemTrimFlow, "flow ")
	if len(confRead.flow) == 0 {
		confRead.flow = append(confRead.flow, map[string]interface{}{
			"advanced_options":                make([]map[string]interface{}, 0),
			"aging":                           make([]map[string]interface{}, 0),
			"allow_dns_reply":                 false,
			"allow_embedded_icmp":             false,
			"allow_reverse_ecmp":              false,
			"ethernet_switching":              make([]map[string]interface{}, 0),
			"force_ip_reassembly":             false,
			"ipsec_performance_acceleration":  false,
			"mcast_buffer_enhance":            false,
			"pending_sess_queue_length":       "",
			"preserve_incoming_fragment_size": false,
			"route_change_timeout":            0,
			"syn_flood_protection_mode":       "",
			"sync_icmp_session":               false,
			"tcp_mss":                         make([]map[string]interface{}, 0),
			"tcp_session":                     make([]map[string]interface{}, 0),
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "advanced-options"):
		if len(confRead.flow[0]["advanced_options"].([]map[string]interface{})) == 0 {
			confRead.flow[0]["advanced_options"] = append(
				confRead.flow[0]["advanced_options"].([]map[string]interface{}), map[string]interface{}{
					"drop_matching_reserved_ip_address": false,
					"drop_matching_link_local_address":  false,
					"reverse_route_packet_mode_vr":      false,
				})
		}
		switch {
		case itemTrim == "advanced-options drop-matching-reserved-ip-address":
			confRead.flow[0]["advanced_options"].([]map[string]interface{})[0]["drop_matching_reserved_ip_address"] = true
		case itemTrim == "advanced-options drop-matching-link-local-address":
			confRead.flow[0]["advanced_options"].([]map[string]interface{})[0]["drop_matching_link_local_address"] = true
		case itemTrim == "advanced-options reverse-route-packet-mode-vr":
			confRead.flow[0]["advanced_options"].([]map[string]interface{})[0]["reverse_route_packet_mode_vr"] = true
		}
	case strings.HasPrefix(itemTrim, "aging"):
		if len(confRead.flow[0]["aging"].([]map[string]interface{})) == 0 {
			confRead.flow[0]["aging"] = append(
				confRead.flow[0]["aging"].([]map[string]interface{}), map[string]interface{}{
					"early_ageout":   0,
					"high_watermark": -1,
					"low_watermark":  -1,
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "aging early-ageout "):
			var err error
			confRead.flow[0]["aging"].([]map[string]interface{})[0]["early_ageout"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "aging early-ageout "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "aging high-watermark"):
			var err error
			confRead.flow[0]["aging"].([]map[string]interface{})[0]["high_watermark"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "aging high-watermark "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "aging low-watermark"):
			var err error
			confRead.flow[0]["aging"].([]map[string]interface{})[0]["low_watermark"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "aging low-watermark "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case itemTrim == "allow-dns-reply":
		confRead.flow[0]["allow_dns_reply"] = true
	case itemTrim == "allow-embedded-icmp":
		confRead.flow[0]["allow_embedded_icmp"] = true
	case itemTrim == "allow-reverse-ecmp":
		confRead.flow[0]["allow_reverse_ecmp"] = true
	case itemTrim == "enable-reroute-uniform-link-check nat":
		confRead.flow[0]["enable_reroute_uniform_link_check_nat"] = true
	case strings.HasPrefix(itemTrim, "ethernet-switching"):
		if len(confRead.flow[0]["ethernet_switching"].([]map[string]interface{})) == 0 {
			confRead.flow[0]["ethernet_switching"] = append(
				confRead.flow[0]["ethernet_switching"].([]map[string]interface{}), map[string]interface{}{
					"block_non_ip_all":      false,
					"bypass_non_ip_unicast": false,
					"bpdu_vlan_flooding":    false,
					"no_packet_flooding":    make([]map[string]interface{}, 0),
				})
		}
		switch {
		case itemTrim == "ethernet-switching block-non-ip-all":
			confRead.flow[0]["ethernet_switching"].([]map[string]interface{})[0]["block_non_ip_all"] = true
		case itemTrim == "ethernet-switching bypass-non-ip-unicast":
			confRead.flow[0]["ethernet_switching"].([]map[string]interface{})[0]["bypass_non_ip_unicast"] = true
		case itemTrim == "ethernet-switching bpdu-vlan-flooding":
			confRead.flow[0]["ethernet_switching"].([]map[string]interface{})[0]["bpdu_vlan_flooding"] = true
		case strings.HasPrefix(itemTrim, "ethernet-switching no-packet-flooding"):
			if len(confRead.flow[0]["ethernet_switching"].([]map[string]interface{})[0]["no_packet_flooding"].([]map[string]interface{})) == 0 { // nolint: lll
				confRead.flow[0]["ethernet_switching"].([]map[string]interface{})[0]["no_packet_flooding"] = append(
					confRead.flow[0]["ethernet_switching"].([]map[string]interface{})[0]["no_packet_flooding"].([]map[string]interface{}), // nolint: lll
					map[string]interface{}{
						"no_trace_route": false,
					})
			}
			if itemTrim == "ethernet-switching no-packet-flooding no-trace-route" {
				confRead.flow[0]["ethernet_switching"].([]map[string]interface{})[0]["no_packet_flooding"].([]map[string]interface{})[0]["no_trace_route"] = true // nolint: lll
			}
		}
	case itemTrim == "force-ip-reassembly":
		confRead.flow[0]["force_ip_reassembly"] = true
	case itemTrim == "ipsec-performance-acceleration":
		confRead.flow[0]["ipsec_performance_acceleration"] = true
	case itemTrim == "mcast-buffer-enhance":
		confRead.flow[0]["mcast_buffer_enhance"] = true
	case strings.HasPrefix(itemTrim, "pending-sess-queue-length "):
		confRead.flow[0]["pending_sess_queue_length"] = strings.TrimPrefix(itemTrim, "pending-sess-queue-length ")
	case itemTrim == "preserve-incoming-fragment-size":
		confRead.flow[0]["preserve_incoming_fragment_size"] = true
	case strings.HasPrefix(itemTrim, "route-change-timeout "):
		var err error
		confRead.flow[0]["route_change_timeout"], err = strconv.Atoi(
			strings.TrimPrefix(itemTrim, "route-change-timeout "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "syn-flood-protection-mode "):
		confRead.flow[0]["syn_flood_protection_mode"] = strings.TrimPrefix(itemTrim, "syn-flood-protection-mode ")
	case itemTrim == "sync-icmp-session":
		confRead.flow[0]["sync_icmp_session"] = true
	case strings.HasPrefix(itemTrim, "tcp-mss "):
		if len(confRead.flow[0]["tcp_mss"].([]map[string]interface{})) == 0 {
			confRead.flow[0]["tcp_mss"] = append(
				confRead.flow[0]["tcp_mss"].([]map[string]interface{}), map[string]interface{}{
					"all_tcp_mss": 0,
					"gre_in":      make([]map[string]interface{}, 0),
					"gre_out":     make([]map[string]interface{}, 0),
					"ipsec_vpn":   make([]map[string]interface{}, 0),
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "tcp-mss all-tcp mss "):
			var err error
			confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["all_tcp_mss"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "tcp-mss all-tcp mss "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "tcp-mss gre-in"):
			if len(confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["gre_in"].([]map[string]interface{})) == 0 {
				confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["gre_in"] = append(
					confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["gre_in"].([]map[string]interface{}),
					map[string]interface{}{
						"mss": 0,
					})
			}
			if strings.HasPrefix(itemTrim, "tcp-mss gre-in mss ") {
				var err error
				confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["gre_in"].([]map[string]interface{})[0]["mss"],
					err = strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp-mss gre-in mss "))
				if err != nil {
					return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		case strings.HasPrefix(itemTrim, "tcp-mss gre-out"):
			if len(confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["gre_out"].([]map[string]interface{})) == 0 {
				confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["gre_out"] = append(
					confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["gre_out"].([]map[string]interface{}),
					map[string]interface{}{
						"mss": 0,
					})
			}
			if strings.HasPrefix(itemTrim, "tcp-mss gre-out mss ") {
				var err error
				confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["gre_out"].([]map[string]interface{})[0]["mss"],
					err = strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp-mss gre-out mss "))
				if err != nil {
					return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		case strings.HasPrefix(itemTrim, "tcp-mss ipsec-vpn"):
			if len(confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["ipsec_vpn"].([]map[string]interface{})) == 0 {
				confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["ipsec_vpn"] = append(
					confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["ipsec_vpn"].([]map[string]interface{}),
					map[string]interface{}{
						"mss": 0,
					})
			}
			if strings.HasPrefix(itemTrim, "tcp-mss ipsec-vpn mss ") {
				var err error
				confRead.flow[0]["tcp_mss"].([]map[string]interface{})[0]["ipsec_vpn"].([]map[string]interface{})[0]["mss"],
					err = strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp-mss ipsec-vpn mss "))
				if err != nil {
					return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		}
	case strings.HasPrefix(itemTrim, "tcp-session "):
		if len(confRead.flow[0]["tcp_session"].([]map[string]interface{})) == 0 {
			confRead.flow[0]["tcp_session"] = append(
				confRead.flow[0]["tcp_session"].([]map[string]interface{}), map[string]interface{}{
					"fin_invalidate_session": false,
					"maximum_window":         "",
					"no_sequence_check":      false,
					"no_syn_check":           false,
					"no_syn_check_in_tunnel": false,
					"rst_invalidate_session": false,
					"rst_sequence_check":     false,
					"strict_syn_check":       false,
					"tcp_initial_timeout":    0,
					"time_wait_state":        make([]map[string]interface{}, 0),
				})
		}
		switch {
		case itemTrim == "tcp-session fin-invalidate-session":
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["fin_invalidate_session"] = true
		case strings.HasPrefix(itemTrim, "tcp-session maximum-window "):
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["maximum_window"] =
				strings.TrimPrefix(itemTrim, "tcp-session maximum-window ")
		case itemTrim == "tcp-session no-sequence-check":
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["no_sequence_check"] = true
		case itemTrim == "tcp-session no-syn-check":
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["no_syn_check"] = true
		case itemTrim == "tcp-session no-syn-check-in-tunnel":
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["no_syn_check_in_tunnel"] = true
		case itemTrim == "tcp-session rst-invalidate-session":
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["rst_invalidate_session"] = true
		case itemTrim == "tcp-session rst-sequence-check":
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["rst_sequence_check"] = true
		case itemTrim == "tcp-session strict-syn-check":
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["strict_syn_check"] = true
		case strings.HasPrefix(itemTrim, "tcp-session tcp-initial-timeout "):
			var err error
			confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["tcp_initial_timeout"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp-session tcp-initial-timeout "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "tcp-session time-wait-state"):
			if len(confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["time_wait_state"].([]map[string]interface{})) == 0 { // nolint: lll
				confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["time_wait_state"] = append(
					confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["time_wait_state"].([]map[string]interface{}),
					map[string]interface{}{
						"apply_to_half_close_state": false,
						"session_ageout":            false,
						"session_timeout":           0,
					})
			}
			switch {
			case itemTrim == "tcp-session time-wait-state apply-to-half-close-state":
				confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["time_wait_state"].([]map[string]interface{})[0]["apply_to_half_close_state"] = true // nolint: lll
			case itemTrim == "tcp-session time-wait-state session-ageout":
				confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["time_wait_state"].([]map[string]interface{})[0]["session_ageout"] = true // nolint: lll
			case strings.HasPrefix(itemTrim, "tcp-session time-wait-state session-timeout "):
				var err error
				confRead.flow[0]["tcp_session"].([]map[string]interface{})[0]["time_wait_state"].([]map[string]interface{})[0]["session_timeout"], err = // nolint: lll
					strconv.Atoi(strings.TrimPrefix(itemTrim, "tcp-session time-wait-state session-timeout "))
				if err != nil {
					return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			}
		}
	}

	return nil
}
func readSecurityForwardingOpts(confRead *securityOptions, itemTrimFwOpts string) {
	itemTrim := strings.TrimPrefix(itemTrimFwOpts, "forwarding-options ")
	if len(confRead.forwardingOpts) == 0 {
		confRead.forwardingOpts = append(confRead.forwardingOpts, map[string]interface{}{
			"inet6_mode":            "",
			"iso_mode_packet_based": false,
			"mpls_mode":             "",
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "family inet6 mode "):
		confRead.forwardingOpts[0]["inet6_mode"] = strings.TrimPrefix(itemTrim, "family inet6 mode ")
	case itemTrim == "family iso mode packet-based":
		confRead.forwardingOpts[0]["iso_mode_packet_based"] = true
	case strings.HasPrefix(itemTrim, "family mpls mode "):
		confRead.forwardingOpts[0]["mpls_mode"] = strings.TrimPrefix(itemTrim, "family mpls mode ")
	}
}
func readSecurityIkeTraceOptions(confRead *securityOptions, itemTrimIkeTraceOpts string) error {
	itemTrim := strings.TrimPrefix(itemTrimIkeTraceOpts, "ike traceoptions ")
	if len(confRead.ikeTraceoptions) == 0 {
		confRead.ikeTraceoptions = append(confRead.ikeTraceoptions, map[string]interface{}{
			"file":            make([]map[string]interface{}, 0),
			"flag":            make([]string, 0),
			"no_remote_trace": false,
			"rate_limit":      -1,
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "file"):
		if len(confRead.ikeTraceoptions[0]["file"].([]map[string]interface{})) == 0 {
			confRead.ikeTraceoptions[0]["file"] = append(
				confRead.ikeTraceoptions[0]["file"].([]map[string]interface{}), map[string]interface{}{
					"name":              "",
					"files":             0,
					"match":             "",
					"size":              0,
					"world_readable":    false,
					"no_world_readable": false,
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "file files"):
			var err error
			confRead.ikeTraceoptions[0]["file"].([]map[string]interface{})[0]["files"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "file files "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "file match"):
			confRead.ikeTraceoptions[0]["file"].([]map[string]interface{})[0]["match"] = strings.Trim(
				strings.TrimPrefix(itemTrim, "file match "), "\"")
		case strings.HasPrefix(itemTrim, "file size"):
			var err error
			confRead.ikeTraceoptions[0]["file"].([]map[string]interface{})[0]["size"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "file size "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case itemTrim == "file world-readable":
			confRead.ikeTraceoptions[0]["file"].([]map[string]interface{})[0]["world_readable"] = true
		case itemTrim == "file no-world-readable":
			confRead.ikeTraceoptions[0]["file"].([]map[string]interface{})[0]["no_world_readable"] = true
		case strings.HasPrefix(itemTrim, "file "):
			confRead.ikeTraceoptions[0]["file"].([]map[string]interface{})[0]["name"] = strings.Trim(
				strings.TrimPrefix(itemTrim, "file "), "\"")
		}
	case strings.HasPrefix(itemTrim, "flag"):
		confRead.ikeTraceoptions[0]["flag"] = append(confRead.ikeTraceoptions[0]["flag"].([]string),
			strings.TrimPrefix(itemTrim, "flag "))
	case itemTrim == "no-remote-trace":
		confRead.ikeTraceoptions[0]["no_remote_trace"] = true
	case strings.HasPrefix(itemTrim, "rate-limit"):
		var err error
		confRead.ikeTraceoptions[0]["rate_limit"], err = strconv.Atoi(
			strings.TrimPrefix(itemTrim, "rate-limit "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}
func readSecurityLog(confRead *securityOptions, itemTrimLog string) error {
	itemTrim := strings.TrimPrefix(itemTrimLog, "log ")
	if len(confRead.log) == 0 {
		confRead.log = append(confRead.log, map[string]interface{}{
			"disable":             false,
			"event_rate":          -1,
			"facility_override":   "",
			"file":                make([]map[string]interface{}, 0),
			"format":              "",
			"max_database_record": -1,
			"mode":                "",
			"rate_cap":            -1,
			"report":              false,
			"source_address":      "",
			"source_interface":    "",
			"transport":           make([]map[string]interface{}, 0),
			"utc_timestamp":       false,
		})
	}
	switch {
	case itemTrim == disableW:
		confRead.log[0]["disable"] = true
	case strings.HasPrefix(itemTrim, "event-rate "):
		var err error
		confRead.log[0]["event_rate"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "event-rate "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "facility-override "):
		confRead.log[0]["facility_override"] = strings.TrimPrefix(itemTrim, "facility-override ")
	case strings.HasPrefix(itemTrim, "file"):
		if len(confRead.log[0]["file"].([]map[string]interface{})) == 0 {
			confRead.log[0]["file"] = append(confRead.log[0]["file"].([]map[string]interface{}), map[string]interface{}{
				"files": 0,
				"name":  "",
				"path":  "",
				"size":  0,
			})
		}
		switch {
		case strings.HasPrefix(itemTrim, "file files "):
			var err error
			confRead.log[0]["file"].([]map[string]interface{})[0]["files"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "file files "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "file name "):
			confRead.log[0]["file"].([]map[string]interface{})[0]["name"] = strings.TrimPrefix(itemTrim, "file name ")
		case strings.HasPrefix(itemTrim, "file path "):
			confRead.log[0]["file"].([]map[string]interface{})[0]["path"] = strings.TrimPrefix(itemTrim, "file path ")
		case strings.HasPrefix(itemTrim, "file size "):
			var err error
			confRead.log[0]["file"].([]map[string]interface{})[0]["size"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "file size "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "format "):
		confRead.log[0]["format"] = strings.TrimPrefix(itemTrim, "format ")
	case strings.HasPrefix(itemTrim, "max-database-record "):
		var err error
		confRead.log[0]["max_database_record"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-database-record "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "mode "):
		confRead.log[0]["mode"] = strings.TrimPrefix(itemTrim, "mode ")
	case strings.HasPrefix(itemTrim, "rate-cap "):
		var err error
		confRead.log[0]["rate_cap"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "rate-cap "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case itemTrim == "report":
		confRead.log[0]["report"] = true
	case strings.HasPrefix(itemTrim, "source-address "):
		confRead.log[0]["source_address"] = strings.TrimPrefix(itemTrim, "source-address ")
	case strings.HasPrefix(itemTrim, "source-interface "):
		confRead.log[0]["source_interface"] = strings.TrimPrefix(itemTrim, "source-interface ")
	case strings.HasPrefix(itemTrim, "transport"):
		if len(confRead.log[0]["transport"].([]map[string]interface{})) == 0 {
			confRead.log[0]["transport"] = append(
				confRead.log[0]["transport"].([]map[string]interface{}), map[string]interface{}{
					"protocol":        "",
					"tcp_connections": 0,
					"tls_profile":     "",
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "transport protocol "):
			confRead.log[0]["transport"].([]map[string]interface{})[0]["protocol"] =
				strings.TrimPrefix(itemTrim, "transport protocol ")
		case strings.HasPrefix(itemTrim, "transport tcp-connections "):
			var err error
			confRead.log[0]["transport"].([]map[string]interface{})[0]["tcp_connections"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "transport tcp-connections "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "transport tls-profile "):
			confRead.log[0]["transport"].([]map[string]interface{})[0]["tls_profile"] =
				strings.TrimPrefix(itemTrim, "transport tls-profile ")
		}
	case itemTrim == "utc-timestamp":
		confRead.log[0]["utc_timestamp"] = true
	}

	return nil
}

func fillSecurity(d *schema.ResourceData, securityOptions securityOptions) {
	if tfErr := d.Set("alg", securityOptions.alg); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("flow", securityOptions.flow); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forwarding_options", securityOptions.forwardingOpts); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ike_traceoptions", securityOptions.ikeTraceoptions); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("log", securityOptions.log); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("utm", securityOptions.utm); tfErr != nil {
		panic(tfErr)
	}
}
