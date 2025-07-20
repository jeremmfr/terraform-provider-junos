resource "junos_routing_options" "testacc_bgpneighbor" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_routing_instance" "testacc_bgpneighbor2" {
  name = "testacc_bgpneighbor2"
  as   = "65000"
}
resource "junos_bgp_group" "testacc_bgpneighbor2" {
  name             = "testacc_bgpneighbor2"
  routing_instance = junos_routing_instance.testacc_bgpneighbor2.name
  type             = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor2" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  ip                            = "192.0.2.4"
  routing_instance              = junos_routing_instance.testacc_bgpneighbor2.name
  group                         = junos_bgp_group.testacc_bgpneighbor2.name
  advertise_external            = true
  accept_remote_nexthop         = true
  multihop                      = true
  local_as                      = "65000"
  local_as_no_prepend_global_as = true
  metric_out_minimum_igp_offset = -10
}
resource "junos_bgp_group" "testacc_bgpneighbor2b" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  name = "testacc_bgpneighbor2b"
  type = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor2b" {
  ip    = "192.0.2.5"
  group = junos_bgp_group.testacc_bgpneighbor2b.name
  family_evpn {
    accepted_prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
    prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
  }
  bgp_error_tolerance {
    malformed_route_limit         = 234
    malformed_update_log_interval = 567
  }
}
