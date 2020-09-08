---
layout: "junos"
page_title: "Junos: junos_interface"
sidebar_current: "docs-junos-data-source-interface"
description: |-
  Get information on an Interface (as with an junos_interface resource import)
---

# junos_interface

Get information on an Interface

## Example Usage

```hcl
# Search interface with IP
data junos_interface "demo_ip" {
  match = "192.0.2.2/"
}
# Search interface with name
data junos_interface "interface_fw_demo" {
  config_interface         = "ge-0/0/3.0"
}
```

## Argument Reference

The following arguments are supported:

* `config_interface` - (Optional)(`String`) Specifies the interface part for search. Command is 'show configuration interfaces <config_interface>'
* `match` - (Optional)(`String`) Regex string to filter lines and find only one interface.

~> **NOTE:** If more or less than a single match is returned by the search, Terraform will fail.

## Attributes Reference

* `id` - Like resource it's the `name` of interface
* `name` - Name of interface or unit interface (with dot).
* `description` - Description for interface.
* `vlan_tagging` - 802.1q VLAN tagging support.
* `vlan_taggind_id` - 802.1q VLAN ID for unit interface.
* `inet` - Family inet enabled.
* `inet6` - Family inet6 enabled.
* `inet_address` - List of `family inet` `address` and with each vrrp-group set.
  * `inet_address.#.address` - IPv4 address with mask.
  * `inet_address.#.vrrp_group` - See [`vrrp_group` attributes for inet_address](#vrrp_group-attributes-for-inet_address)
* `inet6_address` - List of `family inet6` `address` and with each vrrp-group set.
  * `inet6_address.#.address` - IPv6 address with mask.
  * `inet6_address.#.vrrp_group` -  See [`vrrp_group` attributes for inet6_address](#vrrp_group-attributes-for-inet6_address)
* `inet_mtu` - MTU for family inet
* `inet6_mtu` - MTU for family inet6
* `inet_filter_input` - Filter applied to received packets for family inet.
* `inet_filter_output` - Filter applied to transmitted packets for family inet.
* `inet6_filter_input` - Filter applied to received packets for family inet6.
* `inet6_filter_output` - Filter applied to transmitted packets for family inet6.
* `ether802_3ad` - Link of 802.3ad interface.
* `trunk` - Interface mode is trunk.
* `vlan_members` - List of vlan membership for this interface.
* `vlan_native` - Vlan for untagged frames.
* `ae_lacp` - LACP option in aggregated-ether-options.
* `ae_link_speed` - Link speed of individual interface that joins the AE.
* `ae_minimum_links` - Minimum number of aggregated links (1..8).
* `security_zone` - Security zone where the interface is
* `routing_instance` - Routing_instance where the interface is (if not default instance)

#### vrrp_group attributes for inet_address
* `identifier` - ID for vrrp
* `virtual_address` - List of address IP v4.
* `accept_data` - Accept packets destined for virtual IP address.
* `advertise_interval` - Advertisement interval (seconds)
* `advertisements_threshold` - Number of vrrp advertisements missed before declaring master down.
* `authentication_key` - Authentication key
* `authentication_type` - Authentication type.
* `no_accept_data` - Don't accept packets destined for virtual IP address.
* `no_preempt` - Preemption not allowed.
* `preempt` - Preemption allowed.
* `priority` - Virtual router election priority.
* `track_interface` - List of track_interface.
  * `interface` - Interface tracked.
  * `priority_cost` - Value to subtract from priority when interface is down
* `track_route` - List of track_route.
  * `route` - Route address tracked.
  * `routing_instance` - Routing instance to which route belongs.
  * `priority_cost` - Value to subtract from priority when route is down.

#### vrrp_group attributes for inet6_address
Same as [`vrrp_group` attributes for inet_address](#vrrp_group-attributes-for-inet_address) block but without `authentication_key`, `authentication_type` and with

 * `virtual_link_local_address` - Address IPv6 for Virtual link-local addresses.
