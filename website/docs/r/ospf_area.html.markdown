---
layout: "junos"
page_title: "Junos: junos_ospf_area"
sidebar_current: "docs-junos-resource-ospf-area"
description: |-
  Create a ospf area
---

# junos_ospf_area

Provides a ospf area resource.

## Example Usage

```hcl
# Add a ospf area
resource junos_ospf_area "demo_area" {
  area_id = "0.0.0.0"
  interface {
    name = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `area_id` - (Required, Forces new resource)(`String`) The id of ospf area.
* `routing_instance` - (Optional)(`String`) Routing instance for area. Need to be 'default' or name of routing instance. Defaults to `default`.
* `version` - (Optional)(`String`) Version of ospf. Need to be 'v2' or 'v3'. Defaults to `v2`.
* `interface` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each interface or interface-range to declare.
  * `name` - (Required)(`String`) Name of interface or interface-range.
  * `dead_interval` - (Optional)(`Int`) Dead interval (seconds).
  * `disable` - (Optional)(`Bool`) Disable OSPF on this interface.
  * `hello_interval` - (Optional)(`Int`) Hello interval (seconds).
  * `metric` - (Optional)(`Int`) Interface metric.
  * `passive` - (Optional)(`Bool`) Do not run OSPF, but advertise it.
  * `retransmit_interval` - (Optional)(`Int`) Retransmission interval (seconds).

## Import

Junos ospf area can be imported using an id made up of `<aread_id>_-_<version>_-_<routing_instance>`, e.g.

```
$ terraform import junos_ospf_area.demo_area 0.0.0.0_-_v2_-_default
```
