resource "junos_interface_physical" "testacc_interface" {
  name        = var.interface
  description = "testacc_interface"
  gigether_opts {
    ae_8023ad = var.interfaceAE
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  depends_on = [junos_iccp_peer.testacc_interfaceAE]

  name        = var.interfaceAE
  description = "testacc_interfaceAE"
  esi {
    identifier = "00:11:11:11:11:11:11:11:11:11"
    mode       = "all-active"
  }
  parent_ether_opts {
    lacp {
      mode      = "passive"
      admin_key = 1
      system_id = "ab:cd:ef:fe:dc:ba"
    }
    mc_ae {
      chassis_id     = 0
      mc_ae_id       = 200
      mode           = "active-standby"
      status_control = "active"
      events_iccp_peer_down {}
      revert_time     = 5
      switchover_mode = "revertive"
    }
  }
  vlan_tagging = true
}

resource "junos_iccp" "testacc_interfaceAE" {
  local_ip_addr = "192.0.2.1"
}

resource "junos_iccp_peer" "testacc_interfaceAE" {
  depends_on = [junos_iccp.testacc_interfaceAE]

  ip_address               = "192.0.2.2"
  redundancy_group_id_list = [300]

  liveness_detection {
    minimum_interval = 600
  }
}
