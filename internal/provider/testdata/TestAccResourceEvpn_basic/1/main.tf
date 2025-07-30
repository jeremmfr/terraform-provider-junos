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
resource "junos_policyoptions_community" "testacc_evpn" {
  lifecycle {
    create_before_destroy = true
  }
  name    = "testacc_evpn"
  members = ["target:65000:100"]
}
resource "junos_policyoptions_policy_statement" "testacc_evpn" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_evpn1"
  from {
    bgp_community = [junos_policyoptions_community.testacc_evpn.name]
  }
  then {
    action = "accept"
  }
}
resource "junos_evpn" "testacc_evpn_default" {
  depends_on = [
    junos_switch_options.testacc_evpn,
  ]
  encapsulation     = "vxlan"
  no_core_isolation = true
  switch_or_ri_options {
    route_distinguisher = "20:1"
    vrf_target          = "target:20:2"
    vrf_import          = [junos_policyoptions_policy_statement.testacc_evpn.name]
    vrf_export          = [junos_policyoptions_policy_statement.testacc_evpn.name]
  }
}
resource "junos_routing_instance" "testacc_evpn_ri1" {
  name                  = "testacc_evpn_ri1"
  type                  = "virtual-switch"
  route_distinguisher   = "1:1"
  vrf_target            = "target:1:2"
  vrf_import            = [junos_policyoptions_policy_statement.testacc_evpn.name]
  vrf_export            = [junos_policyoptions_policy_statement.testacc_evpn.name]
  vtep_source_interface = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_ri1" {
  routing_instance = junos_routing_instance.testacc_evpn_ri1.name
  encapsulation    = "vxlan"
  duplicate_mac_detection {
    auto_recovery_time = 5
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
  default_gateway  = "advertise"
  switch_or_ri_options {
    route_distinguisher = "10:1"
    vrf_import          = [junos_policyoptions_policy_statement.testacc_evpn.name]
    vrf_export          = [junos_policyoptions_policy_statement.testacc_evpn.name]
    vrf_target          = "target:10:2"
  }
}
