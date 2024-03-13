package providersdk

import (
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

func schemaForwardingOptionsDhcpRelayAuthUsernameInclude() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"client_id_exclude_headers": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"client_id_use_automatic_ascii_hex_encoding": {
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
		"option_60": { // only dhcpv4
			Type:     schema.TypeBool,
			Optional: true,
		},
		"option_82": { // only dhcpv4
			Type:     schema.TypeBool,
			Optional: true,
		},
		"option_82_circuit_id": { // only dhcpv4
			Type:         schema.TypeBool,
			Optional:     true,
			RequiredWith: []string{"authentication_username_include.0.option_82"},
		},
		"option_82_remote_id": { // only dhcpv4
			Type:         schema.TypeBool,
			Optional:     true,
			RequiredWith: []string{"authentication_username_include.0.option_82"},
		},
		"relay_agent_interface_id": { // only dhcpv6
			Type:     schema.TypeBool,
			Optional: true,
		},
		"relay_agent_remote_id": { // only dhcpv6
			Type:     schema.TypeBool,
			Optional: true,
		},
		"relay_agent_subscriber_id": { // only dhcpv6
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
	}
}

func schemaForwardingOptionsDhcpRelayOverridesV4() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_no_end_option": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"allow_snooped_clients": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"always_write_giaddr": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"always_write_option_82": {
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
		"delay_authentication": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"delete_binding_on_renegotiation": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"disable_relay": {
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
		"layer2_unicast_replies": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"no_allow_snooped_clients": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"no_bind_on_request": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"no_unicast_replies": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"proxy_mode": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"relay_source": {
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
		"replace_ip_source_with_giaddr": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"send_release_on_delete": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"trust_option_82": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"user_defined_option_82": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func schemaForwardingOptionsDhcpRelayOverridesV6() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_snooped_clients": {
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
		"delay_authentication": {
			Type:     schema.TypeBool,
			Optional: true,
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
		"no_allow_snooped_clients": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"no_bind_on_request": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"relay_source": {
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
		"send_release_on_delete": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func schemaForwardingOptionsDhcpRelayAgentID(keepIncomingIDStrict bool) map[string]*schema.Schema {
	r := map[string]*schema.Schema{
		"include_irb_and_l2": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"keep_incoming_id": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"no_vlan_interface_name": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"prefix_host_name": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"prefix_routing_instance_name": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"use_interface_description": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"device", "logical"}, false),
		},
		"use_option_82": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"use_option_82_strict": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"use_vlan_id": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
	if keepIncomingIDStrict {
		r["keep_incoming_id_strict"] = &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		}
	}

	return r
}

func schemaForwardingOptionsDhcpRelayOption() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"option_15": { // only dhcpv6
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
						ValidateFunc: validation.StringInSlice([]string{"drop", "forward-only", "relay-server-group"}, false),
					},
					"group": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"option_15_default_action": { // only dhcpv6
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"drop", "forward-only", "relay-server-group"}, false),
					},
					"group": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"option_16": { // only dhcpv6
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
						ValidateFunc: validation.StringInSlice([]string{"drop", "forward-only", "relay-server-group"}, false),
					},
					"group": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"option_16_default_action": { // only dhcpv6
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"drop", "forward-only", "relay-server-group"}, false),
					},
					"group": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"option_60": { // only dhcpv4
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
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"drop",
							"forward-only",
							"local-server-group",
							"relay-server-group",
						}, false),
					},
					"group": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"option_60_default_action": { // only dhcpv4
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"drop",
							"forward-only",
							"local-server-group",
							"relay-server-group",
						}, false),
					},
					"group": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"option_77": { // only dhcpv4
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
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"drop",
							"forward-only",
							"local-server-group",
							"relay-server-group",
						}, false),
					},
					"group": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"option_77_default_action": { // only dhcpv4
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"drop",
							"forward-only",
							"local-server-group",
							"relay-server-group",
						}, false),
					},
					"group": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"option_order": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"15", "16", "60", "77"}, false),
			},
		},
	}
}

