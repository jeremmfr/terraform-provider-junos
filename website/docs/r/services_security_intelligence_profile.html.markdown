---
layout: "junos"
page_title: "Junos: junos_services_security_intelligence_profile"
sidebar_current: "docs-junos-resource-services-security-intelligence-profile"
description: |-
  Create a services security-intelligence profile
---

# junos_services_security_intelligence_profile

Provides a services security-intelligence profile resource.

## Example Usage

```hcl
# Add a services security-intelligence profile
resource "junos_services_security_intelligence_profile" "demo" {
  name     = "demo"
  category = "CC"
  rule {
    name = "rule_1"
    match {
      threat_level = [10]
    }
    then_action = "block close http redirect-url http://www.test.com/url1.html"
    then_log    = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Security intelligence profile name.
* `category` - (Required)(`String`) Profile category name.
* `rule` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure a rule. Can be specified multiple times for each rule name. See the [`rule` arguments] (#rule-arguments) block.
* `default_rule_then` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare profile default rule.
  * `action` - (Required)(`String`) Security intelligence profile action. Need to be 'permit', 'recommended', 'block drop', 'block close' or 'block close http (file|message|redirect-url) ...'.
  * `log` - (Optional)(`Bool`) Log security intelligence block action.
  * `no_log` - (Optional)(`Bool`) Don't log security intelligence block action.
* `description` - (Optional)(`String`) Text description of profile.

---
#### rule arguments
* `name` - (Required)(`String`) Profile rule name.
* `match` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to configure profile matching feed name and threat levels.
  * `threat_level` - (Required)(`ListOfInt`) Profile matching threat levels, higher number is more severe (1..10).
  * `feed_name` - (Optional)(`ListOfString`) Profile matching feed name.
* `then_action` - (Required)(`String`) Security intelligence profile action. Need to be 'permit', 'recommended', 'block drop', 'block close' or 'block close http (file|message|redirect-url) ...'.
* `then_log` - (Optional)(`Bool`) Log security intelligence block action.

## Import

Junos services security-intelligence profile can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_services_security_intelligence_profile.demo demo
```
