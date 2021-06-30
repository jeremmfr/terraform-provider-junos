---
layout: "junos"
page_title: "Junos: junos_security_idp_custom_attack_group"
sidebar_current: "docs-junos-resource-security-idp-custom-attack-group"
description: |-
  Create a security idp custom-attack-group (when Junos device supports it)
---

# junos_security_idp_custom_attack_group

Provides a security idp custom-attack-group resource.

## Example Usage

```hcl
# Add a idp custom-attack_group
resource junos_security_idp_custom_attack_group "demo_idp_custom_attack_group" {
  name   = "group_of_Attacks"
  member = ["custom_attack_1", "custom_attack_2"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of idp custom-attack-group.
* `member` - (Optional)(`ListOfString`) List of attacks/attack groups belonging to this group.

## Import

Junos security idp custom-attack-group can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_idp_custom_attack_group.demo_idp_custom_attack_group group_of_Attacks
```
