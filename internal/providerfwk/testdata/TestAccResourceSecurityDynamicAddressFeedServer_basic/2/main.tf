resource "junos_security_dynamic_address_feed_server" "testacc_dyn_add_feed_srv" {
  name        = "tfacc_dafeedsrv"
  hostname    = "example.com/?test=#1"
  description = "testacc junos_security_dynamic_address_feed_server"
  feed_name {
    name            = "feed_b"
    path            = "/srx/"
    description     = "testacc junos_security_dynamic_address_feed_server feed_b"
    hold_interval   = 1110
    update_interval = 120
  }
  feed_name {
    name          = "feed_a"
    path          = "/srx/"
    hold_interval = 0
  }
  feed_name {
    name            = "feed_0"
    path            = "/srx/"
    description     = "testacc junos_security_dynamic_address_feed_server feed_0"
    hold_interval   = 1130
    update_interval = 140
  }
  hold_interval   = 1150
  update_interval = 160
}
