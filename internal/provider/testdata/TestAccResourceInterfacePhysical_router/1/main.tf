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
    identifier = "00:01:11:11:11:11:11:11:11:11"
    mode       = "all-active"
  }
  parent_ether_opts {
    source_address_filter = ["00:11:22:33:44:55"]
    source_filtering      = true
  }
  vlan_tagging = true
}
