resource "junos_vlan" "testacc_community" {
  name         = "testacc_community"
  vlan_id      = 20
  private_vlan = "community"
}

resource "junos_vlan" "testacc_community-one" {
  name         = "testacc_community-one"
  vlan_id      = 30
  private_vlan = "community"
}

resource "junos_vlan" "testacc_isolated" {
  name         = "testacc_isolated"
  vlan_id      = 200
  private_vlan = "isolated"
}

resource "junos_vlan" "testacc_pvlan" {
  name          = "testacc_pvlan"
  vlan_id       = "2000"
  isolated_vlan = junos_vlan.testacc_isolated.name
  community_vlans = [
    junos_vlan.testacc_community.name,
    junos_vlan.testacc_community-one.name,
  ]
}

resource "junos_vlan" "testacc_none" {
  name    = "testacc_none"
  vlan_id = "none"
}

resource "junos_routing_instance" "testacc_vlan_ri" {
  name = "testacc_vlan_ri"
  type = "virtual-switch"
}

resource "junos_vlan" "testacc_vlan_ri" {
  name             = "testacc_vlan_ri"
  routing_instance = junos_routing_instance.testacc_vlan_ri.name
  description      = "testacc_vlan_ri"
  vlan_id          = 200
}
