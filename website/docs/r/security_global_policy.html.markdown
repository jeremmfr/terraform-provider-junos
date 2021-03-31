---
layout: "junos"
page_title: "Junos: junos_security_global_policy"
sidebar_current: "docs-junos-resource-security-global-policy"
description: |-
  Configure static configuration in security policies global block
---

# junos_security_global_policy

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `security policies global` block.

Configure static configuration in `security policies global` block

## Example Usage

```hcl
# Configure security policies global
resource junos_security_global_policy "global" {
  policy {
    name                      = "test"
    match_source_address      = ["blue"]
    match_destination_address = ["green"]
    match_application         = ["any"]
    match_from_zone           = ["any"]
    match_to_zone             = ["any"]
  }
  policy {
    name                      = "drop"
    match_source_address      = ["blue"]
    match_destination_address = ["any"]
    match_application         = ["any"]
    match_from_zone           = ["any"]
    match_to_zone             = ["any"]
    then                      = "deny"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each policy.
  * `name` - (Required)(`String`) Security policy name.
  * `match_source_address` - (Required)(`ListOfString`) List of source address match.
  * `match_destination_address` - (Required)(`ListOfString`) List of destination address match.
  * `match_application` - (Required)(`ListOfString`) List of applications match.
  * `match_from_zone` - (Required)(`ListOfString`) Match multiple source zone.
  * `match_to_zone` - (Required)(`ListOfString`) Match multiple destination zone.
  * `then` - (Optional)(`String`) Action of policy. Defaults to `permit`.
  * `count` - (Optional)(`Bool`) Enable count.
  * `log_init` - (Optional)(`Bool`) Log at session init time.
  * `log_close` - (Optional)(`Bool`) Log at session close time.
  * `match_destination_address_excluded` - (Optional)(`Bool`) Exclude destination addresses.
  * `match_dynamic_application` - (Optional)(`ListOfString`) List of dynamic application or group match.
  * `match_source_address_excluded` - (Optional)(`Bool`) Exclude source addresses.
  * `permit_application_services` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'permit application-services' configuration. See the [`permit_application_services` arguments] (#permit_application_services-arguments) block.

---
#### permit_application_services arguments
* `application_firewall_rule_set` - (Optional)(`String`) Servie rule-set Name for Application firewall.
* `application_traffic_control_rule_set` - (Optional)(`String`) Service rule-set name Application traffic control.
* `gprs_gtp_profile` - (Optional)(`String`) Specify GPRS Tunneling Protocol profile name.
* `gprs_sctp_profile` - (Optional)(`String`) Specify GPRS stream control protocol profile name.
* `idp` - (Optional)(`Bool`) Enable Intrusion detection and prevention.
* `redirect_wx` - (Optional)(`Bool`) Set WX redirection.
* `reverse_redirect_wx` - (Optional)(`Bool`) Set WX reverse redirection.
* `security_intelligence_policy` - (Optional)(`String`) Specify security-intelligence policy name.
* `ssl_proxy` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable SSL Proxy.
  * `profile_name` - (Optional)(`String`) Specify SSL proxy service profile name.
* `uac_policy` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable unified access control enforcement.
  * `captive_portal` - (Optional)(`String`) Specify captive portal.
* `utm_policy` - (Optional)(`String`) Specify utm policy name.

## Import

Junos security global policies can be imported using any id, e.g.

```
$ terraform import junos_security_global_policy.global random
```
