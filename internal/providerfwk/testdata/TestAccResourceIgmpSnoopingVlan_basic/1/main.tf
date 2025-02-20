resource "junos_routing_instance" "testacc_igmp_snooping_vlan" {
  name = "testacc_igmp_snooping_vlan"
  type = "vpls"
}
resource "junos_igmp_snooping_vlan" "all" {
  name = "all"
}
resource "junos_igmp_snooping_vlan" "vlan10" {
  name             = "vlan10"
  routing_instance = junos_routing_instance.testacc_igmp_snooping_vlan.name
  immediate_leave  = true
  interface {
    name                = "${var.interface}.1"
    host_only_interface = true
  }
  proxy = true
}
