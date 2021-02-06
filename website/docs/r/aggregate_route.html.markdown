---
layout: "junos"
page_title: "Junos: junos_aggregate_route"
sidebar_current: "docs-junos-resource-aggregate-route"
description: |-
  Create a aggregate route for destination
---

# junos_aggregate_route

Provides a aggregate route resource for destination.

## Example Usage

```hcl
# Add a aggregate route
resource junos_aggregate_route "demo_aggregate_route" {
  destination      = "192.0.2.0/25"
  routing_instance = "prod-vr"
  brief            = true
}
```

## Argument Reference

The following arguments are supported:

* `destination` - (Required, Forces new resource)(`String`) The destination for aggregate route.
* `routing_instance` - (Optional, Forces new resource)(`String`) Routing instance for route. Need to be default or name of routing instance. Defaults to `default`
* `active` - (Optional)(`Bool`) Remove inactive route from forwarding table.
* `brief` - (Optional)(`Bool`) Include longest common sequences from contributing paths.
* `community` - (Optional)(`ListOfString`) List of BGP community.
* `discard` - (Optional)(`Bool`) Drop packets to destination; send no ICMP unreachables.
* `full` - (Optional)(`Bool`) Include all AS numbers from all contributing paths.
* `metric` - (Optional)(`Int`) Metric for aggregate route.
* `passive` - (Optional)(`Bool`) Retain inactive route in forwarding table.
* `policy` - (Optional)(`ListOfString`) List of Policy filter.
* `preference` - (Optional)(`Int`) Preference for aggregate route.

## Import

Junos aggregate route can be imported using an id made up of `<destination>_-_<routing_instance>`, e.g.

```
$ terraform import junos_aggregate_route.demo_aggregate_route 192.0.2.0/25_-_prod-vr
```
