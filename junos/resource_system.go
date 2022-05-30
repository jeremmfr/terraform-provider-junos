package junos

import (
	"context"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	jdecode "github.com/jeremmfr/junosdecode"
)

type systemOptions struct {
	autoSnapshot                         bool
	defaultAddressSelection              bool
	noMulticastEcho                      bool
	noPingRecordRoute                    bool
	noPingTimeStamp                      bool
	noRedirects                          bool
	noRedirectsIPv6                      bool
	radiusOptionsEnhancedAccounting      bool
	radiusOptionsPasswodProtoclMsChapV2  bool
	maxConfigurationRollbacks            int
	maxConfigurationsOnFlash             int
	domainName                           string
	hostName                             string
	radiusOptionsAttributesNasIPAddress  string
	timeZone                             string
	tracingDestinationOverrideSyslogHost string
	authenticationOrder                  []string
	nameServer                           []string
	archivalConfiguration                []map[string]interface{}
	inet6BackupRouter                    []map[string]interface{}
	internetOptions                      []map[string]interface{}
	license                              []map[string]interface{}
	login                                []map[string]interface{}
	ntp                                  []map[string]interface{}
	ports                                []map[string]interface{}
	services                             []map[string]interface{}
	syslog                               []map[string]interface{}
}

func resourceSystem() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSystemCreate,
		ReadWithoutTimeout:   resourceSystemRead,
		UpdateWithoutTimeout: resourceSystemUpdate,
		DeleteWithoutTimeout: resourceSystemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemImport,
		},
		Schema: map[string]*schema.Schema{
			"archival_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"archive_site": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringDoesNotContainAny(" "),
									},
									"password": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
								},
							},
						},
						"transfer_interval": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(15, 2880),
							ConflictsWith: []string{"archival_configuration.0.transfer_on_commit"},
						},
						"transfer_on_commit": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"archival_configuration.0.transfer_interval"},
						},
					},
				},
			},
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
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
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
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"drop-all-tcp", "drop-tcp-with-syn-only"}, false),
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
			"license": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoupdate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"autoupdate_password": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							RequiredWith: []string{"license.0.autoupdate_url"},
						},
						"autoupdate_url": {
							Type:         schema.TypeString,
							Optional:     true,
							RequiredWith: []string{"license.0.autoupdate"},
						},
						"renew_before_expiration": {
							Type:         schema.TypeInt,
							Default:      -1,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 60),
							RequiredWith: []string{"license.0.renew_interval"},
						},
						"renew_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 336),
							RequiredWith: []string{"license.0.renew_before_expiration"},
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
							Type:     schema.TypeSet,
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
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"character-sets", "set-transitions"}, false),
									},
									"format": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"sha1", "sha256", "sha512"}, false),
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
			"ntp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot_server": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"broadcast_client": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"interval_range": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 3),
						},
						"multicast_client": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"multicast_client_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"threshold_action": {
							Type:         schema.TypeString,
							Optional:     true,
							RequiredWith: []string{"ntp.0.threshold_value"},
							ValidateFunc: validation.StringInSlice([]string{"accept", "reject"}, false),
						},
						"threshold_value": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"ntp.0.threshold_action"},
							ValidateFunc: validation.IntBetween(1, 600),
						},
					},
				},
			},
			"ports": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auxiliary_authentication_order": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"auxiliary_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auxiliary_insecure": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auxiliary_logout_on_disconnect": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auxiliary_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ansi", "small-xterm", "vt100", "xterm"}, false),
						},
						"console_authentication_order": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"console_disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"console_insecure": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"console_logout_on_disconnect": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"console_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ansi", "small-xterm", "vt100", "xterm"}, false),
						},
					},
				},
			},
			"radius_options_attributes_nas_ipaddress": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"radius_options_enhanced_accounting": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"radius_options_password_protocol_mschapv2": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"netconf_ssh": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"rate_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 250),
									},
								},
							},
						},
						"netconf_traceoptions": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringDoesNotContainAny("/% "),
									},
									"file_files": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(2, 1000),
									},
									"file_match": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"file_no_world_readable": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"services.0.netconf_traceoptions.0.file_world_readable"},
									},
									"file_size": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(10240, 1073741824),
									},
									"file_world_readable": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"services.0.netconf_traceoptions.0.file_no_world_readable"},
									},
									"flag": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"no_remote_trace": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"on_demand": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
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
										Type:     schema.TypeSet,
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
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"md5", "sha2-256"}, false),
									},
									"hostkey_algorithm": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"key_exchange": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"log_key_changes": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"macs": {
										Type:     schema.TypeSet,
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
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"rate_limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 250),
									},
									"root_login": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"allow", "deny", "deny-password"}, false),
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
						"web_management_http": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interface": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
								},
							},
						},
						"web_management_https": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interface": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"local_certificate": {
										Type:     schema.TypeString,
										Optional: true,
										ExactlyOneOf: []string{
											"services.0.web_management_https.0.local_certificate",
											"services.0.web_management_https.0.pki_local_certificate",
											"services.0.web_management_https.0.system_generated_certificate",
										},
									},
									"pki_local_certificate": {
										Type:     schema.TypeString,
										Optional: true,
										ExactlyOneOf: []string{
											"services.0.web_management_https.0.local_certificate",
											"services.0.web_management_https.0.pki_local_certificate",
											"services.0.web_management_https.0.system_generated_certificate",
										},
									},
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"system_generated_certificate": {
										Type:     schema.TypeBool,
										Optional: true,
										ExactlyOneOf: []string{
											"services.0.web_management_https.0.local_certificate",
											"services.0.web_management_https.0.pki_local_certificate",
											"services.0.web_management_https.0.system_generated_certificate",
										},
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
						"console": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"any_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"authorization_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"changelog_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"conflictlog_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"daemon_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"dfc_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"external_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"firewall_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"ftp_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"interactivecommands_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"kernel_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"ntp_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"pfe_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"security_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
									},
									"user_severity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listOfSyslogSeverity(), false),
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
						"time_format_millisecond": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"time_format_year": {
							Type:     schema.TypeBool,
							Optional: true,
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
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSystem(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("system")

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := setSystem(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_system", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId("system")

	return append(diagWarns, resourceSystemReadWJunSess(d, clt, junSess)...)
}

func resourceSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSystemReadWJunSess(d, clt, junSess)
}

func resourceSystemReadWJunSess(d *schema.ResourceData, clt *Client, junSess *junosSession) diag.Diagnostics {
	mutex.Lock()
	systemOptions, err := readSystem(clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillSystem(d, systemOptions)

	return nil
}

func resourceSystemUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSystem(clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystem(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)
	if err := clt.configLock(ctx, junSess); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSystem(clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystem(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_system", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemReadWJunSess(d, clt, junSess)...)
}

func resourceSystemDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSystemImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	systemOptions, err := readSystem(clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSystem(d, systemOptions)
	d.SetId("system")
	result[0] = d

	return result, nil
}

func setSystem(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := "set system "
	configSet := make([]string, 0)

	for _, v := range d.Get("archival_configuration").([]interface{}) {
		archivalConfig := v.(map[string]interface{})
		archiveSiteURLList := make([]string, 0)
		for _, v2 := range archivalConfig["archive_site"].([]interface{}) {
			archiveSite := v2.(map[string]interface{})
			if bchk.StringInSlice(archiveSite["url"].(string), archiveSiteURLList) {
				return fmt.Errorf("multiple blocks archive_site with the same url %s", archiveSite["url"].(string))
			}
			archiveSiteURLList = append(archiveSiteURLList, archiveSite["url"].(string))
			configSet = append(configSet, setPrefix+"archival configuration archive-sites \""+archiveSite["url"].(string)+"\"")
			if pass := archiveSite["password"].(string); pass != "" {
				configSet = append(configSet,
					setPrefix+"archival configuration archive-sites \""+archiveSite["url"].(string)+"\" password \""+pass+"\"")
			}
		}
		switch {
		case archivalConfig["transfer_interval"].(int) != 0:
			configSet = append(configSet, setPrefix+"archival configuration transfer-interval "+
				strconv.Itoa(archivalConfig["transfer_interval"].(int)))
		case archivalConfig["transfer_on_commit"].(bool):
			configSet = append(configSet, setPrefix+"archival configuration transfer-on-commit")
		default:
			return fmt.Errorf("transfer_interval or transfer_on_commit missing for archival_configuration")
		}
	}
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
		for _, dest := range sortSetOfString(inet6BackupRouter["destination"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+"inet6-backup-router destination "+dest)
		}
	}
	if err := setSystemInternetOptions(d, clt, junSess); err != nil {
		return err
	}
	for _, v := range d.Get("license").([]interface{}) {
		setPrefixLicense := setPrefix + "license "
		license := v.(map[string]interface{})
		if license["autoupdate"].(bool) {
			configSet = append(configSet, setPrefixLicense+"autoupdate")
			if license["autoupdate_url"].(string) != "" {
				setPrefixLicenseUpdate := setPrefixLicense + "autoupdate url \"" + license["autoupdate_url"].(string) + "\""
				if license["autoupdate_password"].(string) != "" {
					configSet = append(configSet, setPrefixLicenseUpdate+" password \""+
						license["autoupdate_password"].(string)+"\"")
				} else {
					configSet = append(configSet, setPrefixLicenseUpdate)
				}
			}
		} else if license["autoupdate_url"].(string) != "" {
			return fmt.Errorf("license.0.autoupdate need to be true")
		}
		if license["renew_before_expiration"].(int) != -1 {
			configSet = append(configSet, setPrefixLicense+"renew before-expiration "+
				strconv.Itoa(license["renew_before_expiration"].(int)))
		}
		if license["renew_interval"].(int) > 0 {
			configSet = append(configSet, setPrefixLicense+"renew interval "+
				strconv.Itoa(license["renew_interval"].(int)))
		}
		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefixLicense) {
			return fmt.Errorf("license block is empty")
		}
	}
	if err := setSystemLogin(d, clt, junSess); err != nil {
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
	for _, vi := range d.Get("ntp").([]interface{}) {
		setPrefixNtp := setPrefix + "ntp "
		ntp := vi.(map[string]interface{})
		if v := ntp["boot_server"].(string); v != "" {
			configSet = append(configSet, setPrefixNtp+"boot-server "+v)
		}
		if ntp["broadcast_client"].(bool) {
			configSet = append(configSet, setPrefixNtp+"broadcast-client")
		}
		if v := ntp["interval_range"].(int); v != -1 {
			configSet = append(configSet, setPrefixNtp+"interval-range "+strconv.Itoa(v))
		}
		if ntp["multicast_client"].(bool) {
			configSet = append(configSet, setPrefixNtp+"multicast-client")
			if v := ntp["multicast_client_address"].(string); v != "" {
				configSet = append(configSet, setPrefixNtp+"multicast-client "+v)
			}
		} else if ntp["multicast_client_address"].(string) != "" {
			return fmt.Errorf("ntp.0.multicast_client need to be true with multicast_client_address")
		}
		if v := ntp["threshold_value"].(int); v != 0 {
			configSet = append(configSet, setPrefixNtp+"threshold "+strconv.Itoa(v)+" action "+ntp["threshold_action"].(string))
		}
		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefixNtp) {
			return fmt.Errorf("ntp block is empty")
		}
	}
	for _, p := range d.Get("ports").([]interface{}) {
		if p == nil {
			return fmt.Errorf("ports block is empty")
		}
		ports := p.(map[string]interface{})
		for _, v := range ports["auxiliary_authentication_order"].([]interface{}) {
			configSet = append(configSet, setPrefix+"ports auxiliary authentication-order "+v.(string))
		}
		if ports["auxiliary_disable"].(bool) {
			configSet = append(configSet, setPrefix+"ports auxiliary disable")
		}
		if ports["auxiliary_insecure"].(bool) {
			configSet = append(configSet, setPrefix+"ports auxiliary insecure")
		}
		if ports["auxiliary_logout_on_disconnect"].(bool) {
			configSet = append(configSet, setPrefix+"ports auxiliary log-out-on-disconnect")
		}
		if v := ports["auxiliary_type"].(string); v != "" {
			configSet = append(configSet, setPrefix+"ports auxiliary type "+v)
		}
		for _, v := range ports["console_authentication_order"].([]interface{}) {
			configSet = append(configSet, setPrefix+"ports console authentication-order "+v.(string))
		}
		if ports["console_disable"].(bool) {
			configSet = append(configSet, setPrefix+"ports console disable")
		}
		if ports["console_insecure"].(bool) {
			configSet = append(configSet, setPrefix+"ports console insecure")
		}
		if ports["console_logout_on_disconnect"].(bool) {
			configSet = append(configSet, setPrefix+"ports console log-out-on-disconnect")
		}
		if v := ports["console_type"].(string); v != "" {
			configSet = append(configSet, setPrefix+"ports console type "+v)
		}
		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"ports ") {
			return fmt.Errorf("ports block is empty")
		}
	}
	if v := d.Get("radius_options_attributes_nas_ipaddress").(string); v != "" {
		configSet = append(configSet, setPrefix+"radius-options attributes nas-ip-address "+v)
	}
	if d.Get("radius_options_enhanced_accounting").(bool) {
		configSet = append(configSet, setPrefix+"radius-options enhanced-accounting")
	}
	if d.Get("radius_options_password_protocol_mschapv2").(bool) {
		configSet = append(configSet, setPrefix+"radius-options password-protocol mschap-v2")
	}
	if err := setSystemServices(d, clt, junSess); err != nil {
		return err
	}
	if err := setSystemSyslog(d, clt, junSess); err != nil {
		return err
	}
	if d.Get("time_zone").(string) != "" {
		configSet = append(configSet, setPrefix+"time-zone "+d.Get("time_zone").(string))
	}
	if d.Get("tracing_dest_override_syslog_host").(string) != "" {
		configSet = append(configSet, setPrefix+"tracing destination-override syslog host "+
			d.Get("tracing_dest_override_syslog_host").(string))
	}

	return clt.configSet(configSet, junSess)
}

func setSystemServices(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := "set system services "
	configSet := make([]string, 0)

	for _, services := range d.Get("services").([]interface{}) {
		if services == nil {
			return fmt.Errorf("services block is empty")
		}
		servicesM := services.(map[string]interface{})
		for _, servicesNetconfSSH := range servicesM["netconf_ssh"].([]interface{}) {
			netconfSSH := servicesNetconfSSH.(map[string]interface{})
			if v := netconfSSH["client_alive_count_max"].(int); v > -1 {
				configSet = append(configSet, setPrefix+"netconf ssh client-alive-count-max "+strconv.Itoa(v))
			}
			if v := netconfSSH["client_alive_interval"].(int); v > -1 {
				configSet = append(configSet, setPrefix+"netconf ssh client-alive-interval "+strconv.Itoa(v))
			}
			if v := netconfSSH["connection_limit"].(int); v > 0 {
				configSet = append(configSet, setPrefix+"netconf ssh connection-limit "+strconv.Itoa(v))
			}
			if v := netconfSSH["rate_limit"].(int); v > 0 {
				configSet = append(configSet, setPrefix+"netconf ssh rate-limit "+strconv.Itoa(v))
			}
			if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"netconf ssh ") {
				return fmt.Errorf("services.0.netconf_ssh block is empty")
			}
		}
		for _, servicesNetconfTraceOpts := range servicesM["netconf_traceoptions"].([]interface{}) {
			if servicesNetconfTraceOpts == nil {
				return fmt.Errorf("services.0.netconf_traceoptions block is empty")
			}
			netconfTraceOpts := servicesNetconfTraceOpts.(map[string]interface{})
			if v := netconfTraceOpts["file_name"].(string); v != "" {
				configSet = append(configSet, setPrefix+"netconf traceoptions file \""+v+"\"")
			}
			if v := netconfTraceOpts["file_files"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"netconf traceoptions file files "+strconv.Itoa(v))
			}
			if v := netconfTraceOpts["file_match"].(string); v != "" {
				configSet = append(configSet, setPrefix+"netconf traceoptions file match \""+v+"\"")
			}
			if v := netconfTraceOpts["file_size"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"netconf traceoptions file size "+strconv.Itoa(v))
			}
			if netconfTraceOpts["file_no_world_readable"].(bool) {
				configSet = append(configSet, setPrefix+"netconf traceoptions file no-world-readable")
			}
			if netconfTraceOpts["file_world_readable"].(bool) {
				configSet = append(configSet, setPrefix+"netconf traceoptions file world-readable")
			}
			for _, v := range sortSetOfString(netconfTraceOpts["flag"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"netconf traceoptions flag "+v)
			}
			if netconfTraceOpts["no_remote_trace"].(bool) {
				configSet = append(configSet, setPrefix+"netconf traceoptions no-remote-trace")
			}
			if netconfTraceOpts["on_demand"].(bool) {
				configSet = append(configSet, setPrefix+"netconf traceoptions on-demand")
			}
		}
		for _, servicesSSH := range servicesM["ssh"].([]interface{}) {
			servicesSSHM := servicesSSH.(map[string]interface{})
			for _, auth := range servicesSSHM["authentication_order"].([]interface{}) {
				configSet = append(configSet, setPrefix+"ssh authentication-order "+auth.(string))
			}
			for _, ciphers := range sortSetOfString(servicesSSHM["ciphers"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"ssh ciphers \""+ciphers+"\"")
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
			for _, hostkeyAlgo := range sortSetOfString(servicesSSHM["hostkey_algorithm"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"ssh hostkey-algorithm "+hostkeyAlgo)
			}
			for _, keyExchange := range sortSetOfString(servicesSSHM["key_exchange"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"ssh key-exchange "+keyExchange)
			}
			if servicesSSHM["log_key_changes"].(bool) {
				configSet = append(configSet, setPrefix+"ssh log-key-changes")
			}
			for _, macs := range sortSetOfString(servicesSSHM["macs"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"ssh macs "+macs)
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
			for _, version := range sortSetOfString(servicesSSHM["protocol_version"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"ssh protocol-version "+version)
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
			if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"ssh") {
				return fmt.Errorf("services.0.ssh block is empty")
			}
		}
		for _, http := range servicesM["web_management_http"].([]interface{}) {
			configSet = append(configSet, setPrefix+"web-management http")
			if http != nil {
				httpOptions := http.(map[string]interface{})
				for _, interf := range sortSetOfString(httpOptions["interface"].(*schema.Set).List()) {
					configSet = append(configSet, setPrefix+"web-management http interface "+interf)
				}
				if httpOptions["port"].(int) > 0 {
					configSet = append(configSet, setPrefix+"web-management http port "+
						strconv.Itoa(httpOptions["port"].(int)))
				}
			}
		}
		for _, https := range servicesM["web_management_https"].([]interface{}) {
			httpsOptions := https.(map[string]interface{})
			for _, interf := range sortSetOfString(httpsOptions["interface"].(*schema.Set).List()) {
				configSet = append(configSet, setPrefix+"web-management https interface "+interf)
			}
			if httpsOptions["local_certificate"].(string) != "" {
				configSet = append(configSet,
					setPrefix+"web-management https local-certificate \""+httpsOptions["local_certificate"].(string)+"\"")
			}
			if httpsOptions["pki_local_certificate"].(string) != "" {
				configSet = append(configSet,
					setPrefix+"web-management https pki-local-certificate \""+httpsOptions["pki_local_certificate"].(string)+"\"")
			}
			if httpsOptions["port"].(int) > 0 {
				configSet = append(configSet, setPrefix+"web-management https port "+
					strconv.Itoa(httpsOptions["port"].(int)))
			}
			if httpsOptions["system_generated_certificate"].(bool) {
				configSet = append(configSet, setPrefix+"web-management https system-generated-certificate")
			}
		}
	}

	return clt.configSet(configSet, junSess)
}

