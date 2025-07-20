resource "junos_security_utm_policy" "testacc_Policy" {
  name                   = "testacc Policy"
  anti_spam_smtp_profile = "junos-as-defaults"
  anti_virus {
    http_profile = "junos-sophos-av-defaults"
    imap_profile = "junos-av-defaults"
    pop3_profile = "junos-av-defaults"
    smtp_profile = "junos-av-defaults"
  }
  traffic_sessions_per_client {
    limit      = 1000
    over_limit = "log-and-permit"
  }
  web_filtering_profile = "junos-wf-local-default"
}
