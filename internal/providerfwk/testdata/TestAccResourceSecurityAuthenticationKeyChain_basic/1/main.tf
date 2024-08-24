resource "junos_security_authentication_key_chain" "testacc_secauthKeyChain" {
  name = "testacc_secauthKeyChain#1"
  key {
    id         = 11
    secret     = "aS3cret#1"
    start_time = "2021-12-11.10:09:08"
  }
}
