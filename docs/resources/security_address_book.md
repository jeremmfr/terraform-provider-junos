---
page_title: "Junos: junos_security_address_book"
---

# junos_security_address_book

Provides a security address book resource.

## Example Usage

```hcl
# Add an address book with entries
resource "junos_security_address_book" "testAddressBook" {
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

- **name** (Optional, String, Forces new resource)  
  The name of address book.  
  Defaults to `global`.
- **description** (Optional, String)  
  The description of the address book.
- **attach_zone** (Optional, List of String)  
  List of zones to attach address book to.  
  **NOTE:** Cannot be set on global address book.
- **network_address** (Optional, Block Set)  
  For each name of network address.
  - **name** (Required, String)  
    Name of network address.
  - **value** (Required, String)  
    CIDR value of network address (`192.0.0.0/24`).
  - **description** (Optional, String)  
    Description of network address.
- **dns_name** (Optional, Block Set)  
  For each name of dns name address.
  - **name** (Required, String)  
    Name of dns name address.
  - **value** (Required, String)  
    DNS name string value (`juniper.net`).
  - **description** (Optional, String)  
    Description of dns name address.
  - **ipv4_only** (Optional, Boolean)  
    IPv4 dns address.
  - **ipv6_only** (Optional, Boolean)  
    IPv6 dns address.
- **range_address** (Optional, Block Set)  
   For each name of range address.
  - **name** (Required, String)  
    Name of range address.
  - **from** (Required, String)  
    IP address of start of range.
  - **to** (Required, String)  
    IP address of end of range.
  - **description** (Optional, String)  
    Description of range address.
- **wildcard_address** (Optional, Block Set)  
  For each name of wildcard address.
  - **name** (Required, String)  
    Name of wildcard address.
  - **value** (Required, String)  
    Network and mask of wildcard address (`192.0.0.0/255.255.0.255`).
  - **description** (Optional, String)  
    Description of wildcard address.
- **address_set** (Optional, Block Set)  
  For each name of address-set to declare.
  - **name** (Required, String)  
    Name of address-set.
  - **address** (Optional, Set of String)  
    List of address names.
  - **address_set** (Optional, Set of String)  
    List of address-set names.
  - **description** (Optional, String)  
    Description of address-set.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security address book can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_address_book.global global
```
