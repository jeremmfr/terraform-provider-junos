---
page_title: "Junos: junos_security_nat_destination"
---

# junos_security_nat_destination

Provides a security destination nat resource.

## Example Usage

```hcl
# Add a destination nat
resource "junos_security_nat_destination" "demo_dnat" {
  name = "dnat_from_untrust"
  from {
    type  = "zone"
    value = ["untrust"]
  }
  rule {
    name                = "nat_192_0_2_129"
    destination_address = "192.0.2.129/32"
    then {
      type = "pool"
      pool = "pool_trust"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Destination nat rule-set name.
- **from** (Required, Block)  
  Declare where is the traffic from.
  - **type** (Required, String)  
    Type of traffic source.  
    Need to be `interface`, `routing-instance` or `zone`
  - **value** (Required, Set of String)  
    Name of interface, routing-instance or zone for traffic source.
- **rule** (Required, Block List)  
  For each name of destination nat rule to declare.  
  See [below for nested schema](#rule-arguments).
- **description** (Optional, String)  
  Text description of destination nat rule-set.

---

### rule arguments

-> **Note:** One of `destination_address` or `destination_address_name` arguments is required.

- **name** (Required, String)  
  Rule name.
- **destination_address** (Optional, String)  
  CIDR destination address to match.
- **destination_address_name** (Optional, String)  
  Destination address from address book to match.
- **application** (Optional, Set of String)  
  Application or application-set name to match.
- **destination_port** (Optional, Set of String)  
  Destination port to match.  
  Format need to be `x` or `x to y`.
- **protocol** (Optional, Set of String)  
  IP Protocol to match.
- **source_address** (Optional, Set of String)  
  CIDR source address to match.
- **source_address_name** (Optional, Set of String)  
  Source address from address book to match.
- **then** (Required, Block)  
  Declare `then` action.
  - **type** (Required, String)  
    Type of destination nat.  
    Need to be `pool` or `off`
  - **pool** (Optional, String)  
    Name of destination nat pool when type is pool.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security nat destination can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_nat_destination.demo_dnat dnat_from_untrust
```