func schemaForwardingOptionsDhcpRelayOption82() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"circuit_id": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"include_irb_and_l2": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"keep_incoming_circuit_id": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"no_vlan_interface_name": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"prefix_host_name": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"prefix_routing_instance_name": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"use_interface_description": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"device", "logical"}, false),
					},
					"use_vlan_id": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"user_defined": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"vlan_id_only": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"exclude_relay_agent_identifier": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"link_selection": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"remote_id": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hostname_only": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"include_irb_and_l2": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"keep_incoming_remote_id": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"no_vlan_interface_name": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"prefix_host_name": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"prefix_routing_instance_name": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"use_interface_description": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"device", "logical"}, false),
					},
					"use_string": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"use_vlan_id": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"server_id_override": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"vendor_specific_host_name": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"vendor_specific_location": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func setForwardingOptionsDhcpRelayAuthUsernameInclude(
	authenticationUsernameInclude map[string]interface{}, setPrefixSrc, version string,
) ([]string, error) {
	configSet := make([]string, 0)

	setPrefix := setPrefixSrc + "authentication username-include "
	if authenticationUsernameInclude["circuit_type"].(bool) {
		configSet = append(configSet, setPrefix+"circuit-type")
	}
	if authenticationUsernameInclude["client_id"].(bool) {
		configSet = append(configSet, setPrefix+"client-id")
		if authenticationUsernameInclude["client_id_exclude_headers"].(bool) {
			configSet = append(configSet, setPrefix+"client-id exclude-headers")
		}
		if authenticationUsernameInclude["client_id_use_automatic_ascii_hex_encoding"].(bool) {
			configSet = append(configSet, setPrefix+"client-id use-automatic-ascii-hex-encoding")
		}
	} else if authenticationUsernameInclude["client_id_exclude_headers"].(bool) ||
		authenticationUsernameInclude["client_id_use_automatic_ascii_hex_encoding"].(bool) {
		return configSet, errors.New("authentication_username_include.0.client_id need to be true with " +
			"client_id_exclude_headers or client_id_use_automatic_ascii_hex_encoding")
	}
	if v := authenticationUsernameInclude["delimiter"].(string); v != "" {
		configSet = append(configSet, setPrefix+"delimiter \""+v+"\"")
	}
	if v := authenticationUsernameInclude["domain_name"].(string); v != "" {
		configSet = append(configSet, setPrefix+"domain-name \""+v+"\"")
	}
	if v := authenticationUsernameInclude["interface_description"].(string); v != "" {
		configSet = append(configSet, setPrefix+"interface-description "+v)
	}
	if authenticationUsernameInclude["interface_name"].(bool) {
		configSet = append(configSet, setPrefix+"interface-name")
	}
	if authenticationUsernameInclude["mac_address"].(bool) {
		configSet = append(configSet, setPrefix+"mac-address")
	}
	if authenticationUsernameInclude["option_60"].(bool) {
		if version == "v6" {
			return configSet, errors.New("authentication_username_include.0.option_60 not compatible when version = v6")
		}
		configSet = append(configSet, setPrefix+"option-60")
	}
	if authenticationUsernameInclude["option_82"].(bool) {
		if version == "v6" {
			return configSet, errors.New("authentication_username_include.0.option_82 not compatible when version = v6")
		}
		configSet = append(configSet, setPrefix+"option-82")
		if authenticationUsernameInclude["option_82_circuit_id"].(bool) {
			configSet = append(configSet, setPrefix+"option-82 circuit-id")
		}
		if authenticationUsernameInclude["option_82_remote_id"].(bool) {
			configSet = append(configSet, setPrefix+"option-82 remote-id")
		}
	} else if authenticationUsernameInclude["option_82_circuit_id"].(bool) ||
		authenticationUsernameInclude["option_82_remote_id"].(bool) {
		return configSet, errors.New("authentication_username_include.0.option_82 need to be true with " +
			"option_82_circuit_id or option_82_remote_id")
	}
	if authenticationUsernameInclude["relay_agent_interface_id"].(bool) {
		if version == "v4" {
			return configSet, errors.New("authentication_username_include.0.relay_agent_interface_id" +
				" not compatible when version = v4")
		}
		configSet = append(configSet, setPrefix+"relay-agent-interface-id")
	}
	if authenticationUsernameInclude["relay_agent_remote_id"].(bool) {
		if version == "v4" {
			return configSet, errors.New("authentication_username_include.0.relay_agent_remote_id" +
				" not compatible when version = v4")
		}
		configSet = append(configSet, setPrefix+"relay-agent-remote-id")
	}
	if authenticationUsernameInclude["relay_agent_subscriber_id"].(bool) {
		if version == "v4" {
			return configSet, errors.New("authentication_username_include.0.relay_agent_subscriber_id" +
				" not compatible when version = v4")
		}
		configSet = append(configSet, setPrefix+"relay-agent-subscriber-id")
	}
	if authenticationUsernameInclude["routing_instance_name"].(bool) {
		configSet = append(configSet, setPrefix+"routing-instance-name")
	}
	if v := authenticationUsernameInclude["user_prefix"].(string); v != "" {
		configSet = append(configSet, setPrefix+"user-prefix \""+v+"\"")
	}
	if authenticationUsernameInclude["vlan_tags"].(bool) {
		configSet = append(configSet, setPrefix+"vlan-tags")
	}

	return configSet, nil
}

