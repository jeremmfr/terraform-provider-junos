resource "junos_rstp_interface" "all" {
  name = "all"
}

resource "junos_routing_instance" "testacc_rstp_interface" {
  name = "testacc_rstp_intface"
  type = "virtual-switch"
}

resource "junos_rstp_interface" "all2" {
  name             = "all"
  routing_instance = junos_routing_instance.testacc_rstp_interface.name
}
