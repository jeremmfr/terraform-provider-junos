---
page_title: "Junos: junos_iccp_peer"
---

# junos_iccp_peer

Provides an ICCP peer resource.

## Example Usage

```hcl
# Add an ICCP peer
resource "junos_iccp_peer" "peer1" {
  ip_address = "192.0.2.1"
}
```

## Argument Reference

The following arguments are supported:

- **ip_address** (Required, String)  
  IP address for this peer.
- **liveness_detection** (Required, Block)  
  Bidirectional Forwarding Detection options for the peer.
  - **detection_time_threshold** (Optional, Number)  
    High detection-time triggering a trap (milliseconds).
  - **minimum_interval** (Optional, Number)  
    Minimum transmit and receive interval (1..255000 milliseconds).
  - **minimum_receive_interval** (Optional, Number)  
    Minimum receive interval (1..255000 milliseconds).
  - **multiplier** (Optional, Number)  
    Detection time multiplier (1..255).
  - **no_adaptation** (Optional, Boolean)  
    Disable adaptation.
  - **transmit_interval_minimum_interval** (Optional, Number)  
    Minimum transmit interval (1..255000 milliseconds).
  - **transmit_interval_threshold** (Optional, Number)  
    High transmit interval triggering a trap (milliseconds).
  - **version** (Optional, String)  
    BFD protocol version number.  
    Need to be `0`, `1` or `automatic`.
- **redundancy_group_id_list** (Required, Set of Number)  
  List of redundancy groups this peer is part of.
- **authentication_key** (Optional, String, Sensitive)  
  MD5 authentication key.  
  Conflict with `authentication_key_wo`.
- **authentication_key_wo** (Optional, String, Sensitive, Write-only)  
  MD5 authentication key, not stored in state.  
  Requires `authentication_key_wo_version` and Terraform 1.11 or later.  
  Conflict with `authentication_key`.
- **authentication_key_wo_version** (Optional, Number)  
  Version of `authentication_key_wo` to trigger the sending of its value.  
  Increment it to send the current value of `authentication_key_wo` to the device.  
  Requires `authentication_key_wo`.
- **backup_liveness_detection** (Optional, Block)  
  Backup liveness detection.
  - **backup_peer_ip** (Optional, String)  
    Backup liveness detection peer's IP address.
- **local_ip_addr** (Optional, String)  
  Local IP address to use for this peer alone.
- **session_establishment_hold_time** (Optional, Number)  
  Time within which connection must succeed with this peer (45..600 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<ip_address>`.

## Import

Junos ICCP peer can be imported using an id made up of `<ip_address>`, e.g.

```shell
$ terraform import junos_iccp_peer.peer1 192.0.2.1
```

!> **Warning**
  Write-only arguments cannot be filled by an import, so the MD5 authentication key read on the
  device is stored in `authentication_key`, and therefore in the Terraform state.  
  When the configuration uses `authentication_key_wo`, the next apply removes it from the state.
