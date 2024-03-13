---
page_title: "Junos: junos_snmp_clientlist"
---

# junos_snmp_clientlist

Provides a snmp client-list resource.

## Example Usage

```hcl
# Add a snmp clientlist
resource "junos_snmp_clientlist" "list1" {
  name   = "list1"
  prefix = ["192.0.2.0/24"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of snmp client-list.
- **prefix** (Optional, Set of String)  
  Address or prefix.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos snmp client-list can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_snmp_clientlist.list1 list1
```
