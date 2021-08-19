---
layout: "junos"
page_title: "Junos: junos_security_utm_custom_url_category"
sidebar_current: "docs-junos-resource-security-utm-custom-url-category"
description: |-
  Create a security utm custom-object custom-url-category (when Junos device supports it)
---

# junos_security_utm_custom_url_category

Provides a security utm custom-object custom-url-category resource.

## Example Usage

```hcl
# Add a security utm custom-object custom-url-category
resource junos_security_utm_custom_url_category "demo_url_category" {
  name = "custom-category"
  value = [
    "custompattern1",
  ]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of security utm custom-object custom-url-category.
- **value** (Required, List of String)  
  List of url patterns for security utm custom-object custom-url-category.

## Import

Junos security utm custom-object url-category can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_custom_url_category.demo_url_category custom-category
```
