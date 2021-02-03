---
layout: "junos"
page_title: "Junos: junos_interface_physical"
sidebar_current: "docs-junos-resource-interface-physical"
description: |-
  Create/configure a physical interface
---

# junos_interface_physical

Provides a physical interface resource.

## Example Usage

```hcl
# Configure interface of switch
resource junos_interface_physical "interface_switch_demo" {
  name         = "ge-0/0/0"
  description  = "interfaceSwitchDemo"
  trunk        = true
  vlan_members = ["100"]
}
# Configure interface for L3 logical interface on Junos Router or firewall
resource junos_interface_physical "interface_fw_demo" {
  name         = "ge-0/0/1"
  description  = "interfaceFwDemo"
  vlan_tagging = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Name of physical interface (without dot).
* `no_disable_on_destroy` - (Optional)(`Bool`) When destroy this resource, delete all configurations => do not add `disable` + `descrition NC` or `apply-groups` with `group_interface_delete` provider argument on interface.
* `ae_lacp` - (Optional)(`String`) Add lacp option in aggregated-ether-options. Need to be 'active' or 'passive' for initiate transmission or respond.
* `ae_link_speed` - (Optional)(`String`) Link speed of individual interface that joins the AE.
* `ae_minimum_links` - (Optional)(`Int`) Minimum number of aggregated links (1..8).
* `description` - (Optional)(`String`) Description for interface.
* `ether802_3ad` - (Optional)(`String`) Name of aggregated device for add this interface to link of 802.3ad interface.
* `trunk` - (Optional)(`Bool`) Interface mode is trunk.
* `vlan_members` - (Optional)(`ListOfString`) List of vlan for membership for this interface.
* `vlan_native` - (Optional)(`Int`) Vlan for untagged frames.
* `vlan_tagging` - (Optional)(`Bool`) Add 802.1q VLAN tagging support.

## Import

Junos interface can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_interface_physical.interface_switch_demo ge-0/0/0
$ terraform import junos_interface_physical.interface_fw_demo_100 ge-0/0/1
```
