---
layout: "junos"
page_title: "Junos: junos_security_idp_policy"
sidebar_current: "docs-junos-resource-security-idp-policy"
description: |-
  Create a security idp policy (when Junos device supports it)
---

# junos_security_idp_policy

Provides a security idp policy resource.

## Example Usage

```hcl
# Add a idp policy
resource junos_security_idp_policy "demo_idp_policy" {
  name = "Idp-Policy"
  ips_rule {
    name        = "rules_1"
    description = "rules n1"
    match {
      application         = "junos:telnet"
      destination_address = ["192.0.2.0/24"]
    }
    then {
      action    = "drop-connection"
      ip_action = "ip-close"
      severity  = "info"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of idp policy.
* `exempt_rule` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each 'rulebase-exempt rule' to declare.
  * `name` - (Required)(`String`) The name of the rulebase-exempt rule.
  * `match` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'match' configuration. See the [`match` arguments for exempt_rule and ips_rule] (#match-arguments-for-exempt_rule-and-ips_rule) block but without `application` argument.
  * `description` - (Optional)(`String`) Rule description.
* `ips_rule` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each 'rulebase-ips rule' to declare.
  * `name` - (Required)(`String`) The name of the rulebase-ips rule.
  * `match` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'match' configuration. See the [`match` arguments for exempt_rule and ips_rule] (#match-arguments-for-exempt_rule-and-ips_rule) block.
  * `then` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'then' configuration. See the [`then` arguments] (#then-arguments) block.
  * `description` - (Optional)(`String`) Rule description.
  * `terminal` - (Optional)(`Bool`) Set/Unset terminal flag.

---

### match arguments for exempt_rule and ips_rule

* `application` - (Optional)(`String`) Specify application or application-set name to match. Only with `ips_rule`.
* `custom_attack_group` - (Optional)(`ListOfString`) Match custom attack groups.
* `custom_attack` - (Optional)(`ListOfString`) Match custom attacks.
* `destination_address` - (Optional)(`ListOfString`) Match destination address.
* `destination_address_except` - (Optional)(`ListOfString`) Don't match destination address.
* `dynamic_attack_group` - (Optional)(`ListOfString`) Match dynamic attack groups.
* `from_zone` - (Optional)(`String`) Match from zone.
* `predefined_attack_group` - (Optional)(`ListOfString`) Match predefined attack groups.
* `predefined_attack` - (Optional)(`ListOfString`) Match predefined attacks.
* `source_address` - (Optional)(`ListOfString`) Match source address.
* `source_address_except` - (Optional)(`ListOfString`) Don't match source address.
* `to_zone` - (Optional)(`String`) Match to zone.

---

### then arguments

* `action` - (Required)(`String`) Action. Need to be 'class-of-service', 'close-client', 'close-client-and-server', 'close-server', 'drop-connection', 'drop-packet', 'ignore-connection', 'mark-diffserv', 'no-action' or 'recommended'.
* `action_cos_forwarding_class` - (Optional)(`String`) Forwarding class for outgoing packets. `action` need to be 'class-of-service'.
* `action_dscp_code_point` - (Optional)(`Int`) Codepoint value (0..63). `action` need to be 'class-of-service' or 'mark-diffserv'.
* `ip_action` - (Optional)(`String`) IP-action. Need to be 'ip-block', 'ip-close' or 'ip-notify'.
* `ip_action_log` - (Optional)(`Bool`) Log IP action taken.
* `ip_action_log_create` - (Optional)(`Bool`) Log IP action creation.
* `ip_action_refresh_timeout` - (Optional)(`Bool`) Refresh timeout when future connections match installed ip-action filter.
* `ip_action_target` - (Optional)(`String`) IP-action target. Need to be 'destination-address', 'service', 'source-address', 'source-zone', 'source-zone-address' or 'zone-service'.
* `ip_action_timeout` - (Optional)(`Int`) Number of seconds IP action should remain effective (0..64800).
* `notification` - (Optional)(`Bool`) Configure notification.
* `notification_log_attacks` - (Optional)(`Bool`) Enable attack logging.
* `notification_log_attacks_alert` - (Optional)(`Bool`) Set alert flag in attack log.
* `notification_packet_log` - (Optional)(`Bool`) Enable packet-log.
* `notification_packet_log_post_attack` - (Optional)(`Int`) No of packets to capture after attack (0..255).
* `notification_packet_log_post_attack_timeout` - (Optional)(`Int`) Timeout (seconds) after attack before stopping packet capture (0..1800).
* `notification_packet_log_pre_attack` - (Optional)(`Int`) No of packets to capture before attack (1..255).
* `severity` - (Optional)(`String`) Set rule severity level. Need to be 'critical', 'info', 'major', 'minor' or 'warning'.

## Import

Junos security idp policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_idp_policy.demo_idp_policy Idp-Policy
```
