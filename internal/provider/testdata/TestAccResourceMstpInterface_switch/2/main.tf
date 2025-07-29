resource "junos_interface_physical" "testacc_mstp_interface" {
  name         = var.interface
  vlan_members = ["default"]
}

resource "junos_mstp_interface" "testacc_mstp_interface" {
  name         = junos_interface_physical.testacc_mstp_interface.name
  no_root_port = true
}

resource "junos_routing_instance" "testacc_mstp_interface" {
  name = "testacc_mstp_interface"
  type = "virtual-switch"
}

resource "junos_mstp_interface" "all" {
  name                      = "all"
  routing_instance          = junos_routing_instance.testacc_mstp_interface.name
  access_trunk              = true
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 16
  edge                      = true
  mode                      = "shared"
  priority                  = 240
}
