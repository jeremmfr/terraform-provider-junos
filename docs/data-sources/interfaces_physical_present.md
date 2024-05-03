---
page_title: "Junos: junos_interfaces_physical_present"
---

# junos_interfaces_physical_present

Get list of all of filtered physical interfaces present on the Junos device and their
admin/operational statuses.

## Example Usage

```hcl
# All interfaces that begin with 'ge-'
data "junos_interfaces_physical_present" "interfaces_ge" {
  match_name = "^ge-.*$"
}
```

## Argument Reference

The following arguments are supported:

- **match_name** (Optional, String)  
  A regexp to apply filter on name.  
  Need to be a valid regexp.
- **match_admin_up** (Optional, Bool)  
  Filter on interfaces that have admin status `up`.
- **match_oper_up** (Optional, Bool)  
  Filter on interfaces that have operational status `up`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source.
- **interface_names** (List of String)  
  Found interface names.
- **interfaces** (Block Map)  
  Dictionary of found interfaces with interface name as key.
  - **name** (String)  
    Interface name (as the map key).
  - **admin_status** (String)  
    Admin status.
  - **oper_status** (String)  
    Operational status.
  - **logical_interface_names** (List of String)  
    Logical interface names under this physical interface.
- **interface_statuses** (Block List, **Deprecated**)  
  For each found interface name, its status.  
  Deprecated attribute, use the `interfaces` attribute instead.
  - **name** (String)  
    Interface name.
  - **admin_status** (String)  
    Admin status.
  - **oper_status** (String)  
    Operational status.
