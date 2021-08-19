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
  Name of unit interface (with dot).
- **description** (String)  
  Description for interface.
- **family_inet** (Block)  
  Family inet enabled and possible configuration.
  - **address** (Block List)  
    List of address.  
    See [below for nested schema](#address-attributes-for-family_inet).
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
    List of address.  
    See [below for nested schema](#address-attributes-for-family_inet6).
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
- **security_inbound_protocols** (List of String)  
  The inbound protocols allowed.
- **security_inbound_services** (List of String)  
  The inbound services allowed.
- **security_zone** (String)  
  Security zone where the interface is.
- **vlan_id** (Number)  
  802.1q VLAN ID for unit interface.

---

### address attributes for family_inet

- **cidr_ip** (String)  
  Address IP/Mask v4.
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

### address attributes for family_inet6

- **cidr_ip** (String)  
  Address IP/Mask v6.
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

### rpf_check attributes

- **fail_filter** (String)  
  Name of filter applied to packets failing RPF check.
- **mode_loose** (Boolean)  
  Reverse-path-forwarding use loose mode instead the strict mode.
