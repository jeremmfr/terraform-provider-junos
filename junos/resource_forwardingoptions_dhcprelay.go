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

type fwdOptsDhcpRelayOptions struct {
	activeServerGroupAllowServerChange   bool
	arpInspection                        bool
	duplicateClientsIncomingInterface    bool
	dynamicProfileAggregateClients       bool
	excludeRelayAgentIdentifier          bool
	forwardOnly                          bool
	forwardOnlyReplies                   bool
	noSnoop                              bool
	persistentStorageAutomatic           bool
	relayAgentOption79                   bool
	remoteIDMismatchDisconnect           bool
	routeSuppressionAccess               bool
	routeSuppressionAccessInternal       bool
	routeSuppressionDestination          bool
	sourceIPChange                       bool
	vendorSpecificInformationHostName    bool
	vendorSpecificInformationLocation    bool
	clientResponseTTL                    int
	maximumHopCount                      int
	minimumWaitTime                      int
	serverResponseTime                   int
	shortCycleProtectionLockoutMaxTime   int
	shortCycleProtectionLockoutMinTime   int
	accessProfile                        string
	activeServerGroup                    string
	authenticationPassword               string
	duplicateClientsInSubnet             string
	dynamicProfile                       string
	dynamicProfileAggregateClientsAction string
	dynamicProfileUsePrimary             string
	forwardOnlyRoutingInstance           string
	forwardSnoopedClients                string
	livenessDetectionFailureAction       string
	routingInstance                      string
	serverMatchDefaultAction             string
	serviceProfile                       string
	version                              string
	activeLeasequery                     []map[string]interface{}
	authenticationUsernameInclude        []map[string]interface{}
	bulkLeasequery                       []map[string]interface{}
	leaseTimeValidation                  []map[string]interface{}
	leasequery                           []map[string]interface{}
	livenessDetectionMethodBfd           []map[string]interface{}
	livenessDetectionMethodLayer2        []map[string]interface{}
	overridesV4                          []map[string]interface{}
	overridesV6                          []map[string]interface{}
	relayAgentInterfaceID                []map[string]interface{}
	relayAgentRemoteID                   []map[string]interface{}
	relayOption                          []map[string]interface{}
	relayOption82                        []map[string]interface{}
	serverMatchAddress                   []map[string]interface{}
	serverMatchDuid                      []map[string]interface{}
}

