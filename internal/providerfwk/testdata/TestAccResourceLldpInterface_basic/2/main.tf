resource "junos_lldp_interface" "testacc_all" {
  name                     = "all"
  enable                   = true
  trap_notification_enable = true

}
resource "junos_lldp_interface" "testacc_interface" {
  name                      = var.interface
  disable                   = true
  trap_notification_disable = true
}
