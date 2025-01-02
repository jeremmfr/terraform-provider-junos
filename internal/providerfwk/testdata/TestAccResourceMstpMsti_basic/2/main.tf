resource "junos_interface_physical" "testacc_mstp_msti" {
  name         = var.interface
  vlan_members = ["default"]
}

resource "junos_interface_physical" "testacc_mstp_msti2" {
  name         = var.interface2
  vlan_members = ["default"]
}

resource "junos_mstp_msti" "testacc" {
  msti_id                = 17
  vlan                   = ["35-37", "32"]
  backup_bridge_priority = "28k"
  bridge_priority        = "24k"
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
