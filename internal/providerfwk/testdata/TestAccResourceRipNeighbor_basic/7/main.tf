resource "junos_rip_group" "testacc_ripneigh" {
  name = "test_rip_group#1"
}
resource "junos_routing_instance" "testacc_ripneigh2" {
  name = "testacc_ripneigh2"
}
resource "junos_rip_group" "testacc_ripneigh2" {
  name             = "test_rip_group#2"
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
}
resource "junos_rip_group" "testacc_ripngneigh" {
  name = "test_ripng_group#1"
  ng   = true
}
resource "junos_rip_group" "testacc_ripngneigh2" {
  name             = "test_ripng_group#2"
  ng               = true
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
}
resource "junos_rip_neighbor" "testacc_ripneigh" {
  name                = "ae0.0"
  group               = junos_rip_group.testacc_ripneigh.name
  authentication_type = "none"
}
resource "junos_rip_neighbor" "testacc_ripngneigh" {
  name  = "ae0.0"
  ng    = true
  group = junos_rip_group.testacc_ripngneigh.name
}
resource "junos_rip_neighbor" "testacc_ripneigh_all" {
  name  = "all"
  group = junos_rip_group.testacc_ripneigh.name
}
resource "junos_interface_physical" "testacc_ripneigh2" {
  name = var.interface
}
resource "junos_interface_logical" "testacc_ripneigh2" {
  name             = "${junos_interface_physical.testacc_ripneigh2.name}.0"
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
}
resource "junos_rip_neighbor" "testacc_ripneigh2" {
  name             = junos_interface_logical.testacc_ripneigh2.name
  group            = junos_rip_group.testacc_ripneigh2.name
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
}
resource "junos_rip_neighbor" "testacc_ripngneigh2" {
  name             = junos_interface_logical.testacc_ripneigh2.name
  ng               = true
  group            = junos_rip_group.testacc_ripngneigh2.name
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
}
