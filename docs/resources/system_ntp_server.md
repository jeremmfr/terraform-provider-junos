---
page_title: "Junos: junos_system_ntp_server"
---

# junos_system_ntp_server

Configure a system ntp server.

## Example Usage

```hcl
# Add a system ntp server
resource "junos_system_ntp_server" "demo_ntp_server" {
  address = "192.0.2.1"
  prefer  = true
}
```

## Argument Reference

The following arguments are supported:

- **address** (Required, String, Forces new resource)  
  Address of server.
- **key** (Optional, Number)  
  Authentication key (1..65534).
- **prefer** (Optional, Boolean)  
  Prefer this peer_serv.
- **routing_instance** (Optional, String)  
  Routing instance through which server is reachable.
- **version** (Optional, Number)  
  NTP version to use (1..4).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<address>`.

## Import

Junos system ntp server can be imported using an id made up of `<address>`, e.g.

```shell
$ terraform import junos_system_ntp_server.demo_ntp_server 192.0.2.1
```
