---
layout: "junos"
page_title: "Junos: junos_evpn"
sidebar_current: "docs-junos-resource-evpn"
description: |-
  Configure static configuration in protocols evpn block in root or routing-instance level and options potentially required in switch-options or on routing-instance.
---

# junos_evpn

-> **Note:** This resource should only be created **once** for root level or each routing-instance. It's used to configure static (not object) options in `protocols evpn` block in root or routing-instance level, in `switch-options` and potentially also in `routing-instance` level directly.

Configure static configuration in `protocols evpn` block for root ou routing-instance level and the various options potentially required in `switch-options` or on the `routing-instance` in same commit.

## Example Usage

```hcl
# Configure evpn
resource junos_evpn "default" {
  encapsulation = "vxlan"
  switch_or_ri_options {
    route_distinguisher = "20:1"
    vrf_target          = "target:20:2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `routing_instance` - (Optional)(`String`) Routing instance. Need to be 'default' (for root level) or the name of routing instance. Defaults to `default`.
* `routing_instance_evpn` - (Optional)(`Bool`) Configure routing-instance is an evpn instance-type.
* `encapsulation` - (Required)(`String`) Encapsulation type for EVPN. Need to be `mpls` or `vxlan`.
* `multicast_mode` - (Optional)(`String`) Multicast mode for EVPN.
* `switch_or_ri_options` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare `switch-options` or `routing-instance` configuration. Need to be set if `routing_instance` = `default` or `routing_instance_evpn` = true. To avoid conflict with same options on `routing_instance` resource, add him `configure_rd_vrfopts_singly`.
  * `route_distinguisher` - (Required)(`String`) Route distinguisher for this instance.
  * `vrf_export` - (Optional)(`ListOfString`) Export policy for VRF instance RIBs.
  * `vrf_import` - (Optional)(`ListOfString`) Import policy for VRF instance RIBs.
  * `vrf_target` - (Optional)(`String`) VRF target community configuration.
  * `vrf_target_auto` - (Optional)(`Bool`) Auto derive import and export target community from BGP AS & L2.
  * `vrf_target_export` - (Optional)(`String`) Target community to use when marking routes on export.
  * `vrf_target_import` - (Optional)(`String`) Target community to use when filtering on import.

## Import

Junos evpn can be imported using an id made up of `<routing_instance>`, e.g.

```shell
$ terraform import junos_evpn.default default
```

If `routing_instance` != `default`, `switch_or_ri_options` is not imported. Add the internal delimiter and a random word to import it, e.g.

```shell
$ terraform import junos_evpn.ri ri_name_-_random
```
