---
layout: "junos"
page_title: "Junos: junos_security_routing_instance"
sidebar_current: "docs-junos-resource-security-routing-instance"
description: |-
  Create a routing instance
---

# junos_security_routing_instance

Provides a routing instance resource.

## Example Usage

```hcl
# Add a routing instance
resource "junos_security_routing_instance" "DemoRI" {
  name = "prod-vr"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of routing instance.
* `type` - (Optional)(`String`) Type of routing instance. Default to `virtual-router`
* `as` - (Optional)(`String`) Autonomous system number in plain number or 'higher 16bits'.'Lower 16 bits' (asdot notation) format.

## Import

Junos routing instance can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_routing_instance.DemoRI prod-vr
```
