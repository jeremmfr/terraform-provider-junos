resource "junos_security_utm_profile_web_filtering_websense_redirect" "testacc_ProfileWebFWebS" {
  name    = "testacc ProfileWebFWebS"
  account = "test ACC"
  category {
    name   = junos_security_utm_custom_url_category.testacc_ProfileWebFWebS.name
    action = "permit"
  }
  custom_block_message = "Blocked by Juniper"
  no_safe_search       = true
  fallback_settings {
    default             = "log-and-permit"
    server_connectivity = "log-and-permit"
    timeout             = "log-and-permit"
  }
}

resource "junos_security_utm_custom_url_pattern" "testacc_ProfileWebFWebS" {
  name  = "testacc_ProfileWebFWebS"
  value = ["*.google.com"]
}
resource "junos_security_utm_custom_url_category" "testacc_ProfileWebFWebS" {
  name = "testacc_ProfileWebFWebS"
  value = [
    junos_security_utm_custom_url_pattern.testacc_ProfileWebFWebS.name,
  ]
}
