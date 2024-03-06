---
page_title: "Junos: junos_security_idp_custom_attack_group"
---

# junos_security_idp_custom_attack_group

Provides a security idp custom-attack-group resource.

## Example Usage

```hcl
# Add an idp custom-attack_group
resource "junos_security_idp_custom_attack_group" "demo_idp_custom_attack_group" {
  name   = "group_of_Attacks"
  member = ["custom_attack_1", "custom_attack_2"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of idp custom-attack-group.
- **member** (Optional, Set of String)  
  List of attacks/attack groups belonging to this group.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security idp custom-attack-group can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_idp_custom_attack_group.demo_idp_custom_attack_group group_of_Attacks
```
