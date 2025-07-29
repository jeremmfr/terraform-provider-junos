resource "junos_routing_instance" "testacc_staticRoute" {
  name = "testacc_staticRoute"
}
resource "junos_routing_instance" "testacc_staticRoute2" {
  name = "testacc_staticRoute2"
}
resource "junos_static_route" "testacc_staticRoute_instance" {
  destination      = "192.0.2.0/25"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  receive          = true
  resolve          = true
}
resource "junos_static_route" "testacc_staticRoute_ipv6_default" {
  destination = "2001:db8:85a3::/50"
  receive     = true
  resolve     = true
}
resource "junos_static_route" "testacc_staticRoute_instance2" {
  destination      = "192.0.2.0/26"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  discard          = true
}
resource "junos_static_route" "testacc_staticRoute_ipv6_default2" {
  destination = "2001:db8:85a3::/52"
  discard     = true
}
resource "junos_static_route" "testacc_staticRoute_default" {
  destination = "192.0.2.0/27"
  reject      = true
}
resource "junos_static_route" "testacc_staticRoute_ipv6_instance" {
  destination      = "2001:db8:85a3::/54"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  reject           = true
}
resource "junos_static_route" "testacc_staticRoute_default2" {
  destination = "192.0.2.0/28"
  next_table  = "${junos_routing_instance.testacc_staticRoute2.name}.inet.0"
}
resource "junos_static_route" "testacc_staticRoute_ipv6_instance2" {
  destination      = "2001:db8:85a3::/56"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  next_table       = "${junos_routing_instance.testacc_staticRoute2.name}.inet6.0"
}
