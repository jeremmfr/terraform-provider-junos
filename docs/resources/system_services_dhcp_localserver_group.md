---
page_title: "Junos: junos_system_services_dhcp_localserver_group"
---

# junos_system_services_dhcp_localserver_group

Provides a DHCP (or DHCPv6) local server group.

## Example Usage

```hcl
# Add a DHCP local server group
resource "junos_system_services_dhcp_localserver_group" "demo_dhcp_group" {
  name = "demo_dhcp_group"
  interface {
    name = "ge-0/0/3.1"
  }
}
# Add a DHCPv6 local server group
resource "junos_system_services_dhcp_localserver_group" "demo_dhcp_group_v6" {
  name    = "demo_dhcp_group"
  version = "v6"
  interface {
    name = "ge-0/0/3.1"
  }
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
- **dynamic_profile** (Optional, String)  
  Dynamic profile to use.
- **dynamic_profile_use_primary** (Optional, String)  
  Dynamic profile to use on the primary interface.  
  `dynamic_profile` need to be set.  
  Conflict with `dynamic_profile_aggregate_clients`.
- **dynamic_profile_aggregate_clients** (Optional, Boolean)  
  Aggregate client profiles.
- **dynamic_profile_aggregate_clients_action** (Optional, String)  
  Merge or replace the client dynamic profiles.  
  Need to be `merge` or `replace`.  
  `dynamic_profile_aggregate_clients` need to be true.
- **interface** (Optional, Block Set)  
  For each name of interface to declare.  
  See [below for nested schema](#interface-arguments).
- **lease_time_validation** (Optional, Block)  
  Configure lease time violation validation.  
  - **lease_time_threshold** (Optional, Number)  
    Threshold for lease time violation seconds (60..2147483647 seconds).
  - **violation_action** (Optional, String)  
    Lease time validation violation action.  
    Need to be `override-lease` or `strict`.
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
- **overrides_v4** (Optional, Block)  
  DHCP override processing.  
  See [below for nested schema](#overrides_v4-arguments).
- **overrides_v6** (Optional, Block)  
  DHCPv6 override processing.  
  See [below for nested schema](#overrides_v6-arguments).
- **reauthenticate_lease_renewal** (Optional, Boolean)  
  Reauthenticate on each renew, rebind, DISCOVER or SOLICIT.
- **reauthenticate_remote_id_mismatch** (Optional, Boolean)  
  Reauthenticate on remote-id mismatch for renew, rebind and re-negotiation.
- **reconfigure** (Optional, Block)  
  DHCP reconfigure processing.
  - **attempts** (Optional, Number)  
    Number of reconfigure attempts before aborting (1..10).
  - **clear_on_abort** (Optional, Boolean)  
    Delete client on reconfiguration abort.
  - **support_option_pd_exclude** (Optional, Boolean)  
    Request prefix exclude option in reconfigure message
  - **timeout** (Optional, Number)  
    Initial timeout value for retry (1..10).
  - **token** (Optional, String)  
    Reconfigure token.
  - **trigger_radius_disconnect** (Optional, Boolean)  
    Trigger DHCP reconfigure by radius initiated disconnect.
- **remote_id_mismatch_disconnect** (Optional, Boolean)  
  Disconnect session on remote-id mismatch.
- **route_suppression_access** (Optional, Boolean)  
  Suppress access route addition.
- **route_suppression_access_internal** (Optional, Boolean)  
  Suppress access-internal route addition.
- **route_suppression_destination** (Optional, Boolean)  
  Suppress destination route addition.
- **service_profile** (Optional, String)  
  Dynamic profile to use for default service activation.
- **short_cycle_protection_lockout_max_time** (Optional, Number)  
  Short cycle lockout max time in seconds (1..86400).  
  `short_cycle_protection_lockout_min_time` need to be set.
- **short_cycle_protection_lockout_min_time** (Optional, Number)  
  Short cycle lockout min time in seconds (1..86400).  
  `short_cycle_protection_lockout_max_time` need to be set.

---

### interface arguments

- **name** (Required, String)  
  Interface name.  
  Need to be a logical interface or `all`.
- **access_profile** (Optional, String)  
  Access profile to use for AAA services.
- **dynamic_profile** (Optional, String)  
  Dynamic profile to use.
- **dynamic_profile_use_primary** (Optional, String)  
  Dynamic profile to use on the primary interface.  
  `dynamic_profile` need to be set.
  Conflict with `dynamic_profile_aggregate_clients`.
- **dynamic_profile_aggregate_clients** (Optional, Boolean)  
  Aggregate client profiles.
- **dynamic_profile_aggregate_clients_action** (Optional, String)  
  Merge or replace the client dynamic profiles.  
  Need to be `merge` or `replace`.  
  `dynamic_profile_aggregate_clients` need to be true.
- **exclude** (Optional, Boolean)  
  Exclude this interface range.
- **overrides_v4** (Optional, Block)  
  DHCP override processing.  
  See [below for nested schema](#overrides_v4-arguments).
- **overrides_v6** (Optional, Block)  
  DHCPv6 override processing.  
  See [below for nested schema](#overrides_v6-arguments).
- **service_profile** (Optional, String)  
  Dynamic profile to use for default service activation.
- **short_cycle_protection_lockout_max_time** (Optional, Number)  
  Short cycle lockout max time in seconds (1..86400).  
  `short_cycle_protection_lockout_min_time` need to be set.
- **short_cycle_protection_lockout_min_time** (Optional, Number)  
  Short cycle lockout min time in seconds (1..86400).  
  `short_cycle_protection_lockout_max_time` need to be set.
- **trace** (Optional, Boolean)  
  Enable tracing for this interface.
- **upto** (Optional, String)  
  Interface up to.  
  Need to be a logical interface.

---

### overrides_v4 arguments

- **allow_no_end_option** (Optional, Boolean)  
  Allow packets without end-of-option.
- **asymmetric_lease_time** (Optional, Number)  
  Use a reduced lease time for the client. In seconds (600..86400 seconds).
- **bootp_support** (Optional, Boolean)  
  Allow processing of bootp requests.
- **client_discover_match** (Optional, String)  
  Use incoming interface or option 60 and option 82 match criteria for DISCOVER PDU.  
  Need to be `incoming-interface` or `option60-and-option82`.
- **delay_offer_based_on** (Optional, Block Set)  
  For each combination of block arguments, filter options for dhcp-server.
  - **option** (Required, String)  
    Option.  
    Need to be `option-60`, `option-77` or `option-82`.
  - **compare** (Required, String)  
    How to compare.  
    Need to be `equals`, `not-equals` or `starts-with`.
  - **value_type** (Required, String)  
    Type of string.  
    Need to be `ascii` or `hexadecimal`.
  - **value** (Required, String)  
    String to compare.
- **delay_offer_delay_time** (Optional, Number)  
  Time delay between discover and offer (1..30 seconds).  
  `delay_offer_based_on` need to be set at least one time.
- **delete_binding_on_renegotiation** (Optional, Boolean)  
  Delete binding on renegotiation.
- **dual_stack** (Optional, String)  
  Dual stack group to use.
- **include_option_82_forcerenew** (Optional, Boolean)  
  Include option-82 in FORCERENEW.
- **include_option_82_nak** (Optional, Boolean)  
  Include option-82 in NAK.
- **interface_client_limit** (Optional, Number)  
  Limit the number of clients allowed on an interface (1..500000).
- **process_inform** (Optional, Boolean)  
  Process INFORM PDUs.
- **process_inform_pool** (Optional, String)  
  Pool name for family inet.
- **protocol_attributes** (Optional, String)  
  DHCPv4 attributes to use as defined under access protocol-attributes.

### overrides_v6 arguments

- **always_add_option_dns_server** (Optional, Boolean)  
  Add option-23, DNS recursive name server in Advertise and Reply.
- **always_process_option_request_option** (Optional, Boolean)  
  Always process option even after address allocation failure.
- **asymmetric_lease_time** (Optional, Number)  
  Use a reduced lease time for the client. In seconds (600..86400 seconds).
- **asymmetric_prefix_lease_time** (Optional, Number)  
  Use a reduced prefix lease time for the client. In seconds (600..86400 seconds).
- **client_negotiation_match_incoming_interface** (Optional, Boolean)  
  Use incoming interface match criteria for SOLICIT PDU
- **delay_advertise_based_on** (Optional, Block Set)  
  For each combination of block arguments, filter options for dhcp-server.
  - **option** (Required, String)  
    Option.  
    Need to be `option-60`, `option-77` or `option-82`.
  - **compare** (Required, String)  
    How to compare.  
    Need to be `equals`, `not-equals` or `starts-with`.
  - **value_type** (Required, String)  
    Type of string.  
    Need to be `ascii` or `hexadecimal`.
  - **value** (Required, String)  
    String to compare.
- **delay_advertise_delay_time** (Optional, Number)  
  Time delay between solicit and advertise (1..30 seconds).  
  `delay_advertise_based_on` need to be set at least one time.
- **delegated_pool** (Optional, String)  
  Delegated pool name for inet6.
- **delete_binding_on_renegotiation** (Optional, Boolean)  
  Delete binding on renegotiation.
- **dual_stack** (Optional, String)  
  Dual stack group to use.
- **interface_client_limit** (Optional, Number)  
  Limit the number of clients allowed on an interface (1..500000).
- **multi_address_embedded_option_response** (Optional, Boolean)  
  If the client requests multiple addresses place the options in each address.
- **process_inform** (Optional, Boolean)  
  Process INFORM PDUs.
- **process_inform_pool** (Optional, String)  
  Pool name for family inet.
- **protocol_attributes** (Optional, String)  
  DHCPv6 attributes to use as defined under access protocol-attributes.
- **rapid_commit** (Optional, Boolean)  
  Enable rapid commit processing.
- **top_level_status_code** (Optional, Boolean)  
  A top level status code option rather than encapsulated in IA for NoAddrsAvail in Advertise PDUs.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>_-_<version>`.

## Import

Junos system DHCP local server group can be imported using an id made up of
`<name>_-_<routing_instance>_-_<version>`, e.g.

```shell
$ terraform import junos_system_services_dhcp_localserver_group.demo_dhcp_group demo_dhcp_group_-_default_-_v4
```
