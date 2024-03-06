---
page_title: "Junos: junos_security_ipsec_policy"
---

# junos_security_ipsec_policy

Provides a security IPSec policy resource.

## Example Usage

```hcl
# Add an IPSec policy
resource "junos_security_ipsec_policy" "demo_vpn_policy" {
  name      = "ipsec-policy"
  proposals = ["ipsec-proposal"]
  pfs_keys  = "group2"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of IPSec policy.
- **description** (Optional, String)  
  Text description of IPSec policy.
- **pfs_keys** (Optional, String)  
  Diffie-Hellman Group.
- **proposals** (Optional, List of String)  
  IPSec proposal list.
- **proposal_set** (Optional, String)  
  Types of default IPSec proposal-set.  
  Need to be `basic`, `compatible`, `prime-128`, `prime-256`, `standard`, `suiteb-gcm-128` or `suiteb-gcm-256`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security IPSec policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_ipsec_policy.demo_vpn_policy ipsec-policy
```
