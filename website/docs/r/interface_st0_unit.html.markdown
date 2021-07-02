---
layout: "junos"
page_title: "Junos: junos_interface_st0_unit"
sidebar_current: "docs-junos-resource-interface-st0-unit"
description: |-
  Find st0 unit available and create interface.
---

# junos_interface_st0_unit

Find st0 unit available and create interface.

It's useful for bind_interface in `junos_security_ipsec_vpn` resource.  
New st0 unit interface can be configured with `junos_interface_logical` resource.

## Example Usage

```hcl
resource junos_interface_st0_unit "demo" {}
```

## Attributes Reference

* `id` - Name of interface found and created.

## Import

Junos st0 unit interface can be imported using an id made up of the name of interface, e.g.

```shell
$ terraform import junos_interface_st0_unit.demo st0.0
```
