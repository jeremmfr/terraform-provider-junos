---
layout: "junos"
page_title: "Junos: junos_security_nat_destination"
sidebar_current: "docs-junos-resource-security-nat-destination"
description: |-
  Create a security nat destination (when Junos device supports it)
---

# junos_security_nat_destination

Provides a security destination nat resource.

## Example Usage

```hcl
# Add a destination nat
resource junos_security_nat_destination "demo_dnat" {
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
  The name of destination nat.
- **from** (Required, Block)  
  Declare `from` configuration.
  - **type** (Required, String)  
    Type of from options.  
    Need to be `interface`, `routing-instance` or `zone`
  - **value** (Required, Set of String)  
    Name of interface, routing-instance or zone for from options
- **rule** (Required, Block List)  
  For each name of rule to declare.  
  See [below for nested schema](#rule-arguments).
- **description** (Optional, String)  
  Text description of rule set

---

### rule arguments

-> **Note:** One of `destination_address` or `destination_address_name` arguments is required.

- **name** (Required, String)  
  Name of rule
- **destination_address** (Optional, String)  
  CIDR for match destination address
- **destination_address_name** (Optional, String)  
  Destination address from address book for rule match.
- **application** (Optional, Set of String)  
  Specify application or application-set name for rule match.
- **destination_port** (Optional, Set of String)  
  List of destination port for rule match.  
  Format need to be `x` or `x to y`.
- **protocol** (Optional, Set of String)  
  IP Protocol for rule match.
- **source_address** (Optional, Set of String)  
  List of CIDR source address for rule match.
- **source_address_name** (Optional, Set of String)  
  List of source address from address book for rule match.
- **then** (Required, Block)  
  Declare `then` action.
  - **type** (Required, String)  
    Type of destination nat.  
    Need to be `pool` or `off`
  - **pool** (Optional, String)  
    Name of nat destination pool when type pool

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security nat destination can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_nat_destination.demo_dnat dnat_from_untrust
```
