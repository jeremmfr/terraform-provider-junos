resource "junos_security_authentication_key_chain" "testacc_secauthKeyChain" {
  name        = "testacc_secauthKeyChain#1"
  description = "testacc secauthKeyChain"
  tolerance   = 600
  key {
    id         = 11
    secret     = "aS3cret#1"
    start_time = "2021-12-11.10:09:08"
  }
  key {
    id         = 5
    secret     = "aSecret#1234"
    start_time = "2024-12-11.10:09:08"

    algorithm = "md5"
    key_name  = "ffaa1234"
    options   = "basic"
  }
}

resource "junos_security_authentication_key_chain" "testacc_secauthKeyChainAO" {
  name = "testacc_secauthKeyChainAO"
  key {
    id         = 5
    secret     = "secret aa"
    start_time = "2025-12-11.10:09:08"

    algorithm                  = "ao"
    ao_cryptographic_algorithm = "aes-128-cmac-96"
    ao_recv_id                 = 100
    ao_send_id                 = 150
    ao_tcp_ao_option           = "enabled"
  }
}
