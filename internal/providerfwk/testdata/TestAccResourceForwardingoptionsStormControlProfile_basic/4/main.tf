resource "junos_forwardingoptions_storm_control_profile" "testacc_stormControlProfile" {
  name = "testacc_stormControlProfile @1"
  all {
    no_unregistered_multicast = true
  }
}
