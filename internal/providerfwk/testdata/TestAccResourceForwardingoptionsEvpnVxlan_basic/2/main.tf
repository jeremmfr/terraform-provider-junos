resource "junos_forwardingoptions_evpn_vxlan" "ri" {
  routing_instance = junos_routing_instance.fwOpts_evpn_vxlan.name
  shared_tunnels   = true
}

resource "junos_routing_instance" "fwOpts_evpn_vxlan" {
  name = "fwOpts_evpn_vxlan"
}
