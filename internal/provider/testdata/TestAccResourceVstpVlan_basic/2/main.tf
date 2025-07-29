resource "junos_vstp_vlan" "testacc_vstp_vlan" {
  vlan_id                = "10"
  backup_bridge_priority = "8k"
  bridge_priority        = "4k"
}
resource "junos_vstp_vlan" "testacc_vstp_vlan_all" {
  vlan_id    = "all"
  hello_time = 2
}
resource "junos_routing_instance" "testacc_vstp_vlan" {
  name = "testacc_vstp_vlan"
  type = "virtual-switch"
}
resource "junos_vstp_vlan" "testacc_ri_vstp_vlan" {
  routing_instance       = junos_routing_instance.testacc_vstp_vlan.name
  vlan_id                = "11"
  backup_bridge_priority = "20k"
  bridge_priority        = 0
  forward_delay          = 22
  hello_time             = 3
  max_age                = 24
}
resource "junos_vstp_vlan" "testacc_ri_vstp_vlan_all" {
  routing_instance = junos_routing_instance.testacc_vstp_vlan.name
  vlan_id          = "all"
}
