resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v4_default" {
  name = "testacc_dhcpgroup_v4_default"

  dynamic_profile = "junos-default-profile"
}

resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v6_default" {
  name    = "testacc_dhcpgroup_v6_default"
  version = "v6"

  interface {
    name    = "ge-0/0/3.1"
    upto    = "ge-0/0/3.3"
    exclude = true
  }
}

resource "junos_routing_instance" "testacc_dhcpgroup" {
  name = "testacc_dhcpgroup"
}
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v4_ri" {
  name             = "testacc_dhcpgroup_v4_ri"
  routing_instance = junos_routing_instance.testacc_dhcpgroup.name

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
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v6_ri" {
  name             = "testacc_dhcpgroup_v6_ri"
  routing_instance = junos_routing_instance.testacc_dhcpgroup.name
  version          = "v6"

  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}
