resource "junos_routing_instance" "testacc_ntpServer" {
  name = "testacc_ntpServer"
}
resource "junos_system_ntp_server" "testacc_ntpServer" {
  address          = "192.0.2.1"
  prefer           = true
  version          = 4
  routing_instance = junos_routing_instance.testacc_ntpServer.name
}
