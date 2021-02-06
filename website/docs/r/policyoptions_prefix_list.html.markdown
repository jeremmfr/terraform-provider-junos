---
layout: "junos"
page_title: "Junos: junos_policyoptions_prefix_list"
sidebar_current: "docs-junos-resource-policyoptions-prefix-list"
description: |-
  Create a prefix list
---

# junos_policyoptions_prefix_list

Provides a prefix list resource.

## Example Usage

```hcl
# Add a prefix list
resource junos_policyoptions_prefix_list "demo_plist" {
  name   = "DemoPList"
  prefix = ["192.0.2.0/25"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of prefix list.
* `apply_path` - (Optional)(`String`) Apply IP prefixes from a configuration statement.
* `dynamic_db` - (Optional)(`Bool`) Object may exist in dynamic database.
* `prefix` - (Optional)(`ListOfString`) List of CIDR.

## Import

Junos prefix list can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_policyoptions_prefix_list.demo_plist DemoPList
```
