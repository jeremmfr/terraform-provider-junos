---
layout: "junos"
page_title: "Junos: junos_security_dynamic_address_feed_server"
sidebar_current: "docs-junos-resource-security-dynamic-address-feed-server"
description: |-
  Create a security dynamic-address feed-server (when Junos device supports it)
---

# junos_security_dynamic_address_feed_server

Provides a security dynamic-address feed-server resource.

## Example Usage

```hcl
# Add a security dynamic-address feed-server
resource "junos_security_dynamic_address_feed_server" "demo_feed_srv" {
  name        = "demo"
  hostname    = "example.com"
  description = "demo junos_security_dynamic_address_feed_server"
  feed_name {
    name        = "bad_ips"
    path        = "/srx/"
    description = "demo feed_name junos_security_dynamic_address_feed_server"
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Security dynamic address feed-server name.
- **hostname** (Required, String)  
  Hostname or IP address of feed-server
- **description** (Optional, String)  
  Text description of feed-server.
- **feed_name** (Optional, Block List)  
  For each feed name.
  - **name** (Required, String)  
    Security dynamic address feed name in feed-server.
  - **path** (Required, String)  
    Path of feed, appended to feed-server to form a complete URL.
  - **description** (Optional, String)  
    Text description of feed in feed-server.
  - **hold_interval** (Optional, Number)  
    Time to keep IP entry when update failed (0..4294967295 seconds).
  - **update_interval** (Optional, Number)  
    Interval to retrieve update (30..4294967295 seconds).
- **hold_interval** (Optional, Number)  
  Time to keep IP entry when update failed (0..4294967295 seconds).
- **update_interval** (Optional, Number)  
  Interval to retrieve update (30..4294967295 seconds).

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `<name>`.

## Import

Junos security dynamic-address feed-server can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_dynamic_address_feed_server.demo_feed_srv demo
```
