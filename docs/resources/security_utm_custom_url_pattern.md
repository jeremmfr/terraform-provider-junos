---
page_title: "Junos: junos_security_utm_custom_url_pattern"
---

# junos_security_utm_custom_url_pattern

Provides a security utm custom-object url-pattern resource.

## Example Usage

```hcl
# Add a security utm custom-object url-pattern
resource "junos_security_utm_custom_url_pattern" "demo_url_pattern" {
  name = "Global_Whitelisted"
  value = [
    "http://*.google.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of security utm custom-object url-pattern.
- **value** (Required, List of String)  
  List of url for security utm custom-object url-pattern.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security utm custom-object url-pattern can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_custom_url_pattern.demo_url_pattern Global_Whitelisted
```
