---
page_title: "Junos: junos_routing_instance"
---

# junos_routing_instance

Get information on a routing instance.

## Example Usage

```hcl
# Read routing instance
data junos_routing_instance "demo_ri" {
  name = "prod-vr"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  The name of routing instance.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.
- **type** (String)  
  Type of routing instance.  
- **as** (String)  
  Autonomous system number in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format.
- **description** (String)  
  Text description of routing instance.
- **instance_export** (List of String)  
  Export policy for instance RIBs
- **instance_import** (List of String)  
  Import policy for instance RIBs
- **interface** (Set of String)  
  List of interfaces in routing-instance.
- **route_distinguisher** (String)  
  Route distinguisher for this instance.
- **vrf_export** (List of String)  
  Export policy for VRF instance RIBs.
- **vrf_import** (List of String)  
  Import policy for VRF instance RIBs.
- **vrf_target** (String)  
  Target community to use in import and export.
- **vrf_target_auto** (Boolean)  
  Auto derive import and export target community from BGP AS & L2.
- **vrf_target_export** (String)  
  Target community to use when marking routes on export.
- **vrf_target_import** (String)  
  Target community to use when filtering on import.
- **vtep_source_interface** (String)  
  Source layer-3 IFL for VXLAN.
