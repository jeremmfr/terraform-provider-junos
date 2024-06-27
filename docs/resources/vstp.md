---
page_title: "Junos: junos_vstp"
---

# junos_vstp

-> **Note:** This resource should only be created **once** for root level or each
routing-instance. It's used to configure static (not object) options in `protocols vstp` block
in root or routing-instance level.

Configure static configuration in `protocols vstp` block for root or routing-instance level.

## Example Usage

```hcl
# Configure vstp
resource "junos_vstp" "vstp" {
  bpdu_block_on_edge = true
}
```

## Argument Reference

The following arguments are supported:

- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance.  
  Need to be `default` (for root level) or the name of routing instance.  
  Defaults to `default`.
- **bpdu_block_on_edge** (Optional, Boolean)  
  Block BPDU on all interfaces configured as edge (BPDU Protect).
- **disable** (Optional, Boolean)  
  Disable STP.
- **force_version_stp** (Optional, Boolean)  
  Force protocol version STP.
- **priority_hold_time** (Optional, Number)  
  Hold time before switching to primary priority when core domain becomes up (1..255 seconds).
- **system_id** (Optional, Block Set)  
  For each ID, System ID to IP mapping.
  - **id** (Required, String)  
    System ID.  
    Format need to be `<mac-address>`
  - **ip_address** (Optional, String)  
    Peer ID (IP Address).
- **vpls_flush_on_topology_change** (Optional, Boolean)  
  Enable VPLS MAC flush on root protected CE interface receiving topology change.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<routing_instance>`.

## Import

Junos vstp can be imported using an id made up of `<routing_instance>`, e.g.

```shell
$ terraform import junos_vstp.vstp default
```
