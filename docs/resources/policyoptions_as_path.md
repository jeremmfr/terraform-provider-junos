---
page_title: "Junos: junos_policyoptions_as_path"
---

# junos_policyoptions_as_path

Provides a policy-options as-path resource.

## Example Usage

```hcl
# Add a policy-options as-path
resource "junos_policyoptions_as_path" "github" {
  name = "github"
  path = ".* 36459"
}
```

## Argument Reference

-> **Note:** One of `dynamic_db` or `path` arguments is required.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name to identify AS path regular expression.
- **dynamic_db** (Optional, Boolean)  
  Object may exist in dynamic database.
- **path** (Optional, String)  
  AS path regular expression.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos as-path can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_as_path.github github
```
