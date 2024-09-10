---
page_title: "Junos: junos_rip_group"
---

# junos_rip_group

Provides a RIP or RIPng group resource.

## Example Usage

```hcl
# Add a RIPng group
resource "junos_rip_group" "demo_rip" {
  name = "group1"
  ng   = true
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of group.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for group.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`.
- **ng** (Optional, Boolean, Forces new resource)  
  Protocol `ripng` instead of `rip`.
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
- **demand_circuit** (Optional, Boolean)  
  Enable demand circuit.  
  Conflict with `ng`.
- **export** (Optional, List of String)  
  Export policy.
- **import** (Optional, List of String)  
  Import policy.
- **max_retrans_time** (Optional, Number)  
  Maximum time to re-transmit a message in demand-circuit (5..180).  
  Conflict with `ng`.
- **metric_out** (Optional, Number)  
  Default metric of exported routes (1..15).
- **preference** (Optional, Number)  
  Preference of routes learned by this group.
- **route_timeout** (Optional, Number)  
  Delay before routes time out (30..360 seconds).
- **update_interval** (Optional, Number)  
  Interval between regular route updates (10..60 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>`  
  or `<name>_-_ng_-_<routing_instance>` if `ng` is set to `true`.

## Import

Junos RIP or RIPng group can be imported using an id made up of
`<name>_-_<routing_instance>` or `<name>_-_ng_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_rip_group.demo_rip group1_-_ng_-_default
```
