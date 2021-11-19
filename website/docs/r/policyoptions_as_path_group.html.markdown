---
layout: "junos"
page_title: "Junos: junos_policyoptions_as_path_group"
sidebar_current: "docs-junos-resource-policyoptions-as-path-group"
description: |-
  Create an as-path group
---

# junos_policyoptions_as_path_group

Provides an as-path group resource.

## Example Usage

```hcl
# Add an as-path group
resource junos_policyoptions_as_path_group "via_century_link" {
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

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of as-path group.
- **as_path** (Optional, Block List)  
  For each name of as-path to declare.
  - **name** (Required, String)  
    Name of as-path
  - **path** (Required, String)  
    As-path
- **dynamic_db** (Optional, Boolean)  
  Add `dynamic-db` parameter.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos as-path group can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_as_path_group.via_century_link viaCenturyLink
```
