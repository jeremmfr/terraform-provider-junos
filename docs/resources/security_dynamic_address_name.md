---
page_title: "Junos: junos_security_dynamic_address_name"
---

# junos_security_dynamic_address_name

Provides a security dynamic-address address-name resource.

## Example Usage

```hcl
# Add a security dynamic-address address-name
resource "junos_security_dynamic_address_name" "demo_address_name" {
  name        = "demo"
  description = "demo junos_security_dynamic_address_name"
  profile_category {
    name = "GeoIP"
    property {
      name   = "country"
      string = ["AU", "CN"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Security dynamic address name.
- **description** (Optional, String)  
  Text description of dynamic address.
- **profile_feed_name** (Optional, String)  
  Name of feed in feed-server for this dynamic address.  
  Need to set one of `profile_feed_name` or `profile_category`.
- **profile_category** (Optional, Block)  
  Declare `profile category` configuration to categorize feed data into this dynamic address.  
  Need to set one of `profile_feed_name` or `profile_category`.  
  See [below for nested schema](#profile_category-arguments).

### profile_category arguments

- **name** (Required, String)  
  Name of category.
- **feed** (Optional, String)  
  Name of feed under category.
- **property** (Optional, Block List, Max: 3)  
  For each name of property to match.
  - **name** (Required, String)  
    Name of property.
  - **string** (Required, List of String)  
    String value.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `<name>`.

## Import

Junos security dynamic-address address-name can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_dynamic_address_name.demo_feed_srv demo
```
