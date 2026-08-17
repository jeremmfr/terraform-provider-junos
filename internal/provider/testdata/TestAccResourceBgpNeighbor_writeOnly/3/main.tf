resource "junos_routing_options" "testacc_bgpneighbor_wo" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
}
resource "junos_bgp_group" "testacc_bgpneighbor_wo" {
  name = "testacc_bgpneighbor_wo"
  type = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor_wo" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor_wo
  ]
  ip                            = "192.0.2.4"
  group                         = junos_bgp_group.testacc_bgpneighbor_wo.name
  authentication_key_wo         = "password2"
  authentication_key_wo_version = 2
}
