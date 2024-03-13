---
page_title: "Junos: junos_services_user_identification_ad_access_domain"
---

# junos_services_user_identification_ad_access_domain

Provides a services user-identification active-directory-access domain resource.

## Example Usage

```hcl
# Add a services user-identification active-directory-access domain
resource "junos_services_user_identification_ad_access_domain" "demo" {
  name          = "example.com"
  user_name     = "user_dom"
  user_password = "user_pass"
  domain_controller {
    name    = "server1"
    address = "192.0.2.3"
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Domain name.
- **user_name** (Required, String)  
  User name.
- **user_password** (Required, String, Sensitive)  
  Password string.  
- **domain_controller** (Optional, Block List)  
  For each name of domain-controller, configure address.
  - **name** (Required, String)  
    Domain controller name.
  - **address** (Required, String)  
    Address of domain controller.
- **ip_user_mapping_discovery_wmi** (Optional, Block)  
  Enable `ip-user-mapping discovery-method wmi`.
  - **event_log_scanning_interval** (Optional, Number)  
    Interval of event log scanning (5..60 seconds).
  - **initial_event_log_timespan** (Optional, Number)  
    Event log scanning timespan (1..168 hours).
- **user_group_mapping_ldap** (Optional, Block)  
  Declare `user-group-mapping ldap` configuration.
  - **base** (Required, String)  
    Base distinguished name.
  - **address** (Optional, List of String)  
    Address of LDAP server.
  - **auth_algo_simple** (Optional, Boolean)  
    Authentication-algorithm simple.
  - **ssl** (Optional, Boolean)  
    SSL.
  - **user_name** (Optional, String)  
    User name.
  - **user_password** (Optional, String, Sensitive)  
    Password string.  

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos services user-identification active-directory-access domain can be imported using an
id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_user_identification_ad_access_domain.demo example.com
```
