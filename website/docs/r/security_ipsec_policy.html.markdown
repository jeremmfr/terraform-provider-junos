---
layout: "junos"
page_title: "Junos: junos_security_ipsec_policy"
sidebar_current: "docs-junos-resource-security-ipsec-policy"
description: |-
  Create a security ipsec policy (when Junos device supports it)
---

# junos_security_ipsec_policy

Provides a security ipsec policy resource.

## Example Usage

```hcl
# Add a ipsec policy
resource junos_security_ipsec_policy "demo_vpn_policy" {
  name      = "ipsec-policy"
  proposals = ["ipsec-proposal"]
  pfs_keys  = "group2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of ipsec policy.
* `proposals` - (Required)(`ListOfString`) Ipsec proposal list.
* `pfs_keys` - (Optional)(`String`) Diffie-Hellman Group.

## Import

Junos security ipsec policy can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_ipsec_policy.demo_vpn_policy ipsec-policy
```
