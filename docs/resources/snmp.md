---
page_title: "Junos: junos_snmp"
---

# junos_snmp

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `snmp` block.  
By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.

Configure static configuration in `snmp` block

## Example Usage

```hcl
# Configure snmp
resource "junos_snmp" "snmp" {
  filter_duplicates          = true
  filter_internal_interfaces = true
  health_monitor {}
  location = "Paris, France"
}
```

## Argument Reference

The following arguments are supported:

- **clean_on_destroy** (Optional, Boolean)  
  Clean supported lines when destroy this resource.
- **arp** (Optional, Boolean)  
  JVision ARP.
- **arp_host_name_resolution** (Optional, Boolean)  
  Enable host name resolution for JVision ARP.
- **contact** (Optional, String)  
  Contact information for administrator.
- **description** (Optional, String)  
  System description.
- **engine_id** (Optional, String)  
  SNMPv3 engine ID.  
  Need to be `use-default-ip-address`, `use-mac-address` or `local ...`
- **filter_duplicates** (Optional, Boolean)  
  Filter requests with duplicate source address/port and request ID.
- **filter_interfaces** (Optional, Set of String)  
  Regular expressions to list of interfaces that needs to be filtered.
- **filter_internal_interfaces** (Optional, Boolean)  
  Filter all internal interfaces.
- **health_monitor** (Optional, Block)  
  Enable `health-monitor`.
  - **falling_threshold** (Optional, Number)  
    Falling threshold applied to all monitored objects (0..100).
  - **idp** (Optional, Boolean)  
    Enable IDP health monitor.
  - **idp_falling_threshold** (Optional, Number)  
    Falling threshold applied to all idp monitored objects (0..100).
  - **idp_interval** (Optional, Number)  
    Interval between idp samples (1..2147483647).
  - **idp_rising_threshold** (Optional, Number)  
    Rising threshold applied to all monitored idp objects (0..100).
  - **interval** (Optional, Number)  
    Interval between samples (1..2147483647).
  - **rising_threshold** (Optional, Number)  
    Rising threshold applied to all monitored objects (1..100).
- **if_count_with_filter_interfaces** (Optional, Boolean)  
  Filter interfaces config for ifNumber and ipv6Interfaces.
- **interface** (Optional, Set of String)  
  Restrict SNMP requests to interfaces.
- **location** (Optional, String)  
  Physical location of system.
- **routing_instance_access** (Optional, Boolean)  
  Enable SNMP routing instance.
- **routing_instance_access_list** (Optional, Set of String)  
  Allow/Deny SNMP access to routing instances.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `snmp`.

## Import

Junos snmp can be imported using any id, e.g.

```shell
$ terraform import junos_snmp.snmp random
```
