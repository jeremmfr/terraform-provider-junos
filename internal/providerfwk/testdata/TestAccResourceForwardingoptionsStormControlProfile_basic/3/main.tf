resource "junos_forwardingoptions_storm_control_profile" "testacc_stormControlProfile" {
  name = "testacc_stormControlProfile @1"
  all {
    bandwidth_percentage    = 50
    no_registered_multicast = true
  }
}
