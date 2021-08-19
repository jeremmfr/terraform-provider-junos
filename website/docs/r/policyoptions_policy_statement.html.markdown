---
layout: "junos"
page_title: "Junos: junos_policyoptions_policy_statement"
sidebar_current: "docs-junos-resource-policyoptions-policy-statement"
description: |-
  Create a routing policy
---

# junos_policyoptions_policy_statement

Provides a routing policy resource.

## Example Usage

```hcl
# Add a policy
resource junos_policyoptions_policy_statement "demo_policy" {
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
  The name of routing policy.
- **add_it_to_forwarding_table_export** (Optional, Boolean)  
  Add this policy in `routing-options forwarding-table export` list.
- **from** (Optional, Block)  
  Declare from filter.  
  See [below for nested schema](#from-arguments).
- **to** (Optional, Block)  
  Declare to filter.  
  See [below for nested schema](#to-arguments).
- **then** (Optional, Block)  
  Declare then actions.  
  See [below for nested schema](#then-arguments).
- **term** (Optional, Block List)  
  For each name of term.
  - **name** (Required, String)  
    Name of policy
  - **from** (Optional, Block)  
    Declare from filter.  
    See [below for nested schema](#from-arguments).
  - **to** (Optional, Block)  
    Declare to filter.  
    See [below for nested schema](#to-arguments).
  - **then** (Optional, Block)  
    Declare then actions.  
    See [below for nested schema](#then-arguments).

---

### from arguments

- **aggregate_contributor** (Optional, Boolean)  
  Match more specifics of an aggregate.
- **bgp_as_path** (Optional, List of String)  
  Name of AS path regular expression.  
  See resource `junos_policyoptions_as_path`.
- **bgp_as_path_group** (Optional, List of String)  
  Name of AS path group.  
  See resource `junos_policyoptions_as_path_group`.
- **bgp_community** (Optional, List of String)  
  BGP community.  
  See resource `junos_policyoptions_community`.
- **bgp_origin** (Optional, String)  
  BGP origin attribute.  
  Need to be `egp`, `igp` or `incomplete`.
- **family** (Optional, String)  
  IP family.
- **local_preference** (Optional, Number)  
  Local preference associated with a route.
- **routing_instance** (Optional, String)  
  Routing protocol instance.
- **interface** (Optional, List of String)  
  List of interface name
- **metric** (Optional, Number)  
  Metric value
- **neighbor** (Optional, List of String)  
  Neighboring router
- **next_hop** (Optional, List of String)  
  Next-hop router
- **ospf_area** (Optional, String)  
  OSPF area identifier
- **policy** (Optional, List of String)  
  Name of policy to evaluate
- **preference** (Optional, Number)  
  Preference value
- **prefix_list** (Optional, List of String)  
  List of prefix-lists of routes to match.  
  See resource `junos_policyoptions_prefix_list`.
- **protocol** (Optional, List of String)  
  Protocol from which route was learned
- **route_filter** (Optional, Block List)  
  For each filter to declare.
  - **route** (Required, String)  
    IP address
  - **option** (Required, String)  
    Mask option.  
    Need to be `address-mask`, `exact`, `longer`, `orlonger`, `prefix-length-range`, `through` or `upto`.
  - **option_value** (Optional, String)  
    For options that need an argument

---

### to arguments

- **bgp_as_path** (Optional, List of String)  
  Name of AS path regular expression.  
  See resource `junos_policyoptions_as_path`.
- **bgp_as_path_group** (Optional, List of String)  
  Name of AS path group.  
  See resource `junos_policyoptions_as_path_group`.
- **bgp_community** (Optional, List of String)  
  BGP community.  
  See resource `junos_policyoptions_community`.
- **bgp_origin** (Optional, String)  
  BGP origin attribute.  
  Need to be `egp`, `igp` or `incomplete`.
- **family** (Optional, String)  
  IP family.
- **local_preference** (Optional, Number)  
  Local preference associated with a route.
- **routing_instance** (Optional, String)  
  Routing protocol instance.
- **interface** (Optional, List of String)  
  List of interface name
- **metric** (Optional, Number)  
  Metric value
- **neighbor** (Optional, List of String)  
  Neighboring router
- **next_hop** (Optional, List of String)  
  Next-hop router
- **ospf_area** (Optional, String)  
  OSPF area identifier
- **policy** (Optional, List of String)  
  Name of policy to evaluate
- **preference** (Optional, Number)  
  Preference value
- **protocol** (Optional, List of String)  
  Protocol from which route was learned

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
    Value for action
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
    Value for action
- **next** (Optional, String)  
  Skip to next `policy` or `term`.
- **next_hop** (Optional, String)  
  Set the address of the next-hop router
- **metric** (Optional, Block)  
  Declare metric action.
  - **action** (Required, String)  
    Action on metric.  
    Need to be `add`, `subtract` or `none`.
  - **value** (Required, String)  
    Value for action
- **origin** (Optional, String)  
  BGP path origin
- **preference** (Optional, Block)  
  Declare preference action.
  - **action** (Required, String)  
    Action on preference.  
    Need to be `add`, `subtract` or `none`.
  - **value** (Required, String)  
    Value for action

## Import

Junos policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_policy_statement.demo_policy DemoPolicy
```
