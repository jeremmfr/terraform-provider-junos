---
page_title: "Junos: junos_igmp_snooping_vlan"
---

# junos_igmp_snooping_vlan

Provides a IGMP snooping vlan resource.

## Example Usage

```hcl
resource "junos_igmp_snooping_vlan" "all" {
  name = "all"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  VLAN name or `all`.
- **routing_instance** (Optional, String, Forces new resource)  
  Configure IGMP snooping vlan in routing instance.  
  Need to be `default` (for root level) or name of routing instance.  
  Defaults to `default`.
- **immediate_leave** (Optional, Boolean)  
  Enable immediate group leave on interfaces.
- **interface** (Optional, Block List)  
  For each interface name, configure interface options for IGMP.  
  See [below for nested schema](#interface-arguments).
- **l2_querier_source_address** (Optional, String)  
  Enable L2 querier mode with source address.
- **proxy** (Optional, Boolean)  
  Enable proxy mode.
- **proxy_source_address** (Optional, String)  
  Source IP address to use for proxy.  
  `proxy` need to be true.
- **query_interval** (Optional, Number)  
  When to send host query messages (1..1024 seconds).
- **query_last_member_interval** (Optional, String)  
  When to send group query messages (seconds).
- **query_response_interval** (Optional, String)  
  How long to wait for a host query response (seconds).
- **robust_count** (Optional, Number)  
  Expected packet loss on a subnet (2..10).

---

### interface arguments

- **name** (Required, String)  
  Interface name.  
  Need to be a logical interface (with dot).
- **group_limit** (Optional, Number)  
  Maximum number of groups an interface can join.
- **host_only_interface** (Optional, Boolean)  
  Enable interface to be treated as host-side interface.
- **immediate_leave** (Optional, Boolean)  
  Enable immediate group leave on interface.
- **multicast_router_interface** (Optional, Boolean)  
  Enabling multicast-router-interface on the interface.
- **static_group** (Optional, Block Set)  
  For each static group address.
  - **address** (Required, String)  
    IP multicast group address.
  - **source** (Optional, String)  
    IP multicast source address.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>`.

## Import

Junos igmp-snooping vlan can be imported using an id made up of `<name>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_igmp_snooping_vlan.all all_-_default
```
