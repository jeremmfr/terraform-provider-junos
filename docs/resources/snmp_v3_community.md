---
page_title: "Junos: junos_snmp_v3_community"
---

# junos_snmp_v3_community

Provides a snmp v3 community resource.

## Example Usage

```hcl
# Add a snmp v3 community
resource "junos_snmp_v3_community" "index1" {
  community_index = "index1"
  security_name   = "john"
}
```

## Argument Reference

The following arguments are supported:

- **community_index** (Required, String, Forces new resource)  
  Unique index value in this community table entry.
- **security_name** (Required, String)  
  Security name used when performing access control.  
- **community_name** (Optional, String)  
  SNMPv1/v2c community name (default is same as community-index).
- **context** (Optional, String)  
  Context used when performing access control.
- **tag** (Optional, String)  
  Tag identifier for set of targets allowed to use this community string.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<community_index>`.

## Import

Junos snmp v3 community can be imported using an id made up of `<community_index>`, e.g.

```shell
$ terraform import junos_snmp_v3_community.index1 index1
```
