---
page_title: "Junos: junos_security_utm_profile_web_filtering_websense_redirect"
---

# junos_security_utm_profile_web_filtering_websense_redirect

Provides a security utm feature-profile web-filtering websense-redirect profile resource.

## Example Usage

```hcl
# Add a security utm feature-profile web-filtering websense-redirect profile
resource "junos_security_utm_profile_web_filtering_websense_redirect" "demo_profile" {
  name                 = "Default Webfilter3"
  custom_block_message = "Blocked by Juniper"
  server {
    host = "10.0.0.1"
    port = 1024
  }
  timeout = 3
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of security utm feature-profile web-filtering websense-redirect profile.
- **account** (Optional, String)  
  Set websense redirect account.
- **custom_block_message** (Optional, String)  
  Custom block message sent to HTTP client.
- **fallback_settings** (Optional, Block)  
  Configure fallback settings.
  - **default** (Optional, String)  
    Default action.  
    Need to be `block` or `log-and-permit`.
  - **server_connectivity** (Optional, String)  
    Action when device cannot connect to server.  
    Need to be `block` or `log-and-permit`.
  - **timeout** (Optional, String)  
    Action when connection to server timeout.  
    Need to be `block` or `log-and-permit`.
  - **too_many_requests** (Optional, String)  
    Action when requests exceed the limit of engine.  
    Need to be `block` or `log-and-permit`.
- **no_safe_search** (Optional, Boolean)  
  Do not perform safe-search for websense redirect protocol.
- **server** (Optional, Block)  
  Configure server settings.
  - **host** (Optional, String)  
    Server host IP address or string host name.
  - **port** (Optional, Number)  
    Server port (1024..65535).
  - **routing_instance** (Optional, String)  
    Routing instance used to connect server.
  - **source_address** (Optional, String)  
    Source ip address used to connect server.
- **socket** (Optional, Number)  
  Set sockets number (1..32).
- **timeout** (Optional, Number)  
  Set timeout (1..1800).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security utm feature-profile web-filtering websense-redirect profile can be imported using an
id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_profile_web_filtering_websense_redirect.demo_profile "Default Webfilter3"
```
