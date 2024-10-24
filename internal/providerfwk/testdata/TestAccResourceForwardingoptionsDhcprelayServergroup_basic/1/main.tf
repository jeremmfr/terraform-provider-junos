resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6" {
  name    = "testacc_dhcprelay_servergroup_v6"
  version = "v6"
}

resource "junos_routing_instance" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_ri" {
  name             = "testacc_dhcprelay_servergroup_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelay_servergroup.name
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6_ri" {
  name             = "testacc_dhcprelay_servergroup_v6_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelay_servergroup.name
  version          = "v6"
}
