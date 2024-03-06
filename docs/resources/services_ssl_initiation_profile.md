---
page_title: "Junos: junos_services_ssl_initiation_profile"
---

# junos_services_ssl_initiation_profile

Provides a services ssl initiation profile

## Example Usage

```hcl
# Add a services services ssl initiation profile
resource "junos_services_ssl_initiation_profile" "demo" {
  name              = "demo"
  preferred_ciphers = "medium"
  protocol_version  = "tls12"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Profile name (Profile identifier).
- **actions** (Optional, Block)  
  Declare `actions` configuration.
  - **crl_disable** (Optional, Boolean)  
    Disable CRL validation.
  - **crl_if_not_present** (Optional, String)  
    Action if CRL information is not present.  
    Need to be `allow` or `drop`.
  - **crl_ignore_hold_instruction_code** (Optional, Boolean)  
    Ignore 'Hold Instruction Code' present in the CRL entry.
  - **ignore_server_auth_failure** (Optional, Boolean)  
    Ignore server authentication failure.
- **client_certificate** (Optional, String)  
  Local certificate identifier.
- **custom_ciphers** (Optional, Set of String)  
  Custom cipher list.
- **enable_flow_tracing** (Optional, Boolean)  
  Enable flow tracing for the profile.
- **enable_session_cache** (Optional, Boolean)  
  Enable SSL session cache.
- **preferred_ciphers** (Optional, String)  
  Select preferred ciphers.  
  Need to be `custom`, `medium`, `strong` or `weak`.
- **protocol_version** (Optional, String)  
  Protocol SSL version accepted.
- **trusted_ca** (Optional, Set of String)  
  List of trusted certificate authority profiles.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos services ssl initiation profile can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_ssl_initiation_profile.demo demo
```
