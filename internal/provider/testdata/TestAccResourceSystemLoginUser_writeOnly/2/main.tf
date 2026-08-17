resource "junos_system_login_user" "testacc_wo" {
  name  = "testacc_wo"
  class = "unauthorized"
  authentication {
    encrypted_password_wo         = "test"
    encrypted_password_wo_version = 1
  }
}
# read the configuration to check that the write-only password has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_wo" {
  format = "set"
}
