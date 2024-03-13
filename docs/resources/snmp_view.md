---
page_title: "Junos: junos_snmp_view"
---

# junos_snmp_view

Provides a snmp view resource.

## Example Usage

```hcl
# Add a snmp view
resource "junos_snmp_view" "view1" {
  name        = "view1"
  oid_include = [".1"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of snmp view.
- **oid_include** (Optional, Set of String)  
  OID include list.
- **oid_exclude** (Optional, Set of String)  
  OID exclude list.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos snmp view can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_snmp_view.view1 view1
```
