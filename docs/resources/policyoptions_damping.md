---
page_title: "Junos: junos_policyoptions_damping"
---

# junos_policyoptions_damping

Provides a policy-options damping resource.

## Example Usage

```hcl
# Add a damping parameters
resource "junos_policyoptions_damping" "demo_damping" {
  name         = "DemoDamping"
  half_life    = 15
  max_suppress = 60
  reuse        = 750
  suppress     = 3000
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name to identify route flap damping parameters.
- **disable** (Optional, Boolean)  
  Disable damping.
- **half_life** (Optional, Number)  
  Decay half-life (1..45 minutes).
- **max_suppress** (Optional, Number)  
  Maximum hold-down time (1..720 minutes).  
  Must be greater than `half_life`.
- **reuse** (Optional, Number)  
  Reuse threshold (1..20000).
- **suppress** (Optional, Number)  
  Cutoff threshold (1..20000).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos policy-options damping can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_damping.demo_damping DemoDamping
```
