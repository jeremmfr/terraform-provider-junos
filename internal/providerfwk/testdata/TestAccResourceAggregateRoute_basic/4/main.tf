resource "junos_routing_instance" "testacc_aggregateRoute" {
  name = "testacc_aggregateRoute"
}

resource "junos_aggregate_route" "testacc_aggregateRoute" {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_aggregateRoute.name
  passive          = true
  brief            = true
}
resource "junos_aggregate_route" "testacc_aggregateRoute6" {
  destination      = "2001:db8:85a3::/48"
  routing_instance = junos_routing_instance.testacc_aggregateRoute.name
  passive          = true
  brief            = true
}
