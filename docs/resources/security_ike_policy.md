---
page_title: "Junos: junos_security_ike_policy"
---

# junos_security_ike_policy

Provides a security IKE policy resource.

## Example Usage

```hcl
# Add an ike policy
resource "junos_security_ike_policy" "demo_vpn_policy" {
  name                = "ike-policy"
  proposals           = ["ike-proposal"]
  pre_shared_key_text = "theKey"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of IKE policy.
- **proposals** (Optional, List of String)  
  IKE proposals list.
- **proposal_set** (Optional, String)  
  Types of default IKE proposal-set.  
  Need to be `basic`, `compatible`, `prime-128`, `prime-256`, `standard`, `suiteb-gcm-128` or `suiteb-gcm-256`.
- **description** (Optional, String)  
  Text description of IKE policy.
- **mode** (Optional, String)  
  IKE mode for Phase 1.  
  Need to `main` or `aggressive`.  
  Defaults to `main`.
- **pre_shared_key_text** (Optional, String, Sensitive)  
  Preshared key wit format as text.
- **pre_shared_key_hexa** (Optional, String, Sensitive)  
  Preshared key with format as hexadecimal.
- **reauth_frequency** (Optional, Number)  
  Re-auth Peer after reauth-frequency times hard lifetime. (0-100)

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security IKE policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_ike_policy.demo_vpn_policy ike-policy
```
