---
layout: "junos"
page_title: "Junos: junos_security_policy_tunnel_pair_policy"
sidebar_current: "docs-junos-resource-security-policy-tunnel-pair-policy"
description: |-
  Add a tunnel pair policy options in each policy (when Junos device supports it)
---

# junos_security_policy

Provides a tunnel pair policy resource options in each policy.

## Example Usage

```hcl
# Add a tunnel pair policy
resource "junos_security_policy_tunnel_pair_policy" "DemoPair" {
  zone_a        = "trust"
  zone_b        = "untrust"
  policy_a_to_b = "trust_to_untrust"
  policy_b_to_a = "untrust_to_trust"
}
```

## Argument Reference

The following arguments are supported:

* `zone_a` - (Required)(`String`) The name of first zone.
* `zone_b` - (Required)(`String`) The name of second zone.
* `policy_a_to_b` - (Required)(`String`) The name of policy when from zone zone_a to zone zone_b.
* `policy_b_to_a` - (Required)(`String`) The name of policy when from zone zone_b to zone zone_a.

All arguments forces new resource

## Import

Junos security policy can be imported using an id made up of `<zone_a>_-_<policy_a_to_b>_-_<zone_b>_-_<policy_b_to_a>`, e.g.

```
$ terraform import junos_security_policy_tunnel_pair_policy.DemoDemoPair trust_-_trust_to_untrust_-_untrust_-_untrust_to_trust
```
