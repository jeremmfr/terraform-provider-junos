resource "junos_snmp_v3_usm_user" "testacc_snmpv3user_3" {
  name                    = "testacc_snmpv3user#3"
  engine_type             = "remote"
  engine_id               = "engine#ID"
  authentication_type     = "authentication-sha"
  authentication_password = "pass1234"
  privacy_type            = "privacy-des"
  privacy_password        = "aPasswordAA"
}

resource "junos_snmp_v3_usm_user" "testacc_snmpv3user_3_copy" {
  depends_on              = [junos_snmp_v3_usm_user.testacc_snmpv3user_3]
  name                    = "testacc_snmpv3user#3"
  engine_type             = "remote"
  engine_id               = "engine#ID"
  authentication_type     = "authentication-sha"
  authentication_password = "pass4321"
  privacy_type            = "privacy-des"
  privacy_password        = "aPasswordBB"
}
