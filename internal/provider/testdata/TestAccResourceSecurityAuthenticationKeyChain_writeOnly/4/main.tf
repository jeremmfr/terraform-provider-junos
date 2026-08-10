# the key 4 is migrated from the standard argument to the write-only one
resource "junos_security_authentication_key_chain" "testacc_secauthKeyChain_wo" {
  name = "testacc_secauthKeyChainWO"
  key {
    id         = 4
    start_time = "2021-12-11.10:09:08"
  }
  key {
    id         = 5
    start_time = "2022-11-21.10:09:08"
  }
  key {
    id         = 6
    start_time = "2022-12-22.10:09:08"
  }
  key_secret_wo = {
    "4" = {
      value   = "aS3cret#4bis"
      version = 1
    }
    "5" = {
      value   = "aS3cret#5bis"
      version = 2
    }
    "6" = {
      value   = "aS3cret#6"
      version = 1
    }
  }
}
