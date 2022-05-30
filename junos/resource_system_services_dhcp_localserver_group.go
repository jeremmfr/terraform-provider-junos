package junos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type systemServicesDhcpLocalServerGroupOptions struct {
	dynamicProfileAggregateClients       bool
	reauthenticateLeaseRenewal           bool
	reauthenticateRemoteIDMismatch       bool
	remoteIDMismatchDisconnect           bool
	routeSuppressionAccess               bool
	routeSuppressionAccessInternal       bool
	routeSuppressionDestination          bool
	shortCycleProtectionLockoutMaxTime   int
	shortCycleProtectionLockoutMinTime   int
	accessProfile                        string
	authenticationPassword               string
	dynamicProfile                       string
	dynamicProfileUsePrimary             string
	dynamicProfileAggregateClientsAction string
	livenessDetectionFailureAction       string
	name                                 string
	routingInstance                      string
	serviceProfile                       string
	version                              string
	authenticationUsernameInclude        []map[string]interface{}
	interFace                            []map[string]interface{}
	leaseTimeValidation                  []map[string]interface{}
	livenessDetectionMethodBfd           []map[string]interface{}
	livenessDetectionMethodLayer2        []map[string]interface{}
	overridesV4                          []map[string]interface{}
	overridesV6                          []map[string]interface{}
	reconfigure                          []map[string]interface{}
}

func resourceSystemServicesDhcpLocalServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSystemServicesDhcpLocalServerGroupCreate,
		ReadWithoutTimeout:   resourceSystemServicesDhcpLocalServerGroupRead,
		UpdateWithoutTimeout: resourceSystemServicesDhcpLocalServerGroupUpdate,
		DeleteWithoutTimeout: resourceSystemServicesDhcpLocalServerGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemServicesDhcpLocalServerGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
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
				ForceNew:     true,
				Default:      "v4",
				ValidateFunc: validation.StringInSlice([]string{"v4", "v6"}, false),
			},
			"access_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authentication_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authentication_username_include": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"circuit_type": {
							Type:     schema.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{
								"authentication_username_include.0.circuit_type",
								"authentication_username_include.0.client_id",
								"authentication_username_include.0.delimiter",
								"authentication_username_include.0.domain_name",
								"authentication_username_include.0.interface_description",
								"authentication_username_include.0.interface_name",
								"authentication_username_include.0.mac_address",
								"authentication_username_include.0.option_60",
								"authentication_username_include.0.option_82",
								"authentication_username_include.0.relay_agent_interface_id",
								"authentication_username_include.0.relay_agent_remote_id",
								"authentication_username_include.0.relay_agent_subscriber_id",
								"authentication_username_include.0.routing_instance_name",
								"authentication_username_include.0.user_prefix",
								"authentication_username_include.0.vlan_tags",
							},
						},
						"client_id": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"delimiter": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 1),
						},
						"domain_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"interface_description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"device", "logical"}, false),
						},
						"interface_name": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"mac_address": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"option_60": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"option_82": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"option_82_circuit_id": {
							Type:         schema.TypeBool,
							Optional:     true,
							RequiredWith: []string{"authentication_username_include.0.option_82"},
						},
						"option_82_remote_id": {
							Type:         schema.TypeBool,
							Optional:     true,
							RequiredWith: []string{"authentication_username_include.0.option_82"},
						},
						"relay_agent_interface_id": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"relay_agent_remote_id": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"relay_agent_subscriber_id": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"routing_instance_name": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"user_prefix": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"vlan_tags": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"dynamic_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dynamic_profile_use_primary": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"dynamic_profile"},
				ConflictsWith: []string{"dynamic_profile_aggregate_clients"},
			},
			"dynamic_profile_aggregate_clients": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"dynamic_profile"},
			},
			"dynamic_profile_aggregate_clients_action": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"dynamic_profile_aggregate_clients"},
				ValidateFunc: validation.StringInSlice([]string{"merge", "replace"}, false),
			},
			"interface": {
				Type:     schema.TypeSet,
				Optional: true,
				AtLeastOneOf: []string{
					"access_profile",
					"authentication_password",
					"authentication_username_include",
					"dynamic_profile",
					"interface",
					"lease_time_validation",
					"liveness_detection_failure_action",
					"liveness_detection_method_bfd",
					"liveness_detection_method_layer2",
					"overrides_v4",
					"overrides_v6",
					"reauthenticate_lease_renewal",
					"reauthenticate_remote_id_mismatch",
					"reconfigure",
					"remote_id_mismatch_disconnect",
					"route_suppression_access",
					"route_suppression_access_internal",
					"route_suppression_destination",
					"service_profile",
					"short_cycle_protection_lockout_max_time",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if strings.Count(value, ".") != 1 {
									errors = append(errors, fmt.Errorf(
										"%q in %q need to have 1 dot", value, k))
								}

								return
							},
						},
						"access_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dynamic_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dynamic_profile_use_primary": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dynamic_profile_aggregate_clients": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"dynamic_profile_aggregate_clients_action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"merge", "replace"}, false),
						},
						"exclude": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"overrides_v4": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: schemaSystemServicesDhcpLocalServerGroupOverridesV4(),
							},
						},
						"overrides_v6": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: schemaSystemServicesDhcpLocalServerGroupOverridesV6(),
							},
						},
						"service_profile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"short_cycle_protection_lockout_max_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 86400),
						},
						"short_cycle_protection_lockout_min_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 86400),
						},
						"trace": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"upto": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if strings.Count(value, ".") != 1 {
									errors = append(errors, fmt.Errorf(
										"%q in %q need to have 1 dot", value, k))
								}

								return
							},
						},
					},
				},
			},
			"lease_time_validation": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lease_time_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 2147483647),
						},
						"violation_action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"override-lease", "strict"}, false),
						},
					},
				},
			},
			"liveness_detection_failure_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"clear-binding",
					"clear-binding-if-interface-up",
					"log-only",
				}, false),
			},
			"liveness_detection_method_bfd": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"liveness_detection_method_layer2"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"detection_time_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							AtLeastOneOf: []string{
								"liveness_detection_method_bfd.0.detection_time_threshold",
								"liveness_detection_method_bfd.0.holddown_interval",
								"liveness_detection_method_bfd.0.minimum_interval",
								"liveness_detection_method_bfd.0.minimum_receive_interval",
								"liveness_detection_method_bfd.0.multiplier",
								"liveness_detection_method_bfd.0.no_adaptation",
								"liveness_detection_method_bfd.0.session_mode",
								"liveness_detection_method_bfd.0.transmit_interval_minimum",
								"liveness_detection_method_bfd.0.transmit_interval_threshold",
								"liveness_detection_method_bfd.0.version",
							},
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
						"holddown_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 255000),
						},
						"minimum_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30000, 255000),
						},
						"minimum_receive_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30000, 255000),
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
						"session_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"automatic", "multihop", "single-hop"}, false),
						},
						"transmit_interval_minimum": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30000, 255000),
						},
						"transmit_interval_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
						"version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"0", "1", "automatic"}, false),
						},
					},
				},
			},
			"liveness_detection_method_layer2": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"liveness_detection_method_bfd"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_consecutive_retries": {
							Type:     schema.TypeInt,
							Optional: true,
							AtLeastOneOf: []string{
								"liveness_detection_method_layer2.0.max_consecutive_retries",
								"liveness_detection_method_layer2.0.transmit_interval",
							},
							ValidateFunc: validation.IntBetween(3, 6),
						},
						"transmit_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(300, 1800),
						},
					},
				},
			},
			"overrides_v4": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaSystemServicesDhcpLocalServerGroupOverridesV4(),
				},
			},
			"overrides_v6": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaSystemServicesDhcpLocalServerGroupOverridesV6(),
				},
			},
			"reauthenticate_lease_renewal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reauthenticate_remote_id_mismatch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reconfigure": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attempts": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"clear_on_abort": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"support_option_pd_exclude": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"token": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"trigger_radius_disconnect": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"remote_id_mismatch_disconnect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"route_suppression_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"route_suppression_access_internal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"route_suppression_destination": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"service_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"short_cycle_protection_lockout_max_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"short_cycle_protection_lockout_min_time"},
				ValidateFunc: validation.IntBetween(1, 86400),
			},
			"short_cycle_protection_lockout_min_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"short_cycle_protection_lockout_max_time"},
				ValidateFunc: validation.IntBetween(1, 86400),
			},
		},
	}
}

