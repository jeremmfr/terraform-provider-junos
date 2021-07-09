---
layout: "junos"
page_title: "Junos: junos_generate_route"
sidebar_current: "docs-junos-resource-generate-route"
description: |-
  Create a generate route for destination
---

# junos_generate_route

Provides a generate route resource for destination.

## Example Usage

```hcl
# Add a generate route
resource junos_generate_route "demo_generate_route" {
  destination      = "192.0.2.0/25"
  routing_instance = "prod-vr"
  brief            = true
}
```

## Argument Reference

The following arguments are supported:

* `destination` - (Required, Forces new resource)(`String`) The destination for generate route.
* `routing_instance` - (Optional, Forces new resource)(`String`) Routing instance for route. Need to be default or name of routing instance. Defaults to `default`
* `active` - (Optional)(`Bool`) Remove inactive route from forwarding table.
* `as_path_aggregator_address` - (Optional)(`String`) Address of BGP system to add AGGREGATOR path attribute to route.
* `as_path_aggregator_as_number` - (Optional)(`String`) AS number to add AGGREGATOR path attribute to route.
* `as_path_atomic_aggregate` - (Optional)(`Bool`) Add ATOMIC_AGGREGATE path attribute to route.
* `as_path_origin` - (Optional)(`String`) Define origin. Need to be 'egp', 'igp' or 'incomplete'.
* `as_path_path` - (Optional)(`String`) Path to as-path.
* `brief` - (Optional)(`Bool`) Include longest common sequences from contributing paths.
* `community` - (Optional)(`ListOfString`) List of BGP community.
* `discard` - (Optional)(`Bool`) Drop packets to destination; send no ICMP unreachables.
* `full` - (Optional)(`Bool`) Include all AS numbers from all contributing paths.
* `metric` - (Optional)(`Int`) Metric for generate route.
* `next_table` - (Optional)(`String`)  Next hop to another table. Conflict with `discard`.
* `passive` - (Optional)(`Bool`) Retain inactive route in forwarding table.
* `policy` - (Optional)(`ListOfString`) List of Policy filter.
* `preference` - (Optional)(`Int`) Preference for generate route.

## Import

Junos generate route can be imported using an id made up of `<destination>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_generate_route.demo_generate_route 192.0.2.0/25_-_prod-vr
```
