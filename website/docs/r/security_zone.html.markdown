---
layout: "junos"
page_title: "Junos: junos_security_zone"
sidebar_current: "docs-junos-resource-security-zone"
description: |-
  Create a security zone (when Junos device supports it)
---

# junos_security_zone

Provides a security zone resource.

## Example Usage

```hcl
# Add a security zone
resource "junos_security_zone" "DemoZone" {
  name              = "DemoZone"
  inbound_protocols = ["bgp"]
  address_book {
    name    = "DemoAddress"
    network = "192.0.2.0/25"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of security zone.
* `inbound_services` - (Optional)(`ListOfString`) The inbound services allowed. Must be a list of Junos services
* `inbound_protocols` - (Optional)(`ListOfString`) The inbound protocols allowed. Must be a list of Junos protocols
* `address_book` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each address to declare.
  The `address_book` block supports:
  * `name` - (Required)(`String`) Name of address
  * `network` - (Required)(`String`) CIDR of address
* `address_book_set` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each address-set to declare.
    The `address_book` block supports:
  * `name` - (Required)(`String`) Name of address-set
  * `address` - (Required)(`ListOfString`) List of address names

## Import

Junos security zone can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_zone.DemoZone DemoZone
```
