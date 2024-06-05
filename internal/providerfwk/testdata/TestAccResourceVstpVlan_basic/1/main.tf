resource "junos_vstp_vlan" "testacc_vstp_vlan" {
  vlan_id = "10"
}
resource "junos_vstp_vlan" "testacc_vstp_vlan_all" {
  vlan_id = "all"
}
resource "junos_routing_instance" "testacc_vstp_vlan" {
  name = "testacc_vstp_vlan"
  type = "virtual-switch"
}
resource "junos_vstp_vlan" "testacc_ri_vstp_vlan" {
  routing_instance = junos_routing_instance.testacc_vstp_vlan.name
  vlan_id          = "11"
  bridge_priority  = "16k"
}
resource "junos_vstp_vlan" "testacc_ri_vstp_vlan_all" {
  routing_instance  = junos_routing_instance.testacc_vstp_vlan.name
  vlan_id           = "all"
  system_identifier = "00:aa:bc:ed:ff:11"
}
