
resource "junos_system_syslog_file" "testacc_syslogFile" {
  filename                     = "testacc_syslogFile"
  allow_duplicates             = true
  match                        = "match testacc"
  any_severity                 = "emergency"
  changelog_severity           = "critical"
  conflictlog_severity         = "error"
  daemon_severity              = "warning"
  dfc_severity                 = "alert"
  external_severity            = "any"
  firewall_severity            = "info"
  ftp_severity                 = "none"
  interactivecommands_severity = "notice"
  kernel_severity              = "error"
  ntp_severity                 = "error"
  pfe_severity                 = "error"
  security_severity            = "error"
  user_severity                = "error"
  structured_data {
    brief = true
  }
  archive {
    binary_data       = true
    world_readable    = true
    size              = 1073741823
    files             = 5
    transfer_interval = 5
    sites {
      url      = "example.com"
      password = "password"
    }
  }
}
