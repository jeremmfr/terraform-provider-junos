resource "junos_security_address_book" "testacc_securityGlobalAddressBook" {
  description = "testacc global description"
  network_address {
    name        = "testacc_network"
    description = "testacc_network description"
    value       = "10.1.0.0/24"
  }
  dns_name {
    name        = "testacc_dns"
    description = "testacc_dns description"
    value       = "google.com"
    ipv4_only   = true
  }
  dns_name {
    name      = "testacc_dns6"
    value     = "google.com"
    ipv6_only = true
  }
  address_set {
    name        = "testacc_addressSet"
    description = "testacc_addressSet description"
    address     = ["testacc_network", "testacc_dns"]
  }
  address_set {
    name        = "testacc_addressSet2"
    address     = ["testacc_dns"]
    address_set = ["testacc_addressSet"]
  }
}

resource "junos_security_address_book" "testacc_securityNamedAddressBook" {
  name = "testacc_secAddrBook"
  network_address {
    name  = "testacc_network"
    value = "10.1.2.4/32"
  }
}
