---
layout: "junos"
page_title: "Junos: junos_security_zone_book_address_set"
sidebar_current: "docs-junos-resource-security-zone-book-address-set"
description: |-
  Create an address-set in address-book of security zone (when Junos device supports it)
---

# junos_security_zone_book_address_set

Provides an address-set resource in address-book of security zone.

-> **Note:** The `security_zone` resource needs to have `address_book_configure_singly` set to true otherwise there will be a conflict between resources.

## Example Usage

```hcl
# Add an address-set
resource junos_security_zone_book_address_set "demo" {
  name    = "addressSet1"
  zone    = "theZone"
  address = ["address1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of address-set.
* `zone` - (Required, Forces new resource)(`String`) The name of security zone.
* `address` - (Required)(`ListOfString`) Address to be included in this set.
* `description` - (Optional)(`String`) Description of address-set.

## Import

Junos address-set in address-book of security zone can be imported using an id made up of `<zone>_-_<name>`, e.g.

```shell
$ terraform import junos_security_zone_book_address_set.demo theZone_-_addressSet1
```
