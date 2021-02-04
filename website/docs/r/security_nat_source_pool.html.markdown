---
layout: "junos"
page_title: "Junos: junos_security_nat_source_pool"
sidebar_current: "docs-junos-resource-security-nat-source-pool"
description: |-
  Create a security nat source pool (when Junos device supports it)
---

# junos_security_nat_source_pool

Provides a security pool resource for source nat.

## Example Usage

```hcl
# Add a source nat pool
resource junos_security_nat_source_pool "demo_snat_pool" {
  name    = "ip_external"
  address = ["192.0.2.129/32"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of source nat pool.
* `address` - (Required)(`ListofString`) List of IP/mask for source nat pool.
* `port_no_translation` - (Optional)(`Bool`) Do not perform port translation.
* `port_overloading_factor` - (Optional)(`Int`) Port overloading factor for each IP.
* `port_range` - (Optional)(`String`) Range of port for source nat.
* `routing_instance` - (Optional)(`String`) Name of routing instance for switch with nat.

## Import

Junos security nat source pool can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_nat_source_pool.demo_snat_pool ip_external
```
