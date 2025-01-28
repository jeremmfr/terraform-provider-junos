resource "junos_security_dynamic_address_feed_server" "testacc_dyn_add_feed_srv" {
  name     = "tfacc_dafeedsrv"
  hostname = "example.com"
  feed_name {
    name = "feed_b"
    path = "/srx/"
  }
  feed_name {
    name = "feed_a"
    path = "/srx/"
  }
}
resource "junos_security_dynamic_address_feed_server" "testacc_dyn_add_feed_srv2" {
  name = "tfacc_dafeedsr2"
  url  = "https://example.com/?test=#2"
}
