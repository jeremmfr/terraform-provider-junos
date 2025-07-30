resource "junos_system_syslog_host" "testacc_syslogHost" {
  host                         = "192.0.2.1"
  port                         = 514
  allow_duplicates             = true
  exclude_hostname             = true
  explicit_priority            = true
  facility_override            = "local3"
  log_prefix                   = "prefix"
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
  structured_data {}
  source_address = "192.0.2.2"
}
