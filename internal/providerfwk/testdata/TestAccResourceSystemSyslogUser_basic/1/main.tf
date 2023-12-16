resource "junos_system_syslog_user" "testacc_syslogUser" {
  username = "testacc"

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
}

resource "local_file" "clear_user_all" {
  content  = "delete system syslog user *"
  filename = "/tmp/testacc/terraform-provider-junos/ResourceSystemSyslogUser_basic_clear_user_all.txt"
}

resource "junos_null_commit_file" "clear_user_all" {
  filename = local_file.clear_user_all.filename
}

resource "junos_system_syslog_user" "testacc_syslogUserAll" {
  depends_on = [junos_null_commit_file.clear_user_all]

  username = "*"

  any_severity                 = "emergency"
  changelog_severity           = "emergency"
  conflictlog_severity         = "emergency"
  daemon_severity              = "emergency"
  dfc_severity                 = "emergency"
  external_severity            = "emergency"
  firewall_severity            = "emergency"
  ftp_severity                 = "emergency"
  interactivecommands_severity = "emergency"
  kernel_severity              = "emergency"
  ntp_severity                 = "emergency"
  pfe_severity                 = "emergency"
  security_severity            = "emergency"
  user_severity                = "emergency"
}
