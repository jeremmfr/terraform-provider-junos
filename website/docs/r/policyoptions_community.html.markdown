---
layout: "junos"
page_title: "Junos: junos_policyoptions_community"
sidebar_current: "docs-junos-resource-policyoptions-community"
description: |-
  Create a community
---

# junos_policyoptions_community

Provides a community BGP resource.

## Example Usage

```hcl
# Add a community
resource junos_policyoptions_community "communityDemo" {
  name    = "communityDemo"
  members = [ "65000:100" ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of community.
* `members` - (Required)(`ListOfString`) List of community.
* `invert_match` - (Optional)(`Bool`) Add 'invert-match' parameter.

## Import

Junos community can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_policyoptions_community.communityDemo communityDemo
```
