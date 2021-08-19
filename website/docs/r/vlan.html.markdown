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

- **name** (Required, String, Forces new resource)  
  The name of vlan.
- **community_vlans** (Optional, Set of Number)  
  List of ID community vlan for primary vlan (when Junos device supports it).
- **description** (Optional, String)  
  A description for vlan.
- **forward_filter_input** (Optional, String)  
  Input filter to apply for forwarded packets (when Junos device supports it).
- **forward_filter_output** (Optional, String)  
  Output filter to apply for forwarded packets (when Junos device supports it).
- **forward_flood_input** (Optional, String)  
  Input filter to apply for ethernet switching flood packets (when Junos device supports it).
- **l3_interface** (Optional, String)  
  L3 interface name for this vlans.  
  Must be start with `irb.` or `vlan.`.
- **isolated-vlan** (Optional, Number)  
  Declare ID isolated vlan for primary vlan (when Junos device supports it).
- **private_vlan** (Optional, String)  
  Type of secondary vlan for private vlan (when Junos device supports it).  
  Must be `community` or `isolated`.
- **service_id** (Optional, Number)  
  Service id (when Junos device supports it).
- **vlan_id** (Optional, Number)  
  802.1q VLAN identifier.  
  Conflict with `vlan_id_list`.
- **vlan_id_list** (Optional, Set of String)  
  List of vlan ID.  
  Can be a ID or range (exemple: 10-20).  
  Conflict with `vlan_id`.
- **vxlan** (Optional, Block)  
  Declare vxlan configuration (when Junos device supports it).
  - **vni** (Required, Number)  
    VXLAN identifier (0..16777214).
  - **encapsulate_inner_vlan** (Optional, Boolean)  
    Retain inner VLAN in the packet.
  - **ingress_node_replication** (Optional, Boolean)  
    Enable ingress node replication.
  - **multicast_group** (Optional, String)  
    CIDR for Multicast group registered for VXLAN segment.
  - **ovsdb_managed** (Optional, Boolean)  
    Bridge-domain is managed remotely via VXLAN OVSDB Controller.
  - **vni_extend_evpn** (Optional, Boolean)  
    Extend VNI to EVPN.
  - **unreachable_vtep_aging_timer** (Optional, Number)  
    Unreachable VXLAN tunnel endpoint removal timer (300..1800 seconds).

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos vlan can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_vlan.blue blue
```
