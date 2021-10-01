---
layout: "junos"
page_title: "Junos: junos_security_nat_static"
sidebar_current: "docs-junos-resource-security-nat-static"
description: |-
  Create a security nat static (when Junos device supports it)
---

# junos_security_nat_static

Provides a security static nat resource.

## Example Usage

```hcl
# Add a static nat
resource junos_security_nat_static "demo_nat" {
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

- **name** (Required, String, Forces new resource)  
  The name of static nat.
- **from** (Required, Block)  
  Declare `from` configuration.
  - **type** (Required, String)  
    Type of from options.  
    Need to be `interface`, `routing-instance` or `zone`.
  - **value** (Required, Set of String)  
    Name of interface, routing-instance or zone for from options.
- **rule** (Required, Block List)  
  For each name of rule to declare.  
  See [below for nested schema](#rule-arguments).

---

### rule arguments

-> **Note:** One of `destination_address` or `destination_address_name` arguments is required.

- **name** (Required, String)  
  Name of rule.
- **destination_address** (Optional, String)  
  CIDR of destination address for rule match.
- **destination_address_name** (Optional, String)  
  Destination address from address book for rule match.
- **destination_port** (Optional, Number)  
  Destination port or lower limit of port range for rule match.
- **destination_port_to** (Optional, Number)  
  Port range upper limit for rule match.
- **source_address** (Optional, Set of String)  
  List of CIDR source address for rule match.
- **source_address_name** (Optional, Set of String)  
  List of source address from address book for rule match.
- **source_port** (Optional, Set of String)  
  List of source port for rule match.  
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
    Change routing_instance with nat.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security nat static can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_nat_static.demo_nat nat_from_trust
```
