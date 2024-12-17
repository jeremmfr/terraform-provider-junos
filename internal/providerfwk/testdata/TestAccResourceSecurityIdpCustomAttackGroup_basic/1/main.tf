resource "junos_security_idp_custom_attack" "testacc_idpCustomAttackGroup" {
  name               = "testacc/#1_Group"
  recommended_action = "ignore"
  severity           = "info"
  attack_type_anomaly {
    direction = "any"
    service   = "TELNET"
    test      = "SUBOPTION_OVERFLOW"
    shellcode = "all"
  }
}
resource "junos_security_idp_custom_attack_group" "testacc_idpCustomAttackGroup" {
  name = "testacc/#1_CustomAttackGroup"
  member = [
    junos_security_idp_custom_attack.testacc_idpCustomAttackGroup.name,
  ]
}
