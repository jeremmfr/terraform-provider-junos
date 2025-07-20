resource "junos_iccp" "testacc_iccp_peer" {
  local_ip_addr = "192.0.2.1"
}

resource "junos_iccp_peer" "testacc_iccp_peer" {
  depends_on = [junos_iccp.testacc_iccp_peer]

  ip_address               = "192.0.2.2"
  redundancy_group_id_list = [101, 100]

  liveness_detection {
    minimum_interval = 600
  }
}
