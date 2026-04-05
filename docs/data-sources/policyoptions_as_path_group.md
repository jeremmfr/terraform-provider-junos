---
page_title: "Junos: junos_policyoptions_as_path_group"
---

# junos_policyoptions_as_path_group

Get configuration from a policy-options as-path-group.

## Example Usage

```hcl
# Read a policy-options as-path-group configuration
data "junos_policyoptions_as_path_group" "via_century_link" {
  name = "viaCenturyLink"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  Name to identify AS path group.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **as_path** (Block List)  
  List of as-path entries in this group.
  - **name** (String)  
    Name to identify AS path regular expression.
  - **path** (String)  
    AS path regular expression.
- **dynamic_db** (Boolean)  
  Object may exist in dynamic database.
