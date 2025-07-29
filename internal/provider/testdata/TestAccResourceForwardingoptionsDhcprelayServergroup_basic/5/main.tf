resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
  ip_address = [
    "fe80::b",
    "192.0.2.8",
  ]
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6" {
  name    = "testacc_dhcprelay_servergroup_v6"
  version = "v6"
  ip_address = [
    "fe80::b",
    "192.0.2.9",
    "fe80::a",
  ]
}
