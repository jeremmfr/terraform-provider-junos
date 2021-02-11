---
layout: "junos"
page_title: "Junos: junos_interface_physical"
sidebar_current: "docs-junos-data-source-interface-physical"
description: |-
  Get information on a physical interface (as with an junos_interface_physcial resource import)
---

# junos_interface_physical

Get information on a physical interface.

## Example Usage

```hcl
# Search interface with name
data junos_interface_physical "interface_physical_demo" {
  config_interface = "ge-0/0/3"
}
```

## Argument Reference

The following arguments are supported:

* `config_interface` - (Optional)(`String`) Specifies the interface part for search. Command is 'show configuration interfaces <config_interface>'
* `match` - (Optional)(`String`) Regex string to filter lines and find only one interface.

~> **NOTE:** If more or less than a single match is returned by the search, Terraform will fail.

## Attributes Reference

* `id` - Like resource it's the `name` of interface
* `name` - Name of physical interface (without dot).
* `ae_lacp` - LACP option in aggregated-ether-options.
* `ae_link_speed` - Link speed of individual interface that joins the AE.
* `ae_minimum_links` - Minimum number of aggregated links (1..8).
* `description` - Description for interface.
* `ether802_3ad` - Link of 802.3ad interface.
* `trunk` - Interface mode is trunk.
* `vlan_members` - List of vlan membership for this interface.
* `vlan_native` - Vlan for untagged frames.
* `vlan_tagging` - 802.1q VLAN tagging support.
