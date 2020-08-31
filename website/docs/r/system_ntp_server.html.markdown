---
layout: "junos"
page_title: "Junos: junos_system_ntp_server"
sidebar_current: "docs-junos-resource-system-ntp-server"
description: |-
  Configure a system ntp server
---

# junos_system_ntp_server

Configure a system ntp server.

## Example Usage

```hcl
# Add a system ntp server
resource junos_system_ntp_server "demo_ntp_server" {
  address = "192.0.2.1"
  prefer  = true
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Required, Forces new resource)(`String`) The name or address of server.
* `key` - (Optional)(`Int`) Authentication key (1..65534).
* `prefer` - (Optional)(`Bool`) Prefer this peer_serv.
* `routing_instance` - (Optional)(`String`) Routing instance through which server is reachable.
* `version` - (Optional)(`Int`) NTP version to use (1..4).

## Import

Junos system ntp server can be imported using an id made up of `<address>`, e.g.

```
$ terraform import junos_system_ntp_server.demo_ntp_server 192.0.2.1
```
