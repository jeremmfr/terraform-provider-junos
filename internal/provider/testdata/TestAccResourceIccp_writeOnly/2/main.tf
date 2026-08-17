resource "junos_iccp" "testacc_iccp_wo" {
  local_ip_addr                 = "192.0.2.1"
  authentication_key_wo         = "a@Key"
  authentication_key_wo_version = 1
}
# read the configuration to check that the write-only key has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_iccp_wo" {
  format = "set"
}
