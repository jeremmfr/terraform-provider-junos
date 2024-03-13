resource "junos_routing_instance" "testacc_tacplusServer" {
  name = "testacc_tacplusServer"
}
resource "junos_system_tacplus_server" "testacc_tacplusServer" {
  address = "192.0.2.1"
  port    = 49
}
