---
layout: "junos"
page_title: "Junos: junos_chassis_cluster"
sidebar_current: "docs-junos-resource-chassis-cluster"
description: |-
  Configure chassis cluster configuration (when Junos device supports it)
---

# junos_chassis_cluster

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `chassis cluster` block and `interfaces fab0/1`.

Configure static configuration in `chassis cluster` block and `interfaces fab0/1`.

## Example Usage

```hcl
# Configure chassis cluster
resource "junos_chassis_cluster" "cluster" {
  fab0 {
    member_interfaces = ["ge-0/0/3"]
  }
  redundancy_group { # id 0
    node0_priority = 100
    node1_priority = 50
  }
  redundancy_group { # id 1
    node0_priority = 100
    node1_priority = 50
    interface_monitor {
      name   = "ge-0/0/4"
      weight = 255
    }
    interface_monitor {
      name   = "ge-0/0/5"
      weight = 254
    }
  }
  reth_count = 2
}
```

## Argument Reference

The following arguments are supported:

* `fab0` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'interfaces fab0' configuration.
  * `member_interfaces` - (Optional)(`ListOfString`) Member interfaces for the fabric interface.
  * `description` - (Optional)(`String`) Text description of interface.
* `fab1` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'interfaces fab1' configuration.
  * `member_interfaces` - (Optional)(`ListOfString`) Member interfaces for the fabric interface.
  * `description` - (Optional)(`String`) Text description of interface.
* `redundancy_group` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each redundancy-group to declare. First in list have id=0, second id=1, etc.
  * `node0_priority` - (Required)(`Int`) Priority of the node in the redundancy-group (1..254).
  * `node1_priority` - (Required)(`Int`) Priority of the node in the redundancy-group (1..254).
  * `gratuitous_arp_count` - (Optional)(`Int`) Number of gratuitous ARPs to send on an active interface after failover (1..16).
  * `hold_down_interval` - (Optional)(`Int`) RG failover interval. RG0(300-1800) RG1+(0-1800) (0..1800 seconds)
  * `interface_monitor` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each monitoring interface to declare. See the [`interface_monitor` arguments for redundancy_group] (#interface_monitor-arguments-for-redundancy_group) block.
  * `preempt` - (Optional)(`Bool`) Allow preemption of primaryship based on priority.
* `reth_count` - (Required)(`Int`) Number of redundant ethernet interfaces (1..128)
* `config_sync_no_secondary_bootup_auto` - (Optional)(`Bool`) Disable auto configuration synchronize on secondary bootup.
* `control_link_recovery` - (Optional)(`Bool`) Enable automatic control link recovery.
* `heartbeat_interval` - (Optional)(`Int`) Interval between successive heartbeats (1000..2000 milliseconds).
* `heartbeat_threshold` - (Optional)(`Int`) Number of consecutive missed heartbeats to indicate device failure (3..8).

---

### interface_monitor arguments for redundancy_group

* `name` - (Required)(`String`) Name of the interface to monitor.
* `weight` - (Required)(`Int`) Weight assigned to this interface that influences failover (0..255).

## Import

Junos chassis cluster can be imported using any id, e.g.

```shell
$ terraform import junos_chassis_cluster.cluster random
```
