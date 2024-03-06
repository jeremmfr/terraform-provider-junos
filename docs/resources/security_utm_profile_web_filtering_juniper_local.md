---
page_title: "Junos: junos_security_utm_profile_web_filtering_juniper_local"
---

# junos_security_utm_profile_web_filtering_juniper_local

Provides a security utm feature-profile web-filtering juniper-local profile resource.

## Example Usage

```hcl
# Add a security utm feature-profile web-filtering juniper-local profile
resource "junos_security_utm_profile_web_filtering_juniper_local" "demo_profile" {
  name                 = "Default Webfilter2"
  default_action       = "log-and-permit"
  custom_block_message = "Blocked by Juniper"
  timeout              = 3
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of security utm feature-profile web-filtering juniper-local profile.
- **custom_block_message** (Optional, String)  
  Custom block message sent to HTTP client.
- **default_action** (Optional, String)  
  Default action.  
  Need to be `block`, `log-and-permit` or `permit`.
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
- **timeout** (Optional, Number)  
  Set timeout.  
  Need to be between 1 and 1800.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security utm feature-profile web-filtering juniper-local profile can be imported using an
id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_profile_web_filtering_juniper_local.demo_profile "Default Webfilter2"
```
