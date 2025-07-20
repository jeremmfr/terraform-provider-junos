resource "junos_routing_options" "testacc_bgpneighbor" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_routing_instance" "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  as   = "65000"
}
resource "junos_bgp_group" "testacc_bgpneighbor" {
  name             = "testacc_bgpneighbor"
  routing_instance = junos_routing_instance.testacc_bgpneighbor.name
  type             = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  ip                              = "192.0.2.4"
  routing_instance                = junos_routing_instance.testacc_bgpneighbor.name
  group                           = junos_bgp_group.testacc_bgpneighbor.name
  description                     = "peer 2.4"
  advertise_external_conditional  = true
  keep_none                       = true
  no_advertise_peer_as            = true
  metric_out_igp_offset           = -10
  metric_out_igp_delay_med_update = true
  authentication_key              = "password"
  bgp_multipath {
    multiple_as = true
  }
  graceful_restart {
    restart_time     = 10
    stale_route_time = 10
  }
  tcp_aggressive_transmission = true
  bgp_error_tolerance {}
}
