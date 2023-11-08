resource "junos_application_set" "testacc_app_set" {
  name         = "testacc_app_set"
  applications = ["junos-ssh", "junos-telnet"]
  application_set = [
    junos_application_set.testacc_app_set2.name
  ]
}
resource "junos_application_set" "testacc_app_set2" {
  name         = "testacc_app_set2"
  applications = ["junos-ftp"]
  description  = "testacc appset2"
}
