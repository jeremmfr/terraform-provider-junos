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
  ip                     = "192.0.2.4"
  routing_instance       = junos_routing_instance.testacc_bgpneighbor2.name
  group                  = junos_bgp_group.testacc_bgpneighbor2.name
  local_as               = "65000"
  local_as_alias         = true
  metric_out_minimum_igp = true
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
  family_evpn {}
  bgp_error_tolerance {
    no_malformed_route_limit = true
  }
}
