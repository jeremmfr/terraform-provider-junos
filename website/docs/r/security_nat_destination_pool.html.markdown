---
layout: "junos"
page_title: "Junos: junos_security_nat_destination_pool"
sidebar_current: "docs-junos-resource-security-nat-destination-pool"
description: |-
  Create a security nat destination pool (when Junos device supports it)
---

# junos_security_nat_destination_pool

Provides a security pool resource for destination nat.

## Example Usage

```hcl
# Add a destination nat pool
resource junos_security_nat_destination_pool "demo_dnat_pool" {
  name    = "ip_internal"
  address = "192.0.2.2/32"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of destination nat pool.
* `address` - (Required)(`String`) IP/mask for destination nat pool.
* `address_port` - (Optional)(`Int`) Port change too with destination nat. Conflict with `address_to`.
* `address_to` - (Optional)(`String`) IP/mask for range of destination nat pool (range = `address` to `address_to`).
* `routing_instance` - (Optional)(`String`) Name of routing instance for switch instance with nat.

## Import

Junos security nat destination pool can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_nat_destination_pool.demo_dnat_pool ip_internal
```
