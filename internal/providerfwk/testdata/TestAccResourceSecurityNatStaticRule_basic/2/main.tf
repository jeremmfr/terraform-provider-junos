resource "junos_security_nat_static" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATSttRule.name]
  }
  configure_rules_singly = true
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule" {
  name                = "testacc_securityNATSttRule"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "192.0.2.0/27"
  then {
    type             = "prefix"
    routing_instance = junos_routing_instance.testacc_securityNATSttRule.name
    prefix           = "192.0.2.64/27"
  }
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule2" {
  depends_on = [
    junos_security_address_book.testacc_securityNATSttRule
  ]
  name                     = "testacc_securityNATSttRule2"
  rule_set                 = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address_name = "testacc_securityNATSttRule2"
  source_address = [
    "192.0.2.144/28"
  ]
  source_port = [
    "1024",
    "1025 to 1026",
  ]
  then {
    routing_instance = junos_routing_instance.testacc_securityNATSttRule.name
    type             = "prefix-name"
    prefix           = "testacc_securityNATSttRule-prefix"
  }
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule3" {
  depends_on = [
    junos_security_address_book.testacc_securityNATSttRule
  ]
  name                = "testacc_securityNATSttRule3"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "192.0.3.1/32"
  source_address_name = [
    "testacc_securityNATSttRule-src"
  ]
  destination_port    = 81
  destination_port_to = 82
  then {
    routing_instance = junos_routing_instance.testacc_securityNATSttRule.name
    type             = "prefix"
    prefix           = "192.0.3.2/32"
    mapped_port      = 8081
    mapped_port_to   = 8082
  }
}

resource "junos_security_zone" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
}
resource "junos_routing_instance" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
}
resource "junos_security_address_book" "testacc_securityNATSttRule" {
  network_address {
    name  = "testacc_securityNATSttRule2"
    value = "192.0.2.160/28"
  }
  network_address {
    name  = "testacc_securityNATSttRule-prefix"
    value = "192.0.2.176/28"
  }
  network_address {
    name  = "testacc_securityNATSttRule-src"
    value = "192.0.2.224/27"
  }
}
