resource "junos_security_utm_profile_web_filtering_websense_redirect" "testacc_ProfileWebFWebS" {
  name                 = "testacc ProfileWebFWebS"
  account              = "test ACC"
  custom_block_message = "Blocked by Juniper"
  fallback_settings {
    default             = "log-and-permit"
    server_connectivity = "log-and-permit"
    timeout             = "log-and-permit"
  }
}
