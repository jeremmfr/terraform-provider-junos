---
page_title: "Junos: junos_policyoptions_as_path"
---

# junos_policyoptions_as_path

Get configuration from a policy-options as-path.

## Example Usage

```hcl
# Read a policy-options as-path configuration
data "junos_policyoptions_as_path" "github" {
  name = "github"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  Name to identify AS path regular expression.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **dynamic_db** (Boolean)  
  Object may exist in dynamic database.
- **path** (String)  
  AS path regular expression.
