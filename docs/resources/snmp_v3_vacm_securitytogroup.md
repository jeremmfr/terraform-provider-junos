---
page_title: "Junos: junos_snmp_v3_vacm_securitytogroup"
---

# junos_snmp_v3_vacm_securitytogroup

Provides a snmp v3 VACM security name assignment to group resource.

## Example Usage

```hcl
# Assigns security names to group
resource "junos_snmp_v3_vacm_securitytogroup" "read" {
  model = "usm"
  name  = "read"
  group = "group1"
}
```

## Argument Reference

The following arguments are supported:

- **model** (Required, String, Forces new resource)  
  Security model context for group assignment.  
  Need to be `usm`, `v1` or `v2c`.
- **name** (Required, String, Forces new resource)  
  Security name to assign to group.
- **group** (Required, String)  
  Group to which to assign security name.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<model>_-_<name>`.

## Import

Junos snmp v3 VACM security name assignment to group can be imported using an id made up of
`<model>_-_<name>`, e.g.

```shell
$ terraform import junos_snmp_v3_vacm_securitytogroup.read usm_-_read
```
