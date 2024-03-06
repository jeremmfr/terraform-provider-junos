---
page_title: "Junos: junos_chassis_redundancy"
---

# junos_chassis_redundancy

-> **Note:** This resource should only be created **once**.
It's used to configure options in `chassis redundancy` block.  

Configure `chassis redundancy` block

## Example Usage

```hcl
# Configure chassis redundancy
resource "junos_chassis_redundancy" "chassis_redundancy" {
  graceful_switchover = true
}
```

## Argument Reference

The following arguments are supported:

- **failover_disk_read_threshold** (Optional, Number)  
  To failover, read threshold (ms) on disk underperform monitoring (1000..10000).
- **failover_disk_write_threshold** (Optional, Number)  
  To failover, write threshold (ms) on disk underperform monitoring (1000..10000).
- **failover_not_on_disk_underperform** (Optional, Boolean)  
  Prevent gstatd from initiating failovers in response to slow disks.
- **failover_on_disk_failure** (Optional, Boolean)  
  Failover on disk failure.
- **failover_on_loss_of_keepalives** (Optional, Boolean)  
  Failover on loss of keepalives.
- **graceful_switchover** (Optional, Boolean)  
  Enable graceful switchover on supported hardware.
- **keepalive_time** (Optional, Number)  
  Time before Routing Engine failover (2..10000 seconds).
- **routing_engine** (Optional, Block Set)  
  For each slot, redundancy options.
  - **slot** (Required, Number)  
    Routing Engine slot number (0..1).
  - **role** (Required, String)  
    Define role.  
    Need to be `backup`, `disabled` or `master`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `redundancy`.

## Import

Junos chassis redundancy can be imported using any id, e.g.

```shell
$ terraform import junos_chassis_redundancy.chassis_redundancy random
```
