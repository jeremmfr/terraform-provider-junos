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

* `name` - (Required, Forces new resource)(`String`) The name of destination nat.
* `from` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'from' configuration.
  * `type` - (Required)(`String`) Type of from options. Need to be 'interface', 'routing-instance' or 'zone'
  * `value`  - (Required)(`String`) Name of interface, routing-instance or zone for from options
* `rule` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each rule to declare. See the [`rule` arguments](#rule-arguments) block.

---
#### rule arguments
* `name` - (Required)(`String`) Name of rule
* `destination_address` - (Required)(`String`) CIDR for match destination address
* `then` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'then' action.
  * `type` - (Required)(`String`) Type of destination nat. Need to be 'pool' or 'off'
  * `pool` - (Optional)(`String`) Name of nat destination pool when type pool

## Import

Junos security nat destination can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_nat_destination.demo_dnat dnat_from_untrust
```
