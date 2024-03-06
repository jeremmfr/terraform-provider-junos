---
page_title: "Junos: junos_snmp_v3_vacm_accessgroup"
---

# junos_snmp_v3_vacm_accessgroup

Provides a snmp v3 VACM access group resource.

## Example Usage

```hcl
# Add a snmpv3 VACM access group
resource "junos_snmp_v3_vacm_accessgroup" "group1" {
  name = "group1"
  default_context_prefix {
    model     = "any"
    level     = "none"
    read_view = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

-> **Note:** At least one of `context_prefix` or `default_context_prefix` need to be set

- **name** (Required, String, Forces new resource)  
  SNMPv3 VACM group name.
- **context_prefix** (Optional, Block List)  
  For each prefix of context-prefix access configuration.
  - **prefix** (Required, String)  
    SNMPv3 VACM context prefix.
  - **access_config** (Optional, Block Set)  
    For each combination of `model` and `level`, define context-prefix access configuration.  
    See [below for nested schema](#access_config-or-default_context_prefix-arguments).
- **default_context_prefix** (Optional, Block Set)  
  For each combination of `model` and `level`, define default context-prefix access configuration.  
  See [below for nested schema](#access_config-or-default_context_prefix-arguments).

---

### access_config or default_context_prefix arguments

- **model** (Required, String)  
  Security model access configuration.  
  Need to be `any`, `usm`, `v1` or `v2c`.
- **level** (Required, String)  
  Security level access configuration.  
  Need to be `authentication`, `none` or `privacy`.
- **context_match** (Optional, String)  
  Type of match to perform on context-prefix.  
  Need to be `exact` or `prefix`.
- **notify_view** (Optional, String)  
  View used to notifications.
- **read_view** (Optional, String)  
  View used for read access.
- **write_view** (Optional, String)  
  View used for write access.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos snmp v3 VACM access group can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_snmp_v3_vacm_accessgroup.group1 group1
```
