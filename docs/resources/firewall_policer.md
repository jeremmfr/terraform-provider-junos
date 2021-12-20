---
layout: "junos"
page_title: "Junos: junos_firewall_policer"
sidebar_current: "docs-junos-resource-firewall-policer"
description: |-
  Create firewall policer
---

# junos_firewall_policer

Provides a firewall policer resource.

## Example Usage

```hcl
# Configure a firewall policer
resource junos_firewall_policer "policer_demo" {
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

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of policer.
- **filter_specific** (Optional, Boolean)  
  Policer is filter-specific.
- **if_exceeding** (Required, Block)  
  Define rate limits options.
  - **burst_size_limit** (Required, String)  
    Burst size limit in bytes.
  - **bandwidth_percent** (Optional, Number)  
    Bandwidth limit in percentage.
  - **bandwidth_limit** (Optional, String)  
    Bandwidth limit in bits/second.
- **then** (Required, Block)  
  Define action to take if the rate limits are exceeded.
  - **discard** (Optional, Boolean)  
    Discard the packet.
  - **forwarding_class** (Optional, String)  
    Classify packet to forwarding class.
  - **loss_priority** (Optional, String)  
    Packet's loss priority.
  - **out_of_profile** (Optional, Boolean)  
     Discard packets only if both congested and over threshold.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos firewall policer can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_firewall_policer.policer_demo policerDemo
```
