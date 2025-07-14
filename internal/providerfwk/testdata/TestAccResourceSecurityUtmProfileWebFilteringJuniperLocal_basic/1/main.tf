resource "junos_security_utm_profile_web_filtering_juniper_local" "testacc_ProfileWebFL" {
  name = "testacc ProfileWebFL"
  category {
    name   = junos_security_utm_custom_url_category.testacc_ProfileWebFL.name
    action = "permit"
  }
  custom_block_message = "Blocked by Juniper"
  default_action       = "log-and-permit"
  no_safe_search       = true
  fallback_settings {
    default             = "log-and-permit"
    server_connectivity = "log-and-permit"
    timeout             = "log-and-permit"
  }
}

resource "junos_security_utm_custom_url_pattern" "testacc_ProfileWebFL" {
  name  = "testacc-ProfileWebFL"
  value = ["*.google.com"]
}
resource "junos_security_utm_custom_url_category" "testacc_ProfileWebFL" {
  name = "testacc-ProfileWebFL"
  value = [
    junos_security_utm_custom_url_pattern.testacc_ProfileWebFL.name,
  ]
}
