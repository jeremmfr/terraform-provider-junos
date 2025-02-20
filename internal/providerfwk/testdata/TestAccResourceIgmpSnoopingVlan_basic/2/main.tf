resource "junos_routing_instance" "testacc_igmp_snooping_vlan" {
  name = "testacc_igmp_snooping_vlan"
  type = "vpls"
}
resource "junos_igmp_snooping_vlan" "vlan10" {
  name             = "vlan10"
  routing_instance = junos_routing_instance.testacc_igmp_snooping_vlan.name
  immediate_leave  = true
  interface {
    name = "${var.interface}.1"
  }
  interface {
    name                       = "${var.interface}.0"
    group_limit                = 32
    immediate_leave            = true
    multicast_router_interface = true
    static_group {
      address = "224.255.0.2"
    }
    static_group {
      address = "224.255.0.1"
      source  = "192.0.2.1"
    }
  }
  proxy                      = true
  proxy_source_address       = "192.0.2.11"
  query_interval             = 33
  query_last_member_interval = "32.1"
  query_response_interval    = "31.0"
  robust_count               = 5
}
