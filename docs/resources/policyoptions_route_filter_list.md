---
page_title: "Junos: junos_policyoptions_route_filter_list"
---

# junos_policyoptions_route_filter_list

Provides a route filter list resource.

## Example Usage

```hcl
# Add a route filter list
resource "junos_policyoptions_route_filter_list" "demo_rflist" {
  name = "DemoRFList"
  address {
    address = "192.0.2.0/25"
    option  = "orlonger"
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Route filter list name.
- **address** (Optional, Block Set)  
  List of addresses.  
  Conflict with `dynamic_db`.
  - **address** (Required, String)  
    IP address.
  - **option** (Required, String)  
    Mask option.  
    Need to be `address-mask`, `exact`, `longer`, `orlonger`, `prefix-length-range`, `through` or `upto`.
  - **option_value** (Optional, String)  
    For options that need an argument.
- **dynamic_db** (Optional, Boolean)  
  Object may exist in dynamic database.  
  Conflict with `address`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos route filter list can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_route_filter_list.demo_rflist DemoRFList
```
