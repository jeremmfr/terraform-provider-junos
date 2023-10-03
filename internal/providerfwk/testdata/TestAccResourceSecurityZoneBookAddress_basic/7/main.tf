resource "junos_security_zone" "testacc_szone_bookaddress" {
  name                          = "testacc_szone_bookaddress"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress1" {
  name        = "testacc_szone_bookaddress1"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  cidr        = "192.0.2.128/25"
  description = "testacc szone address1"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress2" {
  name        = "testacc_szone_bookaddress2"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  dns_name    = "test.fr"
  description = "testacc szone address2"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress3" {
  name        = "testacc_szone_bookaddress3"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  dns_name    = "test.net"
  description = "testacc szone address3"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress5" {
  name        = "testacc_szone_bookaddress5"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  range_from  = "192.0.2.20"
  range_to    = "192.0.2.22"
  description = "testacc szone ddress5"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress6" {
  name        = "testacc_szone_bookaddress6"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  wildcard    = "192.0.4.0/255.0.255.255"
  description = "testacc szone address6"
}
