resource "junos_application_set" "testacc_app_set" {
  name         = "testacc_app_set"
  applications = ["junos-ssh", "junos-telnet"]
}
resource "junos_application_set" "testacc_app_set2" {
  name            = "testacc_app_set2"
  application_set = [junos_application_set.testacc_app_set.name]
  description     = "test-data-source-appSet"
}
