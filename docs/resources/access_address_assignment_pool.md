---
page_title: "Junos: junos_access_address_assignment_pool"
---

# junos_access_address_assignment_pool

Provides an access address-assignment pool.

## Example Usage

```hcl
# Add an access address-assignment pool
resource "junos_access_address_assignment_pool" "demo_dhcp_pool" {
  name = "demo_dhcp_pool"
  family {
    type    = "inet"
    network = "192.0.2.128/25"
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Address pool name.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for pool.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`
- **family** (Required, Block)  
  Configure address family (`inet` or `inet6`).  
  See [below for nested schema](#family-arguments).
- **active_drain** (Optional, Boolean)  
  Notify client of pool active drain mode.
- **hold_down** (Optional, Boolean)  
  Place pool in passive drain mode.
- **link** (Optional, String)  
  Address pool link name.

---

### family arguments

- **type** (Required, String)  
  Type of family.  
  Need to be `inet` or `inet6`.
- **network** (Required, String)  
  Network address of pool.
- **dhcp_attributes** (Optional, Block)  
  DHCP options and match criteria.  
  See [below for nested schema](#dhcp_attributes-arguments).
- **excluded_address** (Optional, Set of String)  
  Excluded Addresses.  
  Need to be valid IP addresses.
- **excluded_range** (Optional, Block List)  
  For each name of excluded address range to declare.
  - **name** (Required, String)  
    Range name.
  - **low** (Required, String)  
    Lower limit of excluded address range.  
    Need to be a valid IP address.
  - **high** (Required, String)  
    Upper limit of excluded address range.  
    Need to be a valid IP address.
- **host** (Optional, Block List)  
  For each name of host to declare.  
  `type` need to be `inet`.
  - **name** (Required, String)  
    Hostname.
  - **hardware_address** (Required, String)  
    Hardware address.  
    Need to be a valid MAC address.
  - **Reserved address** (Required, String)  
    Hardware address.  
    Need to be a valid IPv4 address.
- **inet_range** (Optional, Block List)  
  For each name of address range to declare.  
  `type` need to be `inet`.
  - **name** (Required, String)  
    Range name.
  - **low** (Required, String)  
    Lower limit of address range.  
    Need to be a valid IPv4 address.
  - **high** (Required, String)  
    Upper limit of address range.  
    Need to be a valid IPv4 address.
- **inet6_range** (Optional, Block List)  
  For each name of address range to declare.  
  Need to set one of `prefix_length` or `low` + `high`.  
  `type` need to be `inet6`.
  - **name** (Required, String)  
    Range name.
  - **low** (Optional, String)  
    Lower limit of IPv6 address range.  
    Need to be a valid IPv6 address with mask.
  - **high** (Optional, String)  
    Upper limit of IPv6 address range.  
    Need to be a valid IPv6 address with mask.
  - **prefix_length** (Optional, Number)  
    IPv6 delegated prefix length (1..128).
- **xauth_attributes_primary_dns** (Optional, String)  
  Specify the primary-dns IP address.  
  Need to be a valid IPv4 address with mask.  
  `type` need to be `inet`.
- **xauth_attributes_primary_wins** (Optional, String)  
  Specify the primary-wins IP address.  
  Need to be a valid IPv4 address with mask.  
  `type` need to be `inet`.
- **xauth_attributes_secondary_dns** (Optional, String)  
  Specify the secondary-dns IP address.  
  Need to be a valid IPv4 address with mask.  
  `type` need to be `inet`.
- **xauth_attributes_secondary_wins** (Optional, String)  
  Specify the secondary-wins IP address.  
  Need to be a valid IPv4 address with mask.  
  `type` need to be `inet`.

---

### dhcp_attributes arguments

- **boot_file** (Optional, String)  
  Boot filename advertised to clients.
- **boot_server** (Optional, String)  
  Boot server advertised to clients.
- **dns_server** (Optional, List of String)  
  IPv6 domain name servers available to the client.  
  `type` need to be `inet6`.
- **domain_name** (Optional, String)  
  Domain name advertised to clients.
- **exclude_prefix_len** (Optional, Number)  
  Length of IPv6 prefix to be excluded from delegated prefix (1..128).  
  `type` need to be `inet6`.
- **grace_period** (Optional, Number)  
  Grace period for leases (seconds).
- **maximum_lease_time** (Optional, Number)  
  Maximum lease time advertised to clients (seconds).  
  Conflict with `maximum_lease_time_infinite`, `preferred_lifetime*`, `valid_lifetime*`.
- **maximum_lease_time_infinite** (Optional, Boolean)  
  Lease time can be infinite.  
  Conflict with `maximum_lease_time`, `preferred_lifetime*`, `valid_lifetime*`.
- **name_server** (Optional, List of String)  
  IPv4 domain name servers available to the client.  
  Need to be valid IPv4 addresses.  
- **netbios_node_type** (Optional, String)  
  Type of NETBIOS node advertised to clients.  
  Need to be `b-node`, `h-node`, `m-node` or `p-node`.
- **next_server** (Optional, String)  
  Next server that clients need to contact.  
  Need to be a valid IPv4 address.
<!-- markdownlint-disable -->
- **option** (Optional, String)  
  DHCP option.  
  Format need to match `^\d+ (array )?(byte|flag|hex-string|integer|ip-address|short|string|unsigned-integer|unsigned-short) .*$`.
<!-- markdownlint-restore -->
- **option_match_82_circuit_id** (Optional, Block List)  
  For each value to declare, circuit ID portion of the option 82.
  - **value** (Required, String)  
    Match value.
  - **range** (Required, String)  
    Range name.
- **option_match_82_remote_id** (Optional, Block List)  
  For each value to declare, remote ID portion of the option 82.
  - **value** (Required, String)  
    Match value.
  - **range** (Required, String)  
    Range name.
- **preferred_lifetime** (Optional, Number)  
  Preferred lifetime advertised to clients (seconds).  
  `type` need to be `inet6`.  
  Conflict with `preferred_lifetime_infinite`, `maximum_lease_time*`.
- **preferred_lifetime_infinite** (Optional, Boolean)
  Lease time can be infinite.  
  `type` need to be `inet6`.  
  Conflict with `preferred_lifetime`, `maximum_lease_time*`.
- **propagate_ppp_settings** (Optional, Set of String)  
  PPP interface name for propagating DNS/WINS settings.
- **propagate_settings** (Optional, String)  
  Interface name for propagating TCP/IP Settings to pool.
- **router** (Optional, String)  
  Routers advertised to clients.  
  Need to be valid IPv4 addresses.
- **server_identifier** (Optional, String)  
  Server Identifier - IP address value.  
  Need to be a valid IPv4 address.
- **sip_server_inet_address** (Optional, List of String)  
  SIP servers list of IPv4 addresses available to the client.  
  Need to be valid IPv4 addresses.
- **sip_server_inet_domain_name** (Optional, List of String)  
  SIP server domain name available to clients.
- **sip_server_inet6_address** (Optional, List of String)  
  SIP Servers list of IPv6 addresses available to the client.  
  Need to be valid IPv6 addresses.  
  `type` need to be `inet6`.  
- **sip_server_inet6_domain_name** (Optional, String)  
  SIP server domain name available to clients.  
  `type` need to be `inet6`.  
- **t1_percentage** (Optional, Number)  
  T1 time as percentage of preferred lifetime or max lease (0..100 percent).  
  Conflict with `t(1|2)_(renewal|rebinding)_time`.
- **t1_renewal_time** (Optional, Number)  
  T1 renewal time (seconds).  
  Conflict with `t(1|2)_percentage`.
- **t2_percentage** (Optional, Number)  
  T2 time as percentage of preferred lifetime or max lease (0..100 percent).  
  Conflict with `t(1|2)_(renewal|rebinding)_time`.
- **t2_rebinding_time**(Optional, Number)  
  T2 rebinding time (seconds).  
  Conflict with `t(1|2)_percentage`.
- **tftp_server** (Optional, String)  
  TFTP server IP address advertised to clients.  
  Need to be a valid IPv4 address.
- **valid_lifetime** (Optional, Number)  
  Valid lifetime advertised to clients (seconds).  
  `type` need to be `inet6`.  
  Conflict with `valid_lifetime_infinite`, `maximum_lease_time*`.
- **valid_lifetime_infinite** (Optional, Boolean)  
  Lease time can be infinite.  
  `type` need to be `inet6`.  
  Conflict with `valid_lifetime`, `maximum_lease_time*`.
- **wins_server** (Optional, List of String)  
  WINS name servers.  
  Need to be valid IPv4 addresses.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>`.

## Import

Junos access address-assignment pool can be imported using an id made up of
`<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_access_address_assignment_pool.demo_dhcp_pool demo_dhcp_pool_-_default
```
