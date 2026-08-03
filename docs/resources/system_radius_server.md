---
page_title: "Junos: junos_system_radius_server"
---

# junos_system_radius_server

Configure a system radius-server.

## Example Usage

```hcl
# Add a system radius-server
resource "junos_system_radius_server" "demo_radius_server" {
  address = "192.0.2.1"
  secret  = "password"
}
```

## Argument Reference

The following arguments are supported:

-> **Note**
  One of `secret` or `secret_wo` arguments is required.

- **address** (Required, String, Forces new resource)  
  RADIUS server address.
- **secret** (Optional, String, Sensitive)  
  Shared secret with the RADIUS server.  
  Conflict with `secret_wo`.
- **secret_wo** (Optional, String, Sensitive, Write-only)  
  Shared secret with the RADIUS server, not stored in state.  
  Requires `secret_wo_version` and Terraform 1.11 or later.  
  Conflict with `secret`.
- **secret_wo_version** (Optional, Number)  
  Version of `secret_wo` to trigger the sending of its value.  
  Increment it to send the current value of `secret_wo` to the device.  
  Requires `secret_wo`.
- **accounting_port** (Optional, Number)  
  RADIUS server accounting port number (1..65535).
- **accounting_retry** (Optional, Number)  
  Accounting retry attempts (0..100).
- **accounting_timeout** (Optional, Number)  
  Accounting request timeout period (0..1000 seconds).
- **dynamic_request_port** (Optional, Number)  
  RADIUS client dynamic request port number (1..65535).
- **max_outstanding_requests** (Optional, Number)  
  Maximum requests in flight to server (0..2000).
- **port** (Optional, Number)  
  RADIUS server authentication port number (1..65535).
- **preauthentication_port** (Optional, Number)  
  RADIUS server preauthentication port number (1..65535).
- **preauthentication_secret** (Optional, String, Sensitive)  
  Preauthentication shared secret with the RADIUS server.  
  Conflict with `preauthentication_secret_wo`.
- **preauthentication_secret_wo** (Optional, String, Sensitive, Write-only)  
  Preauthentication shared secret with the RADIUS server, not stored in state.  
  Requires `preauthentication_secret_wo_version` and Terraform 1.11 or later.  
  Conflict with `preauthentication_secret`.
- **preauthentication_secret_wo_version** (Optional, Number)  
  Version of `preauthentication_secret_wo` to trigger the sending of its value.  
  Increment it to send the current value of `preauthentication_secret_wo` to the device.  
  Requires `preauthentication_secret_wo`.
- **retry** (Optional, Number)  
  Retry attempts (1..100).
- **routing_instance** (Optional, String)  
  Routing instance.
- **source_address** (Optional, String)  
  Use specified address as source address.
- **timeout** (Optional, Number)  
  Request timeout period (1..1000 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<address>`.

## Import

Junos system radius-server can be imported using an id made up of `<address>`, e.g.

```shell
$ terraform import junos_system_radius_server.demo_radius_server 192.0.2.1
```

!> **Warning**
  Write-only arguments cannot be filled by an import, so the shared secrets read on the device
  are stored in `secret` and `preauthentication_secret`, and therefore in the Terraform state.  
  When the configuration uses the write-only arguments, the next apply removes them from the
  state.
