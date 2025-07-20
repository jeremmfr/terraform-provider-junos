resource "junos_security_utm_policy" "testacc_Policy" {
  name = "testacc Policy"
  anti_virus {
    ftp_download_profile = "junos-av-defaults"
    ftp_upload_profile   = "junos-av-defaults"
  }
  content_filtering {
    ftp_download_profile = "junos-cf-defaults"
    ftp_upload_profile   = "junos-cf-defaults"
    http_profile         = "junos-cf-defaults"
    imap_profile         = "junos-cf-defaults"
    pop3_profile         = "junos-cf-defaults"
    smtp_profile         = "junos-cf-defaults"
  }
  web_filtering_profile = "junos-wf-enhanced-default"
}
