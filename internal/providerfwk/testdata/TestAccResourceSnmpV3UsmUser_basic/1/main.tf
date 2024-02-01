resource "junos_snmp_v3_usm_user" "testacc_snmpv3user" {
  name = "testacc_snmpv3user"
}
resource "junos_snmp_v3_usm_user" "testacc_snmpv3user_2" {
  name                = "testacc_snmpv3user#2"
  authentication_type = "authentication-md5"
  authentication_key  = "keymd5"
  privacy_type        = "privacy-3des"
  privacy_key         = "key3des"
}
resource "junos_snmp_v3_usm_user" "testacc_snmpv3user_3" {
  name                    = "testacc_snmpv3user#3"
  engine_type             = "remote"
  engine_id               = "engineID"
  authentication_type     = "authentication-sha"
  authentication_password = "pass1234"
  privacy_type            = "privacy-aes128"
  privacy_password        = "pass5678"
}
