data "junos_interface_physical" "testacc_interface" {
  config_interface = var.interface
}
resource "junos_interface_physical" "testacc_interface" {
  name          = var.interface
  description   = "testacc_interface"
  storm_control = junos_forwardingoptions_storm_control_profile.testacc_interface.name
  trunk         = true
  vlan_native   = 100
  vlan_members  = ["100-110"]
}

resource "junos_forwardingoptions_storm_control_profile" "testacc_interface" {
  name = "testacc interface"
  all {}
}
