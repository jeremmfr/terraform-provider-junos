resource "junos_system_syslog_file" "testacc_syslogFile_wo" {
  filename     = "testacc_wo"
  any_severity = "emergency"
  archive {
    files = 5
    sites {
      url = "192.0.2.1"
    }
    sites {
      url                 = "192.0.2.2"
      password_wo         = "aPassword2"
      password_wo_version = 2
    }
  }
}
