resource "junos_lldpmed_interface" "testacc_all" {
  name    = "all"
  disable = true
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name   = var.interface
  enable = true
  location {}
}
