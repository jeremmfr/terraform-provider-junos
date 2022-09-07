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

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source.
- **interface_names** (List of String)  
  List of interface names found.
- **interface_statuses** (Block List)  
  For each interface name.
  - **name** (String)  
    Interface name.
  - **admin_status** (String)  
    Admin status.
  - **oper_status** (String)  
    Operational status.