func setSystemInternetOptions(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := "set system internet-options "
	configSet := make([]string, 0)
	for _, v := range d.Get("internet_options").([]interface{}) {
		internetOptions := v.(map[string]interface{})
		if internetOptions["gre_path_mtu_discovery"].(bool) {
			configSet = append(configSet, setPrefix+"gre-path-mtu-discovery")
		}
		for _, v2 := range internetOptions["icmpv4_rate_limit"].([]interface{}) {
			icmpv4RL := v2.(map[string]interface{})
			if icmpv4RL["bucket_size"].(int) != -1 {
				configSet = append(configSet,
					setPrefix+"icmpv4-rate-limit bucket-size "+strconv.Itoa(icmpv4RL["bucket_size"].(int)))
			}
			if icmpv4RL["packet_rate"].(int) != -1 {
				configSet = append(configSet,
					setPrefix+"icmpv4-rate-limit packet-rate "+strconv.Itoa(icmpv4RL["packet_rate"].(int)))
			}
			if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"icmpv4-rate-limit") {
				return fmt.Errorf("internet_options.0.icmpv4_rate_limit block is empty")
			}
		}
		for _, v2 := range internetOptions["icmpv6_rate_limit"].([]interface{}) {
			icmpv6RL := v2.(map[string]interface{})
			if icmpv6RL["bucket_size"].(int) != -1 {
				configSet = append(configSet,
					setPrefix+"icmpv6-rate-limit bucket-size "+strconv.Itoa(icmpv6RL["bucket_size"].(int)))
			}
			if icmpv6RL["packet_rate"].(int) != -1 {
				configSet = append(configSet,
					setPrefix+"icmpv6-rate-limit packet-rate "+strconv.Itoa(icmpv6RL["packet_rate"].(int)))
			}
			if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"icmpv6-rate-limit") {
				return fmt.Errorf("internet_options.0.icmpv6_rate_limit block is empty")
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
	if len(configSet) == 0 && len(d.Get("internet_options").([]interface{})) != 0 {
		return fmt.Errorf("internet_options block is empty")
	}

	return clt.configSet(configSet, junSess)
}

func setSystemLogin(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
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
		for _, denySrcAddress := range sortSetOfString(login["deny_sources_address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefix+"deny-sources address "+denySrcAddress)
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

	return clt.configSet(configSet, junSess)
}

func setSystemSyslog(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	setPrefix := "set system syslog "
	configSet := make([]string, 0)
	for _, syslog := range d.Get("syslog").([]interface{}) {
		if syslog == nil {
			return fmt.Errorf("syslog block is empty")
		}
		syslogM := syslog.(map[string]interface{})
		for _, archive := range syslogM["archive"].([]interface{}) {
			configSet = append(configSet, setPrefix+"archive")
			if archive != nil {
				archiveM := archive.(map[string]interface{})
				if archiveM["binary_data"].(bool) && archiveM["no_binary_data"].(bool) {
					return fmt.Errorf("conflict between 'binary_data' and 'no_binary_data' for syslog archive")
				}
				if archiveM["binary_data"].(bool) {
					configSet = append(configSet, setPrefix+"archive binary-data")
				}
				if archiveM["no_binary_data"].(bool) {
					configSet = append(configSet, setPrefix+"archive no-binary-data")
				}
				if archiveM["files"].(int) > 0 {
					configSet = append(configSet, setPrefix+"archive files "+strconv.Itoa(archiveM["files"].(int)))
				}
				if archiveM["size"].(int) > 0 {
					configSet = append(configSet, setPrefix+"archive size "+strconv.Itoa(archiveM["size"].(int)))
				}
				if archiveM["no_world_readable"].(bool) && archiveM["world_readable"].(bool) {
					return fmt.Errorf("conflict between 'world_readable' and 'no_world_readable' for syslog archive")
				}
				if archiveM["no_world_readable"].(bool) {
					configSet = append(configSet, setPrefix+"archive no-world-readable")
				}
				if archiveM["world_readable"].(bool) {
					configSet = append(configSet, setPrefix+"archive world-readable")
				}
			}
		}
		for _, consoleSchema := range syslogM["console"].([]interface{}) {
			if consoleSchema == nil {
				return fmt.Errorf("syslog.0.console block is empty")
			}
			console := consoleSchema.(map[string]interface{})
			if v := console["any_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console any "+v)
			}
			if v := console["authorization_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console authorization "+v)
			}
			if v := console["changelog_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console change-log "+v)
			}
			if v := console["conflictlog_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console conflict-log "+v)
			}
			if v := console["daemon_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console daemon "+v)
			}
			if v := console["dfc_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console dfc "+v)
			}
			if v := console["external_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console external "+v)
			}
			if v := console["firewall_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console firewall "+v)
			}
			if v := console["ftp_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console ftp "+v)
			}
			if v := console["interactivecommands_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console interactive-commands "+v)
			}
			if v := console["kernel_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console kernel "+v)
			}
			if v := console["ntp_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console ntp "+v)
			}
			if v := console["pfe_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console pfe "+v)
			}
			if v := console["security_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console security "+v)
			}
			if v := console["user_severity"].(string); v != "" {
				configSet = append(configSet, setPrefix+"console user "+v)
			}
			if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"console") {
				return fmt.Errorf("syslog.0.console block is empty")
			}
		}
		if syslogM["log_rotate_frequency"].(int) > 0 {
			configSet = append(configSet, setPrefix+"log-rotate-frequency "+
				strconv.Itoa(syslogM["log_rotate_frequency"].(int)))
		}
		if syslogM["source_address"].(string) != "" {
			configSet = append(configSet, setPrefix+"source-address "+syslogM["source_address"].(string))
		}
		if syslogM["time_format_millisecond"].(bool) {
			configSet = append(configSet, setPrefix+"time-format millisecond")
		}
		if syslogM["time_format_year"].(bool) {
			configSet = append(configSet, setPrefix+"time-format year")
		}
	}

	return clt.configSet(configSet, junSess)
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

