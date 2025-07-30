resource "junos_security_nat_static" "testacc_securityNATSttRuleInet" {
  name = "testacc_securityNATSttRuleInet"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATSttRuleInet.name]
  }
  configure_rules_singly = true
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRuleInet" {
  name                = "testacc_securityNATSttRuleInet"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRuleInet.name
  destination_address = "64:ff9b::/96"
  then {
    type             = "inet"
    routing_instance = junos_routing_instance.testacc_securityNATSttRuleInet.name
  }
}
resource "junos_security_zone" "testacc_securityNATSttRuleInet" {
  name = "testacc_securityNATSttRuleInet"
}
resource "junos_routing_instance" "testacc_securityNATSttRuleInet" {
  name = "testacc_securityNATSttRuleInet"
}
