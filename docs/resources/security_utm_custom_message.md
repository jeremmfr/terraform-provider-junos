---
page_title: "Junos: junos_security_utm_custom_message"
---

# junos_security_utm_custom_message

Provides a security utm custom-object custom-message resource.

## Example Usage

```hcl
# Add a security utm custom-object custom-message
resource "junos_security_utm_custom_message" "demo_message" {
  name    = "demo"
  type    = "user-message"
  content = "Demo Message"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of security utm custom-object custom-message.
- **type** (Required, String)  
  Type of custom message.  
  Need to be `custom-page`, `redirect-url` or `user-message`.
- **content** (Optional, String)  
  Content of custom message.  
  `type` need to be `redirect-url` or `user-message`.
- **custom_page_file** (Optional, String)  
  Name of custom page file.  
  `type` need to be `custom-page`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security utm custom-object custom-message can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_custom_message.demo_message demo
```
