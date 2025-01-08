resource "junos_routing_instance" "testacc_mstp_msti" {
  name = "testacc_mstp_msti"
  type = "virtual-switch"
}

resource "junos_interface_physical" "testacc_mstp_msti" {
  lifecycle {
    create_before_destroy = true
  }

  name         = var.interface
  vlan_members = ["default"]
}

resource "junos_interface_physical" "testacc_mstp_msti2" {
  lifecycle {
    create_before_destroy = true
  }

  name         = var.interface2
  vlan_members = ["default"]
}

resource "junos_mstp_msti" "testacc" {
  msti_id                = 3377
  routing_instance       = junos_routing_instance.testacc_mstp_msti.name
  vlan                   = ["302", "333", "310-313"]
  backup_bridge_priority = "8k"
  bridge_priority        = "4k"
  interface {
    name     = junos_interface_physical.testacc_mstp_msti2.name
    cost     = 42
    priority = 64
  }
  interface {
    name     = junos_interface_physical.testacc_mstp_msti.name
    cost     = 52
    priority = 80
  }
}
