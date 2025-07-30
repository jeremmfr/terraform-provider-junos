resource "junos_application_set" "testacc_app_set" {
  name         = "testacc_app_set"
  applications = ["junos-ssh", "junos-telnet"]
}

data "junos_application_sets" "testacc_ssh_without_telnet" {
  match_applications = ["junos-ssh"]
}
data "junos_application_sets" "testacc_ssh_with_telnet" {
  match_applications = ["junos-telnet", "junos-ssh"]
}
data "junos_application_sets" "testacc_default_cifs" {
  match_applications = ["junos-netbios-session", "junos-smb-session"]
}
data "junos_application_sets" "testacc_name" {
  match_name = "testacc_.*"
}
data "junos_application_sets" "testacc_appsets" {
  match_application_sets = ["testacc_app_set"]
}
