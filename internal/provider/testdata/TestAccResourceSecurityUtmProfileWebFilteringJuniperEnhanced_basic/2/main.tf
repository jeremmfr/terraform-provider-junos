resource "junos_security_utm_profile_web_filtering_juniper_enhanced" "testacc_ProfileWebFE" {
  name = "testacc ProfileWebFE"
  category {
    name           = "Enhanced_Network_Errors"
    action         = "block"
    custom_message = junos_security_utm_custom_message.testacc_ProfileWebFE.name
  }
  custom_block_message      = "Blocked by Juniper"
  custom_message            = junos_security_utm_custom_message.testacc_ProfileWebFE.name
  default_action            = "log-and-permit"
  no_safe_search            = true
  quarantine_custom_message = "Quarantine by Juniper"
  quarantine_message {
    url                      = "quarantine.local"
    type_custom_redirect_url = true
  }
  site_reputation_action {
    site_reputation = "harmful"
    action          = "block"
  }
  timeout = 3
}

resource "junos_security_utm_custom_message" "testacc_ProfileWebFE" {
  name    = "testacc-profilewebfe"
  type    = "user-message"
  content = "testacc_ProfileWebFE"
}
