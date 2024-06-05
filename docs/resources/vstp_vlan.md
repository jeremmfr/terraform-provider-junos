---
page_title: "Junos: junos_vstp_vlan"
---

# junos_vstp_vlan

Provides a VSTP vlan resource.

## Example Usage

```hcl
resource "junos_vstp_vlan" "all" {
  vlan_id = "all"
}
```

## Argument Reference

The following arguments are supported:

- **vlan_id** (Required, String, Forces new resource)  
  VLAN id or `all`.
- **routing_instance** (Optional, String, Forces new resource)  
  Configure VSTP vlan in routing instance.  
  Need to be `default` (for root level) or name of routing instance.  
  Defaults to `default`.
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
  An identifier for the resource with format `<vlan_id>_-_<routing_instance>`.

## Import

Junos vstp vlan can be imported using an id made up of `<vlan_id>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_vstp_vlan.all all_-_default
```
