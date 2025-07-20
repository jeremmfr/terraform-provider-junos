resource "junos_vstp_interface" "testacc_vstp_interface" {
  name                      = "all"
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 101
  edge                      = true
  priority                  = 32
}
resource "junos_vstp_vlan" "testacc_vstp_interface2" {
  vlan_id = "10"
}
resource "junos_interface_physical" "testacc_vstp_interface2" {
  name         = var.interface
  description  = "testacc_vstp_interface2"
  vlan_members = ["15"]
}
resource "junos_vstp_interface" "testacc_vstp_interface2" {
  name         = junos_interface_physical.testacc_vstp_interface2.name
  access_trunk = true
  mode         = "shared"
  no_root_port = true
  vlan         = junos_vstp_vlan.testacc_vstp_interface2.vlan_id
}
resource "junos_vstp_vlan_group" "testacc_vstp_interface3" {
  name = "testacc_vstp_interface2"
  vlan = ["11"]
}
resource "junos_vstp_interface" "testacc_vstp_interface3" {
  name       = "all"
  priority   = 32
  vlan_group = junos_vstp_vlan_group.testacc_vstp_interface3.name
}
resource "junos_routing_instance" "testacc_vstp_interface" {
  name = "testacc_vstp_intface"
  type = "virtual-switch"
}
resource "junos_vstp_interface" "testacc_vstp_interface4" {
  name             = "all"
  edge             = true
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
}
resource "junos_vstp_vlan" "testacc_vstp_interface5" {
  vlan_id          = "all"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
}
resource "junos_vstp_interface" "testacc_vstp_interface5" {
  name             = "all"
  mode             = "point-to-point"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan             = junos_vstp_vlan.testacc_vstp_interface5.vlan_id
}
resource "junos_vstp_vlan_group" "testacc_vstp_interface6" {
  name             = "testacc_vstp_interface6"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan             = ["13"]
}
resource "junos_vstp_interface" "testacc_vstp_interface6" {
  name             = var.interface
  no_root_port     = true
  priority         = 64
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan_group       = junos_vstp_vlan_group.testacc_vstp_interface6.name
}
