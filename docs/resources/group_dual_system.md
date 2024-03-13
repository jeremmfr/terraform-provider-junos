---
page_title: "Junos: junos_group_dual_system"
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

- **name** (Required, String, Forces new resource)  
  Name of group.
- **apply_groups** (Optional, Boolean)  
  Apply the group.  
  Defaults to `true`.
- **interface_fxp0** (Optional, Block)  
  Configure `fxp0` interface.  
  See [below for nested schema](#interface_fxp0-arguments).
- **routing_options** (Optional, Block)  
  Configure `routing-options` block.  
  See [below for nested schema](#routing_options-arguments).
- **security** (Optional, Block)  
  Configure `security` block.
  - **log_source_address** (Required, String)  
    Source ip address used when exporting security logs.
- **system** (Optional, Block)  
  Configure `system` block.
  - **host_name** (Optional, String)  
    Hostname.
  - **backup_router_address** (Optional, String)  
    IPv4 address backup router.
  - **backup_router_destination** (Optional, Set of String)  
    Destinations network reachable through the IPv4 router.
  - **inet6_backup_router_address** (Optional, String)  
    IPv6 address backup router.
  - **inet6_backup_router_destination** (Optional, Set of String)  
    Destinations network reachable through the IPv6 router.

---

### interface_fxp0 arguments

- **description** (Optional, String)  
  Description for interface.
- **family_inet_address** (Optional, Block List)  
  For each ip address to declare.
  - **cidr_ip** (Required, String)  
    Address IP/Mask v4.
  - **master_only** (Optional, Boolean)  
    Master management IP address.
  - **preferred** (Optional, Boolean)  
    Preferred address on interface.
  - **primary** (Optional, Boolean)  
    Candidate for primary address in system.
- **family_inet6_address** (Optional, Block List)  
  For each ip v6 address to declare.
  - **cidr_ip** (Required, String)  
    Address IP/Mask v6.
  - **master_only** (Optional, Boolean)  
    Master management IP address.
  - **preferred** (Optional, Boolean)  
    Preferred address on interface.
  - **primary** (Optional, Boolean)  
    Candidate for primary address in system.

---

### routing_options arguments

- **static_route** (Required, Block List)  
  For each destination to declare.
  - **destination** (Required, String)  
    The destination for static route.
  - **next_hop** (Required, List of String)  
    List of next-hop.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos group can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_group_dual_system.node0 node0
```