func setForwardingOptionsDhcpRelayOverridesV4(overrides map[string]interface{}, setPrefix string,
) ([]string, error) {
	configSet := make([]string, 0)
	setPrefix += "overrides "

	if overrides["allow_no_end_option"].(bool) {
		configSet = append(configSet, setPrefix+"allow-no-end-option")
	}
	if overrides["allow_snooped_clients"].(bool) {
		configSet = append(configSet, setPrefix+"allow-snooped-clients")
	}
	if overrides["always_write_giaddr"].(bool) {
		configSet = append(configSet, setPrefix+"always-write-giaddr")
	}
	if overrides["always_write_option_82"].(bool) {
		configSet = append(configSet, setPrefix+"always-write-option-82")
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
	if overrides["delay_authentication"].(bool) {
		configSet = append(configSet, setPrefix+"delay-authentication")
	}
	if overrides["delete_binding_on_renegotiation"].(bool) {
		configSet = append(configSet, setPrefix+"delete-binding-on-renegotiation")
	}
	if overrides["disable_relay"].(bool) {
		configSet = append(configSet, setPrefix+"disable-relay")
	}
	if v := overrides["dual_stack"].(string); v != "" {
		configSet = append(configSet, setPrefix+"dual-stack \""+v+"\"")
	}
	if v := overrides["interface_client_limit"].(int); v != 0 {
		configSet = append(configSet, setPrefix+"interface-client-limit "+strconv.Itoa(v))
	}
	if overrides["layer2_unicast_replies"].(bool) {
		configSet = append(configSet, setPrefix+"layer2-unicast-replies")
	}
	if overrides["no_allow_snooped_clients"].(bool) {
		if overrides["allow_snooped_clients"].(bool) {
			return configSet, errors.New("allow_snooped_clients and no_allow_snooped_clients can't be true in same time")
		}
		configSet = append(configSet, setPrefix+"no-allow-snooped-clients")
	}
	if overrides["no_bind_on_request"].(bool) {
		configSet = append(configSet, setPrefix+"no-bind-on-request")
	}
	if overrides["no_unicast_replies"].(bool) {
		if overrides["layer2_unicast_replies"].(bool) {
			return configSet, errors.New("layer2_unicast_replies and no_unicast_replies can't be true in same time")
		}
		configSet = append(configSet, setPrefix+"no-unicast-replies")
	}
	if overrides["proxy_mode"].(bool) {
		configSet = append(configSet, setPrefix+"proxy-mode")
	}
	if v := overrides["relay_source"].(string); v != "" {
		configSet = append(configSet, setPrefix+"relay-source "+v)
	}
	if overrides["replace_ip_source_with_giaddr"].(bool) {
		configSet = append(configSet, setPrefix+"replace-ip-source-with giaddr")
	}
	if overrides["send_release_on_delete"].(bool) {
		configSet = append(configSet, setPrefix+"send-release-on-delete")
	}
	if overrides["trust_option_82"].(bool) {
		configSet = append(configSet, setPrefix+"trust-option-82")
	}
	if v := overrides["user_defined_option_82"].(string); v != "" {
		configSet = append(configSet, setPrefix+"user-defined-option-82 \""+v+"\"")
	}

	if len(configSet) == 0 {
		return configSet, errors.New("an overrides_v4 block is empty")
	}

	return configSet, nil
}

func setForwardingOptionsDhcpRelayOverridesV6(overrides map[string]interface{}, setPrefix string,
) ([]string, error) {
	configSet := make([]string, 0)
	setPrefix += "overrides "

	if overrides["allow_snooped_clients"].(bool) {
		configSet = append(configSet, setPrefix+"allow-snooped-clients")
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
	if overrides["delay_authentication"].(bool) {
		configSet = append(configSet, setPrefix+"delay-authentication")
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
	if overrides["no_allow_snooped_clients"].(bool) {
		if overrides["allow_snooped_clients"].(bool) {
			return configSet, errors.New("allow_snooped_clients and no_allow_snooped_clients can't be true in same time")
		}
		configSet = append(configSet, setPrefix+"no-allow-snooped-clients")
	}
	if overrides["no_bind_on_request"].(bool) {
		configSet = append(configSet, setPrefix+"no-bind-on-request")
	}
	if v := overrides["relay_source"].(string); v != "" {
		configSet = append(configSet, setPrefix+"relay-source "+v)
	}
	if overrides["send_release_on_delete"].(bool) {
		configSet = append(configSet, setPrefix+"send-release-on-delete")
	}

	if len(configSet) == 0 {
		return configSet, errors.New("an overrides_v6 block is empty")
	}

	return configSet, nil
}

func setForwardingOptionsDhcpRelayAgentID(relayAgentID map[string]interface{}, setPrefix, blockName string,
) ([]string, error) {
	configSet := make([]string, 0)

	if relayAgentID["include_irb_and_l2"].(bool) {
		configSet = append(configSet, setPrefix+"include-irb-and-l2")
	}
	if relayAgentID["keep_incoming_id"].(bool) {
		switch blockName {
		case "relay_agent_interface_id":
			configSet = append(configSet, setPrefix+"keep-incoming-interface-id")
			if relayAgentID["keep_incoming_id_strict"].(bool) {
				configSet = append(configSet, setPrefix+"keep-incoming-interface-id strict")
			}
		case "relay_agent_remote_id":
			configSet = append(configSet, setPrefix+"keep-incoming-remote-id")
		}
	} else if blockName == "relay_agent_interface_id" && relayAgentID["keep_incoming_id_strict"].(bool) {
		return configSet, errors.New("keep_incoming_id need to be true with keep_incoming_id_strict")
	}
	if relayAgentID["no_vlan_interface_name"].(bool) {
		configSet = append(configSet, setPrefix+"no-vlan-interface-name")
	}
	if relayAgentID["prefix_host_name"].(bool) {
		configSet = append(configSet, setPrefix+"prefix host-name")
	}
	if relayAgentID["prefix_routing_instance_name"].(bool) {
		configSet = append(configSet, setPrefix+"prefix routing-instance-name")
	}
	if v := relayAgentID["use_interface_description"].(string); v != "" {
		configSet = append(configSet, setPrefix+"use-interface-description "+v)
	}
	if relayAgentID["use_option_82"].(bool) {
		configSet = append(configSet, setPrefix+"use-option-82")
		if relayAgentID["use_option_82_strict"].(bool) {
			configSet = append(configSet, setPrefix+"use-option-82 strict")
		}
	} else if relayAgentID["use_option_82_strict"].(bool) {
		return configSet, errors.New("use_option_82 need to be true with use_option_82_strict")
	}
	if relayAgentID["use_vlan_id"].(bool) {
		configSet = append(configSet, setPrefix+"use-vlan-id")
	}

	if len(configSet) == 0 {
		return configSet, fmt.Errorf("an %s block is empty", blockName)
	}

	return configSet, nil
}

func setForwardingOptionsDhcpRelayOption(relayOption map[string]interface{}, setPrefixSrc, version string,
) ([]string, error) {
	configSet := make([]string, 0)
	setPrefix := setPrefixSrc + "relay-option "

	for _, block := range relayOption["option_15"].(*schema.Set).List() {
		if version == "v4" {
			return configSet, errors.New("relay_option.0.option_15 not compatible when version = v4")
		}
		option15 := block.(map[string]interface{})
		if action := option15["action"].(string); action == "relay-server-group" {
			if option15["group"].(string) == "" {
				return configSet, errors.New("group must be set when " +
					"action = relay-server-group in option_15 block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-15 "+
				option15["compare"].(string)+" "+
				option15["value_type"].(string)+" "+
				"\""+option15["value"].(string)+"\" "+
				action+" "+
				"\""+option15["group"].(string)+"\"")
		} else {
			if option15["group"].(string) != "" {
				return configSet, errors.New("group must be set only with " +
					"action = relay-server-group in option_15 block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-15 "+
				option15["compare"].(string)+" "+
				option15["value_type"].(string)+" "+
				"\""+option15["value"].(string)+"\" "+
				action)
		}
	}
	for _, block := range relayOption["option_15_default_action"].([]interface{}) {
		if version == "v4" {
			return configSet, errors.New("relay_option.0.option_15_default_action not compatible when version = v4")
		}
		option15DefAction := block.(map[string]interface{})
		if action := option15DefAction["action"].(string); action == "relay-server-group" {
			if option15DefAction["group"].(string) == "" {
				return configSet, errors.New("group must be set when " +
					"action = relay-server-group in option_15_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-15 default-action "+action+
				" \""+option15DefAction["group"].(string)+"\"")
		} else {
			if option15DefAction["group"].(string) != "" {
				return configSet, errors.New("group must be set only with " +
					"action = relay-server-group in option_15_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-15 default-action "+action)
		}
	}
	for _, block := range relayOption["option_16"].(*schema.Set).List() {
		if version == "v4" {
			return configSet, errors.New("relay_option.0.option_16 not compatible when version = v4")
		}
		option16 := block.(map[string]interface{})
		if action := option16["action"].(string); action == "relay-server-group" {
			if option16["group"].(string) == "" {
				return configSet, errors.New("group must be set when " +
					"action = relay-server-group in option_16 block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-16 "+
				option16["compare"].(string)+" "+
				option16["value_type"].(string)+" "+
				"\""+option16["value"].(string)+"\" "+
				action+" "+
				"\""+option16["group"].(string)+"\"")
		} else {
			if option16["group"].(string) != "" {
				return configSet, errors.New("group must be set only with " +
					"action = relay-server-group in option_16 block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-16 "+
				option16["compare"].(string)+" "+
				option16["value_type"].(string)+" "+
				"\""+option16["value"].(string)+"\" "+
				action)
		}
	}
	for _, block := range relayOption["option_16_default_action"].([]interface{}) {
		if version == "v4" {
			return configSet, errors.New("relay_option.0.option_16_default_action not compatible when version = v4")
		}
		option16DefAction := block.(map[string]interface{})
		if action := option16DefAction["action"].(string); action == "relay-server-group" {
			if option16DefAction["group"].(string) == "" {
				return configSet, errors.New("group must be set when " +
					"action = relay-server-group in option_16_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-16 default-action "+action+
				" \""+option16DefAction["group"].(string)+"\"")
		} else {
			if option16DefAction["group"].(string) != "" {
				return configSet, errors.New("group must be set only with " +
					"action = relay-server-group in option_16_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-16 default-action "+action)
		}
	}
	for _, block := range relayOption["option_60"].(*schema.Set).List() {
		if version == "v6" {
			return configSet, errors.New("relay_option.0.option_60 not compatible when version = v6")
		}
		option60 := block.(map[string]interface{})
		if action := option60["action"].(string); action == "local-server-group" || action == "relay-server-group" {
			if option60["group"].(string) == "" {
				return configSet, errors.New("group must be set when " +
					"action = local-server-group or relay-server-group in option_60 block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-60 "+
				option60["compare"].(string)+" "+
				option60["value_type"].(string)+" "+
				"\""+option60["value"].(string)+"\" "+
				action+" "+
				"\""+option60["group"].(string)+"\"")
		} else {
			if option60["group"].(string) != "" {
				return configSet, errors.New("group must be set only with " +
					"action = local-server-group or relay-server-group in option_60 block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-60 "+
				option60["compare"].(string)+" "+
				option60["value_type"].(string)+" "+
				"\""+option60["value"].(string)+"\" "+
				action)
		}
	}
	for _, block := range relayOption["option_60_default_action"].([]interface{}) {
		if version == "v6" {
			return configSet, errors.New("relay_option.0.option_60_default_action not compatible when version = v6")
		}
		option60DefAction := block.(map[string]interface{})
		if action := option60DefAction["action"].(string); action == "local-server-group" || action == "relay-server-group" {
			if option60DefAction["group"].(string) == "" {
				return configSet, errors.New("group must be set when " +
					"action = local-server-group or relay-server-group in option_60_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-60 default-action "+action+
				" \""+option60DefAction["group"].(string)+"\"")
		} else {
			if option60DefAction["group"].(string) != "" {
				return configSet, errors.New("group must be set only with " +
					"action = local-server-group or relay-server-group in option_60_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-60 default-action "+action)
		}
	}
	for _, block := range relayOption["option_77"].(*schema.Set).List() {
		if version == "v6" {
			return configSet, errors.New("relay_option.0.option_77 not compatible when version = v6")
		}
		option77 := block.(map[string]interface{})
		if action := option77["action"].(string); action == "local-server-group" || action == "relay-server-group" {
			if option77["group"].(string) == "" {
				return configSet, errors.New("group must be set when " +
					"action = local-server-group or relay-server-group in option_77 block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-77 "+
				option77["compare"].(string)+" "+
				option77["value_type"].(string)+" "+
				"\""+option77["value"].(string)+"\" "+
				action+" "+
				"\""+option77["group"].(string)+"\"")
		} else {
			if option77["group"].(string) != "" {
				return configSet, errors.New("group must be set only with " +
					"action = local-server-group or relay-server-group in option_77 block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-77 "+
				option77["compare"].(string)+" "+
				option77["value_type"].(string)+" "+
				"\""+option77["value"].(string)+"\" "+
				action)
		}
	}
	for _, block := range relayOption["option_77_default_action"].([]interface{}) {
		if version == "v6" {
			return configSet, errors.New("relay_option.0.option_77_default_action not compatible when version = v6")
		}
		option77DefAction := block.(map[string]interface{})
		if action := option77DefAction["action"].(string); action == "local-server-group" || action == "relay-server-group" {
			if option77DefAction["group"].(string) == "" {
				return configSet, errors.New("group must be set when " +
					"action = local-server-group or relay-server-group in option_77_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-77 default-action "+action+
				" \""+option77DefAction["group"].(string)+"\"")
		} else {
			if option77DefAction["group"].(string) != "" {
				return configSet, errors.New("group must be set only with " +
					"action = local-server-group or relay-server-group in option_77_default_action block in relay_option block")
			}
			configSet = append(configSet, setPrefix+"option-77 default-action "+action)
		}
	}
	for _, elem := range relayOption["option_order"].([]interface{}) {
		opt := elem.(string)
		if version == "v4" && (opt == "15" || opt == "16") {
			return configSet, errors.New("15 & 16 for value in relay_option.0.option_order not compatible when version = v4")
		}
		if version == "v6" && (opt == "60" || opt == "77") {
			return configSet, errors.New("60 & 77 for value in relay_option.0.option_order not compatible when version = v4")
		}
		configSet = append(configSet, setPrefix+"option-order "+opt)
	}

	return configSet, nil
}

func setForwardingOptionsDhcpRelayOption82(relayOption map[string]interface{}, setPrefixSrc string) []string {
	configSet := make([]string, 0)
	setPrefix := setPrefixSrc + "relay-option-82 "

	for _, block := range relayOption["circuit_id"].([]interface{}) {
		setPrefixCircuitID := setPrefix + "circuit-id "
		configSet = append(configSet, setPrefixCircuitID)
		if block != nil {
			circuitID := block.(map[string]interface{})
			if circuitID["include_irb_and_l2"].(bool) {
				configSet = append(configSet, setPrefixCircuitID+"include-irb-and-l2")
			}
			if circuitID["keep_incoming_circuit_id"].(bool) {
				configSet = append(configSet, setPrefixCircuitID+"keep-incoming-circuit-id")
			}
			if circuitID["no_vlan_interface_name"].(bool) {
				configSet = append(configSet, setPrefixCircuitID+"no-vlan-interface-name")
			}
			if circuitID["prefix_host_name"].(bool) {
				configSet = append(configSet, setPrefixCircuitID+"prefix host-name")
			}
			if circuitID["prefix_routing_instance_name"].(bool) {
				configSet = append(configSet, setPrefixCircuitID+"prefix routing-instance-name")
			}
			if v := circuitID["use_interface_description"].(string); v != "" {
				configSet = append(configSet, setPrefixCircuitID+"use-interface-description "+v)
			}
			if circuitID["use_vlan_id"].(bool) {
				configSet = append(configSet, setPrefixCircuitID+"use-vlan-id")
			}
			if circuitID["user_defined"].(bool) {
				configSet = append(configSet, setPrefixCircuitID+"user-defined")
			}
			if circuitID["vlan_id_only"].(bool) {
				configSet = append(configSet, setPrefixCircuitID+"vlan-id-only")
			}
		}
	}
	if relayOption["exclude_relay_agent_identifier"].(bool) {
		configSet = append(configSet, setPrefix+"exclude-relay-agent-identifier")
	}
	if relayOption["link_selection"].(bool) {
		configSet = append(configSet, setPrefix+"link-selection")
	}
	for _, block := range relayOption["remote_id"].([]interface{}) {
		setPrefixRemoteID := setPrefix + "remote-id "
		configSet = append(configSet, setPrefixRemoteID)
		if block != nil {
			remoteID := block.(map[string]interface{})
			if remoteID["hostname_only"].(bool) {
				configSet = append(configSet, setPrefixRemoteID+"hostname-only")
			}
			if remoteID["include_irb_and_l2"].(bool) {
				configSet = append(configSet, setPrefixRemoteID+"include-irb-and-l2")
			}
			if remoteID["keep_incoming_remote_id"].(bool) {
				configSet = append(configSet, setPrefixRemoteID+"keep-incoming-remote-id")
			}
			if remoteID["no_vlan_interface_name"].(bool) {
				configSet = append(configSet, setPrefixRemoteID+"no-vlan-interface-name")
			}
			if remoteID["prefix_host_name"].(bool) {
				configSet = append(configSet, setPrefixRemoteID+"prefix host-name")
			}
			if remoteID["prefix_routing_instance_name"].(bool) {
				configSet = append(configSet, setPrefixRemoteID+"prefix routing-instance-name")
			}
			if v := remoteID["use_interface_description"].(string); v != "" {
				configSet = append(configSet, setPrefixRemoteID+"use-interface-description "+v)
			}
			if v := remoteID["use_string"].(string); v != "" {
				configSet = append(configSet, setPrefixRemoteID+"use-string \""+v+"\"")
			}
			if remoteID["use_vlan_id"].(bool) {
				configSet = append(configSet, setPrefixRemoteID+"use-vlan-id")
			}
		}
	}
	if relayOption["server_id_override"].(bool) {
		configSet = append(configSet, setPrefix+"server-id-override")
	}
	if relayOption["vendor_specific_host_name"].(bool) {
		configSet = append(configSet, setPrefix+"vendor-specific host-name")
	}
	if relayOption["vendor_specific_location"].(bool) {
		configSet = append(configSet, setPrefix+"vendor-specific location")
	}

	return configSet
}

func genForwardingOptionsDhcpRelayAuthUsernameInclude() map[string]interface{} {
	return map[string]interface{}{
		"circuit_type":              false,
		"client_id":                 false,
		"client_id_exclude_headers": false,
		"client_id_use_automatic_ascii_hex_encoding": false,
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
	}
}

func genForwardingOptionsDhcpRelayOverridesV4() map[string]interface{} {
	return map[string]interface{}{
		"allow_no_end_option":             false,
		"allow_snooped_clients":           false,
		"always_write_giaddr":             false,
		"always_write_option_82":          false,
		"asymmetric_lease_time":           0,
		"bootp_support":                   false,
		"client_discover_match":           "",
		"delay_authentication":            false,
		"delete_binding_on_renegotiation": false,
		"disable_relay":                   false,
		"dual_stack":                      "",
		"interface_client_limit":          0,
		"layer2_unicast_replies":          false,
		"no_allow_snooped_clients":        false,
		"no_bind_on_request":              false,
		"no_unicast_replies":              false,
		"proxy_mode":                      false,
		"relay_source":                    "",
		"replace_ip_source_with_giaddr":   false,
		"send_release_on_delete":          false,
		"trust_option_82":                 false,
		"user_defined_option_82":          "",
	}
}

func genForwardingOptionsDhcpRelayOverridesV6() map[string]interface{} {
	return map[string]interface{}{
		"allow_snooped_clients":                       false,
		"always_process_option_request_option":        false,
		"asymmetric_lease_time":                       0,
		"asymmetric_prefix_lease_time":                0,
		"client_negotiation_match_incoming_interface": false,
		"delay_authentication":                        false,
		"delete_binding_on_renegotiation":             false,
		"dual_stack":                                  "",
		"interface_client_limit":                      0,
		"no_allow_snooped_clients":                    false,
		"no_bind_on_request":                          false,
		"relay_source":                                "",
		"send_release_on_delete":                      false,
	}
}

func genForwardingOptionsDhcpRelayAgentID(keepIncomingIDStrict bool) map[string]interface{} {
	r := map[string]interface{}{
		"include_irb_and_l2":           false,
		"keep_incoming_id":             false,
		"no_vlan_interface_name":       false,
		"prefix_host_name":             false,
		"prefix_routing_instance_name": false,
		"use_interface_description":    "",
		"use_option_82":                false,
		"use_option_82_strict":         false,
		"use_vlan_id":                  false,
	}
	if keepIncomingIDStrict {
		r["keep_incoming_id_strict"] = false
	}

	return r
}

func genForwardingOptionsDhcpRelayOption() map[string]interface{} {
	return map[string]interface{}{
		"option_15":                make([]map[string]interface{}, 0),
		"option_15_default_action": make([]map[string]interface{}, 0),
		"option_16":                make([]map[string]interface{}, 0),
		"option_16_default_action": make([]map[string]interface{}, 0),
		"option_60":                make([]map[string]interface{}, 0),
		"option_60_default_action": make([]map[string]interface{}, 0),
		"option_77":                make([]map[string]interface{}, 0),
		"option_77_default_action": make([]map[string]interface{}, 0),
		"option_order":             make([]string, 0),
	}
}

func genForwardingOptionsDhcpRelayOption82() map[string]interface{} {
	return map[string]interface{}{
		"circuit_id":                     make([]map[string]interface{}, 0),
		"exclude_relay_agent_identifier": false,
		"link_selection":                 false,
		"remote_id":                      make([]map[string]interface{}, 0),
		"server_id_override":             false,
		"vendor_specific_host_name":      false,
		"vendor_specific_location":       false,
	}
}

func readForwardingOptionsDhcpRelayAuthUsernameInclude(itemTrim string, authUsernameInclude map[string]interface{}) {
	switch {
	case itemTrim == "circuit-type":
		authUsernameInclude["circuit_type"] = true
	case balt.CutPrefixInString(&itemTrim, "client-id"):
		authUsernameInclude["client_id"] = true
		switch {
		case itemTrim == " exclude-headers":
			authUsernameInclude["client_id_exclude_headers"] = true
		case itemTrim == " use-automatic-ascii-hex-encoding":
			authUsernameInclude["client_id_use_automatic_ascii_hex_encoding"] = true
		}
	case balt.CutPrefixInString(&itemTrim, "delimiter "):
		authUsernameInclude["delimiter"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "domain-name "):
		authUsernameInclude["domain_name"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "interface-description "):
		authUsernameInclude["interface_description"] = itemTrim
	case itemTrim == "interface-name":
		authUsernameInclude["interface_name"] = true
	case itemTrim == "mac-address":
		authUsernameInclude["mac_address"] = true
	case itemTrim == "option-60":
		authUsernameInclude["option_60"] = true
	case balt.CutPrefixInString(&itemTrim, "option-82"):
		authUsernameInclude["option_82"] = true
		switch {
		case itemTrim == " circuit-id":
			authUsernameInclude["option_82_circuit_id"] = true
		case itemTrim == " remote-id":
			authUsernameInclude["option_82_remote_id"] = true
		}
	case itemTrim == "relay-agent-interface-id":
		authUsernameInclude["relay_agent_interface_id"] = true
	case itemTrim == "relay-agent-remote-id":
		authUsernameInclude["relay_agent_remote_id"] = true
	case itemTrim == "relay-agent-subscriber-id":
		authUsernameInclude["relay_agent_subscriber_id"] = true
	case itemTrim == "routing-instance-name":
		authUsernameInclude["routing_instance_name"] = true
	case balt.CutPrefixInString(&itemTrim, "user-prefix "):
		authUsernameInclude["user_prefix"] = strings.Trim(itemTrim, "\"")
	case itemTrim == "vlan-tags":
		authUsernameInclude["vlan_tags"] = true
	}
}

func readForwardingOptionsDhcpRelayOverridesV4(itemTrim string, overrides map[string]interface{}) (err error) {
	switch {
	case itemTrim == "allow-no-end-option":
		overrides["allow_no_end_option"] = true
	case itemTrim == "allow-snooped-clients":
		overrides["allow_snooped_clients"] = true
	case itemTrim == "always-write-giaddr":
		overrides["always_write_giaddr"] = true
	case itemTrim == "always-write-option-82":
		overrides["always_write_option_82"] = true
	case balt.CutPrefixInString(&itemTrim, "asymmetric-lease-time "):
		overrides["asymmetric_lease_time"], err = strconv.Atoi(itemTrim)
	case itemTrim == "bootp-support":
		overrides["bootp_support"] = true
	case balt.CutPrefixInString(&itemTrim, "client-discover-match "):
		overrides["client_discover_match"] = itemTrim
	case itemTrim == "delay-authentication":
		overrides["delay_authentication"] = true
	case itemTrim == "delete-binding-on-renegotiation":
		overrides["delete_binding_on_renegotiation"] = true
	case itemTrim == "disable-relay":
		overrides["disable_relay"] = true
	case balt.CutPrefixInString(&itemTrim, "dual-stack "):
		overrides["dual_stack"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "interface-client-limit "):
		overrides["interface_client_limit"], err = strconv.Atoi(itemTrim)
	case itemTrim == "layer2-unicast-replies":
		overrides["layer2_unicast_replies"] = true
	case itemTrim == "no-allow-snooped-clients":
		overrides["no_allow_snooped_clients"] = true
	case itemTrim == "no-bind-on-request":
		overrides["no_bind_on_request"] = true
	case itemTrim == "no-unicast-replies":
		overrides["no_unicast_replies"] = true
	case itemTrim == "proxy-mode":
		overrides["proxy_mode"] = true
	case balt.CutPrefixInString(&itemTrim, "relay-source "):
		overrides["relay_source"] = itemTrim
	case itemTrim == "replace-ip-source-with giaddr":
		overrides["replace_ip_source_with_giaddr"] = true
	case itemTrim == "send-release-on-delete":
		overrides["send_release_on_delete"] = true
	case itemTrim == "trust-option-82":
		overrides["trust_option_82"] = true
	case balt.CutPrefixInString(&itemTrim, "user-defined-option-82 "):
		overrides["user_defined_option_82"] = strings.Trim(itemTrim, "\"")
	}
	if err != nil {
		return fmt.Errorf(failedConvAtoiError, itemTrim, err)
	}

	return nil
}

func readForwardingOptionsDhcpRelayOverridesV6(itemTrim string, overrides map[string]interface{}) (err error) {
	switch {
	case itemTrim == "allow-snooped-clients":
		overrides["allow_snooped_clients"] = true
	case itemTrim == "always-process-option-request-option":
		overrides["always_process_option_request_option"] = true
	case balt.CutPrefixInString(&itemTrim, "asymmetric-lease-time "):
		overrides["asymmetric_lease_time"], err = strconv.Atoi(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "asymmetric-prefix-lease-time "):
		overrides["asymmetric_prefix_lease_time"], err = strconv.Atoi(itemTrim)
	case itemTrim == "client-negotiation-match incoming-interface":
		overrides["client_negotiation_match_incoming_interface"] = true
	case itemTrim == "delay-authentication":
		overrides["delay_authentication"] = true
	case itemTrim == "delete-binding-on-renegotiation":
		overrides["delete_binding_on_renegotiation"] = true
	case balt.CutPrefixInString(&itemTrim, "dual-stack "):
		overrides["dual_stack"] = strings.Trim(itemTrim, "\"")
	case balt.CutPrefixInString(&itemTrim, "interface-client-limit "):
		overrides["interface_client_limit"], err = strconv.Atoi(itemTrim)
	case itemTrim == "no-allow-snooped-clients":
		overrides["no_allow_snooped_clients"] = true
	case itemTrim == "no-bind-on-request":
		overrides["no_bind_on_request"] = true
	case balt.CutPrefixInString(&itemTrim, "relay-source "):
		overrides["relay_source"] = itemTrim
	case itemTrim == "send-release-on-delete":
		overrides["send_release_on_delete"] = true
	}
	if err != nil {
		return fmt.Errorf(failedConvAtoiError, itemTrim, err)
	}

	return nil
}

func readForwardingOptionsDhcpRelayAgentID(itemTrim string, agentID map[string]interface{}) {
	switch {
	case itemTrim == "include-irb-and-l2":
		agentID["include_irb_and_l2"] = true
	case itemTrim == "keep-incoming-interface-id":
		agentID["keep_incoming_id"] = true
	case itemTrim == "keep-incoming-remote-id":
		agentID["keep_incoming_id"] = true
	case itemTrim == "keep-incoming-interface-id strict":
		agentID["keep_incoming_id_strict"] = true
		agentID["keep_incoming_id"] = true
	case itemTrim == "no-vlan-interface-name":
		agentID["no_vlan_interface_name"] = true
	case itemTrim == "prefix host-name":
		agentID["prefix_host_name"] = true
	case itemTrim == "prefix routing-instance-name":
		agentID["prefix_routing_instance_name"] = true
	case balt.CutPrefixInString(&itemTrim, "use-interface-description "):
		agentID["use_interface_description"] = itemTrim
	case itemTrim == "use-option-82":
		agentID["use_option_82"] = true
	case itemTrim == "use-option-82 strict":
		agentID["use_option_82_strict"] = true
		agentID["use_option_82"] = true
	case itemTrim == "use-vlan-id":
		agentID["use_vlan_id"] = true
	}
}

func readForwardingOptionsDhcpRelayOption(itemTrim string, relayOption map[string]interface{}) error {
	switch {
	case balt.CutPrefixInString(&itemTrim, "option-15 default-action "):
		itemTrimFields := strings.Split(itemTrim, " ")
		defAction := map[string]interface{}{
			"action": itemTrimFields[0],
			"group":  "",
		}
		if len(itemTrimFields) > 1 { // <action> <group>
			defAction["group"] = strings.Trim(strings.Join(itemTrimFields[1:], " "), "\"")
		}
		relayOption["option_15_default_action"] = append(
			relayOption["option_15_default_action"].([]map[string]interface{}),
			defAction,
		)
	case balt.CutPrefixInString(&itemTrim, "option-15 "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 4 { // <compare> <value_type> <value> <action> <group>?
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "option-15", itemTrim)
		}
		value := itemTrimFields[2]
		actionIndex := 3
		if (strings.HasPrefix(itemTrimFields[2], "\"") && !strings.HasSuffix(itemTrimFields[2], "\"")) ||
			itemTrimFields[2] == "\"" {
			for k, v := range itemTrimFields[3:] {
				value += " " + v
				if strings.Contains(v, "\"") {
					actionIndex = 3 + k + 1

					break
				}
			}
		}
		value = html.UnescapeString(strings.Trim(value, "\""))
		action := itemTrimFields[actionIndex]
		option15 := map[string]interface{}{
			"compare":    itemTrimFields[0],
			"value_type": itemTrimFields[1],
			"value":      value,
			"action":     action,
			"group":      "",
		}
		if len(itemTrimFields) > actionIndex+1 {
			option15["group"] = strings.Trim(strings.Join(itemTrimFields[actionIndex+1:], " "), "\"")
		}
		relayOption["option_15"] = append(
			relayOption["option_15"].([]map[string]interface{}),
			option15,
		)
	case balt.CutPrefixInString(&itemTrim, "option-16 default-action "):
		itemTrimFields := strings.Split(itemTrim, " ")
		defAction := map[string]interface{}{
			"action": itemTrimFields[0],
			"group":  "",
		}
		if len(itemTrimFields) > 1 { // <action> <group>
			defAction["group"] = strings.Trim(strings.Join(itemTrimFields[1:], " "), "\"")
		}
		relayOption["option_16_default_action"] = append(
			relayOption["option_16_default_action"].([]map[string]interface{}),
			defAction,
		)
	case balt.CutPrefixInString(&itemTrim, "option-16 "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 4 { // <compare> <value_type> <value> <action> <group>?
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "option-16", itemTrim)
		}
		value := itemTrimFields[2]
		actionIndex := 3
		if (strings.HasPrefix(itemTrimFields[2], "\"") && !strings.HasSuffix(itemTrimFields[2], "\"")) ||
			itemTrimFields[2] == "\"" {
			for k, v := range itemTrimFields[3:] {
				value += " " + v
				if strings.Contains(v, "\"") {
					actionIndex = 3 + k + 1

					break
				}
			}
		}
		value = html.UnescapeString(strings.Trim(value, "\""))
		action := itemTrimFields[actionIndex]
		option16 := map[string]interface{}{
			"compare":    itemTrimFields[0],
			"value_type": itemTrimFields[1],
			"value":      value,
			"action":     action,
			"group":      "",
		}
		if len(itemTrimFields) > actionIndex+1 {
			option16["group"] = strings.Trim(strings.Join(itemTrimFields[actionIndex+1:], " "), "\"")
		}
		relayOption["option_16"] = append(
			relayOption["option_16"].([]map[string]interface{}),
			option16,
		)
	case balt.CutPrefixInString(&itemTrim, "option-60 default-action "):
		itemTrimFields := strings.Split(itemTrim, " ")
		defAction := map[string]interface{}{
			"action": itemTrimFields[0],
			"group":  "",
		}
		if len(itemTrimFields) > 1 { // <action> <group>
			defAction["group"] = strings.Trim(strings.Join(itemTrimFields[1:], " "), "\"")
		}
		relayOption["option_60_default_action"] = append(
			relayOption["option_60_default_action"].([]map[string]interface{}),
			defAction,
		)
	case balt.CutPrefixInString(&itemTrim, "option-60 "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 4 { // <compare> <value_type> <value> <action> <group>?
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "option-60", itemTrim)
		}
		value := itemTrimFields[2]
		actionIndex := 3
		if (strings.HasPrefix(itemTrimFields[2], "\"") && !strings.HasSuffix(itemTrimFields[2], "\"")) ||
			itemTrimFields[2] == "\"" {
			for k, v := range itemTrimFields[3:] {
				value += " " + v
				if strings.Contains(v, "\"") {
					actionIndex = 3 + k + 1

					break
				}
			}
		}
		value = html.UnescapeString(strings.Trim(value, "\""))
		action := itemTrimFields[actionIndex]
		option60 := map[string]interface{}{
			"compare":    itemTrimFields[0],
			"value_type": itemTrimFields[1],
			"value":      value,
			"action":     action,
			"group":      "",
		}
		if len(itemTrimFields) > actionIndex+1 {
			option60["group"] = strings.Trim(strings.Join(itemTrimFields[actionIndex+1:], " "), "\"")
		}
		relayOption["option_60"] = append(
			relayOption["option_60"].([]map[string]interface{}),
			option60,
		)
	case balt.CutPrefixInString(&itemTrim, "option-77 default-action "):
		itemTrimFields := strings.Split(itemTrim, " ")
		defAction := map[string]interface{}{
			"action": itemTrimFields[0],
			"group":  "",
		}
		if len(itemTrimFields) > 1 { // <action> <group>
			defAction["group"] = strings.Trim(strings.Join(itemTrimFields[1:], " "), "\"")
		}
		relayOption["option_77_default_action"] = append(
			relayOption["option_77_default_action"].([]map[string]interface{}),
			defAction,
		)
	case balt.CutPrefixInString(&itemTrim, "option-77 "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 4 { // <compare> <value_type> <value> <action> <group>?
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "option-77", itemTrim)
		}
		value := itemTrimFields[2]
		actionIndex := 3
		if (strings.HasPrefix(itemTrimFields[2], "\"") && !strings.HasSuffix(itemTrimFields[2], "\"")) ||
			itemTrimFields[2] == "\"" {
			for k, v := range itemTrimFields[3:] {
				value += " " + v
				if strings.Contains(v, "\"") {
					actionIndex = 3 + k + 1

					break
				}
			}
		}
		value = html.UnescapeString(strings.Trim(value, "\""))
		action := itemTrimFields[actionIndex]
		option77 := map[string]interface{}{
			"compare":    itemTrimFields[0],
			"value_type": itemTrimFields[1],
			"value":      value,
			"action":     action,
			"group":      "",
		}
		if len(itemTrimFields) > actionIndex+1 {
			option77["group"] = strings.Trim(strings.Join(itemTrimFields[actionIndex+1:], " "), "\"")
		}
		relayOption["option_77"] = append(
			relayOption["option_77"].([]map[string]interface{}),
			option77,
		)
	case balt.CutPrefixInString(&itemTrim, "option-order "):
		relayOption["option_order"] = append(relayOption["option_order"].([]string), itemTrim)
	}

	return nil
}

func readForwardingOptionsDhcpRelayOption82(itemTrim string, relayOption82 map[string]interface{}) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "circuit-id"):
		if len(relayOption82["circuit_id"].([]map[string]interface{})) == 0 {
			relayOption82["circuit_id"] = append(relayOption82["circuit_id"].([]map[string]interface{}), map[string]interface{}{
				"include_irb_and_l2":           false,
				"keep_incoming_circuit_id":     false,
				"no_vlan_interface_name":       false,
				"prefix_host_name":             false,
				"prefix_routing_instance_name": false,
				"use_interface_description":    "",
				"use_vlan_id":                  false,
				"user_defined":                 false,
				"vlan_id_only":                 false,
			})
		}
		if balt.CutPrefixInString(&itemTrim, " ") {
			circuitID := relayOption82["circuit_id"].([]map[string]interface{})[0]
			switch {
			case itemTrim == "include-irb-and-l2":
				circuitID["include_irb_and_l2"] = true
			case itemTrim == "keep-incoming-circuit-id":
				circuitID["keep_incoming_circuit_id"] = true
			case itemTrim == "no-vlan-interface-name":
				circuitID["no_vlan_interface_name"] = true
			case itemTrim == "prefix host-name":
				circuitID["prefix_host_name"] = true
			case itemTrim == "prefix routing-instance-name":
				circuitID["prefix_routing_instance_name"] = true
			case balt.CutPrefixInString(&itemTrim, "use-interface-description "):
				circuitID["use_interface_description"] = itemTrim
			case itemTrim == "use-vlan-id":
				circuitID["use_vlan_id"] = true
			case itemTrim == "user-defined":
				circuitID["user_defined"] = true
			case itemTrim == "vlan-id-only":
				circuitID["vlan_id_only"] = true
			}
		}
	case itemTrim == "exclude-relay-agent-identifier":
		relayOption82["exclude_relay_agent_identifier"] = true
	case itemTrim == "link-selection":
		relayOption82["link_selection"] = true
	case balt.CutPrefixInString(&itemTrim, "remote-id"):
		if len(relayOption82["remote_id"].([]map[string]interface{})) == 0 {
			relayOption82["remote_id"] = append(relayOption82["remote_id"].([]map[string]interface{}), map[string]interface{}{
				"hostname_only":                false,
				"include_irb_and_l2":           false,
				"keep_incoming_remote_id":      false,
				"no_vlan_interface_name":       false,
				"prefix_host_name":             false,
				"prefix_routing_instance_name": false,
				"use_interface_description":    "",
				"use_string":                   "",
				"use_vlan_id":                  false,
			})
		}
		if balt.CutPrefixInString(&itemTrim, " ") {
			remoteID := relayOption82["remote_id"].([]map[string]interface{})[0]
			switch {
			case itemTrim == "hostname-only":
				remoteID["hostname_only"] = true
			case itemTrim == "include-irb-and-l2":
				remoteID["include_irb_and_l2"] = true
			case itemTrim == "keep-incoming-remote-id":
				remoteID["keep_incoming_remote_id"] = true
			case itemTrim == "no-vlan-interface-name":
				remoteID["no_vlan_interface_name"] = true
			case itemTrim == "prefix host-name":
				remoteID["prefix_host_name"] = true
			case itemTrim == "prefix routing-instance-name":
				remoteID["prefix_routing_instance_name"] = true
			case balt.CutPrefixInString(&itemTrim, "use-interface-description "):
				remoteID["use_interface_description"] = itemTrim
			case balt.CutPrefixInString(&itemTrim, "use-string "):
				remoteID["use_string"] = strings.Trim(itemTrim, "\"")
			case itemTrim == "use-vlan-id":
				remoteID["use_vlan_id"] = true
			}
		}
	case itemTrim == "server-id-override":
		relayOption82["server_id_override"] = true
	case itemTrim == "vendor-specific host-name":
		relayOption82["vendor_specific_host_name"] = true
	case itemTrim == "vendor-specific location":
		relayOption82["vendor_specific_location"] = true
	}
}
