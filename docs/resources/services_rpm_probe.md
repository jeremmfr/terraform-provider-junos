---
page_title: "Junos: junos_services_rpm_probe"
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

- **name** (Required, String, Forces new resource)  
  Name of owner.
- **delegate_probes** (Optional, Boolean)  
  Offload real-time performance monitoring probes to MS-MIC/MS-MPC card.
- **test** (Optional, Block List)  
  For each name of test, configure a test.  
  See [below for nested schema](#test-arguments).

---

### test arguments

- **name** (Required, String)  
  Name of test.
- **data_fill** (Optional, String)  
  Define contents of the data portion of the probes.
- **data_size** (Optional, Number)  
  Size of the data portion of the probes (0..65400).
- **destination_interface** (Optional, String)  
  Name of output interface for probes.
- **destination_port** (Optional, Number)  
  TCP/UDP port number (7..65535).
- **dscp_code_points** (Optional, String)  
  Differentiated Services code point bits or alias.
- **hardware_timestamp** (Optional, Boolean)  
  Packet Forwarding Engine updates timestamps.
- **history_size** (Optional, Number)  
  Number of stored history entries (0..512).
- **inet6_source_address** (Optional, String)  
  Inet6 Source Address of the probe.
- **moving_average_size** (Optional, Number)  
  Number of samples used for moving average (0..1024).
- **one_way_hardware_timestamp** (Optional, Boolean)  
  Enable hardware timestamps for one-way measurements.
- **probe_count** (Optional, Number)  
  Total number of probes per test (1..15).
- **probe_interval** (Optional, Number)  
  Delay between probes (1..255 seconds).
- **probe_type** (Optional, String)  
  Probe request type.
- **routing_instance** (Optional, String)  
  Routing instance used by probes.
- **rpm_scale** (Optional, Block)  
  Configure real-time performance monitoring scale tests.
  - **tests_count** (Required, Number)  
    Number of probe-tests generated using scale config (1..500000).
  - **destination_interface** (Optional, String)  
    Base destination interface for scale test.
  - **destination_subunit_cnt** (Optional, Number)  
    Subunit count for destination interface for scale test (1..500000)
  - **source_address_base** (Optional, String)  
    Source base address of host in `a.b.c.d` format.
  - **source_count** (Optional, Number)  
    Source-address count (1..500000).
  - **source_step** (Optional, String)  
    Steps to increment src address in `a.b.c.d` format.
  - **source_inet6_address_base** (Optional, String)  
    Source base inet6 address of host in `a:b:c:d:e:f:g:h` format.
  - **source_inet6_count** (Optional, Number)  
    Source-inet6-address count (1..500000).
  - **source_inet6_step** (Optional, String)  
    Steps to increment src inet6 address in `a:b:c:d:e:f:g:h` format.
  - **target_address_base** (Optional, String)  
    Base address of target host in `a.b.c.d` format.
  - **target_count** (Optional, Number)  
    Target address count (1..500000).
  - **target_step** (Optional, String)  
    Steps to increment target address in `a.b.c.d` format.
  - **target_inet6_address_base** (Optional, String)  
    Base inet6 address of target host in `a:b:c:d:e:f:g:h` format.
  - **target_inet6_count** (Optional, Number)  
    Target inet6 address count (1..500000)
  - **target_inet6_step** (Optional, String)  
    Steps to increment target inet6 address in `a:b:c:d:e:f:g:h` format
- **source_address** (Optional, String)  
  Source address for probe.
- **target_type** (Optional, String)  
  Type of target destination for probe.
- **target_value** (Optional, String)  
  Target destination for probe.
- **test_interval** (Optional, Number)  
  Delay between tests (0..86400 seconds).
- **thresholds** (Optional, Block)  
  Declare `thresholds` configuration.  
  Set 0 to disable respective threshold.
  - **egress_time** (Optional, Number)  
    Maximum source to destination time per probe (0..60000000 microseconds).
  - **ingress_time** (Optional, Number)  
    Maximum destination to source time per probe (0..60000000 microseconds).
  - **jitter_egress** (Optional, Number)  
    Maximum source to destination jitter per test (0..60000000 microseconds).
  - **jitter_ingress** (Optional, Number)  
    Maximum destination to source jitter per test (0..60000000 microseconds).
  - **jitter_rtt** (Optional, Number)  
    Maximum jitter per test (0..60000000 microseconds).
  - **rtt** (Optional, Number)  
    Maximum round trip time per probe (0..60000000 microseconds).
  - **std_dev_egress** (Optional, Number)  
    Maximum source to destination standard deviation per test (0..60000000 microseconds).
  - **std_dev_ingress** (Optional, Number)  
    Maximum destination to source standard deviation per test (0..60000000 microseconds).
  - **std_dev_rtt** (Optional, Number)  
    Maximum standard deviation per test (0..60000000 microseconds).
  - **successive_loss** (Optional, Number)  
    Successive probe loss count indicating probe failure (0..15).
  - **total_loss** (Optional, Number)  
    Total probe loss count indicating test failure (0..15).
- **traps** (Optional, Set of String)  
  Trap to send if threshold is met or exceeded.
- **ttl** (Optional, Number)  
  Time to Live (hop-limit) value for an RPM IPv4(IPv6) packet (1..254).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos services rpm probe can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_rpm_probe.demo demo
```
