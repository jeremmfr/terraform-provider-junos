resource "junos_policyoptions_policy_statement" "testacc_routing_options" {
  name                              = " testacc_routing_options "
  add_it_to_forwarding_table_export = true
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
  clean_on_destroy                         = true
  forwarding_table_export_configure_singly = true
  forwarding_table {
    no_ecmp_fast_reroute                         = true
    no_indirect_next_hop                         = true
    no_indirect_next_hop_change_acknowledgements = true
    unicast_reverse_path                         = "feasible-paths"
  }
  graceful_restart {}
  instance_export = [
    junos_policyoptions_policy_statement.testacc_routing_options3.name,
    junos_policyoptions_policy_statement.testacc_routing_options2.name,
  ]
  instance_import = [
    junos_policyoptions_policy_statement.testacc_routing_options2.name,
    junos_policyoptions_policy_statement.testacc_routing_options3.name,
  ]
}
