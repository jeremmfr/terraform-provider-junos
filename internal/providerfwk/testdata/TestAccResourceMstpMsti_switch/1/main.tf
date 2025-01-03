resource "junos_routing_instance" "testacc_mstp_msti" {
  name = "testacc_mstp_msti"
  type = "virtual-switch"
}

resource "junos_mstp_msti" "testacc" {
  msti_id          = 3377
  routing_instance = junos_routing_instance.testacc_mstp_msti.name
  vlan             = ["302", "333"]
}
