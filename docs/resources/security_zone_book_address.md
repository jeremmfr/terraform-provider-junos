---
page_title: "Junos: junos_security_zone_book_address"
---

# junos_security_zone_book_address

Provides an address resource in address-book of security zone.

-> **Note:** The `junos_security_zone` resource needs to have `address_book_configure_singly` set to
true otherwise there will be a conflict between resources.

## Example Usage

```hcl
# Add an address
resource "junos_security_zone_book_address" "demo" {
  name = "address1"
  zone = "theZone"
  cidr = "192.0.2.0/25"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of address.
- **zone** (Required, String, Forces new resource)  
  The name of security zone.
- **cidr** (Optional, String)  
  CIDR of address.
- **description** (Optional, String)  
  Description of address.
- **dns_ipv4_only** (Optional, Boolean)  
  IPv4 dns address.
- **dns_ipv6_only** (Optional, Boolean)  
  IPv6 dns address.
- **dns_name** (Optional, String)  
  DNS address name.
- **range_from** (Optional, String)  
  Lower limit of address range.
- **range_to** (Optional, String)  
  Upper limit of address range.
- **wildcard** (Optional, String)  
  Numeric IPv4 wildcard address in the form of `a.d.d.r/netmask`.

-> **Note:** One of `cidr`, `dns_name`, `range_from` or `wildcard` arguments need to be set.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<zone>_-_<name>`.

## Import

Junos address in address-book of security zone can be imported using an id made up of
`<zone>_-_<name>`, e.g.

```shell
$ terraform import junos_security_zone_book_address.demo theZone_-_address1
```
