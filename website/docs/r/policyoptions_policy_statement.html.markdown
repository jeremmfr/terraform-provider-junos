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

* `name` - (Required, Forces new resource)(`String`) The name of routing policy.
* `from` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Declare from filter. See the [`from` arguments](#from-arguments) block. Max of 1.
* `to` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Declare to filter. See the [`to` arguments](#to-arguments) block. Max of 1.
* `then` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Declare then actions. See the [`then` arguments](#then-arguments) block. Max of 1.
* `term` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each term declaration.
  * `name`(Required)(`String`) Name of policy
  * `from` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Declare from filter. See the [`from` arguments](#from-arguments) block. Max of 1.
  * `to` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Declare to filter. See the [`to` arguments](#to-arguments) block. Max of 1.
  * `then` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Declare then actions. See the [`then` arguments](#then-arguments) block. Max of 1.

#### from arguments
* `aggregate_contributor` - (Optional)(`Bool`) Match more specifics of an aggregate.
* `bgp_as_path` - (Optional)(`ListOfString`) Name of AS path regular expression. See resource `junos_policyoptions_as_path`.
* `bgp_as_path_group` - (Optional)(`ListOfString`) Name of AS path group. See resource `junos_policyoptions_as_path_group`.
* `bgp_community` - (Optional)(`ListOfString`) BGP community. See resource `junos_policyoptions_community`.
* `bgp_origin` - (Optional)(`String`) BGP origin attribute. Need to be 'egp', 'igp' or 'incomplete'.
* `family` - (Optional)(`String`) IP family.
* `local_preference` - (Optional)(`Int`) Local preference associated with a route.
* `routing_instance` - (Optional)(`String`) Routing protocol instance.
* `interface` - (Optional)(`ListOfString`) List of interface name
* `metric` - (Optional)(`Int`) Metric value
* `neighbor` - (Optional)(`ListOfString`) Neighboring router
* `next_hop` - (Optional)(`ListOfString`) Next-hop router
* `ospf_area` - (Optional)(`String`) OSPF area identifier
* `policy` - (Optional)(`ListOfString`) Name of policy to evaluate
* `preference` - (Optional)(`Int`) Preference value
* `prefix_list` - (Optional)(`ListOfString`) List of prefix-lists of routes to match. See resource `junos_policyoptions_prefix_list`.
* `protocol` - (Optional)(`ListOfString`) Protocol from which route was learned
* `route_filter` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each filter to declare.
  * `route` - (Required)(`String`) IP address
  * `option` - (Required)(`String`): Mask option. Need to be 'address-mask', 'exact', 'longer', 'orlonger', 'prefix-length-range', 'through' or 'upto'.
  * `option_value` - (Optional)(`String`) For options that need an argument

#### to arguments
* `bgp_as_path` - (Optional)(`ListOfString`) Name of AS path regular expression. See resource `junos_policyoptions_as_path`.
* `bgp_as_path_group` - (Optional)(`ListOfString`) Name of AS path group. See resource `junos_policyoptions_as_path_group`.
* `bgp_community` - (Optional)(`ListOfString`) BGP community. See resource `junos_policyoptions_community`.
* `bgp_origin` - (Optional)(`String`) BGP origin attribute. Need to be egp, igp or incomplete.
* `family` - (Optional)(`String`) IP family.
* `local_preference` - (Optional)(`Int`) Local preference associated with a route.
* `routing_instance` - (Optional)(`String`) Routing protocol instance.
* `interface` - (Optional)(`ListOfString`) List of interface name
* `metric` - (Optional)(`Int`) Metric value
* `neighbor` - (Optional)(`ListOfString`) Neighboring router
* `next_hop` - (Optional)(`ListOfString`) Next-hop router
* `ospf_area` - (Optional)(`String`) OSPF area identifier
* `policy` - (Optional)(`ListOfString`) Name of policy to evaluate
* `preference` - (Optional)(`Int`) Preference value
* `protocol` - (Optional)(`ListOfString`) Protocol from which route was learned

#### then arguments
* `action` - (Optional)(`String`) Action 'accept' or 'reject'.
* `as_path_expand` - (Optional)(`String`) Prepend AS numbers prior to adding local-as.
* `as_path_prepend` - (Optional)(`String`) Prepend AS numbers to an AS path.
* `community` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each community action.
  * `action` - (Required)(`String`) Action on BGP community. Need to be 'add', 'delete' or 'set'.
  * `value` - (Required)(`String`) Value for action
* `default_action` - (Optional)(`String`) Set default policy action. Can be 'accept' or 'reject'.
* `load_balance` - (Optional)(`String`) Type of load balancing in forwarding table. Can be 'per-packet' or 'consistent-hash'.
* `local_preference` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare local-preference action.
  * `action` - (Required)(`String`) Action on local-preference. Need to be 'add', 'subtract' or 'none'.
  * `value` - (Required)(`String`) Value for action
* `next` - (Optional)(`String`) Skip to next 'policy' or 'term'.
* `next_hop` - (Optional)(`String`) Set the address of the next-hop router
* `metric` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare metric action.
  * `action` - (Required)(`String`) Action on metric. Need to be 'add', 'subtract' or 'none'.
  * `value` - (Required)(`String`) Value for action
* `origin` - (Optional)(`String`) BGP path origin
* `preference` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare preference action.
  * `action` - (Required)(`String`) Action on preference. Need to be 'add', 'subtract' or 'none'.
  * `value` - (Required)(`String`) Value for action

## Import

Junos policy can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_policyoptions_policy_statement.demo_policy DemoPolicy
```
