resource "junos_security_utm_profile_web_filtering_juniper_local" "testacc_ProfileWebFL" {
  name                 = "testacc ProfileWebFL"
  custom_block_message = "Blocked by Juniper"
  default_action       = "log-and-permit"
  fallback_settings {
    default             = "log-and-permit"
    server_connectivity = "log-and-permit"
    timeout             = "log-and-permit"
  }
}
