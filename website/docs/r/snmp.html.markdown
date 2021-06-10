---
layout: "junos"
page_title: "Junos: junos_snmp"
sidebar_current: "docs-junos-resource-snmp"
description: |-
  Configure static configuration in snmp block
---

# junos_snmp

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `snmp` block. By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.  

Configure static configuration in `snmp` block

## Example Usage

```hcl
# Configure snmp
resource junos_snmp "snmp" {
  filter_duplicates          = true
  filter_internal_interfaces = true
  health_monitor {}
  location = "Paris, France"
}
```

## Argument Reference

The following arguments are supported:

* `clean_on_destroy` - (Optional)(`Bool`) Clean supported lines when destroy this resource.
* `arp` - (Optional)(`Bool`) JVision ARP.
* `arp_host_name_resolution` - (Optional)(`Bool`) Enable host name resolution for JVision ARP.
* `contact` - (Optional)(`String`) Contact information for administrator.
* `description` - (Optional)(`String`) System description.
* `filter_duplicates` - (Optional)(`Bool`) Filter requests with duplicate source address/port and request ID.
* `filter_interfaces` - (Optional)(`ListOfString`) List of interfaces that needs to be filtered.
* `filter_internal_interfaces` - (Optional)(`Bool`) Filter all internal interfaces.
* `health_monitor` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable 'health-monitor'.
  * `falling_threshold` - (Optional)(`Int`) Falling threshold applied to all monitored objects (0..100).
  * `idp` - (Optional)(`Bool`) Enable IDP health monitor.
  * `idp_falling_threshold` - (Optional)(`Int`) Falling threshold applied to all idp monitored objects (0..100).
  * `idp_interval` - (Optional)(`Int`) Interval between idp samples (1..2147483647).
  * `idp_rising_threshold` - (Optional)(`Int`) Rising threshold applied to all monitored idp objects(0..100).
  * `interval` - (Optional)(`Int`) Interval between samples (1..2147483647).
  * `rising_threshold` - (Optional)(`Int`) Rising threshold applied to all monitored objects (1..100).
* `if_count_with_filter_interfaces` - (Optional)(`Bool`) Filter interfaces config for ifNumber and ipv6Interfaces.
* `interface` - (Optional)(`ListOfString`) Restrict SNMP requests to interfaces.
* `location` - (Optional)(`String`) Physical location of system.
* `routing_instance_access` - (Optional)(`Bool`) Enable SNMP routing-instance.
* `routing_instance_access_list` - (Optional)(`ListOfString`) Allow/Deny SNMP access to routing-instances.

## Import

Junos snmp can be imported using any id, e.g.

```
$ terraform import junos_snmp.snmp random
```
