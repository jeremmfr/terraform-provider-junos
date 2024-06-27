---
page_title: "Junos: junos_interface_logical"
---

# junos_interface_logical

Get configuration from a logical interface  
(as with a junos_interface_logical resource import).

## Example Usage

```hcl
# Search interface with IP
data "junos_interface_logical" "demo_ip" {
  match = "192.0.2.2/"
}
# Search interface with name
data "junos_interface_logical" "interface_fw_demo" {
  config_interface = "ge-0/0/3.0"
}
```

## Argument Reference

The following arguments are supported:

- **config_interface** (Optional, String)  
  Specifies the interface part for search.  
  Command is `show configuration interfaces <config_interface>`
- **match** (Optional, String)  
  Regex string to filter lines and find only one interface.

~> **NOTE:** If more or less than a single match is returned by the search, Terraform will fail.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **name** (String)  
  Name of logical interface (with dot).
- **description** (String)  
  Description for interface.
- **disable** (Boolean)  
  Interface disabled.
- **encapsulation** (String)  
  Logical link-layer encapsulation.
- **family_inet** (Block)  
  Family inet enabled and possible configuration.
  - **address** (Block List)  
    For each IPv4 address declared.  
    See [below for nested schema](#address-attributes-for-family_inet).
  - **dhcp** (Block)  
    Enable DHCP client and configuration.  
    See [below for nested schema](#dhcp-attributes-for-family_inet).
  - **filter_input** (String)  
    Filter applied to received packets.
  - **filter_output** (String)  
    Filter applied to transmitted packets.
  - **mtu** (Number)  
    Maximum transmission unit.
  - **rpf_check** (Block)  
    Reverse-path-forwarding checks enabled and possible configuration.  
    See [below for nested schema](#rpf_check-attributes).
  - **sampling_input** (Boolean)  
    Sample all packets input on this interface.
  - **sampling_output** (Boolean)  
    Sample all packets output on this interface.
- **family_inet6** (Block)  
  Family inet6 enabled and possible configuration.
  - **address** (Block List)  
    For each IPv6 address declared.  
    See [below for nested schema](#address-attributes-for-family_inet6).
  - **dad_disable** (Boolean)  
    Disable duplicate-address-detection.
  - **dhcpv6_client** (Block)  
    Enable DHCP client and configuration.  
    See [below for nested schema](#dhcpv6_client-attributes-for-family_inet6).
  - **filter_input** (String)  
    Filter applied to received packets.
  - **filter_output** (String)  
    Filter applied to transmitted packets.
  - **mtu** (Number)  
    Maximum transmission unit.
  - **rpf_check** (Block)  
    Reverse-path-forwarding checks enabled and possible configuration.  
    See [below for nested schema](#rpf_check-attributes).
  - **sampling_input** (Boolean)  
    Sample all packets input on this interface.
  - **sampling_output** (Boolean)  
    Sample all packets output on this interface.
- **routing_instance** (String)  
  Routing_instance where the interface is (if not default instance).
- **security_inbound_protocols** (Set of String)  
  The inbound protocols allowed.
- **security_inbound_services** (Set of String)  
  The inbound services allowed.
- **security_zone** (String)  
  Security zone where the interface is.
- **tunnel** (Block)  
  Tunnel parameters.  
  See [below for nested schema](#tunnel-arguments).
- **vlan_id** (Number)  
  Virtual LAN identifier value for 802.1q VLAN tags.

---

### tunnel arguments

- **destination** (String)  
  Tunnel destination.
- **source** (String)  
  Tunnel source.
- **allow_fragmentation** (Boolean)  
  Do not set DF bit on packets.
- **do_not_fragment** (Boolean)  
  Set DF bit on packets.
- **flow_label** (Number)  
  Flow label field of IP6-header (0..1048575).
- **no_path_mtu_discovery** (Boolean)  
  Don't enable path MTU discovery for tunnels.
- **path_mtu_discovery** (Boolean)  
  Enable path MTU discovery for tunnels.
- **routing_instance_destination** (String)  
  Routing instance to which tunnel ends belong.
- **traffic_class** (Number)  
  TOS/Traffic class field of IP-header (0..255).
- **ttl** (Number)  
  Time to live (1..255).

---

### address attributes for family_inet

- **cidr_ip** (String)  
  IPv4 address in CIDR format.
- **preferred** (Boolean)  
  Preferred address on interface.
- **primary** (Boolean)  
  Candidate for primary address in system.
- **vrrp_group** (Block List)  
  List of vrrp group configurations.  
  See [below for nested schema](#vrrp_group-attributes-for-address-in-family_inet).

---

### vrrp_group attributes for address in family_inet

- **identifier** (Number)  
  ID for vrrp.
- **virtual_address** (List of String)  
  List of address IP v4.
- **accept_data** (Boolean)  
  Accept packets destined for virtual IP address.
- **advertise_interval** (Number)  
  Advertisement interval (seconds).
- **advertisements_threshold** (Number)  
  Number of vrrp advertisements missed before declaring master down.
- **authentication_key** (String, Sensitive)  
  Authentication key.
- **authentication_type** (String)  
  Authentication type.
- **no_accept_data** (Boolean)  
  Don't accept packets destined for virtual IP address.
- **no_preempt** (Boolean)  
  Preemption not allowed.
- **preempt** (Boolean)  
  Preemption allowed.
- **priority** (Number)  
  Virtual router election priority.
- **track_interface** (Block List)  
  List of track_interface.
  - **interface** (String)  
    Interface tracked.
  - **priority_cost** (Number)  
    Value to subtract from priority when interface is down.
- **track_route** (Block List)  
  List of track_route.
  - **route** (String)  
    Route address tracked.
  - **routing_instance** (String)  
    Routing instance to which route belongs.
  - **priority_cost** (Number)  
    Value to subtract from priority when route is down.

---

### dhcp attributes for family_inet

- **srx_old_option_name** (Boolean)  
  For configuration, use the old option name `dhcp-client` instead of `dhcp`.  
  This is used for SRX devices with an older version of Junos.
- **client_identifier_ascii** (String)  
  Client identifier as an ASCII string.  
- **client_identifier_hexadecimal** (String)  
  Client identifier as a hexadecimal string.  
- **client_identifier_prefix_hostname** (Boolean)  
  Add prefix router host name to client-id option.
- **client_identifier_prefix_routing_instance_name** (Boolean)  
  Add prefix routing instance name to client-id option.
- **client_identifier_use_interface_description** (Boolean)  
  Use the interface description.  
- **client_identifier_userid_ascii** (String)  
  Add user id as an ASCII string to client-id option.
- **client_identifier_userid_hexadecimal** (String)  
  Add user id as a hexadecimal string to client-id option.
- **force_discover** (Boolean)  
  Send DHCPDISCOVER after DHCPREQUEST retransmission failure.
- **lease_time** (Number)  
  Lease time in seconds requested in DHCP client protocol packet (60..2147483647 seconds).  
- **lease_time_infinite** (Boolean)  
  Lease never expires.  
- **metric** (Number)  
  Client initiated default-route metric (0..255).
- **no_dns_install** (Boolean)  
  Do not install DNS information learned from DHCP server.
- **options_no_hostname** (Boolean)  
  Do not carry hostname (RFC option code is 12) in packet.
- **retransmission_attempt** (Number)  
  Number of attempts to retransmit the DHCP client protocol packet (0..50000).
- **retransmission_interval** (Number)  
  Number of seconds between successive retransmission (4..64 seconds).
- **server_address** (String)  
  DHCP Server-address.
- **update_server** (Boolean)  
  Propagate TCP/IP settings to DHCP server.
- **vendor_id** (String)  
  Vendor class id for the DHCP Client.

---

### address attributes for family_inet6

- **cidr_ip** (String)  
  IPv6 address in CIDR format.
- **preferred** (Boolean)  
  Preferred address on interface.
- **primary** (Boolean)  
  Candidate for primary address in system.
- **vrrp_group** (Block List)  
  List of vrrp group configurations.  
  See [below for nested schema](#vrrp_group-attributes-for-address-in-family_inet6).

---

### vrrp_group attributes for address in family_inet6

Same as [`vrrp_group` attributes for address in family_inet](#vrrp_group-attributes-for-address-in-family_inet)
block but without `authentication_key`, `authentication_type` and with  

- **virtual_link_local_address** (String)  
  Address IPv6 for Virtual link-local addresses.

---

### dhcpv6_client attributes for family_inet6

- **client_identifier_duid_type** (String)  
  DUID identifying a client.  
- **client_type** (String)  
  DHCPv6 client type.  
- **client_ia_type_na** (Boolean)  
  DHCPv6 client identity association type Non-temporary Address.  
- **client_ia_type_pd** (Boolean)  
  DHCPv6 client identity association type Prefix Address.  
- **no_dns_install** (Boolean)  
  Do not install DNS information learned from DHCP server.
- **prefix_delegating_preferred_prefix_length** (Number)  
  Client preferred prefix length (0..64).
- **prefix_delegating_sub_prefix_length** (Number)  
  The sub prefix length for LAN interfaces (1..127).
- **rapid_commit** (Boolean)  
  Option is used to signal the use of the two message exchange for address assignment.
- **req_option** (Set of String)  
  DHCPV6 client requested option configuration.
- **retransmission_attempt** (Number)  
  Number of attempts to retransmit the DHCPV6 client protocol packet (0..9).
- **update_router_advertisement_interface** (Set of String)  
  Interfaces on which to delegate prefix.
- **update_server** (Boolean)  
  Propagate TCP/IP settings to DHCP server.

---

### rpf_check attributes

- **fail_filter** (String)  
  Name of filter applied to packets failing RPF check.
- **mode_loose** (Boolean)  
  Reverse-path-forwarding use loose mode instead the strict mode.
