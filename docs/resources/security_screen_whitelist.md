---
page_title: "Junos: junos_security_screen_whitelist"
---

# junos_security_screen_whitelist

Provides a security screen white-list resource.

## Example Usage

```hcl
# Add a security screen white-list
resource "junos_security_screen_whitelist" "demo_screen_whitelist" {
  name = "demo_screen_whitelist"
  address = [
    "192.0.2.128/26",
    "192.0.2.64/26",
  ]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of screen.
- **address** (Required, Set of String)  
  List of address.  
  Need to be a valid CIDR network.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security screen white-list can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_screen_whitelist.demo_screen_whitelist demo_screen_whitelist
```
