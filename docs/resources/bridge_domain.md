---
page_title: "Junos: junos_bridge_domain"
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

-> **Note:** At least one of arguments need to be set
(in addition to `name` and `routing_instance`).

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Bridge domain name.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance.  
  Need to be `default` (for root level) or the name of routing instance.
  Defaults to `default`.
- **community_vlans** (Optional, Set of String)  
  List of Community VLANs for private vlan bridge domain.
- **description** (Optional, String)  
  Text description of bridge domain.
- **domain_id** (Optional, Number)  
  Domain-id for auto derived Route Target (1..15).
- **domain_type_bridge** (Optional, Boolean)  
  Forwarding instance.
- **interface** (Optional, Set of String)  
  Interface for this bridge domain.
- **isolated_vlan** (Optional, Number)  
  Isolated VLAN ID for private vlan bridge domain (1..4094).
- **routing_interface** (Optional, String)  
  Routing interface name for this bridge-domain.
- **service_id** (Optional, Number)  
  Service id required if bridge-domain is of type MC-AE and
  vlan-id all or vlan-id none or vlan-tags (1..65535).
- **vlan_id** (Optional, Number)  
  IEEE 802.1q VLAN identifier for bridging domain (1..4094).
- **vlan_id_list** (Optional, Set of String)  
  Create bridge-domain for each of the vlan-id specified in the vlan-id-list.
- **vxlan** (Optional, Block)  
  Declare vxlan options.
  - **vni** (Required, Number)  
    VXLAN identifier (0..16777214).
  - **vni_extend_evpn** (Optional, Boolean)  
    Extend VNI to EVPN.
  - **decapsulate_accept_inner_vlan** (Optional, Boolean)  
    Accept VXLAN packets with inner VLAN.
  - **encapsulate_inner_vlan** (Optional, Boolean)  
    Retain inner VLAN in the packet.
  - **ingress_node_replication** (Optional, Boolean)  
    Enable ingress node replication.
  - **multicast_group** (Optional, String)  
    CIDR for Multicast group registered for VXLAN segment.
  - **ovsdb_managed** (Optional, Boolean)  
    Bridge-domain is managed remotely via VXLAN OVSDB Controller.
  - **static_remote_vtep_list** (Optional, Set of String)  
    Configure bridge domain specific static remote VXLAN tunnel endpoints.
  - **unreachable_vtep_aging_timer** (Optional, Number)  
    Unreachable VXLAN tunnel endpoint removal timer (300..1800 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>`.

## Import

Junos bridge domain can be imported using an id made up of `<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_bridge_domain.demo demo_-_default
```