func schemaSystemServicesDhcpLocalServerGroupOverridesV4() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_no_end_option": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"asymmetric_lease_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(600, 86400),
		},
		"bootp_support": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"client_discover_match": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"incoming-interface", "option60-and-option82"}, false),
		},
		"delay_offer_based_on": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"option": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"option-60", "option-77", "option-82"}, false),
					},
					"compare": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"equals", "not-equals", "starts-with"}, false),
					},
					"value_type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"ascii", "hexadecimal"}, false),
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"delay_offer_delay_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 30),
		},
		"delete_binding_on_renegotiation": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"dual_stack": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"include_option_82_forcerenew": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"include_option_82_nak": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"interface_client_limit": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 500000),
		},
		"process_inform": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"process_inform_pool": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},
		"protocol_attributes": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},
	}
}

func schemaSystemServicesDhcpLocalServerGroupOverridesV6() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"always_add_option_dns_server": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"always_process_option_request_option": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"asymmetric_lease_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(600, 86400),
		},
		"asymmetric_prefix_lease_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(600, 86400),
		},
		"client_negotiation_match_incoming_interface": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"delay_advertise_based_on": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"option": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"option-15", "option-16", "option-18", "option-37"}, false),
					},
					"compare": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"equals", "not-equals", "starts-with"}, false),
					},
					"value_type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"ascii", "hexadecimal"}, false),
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"delay_advertise_delay_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 30),
		},
		"delegated_pool": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},
		"delete_binding_on_renegotiation": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"dual_stack": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"interface_client_limit": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 500000),
		},
		"multi_address_embedded_option_response": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"process_inform": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"process_inform_pool": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},
		"protocol_attributes": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},
		"rapid_commit": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"top_level_status_code": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func resourceSystemServicesDhcpLocalServerGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setSystemServicesDhcpLocalServerGroup(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string) +
			idSeparator + d.Get("routing_instance").(string) +
			idSeparator + d.Get("version").(string))

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
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", d.Get("routing_instance").(string)))...)
		}
	}
	systemServicesDhcpLocalServerGroupExists, err := checkSystemServicesDhcpLocalServerGroupExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("version").(string),
		clt, junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if systemServicesDhcpLocalServerGroupExists {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))
		if d.Get("version").(string) == "v6" {
			return append(diagWarns, diag.FromErr(
				fmt.Errorf("system services dhcp-local-server dhcpv6 group %v already exists in routing-instance %s",
					d.Get("name").(string), d.Get("routing_instance").(string)))...)
		}

		return append(diagWarns, diag.FromErr(
			fmt.Errorf("system services dhcp-local-server group %v already exists in routing-instance %s",
				d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}
	if err := setSystemServicesDhcpLocalServerGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("create resource junos_system_services_dhcp_localserver_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	systemServicesDhcpLocalServerGroupExists, err = checkSystemServicesDhcpLocalServerGroupExists(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("version").(string),
		clt, junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if systemServicesDhcpLocalServerGroupExists {
		d.SetId(d.Get("name").(string) +
			idSeparator + d.Get("routing_instance").(string) +
			idSeparator + d.Get("version").(string))
	} else {
		if d.Get("version").(string) == "v6" {
			return append(diagWarns,
				diag.FromErr(fmt.Errorf("system services dhcp-local-server dhcpv6 group %v "+
					"not exists in routing_instance %s after commit => check your config",
					d.Get("name").(string), d.Get("routing_instance").(string)))...)
		}

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("system services dhcp-local-server group %v "+
				"not exists in routing_instance %s after commit => check your config",
				d.Get("name").(string), d.Get("routing_instance").(string)))...)
	}

	return append(diagWarns, resourceSystemServicesDhcpLocalServerGroupReadWJunSess(d, clt, junSess)...)
}

func resourceSystemServicesDhcpLocalServerGroupRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceSystemServicesDhcpLocalServerGroupReadWJunSess(d, clt, junSess)
}