func listLinesLicense() []string {
	return []string{
		"license autoupdate",
		"license renew",
	}
}

func listLinesNtp() []string {
	return []string{
		"ntp boot-server",
		"ntp broadcast-client",
		"ntp interval-range",
		"ntp multicast-client",
		"ntp source-address",
		"ntp threshold",
	}
}

func listLinesServices() []string {
	ls := make([]string, 0)
	ls = append(ls, "services netconf traceoptions")
	ls = append(ls, listLinesServicesNetconfSSH()...)
	ls = append(ls, listLinesServicesSSH()...)
	ls = append(ls, listLinesServicesWebManagement()...)

	return ls
}

func listLinesServicesNetconfSSH() []string {
	return []string{
		"services netconf ssh client-alive-count-max",
		"services netconf ssh client-alive-interval",
		"services netconf ssh connection-limit",
		"services netconf ssh rate-limit",
	}
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

func listLinesServicesWebManagement() []string {
	return []string{
		"services web-management http",
		"services web-management https",
	}
}

func listLinesSyslog() []string {
	return []string{
		"syslog archive",
		"syslog console ",
		"syslog log-rotate-frequency",
		"syslog source-address",
		"syslog time-format ",
	}
}

func delSystem(clt *Client, junSess *junosSession) error {
	listLinesToDelete := make([]string, 0)
	listLinesToDelete = append(listLinesToDelete, "archival configuration")
	listLinesToDelete = append(listLinesToDelete, "authentication-order")
	listLinesToDelete = append(listLinesToDelete, "auto-snapshot")
	listLinesToDelete = append(listLinesToDelete, "default-address-selection")
	listLinesToDelete = append(listLinesToDelete, "domain-name")
	listLinesToDelete = append(listLinesToDelete, "host-name")
	listLinesToDelete = append(listLinesToDelete, "inet6-backup-router")
	listLinesToDelete = append(listLinesToDelete, "internet-options")
	listLinesToDelete = append(listLinesToDelete, listLinesLicense()...)
	listLinesToDelete = append(listLinesToDelete, listLinesLogin()...)
	listLinesToDelete = append(listLinesToDelete, "max-configuration-rollbacks")
	listLinesToDelete = append(listLinesToDelete, "max-configurations-on-flash")
	listLinesToDelete = append(listLinesToDelete, listLinesNtp()...)
	listLinesToDelete = append(listLinesToDelete, "name-server")
	listLinesToDelete = append(listLinesToDelete, "no-multicast-echo")
	listLinesToDelete = append(listLinesToDelete, "no-ping-record-route")
	listLinesToDelete = append(listLinesToDelete, "no-ping-time-stamp")
	listLinesToDelete = append(listLinesToDelete, "no-redirects")
	listLinesToDelete = append(listLinesToDelete, "no-redirects-ipv6")
	listLinesToDelete = append(listLinesToDelete, "ports")
	listLinesToDelete = append(listLinesToDelete, "radius-options")
	listLinesToDelete = append(listLinesToDelete, listLinesServices()...)
	listLinesToDelete = append(listLinesToDelete, listLinesSyslog()...)
	listLinesToDelete = append(listLinesToDelete, "time-zone")
	listLinesToDelete = append(listLinesToDelete,
		"tracing destination-override syslog host",
	)

	configSet := make([]string, 0)
	delPrefix := "delete system "
	for _, line := range listLinesToDelete {
		configSet = append(configSet,
			delPrefix+line)
	}

	return clt.configSet(configSet, junSess)
}

func readSystem(clt *Client, junSess *junosSession) (systemOptions, error) {
	var confRead systemOptions
	// default -1
	confRead.maxConfigurationRollbacks = -1
	confRead.maxConfigurationsOnFlash = -1

	showConfig, err := clt.command(cmdShowConfig+"system"+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	if showConfig != emptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "archival configuration "):
				if len(confRead.archivalConfiguration) == 0 {
					confRead.archivalConfiguration = append(confRead.archivalConfiguration, map[string]interface{}{
						"archive_site":       make([]map[string]interface{}, 0),
						"transfer_interval":  0,
						"transfer_on_commit": false,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "archival configuration archive-sites "):
					archiveSiteSplit := strings.Split(strings.TrimPrefix(itemTrim, "archival configuration archive-sites "), " ")
					if len(archiveSiteSplit) == 1 {
						confRead.archivalConfiguration[0]["archive_site"] = append(
							confRead.archivalConfiguration[0]["archive_site"].([]map[string]interface{}), map[string]interface{}{
								"url":      strings.Trim(archiveSiteSplit[0], "\""),
								"password": "",
							})
					} else {
						passWord, err := jdecode.Decode(strings.Trim(archiveSiteSplit[2], "\""))
						if err != nil {
							return confRead, fmt.Errorf("failed to decode archive-site password: %w", err)
						}
						confRead.archivalConfiguration[0]["archive_site"] = append(
							confRead.archivalConfiguration[0]["archive_site"].([]map[string]interface{}), map[string]interface{}{
								"url":      strings.Trim(archiveSiteSplit[0], "\""),
								"password": passWord,
							})
					}
				case strings.HasPrefix(itemTrim, "archival configuration transfer-interval "):
					var err error
					confRead.archivalConfiguration[0]["transfer_interval"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "archival configuration transfer-interval "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case itemTrim == "archival configuration transfer-on-commit":
					confRead.archivalConfiguration[0]["transfer_on_commit"] = true
				}
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
			case bchk.StringHasOneOfPrefixes(itemTrim, listLinesLicense()):
				if err := readSystemLicense(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, listLinesLogin()):
				if err := readSystemLogin(&confRead, itemTrim); err != nil {
					return confRead, err
				}
			case strings.HasPrefix(itemTrim, "max-configuration-rollbacks "):
				var err error
				confRead.maxConfigurationRollbacks, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-configuration-rollbacks "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "max-configurations-on-flash "):
				var err error
				confRead.maxConfigurationsOnFlash, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "max-configurations-on-flash "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, listLinesNtp()):
				if len(confRead.ntp) == 0 {
					confRead.ntp = append(confRead.ntp, map[string]interface{}{
						"boot_server":              "",
						"broadcast_client":         false,
						"interval_range":           -1,
						"multicast_client":         false,
						"multicast_client_address": "",
						"threshold_action":         "",
						"threshold_value":          0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "ntp boot-server "):
					confRead.ntp[0]["boot_server"] = strings.TrimPrefix(itemTrim, "ntp boot-server ")
				case itemTrim == "ntp broadcast-client":
					confRead.ntp[0]["broadcast_client"] = true
				case strings.HasPrefix(itemTrim, "ntp interval-range "):
					var err error
					confRead.ntp[0]["interval_range"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "ntp interval-range "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "ntp multicast-client"):
					confRead.ntp[0]["multicast_client"] = true
					if strings.HasPrefix(itemTrim, "ntp multicast-client ") {
						confRead.ntp[0]["multicast_client_address"] = strings.TrimPrefix(itemTrim, "ntp multicast-client ")
					}
				case strings.HasPrefix(itemTrim, "ntp threshold action "):
					confRead.ntp[0]["threshold_action"] = strings.TrimPrefix(itemTrim, "ntp threshold action ")
				case strings.HasPrefix(itemTrim, "ntp threshold "):
					var err error
					confRead.ntp[0]["threshold_value"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "ntp threshold "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case strings.HasPrefix(itemTrim, "ports "):
				if len(confRead.ports) == 0 {
					confRead.ports = append(confRead.ports, map[string]interface{}{
						"auxiliary_authentication_order": make([]string, 0),
						"auxiliary_disable":              false,
						"auxiliary_insecure":             false,
						"auxiliary_logout_on_disconnect": false,
						"auxiliary_type":                 "",
						"console_authentication_order":   make([]string, 0),
						"console_disable":                false,
						"console_insecure":               false,
						"console_logout_on_disconnect":   false,
						"console_type":                   "",
					})
				}
				itemTrimPorts := strings.TrimPrefix(itemTrim, "ports ")
				switch {
				case strings.HasPrefix(itemTrimPorts, "auxiliary authentication-order "):
					confRead.ports[0]["auxiliary_authentication_order"] = append(
						confRead.ports[0]["auxiliary_authentication_order"].([]string),
						strings.TrimPrefix(itemTrimPorts, "auxiliary authentication-order "))
				case itemTrimPorts == "auxiliary disable":
					confRead.ports[0]["auxiliary_disable"] = true
				case itemTrimPorts == "auxiliary insecure":
					confRead.ports[0]["auxiliary_insecure"] = true
				case itemTrimPorts == "auxiliary log-out-on-disconnect":
					confRead.ports[0]["auxiliary_logout_on_disconnect"] = true
				case strings.HasPrefix(itemTrimPorts, "auxiliary type "):
					confRead.ports[0]["auxiliary_type"] = strings.TrimPrefix(itemTrimPorts, "auxiliary type ")
				case strings.HasPrefix(itemTrimPorts, "console authentication-order "):
					confRead.ports[0]["console_authentication_order"] = append(
						confRead.ports[0]["console_authentication_order"].([]string),
						strings.TrimPrefix(itemTrimPorts, "console authentication-order "))
				case itemTrimPorts == "console disable":
					confRead.ports[0]["console_disable"] = true
				case itemTrimPorts == "console insecure":
					confRead.ports[0]["console_insecure"] = true
				case itemTrimPorts == "console log-out-on-disconnect":
					confRead.ports[0]["console_logout_on_disconnect"] = true
				case strings.HasPrefix(itemTrimPorts, "console type "):
					confRead.ports[0]["console_type"] = strings.TrimPrefix(itemTrimPorts, "console type ")
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
			case strings.HasPrefix(itemTrim, "radius-options attributes nas-ip-address "):
				confRead.radiusOptionsAttributesNasIPAddress = strings.TrimPrefix(itemTrim,
					"radius-options attributes nas-ip-address ")
			case itemTrim == "radius-options enhanced-accounting":
				confRead.radiusOptionsEnhancedAccounting = true
			case itemTrim == "radius-options password-protocol mschap-v2":
				confRead.radiusOptionsPasswodProtoclMsChapV2 = true
			case bchk.StringHasOneOfPrefixes(itemTrim, listLinesServices()):
				if len(confRead.services) == 0 {
					confRead.services = append(confRead.services, map[string]interface{}{
						"netconf_ssh":          make([]map[string]interface{}, 0),
						"netconf_traceoptions": make([]map[string]interface{}, 0),
						"ssh":                  make([]map[string]interface{}, 0),
						"web_management_http":  make([]map[string]interface{}, 0),
						"web_management_https": make([]map[string]interface{}, 0),
					})
				}
				if bchk.StringHasOneOfPrefixes(itemTrim, listLinesServicesNetconfSSH()) {
					if err := readSystemServicesNetconfSSH(&confRead, itemTrim); err != nil {
						return confRead, err
					}
				}
				if strings.HasPrefix(itemTrim, "services netconf traceoptions ") {
					if err := readSystemServicesNetconfTraceOpts(&confRead, itemTrim); err != nil {
						return confRead, err
					}
				}
				if bchk.StringHasOneOfPrefixes(itemTrim, listLinesServicesSSH()) {
					if err := readSystemServicesSSH(&confRead, itemTrim); err != nil {
						return confRead, err
					}
				}
				if bchk.StringHasOneOfPrefixes(itemTrim, listLinesServicesWebManagement()) {
					if err := readSystemServicesWebManagement(&confRead, itemTrim); err != nil {
						return confRead, err
					}
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, listLinesSyslog()):
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
		confRead.login[0]["announcement"] = html.UnescapeString(strings.Trim(strings.TrimPrefix(
			itemTrim, "login announcement "), "\""))
	case strings.HasPrefix(itemTrim, "login deny-sources address "):
		confRead.login[0]["deny_sources_address"] = append(confRead.login[0]["deny_sources_address"].([]string),
			strings.TrimPrefix(itemTrim, "login deny-sources address "))
	case strings.HasPrefix(itemTrim, "login idle-timeout "):
		var err error
		confRead.login[0]["idle_timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "login idle-timeout "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
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
		password := confRead.login[0]["password"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "login password change-type "):
			password["change_type"] = strings.TrimPrefix(itemTrim, "login password change-type ")
		case strings.HasPrefix(itemTrim, "login password format "):
			password["format"] = strings.TrimPrefix(itemTrim, "login password format ")
		case strings.HasPrefix(itemTrim, "login password maximum-length "):
			var err error
			password["maximum_length"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password maximum-length "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-changes "):
			var err error
			password["minimum_changes"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password minimum-changes "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-character-changes "):
			var err error
			password["minimum_character_changes"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password minimum-character-changes "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-length "):
			var err error
			password["minimum_length"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password minimum-length "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-lower-cases "):
			var err error
			password["minimum_lower_cases"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password minimum-lower-cases "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-numerics "):
			var err error
			password["minimum_numerics"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password minimum-numerics "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-punctuations "):
			var err error
			password["minimum_punctuations"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password minimum-punctuations "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-reuse "):
			var err error
			password["minimum_reuse"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password minimum-reuse "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login password minimum-upper-cases "):
			var err error
			password["minimum_upper_cases"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login password minimum-upper-cases "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
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
		retryOptions := confRead.login[0]["retry_options"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "login retry-options backoff-factor "):
			var err error
			retryOptions["backoff_factor"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login retry-options backoff-factor "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options backoff-threshold "):
			var err error
			retryOptions["backoff_threshold"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login retry-options backoff-threshold "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options lockout-period "):
			var err error
			retryOptions["lockout_period"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login retry-options lockout-period "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options maximum-time "):
			var err error
			retryOptions["maximum_time"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login retry-options maximum-time "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options minimum-time "):
			var err error
			retryOptions["minimum_time"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login retry-options minimum-time "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "login retry-options tries-before-disconnect "):
			var err error
			retryOptions["tries_before_disconnect"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "login retry-options tries-before-disconnect "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
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
		icmpV4RateLimit := confRead.internetOptions[0]["icmpv4_rate_limit"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "internet-options icmpv4-rate-limit bucket-size "):
			var err error
			icmpV4RateLimit["bucket_size"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "internet-options icmpv4-rate-limit bucket-size "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "internet-options icmpv4-rate-limit packet-rate "):
			var err error
			icmpV4RateLimit["packet_rate"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "internet-options icmpv4-rate-limit packet-rate "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
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
		icmpV6RateLimit := confRead.internetOptions[0]["icmpv6_rate_limit"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "internet-options icmpv6-rate-limit bucket-size "):
			var err error
			icmpV6RateLimit["bucket_size"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "internet-options icmpv6-rate-limit bucket-size "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "internet-options icmpv6-rate-limit packet-rate "):
			var err error
			icmpV6RateLimit["packet_rate"], err = strconv.Atoi(strings.TrimPrefix(
				itemTrim, "internet-options icmpv6-rate-limit packet-rate "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	case itemTrim == "internet-options ipip-path-mtu-discovery":
		confRead.internetOptions[0]["ipip_path_mtu_discovery"] = true
	case strings.HasPrefix(itemTrim, "internet-options ipv6-duplicate-addr-detection-transmits "):
		var err error
		confRead.internetOptions[0]["ipv6_duplicate_addr_detection_transmits"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "internet-options ipv6-duplicate-addr-detection-transmits "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "internet-options ipv6-path-mtu-discovery":
		confRead.internetOptions[0]["ipv6_path_mtu_discovery"] = true
	case strings.HasPrefix(itemTrim, "internet-options ipv6-path-mtu-discovery-timeout "):
		var err error
		confRead.internetOptions[0]["ipv6_path_mtu_discovery_timeout"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "internet-options ipv6-path-mtu-discovery-timeout "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
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
		confRead.internetOptions[0]["source_port_upper_limit"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "internet-options source-port upper-limit "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "internet-options source-quench":
		confRead.internetOptions[0]["source_quench"] = true
	case itemTrim == "internet-options tcp-drop-synfin-set":
		confRead.internetOptions[0]["tcp_drop_synfin_set"] = true
	case strings.HasPrefix(itemTrim, "internet-options tcp-mss "):
		var err error
		confRead.internetOptions[0]["tcp_mss"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "internet-options tcp-mss "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	}

	return nil
}

func readSystemLicense(confRead *systemOptions, itemTrim string) error {
	if len(confRead.license) == 0 {
		confRead.license = append(confRead.license, map[string]interface{}{
			"autoupdate":              false,
			"autoupdate_password":     "",
			"autoupdate_url":          "",
			"renew_before_expiration": -1,
			"renew_interval":          0,
		})
	}
	switch {
	case itemTrim == "license autoupdate":
		confRead.license[0]["autoupdate"] = true
	case strings.HasPrefix(itemTrim, "license autoupdate url "):
		confRead.license[0]["autoupdate"] = true
		itemTrimAutoupdateSplit := strings.Split(strings.TrimPrefix(itemTrim, "license autoupdate url "), " ")
		confRead.license[0]["autoupdate_url"] = strings.Trim(itemTrimAutoupdateSplit[0], "\"")

		itemTrimPassword := strings.TrimPrefix(itemTrim, "license autoupdate url "+itemTrimAutoupdateSplit[0]+" ")
		if strings.HasPrefix(itemTrimPassword, "password ") {
			var err error
			confRead.license[0]["autoupdate_password"], err = jdecode.Decode(strings.Trim(strings.TrimPrefix(
				itemTrimPassword, "password "), "\""))
			if err != nil {
				return fmt.Errorf("failed to decode password: %w", err)
			}
		}
	case strings.HasPrefix(itemTrim, "license renew before-expiration "):
		var err error
		confRead.license[0]["renew_before_expiration"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "license renew before-expiration "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "license renew interval "):
		var err error
		confRead.license[0]["renew_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "license renew interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	}

	return nil
}

func readSystemServicesNetconfTraceOpts(confRead *systemOptions, itemTrimNetconfTraceOpts string) error {
	if len(confRead.services[0]["netconf_traceoptions"].([]map[string]interface{})) == 0 {
		confRead.services[0]["netconf_traceoptions"] = append(
			confRead.services[0]["netconf_traceoptions"].([]map[string]interface{}),
			map[string]interface{}{
				"file_name":              "",
				"file_files":             0,
				"file_match":             "",
				"file_no_world_readable": false,
				"file_size":              0,
				"file_world_readable":    false,
				"flag":                   make([]string, 0),
				"no_remote_trace":        false,
				"on_demand":              false,
			})
	}
	netconfTraceOpts := confRead.services[0]["netconf_traceoptions"].([]map[string]interface{})[0]
	itemTrim := strings.TrimPrefix(itemTrimNetconfTraceOpts, "services netconf traceoptions ")
	switch {
	case strings.HasPrefix(itemTrim, "file files "):
		var err error
		netconfTraceOpts["file_files"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "file files "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "file match "):
		netconfTraceOpts["file_match"] = strings.Trim(strings.TrimPrefix(itemTrim, "file match "), "\"")
	case itemTrim == "file no-world-readable":
		netconfTraceOpts["file_no_world_readable"] = true
	case strings.HasPrefix(itemTrim, "file size "):
		var err error
		switch {
		case strings.HasSuffix(itemTrim, "k"):
			netconfTraceOpts["file_size"], err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(
				itemTrim, "file size "), "k"))
			netconfTraceOpts["file_size"] = netconfTraceOpts["file_size"].(int) * 1024
		case strings.HasSuffix(itemTrim, "m"):
			netconfTraceOpts["file_size"], err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(
				itemTrim, "file size "), "m"))
			netconfTraceOpts["file_size"] = netconfTraceOpts["file_size"].(int) * 1024 * 1024
		case strings.HasSuffix(itemTrim, "g"):
			netconfTraceOpts["file_size"], err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(
				itemTrim, "file size "), "g"))
			netconfTraceOpts["file_size"] = netconfTraceOpts["file_size"].(int) * 1024 * 1024 * 1024
		default:
			netconfTraceOpts["file_size"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "file size "))
		}
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "file world-readable":
		netconfTraceOpts["file_world_readable"] = true
	case strings.HasPrefix(itemTrim, "file "):
		netconfTraceOpts["file_name"] = strings.Trim(strings.TrimPrefix(itemTrim, "file "), "\"")
	case strings.HasPrefix(itemTrim, "flag "):
		netconfTraceOpts["flag"] = append(netconfTraceOpts["flag"].([]string), strings.TrimPrefix(itemTrim, "flag "))
	case itemTrim == "no-remote-trace":
		netconfTraceOpts["no_remote_trace"] = true
	case itemTrim == "on-demand":
		netconfTraceOpts["on_demand"] = true
	}

	return nil
}

func readSystemServicesNetconfSSH(confRead *systemOptions, itemTrim string) error {
	if len(confRead.services[0]["netconf_ssh"].([]map[string]interface{})) == 0 {
		confRead.services[0]["netconf_ssh"] = append(confRead.services[0]["netconf_ssh"].([]map[string]interface{}),
			map[string]interface{}{
				"client_alive_count_max": -1,
				"client_alive_interval":  -1,
				"connection_limit":       0,
				"rate_limit":             0,
			})
	}
	netconfSSH := confRead.services[0]["netconf_ssh"].([]map[string]interface{})[0]
	switch {
	case strings.HasPrefix(itemTrim, "services netconf ssh client-alive-count-max "):
		var err error
		netconfSSH["client_alive_count_max"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "services netconf ssh client-alive-count-max "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services netconf ssh client-alive-interval "):
		var err error
		netconfSSH["client_alive_interval"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "services netconf ssh client-alive-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services netconf ssh connection-limit "):
		var err error
		netconfSSH["connection_limit"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "services netconf ssh connection-limit "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services netconf ssh rate-limit "):
		var err error
		netconfSSH["rate_limit"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "services netconf ssh rate-limit "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
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
	ssh := confRead.services[0]["ssh"].([]map[string]interface{})[0]
	switch {
	case strings.HasPrefix(itemTrim, "services ssh authentication-order "):
		ssh["authentication_order"] = append(ssh["authentication_order"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh authentication-order "))
	case strings.HasPrefix(itemTrim, "services ssh ciphers "):
		ssh["ciphers"] = append(ssh["ciphers"].([]string),
			strings.Trim(strings.TrimPrefix(itemTrim, "services ssh ciphers "), "\""))
	case strings.HasPrefix(itemTrim, "services ssh client-alive-count-max "):
		var err error
		ssh["client_alive_count_max"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "services ssh client-alive-count-max "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh client-alive-interval "):
		var err error
		ssh["client_alive_interval"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh client-alive-interval "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh connection-limit "):
		var err error
		ssh["connection_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh connection-limit "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh fingerprint-hash "):
		ssh["fingerprint_hash"] = strings.TrimPrefix(itemTrim, "services ssh fingerprint-hash ")
	case strings.HasPrefix(itemTrim, "services ssh hostkey-algorithm "):
		ssh["hostkey_algorithm"] = append(ssh["hostkey_algorithm"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh hostkey-algorithm "))
	case strings.HasPrefix(itemTrim, "services ssh key-exchange "):
		ssh["key_exchange"] = append(ssh["key_exchange"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh key-exchange "))
	case itemTrim == "services ssh log-key-changes":
		ssh["log_key_changes"] = true
	case strings.HasPrefix(itemTrim, "services ssh macs "):
		ssh["macs"] = append(ssh["macs"].([]string), strings.TrimPrefix(itemTrim, "services ssh macs "))
	case strings.HasPrefix(itemTrim, "services ssh max-pre-authentication-packets "):
		var err error
		ssh["max_pre_authentication_packets"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "services ssh max-pre-authentication-packets "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh max-sessions-per-connection "):
		var err error
		ssh["max_sessions_per_connection"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "services ssh max-sessions-per-connection "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case itemTrim == "services ssh no-passwords":
		ssh["no_passwords"] = true
	case itemTrim == "services ssh no-public-keys":
		ssh["no_public_keys"] = true
	case strings.HasPrefix(itemTrim, "services ssh port "):
		var err error
		ssh["port"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh port "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh protocol-version "):
		ssh["protocol_version"] = append(ssh["protocol_version"].([]string),
			strings.TrimPrefix(itemTrim, "services ssh protocol-version "))
	case strings.HasPrefix(itemTrim, "services ssh rate-limit "):
		var err error
		ssh["rate_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "services ssh rate-limit "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "services ssh root-login "):
		ssh["root_login"] = strings.TrimPrefix(itemTrim, "services ssh root-login ")
	case itemTrim == "services ssh no-tcp-forwarding":
		ssh["no_tcp_forwarding"] = true
	case itemTrim == "services ssh tcp-forwarding":
		ssh["tcp_forwarding"] = true
	}

	return nil
}

func readSystemServicesWebManagement(confRead *systemOptions, itemTrim string) error {
	switch {
	case strings.HasPrefix(itemTrim, "services web-management https "):
		if len(confRead.services[0]["web_management_https"].([]map[string]interface{})) == 0 {
			confRead.services[0]["web_management_https"] = append(
				confRead.services[0]["web_management_https"].([]map[string]interface{}),
				map[string]interface{}{
					"interface":                    make([]string, 0),
					"port":                         0,
					"local_certificate":            "",
					"pki_local_certificate":        "",
					"system_generated_certificate": false,
				})
		}
		webMHTTPS := confRead.services[0]["web_management_https"].([]map[string]interface{})[0]
		if strings.HasPrefix(itemTrim, "services web-management https interface ") {
			webMHTTPS["interface"] = append(webMHTTPS["interface"].([]string),
				strings.TrimPrefix(itemTrim, "services web-management https interface "))
		}
		if strings.HasPrefix(itemTrim, "services web-management https port ") {
			var err error
			webMHTTPS["port"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "services web-management https port "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
		if strings.HasPrefix(itemTrim, "services web-management https local-certificate ") {
			webMHTTPS["local_certificate"] = strings.Trim(strings.TrimPrefix(itemTrim,
				"services web-management https local-certificate "), "\"")
		}
		if strings.HasPrefix(itemTrim, "services web-management https pki-local-certificate ") {
			webMHTTPS["pki_local_certificate"] = strings.Trim(strings.TrimPrefix(itemTrim,
				"services web-management https pki-local-certificate "), "\"")
		}
		if itemTrim == "services web-management https system-generated-certificate" {
			webMHTTPS["system_generated_certificate"] = true
		}
	case strings.HasPrefix(itemTrim, "services web-management http"):
		if len(confRead.services[0]["web_management_http"].([]map[string]interface{})) == 0 {
			confRead.services[0]["web_management_http"] = append(
				confRead.services[0]["web_management_http"].([]map[string]interface{}),
				map[string]interface{}{
					"interface": make([]string, 0),
					"port":      0,
				})
		}
		webMHTTP := confRead.services[0]["web_management_http"].([]map[string]interface{})[0]
		if strings.HasPrefix(itemTrim, "services web-management http interface ") {
			webMHTTP["interface"] = append(webMHTTP["interface"].([]string),
				strings.TrimPrefix(itemTrim, "services web-management http interface "))
		}
		if strings.HasPrefix(itemTrim, "services web-management http port ") {
			var err error
			webMHTTP["port"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "services web-management http port "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		}
	}

	return nil
}

func readSystemSyslog(confRead *systemOptions, itemTrim string) error {
	if len(confRead.syslog) == 0 {
		confRead.syslog = append(confRead.syslog, map[string]interface{}{
			"archive":                 make([]map[string]interface{}, 0),
			"console":                 make([]map[string]interface{}, 0),
			"log_rotate_frequency":    0,
			"source_address":          "",
			"time_format_millisecond": false,
			"time_format_year":        false,
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
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case strings.HasPrefix(itemTrim, "syslog archive size "):
			var err error
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["size"], err = strconv.Atoi(
				strings.TrimPrefix(itemTrim, "syslog archive size "))
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case itemTrim == "syslog archive no-world-readable":
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["no_world_readable"] = true
		case itemTrim == "syslog archive world-readable":
			confRead.syslog[0]["archive"].([]map[string]interface{})[0]["world_readable"] = true
		}
	case strings.HasPrefix(itemTrim, "syslog console "):
		if len(confRead.syslog[0]["console"].([]map[string]interface{})) == 0 {
			confRead.syslog[0]["console"] = append(confRead.syslog[0]["console"].([]map[string]interface{}),
				map[string]interface{}{
					"any_severity":                 "",
					"authorization_severity":       "",
					"changelog_severity":           "",
					"conflictlog_severity":         "",
					"daemon_severity":              "",
					"dfc_severity":                 "",
					"external_severity":            "",
					"firewall_severity":            "",
					"ftp_severity":                 "",
					"interactivecommands_severity": "",
					"kernel_severity":              "",
					"ntp_severity":                 "",
					"pfe_severity":                 "",
					"security_severity":            "",
					"user_severity":                "",
				})
		}
		console := confRead.syslog[0]["console"].([]map[string]interface{})[0]
		switch {
		case strings.HasPrefix(itemTrim, "syslog console any "):
			console["any_severity"] = strings.TrimPrefix(itemTrim, "syslog console any ")
		case strings.HasPrefix(itemTrim, "syslog console authorization "):
			console["authorization_severity"] = strings.TrimPrefix(itemTrim, "syslog console authorization ")
		case strings.HasPrefix(itemTrim, "syslog console change-log "):
			console["changelog_severity"] = strings.TrimPrefix(itemTrim, "syslog console change-log ")
		case strings.HasPrefix(itemTrim, "syslog console conflict-log "):
			console["conflictlog_severity"] = strings.TrimPrefix(itemTrim, "syslog console conflict-log ")
		case strings.HasPrefix(itemTrim, "syslog console daemon "):
			console["daemon_severity"] = strings.TrimPrefix(itemTrim, "syslog console daemon ")
		case strings.HasPrefix(itemTrim, "syslog console dfc "):
			console["dfc_severity"] = strings.TrimPrefix(itemTrim, "syslog console dfc ")
		case strings.HasPrefix(itemTrim, "syslog console external "):
			console["external_severity"] = strings.TrimPrefix(itemTrim, "syslog console external ")
		case strings.HasPrefix(itemTrim, "syslog console firewall "):
			console["firewall_severity"] = strings.TrimPrefix(itemTrim, "syslog console firewall ")
		case strings.HasPrefix(itemTrim, "syslog console ftp "):
			console["ftp_severity"] = strings.TrimPrefix(itemTrim, "syslog console ftp ")
		case strings.HasPrefix(itemTrim, "syslog console interactive-commands "):
			console["interactivecommands_severity"] = strings.TrimPrefix(itemTrim, "syslog console interactive-commands ")
		case strings.HasPrefix(itemTrim, "syslog console kernel "):
			console["kernel_severity"] = strings.TrimPrefix(itemTrim, "syslog console kernel ")
		case strings.HasPrefix(itemTrim, "syslog console ntp "):
			console["ntp_severity"] = strings.TrimPrefix(itemTrim, "syslog console ntp ")
		case strings.HasPrefix(itemTrim, "syslog console pfe "):
			console["pfe_severity"] = strings.TrimPrefix(itemTrim, "syslog console pfe ")
		case strings.HasPrefix(itemTrim, "syslog console security "):
			console["security_severity"] = strings.TrimPrefix(itemTrim, "syslog console security ")
		case strings.HasPrefix(itemTrim, "syslog console user "):
			console["user_severity"] = strings.TrimPrefix(itemTrim, "syslog console user ")
		}
	case strings.HasPrefix(itemTrim, "syslog log-rotate-frequency "):
		var err error
		confRead.syslog[0]["log_rotate_frequency"], err = strconv.Atoi(
			strings.TrimPrefix(itemTrim, "syslog log-rotate-frequency "))
		if err != nil {
			return fmt.Errorf(failedConvAtoiError, itemTrim, err)
		}
	case strings.HasPrefix(itemTrim, "syslog source-address "):
		confRead.syslog[0]["source_address"] = strings.TrimPrefix(
			itemTrim, "syslog source-address ")
	case itemTrim == "syslog time-format millisecond":
		confRead.syslog[0]["time_format_millisecond"] = true
	case itemTrim == "syslog time-format year":
		confRead.syslog[0]["time_format_year"] = true
	}

	return nil
}

func fillSystem(d *schema.ResourceData, systemOptions systemOptions) {
	if tfErr := d.Set("archival_configuration", systemOptions.archivalConfiguration); tfErr != nil {
		panic(tfErr)
	}
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
	if tfErr := d.Set("license", systemOptions.license); tfErr != nil {
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
	if tfErr := d.Set("ntp", systemOptions.ntp); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ports", systemOptions.ports); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("radius_options_attributes_nas_ipaddress",
		systemOptions.radiusOptionsAttributesNasIPAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("radius_options_enhanced_accounting",
		systemOptions.radiusOptionsEnhancedAccounting); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("radius_options_password_protocol_mschapv2",
		systemOptions.radiusOptionsPasswodProtoclMsChapV2); tfErr != nil {
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
