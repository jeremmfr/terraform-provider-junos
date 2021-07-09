---
layout: "junos"
page_title: "Junos: junos_security_ike_proposal"
sidebar_current: "docs-junos-resource-security-ike-proposal"
description: |-
  Create a security ike proposal (when Junos device supports it)
---

# junos_security_ike_proposal

Provides a security ike proposal resource.

## Example Usage

```hcl
# Add a ike proposal
resource junos_security_ike_proposal "demo_vpn_proposal" {
  name                     = "ike-proposal"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-128-cbc"
  dh_group                 = "group2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of ike proposal.
* `authentication_algorithm` - (Optional)(`String`) Authentication algorithm.
* `authentication_method` - (Optional)(`String`) Authentication method. Defaults to `pre-shared-keys`.
* `dh_group` - (Optional)(`String`) Diffie-Hellman Group.
* `encryption_algorithm` - (Optional)(`String`) Encryption algorithm.
* `lifetime_seconds` - (Optional)(`Int`) Lifetime, in seconds.

## Import

Junos security ike proposal can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_ike_proposal.demo_vpn_proposal ike-proposal
```
