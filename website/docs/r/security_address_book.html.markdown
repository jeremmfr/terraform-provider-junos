---
layout: "junos"
page_title: "Junos: junos_security_address_book"
sidebar_current: "docs-junos-resource-security-address-book"
description: |-
  Create a security address book (when Junos device supports it)
---

# junos_security_address_book

Provides a security address book resource.

## Example Usage

```hcl
# Add an address book with entries
resource junos_security_address_book "testAddressBook" {
  name        = "testAddressBook"
  attach_zone = ["SecurityZone"]
  network_address {
    name        = "DemoNetworkAddress"
    description = "Test Description"
    value       = "192.0.0.0/24"
  }
  network_address {
    name        = "DemoNetworkAddress2"
    description = "Test Description 2"
    value       = "192.1.0.0/24"
  }
  dns_name {
    name  = "DemoDnsName"
    value = "juniper.net"
  }
  range_address {
    name = "DemoRangeAddress"
    from = "192.0.0.1"
    to   = "192.0.0.10"
  }
  wildcard_address {
    name  = "DemoWildcardAddress"
    value = "juniper.net"
  }
  address_set {
    name    = "DemoAddressSet"
    address = ["DemoDnsName", "DemoWildcardAddress", "DemoRangeAddress"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, Forces new resource)(`String`) The name of address book. Defaults to `global`.
* `description` - (Optional)(`String`) The description of the address book.
* `attach_zone` - (Optional)(`ListOfString`) List of zones to attach address book to. **NOTE:** Cannot be set on global address book.
* `network_address` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each network address.
  * `name` - (Required)(`String`) Name of network address.
  * `description` - (Optional)(`String`) Description of network address.
  * `value` - (Required)(`String`) CIDR value of network address. `192.0.0.0/24`
* `wildcard_address` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each wildcard address.
  * `name` - (Required)(`String`) Name of wildcard address.
  * `description` - (Optional)(`String`) Description of network address.
  * `value` - (Required)(`String`) Nework and Mask of wildcard address. `192.0.0.0/255.255.0.255`
* `dns_name` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each dns name address.
  * `name` - (Required)(`String`) Name of dns name address.
  * `description` - (Optional)(`String`) Description of dns name address.
  * `value` - (Required)(`String`) DNS name string value. `juniper.net`
* `range_address` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html))  Can be specified multiple times for each range address.
  * `name` - (Required)(`String`) Name of range address.
  * `description` - (Optional)(`String`) Description of range address.
  * `from` - (Required)(`String`) IP address of start of range.
  * `to` - (Required)(`String`) IP address of end of range.
* `address_book_set` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each address-set to declare.
  * `name` - (Required)(`String`) Name of address-set.
  * `description` - (Optional)(`String`) Description of address-set.
  * `address` - (Required)(`ListOfString`) List of address names.

## Import

Junos security address book can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_address_book.global global
```
