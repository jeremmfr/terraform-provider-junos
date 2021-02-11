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

type systemOptions struct {
	autoSnapshot                         bool
	defaultAddressSelection              bool
	noMulticastEcho                      bool
	noPingRecordRoute                    bool
	noPingTimeStamp                      bool
	noRedirects                          bool
	noRedirectsIPv6                      bool
	maxConfigurationRollbacks            int
	maxConfigurationsOnFlash             int
	domainName                           string
	hostName                             string
	timeZone                             string
	tracingDestinationOverrideSyslogHost string
	authenticationOrder                  []string
	nameServer                           []string
	inet6BackupRouter                    []map[string]interface{}
	internetOptions                      []map[string]interface{}
	login                                []map[string]interface{}
	services                             []map[string]interface{}
	syslog                               []map[string]interface{}
}

func resourceSystem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemCreate,
		ReadContext:   resourceSystemRead,
		UpdateContext: resourceSystemUpdate,
		DeleteContext: resourceSystemDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSystemImport,
		},
		Schema: map[string]*schema.Schema{
			"authentication_order": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_address_selection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inet6_backup_router": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"destination": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"internet_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gre_path_mtu_discovery": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.no_gre_path_mtu_discovery"},
						},
						"icmpv4_rate_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_size": {
										Type:         schema.TypeInt,
										Default:      -1,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"packet_rate": {
										Type:         schema.TypeInt,
										Default:      -1,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
								},
							},
						},
						"icmpv6_rate_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_size": {
										Type:         schema.TypeInt,
										Default:      -1,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"packet_rate": {
										Type:         schema.TypeInt,
										Default:      -1,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
								},
							},
						},
						"ipip_path_mtu_discovery": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.no_ipip_path_mtu_discovery"},
						},
						"ipv6_duplicate_addr_detection_transmits": {
							Type:         schema.TypeInt,
							Default:      -1,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 20),
						},
						"ipv6_path_mtu_discovery": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.no_ipv6_path_mtu_discovery"},
						},
						"ipv6_path_mtu_discovery_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(5, 71582788),
						},
						"ipv6_reject_zero_hop_limit": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.no_ipv6_reject_zero_hop_limit"},
						},
						"no_gre_path_mtu_discovery": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.gre_path_mtu_discovery"},
						},
						"no_ipip_path_mtu_discovery": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.ipip_path_mtu_discovery"},
						},
						"no_ipv6_path_mtu_discovery": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.ipv6_path_mtu_discovery"},
						},
						"no_ipv6_reject_zero_hop_limit": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.ipv6_reject_zero_hop_limit"},
						},
						"no_path_mtu_discovery": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.path_mtu_discovery"},
						},
						"no_source_quench": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.source_quench"},
						},
						"no_tcp_reset": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"drop-all-tcp", "drop-tcp-with-syn-only"}, false),
						},
						"no_tcp_rfc1323": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"no_tcp_rfc1323_paws": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"path_mtu_discovery": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.no_path_mtu_discovery"},
						},
						"source_port_upper_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(5000, 65535),
						},
						"source_quench": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_options.0.no_source_quench"},
						},
						"tcp_drop_synfin_set": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tcp_mss": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(64, 65535),
						},
					},
				},
			},
			"login": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"announcement": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"deny_sources_address": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"idle_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 60),
						},
						"message": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"change_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"character-sets", "set-transitions"}, false),
									},
									"format": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"sha1", "sha256", "sha512"}, false),
									},
									"maximum_length": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(20, 128),
									},
									"minimum_changes": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 128),
									},
									"minimum_character_changes": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(4, 15),
									},
									"minimum_length": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(6, 20),
									},
									"minimum_lower_cases": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 128),
									},
									"minimum_numerics": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 128),
									},
									"minimum_punctuations": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 128),
									},
									"minimum_reuse": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 20),
									},
									"minimum_upper_cases": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 128),
									},
								},
							},
						},
						"retry_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backoff_factor": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(5, 10),
									},
									"backoff_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 3),
									},
									"lockout_period": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 43200),
									},
									"maximum_time": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(20, 300),
									},
									"minimum_time": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(20, 60),
									},
									"tries_before_disconnect": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(2, 10),
									},
								},
							},
						},
					},
				},
			},
			"max_configuration_rollbacks": {
				Type:         schema.TypeInt,
				Default:      -1,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 49),
			},
			"max_configurations_on_flash": {
				Type:         schema.TypeInt,
				Default:      -1,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 49),
			},
			"name_server": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"no_multicast_echo": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_ping_record_route": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_ping_time_stamp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_redirects": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_redirects_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssh": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication_order": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"ciphers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"client_alive_count_max": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 255),
									},
									"client_alive_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 65535),
									},
									"connection_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 250),
									},
									"fingerprint_hash": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"md5", "sha2-256"}, false),
									},
									"hostkey_algorithm": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"key_exchange": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"log_key_changes": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"macs": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"max_pre_authentication_packets": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(20, 2147483647),
									},
									"max_sessions_per_connection": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"no_passwords": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"services.0.ssh.0.no_public_keys"},
									},
									"no_public_keys": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"services.0.ssh.0.no_passwords"},
									},
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"protocol_version": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"rate_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 250),
									},
									"root_login": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"allow", "deny", "deny-password"}, false),
									},
									"no_tcp_forwarding": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"tcp_forwarding": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"syslog": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"archive": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"binary_data": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"no_binary_data": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"files": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 1000),
									},
									"size": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(65536, 1073741824),
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
						"log_rotate_frequency": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 59),
						},
						"source_address": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateAddress(),
						},
					},
				},
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tracing_dest_override_syslog_host": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateAddress(),
			},
		},
	}
}

func resourceSystemCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)

	if err := setSystem(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("create resource junos_system", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("system")

	return append(diagWarns, resourceSystemReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)

	return resourceSystemReadWJnprSess(d, m, jnprSess)
}
func resourceSystemReadWJnprSess(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) diag.Diagnostics {
	mutex.Lock()
	systemOptions, err := readSystem(m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSystem(d, systemOptions)

	return nil
}
func resourceSystemUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer sess.closeSession(jnprSess)
	sess.configLock(jnprSess)
	if err := delSystem(m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	if err := setSystem(d, m, jnprSess); err != nil {
		sess.configClear(jnprSess)

		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	warns, err := sess.commitConf("update resource junos_system", jnprSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		sess.configClear(jnprSess)

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemReadWJnprSess(d, m, jnprSess)...)
}
func resourceSystemDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
func resourceSystemImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	systemOptions, err := readSystem(m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillSystem(d, systemOptions)
	d.SetId("system")
	result[0] = d

	return result, nil
}

func setSystem(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)

	setPrefix := "set system "
	configSet := make([]string, 0)

	for _, v := range d.Get("authentication_order").([]interface{}) {
		configSet = append(configSet, setPrefix+"authentication-order "+v.(string))
	}
	if d.Get("auto_snapshot").(bool) {
		configSet = append(configSet, setPrefix+"auto-snapshot")
	}
	if d.Get("default_address_selection").(bool) {
		configSet = append(configSet, setPrefix+"default-address-selection")
	}
	if d.Get("domain_name").(string) != "" {
		configSet = append(configSet, setPrefix+"domain-name "+d.Get("domain_name").(string))
	}
	if d.Get("host_name").(string) != "" {
		configSet = append(configSet, setPrefix+"host-name "+d.Get("host_name").(string))
	}
	for _, v := range d.Get("inet6_backup_router").([]interface{}) {
		inet6BackupRouter := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"inet6-backup-router "+inet6BackupRouter["address"].(string))
		for _, dest := range inet6BackupRouter["destination"].([]interface{}) {
			configSet = append(configSet, setPrefix+"inet6-backup-router destination "+dest.(string))
		}
	}
	if err := setSystemInternetOptions(d, m, jnprSess); err != nil {
		return err
	}
	if err := setSystemLogin(d, m, jnprSess); err != nil {
		return err
	}
	if d.Get("max_configuration_rollbacks").(int) != -1 {
		configSet = append(configSet, setPrefix+
			"max-configuration-rollbacks "+strconv.Itoa(d.Get("max_configuration_rollbacks").(int)))
	}
	if d.Get("max_configurations_on_flash").(int) != -1 {
		configSet = append(configSet, setPrefix+
			"max-configurations-on-flash "+strconv.Itoa(d.Get("max_configurations_on_flash").(int)))
	}
	for _, nameServer := range d.Get("name_server").([]interface{}) {
		configSet = append(configSet, setPrefix+"name-server "+nameServer.(string))
	}
	if d.Get("no_multicast_echo").(bool) {
		configSet = append(configSet, setPrefix+"no-multicast-echo")
	}
	if d.Get("no_ping_record_route").(bool) {
		configSet = append(configSet, setPrefix+"no-ping-record-route")
	}
	if d.Get("no_ping_time_stamp").(bool) {
		configSet = append(configSet, setPrefix+"no-ping-time-stamp")
	}
	if d.Get("no_redirects").(bool) {
		configSet = append(configSet, setPrefix+"no-redirects")
	}
	if d.Get("no_redirects_ipv6").(bool) {
		configSet = append(configSet, setPrefix+"no-redirects-ipv6")
	}
	if err := setSystemServices(d, m, jnprSess); err != nil {
		return err
	}
	for _, syslog := range d.Get("syslog").([]interface{}) {
		if syslog == nil {
			return fmt.Errorf("syslog block is empty")
		}
		syslogM := syslog.(map[string]interface{})
		for _, archive := range syslogM["archive"].([]interface{}) {
			configSet = append(configSet, setPrefix+"syslog archive")
			if archive != nil {
				archiveM := archive.(map[string]interface{})
				if archiveM["binary_data"].(bool) && archiveM["no_binary_data"].(bool) {
					return fmt.Errorf("conflict between 'binary_data' and 'no_binary_data' for syslog archive")
				}
				if archiveM["binary_data"].(bool) {
					configSet = append(configSet, setPrefix+"syslog archive binary-data")
				}
				if archiveM["no_binary_data"].(bool) {
					configSet = append(configSet, setPrefix+"syslog archive no-binary-data")
				}
				if archiveM["files"].(int) > 0 {
					configSet = append(configSet, setPrefix+"syslog archive files "+strconv.Itoa(archiveM["files"].(int)))
				}
				if archiveM["size"].(int) > 0 {
					configSet = append(configSet, setPrefix+"syslog archive size "+strconv.Itoa(archiveM["size"].(int)))
				}
				if archiveM["no_world_readable"].(bool) && archiveM["world_readable"].(bool) {
					return fmt.Errorf("conflict between 'world_readable' and 'no_world_readable' for syslog archive")
				}
				if archiveM["no_world_readable"].(bool) {
					configSet = append(configSet, setPrefix+"syslog archive no-world-readable")
				}
				if archiveM["world_readable"].(bool) {
					configSet = append(configSet, setPrefix+"syslog archive world-readable")
				}
			}
		}
		if syslogM["log_rotate_frequency"].(int) > 0 {
			configSet = append(configSet, setPrefix+"syslog log-rotate-frequency "+
				strconv.Itoa(syslogM["log_rotate_frequency"].(int)))
		}
		if syslogM["source_address"].(string) != "" {
			configSet = append(configSet, setPrefix+"syslog source-address "+syslogM["source_address"].(string))
		}
	}
	if d.Get("time_zone").(string) != "" {
		configSet = append(configSet, setPrefix+"time-zone "+d.Get("time_zone").(string))
	}
	if d.Get("tracing_dest_override_syslog_host").(string) != "" {
		configSet = append(configSet, setPrefix+"tracing destination-override syslog host "+
			d.Get("tracing_dest_override_syslog_host").(string))
	}

	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func setSystemServices(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	setPrefix := "set system services "
	configSet := make([]string, 0)

	for _, services := range d.Get("services").([]interface{}) {
		if services == nil {
			return fmt.Errorf("services block is empty")
		}
		servicesM := services.(map[string]interface{})
		for _, servicesSSH := range servicesM["ssh"].([]interface{}) {
			if servicesSSH == nil {
				return fmt.Errorf("services.0.ssh block is empty")
			}
			servicesSSHM := servicesSSH.(map[string]interface{})
			for _, auth := range servicesSSHM["authentication_order"].([]interface{}) {
				configSet = append(configSet, setPrefix+"ssh authentication-order "+auth.(string))
			}
			for _, ciphers := range servicesSSHM["ciphers"].([]interface{}) {
				configSet = append(configSet, setPrefix+"ssh ciphers "+ciphers.(string))
			}
			if servicesSSHM["client_alive_count_max"].(int) > -1 {
				configSet = append(configSet, setPrefix+"ssh client-alive-count-max "+
					strconv.Itoa(servicesSSHM["client_alive_count_max"].(int)))
			}
			if servicesSSHM["client_alive_interval"].(int) > -1 {
				configSet = append(configSet, setPrefix+"ssh client-alive-interval "+
					strconv.Itoa(servicesSSHM["client_alive_interval"].(int)))
			}
			if servicesSSHM["connection_limit"].(int) > 0 {
				configSet = append(configSet, setPrefix+"ssh connection-limit "+
					strconv.Itoa(servicesSSHM["connection_limit"].(int)))
			}
			if servicesSSHM["fingerprint_hash"].(string) != "" {
				configSet = append(configSet, setPrefix+"ssh fingerprint-hash "+
					servicesSSHM["fingerprint_hash"].(string))
			}
			for _, hostkeyAlgo := range servicesSSHM["hostkey_algorithm"].([]interface{}) {
				configSet = append(configSet, setPrefix+"ssh hostkey-algorithm "+hostkeyAlgo.(string))
			}
			for _, keyExchange := range servicesSSHM["key_exchange"].([]interface{}) {
				configSet = append(configSet, setPrefix+"ssh key-exchange "+keyExchange.(string))
			}
			if servicesSSHM["log_key_changes"].(bool) {
				configSet = append(configSet, setPrefix+"ssh log-key-changes")
			}
			for _, macs := range servicesSSHM["macs"].([]interface{}) {
				configSet = append(configSet, setPrefix+"ssh macs "+macs.(string))
			}
			if servicesSSHM["max_pre_authentication_packets"].(int) > 0 {
				configSet = append(configSet, setPrefix+"ssh max-pre-authentication-packets "+
					strconv.Itoa(servicesSSHM["max_pre_authentication_packets"].(int)))
			}
			if servicesSSHM["max_sessions_per_connection"].(int) > 0 {
				configSet = append(configSet, setPrefix+"ssh max-sessions-per-connection "+
					strconv.Itoa(servicesSSHM["max_sessions_per_connection"].(int)))
			}
			if servicesSSHM["no_passwords"].(bool) {
				configSet = append(configSet, setPrefix+"ssh no-passwords")
			}
			if servicesSSHM["no_public_keys"].(bool) {
				configSet = append(configSet, setPrefix+"ssh no-public-keys")
			}
			if servicesSSHM["port"].(int) > 0 {
				configSet = append(configSet, setPrefix+"ssh port "+
					strconv.Itoa(servicesSSHM["port"].(int)))
			}
			for _, version := range servicesSSHM["protocol_version"].([]interface{}) {
				configSet = append(configSet, setPrefix+"ssh protocol-version "+version.(string))
			}
			if servicesSSHM["rate_limit"].(int) > 0 {
				configSet = append(configSet, setPrefix+"ssh rate-limit "+
					strconv.Itoa(servicesSSHM["rate_limit"].(int)))
			}
			if servicesSSHM["root_login"].(string) != "" {
				configSet = append(configSet, setPrefix+"ssh root-login "+servicesSSHM["root_login"].(string))
			}
			if servicesSSHM["no_tcp_forwarding"].(bool) && servicesSSHM["tcp_forwarding"].(bool) {
				return fmt.Errorf("conflict between 'no_tcp_forwarding' and 'tcp_forwarding' for services ssh")
			}
			if servicesSSHM["no_tcp_forwarding"].(bool) {
				configSet = append(configSet, setPrefix+"ssh no-tcp-forwarding")
			}
			if servicesSSHM["tcp_forwarding"].(bool) {
				configSet = append(configSet, setPrefix+"ssh tcp-forwarding")
			}
		}
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func setSystemInternetOptions(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	setPrefix := "set system internet-options "
	configSet := make([]string, 0)
	for _, v := range d.Get("internet_options").([]interface{}) {
		if v == nil {
			return fmt.Errorf("internet_options block is empty")
		}
		internetOptions := v.(map[string]interface{})
		if internetOptions["gre_path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"gre-path-mtu-discovery")
		}
		for _, v2 := range internetOptions["icmpv4_rate_limit"].([]interface{}) {
			if v2 == nil {
				return fmt.Errorf("internet_options.0.icmpv4_rate_limit block is empty")
			}
			icmpv4RL := v2.(map[string]interface{})
			if icmpv4RL["bucket_size"].(int) != -1 {
				configSet = append(configSet,
					setPrefix+"icmpv4-rate-limit bucket-size "+strconv.Itoa(icmpv4RL["bucket_size"].(int)))
			}
			if icmpv4RL["packet_rate"].(int) != -1 {
				configSet = append(configSet,
					setPrefix+"icmpv4-rate-limit packet-rate "+strconv.Itoa(icmpv4RL["packet_rate"].(int)))
			}
		}
		for _, v2 := range internetOptions["icmpv6_rate_limit"].([]interface{}) {
			if v2 == nil {
				return fmt.Errorf("internet_options.0.icmpv6_rate_limit block is empty")
			}
			icmpv6RL := v2.(map[string]interface{})
			if icmpv6RL["bucket_size"].(int) != -1 {
				configSet = append(configSet,
					setPrefix+"icmpv6-rate-limit bucket-size "+strconv.Itoa(icmpv6RL["bucket_size"].(int)))
			}
			if icmpv6RL["packet_rate"].(int) != -1 {
				configSet = append(configSet,
					setPrefix+"icmpv6-rate-limit packet-rate "+strconv.Itoa(icmpv6RL["packet_rate"].(int)))
			}
		}
		if internetOptions["ipip_path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"ipip-path-mtu-discovery")
		}
		if internetOptions["ipv6_duplicate_addr_detection_transmits"].(int) != -1 {
			configSet = append(configSet,
				setPrefix+"ipv6-duplicate-addr-detection-transmits "+
					strconv.Itoa(internetOptions["ipv6_duplicate_addr_detection_transmits"].(int)))
		}
		if internetOptions["ipv6_path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-path-mtu-discovery")
		}
		if internetOptions["ipv6_path_mtu_discovery_timeout"].(int) != -1 {
			configSet = append(configSet,
				setPrefix+"ipv6-path-mtu-discovery-timeout "+strconv.Itoa(internetOptions["ipv6_path_mtu_discovery_timeout"].(int)))
		}
		if internetOptions["ipv6_reject_zero_hop_limit"].(bool) {
			configSet = append(configSet, setPrefix+"ipv6-reject-zero-hop-limit")
		}
		if internetOptions["no_gre_path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"no-gre-path-mtu-discovery")
		}
		if internetOptions["no_ipip_path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"no-ipip-path-mtu-discovery")
		}
		if internetOptions["no_ipv6_path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"no-ipv6-path-mtu-discovery")
		}
		if internetOptions["no_ipv6_reject_zero_hop_limit"].(bool) {
			configSet = append(configSet, setPrefix+"no-ipv6-reject-zero-hop-limit")
		}
		if internetOptions["no_path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"no-path-mtu-discovery")
		}
		if internetOptions["no_source_quench"].(bool) {
			configSet = append(configSet, setPrefix+"no-source-quench")
		}
		if internetOptions["no_tcp_reset"].(string) != "" {
			configSet = append(configSet, setPrefix+"no-tcp-reset "+internetOptions["no_tcp_reset"].(string))
		}
		if internetOptions["no_tcp_rfc1323"].(bool) {
			configSet = append(configSet, setPrefix+"no-tcp-rfc1323")
		}
		if internetOptions["no_tcp_rfc1323_paws"].(bool) {
			configSet = append(configSet, setPrefix+"no-tcp-rfc1323-paws")
		}
		if internetOptions["path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"path-mtu-discovery")
		}
		if internetOptions["source_port_upper_limit"].(int) != 0 {
			configSet = append(configSet,
				setPrefix+"source-port upper-limit "+strconv.Itoa(internetOptions["source_port_upper_limit"].(int)))
		}
		if internetOptions["source_quench"].(bool) {
			configSet = append(configSet, setPrefix+"source-quench")
		}
		if internetOptions["tcp_drop_synfin_set"].(bool) {
			configSet = append(configSet, setPrefix+"tcp-drop-synfin-set")
		}
		if internetOptions["tcp_mss"].(int) != 0 {
			configSet = append(configSet, setPrefix+"tcp-mss "+strconv.Itoa(internetOptions["tcp_mss"].(int)))
		}
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func setSystemLogin(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	setPrefix := "set system login "
	configSet := make([]string, 0)
	for _, v := range d.Get("login").([]interface{}) {
		if v == nil {
			return fmt.Errorf("login block is empty")
		}
		login := v.(map[string]interface{})
		if login["announcement"].(string) != "" {
			configSet = append(configSet, setPrefix+"announcement \""+login["announcement"].(string)+"\"")
		}
		for _, denySrcAddress := range login["deny_sources_address"].([]interface{}) {
			configSet = append(configSet, setPrefix+"deny-sources address "+denySrcAddress.(string))
		}
		if login["idle_timeout"].(int) != 0 {
			configSet = append(configSet, setPrefix+"idle-timeout "+strconv.Itoa(login["idle_timeout"].(int)))
		}
		if login["message"].(string) != "" {
			configSet = append(configSet, setPrefix+"message \""+login["message"].(string)+"\"")
		}
		for _, v2 := range login["password"].([]interface{}) {
			if v2 == nil {
				return fmt.Errorf("login.0.password block is empty")
			}
			loginPassword := v2.(map[string]interface{})
			if loginPassword["change_type"].(string) != "" {
				configSet = append(configSet,
					setPrefix+"password change-type "+loginPassword["change_type"].(string))
			}
			if loginPassword["format"].(string) != "" {
				configSet = append(configSet,
					setPrefix+"password format "+loginPassword["format"].(string))
			}
			if loginPassword["maximum_length"].(int) != 0 {
				configSet = append(configSet,
					setPrefix+"password maximum-length "+strconv.Itoa(loginPassword["maximum_length"].(int)))
			}
			if loginPassword["minimum_changes"].(int) != 0 {
				configSet = append(configSet,
					setPrefix+"password minimum-changes "+strconv.Itoa(loginPassword["minimum_changes"].(int)))
			}
			if loginPassword["minimum_character_changes"].(int) != 0 {
				configSet = append(configSet, setPrefix+"password minimum-character-changes "+
					strconv.Itoa(loginPassword["minimum_character_changes"].(int)))
			}
			if loginPassword["minimum_length"].(int) != 0 {
				configSet = append(configSet, setPrefix+"password minimum-length "+
					strconv.Itoa(loginPassword["minimum_length"].(int)))
			}
			if loginPassword["minimum_lower_cases"].(int) != 0 {
				configSet = append(configSet, setPrefix+"password minimum-lower-cases "+
					strconv.Itoa(loginPassword["minimum_lower_cases"].(int)))
			}
			if loginPassword["minimum_numerics"].(int) != 0 {
				configSet = append(configSet, setPrefix+"password minimum-numerics "+
					strconv.Itoa(loginPassword["minimum_numerics"].(int)))
			}
			if loginPassword["minimum_punctuations"].(int) != 0 {
				configSet = append(configSet, setPrefix+"password minimum-punctuations "+
					strconv.Itoa(loginPassword["minimum_punctuations"].(int)))
			}
			if loginPassword["minimum_reuse"].(int) != 0 {
				configSet = append(configSet, setPrefix+"password minimum-reuse "+
					strconv.Itoa(loginPassword["minimum_reuse"].(int)))
			}
			if loginPassword["minimum_upper_cases"].(int) != 0 {
				configSet = append(configSet, setPrefix+"password minimum-upper-cases "+
					strconv.Itoa(loginPassword["minimum_upper_cases"].(int)))
			}
		}
		for _, v2 := range login["retry_options"].([]interface{}) {
			if v2 == nil {
				return fmt.Errorf("login.0.retry_options block is empty")
			}
			loginRetryOptions := v2.(map[string]interface{})
			if loginRetryOptions["backoff_factor"].(int) != 0 {
				configSet = append(configSet, setPrefix+"retry-options backoff-factor "+
					strconv.Itoa(loginRetryOptions["backoff_factor"].(int)))
			}
			if loginRetryOptions["backoff_threshold"].(int) != 0 {
				configSet = append(configSet, setPrefix+"retry-options backoff-threshold "+
					strconv.Itoa(loginRetryOptions["backoff_threshold"].(int)))
			}
			if loginRetryOptions["lockout_period"].(int) != 0 {
				configSet = append(configSet, setPrefix+"retry-options lockout-period "+
					strconv.Itoa(loginRetryOptions["lockout_period"].(int)))
			}
			if loginRetryOptions["maximum_time"].(int) != 0 {
				configSet = append(configSet, setPrefix+"retry-options maximum-time "+
					strconv.Itoa(loginRetryOptions["maximum_time"].(int)))
			}
			if loginRetryOptions["minimum_time"].(int) != 0 {
				configSet = append(configSet, setPrefix+"retry-options minimum-time "+
					strconv.Itoa(loginRetryOptions["minimum_time"].(int)))
			}
			if loginRetryOptions["tries_before_disconnect"].(int) != 0 {
				configSet = append(configSet, setPrefix+"retry-options tries-before-disconnect "+
					strconv.Itoa(loginRetryOptions["tries_before_disconnect"].(int)))
			}
		}
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}

func listLinesLogin() []string {
	return []string{
		"login announcement",
		"login deny-sources",
		"login idle-timeout",
		"login message",
		"login password",
		"login retry-options",
	}
}
func listLinesServices() []string {
	ls := make([]string, 0)
	ls = append(ls, listLinesServicesSSH()...)

	return ls
}
func listLinesServicesSSH() []string {
	return []string{
		"services ssh authentication-order",
		"services ssh ciphers",
		"services ssh client-alive-count-max",
		"services ssh client-alive-interval",
		"services ssh connection-limit",
		"services ssh fingerprint-hash",
		"services ssh hostkey-algorithm",
		"services ssh key-exchange",
		"services ssh log-key-changes",
		"services ssh macs",
		"services ssh max-pre-authentication-packets",
		"services ssh max-sessions-per-connection",
		"services ssh no-passwords",
		"services ssh no-public-keys",
		"services ssh port",
		"services ssh protocol-version",
		"services ssh rate-limit",
		"services ssh root-login",
		"services ssh no-tcp-forwarding",
		"services ssh tcp-forwarding",
	}
}
func listLinesSyslog() []string {
	return []string{
		"syslog archive",
		"syslog log-rotate-frequency",
		"syslog source-address",
	}
}
func delSystem(m interface{}, jnprSess *NetconfObject) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, "authentication-order")
	listLinesToDelete = append(listLinesToDelete, "auto-snapshot")
	listLinesToDelete = append(listLinesToDelete, "default-address-selection")
	listLinesToDelete = append(listLinesToDelete, "domain-name")
	listLinesToDelete = append(listLinesToDelete, "host-name")
	listLinesToDelete = append(listLinesToDelete, "inet6-backup-router")
	listLinesToDelete = append(listLinesToDelete, "internet-options")
	listLinesToDelete = append(listLinesToDelete, listLinesLogin()...)
	listLinesToDelete = append(listLinesToDelete, "max-configuration-rollbacks")
	listLinesToDelete = append(listLinesToDelete, "max-configurations-on-flash")
	listLinesToDelete = append(listLinesToDelete, "name-server")
	listLinesToDelete = append(listLinesToDelete, "no-multicast-echo")
	listLinesToDelete = append(listLinesToDelete, "no-ping-record-route")
	listLinesToDelete = append(listLinesToDelete, "no-ping-time-stamp")
	listLinesToDelete = append(listLinesToDelete, "no-redirects")
	listLinesToDelete = append(listLinesToDelete, "no-redirects-ipv6")
	listLinesToDelete = append(listLinesToDelete, listLinesServices()...)
	listLinesToDelete = append(listLinesToDelete, listLinesSyslog()...)
	listLinesToDelete = append(listLinesToDelete, "time-zone")
	listLinesToDelete = append(listLinesToDelete,
		"tracing destination-override syslog host",
	)
	sess := m.(*Session)
	configSet := make([]string, 0)
	delPrefix := "delete system "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}
	if err := sess.configSet(configSet, jnprSess); err != nil {
		return err
	}

	return nil
}
func readSystem(m interface{}, jnprSess *NetconfObject) (systemOptions, error) {
	sess := m.(*Session)
	var confRead systemOptions
	// default -1
	confRead.maxConfigurationRollbacks = -1
	confRead.maxConfigurationsOnFlash = -1

	systemConfig, err := sess.command("show configuration system"+
		" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if systemConfig != emptyWord {
		for _, item := range strings.Split(systemConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "authentication-order "):
				confRead.authenticationOrder = append(confRead.authenticationOrder,
					strings.TrimPrefix(itemTrim, "authentication-order "))
			case itemTrim == "auto-snapshot":
				confRead.autoSnapshot = true
			case itemTrim == "default-address-selection":
				confRead.defaultAddressSelection = true
			case strings.HasPrefix(itemTrim, "domain-name "):
				confRead.domainName = strings.TrimPrefix(itemTrim, "domain-name ")
			case strings.HasPrefix(itemTrim, "host-name "):
				confRead.hostName = strings.TrimPrefix(itemTrim, "host-name ")
			case strings.HasPrefix(itemTrim, "inet6-backup-router "):
				if len(confRead.inet6BackupRouter) == 0 {
					confRead.inet6BackupRouter = append(confRead.inet6BackupRouter, map[string]interface{}{
						"address":     "",
						"destination": make([]string, 0),
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "inet6-backup-router destination "):
					confRead.inet6BackupRouter[0]["destination"] = append(confRead.inet6BackupRouter[0]["destination"].([]string),
						strings.TrimPrefix(itemTrim, "inet6-backup-router destination "))
				default:
					confRead.inet6BackupRouter[0]["address"] = strings.TrimPrefix(itemTrim, "inet6-backup-router ")
				}
			case strings.HasPrefix(itemTrim, "internet-options "):
				if err := readSystemInternetOptions(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case checkStringHasPrefixInList(itemTrim, listLinesLogin()):
				if err := readSystemLogin(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "max-configuration-rollbacks "):
				var err error
				confRead.maxConfigurationRollbacks, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-configuration-rollbacks "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "max-configurations-on-flash "):
				var err error
				confRead.maxConfigurationsOnFlash, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-configurations-on-flash "))
				if err != nil {
					return confRead, fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "name-server "):
				confRead.nameServer = append(confRead.nameServer, strings.TrimPrefix(itemTrim, "name-server "))
			case itemTrim == "no-multicast-echo":
				confRead.noMulticastEcho = true
			case itemTrim == "no-ping-record-route":
				confRead.noPingRecordRoute = true
			case itemTrim == "no-ping-time-stamp":
				confRead.noPingTimeStamp = true
			case itemTrim == "no-redirects":
				confRead.noRedirects = true
			case itemTrim == "no-redirects-ipv6":
				confRead.noRedirectsIPv6 = true
			case checkStringHasPrefixInList(itemTrim, listLinesServices()):
				if len(confRead.services) == 0 {
					confRead.services = append(confRead.services, map[string]interface{}{
						"ssh": make([]map[string]interface{}, 0),
					})
				}
				if checkStringHasPrefixInList(itemTrim, listLinesServicesSSH()) {
					if err := readSystemServicesSSH(&confRead, itemTrim); err != nil {
						return confRead, err
					}
				}
			case checkStringHasPrefixInList(itemTrim, listLinesSyslog()):
				if err := readSystemSyslog(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "time-zone "):
				confRead.timeZone = strings.TrimPrefix(itemTrim, "time-zone ")
			case strings.HasPrefix(itemTrim, "tracing destination-override syslog host "):
				confRead.tracingDestinationOverrideSyslogHost = strings.TrimPrefix(itemTrim,
					"tracing destination-override syslog host ")
			}
		}
	}

	return confRead, nil
}

func readSystemLogin(confRead *systemOptions, itemTrim string) error {
	if len(confRead.login) == 0 {
		confRead.login = append(confRead.login, map[string]interface{}{
			"announcement":         "",
			"deny_sources_address": make([]string, 0),
			"idle_timeout":         0,
			"message":              "",
			"password":             make([]map[string]interface{}, 0),
			"retry_options":        make([]map[string]interface{}, 0),
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "login announcement "):
		confRead.login[0]["announcement"] = strings.Trim(strings.TrimPrefix(itemTrim, "login announcement "), "\"")
	case strings.HasPrefix(itemTrim, "login deny-sources address "):
		confRead.login[0]["deny_sources_address"] = append(confRead.login[0]["deny_sources_address"].([]string),
			strings.TrimPrefix(itemTrim, "login deny-sources address "))
	case strings.HasPrefix(itemTrim, "login idle-timeout "):
		var err error
		confRead.login[0]["idle_timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "login idle-timeout "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "login message "):
		confRead.login[0]["message"] = strings.Trim(strings.TrimPrefix(itemTrim, "login message "), "\"")
	case strings.HasPrefix(itemTrim, "login password "):
		if len(confRead.login[0]["password"].([]map[string]interface{})) == 0 {
			confRead.login[0]["password"] = append(confRead.login[0]["password"].([]map[string]interface{}),
				map[string]interface{}{
					"change_type":               "",
					"format":                    "",
					"maximum_length":            0,
					"minimum_changes":           0,
					"minimum_character_changes": 0,
					"minimum_length":            0,
					"minimum_lower_cases":       0,
					"minimum_numerics":          0,
					"minimum_punctuations":      0,
					"minimum_reuse":             0,
					"minimum_upper_cases":       0,
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "login password change-type "):
			confRead.login[0]["password"].([]map[string]interface{})[0]["change_type"] =
				strings.TrimPrefix(itemTrim, "login password change-type ")
		case strings.HasPrefix(itemTrim, "login password format "):
			confRead.login[0]["password"].([]map[string]interface{})[0]["format"] =
				strings.TrimPrefix(itemTrim, "login password format ")
		case strings.HasPrefix(itemTrim, "login password maximum-length "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["maximum_length"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password maximum-length "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-changes "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["minimum_changes"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password minimum-changes "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-character-changes "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["minimum_character_changes"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password minimum-character-changes "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-length "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["minimum_length"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password minimum-length "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-lower-cases "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["minimum_lower_cases"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password minimum-lower-cases "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-numerics "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["minimum_numerics"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password minimum-numerics "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-punctuations "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["minimum_punctuations"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password minimum-punctuations "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-reuse "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["minimum_reuse"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password minimum-reuse "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-upper-cases "):
			var err error
			confRead.login[0]["password"].([]map[string]interface{})[0]["minimum_upper_cases"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login password minimum-upper-cases "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "login retry-options "):
		if len(confRead.login[0]["retry_options"].([]map[string]interface{})) == 0 {
			confRead.login[0]["retry_options"] = append(confRead.login[0]["retry_options"].([]map[string]interface{}),
				map[string]interface{}{
					"backoff_factor":          0,
					"backoff_threshold":       0,
					"lockout_period":          0,
					"maximum_time":            0,
					"minimum_time":            0,
					"tries_before_disconnect": 0,
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "login retry-options backoff-factor "):
			var err error
			confRead.login[0]["retry_options"].([]map[string]interface{})[0]["backoff_factor"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login retry-options backoff-factor "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options backoff-threshold "):
			var err error
			confRead.login[0]["retry_options"].([]map[string]interface{})[0]["backoff_threshold"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login retry-options backoff-threshold "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options lockout-period "):
			var err error
			confRead.login[0]["retry_options"].([]map[string]interface{})[0]["lockout_period"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login retry-options lockout-period "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options maximum-time "):
			var err error
			confRead.login[0]["retry_options"].([]map[string]interface{})[0]["maximum_time"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login retry-options maximum-time "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options minimum-time "):
			var err error
			confRead.login[0]["retry_options"].([]map[string]interface{})[0]["minimum_time"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login retry-options minimum-time "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options tries-before-disconnect "):
			var err error
			confRead.login[0]["retry_options"].([]map[string]interface{})[0]["tries_before_disconnect"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "login retry-options tries-before-disconnect "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	}

	return nil
}

func readSystemInternetOptions(confRead *systemOptions, itemTrim string) error {
	if len(confRead.internetOptions) == 0 {
		confRead.internetOptions = append(confRead.internetOptions, map[string]interface{}{
			"gre_path_mtu_discovery":                  false,
			"icmpv4_rate_limit":                       make([]map[string]interface{}, 0),
			"icmpv6_rate_limit":                       make([]map[string]interface{}, 0),
			"ipip_path_mtu_discovery":                 false,
			"ipv6_duplicate_addr_detection_transmits": -1,
			"ipv6_path_mtu_discovery":                 false,
			"ipv6_path_mtu_discovery_timeout":         0,
			"ipv6_reject_zero_hop_limit":              false,
			"no_gre_path_mtu_discovery":               false,
			"no_ipip_path_mtu_discovery":              false,
			"no_ipv6_path_mtu_discovery":              false,
			"no_ipv6_reject_zero_hop_limit":           false,
			"no_path_mtu_discovery":                   false,
			"no_source_quench":                        false,
			"no_tcp_reset":                            "",
			"no_tcp_rfc1323":                          false,
			"no_tcp_rfc1323_paws":                     false,
			"path_mtu_discovery":                      false,
			"source_port_upper_limit":                 0,
			"source_quench":                           false,
			"tcp_drop_synfin_set":                     false,
			"tcp_mss":                                 0,
		})
	}
	switch {
	case itemTrim == "internet-options gre-path-mtu-discovery":
		confRead.internetOptions[0]["gre_path_mtu_discovery"] = true
	case strings.HasPrefix(itemTrim, "internet-options icmpv4-rate-limit"):
		if len(confRead.internetOptions[0]["icmpv4_rate_limit"].([]map[string]interface{})) == 0 {
			confRead.internetOptions[0]["icmpv4_rate_limit"] = append(
				confRead.internetOptions[0]["icmpv4_rate_limit"].([]map[string]interface{}), map[string]interface{}{
					"bucket_size": -1,
					"packet_rate": -1,
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "internet-options icmpv4-rate-limit bucket-size "):
			var err error
			confRead.internetOptions[0]["icmpv4_rate_limit"].([]map[string]interface{})[0]["bucket_size"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options icmpv4-rate-limit bucket-size "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "internet-options icmpv4-rate-limit packet-rate "):
			var err error
			confRead.internetOptions[0]["icmpv4_rate_limit"].([]map[string]interface{})[0]["packet_rate"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options icmpv4-rate-limit packet-rate "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case strings.HasPrefix(itemTrim, "internet-options icmpv6-rate-limit"):
		if len(confRead.internetOptions[0]["icmpv6_rate_limit"].([]map[string]interface{})) == 0 {
			confRead.internetOptions[0]["icmpv6_rate_limit"] = append(
				confRead.internetOptions[0]["icmpv6_rate_limit"].([]map[string]interface{}), map[string]interface{}{
					"bucket_size": -1,
					"packet_rate": -1,
				})
		}
		switch {
		case strings.HasPrefix(itemTrim, "internet-options icmpv6-rate-limit bucket-size "):
			var err error
			confRead.internetOptions[0]["icmpv6_rate_limit"].([]map[string]interface{})[0]["bucket_size"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options icmpv6-rate-limit bucket-size "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "internet-options icmpv6-rate-limit packet-rate "):
			var err error
			confRead.internetOptions[0]["icmpv6_rate_limit"].([]map[string]interface{})[0]["packet_rate"], err =
				strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options icmpv6-rate-limit packet-rate "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		}
	case itemTrim == "internet-options ipip-path-mtu-discovery":
		confRead.internetOptions[0]["ipip_path_mtu_discovery"] = true
	case strings.HasPrefix(itemTrim, "internet-options ipv6-duplicate-addr-detection-transmits "):
		var err error
		confRead.internetOptions[0]["ipv6_duplicate_addr_detection_transmits"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options ipv6-duplicate-addr-detection-transmits "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case itemTrim == "internet-options ipv6-path-mtu-discovery":
		confRead.internetOptions[0]["ipv6_path_mtu_discovery"] = true
	case strings.HasPrefix(itemTrim, "internet-options ipv6-path-mtu-discovery-timeout "):
		var err error
		confRead.internetOptions[0]["ipv6_path_mtu_discovery_timeout"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options ipv6-path-mtu-discovery-timeout "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case itemTrim == "internet-options ipv6-reject-zero-hop-limit":
		confRead.internetOptions[0]["ipv6_reject_zero_hop_limit"] = true
	case itemTrim == "internet-options no-gre-path-mtu-discovery":
		confRead.internetOptions[0]["no_gre_path_mtu_discovery"] = true
	case itemTrim == "internet-options no-ipip-path-mtu-discovery":
		confRead.internetOptions[0]["no_ipip_path_mtu_discovery"] = true
	case itemTrim == "internet-options no-ipv6-path-mtu-discovery":
		confRead.internetOptions[0]["no_ipv6_path_mtu_discovery"] = true
	case itemTrim == "internet-options no-ipv6-reject-zero-hop-limit":
		confRead.internetOptions[0]["no_ipv6_reject_zero_hop_limit"] = true
	case itemTrim == "internet-options no-path-mtu-discovery":
		confRead.internetOptions[0]["no_path_mtu_discovery"] = true
	case itemTrim == "internet-options no-source-quench":
		confRead.internetOptions[0]["no_source_quench"] = true
	case strings.HasPrefix(itemTrim, "internet-options no-tcp-reset "):
		confRead.internetOptions[0]["no_tcp_reset"] = strings.TrimPrefix(itemTrim, "internet-options no-tcp-reset ")
	case itemTrim == "internet-options no-tcp-rfc1323":
		confRead.internetOptions[0]["no_tcp_rfc1323"] = true
	case itemTrim == "internet-options no-tcp-rfc1323-paws":
		confRead.internetOptions[0]["no_tcp_rfc1323_paws"] = true
	case itemTrim == "internet-options path-mtu-discovery":
		confRead.internetOptions[0]["path_mtu_discovery"] = true
	case strings.HasPrefix(itemTrim, "internet-options source-port upper-limit "):
		var err error
		confRead.internetOptions[0]["source_port_upper_limit"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options source-port upper-limit "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case itemTrim == "internet-options source-quench":
		confRead.internetOptions[0]["source_quench"] = true
	case itemTrim == "internet-options tcp-drop-synfin-set":
		confRead.internetOptions[0]["tcp_drop_synfin_set"] = true
	case strings.HasPrefix(itemTrim, "internet-options tcp-mss "):
		var err error
		confRead.internetOptions[0]["tcp_mss"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options tcp-mss "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	}

	return nil
}

func readSystemServicesSSH(confRead *systemOptions, itemTrim string) error {
	if len(confRead.services[0]["ssh"].([]map[string]interface{})) == 0 {
		confRead.services[0]["ssh"] = append(confRead.services[0]["ssh"].([]map[string]interface{}),
			map[string]interface{}{
				"authentication_order":           make([]string, 0),
				"ciphers":                        make([]string, 0),
				"client_alive_count_max":         -1,
				"client_alive_interval":          -1,
				"connection_limit":               0,
				"fingerprint_hash":               "",
				"hostkey_algorithm":              make([]string, 0),
				"key_exchange":                   make([]string, 0),
				"log_key_changes":                false,
				"macs":                           make([]string, 0),
				"max_pre_authentication_packets": 0,
				"max_sessions_per_connection":    0,
				"no_passwords":                   false,
				"no_public_keys":                 false,
				"port":                           0,
				"protocol_version":               make([]string, 0),
				"rate_limit":                     0,
				"root_login":                     "",
				"no_tcp_forwarding":              false,
				"tcp_forwarding":                 false,
			})
	}
	switch {
	case strings.HasPrefix(itemTrim, "services ssh authentication-order "):
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["authentication_order"] = append(
			confRead.services[0]["ssh"].([]map[string]interface{})[0]["authentication_order"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh authentication-order "))
	case strings.HasPrefix(itemTrim, "services ssh ciphers "):
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["ciphers"] = append(
			confRead.services[0]["ssh"].([]map[string]interface{})[0]["ciphers"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh ciphers "))
	case strings.HasPrefix(itemTrim, "services ssh client-alive-count-max "):
		var err error
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["client_alive_count_max"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh client-alive-count-max "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh client-alive-interval "):
		var err error
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["client_alive_interval"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh client-alive-interval "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh connection-limit "):
		var err error
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["connection_limit"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh connection-limit "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh fingerprint-hash "):
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["fingerprint_hash"] = strings.TrimPrefix(
			itemTrim, "services ssh fingerprint-hash ")
	case strings.HasPrefix(itemTrim, "services ssh hostkey-algorithm "):
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["hostkey_algorithm"] = append(
			confRead.services[0]["ssh"].([]map[string]interface{})[0]["hostkey_algorithm"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh hostkey-algorithm "))
	case strings.HasPrefix(itemTrim, "services ssh key-exchange "):
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["key_exchange"] = append(
			confRead.services[0]["ssh"].([]map[string]interface{})[0]["key_exchange"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh key-exchange "))
	case itemTrim == "services ssh log-key-changes":
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["log_key_changes"] = true
	case strings.HasPrefix(itemTrim, "services ssh macs "):
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["macs"] = append(
			confRead.services[0]["ssh"].([]map[string]interface{})[0]["macs"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh macs "))
	case strings.HasPrefix(itemTrim, "services ssh max-pre-authentication-packets "):
		var err error
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["max_pre_authentication_packets"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh max-pre-authentication-packets "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh max-sessions-per-connection "):
		var err error
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["max_sessions_per_connection"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh max-sessions-per-connection "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case itemTrim == "services ssh no-passwords":
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["no_passwords"] = true
	case itemTrim == "services ssh no-public-keys":
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["no_public_keys"] = true
	case strings.HasPrefix(itemTrim, "services ssh port "):
		var err error
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["port"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh port "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh protocol-version "):
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["protocol_version"] = append(
			confRead.services[0]["ssh"].([]map[string]interface{})[0]["protocol_version"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh protocol-version "))
	case strings.HasPrefix(itemTrim, "services ssh rate-limit "):
		var err error
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["rate_limit"], err =
			strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh rate-limit "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh root-login "):
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["root_login"] =
			strings.TrimPrefix(itemTrim, "services ssh root-login ")
	case itemTrim == "services ssh no-tcp-forwarding":
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["no_tcp_forwarding"] = true
	case itemTrim == "services ssh tcp-forwarding":
		confRead.services[0]["ssh"].([]map[string]interface{})[0]["tcp_forwarding"] = true
	}

	return nil
}

func readSystemSyslog(confRead *systemOptions, itemTrim string) error {
	if len(confRead.syslog) == 0 {
		confRead.syslog = append(confRead.syslog, map[string]interface{}{
			"archive":              make([]map[string]interface{}, 0),
			"log_rotate_frequency": 0,
			"source_address":       "",
		})
	}
	switch {
	case strings.HasPrefix(itemTrim, "syslog archive"):
		if len(confRead.syslog[0]["archive"].([]map[string]interface{})) == 0 {
			confRead.syslog[0]["archive"] = append(confRead.syslog[0]["archive"].([]map[string]interface{}),
				map[string]interface{}{
					"binary_data":       false,
					"no_binary_data":    false,
					"files":             0,
					"size":              0,
					"no_world_readable": false,
					"world_readable":    false,
				})
		}
		switch {
		case itemTrim == "syslog archive binary-data":
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["binary_data"] = true
		case itemTrim == "syslog archive no-binary-data":
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["no_binary_data"] = true
		case strings.HasPrefix(itemTrim, "syslog archive files "):
			var err error
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["files"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "syslog archive files "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "syslog archive size "):
			var err error
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["size"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "syslog archive size "))
			if err != nil {
				return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
			}
		case itemTrim == "syslog archive no-world-readable":
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["no_world_readable"] = true
		case itemTrim == "syslog archive world-readable":
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["world_readable"] = true
		}
	case strings.HasPrefix(itemTrim, "syslog log-rotate-frequency "):
		var err error
		confRead.syslog[0]["log_rotate_frequency"], err = strconv.Atoi(
			strings.TrimPrefix(itemTrim, "syslog log-rotate-frequency "))
		if err != nil {
			return fmt.Errorf("failed to convert value from '%s' to integer : %w", itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "syslog source-address "):
		confRead.syslog[0]["source_address"] = strings.TrimPrefix(
			itemTrim, "syslog source-address ")
	}

	return nil
}

func fillSystem(d *schema.ResourceData, systemOptions systemOptions) {
	if tfErr := d.Set("authentication_order", systemOptions.authenticationOrder); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("auto_snapshot", systemOptions.autoSnapshot); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("default_address_selection", systemOptions.defaultAddressSelection); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("domain_name", systemOptions.domainName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("host_name", systemOptions.hostName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("inet6_backup_router", systemOptions.inet6BackupRouter); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("internet_options", systemOptions.internetOptions); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("login", systemOptions.login); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_configuration_rollbacks", systemOptions.maxConfigurationRollbacks); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_configurations_on_flash", systemOptions.maxConfigurationsOnFlash); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("name_server", systemOptions.nameServer); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_multicast_echo", systemOptions.noMulticastEcho); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_ping_record_route", systemOptions.noPingRecordRoute); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_ping_time_stamp", systemOptions.noPingTimeStamp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_redirects", systemOptions.noRedirects); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("no_redirects_ipv6", systemOptions.noRedirectsIPv6); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("services", systemOptions.services); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("syslog", systemOptions.syslog); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("time_zone", systemOptions.timeZone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("tracing_dest_override_syslog_host",
		systemOptions.tracingDestinationOverrideSyslogHost); tfErr != nil {
		panic(tfErr)
	}
}
