resource "junos_security_nat_destination" "testacc_securityDNAT" {
  name        = "testacc_securityDNAT"
  description = "testacc securityDNAT"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityDNAT.name]
  }
  rule {
    name                = "testacc_securityDNATRule"
    destination_address = "192.0.2.1/32"
    then {
      type = "pool"
      pool = junos_security_nat_destination_pool.testacc_securityDNATPool.name
    }
  }
}
resource "junos_security_nat_destination_pool" "testacc_securityDNATPool" {
  name             = "testacc_securityDNATPool"
  description      = "testacc securityDNATPool"
  address          = "192.0.2.1/32"
  address_to       = "192.0.2.2/32"
  routing_instance = junos_routing_instance.testacc_securityDNAT.name
}
resource "junos_security_nat_destination_pool" "testacc_securityDNATPool2" {
  name             = "testacc_securityDNATPool2"
  address          = "192.0.2.1/32"
  address_port     = 80
  routing_instance = junos_routing_instance.testacc_securityDNAT.name
}

resource "junos_security_zone" "testacc_securityDNAT" {
  name = "testacc_securityDNAT"
}
resource "junos_routing_instance" "testacc_securityDNAT" {
  name = "testacc_securityDNAT"
}
