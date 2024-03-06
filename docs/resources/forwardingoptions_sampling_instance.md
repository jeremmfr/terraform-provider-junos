---
page_title: "Junos: junos_forwardingoptions_sampling_instance"
---

# junos_forwardingoptions_sampling_instance

Provides a forwarding-options sampling instance resource.

## Example Usage

```hcl
# Add a forwarding-options sampling instance
resource "junos_forwardingoptions_sampling_instance" "demo" {
  name = "demo"
  family_inet_input {
    rate = 1
  }
  family_inet_output {
    inline_jflow_source_address = "192.0.2.2"
    flow_server {
      hostname               = "192.0.2.1"
      port                   = 3000
      version_ipfix_template = "demo"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name for sampling instance.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance if not root level.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`
- **disable** (Optional, Boolean)  
  Disable sampling instance.
- **family_inet_input** (Optional, Block)  
  Declare `family inet input` configuration.  
  See [below for nested schema](#input-arguments).
- **family_inet_output** (Optional, Block)  
  Declare `family inet output` configuration.  
  See [below for nested schema](#output-arguments).
- **family_inet6_input** (Optional, Block)  
  Declare `family inet6 input` configuration.  
  See [below for nested schema](#input-arguments).
- **family_inet6_output** (Optional, Block)  
  Declare `family inet6 output` configuration.  
  See [below for nested schema](#output-arguments).
- **family_mpls_input** (Optional, Block)  
  Declare `family mpls input` configuration.  
  See [below for nested schema](#input-arguments).
- **family_mpls_output** (Optional, Block)  
  Declare `family mpls output` configuration.  
  See [below for nested schema](#output-arguments).
- **input** (Optional, Block)  
  Declare `input` configuration.  
  See [below for nested schema](#input-arguments).

---

### input arguments

- **max_packets_per_second** (Optional, Number)  
  Threshold of samples per second before dropping.
- **maximum_packet_length** (Optional, Number)  
  Maximum length of the sampled packet (0..9192 bytes).
- **rate** (Optional, Number)  
  Ratio of packets to be sampled (1 out of N) (1..16000000).
- **run_length** (Optional, Number)  
  Number of samples after initial trigger (0..20).

---

### output arguments

- **aggregate_export_interval** (Optional, Number)  
  Interval of exporting aggregate accounting information (90..1800 seconds).
- **extension_service** (Optional, List of String)  
  Define the customer specific sampling configuration.  
  **Not available for `family mpls`**.
- **flow_active_timeout** (Optional, Number)  
  Interval after which an active flow is exported (60..1800 seconds).
- **flow_inactive_timeout** (Optional, Number)  
  Interval of inactivity that marks a flow inactive (15..1800 seconds).
- **flow_server** (Optional, Block Set)  
  For each hostname, configure sending traffic aggregates in cflowd format.
  - **hostname** (Required, String)  
    Name of host collecting cflowd packets.
  - **port** (Required, Number)  
    UDP port number on host collecting cflowd packets (1..65535).
  - **aggregation_autonomous_system** (Optional, Boolean)  
    Aggregate by autonomous system number.
  - **aggregation_destination_prefix** (Optional, Boolean)  
    Aggregate by destination prefix.
  - **aggregation_protocol_port** (Optional, Boolean)  
    Aggregate by protocol and port number.
  - **aggregation_source_destination_prefix** (Optional, Boolean)  
    Aggregate by source and destination prefix.
  - **aggregation_source_destination_prefix_caida_compliant** (Optional, Boolean)  
    Compatible with Caida record format for prefix aggregation (v8).
  - **aggregation_source_prefix** (Optional, Boolean)  
    Aggregate by source prefix.
  - **autonomous_system_type** (Optional, String)  
    Type of autonomous system number to export.  
    Need to be `origin` or `peer`.
  - **dscp** (Optional, Number)  
    Numeric DSCP value in the range 0 to 63 (0..63).
  - **forwarding_class** (Optional, String)  
    Forwarding-class for exported jflow packets, applicable only for inline-jflow.
  - **local_dump** (Optional, Boolean)  
    Dump cflowd records to log file before exporting.
  - **no_local_dump** (Optional, Boolean)  
    Don't dump cflowd records to log file before exporting.
  - **routing_instance** (Optional, String)  
    Name of routing instance on which flow collector is reachable.
  - **source_address** (Optional, String)  
    Source IPv4 address for cflowd packets.
  - **version** (Optional, Number)  
    Format of exported cflowd aggregates.  
    Need to be `5` or `8`.  
    **Only available for `family inet`**.
  - **version9_template** (Optional, String)  
    Template to export data in version 9 format.
  - **version_ipfix_template** (Optional, String)  
    Template to export data in version ipfix format.
- **inline_jflow_export_rate** (Optional, Number)  
  Inline processing of sampled packets with flow export rate of monitored packets in kpps (1..3200).
- **inline_jflow_source_address** (Optional, String)  
  Inline processing of sampled packets with address to use for generating monitored packets.
- **interface** (Optional, Block List)  
  For each name of interface, configure interfaces used to send monitored information.
  - **name** (Required, String)  
    Name of interface.
  - **engine_id** (Optional, Number)  
    Identity (number) of this accounting interface (0..255).
  - **engine_type** (Optional, Number)  
    Type (number) of this accounting interface (0..255).
  - **source_address** (Optional, String)  
    Address to use for generating monitored packets.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>`.

## Import

Junos forwarding-options sampling instance can be imported using an id made up of
`<name>` or `<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_forwardingoptions_sampling_instance.demo demo
```
