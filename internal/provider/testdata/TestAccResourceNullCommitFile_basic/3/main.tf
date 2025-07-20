resource "junos_interface_physical" "testacc_nullcommitfile" {
  provider = junos.fake

  name         = var.interface
  description  = "testacc_null"
  vlan_tagging = true
}