func resourceSystemServicesDhcpLocalServerGroupReadWJunSess(
	d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	systemServicesDhcpLocalServerGroupOptions, err := readSystemServicesDhcpLocalServerGroup(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("version").(string),
		clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if systemServicesDhcpLocalServerGroupOptions.name == "" {
		d.SetId("")
	} else {
		fillSystemServicesDhcpLocalServerGroupData(d, systemServicesDhcpLocalServerGroupOptions)
	}

	return nil
}

func resourceSystemServicesDhcpLocalServerGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delSystemServicesDhcpLocalServerGroup(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			d.Get("version").(string),
			clt, nil,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := setSystemServicesDhcpLocalServerGroup(d, clt, nil); err != nil {
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
	if err := delSystemServicesDhcpLocalServerGroup(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("version").(string),
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSystemServicesDhcpLocalServerGroup(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("update resource junos_system_services_dhcp_localserver_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSystemServicesDhcpLocalServerGroupReadWJunSess(d, clt, junSess)...)
}

func resourceSystemServicesDhcpLocalServerGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delSystemServicesDhcpLocalServerGroup(
			d.Get("name").(string),
			d.Get("routing_instance").(string),
			d.Get("version").(string),
			clt, nil,
		); err != nil {
			return diag.FromErr(err)
		}

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
	if err := delSystemServicesDhcpLocalServerGroup(
		d.Get("name").(string),
		d.Get("routing_instance").(string),
		d.Get("version").(string),
		clt, junSess,
	); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_system_services_dhcp_localserver_group", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSystemServicesDhcpLocalServerGroupImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 3 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	if idSplit[2] != "v4" && idSplit[2] != "v6" {
		return nil, fmt.Errorf("bad version '%s' in id, need to be 'v4' or 'v6' (id must be "+
			"<name>"+idSeparator+"<routing_instance>"+idSeparator+"<version>)", idSplit[2])
	}
	systemServicesDhcpLocalServerGroupExists, err := checkSystemServicesDhcpLocalServerGroupExists(
		idSplit[0],
		idSplit[1],
		idSplit[2],
		clt, junSess)
	if err != nil {
		return nil, err
	}
	if !systemServicesDhcpLocalServerGroupExists {
		if idSplit[2] == "v6" {
			return nil, fmt.Errorf("don't find system services dhcp-local-server dhcpv6 group with id '%v' (id must be "+
				"<name>"+idSeparator+"<routing_instance>"+idSeparator+"<version>)", d.Id())
		}

		return nil, fmt.Errorf("don't find system services dhcp-local-server group with id '%v' (id must be "+
			"<name>"+idSeparator+"<routing_instance>"+idSeparator+"<version>)", d.Id())
	}
	systemServicesDhcpLocalServerGroupOptions, err := readSystemServicesDhcpLocalServerGroup(
		idSplit[0],
		idSplit[1],
		idSplit[2],
		clt, junSess)
	if err != nil {
		return nil, err
	}
	fillSystemServicesDhcpLocalServerGroupData(d, systemServicesDhcpLocalServerGroupOptions)

	result[0] = d

	return result, nil
}

func checkSystemServicesDhcpLocalServerGroupExists(
	name, instance, version string, clt *Client, junSess *junosSession,
) (bool, error) {
	var showConfig string
	var err error
	showCmd := cmdShowConfig
	if instance != defaultW {
		showCmd += routingInstancesWS + instance + " "
	}
	showCmd += "system services dhcp-local-server "
	if version == "v6" {
		showCmd += "dhcpv6 group " + name
	} else {
		showCmd += "group " + name
	}

	showConfig, err = clt.command(showCmd+pipeDisplaySet, junSess)
	if err != nil {
		return false, err
	}

	if showConfig == emptyW {
		return false, nil
	}

	return true, nil
}

func setSystemServicesDhcpLocalServerGroup(d *schema.ResourceData, clt *Client, junSess *junosSession) error {
	configSet := make([]string, 0)

	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	if d.Get("version").(string) == "v6" {
		setPrefix += "system services dhcp-local-server dhcpv6 group " + d.Get("name").(string) + " "
	} else {
		setPrefix += "system services dhcp-local-server group " + d.Get("name").(string) + " "
	}

	if v := d.Get("access_profile").(string); v != "" {
		configSet = append(configSet, setPrefix+"access-profile \""+v+"\"")
	}
	if v := d.Get("authentication_password").(string); v != "" {
		configSet = append(configSet, setPrefix+"authentication password \""+v+"\"")
	}
	for _, vBlock := range d.Get("authentication_username_include").([]interface{}) {
		authenticationUsernameInclude := vBlock.(map[string]interface{})
		if authenticationUsernameInclude["circuit_type"].(bool) {
			configSet = append(configSet, setPrefix+"authentication username-include circuit-type")
		}
		if authenticationUsernameInclude["client_id"].(bool) {
			if d.Get("version").(string) == "v4" {
				return fmt.Errorf("client_id not compatible when version = v4")
			}
			configSet = append(configSet, setPrefix+"authentication username-include client-id")
		}
		if v := authenticationUsernameInclude["delimiter"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication username-include delimiter \""+v+"\"")
		}
		if v := authenticationUsernameInclude["domain_name"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication username-include domain-name \""+v+"\"")
		}
		if v := authenticationUsernameInclude["interface_description"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication username-include interface-description "+v)
		}
		if authenticationUsernameInclude["interface_name"].(bool) {
			configSet = append(configSet, setPrefix+"authentication username-include interface-name")
		}
		if authenticationUsernameInclude["mac_address"].(bool) {
			configSet = append(configSet, setPrefix+"authentication username-include mac-address")
		}
		if authenticationUsernameInclude["option_60"].(bool) {
			if d.Get("version").(string) == "v6" {
				return fmt.Errorf("option_60 not compatible when version = v6")
			}
			configSet = append(configSet, setPrefix+"authentication username-include option-60")
		}
		if authenticationUsernameInclude["option_82"].(bool) {
			if d.Get("version").(string) == "v6" {
				return fmt.Errorf("option_82 not compatible when version = v6")
			}
			configSet = append(configSet, setPrefix+"authentication username-include option-82")
			if authenticationUsernameInclude["option_82_circuit_id"].(bool) {
				configSet = append(configSet, setPrefix+"authentication username-include option-82 circuit-id")
			}
			if authenticationUsernameInclude["option_82_remote_id"].(bool) {
				configSet = append(configSet, setPrefix+"authentication username-include option-82 remote-id")
			}
		} else if authenticationUsernameInclude["option_82_circuit_id"].(bool) ||
			authenticationUsernameInclude["option_82_remote_id"].(bool) {
			return fmt.Errorf("authentication_username_include.0.option_82 need to be true with " +
				"option_82_circuit_id or option_82_remote_id")
		}
		if authenticationUsernameInclude["relay_agent_interface_id"].(bool) {
			if d.Get("version").(string) == "v4" {
				return fmt.Errorf("relay_agent_interface_id not compatible when version = v4")
			}
			configSet = append(configSet, setPrefix+"authentication username-include relay-agent-interface-id")
		}
		if authenticationUsernameInclude["relay_agent_remote_id"].(bool) {
			if d.Get("version").(string) == "v4" {
				return fmt.Errorf("relay_agent_remote_id not compatible when version = v4")
			}
			configSet = append(configSet, setPrefix+"authentication username-include relay-agent-remote-id")
		}
		if authenticationUsernameInclude["relay_agent_subscriber_id"].(bool) {
			if d.Get("version").(string) == "v4" {
				return fmt.Errorf("relay_agent_subscriber_id not compatible when version = v4")
			}
			configSet = append(configSet, setPrefix+"authentication username-include relay-agent-subscriber-id")
		}
		if authenticationUsernameInclude["routing_instance_name"].(bool) {
			configSet = append(configSet, setPrefix+"authentication username-include routing-instance-name")
		}
		if v := authenticationUsernameInclude["user_prefix"].(string); v != "" {
			configSet = append(configSet, setPrefix+"authentication username-include user-prefix \""+v+"\"")
		}
		if authenticationUsernameInclude["vlan_tags"].(bool) {
			configSet = append(configSet, setPrefix+"authentication username-include vlan-tags")
		}
	}
	if dynProfile := d.Get("dynamic_profile").(string); dynProfile != "" {
		configSet = append(configSet, setPrefix+"dynamic-profile \""+dynProfile+"\"")
		if v := d.Get("dynamic_profile_use_primary").(string); v != "" {
			configSet = append(configSet, setPrefix+"dynamic-profile use-primary \""+v+"\"")
		}
		if d.Get("dynamic_profile_aggregate_clients").(bool) {
			configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients")
			if v := d.Get("dynamic_profile_aggregate_clients_action").(string); v != "" {
				configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients "+v)
			}
		} else if d.Get("dynamic_profile_aggregate_clients_action").(string) != "" {
			return fmt.Errorf("dynamic_profile_aggregate_clients need to be true with " +
				"dynamic_profile_aggregate_clients_action")
		}
	} else if d.Get("dynamic_profile_use_primary").(string) != "" ||
		d.Get("dynamic_profile_aggregate_clients").(bool) ||
		d.Get("dynamic_profile_aggregate_clients_action").(string) != "" {
		return fmt.Errorf("dynamic_profile need to be set with " +
			"dynamic_profile_use_primary, dynamic_profile_aggregate_clients " +
			"and dynamic_profile_aggregate_clients_action")
	}
	interfaceNameList := make([]string, 0)
	for _, v := range d.Get("interface").(*schema.Set).List() {
		interFace := v.(map[string]interface{})
		if bchk.StringInSlice(interFace["name"].(string), interfaceNameList) {
			return fmt.Errorf("multiple blocks interface with the same name %s", interFace["name"].(string))
		}
		interfaceNameList = append(interfaceNameList, interFace["name"].(string))
		configSetInterface, err := setSystemServicesDhcpLocalServerGroupInterface(
			interFace, setPrefix, d.Get("version").(string))
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetInterface...)
	}
	for _, lTVal := range d.Get("lease_time_validation").([]interface{}) {
		configSet = append(configSet, setPrefix+"lease-time-validation")
		if lTVal != nil {
			leaseTimeValidation := lTVal.(map[string]interface{})
			if v := leaseTimeValidation["lease_time_threshold"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"lease-time-validation lease-time-threshold "+strconv.Itoa(v))
			}
			if v := leaseTimeValidation["violation_action"].(string); v != "" {
				configSet = append(configSet, setPrefix+"lease-time-validation violation-action "+v)
			}
		}
	}
	if v := d.Get("liveness_detection_failure_action").(string); v != "" {
		configSet = append(configSet, setPrefix+"liveness-detection failure-action "+v)
	}
	for _, ldmBfd := range d.Get("liveness_detection_method_bfd").([]interface{}) {
		liveDetectMethBfd := ldmBfd.(map[string]interface{})
		setPrefixLDMBfd := setPrefix + "liveness-detection method bfd "
		if v := liveDetectMethBfd["detection_time_threshold"].(int); v != -1 {
			configSet = append(configSet, setPrefixLDMBfd+"detection-time threshold "+strconv.Itoa(v))
		}
		if v := liveDetectMethBfd["holddown_interval"].(int); v != -1 {
			configSet = append(configSet, setPrefixLDMBfd+"holddown-interval "+strconv.Itoa(v))
		}
		if v := liveDetectMethBfd["minimum_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixLDMBfd+"minimum-interval "+strconv.Itoa(v))
		}
		if v := liveDetectMethBfd["minimum_receive_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixLDMBfd+"minimum-receive-interval "+strconv.Itoa(v))
		}
		if v := liveDetectMethBfd["multiplier"].(int); v != 0 {
			configSet = append(configSet, setPrefixLDMBfd+"multiplier "+strconv.Itoa(v))
		}
		if liveDetectMethBfd["no_adaptation"].(bool) {
			configSet = append(configSet, setPrefixLDMBfd+"no-adaptation")
		}
		if v := liveDetectMethBfd["session_mode"].(string); v != "" {
			configSet = append(configSet, setPrefixLDMBfd+"session-mode "+v)
		}
		if v := liveDetectMethBfd["transmit_interval_minimum"].(int); v != 0 {
			configSet = append(configSet, setPrefixLDMBfd+"transmit-interval minimum-interval "+strconv.Itoa(v))
		}
		if v := liveDetectMethBfd["transmit_interval_threshold"].(int); v != -1 {
			configSet = append(configSet, setPrefixLDMBfd+"transmit-interval threshold "+strconv.Itoa(v))
		}
		if v := liveDetectMethBfd["version"].(string); v != "" {
			configSet = append(configSet, setPrefixLDMBfd+"version "+v)
		}

		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefixLDMBfd) {
			return fmt.Errorf("liveness_detection_method_bfd block is empty")
		}
	}
	for _, ldmLayer2 := range d.Get("liveness_detection_method_layer2").([]interface{}) {
		if ldmLayer2 != nil {
			liveDetectMethLayer2 := ldmLayer2.(map[string]interface{})
			setPrefixLDMLayer2 := setPrefix + "liveness-detection method layer2-liveness-detection "
			if v := liveDetectMethLayer2["max_consecutive_retries"].(int); v != 0 {
				configSet = append(configSet, setPrefixLDMLayer2+"max-consecutive-retries "+strconv.Itoa(v))
			}
			if v := liveDetectMethLayer2["transmit_interval"].(int); v != 0 {
				configSet = append(configSet, setPrefixLDMLayer2+"transmit-interval "+strconv.Itoa(v))
			}
		} else {
			return fmt.Errorf("liveness_detection_method_layer2 block is empty")
		}
	}
	for _, v := range d.Get("overrides_v4").([]interface{}) {
		if d.Get("version").(string) == "v6" {
			return fmt.Errorf("overrides_v4 not compatible if version = v6")
		}
		configSetOverrides, err := setSystemServicesDhcpLocalServerGroupFamilyDhcpOverridesV4(
			v.(map[string]interface{}), setPrefix)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetOverrides...)
	}
	for _, v := range d.Get("overrides_v6").([]interface{}) {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("overrides_v6 not compatible if version = v4")
		}
		configSetOverrides, err := setSystemServicesDhcpLocalServerGroupFamilyDhcpOverridesV6(
			v.(map[string]interface{}), setPrefix)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetOverrides...)
	}
	if d.Get("reauthenticate_lease_renewal").(bool) {
		configSet = append(configSet, setPrefix+"reauthenticate lease-renewal")
	}
	if d.Get("reauthenticate_remote_id_mismatch").(bool) {
		configSet = append(configSet, setPrefix+"reauthenticate remote-id-mismatch")
	}
	for _, rec := range d.Get("reconfigure").([]interface{}) {
		configSet = append(configSet, setPrefix+"reconfigure")
		if rec != nil {
			reconfigure := rec.(map[string]interface{})
			if v := reconfigure["attempts"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"reconfigure attempts "+strconv.Itoa(v))
			}
			if reconfigure["clear_on_abort"].(bool) {
				configSet = append(configSet, setPrefix+"reconfigure clear-on-abort")
			}
			if reconfigure["support_option_pd_exclude"].(bool) {
				configSet = append(configSet, setPrefix+"reconfigure support-option-pd-exclude")
			}
			if v := reconfigure["timeout"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"reconfigure timeout "+strconv.Itoa(v))
			}
			if v := reconfigure["token"].(string); v != "" {
				configSet = append(configSet, setPrefix+"reconfigure token \""+v+"\"")
			}
			if reconfigure["trigger_radius_disconnect"].(bool) {
				configSet = append(configSet, setPrefix+"reconfigure trigger radius-disconnect")
			}
		}
	}
	if d.Get("remote_id_mismatch_disconnect").(bool) {
		configSet = append(configSet, setPrefix+"remote-id-mismatch disconnect")
	}
	if d.Get("route_suppression_access").(bool) {
		configSet = append(configSet, setPrefix+"route-suppression access")
	}
	if d.Get("route_suppression_access_internal").(bool) {
		configSet = append(configSet, setPrefix+"route-suppression access-internal")
	}
	if d.Get("route_suppression_destination").(bool) {
		configSet = append(configSet, setPrefix+"route-suppression destination")
	}
	if v := d.Get("service_profile").(string); v != "" {
		configSet = append(configSet, setPrefix+"service-profile \""+v+"\"")
	}
	if v := d.Get("short_cycle_protection_lockout_max_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-max-time "+strconv.Itoa(v))
	}
	if v := d.Get("short_cycle_protection_lockout_min_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-min-time "+strconv.Itoa(v))
	}

	return clt.configSet(configSet, junSess)
}

