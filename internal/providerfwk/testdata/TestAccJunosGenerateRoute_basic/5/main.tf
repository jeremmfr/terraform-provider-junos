resource "junos_routing_instance" "testacc_generateRoute" {
  name = "testacc_generateRoute"
}
resource "junos_routing_instance" "testacc_generateRoute2" {
  name = "testacc_generateRoute2"
}

resource "junos_generate_route" "testacc_generateRoute" {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_generateRoute.name
  next_table       = "${junos_routing_instance.testacc_generateRoute2.name}.inet.0"
}
