resource "junos_vstp" "testacc_vstp" {
  disable = true
}
resource "junos_routing_instance" "testacc_vstp" {
  name = "testacc_vstp"
  type = "virtual-switch"
}
resource "junos_vstp" "testacc_ri_vstp" {
  routing_instance   = junos_routing_instance.testacc_vstp.name
  bpdu_block_on_edge = true
  force_version_stp  = true
  priority_hold_time = 10
  system_id {
    id = "00:11:22:33:44:55"
  }
  system_id {
    id         = "00:22:33:44:55:aa"
    ip_address = "192.0.2.4/31"
  }
  vpls_flush_on_topology_change = true
}
