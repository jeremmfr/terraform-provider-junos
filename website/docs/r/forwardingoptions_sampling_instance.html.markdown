---
layout: "junos"
page_title: "Junos: junos_forwardingoptions_sampling_instance"
sidebar_current: "docs-junos-resource-forwardingoptions-sampling-instance"
description: |-
  Create a forwarding-options sampling-instance
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

* `name` - (Required, Forces new resource)(`String`) Name for sampling instance.
* `disable` - (Optional)(`Bool`) Disable sampling instance
* `family_inet_input` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'family inet input' configuration. See the [`input` arguments] (#input-arguments) block.
* `family_inet_output` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'family inet output' configuration. See the [`output` arguments] (#output-arguments) block.
* `family_inet6_input` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'family inet6 input' configuration. See the [`input` arguments] (#input-arguments) block.
* `family_inet6_output` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'family inet6 output' configuration. See the [`output` arguments] (#output-arguments) block.
* `family_mpls_input` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'family mpls input' configuration. See the [`input` arguments] (#input-arguments) block.
* `family_mpls_output` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'family inet6 output' configuration. See the [`output` arguments] (#output-arguments) block.
* `input` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'input' configuration. See the [`input` arguments] (#input-arguments) block.

---
#### input arguments
* `max_packets_per_second` - (Optional)(`Int`) Threshold of samples per second before dropping.
* `maximum_packet_length` - (Optional)(`Int`) Maximum length of the sampled packet (0..9192 bytes)
* `rate` - (Optional)(`Int`) Ratio of packets to be sampled (1 out of N) (1..16000000)
* `run_length` - (Optional)(`Int`) Number of samples after initial trigger (0..20)

---
#### output arguments
* `aggregate_export_interval` - (Optional)(`Int`) Interval of exporting aggregate accounting information (90..1800 seconds).
* `extension_service` - (Optional)(`ListOfString`) Define the customer specific sampling configuration. **Not available for `family mpls`**.
* `flow_active_timeout` - (Optional)(`Int`) Interval after which an active flow is exported (60..1800 seconds).
* `flow_inactive_timeout` - (Optional)(`Int`) Interval of inactivity that marks a flow inactive (15..1800 seconds).
* `flow_server` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure sending traffic aggregates in cflowd format. Can be specified multiple times for each hostname.
  * `hostname` - (Required)(`String`) Name of host collecting cflowd packets.
  * `port` - (Required)(`Int`) UDP port number on host collecting cflowd packets (1..65535).
  * `aggregation_autonomous_system` - (Optional)(`Bool`) Aggregate by autonomous system number.
  * `aggregation_destination_prefix` - (Optional)(`Bool`) Aggregate by destination prefix
  * `aggregation_protocol_port` - (Optional)(`Bool`) Aggregate by protocol and port number.
  * `aggregation_source_destination_prefix` - (Optional)(`Bool`) Aggregate by source and destination prefix.
  * `aggregation_source_destination_prefix_caida_compliant` - (Optional)(`Bool`) Compatible with Caida record format for prefix aggregation (v8).
  * `aggregation_source_prefix` - (Optional)(`Bool`) Aggregate by source prefix
  * `autonomous_system_type` - (Optional)(`String`) Type of autonomous system number to export. Need to be 'origin' or 'peer'.
  * `dscp` - (Optional)(`Int`) Numeric DSCP value in the range 0 to 63 (0..63).
  * `forwarding_class` - (Optional)(`String`) Forwarding-class for exported jflow packets, applicable only for inline-jflow.
  * `local_dump` - (Optional)(`Bool`) Dump cflowd records to log file before exporting.
  * `no_local_dump` - (Optional)(`Bool`) Don't dump cflowd records to log file before exporting.
  * `routing_instance` - (Optional)(`String`) Name of routing instance on which flow collector is reachable.
  * `source_address` - (Optional)(`String`) Source IPv4 address for cflowd packets.
  * `version` - (Optional)(`Int`) Format of exported cflowd aggregates. Need to be '5' or '8'. **Only available for `family inet`**. 
  * `version_ipfix_template` - (Optional)(`String`) Template to export data in version ipfix format.
  * `version9_template` - (Optional)(`String`) Template to export data in version 9 format.
* `inline_jflow_export_rate` - (Optional)(`Int`) Inline processing of sampled packets with flow export rate of monitored packets in kpps (1..3200)
* `inline_jflow_source_address` - (Optional)(`String`) Inline processing of sampled packets with address to use for generating monitored packets.
* `interface` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure Interfaces used to send monitored information. Can be specified multiple times for each interface.
  * `name` - (Required)(`String`) Name of interface.
  * `engine_id` - (Optional)(`Int`) Identity (number) of this accounting interface (0..255).
  * `engine_type` - (Optional)(`Int`) Type (number) of this accounting interface (0..255).
  * `source_address` - (Optional)(`String`) Address to use for generating monitored packets.

## Import

Junos forwarding-options sampling instance can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_forwardingoptions_sampling_instance.demo demo
```
