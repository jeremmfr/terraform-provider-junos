resource "junos_system_syslog_file" "testacc_syslogFile" {
  filename                     = "testacc"
  allow_duplicates             = true
  explicit_priority            = true
  match                        = "match testacc"
  match_strings                = ["match testacc"]
  any_severity                 = "emergency"
  changelog_severity           = "critical"
  conflictlog_severity         = "error"
  daemon_severity              = "warning"
  dfc_severity                 = "alert"
  external_severity            = "any"
  firewall_severity            = "info"
  ftp_severity                 = "none"
  interactivecommands_severity = "notice"
  kernel_severity              = "emergency"
  ntp_severity                 = "emergency"
  pfe_severity                 = "emergency"
  security_severity            = "emergency"
  user_severity                = "emergency"
}

resource "junos_system_syslog_file" "testacc_syslogFile2" {
  filename     = "test_acc.2"
  any_severity = "emergency"
}
