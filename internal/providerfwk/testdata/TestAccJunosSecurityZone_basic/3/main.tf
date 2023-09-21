resource "junos_security_zone" "testacc_securityZone" {
  name = "testacc_securityZone"
  address_book {
    name    = "testacc_zone1"
    network = "192.0.2.0/25"
  }
  address_book {
    name    = "testacc_zone2"
    network = "192.0.2.128/25"
  }
  address_book_set {
    name    = "testacc_zoneSet"
    address = ["testacc_zone1", "testacc_zone2"]
  }
  inbound_protocols = ["bgp"]
  inbound_services  = ["ssh"]
}
resource "junos_interface_logical" "testacc_securityZone" {
  name                       = "${var.interface}.0"
  security_zone              = junos_security_zone.testacc_securityZone.name
  security_inbound_protocols = ["bgp"]
  security_inbound_services  = ["ssh"]
}
