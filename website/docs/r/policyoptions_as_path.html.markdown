---
layout: "junos"
page_title: "Junos: junos_policyoptions_as_path"
sidebar_current: "docs-junos-resource-policyoptions-as-path"
description: |-
  Create an as-path
---

# junos_policyoptions_as_path

Provides an as-path resource.

## Example Usage

```hcl
# Add an as-path
resource junos_policyoptions_as_path "github" {
  name = "github"
  path = ".* 36459"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of as-path.
- **dynamic_db** (Optional, Boolean)  
  Add `dynamic-db` parameter.
- **path** (Optional, String)  
  As-path.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos as-path can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_as_path.github github
```
