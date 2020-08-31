---
layout: "junos"
page_title: "Junos: junos_system_radius_server"
sidebar_current: "docs-junos-resource-system-radius-server"
description: |-
  Configure a system radius-server
---

# junos_system_radius_server

Configure a system radius-server.

## Example Usage

```hcl
# Add a system radius-server
resource junos_system_radius_server "demo_radius_server" {
  address = "192.0.2.1"
  secret  = "password"
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Required, Forces new resource)(`String`) RADIUS server address.
* `secret` - (Required)(`String`) Shared secret with the RADIUS server.
**WARNING** Clear in tfstate.
* `preauthentication_secret` - (Optional)(`String`) Preauthentication shared secret with the RADIUS server.
**WARNING** Clear in tfstate.
* `source_address` - (Optional)(`String`) Use specified address as source address.
* `port` - (Optional)(`Int`) RADIUS server authentication port number (1..65535).
* `accounting_port` - (Optional)(`Int`) RADIUS server accounting port number (1..65535).
* `dynamic_request_port` - (Optional)(`Int`) RADIUS client dynamic request port number (1..65535).
* `preauthentication_port` - (Optional)(`Int`) RADIUS server preauthentication port number (1..65535).
* `timeout` - (Optional)(`Int`) Request timeout period (1..1000 seconds).
* `accouting_timeout` - (Optional)(`Int`) Accounting request timeout period (0..1000 seconds).
* `retry` - (Optional)(`Int`) Retry attempts (1..100)
* `accounting_retry` - (Optional)(`Int`) Accounting retry attempts (0..100).
* `max_outstanding_requests` - (Optional)(`Int`) Maximum requests in flight to server (0..2000).
* `routing_instance` - (Optional)(`String`) Routing instance.

## Import

Junos system radius-server can be imported using an id made up of `<address>`, e.g.

```
$ terraform import junos_system_radius_server.demo_radius_server 192.0.2.1
```
