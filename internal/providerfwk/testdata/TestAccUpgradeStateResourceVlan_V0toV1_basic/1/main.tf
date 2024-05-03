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
  vlan_id       = 2000
  isolated_vlan = junos_vlan.testacc_isolated.vlan_id
  community_vlans = [
    junos_vlan.testacc_community.vlan_id,
    junos_vlan.testacc_community-one.vlan_id,
  ]
}

resource "junos_firewall_filter" "testacc_vlansw" {
  lifecycle {
    create_before_destroy = true
  }
  name   = "testacc_vlansw"
  family = "ethernet-switching"
  term {
    name = "testacc_vlansw_term1"
    then {
      action = "accept"
    }
  }
}
resource "junos_interface_logical" "testacc_vlansw" {
  lifecycle {
    create_before_destroy = true
  }
  name = "irb.1000"
}
resource "junos_vlan" "testacc_vlansw" {
  name                  = "testacc_vlansw"
  description           = "testacc_vlansw"
  vlan_id               = 1000
  service_id            = 1000
  l3_interface          = junos_interface_logical.testacc_vlansw.name
  forward_filter_input  = junos_firewall_filter.testacc_vlansw.name
  forward_filter_output = junos_firewall_filter.testacc_vlansw.name
  forward_flood_input   = junos_firewall_filter.testacc_vlansw.name
}
