resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name               = "testacc/#1_"
  recommended_action = "ignore"
  severity           = "info"
  attack_type_anomaly {
    direction = "any"
    service   = "TELNET"
    test      = "SUBOPTION_OVERFLOW"
    shellcode = "all"
  }
}
