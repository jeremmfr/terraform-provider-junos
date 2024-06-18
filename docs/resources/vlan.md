---
page_title: "Junos: junos_vlan"
---

# junos_vlan

Provides a vlan resource.

## Example Usage

```hcl
# Add a vlan
resource "junos_vlan" "blue" {
  name        = "blue"
  description = "blue-10"
  vlan_id     = 10
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of VLAN.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance if not root level.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`
- **community_vlans** (Optional, Set of String)  
  List of VLAN id or name of community vlans for primary vlan.
- **description** (Optional, String)  
  Text description of VLAN.
- **forward_filter_input** (Optional, String)  
  Input filter to apply for forwarded packets.
- **forward_filter_output** (Optional, String)  
  Output filter to apply for forwarded packets.
- **forward_flood_input** (Optional, String)  
  Input filter to apply for ethernet switching flood packets.
- **isolated_vlan** (Optional, String)  
  VLAN id or name of isolated vlan for primary vlan.
- **l3_interface** (Optional, String)  
  L3 interface name for this VLAN.  
  Must be start with `irb.` or `vlan.`.
- **no_arp_suppression** (Optional, Boolean)  
  Turn off ARP suppression.
- **private_vlan** (Optional, String)  
  Type of secondary VLAN for private vlan.  
  Must be `community` or `isolated`.
- **service_id** (Optional, Number)  
  Service id.
- **vlan_id** (Optional, String)  
  802.1q VLAN id or `all` or `none`.  
  Conflict with `vlan_id_list`.
- **vlan_id_list** (Optional, Set of String)  
  List of 802.1q VLAN id.  
  Can be a VLAN id or range of VLAN id (example: 10-20).  
  Conflict with `vlan_id`.
- **vxlan** (Optional, Block)  
  Declare vxlan configuration.
  - **vni** (Required, Number)  
    VXLAN identifier (0..16777214).
  - **vni_extend_evpn** (Optional, Boolean)  
    Extend VNI to EVPN.
  - **encapsulate_inner_vlan** (Optional, Boolean)  
    Retain inner VLAN in the packet.
  - **ingress_node_replication** (Optional, Boolean)  
    Enable ingress node replication.
  - **multicast_group** (Optional, String)  
    Multicast group registered for VXLAN segment.
  - **ovsdb_managed** (Optional, Boolean)  
    Bridge-domain is managed remotely via VXLAN OVSDB Controller.
  - **static_remote_vtep_list** (Optional, Set of String)  
    Configure vlan specific static remote VXLAN tunnel endpoints.
  - **translation_vni** (Optional, Number)  
    Translated VXLAN identifier (1..16777214).
  - **unreachable_vtep_aging_timer** (Optional, Number)  
    Unreachable VXLAN tunnel endpoint removal timer (300..1800 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>`.

## Import

Junos vlan can be imported using an id made up of
`<name>` or `<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_vlan.blue blue
```
