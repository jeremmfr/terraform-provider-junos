resource "junos_routing_instance" "testacc_policyOptions" {
  name = "testacc_policyOptions"
}
resource "junos_policyoptions_as_path" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  path = "5|15"
}
resource "junos_policyoptions_as_path_group" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  as_path {
    name = "testacc_policyOptions"
    path = "5|15"
  }
}
resource "junos_policyoptions_community" "testacc_policyOptions" {
  name         = "testacc_policyOptions"
  members      = ["65000:200"]
  invert_match = true
}
resource "junos_policyoptions_prefix_list" "testacc_policyOptions" {
  name   = "testacc_policyOptions"
  prefix = ["192.0.2.0/26", "192.0.2.64/26"]
}
resource "junos_policyoptions_prefix_list" "testacc_policyOptions2" {
  name       = "testacc_policyOptions2"
  apply_path = "system radius-server <*>"
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
    bgp_as_path_calc_length {
      count = 3
      match = "equal"
    }
    bgp_as_path_unique_count {
      count = 3
      match = "equal"
    }
    bgp_as_path_unique_count {
      count = 2
      match = "orhigher"
    }
    bgp_community = [junos_policyoptions_community.testacc_policyOptions.name]
    bgp_community_count {
      count = 6
      match = "orhigher"
    }
    bgp_community_count {
      count = 5
      match = "equal"
    }
    bgp_origin             = "igp"
    bgp_srte_discriminator = 30

    evpn_esi             = ["00:11:11:11:11:11:11:11:11:33", "00:11:11:11:11:11:11:11:11:32"]
    evpn_mac_route       = "mac-only"
    evpn_tag             = [36, 35, 33]
    family               = "evpn"
    local_preference     = 100
    routing_instance     = junos_routing_instance.testacc_policyOptions.name
    interface            = ["st0.0"]
    metric               = 5
    neighbor             = ["192.0.2.4"]
    next_hop             = ["192.0.2.4"]
    next_hop_type_merged = true
    next_hop_weight {
      match  = "greater-than-equal"
      weight = 500
    }
    next_hop_weight {
      match  = "equal"
      weight = 200
    }
    ospf_area  = "0.0.0.0"
    preference = 100
    prefix_list = [junos_policyoptions_prefix_list.testacc_policyOptions.name,
      junos_policyoptions_prefix_list.testacc_policyOptions2.name,
    ]
    protocol = ["bgp"]
    route_filter {
      route  = "192.0.2.0/25"
      option = "exact"
    }
    route_type          = "internal"
    srte_color          = 39
    state               = "active"
    tunnel_type         = ["ipip"]
    validation_database = "valid"
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
    protocol         = ["ospf"]
  }
  then {
    action          = "accept"
    as_path_expand  = "65000 65000"
    as_path_prepend = "65000 65000"
    community {
      action = "set"
      value  = junos_policyoptions_community.testacc_policyOptions.name
    }
    default_action = "reject"
    load_balance   = "per-packet"
    next           = "policy"
    next_hop       = "192.0.2.4"
    origin         = "igp"
  }
  term {
    name = "term"
    from {
      aggregate_contributor = true
      bgp_as_path           = [junos_policyoptions_as_path.testacc_policyOptions.name]
      bgp_as_path_unique_count {
        count = 4
        match = "orlower"
      }
      bgp_community    = [junos_policyoptions_community.testacc_policyOptions.name]
      bgp_origin       = "igp"
      family           = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface        = ["st0.0"]
      metric           = 5
      neighbor         = ["192.0.2.4"]
      next_hop         = ["192.0.2.4"]
      ospf_area        = "0.0.0.0"
      policy           = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
      preference       = 100
      prefix_list      = [junos_policyoptions_prefix_list.testacc_policyOptions.name]
      protocol         = ["bgp"]
      route_filter {
        route  = "192.0.2.0/25"
        option = "exact"
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
      protocol         = ["ospf"]
    }
    then {
      action          = "accept"
      as_path_expand  = "last-as count 1"
      as_path_prepend = "65000 65000"
      default_action  = "accept"
      load_balance    = "per-packet"
      next            = "policy"
      next_hop        = "192.0.2.4"
      origin          = "igp"
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
      action = "none"
      value  = 10
    }
    metric {
      action = "none"
      value  = 10
    }
    preference {
      action = "none"
      value  = 10
    }
    action = "accept"
  }
  term {
    name = "term"
    then {
      local_preference {
        action = "none"
        value  = 10
      }
      metric {
        action = "none"
        value  = 10
      }
      preference {
        action = "none"
        value  = 10
      }
    }
  }
}
