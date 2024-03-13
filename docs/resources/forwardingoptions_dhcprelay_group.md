---
page_title: "Junos: junos_forwardingoptions_dhcprelay_group"
---

# junos_forwardingoptions_dhcprelay_group

Provides a DHCP (or DHCPv6) relay group.

## Example Usage

```hcl
# Add a dhcp-relay group
resource "junos_forwardingoptions_dhcprelay_group" "demo" {
  name                = "demo"
  active_server_group = junos_forwardingoptions_dhcprelay_servergroup.demo.name
}

resource "junos_forwardingoptions_dhcprelay_servergroup" "demo" {
  name = "demo"
  ip_address = [
    "192.0.2.8",
  ]
}
```

## Argument Reference

-> **Note:** At least one of arguments need to be set
(in addition to `name`, `routing_instance` and `version`).

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Group name.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for group.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`
- **version** (Optional, String, Forces new resource)  
  Version for DHCP or DHCPv6.  
  Need to be `v4` or `v6`.
- **access_profile** (Optional, String)  
  Access profile to use for AAA services.
- **active_server_group** (Optional, String)  
  Name of DHCP server group.
- **active_server_group_allow_server_change** (Optional, Boolean)  
  Accept DHCP-ACK from any server in this group.  
  `version` need to be `v4`.
- **authentication_password** (Optional, String)  
  DHCP authentication, username password to use.
- **authentication_username_include** (Optional, Block)  
  DHCP authentication, add username options.  
  At least one of arguments of block need to be set.
  - **circuit_type** (Optional, Boolean)  
    Include circuit type.
  - **client_id** (Optional, Boolean)  
    Include client ID.
  - **client_id_exclude_headers** (Optional, Boolean)  
    Exclude all the headers.  
    `client_id` need to be true.
  - **client_id_use_automatic_ascii_hex_encoding** (Optional, Boolean)  
    Use automatic ascii hex username encoding.  
    `client_id` need to be true.
  - **delimiter** (Optional, String)  
    Change delimiter/separator character.  
    One character maximum.
  - **domain_name** (Optional, String)  
    Add domain name.
  - **interface_description** (Optional, String)  
    Include interface description.  
    Need to be `device` or `logical`.
  - **interface_name** (Optional, Boolean)  
    Include interface name.
  - **mac_address** (Optional, Boolean)  
    Include MAC address.
  - **option_60** (Optional, Boolean)  
    Include option 60.  
    `version` need to be `v4`.
  - **option_82** (Optional, Boolean)  
    Include option 82.  
    `version` need to be `v4`.
  - **option_82_circuit_id** (Optional, Boolean)  
    Include option 82 circuit-id (sub option 1).  
    `option_82` need to be true.
  - **option_82_remote_id** (Optional, Boolean)  
    Include option 82 remote-id (sub option 2).  
    `option_82` need to be true.
  - **relay_agent_interface_id** (Optional, Boolean)  
    Include the relay agent interface ID.  
    `version` need to be `v6`.
  - **relay_agent_remote_id** (Optional, Boolean)  
    Include the relay agent remote ID.  
    `version` need to be `v6`.
  - **relay_agent_subscriber_id** (Optional, Boolean)  
    Include the relay agent subscriber ID.  
    `version` need to be `v6`.
  - **routing_instance_name** (Optional, Boolean)  
    Include routing instance name.
  - **user_prefix** (Optional, String)  
    Add user defined prefix.
  - **vlan_tags** (Optional, Boolean)  
    Include the vlan tag(s).
- **client_response_ttl** (Optional, Number)  
  IP time-to-live value to set in responses to client (1..255).  
  `version` need to be `v4`.
- **description** (Optional, String)  
  Description.
- **dynamic_profile** (Optional, String)  
  Dynamic profile to use.
- **dynamic_profile_aggregate_clients** (Optional, Boolean)  
  Aggregate client profiles.  
  `dynamic_profile` need to be set.
- **dynamic_profile_aggregate_clients_action** (Optional, String)  
  Merge or replace the client dynamic profiles.  
  Need to be `merge` or `replace`.  
  `dynamic_profile_aggregate_clients` need to be true.
- **dynamic_profile_use_primary** (Optional, String)  
  Dynamic profile to use on the primary interface.  
  `dynamic_profile` need to be set.  
  Conflict with `dynamic_profile_aggregate_clients`.
- **forward_only** (Optional, Boolean)  
  Forward DHCP packets without creating binding.
- **forward_only_routing_instance** (Optional, String)  
  Name of routing instance to forward-only.
- **interface** (Optional, Block Set)  
  For each name of interface to declare.  
  - **name** (Required, String)  
    Interface name.
  - **access_profile** (Optional, String)  
    Access profile to use for AAA services.
  - **dynamic_profile** (Optional, String)  
    Dynamic profile to use.
  - **dynamic_profile_aggregate_clients** (Optional, String)  
    Aggregate client profiles.  
    `dynamic_profile` need to be set.
  - **dynamic_profile_aggregate_clients_action** (Optional, String)  
    Merge or replace the client dynamic profiles.  
    Need to be `merge` or `replace`.  
    `dynamic_profile_aggregate_clients` need to be set.
  - **dynamic_profile_use_primary** (Optional, String)  
    Dynamic profile to use on the primary interface.  
    `dynamic_profile` need to be set.  
    Conflict with `dynamic_profile_aggregate_clients`.
  - **exclude** (Optional, Boolean)  
    Exclude this interface range.
  - **overrides_v4** (Optional, Block)  
    DHCP override processing.  
    `version` need to be `v4`.  
    See [below for nested schema](#overrides_v4-arguments).
  - **overrides_v6** (Optional, Block)  
    DHCPv6 override processing.  
    `version` need to be `v6`.  
    See [below for nested schema](#overrides_v6-arguments).
  - **service_profile** (Optional, String)  
    Dynamic profile to use for default service activation.
  - **short_cycle_protection_lockout_max_time** (Optional, Number)  
    Short cycle lockout max time in seconds (1..86400).
  - **short_cycle_protection_lockout_min_time** (Optional, Number)  
    Short cycle lockout min time in seconds (1..86400).
  - **trace** (Optional, Boolean)  
    Enable tracing for this interface.
  - **upto** (Optional, String)  
    Interface up to.
- **lease_time_validation** (Optional, Block)  
  Configure lease time violation validation.  
  - **lease_time_threshold** (Optional, Number)  
    Threshold for lease time violation seconds (60..2147483647 seconds).
  - **violation_action_drop** (Optional, Boolean)  
    Lease time validation violation action is drop.
- **liveness_detection_failure_action** (Optional, String)  
  Liveness detection failure action options.  
  Need to be `clear-binding`, `clear-binding-if-interface-up` or `log-only`.
- **liveness_detection_method_bfd** (Optional, Block)  
  Liveness detection method BFD options.  
  At least one of arguments of block need to be set.  
  Conflict with `liveness_detection_method_layer2`.
  - **detection_time_threshold** (Optional, Number)  
    High detection-time triggering a trap (milliseconds).
  - **holddown_interval** (Optional, Number)  
    Time to hold the session-UP notification to the client (0..255000 milliseconds).
  - **minimum_interval** (Optional, Number)  
    Minimum transmit and receive interval (30000..255000 milliseconds).
  - **minimum_receive_interval** (Optional, Number)  
    Minimum receive interval (30000..255000 milliseconds).
  - **multiplier** (Optional, Number)  
    Detection time multiplier (1..255).
  - **no_adaptation** (Optional, Boolean)  
    Disable adaptation.
  - **session_mode** (Optional, String)  
    BFD single-hop or multihop session-mode.  
    Need to be `automatic`, `multihop` or `single-hop`.
  - **transmit_interval_minimum** (Optional, Number)  
    Minimum transmit interval (30000..255000 milliseconds).
  - **transmit_interval_threshold** (Optional, Number)  
    High transmit interval triggering a trap (milliseconds)
  - **version** (Optional, String)  
    BFD protocol version number.  
    Need to be `0`, `1` or `automatic`.
- **liveness_detection_method_layer2** (Optional, Block)  
  Liveness detection method address resolution options.  
  At least one of arguments of block need to be set.  
  Conflict with `liveness_detection_method_bfd`.
  - **max_consecutive_retries** (Optional, Number)  
    Retry attempts (3..6).
  - **transmit_interval** (Optional, Number)  
    Transmit interval for address resolution (300..1800 seconds).
- **maximum_hop_count** (Optional, Number)  
  Maximum number of hops per packet (1..16)  
  `version` need to be `v4`.
- **minimum_wait_time** (Optional, Number)  
  Minimum number of seconds before requests are forwarded (0..30000).  
  `version` need to be `v4`.
- **overrides_v4** (Optional, Block)  
  DHCP override processing.  
  `version` need to be `v4`.  
  See [below for nested schema](#overrides_v4-arguments).
- **overrides_v6** (Optional, Block)  
  DHCPv6 override processing.  
  `version` need to be `v6`.  
  See [below for nested schema](#overrides_v6-arguments).
- **relay_agent_interface_id** (Optional, Block)  
  DHCPv6 interface-id option processing.  
  `version` need to be `v6`.  
  See [below for nested schema](#relay_agent_interface_id-or-relay_agent_remote_id-arguments).
- **relay_agent_option_79** (Optional, Boolean)  
  Add the client MAC address to the Relay Forward header.  
  `version` need to be `v6`.
- **relay_agent_remote_id** (Optional, Block)  
  DHCPv6 remote-id option processing.  
  `version` need to be `v6`.  
  See [below for nested schema](#relay_agent_interface_id-or-relay_agent_remote_id-arguments)
  but without `keep_incoming_id_strict`.
- **relay_option** (Optional, Block)  
  DHCP option processing.  
  See [below for nested schema](#relay_option-arguments).
- **relay_option_82** (Optional, Block)  
  DHCP option-82 processing.  
  `version` need to be `v4`.  
  See [below for nested schema](#relay_option_82-arguments).
- **remote_id_mismatch_disconnect** (Optional, Boolean)  
  Disconnect session on remote-id mismatch.
- **route_suppression_access** (Optional, Boolean)  
  Suppress access route addition.  
  `version` need to be `v6`.
- **route_suppression_access_internal** (Optional, Boolean)  
  Suppress access-internal route addition.
- **route_suppression_destination** (Optional, Boolean)  
  Suppress destination route addition.  
  `version` need to be `v4`.
- **server_match_address** (Optional, Block Set)  
  For each `address`, server match processing.
  - **address** (Required, String)  
    Server address.
  - **action**  (Required, String)  
    Action on address.  
    Need to be `create-relay-entry` or `forward-only`.
- **server_match_default_action** (Optional, String)  
  Server match default action.  
  Need to be `create-relay-entry` or `forward-only`.
- **server_match_duid** (Optional, Block Set)  
  For each combination of `compare`, `value_type` and `value` arguments,  match duid processing.  
  `version` need to be `v6`.
  - **compare** (Required, String)  
    How to compare.  
    Need to be `equals` or `starts-with`.
  - **value_type** (Required, String)  
    Type of string.  
    Need to be `ascii` or `hexadecimal`.
  - **value** (Required, String)  
    String to compare.
  - **action** (Required, String)  
    Action on match.  
    Need to be `create-relay-entry` or `forward-only`.
- **service_profile** (Optional, String)  
  Dynamic profile to use for default service activation.
- **short_cycle_protection_lockout_max_time** (Optional, Number)  
  Short cycle lockout max time in seconds (1..86400).  
  `short_cycle_protection_lockout_min_time` need to be set.
- **short_cycle_protection_lockout_min_time** (Optional, Number)  
  Short cycle lockout min time in seconds (1..86400).  
  `short_cycle_protection_lockout_max_time` need to be set.
- **source_ip_change** (Optional, Boolean)  
  Use address of egress interface as source ip.  
  `version` need to be `v4`.
- **vendor_specific_information_host_name** (Optional, Boolean)  
  DHCPv6 option 17 vendor-specific processing, add router host name.  
  `version` need to be `v6`.
- **vendor_specific_information_location** (Optional, Boolean)  
  DHCPv6 option 17 vendor-specific processing,
  add location information expressed as interface name format.  
  `version` need to be `v6`.

---

### overrides_v4 arguments

- **allow_no_end_option** (Optional, Boolean)  
  Allow packets without end-of-option.
- **allow_snooped_clients** (Optional, Boolean)  
  Allow client creation from snooped PDUs.
- **always_write_giaddr** (Optional, Boolean)  
  Overwrite existing 'giaddr' field, when present.
- **always_write_option_82** (Optional, Boolean)  
  Overwrite existing value of option 82, when present.
- **asymmetric_lease_time** (Optional, Number)  
  Use a reduced lease time for the client. In seconds (600..86400 seconds).
- **bootp_support** (Optional, Boolean)  
  Allows relay of bootp req and reply.
- **client_discover_match** (Optional, String)  
  Use secondary match criteria for DISCOVER PDU.  
  Need to be `incoming-interface` or `option60-and-option82`.
- **delay_authentication** (Optional, Boolean)  
  Delay subscriber authentication in DHCP protocol processing until request packet.
- **delete_binding_on_renegotiation** (Optional, Boolean)  
  Delete binding on rengotiation.
- **disable_relay** (Optional, Boolean)  
  Disable DHCP relay processing.
- **dual_stack** (Optional, String)  
  Dual stack group to use.
- **interface_client_limit** (Optional, Number)  
  Limit the number of clients allowed on an interface (1..500000).
- **layer2_unicast_replies** (Optional, Boolean)  
  Do not broadcast client responses.
- **no_allow_snooped_clients** (Optional, Boolean)  
  Don't allow client creation from snooped PDUs.
- **no_bind_on_request** (Optional, Boolean)  
  Do not bind if stray DHCP request is received.
- **no_unicast_replies** (Optional, Boolean)  
  Overwrite unicast bit in incoming packet, when present.
- **proxy_mode** Optional, Boolean)  
  Put the relay in proxy mode.
- **relay_source** (Optional, String)  
  Interface for relay source.
- **replace_ip_source_with_giaddr** (Optional, Boolean)  
  Replace IP source address in request and release packets.
- **send_release_on_delete** (Optional, Boolean)  
  Always send RELEASE to the server when a binding is deleted.
- **trust_option_82** (Optional, Boolean)  
  Trust options-82 option.
- **user_defined_option_82** (Optional, String)  
  Set user defined description for option-82.

### overrides_v6 arguments

- **allow_snooped_clients** (Optional, Boolean)  
  Allow client creation from snooped PDUs.
- **always_process_option_request_option** (Optional, Boolean)  
  Always process option even after address allocation failure.
- **asymmetric_lease_time** (Optional, Number)  
  Use a reduced lease time for the client. In seconds (600..86400 seconds).
- **asymmetric_prefix_lease_time** (Optional, Number)  
  Use a reduced prefix lease time for the client. In seconds (600..86400 seconds).
- **client_negotiation_match_incoming_interface** (Optional, Boolean)  
  Use incoming interface match criteria for SOLICIT PDU.
- **delay_authentication** (Optional, Boolean)  
  Delay subscriber authentication in DHCP protocol processing until request packet.
- **delete_binding_on_renegotiation** (Optional, Boolean)  
  Delete binding on rengotiation.
- **dual_stack** (Optional, String)  
  Dual stack group to use.
- **interface_client_limit** (Optional, Number)  
  Limit the number of clients allowed on an interface (1..500000).
- **no_allow_snooped_clients** (Optional, Boolean)  
  Don't allow client creation from snooped PDUs.
- **no_bind_on_request** (Optional, Boolean)  
  Do not bind if stray DHCPv6 RENEW, REBIND is received.
- **relay_source** (Optional, String)  
  Interface for relay source.
- **send_release_on_delete** (Optional, Boolean)  
  Always send RELEASE to the server when a binding is deleted.

### relay_agent_interface_id or relay_agent_remote_id arguments

- **include_irb_and_l2** (Optional, Boolean)  
  Include IRB and L2 interface name.
- **keep_incoming_id** (Optional, Boolean)  
  Keep incoming interface identifier.
- **keep_incoming_id_strict** (Optional, Boolean)  
  Drop packet if interface identifier not present.  
  Only on `relay_agent_interface_id` block.
- **no_vlan_interface_name** (Optional, Boolean)  
  Not include vlan or interface name.
- **prefix_host_name** (Optional, Boolean)  
  Add router host name to circuit / interface-id or remote-id.
- **prefix_routing_instance_name** (Optional, Boolean)  
  Add routing instance name to circuit / interface-id or remote-id.
- **use_interface_description** (Optional, String)  
  Use interface description instead of circuit identifier.  
  Need to be `device` or `logical`.
- **use_option_82** (Optional, Boolean)  
  Use option-82 circuit-id for interface-id or remote-id.
- **use_option_82_strict** (Optional, Boolean)  
  Drop packet if option-82 circuit-id not present.
- **use_vlan_id** (Optional, Boolean)  
  Use VLAN id instead of name.

### relay_option arguments

- **option_15** (Optional, Block Set)  
  For each combination of `compare`, `value_type` and `value` arguments, add option 15 processing.  
  `version` need to be `v6`.
  - **compare** (Required, String)  
    How to compare.  
    Need to be `equals` or `starts-with`.
  - **value_type** (Required, String)  
    Type of string.  
    Need to be `ascii` or `hexadecimal`.
  - **value** (Required, String)  
    String to compare.
  - **action** (Required, String)  
    Action on match.  
    Need to be `drop`, `forward-only` or `relay-server-group`.
  - **group** (Optional, String)  
    Group for action `relay-server-group`.
- **option_15_default_action** (Optional, Block)  
  Generic option 15 default action.  
  `version` need to be `v6`.
  - **action** (Required, String)  
    Action.  
    Need to be `drop`, `forward-only` or `relay-server-group`.
  - **group** (Optional, String)  
    Group for action `relay-server-group`.
- **option_16** (Optional, Block Set)  
  For each combination of `compare`, `value_type` and `value` arguments, add option 16 processing.  
  `version` need to be `v6`.
  - **compare** (Required, String)  
    How to compare.  
    Need to be `equals` or `starts-with`.
  - **value_type** (Required, String)  
    Type of string.  
    Need to be `ascii` or `hexadecimal`.
  - **value** (Required, String)  
    String to compare.
  - **action** (Required, String)  
    Action on match.  
    Need to be `drop`, `forward-only` or `relay-server-group`.
  - **group** (Optional, String)  
    Group for action `relay-server-group`.
- **option_16_default_action** (Optional, Block)  
  Generic option 16 default action.  
  `version` need to be `v6`.
  - **action** (Required, String)  
    Action.  
    Need to be `drop`, `forward-only` or `relay-server-group`.
  - **group** (Optional, String)  
    Group for action `relay-server-group`.
- **option_60** (Optional, Block Set)  
  For each combination of `compare`, `value_type` and `value` arguments, add option 60 processing.  
  `version` need to be `v4`.
  - **compare** (Required, String)  
    How to compare.  
    Need to be `equals` or `starts-with`.
  - **value_type** (Required, String)  
    Type of string.  
    Need to be `ascii` or `hexadecimal`.
  - **value** (Required, String)  
    String to compare.
  - **action** (Required, String)  
    Action on match.  
    Need to be `drop`, `forward-only`, `local-server-group` or `relay-server-group`.
  - **group** (Optional, String)  
    Group for action `local-server-group` or `relay-server-group`.
- **option_60_default_action** (Optional, Block)  
  Generic option 60 default action.  
  `version` need to be `v4`.
  - **action** (Required, String)  
    Action.  
    Need to be `drop`, `forward-only`, `local-server-group` or `relay-server-group`.
  - **group** (Optional, String)  
    Group for action `local-server-group` or `relay-server-group`.
- **option_77** (Optional, Block Set)  
  For each combination of `compare`, `value_type` and `value` arguments, add option 77 processing.  
  `version` need to be `v4`.
  - **compare** (Required, String)  
    How to compare.  
    Need to be `equals` or `starts-with`.
  - **value_type** (Required, String)  
    Type of string.  
    Need to be `ascii` or `hexadecimal`.
  - **value** (Required, String)  
    String to compare.
  - **action** (Required, String)  
    Action on match.  
    Need to be `drop`, `forward-only`, `local-server-group` or `relay-server-group`.
  - **group** (Optional, String)  
    Group for action `local-server-group` or `relay-server-group`.
- **option_77_default_action** (Optional, Block)  
  Generic option 77 default action.  
  `version` need to be `v4`.
  - **action** (Required, String)  
    Action.  
    Need to be `drop`, `forward-only`, `local-server-group` or `relay-server-group`.
  - **group** (Optional, String)  
    Group for action `local-server-group` or `relay-server-group`.
- **option_order** (Optional, List of String)  
  Options precedence order.  
  Need to be `60` or `77` with `version` = `v4`.  
  Need to be `15` or `16` with `version` = `v6`.

### relay_option_82 arguments

- **circuit_id** (Optional, Block)  
  Add circuit identifier.
  - **include_irb_and_l2** (Optional, Boolean)  
    Include IRB and L2 interface name.
  - **keep_incoming_circuit_id** (Optional, Boolean)  
    Keep incoming circuit identifier.
  - **no_vlan_interface_name** (Optional, Boolean)  
    Not include vlan or interface name.
  - **prefix_host_name** (Optional, Boolean)  
    Add router host name to circuit / interface-id or remote-id.  
  - **prefix_routing_instance_name** (Optional, Boolean)  
    Add routing instance name to circuit / interface-id or remote-id.
  - **use_interface_description** (Optional, String)  
    Use interface description instead of circuit identifier.  
    Need to be `device` or `logical`.
  - **use_vlan_id** (Optional, Boolean)  
    Use VLAN id instead of name.
  - **user_defined** (Optional, Boolean)  
    Include user defined string.
  - **vlan_id_only** (Optional, Boolean)  
    Use only VLAN id.
- **exclude_relay_agent_identifier** (Optional, Boolean)  
  Exclude relay agent identifier from packets to server.
- **link_selection** (Optional, Boolean)  
  Add link-selection sub-option on packets to server.
- **remote_id** (Optional, Block)  
  Add remote identifier.
  - **hostname_only** (Optional, Boolean)  
    Include hostname only.
  - **include_irb_and_l2** (Optional, Boolean)  
    Include IRB and L2 interface name.
  - **keep_incoming_remote_id** (Optional, Boolean)  
    Keep incoming remote identifier.
  - **no_vlan_interface_name** (Optional, Boolean)  
    Not include vlan or interface name.
  - **prefix_host_name** (Optional, Boolean)  
    Add router host name to circuit / interface-id or remote-id.  
  - **prefix_routing_instance_name** (Optional, Boolean)  
    Add routing instance name to circuit / interface-id or remote-id.
  - **use_interface_description** (Optional, String)  
    Use interface description instead of circuit identifier.  
    Need to be `device` or `logical`.
  - **use_string** (Optional, String)  
    Use raw string instead of the default remote id.
  - **use_vlan_id** Optional, String)  
    Use VLAN id instead of name.
- **server_id_override** (Optional, Boolean)  
  Add link-selection and server-id sub-options on packets to server.
- **vendor_specific_host_name** (Optional, Boolean)  
  Add vendor-specific information, add router host name.
- **vendor_specific_location** (Optional, Boolean)  
  Add vendor-specific information, add location information expressed as interface name format.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>_-_<version>`.

## Import

Junos forwarding-options dhcp-relay group can be imported using an id made up of
`<name>_-_<routing_instance>_-_<version>`, e.g.

```shell
$ terraform import junos_forwardingoptions_dhcprelay_group.demo demo_-_default_-_v4
```
