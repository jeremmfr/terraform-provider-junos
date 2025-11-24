resource "junos_interface_physical" "testacc_actioncommitfile" {
  name         = var.interface
  description  = "testacc_null"
  vlan_tagging = true
}

data "junos_interface_physical" "testacc_actioncommitfile" {
  config_interface = var.interface
}
