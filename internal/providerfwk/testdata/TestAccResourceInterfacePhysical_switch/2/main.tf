resource "junos_interface_physical" "testacc_interface" {
  name         = var.interface
  description  = "testacc_interfaceU"
  vlan_members = ["100"]
}
