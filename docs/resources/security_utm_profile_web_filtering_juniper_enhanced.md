---
page_title: "Junos: junos_security_utm_profile_web_filtering_juniper_enhanced"
---

# junos_security_utm_profile_web_filtering_juniper_enhanced

Provides a security utm feature-profile web-filtering juniper-enhanced profile resource.

## Example Usage

```hcl
# Add a security utm feature-profile web-filtering juniper-enhanced profile
resource "junos_security_utm_profile_web_filtering_juniper_enhanced" "demo_profile" {
  name = "Default Webfilter"
  category {
    name   = "Enhanced_Network_Errors"
    action = "block"
  }
  site_reputation_action {
    site_reputation = "very-safe"
    action          = "permit"
  }
  site_reputation_action {
    site_reputation = "harmful"
    action          = "block"
  }
  default_action       = "log-and-permit"
  custom_block_message = "Blocked by Juniper"
  timeout              = 3
  no_safe_search       = true
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of security utm feature-profile web-filtering juniper-enhanced profile.
- **block_message** (Optional, Block)  
  Configure block message.
  - **url** (Optional, String)  
    URL of block message.
  - **type_custom_redirect_url** (Optional, Boolean)  
    Enable Custom redirect URL server type.
- **category** (Optional, Block List)  
  For each name of category, configure enhanced category actions.  
  See [below for nested schema](#category-arguments).
- **custom_block_message** (Optional, String)  
  Custom block message sent to HTTP client.
- **default_action** (Optional, String)  
  Default action.  
  Need to be `block`, `log-and-permit`, `permit` or `quarantine`.
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
  Do not perform safe-search for Juniper enhanced protocol.
- **quarantine_custom_message** (Optional, String)  
  Quarantine custom message.
- **quarantine_message** (Optional, Block)  
  Configure quarantine message.
  - **url** (Optional, String)  
    URL of quarantine message.
  - **type_custom_redirect_url** (Optional, Boolean)  
    Enable Custom redirect URL server type.
- **site_reputation_action** (Optional, Block List)  
  For each site_reputation, configure action.
  - **site_reputation** (Required, String)  
    Level of reputation.  
    Need to be `fairly-safe`, `harmful`, `moderately-safe`, `suspicious`, `very-safe`.
  - **action** (Required, String)  
    Action for site-reputation.  
    Need to be `block`, `log-and-permit`, `permit` or `quarantine`.
- **timeout** (Optional, Number)  
  Set timeout. Need to be between 1 and 1800.

---

### category arguments

- **name** (Required, String)  
  Name of category.
- **action** (Required, String)  
  Action when web traffic matches category.  
  Need to be `block`, `log-and-permit`, `permit` or `quarantine`.
- **reputation_action** (Optional, Block List)  
  For each site_reputation, configure action.
  - **site_reputation** (Required, String)  
    Level of reputation.  
    Need to be `fairly-safe`, `harmful`, `moderately-safe`, `suspicious`, `very-safe`.
  - **action** (Required, String)  
    Action for site-reputation.  
    Need to be `block`, `log-and-permit`, `permit` or `quarantine`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security utm feature-profile web-filtering juniper-enhanced profile can be imported using an
id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_profile_web_filtering_juniper_enhanced.demo_profile "Default Webfilter"
```
