resource "junos_system_ntp_server" "testacc_ntpServer" {
  address = "192.0.2.1"
  prefer  = true
  version = 4
  key     = 1
}