func resourceForwardingOptionsDhcpRelay() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceForwardingOptionsDhcpRelayCreate,
		ReadWithoutTimeout:   resourceForwardingOptionsDhcpRelayRead,
		UpdateWithoutTimeout: resourceForwardingOptionsDhcpRelayUpdate,
		DeleteWithoutTimeout: resourceForwardingOptionsDhcpRelayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceForwardingOptionsDhcpRelayImport,
		},
		Schema: map[string]*schema.Schema{
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
			"active_leasequery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idle_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 3600),
						},
						"peer_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 3600),
						},
						"topology_discover": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"active_server_group": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"active_server_group_allow_server_change": { // only dhcpv4
				Type:     schema.TypeBool,
				Optional: true,
			},
			"arp_inspection": { // only dhcpv4
				Type:     schema.TypeBool,
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
					Schema: schemaForwardingOptionsDhcpRelayAuthUsernameInclude(),
				},
			},
			"bulk_leasequery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attempts": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 720),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"trigger_automatic": { // only dhcpv6
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"client_response_ttl": { // only dhcpv4
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"duplicate_clients_in_subnet": { // only dhcpv4
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"incoming-interface", "option-82"}, false),
			},
			"duplicate_clients_incoming_interface": { // only dhcpv6
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dynamic_profile": {
				Type:     schema.TypeString,
				Optional: true,
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
			"dynamic_profile_use_primary": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"dynamic_profile"},
			},
			"exclude_relay_agent_identifier": { // only dhcpv6
				Type:     schema.TypeBool,
				Optional: true,
			},
			"forward_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"forward_only_replies": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"forward_only_routing_instance": {
				Type:             schema.TypeString,
				Optional:         true,
				RequiredWith:     []string{"forward_only"},
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"forward_snooped_clients": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"all-interfaces",
					"configured-interfaces",
					"non-configured-interfaces",
				}, false),
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
						"violation_action_drop": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"leasequery": {
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
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10),
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
			"maximum_hop_count": { // only dhcpv4
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 16),
			},
			"minimum_wait_time": { // only dhcpv4
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 30000),
			},
			"no_snoop": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"overrides_v4": { // only dhcpv4
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaForwardingOptionsDhcpRelayOverridesV4(),
				},
			},
			"overrides_v6": { // only dhcpv6
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaForwardingOptionsDhcpRelayOverridesV6(),
				},
			},
			"persistent_storage_automatic": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"relay_agent_interface_id": { // only dhcpv6
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaForwardingOptionsDhcpRelayAgentID(true),
				},
			},
			"relay_agent_option_79": { // only dhcpv6
				Type:     schema.TypeBool,
				Optional: true,
			},
			"relay_agent_remote_id": { // only dhcpv6
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaForwardingOptionsDhcpRelayAgentID(false),
				},
			},
			"relay_option": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaForwardingOptionsDhcpRelayOption(),
				},
			},
			"relay_option_82": { // only dhcpv4
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: schemaForwardingOptionsDhcpRelayOption82(),
				},
			},
			"remote_id_mismatch_disconnect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"route_suppression_access": { // only dhcpv6
				Type:     schema.TypeBool,
				Optional: true,
			},
			"route_suppression_access_internal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"route_suppression_destination": { // only dhcpv4
				Type:     schema.TypeBool,
				Optional: true,
			},
			"server_match_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"create-relay-entry", "forward-only"}, false),
						},
					},
				},
			},
			"server_match_default_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"create-relay-entry", "forward-only"}, false),
			},
			"server_match_duid": { // only dhcpv6
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compare": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"equals", "starts-with"}, false),
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
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"create-relay-entry", "forward-only"}, false),
						},
					},
				},
			},
			"server_response_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(0, 4294967295),
			},
			"service_profile": { // only dhcpv6
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
			"source_ip_change": { // only dhcpv4
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vendor_specific_information_host_name": { // only dhcpv6
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vendor_specific_information_location": { // only dhcpv6
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceForwardingOptionsDhcpRelayCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	routingInstanceArg := d.Get("routing_instance").(string)
	versionArg := d.Get("version").(string)
	clt := m.(*Client)
	if clt.fakeCreateSetFile != "" {
		if err := setForwardingOptionsDhcpRelay(d, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(routingInstanceArg + idSeparator + versionArg)

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
	if routingInstanceArg != defaultW {
		instanceExists, err := checkRoutingInstanceExists(routingInstanceArg, clt, junSess)
		if err != nil {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns, diag.FromErr(err)...)
		}
		if !instanceExists {
			appendDiagWarns(&diagWarns, clt.configClear(junSess))

			return append(diagWarns,
				diag.FromErr(fmt.Errorf("routing instance %v doesn't exist", routingInstanceArg))...)
		}
	}
	if err := setForwardingOptionsDhcpRelay(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	warns, err := clt.commitConf("create resource junos_forwardingoptions_dhcprelay", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.SetId(routingInstanceArg + idSeparator + versionArg)

	return append(diagWarns, resourceForwardingOptionsDhcpRelayReadWJunSess(d, clt, junSess)...)
}

func resourceForwardingOptionsDhcpRelayRead(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer clt.closeSession(junSess)

	return resourceForwardingOptionsDhcpRelayReadWJunSess(d, clt, junSess)
}

func resourceForwardingOptionsDhcpRelayReadWJunSess(
	d *schema.ResourceData, clt *Client, junSess *junosSession,
) diag.Diagnostics {
	mutex.Lock()
	if d.Get("routing_instance").(string) != defaultW {
		instanceExists, err := checkRoutingInstanceExists(d.Get("routing_instance").(string), clt, junSess)
		if err != nil {
			mutex.Unlock()

			return diag.FromErr(err)
		}
		if !instanceExists {
			mutex.Unlock()
			d.SetId("")

			return nil
		}
	}
	fwdOptsDhcpRelayOptions, err := readForwardingOptionsDhcpRelay(
		d.Get("routing_instance").(string),
		d.Get("version").(string),
		clt, junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	fillForwardingOptionsDhcpRelayData(d, fwdOptsDhcpRelayOptions)

	return nil
}

func resourceForwardingOptionsDhcpRelayUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	routingInstanceArg := d.Get("routing_instance").(string)
	versionArg := d.Get("version").(string)
	clt := m.(*Client)
	if clt.fakeUpdateAlso {
		if err := delForwardingOptionsDhcpRelay(routingInstanceArg, versionArg, clt, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := setForwardingOptionsDhcpRelay(d, clt, nil); err != nil {
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
	if err := delForwardingOptionsDhcpRelay(routingInstanceArg, versionArg, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setForwardingOptionsDhcpRelay(d, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	warns, err := clt.commitConf("update resource junos_forwardingoptions_dhcprelay", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceForwardingOptionsDhcpRelayReadWJunSess(d, clt, junSess)...)
}

func resourceForwardingOptionsDhcpRelayDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	routingInstanceArg := d.Get("routing_instance").(string)
	versionArg := d.Get("version").(string)
	clt := m.(*Client)
	if clt.fakeDeleteAlso {
		if err := delForwardingOptionsDhcpRelay(routingInstanceArg, versionArg, clt, nil); err != nil {
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
	if err := delForwardingOptionsDhcpRelay(routingInstanceArg, versionArg, clt, junSess); err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := clt.commitConf("delete resource junos_forwardingoptions_dhcprelay", junSess)
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, clt.configClear(junSess))

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceForwardingOptionsDhcpRelayImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*Client)
	junSess, err := clt.startNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.closeSession(junSess)
	result := make([]*schema.ResourceData, 1)
	idSplit := strings.Split(d.Id(), idSeparator)
	if len(idSplit) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", idSeparator)
	}
	if idSplit[1] != "v4" && idSplit[1] != "v6" {
		return nil, fmt.Errorf("bad version '%s' in id, need to be 'v4' or 'v6' (id must be "+
			"<routing_instance>"+idSeparator+"<version>)", idSplit[2])
	}
	if idSplit[0] != defaultW {
		routingInstanceExists, err := checkRoutingInstanceExists(idSplit[0], clt, junSess)
		if err != nil {
			return nil, err
		}
		if !routingInstanceExists {
			return nil, fmt.Errorf("don't find routing_instance with id '%v' (id must be "+
				"<routing_instance>"+idSeparator+"<version>)", d.Id())
		}
	}
	fwdOptsDhcpRelayOptions, err := readForwardingOptionsDhcpRelay(idSplit[0], idSplit[1], clt, junSess)
	if err != nil {
		return nil, err
	}
	fillForwardingOptionsDhcpRelayData(d, fwdOptsDhcpRelayOptions)

	result[0] = d

	return result, nil
}

func setForwardingOptionsDhcpRelay( //nolint: gocognit, gocyclo
	d *schema.ResourceData, clt *Client, junSess *junosSession,
) error {
	configSet := make([]string, 0)

	setPrefix := setLS
	if d.Get("routing_instance").(string) != defaultW {
		setPrefix = setRoutingInstances + d.Get("routing_instance").(string) + " "
	}
	if d.Get("version").(string) == "v6" {
		setPrefix += "forwarding-options dhcp-relay dhcpv6 "
	} else {
		setPrefix += "forwarding-options dhcp-relay "
	}

	if v := d.Get("access_profile").(string); v != "" {
		configSet = append(configSet, setPrefix+"access-profile \""+v+"\"")
	}
	for _, block := range d.Get("active_leasequery").([]interface{}) {
		configSet = append(configSet, setPrefix+"active-leasequery")
		if block != nil {
			activeLeasequery := block.(map[string]interface{})
			if v := activeLeasequery["idle_timeout"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"active-leasequery idle-timeout "+strconv.Itoa(v))
			}
			if v := activeLeasequery["peer_address"].(string); v != "" {
				configSet = append(configSet, setPrefix+"active-leasequery peer-address "+v)
			}
			if v := activeLeasequery["timeout"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"active-leasequery timeout "+strconv.Itoa(v))
			}
			if activeLeasequery["topology_discover"].(bool) {
				configSet = append(configSet, setPrefix+"active-leasequery topology-discover")
			}
		}
	}
	if v := d.Get("active_server_group").(string); v != "" {
		configSet = append(configSet, setPrefix+"active-server-group "+v)
	}
	if d.Get("active_server_group_allow_server_change").(bool) {
		if d.Get("version").(string) == "v6" {
			return fmt.Errorf("active_server_group_allow_server_change not compatible when version = v6")
		}
		configSet = append(configSet, setPrefix+"active-server-group allow-server-change")
	}
	if d.Get("arp_inspection").(bool) {
		configSet = append(configSet, setPrefix+"arp-inspection")
	}
	if v := d.Get("authentication_password").(string); v != "" {
		configSet = append(configSet, setPrefix+"authentication password \""+v+"\"")
	}
	for _, vBlock := range d.Get("authentication_username_include").([]interface{}) {
		authenticationUsernameInclude := vBlock.(map[string]interface{})
		configSetAuthUsernameInclude, err := setForwardingOptionsDhcpRelayAuthUsernameInclude(
			authenticationUsernameInclude,
			setPrefix,
			d.Get("version").(string),
		)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetAuthUsernameInclude...)
	}
	for _, block := range d.Get("bulk_leasequery").([]interface{}) {
		configSet = append(configSet, setPrefix+"bulk-leasequery")
		if block != nil {
			bulkLeasequery := block.(map[string]interface{})
			if v := bulkLeasequery["attempts"].(int); v != 0 {
				if d.Get("version") == "v6" && v > 10 {
					return fmt.Errorf("expected bulk_leasequery.0.attempts to be in the range (1 - 10) with version = v6, got %d", v)
				}
				configSet = append(configSet, setPrefix+"bulk-leasequery attempts "+strconv.Itoa(v))
			}
			if v := bulkLeasequery["timeout"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"bulk-leasequery timeout "+strconv.Itoa(v))
			}
			if bulkLeasequery["trigger_automatic"].(bool) {
				if d.Get("version").(string) != "v6" {
					return fmt.Errorf("bulk_leasequery.0.trigger_automatic only compatible when version = v6")
				}
				configSet = append(configSet, setPrefix+"bulk-leasequery trigger automatic")
			}
		}
	}
	if v := d.Get("client_response_ttl").(int); v != 0 {
		if d.Get("version").(string) != "v4" {
			return fmt.Errorf("client_response_ttl only compatible when version = v4")
		}
		configSet = append(configSet, setPrefix+"client-response-ttl "+strconv.Itoa(v))
	}
	if v := d.Get("duplicate_clients_in_subnet").(string); v != "" {
		if d.Get("version").(string) != "v4" {
			return fmt.Errorf("duplicate_clients_in_subnet only compatible when version = v4")
		}
		configSet = append(configSet, setPrefix+"duplicate-clients-in-subnet "+v)
	}
	if d.Get("duplicate_clients_incoming_interface").(bool) {
		if d.Get("version").(string) != "v6" {
			return fmt.Errorf("duplicate_clients_incoming_interface only compatible when version = v6")
		}
		configSet = append(configSet, setPrefix+"duplicate-clients incoming-interface")
	}
	if dynProfile := d.Get("dynamic_profile").(string); dynProfile != "" {
		configSet = append(configSet, setPrefix+"dynamic-profile \""+dynProfile+"\"")
		if d.Get("dynamic_profile_aggregate_clients").(bool) {
			configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients")
			if v := d.Get("dynamic_profile_aggregate_clients_action").(string); v != "" {
				configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients "+v)
			}
		} else if d.Get("dynamic_profile_aggregate_clients_action").(string) != "" {
			return fmt.Errorf("dynamic_profile_aggregate_clients need to be true with " +
				"dynamic_profile_aggregate_clients_action")
		}
		if v := d.Get("dynamic_profile_use_primary").(string); v != "" {
			configSet = append(configSet, setPrefix+"dynamic-profile use-primary \""+v+"\"")
		}
	} else if d.Get("dynamic_profile_aggregate_clients").(bool) ||
		d.Get("dynamic_profile_aggregate_clients_action").(string) != "" ||
		d.Get("dynamic_profile_use_primary").(string) != "" {
		return fmt.Errorf("dynamic_profile need to be set with " +
			"dynamic_profile_use_primary, dynamic_profile_aggregate_clients " +
			"and dynamic_profile_aggregate_clients_action")
	}
	if d.Get("exclude_relay_agent_identifier").(bool) {
		if d.Get("version").(string) != "v6" {
			return fmt.Errorf("exclude_relay_agent_identifier only compatible when version = v6")
		}
		configSet = append(configSet, setPrefix+"exclude-relay-agent-identifier")
	}
	if d.Get("forward_only").(bool) {
		configSet = append(configSet, setPrefix+"forward-only")
		if v := d.Get("forward_only_routing_instance").(string); v != "" {
			configSet = append(configSet, setPrefix+"forward-only routing-instance "+v)
		}
	} else if d.Get("forward_only_routing_instance").(string) != "" {
		return fmt.Errorf("forward_only need to be true with forward_only_routing_instance")
	}
	if d.Get("forward_only_replies").(bool) {
		configSet = append(configSet, setPrefix+"forward-only-replies")
	}
	if v := d.Get("forward_snooped_clients").(string); v != "" {
		configSet = append(configSet, setPrefix+"forward-snooped-clients "+v)
	}
	for _, vBlock := range d.Get("lease_time_validation").([]interface{}) {
		configSet = append(configSet, setPrefix+"lease-time-validation")
		if vBlock != nil {
			leaseTimeValidation := vBlock.(map[string]interface{})
			if v := leaseTimeValidation["lease_time_threshold"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"lease-time-validation lease-time-threshold "+strconv.Itoa(v))
			}
			if leaseTimeValidation["violation_action_drop"].(bool) {
				configSet = append(configSet, setPrefix+"lease-time-validation violation-action drop")
			}
		}
	}
	for _, block := range d.Get("leasequery").([]interface{}) {
		configSet = append(configSet, setPrefix+"leasequery")
		if block != nil {
			leasequery := block.(map[string]interface{})
			if v := leasequery["attempts"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"leasequery attempts "+strconv.Itoa(v))
			}
			if v := leasequery["timeout"].(int); v != 0 {
				configSet = append(configSet, setPrefix+"leasequery timeout "+strconv.Itoa(v))
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
		if ldmLayer2 == nil {
			return fmt.Errorf("liveness_detection_method_layer2 block is empty")
		}
		liveDetectMethLayer2 := ldmLayer2.(map[string]interface{})
		setPrefixLDMLayer2 := setPrefix + "liveness-detection method layer2-liveness-detection "
		if v := liveDetectMethLayer2["max_consecutive_retries"].(int); v != 0 {
			configSet = append(configSet, setPrefixLDMLayer2+"max-consecutive-retries "+strconv.Itoa(v))
		}
		if v := liveDetectMethLayer2["transmit_interval"].(int); v != 0 {
			configSet = append(configSet, setPrefixLDMLayer2+"transmit-interval "+strconv.Itoa(v))
		}
	}
	if v := d.Get("maximum_hop_count").(int); v != 0 {
		if d.Get("version").(string) == "v6" {
			return fmt.Errorf("maximum_hop_count not compatible if version = v6")
		}
		configSet = append(configSet, setPrefix+"maximum-hop-count "+strconv.Itoa(v))
	}
	if v := d.Get("minimum_wait_time").(int); v != -1 {
		if d.Get("version").(string) == "v6" {
			return fmt.Errorf("minimum_wait_time not compatible if version = v6")
		}
		configSet = append(configSet, setPrefix+"minimum-wait-time "+strconv.Itoa(v))
	}
	if d.Get("no_snoop").(bool) {
		configSet = append(configSet, setPrefix+"no-snoop")
	}
	for _, v := range d.Get("overrides_v4").([]interface{}) {
		if d.Get("version").(string) == "v6" {
			return fmt.Errorf("overrides_v4 not compatible if version = v6")
		}
		if v == nil {
			return fmt.Errorf("overrides_v4 block is empty")
		}
		configSetOverrides, err := setForwardingOptionsDhcpRelayOverridesV4(
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
		if v == nil {
			return fmt.Errorf("overrides_v6 block is empty")
		}
		configSetOverrides, err := setForwardingOptionsDhcpRelayOverridesV6(
			v.(map[string]interface{}), setPrefix)
		if err != nil {
			return err
		}
		configSet = append(configSet, configSetOverrides...)
	}
	if d.Get("persistent_storage_automatic").(bool) {
		configSet = append(configSet, setPrefix+"persistent-storage automatic")
	}
	for _, v := range d.Get("relay_agent_interface_id").([]interface{}) {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("relay_agent_interface_id not compatible if version = v4")
		}
		configSet = append(configSet, setPrefix+"relay-agent-interface-id")
		if v != nil {
			configSetRelayAgentID, err := setForwardingOptionsDhcpRelayAgentID(
				v.(map[string]interface{}), setPrefix+"relay-agent-interface-id ", "relay_agent_interface_id")
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetRelayAgentID...)
		}
	}
	if d.Get("relay_agent_option_79").(bool) {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("relay_agent_interface_id not compatible if version = v4")
		}
		configSet = append(configSet, setPrefix+"relay-agent-option-79")
	}
	for _, v := range d.Get("relay_agent_remote_id").([]interface{}) {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("relay_agent_interface_id not compatible if version = v4")
		}
		configSet = append(configSet, setPrefix+"relay-agent-remote-id")
		if v != nil {
			configSetRelayAgentID, err := setForwardingOptionsDhcpRelayAgentID(
				v.(map[string]interface{}), setPrefix+"relay-agent-remote-id ", "relay_agent_remote_id")
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetRelayAgentID...)
		}
	}
	for _, v := range d.Get("relay_option").([]interface{}) {
		configSet = append(configSet, setPrefix+"relay-option")
		if v != nil {
			configSetRelayOption, err := setForwardingOptionsDhcpRelayOption(
				v.(map[string]interface{}), setPrefix, d.Get("version").(string))
			if err != nil {
				return err
			}
			configSet = append(configSet, configSetRelayOption...)
		}
	}
	for _, v := range d.Get("relay_option_82").([]interface{}) {
		if d.Get("version").(string) == "v6" {
			return fmt.Errorf("relay_option_82 not compatible if version = v6")
		}
		configSet = append(configSet, setPrefix+"relay-option-82")
		if v != nil {
			configSet = append(configSet, setForwardingOptionsDhcpRelayOption82(v.(map[string]interface{}), setPrefix)...)
		}
	}
	if d.Get("remote_id_mismatch_disconnect").(bool) {
		configSet = append(configSet, setPrefix+"remote-id-mismatch disconnect")
	}
	if d.Get("route_suppression_access").(bool) {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("route_suppression_access not compatible if version = v4")
		}
		configSet = append(configSet, setPrefix+"route-suppression access")
	}
	if d.Get("route_suppression_access_internal").(bool) {
		configSet = append(configSet, setPrefix+"route-suppression access-internal")
	}
	if d.Get("route_suppression_destination").(bool) {
		if d.Get("version").(string) == "v6" {
			return fmt.Errorf("route_suppression_destination not compatible if version = v6")
		}
		configSet = append(configSet, setPrefix+"route-suppression destination")
	}
	serverMatchAddressList := make([]string, 0)
	for _, v := range d.Get("server_match_address").(*schema.Set).List() {
		serverMatchAddress := v.(map[string]interface{})
		if bchk.StringInSlice(serverMatchAddress["address"].(string), serverMatchAddressList) {
			return fmt.Errorf("multiple blocks server_match_address with the same address %s",
				serverMatchAddress["address"].(string))
		}
		serverMatchAddressList = append(serverMatchAddressList, serverMatchAddress["address"].(string))
		configSet = append(configSet, setPrefix+"server-match address "+
			serverMatchAddress["address"].(string)+" "+serverMatchAddress["action"].(string))
	}
	if v := d.Get("server_match_default_action").(string); v != "" {
		configSet = append(configSet, setPrefix+"server-match default-action "+v)
	}
	serverMatchDuidList := make([]string, 0)
	for _, v := range d.Get("server_match_duid").(*schema.Set).List() {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("server_match_duid not compatible if version = v4")
		}
		serverMatchDuid := v.(map[string]interface{})
		serverMatchDuidCompare := serverMatchDuid["compare"].(string)
		serverMatchDuidValueType := serverMatchDuid["value_type"].(string)
		serverMatchDuidValue := serverMatchDuid["value"].(string)
		if bchk.StringInSlice(
			serverMatchDuidCompare+idSeparator+serverMatchDuidValueType+idSeparator+serverMatchDuidValue,
			serverMatchDuidList,
		) {
			return fmt.Errorf("multiple blocks server_match_duid with the same compare %s, value_type %s, value %s",
				serverMatchDuidCompare, serverMatchDuidValueType, serverMatchDuidValue)
		}
		serverMatchDuidList = append(
			serverMatchDuidList,
			serverMatchDuidCompare+idSeparator+serverMatchDuidValueType+idSeparator+serverMatchDuidValue,
		)
		configSet = append(configSet, setPrefix+"server-match duid "+
			serverMatchDuidCompare+" "+
			serverMatchDuidValueType+" "+
			"\""+serverMatchDuidValue+"\" "+
			serverMatchDuid["action"].(string))
	}
	if v := d.Get("server_response_time").(int); v != -1 {
		configSet = append(configSet, setPrefix+"server-response-time "+strconv.Itoa(v))
	}
	if v := d.Get("service_profile").(string); v != "" {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("service_profile not compatible if version = v4")
		}
		configSet = append(configSet, setPrefix+"service-profile \""+v+"\"")
	}
	if v := d.Get("short_cycle_protection_lockout_max_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-max-time "+strconv.Itoa(v))
	}
	if v := d.Get("short_cycle_protection_lockout_min_time").(int); v != 0 {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-min-time "+strconv.Itoa(v))
	}
	if d.Get("source_ip_change").(bool) {
		if d.Get("version").(string) == "v6" {
			return fmt.Errorf("source_ip_change not compatible if version = v6")
		}
		configSet = append(configSet, setPrefix+"source-ip-change")
	}
	if d.Get("vendor_specific_information_host_name").(bool) {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("vendor_specific_information_host_name not compatible if version = v4")
		}
		configSet = append(configSet, setPrefix+"vendor-specific-information host-name")
	}
	if d.Get("vendor_specific_information_location").(bool) {
		if d.Get("version").(string) == "v4" {
			return fmt.Errorf("vendor_specific_information_location not compatible if version = v4")
		}
		configSet = append(configSet, setPrefix+"vendor-specific-information location")
	}

	return clt.configSet(configSet, junSess)
}

