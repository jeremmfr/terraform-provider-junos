---
layout: "junos"
page_title: "Junos: junos_services_user_identification_ad_access_domain"
sidebar_current: "docs-junos-resource-services-user-identification-ad-access-domain"
description: |-
  Create a services user-identification active-directory-access domain
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

* `name` - (Required, Forces new resource)(`String`) Domain name.
* `user_name` - (Required)(`String`) User name.
* `user_password` - (Required)(`String`) Password string.  
**WARNING** Clear in tfstate.
* `domain_controller` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure a domain-controller. Can be specified multiple times for each controller.
  * `name` - (Required)(`String`) Domain controller name.
  * `address` - (Required)(`String`) Address of domain controller.
* `ip_user_mapping_discovery_wmi` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable 'ip-user-mapping discovery-method wmi'.
  * `event_log_scanning_interval` - (Optional)(`Int`) Interval of event log scanning (5..60 seconds).
  * `initial_event_log_timespan` - (Optional)(`Int`) Event log scanning timespan (1..168 hours).
* `user_group_mapping_ldap` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'user-group-mapping ldap' configuration.
  * `base` - (Required)(`String`) Base distinguished name.
  * `address` - (Optional)(`ListOfString`) Address of LDAP server.
  * `auth_algo_simple` - (Optional)(`Bool`) Authentication-algorithm simple.
  * `ssl` - (Optional)(`Bool`) SSL.
  * `user_name` - (Optional)(`String`) User name.
  * `user_password` - (Optional)(`String`) Password string.  
  **WARNING** Clear in tfstate.

## Import

Junos services user-identification active-directory-access domain can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_user_identification_ad_access_domain.demo example.com
```
