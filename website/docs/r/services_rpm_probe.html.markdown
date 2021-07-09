---
layout: "junos"
page_title: "Junos: junos_services_rpm_probe"
sidebar_current: "docs-junos-resource-services-rpm-probe"
description: |-
  Create a services rpm probe
---

# junos_services_rpm_probe

Provides a services rpm probe resource.

## Example Usage

```hcl
# Add a services rpm probe
resource "junos_services_rpm_probe" "demo" {
  name = "demo"
  test {
    name           = "test_ping"
    target_type    = "address"
    target_value   = "192.0.2.33"
    source_address = "192.0.2.32"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Name of owner.
* `delegate_probes` - (Optional)(`Bool`) Offload real-time performance monitoring probes to MS-MIC/MS-MPC card.
* `test` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure a test. Can be specified multiple times for each test name. See the [`test` arguments] (#test-arguments) block.

---

### test arguments

* `name` - (Required)(`String`) Name of test.
* `data_fill` - (Optional)(`String`) Define contents of the data portion of the probes.
* `data_size` - (Optional)(`Int`) Size of the data portion of the probes (0..65400).
* `destination_interface` - (Optional)(`String`) Name of output interface for probes.
* `destination_port` - (Optional)(`Int`) TCP/UDP port number (7..65535).
* `dscp_code_points` - (Optional)(`String`) Differentiated Services code point bits or alias.
* `hardware_timestamp` - (Optional)(`Bool`) Packet Forwarding Engine updates timestamps.
* `history_size` - (Optional)(`Int`) Number of stored history entries (0..512).
* `inet6_source_address` - (Optional)(`String`) Inet6 Source Address of the probe.
* `moving_average_size` - (Optional)(`Int`) Number of samples used for moving average (0..1024).
* `one_way_hardware_timestamp` - (Optional)(`Bool`) Enable hardware timestamps for one-way measurements.
* `probe_count` - (Optional)(`Int`) Total number of probes per test (1..15).
* `probe_interval` - (Optional)(`Int`) Delay between probes (1..255 seconds).
* `probe_type` - (Optional)(`String`) Probe request type.
* `routing_instance` - (Optional)(`String`) Routing instance used by probes.
* `rpm_scale` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to configure real-time performance monitoring scale tests.
  * `tests_count` - (Required)(`Int`) Number of probe-tests generated using scale config (1..500000).
  * `destination_interface` - (Optional)(`String`) Base destination interface for scale test.
  * `destination_subunit_cnt` - (Optional)(`Int`) Subunit count for destination interface for scale test (1..500000)
  * `source_address_base` - (Optional)(`String`) Source base address of host in a.b.c.d format.
  * `source_count` - (Optional)(`Int`) Source-address count (1..500000).
  * `source_step` - (Optional)(`String`) Steps to increment src address in a.b.c.d format.
  * `source_inet6_address_base` - (Optional)(`String`) Source base inet6 address of host in a:b:c:d:e:f:g:h format.
  * `source_inet6_count` - (Optional)(`Int`) Source-inet6-address count (1..500000).
  * `source_inet6_step` - (Optional)(`String`) Steps to increment src inet6 address in a:b:c:d:e:f:g:h format.
  * `target_address_base` - (Optional)(`String`) Base address of target host in a.b.c.d format.
  * `target_count` - (Optional)(`Int`) Target address count (1..500000).
  * `target_step` - (Optional)(`String`) Steps to increment target address in a.b.c.d format.
  * `target_inet6_address_base` - (Optional)(`String`) Base inet6 address of target host in a:b:c:d:e:f:g:h format.
  * `target_inet6_count` - (Optional)(`Int`) Target inet6 address count (1..500000)
  * `target_inet6_step` - (Optional)(`String`) Steps to increment target inet6 address in a:b:c:d:e:f:g:h format
* `source_address` - (Optional)(`String`) Source address for probe.
* `target_type` - (Optional)(`String`) Type of target destination for probe.
* `target_value` - (Optional)(`String`) Target destination for probe.
* `test_interval` - (Optional)(`Int`) Delay between tests (0..86400 seconds).
* `thresholds` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'thresholds' configuration. Set 0 to disable respective threshold.
  * `egress_time` - (Optional)(`Int`) Maximum source to destination time per probe (0..60000000 microseconds).
  * `ingress_time` - (Optional)(`Int`) Maximum destination to source time per probe (0..60000000 microseconds).
  * `jitter_egress` - (Optional)(`Int`) Maximum source to destination jitter per test (0..60000000 microseconds).
  * `jitter_ingress` - (Optional)(`Int`) Maximum destination to source jitter per test (0..60000000 microseconds).
  * `jitter_rtt` - (Optional)(`Int`) Maximum jitter per test (0..60000000 microseconds).
  * `rtt` - (Optional)(`Int`) Maximum round trip time per probe (0..60000000 microseconds).
  * `std_dev_egress` - (Optional)(`Int`) Maximum source to destination standard deviation per test (0..60000000 microseconds).
  * `std_dev_ingress` - (Optional)(`Int`) Maximum destination to source standard deviation per test (0..60000000 microseconds).
  * `std_dev_rtt` - (Optional)(`Int`) Maximum standard deviation per test (0..60000000 microseconds).
  * `successive_loss` - (Optional)(`Int`) Successive probe loss count indicating probe failure (0..15).
  * `total_loss` - (Optional)(`Int`) Total probe loss count indicating test failure (0..15).
* `traps` - (Optional)(`ListOfString`) Trap to send if threshold is met or exceeded.
* `ttl` - (Optional)(`Int`) Time to Live (hop-limit) value for an RPM IPv4(IPv6) packet (1..254).

## Import

Junos services rpm probe can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_rpm_probe.demo demo
```
