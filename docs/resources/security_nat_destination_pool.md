---
page_title: "Junos: junos_security_nat_destination_pool"
---

# junos_security_nat_destination_pool

Provides a security pool resource for destination nat.

## Example Usage

```hcl
# Add a destination nat pool
resource "junos_security_nat_destination_pool" "demo_dnat_pool" {
  name    = "ip_internal"
  address = "192.0.2.2/32"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Pool name.
- **address** (Required, String)  
  CIDR address to destination nat pool.
- **address_port** (Optional, Number)  
  Port change too with destination nat.  
  Conflict with `address_to`.
- **address_to** (Optional, String)  
  CIDR to define range of address to destination nat pool (range = `address` to `address_to`).
- **description** (Optional, String)  
  Text description of pool.
- **routing_instance** (Optional, String)  
  Name of routing instance to switch instance with nat.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security nat destination pool can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_nat_destination_pool.demo_dnat_pool ip_internal
```
