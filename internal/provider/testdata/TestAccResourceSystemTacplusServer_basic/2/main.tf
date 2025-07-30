resource "junos_routing_instance" "testacc_tacplusServer" {
  name = "testacc_tacplusServer"
}
resource "junos_system_tacplus_server" "testacc_tacplusServer" {
  address           = "192.0.2.1"
  secret            = "password"
  source_address    = "192.0.2.2"
  port              = 49
  timeout           = 10
  single_connection = true
  routing_instance  = junos_routing_instance.testacc_tacplusServer.name
}
