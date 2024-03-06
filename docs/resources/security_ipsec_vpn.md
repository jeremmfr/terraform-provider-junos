---
page_title: "Junos: junos_security_ipsec_vpn"
---

# junos_security_ipsec_vpn

Provides a security IPSec vpn resource.

## Example Usage

```hcl
# Add a route-based IPSec vpn
resource "junos_security_ipsec_vpn" "demo_vpn" {
  name              = "first-vpn"
  establish_tunnels = "immediately"
  bind_interface    = junos_interface_st0_unit.demo.id
  ike {
    gateway = "ike-gateway"
    policy  = "ipsec-policy"
  }
}
resource "junos_interface_st0_unit" "demo" {}
```

## Argument Reference

-> **Note:** One of `ike` or `manual` arguments is required.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of vpn.
- **bind_interface** (Optional, String)  
  Interface to bind vpn for route-based vpn.  
- **copy_outer_dscp** (Optional, Boolean)  
  Enable copying outer IP header DSCP and ECN to inner IP header.
- **df_bit** (Optional, String)  
  Specifies how to handle the Don't Fragment bit.  
  Need to be `clear`, `copy` or `set`.
- **establish_tunnels** (Optional, String)  
  When the VPN comes up.  
  Need to be `immediately` or `on-traffic`.
- **ike** (Optional, Block)  
  Declare IKE-keyed configuration.
  - **gateway** (Required, String)  
    The name of security IKE gateway (phase-1).
  - **policy** (Required, String)  
    The name of IPSec policy.
  - **identity_local** (Optional, String)  
    IPSec proxy-id local parameter.
  - **identity_remote** (Optional, String)  
    IPSec proxy-id remote parameter.
  - **identity_service** (Optional, String)  
    IPSec proxy-id service parameter.
- **manual** (Optional, Block)  
  Define a manual security association.
  - **external_interface** (Required, String)  
    External interface for the security association.
  - **protocol** (Required, String)  
    Define an IPSec protocol for the security association.  
    Need to be `ah` or `esp`.
  - **spi** (Required, Number)  
    Define security parameter index (256..16639).
  - **authentication_algorithm** (Optional, String)  
    Define authentication algorithm.
  - **authentication_key_hexa** (Optional, String, Sensitive)  
    Define an authentication key with format as hexadecimal.
  - **authentication_key_hexa** (Optional, String, Sensitive)  
    Define an authentication key with format as text.
  - **encryption_algorithm** (Optional, String)  
    Define encryption algorithm.
  - **encryption_key_hexa** (Optional, String, Sensitive)  
    Define an encryption key with format as hexadecimal.
  - **encryption_key_text** (Optional, String, Sensitive)  
    Define an encryption key with format as text.
  - **gateway** (Optional, String)  
    Define the IPSec peer.
- **multi_sa_forwarding_class** (Optional, Set of String)  
  Negotiate multiple SAs with forwarding-classes.
- **traffic_selector** (Optional, Block List)  
  For each name of traffic-selector to declare.
  - **name** (Required, String)  
    Name of traffic-selector.
  - **local_ip** (Required, String)  
    CIDR for IP addresses of local traffic-selector.
  - **remote_ip** (Required, String)  
    CIDR for IP addresses of remote traffic-selector.
- **udp_encapsulate** (Optional, Block)  
  UDP encapsulation of IPsec data traffic.
  - **dest_port** (Optional, Number)  
    UDP destination port (1025..65536).
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

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security IPSec vpn can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_ipsec_vpn.demo_vpn first-vpn
```
