resource "junos_interface_physical" "testacc_datainterfaceP" {
  name         = var.interface
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}