func setSystemServicesDhcpLocalServerGroupInterface(
	interFace map[string]interface{}, setPrefixInterface, version string,
) ([]string, error) {
	configSet := make([]string, 0)

	setPrefix := setPrefixInterface + "interface " + interFace["name"].(string) + " "

	configSet = append(configSet, setPrefix)
	if v := interFace["access_profile"].(string); v != "" {
		configSet = append(configSet, setPrefix+"access-profile \""+v+"\"")
	}
	if dynProfile := interFace["dynamic_profile"].(string); dynProfile != "" {
		configSet = append(configSet, setPrefix+"dynamic-profile \""+dynProfile+"\"")
		if interFace["dynamic_profile_use_primary"].(string) != "" &&
			interFace["dynamic_profile_aggregate_clients"].(bool) {
			return configSet, fmt.Errorf("conflict between "+
				"dynamic_profile_use_primary and dynamic_profile_aggregate_clients in interface %s", interFace["name"].(string))
		}
		if v := interFace["dynamic_profile_use_primary"].(string); v != "" {
			configSet = append(configSet, setPrefix+"dynamic-profile use-primary \""+v+"\"")
		}
		if interFace["dynamic_profile_aggregate_clients"].(bool) {
			configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients")
			if v := interFace["dynamic_profile_aggregate_clients_action"].(string); v != "" {
				configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients "+v)
			}
		} else if interFace["dynamic_profile_aggregate_clients_action"].(string) != "" {
			return configSet, fmt.Errorf("dynamic_profile_aggregate_clients need to be true with "+
				"dynamic_profile_aggregate_clients_action in interface %s", interFace["name"].(string))
		}
	} else if interFace["dynamic_profile_use_primary"].(string) != "" ||
		interFace["dynamic_profile_aggregate_clients"].(bool) ||
		interFace["dynamic_profile_aggregate_clients_action"].(string) != "" {
		return configSet, fmt.Errorf("dynamic_profile need to be set with "+
			"dynamic_profile_use_primary, dynamic_profile_aggregate_clients "+
			"or dynamic_profile_aggregate_clients_action in interface %s", interFace["name"].(string))
	}
	if interFace["exclude"].(bool) {
		configSet = append(configSet, setPrefix+"exclude")
	}
	for _, v := range interFace["overrides_v4"].([]interface{}) {
		if version == "v6" {
			return configSet, fmt.Errorf("overrides_v4 not compatible if version = v6")
		}
		configSetOverrides, err := setSystemServicesDhcpLocalServerGroupFamilyDhcpOverridesV4(
			v.(map[string]interface{}), setPrefix)
		if err != nil {
			return configSet, err
		}
		configSet = append(configSet, configSetOverrides...)
	}
	for _, v := range interFace["overrides_v6"].([]interface{}) {
		if version == "v4" {
			return configSet, fmt.Errorf("overrides_v6 not compatible if version = v4")
		}
		configSetOverrides, err := setSystemServicesDhcpLocalServerGroupFamilyDhcpOverridesV6(
			v.(map[string]interface{}), setPrefix)
		if err != nil {
			return configSet, err
		}
		configSet = append(configSet, configSetOverrides...)
	}
	if v := interFace["service_profile"].(string); v != "" {
		configSet = append(configSet, setPrefix+"service-profile \""+v+"\"")
	}
	if v := interFace["short_cycle_protection_lockout_max_time"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-max-time "+strconv.Itoa(v))
	}
	if v := interFace["short_cycle_protection_lockout_min_time"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-min-time "+strconv.Itoa(v))
	}
	if interFace["trace"].(bool) {
		configSet = append(configSet, setPrefix+"trace")
	}
	if v := interFace["upto"].(string); v != "" {
		configSet = append(configSet, setPrefix+"upto "+v)
	}

	return configSet, nil
}

func setSystemServicesDhcpLocalServerGroupFamilyDhcpOverridesV4(overrides map[string]interface{}, setPrefix string,
) ([]string, error) {
	configSet := make([]string, 0)
	setPrefix += "overrides "

	if overrides["allow_no_end_option"].(bool) {
		configSet = append(configSet, setPrefix+"allow-no-end-option")
	}
	if v := overrides["asymmetric_lease_time"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"asymmetric-lease-time "+strconv.Itoa(v))
	}
	if overrides["bootp_support"].(bool) {
		configSet = append(configSet, setPrefix+"bootp-support")
	}
	if v := overrides["client_discover_match"].(string); v != "" {
		configSet = append(configSet, setPrefix+"client-discover-match "+v)
	}
	for _, dobason := range overrides["delay_offer_based_on"].(*schema.Set).List() {
		delayOfferBasedOn := dobason.(map[string]interface{})
		configSet = append(configSet, setPrefix+"delay-offer based-on "+
			delayOfferBasedOn["option"].(string)+" "+delayOfferBasedOn["compare"].(string)+" "+
			delayOfferBasedOn["value_type"].(string)+" \""+delayOfferBasedOn["value"].(string)+"\"")
	}
	if len(overrides["delay_offer_based_on"].(*schema.Set).List()) == 0 &&
		overrides["delay_offer_delay_time"].(int) != 0 {
		return configSet, fmt.Errorf("delay_offer_based_on need to be set with delay_offer_delay_time")
	}
	if v := overrides["delay_offer_delay_time"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"delay-offer delay-time "+strconv.Itoa(v))
	}
	if overrides["delete_binding_on_renegotiation"].(bool) {
		configSet = append(configSet, setPrefix+"delete-binding-on-renegotiation")
	}
	if v := overrides["dual_stack"].(string); v != "" {
		configSet = append(configSet, setPrefix+"dual-stack \""+v+"\"")
	}
	if overrides["include_option_82_forcerenew"].(bool) {
		configSet = append(configSet, setPrefix+"include-option-82 forcerenew")
	}
	if overrides["include_option_82_nak"].(bool) {
		configSet = append(configSet, setPrefix+"include-option-82 nak")
	}
	if v := overrides["interface_client_limit"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"interface-client-limit "+strconv.Itoa(v))
	}
	if overrides["process_inform"].(bool) {
		configSet = append(configSet, setPrefix+"process-inform")
		if v := overrides["process_inform_pool"].(string); v != "" {
			configSet = append(configSet, setPrefix+"process-inform pool \""+v+"\"")
		}
	} else if overrides["process_inform_pool"].(string) != "" {
		return configSet, fmt.Errorf("process_inform need to be true with process_inform_pool")
	}
	if v := overrides["protocol_attributes"].(string); v != "" {
		configSet = append(configSet, setPrefix+"protocol-attributes \""+v+"\"")
	}

	if len(configSet) == 0 {
		return configSet, fmt.Errorf("an overrides_v4 block is empty")
	}

	return configSet, nil
}

