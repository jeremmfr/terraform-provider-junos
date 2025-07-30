resource "junos_vstp" "testacc_vstp" {
  bpdu_block_on_edge = true
}
resource "junos_routing_instance" "testacc_vstp" {
  name = "testacc_vstp"
  type = "virtual-switch"
}
resource "junos_vstp" "testacc_ri_vstp" {
  routing_instance = junos_routing_instance.testacc_vstp.name
  system_id {
    id = "00:11:22:33:44:56"
  }
}
