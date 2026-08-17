resource "junos_iccp" "testacc_iccp_wo" {
  local_ip_addr                 = "192.0.2.1"
  authentication_key_wo         = "a@Key"
  authentication_key_wo_version = 1
}