func setSystemServicesDhcpLocalServerGroupFamilyDhcpOverridesV6(overrides map[string]interface{}, setPrefix string,
) ([]string, error) {
	configSet := make([]string, 0)
	setPrefix += "overrides "

	if overrides["always_add_option_dns_server"].(bool) {
		configSet = append(configSet, setPrefix+"always-add-option-dns-server")
	}
	if overrides["always_process_option_request_option"].(bool) {
		configSet = append(configSet, setPrefix+"always-process-option-request-option")
	}
	if v := overrides["asymmetric_lease_time"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"asymmetric-lease-time "+strconv.Itoa(v))
	}
	if v := overrides["asymmetric_prefix_lease_time"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"asymmetric-prefix-lease-time "+strconv.Itoa(v))
	}
	if overrides["client_negotiation_match_incoming_interface"].(bool) {
		configSet = append(configSet, setPrefix+"client-negotiation-match incoming-interface")
	}
	for _, dobason := range overrides["delay_advertise_based_on"].(*schema.Set).List() {
		delayOfferBasedOn := dobason.(map[string]interface{})
		configSet = append(configSet, setPrefix+"delay-advertise based-on "+
			delayOfferBasedOn["option"].(string)+" "+delayOfferBasedOn["compare"].(string)+" "+
			delayOfferBasedOn["value_type"].(string)+" \""+delayOfferBasedOn["value"].(string)+"\"")
	}
	if len(overrides["delay_advertise_based_on"].(*schema.Set).List()) == 0 &&
		overrides["delay_advertise_delay_time"].(int) != 0 {
		return configSet, fmt.Errorf("delay_offer_based_on need to be set with delay_offer_delay_time")
	}
	if v := overrides["delay_advertise_delay_time"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"delay-advertise delay-time "+strconv.Itoa(v))
	}
	if v := overrides["delegated_pool"].(string); v != "" {
		configSet = append(configSet, setPrefix+"delegated-pool \""+v+"\"")
	}
	if overrides["delete_binding_on_renegotiation"].(bool) {
		configSet = append(configSet, setPrefix+"delete-binding-on-renegotiation")
	}
	if v := overrides["dual_stack"].(string); v != "" {
		configSet = append(configSet, setPrefix+"dual-stack \""+v+"\"")
	}
	if v := overrides["interface_client_limit"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"interface-client-limit "+strconv.Itoa(v))
	}
	if overrides["multi_address_embedded_option_response"].(bool) {
		configSet = append(configSet, setPrefix+"multi-address-embedded-option-response")
	}
	if overrides["process_inform"].(bool) {
		configSet = append(configSet, setPrefix+"process-inform")
		if v := overrides["process_inform_pool"].(string); v != "" {
			configSet = append(configSet, setPrefix+"process-inform pool \""+v+"\"")
		}
	} else if overrides["process_inform_pool"].(string) != "" {
		return configSet, fmt.Errorf("process_inform need to be true with process_inform_pool")
	}
	if v := overrides["protocol_attributes"].(string); v != "" {
		configSet = append(configSet, setPrefix+"protocol-attributes \""+v+"\"")
	}
	if overrides["rapid_commit"].(bool) {
		configSet = append(configSet, setPrefix+"rapid-commit")
	}
	if overrides["top_level_status_code"].(bool) {
		configSet = append(configSet, setPrefix+"top-level-status-code")
	}

	if len(configSet) == 0 {
		return configSet, fmt.Errorf("an overrides_v6 block is empty")
	}

	return configSet, nil
}

