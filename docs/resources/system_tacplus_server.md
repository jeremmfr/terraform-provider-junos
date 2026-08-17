---
page_title: "Junos: junos_system_tacplus_server"
---

# junos_system_tacplus_server

Configure a system tacplus-server.

## Example Usage

```hcl
# Add a system tacplus-server
resource "junos_system_tacplus_server" "demo_tacplus_server" {
  address = "192.0.2.1"
}
```

## Argument Reference

The following arguments are supported:

- **address** (Required, String, Forces new resource)  
  TACACS+ authentication server address.
- **port** (Optional, Number)  
  TACACS+ authentication server port number (1..65535).
- **routing_instance** (Optional, String)  
  Routing instance.
- **secret** (Optional, String, Sensitive)  
  Shared secret with the authentication server.  
  Conflict with `secret_wo`.
- **secret_wo** (Optional, String, Sensitive, Write-only)  
  Shared secret with the authentication server, not stored in state.  
  Requires `secret_wo_version` and Terraform 1.11 or later.  
  Conflict with `secret`.
- **secret_wo_version** (Optional, Number)  
  Version of `secret_wo` to trigger the sending of its value.  
  Increment it to send the current value of `secret_wo` to the device.  
  Requires `secret_wo`.
- **single_connection** (Optional, Boolean)  
  Optimize TCP connection attempts.
- **source_address** (Optional, String)  
  Use specified address as source address.
- **timeout** (Optional, Number)  
  Request timeout period (1..90 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<address>`.

## Import

Junos system tacplus-server can be imported using an id made up of `<address>`, e.g.

```shell
$ terraform import junos_system_tacplus_server.demo_tacplus_server 192.0.2.1
```

!> **Warning**
  Write-only arguments cannot be filled by an import, so the shared secret read on the device
  is stored in `secret`, and therefore in the Terraform state.  
  When the configuration uses `secret_wo`, the next apply removes it from the state.
