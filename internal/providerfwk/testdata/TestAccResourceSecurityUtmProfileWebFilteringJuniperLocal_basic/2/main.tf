resource "junos_security_utm_profile_web_filtering_juniper_local" "testacc_ProfileWebFL" {
  name                 = "testacc ProfileWebFL"
  custom_block_message = "Blocked by Juniper"
  custom_message       = junos_security_utm_custom_message.testacc_ProfileWebFL.name
  default_action       = "log-and-permit"
  timeout              = 3
}

resource "junos_security_utm_custom_message" "testacc_ProfileWebFL" {
  name    = "testacc-profilewebfl"
  type    = "user-message"
  content = "testacc_ProfileWebFL"
}
