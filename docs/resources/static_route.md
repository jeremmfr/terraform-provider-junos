---
page_title: "Junos: junos_static_route"
---

# junos_static_route

Provides a static route resource for destination.

## Example Usage

```hcl
# Add a static route
resource "junos_static_route" "demo_static_route" {
  destination      = "192.0.2.0/25"
  routing_instance = "prod-vr"
  next_hop         = ["st0.0"]
}
```

## Argument Reference

The following arguments are supported:

- **destination** (Required, String, Forces new resource)  
  Destination prefix.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for static route.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`.
- **active** (Optional, Boolean)  
  Remove inactive route from forwarding table.  
  Conflict with `passive`.
- **as_path_aggregator_address** (Optional, String)  
  Address of BGP system to add AGGREGATOR path attribute to route.
- **as_path_aggregator_as_number** (Optional, String)  
  AS number to add AGGREGATOR path attribute to route.
- **as_path_atomic_aggregate** (Optional, Boolean)  
  Add ATOMIC_AGGREGATE path attribute to route.
- **as_path_origin** (Optional, String)  
  Define origin.  
  Need to be `egp`, `igp` or `incomplete`.
- **as_path_path** (Optional, String)  
  Path to as-path.
- **community** (Optional, List of String)  
  BGP community.
- **discard** (Optional, Boolean)  
  Drop packets to destination; send no ICMP unreachables.  
  Conflict with `next_hop`, `next_table`, `qualified_next_hop`, `receive` and `reject`.
- **install** (Optional, Boolean)  
  Install route into forwarding table.  
  Conflict with `no_install`.
- **no_install** (Optional, Boolean)  
  Don't install route into forwarding table.  
  Conflict with `install`.
- **metric** (Optional, Number)  
  Metric for static route.
- **next_hop** (Optional, List of String)  
  Next-hop to destination.  
  Conflict with `discard`, `next_table`, `receive` and `reject`.
- **next_table** (Optional, String)  
  Next hop to another table.  
  Conflict with `discard`, `next_hop`, `qualified_next_hop`, `receive` and `reject`.
- **passive** (Optional, Boolean)  
  Retain inactive route in forwarding table.  
  Conflict with `active`.
- **preference** (Optional, Number)  
  Preference for static route.
- **qualified_next_hop** (Optional, Block List)  
  For each `next_hop` with qualifiers.  
  Conflict with `discard`, `next_table`, `receive` and `reject`.
  - **next_hop** (Required, String)  
    Next-hop with qualifiers to destination.
  - **interface** (Optional, String)  
    Interface of qualified next hop (Cannot be used with interface set as next-hop).
  - **metric** (Optional, Number)  
    Metric of qualified next hop.
  - **preference** (Optional, Number)  
    Preference of qualified next hop.
- **readvertise** (Optional, Boolean)  
  Mark route as eligible to be readvertised.  
  Conflict with `no_readvertise`.
- **no_readvertise** (Optional, Boolean)  
  Don't mark route as eligible to be readvertised.  
  Conflict with `readvertise`.
- **receive** (Optional, Boolean)  
  Install a receive route for the destination.  
  Conflict with `discard`, `next_hop`, `next_table`, `qualified_next_hop` and `reject`.
- **reject** (Optional, Boolean)  
  Drop packets to destination; send ICMP unreachables.  
  Conflict with `discard`, `next_hop`, `next_table`, `qualified_next_hop` and `receive`.
- **resolve** (Optional, Boolean)  
  Allow resolution of indirectly connected next hops.  
  Conflict with `no_resolve`, `retain` and `no_retain`.
- **no_resolve** (Optional, Boolean)  
  Don't allow resolution of indirectly connected next hops.  
  Conflict with `resolve`.
- **retain** (Optional, Boolean)  
  Always keep route in forwarding table.  
  Conflict with `resolve` and `no_retain`.
- **no_retain** (Optional, Boolean)  
  Don't always keep route in forwarding table.  
  Conflict with `resolve` and `retain`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<destination>_-_<routing_instance>`.

## Import

Junos static route can be imported using an id made up of `<destination>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_static_route.demo_static_route 192.0.2.0/25_-_prod-vr
```
