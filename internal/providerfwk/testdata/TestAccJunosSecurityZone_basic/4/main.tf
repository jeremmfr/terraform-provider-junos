
resource "junos_security_zone" "testacc_securityZone" {
  name = "testacc_securityZone"
}
resource "junos_interface_logical" "testacc_securityZone" {
  name                       = "${var.interface}.0"
  security_zone              = junos_security_zone.testacc_securityZone.name
  security_inbound_protocols = ["bgp"]
  security_inbound_services  = ["ssh"]
}
