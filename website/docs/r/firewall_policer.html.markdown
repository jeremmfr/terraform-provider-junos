---
layout: "junos"
page_title: "Junos: junos_firewall_policer"
sidebar_current: "docs-junos-resource-firewall-policer"
description: |-
  Create firewall policer
---

# junos_firewall_policer

Provides a firewall policer resource.

## Example Usage

```hcl
# Configure a firewall policer
resource junos_firewall_policer "policer_demo" {
  name            = "policerDemo"
  filter_specific = true
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit  = "50k"
  }
  then {
    discard = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Name of policer.
* `filter_specific` - (Optional)(`Bool`) Policer is filter-specific.
* `if_exceeding` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for define rate limits options.
  * `bandwidth_percent` - (Optional)(`Int`) Bandwidth limit in percentage.
  * `bandwidth_limit` - (Optional)(`String`) Bandwidth limit in bits/second.
  * `burst_size_limit` - (Required)(`String`) Burst size limit in bytes.
* `then` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for define action to take if the rate limits are exceeded.
  * `discard` - (Optional)(`Bool`) Discard the packet.
  * `forwarding_class` - (Optional)(`String`) Classify packet to forwarding class.
  * `loss_priority` - (Optional)(`String`) Packet's loss priority.
  * `out_of_profile` - (Optional)(`Bool`)  Discard packets only if both congested and over threshold.

## Import

Junos firewall policer can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_firewall_policer.policer_demo policerDemo
```
