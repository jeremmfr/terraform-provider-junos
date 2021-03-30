---
layout: "junos"
page_title: "Junos: junos_services_security_intelligence_policy"
sidebar_current: "docs-junos-resource-services-security-intelligence-policy"
description: |-
  Create a services security-intelligence policy
---

# junos_services_security_intelligence_policy

Provides a services security-intelligence policy resource.

## Example Usage

```hcl
# Add a services security-intelligence policy
resource "junos_services_security_intelligence_policy" "demo" {
  name = "demo"
  category {
    name         = "CC"
    profile_name = "profile_CC"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Security intelligence policy name.
* `category` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure a profile for a category. Can be specified multiple times for each category name.
  * `name` - (Required)(`String`) Name of security intelligence category.
  * `profile_name` - (Required)(`String`) Name of profile.
* `description` - (Optional)(`String`) Text description of policy.

## Import

Junos services security-intelligence policy can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_services_security_intelligence_policy.demo demo
```
