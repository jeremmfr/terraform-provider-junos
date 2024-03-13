---
page_title: "Junos: junos_policyoptions_as_path_group"
---

# junos_policyoptions_as_path_group

Provides a policy-options as-path-group resource.

## Example Usage

```hcl
# Add a policy-options as-path-group
resource "junos_policyoptions_as_path_group" "via_century_link" {
  name = "viaCenturyLink"
  as_path {
    name = "qwest"
    path = ".* 209 .*"
  }
  as_path {
    name = "level3"
    path = ".* 3356 .*"
  }
  as_path {
    name = "level3-bis"
    path = ".* 3549 .*"
  }
  as_path {
    name = "twtc"
    path = ".* 4323 .*"
  }
}
```

## Argument Reference

-> **Note:** One of `dynamic_db` or `as_path` arguments is required.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name to identify AS path group.
- **as_path** (Optional, Block List)  
  For each name of as-path to declare.
  - **name** (Required, String)  
    Name to identify AS path regular expression.
  - **path** (Required, String)  
    AS path regular expression.
- **dynamic_db** (Optional, Boolean)  
  Object may exist in dynamic database.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos as-path group can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_as_path_group.via_century_link viaCenturyLink
```
