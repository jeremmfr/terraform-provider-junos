---
page_title: "Junos: junos_security_zone_book_address_set"
---

# junos_security_zone_book_address_set

Provides an address-set resource in address-book of security zone.

-> **Note:** The `junos_security_zone` resource needs to have `address_book_configure_singly` set to
true otherwise there will be a conflict between resources.

## Example Usage

```hcl
# Add an address-set
resource "junos_security_zone_book_address_set" "demo" {
  name    = "addressSet1"
  zone    = "theZone"
  address = ["address1"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of address-set.
- **zone** (Required, String, Forces new resource)  
  The name of security zone.
- **address** (Optional, Set of String)  
  Address to be included in this set.
- **address_set** (Optional, Set of String)  
  Address-set to be included in this set.
- **description** (Optional, String)  
  Description of address-set.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<zone>_-_<name>`.

## Import

Junos address-set in address-book of security zone can be imported using an id made up of
`<zone>_-_<name>`, e.g.

```shell
$ terraform import junos_security_zone_book_address_set.demo theZone_-_addressSet1
```
