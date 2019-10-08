---
layout: "junos"
page_title: "Junos: junos_policyoptions_as_path_group"
sidebar_current: "docs-junos-resource-policyoptions-as-path-group"
description: |-
  Create a as-path group
---

# junos_policyoptions_as_path_group

Provides a as-path group resource.

## Example Usage

```hcl
# Add a as-path group
resource junos_policyoptions_as_path_group "viaCenturyLink" {
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

* `name` - (Required, Forces new resource)(`String`) The name of as-path group.
* `as_path` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each as-path to declare.
  * `name` - (Required)(`String`) Name of as-path
  * `path` - (Required)(`String`) As-path
* `dynamic_db` - (Optional)(`Bool`) Add 'dynamic-db' parameter.

## Import

Junos as-path group can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_policyoptions_as_path_group.viaCenturyLink viaCenturyLink
```
