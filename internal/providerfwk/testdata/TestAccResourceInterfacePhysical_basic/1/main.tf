resource "junos_interface_physical" "testacc_interface" {
  name        = var.interface
  description = "testacc_interface"
  disable     = true
  gigether_opts {
    ae_8023ad = var.interfaceAE
  }
}
resource "junos_interface_physical" "testacc_interface2" {
  name           = var.interface2
  description    = "testacc_interface2"
  hold_time_down = 6000
  hold_time_up   = 7000
  gigether_opts {
    flow_control     = true
    loopback         = true
    auto_negotiation = true
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  depends_on = [
    junos_interface_physical.testacc_interface,
  ]
  name        = var.interfaceAE
  description = "testacc_interfaceAE"
  parent_ether_opts {
    flow_control = true
    lacp {
      mode            = "active"
      admin_key       = 1
      periodic        = "slow"
      sync_reset      = "disable"
      system_id       = "00:00:01:00:01:00"
      system_priority = 250
    }
    loopback      = true
    link_speed    = "1g"
    minimum_links = 1
  }
  vlan_tagging = true
}
