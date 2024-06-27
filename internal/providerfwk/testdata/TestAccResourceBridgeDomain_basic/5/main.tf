resource "junos_bridge_domain" "testacc_default" {
  name        = "testacc_bd_def"
  description = "testacc bridge domain default update"
  vlan_id     = 8
}
resource "junos_bridge_domain" "testacc_default2" {
  name         = "testacc_bd_def2"
  vlan_id_list = [9]
}

resource "junos_interface_logical" "testacc_bridge_ri" {
  name = "lo0.1"
  family_inet {
    address {
      cidr_ip = "${junos_routing_options.testacc_bridge_ri.router_id}/32"
    }
  }
}
resource "junos_routing_options" "testacc_bridge_ri" {
  clean_on_destroy = true
  router_id        = "192.0.2.5"
}

resource "junos_routing_instance" "testacc_bridge_ri" {
  name                  = "testacc_bridge_ri"
  type                  = "virtual-switch"
  route_distinguisher   = "10:11"
  vrf_target            = "target:1:200"
  vtep_source_interface = junos_interface_logical.testacc_bridge_ri.name
  remote_vtep_list      = ["192.0.2.136", "192.0.2.36"]
}
resource "junos_evpn" "testacc_bridge_ri" {
  routing_instance = junos_routing_instance.testacc_bridge_ri.name
  encapsulation    = "vxlan"
  multicast_mode   = "ingress-replication"
}
resource "junos_bridge_domain" "testacc_bridge_ri" {
  depends_on = [
    junos_evpn.testacc_bridge_ri
  ]
  name              = "testacc_bd_ri"
  routing_instance  = junos_routing_instance.testacc_bridge_ri.name
  description       = "testacc bridge domain routing instance"
  routing_interface = "irb.13"
  interface = [
    "${var.interface}.0",
  ]
  service_id = 12
  vlan_id    = 13
  vxlan {
    vni                     = 15
    static_remote_vtep_list = ["192.0.2.36"]
  }
}
