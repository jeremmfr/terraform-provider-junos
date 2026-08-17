resource "junos_snmp_v3_usm_user" "testacc_snmpv3user_wo" {
  name                          = "testacc_snmpv3user_wo"
  authentication_type           = "authentication-md5"
  authentication_key_wo         = "keymd5"
  authentication_key_wo_version = 1
  privacy_type                  = "privacy-3des"
  privacy_key_wo                = "key3des"
  privacy_key_wo_version        = 1
}
resource "junos_snmp_v3_usm_user" "testacc_snmpv3user_wo_2" {
  name                               = "testacc_snmpv3user_wo#2"
  authentication_type                = "authentication-sha"
  authentication_password_wo         = "pass1234"
  authentication_password_wo_version = 1
  privacy_type                       = "privacy-aes128"
  privacy_password_wo                = "pass5678"
  privacy_password_wo_version        = 1
}
