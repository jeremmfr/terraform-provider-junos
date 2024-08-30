---
page_title: "Junos: junos_applications"
---

# junos_applications

~> **Note**
  This resource should only be created **once**.  
  It's used to configure entirely `applications` block.  

!> **Warning**
  Don't use `junos_application` or `junos_application_set` resources to avoid conflict with this resource.

Configure entirely `applications` block.

## Example Usage

```hcl
# Configure applications
resource "junos_applications" "applications" {
  application {
    name                 = "customPort555"
    application_protocol = "tcp"
    destination_port     = 555
  }
}
```

## Argument Reference

The following arguments are supported:

- **application** (Optional, Block Set)  
  For each name, define an application.  
  See [arguments of application resource](application#argument-reference) for nested schema.
- **application_set** (Optional, Block Set)  
  For each name, define an applications set.  
  See [arguments of application_set resource](application_set#argument-reference) for nested schema.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `applications`.

## Import

Junos applications can be imported using any id, e.g.

```shell
$ terraform import junos_applications.applications random
```
