---
page_title: "Junos: junos_security_nat_static_rule"
---

# junos_security_nat_static_rule

Provides a security nat static rule resource in rule-set.

-> **Note:** Security nat static rule-set can be created with `junos_security_nat_static` resource.
This resource needs to have `configure_rules_singly` set to true otherwise there will be a conflict
between resources.

## Example Usage

```hcl
# Add a static nat rule
resource "junos_security_nat_static_rule" "demo_nat_rule" {
  name                = "nat_192_0_2_0_25"
  rule_set            = "nat_from_trust"
  destination_address = "192.0.2.0/25"
  then {
    type   = "prefix"
    prefix = "192.0.2.128/25"
  }
}
```

## Argument Reference

The following arguments are supported:

-> **Note:** One of `destination_address` or `destination_address_name` arguments is required.

- **name** (Required, String, Forces new resource)  
  Static Rule name.
- **rule_set** (Required, String, Forces new resource)  
  Static nat rule-set name.
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
  An identifier for the resource with format `<rule_set>_-_<name>`.

## Import

Junos security nat static rule can be imported using an id made up of `<rule_set>_-_<name>`, e.g.

```shell
$ terraform import junos_security_nat_static_rule.demo_nat_rule nat_from_trust_-_nat_192_0_2_0_25
```
