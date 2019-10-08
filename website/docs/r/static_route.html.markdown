---
layout: "junos"
page_title: "Junos: junos_security_static_route"
sidebar_current: "docs-junos-resource-security-static-route"
description: |-
  Create a static route for destination
---

# junos_security_static_route

Provides a static route resource for destination.

## Example Usage

```hcl
# Add a static route
resource "junos_security_static_route" "DemoStaticRoute" {
  destination      = "192.0.2.0/25"
  routing_instance = "prod-vr"
  next_hop         = ["st0.0"]
}
```

## Argument Reference

The following arguments are supported:

* `destination` - (Required, Forces new resource)(`String`) The destination for static route.
* `routing_instance` - (Optional, Forces new resource)(`String`) Routing instance for route. Need to be default or name of routing instance. Default to `default`
* `preference` - (Optional)(`Int`) Preference for static route
* `metric` - (Optional)(`Int`) Metric for static route
* `next_hop` - (Optional)(`ListOfString`) List of next-hop
* `qualified_next_hop` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) List of qualified-next-hop with options. Can be specified multiple times for each qualified-next-hop.
  * `next_hop` - (Required)(`String`) Target for qualified-next-hop
  * `preference` - (Optional)(`Int`) Preference of qualified next hop
  * `metric` - (Optional)(`Int`) Metric of qualified next hop

## Import

Junos static route can be imported using an id made up of `<destination>_-_<routing_instance>`, e.g.

```
$ terraform import junos_security_static_route.DemoStaticRoute 192.0.2.0/25_-_prod-vr
```
