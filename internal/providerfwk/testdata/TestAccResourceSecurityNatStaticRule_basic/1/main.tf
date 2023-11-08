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
  destination_address = "192.0.2.0/28"
  then {
    type             = "prefix"
    routing_instance = junos_routing_instance.testacc_securityNATSttRule.name
    prefix           = "192.0.2.128/28"
  }
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRuleInet" {
  name                = "testacc_securityNATSttRuleInet"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "64:ff9b::/96"
  then {
    type = "inet"
  }
}

resource "junos_security_zone" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
}
resource "junos_routing_instance" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
}
