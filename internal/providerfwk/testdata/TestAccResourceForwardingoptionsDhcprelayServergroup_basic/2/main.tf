resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
  ip_address = [
    "192.0.2.8",
  ]
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6" {
  name    = "testacc_dhcprelay_servergroup_v6"
  version = "v6"
  ip_address = [
    "fe80::b",
  ]
}

resource "junos_routing_instance" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_ri" {
  name             = "testacc_dhcprelay_servergroup_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelay_servergroup.name
  ip_address = [
    "192.0.2.88",
  ]
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6_ri" {
  name             = "testacc_dhcprelay_servergroup_v6_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelay_servergroup.name
  version          = "v6"
  ip_address = [
    "fe80::bb",
  ]
}
