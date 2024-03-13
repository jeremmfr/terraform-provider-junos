---
page_title: "Junos: junos_services_proxy_profile"
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

- **name** (Required, String, Forces new resource)  
  Proxy profile name.
- **protocol_http_host** (Required, String)  
  Proxy server name or IP address.
- **protocol_http_port** (Optional, Number)  
  Proxy server port (1..65535).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos services proxy profile can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_proxy_profile.demo demo
```
