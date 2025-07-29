resource "junos_rstp" "testacc_rstp" {
  bpdu_block_on_edge = true
}

resource "junos_routing_instance" "testacc_rstp" {
  name = "testacc_rstp"
  type = "virtual-switch"
}

resource "junos_rstp" "testacc_ri_rstp" {
  routing_instance = junos_routing_instance.testacc_rstp.name
  bridge_priority  = 0
  system_id {
    id = "00:11:22:33:44:56"
  }
}
