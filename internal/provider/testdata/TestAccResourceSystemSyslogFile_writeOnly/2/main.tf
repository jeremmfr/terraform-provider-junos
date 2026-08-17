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
      password_wo         = "aPassword"
      password_wo_version = 1
    }
  }
}

# read the configuration to check that the write-only password has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_syslogFile_wo" {
  format = "set"
}
