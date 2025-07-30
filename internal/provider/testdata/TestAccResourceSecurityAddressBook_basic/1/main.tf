resource "junos_security_zone" "testacc_secZoneAddr1" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_secZoneAddr1"
}
resource "junos_security_zone" "testacc_secZoneAddr2" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_secZoneAddr2"
}

resource "junos_security_address_book" "testacc_securityGlobalAddressBook" {
  description = "testacc global description"
  network_address {
    name        = "testacc_network"
    description = "testacc_network description"
    value       = "10.0.0.0/24"
  }
  wildcard_address {
    name        = "testacc_wildcard"
    description = "testacc_wildcard description"
    value       = "10.0.0.0/255.255.0.255"
  }
  network_address {
    name        = "testacc_network2"
    description = "testacc_network description2"
    value       = "10.1.0.0/24"
  }
  range_address {
    name        = "testacc_range"
    description = "testacc_range description"
    from        = "10.1.1.1"
    to          = "10.1.1.5"
  }
  dns_name {
    name  = "testacc_dns"
    value = "google.com"
  }
  address_set {
    name    = "testacc_addressSet"
    address = ["testacc_network", "testacc_wildcard", "testacc_network2"]
  }
}

resource "junos_security_address_book" "testacc_securityNamedAddressBook" {
  name        = "testacc_secAddrBook"
  attach_zone = [junos_security_zone.testacc_secZoneAddr1.name, junos_security_zone.testacc_secZoneAddr2.name]
  network_address {
    name  = "testacc_network"
    value = "10.1.2.3/32"
  }
}
