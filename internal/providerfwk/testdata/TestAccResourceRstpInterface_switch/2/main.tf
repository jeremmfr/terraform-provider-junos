resource "junos_rstp_interface" "all" {
  name                      = "all"
  access_trunk              = true
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 16
  edge                      = true
  mode                      = "shared"
  priority                  = 240
}

resource "junos_interface_physical" "testacc_rstp_interface" {
  name         = var.interface
  vlan_members = ["default"]
}

resource "junos_rstp_interface" "testacc_rstp_interface" {
  name         = junos_interface_physical.testacc_rstp_interface.name
  no_root_port = true
}

resource "junos_routing_instance" "testacc_rstp_interface" {
  name = "testacc_rstp_interface"
  type = "virtual-switch"
}

resource "junos_rstp_interface" "all2" {
  name                      = "all"
  routing_instance          = junos_routing_instance.testacc_rstp_interface.name
  access_trunk              = true
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 16
  edge                      = true
  mode                      = "shared"
  priority                  = 240
}
