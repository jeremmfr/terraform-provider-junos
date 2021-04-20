---
layout: "junos"
page_title: "Junos: junos_snmp_clientlist"
sidebar_current: "docs-junos-resource-snmp-clientlist"
description: |-
  Create a snmp client-list
---

# junos_snmp_clientlist

Provides a snmp client-list resource.

## Example Usage

```hcl
# Add a snmp clientlist
resource junos_snmp_clientlist "list1" {
  name   = "list1"
  prefix = ["192.0.2.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of snmp client-list.
* `prefix` - (Optional)(`ListOfString`) Address or prefix.

## Import

Junos snmp client-list can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_snmp_clientlist.list1 list1
```
