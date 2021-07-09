---
layout: "junos"
page_title: "Junos: junos_bridge_domain"
sidebar_current: "docs-junos-resource-bridge-domain"
description: |-
  Create an bridge domain on root level or routing-instance (when Junos device supports it: MX, vMX)
---

# junos_bridge_domain

Provides an bridge domain on root level or routing-instance

## Example Usage

```hcl
# Add an bridge domain
resource "junos_bridge_domain" "demo" {
  name              = "demo"
  routing_interface = "irb.8"
  vlan_id           = 8
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of bridge domain.
* `routing_instance` - (Optional, Forces new resource)(`String`) Routing instance. Need to be 'default' (for root level) or the name of routing instance. Defaults to `default`.
* `community_vlans` - (Optional)(`ListOfString`) List of Community VLANs for private vlan bridge domain.
* `description` - (Optional)(`String`) Text description of bridge domain.
* `domain_id` - (Optional)(`Int`) Domain-id for auto derived Route Target (1..15).
* `domain_type_bridge` - (Optional)(`Bool`) Forwarding instance.
* `isolated_vlan` - (Optional)(`Int`) Isolated VLAN ID for private vlan bridge domain (1..4094).
* `routing_interface` - (Optional)(`String`) Routing interface name for this bridge-domain.
* `service_id` - (Optional)(`Int`) Service id required if bridge-domain is of type MC-AE and vlan-id all or vlan-id none or vlan-tags (1..65535).
* `vlan_id` - (Optional)(`Int`) IEEE 802.1q VLAN identifier for bridging domain (1..4094).
* `vlan_id_list` - (Optional)(`ListOfString`) Create bridge-domain for each of the vlan-id specified in the vlan-id-list.
* `vxlan` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare vxlan configuration.
  * `vni` - (Required)(`Int`) VXLAN identifier (0..16777214).
  * `decapsulate_accept_inner_vlan` - (Optional)(`Bool`) Accept VXLAN packets with inner VLAN.
  * `encapsulate_inner_vlan` - (Optional)(`Bool`) Retain inner VLAN in the packet.
  * `ingress_node_replication` - (Optional)(`Bool`) Enable ingress node replication.
  * `multicast_group` - (Optional)(`String`) CIDR for Multicast group registered for VXLAN segment.
  * `ovsdb_managed` - (Optional)(`Bool`) Bridge-domain is managed remotely via VXLAN OVSDB Controller.
  * `vni_extend_evpn` - (Optional)(`Bool`) Extend VNI to EVPN.
  * `unreachable_vtep_aging_timer` - (Optional)(`Int`) Unreachable VXLAN tunnel endpoint removal timer (300..1800 seconds).

## Import

Junos bridge domain can be imported using an id made up of `<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_bridge_domain.demo demo_-_default
```
