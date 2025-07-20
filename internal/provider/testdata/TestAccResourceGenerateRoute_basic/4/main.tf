resource "junos_routing_instance" "testacc_generateRoute" {
  name = "testacc_generateRoute"
}

resource "junos_generate_route" "testacc_generateRoute" {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_generateRoute.name
  passive          = true
  brief            = true
}
resource "junos_generate_route" "testacc_generateRoute6" {
  destination      = "2001:db8:85a3::/48"
  routing_instance = junos_routing_instance.testacc_generateRoute.name
  passive          = true
  brief            = true
}
