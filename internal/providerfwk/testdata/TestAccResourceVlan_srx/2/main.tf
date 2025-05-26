resource "junos_interface_physical" "testacc_vlan" {
  name         = var.interface
  vlan_tagging = true
}

resource "junos_interface_logical" "testacc_vlan" {
  name    = "${junos_interface_physical.testacc_vlan.name}.10"
  vlan_id = 10
}

resource "junos_interface_logical" "testacc_vlan2" {
  name    = "${junos_interface_physical.testacc_vlan.name}.11"
  vlan_id = 11
}


resource "junos_vlan" "testacc_vlan" {
  name    = "vlan10"
  vlan_id = 10
  interface = [
    junos_interface_logical.testacc_vlan2.name,
    junos_interface_logical.testacc_vlan.name,
  ]
}
