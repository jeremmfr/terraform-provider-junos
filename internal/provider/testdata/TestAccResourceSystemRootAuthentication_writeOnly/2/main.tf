resource "junos_system_root_authentication" "root_auth_wo" {
  encrypted_password_wo         = "$6$XXXX"
  encrypted_password_wo_version = 1
}
# read the configuration to check that the write-only password has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "root_auth_wo" {
  format = "set"
}
