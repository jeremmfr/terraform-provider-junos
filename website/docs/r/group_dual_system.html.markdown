---
layout: "junos"
page_title: "Junos: junos_group_dual_system"
sidebar_current: "docs-junos-group-dual-system"
description: |-
  Create a group for member of dual system (node/re) 
---

# junos_group_dual_system

Provides a group resource for member of dual system (node/re).

## Example Usage

```hcl
# Configure a group for member of dual system.
resource "junos_group_dual_system" "node0" {
  name = "node0"
  interface_fxp0 {
    description = "demo"
    family_inet_address {
      cidr_ip = "192.0.2.193/26"
    }
    family_inet_address {
      cidr_ip     = "192.0.2.194/26"
      master_only = true
    }
  }
  routing_options {
    static_route {
      destination = "192.0.2.0/26"
      next_hop    = ["192.0.2.254"]
    }
  }
  system {
    host_name                 = "test_node"
    backup_router_address     = "192.0.2.254"
    backup_router_destination = ["192.0.2.0/26"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Name of group.
* `apply_groups` - (Optional)(`Bool`) Apply the group. Defaults to `true`.
* `interface_fxp0` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for configure `fxp0` interface. See the [`interface_fxp0` arguments] (#interface_fxp0-arguments) block.
* `routing_options` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for configure `routing-options` block. See the [`routing_options` arguments] (#routing_options-arguments) block.
* `security` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for configure `security` block. 
  * `log_source_address` - (Required)(`String`) Source ip address used when exporting security logs. 
* `system` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) 
Can be specified only once for configure `system` block.
  * `host_name` - (Optional)(`String`) Hostname.
  * `backup_router_address` - (Optional)(`String`) IPv4 address backup router.
  * `backup_router_destination` - (Optional)(`ListOfString`) Destinations network reachable through the router.

---
#### interface_fxp0 arguments
* `description` - (Optional)(`String`) Description for interface.
* `family_inet_address` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each ip address to declare.
  * `cidr_ip` - (Required)(`String`) Address IP/Mask v4.
  * `master_only` - (Optional)(`Bool`) Master management IP address.
  * `preferred` - (Optional)(`Bool`) Preferred address on interface.
  * `primary` - (Optional)(`Bool`) Candidate for primary address in system.
* `family_inet6_address` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each ip v6 address to declare.
  * `cidr_ip` - (Required)(`String`) Address IP/Mask v6.
  * `master_only` - (Optional)(`Bool`) Master management IP address.
  * `preferred` - (Optional)(`Bool`) Preferred address on interface.
  * `primary` - (Optional)(`Bool`) Candidate for primary address in system.

---
#### routing_options arguments
* `static_route` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each destination to declare.
  * `destination` - (Required)(`String`) The destination for static route.
  * `next_hop` - (Required)(`ListOfString`) List of next-hop.

## Import

Junos group can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_group_dual_system.node0 node0
```