func readForwardingOptionsDhcpRelay( //nolint: gocognit, gocyclo
	instance, version string, clt *Client, junSess *junosSession,
) (fwdOptsDhcpRelayOptions, error) {
	var confRead fwdOptsDhcpRelayOptions
	confRead.minimumWaitTime = -1    // default = -1
	confRead.serverResponseTime = -1 // default = -1

	showCmd := cmdShowConfig
	if instance != defaultW {
		showCmd += routingInstancesWS + instance + " "
	}
	showCmd += "forwarding-options dhcp-relay "
	if version == "v6" {
		showCmd += "dhcpv6 "
	}

	showConfig, err := clt.command(showCmd+pipeDisplaySetRelative, junSess)
	if err != nil {
		return confRead, err
	}
	confRead.routingInstance = instance
	confRead.version = version
	if showConfig != emptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, setLS)
			if strings.Contains(item, xmlStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, xmlEndTagConfigOut) {
				break
			}
			switch {
			case strings.HasPrefix(itemTrim, "access-profile "):
				confRead.accessProfile = strings.Trim(strings.TrimPrefix(itemTrim, "access-profile "), "\"")
			case strings.HasPrefix(itemTrim, "active-leasequery"):
				if len(confRead.activeLeasequery) == 0 {
					confRead.activeLeasequery = append(confRead.activeLeasequery, map[string]interface{}{
						"idle_timeout":      0,
						"peer_address":      "",
						"timeout":           0,
						"topology_discover": false,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "active-leasequery idle-timeout "):
					var err error
					confRead.activeLeasequery[0]["idle_timeout"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "active-leasequery idle-timeout "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "active-leasequery peer-address "):
					confRead.activeLeasequery[0]["peer_address"] = strings.TrimPrefix(
						itemTrim, "active-leasequery peer-address ")
				case strings.HasPrefix(itemTrim, "active-leasequery timeout "):
					var err error
					confRead.activeLeasequery[0]["timeout"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "active-leasequery timeout "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case itemTrim == "active-leasequery topology-discover":
					confRead.activeLeasequery[0]["topology_discover"] = true
				}
			case itemTrim == "active-server-group allow-server-change":
				confRead.activeServerGroupAllowServerChange = true
			case strings.HasPrefix(itemTrim, "active-server-group "):
				confRead.activeServerGroup = strings.TrimPrefix(itemTrim, "active-server-group ")
			case itemTrim == "arp-inspection":
				confRead.arpInspection = true
			case strings.HasPrefix(itemTrim, "authentication password "):
				confRead.authenticationPassword = strings.Trim(strings.TrimPrefix(itemTrim, "authentication password "), "\"")
			case strings.HasPrefix(itemTrim, "authentication username-include "):
				if len(confRead.authenticationUsernameInclude) == 0 {
					confRead.authenticationUsernameInclude = append(confRead.authenticationUsernameInclude,
						genForwardingOptionsDhcpRelayAuthUsernameInclude())
				}
				readForwardingOptionsDhcpRelayAuthUsernameInclude(
					strings.TrimPrefix(itemTrim, "authentication username-include "),
					confRead.authenticationUsernameInclude[0],
				)
			case strings.HasPrefix(itemTrim, "bulk-leasequery"):
				if len(confRead.bulkLeasequery) == 0 {
					confRead.bulkLeasequery = append(confRead.bulkLeasequery, map[string]interface{}{
						"attempts":          0,
						"timeout":           0,
						"trigger_automatic": false,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "bulk-leasequery attempts "):
					var err error
					confRead.bulkLeasequery[0]["attempts"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "bulk-leasequery attempts "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "bulk-leasequery timeout "):
					var err error
					confRead.bulkLeasequery[0]["timeout"], err = strconv.Atoi(strings.TrimPrefix(
						itemTrim, "bulk-leasequery timeout "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case itemTrim == "bulk-leasequery trigger automatic":
					confRead.bulkLeasequery[0]["trigger_automatic"] = true
				}
			case strings.HasPrefix(itemTrim, "client-response-ttl "):
				var err error
				confRead.clientResponseTTL, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "client-response-ttl "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "duplicate-clients-in-subnet "):
				confRead.duplicateClientsInSubnet = strings.TrimPrefix(itemTrim, "duplicate-clients-in-subnet ")
			case itemTrim == "duplicate-clients incoming-interface":
				confRead.duplicateClientsIncomingInterface = true
			case strings.HasPrefix(itemTrim, "dynamic-profile aggregate-clients"):
				confRead.dynamicProfileAggregateClients = true
				if strings.HasPrefix(itemTrim, "dynamic-profile aggregate-clients ") {
					confRead.dynamicProfileAggregateClientsAction = strings.TrimPrefix(itemTrim, "dynamic-profile aggregate-clients ")
				}
			case strings.HasPrefix(itemTrim, "dynamic-profile use-primary "):
				confRead.dynamicProfileUsePrimary = strings.Trim(strings.TrimPrefix(itemTrim, "dynamic-profile use-primary "), "\"")
			case strings.HasPrefix(itemTrim, "dynamic-profile "):
				confRead.dynamicProfile = strings.Trim(strings.TrimPrefix(itemTrim, "dynamic-profile "), "\"")
			case itemTrim == "exclude-relay-agent-identifier":
				confRead.excludeRelayAgentIdentifier = true
			case itemTrim == "forward-only-replies":
				confRead.forwardOnlyReplies = true
			case strings.HasPrefix(itemTrim, "forward-only"):
				confRead.forwardOnly = true
				if strings.HasPrefix(itemTrim, "forward-only routing-instance ") {
					confRead.forwardOnlyRoutingInstance = strings.Trim(strings.TrimPrefix(
						itemTrim, "forward-only routing-instance "), "\"")
				}
			case strings.HasPrefix(itemTrim, "forward-snooped-clients "):
				confRead.forwardSnoopedClients = strings.TrimPrefix(itemTrim, "forward-snooped-clients ")
			case strings.HasPrefix(itemTrim, "lease-time-validation"):
				if len(confRead.leaseTimeValidation) == 0 {
					confRead.leaseTimeValidation = append(confRead.leaseTimeValidation, map[string]interface{}{
						"lease_time_threshold":  0,
						"violation_action_drop": false,
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
				case itemTrim == "lease-time-validation violation-action drop":
					confRead.leaseTimeValidation[0]["violation_action_drop"] = true
				}
			case strings.HasPrefix(itemTrim, "leasequery"):
				if len(confRead.leasequery) == 0 {
					confRead.leasequery = append(confRead.leasequery, map[string]interface{}{
						"attempts": 0,
						"timeout":  0,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "leasequery attempts "):
					var err error
					confRead.leasequery[0]["attempts"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "leasequery attempts "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
				case strings.HasPrefix(itemTrim, "leasequery timeout "):
					var err error
					confRead.leasequery[0]["timeout"], err = strconv.Atoi(strings.TrimPrefix(itemTrim, "leasequery timeout "))
					if err != nil {
						return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
					}
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
			case strings.HasPrefix(itemTrim, "maximum-hop-count "):
				var err error
				confRead.maximumHopCount, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "maximum-hop-count "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case strings.HasPrefix(itemTrim, "minimum-wait-time "):
				var err error
				confRead.minimumWaitTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "minimum-wait-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case itemTrim == "no-snoop":
				confRead.noSnoop = true
			case strings.HasPrefix(itemTrim, "overrides "):
				if version == "v4" {
					if len(confRead.overridesV4) == 0 {
						confRead.overridesV4 = append(confRead.overridesV4,
							genForwardingOptionsDhcpRelayOverridesV4())
					}
					if err := readForwardingOptionsDhcpRelayOverridesV4(
						strings.TrimPrefix(itemTrim, "overrides "),
						confRead.overridesV4[0],
					); err != nil {
						return confRead, err
					}
				} else if version == "v6" {
					if len(confRead.overridesV6) == 0 {
						confRead.overridesV6 = append(confRead.overridesV6,
							genForwardingOptionsDhcpRelayOverridesV6())
					}
					if err := readForwardingOptionsDhcpRelayOverridesV6(
						strings.TrimPrefix(itemTrim, "overrides "),
						confRead.overridesV6[0],
					); err != nil {
						return confRead, err
					}
				}
			case itemTrim == "persistent-storage automatic":
				confRead.persistentStorageAutomatic = true
			case strings.HasPrefix(itemTrim, "relay-agent-interface-id"):
				if len(confRead.relayAgentInterfaceID) == 0 {
					confRead.relayAgentInterfaceID = append(confRead.relayAgentInterfaceID,
						genForwardingOptionsDhcpRelayAgentID(true))
				}
				if strings.HasPrefix(itemTrim, "relay-agent-interface-id ") {
					readForwardingOptionsDhcpRelayAgentID(
						strings.TrimPrefix(itemTrim, "relay-agent-interface-id "),
						confRead.relayAgentInterfaceID[0],
					)
				}
			case itemTrim == "relay-agent-option-79":
				confRead.relayAgentOption79 = true
			case strings.HasPrefix(itemTrim, "relay-agent-remote-id"):
				if len(confRead.relayAgentRemoteID) == 0 {
					confRead.relayAgentRemoteID = append(confRead.relayAgentRemoteID,
						genForwardingOptionsDhcpRelayAgentID(false))
				}
				if strings.HasPrefix(itemTrim, "relay-agent-remote-id ") {
					readForwardingOptionsDhcpRelayAgentID(
						strings.TrimPrefix(itemTrim, "relay-agent-remote-id "),
						confRead.relayAgentRemoteID[0],
					)
				}
			case strings.HasPrefix(itemTrim, "relay-option-82"):
				if len(confRead.relayOption82) == 0 {
					confRead.relayOption82 = append(confRead.relayOption82,
						genForwardingOptionsDhcpRelayOption82())
				}
				if strings.HasPrefix(itemTrim, "relay-option-82 ") {
					readForwardingOptionsDhcpRelayOption82(
						strings.TrimPrefix(itemTrim, "relay-option-82 "),
						confRead.relayOption82[0],
					)
				}
			case strings.HasPrefix(itemTrim, "relay-option"):
				if len(confRead.relayOption) == 0 {
					confRead.relayOption = append(confRead.relayOption,
						genForwardingOptionsDhcpRelayOption())
				}
				if strings.HasPrefix(itemTrim, "relay-option ") {
					if err := readForwardingOptionsDhcpRelayOption(
						strings.TrimPrefix(itemTrim, "relay-option "),
						confRead.relayOption[0],
					); err != nil {
						return confRead, err
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
			case strings.HasPrefix(itemTrim, "server-match address "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "server-match address "), " ")
				if len(itemTrimSplit) < 2 {
					return confRead, fmt.Errorf("can't read values for server_match_address in '%s'", itemTrim)
				}
				confRead.serverMatchAddress = append(confRead.serverMatchAddress, map[string]interface{}{
					"address": itemTrimSplit[0],
					"action":  itemTrimSplit[1],
				})
			case strings.HasPrefix(itemTrim, "server-match default-action "):
				confRead.serverMatchDefaultAction = strings.TrimPrefix(itemTrim, "server-match default-action ")
			case strings.HasPrefix(itemTrim, "server-match duid "):
				itemTrimSplit := strings.Split(strings.TrimPrefix(itemTrim, "server-match duid "), " ")
				if len(itemTrimSplit) < 4 {
					return confRead, fmt.Errorf("can't read values for server_match_duid in '%s'", itemTrim)
				}
				if strings.Contains(itemTrimSplit[2], "\"") {
					action := itemTrimSplit[len(itemTrimSplit)-1]
					value := strings.Trim(strings.Join(itemTrimSplit[2:len(itemTrimSplit)-1], " "), "\"")
					itemTrimSplit[2] = value
					itemTrimSplit[3] = action
				}
				confRead.serverMatchDuid = append(confRead.serverMatchDuid, map[string]interface{}{
					"compare":    itemTrimSplit[0],
					"value_type": itemTrimSplit[1],
					"value":      itemTrimSplit[2],
					"action":     itemTrimSplit[3],
				})
			case strings.HasPrefix(itemTrim, "server-response-time "):
				var err error
				confRead.serverResponseTime, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "server-response-time "))
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
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
			case itemTrim == "source-ip-change":
				confRead.sourceIPChange = true
			case itemTrim == "vendor-specific-information host-name":
				confRead.vendorSpecificInformationHostName = true
			case itemTrim == "vendor-specific-information location":
				confRead.vendorSpecificInformationLocation = true
			}
		}
	}

	return confRead, nil
}

func delForwardingOptionsDhcpRelay(instance, version string, clt *Client, junSess *junosSession) error {
	var delPrefix string
	switch {
	case instance == defaultW && version == "v6":
		delPrefix = deleteLS + "forwarding-options dhcp-relay dhcpv6 "
	case instance == defaultW && version == "v4":
		delPrefix = deleteLS + "forwarding-options dhcp-relay "
	case instance != defaultW && version == "v6":
		delPrefix = delRoutingInstances + instance + " forwarding-options dhcp-relay dhcpv6 "
	case instance != defaultW && version == "v4":
		delPrefix = delRoutingInstances + instance + " forwarding-options dhcp-relay "
	}
	listLinesToDelete := []string{
		"access-profile",
		"active-leasequery",
		"active-server-group",
		"arp-inspection",
		"authentication",
		"bulk-leasequery",
		"client-response-ttl",
		"duplicate-clients-in-subnet",
		"duplicate-clients",
		"dynamic-profile",
		"exclude-relay-agent-identifier",
		"forward-only",
		"forward-only-replies",
		"forward-snooped-clients",
		"lease-time-validation",
		"leasequery",
		"liveness-detection",
		"maximum-hop-count",
		"minimum-wait-time",
		"no-snoop",
		"overrides",
		"persistent-storage",
		"relay-agent-interface-id",
		"relay-agent-option-79",
		"relay-agent-remote-id",
		"relay-option",
		"relay-option-82",
		"remote-id-mismatch",
		"route-suppression",
		"server-match",
		"server-response-time",
		"service-profile",
		"short-cycle-protection",
		"source-ip-change",
		"vendor-specific-information",
	}
	configSet := make([]string, len(listLinesToDelete))
	for k, line := range listLinesToDelete {
		configSet[k] = delPrefix + line
	}

	return clt.configSet(configSet, junSess)
}

func fillForwardingOptionsDhcpRelayData(
	d *schema.ResourceData, fwdOptsDhcpRelayOptions fwdOptsDhcpRelayOptions,
) {
	if tfErr := d.Set(
		"routing_instance",
		fwdOptsDhcpRelayOptions.routingInstance,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"version",
		fwdOptsDhcpRelayOptions.version,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"access_profile",
		fwdOptsDhcpRelayOptions.accessProfile,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"active_leasequery",
		fwdOptsDhcpRelayOptions.activeLeasequery,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"active_server_group",
		fwdOptsDhcpRelayOptions.activeServerGroup,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"active_server_group_allow_server_change",
		fwdOptsDhcpRelayOptions.activeServerGroupAllowServerChange,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"arp_inspection",
		fwdOptsDhcpRelayOptions.arpInspection,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"authentication_password",
		fwdOptsDhcpRelayOptions.authenticationPassword,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"authentication_username_include",
		fwdOptsDhcpRelayOptions.authenticationUsernameInclude,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"bulk_leasequery",
		fwdOptsDhcpRelayOptions.bulkLeasequery,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"client_response_ttl",
		fwdOptsDhcpRelayOptions.clientResponseTTL,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"duplicate_clients_in_subnet",
		fwdOptsDhcpRelayOptions.duplicateClientsInSubnet,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"duplicate_clients_incoming_interface",
		fwdOptsDhcpRelayOptions.duplicateClientsIncomingInterface,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"dynamic_profile",
		fwdOptsDhcpRelayOptions.dynamicProfile,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"dynamic_profile_aggregate_clients",
		fwdOptsDhcpRelayOptions.dynamicProfileAggregateClients,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"dynamic_profile_aggregate_clients_action",
		fwdOptsDhcpRelayOptions.dynamicProfileAggregateClientsAction,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"dynamic_profile_use_primary",
		fwdOptsDhcpRelayOptions.dynamicProfileUsePrimary,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"exclude_relay_agent_identifier",
		fwdOptsDhcpRelayOptions.excludeRelayAgentIdentifier,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"forward_only",
		fwdOptsDhcpRelayOptions.forwardOnly,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"forward_only_replies",
		fwdOptsDhcpRelayOptions.forwardOnlyReplies,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"forward_only_routing_instance",
		fwdOptsDhcpRelayOptions.forwardOnlyRoutingInstance,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"forward_snooped_clients",
		fwdOptsDhcpRelayOptions.forwardSnoopedClients,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"lease_time_validation",
		fwdOptsDhcpRelayOptions.leaseTimeValidation,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"leasequery",
		fwdOptsDhcpRelayOptions.leasequery,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"liveness_detection_failure_action",
		fwdOptsDhcpRelayOptions.livenessDetectionFailureAction,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"liveness_detection_method_bfd",
		fwdOptsDhcpRelayOptions.livenessDetectionMethodBfd,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"liveness_detection_method_layer2",
		fwdOptsDhcpRelayOptions.livenessDetectionMethodLayer2,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"maximum_hop_count",
		fwdOptsDhcpRelayOptions.maximumHopCount,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"minimum_wait_time",
		fwdOptsDhcpRelayOptions.minimumWaitTime,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"no_snoop",
		fwdOptsDhcpRelayOptions.noSnoop,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"overrides_v4",
		fwdOptsDhcpRelayOptions.overridesV4,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"overrides_v6",
		fwdOptsDhcpRelayOptions.overridesV6,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"persistent_storage_automatic",
		fwdOptsDhcpRelayOptions.persistentStorageAutomatic,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"relay_agent_interface_id",
		fwdOptsDhcpRelayOptions.relayAgentInterfaceID,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"relay_agent_option_79",
		fwdOptsDhcpRelayOptions.relayAgentOption79,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"relay_agent_remote_id",
		fwdOptsDhcpRelayOptions.relayAgentRemoteID,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"relay_option",
		fwdOptsDhcpRelayOptions.relayOption,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"relay_option_82",
		fwdOptsDhcpRelayOptions.relayOption82,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"remote_id_mismatch_disconnect",
		fwdOptsDhcpRelayOptions.remoteIDMismatchDisconnect,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"route_suppression_access",
		fwdOptsDhcpRelayOptions.routeSuppressionAccess,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"route_suppression_access_internal",
		fwdOptsDhcpRelayOptions.routeSuppressionAccessInternal,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"route_suppression_destination",
		fwdOptsDhcpRelayOptions.routeSuppressionDestination,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"server_match_address",
		fwdOptsDhcpRelayOptions.serverMatchAddress,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"server_match_default_action",
		fwdOptsDhcpRelayOptions.serverMatchDefaultAction,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"server_match_duid",
		fwdOptsDhcpRelayOptions.serverMatchDuid,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"server_response_time",
		fwdOptsDhcpRelayOptions.serverResponseTime,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"service_profile",
		fwdOptsDhcpRelayOptions.serviceProfile,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"short_cycle_protection_lockout_max_time",
		fwdOptsDhcpRelayOptions.shortCycleProtectionLockoutMaxTime,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"short_cycle_protection_lockout_min_time",
		fwdOptsDhcpRelayOptions.shortCycleProtectionLockoutMinTime,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"source_ip_change",
		fwdOptsDhcpRelayOptions.sourceIPChange,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"vendor_specific_information_host_name",
		fwdOptsDhcpRelayOptions.vendorSpecificInformationHostName,
	); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set(
		"vendor_specific_information_location",
		fwdOptsDhcpRelayOptions.vendorSpecificInformationLocation,
	); tfErr != nil {
		panic(tfErr)
	}
}
