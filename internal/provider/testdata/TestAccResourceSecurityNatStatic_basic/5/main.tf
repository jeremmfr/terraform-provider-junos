resource "junos_security_nat_static" "testacc_securityNATStt" {
  name = "testacc_securityNATStt"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATStt.name]
  }
  rule {
    name                = "testacc_securityNATSttRule"
    destination_address = "64:ff9b::/96"
    then {
      type             = "inet"
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
    }
  }
}

resource "junos_security_zone" "testacc_securityNATStt" {
  name = "testacc_securityNATStt"
}
resource "junos_routing_instance" "testacc_securityNATStt" {
  name = "testacc_securityNATStt"
}

resource "junos_security_address_book" "testacc_securityNATStt" {
  network_address {
    name  = "testacc_securityNATSttRule2"
    value = "192.0.2.128/27"
  }
  network_address {
    name  = "testacc_securityNATStt-prefix"
    value = "192.0.2.160/27"
  }
  network_address {
    name  = "testacc_securityNATStt-src"
    value = "192.0.2.224/27"
  }
}
