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
  parent_ether_opts {
    lacp {
      mode      = "active"
      admin_key = 1
      system_id = "ab:cd:ef:fe:dc:ba"
    }
    mc_ae {
      chassis_id     = 0
      mc_ae_id       = 200
      mode           = "active-active"
      status_control = "active"
    }
  }
  vlan_tagging = true
}
