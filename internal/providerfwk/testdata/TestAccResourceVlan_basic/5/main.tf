
resource "junos_routing_options" "testacc_vlan_vxlan" {
  clean_on_destroy = true
  router_id        = "192.0.2.18"
}
resource "junos_interface_logical" "testacc_vlan_vxlan" {
  depends_on = [
    junos_routing_options.testacc_vlan_vxlan,
  ]
  name        = "lo0.0"
  description = "testacc_vlan_vxlan"
  family_inet {
    address {
      cidr_ip = "192.0.2.18/32"
    }
  }
}
resource "junos_policyoptions_community" "testacc_vlan_vxlan" {
  lifecycle {
    create_before_destroy = true
  }
  name    = "testacc_vlan_vxlan"
  members = ["target:65000:100"]
}
resource "junos_policyoptions_policy_statement" "testacc_vlan_vxlan" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_vlan_vxlan"
  from {
    bgp_community = [junos_policyoptions_community.testacc_vlan_vxlan.name]
  }
  then {
    action = "accept"
  }
}
resource "junos_switch_options" "testacc_vlan_vxlan" {
  clean_on_destroy      = true
  vtep_source_interface = junos_interface_logical.testacc_vlan_vxlan.name
}
resource "junos_evpn" "testacc_vlan_vxlan" {
  depends_on = [
    junos_switch_options.testacc_vlan_vxlan,
  ]
  encapsulation = "vxlan"
  switch_or_ri_options {
    route_distinguisher = "20:1"
    vrf_target          = "target:20:2"
    vrf_import          = [junos_policyoptions_policy_statement.testacc_vlan_vxlan.name]
    vrf_export          = [junos_policyoptions_policy_statement.testacc_vlan_vxlan.name]
  }
}

resource "junos_vlan" "testacc_vlan_vxlan" {
  depends_on = [
    junos_evpn.testacc_vlan_vxlan,
  ]
  name    = "testacc_vlan_vxlan"
  vlan_id = 1020
  vxlan {
    vni                          = 102010
    vni_extend_evpn              = true
    encapsulate_inner_vlan       = true
    ingress_node_replication     = true
    unreachable_vtep_aging_timer = 900
  }
}

resource "junos_routing_instance" "testacc_vlan_ri" {
  name                  = "testacc_vlan_ri"
  type                  = "virtual-switch"
  route_distinguisher   = "11:1"
  vrf_target            = "target:11:2"
  vtep_source_interface = junos_interface_logical.testacc_vlan_vxlan.name
}

resource "junos_evpn" "testacc_vlan_ri" {
  routing_instance = junos_routing_instance.testacc_vlan_ri.name
  encapsulation    = "vxlan"
}

resource "junos_vlan" "testacc_vlan_ri" {
  depends_on = [
    junos_evpn.testacc_vlan_ri,
  ]

  name               = "testacc_vlan_ri"
  routing_instance   = junos_routing_instance.testacc_vlan_ri.name
  description        = "testacc_vlan_ri"
  vlan_id            = 1030
  no_arp_suppression = true
  vxlan {
    vni             = 103010
    vni_extend_evpn = true
    translation_vni = 1103010
  }
}
