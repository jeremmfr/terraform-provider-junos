---
layout: "junos"
page_title: "Junos: junos_services_flowmonitoring_vipfix_template"
sidebar_current: "docs-junos-resource-services-flowmonitoring-vipfix-template"
description: |-
  Create a services flow-monitoring version-ipfix template
---

# junos_services_flowmonitoring_vipfix_template

Provides a services flow-monitoring version-ipfix template resource.

## Example Usage

```hcl
# Add a services flow-monitoring version-ipfix template
resource "junos_services_flowmonitoring_vipfix_template" "demo" {
  name = "demo"
  type = "ipv4-template"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Name of flow-monitoring ip-fix template.
* `type` - (Required)(`String`) Type of template. Need to be 'ipv4-template', 'ipv6-template' or 'mpls-template'.
* `flow_active_timeout` - (Optional)(`Int`) Interval after which active flow is exported (10..600).
* `flow_inactive_timeout` - (Optional)(`Int`) Period of inactivity that marks a flow inactive (10..600).
* `flow_key_flow_direction` - (Optional)(`Bool`) Include flow direction.
* `flow_key_vlan_id` - (Optional)(`Bool`) Include vlan ID.
* `ip_template_export_extension` - (Optional)(`ListOfString`) Export-extension for 'ipv4-template', 'ipv6-template' type.
* `nexthop_learning_enable` - (Optional)(`Bool`) Enable nexthop learning.
* `nexthop_learning_disable` - (Optional)(`Bool`) Disable nexthop learning.
* `observation_domain_id` - (Optional)(`Int`) Observation Domain Id (0..255).
* `option_refresh_rate` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'option-refresh-rate' configuration.
  * `packets` - (Optional)(`Int`) In number of packets (1..480000)
  * `seconds` - (Optional)(`Int`) In number of seconds (10..600)
* `option_template_id` - (Optional)(`Int`) Options template id (1024..65535).
* `template_id` - (Optional)(`Int`) Template id (1024..65535).
* `tunnel_observation_ipv4` - (Optional)(`Bool`) Tunnel observation IPv4.
* `tunnel_observation_ipv6` - (Optional)(`Bool`) Tunnel observation IPv6.

## Import

Junos services flow-monitoring version-ipfix template can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_services_flowmonitoring_vipfix_template.demo demo
```
