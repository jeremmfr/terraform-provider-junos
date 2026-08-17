resource "junos_security_authentication_key_chain" "testacc_secauthKeyChain_wo" {
  name = "testacc_secauthKeyChainWO"
  key {
    id         = 4
    secret     = "aS3cret#4"
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
    "5" = {
      value   = "aS3cret#5"
      version = 1
    }
    "6" = {
      value   = "aS3cret#6"
      version = 1
    }
  }
}

# read the configuration to check that the write-only secrets have been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_secauthKeyChain_wo" {
  format = "set"
}
