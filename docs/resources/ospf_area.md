---
page_title: "Junos: junos_ospf_area"
---

# junos_ospf_area

Provides an ospf area resource.

## Example Usage

```hcl
# Add an ospf area
resource "junos_ospf_area" "demo_area" {
  area_id = "0.0.0.0"
  interface {
    name = "all"
  }
}
```

## Argument Reference

-> **Note** Some arguments are not compatible with version `v3` of ospf.

The following arguments are supported:

- **area_id** (Required, String, Forces new resource)  
  The id of ospf area.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for area.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`.
- **realm** (Optional, String, Forces new resource)  
  OSPFv3 realm configuration.  
  Need to be `ipv4-unicast`, `ipv4-multicast` or `ipv6-multicast`.  
  `version` need to be `v3`.
- **version** (Optional, String, Forces new resource)  
  Version of ospf.  
  Need to be `v2` or `v3`.  
  Defaults to `v2`.
- **interface** (Required, Block List)  
  For each interface or interface-range to declare.
  - **name** (Required, String)  
    Name of interface or interface-range.
  - **authentication_simple_password** (Optional, String, Sensitive)  
    Authentication key.  
    Conflict with `authentication_md5`.
  - **authentication_md5** (Optional, Block List)  
    For each key_id, MD5 authentication key.  
    See [below for nested schema](#authentication_md5-arguments).  
    Conflict with `authentication_simple_password`.
  - **bandwidth_based_metrics** (Optional, Block Set)  
    For each bandwidth, configure bandwidth based metrics.  
    See [below for nested schema](#bandwidth_based_metrics-arguments).
  - **bfd_liveness_detection** (Optional, Block)  
    Bidirectional Forwarding Detection options.  
    See [below for nested schema](#bfd_liveness_detection-arguments).
  - **dead_interval** (Optional, Number)  
    Dead interval (seconds).
  - **demand_circuit** (Optional, Boolean)  
    Interface functions as a demand circuit.
  - **disable** (Optional, Boolean)  
    Disable OSPF on this interface.
  - **dynamic_neighbors** (Optional, Boolean)  
    Learn neighbors dynamically on a p2mp interface.
  - **flood_reduction** (Optional, Boolean)  
    Enable flood reduction.
  - **hello_interval** (Optional, Number)  
    Hello interval (seconds).
  - **interface_type** (Optional, String)  
    Type of interface.
  - **ipsec_sa** (Optional, String)  
    IPSec security association name.
  - **ipv4_adjacency_segment_protected_type** (Optional, String)  
    Type to define adjacency SID is eligible for protection.  
    Need to be `dynamic`, `index` or `label`.
  - **ipv4_adjacency_segment_protected_value** (Optional, String)  
    Value for index or label to define adjacency SID is eligible for protection.
  - **ipv4_adjacency_segment_unprotected_type** (Optional, String)  
    Type to define adjacency SID uneligible for protection.  
    Need to be `dynamic`, `index` or `label`.
  - **ipv4_adjacency_segment_unprotected_value** (Optional, String)  
    Value for index or label to define adjacency SID uneligible for protection.
  - **link_protection** (Optional, Boolean)  
    Protect interface from link faults only.
  - **metric** (Optional, Number)  
    Interface metric.
  - **mtu** (Optional, Number)  
    Maximum OSPF packet size (128..65535).
  - **neighbor** (Optional, Block Set)  
    For each address, configure NBMA neighbor.  
    See [below for nested schema](#neighbor-arguments).
  - **no_advertise_adjacency_segment** (Optional, Boolean)  
    Do not advertise an adjacency segment for this interface.
  - **no_eligible_backup** (Optional, Boolean)  
    Not eligible to backup traffic from protected interfaces.
  - **no_eligible_remote_backup** (Optional, Boolean)  
    Not eligible for Remote-LFA backup traffic from protected interfaces.
  - **no_interface_state_traps** (Optional, Boolean)  
    Do not send interface state change traps.
  - **no_neighbor_down_notification** (Optional, Boolean)  
    Don't inform other protocols about neighbor down events.
  - **node_link_protection** (Optional, Boolean)  
    Protect interface from both link and node faults.
  - **passive** (Optional, Boolean)  
    Do not run OSPF, but advertise it.
  - **passive_traffic_engineering_remote_node_id** (Optional, String)  
    Advertise TE link information, remote address of the link.
  - **passive_traffic_engineering_remote_node_router_id** (Optional, String)  
    Advertise TE link information, TE Router-ID of the remote node.
  - **poll_interval** (Optional, Number)  
    Poll interval for NBMA interfaces (1..65535).
  - **priority** (Optional, Number)  
    Designated router priority (0..255).
  - **retransmit_interval** (Optional, Number)  
    Retransmission interval (seconds).
  - **secondary** (Optional, Boolean)  
    Treat interface as secondary.
  - **strict_bfd** (Optional, Boolean)  
    Enable strict bfd over this interface
  - **te_metric** (Optional, Number)  
    Traffic engineering metric (1..4294967295).
  - **transit_delay** (Optional, Number)  
    Transit delay (seconds) (1..65535).

---

### authentication_md5 arguments

- **key_id** (Required, Number)  
  Key ID for MD5 authentication (0..255).
- **key** (Required, String, Sensitive)  
  MD5 authentication key value.
- **start_time** (Optional, String)  
  Start time for key transmission (YYYY-MM-DD.HH:MM:SS).

---

### bandwidth_based_metrics arguments

- **bandwidth** (Required, String)  
  Bandwidth threshold.  
  Format need to be `(\d)+(m|k|g)?`
- **metric** (Required, Number)  
  Metric associated with specified bandwidth (1..65535).

---

### bfd_liveness_detection arguments

- **authentication_algorithm** (Optional, String)  
  Authentication algorithm name.
- **authentication_key_chain** (Optional, String)  
  Authentication key chain name.
- **authentication_loose_check** (Optional, Boolean)  
  Verify authentication only if authentication is negotiated.
- **detection_time_threshold** (Optional, Number)  
  High detection-time triggering a trap (milliseconds).
- **full_neighbors_only** (Optional, Boolean)  
  Setup BFD sessions only to Full neighbors.
- **holddown_interval** (Optional, Number)  
  Time to hold the session-UP notification to the client (0..255000 milliseconds).
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

---

### neighbor arguments

- **address** (Required, String)  
  Address of neighbor.
- **eligible** (Optional, Boolean)  
  Eligible to be DR on an NBMA network.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format  
  `<aread_id>_-_<version>_-_<routing_instance>`  
  or `<aread_id>_-_<version>_-_<realm>_-_<routing_instance>` if realm is set.

## Import

Junos ospf area can be imported using an id made up of  
`<aread_id>_-_<version>_-_<routing_instance>` or  
`<aread_id>_-_<version>_-_<realm>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_ospf_area.demo_area 0.0.0.0_-_v2_-_default
$ terraform import junos_ospf_area.demo_area2 0.0.0.0_-_v3_-_ipv4-unicast_-_default
```
