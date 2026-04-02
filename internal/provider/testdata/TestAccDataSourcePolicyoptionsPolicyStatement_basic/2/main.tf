resource "junos_policyoptions_as_path" "testacc_dataPolicyStatement" {
  name = "testacc_dataPolicyStatement"
  path = "5|12|18"
}
resource "junos_policyoptions_as_path_group" "testacc_dataPolicyStatement" {
  name = "testacc_dataPolicyStatement"
  as_path {
    name = "testacc_dataPolicyStatement"
    path = "5|12|18"
  }
}
resource "junos_policyoptions_community" "testacc_dataPolicyStatement" {
  name    = "testacc_dataPolicyStatement"
  members = ["65000:100"]
}
resource "junos_policyoptions_prefix_list" "testacc_dataPolicyStatement" {
  name   = "testacc_dataPolicyStatement"
  prefix = ["192.0.2.0/25"]
}
resource "junos_routing_instance" "testacc_dataPolicyStatement" {
  name = "testacc_dataPolicyStatement"
}
resource "junos_policyoptions_policy_statement" "testacc_dataPolicyStatementRef" {
  name = "testacc_dataPolicyStatementRef"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_dataPolicyStatement" {
  name = "testacc_dataPolicyStatement"
  from {
    aggregate_contributor = true
    bgp_as_path           = [junos_policyoptions_as_path.testacc_dataPolicyStatement.name]
    bgp_as_path_calc_length {
      count = 4
      match = "orhigher"
    }
    bgp_as_path_calc_length {
      count = 3
      match = "equal"
    }
    bgp_community = [junos_policyoptions_community.testacc_dataPolicyStatement.name]
    bgp_community_count {
      count = 6
      match = "orhigher"
    }
    bgp_community_count {
      count = 5
      match = "equal"
    }
    bgp_origin       = "igp"
    color            = 31
    family           = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_dataPolicyStatement.name
    interface        = ["st0.0"]
    metric           = 5
    neighbor         = ["192.0.2.4"]
    next_hop         = ["192.0.2.4"]
    ospf_area        = "0.0.0.0"
    preference       = 100
    prefix_list      = [junos_policyoptions_prefix_list.testacc_dataPolicyStatement.name]
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
    bgp_as_path      = [junos_policyoptions_as_path.testacc_dataPolicyStatement.name]
    bgp_community    = [junos_policyoptions_community.testacc_dataPolicyStatement.name]
    bgp_origin       = "igp"
    family           = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_dataPolicyStatement.name
    interface        = ["st0.0"]
    metric           = 5
    neighbor         = ["192.0.2.5"]
    next_hop         = ["192.0.2.5"]
    ospf_area        = "0.0.0.0"
    policy           = [junos_policyoptions_policy_statement.testacc_dataPolicyStatementRef.name]
    preference       = 100
    protocol         = ["ospf"]
  }
  then {
    action          = "accept"
    as_path_expand  = "65000 65000"
    as_path_prepend = "65000 65000"
    community {
      action = "set"
      value  = junos_policyoptions_community.testacc_dataPolicyStatement.name
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
      bgp_as_path           = [junos_policyoptions_as_path.testacc_dataPolicyStatement.name]
      bgp_as_path_unique_count {
        count = 4
        match = "orlower"
      }
      bgp_community    = [junos_policyoptions_community.testacc_dataPolicyStatement.name]
      bgp_origin       = "igp"
      family           = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_dataPolicyStatement.name
      interface        = ["st0.0"]
      metric           = 5
      neighbor         = ["192.0.2.4"]
      next_hop         = ["192.0.2.4"]
      ospf_area        = "0.0.0.0"
      policy           = [junos_policyoptions_policy_statement.testacc_dataPolicyStatementRef.name]
      preference       = 100
      prefix_list      = [junos_policyoptions_prefix_list.testacc_dataPolicyStatement.name]
      protocol         = ["bgp"]
      route_filter {
        route  = "192.0.2.0/25"
        option = "exact"
      }
    }
    to {
      bgp_as_path      = [junos_policyoptions_as_path.testacc_dataPolicyStatement.name]
      bgp_community    = [junos_policyoptions_community.testacc_dataPolicyStatement.name]
      bgp_origin       = "igp"
      family           = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_dataPolicyStatement.name
      interface        = ["st0.0"]
      metric           = 5
      neighbor         = ["192.0.2.5"]
      next_hop         = ["192.0.2.5"]
      ospf_area        = "0.0.0.0"
      policy           = [junos_policyoptions_policy_statement.testacc_dataPolicyStatementRef.name]
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
resource "junos_policyoptions_policy_statement" "testacc_dataPolicyStatement2" {
  name = "testacc_dataPolicyStatement2"
  from {
    bgp_as_path_group = [junos_policyoptions_as_path_group.testacc_dataPolicyStatement.name]
  }
  to {
    bgp_as_path_group = [junos_policyoptions_as_path_group.testacc_dataPolicyStatement.name]
  }
  then {
    action = "accept"
    local_preference {
      action = "add"
      value  = 10
    }
    metric {
      action = "add"
      value  = 10
    }
    preference {
      action = "add"
      value  = 10
    }
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
resource "junos_policyoptions_policy_statement" "testacc_dataPolicyStatement3" {
  name = "testacc_dataPolicyStatement3"
  from {
    bgp_as_path_unique_count {
      count = 3
      match = "equal"
    }
    bgp_as_path_unique_count {
      count = 2
      match = "orhigher"
    }
    next_hop_type_merged = true
    next_hop_weight {
      match  = "greater-than-equal"
      weight = 500
    }
    next_hop_weight {
      match  = "equal"
      weight = 200
    }
  }
  then {
    action = "accept"
  }
  term {
    name = "term"
    from {
      bgp_as_path_calc_length {
        count = 4
        match = "orhigher"
      }
      bgp_as_path_unique_count {
        count = 4
        match = "orlower"
      }
    }
    then {
      action = "reject"
    }
  }
}

data "junos_policyoptions_policy_statement" "testacc_dataPolicyStatement" {
  name = junos_policyoptions_policy_statement.testacc_dataPolicyStatement.name
}
data "junos_policyoptions_policy_statement" "testacc_dataPolicyStatement2" {
  name = junos_policyoptions_policy_statement.testacc_dataPolicyStatement2.name
}
data "junos_policyoptions_policy_statement" "testacc_dataPolicyStatement3" {
  name = junos_policyoptions_policy_statement.testacc_dataPolicyStatement3.name
}
