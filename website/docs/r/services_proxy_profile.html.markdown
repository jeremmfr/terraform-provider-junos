---
layout: "junos"
page_title: "Junos: junos_services_proxy_profile"
sidebar_current: "docs-junos-resource-services-proxy-profile"
description: |-
  Create a services proxy profile
---

# junos_services_proxy_profile

Provides a services proxy profile resource.

## Example Usage

```hcl
# Add a services proxy profile
resource "junos_services_proxy_profile" "demo" {
  name               = "demo"
  protocol_http_host = "192.0.2.1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Proxy profile name.
* `protocol_http_host` - (Required)(`String`) Proxy server name or IP address.
* `protocol_http_port` - (Optional)(`Int`) Proxy server port (1..65535).

## Import

Junos services proxy profile can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_proxy_profile.demo demo
```
