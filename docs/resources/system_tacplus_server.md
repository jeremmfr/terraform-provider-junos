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
