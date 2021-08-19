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

- **area_id** (Required, String, Forces new resource)  
  The id of ospf area.
- **routing_instance** (Optional, String)  
  Routing instance for area.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`.
- **version** (Optional, String)  
  Version of ospf.  
  Need to be `v2` or `v3`.  
  Defaults to `v2`.
- **interface** (Required, Block List)  
  For each interface or interface-range to declare.
  - **name** (Required, String)  
    Name of interface or interface-range.
  - **dead_interval** (Optional, Number)  
    Dead interval (seconds).
  - **disable** (Optional, Boolean)  
    Disable OSPF on this interface.
  - **hello_interval** (Optional, Number)  
    Hello interval (seconds).
  - **metric** (Optional, Number)  
    Interface metric.
  - **passive** (Optional, Boolean)  
    Do not run OSPF, but advertise it.
  - **retransmit_interval** (Optional, Number)  
    Retransmission interval (seconds).

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<aread_id>_-_<version>_-_<routing_instance>`.

## Import

Junos ospf area can be imported using an id made up of
`<aread_id>_-_<version>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_ospf_area.demo_area 0.0.0.0_-_v2_-_default
```
