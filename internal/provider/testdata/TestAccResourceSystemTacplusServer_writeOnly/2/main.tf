resource "junos_system_tacplus_server" "testacc_tacplusServer_wo" {
  address           = "192.0.2.11"
  secret_wo         = "a@Secret"
  secret_wo_version = 1
}
# read the configuration to check that the write-only secret has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_tacplusServer_wo" {
  format = "set"
}
