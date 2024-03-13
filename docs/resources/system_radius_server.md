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

- **address** (Required, String, Forces new resource)  
  RADIUS server address.
- **secret** (Required, String, Sensitive)  
  Shared secret with the RADIUS server.
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
