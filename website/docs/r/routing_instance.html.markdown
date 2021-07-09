---
layout: "junos"
page_title: "Junos: junos_routing_instance"
sidebar_current: "docs-junos-resource-routing-instance"
description: |-
  Create a routing instance
---

# junos_routing_instance

Provides a routing instance resource.

## Example Usage

```hcl
# Add a routing instance
resource junos_routing_instance "demo_ri" {
  name = "prod-vr"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of routing instance.
* `configure_rd_vrfopts_singly` - (Optional, Forces new resource)(`Bool`) Configure `route-distinguisher` and `vrf-*` options in other resource (like `junos_evpn`).
* `configure_type_singly` - (Optional, Forces new resource)(`Bool`) Configure `instance-type` option in other resource (like `junos_evpn`). `type` argument need to be set to "" when true (to avoid confusion).
* `type` - (Optional)(`String`) Type of routing instance. Defaults to `virtual-router`.
* `as` - (Optional)(`String`) Autonomous system number in plain number or 'higher 16bits'.'Lower 16 bits' (asdot notation) format.
* `description` - (Optional)(`String`) Text description of routing instance.
* `route_distinguisher` - (Optional)(`String`) Route distinguisher for this instance.
* `vrf_export` - (Optional)(`ListOfString`) Export policy for VRF instance RIBs;
* `vrf_import` - (Optional)(`ListOfString`) Import policy for VRF instance RIBs.
* `vrf_target` - (Optional)(`String`) Target community to use in import and export.
* `vrf_target_auto` - (Optional)(`Bool`) Auto derive import and export target community from BGP AS & L2.
* `vrf_target_export` - (Optional)(`String`) Target community to use when marking routes on export.
* `vrf_target_import` - (Optional)(`String`) Target community to use when filtering on import.
* `vtep_source_interface` - (Optional)(`String`) Source layer-3 IFL for VXLAN.

## Import

Junos routing instance can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_routing_instance.demo_ri prod-vr
```
