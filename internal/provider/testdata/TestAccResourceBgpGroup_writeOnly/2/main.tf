resource "junos_routing_options" "testacc_bgpgroup_wo" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
}
resource "junos_bgp_group" "testacc_bgpgroup_wo" {
  depends_on = [
    junos_routing_options.testacc_bgpgroup_wo
  ]
  name                          = "testacc_bgpgroup_wo"
  authentication_key_wo         = "password"
  authentication_key_wo_version = 1
}
# read the configuration to check that the write-only key has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_bgpgroup_wo" {
  format = "set"
}
