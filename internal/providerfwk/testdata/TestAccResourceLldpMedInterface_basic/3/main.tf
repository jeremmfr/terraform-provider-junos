resource "junos_lldpmed_interface" "testacc_all" {
  name = "all"
  location {
    elin = "testacc_lldpmed"
  }
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name = var.interface
  location {
    co_ordinate_latitude  = 180
    co_ordinate_longitude = 180
  }
}
