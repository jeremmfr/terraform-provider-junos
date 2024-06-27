---
page_title: "Junos: junos_routing_instance"
---

# junos_routing_instance

Provides a routing instance resource.

## Example Usage

```hcl
# Add a routing instance
resource "junos_routing_instance" "demo_ri" {
  name = "prod-vr"
}
```

## Argument Reference

-> **Note:** The interfaces can be configured with the `junos_interface_logical` resource and
the `routing_instance` argument.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of routing instance.
- **configure_rd_vrfopts_singly** (Optional, Boolean, Forces new resource)  
  Configure `route-distinguisher` and `vrf-*` options in other resource (like `junos_evpn`).
- **configure_type_singly** (Optional, Boolean, Forces new resource)  
  Configure `instance-type` option in other resource (like `junos_evpn`).  
  `type` argument need to be set to "" when true (to avoid confusion).
- **type** (Optional, String)  
  Type of routing instance.  
  Defaults to `virtual-router`.
- **as** (Optional, String)  
  Autonomous system number in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format.
- **description** (Optional, String)  
  Text description of routing instance.
- **instance_export** (Optional, List of String)  
  Export policy for instance RIBs.
- **instance_import** (Optional, List of String)  
  Import policy for instance RIBs.
- **remote_vtep_list** (Optional, Set of String)  
  Configure static remote VXLAN tunnel endpoints.
- **remote_vtep_v6_list** (Optional, Set of String)  
  Configure static ipv6 remote VXLAN tunnel endpoints.
- **route_distinguisher** (Optional, String)  
  Route distinguisher for this instance.
- **router_id** (Optional, String)  
  Router identifier.
- **vrf_export** (Optional, List of String)  
  Export policy for VRF instance RIBs.
- **vrf_import** (Optional, List of String)  
  Import policy for VRF instance RIBs.
- **vrf_target** (Optional, String)  
  Target community to use in import and export.
- **vrf_target_auto** (Optional, Boolean)  
  Auto derive import and export target community from BGP AS & L2.
- **vrf_target_export** (Optional, String)  
  Target community to use when marking routes on export.
- **vrf_target_import** (Optional, String)  
  Target community to use when filtering on import.
- **vtep_source_interface** (Optional, String)  
  Source layer-3 IFL for VXLAN.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos routing instance can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_routing_instance.demo_ri prod-vr
```
