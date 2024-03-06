---
page_title: "Junos: junos_security_zone"
---

# junos_security_zone

Get configuration from a security zone.

## Example Usage

```hcl
# Read security zone configuration
data "junos_security_zone" "demo_zone" {
  name = "DemoZone"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  The name of security zone.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **address_book** (Block Set)  
  For each name of address.
  - **name** (String)  
    Name of network address.
  - **network** (String)  
    CIDR value of network address.
  - **description** (String)  
    Description of network address.
- **address_book_dns** (Block Set)  
  For each name of dns-name address.
  - **name** (String)  
    Name of dns name address.
  - **fqdn** (String)  
    Fully qualified domain name.
  - **description** (String)  
    Description of dns name address.
  - **ipv4_only** (Boolean)  
    IPv4 dns address.
  - **ipv6_only** (Boolean)  
    IPv6 dns address.
- **address_book_range** (Block Set)  
  For each name of range-address.
  - **name** (String)  
    Name of range address.
  - **from** (String)  
    Lower limit of address range.
  - **to** (String)  
    Upper limit of address range.
  - **description** (String)  
    Description of range address.
- **address_book_set** (Block Set)  
  For each name of address-set.
  - **name** (String)  
    Name of address-set.
  - **address** (Set of String)  
    List of address names.
  - **address_set** (Set of String)  
    List of address-set names.
  - **description** (String)  
    Description of address-set.
- **address_book_wildcard** (Block Set)  
  For each name of wildcard-address.
  - **name** (String)  
    Name of wildcard address.
  - **network** (String)  
    Numeric IPv4 wildcard address with in the form of a.d.d.r/netmask.
  - **description** (String)  
    Description of wildcard address.
- **advance_policy_based_routing_profile** (String)  
  Enable Advance Policy Based Routing on this zone with a profile.
- **application_tracking** (Boolean)  
  Enable Application tracking support for this zone.
- **description** (String)  
  Text description of zone.
- **inbound_protocols** (Set of String)  
  The inbound protocols allowed.  
- **inbound_services** (Set of String)  
  The inbound services allowed.  
- **interface** (Block Set)  
  List of interfaces in security-zone.  
  - **name** (String)  
    Interface name.
  - **inbound_protocols** (Set of String)  
    Protocol type of incoming traffic to accept.
  - **inbound_services** (Set of String)  
    Type of incoming system-service traffic to accept.
- **reverse_reroute** (Boolean)  
  Enable Reverse route lookup when there is change in ingress interface.
- **screen** (String)  
  Name of ids option object (screen) applied to the zone.
- **source_identity_log** (Boolean)  
  Show user and group info in session log for this zone.
- **tcp_rst** (Boolean)  
  Send RST for NON-SYN packet not matching TCP session.
