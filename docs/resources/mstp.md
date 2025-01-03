---
page_title: "Junos: junos_mstp"
---

# junos_mstp

~> **Note**
  This resource should only be created **once** for root level or each routing-instance.  
  It's used to configure static (not object) options in `protocols mstp` block
  in root or routing-instance level.

Configure static configuration in `protocols mstp` block for root or routing-instance level.

## Example Usage

```hcl
# Configure mstp
resource "junos_mstp" "mstp" {
  configuration_name = "abc"
  revision_level     = 1
  bridge_priority    = "32k"
}
```

## Argument Reference

The following arguments are supported:

- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance.  
  Need to be `default` (for root level) or the name of routing instance.  
  Defaults to `default`.
- **backup_bridge_priority** (Optional, String)  
  Priority of the bridge (in increments of 4k - 4k,8k,..60k).
- **bpdu_block_on_edge** (Optional, Boolean)  
  Block BPDU on all interfaces configured as edge (BPDU Protect).
- **bpdu_destination_mac_address_provider_bridge_group** (Optional, Boolean)  
  Destination MAC address in the spanning tree BPDUs is 802.1ad provider bridge group address.
- **bridge_priority** (Optional, String)  
  Priority of the bridge (in increments of 4k - 0,4k,8k,..60k).
- **configuration_name** (Optional, String)  
  Configuration name (part of MST configuration identifier).
- **disable** (Optional, Boolean)  
  Disable MSTP.
- **forward_delay** (Optional, Number)  
  Time spent in listening or learning state (4..30 seconds).
- **hello_time** (Optional, Number)  
  Time interval between configuration BPDUs (1..10 seconds).
- **max_age** (Optional, Number)  
  Maximum age of received protocol bpdu (6..40 seconds).
- **max_hops** (Optional, Number)  
  Maximum number of hops (1..255).
- **priority_hold_time** (Optional, Number)  
  Hold time before switching to primary priority when core domain becomes up (1..255 seconds).
- **revision_level** (Optional, Number)  
  Revision level (part of MST configuration identifier).
- **system_id** (Optional, Block Set)  
  For each ID, System ID to IP mapping.
  - **id** (Required, String)  
    System ID.  
    Format need to be `<mac-address>`
  - **ip_address** (Optional, String)  
    Peer ID (IP Address).
- **system_identifier** (Optional, String)  
  System identifier to represent this node.
- **vpls_flush_on_topology_change** (Optional, Boolean)  
  Enable VPLS MAC flush on root protected CE interface receiving topology change.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<routing_instance>`.

## Import

Junos mstp can be imported using an id made up of `<routing_instance>`, e.g.

```shell
$ terraform import junos_mstp.mstp default
```
