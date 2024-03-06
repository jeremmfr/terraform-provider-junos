---
page_title: "Junos: junos_aggregate_route"
---

# junos_aggregate_route

Provides an aggregate route resource for destination.

## Example Usage

```hcl
# Add an aggregate route
resource "junos_aggregate_route" "demo_aggregate_route" {
  destination      = "192.0.2.0/25"
  routing_instance = "prod-vr"
  brief            = true
}
```

## Argument Reference

The following arguments are supported:

- **destination** (Required, String, Forces new resource)  
  Destination prefix.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for aggregate route.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`
- **active** (Optional, Boolean)  
  Remove inactive route from forwarding table.
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
- **brief** (Optional, Boolean)  
  Include longest common sequences from contributing paths.
- **community** (Optional, List of String)  
  BGP community.
- **discard** (Optional, Boolean)  
  Drop packets to destination; send no ICMP unreachables.
- **full** (Optional, Boolean)  
  Include all AS numbers from all contributing paths.
- **metric** (Optional, Number)  
  Metric for aggregate route.
- **passive** (Optional, Boolean)  
  Retain inactive route in forwarding table.
- **policy** (Optional, List of String)  
  Policy filter.
- **preference** (Optional, Number)  
  Preference for aggregate route.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<destination>_-_<routing_instance>`.

## Import

Junos aggregate route can be imported using an id made up of `<destination>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_aggregate_route.demo_aggregate_route 192.0.2.0/25_-_prod-vr
```
