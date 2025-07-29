resource "junos_system_syslog_user" "testacc_syslogUser" {
  username = "testacc"

  match                        = "testacc"
  match_strings                = ["testacc"]
  any_severity                 = "emergency"
  changelog_severity           = "emergency"
  conflictlog_severity         = "emergency"
  daemon_severity              = "emergency"
  dfc_severity                 = "emergency"
  external_severity            = "emergency"
  firewall_severity            = "emergency"
  ftp_severity                 = "emergency"
  interactivecommands_severity = "error"
  kernel_severity              = "emergency"
  ntp_severity                 = "emergency"
  pfe_severity                 = "emergency"
  security_severity            = "emergency"
  user_severity                = "emergency"
}

resource "junos_system_syslog_user" "testacc_syslogUserAll" {
  username = "*"

  any_severity = "emergency"
}
