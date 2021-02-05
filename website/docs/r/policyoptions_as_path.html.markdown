---
layout: "junos"
page_title: "Junos: junos_policyoptions_as_path"
sidebar_current: "docs-junos-resource-policyoptions-as-path"
description: |-
  Create a as-path
---

# junos_policyoptions_as_path

Provides a as-path resource.

## Example Usage

```hcl
# Add a as-path
resource junos_policyoptions_as_path "github" {
  name = "github"
  path = ".* 36459"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of as-path.
* `dynamic_db` - (Optional)(`Bool`) Add 'dynamic-db' parameter.
* `path` - (Optional)(`String`) As-path.

## Import

Junos as-path can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_policyoptions_as_path.github github
```
