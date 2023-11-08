resource "junos_security_zone" "testacc_szone_bookaddressset" {
  name                          = "testacc_szone_bookaddressset"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress_set" {
  name = "testacc_szone_bookaddress_set1"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  cidr = "192.0.2.0/25"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress_set2" {
  name = "testacc_szone_bookaddress_set2"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  cidr = "192.0.2.128/25"
}
resource "junos_security_zone_book_address_set" "testacc_szone_bookaddress_set" {
  name = "testacc_szone_bookaddress_set"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  address = [
    junos_security_zone_book_address.testacc_szone_bookaddress_set.name,
    junos_security_zone_book_address.testacc_szone_bookaddress_set2.name,
  ]
}
resource "junos_security_zone_book_address_set" "testacc_szone_bookaddress_setS2" {
  name = "testacc_szone_bookaddress_setS2"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  address = [
    junos_security_zone_book_address.testacc_szone_bookaddress_set2.name,
  ]
  address_set = [
    junos_security_zone_book_address_set.testacc_szone_bookaddress_set.name,
  ]
}
