---
page_title: "Junos: junos_policyoptions_policy_statement"
---

# junos_policyoptions_policy_statement

Get configuration from a policy-options policy-statement.

## Example Usage

```hcl
# Read a policy-options policy-statement configuration
data "junos_policyoptions_policy_statement" "demo_policy" {
  name = "DemoPolicy"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  Name to identify the policy.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **dynamic_db** (Boolean)  
  Object may exist in dynamic database.
- **from** (Block)  
  Conditions to match the source of a route.  
  See [below for nested schema](#from-attributes).
- **to** (Block)  
  Conditions to match the destination of a route.  
  See [below for nested schema](#to-attributes).
- **then** (Block)  
  Actions to take if 'from' and 'to' conditions match.  
  See [below for nested schema](#then-attributes).
- **term** (Block List)  
  For each policy term.
  - **name** (String)  
    Name of term.
  - **from** (Block)  
    Conditions to match the source of a route.  
    See [below for nested schema](#from-attributes).
  - **to** (Block)  
    Conditions to match the destination of a route.  
    See [below for nested schema](#to-attributes).
  - **then** (Block)  
    Actions to take if 'from' and 'to' conditions match.  
    See [below for nested schema](#then-attributes).

---

### from attributes

- **aggregate_contributor** (Boolean)  
  Match more specifics of an aggregate.
- **bgp_as_path** (Set of String)  
  Name of AS path regular expression.
- **bgp_as_path_calc_length** (Block Set)  
  Number of BGP ASes excluding confederations.
  - **count** (Number)  
    Number of ASes (0..1024).
  - **match** (String)  
    Type of match: equal values, higher or equal values, lower or equal values.
- **bgp_as_path_group** (Set of String)  
  Name of AS path group.
- **bgp_as_path_unique_count** (Block Set)  
  Number of unique BGP ASes excluding confederations.
  - **count** (Number)  
    Number of ASes (0..1024).
  - **match** (String)  
    Type of match: equal values, higher or equal values, lower or equal values.
- **bgp_community** (Set of String)  
  BGP community.
- **bgp_community_count** (Block Set)  
  Number of BGP communities.
  - **count** (Number)  
    Number of communities (0..1024).
  - **match** (String)  
    Type of match: equal values, higher or equal values, lower or equal values.
- **bgp_origin** (String)  
  BGP origin attribute.
- **bgp_srte_discriminator** (Number)  
  Srte discriminator.
- **color** (Number)  
  Color (preference) value.
- **evpn_esi** (Set of String)  
  ESI in EVPN Route.
- **evpn_mac_route** (String)  
  EVPN Mac Route type.
- **evpn_tag** (Set of Number)  
  Tag in EVPN Route (0..4294967295).
- **family** (String)  
  Family.
- **local_preference** (Number)  
  Local preference associated with a route.
- **interface** (Set of String)  
  Interface name or address.
- **metric** (Number)  
  Metric value.
- **neighbor** (Set of String)  
  Neighboring router.
- **next_hop** (Set of String)  
  Next-hop router.
- **next_hop_type_merged** (Boolean)  
  Merged next hop.
- **next_hop_weight** (Block Set)  
  Weight of the gateway.
  - **match** (String)  
    Type of match for weight.
  - **weight** (Number)  
    Weight of the gateway (1..65535).
- **ospf_area** (String)  
  OSPF area identifier.
- **policy** (List of String)  
  Name of policy to evaluate.
- **preference** (Number)  
  Preference value.
- **prefix_list** (Set of String)  
  Prefix-lists of routes to match.
- **protocol** (Set of String)  
  Protocol from which route was learned.
- **route_filter** (Block List)  
  Routes to match.
  - **route** (String)  
    IP address.
  - **option** (String)  
    Mask option.
  - **option_value** (String)  
    For options that need an argument.
- **route_type** (String)  
  Route type.
- **routing_instance** (String)  
  Routing protocol instance.
- **srte_color** (Number)  
  Srte color.
- **state** (String)  
  Route state.
- **tunnel_type** (Set of String)  
  Tunnel type.
- **validation_database** (String)  
  Name to identify a validation-state.

---

### to attributes

- **bgp_as_path** (Set of String)  
  Name of AS path regular expression.
- **bgp_as_path_group** (Set of String)  
  Name of AS path group.
- **bgp_community** (Set of String)  
  BGP community.
- **bgp_origin** (String)  
  BGP origin attribute.
- **family** (String)  
  Family.
- **local_preference** (Number)  
  Local preference associated with a route.
- **interface** (Set of String)  
  Interface name or address.
- **metric** (Number)  
  Metric value.
- **neighbor** (Set of String)  
  Neighboring router.
- **next_hop** (Set of String)  
  Next-hop router.
- **ospf_area** (String)  
  OSPF area identifier.
- **policy** (List of String)  
  Name of policy to evaluate.
- **preference** (Number)  
  Preference value.
- **protocol** (Set of String)  
  Protocol from which route was learned.
- **routing_instance** (String)  
  Routing protocol instance.

---

### then attributes

- **action** (String)  
  Action `accept` or `reject`.
- **as_path_expand** (String)  
  Prepend AS numbers prior to adding local-as.
- **as_path_prepend** (String)  
  Prepend AS numbers to an AS path.
- **community** (Block List)  
  For each community action.
  - **action** (String)  
    Action on BGP community.
  - **value** (String)  
    Name to identify a BGP community.
- **default_action** (String)  
  Set default policy action.
- **load_balance** (String)  
  Type of load balancing in forwarding table.
- **local_preference** (Block)  
  Declare local-preference action.
  - **action** (String)  
    Action on local-preference.
  - **value** (Number)  
    Value for action (local-preference, constant).
- **metric** (Block)  
  Declare metric action.
  - **action** (String)  
    Action on metric.
  - **value** (Number)  
    Value for action (metric, constant).
- **next** (String)  
  Skip to next `policy` or `term`.
- **next_hop** (String)  
  Set the address of the next-hop router.
- **origin** (String)  
  BGP path origin.
- **preference** (Block)  
  Declare preference action.
  - **action** (String)  
    Action on preference.
  - **value** (Number)  
    Value for action (preference, constant).
