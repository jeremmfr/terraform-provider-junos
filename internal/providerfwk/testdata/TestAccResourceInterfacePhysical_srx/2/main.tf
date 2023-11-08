resource "junos_interface_physical" "testacc_interface" {
  name                  = var.interface
  description           = "testacc_interface2"
  encapsulation         = "vlan-vpls"
  speed                 = "1g"
  flexible_vlan_tagging = true
}
