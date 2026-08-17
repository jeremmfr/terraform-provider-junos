resource "junos_services" "testacc_wo" {
  clean_on_destroy = true
  application_identification {
    no_application_system_cache = true
  }
}
