resource "junos_interface_physical" "testacc_dataSecurityZone" {
  name         = var.interface
  description  = "testacc_dataSecurityZone"
  vlan_tagging = true
}
resource "junos_security_zone" "testacc_dataSecurityZone" {
  name                          = "testacc_dataSecurityZone"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_dataSecurityZone" {
  name = "testacc_dataSecurityZone"
  zone = junos_security_zone.testacc_dataSecurityZone.name
  cidr = "192.0.2.0/25"
}
resource "junos_interface_logical" "testacc_dataSecurityZone" {
  name                      = "${junos_interface_physical.testacc_dataSecurityZone.name}.100"
  description               = "testacc_dataSecurityZone"
  security_zone             = junos_security_zone.testacc_dataSecurityZone.name
  security_inbound_services = ["ssh"]
}
