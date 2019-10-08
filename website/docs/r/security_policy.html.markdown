---
layout: "junos"
page_title: "Junos: junos_security_policy"
sidebar_current: "docs-junos-resource-security-policy"
description: |-
  Create a security policy (when Junos device supports it)
---

# junos_security_policy

Provides a security policy resource.

## Example Usage

```hcl
# Add a security policy
resource "junos_security_policy" "DemoPolicy" {
  from_zone = "trust"
  to_zone   = "untrust"
  policy {
    name                      = "allow_trust"
    match_source_address      = ["any"]
    match_destination_address = ["any"]
    match_application         = ["any"]
    then                      = "permit"
  }
}
```

## Argument Reference

The following arguments are supported:

* `from_zone` - (Required, Forces new resource)(`String`) The name of source zone.
* `to_zone` - (Required, Forces new resource)(`String`) The name of destination zone.
* `policy` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) List of policy with options. Can be specified multiple times for each policy.
  * `name`  - (Required)(`String`) The name of policy
  * `match_source_address` - (Required)(`ListOfString`) List of source address match
  * `match_destination_address` - (Required)(`ListOfString`) List of destination address match
  * `match_application` - (Required)(`ListOfString`) List of applications match
  * `then` - (Optional)(`String`) action of policy. Defaults to `permit`
  * `permit_tunnel_ipsec_vpn` - (Optional)(`String`) Name of vpn to permit with a tunnel ipsec
  * `count` - (Optional)(`Bool`) Enable count
  * `log_init` - (Optional)(`Bool`) Log at session init time
  * `log_close` - (Optional)(`Bool`) Log at session close time

## Import

Junos security policy can be imported using an id made up of `<from_zone>_-_<to_zone>`, e.g.

```
$ terraform import junos_security_zone.DemoPolicy trust_-_untrust
```
