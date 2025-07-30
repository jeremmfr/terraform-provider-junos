resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v4_default" {
  active_leasequery {}
  bulk_leasequery {}
}

resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v6_default" {
  version = "v6"
}

resource "junos_routing_instance" "testacc_dhcprelay" {
  name = "testacc_dhcprelay"
}
resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v4_ri" {
  routing_instance = junos_routing_instance.testacc_dhcprelay.name
  leasequery {}
  overrides_v4 {
    always_write_option_82 = true
  }
  relay_option_82 {
    circuit_id {}
  }
}
resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v6_ri" {
  routing_instance = junos_routing_instance.testacc_dhcprelay.name
  version          = "v6"

  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}
