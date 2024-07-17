---
page_title: "Junos: junos_rip_neighbor"
---

# junos_rip_neighbor

Provides a RIP or RIPng neighbor resource.

## Example Usage

```hcl
# Add a RIPng neighbor
resource "junos_rip_neighbor" "demo_rip" {
  name  = "ae0.0"
  group = "group1"
  ng    = true
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Interface name.  
  Need to be a logical interface or `all`.
- **group** (Required, String, Forces new resource)  
  Name of RIP or RIPng group for this neighbor.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for neighbor.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`.
- **ng** (Optional, Boolean, Forces new resource)  
  Protocol `ripng` instead of `rip`.
- **any_sender** (Optional, Boolean)  
  Disable strict checks on sender address.  
  Conflict with `ng`.
- **authentication_key** (Optional, String, Sensitive)  
  Authentication key (password).  
  Conflict with `authentication_selective_md5`, `ng`.
- **authentication_selective_md5** (Optional, Block List)  
  For each key_id, MD5 authentication key.  
  Conflict with `authentication_key`, `authentication_type`, `ng`.
  - **key_id** (Required, Number)  
    Key ID for MD5 authentication (0..255).
  - **key** (Required, String, Sensitive)  
    MD5 authentication key value.
  - **start_time** (Optional, String)  
    Start time for key transmission (YYYY-MM-DD.HH:MM:SS).
- **authentication_type** (Optional, String)  
  Authentication type.  
  Need to be `md5`, `none` or `simple`.  
  Conflict with `authentication_selective_md5`, `ng`.
- **bfd_liveness_detection** (Optional, Block)  
  Bidirectional Forwarding Detection options.  
  Conflict with `ng`.
  - **authentication_algorithm** (Optional, String)  
    Authentication algorithm name.
  - **authentication_key_chain** (Optional, String)  
    Authentication key chain name.
  - **authentication_loose_check** (Optional, Boolean)  
    Verify authentication only if authentication is negotiated.
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
- **check_zero** (Optional, Boolean)  
  Check reserved fields on incoming RIPv1 packets.  
  Conflict with `no_check_zero`, `ng`.
- **no_check_zero** (Optional, Booelan)  
  Don't check reserved fields on incoming RIPv1 packets.  
  Conflict with `check_zero`, `ng`.
- **demand_circuit** (Optional, Boolean)  
  Enable demand circuit.  
  Conflict with `ng`.
- **dynamic_peers** (Optional, Boolean)  
  Learn peers dynamically on a p2mp interface.  
  Conflict with `ng`.  
  `interface_type_p2mp` need to be true.
- **import** (Optional, List of String)  
  Import policy.
- **interface_type_p2mp** (Optional, Boolean)  
  Point-to-multipoint link.
- **max_retrans_time** (Optional, Number)  
  Maximum time to re-transmit a message in demand-circuit (5..180).  
  Conflict with `ng`.
- **message_size** (Optional, Number)  
  Number of route entries per update message (25..255).  
  Conflict with `ng`.
- **metric_in** (Optional, Number)  
  Metric value to add to incoming routes (1..15).
- **peer** (Optional, Set of String)  
  P2MP peer.  
  Conflict with `ng`.  
  `interface_type_p2mp` need to be true.  
  Need to be valid IP addresses.
- **receive** (Optional, String)  
  Configure RIP receive options.  
  Need to be `both`, `none`, `version-1` or `version-2`.
- **route_timeout** (Optional, Number)  
  Delay before routes time out (30..360 seconds).
- **send** (Optional, String)  
  Configure RIP send options.  
  Need to be `broadcast`, `multicast`, `none` or `version-1`.
- **update_interval** (Optional, Number)  
  Interval between regular route updates (10..60 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<group>_-_<routing_instance>`  
  or `<name>_-_<group>_-_ng_-_<routing_instance>` if `ng` is set to `true`.

## Import

Junos RIP or RIPng group can be imported using an id made up of
`<name>_-_<group>_-_<routing_instance>` or `<name>_-_<group>_-_ng_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_rip_neighbor.demo_rip ae0.0_-_group1_-_ng_-_default
```
