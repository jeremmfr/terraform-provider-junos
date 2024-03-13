---
page_title: "Junos: junos_multichassis"
---

# junos_multichassis

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `multi-chassis` block.  
By default (without `clean_on_destroy` = true), destroy this resource has no effect on the Junos configuration.

Configure static configuration in `multi-chassis` block.

## Example Usage

```hcl
# Configure multi-chassis
resource "junos_multichassis" "multichassis" {
  mc_lag_consistency_check = true
}
```

## Argument Reference

The following arguments are supported:

- **clean_on_destroy** (Optional, Boolean)  
  Clean entirely `multi-chassis` block when destroy this resource.  
  It includes potential `junos_multichassis_protection_peer` resources.
- **mc_lag_consistency_check** (Optional, Computed, Boolean)  
  Consistency Check.  
  Computed to set to `true` when `mc_lag_consistency_check_comparison_delay_time` is specified.
- **mc_lag_consistency_check_comparison_delay_time** (Optional, Number)  
  Time after which local and remote config are compared (5..600 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `multichassis`.

## Import

Junos multi-chassis can be imported using any id, e.g.

```shell
$ terraform import junos_multichassis.multichassis random
```
