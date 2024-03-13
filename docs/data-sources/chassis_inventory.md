---
page_title: "Junos: junos_chassis_inventory"
---

# junos_chassis_inventory

Get chassis inventory (like the `show chassis hardware` command).

## Example Usage

```hcl
# Read chassis inventory and display serial-number
data "junos_chassis_inventory" "demo" {}
output "SN" {
  value = data.junos_chassis_inventory.demo.chassis.0.serial_number
}
```

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with value `chassis_inventory`.
- **chassis** (Block List)  
  Chassis inventory for each routing engine.  
  - Attributes defined by [component information schema](#component-information-schema)
  - **re_name** (String)  
    Name of the routing engine (only if there are multiple routing engines).
  - **module** (Block List)  
    For each module, component information.  
    See [below for nested schema](#module-attributes).

### module attributes

- Attributes defined by [component information schema](#component-information-schema)
- **sub_module** (Block List)  
  For each sub-module, component information.  
  - Attributes defined by [component information schema](#component-information-schema)
  - **sub_sub_module** (Block List)  
    For each sub-sub-module, component information.  
    Attributes defined by [component information schema](#component-information-schema)

### component information schema

- **clei_code** (String)  
  Common Language Equipment Identifier code of the component.
- **description** (String)  
  Description of the component.
- **model_number** (String)  
  Model number of the component.
- **name** (String)  
  Name of the component.
- **part_number** (String)  
  Part number of the component.
- **serial_number** (String)  
  Serial number of the component.
- **version** (String)  
  Version of the component.
