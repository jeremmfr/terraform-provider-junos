---
layout: "junos"
page_title: "Junos: junos_security_ike_policy"
sidebar_current: "docs-junos-resource-security-ike-policy"
description: |-
  Create a security ike policy (when Junos device supports it)
---

# junos_security_ike_policy

Provides a security ike policy resource.

## Example Usage

```hcl
# Add a ike policy
resource junos_security_ike_policy "demo_vpn_policy" {
  name                = "ike-policy"
  proposals           = ["ike-proposal"]
  pre_shared_key_text = "theKey"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of ike policy.
- **proposals** (Optional, List of String)  
  Ike proposals list.
- **proposal_set** (Optional, String)  
  Types of default IPSEC proposal-set.  
  Need to be `basic`, `compatible`, `prime-128`, `prime-256`, `standard`, `suiteb-gcm-128` or `suiteb-gcm-256`.
- **mode** (Optional, String)  
  IKE mode for Phase 1.  
  Need to `main` or `aggressive`.  
  Defaults to `main`.
- **pre_shared_key_text** (Optional, String, Sensitive)  
  Preshared key wit format as text.
- **pre_shared_key_hexa** (Optional, String, Sensitive)  
  Preshared key wit format as hexa.

## Import

Junos security ike policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_ike_policy.demo_vpn_policy ike-policy
```
