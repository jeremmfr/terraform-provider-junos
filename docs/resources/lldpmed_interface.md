---
page_title: "Junos: junos_lldpmed_interface"
---

# junos_lldpmed_interface

Provides a LLDP MED interface resource.

## Example Usage

```hcl
resource "junos_lldpmed_interface" "all" {
  name = "all"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Interface name or `all`.
- **disable** (Optional, Boolean)  
  Disable LLDP-MED.
- **enable** (Optional, Boolean)  
  Enable LLDP-MED.
- **location** (Optional, Block)  
  Define location.  
  See [below for nested schema](#location-arguments).

---

### location arguments

- **civic_based_ca_type** (Optional, Block List)  
  For each ca-type, configure civic-based ca-type.  
  `civic_based_country_code` need to be set.
  - **ca_type** (Required, Number)  
    Address element type (0..255).
  - **ca_value** (Optional, String)  
    Address element value.
- **civic_based_country_code** (Optional, String)  
  Two-letter country code.
- **civic_based_what** (Optional, Number)  
  Type of address (0..2).  
  `civic_based_country_code` need to be set.
- **co_ordinate_latitude** (Optional, Number)  
  Latitude value (0..360) to address based on longitude and latitude coordinates.
- **co_ordinate_longitude** (Optional, Number)  
  Longitude value (0..360) to address based on longitude and latitude coordinates.
- **elin** (Optional, String)  
  Emergency line identification (ELIN) string.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos lldp-med interface can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_lldpmed_interface.all all
```
