resource "junos_iccp" "testacc_iccp_peer" {
  local_ip_addr = "192.0.2.1"
}

resource "junos_iccp_peer" "testacc_iccp_peer" {
  depends_on = [junos_iccp.testacc_iccp_peer]

  ip_address               = "192.0.2.2"
  redundancy_group_id_list = [101, 100]

  authentication_key              = "a@key"
  session_establishment_hold_time = 300

  liveness_detection {
    minimum_receive_interval           = 600
    transmit_interval_minimum_interval = 600

    detection_time_threshold    = 1800
    multiplier                  = 2
    no_adaptation               = true
    transmit_interval_threshold = 1800
    version                     = "automatic"
  }
  backup_liveness_detection {
    backup_peer_ip = "192.0.2.3"
  }
}
