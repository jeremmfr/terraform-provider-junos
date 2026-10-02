---
page_title: "Junos: junos_policyoptions_source_address_filter_list"
---

# junos_policyoptions_source_address_filter_list

Provides a source address filter list resource.

## Example Usage

```hcl
# Add a source address filter list
resource "junos_policyoptions_source_address_filter_list" "demo_saflist" {
  name = "DemoSAFList"
  address {
    address = "192.0.2.0/25"
    option  = "orlonger"
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Source address filter list name.
- **address** (Optional, Block Set)  
  List of source addresses.  
  Conflict with `dynamic_db`.
  - **address** (Required, String)  
    IP address.
  - **option** (Required, String)  
    Mask option.  
    Need to be `exact`, `longer`, `orlonger`, `prefix-length-range`, `through` or `upto`.
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

Junos source address filter list can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_source_address_filter_list.demo_saflist DemoSAFList
```
