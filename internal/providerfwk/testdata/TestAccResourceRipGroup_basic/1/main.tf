resource "junos_rip_group" "testacc_ripgroup" {
  name = "test_rip_group#1"
}
resource "junos_routing_instance" "testacc_ripgroup2" {
  name = "testacc_ripgroup2"
}
resource "junos_rip_group" "testacc_ripgroup2" {
  name             = "test_rip_group#2"
  routing_instance = junos_routing_instance.testacc_ripgroup2.name
}
resource "junos_rip_group" "testacc_ripnggroup" {
  name = "test_ripng_group#1"
  ng   = true
}
resource "junos_rip_group" "testacc_ripnggroup2" {
  name             = "test_ripng_group#2"
  ng               = true
  routing_instance = junos_routing_instance.testacc_ripgroup2.name
}
