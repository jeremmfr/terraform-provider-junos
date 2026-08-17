resource "junos_iccp" "testacc_iccp_peer_wo" {
  local_ip_addr = "192.0.2.1"
}

resource "junos_iccp_peer" "testacc_iccp_peer_wo" {
  depends_on = [junos_iccp.testacc_iccp_peer_wo]

  ip_address               = "192.0.2.2"
  redundancy_group_id_list = [101, 100]

  authentication_key_wo         = "a@Key"
  authentication_key_wo_version = 1

  liveness_detection {
    minimum_interval = 600
  }
}
