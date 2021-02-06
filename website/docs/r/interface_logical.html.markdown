---
layout: "junos"
page_title: "Junos: junos_interface_logical"
sidebar_current: "docs-junos-resource-interface-logical"
description: |-
  Create/configure a logical interface
---

# junos_interface_logical

Provides a logical interface resource.

## Example Usage

```hcl
resource junos_interface_logical "interface_fw_demo_100" {
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

* `name` - (Required, Forces new resource)(`String`) Name of unit interface (with dot).
* `st0_also_on_destroy` - (Optional)(`Bool`) When destroy this resource, if the name has prefix 'st0.', delete all configurations (not keep empty st0 interface).  
* `description` - (Optional)(`String`) Description for interface.
(Usually, `st0.x` interfaces are completely deleted with `bind_interface_auto` argument in `junos_security_ipsec_vpn` resource or by `junos_interface_st0_unit` resource because of the dependency, but only if st0.x interface is empty or disable.)
* `family_inet` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Enable family inet and add configurations if specified.
  * `address` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each ip address to declare. See the [`address` arguments for family_inet](#address-arguments-for-family_inet) block.
  * `filter_input` - (Optional)(`String`) Filter to be applied to received packets.
  * `filter_output` - (Optional)(`String`) Filter to be applied to transmitted packets.
  * `mtu` - (Optional)(`Int`) Maximum transmission unit.
  * `rpf_check` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for enable reverse-path-forwarding checks on this interface. See the [`rpf_check` arguments](#rpf_check-arguments) block for optional arguments.
* `family_inet6` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Enable family inet6 and add configurations if specified.
  * `address` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each ipv6 address to declare. See the [`address` arguments for family_inet6](#address-arguments-for-family_inet6) block.
  * `filter_input` - (Optional)(`String`) Filter to be applied to received packets.
  * `filter_output` - (Optional)(`String`) Filter to be applied to transmitted packets.
  * `mtu` - (Optional)(`Int`) Maximum transmission unit.
  * `rpf_check` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for enable reverse-path-forwarding checks on this interface. See the [`rpf_check` arguments](#rpf_check-arguments) block for optional arguments. 
* `routing_instance` - (Optional)(`String`) Add this interface in routing_instance. Need to be created before.
* `security_zone` - (Optional)(`String`) Add this interface in security_zone. Need to be created before.
* `vlan_id` - (Optional,Computed)(`Int`) 802.1q VLAN ID for unit interface. If not set, computed with `name` of interface (ge-0/0/0.100 = 100) except if name has '.0' suffix or 'st0.' prefix.

---
#### address arguments for family_inet
* `cidr_ip` - (Required)(`String`) Address IP/Mask v4.
* `vrrp_group` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each vrrp group to declare. See the [`vrrp_group` arguments for address in family_inet](#vrrp_group-arguments-for-address-in-family_inet) block.

---
#### vrrp_group arguments for address in family_inet
* `identifier` - (Required)(`Int`) ID for vrrp
* `virtual_address` - (Required)(`ListOfString`) List of address IP v4.
* `accept_data` - (Optional)(`Bool`) Accept packets destined for virtual IP address. Conflict with `no_accept_data` when apply.
* `advertise_interval` - (Optional)(`Int`) Advertisement interval (seconds)
* `advertisements_threshold` - (Optional)(`Int`)  Number of vrrp advertisements missed before declaring master down.
* `authentication_key` - (Optional)(`String`) Authentication key
* `authentication_type` - (Optional)(`String`) Authentication type. Need to be 'md5' or 'simple'.
* `no_accept_data` - (Optional)(`Bool`) Don't accept packets destined for virtual IP address. Conflict with `accept_data` when apply.
* `no_preempt` - (Optional)(`Bool`) Don't allow preemption. Conflict with `preempt` when apply.
* `preempt` - (Optional)(`Bool`) Allow preemption. Conflict with `no_preempt` when apply.
* `priority` - (Optional)(`Int`) Virtual router election priority.
* `track_interface` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each track_interface to declare.
  * `interface` - (Required)(`String`) Name of interface.
  * `priority_cost` - (Required)(`Int`) Value to subtract from priority when interface is down
* `track_route` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each track_route to declare.
  * `route` - (Required)(`String`) Route address.
  * `routing_instance` - (Required)(`String`) Routing instance to which route belongs, or 'default'.
  * `priority_cost` - (Required)(`Int`) Value to subtract from priority when route is down.

---
#### address arguments for family_inet6
* `cidr_ip` - (Required)(`String`) Address IP/Mask v6.
* `vrrp_group` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each vrrp group to declare. See the [`vrrp_group` arguments for address in family_inet6](#vrrp_group-arguments-for-address-in-family_inet6) block.

---
#### vrrp_group arguments for address in family_inet6
Same as [`vrrp_group` arguments for address in family_inet](#vrrp_group-arguments-for-address-in-family_inet) block but without `authentication_key`, `authentication_type` and with
* `virtual_link_local_address` - (Required)(`String`) Address IPv6 for Virtual link-local addresses.

---
#### rpf_check arguments
* `fail_filter` - (Optional)(`String`) Name of filter applied to packets failing RPF check.
* `mode_loose` - (Optional)(`Bool`) Use reverse-path-forwarding loose mode instead the strict mode.

## Import

Junos interface can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_interface_logical.interface_fw_demo_100 ae.100
```
