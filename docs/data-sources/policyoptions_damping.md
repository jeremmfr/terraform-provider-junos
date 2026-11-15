---
page_title: "Junos: junos_policyoptions_damping"
---

# junos_policyoptions_damping

Get configuration from a policy-options damping.

## Example Usage

```hcl
# Read a policy-options damping configuration
data "junos_policyoptions_damping" "demo_damping" {
  name = "DemoDamping"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  Name to identify route flap damping parameters.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **disable** (Boolean)  
  Disable damping.
- **half_life** (Number)  
  Decay half-life (minutes).
- **max_suppress** (Number)  
  Maximum hold-down time (minutes).
- **reuse** (Number)  
  Reuse threshold.
- **suppress** (Number)  
  Cutoff threshold.
