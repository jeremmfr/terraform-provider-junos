---
layout: "junos"
page_title: "Junos: junos_security_ipsec_vpn"
sidebar_current: "docs-junos-resource-security-ipsec-vpn"
description: |-
  Create a security ipsec vpn (when Junos device supports it)
---

# junos_security_ipsec_vpn

Provides a security ipsec vpn resource.

## Example Usage

```hcl
# Add a route-based ipsec vpn
resource junos_security_ipsec_vpn "demo_vpn" {
  name              = "first-vpn"
  establish_tunnels = "immediately"
  bind_interface    = junos_interface_st0_unit.demo.id
  ike {
    gateway = "ike-gateway"
    policy  = "ipsec-policy"
  }
}
resource junos_interface_st0_unit demo {}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of vpn.
- **bind_interface** (Optional, Computed, String)  
  Interface st0 to bind vpn for route-based vpn.  
  Computed when `bind_interface_auto` = true.
- **bind_interface_auto** (Optional, Boolean, **Deprecated**)  
  Find st0 available for compute bind_interface automatically.  
  Deprecated argument, use the `junos_interface_st0_unit` resource to find st0 unit available instead.
- **df_bit** (Optional, String)  
  Specifies how to handle the Don't Fragment bit.  
  Need to be `clear`, `copy` or `set`.
- **establish_tunnels** (Optional, String)  
  When the VPN comes up.  
  Need to be `immediately` or `on-traffic`.
- **ike** (Required, Block)  
  Declare ike configuration.
  - **gateway** (Required, String)  
    The name of security ike gateway (phase-1).
  - **policy** (Required, String)  
    The name of ipsec policy.
  - **identity_local** (Optional, String)  
    IPSec proxy-id local parameter.
  - **identity_remote** (Optional, String)  
    IPSec proxy-id remote parameter.
  - **identity_service** (Optional, String)  
    IPSec proxy-id service parameter.
- **traffic_selector** (Optional, Block List)  
  For each name of traffic-selector to declare.
  - **name** (Required, String)  
    Name of traffic-selector.
  - **local_ip** (Required, String)  
    CIDR for IP addresses of local traffic-selector.
  - **remote_ip** (Required, String)  
    CIDR for IP addresses of remote traffic-selector.
- **vpn_monitor** (Optional, Block)  
  Declare VPN monitor liveness configuration.
  - **destination_ip** (Optional, String)  
    IP destination for monitor message.
  - **optimized** (Optional, Boolean)  
    Optimize for scalability.
  - **source_interface** (Optional, Computed, String)  
    Set source interface for monitor message.  
    Compute when `source_interface_auto` = true.
  - **source_interface_auto** (Optional, Boolean)  
    Compute the source_interface to `bind_interface`.

## Import

Junos security ipsec vpn can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_ipsec_vpn.demo_vpn first-vpn
```
