resource "junos_security_nat_static" "testacc_securityNATSttRule" {
  name = "testacc_secNATSttRule_upgrade"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATSttRule.name]
  }
  configure_rules_singly = true
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule" {
  name                = "testacc_secNATSttRule_upgrade"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "192.0.2.0/25"
  then {
    type   = "prefix"
    prefix = "192.0.2.128/25"
  }
}
resource "junos_security_zone" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule_upgrade"
}
