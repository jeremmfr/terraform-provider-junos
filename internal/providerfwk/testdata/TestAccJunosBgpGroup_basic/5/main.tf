resource "junos_routing_options" "testacc_bgpgroup" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_bgp_group" "testacc_bgpgroup" {
  depends_on = [
    junos_routing_options.testacc_bgpgroup
  ]
  name                   = "testacc_bgpgroup"
  local_as               = "65000"
  local_as_alias         = true
  metric_out_minimum_igp = true
  family_evpn {}
  bgp_error_tolerance {
    no_malformed_route_limit = true
  }
}
