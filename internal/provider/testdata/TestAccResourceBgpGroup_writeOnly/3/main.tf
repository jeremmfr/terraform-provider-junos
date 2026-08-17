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
  authentication_key_wo         = "password2"
  authentication_key_wo_version = 2
}
