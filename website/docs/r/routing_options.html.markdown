---
layout: "junos"
page_title: "Junos: junos_routing_options"
sidebar_current: "docs-junos-resource-routing-options"
description: |-
  Configure static configuration in routing-options block
---

# junos_routing_options

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `routing-options` block. Destroy this resource has no effect on the Junos configuration.

Configure static configuration in `routing-options` block

## Example Usage

```hcl
# Configure routing-options
resource junos_routing_options "routing_options" {
  autonomous_system {
    number = "65000"
  }
  graceful_restart {}
}
```

## Argument Reference

The following arguments are supported:

* `autonomous_system` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'autonomous-system' configuration.
  * `number` - (Required)(`String`) Autonomous system number in plain number or 'higher 16bits'.'Lower 16 bits' (asdot notation) format.
  * `asdot_notation` - (Optional)(`Bool`) Use AS-Dot notation to display true 4 byte AS numbers.
  * `loops` - (Optional)(`Int`) Maximum number of times this AS can be in an AS path (1..10).
* `forwarding_table` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'forwarding-table' configuration.
  * `chain_composite_max_label_count` - (Optional)(`Int`) Maximum labels inside chain composite for the platform. (1..8).
  * `chained_composite_next_hop_ingress` - (Optional)(`ListOfString`) Next-hop chaining mode -> Ingress LSP nexthop settings.
  * `chained_composite_next_hop_transit` - (Optional)(`ListOfString`) Next-hop chaining mode -> Transit LSP nexthops settings.
  * `dynamic_list_next_hop` - (Optional)(`Bool`) Dynamic next-hop mode for EVPN.
  * `ecmp_fast_reroute` - (Optional)(`Bool`) Enable fast reroute for ECMP next hops.
  * `no_ecmp_fast_reroute` - (Optional)(`Bool`) Don't enable fast reroute for ECMP next hops.
  * `export` - (Optional)(`ListOfString`) Export policy.
  * `indirect_next_hop` - (Optional)(`Bool`) Install indirect next hops in Packet Forwarding Engine.
  * `no_indirect_next_hop` - (Optional)(`Bool`) Don't install indirect next hops in Packet Forwarding Engine.
  * `indirect_next_hop_change_acknowledgements` - (Optional)(`Bool`) Request acknowledgements for Indirect next hop changes.
  * `no_indirect_next_hop_change_acknowledgements` - (Optional)(`Bool`) Don't request acknowledgements for Indirect next hop changes.
  * `krt_nexthop_ack_timeout` - (Optional)(`Int`) Kernel nexthop ack timeout interval (1..400).
  * `remnant_holdtime` - (Optional)(`Int`) Time to hold inherited routes from FIB (0..10000).
  * `unicast_reverse_path` - (Optional)(`String`)  Unicast reverse path (RP) verification. Need to be 'active-paths' or 'feasible-paths'.
* `graceful_restart` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'graceful-restart' configuration.
  * `disable` - (Optional)(`Bool`) Disable graceful restart.
  * `restart_duration` - (Optional)(`Int`) Maximum time for which router is in graceful restart (120..10000).

## Import

Junos routing_options can be imported using any id, e.g.

```
$ terraform import junos_routing_options.routing_options random
```
