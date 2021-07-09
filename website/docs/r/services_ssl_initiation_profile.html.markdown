---
layout: "junos"
page_title: "Junos: junos_services_ssl_initiation_profile"
sidebar_current: "docs-junos-resource-services-ssl-initiation-profile"
description: |-
  Create a services ssl initiation profile
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

* `name` - (Required, Forces new resource)(`String`) Profile name (Profile identifier).
* `actions` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'actions' configuration.
  * `crl_disable` - (Optional)(`Bool`) Disable CRL validation.
  * `crl_if_not_present` - (Optional)(`String`) Action if CRL information is not present. Need to be 'allow' or 'drop'.
  * `crl_ignore_hold_instruction_code` - (Optional)(`Bool`) Ignore 'Hold Instruction Code' present in the CRL entry.
  * `ignore_server_auth_failure` - (Optional)(`Bool`) Ignore server authentication failure.
* `client_certificate` - (Optional)(`String`) Local certificate identifier.
* `custom_ciphers` - (Optional)(`ListOfString`) Custom cipher list.
* `enable_flow_tracing` - (Optional)(`Bool`) Enable flow tracing for the profile.
* `enable_session_cache` - (Optional)(`Bool`) Enable SSL session cache.
* `preferred_ciphers` - (Optional)(`String`) Select preferred ciphers. Need to be 'custom', 'medium', 'strong' or 'weak'.
* `protocol_version` - (Optional)(`String`) Protocol SSL version accepted.
* `trusted_ca` - (Optional)(`ListOfString`) List of trusted certificate authority profiles.

## Import

Junos services ssl initiation profile can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_ssl_initiation_profile.demo demo
```
