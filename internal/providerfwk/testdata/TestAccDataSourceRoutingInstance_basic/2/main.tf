resource "junos_interface_physical" "testacc_dataRoutingInstance" {
  name         = var.interface
  description  = "testacc_dataRoutingInstance"
  vlan_tagging = true
}
resource "junos_routing_instance" "testacc_dataRoutingInstance" {
  name = "testacc_dataRoutingInstance"
}
resource "junos_interface_logical" "testacc_dataRoutingInstance" {
  name             = "${junos_interface_physical.testacc_dataRoutingInstance.name}.100"
  description      = "testacc_dataRoutingInstance"
  routing_instance = junos_routing_instance.testacc_dataRoutingInstance.name
}

data "junos_routing_instance" "testacc_dataRoutingInstance" {
  name = "testacc_dataRoutingInstance"
}
