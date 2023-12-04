resource "junos_interface_physical" "testacc_interface" {
  name         = var.interface
  description  = "testacc_interfaceU"
  vlan_members = ["100"]
}

resource "junos_forwardingoptions_storm_control_profile" "testacc_interface" {
  name = "testacc interface"
  all {}
}
