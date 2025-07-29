resource "junos_interface_physical" "testacc_nullcommitfile" {
  name         = var.interface
  description  = "testacc_null"
  vlan_tagging = true
}

data "junos_interface_physical" "testacc_nullcommitfile" {
  config_interface = var.interface
}

