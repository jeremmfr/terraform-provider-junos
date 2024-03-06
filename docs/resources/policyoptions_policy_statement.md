---
page_title: "Junos: junos_policyoptions_policy_statement"
---

# junos_policyoptions_policy_statement

Provides a routing policy resource.

## Example Usage

```hcl
# Add a policy
resource "junos_policyoptions_policy_statement" "demo_policy" {
  name = "DemoPolicy"
  from {
    protocol = ["bgp"]
  }
  term {
    name = "term_1"
    from {
      route_filter {
        route  = "192.0.2.0/25"
        option = "orlonger"
      }
    }
    then {
      action = "accept"
    }
  }
  term {
    name = "term_2"
    then {
      action = "reject"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name to identify the policy.
- **add_it_to_forwarding_table_export** (Optional, Boolean)  
  Add this policy in `routing-options forwarding-table export` list.
- **dynamic_db** (Optional, Boolean)  
  Object may exist in dynamic database.
- **from** (Optional, Block)  
  Conditions to match the source of a route.  
  See [below for nested schema](#from-arguments).
- **to** (Optional, Block)  
  Conditions to match the destination of a route.  
  See [below for nested schema](#to-arguments).
- **then** (Optional, Block)  
  Actions to take if 'from' and 'to' conditions match.  
  See [below for nested schema](#then-arguments).
- **term** (Optional, Block List)  
  For each policy term.
  - **name** (Required, String)  
    Name of term.
  - **from** (Optional, Block)  
    Conditions to match the source of a route.  
    See [below for nested schema](#from-arguments).
  - **to** (Optional, Block)  
    Conditions to match the destination of a route.  
    See [below for nested schema](#to-arguments).
  - **then** (Optional, Block)  
    Actions to take if 'from' and 'to' conditions match.  
    See [below for nested schema](#then-arguments).

---

### from arguments

- **aggregate_contributor** (Optional, Boolean)  
  Match more specifics of an aggregate.
- **bgp_as_path** (Optional, Set of String)  
  Name of AS path regular expression.  
  See resource `junos_policyoptions_as_path`.
- **bgp_as_path_calc_length** Optional, Block Set)  
  For each count, number of BGP ASes excluding confederations.
  - **count** (Required, Number)  
    Number of ASes (0..1024).
  - **match** (Required, String)  
    Type of match: equal values, higher or equal values, lower or equal values.  
    Need to `equal`, `orhigher` or `orlower`.
- **bgp_as_path_group** (Optional, Set of String)  
  Name of AS path group.  
  See resource `junos_policyoptions_as_path_group`.
- **bgp_as_path_unique_count** (Optional, Block Set)  
  For each count, number of unique BGP ASes excluding confederations.
  - **count** (Required, Number)  
    Number of ASes (0..1024).
  - **match** (Required, String)  
    Type of match: equal values, higher or equal values, lower or equal values.  
    Need to `equal`, `orhigher` or `orlower`.
- **bgp_community** (Optional, Set of String)  
  BGP community.  
  See resource `junos_policyoptions_community`.
- **bgp_community_count** (Optional, Block Set)  
  For each count, number of BGP communities.
  - **count** (Required, Number)  
    Number of communities (0..1024).
  - **match** (Required, String)  
    Type of match: equal values, higher or equal values, lower or equal values.  
    Need to `equal`, `orhigher` or `orlower`.
- **bgp_origin** (Optional, String)  
  BGP origin attribute.  
  Need to be `egp`, `igp` or `incomplete`.
- **bgp_srte_discriminator** (Optional, Number)  
  Srte discriminator.
- **color** (Optional, Number)  
  Color (preference) value.
- **evpn_esi** (Optional, Set of String)  
  ESI in EVPN Route.
- **evpn_mac_route** (Optional, String)  
  EVPN Mac Route type.  
  Need to be `mac-ipv4`, `mac-ipv6` or `mac-only`.
- **evpn_tag** (Optional, Set of Number)  
  Tag in EVPN Route (0..4294967295).
- **family** (Optional, String)  
  Family.
- **local_preference** (Optional, Number)  
  Local preference associated with a route.
- **interface** (Optional, Set of String)  
  Interface name or address.
- **metric** (Optional, Number)  
  Metric value.
- **neighbor** (Optional, Set of String)  
  Neighboring router.
- **next_hop** (Optional, Set of String)  
  Next-hop router.
- **next_hop_type_merged** (Optional, Boolean)  
  Merged next hop.
- **next_hop_weight** (Optional, Block Set)  
  For each combination of block arguments, weight of the gateway.
  - **match** (Required, String)  
    Type of match for weight.  
    Need to be `equal`, `greater-than`, `greater-than-equal`, `less-than` or `less-than-equal`.
  - **weight** (Required, Weight)  
    Weight of the gateway (1..65535).
- **ospf_area** (Optional, String)  
  OSPF area identifier.
- **policy** (Optional, List of String)  
  Name of policy to evaluate.
- **preference** (Optional, Number)  
  Preference value.
- **prefix_list** (Optional, Set of String)  
  Prefix-lists of routes to match.  
  See resource `junos_policyoptions_prefix_list`.
- **protocol** (Optional, Set of String)  
  Protocol from which route was learned.
- **route_filter** (Optional, Block List)  
  For each routes to match.
  - **route** (Required, String)  
    IP address.
  - **option** (Required, String)  
    Mask option.  
    Need to be `address-mask`, `exact`, `longer`, `orlonger`, `prefix-length-range`, `through` or `upto`.
  - **option_value** (Optional, String)  
    For options that need an argument.
- **route_type** (Optional, String)  
  Route type.  
  Need to be `external` or `internal`.
- **routing_instance** (Optional, String)  
  Routing protocol instance.
- **srte_color** (Optional, Number)  
  Srte color.
- **state** (Optional, String)  
  Route state.  
  Need to be `active` or `inactive`.
- **tunnel_type** (Optional, Set of String)  
  Tunnel type.  
  Element need to be `gre`, `ipip` or `udp`.
- **validation_database** (Optional, String)  
  Name to identify a validation-state.  
  Need to be `invalid`, `unknown` or `valid`.

---

### to arguments

- **bgp_as_path** (Optional, Set of String)  
  Name of AS path regular expression.  
  See resource `junos_policyoptions_as_path`.
- **bgp_as_path_group** (Optional, Set of String)  
  Name of AS path group.  
  See resource `junos_policyoptions_as_path_group`.
- **bgp_community** (Optional, Set of String)  
  BGP community.  
  See resource `junos_policyoptions_community`.
- **bgp_origin** (Optional, String)  
  BGP origin attribute.  
  Need to be `egp`, `igp` or `incomplete`.
- **family** (Optional, String)  
  Family.
- **local_preference** (Optional, Number)  
  Local preference associated with a route.
- **interface** (Optional, Set of String)  
  Interface name or address.
- **metric** (Optional, Number)  
  Metric value.
- **neighbor** (Optional, Set of String)  
  Neighboring router.
- **next_hop** (Optional, Set of String)  
  Next-hop router.
- **ospf_area** (Optional, String)  
  OSPF area identifier.
- **policy** (Optional, List of String)  
  Name of policy to evaluate.
- **preference** (Optional, Number)  
  Preference value.
- **protocol** (Optional, Set of String)  
  Protocol from which route was learned.
- **routing_instance** (Optional, String)  
  Routing protocol instance.

---

### then arguments

- **action** (Optional, String)  
  Action `accept` or `reject`.
- **as_path_expand** (Optional, String)  
  Prepend AS numbers prior to adding local-as.
- **as_path_prepend** (Optional, String)  
  Prepend AS numbers to an AS path.
- **community** (Optional, Block List)  
  For each community action.
  - **action** (Required, String)  
    Action on BGP community.  
    Need to be `add`, `delete` or `set`.
  - **value** (Required, String)  
    Name to identify a BGP community.
- **default_action** (Optional, String)  
  Set default policy action.  
  Need to be `accept` or `reject`.
- **load_balance** (Optional, String)  
  Type of load balancing in forwarding table.  
  Need to be `per-packet` or `consistent-hash`.
- **local_preference** (Optional, Block)  
  Declare local-preference action.
  - **action** (Required, String)  
    Action on local-preference.  
    Need to be `add`, `subtract` or `none`.
  - **value** (Required, String)  
    Value for action (local-preference, constant).
- **metric** (Optional, Block)  
  Declare metric action.
  - **action** (Required, String)  
    Action on metric.  
    Need to be `add`, `subtract` or `none`.
  - **value** (Required, String)  
    Value for action (metric, constant).
- **next** (Optional, String)  
  Skip to next `policy` or `term`.
- **next_hop** (Optional, String)  
  Set the address of the next-hop router.  
  Need to be a valid IP or one of `discard`, `next-table`, `peer-address`, `reject`, `self`.
- **origin** (Optional, String)  
  BGP path origin.
- **preference** (Optional, Block)  
  Declare preference action.
  - **action** (Required, String)  
    Action on preference.  
    Need to be `add`, `subtract` or `none`.
  - **value** (Required, String)  
    Value for action (preference, constant).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_policy_statement.demo_policy DemoPolicy
```
