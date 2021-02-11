---
layout: "junos"
page_title: "Junos: junos_vlan"
sidebar_current: "docs-junos-resource-vlan"
description: |-
  Create a vlan (when Junos device supports it)
---

# junos_vlan

Provides a vlan resource.

## Example Usage

```hcl
# Add a vlan
resource junos_vlan "blue" {
  name        = "blue"
  description = "blue-10"
  vlan_id     = 10
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of vlan.
* `community_vlans` - (Optional)(`ListOfInt`) List of ID community vlan for primary vlan (when Junos device supports it).
* `description` - (Optional)(`String`) A description for vlan.
* `forward_filter_input` - (Optional)(`String`) input filter to apply for forwarded packets (when Junos device supports it).
* `forward_filter_output` - (Optional)(`String`) output filter to apply for forwarded packets (when Junos device supports it).
* `forward_flood_input` - (Optional)(`String`) input filter to apply for ethernet switching flood packets (when Junos device supports it).
* `l3_interface` - (Optional)(`String`) L3 interface name for this vlans. Must be start with irb.
* `isolated-vlan` - (Optional)(`Int`) declare ID isolated vlan for primary vlan (when Junos device supports it).
* `private_vlan` - (Optional)(`String`) Type of secondary vlan for private vlan. Must be 'community' or 'isolated' (when Junos device supports it).
* `service_id` - (Optional)(`Int`) Service id (when Junos device supports it).
* `vlan_id` - (Optional)(`Int`) 802.1q VLAN identifier. Conflict with `vlan_id_list`.
* `vlan_id_list` - (Optional)(`ListOfString`) List of vlan ID. Can be a ID or range (exemple: 10-20). Conflict with `vlan_id`.
* `vxlan` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare vxlan configuration (when Junos device supports it).
  * `vni` - (Required)(`Int`) VXLAN identifier.
  * `encapsulate_inner_vlan` - (Optional)(`Bool`) Retain inner VLAN in the packet.
  * `ingress_node_replication` - (Optional)(`Bool`) Enable ingress node replication.
  * `multicast_group` - (Optional)(`String`) CDIR for Multicast group registered for VXLAN segment.
  * `ovsdb_managed` - (Optional)(`Bool`) Bridge-domain is managed remotely via VXLAN OVSDB Controller.
  * `unreachable_vtep_aging_timer` - (Optional)(`Int`) Unreachable VXLAN tunnel endpoint removal timer.



## Import

Junos vlan can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_vlan.blue blue
```
