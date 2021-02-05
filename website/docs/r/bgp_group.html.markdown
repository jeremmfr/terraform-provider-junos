---
layout: "junos"
page_title: "Junos: junos_bgp_group"
sidebar_current: "docs-junos-resource-bgp-group"
description: |-
  Create a BGP group
---

# junos_bgp_group

Provides a bgp group resource.

## Example Usage

```hcl
# Configure a bgp group
resource junos_bgp_group "groupbgpdemo" {
  name             = "GroupBgpDemo"
  routing_instance = "default"
  peer_as          = "65002"
  local_address    = "192.0.2.3"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Name of group.
* `routing_instance` - (Optional, Forces new resource)(`String`) Routing instance for bgp protocol. Need to be default or name of routing instance. Defaults to `default`.
* `type` - (Optional, Forces new resource)(`String`) Type of peer group. Need to be 'internal' or 'external'. Defaults to `external`.
* `accept_remote_nexthop` - (Optional)(`Bool`) Allow import policy to specify a non-directly connected next-hop.
* `advertise_external` - (Optional)(`Bool`) Advertise best external routes. Conflict with `advertise_external_conditional`.
* `advertise_external_conditional` - (Optional)(`Bool`) Route matches active route upto med-comparison rule. Conflict with `advertise_external`.
* `advertise_inactive` - (Optional)(`Bool`) Advertise inactive routes.
* `advertise_peer_as` - (Optional)(`Bool`) Advertise routes received from the same autonomous system. Conflict with `no_advertise_peer_as`.
* `no_advertise_peer_as` - (Optional)(`Bool`) Don't advertise routes received from the same autonomous system. Conflict with `advertise_peer_as`.
* `as_override` - (Optional)(`Bool`) Replace neighbor AS number with our AS number.
* `authentication_algorithm` - (Optional)(`String`) Authentication algorithm name. Conflict with `authentication_key`.
* `authentication_key` - (Optional)(`String`) MD5 authentication key. Conflict with `authentication_*`.
**WARNING** Clear in tfstate.
* `authentication_key_chain` - (Optional)(`String`) Key chain name. Conflict with `authentication_key`.
* `bfd_liveness_detection` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Define Bidirectional Forwarding Detection (BFD) options. See the [`bfd_liveness_detection` arguments](#bfd_liveness_detection-arguments) block. Max of 1.
* `damping` - (Optional)(`Bool`) Enable route flap damping.
* `export` - (Optional)(`ListOfString`) Export policy list.
* `family_inet` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each nlri_type.
See the [`family_inet` arguments](#family_inet-arguments) block.
* `family_inet6` Same options as [`family_inet` arguments](#family_inet-arguments) but for inet6 family.
* `graceful_restart` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Define BGP graceful restart options.See the [`graceful_restart` arguments](#graceful_restart-arguments) block. Max of 1.
* `hold_time` - (Optional)(`Int`) Hold time used when negotiating with a peer.
* `import` - (Optional)(`ListOfString`) Import policy list.
* `local_address` - (Optional)(`String`) Address of local end of BGP session.
* `local_as` - (Optional)(`String`) Local autonomous system number.
* `local_as_alias` - (Optional)(`Bool`) Treat this AS as an alias to the system AS. Conflict with other local_as options.
* `local_as_loops` - (Optional)(`Int`) Maximum number of times this AS can be in an AS path (1..10).
* `local_as_no_prepend_global_as` - (Optional)(`Bool`) Do not prepend global autonomous-system number in advertised paths. Conflict with other local_as options.
* `local_as_private` - (Optional)(`Bool`) Hide this local AS in paths learned from this peering. Conflict with other local_as options.
* `local_interface` - (Optional)(`String`) Local interface for IPv6 link local EBGP peering.
* `local_preference` - (Optional)(`Int`) Value of LOCAL_PREF path attribute.
* `log_updown` - (Optional)(`Bool`) Log a message for peer state transitions.
* `metric_out` - (Optional)(`Int`) Route metric sent in MED.
* `metric_out_igp` - (Optional)(`Bool`) Track the IGP metric. Conflict with `metric_out` and `metric_out_minimum_*`.
* `metric_out_igp_offset` - (Optional)(`Int`) Metric offset for MED. Conflict with `metric_out` and `metric_out_minimum_*`.
* `metric_out_igp_delay_med_update` - (Optional)(`Bool`) Delay updating MED when IGP metric increases. Conflict with `metric_out` and `metric_out_minimum_*`.
* `metric_out_minimum_igp` - (Optional)(`Bool`) Track the minimum IGP metric. Conflict with `metric_out` and `metric_out_(?!minimum)_*`.
* `metric_out_minimum_igp_offset` - (Optional)(`Bool`) Metric offset for MED. Conflict with `metric_out` and `metric_out_(?!minimum)_*`.
* `mtu_discovery` - (Optional)(`Bool`) Enable TCP path MTU discovery.
* `multihop` - (Optional)(`Bool`) Configure an EBGP multihop session.
* `multipath` - (Optional)(`Bool`) Allow load sharing among multiple BGP paths.
* `out_delay` - (Optional)(`Int`) How long before exporting routes from routing table.
* `passive` - (Optional)(`Bool`) Do not send open messages to a peer.
* `peer_as` - (Optional)(`String`) Autonomous system number.
* `preference` - (Optional)(`Int`) Preference value.
* `remove_private` - (Optional)(`Bool`) Remove well-known private AS numbers.

---
#### bfd_liveness_detection arguments
* `authentication_algorithm` - (Optional)(`String`) Authentication algorithm name.
* `authentication_key_chain` - (Optional)(`String`) Authentication key chain name.
* `authentication_loose_check`  - (Optional)(`Bool`) Verify authentication only if authentication is negotiated.
* `detection_time_threshold` - (Optional)(`Int`) High detection-time triggering a trap (milliseconds).
* `holddown_interval` - (Optional)(`Int`) Time to hold the session-UP notification to the client (1..255000 milliseconds).
* `minimum_interval` - (Optional)(`Int`) Minimum transmit and receive interval (1..255000 milliseconds).
* `minimum_receive_interval` - (Optional)(`Int`) Minimum receive interval (1..255000 milliseconds).
* `multiplier` - (Optional)(`Int`) Detection time multiplier (1..255).
* `session_mode` - (Optional)(`String`) BFD single-hop or multihop session-mode. Need to be 'automatic', 'multihop' or 'single-hop'.
* `transmit_interval_minimum_interval` - (Optional)(`Int`) Minimum transmit interval (1..255000 milliseconds).
* `transmit_interval_threshold` - (Optional)(`Int`) High transmit interval triggering a trap (milliseconds).
* `version` - (Optional)(`String`) BFD protocol version number.

---
#### family_inet arguments
Also for `family_inet6`

* `nlri_type` - (Required)(`String`) NLRI type. Need to be 'any', 'flow', 'labeled-unicast', 'unicast' or 'multicast'.
* `accepted_prefix_limit` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for define maximum number of prefixes accepted from a peer and options.
  * `maximum` - (Required)(`Int`) Maximum number of prefixes accepted from a peer (1..4294967295).
  * `teardown` - (Optional)(`Int`) Clear peer connection on reaching limit with this percentage of prefix-limit to start warnings.
  * `teardown_idle_timeout` - (Optional)(`Int`) Timeout before attempting to restart peer.
  * `teardown_idle_timeout_forever`  - (Optional)(`Bool`) Idle the peer until the user intervenes. Conflict with `teardown_idle_timeout`.
* `prefix_limit` Same options as [`accepted_prefix_limit`](#accepted_prefix_limit) but for limit maximum number of prefixes from a peer

---
#### graceful_restart arguments
* `disable` - (Optional)(`Bool`)Disable graceful restart.
* `restart_time` - (Optional)(`Int`) Restart time used when negotiating with a peer (1..600).
* `stale_route_time` - (Optional)(`Int`) Maximum time for which stale routes are kept (1..600).

## Import

Junos bgp group can be imported using an id made up of `<name>_-_<routing_instance>`, e.g.

```
$ terraform import junos_bgp_group.groupbgpdemo GroupBgpDemo_-_default
```
