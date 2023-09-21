resource "junos_security_nat_destination" "testacc_securityDNAT" {
  depends_on = [
    junos_security_address_book.testacc_securityDNAT
  ]
  name = "testacc_securityDNAT"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityDNAT.name]
  }
  rule {
    name                = "testacc_securityDNATRule"
    destination_address = "192.0.2.1/32"
    source_address_name = [
      "testacc_securityDNAT-src",
    ]
    protocol = ["tcp", "50"]
    then {
      type = "pool"
      pool = junos_security_nat_destination_pool.testacc_securityDNATPool.name
    }
  }
  rule {
    name                     = "testacc_securityDNATRule2"
    destination_address_name = "testacc_securityDNAT"
    destination_port = [
      "81",
      "82 to 83",
    ]
    source_address = [
      "192.0.2.128/26",
    ]
    then {
      type = "pool"
      pool = junos_security_nat_destination_pool.testacc_securityDNATPool2.name
    }
  }
  rule {
    name                = "testacc_securityDNATRule3"
    destination_address = "192.0.2.1/32"
    application = [
      "junos-ssh", "junos-http",
    ]

    then {
      type = "pool"
      pool = junos_security_nat_destination_pool.testacc_securityDNATPool2.name
    }
  }
}
resource "junos_security_nat_destination_pool" "testacc_securityDNATPool" {
  name             = "testacc_securityDNATPool"
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
resource "junos_security_address_book" "testacc_securityDNAT" {
  network_address {
    name  = "testacc_securityDNAT"
    value = "192.0.2.128/27"
  }
  network_address {
    name  = "testacc_securityDNAT-src"
    value = "192.0.2.160/27"
  }
}
