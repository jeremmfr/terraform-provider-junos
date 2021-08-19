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

- **name** (Required, String, Forces new resource)  
  Name of unit interface (with dot).
- **st0_also_on_destroy** (Optional, Boolean)  
  When destroy this resource, if the name has prefix `st0.`,
  delete all configurations (not keep empty st0 interface).  
  Usually, `st0.x` interfaces are completely deleted with `bind_interface_auto` argument in
  `junos_security_ipsec_vpn` resource or by `junos_interface_st0_unit` resource because of
  the dependency, but only if st0.x interface is empty or disable.
- **description** (Optional, String)  
  Description for interface.
- **family_inet** (Optional, Block)  
  Enable family inet and add configurations if specified.
  - **address** (Optional, Block List)  
    For each ip address to declare.  
    See [below for nested schema](#address-arguments-for-family_inet).
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
    For each ipv6 address to declare.  
    See [below for nested schema](#address-arguments-for-family_inet6).
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
- **security_inbound_protocols** (Optional, List of String)  
  The inbound protocols allowed.  
  Must be a list of Junos protocols.  
  `security_zone` need to be set.
- **security_inbound_services** (Optional, List of String)  
  The inbound services allowed.  
  Must be a list of Junos services.  
  `security_zone` need to be set.
- **security_zone** (Optional, String)  
  Add this interface in security_zone.  
  Need to be created before.
- **vlan_id** (Optional, Computed, Number)  
  802.1q VLAN ID for unit interface.  
  If not set, computed with `name` of interface (ge-0/0/0.100 = 100)
  except if name has `.0` suffix or `st0.`, `irb.`, `vlan.` prefix.
- **vlan_no_compute** (Optional, Boolean)  
  Disable the automatic compute of the `vlan_id` argument when not set.  
  Unnecessary if name has `.0` suffix or `st0.`, `irb.`, `vlan.` prefix because it's already disabled.

---

### address arguments for family_inet

- **cidr_ip** (Required, String)  
  Address IP/Mask v4.
- **preferred** (Optional, Boolean)  
  Preferred address on interface.
- **primary** (Optional, Boolean)  
  Candidate for primary address in system.
- **vrrp_group** (Optional, Block List)  
  For each vrrp group to declare.  
  See [below for nested schema](#vrrp_group-arguments-for-address-in-family_inet).

---

### vrrp_group arguments for address in family_inet

- **identifier** (Required, Number)  
  ID for vrrp
- **virtual_address** (Required, List of String)  
  List of address IP v4.
- **accept_data** (Optional, Boolean)  
  Accept packets destined for virtual IP address.  
  Conflict with `no_accept_data` when apply.
- **advertise_interval** (Optional, Number)  
  Advertisement interval (seconds)
- **advertisements_threshold** (Optional, Number)  
   Number of vrrp advertisements missed before declaring master down.
- **authentication_key** (Optional, String, Sensitive)  
  Authentication key.
- **authentication_type** (Optional, String)  
  Authentication type.  
  Need to be `md5` or `simple`.
- **no_accept_data** (Optional, Boolean)  
  Don't accept packets destined for virtual IP address.  
  Conflict with `accept_data` when apply.
- **no_preempt** (Optional, Boolean)  
  Don't allow preemption.  
  Conflict with `preempt` when apply.
- **preempt** (Optional, Boolean)  
  Allow preemption.  
  Conflict with `no_preempt` when apply.
- **priority** (Optional, Number)  
  Virtual router election priority.
- **track_interface** (Optional, Block List)  
  For each interface to declare.
  - **interface** (Required, String)  
    Name of interface.
  - **priority_cost** (Required, Number)  
    Value to subtract from priority when interface is down
- **track_route** (Optional, Block List)  
  For each route to declare.
  - **route** (Required, String)  
    Route address.
  - **routing_instance** (Required, String)  
    Routing instance to which route belongs, or `default`.
  - **priority_cost** (Required, Number)  
    Value to subtract from priority when route is down.

---

### address arguments for family_inet6

- **cidr_ip** (Required, String)  
  Address IP/Mask v6.
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

### rpf_check arguments

- **fail_filter** (Optional, String)  
  Name of filter applied to packets failing RPF check.
- **mode_loose** (Optional, Boolean)  
  Use reverse-path-forwarding loose mode instead the strict mode.

## Import

Junos interface can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_interface_logical.interface_fw_demo_100 ae.100
```
