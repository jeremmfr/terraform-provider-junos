resource "junos_lldp_interface" "testacc_all" {
  name   = "all"
  enable = true
  power_negotiation {
    enable = true
  }

}
resource "junos_lldp_interface" "testacc_interface" {
  name    = var.interface
  disable = true
  power_negotiation {
    disable = true
  }
}
