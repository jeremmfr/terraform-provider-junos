resource "junos_interface_physical" "testacc_interface" {
  name        = var.interface
  description = "testacc_interface"
  gigether_opts {
    ae_8023ad = var.interfaceAE
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  name        = var.interfaceAE
  description = "testacc_interfaceAE"
  esi {
    identifier = "00:11:11:11:11:11:11:11:11:11"
    mode       = "all-active"
  }
  vlan_tagging = true
}
