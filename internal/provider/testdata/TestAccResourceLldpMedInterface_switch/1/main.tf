resource "junos_lldpmed_interface" "testacc_all" {
  name = "all"
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name = var.interface
  location {}
}
