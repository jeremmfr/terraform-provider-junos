---
page_title: "Junos: junos_rstp_interface"
---

# junos_rstp_interface

Provides a RSTP interface resource.

## Example Usage

```hcl
resource "junos_rstp_interface" "all" {
  name = "all"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Interface name or `all`.
- **routing_instance** (Optional, String, Forces new resource)  
  Configure RSTP interface in routing instance.  
  Need to be `default` (for root level) or name of routing instance.  
  Defaults to `default`.
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
  An identifier for the resource with format `<name>_-_<routing_instance>`.

## Import

Junos rstp interface can be imported using an id made up of `<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_rstp_interface.all all_-_default
```
