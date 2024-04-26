---
page_title: "Junos: junos_evpn"
---

# junos_evpn

-> **Note:** This resource should only be created **once** for root level or each routing-instance.
It's used to configure static (not object) options in `protocols evpn` block in root or
routing-instance level, in `switch-options` and potentially also in `routing-instance` level directly.

Configure static configuration in `protocols evpn` block for root ou routing-instance level and the
various options potentially required in `switch-options` or on the `routing-instance` in same commit.

## Example Usage

```hcl
# Configure evpn
resource "junos_evpn" "default" {
  encapsulation = "vxlan"
  switch_or_ri_options {
    route_distinguisher = "20:1"
    vrf_target          = "target:20:2"
  }
}
```

## Argument Reference

The following arguments are supported:

- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance.  
  Need to be `default` (for root level) or the name of routing instance.  
  Defaults to `default`.
- **routing_instance_evpn** (Optional, Boolean, Forces new resource)  
  Configure routing instance is an evpn instance-type.
- **encapsulation** (Required, String)  
  Encapsulation type for EVPN.  
  Need to be `mpls` or `vxlan`.
- **default_gateway** (Optional, String)  
  Default gateway mode.  
  Need to be `advertise`, `do-not-advertise` or `no-gateway-community`.
- **duplicate_mac_detection** (Optional, Block)  
  Duplicate MAC detection settings.  
  An attribute of block need to be set.
  - **auto_recovery_time** (Optional, Number)  
    Automatically unblock duplicate MACs after a time delay (1..360 minutes).
  - **detection_threshold** (Optional, Number)  
    Number of moves to trigger duplicate MAC detection (2..20).
  - **detection_window** (Optional, Number)  
    Time window for detection of duplicate MACs (5..600 seconds).
- **multicast_mode** (Optional, String)  
  Multicast mode for EVPN.
- **no_core_isolation** (Optional, Boolean)  
  Disable EVPN Core isolation.  
  `routing_instance` need to be `default`.
- **switch_or_ri_options** (Optional, Block, Forces new resource)  
  Declare `switch-options` or `routing-instance` configuration.  
  Need to be set if `routing_instance` = `default` or `routing_instance_evpn` = true.  
  To avoid conflict with same options on `routing_instance` resource, add him `configure_rd_vrfopts_singly`.
  - **route_distinguisher** (Required, String)  
    Route distinguisher for this instance.
  - **vrf_export** (Optional, List of String)  
    Export policy for VRF instance RIBs.
  - **vrf_import** (Optional, List of String)  
    Import policy for VRF instance RIBs.
  - **vrf_target** (Optional, String)  
    VRF target community configuration.
  - **vrf_target_auto** (Optional, Boolean)  
    Auto derive import and export target community from BGP AS & L2.
  - **vrf_target_export** (Optional, String)  
    Target community to use when marking routes on export.
  - **vrf_target_import** (Optional, String)  
    Target community to use when filtering on import.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<routing_instance>`.

## Import

Junos evpn can be imported using an id made up of `<routing_instance>`, e.g.

```shell
$ terraform import junos_evpn.default default
```

If `routing_instance` != `default`, `switch_or_ri_options` is not imported.  
Add the internal delimiter and a random word to import it, e.g.

```shell
$ terraform import junos_evpn.ri ri_name_-_random
```
