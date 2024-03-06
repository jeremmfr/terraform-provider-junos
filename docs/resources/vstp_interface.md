---
page_title: "Junos: junos_vstp_interface"
---

# junos_vstp_interface

Provides a VSTP interface resource.

## Example Usage

```hcl
resource "junos_vstp_interface" "all" {
  name = "all"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Interface name or `all`.
- **routing_instance** (Optional, String, Forces new resource)  
  Configure VSTP interface in routing instance.  
  Need to be `default` (for root level) or name of routing instance.  
  Defaults to `default`.
- **vlan** (Optional, String)  
  Configure interface in VSTP vlan.  
  Conflict with `vlan_group`.
- **vlan_group** (Optional, String)  
  Configure interface in VSTP vlan-group.  
  Conflict with `vlan`.
- **access_trunk** (Optional, Boolean)  
  Send/Receive untagged RSTP BPDUs on this interface.
- **bpdu_timeout_action_alarm** (Optional, Boolean)  
  Generate an alarm on BPDU expiry (Loop Protect).
- **bpdu_timeout_action_block** (Optional, Boolean)  
  Block the interface on BPDU expiry (Loop Protect).
- **cost** (Optional, Number)  
  Cost of the interface (1..200000000).
- **edge** (Optional, Boolean)  
  Port is an edge port.
- **mode** (Optional, String)  
  Interface mode (P2P or shared).  
  Need to be `point-to-point` or `shared`.
- **no_root_port** (Optional, Boolean)  
  Do not allow the interface to become root (Root Protect).
- **priority** (Optional, Number)  
  Interface priority (in increments of 16 - 0,16,..240).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format  
  `<name>_-__-_<routing_instance>`,  
  `<name>_-_v_<vlan>_-_<routing_instance>`  
  or `<name>_-_vg_<vlan_group>_-_<routing_instance>`.

## Import

Junos vstp interface can be imported using an id made up of  
`<name>_-__-_<routing_instance>`,  
`<name>_-_v_<vlan>_-_<routing_instance>`  
or `<name>_-_vg_<vlan_group>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_vstp_interface.all all_-__-_default
```
