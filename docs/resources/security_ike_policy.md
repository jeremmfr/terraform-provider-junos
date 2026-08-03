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
  Conflict with the other `pre_shared_key_*` arguments.
- **pre_shared_key_text_wo** (Optional, String, Sensitive, Write-only)  
  Preshared key wit format as text, not stored in state.  
  Requires `pre_shared_key_text_wo_version` and Terraform 1.11 or later.  
  Conflict with the other `pre_shared_key_*` arguments.
- **pre_shared_key_text_wo_version** (Optional, Number)  
  Version of `pre_shared_key_text_wo` to trigger the sending of its value.  
  Increment it to send the current value of `pre_shared_key_text_wo` to the device.  
  Requires `pre_shared_key_text_wo`.
- **pre_shared_key_hexa** (Optional, String, Sensitive)  
  Preshared key with format as hexadecimal.  
  Conflict with the other `pre_shared_key_*` arguments.
- **pre_shared_key_hexa_wo** (Optional, String, Sensitive, Write-only)  
  Preshared key with format as hexadecimal, not stored in state.  
  Requires `pre_shared_key_hexa_wo_version` and Terraform 1.11 or later.  
  Conflict with the other `pre_shared_key_*` arguments.
- **pre_shared_key_hexa_wo_version** (Optional, Number)  
  Version of `pre_shared_key_hexa_wo` to trigger the sending of its value.  
  Increment it to send the current value of `pre_shared_key_hexa_wo` to the device.  
  Requires `pre_shared_key_hexa_wo`.
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

!> **Warning**
  Write-only arguments cannot be filled by an import, so the preshared key read on the device
  is stored in `pre_shared_key_text` or `pre_shared_key_hexa`, and therefore in the Terraform
  state.  
  When the configuration uses a write-only argument, the next apply removes it from the state.
