---
layout: "junos"
page_title: "Junos: junos_static_route"
sidebar_current: "docs-junos-resource-static-route"
description: |-
  Create a static route for destination
---

# junos_static_route

Provides a static route resource for destination.

## Example Usage

```hcl
# Add a static route
resource junos_static_route "demo_static_route" {
  destination      = "192.0.2.0/25"
  routing_instance = "prod-vr"
  next_hop         = ["st0.0"]
}
```

## Argument Reference

The following arguments are supported:

* `destination` - (Required, Forces new resource)(`String`) The destination for static route.
* `routing_instance` - (Optional, Forces new resource)(`String`) Routing instance for route. Need to be default or name of routing instance. Defaults to `default`.
* `active` - (Optional)(`Bool`) Remove inactive route from forwarding table. Conflict with `passive`.
* `community` - (Optional)(`ListOfString`) List of BGP community.
* `discard` - (Optional)(`Bool`) Drop packets to destination; send no ICMP unreachables. Conflict with `next_hop`, `next_table`, `qualified_next_hop`, `receive` and `reject`.
* `install` - (Optional)(`Bool`) Install route into forwarding table. Conflict with `no_install`.
* `no_install` - (Optional)(`Bool`) Don't install route into forwarding table. Conflict with `install`.
* `metric` - (Optional)(`Int`) Metric for static route.
* `next_hop` - (Optional)(`ListOfString`) List of next-hop. Conflict with `discard`, `next_table`, `receive` and `reject`.
* `next_table` - (Optional)(`String`) Next hop to another table. Conflict with `discard`, `next_hop`, `qualified_next_hop`, `receive` and `reject`.
* `passive` - (Optional)(`Bool`) Retain inactive route in forwarding table. Conflict with `active`.
* `preference` - (Optional)(`Int`) Preference for static route.
* `qualified_next_hop` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) List of qualified-next-hop with options. Can be specified multiple times for each qualified-next-hop. Conflict with `discard`, `next_table`, `receive` and `reject`.
  * `next_hop` - (Required)(`String`) Target for qualified-next-hop.
  * `interface` - (Optional)(`String`) Interface of qualified next hop (Cannot be used with interface set as next-hop).
  * `metric` - (Optional)(`Int`) Metric of qualified next hop.
  * `preference` - (Optional)(`Int`) Preference of qualified next hop.
* `readvertise` - (Optional)(`Bool`) Mark route as eligible to be readvertised. Conflict with `no_readvertise`.
* `no_readvertise` - (Optional)(`Bool`) Don't mark route as eligible to be readvertised. Conflict with `readvertise`.
* `receive` - (Optional)(`Bool`) Install a receive route for the destination. Conflict with `discard`, `next_hop`, `next_table`, `qualified_next_hop` and `reject`.
* `reject` - (Optional)(`Bool`) Drop packets to destination; send ICMP unreachables. Conflict with `discard`, `next_hop`, `next_table`, `qualified_next_hop` and `receive`.
* `resolve` - (Optional)(`Bool`) Allow resolution of indirectly connected next hops. Conflict with `no_resolve`, `retain` and `no_retain`.
* `no_resolve` - (Optional)(`Bool`) Don't allow resolution of indirectly connected next hops. Conflict with `resolve`.
* `retain` - (Optional)(`Bool`) Always keep route in forwarding table. Conflict with `resolve` and `no_retain`.
* `no_retain` - (Optional)(`Bool`) Don't always keep route in forwarding table. Conflict with `resolve` and `retain`.

## Import

Junos static route can be imported using an id made up of `<destination>_-_<routing_instance>`, e.g.

```
$ terraform import junos_static_route.demo_static_route 192.0.2.0/25_-_prod-vr
```
