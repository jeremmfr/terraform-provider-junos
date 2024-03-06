---
page_title: "Junos: junos_security_nat_source"
---

# junos_security_nat_source

Provides a security source nat resource.

## Example Usage

```hcl
# Add a source nat
resource "junos_security_nat_source" "demo_snat" {
  name = "nat_from_trust_to_untrust"
  from {
    type  = "zone"
    value = ["trust"]
  }
  to {
    type  = "zone"
    value = ["untrust"]
  }
  rule {
    name = "nat_192_0_2_0_25"
    match {
      source_address = ["192.0.2.0/25"]
    }
    then {
      type = "pool"
      pool = "pool_untrust"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Source nat rule-set name.
- **from** (Required, Block)  
  Declare where is the traffic from.
  - **type** (Required, String)  
    Type of traffic source.  
    Need to be `interface`, `routing-instance` or `zone`.
  - **value** (Required, Set of String)  
    Name of interface, routing-instance or zone for traffic source.
- **to** (Required, Block)  
  Declare where is the traffic to.
  - **type** (Required, String)  
    Type of traffic destination.
    Need to be `interface`, `routing-instance` or `zone`.
  - **value** (Required, Set of String)  
    Name of interface, routing-instance or zone for traffic destination.
- **rule** (Required, Block List)  
  For each name of source nat rule to declare.  
  See [below for nested schema](#rule-arguments).
- **description** (Optional, String)  
  Text description of rule set.

---

### rule arguments

- **name** (Required, String)  
  Rule name.
- **match** (Required, Block)  
  Specify source nat rule match criteria.
  - **application** (Optional, Set of String)  
    Application or application-set name to match.
  - **destination_address** (Optional, Set of String)  
    CIDR destination address to match.
  - **destination_address_name** (Optional, Set of String)  
    Destination address from address book to match.
  - **destination_port** (Optional, Set of String)  
    Destination port to match.  
    Format need to be `x` or `x to y`.
  - **protocol** (Optional, Set of String)  
    IP Protocol to match.
  - **source_address** (Optional, Set of String)  
    CIDR source address to match.
  - **source_address_name** (Optional, Set of String)  
    Source address from address book to match.
  - **source_port** (Optional, Set of String)  
    Source port to match.  
    Format need to be `x` or `x to y`.
- **then** (Required, Block)  
  Declare `then` configuration.
  - **type** (Required, String)  
    Type of source nat.  
    Need to be `interface`, `pool` or `off`.
  - **pool** (Optional, String)  
    Name of nat source pool when type pool.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security nat source can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_nat_source.demo_snat nat_from_trust_to_untrust
```
