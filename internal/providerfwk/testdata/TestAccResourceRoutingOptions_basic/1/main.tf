resource "junos_policyoptions_policy_statement" "testacc_routing_options" {
  lifecycle {
    create_before_destroy = true
  }
  name = " testacc_routing_options "
  from {
    protocol = ["bgp"]
    route_filter {
      route  = "192.0.2.0/28"
      option = "orlonger"
    }
  }
  then {
    load_balance = "per-packet"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_routing_options2" {
  name = " testacc_routing_options #2 "
  from {
    route_filter {
      route  = "192.0.2.0/28"
      option = "orlonger"
    }
  }
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_routing_options3" {
  name = " testacc_routing_options #3@_@_#_#_@_# "
  from {
    route_filter {
      route  = "192.0.2.16/28"
      option = "orlonger"
    }
  }
  then {
    action = "accept"
  }
}

resource "junos_routing_options" "testacc_routing_options" {
  autonomous_system {
    number         = "65000"
    asdot_notation = true
    loops          = 5
  }
  forwarding_table {
    dynamic_list_next_hop                     = true
    ecmp_fast_reroute                         = true
    export                                    = [junos_policyoptions_policy_statement.testacc_routing_options.name]
    indirect_next_hop                         = true
    indirect_next_hop_change_acknowledgements = true
    krt_nexthop_ack_timeout                   = 200
    remnant_holdtime                          = 0
    unicast_reverse_path                      = "active-paths"
  }
  graceful_restart {
    restart_duration = 120
    disable          = true
  }
  instance_export = [junos_policyoptions_policy_statement.testacc_routing_options2.name]
  instance_import = [junos_policyoptions_policy_statement.testacc_routing_options3.name]
  ipv6_router_id  = "2001:db8::4"
  router_id       = "192.0.2.4"
}
