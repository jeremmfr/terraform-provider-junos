resource "junos_routing_instance" "testacc_policyOptions" {
  name = "testacc_policyOptions"
}
resource "junos_policyoptions_as_path" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  path = "5|12|18"
}
resource "junos_policyoptions_as_path_group" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  as_path {
    name = "testacc policyOptions"
    path = "5|12|18"
  }
}
resource "junos_policyoptions_community" "testacc_policyOptions" {
  name    = "testacc_policyOptions"
  members = ["65000:100"]
}
resource "junos_policyoptions_prefix_list" "testacc_policyOptions" {
  name   = "testacc_policyOptions"
  prefix = ["192.0.2.0/25"]
}
resource "junos_policyoptions_prefix_list" "testacc_policyOptions2" {
  name   = "testacc policyOptions2"
  prefix = ["192.0.2.0/25", "fe80::/64"]
}
resource "junos_policyoptions_policy_statement" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  from {
    aggregate_contributor = true
    bgp_as_path           = [junos_policyoptions_as_path.testacc_policyOptions.name]
    bgp_as_path_calc_length {
      count = 4
      match = "orhigher"
    }
    bgp_community = [junos_policyoptions_community.testacc_policyOptions.name]
    bgp_community_count {
      count = 6
      match = "orhigher"
    }
    bgp_origin       = "igp"
    color            = 31
    family           = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface        = ["st0.0"]
    metric           = 5
    neighbor         = ["192.0.2.4"]
    next_hop         = ["192.0.2.4"]
    ospf_area        = "0.0.0.0"
    preference       = 100
    prefix_list      = [junos_policyoptions_prefix_list.testacc_policyOptions.name]
    protocol         = ["bgp"]
    route_filter {
      route  = "192.0.2.0/25"
      option = "exact"
    }
    route_filter {
      route        = "192.0.2.128/25"
      option       = "prefix-length-range"
      option_value = "/26-/27"
    }
  }
  to {
    bgp_as_path      = [junos_policyoptions_as_path.testacc_policyOptions.name]
    bgp_community    = [junos_policyoptions_community.testacc_policyOptions.name]
    bgp_origin       = "igp"
    family           = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface        = ["st0.0"]
    metric           = 5
    neighbor         = ["192.0.2.5"]
    next_hop         = ["192.0.2.5"]
    ospf_area        = "0.0.0.0"
    policy           = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
    preference       = 100
    protocol         = ["bgp"]
  }
  then {
    action          = "accept"
    as_path_expand  = "65000 65000"
    as_path_prepend = "65000 65000"
    community {
      action = "set"
      value  = junos_policyoptions_community.testacc_policyOptions.name
    }
    community {
      action = "delete"
      value  = junos_policyoptions_community.testacc_policyOptions.name
    }
    community {
      action = "add"
      value  = junos_policyoptions_community.testacc_policyOptions.name
    }
    default_action = "reject"
    load_balance   = "per-packet"
    local_preference {
      action = "add"
      value  = 10
    }
    next     = "policy"
    next_hop = "192.0.2.4"
    metric {
      action = "add"
      value  = 10
    }
    origin = "igp"
    preference {
      action = "add"
      value  = 10
    }
  }
  term {
    name = "term"
    from {
      aggregate_contributor = true
      bgp_as_path           = [junos_policyoptions_as_path.testacc_policyOptions.name]
      bgp_community         = [junos_policyoptions_community.testacc_policyOptions.name]
      bgp_origin            = "igp"
      family                = "inet"
      local_preference      = 100
      routing_instance      = junos_routing_instance.testacc_policyOptions.name
      interface             = ["st0.0"]
      metric                = 5
      neighbor              = ["192.0.2.4"]
      next_hop              = ["192.0.2.4"]
      ospf_area             = "0.0.0.0"
      policy                = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
      preference            = 100
      prefix_list           = [junos_policyoptions_prefix_list.testacc_policyOptions.name]
      protocol              = ["bgp"]
      route_filter {
        route  = "192.0.2.0/25"
        option = "exact"
      }
      route_filter {
        route        = "192.0.2.128/25"
        option       = "prefix-length-range"
        option_value = "/26-/27"
      }
    }
    to {
      bgp_as_path      = [junos_policyoptions_as_path.testacc_policyOptions.name]
      bgp_community    = [junos_policyoptions_community.testacc_policyOptions.name]
      bgp_origin       = "igp"
      family           = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface        = ["st0.0"]
      metric           = 5
      neighbor         = ["192.0.2.5"]
      next_hop         = ["192.0.2.5"]
      ospf_area        = "0.0.0.0"
      policy           = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
      preference       = 100
      protocol         = ["bgp"]
    }
    then {
      action          = "accept"
      as_path_expand  = "last-as count 1"
      as_path_prepend = "65000 65000"
      community {
        action = "set"
        value  = junos_policyoptions_community.testacc_policyOptions.name
      }
      community {
        action = "delete"
        value  = junos_policyoptions_community.testacc_policyOptions.name
      }
      community {
        action = "add"
        value  = junos_policyoptions_community.testacc_policyOptions.name
      }
      default_action = "reject"
      load_balance   = "per-packet"
      local_preference {
        action = "add"
        value  = 10
      }
      next     = "policy"
      next_hop = "192.0.2.4"
      metric {
        action = "add"
        value  = 10
      }
      origin = "igp"
      preference {
        action = "add"
        value  = 10
      }
    }
  }
}
resource "junos_policyoptions_policy_statement" "testacc_policyOptions2" {
  name = "testacc_policyOptions2"
  from {
    bgp_as_path_group = [junos_policyoptions_as_path_group.testacc_policyOptions.name]
  }
  to {
    bgp_as_path_group = [junos_policyoptions_as_path_group.testacc_policyOptions.name]
  }
  then {
    local_preference {
      action = "subtract"
      value  = 10
    }
    metric {
      action = "subtract"
      value  = 10
    }
    preference {
      action = "subtract"
      value  = 10
    }
    action = "accept"
  }
  term {
    name = "term"
    then {
      local_preference {
        action = "subtract"
        value  = 10
      }
      metric {
        action = "subtract"
        value  = 10
      }
      preference {
        action = "subtract"
        value  = 10
      }
    }
  }
}
resource "junos_policyoptions_policy_statement" "testacc_policyOptions3" {
  name                              = "testacc_policyOptions3"
  add_it_to_forwarding_table_export = true
  from {
    route_filter {
      route  = "192.0.2.0/25"
      option = "orlonger"
    }
  }
  then {
    load_balance = "per-packet"
  }
}
