resource "junos_security_utm_profile_web_filtering_websense_redirect" "testacc_ProfileWebFWebS" {
  name                 = "testacc ProfileWebFWebS"
  custom_block_message = "Blocked by Juniper"
  timeout              = 3
  fallback_settings {
    default             = "log-and-permit"
    server_connectivity = "log-and-permit"
    timeout             = "log-and-permit"
  }
  server {
    host = "10.0.0.1"
    port = 1024
  }
  sockets = 16
}
