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
