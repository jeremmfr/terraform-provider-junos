resource "junos_routing_instance" "testacc_mstp_interface" {
  name = "testacc_mstp_interface"
  type = "virtual-switch"
}

resource "junos_mstp_interface" "all" {
  name             = "all"
  routing_instance = junos_routing_instance.testacc_mstp_interface.name
}
