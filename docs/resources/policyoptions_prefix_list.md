---
page_title: "Junos: junos_policyoptions_prefix_list"
---

# junos_policyoptions_prefix_list

Provides a prefix list resource.

## Example Usage

```hcl
# Add a prefix list
resource "junos_policyoptions_prefix_list" "demo_plist" {
  name   = "DemoPList"
  prefix = ["192.0.2.0/25"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Prefix list name.
- **apply_path** (Optional, String)  
  Apply IP prefixes from a configuration statement.
- **dynamic_db** (Optional, Boolean)  
  Object may exist in dynamic database.
- **prefix** (Optional, Set of String)  
  Address prefixes.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos prefix list can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_prefix_list.demo_plist DemoPList
```
