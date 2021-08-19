---
layout: "junos"
page_title: "Junos: junos_security_ike_gateway"
sidebar_current: "docs-junos-resource-security-ike-gateway"
description: |-
  Create a security ike gateway (when Junos device supports it)
---

# junos_security_ike_gateway

Provides a security ike gateway resource.

## Example Usage

```hcl
# Add a ike gateway
resource junos_security_ike_gateway "demo_vpn_p1" {
  name               = "first-vpn"
  address            = ["192.0.2.1"]
  policy             = "ike-policy"
  external_interface = "ge-0/0/0.0"
}
```

## Argument Reference

-> **Note:** One of `address` or `dynamic_remote` arguments is required.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of ike phase1.
- **external_interface** (Required, String)  
  Interface for ike negotiations.
- **policy** (Required, String)  
  Ike policy.
- **address** (Optional, List of String)  
  List of Peer IP.  
  Need to set one of `address` or `dynamic_remote`.
- **dynamic_remote** (Optional, Block)  
  Declare site to site peer with dynamic IP address.  
  See [below for nested schema](#dynamic_remote-arguments).  
  Need to set one of `address` or `dynamic_remote`.
- **aaa** (Optional, Block)  
  Declare `aaa` configuration.
  - **access_profile** (Optional, String)  
    Access profile that contains authentication information.  
    Conflict with `aaa.client_*`.
  - **client_password** (Optional, String, Sensitive)  
    AAA client password with 1 to 128 characters.  
    Conflict with `aaa.access_profile`.  
  - **client_username** (Optional, String)  
    AAA client username with 1 to 128 characters.  
    Conflict with `aaa.access_profile`.
- **dead_peer_detection** (Optional, Block)  
  Declare RFC-3706 DPD configuration.  
  See [below for nested schema](#dead_peer_detection-arguments).
- **general_ike_id** (Optional, Boolean)  
  Accept peer IKE-ID in general.
- **local_address** (Optional, String)  
  Local IP for ike negotiations.
- **local_identity** (Optional, Block)  
  Declare local IKE identity configuration.
  - **type** (Required, String)  
    Type of IKE identity.
  - **value** (Optional, String)  
    Value for IKE identity.
- **no_nat_traversal** (Optional, Boolean)  
  Disable IPSec NAT traversal.
- **remote_identity** (Optional, Block)  
  Declare remote IKE identity configuration.
  - **type** (Required, String)  
    Type of IKE identity.
  - **value** (Optional, String)  
    Value for IKE identity.
- **version** (Optional, String)  
  Negotiate using either IKE v1 or IKE v2 protocol.  
  Need to be `v1-only` or `v2-only`.

---

### dynamic_remote arguments

-> **Note:** You can only choose one argument between `distinguished_name`, `hostname`, `inet`,
`inet6` and `user_at_hostname`.

- **connections_limit** (Optional, Number)  
  Maximum number of users connected to gateway.
- **distinguished_name** (Optional, Block)  
  Declare distinguished-name configuration.
  - **container** (Optional, String)  
    Container string for a distinguished name.
  - **wildcard** (Optional, String)  
    Wildcard string for a distinguished name.
- **hostname** (Optional, String)  
  Use a fully-qualified domain name.
- **ike_user_type** (Optional, String)  
  Type of the IKE ID.  
  Need to be `shared-ike-id` or `group-ike-id`.
- **inet** (Optional, String)  
  Use an IPV4 address to identify the dynamic peer.
- **inet6** (Optional, String)  
  Use an IPV6 address to identify the dynamic peer.
- **reject_duplicate_connection** (Optional, Boolean)  
  Reject new connection from duplicate IKE-id.
- **user_at_hostname** (Optional, String)  
  Use an e-mail address.

---

### dead_peer_detection arguments

- **interval** (Optional, Number)  
  The interval at which to send DPD.
- **send_mode** (Optional, String)  
  Specify how probes are sent.  
  Need to be `always-send`, `optimized` or `probe-idle-tunnel`.  
  - always-send -> Send probes periodically regardless of incoming and outgoing data traffic.  
  - optimized -> Send probes only when there is outgoing and no incoming data traffic - RFC3706.
  - probe_idle_tunnel -> Send probes same as in optimized mode and also when there is no outgoing
  & incoming data traffic.
- **threshold** (Optional, Number)  
  Maximum number of DPD retransmissions.

## Import

Junos security ike gateway can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_ike_gateway.demo_vpn_p1 first-vpn
```
