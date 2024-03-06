---
page_title: "Junos: junos_interface_st0_unit"
---

# junos_interface_st0_unit

Find an available st0 logical interface and create it.

It's useful for bind_interface in `junos_security_ipsec_vpn` resource.  
New st0 unit interface can be configured with `junos_interface_logical` resource.

## Example Usage

```hcl
resource "junos_interface_st0_unit" "demo" {}
```

## Attribute Reference

- **id** (String)  
  Name of interface found and created.

## Import

Junos st0 unit interface can be imported using an id made up of the name of interface, e.g.

```shell
$ terraform import junos_interface_st0_unit.demo st0.0
```
