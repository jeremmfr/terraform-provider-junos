---
page_title: "Junos: junos_security_zone"
---

# junos_security_zone

Provides a security zone resource.

## Example Usage

```hcl
# Add a security zone
resource "junos_security_zone" "demo_zone" {
  name              = "DemoZone"
  inbound_protocols = ["bgp"]
  address_book {
    name    = "DemoAddress"
    network = "192.0.2.0/25"
  }
}
```

## Argument Reference

-> **Note** The interfaces can be configured with the `junos_interface_logical` resource and the
  `security_zone`, `security_inbound_protocols` and `security_inbound_services` arguments.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of security zone.
- **address_book** (Optional, Block Set)  
  For each name of network address to declare.
  - **name** (Required, String)  
    Name of network address.
  - **network** (Required, String)  
    CIDR value of network address (`192.0.0.0/24`).
  - **description** (Optional, String)  
    Description of network address.
- **address_book_configure_singly** (Optional, Boolean)  
  Disable management of address-book in this resource to be able to manage them with specific
  resources.  
  Conflict with `address_book_*`.
- **address_book_dns** (Optional, Block Set)  
  For each name of dns-name address to declare.
  - **name** (Required, String)  
    Name of dns name address.
  - **fqdn** (Required, String)  
    Fully qualified domain name.
  - **description** (Optional, String)  
    Description of dns name address.
  - **ipv4_only** (Optional, Boolean)  
    IPv4 dns address.
  - **ipv6_only** (Optional, Boolean)  
    IPv6 dns address.
- **address_book_range** (Optional, Block Set)  
  For each name of range-address to declare.
  - **name** (Required, String)  
    Name of range address.
  - **from** (Required, String)  
    Lower limit of address range.
  - **to** (Required, String)  
    Upper limit of address range.
  - **description** (Optional, String)  
    Description of range address.
- **address_book_set** (Optional, Block Set)  
  For each name of address-set to declare.
  - **name** (Required, String)  
    Name of address-set.
  - **address** (Optional, Set of String)  
    List of address names.
  - **address_set** (Optional, Set of String)  
    List of address-set names.
  - **description** (Optional, String)  
    Description of address-set.
- **address_book_wildcard** (Optional, Block Set)  
  For each name of wildcard-address to declare.
  - **name** (Required, String)  
    Name of wildcard address.
  - **network** (Required, String)  
    Numeric IPv4 wildcard address with in the form of a.d.d.r/netmask.
  - **description** (Optional, String)  
    Description of wildcard address.
- **advance_policy_based_routing_profile** (Optional, String)  
  Enable Advance Policy Based Routing on this zone with a profile.
- **application_tracking** (Optional, Boolean)  
  Enable Application tracking support for this zone.
- **description** (Optional, String)  
  Text description of zone.
- **inbound_protocols** (Optional, Set of String)  
  The inbound protocols allowed.  
  Must be a list of Junos protocols.
- **inbound_services** (Optional, Set of String)  
  The inbound services allowed.  
  Must be a list of Junos services.
- **reverse_reroute** (Optional, Boolean)  
  Enable Reverse route lookup when there is change in ingress interface.
- **screen** (Optional, String)  
  Name of ids option object (screen) applied to the zone.
- **source_identity_log** (Optional, Boolean)  
  Show user and group info in session log for this zone.
- **tcp_rst** (Optional, Boolean)  
  Send RST for NON-SYN packet not matching TCP session.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security zone can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_zone.demo_zone DemoZone
```
