---
page_title: "Junos: junos_iccp"
---

# junos_iccp

~> **Note**
  This resource should only be created **once**.  
  It's used to configure static (not object) options in `protocols iccp` block.

Configure static configuration in `protocols iccp` block.

## Example Usage

```hcl
# Configure protocol ICCP
resource "junos_iccp" "iccp" {
  local_ip_addr                   = "192.0.2.1"
  authentication_key              = "a_key"
  session_establishment_hold_time = 300
}
```

## Argument Reference

The following arguments are supported:

- **local_ip_addr** (Required, String)  
  Local IP address to use by default for all peers.
- **authentication_key** (Optional, String, Sensitive)  
  MD5 authentication key for all peers.  
  Conflict with `authentication_key_wo`.
- **authentication_key_wo** (Optional, String, Sensitive, Write-only)  
  MD5 authentication key for all peers, not stored in state.  
  Requires `authentication_key_wo_version` and Terraform 1.11 or later.  
  Conflict with `authentication_key`.
- **authentication_key_wo_version** (Optional, Number)  
  Version of `authentication_key_wo` to trigger the sending of its value.  
  Increment it to send the current value of `authentication_key_wo` to the device.  
  Requires `authentication_key_wo`.
- **session_establishment_hold_time** (Optional, Number)  
  Time within which connection must succeed with peers (45..600 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `iccp`.

## Import

Junos protocols ICCP can be imported using any id, e.g.

```shell
$ terraform import junos_iccp.iccp random
```

!> **Warning**
  Write-only arguments cannot be filled by an import, so the MD5 authentication key read on the
  device is stored in `authentication_key`, and therefore in the Terraform state.  
  When the configuration uses `authentication_key_wo`, the next apply removes it from the state.
