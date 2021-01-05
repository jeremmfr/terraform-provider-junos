---
layout: "junos"
page_title: "Junos: junos_security_utm_custom_url_pattern"
sidebar_current: "docs-junos-resource-security-utm-custom-url-pattern"
description: |-
  Create a security utm custom-object url-pattern (when Junos device supports it)
---

# junos_security_utm_custom_url_pattern

Provides a security utm custom-object url-pattern resource.

## Example Usage

```hcl
# Add a security utm custom-object url-pattern
resource junos_security_utm_custom_url_pattern "demo_url_pattern" {
  name = "Global_Whitelisted"
  value = [
    "http://*.google.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of security utm custom-object url-pattern.
* `value` - (Required)(`ListofString`) List of url for security utm custom-object url-pattern.

## Import

Junos security utm custom-object url-pattern can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_utm_custom_url_pattern.demo_url_pattern Global_Whitelisted
```
