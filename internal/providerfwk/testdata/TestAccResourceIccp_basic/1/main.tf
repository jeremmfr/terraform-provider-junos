resource "junos_iccp" "testacc_iccp" {
  local_ip_addr                   = "192.0.2.1"
  authentication_key              = "a@Key"
  session_establishment_hold_time = 300
}
