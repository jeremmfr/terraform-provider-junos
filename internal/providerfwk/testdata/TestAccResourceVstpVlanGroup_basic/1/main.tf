resource "junos_vstp_vlan_group" "testacc_vstp_vlan_group" {
  name = "vlanGroup"
  vlan = ["10"]
}
resource "junos_routing_instance" "testacc_vstp_vlan_group" {
  name = "testacc_vstp_vlan_group"
  type = "virtual-switch"
}
resource "junos_vstp_vlan_group" "testacc_ri_vstp_vlan_group" {
  routing_instance = junos_routing_instance.testacc_vstp_vlan_group.name
  name             = "vlanGroupRI"
  vlan             = ["12"]
  bridge_priority  = "16k"
}
