---
page_title: "Junos: junos_routing_options"
---

# junos_routing_options

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `routing-options` block.  
By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.

Configure static configuration in `routing-options` block

## Example Usage

```hcl
# Configure routing-options
resource "junos_routing_options" "routing_options" {
  autonomous_system {
    number = "65000"
  }
  graceful_restart {}
}
```

## Argument Reference

The following arguments are supported:

- **clean_on_destroy** (Optional, Boolean)  
  Clean supported lines when destroy this resource.
- **autonomous_system** (Optional, Block)  
  Declare `autonomous-system` configuration.
  - **number** (Required, String)  
    Autonomous system number in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format.
  - **asdot_notation** (Optional, Boolean)  
    Use AS-Dot notation to display true 4 byte AS numbers.
  - **loops** (Optional, Number)  
    Maximum number of times this AS can be in an AS path (1..10).
- **forwarding_table** (Optional, Block)  
  Declare `forwarding-table` configuration.
  - **chain_composite_max_label_count** (Optional, Number)  
    Maximum labels inside chain composite for the platform. (1..8).
  - **chained_composite_next_hop_ingress** (Optional, Set of String)  
    Next-hop chaining mode -> Ingress LSP nexthop settings.
  - **chained_composite_next_hop_transit** (Optional, Set of String)  
    Next-hop chaining mode -> Transit LSP nexthops settings.
  - **dynamic_list_next_hop** (Optional, Boolean)  
    Dynamic next-hop mode for EVPN.
  - **ecmp_fast_reroute** (Optional, Boolean)  
    Enable fast reroute for ECMP next hops.
  - **no_ecmp_fast_reroute** (Optional, Boolean)  
    Don't enable fast reroute for ECMP next hops.
  - **export** (Optional, List of String)  
    Export policy.
  - **indirect_next_hop** (Optional, Boolean)  
    Install indirect next hops in Packet Forwarding Engine.
  - **no_indirect_next_hop** (Optional, Boolean)  
    Don't install indirect next hops in Packet Forwarding Engine.
  - **indirect_next_hop_change_acknowledgements** (Optional, Boolean)  
    Request acknowledgements for Indirect next hop changes.
  - **no_indirect_next_hop_change_acknowledgements** (Optional, Boolean)  
    Don't request acknowledgements for Indirect next hop changes.
  - **krt_nexthop_ack_timeout** (Optional, Number)  
    Kernel nexthop ack timeout interval (1..400).
  - **remnant_holdtime** (Optional, Number)  
    Time to hold inherited routes from FIB (0..10000).
  - **unicast_reverse_path** (Optional, String)  
    Unicast reverse path (RP) verification.  
    Need to be `active-paths` or `feasible-paths`.
- **forwarding_table_export_configure_singly** (Optional, Boolean)  
  Disable management of `forwarding-table export` in this resource to be able to manage them directly
  from `junos_policyoptions_policy_statement` resources with `add_it_to_forwarding_table_export`
  argument.  
  Conflict with `forwarding_table.0.export`.
- **graceful_restart** (Optional, Block)  
  Declare `graceful-restart` configuration.
  - **disable** (Optional, Boolean)  
    Disable graceful restart.
  - **restart_duration** (Optional, Number)  
    Maximum time for which router is in graceful restart (120..10000).
- **instance_export** (Optional, List of String)  
  Export policy for instance RIBs
- **instance_import** (Optional, List of String)  
  Import policy for instance RIBs
- **router_id** (Optional, String)  
  Router identifier.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `routing_options`.

## Import

Junos routing_options can be imported using any id, e.g.

```shell
$ terraform import junos_routing_options.routing_options random
```
