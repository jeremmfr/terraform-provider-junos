resource "junos_security_nat_static" "testacc_securityNATStt" {
  name        = "testacc_securityNATStt"
  description = "testacc securityNATStt"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATStt.name]
  }
  rule {
    name                = "testacc_securityNATSttRule"
    destination_address = "192.0.2.0/25"
    then {
      type             = "prefix"
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
      prefix           = "192.0.2.128/25"
    }
  }
  rule {
    name                = "testacc_securityNATSttRule2"
    destination_address = "64:ff9b::/96"
    then {
      type = "inet"
    }
  }
}

resource "junos_security_zone" "testacc_securityNATStt" {
  name = "testacc_securityNATStt"
}
resource "junos_routing_instance" "testacc_securityNATStt" {
  name = "testacc_securityNATStt"
}
