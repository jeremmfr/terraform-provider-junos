---
page_title: "Junos: junos_chassis_cluster"
---

# junos_chassis_cluster

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `chassis cluster` block and `interfaces fab0/1`.

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

- **fab0** (Required, Block)  
  Declare `interfaces fab0` configuration.
  - **member_interfaces** (Optional, List of String)  
    Member interfaces for the fabric interface.
  - **description** (Optional, String)  
    Text description of interface.
- **fab1** (Optional, Block)  
  Declare `interfaces fab1` configuration.
  - **member_interfaces** (Optional, List of String)  
    Member interfaces for the fabric interface.
  - **description** (Optional, String)  
    Text description of interface.
- **redundancy_group** (Required, Block List)  
  For each redundancy-group to declare. First in list have id=0, second id=1, etc.
  - **node0_priority** (Required, Number)  
    Priority of the node in the redundancy-group (1..254).
  - **node1_priority** (Required, Number)  
    Priority of the node in the redundancy-group (1..254).
  - **gratuitous_arp_count** (Optional, Number)  
    Number of gratuitous ARPs to send on an active interface after failover (1..16).
  - **hold_down_interval** (Optional, Number)  
    RG failover interval. RG0(300-1800) RG1+(0-1800) (0..1800 seconds)
  - **interface_monitor** (Optional, Block List)  
    For each monitoring interface to declare.  
    See [below for nested schema](#interface_monitor-arguments-for-redundancy_group).
  - **preempt** (Optional, Boolean)  
    Allow preemption of primaryship based on priority.
  - **preempt_delay** (Optional, Number)  
    Time to wait before taking over mastership (1..21600 seconds).  
    `preempt` need to be true.
  - **preempt_limit** (Optional, Number)  
    Max number of preemptive failovers allowed (1..50).  
    `preempt` need to be true.
  - **preempt_period** (Optional, Number)  
    Time period during which the limit is applied (1..1400 seconds).  
    `preempt` need to be true.
- **reth_count** (Required, Number)  
  Number of redundant ethernet interfaces (1..128)
- **config_sync_no_secondary_bootup_auto** (Optional, Boolean)  
  Disable auto configuration synchronize on secondary bootup.
- **control_ports** (Optional, Block Set)  
  For each combination of block arguments,
  enable the specific control port to use as a control link for the chassis cluster.  
  Only available for some higher end Juniper SRX devices.
  - **fpc** (Required, Number)  
    Flexible PIC Concentrator (FPC) slot number.
  - **port** (Required, Number)  
    Port number on which to configure the control port.
- **control_link_recovery** (Optional, Boolean)  
  Enable automatic control link recovery.
- **heartbeat_interval** (Optional, Number)  
  Interval between successive heartbeats (1000..2000 milliseconds).
- **heartbeat_threshold** (Optional, Number)  
  Number of consecutive missed heartbeats to indicate device failure (3..8).

---

### interface_monitor arguments for redundancy_group

- **name** (Required, String)  
  Name of the interface to monitor.
- **weight** (Required, Number)  
  Weight assigned to this interface that influences failover (0..255).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `cluster`.

## Import

Junos chassis cluster can be imported using any id, e.g.

```shell
$ terraform import junos_chassis_cluster.cluster random
```
