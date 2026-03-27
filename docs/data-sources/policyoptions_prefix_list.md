---
page_title: "Junos: junos_policyoptions_prefix_list"
---

# junos_policyoptions_prefix_list

Get configuration from a policy-options prefix-list.

## Example Usage

```hcl
# Read a policy-options prefix-list configuration
data "junos_policyoptions_prefix_list" "demo_plist" {
  name = "DemoPList"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  Prefix list name.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **apply_path** (String)  
  Apply IP prefixes from a configuration statement.
- **dynamic_db** (Boolean)  
  Object may exist in dynamic database.
- **prefix** (Set of String)  
  Address prefixes.
