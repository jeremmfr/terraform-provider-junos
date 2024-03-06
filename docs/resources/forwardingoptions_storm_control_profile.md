---
page_title: "Junos: junos_forwardingoptions_storm_control_profile"
---

# junos_forwardingoptions_storm_control_profile

Provides a forwarding-options storm-control-profile resource.

## Example Usage

```hcl
# Add a forwarding-options storm-control-profile
resource "junos_forwardingoptions_storm_control_profile" "demo" {
  name            = "demo"
  action_shutdown = true
  all {
    bandwidth_percentage = 80
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Storm control profile name.
- **all** (Required, Block)  
  For all BUM traffic.
  - **bandwidth_level** (Optional, Number)  
    Link bandwidth (100..100000000 kbps).
  - **bandwidth_percentage** (Optional, Number)  
    Percentage of link bandwidth (1..100).
  - **burst_size** (Optional, Number)  
    Burst size (1500..100000000 bytes).
  - **no_broadcast** (Optional, Boolean)  
    Disable broadcast storm control.
  - **no_multicast** (Optional, Boolean)  
    Disable multicast storm control.
  - **no_registered_multicast** (Optional, Boolean)  
    Disable registered multicast storm control.
  - **no_unknown_unicast** (Optional, Boolean)  
    Disable unknown unicast storm control.
  - **no_unregistered_multicast** (Optional, Boolean)  
    Disable unregistered multicast storm control.
- **action_shutdown** (Optional, Boolean)  
  Disable port for excessive storm control errors.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos forwarding-options storm-control-profile can be imported using an id made up of
`<name>`, e.g.

```shell
$ terraform import junos_forwardingoptions_storm_control_profile.demo demo
```
