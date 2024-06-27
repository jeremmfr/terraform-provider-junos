
resource "junos_vstp_vlan_group" "testacc_vstp_vlan_group" {
  name                   = "vlanGroup"
  vlan                   = ["10"]
  backup_bridge_priority = "8k"
  bridge_priority        = "4k"
  hello_time             = 2
}
resource "junos_routing_instance" "testacc_vstp_vlan_group" {
  name = "testacc_vstp_vlan_group"
  type = "virtual-switch"
}
resource "junos_vstp_vlan_group" "testacc_ri_vstp_vlan_group" {
  routing_instance       = junos_routing_instance.testacc_vstp_vlan_group.name
  name                   = "vlanGroupRI"
  vlan                   = ["12", "11"]
  backup_bridge_priority = "20k"
  forward_delay          = 22
  hello_time             = 3
  max_age                = 24
  system_identifier      = "00:aa:bc:ed:ff:11"
}
