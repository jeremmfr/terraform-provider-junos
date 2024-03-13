---
page_title: "Junos: junos_bgp_group"
---

# junos_bgp_group

Provides a bgp group resource.

## Example Usage

```hcl
# Configure a bgp group
resource "junos_bgp_group" "groupbgpdemo" {
  name             = "GroupBgpDemo"
  routing_instance = "default"
  peer_as          = "65002"
  local_address    = "192.0.2.3"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of group.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for bgp protocol if not root level.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`.
- **type** (Optional, String, Forces new resource)  
  Type of peer group.  
  Need to be `internal` or `external`.  
  Defaults to `external`.
- **accept_remote_nexthop** (Optional, Boolean)  
  Allow import policy to specify a non-directly connected next-hop.
- **advertise_external** (Optional, Computed, Boolean)  
  Advertise best external routes.  
  Computed to set to `true` when `advertise_external_conditional` is true.
- **advertise_external_conditional** (Optional, Boolean)  
  Route matches active route upto med-comparison rule.
- **advertise_inactive** (Optional, Boolean)  
  Advertise inactive routes.
- **advertise_peer_as** (Optional, Boolean)  
  Advertise routes received from the same autonomous system.  
  Conflict with `no_advertise_peer_as`.
- **no_advertise_peer_as** (Optional, Boolean)  
  Don't advertise routes received from the same autonomous system.  
  Conflict with `advertise_peer_as`.
- **as_override** (Optional, Boolean)  
  Replace neighbor AS number with our AS number.
- **authentication_algorithm** (Optional, String)  
  Authentication algorithm name.  
  Conflict with `authentication_key`.
- **authentication_key** (Optional, String, Sensitive)  
  MD5 authentication key.  
  Conflict with `authentication_*`.
- **authentication_key_chain** (Optional, String)  
  Key chain name.  
  Conflict with `authentication_key`.
