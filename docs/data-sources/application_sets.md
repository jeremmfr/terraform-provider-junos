---
page_title: "Junos: junos_application_sets"
---

# junos_application_sets

Get list of filtered application-sets on the Junos device  
(in `applications` and `group junos-defaults applications` level).

## Example Usage

```hcl
# Find default application-set junos-cifs 
data "junos_application_sets" "default_cifs" {
  match_applications = ["junos-netbios-session", "junos-smb-session"]
}
```

## Argument Reference

The following arguments are supported:

- **match_name** (Optional, String)  
  A regexp to apply a filter on application-sets name.  
  Need to be a valid regexp.
- **match_applications** (Optional, Set of String)  
  List of applications to apply a filter on application-sets.  
  The list needs to be exact to match.
- **match_application_sets** (Optional, Set of String)  
  List of application-sets to apply a filter on application-sets.  
  The list needs to be exact to match.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source.
- **application_sets** (Block List)  
  For each application-set found.
  - **name** (String)  
    Application set name.
  - **applications** (List of String)  
    List of included application names.
  - **application_set** (List of String)  
    List of included application-set names.
  - **description** (String)  
    Description for application-set.
