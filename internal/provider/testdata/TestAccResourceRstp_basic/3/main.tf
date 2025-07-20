resource "junos_rstp" "testacc_rstp" {
  backup_bridge_priority                             = "60k"
  bridge_priority                                    = "4k"
  bpdu_destination_mac_address_provider_bridge_group = true
  extended_system_id                                 = 0
  force_version_stp                                  = true
  forward_delay                                      = 20
  hello_time                                         = 5
  max_age                                            = 22
  priority_hold_time                                 = 100
  system_id {
    id = "00:11:22:33:44:55"
  }
  system_id {
    id         = "00:22:33:44:55:66"
    ip_address = "192.0.2.4/24"
  }
  system_identifier             = "66:55:44:33:22:11"
  vpls_flush_on_topology_change = true
}
