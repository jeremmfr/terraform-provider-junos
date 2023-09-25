resource "junos_oam_gretunnel_interface" "testacc_oam_gretunnel_interface" {
  name           = "gr-3/3/0.3"
  hold_time      = 10
  keepalive_time = 5
}
resource "junos_oam_gretunnel_interface" "testacc_oam_gretunnel_interface2" {
  name      = "gr-3/3/0.2"
  hold_time = 11
}
resource "junos_oam_gretunnel_interface" "testacc_oam_gretunnel_interface3" {
  name           = "gr-3/3/0.1"
  keepalive_time = 2
}