func readSystemServicesDhcpLocalServerGroup(name, instance, version string, clt *Client, junSess *junosSession,
) (systemServicesDhcpLocalServerGroupOptions, error) {
	var confRead systemServicesDhcpLocalServerGroupOptions
	var showConfig string
	var err error
	showCmd := cmdShowConfig
	if instance != defaultW {
		showCmd += routingInstancesWS + instance + " "
	}
	showCmd += "system services dhcp-local-server "
	if version == "v6" {
		showCmd += "dhcpv6 group " + name
	} else {
		showCmd += "group " + name
	}

	showConfig, err = clt.command(showCmd+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}

	if showConfig != emptyW {
		confRead.name = name
		confRead.routingInstance = instance
		confRead.version = version
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLS)
			switch {
			case strings.HasPrefix(itemTrim, "access-profile "):
				confRead.accessProfile = strings.Trim(strings.TrimPrefix(itemTrim, "access-profile "), "\"")
			case strings.HasPrefix(itemTrim, "authentication password "):
				confRead.authenticationPassword = strings.Trim(strings.TrimPrefix(itemTrim, "authentication password "), "\"")
			case strings.HasPrefix(itemTrim, "authentication username-include "):
				if len(confRead.authenticationUsernameInclude) == 0 {
					confRead.authenticationUsernameInclude = append(confRead.authenticationUsernameInclude, map[string]interface{}{
						"circuit_type":              false,
						"client_id":                 false,
						"delimiter":                 "",
						"domain_name":               "",
						"interface_description":     "",
						"interface_name":            false,
						"mac_address":               false,
						"option_60":                 false,
						"option_82":                 false,
						"option_82_circuit_id":      false,
						"option_82_remote_id":       false,
						"relay_agent_interface_id":  false,
						"relay_agent_remote_id":     false,
						"relay_agent_subscriber_id": false,
						"routing_instance_name":     false,
						"user_prefix":               "",
						"vlan_tags":                 false,
					})
				}
				itemTrimAuthUserIncl := strings.TrimPrefix(itemTrim, "authentication username-include ")
				switch {
				case itemTrimAuthUserIncl == "circuit-type":
					confRead.authenticationUsernameInclude[0]["circuit_type"] = true
				case itemTrimAuthUserIncl == "client-id":
					confRead.authenticationUsernameInclude[0]["client_id"] = true
				case strings.HasPrefix(itemTrimAuthUserIncl, "delimiter "):
					confRead.authenticationUsernameInclude[0]["delimiter"] = strings.Trim(strings.TrimPrefix(
						itemTrimAuthUserIncl, "delimiter "), "\"")
				case strings.HasPrefix(itemTrimAuthUserIncl, "domain-name "):
					confRead.authenticationUsernameInclude[0]["domain_name"] = strings.Trim(strings.TrimPrefix(
						itemTrimAuthUserIncl, "domain-name "), "\"")
				case strings.HasPrefix(itemTrimAuthUserIncl, "interface-description "):
					confRead.authenticationUsernameInclude[0]["interface_description"] = strings.TrimPrefix(
						itemTrimAuthUserIncl, "interface-description ")
				case itemTrimAuthUserIncl == "interface-name":
					confRead.authenticationUsernameInclude[0]["interface_name"] = true
				case itemTrimAuthUserIncl == "mac-address":
					confRead.authenticationUsernameInclude[0]["mac_address"] = true
				case itemTrimAuthUserIncl == "option-60":
					confRead.authenticationUsernameInclude[0]["option_60"] = true
				case itemTrimAuthUserIncl == "option-82 circuit-id":
					confRead.authenticationUsernameInclude[0]["option_82_circuit_id"] = true
					confRead.authenticationUsernameInclude[0]["option_82"] = true
				case itemTrimAuthUserIncl == "option-82 remote-id":
					confRead.authenticationUsernameInclude[0]["option_82_remote_id"] = true
					confRead.authenticationUsernameInclude[0]["option_82"] = true
				case itemTrimAuthUserIncl == "option-82":
					confRead.authenticationUsernameInclude[0]["option_82"] = true
				case itemTrimAuthUserIncl == "relay-agent-interface-id":
					confRead.authenticationUsernameInclude[0]["relay_agent_interface_id"] = true
				case itemTrimAuthUserIncl == "relay-agent-remote-id":
					confRead.authenticationUsernameInclude[0]["relay_agent_remote_id"] = true
				case itemTrimAuthUserIncl == "relay-agent-subscriber-id":
					confRead.authenticationUsernameInclude[0]["relay_agent_subscriber_id"] = true
				case itemTrimAuthUserIncl == "routing-instance-name":
					confRead.authenticationUsernameInclude[0]["routing_instance_name"] = true
				case strings.HasPrefix(itemTrimAuthUserIncl, "user-prefix "):
					confRead.authenticationUsernameInclude[0]["user_prefix"] = strings.Trim(strings.TrimPrefix(
						itemTrimAuthUserIncl, "user-prefix "), "\"")
				case itemTrimAuthUserIncl == "vlan-tags":
					confRead.authenticationUsernameInclude[0]["vlan_tags"] = true
				}
			case strings.HasPrefix(itemTrim, "dynamic-profile "):
				switch {
				case strings.HasPrefix(itemTrim, "dynamic-profile use-primary "):
					confRead.dynamicProfileUsePrimary = strings.Trim(strings.TrimPrefix(
						itemTrim, "dynamic-profile use-primary "), "\"")
				case strings.HasPrefix(itemTrim, "dynamic-profile aggregate-clients"):
					confRead.dynamicProfileAggregateClients = true
					if strings.HasPrefix(itemTrim, "dynamic-profile aggregate-clients ") {
						confRead.dynamicProfileAggregateClientsAction = strings.TrimPrefix(
							itemTrim, "dynamic-profile aggregate-clients ")
					}
				default:
					confRead.dynamicProfile = strings.Trim(strings.TrimPrefix(itemTrim, "dynamic-profile "), "\"")
				}
			case strings.HasPrefix(itemTrim, "interface "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "interface "), " ")
				interFace := map[string]interface{}{
					"name":                                     itemTrimSplit[0],
					"access_profile":                           "",
					"dynamic_profile":                          "",
					"dynamic_profile_use_primary":              "",
					"dynamic_profile_aggregate_clients":        false,
					"dynamic_profile_aggregate_clients_action": "",
					"exclude":                                  false,
					"overrides_v4":                             make([]map[string]interface{}, 0),
					"overrides_v6":                             make([]map[string]interface{}, 0),
					"service_profile":                          "",
					"short_cycle_protection_lockout_max_time":  0,
					"short_cycle_protection_lockout_min_time":  0,
					"trace": false,
					"upto":  "",
				}
				confRead.interFace = copyAndRemoveItemMapList("name", interFace, confRead.interFace)
				itemTrimInterface := strings.TrimPrefix(itemTrim, "interface "+itemTrimSplit[0])
				if strings.HasPrefix(itemTrimInterface, " ") {
					if err := readSystemServicesDhcpLocalServerGroupInterface(
						strings.TrimPrefix(itemTrimInterface, " "), version, interFace); err != nil {
						return confRead, err
					}
				}
				confRead.interFace = append(confRead.interFace, interFace)
			case strings.HasPrefix(itemTrim, "lease-time-validation"):
				if len(confRead.leaseTimeValidation) == 0 {
					confRead.leaseTimeValidation = append(confRead.leaseTimeValidation, map[string]interface{}{
						"lease_time_threshold": 0,
						"violation_action":     "",
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "lease-time-validation lease-time-threshold "):
					var err error
					confRead.leaseTimeValidation[0]["lease_time_threshold"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "lease-time-validation lease-time-threshold "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "lease-time-validation violation-action "):
					confRead.leaseTimeValidation[0]["violation_action"] = strings.TrimPrefix(
						itemTrim, "lease-time-validation violation-action ")
				}
			case strings.HasPrefix(itemTrim, "liveness-detection failure-action "):
				confRead.livenessDetectionFailureAction = strings.TrimPrefix(itemTrim, "liveness-detection failure-action ")
			case strings.HasPrefix(itemTrim, "liveness-detection method bfd "):
				if len(confRead.livenessDetectionMethodBfd) == 0 {
					confRead.livenessDetectionMethodBfd = append(confRead.livenessDetectionMethodBfd, map[string]interface{}{
						"detection_time_threshold":    -1,
						"holddown_interval":           -1,
						"minimum_interval":            0,
						"minimum_receive_interval":    0,
						"multiplier":                  0,
						"no_adaptation":               false,
						"session_mode":                "",
						"transmit_interval_minimum":   0,
						"transmit_interval_threshold": -1,
						"version":                     "",
					})
				}
				itemTrimLiveDetMethBfd := strings.TrimPrefix(itemTrim, "liveness-detection method bfd ")
				var err error
				switch {
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "detection-time threshold "):
					confRead.livenessDetectionMethodBfd[0]["detection_time_threshold"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLiveDetMethBfd, "detection-time threshold "))
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "holddown-interval "):
					confRead.livenessDetectionMethodBfd[0]["holddown_interval"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLiveDetMethBfd, "holddown-interval "))
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "minimum-interval "):
					confRead.livenessDetectionMethodBfd[0]["minimum_interval"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLiveDetMethBfd, "minimum-interval "))
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "minimum-receive-interval "):
					confRead.livenessDetectionMethodBfd[0]["minimum_receive_interval"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLiveDetMethBfd, "minimum-receive-interval "))
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "multiplier "):
					confRead.livenessDetectionMethodBfd[0]["multiplier"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLiveDetMethBfd, "multiplier "))
				case itemTrimLiveDetMethBfd == "no-adaptation":
					confRead.livenessDetectionMethodBfd[0]["no_adaptation"] = true
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "session-mode "):
					confRead.livenessDetectionMethodBfd[0]["session_mode"] = strings.TrimPrefix(
						itemTrimLiveDetMethBfd, "session-mode ")
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "transmit-interval minimum-interval "):
					confRead.livenessDetectionMethodBfd[0]["transmit_interval_minimum"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLiveDetMethBfd, "transmit-interval minimum-interval "))
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "transmit-interval threshold "):
					confRead.livenessDetectionMethodBfd[0]["transmit_interval_threshold"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrimLiveDetMethBfd, "transmit-interval threshold "))
				case strings.HasPrefix(itemTrimLiveDetMethBfd, "version "):
					confRead.livenessDetectionMethodBfd[0]["version"] = strings.TrimPrefix(itemTrimLiveDetMethBfd, "version ")
				}
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "liveness-detection method layer2-liveness-detection "):
				if len(confRead.livenessDetectionMethodLayer2) == 0 {
					confRead.livenessDetectionMethodLayer2 = append(confRead.livenessDetectionMethodLayer2, map[string]interface{}{
						"max_consecutive_retries": 0,
						"transmit_interval":       0,
					})
				}
				var err error
				switch {
				case strings.HasPrefix(itemTrim, "liveness-detection method layer2-liveness-detection max-consecutive-retries "):
					confRead.livenessDetectionMethodLayer2[0]["max_consecutive_retries"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "liveness-detection method layer2-liveness-detection max-consecutive-retries "))
				case strings.HasPrefix(itemTrim, "liveness-detection method layer2-liveness-detection transmit-interval "):
					confRead.livenessDetectionMethodLayer2[0]["transmit_interval"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "liveness-detection method layer2-liveness-detection transmit-interval "))
				}
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "overrides "):
				if version == "v4" {
					if len(confRead.overridesV4) == 0 {
						confRead.overridesV4 = append(confRead.overridesV4, genSystemServicesDhcpLocalServerGroupOverridesV4())
					}
					if err := readSystemServicesDhcpLocalServerGroupOverridesV4(
						strings.TrimPrefix(itemTrim, "overrides "),
						confRead.overridesV4[0],
					); err != nil {
						return confRead, err
					}
				} else if version == "v6" {
					if len(confRead.overridesV6) == 0 {
						confRead.overridesV6 = append(confRead.overridesV6, genSystemServicesDhcpLocalServerGroupOverridesV6())
					}
					if err := readSystemServicesDhcpLocalServerGroupOverridesV6(
						strings.TrimPrefix(itemTrim, "overrides "),
						confRead.overridesV6[0],
					); err != nil {
						return confRead, err
					}
				}
			case itemTrim == "reauthenticate lease-renewal":
				confRead.reauthenticateLeaseRenewal = true
			case itemTrim == "reauthenticate remote-id-mismatch":
				confRead.reauthenticateRemoteIDMismatch = true
			case strings.HasPrefix(itemTrim, "reconfigure"):
				if len(confRead.reconfigure) == 0 {
					confRead.reconfigure = append(confRead.reconfigure, map[string]interface{}{
						"attempts":                  0,
						"clear_on_abort":            false,
						"support_option_pd_exclude": false,
						"timeout":                   0,
						"token":                     "",
						"trigger_radius_disconnect": false,
					})
				}
				if strings.HasPrefix(itemTrim, "reconfigure ") {
					var err error
					switch {
					case strings.HasPrefix(itemTrim, "reconfigure attempts "):
						confRead.reconfigure[0]["attempts"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "reconfigure attempts "))
					case itemTrim == "reconfigure clear-on-abort":
						confRead.reconfigure[0]["clear_on_abort"] = true
					case itemTrim == "reconfigure support-option-pd-exclude":
						confRead.reconfigure[0]["support_option_pd_exclude"] = true
					case strings.HasPrefix(itemTrim, "reconfigure timeout "):
						confRead.reconfigure[0]["timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "reconfigure timeout "))
					case strings.HasPrefix(itemTrim, "reconfigure token "):
						confRead.reconfigure[0]["token"] = strings.Trim(strings.TrimPrefix(itemTrim, "reconfigure token "), "\"")
					case itemTrim == "reconfigure trigger radius-disconnect":
						confRead.reconfigure[0]["trigger_radius_disconnect"] = true
					}
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				}
			case itemTrim == "remote-id-mismatch disconnect":
				confRead.remoteIDMismatchDisconnect = true
			case itemTrim == "route-suppression access":
				confRead.routeSuppressionAccess = true
			case itemTrim == "route-suppression access-internal":
				confRead.routeSuppressionAccessInternal = true
			case itemTrim == "route-suppression destination":
				confRead.routeSuppressionDestination = true
			case strings.HasPrefix(itemTrim, "service-profile "):
				confRead.serviceProfile = strings.Trim(strings.TrimPrefix(itemTrim, "service-profile "), "\"")
			case strings.HasPrefix(itemTrim, "short-cycle-protection lockout-max-time "):
				var err error
				confRead.shortCycleProtectionLockoutMaxTime, err = strconv.Atoi(strings.TrimPrefix(
					itemTrim, "short-cycle-protection lockout-max-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "short-cycle-protection lockout-min-time "):
				var err error
				confRead.shortCycleProtectionLockoutMinTime, err = strconv.Atoi(strings.TrimPrefix(
					itemTrim, "short-cycle-protection lockout-min-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		}
	}

	return confRead, nil
}

func readSystemServicesDhcpLocalServerGroupInterface(itemTrim, version string, interFace map[string]interface{}) error {
	var err error
	switch {
	case strings.HasPrefix(itemTrim, "access-profile "):
		interFace["access_profile"] = strings.Trim(strings.TrimPrefix(itemTrim, "access-profile "), "\"")
	case strings.HasPrefix(itemTrim, "dynamic-profile "):
		switch {
		case strings.HasPrefix(itemTrim, "dynamic-profile use-primary "):
			interFace["dynamic_profile_use_primary"] = strings.Trim(strings.TrimPrefix(
				itemTrim, "dynamic-profile use-primary "), "\"")
		case strings.HasPrefix(itemTrim, "dynamic-profile aggregate-clients"):
			interFace["dynamic_profile_aggregate_clients"] = true
			if strings.HasPrefix(itemTrim, "dynamic-profile aggregate-clients ") {
				interFace["dynamic_profile_aggregate_clients_action"] = strings.TrimPrefix(
					itemTrim, "dynamic-profile aggregate-clients ")
			}
		default:
			interFace["dynamic_profile"] = strings.Trim(strings.TrimPrefix(itemTrim, "dynamic-profile "), "\"")
		}
	case itemTrim == "exclude":
		interFace["exclude"] = true
	case strings.HasPrefix(itemTrim, "overrides "):
		if version == "v4" {
			if len(interFace["overrides_v4"].([]map[string]interface{})) == 0 {
				interFace["overrides_v4"] = append(
					interFace["overrides_v4"].([]map[string]interface{}),
					genSystemServicesDhcpLocalServerGroupOverridesV4(),
				)
			}
			if err := readSystemServicesDhcpLocalServerGroupOverridesV4(
				strings.TrimPrefix(itemTrim, "overrides "),
				interFace["overrides_v4"].([]map[string]interface{})[0],
			); err != nil {
				return err
			}
		} else if version == "v6" {
			if len(interFace["overrides_v6"].([]map[string]interface{})) == 0 {
				interFace["overrides_v6"] = append(
					interFace["overrides_v6"].([]map[string]interface{}),
					genSystemServicesDhcpLocalServerGroupOverridesV6(),
				)
			}
			if err := readSystemServicesDhcpLocalServerGroupOverridesV6(
				strings.TrimPrefix(itemTrim, "overrides "),
				interFace["overrides_v6"].([]map[string]interface{})[0],
			); err != nil {
				return err
			}
		}
	case strings.HasPrefix(itemTrim, "service-profile "):
		interFace["service_profile"] = strings.Trim(strings.TrimPrefix(itemTrim, "service-profile "), "\"")
	case strings.HasPrefix(itemTrim, "short-cycle-protection lockout-max-time "):
		interFace["short_cycle_protection_lockout_max_time"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "short-cycle-protection lockout-max-time "))
	case strings.HasPrefix(itemTrim, "short-cycle-protection lockout-min-time "):
		interFace["short_cycle_protection_lockout_min_time"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "short-cycle-protection lockout-min-time "))
	case itemTrim == "trace":
		interFace["trace"] = true
	case strings.HasPrefix(itemTrim, "upto "):
		interFace["upto"] = strings.TrimPrefix(itemTrim, "upto ")
	}
	if err != nil {
		return fmt.Errorf(failedConvAtoiError, itemTrim, err)
	}

	return nil
}

