resource "junos_iccp" "testacc_multichassis_protection_peer" {
  local_ip_addr = "192.0.2.1"
}

resource "junos_iccp_peer" "testacc_multichassis_protection_peer" {
  depends_on = [junos_iccp.testacc_multichassis_protection_peer]

  ip_address               = "192.0.2.2"
  redundancy_group_id_list = [101, 100]

  liveness_detection {
    minimum_interval = 600
  }
}

resource "junos_multichassis" "testacc_multichassis_protection_peer" {
  clean_on_destroy = true
}

resource "junos_multichassis_protection_peer" "testacc_multichassis_protection_peer" {
  depends_on     = [junos_multichassis.testacc_multichassis_protection_peer]
  ip_address     = junos_iccp_peer.testacc_multichassis_protection_peer.ip_address
  interface      = var.interface
  icl_down_delay = 60
}
