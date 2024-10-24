resource "junos_forwardingoptions_dhcprelay_group" "testacc_dhcprelaygroup_v4_default" {
  name = "testacc_dhcprelaygroup_v4_default"

  dynamic_profile = "junos-default-profile"
}

resource "junos_forwardingoptions_dhcprelay_group" "testacc_dhcprelaygroup_v6_default" {
  name    = "testacc_dhcprelaygroup_v6_default"
  version = "v6"

  interface {
    name    = "ge-0/0/3.1"
    upto    = "ge-0/0/3.3"
    exclude = true
  }
}

resource "junos_routing_instance" "testacc_dhcprelaygroup" {
  name = "testacc_dhcprelaygroup"
}
resource "junos_forwardingoptions_dhcprelay_group" "testacc_dhcprelaygroup_v4_ri" {
  name             = "testacc_dhcprelaygroup_v4_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelaygroup.name

  interface {
    name            = "ge-0/0/3.0"
    dynamic_profile = "junos-default-profile"
    trace           = true
  }
  interface {
    name            = "ge-0/0/3.1"
    dynamic_profile = "junos-default-profile"
  }
}
resource "junos_forwardingoptions_dhcprelay_group" "testacc_dhcprelaygroup_v6_ri" {
  name             = "testacc_dhcprelaygroup_v6_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelaygroup.name
  version          = "v6"

  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}
