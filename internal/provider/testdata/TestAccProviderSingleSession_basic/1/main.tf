resource "junos_static_route" "testacc_provider_single_connection_1" {
  destination = "192.0.2.128/28"
  next_hop = [
    "192.0.2.254"
  ]
}
resource "junos_static_route" "testacc_provider_single_connection_2" {
  destination = "192.0.2.144/28"
  next_hop = [
    "192.0.2.254"
  ]
}
resource "junos_static_route" "testacc_provider_single_connection_3" {
  destination = "192.0.2.160/28"
  next_hop = [
    "192.0.2.254"
  ]
}
resource "junos_static_route" "testacc_provider_single_connection_4" {
  destination = "192.0.2.176/28"
  next_hop = [
    "192.0.2.254"
  ]
}
resource "junos_static_route" "testacc_provider_single_connection_5" {
  destination = "192.0.2.192/28"
  next_hop = [
    "192.0.2.254"
  ]
}
