---
page_title: "Junos: junos_security_nat_static"
---

# junos_security_nat_static

Provides a security static nat resource.

## Example Usage

```hcl
# Add a static nat
resource "junos_security_nat_static" "demo_nat" {
  name = "nat_from_trust"
  from {
    type  = "zone"
    value = ["trust"]
  }
  rule {
    name                = "nat_192_0_2_0_25"
    destination_address = "192.0.2.0/25"
    then {
      type   = "prefix"
      prefix = "192.0.2.128/25"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

-> **Note:** One of `rule` or `configure_rules_singly` arguments is required.

- **name** (Required, String, Forces new resource)  
  Static nat rule-set name.
- **from** (Required, Block)  
  Declare where is the traffic from.
  - **type** (Required, String)  
    Type of traffice source.  
    Need to be `interface`, `routing-instance` or `zone`.
  - **value** (Required, Set of String)  
    Name of interface, routing-instance or zone for traffic source.
- **rule** (Optional, Block List)  
  For each name of static nat rule to declare.  
  See [below for nested schema](#rule-arguments).
- **configure_rules_singly** (Optional, Boolean)
  Disable management of rules in this resource to be able to manage them with specific
  resources.
- **description** (Optional, String)  
  Text description of static nat rule-set.

---

### rule arguments

-> **Note:** One of `destination_address` or `destination_address_name` arguments is required.

- **name** (Required, String)  
  Rule name.
- **destination_address** (Optional, String)  
  CIDR destination address to match.
- **destination_address_name** (Optional, String)  
  Destination address from address book to match.
- **destination_port** (Optional, Number)  
  Destination port or lower limit of port range to match.
- **destination_port_to** (Optional, Number)  
  Port range upper limit to match.
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
    Type of static nat.  
    Need to be `inet`, `prefix` or `prefix-name`.
  - **mapped_port** (Optional, Number)  
    Port or lower limit of port range to mapped port.  
    `type` need to be `prefix` or `prefix-name`.
  - **mapped_port_to** (Optional, Number)  
    Port range upper limit to mapped port.  
    `type` need to be `prefix` or `prefix-name`.
  - **prefix** (Optional, String)  
    CIDR or address from address book to prefix static nat.  
    `type` need to be `prefix` or `prefix-name`.  
    CIDR is required if `type` = `prefix`.
  - **routing_instance** (Optional, String)  
    Name of routing instance to switch instance with nat.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security nat static can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_nat_static.demo_nat nat_from_trust
```

By default, all rules are imported. To import only rule-set with `configure_rules_singly` = true and
without `rule` blocks, add suffix `_-_no_rules` at `<name>`, e.g.

```shell
$ terraform import junos_security_nat_static.demo_nat nat_from_trust_-_no_rules
```
