---
page_title: "Junos: junos_policyoptions_source_address_filter_list"
---

# junos_policyoptions_source_address_filter_list

Get configuration from a policy-options source-address-filter-list.

## Example Usage

```hcl
# Read a policy-options source-address-filter-list configuration
data "junos_policyoptions_source_address_filter_list" "demo_saflist" {
  name = "DemoSAFList"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  Source address filter list name.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **address** (Block Set)  
  List of source addresses.
  - **address** (String)  
    IP address.
  - **option** (String)  
    Mask option.
  - **option_value** (String)  
    For options that need an argument.
- **dynamic_db** (Boolean)  
  Object may exist in dynamic database.
