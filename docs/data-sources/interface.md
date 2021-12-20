---
layout: "junos"
page_title: "Junos: junos_interface"
sidebar_current: "docs-junos-data-source-interface"
description: |-
  Get information on an Interface (as with an junos_interface resource import)
---

# junos_interface

Get information on an Interface

!> **NOTE:** Since v1.11.0, this data soure is **deprecated**.  
For more consistency, functionalities of this data source have been splitted in two new data source
`junos_interface_physical` and `junos_interface_logical`.

## Example Usage

```hcl
# Search interface with IP
data junos_interface "demo_ip" {
  match = "192.0.2.2/"
}
# Search interface with name
data junos_interface "interface_fw_demo" {
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

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.
- **name** (String)  
  Name of interface or unit interface (with dot).
- **description** (String)  
  Description for interface.
- **vlan_tagging** (Boolean)  
  802.1q VLAN tagging support.
- **vlan_taggind_id** (Number)  
  802.1q VLAN ID for unit interface.
- **inet** (Boolean)  
  Family inet enabled.
- **inet6** (Boolean)  
  Family inet6 enabled.
- **inet_address** (Block List)  
  List of `family inet` `address` and with each vrrp-group set.
  - **address** (String)  
    IPv4 address with mask.
  - **vrrp_group** (Block List)  
    See [below for nested schema](#vrrp_group-attributes-for-inet_address)
- **inet6_address** (Block List)  
  List of `family inet6` `address` and with each vrrp-group set.
  - **address** (String)  
    IPv6 address with mask.
  - **vrrp_group** (Block List)  
    See [below for nested schema](#vrrp_group-attributes-for-inet6_address)
- **inet_mtu** (Number)  
  MTU for family inet
- **inet6_mtu** (Number)  
  MTU for family inet6
- **inet_filter_input** (String)  
  Filter applied to received packets for family inet.
- **inet_filter_output** (String)  
  Filter applied to transmitted packets for family inet.
- **inet6_filter_input** (String)  
  Filter applied to received packets for family inet6.
- **inet6_filter_output** (String)  
  Filter applied to transmitted packets for family inet6.
- **inet_rpf_check** (Block)  
  Reverse-path-forwarding checks enabled and possible configuration for family inet.
  - **fail_filter** (String)  
    Name of filter applied to packets failing RPF check.
  - **mode_loose** (Boolean)  
    Reverse-path-forwarding use loose mode instead the strict mode.
- **inet6_rpf_check** (Block)  
  Reverse-path-forwarding checks enabled and possible configuration for family inet6.  
  Attributes is same as `inet_rpf_check`.
- **ether802_3ad** (String)  
  Link of 802.3ad interface.
- **trunk** (Boolean)  
  Interface mode is trunk.
- **vlan_members** (List of String)  
  List of vlan membership for this interface.
- **vlan_native** (Number)  
  Vlan for untagged frames.
- **ae_lacp** (String)  
  LACP option in aggregated-ether-options.
- **ae_link_speed** (String)  
  Link speed of individual interface that joins the AE.
- **ae_minimum_links** (Number)  
  Minimum number of aggregated links (1..8).
- **security_zone** (String)  
  Security zone where the interface is
- **routing_instance** (String)  
  Routing_instance where the interface is (if not default instance)

---

### vrrp_group attributes for inet_address

- **identifier** (Number)  
  ID for vrrp
- **virtual_address** (List of String)  
  List of address IP v4.
- **accept_data** (Boolean)  
  Accept packets destined for virtual IP address.
- **advertise_interval** (Number)  
  Advertisement interval (seconds)
- **advertisements_threshold** (Number)  
  Number of vrrp advertisements missed before declaring master down.
- **authentication_key** (String, Sensitive)  
  Authentication key
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
    Value to subtract from priority when interface is down
- **track_route** (Block List)  
  List of track_route.
  - **route** (String)  
    Route address tracked.
  - **routing_instance** (String)  
    Routing instance to which route belongs.
  - **priority_cost** (Number)  
    Value to subtract from priority when route is down.

---

### vrrp_group attributes for inet6_address

Same as [`vrrp_group` attributes for inet_address](#vrrp_group-attributes-for-inet_address) block
but without `authentication_key`, `authentication_type` and with

- **virtual_link_local_address** (String)  
  Address IPv6 for Virtual link-local addresses.
