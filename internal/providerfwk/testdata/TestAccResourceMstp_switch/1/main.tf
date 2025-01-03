resource "junos_routing_instance" "testacc_mstp" {
  name = "testacc_mstp"
  type = "virtual-switch"
}

resource "junos_mstp" "testacc_mstp" {
  routing_instance   = junos_routing_instance.testacc_mstp.name
  bridge_priority    = 0
  bpdu_block_on_edge = true
  system_id {
    id = "00:11:22:33:44:56"
  }
}