func genSystemServicesDhcpLocalServerGroupOverridesV4() map[string]interface{} {
	return map[string]interface{}{
		"allow_no_end_option":             false,
		"asymmetric_lease_time":           0,
		"bootp_support":                   false,
		"client_discover_match":           "",
		"delay_offer_based_on":            make([]map[string]interface{}, 0),
		"delay_offer_delay_time":          0,
		"delete_binding_on_renegotiation": false,
		"dual_stack":                      "",
		"include_option_82_forcerenew":    false,
		"include_option_82_nak":           false,
		"interface_client_limit":          0,
		"process_inform":                  false,
		"process_inform_pool":             "",
		"protocol_attributes":             "",
	}
}

func genSystemServicesDhcpLocalServerGroupOverridesV6() map[string]interface{} {
	return map[string]interface{}{
		"always_add_option_dns_server":                false,
		"always_process_option_request_option":        false,
		"asymmetric_lease_time":                       0,
		"asymmetric_prefix_lease_time":                0,
		"client_negotiation_match_incoming_interface": false,
		"delay_advertise_based_on":                    make([]map[string]interface{}, 0),
		"delay_advertise_delay_time":                  0,
		"delegated_pool":                              "",
		"delete_binding_on_renegotiation":             false,
		"dual_stack":                                  "",
		"interface_client_limit":                      0,
		"multi_address_embedded_option_response":      false,
		"process_inform":                              false,
		"process_inform_pool":                         "",
		"protocol_attributes":                         "",
		"rapid_commit":                                false,
		"top_level_status_code":                       false,
	}
}

