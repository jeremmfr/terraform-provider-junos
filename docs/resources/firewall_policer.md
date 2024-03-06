---
page_title: "Junos: junos_firewall_policer"
---

# junos_firewall_policer

Provides a firewall policer resource.

## Example Usage

```hcl
# Configure a firewall policer
resource "junos_firewall_policer" "policer_demo" {
  name            = "policerDemo"
  filter_specific = true
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit  = "50k"
  }
  then {
    discard = true
  }
}
```

## Argument Reference

-> **Note:** One of `if_exceeding` or `if_exceeding_pps` arguments is required.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Policer name.
- **filter_specific** (Optional, Boolean)  
  Policer is filter-specific.
- **logical_bandwidth_policer** (Optional, Boolean)  
  Policer uses logical interface bandwidth.
- **logical_interface_policer** (Optional, Boolean)  
  Policer is logical interface policer.
- **physical_interface_policer** (Optional, Boolean)  
  Policer is physical interface policer.
- **shared_bandwidth_policer** (Optional, Boolean)  
  Share policer bandwidth among bundle links.
- **if_exceeding** (Optional, Block)  
  Define rate limits options.
  - **burst_size_limit** (Required, String)  
    Burst size limit in bytes.  
    Format need to be `(\d)+(m|k|g)?`
  - **bandwidth_limit** (Optional, String)  
    Bandwidth limit in bits/second.  
    Format need to be `(\d)+(m|k|g)?`
  - **bandwidth_percent** (Optional, Number)  
    Bandwidth limit in percentage.
- **if_exceeding_pps** (Optional, Block)  
  Define pps limits options.
  - **packet_burst** (Required, String)  
    PPS burst size limit.
  - **pps_limit** (Required, String)  
    PPS limit.
- **then** (Required, Block)  
  Define action to take if the rate limits are exceeded.
  - **discard** (Optional, Boolean)  
    Discard the packet.
  - **forwarding_class** (Optional, String)  
    Classify packet to forwarding class.
  - **loss_priority** (Optional, String)  
    Packet's loss priority.  
    Need to be `high`, `low`, `medium-high` or `medium-low`.
  - **out_of_profile** (Optional, Boolean)  
     Discard packets only if both congested and over threshold.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos firewall policer can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_firewall_policer.policer_demo policerDemo
```
