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

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of ike phase1.
* `address` - (Required)(`ListOfString`) List of Peer IP
* `local_address` - (Optional)(`String`) Local IP for ike negotiations.
* `policy` - (Required)(`String`) Ike policy.
* `external_interface` - (Required)(`String`) Interface for ike negotiations.
* `general_ike_id` - (Optional)(`Bool`) Accept peer IKE-ID in general.
* `no_nat_traversal` - (Optional)(`Bool`) Disable IPSec NAT traversal.
* `dead_peer_detection` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare RFC-3706 DPD configuration.
  * `interval` - (Optional)(`Int`) The interval at which to send DPD
  * `threshold` - (Optional)(`Int`) Maximum number of DPD retransmissions
* `local_identity` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare local IKE identity configuration.
  * `type` - (Required)(`String`) Type of IKE identity.
  * `value` - (Optional)(`String`) Value for IKE identity
* `remote_identity` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare remote IKE identity configuration.
  * `type` - (Required)(`String`) Type of IKE identity.
  * `value` - (Optional)(`String`) Value for IKE identity
* `version` - (Optional)(`String`) Negotiate using either IKE v1 or IKE v2 protocol. Need to be 'v1-only' or 'v2-only'.

## Import

Junos security ike gateway can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_ike_gateway.demo_vpn_p1 first-vpn
```
