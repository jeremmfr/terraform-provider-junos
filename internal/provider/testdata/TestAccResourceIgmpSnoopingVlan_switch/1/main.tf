resource "junos_igmp_snooping_vlan" "all" {
  name = "all"
}
resource "junos_igmp_snooping_vlan" "vlan10" {
  name            = "vlan10"
  immediate_leave = true
  interface {
    name                = "${var.interface}.1"
    host_only_interface = true
  }
  proxy = true
}
