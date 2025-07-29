resource "junos_security_utm_policy" "testacc_Policy" {
  name = "testacc Policy"
  anti_virus {
    http_profile = "junos-sophos-av-defaults"
  }
  traffic_sessions_per_client {
    over_limit = "log-and-permit"
  }
  web_filtering_profile = "junos-wf-local-default"
}
