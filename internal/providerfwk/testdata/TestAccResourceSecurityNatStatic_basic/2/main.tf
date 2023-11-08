resource "junos_security_nat_static" "testacc_securityNATStt" {
  depends_on = [
    junos_security_address_book.testacc_securityNATStt
  ]
  name = "testacc_securityNATStt"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATStt.name]
  }
  rule {
    name                = "testacc_securityNATSttRule"
    destination_address = "192.0.2.0/26"
    then {
      type             = "prefix"
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
      prefix           = "192.0.2.64/26"
    }
  }
  rule {
    name                     = "testacc_securityNATSttRule2"
    destination_address_name = "testacc_securityNATSttRule2"
    source_address = [
      "192.0.2.128/26"
    ]
    source_port = [
      "1024",
      "1025 to 1026",
    ]
    then {
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
      type             = "prefix-name"
      prefix           = "testacc_securityNATStt-prefix"
    }
  }
  rule {
    name                = "testacc_securityNATSttRule3"
    destination_address = "192.0.3.1/32"
    source_address_name = [
      "testacc_securityNATStt-src"
    ]
    destination_port    = 81
    destination_port_to = 82
    then {
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
      type             = "prefix"
      prefix           = "192.0.3.2/32"
      mapped_port      = 8081
      mapped_port_to   = 8082
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
resource "junos_security_nat_static" "testacc_securityNATStt_singly" {
  name = "testacc_securityNATStt_singly"
  from {
    type  = "routing-instance"
    value = [junos_routing_instance.testacc_securityNATStt.name]
  }
  configure_rules_singly = true
}
