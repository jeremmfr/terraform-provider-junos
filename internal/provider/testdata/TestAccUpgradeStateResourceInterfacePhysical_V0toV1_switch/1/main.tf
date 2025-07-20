resource "junos_interface_physical" "testacc_interface" {
  name         = var.interface
  description  = "testacc_interface"
  trunk        = true
  vlan_native  = 100
  vlan_members = ["100-110"]
}
