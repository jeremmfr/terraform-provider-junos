---
page_title: "Junos: junos_interface_logical"
---

# junos_interface_logical

Provides a logical interface resource.

## Example Usage

```hcl
resource "junos_interface_logical" "interface_fw_demo_100" {
  name        = "ae0.100"
  description = "interfaceFwDemo100"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of logical interface (with dot).
- **st0_also_on_destroy** (Optional, Boolean)  
  When destroy this resource, if the name has prefix `st0.`,
  delete all configurations (not keep empty st0 interface).  
  Usually, `st0.x` interfaces are completely deleted by `junos_interface_st0_unit` resource
  because of the dependency, but only if st0.x interface is empty or disable.
- **description** (Optional, String)  
  Description for interface.
- **disable** (Optional, Boolean)  
  Disable this logical interface.
- **encapsulation** (Optional, String)  
  Logical link-layer encapsulation.
- **family_inet** (Optional, Block)  
  Enable family inet and add configurations if specified.
  - **address** (Optional, Block List)  
    For each IPv4 address to declare.  
    Conflict with `dhcp`.  
    See [below for nested schema](#address-arguments-for-family_inet).
  - **dhcp** (Optional, Block)  
    Enable DHCP client and configuration.  
    Conflict with `address`.  
    See [below for nested schema](#dhcp-arguments-for-family_inet).
  - **filter_input** (Optional, String)  
    Filter to be applied to received packets.
  - **filter_output** (Optional, String)  
    Filter to be applied to transmitted packets.
  - **mtu** (Optional, Number)  
    Maximum transmission unit.
  - **rpf_check** (Optional, Block)  
    Enable reverse-path-forwarding checks on this interface.  
    See [below for nested schema](#rpf_check-arguments).
  - **sampling_input** (Optional, Boolean)  
    Sample all packets input on this interface.
  - **sampling_output** (Optional, Boolean)  
    Sample all packets output on this interface.
- **family_inet6** (Optional, Block)  
  Enable family inet6 and add configurations if specified.
  - **address** (Optional, Block List)  
    For each IPv6 address to declare.  
    Conflict with `dhcpv6_client`.  
    See [below for nested schema](#address-arguments-for-family_inet6).
  - **dad_disable** (Optional, Boolean)  
    Disable duplicate-address-detection.
  - **dhcpv6_client** (Optional, Block)  
    Enable DHCP client and configuration.  
    Conflict with `address`.  
    See [below for nested schema](#dhcpv6_client-arguments-for-family_inet6).
  - **filter_input** (Optional, String)  
    Filter to be applied to received packets.
  - **filter_output** (Optional, String)  
    Filter to be applied to transmitted packets.
  - **mtu** (Optional, Number)  
    Maximum transmission unit.
  - **rpf_check** (Optional, Block)  
    Enable reverse-path-forwarding checks on this interface.  
    See [below for nested schema](#rpf_check-arguments).
  - **sampling_input** (Optional, Boolean)  
    Sample all packets input on this interface.
  - **sampling_output** (Optional, Boolean)  
    Sample all packets output on this interface.
- **routing_instance** (Optional, String)  
  Add this interface in routing_instance.  
  Need to be created before.
- **security_inbound_protocols** (Optional, Set of String)  
  The inbound protocols allowed.  
  Must be a list of Junos protocols.  
  `security_zone` need to be set.
- **security_inbound_services** (Optional, Set of String)  
  The inbound services allowed.  
  Must be a list of Junos services.  
  `security_zone` need to be set.
- **security_zone** (Optional, String)  
  Add this interface in a security zone.  
  Need to be created before.
- **tunnel** (Optional, Block)  
  Tunnel parameters.  
  See [below for nested schema](#tunnel-arguments).
- **vlan_id** (Optional, Computed, Number)  
  Virtual LAN identifier value for 802.1q VLAN tags.  
  If not set, computed with `name` of interface (ge-0/0/0.100 = 100)
  except if name has `.0` suffix or `st0.`, `irb.`, `vlan.` prefix.
- **vlan_no_compute** (Optional, Boolean)  
  Disable the automatic compute of the `vlan_id` argument when not set.  
  Unnecessary if name has `.0` suffix or `st0.`, `irb.`, `vlan.` prefix because it's already disabled.
- **virtual_gateway_accept_data** (Optional, Boolean)
  Accept packets destined for virtual gateway address

---

### tunnel arguments

- **destination** (Required, String)  
  Tunnel destination.
- **source** (Required, String)  
  Tunnel source.
- **allow_fragmentation** (Optional, Boolean)  
  Do not set DF bit on packets.  
  Conflict with `do_not_fragment`.
- **do_not_fragment** (Optional, Boolean)  
  Set DF bit on packets.  
  Conflict with `allow_fragmentation`.
- **flow_label** (Optional, Number)  
  Flow label field of IP6-header (0..1048575).
- **path_mtu_discovery** (Optional, Boolean)  
  Enable path MTU discovery for tunnels.  
  Conflict with `no_path_mtu_discovery`.
- **no_path_mtu_discovery** (Optional, Boolean)  
  Don't enable path MTU discovery for tunnels.  
  Conflict with `path_mtu_discovery`.
- **routing_instance_destination** (Optional, String)  
  Routing instance to which tunnel ends belong.
- **traffic_class** (Optional, Number)  
  TOS/Traffic class field of IP-header (0..255).
- **ttl** (Optional, Number)  
  Time to live (1..255).

---

### address arguments for family_inet

- **cidr_ip** (Required, String)  
  IPv4 address in CIDR format.
- **preferred** (Optional, Boolean)  
  Preferred address on interface.
- **primary** (Optional, Boolean)  
  Candidate for primary address in system.
- **virtual_gateway_address** (Optional, String)
  IPv4 address of Virtual Gateway.
- **vrrp_group** (Optional, Block List)
  For each vrrp group to declare.
  See [below for nested schema](#vrrp_group-arguments-for-address-in-family_inet).

---

### vrrp_group arguments for address in family_inet

- **identifier** (Required, Number)  
  ID for vrrp.
- **virtual_address** (Required, List of String)  
  Virtual IP addresses.
- **accept_data** (Optional, Boolean)  
  Accept packets destined for virtual IP address.  
  Conflict with `no_accept_data` when apply.
- **no_accept_data** (Optional, Boolean)  
  Don't accept packets destined for virtual IP address.  
  Conflict with `accept_data` when apply.
- **advertise_interval** (Optional, Number)  
  Advertisement interval (seconds).
- **advertisements_threshold** (Optional, Number)  
  Number of vrrp advertisements missed before declaring master down.
- **authentication_key** (Optional, String, Sensitive)  
  Authentication key.
- **authentication_type** (Optional, String)  
  Authentication type.  
  Need to be `md5` or `simple`.
- **preempt** (Optional, Boolean)  
  Allow preemption.  
  Conflict with `no_preempt` when apply.
- **no_preempt** (Optional, Boolean)  
  Don't allow preemption.  
  Conflict with `preempt` when apply.
- **priority** (Optional, Number)  
  Virtual router election priority.
- **track_interface** (Optional, Block List)  
  For each interface to track in VRRP group.
  - **interface** (Required, String)  
    Name of interface.
  - **priority_cost** (Required, Number)  
    Value to subtract from priority when interface is down.
- **track_route** (Optional, Block List)  
  For each route to track in VRRP group.
  - **route** (Required, String)  
    Route address.
  - **routing_instance** (Required, String)  
    Routing instance to which route belongs, or `default`.
  - **priority_cost** (Required, Number)  
    Value to subtract from priority when route is down.

---

### dhcp arguments for family_inet

- **srx_old_option_name** (Optional, Boolean)  
  For configuration, use the old option name `dhcp-client` instead of `dhcp`.  
  This is useful for SRX devices with an older version of Junos.
- **client_identifier_ascii** (Optional, String)  
  Client identifier as an ASCII string.  
  Conflict witch `client_identifier_hexadecimal`.
- **client_identifier_hexadecimal** (Optional, String)  
  Client identifier as a hexadecimal string.  
  Conflict witch `client_identifier_ascii`.
- **client_identifier_prefix_hostname** (Optional, Boolean)  
  Add prefix router host name to client-id option.
- **client_identifier_prefix_routing_instance_name** (Optional, Boolean)  
  Add prefix routing instance name to client-id option.
- **client_identifier_use_interface_description** (Optional, Boolean)  
  Use the interface description.  
  Need to be `device` or `logical`.
- **client_identifier_userid_ascii** (Optional, String)  
  Add user id as an ASCII string to client-id option.
- **client_identifier_userid_hexadecimal** (Optional, String)  
  Add user id as a hexadecimal string to client-id option.
- **force_discover** (Optional, Boolean)  
  Send DHCPDISCOVER after DHCPREQUEST retransmission failure.
- **lease_time** (Optional, Number)  
  Lease time in seconds requested in DHCP client protocol packet (60..2147483647 seconds).  
  Conflict witch `lease_time_infinite`.
- **lease_time_infinite** (Optional, Boolean)  
  Lease never expires.  
  Conflict witch `lease_time`.
- **metric** (Optional, Number)  
  Client initiated default-route metric (0..255).
- **no_dns_install** (Optional, Boolean)  
  Do not install DNS information learned from DHCP server.
- **options_no_hostname** (Optional, Boolean)  
  Do not carry hostname (RFC option code is 12) in packet.
- **retransmission_attempt** (Optional, Number)  
  Number of attempts to retransmit the DHCP client protocol packet (0..50000).
- **retransmission_interval** (Optional, Number)  
  Number of seconds between successive retransmission (4..64 seconds).
- **server_address** (Optional, String)  
  DHCP Server-address.
- **update_server** (Optional, Boolean)  
  Propagate TCP/IP settings to DHCP server.
- **vendor_id** (Optional, String)  
  Vendor class id for the DHCP Client.

---

### address arguments for family_inet6

- **cidr_ip** (Required, String)  
  IPv6 address in CIDR format.
- **preferred** (Optional, Boolean)  
  Preferred address on interface.
- **primary** (Optional, Boolean)  
  Candidate for primary address in system.
- **vrrp_group** (Optional, Block List)  
  For each vrrp group to declare.  
  See [below for nested schema](#vrrp_group-arguments-for-address-in-family_inet6).

---

### vrrp_group arguments for address in family_inet6

Same as [`vrrp_group` arguments for address in family_inet](#vrrp_group-arguments-for-address-in-family_inet)
block but without `authentication_key`, `authentication_type` and with

- **virtual_link_local_address** (Required, String)  
  Address IPv6 for Virtual link-local addresses.

---

### dhcpv6_client arguments for family_inet6

- **client_identifier_duid_type** (Required, String)  
  DUID identifying a client.  
  Need to be `duid-ll`, `duid-llt` or `vendor`.
- **client_type** (Required, String)  
  DHCPv6 client type.  
  Need to be `autoconfig` or `stateful`.
- **client_ia_type_na** (Optional, Boolean)  
  DHCPv6 client identity association type Non-temporary Address.  
  At least one of `client_ia_type_na`, `client_ia_type_pd` need to be true.
- **client_ia_type_pd** (Optional, Boolean)  
  DHCPv6 client identity association type Prefix Address.  
  At least one of `client_ia_type_na`, `client_ia_type_pd` need to be true.
- **no_dns_install** (Optional, Boolean)  
  Do not install DNS information learned from DHCP server.
- **prefix_delegating_preferred_prefix_length** (Optional, Number)  
  Client preferred prefix length (0..64).
- **prefix_delegating_sub_prefix_length** (Optional, Number)  
  The sub prefix length for LAN interfaces (1..127).
- **rapid_commit** (Optional, Boolean)  
  Option is used to signal the use of the two message exchange for address assignment.
- **req_option** (Optional, Set of String)  
  DHCPV6 client requested option configuration.
- **retransmission_attempt** (Optional, Number)  
  Number of attempts to retransmit the DHCPV6 client protocol packet (0..9).
- **update_router_advertisement_interface** (Optional, Set of String)  
  Interfaces on which to delegate prefix.
- **update_server** (Optional, Boolean)  
  Propagate TCP/IP settings to DHCP server.

---

### rpf_check arguments

- **fail_filter** (Optional, String)  
  Name of filter applied to packets failing RPF check.
- **mode_loose** (Optional, Boolean)  
  Use reverse-path-forwarding loose mode instead the strict mode.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos interface can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_interface_logical.interface_fw_demo_100 ae.100
```