func readSystemServicesDhcpLocalServerGroupOverridesV4(itemTrim string, overrides map[string]interface{}) error {
	var err error
	switch {
	case itemTrim == "allow-no-end-option":
		overrides["allow_no_end_option"] = true
	case strings.HasPrefix(itemTrim, "asymmetric-lease-time "):
		overrides["asymmetric_lease_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "asymmetric-lease-time "))
	case itemTrim == "bootp-support":
		overrides["bootp_support"] = true
	case strings.HasPrefix(itemTrim, "client-discover-match "):
		overrides["client_discover_match"] = strings.TrimPrefix(itemTrim, "client-discover-match ")
	case strings.HasPrefix(itemTrim, "delay-offer based-on "):
		itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "delay-offer based-on "), " ")
		if len(itemTrimSplit) < 4 {
			return fmt.Errorf("can't read line delay-offer based-on, not enough element")
		}
		overrides["delay_offer_based_on"] = append(
			overrides["delay_offer_based_on"].([]map[string]interface{}),
			map[string]interface{}{
				"option":     itemTrimSplit[0],
				"compare":    itemTrimSplit[1],
				"value_type": itemTrimSplit[2],
				"value": strings.Trim(strings.TrimPrefix(itemTrim,
					"delay-offer based-on "+itemTrimSplit[0]+" "+itemTrimSplit[1]+" "+itemTrimSplit[2]+" "), "\""),
			})
	case strings.HasPrefix(itemTrim, "delay-offer delay-time "):
		overrides["delay_offer_delay_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "delay-offer delay-time "))
	case itemTrim == "delete-binding-on-renegotiation":
		overrides["delete_binding_on_renegotiation"] = true
	case strings.HasPrefix(itemTrim, "dual-stack "):
		overrides["dual_stack"] = strings.Trim(strings.TrimPrefix(itemTrim, "dual-stack "), "\"")
	case itemTrim == "include-option-82 forcerenew":
		overrides["include_option_82_forcerenew"] = true
	case itemTrim == "include-option-82 nak":
		overrides["include_option_82_nak"] = true
	case strings.HasPrefix(itemTrim, "interface-client-limit "):
		overrides["interface_client_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "interface-client-limit "))
	case strings.HasPrefix(itemTrim, "process-inform pool "):
		overrides["process_inform_pool"] = strings.Trim(strings.TrimPrefix(itemTrim, "process-inform pool "), "\"")
		overrides["process_inform"] = true
	case itemTrim == "process-inform":
		overrides["process_inform"] = true
	case strings.HasPrefix(itemTrim, "protocol-attributes "):
		overrides["protocol_attributes"] = strings.Trim(strings.TrimPrefix(itemTrim, "protocol-attributes "), "\"")
	}
	if err != nil {
		return fmt.Errorf(failedConvAtoiError, itemTrim, err)
	}

	return nil
}

func readSystemServicesDhcpLocalServerGroupOverridesV6(itemTrim string, overrides map[string]interface{}) error {
	var err error
	switch {
	case itemTrim == "always-add-option-dns-server":
		overrides["always_add_option_dns_server"] = true
	case itemTrim == "always-process-option-request-option":
		overrides["always_process_option_request_option"] = true
	case strings.HasPrefix(itemTrim, "asymmetric-lease-time "):
		overrides["asymmetric_lease_time"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "asymmetric-lease-time "))
	case strings.HasPrefix(itemTrim, "asymmetric-prefix-lease-time "):
		overrides["asymmetric_prefix_lease_time"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "asymmetric-prefix-lease-time "))
	case itemTrim == "client-negotiation-match incoming-interface":
		overrides["client_negotiation_match_incoming_interface"] = true
	case strings.HasPrefix(itemTrim, "delay-advertise based-on "):
		itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "delay-advertise based-on "), " ")
		if len(itemTrimSplit) < 4 {
			return fmt.Errorf("can't read line delay-advertise based-on, not enough element")
		}
		overrides["delay_advertise_based_on"] = append(
			overrides["delay_advertise_based_on"].([]map[string]interface{}),
			map[string]interface{}{
				"option":     itemTrimSplit[0],
				"compare":    itemTrimSplit[1],
				"value_type": itemTrimSplit[2],
				"value": strings.Trim(strings.TrimPrefix(itemTrim,
					"delay-advertise based-on "+itemTrimSplit[0]+" "+itemTrimSplit[1]+" "+itemTrimSplit[2]+" "), "\""),
			})
	case strings.HasPrefix(itemTrim, "delay-advertise delay-time "):
		overrides["delay_advertise_delay_time"], err = strconv.Atoi(strings.TrimPrefix(
			itemTrim, "delay-advertise delay-time "))
	case strings.HasPrefix(itemTrim, "delegated-pool "):
		overrides["delegated_pool"] = strings.Trim(strings.TrimPrefix(itemTrim, "delegated-pool "), "\"")
	case itemTrim == "delete-binding-on-renegotiation":
		overrides["delete_binding_on_renegotiation"] = true
	case strings.HasPrefix(itemTrim, "dual-stack "):
		overrides["dual_stack"] = strings.Trim(strings.TrimPrefix(itemTrim, "dual-stack "), "\"")
	case strings.HasPrefix(itemTrim, "interface-client-limit "):
		overrides["interface_client_limit"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "interface-client-limit "))
	case itemTrim == "multi-address-embedded-option-response":
		overrides["multi_address_embedded_option_response"] = true
	case strings.HasPrefix(itemTrim, "process-inform pool "):
		overrides["process_inform_pool"] = strings.Trim(strings.TrimPrefix(itemTrim, "process-inform pool "), "\"")
		overrides["process_inform"] = true
	case itemTrim == "process-inform":
		overrides["process_inform"] = true
	case strings.HasPrefix(itemTrim, "protocol-attributes "):
		overrides["protocol_attributes"] = strings.Trim(strings.TrimPrefix(itemTrim, "protocol-attributes "), "\"")
	case itemTrim == "rapid-commit":
		overrides["rapid_commit"] = true
	case itemTrim == "top-level-status-code":
		overrides["top_level_status_code"] = true
	}
	if err != nil {
		return fmt.Errorf(failedConvAtoiError, itemTrim, err)
	}

	return nil
}

func delSystemServicesDhcpLocalServerGroup(name, instance, version string, clt *Client, junSess *junosSession,
) error {
	configSet := make([]string, 0, 1)
	switch {
	case instance == defaultW && version == "v6":
		configSet = append(configSet, "delete system services dhcp-local-server dhcpv6 group "+name)
	case instance == defaultW && version == "v4":
		configSet = append(configSet, "delete system services dhcp-local-server group "+name)
	case instance != defaultW && version == "v6":
		configSet = append(configSet, delRoutingInstances+instance+" "+
			"system services dhcp-local-server dhcpv6 group "+name)
	case instance != defaultW && version == "v4":
		configSet = append(configSet, delRoutingInstances+instance+" "+
			"system services dhcp-local-server group "+name)
	}

	return clt.configSet(configSet, junSess)
}

func fillSystemServicesDhcpLocalServerGroupData(
	d *schema.ResourceData, systemServicesDhcpLocalServerGroupOptions systemServicesDhcpLocalServerGroupOptions,
) {
	if tfErr := d.Set("name",
		systemServicesDhcpLocalServerGroupOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("routing_instance",
		systemServicesDhcpLocalServerGroupOptions.routingInstance); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("version",
		systemServicesDhcpLocalServerGroupOptions.version); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("access_profile",
		systemServicesDhcpLocalServerGroupOptions.accessProfile); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_password",
		systemServicesDhcpLocalServerGroupOptions.authenticationPassword); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("authentication_username_include",
		systemServicesDhcpLocalServerGroupOptions.authenticationUsernameInclude); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_profile",
		systemServicesDhcpLocalServerGroupOptions.dynamicProfile); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_profile_use_primary",
		systemServicesDhcpLocalServerGroupOptions.dynamicProfileUsePrimary); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_profile_aggregate_clients",
		systemServicesDhcpLocalServerGroupOptions.dynamicProfileAggregateClients); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dynamic_profile_aggregate_clients_action",
		systemServicesDhcpLocalServerGroupOptions.dynamicProfileAggregateClientsAction); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("interface",
		systemServicesDhcpLocalServerGroupOptions.interFace); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("lease_time_validation",
		systemServicesDhcpLocalServerGroupOptions.leaseTimeValidation); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("liveness_detection_failure_action",
		systemServicesDhcpLocalServerGroupOptions.livenessDetectionFailureAction); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("liveness_detection_method_bfd",
		systemServicesDhcpLocalServerGroupOptions.livenessDetectionMethodBfd); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("liveness_detection_method_layer2",
		systemServicesDhcpLocalServerGroupOptions.livenessDetectionMethodLayer2); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("overrides_v4",
		systemServicesDhcpLocalServerGroupOptions.overridesV4); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("overrides_v6",
		systemServicesDhcpLocalServerGroupOptions.overridesV6); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("reauthenticate_lease_renewal",
		systemServicesDhcpLocalServerGroupOptions.reauthenticateLeaseRenewal); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("reauthenticate_remote_id_mismatch",
		systemServicesDhcpLocalServerGroupOptions.reauthenticateRemoteIDMismatch); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("reconfigure",
		systemServicesDhcpLocalServerGroupOptions.reconfigure); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("remote_id_mismatch_disconnect",
		systemServicesDhcpLocalServerGroupOptions.remoteIDMismatchDisconnect); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("route_suppression_access",
		systemServicesDhcpLocalServerGroupOptions.routeSuppressionAccess); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("route_suppression_access_internal",
		systemServicesDhcpLocalServerGroupOptions.routeSuppressionAccessInternal); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("route_suppression_destination",
		systemServicesDhcpLocalServerGroupOptions.routeSuppressionDestination); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("service_profile",
		systemServicesDhcpLocalServerGroupOptions.serviceProfile); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("short_cycle_protection_lockout_max_time",
		systemServicesDhcpLocalServerGroupOptions.shortCycleProtectionLockoutMaxTime); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("short_cycle_protection_lockout_min_time",
		systemServicesDhcpLocalServerGroupOptions.shortCycleProtectionLockoutMinTime); tfErr != nil {
		panic(tfErr)
	}
}
