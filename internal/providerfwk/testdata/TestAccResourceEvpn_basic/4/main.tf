resource "junos_interface_logical" "testacc_evpn" {
  depends_on = [
    junos_routing_options.testacc_evpn,
  ]
  name        = "lo0.0"
  description = "testacc_evpn"
  family_inet {
    address {
      cidr_ip = "192.0.2.18/32"
    }
  }
}
resource "junos_routing_options" "testacc_evpn" {
  clean_on_destroy = true
  router_id        = "192.0.2.18"
}
resource "junos_switch_options" "testacc_evpn" {
  clean_on_destroy      = true
  vtep_source_interface = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_default" {
  depends_on = [
    junos_switch_options.testacc_evpn,
  ]
  encapsulation = "vxlan"
  switch_or_ri_options {
    route_distinguisher = "201:1"
    vrf_target          = "target:201:2"
  }
}
resource "junos_routing_instance" "testacc_evpn_ri1" {
  name                  = "testacc_evpn_ri1"
  type                  = "virtual-switch"
  route_distinguisher   = "11:1"
  vrf_target            = "target:11:2"
  vtep_source_interface = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_ri1" {
  routing_instance = junos_routing_instance.testacc_evpn_ri1.name
  encapsulation    = "vxlan"
  duplicate_mac_detection {
    auto_recovery_time  = 100
    detection_threshold = 10
    detection_window    = 30
  }
}
resource "junos_routing_instance" "testacc_evpn_ri2" {
  name                        = "testacc_evpn_ri2"
  type                        = "virtual-switch"
  configure_rd_vrfopts_singly = true
  vtep_source_interface       = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_ri2" {
  routing_instance = junos_routing_instance.testacc_evpn_ri2.name
  encapsulation    = "vxlan"
  switch_or_ri_options {
    route_distinguisher = "101:1"
    vrf_target          = "target:101:2"
  }
}
