resource "junos_forwardingoptions_storm_control_profile" "testacc_stormControlProfile" {
  name            = "testacc_stormControlProfile @1"
  action_shutdown = true
  all {
    bandwidth_level    = 10240
    burst_size         = 10240
    no_broadcast       = true
    no_multicast       = true
    no_unknown_unicast = true
  }
}
