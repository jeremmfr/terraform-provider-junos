---
layout: "junos"
page_title: "Junos: junos_security_screen_whitelist"
sidebar_current: "docs-junos-resource-security-screen-whitelist"
description: |-
  Create a security screen white-list (when Junos device supports it)
---

# junos_security_screen_whitelist

Provides a security screen white-list resource.

## Example Usage

```hcl
# Add a security screen white-list
resource junos_security_screen_whitelist "demo_screen_whitelist" {
  name = "demo_screen_whitelist"
  address = [
    "192.0.2.128/26",
    "192.0.2.64/26",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of screen.
* `address` - (Required)(`ListOfString`) List of address. Need to be a valid CIDR network.

## Import

Junos security screen white-list can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_screen_whitelist.demo_screen_whitelist demo_screen_whitelist
```
