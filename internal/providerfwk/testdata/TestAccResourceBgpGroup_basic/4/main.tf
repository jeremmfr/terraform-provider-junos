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
  name                          = "testacc_bgpgroup"
  advertise_external            = true
  accept_remote_nexthop         = true
  multihop                      = true
  local_as                      = "65000"
  local_as_no_prepend_global_as = true
  metric_out_minimum_igp_offset = -10
  type                          = "internal"
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
