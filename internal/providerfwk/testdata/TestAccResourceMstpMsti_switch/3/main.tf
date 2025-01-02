resource "junos_routing_instance" "testacc_mstp_msti" {
  name = "testacc_mstp_msti"
  type = "virtual-switch"
}

resource "junos_mstp_msti" "testacc" {
  msti_id                = 3377
  routing_instance       = junos_routing_instance.testacc_mstp_msti.name
  vlan                   = ["333", "302", "310-313"]
  backup_bridge_priority = "8k"
  bridge_priority        = "4k"
  interface {
    name     = "all"
    cost     = 42
    priority = 64
  }
}
