---
layout: "junos"
page_title: "Junos: junos_security_ipsec_proposal"
sidebar_current: "docs-junos-resource-security-ipsec-proposal"
description: |-
  Create a security ipsec proposal (when Junos device supports it)
---

# junos_security_ipsec_proposal

Provides a security ipsec proposal resource.

## Example Usage

```hcl
# Add an ipsec proposal
resource junos_security_ipsec_proposal "demo_vpn_proposal" {
  name                     = "ipsec-proposal"
  authentication_algorithm = "hmac-sha1-96"
  encryption_algorithm     = "aes-128-cbc"
  lifetime_seconds         = 3600
  protocol                 = "esp"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of ipsec proposal.
- **authentication_algorithm** (Optional, String)  
  Authentication algorithm.
- **encryption_algorithm** (Optional, String)  
  Encryption algorithm.
- **lifetime_seconds** (Optional, Number)  
  Lifetime, in seconds.
- **lifetime_kilobytes** (Optional, Number)  
  Lifetime, in kilobytes.
- **protocol** (Optional, String)  
  IPSec protocol.  
  Need to be `esp` or `ah`.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security ipsec proposal can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_ipsec_proposal.demo_vpn_proposal ipsec-proposal
```