- **bfd_liveness_detection** (Optional, Block)  
  Define Bidirectional Forwarding Detection (BFD) options.  
  See [below for nested schema](#bfd_liveness_detection-arguments).
- **bgp_error_tolerance** (Optional, Block)  
  Handle BGP malformed updates softly.
  - **malformed_route_limit** (Optional, Number)  
    Maximum number of malformed routes from a peer (0..4294967295).  
    Conflict with `no_malformed_route_limit`.
  - **malformed_update_log_interval** (Optional, Number)  
    Time used when logging malformed update (10..65535 seconds).
  - **no_malformed_route_limit** (Optional, Boolean)  
    No malformed route limit.  
    Conflict with `malformed_route_limit`.
- **bgp_multipath** (Optional, Block)  
  Allow load sharing among multiple BGP paths.
  - **allow_protection** (Optional, Boolean)  
    Allows the BGP multipath and protection to co-exist.
  - **disable** (Optional, Boolean)  
    Disable Multipath.
  - **multiple_as** (Optional, Boolean)  
    Use paths received from different ASs.
- **cluster** (Optional, String)  
  Cluster identifier.  
  Must be a valid IP address.
- **damping** (Optional, Boolean)  
  Enable route flap damping.
- **description** (Optional, String)  
  Text description.
- **export** (Optional, List of String)  
  Export policy list.
- **family_evpn** (Optional, Block List)  
  For each `nlri_type`, configure EVPN NLRI parameters.
  - **nlri_type** (Optional, String)  
    NLRI type.  
    Need to be `signaling`.  
    Default to `signaling`.
  - other options same as [`family_inet` arguments](#family_inet-arguments).
- **family_inet** (Optional, Block List)  
  For each `nlri_type`, configure IPv4 NLRI parameters.  
  See [below for nested schema](#family_inet-arguments).
- **family_inet6** (Optional, Block List)  
  For each `nlri_type`, configure IPv6 NLRI parameters.  
  Same options as [`family_inet` arguments](#family_inet-arguments) but for inet6 family.
- **graceful_restart** (Optional, Block)  
  Define BGP graceful restart options.
  - **disable** (Optional, Boolean)  
    Disable graceful restart.
  - **restart_time** (Optional, Number)  
    Restart time used when negotiating with a peer (1..600).
  - **stale_route_time** (Optional, Number)  
    Maximum time for which stale routes are kept (1..600).
- **hold_time** (Optional, Number)  
  Hold time used when negotiating with a peer.
- **import** (Optional, List of String)  
  Import policy list.
- **keep_all** (Optional, Boolean)  
  Retain all routes.  
  Conflict with `keep_none`.
- **keep_none** (Optional, Boolean)  
  Retain no routes.  
  Conflict with `keep_all`.
- **local_address** (Optional, String)  
  Address of local end of BGP session.
- **local_as** (Optional, String)  
  Local autonomous system number.
- **local_as_alias** (Optional, Boolean)  
  Treat this AS as an alias to the system AS.  
  Conflict with other local_as options.
- **local_as_loops** (Optional, Number)  
  Maximum number of times this AS can be in an AS path (1..10).
- **local_as_no_prepend_global_as** (Optional, Boolean)  
  Do not prepend global autonomous-system number in advertised paths.  
  Conflict with other local_as options.
- **local_as_private** (Optional, Boolean)  
  Hide this local AS in paths learned from this peering.  
  Conflict with other local_as options.
- **local_interface** (Optional, String)  
  Local interface for IPv6 link local EBGP peering.
- **local_preference** (Optional, Number)  
  Value of LOCAL_PREF path attribute.
- **log_updown** (Optional, Boolean)  
  Log a message for peer state transitions.
- **metric_out** (Optional, Number)  
  Route metric sent in MED.
- **metric_out_igp** (Optional, Computed, Boolean)  
  Track the IGP metric.  
  Computed to set to `true` when `metric_out_igp_offset` or `metric_out_igp_delay_med_update`
  is set.  
  Conflict with `metric_out` and `metric_out_minimum_*`.
- **metric_out_igp_offset** (Optional, Number)  
  Metric offset for MED.  
  Conflict with `metric_out` and `metric_out_minimum_*`.
- **metric_out_igp_delay_med_update** (Optional, Boolean)  
  Delay updating MED when IGP metric increases.  
  Conflict with `metric_out` and `metric_out_minimum_*`.
- **metric_out_minimum_igp** (Optional, Computed, Boolean)  
  Track the minimum IGP metric.  
  Computed to set to `true` when `metric_out_minimum_igp_offset` is set.  
  Conflict with `metric_out` and `metric_out_(?!minimum)_*`.
- **metric_out_minimum_igp_offset** (Optional, Boolean)  
  Metric offset for MED.  
  Conflict with `metric_out` and `metric_out_(?!minimum)_*`.
- **mtu_discovery** (Optional, Boolean)  
  Enable TCP path MTU discovery.
- **multihop** (Optional, Boolean)  
  Configure an EBGP multihop session.
- **no_client_reflect** (Optional, Boolean)  
  Disable intracluster route redistribution.
- **out_delay** (Optional, Number)  
  How long before exporting routes from routing table.
- **passive** (Optional, Boolean)  
  Do not send open messages to a peer.
- **peer_as** (Optional, String)  
  Autonomous system number.
- **preference** (Optional, Number)  
  Preference value.
- **remove_private** (Optional, Boolean)  
  Remove well-known private AS numbers.
- **tcp_aggressive_transmission** (Optional, Boolean)  
  Enable aggressive transmission of pure TCP ACKs and retransmissions

---

### bfd_liveness_detection arguments

- **authentication_algorithm** (Optional, String)  
  Authentication algorithm name.
- **authentication_key_chain** (Optional, String)  
  Authentication key chain name.
- **authentication_loose_check** (Optional, Boolean)  
  Verify authentication only if authentication is negotiated.
- **detection_time_threshold** (Optional, Number)  
  High detection-time triggering a trap (milliseconds).
- **holddown_interval** (Optional, Number)  
  Time to hold the session-UP notification to the client (1..255000 milliseconds).
- **minimum_interval** (Optional, Number)  
  Minimum transmit and receive interval (1..255000 milliseconds).
- **minimum_receive_interval** (Optional, Number)  
  Minimum receive interval (1..255000 milliseconds).
- **multiplier** (Optional, Number)  
  Detection time multiplier (1..255).
- **session_mode** (Optional, String)  
  BFD single-hop or multihop session-mode.  
  Need to be `automatic`, `multihop` or `single-hop`.
- **transmit_interval_minimum_interval** (Optional, Number)  
  Minimum transmit interval (1..255000 milliseconds).
- **transmit_interval_threshold** (Optional, Number)  
  High transmit interval triggering a trap (milliseconds).
- **version** (Optional, String)  
  BFD protocol version number.  
  Need to be `0`, `1` or `automatic`.

---

### family_inet arguments

Also for `family_inet6` and `family_evpn` (except `nlri_type`)

- **nlri_type** (Required, String)  
  NLRI type.  
  Need to be `any`, `flow`, `labeled-unicast`, `unicast` or `multicast`.
- **accepted_prefix_limit** (Optional, Block)  
  Define maximum number of prefixes accepted from a peer.
  - **maximum** (Required, Number)  
    Maximum number of prefixes accepted from a peer (1..4294967295).
  - **teardown** (Optional, Number)  
    Clear peer connection on reaching limit with this percentage of
    prefix-limit to start warnings.
  - **teardown_idle_timeout** (Optional, Number)  
    Timeout before attempting to restart peer.
  - **teardown_idle_timeout_forever** (Optional, Boolean)  
    Idle the peer until the user intervenes.  
    Conflict with `teardown_idle_timeout`.
- **prefix_limit** (Optional, Block)  
  Same options as `accepted_prefix_limit` but for limit maximum number of prefixes from a peer.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>`.

## Import

Junos bgp group can be imported using an id made up of `<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_bgp_group.groupbgpdemo GroupBgpDemo_-_default
```
