resource "junos_lldp_interface" "testacc_all" {
  name = "all"
}
resource "junos_lldp_interface" "testacc_interface" {
  name = var.interface
  power_negotiation {}
}
