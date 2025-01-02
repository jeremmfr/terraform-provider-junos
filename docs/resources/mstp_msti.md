---
page_title: "Junos: junos_mstp_msti"
---

# junos_mstp_msti

Provides a MSTP MSTI resource.

## Example Usage

```hcl
resource "junos_mstp_msti" "instance1" {
  msti_id = 1
  vlan    = ["10-20"]
}
```

## Argument Reference

The following arguments are supported:

- **msti_id** (Required, Number, Forces new resource)  
  MSTI identifier (1..4094).
- **routing_instance** (Optional, String, Forces new resource)  
  Configure MSTP MSTI in routing instance.  
  Need to be `default` (for root level) or name of routing instance.  
  Defaults to `default`.
- **vlan** (Required, Set of String)  
  VLAN ID or VLAN ID range.
- **backup_bridge_priority** (Optional, String)  
  Priority of the bridge (in increments of 4k - 4k,8k,..60k).
- **bridge_priority** (Optional, String)  
  Priority of the bridge (in increments of 4k - 0,4k,8k,..60k).
- **interface** (Optional, Block Set)  
  For each interface, options.  
  At least one of the block arguments need to be set (in addition to `name`).
  - **name** (Required, String)  
    Interface name or `all`.
  - **cost** (Optional, Number)  
    Cost of the interface (1..200000000).
  - **priority** (Optional, Number)  
    Interface priority (in increments of 16 - 0,16,..240).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<msti_id>_-_<routing_instance>`.

## Import

Junos mstp msti can be imported using an id made up of `<msti_id>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_mstp_msti.instance1 1_-_default
```
