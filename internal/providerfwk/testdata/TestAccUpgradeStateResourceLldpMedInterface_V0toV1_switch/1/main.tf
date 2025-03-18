resource "junos_lldpmed_interface" "testacc_all" {
  name = "all"
  location {
    civic_based_country_code = "FR"
    civic_based_what         = 1
  }
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name = var.interface
  location {
    civic_based_country_code = "FR"
    civic_based_ca_type {
      ca_type  = 10
      ca_value = "testacc"
    }
    civic_based_ca_type {
      ca_type = 0
    }
  }
}
