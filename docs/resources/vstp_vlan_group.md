---
page_title: "Junos: junos_vstp_vlan_group"
---

# junos_vstp_vlan_group

Provides a VSTP vlan-group group resource.

## Example Usage

```hcl
resource "junos_vstp_vlan_group" "grp" {
  name = "grp"
  vlan = ["10", "11"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  VLAN group name.
- **routing_instance** (Optional, String, Forces new resource)  
  Configure VSTP vlan-group in routing instance.  
  Need to be `default` (for root level) or name of routing instance.  
  Defaults to `default`.
- **vlan** (Required, Set of String)  
  VLAN IDs or VLAN ID ranges (1..4094).
- **backup_bridge_priority** (Optional, String)  
  Priority of the bridge (in increments of 4k - 4k,8k,..60k).
- **bridge_priority** (Optional, String)  
  Priority of the bridge (in increments of 4k - 0,4k,8k,..60k).
- **forward_delay** (Optional, Number)  
  Time spent in listening or learning state (4..30 seconds).
- **hello_time** (Optional, Number)  
  Time interval between configuration BPDUs (1..10 seconds).
- **max_age** (Optional, Number)  
  Maximum age of received protocol bpdu (6..40 seconds).
- **system_identifier** (Optional, String)  
  System identifier to represent this node.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>`.

## Import

Junos vstp vlan-group can be imported using an id made up of `<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_vstp_vlan_group.grp grp_-_default
```
