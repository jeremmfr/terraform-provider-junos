---
layout: "junos"
page_title: "Junos: junos_interface_logical"
sidebar_current: "docs-junos-data-source-interface-logical"
description: |-
  Get information on a logical interface (as with an junos_interface_logical resource import)
---

# junos_interface_logical

Get information on a logical interface

## Example Usage

```hcl
# Search interface with IP
data junos_interface_logical "demo_ip" {
  match = "192.0.2.2/"
}
# Search interface with name
data junos_interface_logical "interface_fw_demo" {
  config_interface = "ge-0/0/3.0"
}
```

## Argument Reference

The following arguments are supported:

* `config_interface` - (Optional)(`String`) Specifies the interface part for search. Command is 'show configuration interfaces <config_interface>'
* `match` - (Optional)(`String`) Regex string to filter lines and find only one interface.

~> **NOTE:** If more or less than a single match is returned by the search, Terraform will fail.

## Attributes Reference

* `id` - Like resource it's the `name` of interface
* `name` - Name of unit interface (with dot).
* `description` - Description for interface.
* `family_inet` - Family inet enabled and possible configuration.
  * `address` - List of address. See the [`address` attributes for family_inet](#address-attributes-for-family_inet) block.
  * `filter_input` - Filter applied to received packets.
  * `filter_output` - Filter applied to transmitted packets.
  * `mtu` - Maximum transmission unit.
  * `rpf_check` - Reverse-path-forwarding checks enabled and possible configuration. See the [`rpf_check` attributes](#rpf_check-attributes) block for attributes.
* `family_inet6` - Family inet6 enabled and possible configuration.
  * `address` - List of address. See the [`address` attributes for family_inet6](#address-attributes-for-family_inet6) block.
  * `filter_input` - Filter applied to received packets.
  * `filter_output` - Filter applied to transmitted packets.
  * `mtu` - Maximum transmission unit.
  * `rpf_check` - Reverse-path-forwarding checks enabled and possible configuration. See the [`rpf_check` attributes](#rpf_check-attributes) block for attributes.
* `routing_instance` - Routing_instance where the interface is (if not default instance).
* `security_zone` - Security zone where the interface is.
* `vlan_id` - 802.1q VLAN ID for unit interface.

---
#### address attributes for family_inet
* `cidr_ip` - Address IP/Mask v4.
* `vrrp_group` - List of vrrp group configurations. See the [`vrrp_group` attributes for address in family_inet](#vrrp_group-attributes-for-address-in-family_inet) block.

---
#### vrrp_group attributes for address in family_inet
* `identifier` - ID for vrrp.
* `virtual_address` - List of address IP v4.
* `accept_data` - Accept packets destined for virtual IP address.
* `advertise_interval` - Advertisement interval (seconds).
* `advertisements_threshold` - Number of vrrp advertisements missed before declaring master down.
* `authentication_key` - Authentication key.
* `authentication_type` - Authentication type.
* `no_accept_data` - Don't accept packets destined for virtual IP address.
* `no_preempt` - Preemption not allowed.
* `preempt` - Preemption allowed.
* `priority` - Virtual router election priority.
* `track_interface` - List of track_interface.
  * `interface` - Interface tracked.
  * `priority_cost` - Value to subtract from priority when interface is down.
* `track_route` - List of track_route.
  * `route` - Route address tracked.
  * `routing_instance` - Routing instance to which route belongs.
  * `priority_cost` - Value to subtract from priority when route is down.

---
#### address attributes for family_inet6
* `cidr_ip` - Address IP/Mask v6.
* `vrrp_group` - List of vrrp group configurations. See the [`vrrp_group` attributes for address in family_inet6](#vrrp_group-attributes-for-address-in-family_inet6) block.

---
#### vrrp_group attributes for address in family_inet6
Same as [`vrrp_group` attributes for address in family_inet](#vrrp_group-attributes-for-address-in-family_inet) block but without `authentication_key`, `authentication_type` and with  
* `virtual_link_local_address` - Address IPv6 for Virtual link-local addresses.

---
### rpf_check attributes
* `fail_filter` - Name of filter applied to packets failing RPF check.
* `mode_loose` - Reverse-path-forwarding use loose mode instead the strict mode.
