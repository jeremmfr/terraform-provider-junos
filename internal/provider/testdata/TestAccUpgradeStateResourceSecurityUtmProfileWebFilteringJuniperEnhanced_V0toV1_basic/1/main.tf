resource "junos_security_utm_profile_web_filtering_juniper_enhanced" "testacc_ProfileWebFE" {
  name = "testacc ProfileWebFE"
  block_message {
    url                      = "block.local"
    type_custom_redirect_url = true
  }
  category {
    name   = "Enhanced_Network_Errors"
    action = "block"
  }
  category {
    name   = "Enhanced_Suspicious_Content"
    action = "quarantine"
    reputation_action {
      site_reputation = "very-safe"
      action          = "log-and-permit"
    }
    reputation_action {
      site_reputation = "moderately-safe"
      action          = "log-and-permit"
    }
  }
  custom_block_message      = "Blocked by Juniper"
  default_action            = "log-and-permit"
  no_safe_search            = true
  quarantine_custom_message = "Quarantine by Juniper"
  fallback_settings {
    default             = "log-and-permit"
    server_connectivity = "log-and-permit"
    timeout             = "log-and-permit"
  }
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
