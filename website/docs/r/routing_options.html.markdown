---
layout: "junos"
page_title: "Junos: junos_routing_options"
sidebar_current: "docs-junos-resource-routing-options"
description: |-
  Configure static configuration in routing-options block
---

# junos_routing_options

-> **Note:** This resource should only create **once**. It's used to configure static (not object) options in `routing-options` block. Destroy this resource as no effect on Junos configuration.

Configure static configuration in `routing-options` block

## Example Usage

```hcl
# Configure routing-options
resource junos_routing_options "routing_options" {
  autonomous_system {
    number = "65000"
  }
  graceful_restart {}
}
```

## Argument Reference

The following arguments are supported:

* `autonomous_system` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'autonomous-system' configuration.
  * `number` - (Required)(`String`) Autonomous system number in plain number or 'higher 16bits'.'Lower 16 bits' (asdot notation) format.
  * `asdot_notation` - (Optional)(`Bool`) Use AS-Dot notation to display true 4 byte AS numbers.
  * `loops` - (Optional)(`Int`) Maximum number of times this AS can be in an AS path (1..10).
* `graceful_restart` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'graceful-restart' configuration.
  * `disable` - (Optional)(`Bool`) Disable graceful restart.
  * `restart_duration` - (Optional)(`Int`) Maximum time for which router is in graceful restart (120..10000).

## Import

Junos routing_options can be imported using any id, e.g.

```
$ terraform import junos_routing_options.routing_options random
```
