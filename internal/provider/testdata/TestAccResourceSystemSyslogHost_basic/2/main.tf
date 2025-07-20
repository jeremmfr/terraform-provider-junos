resource "junos_system_syslog_host" "testacc_syslogHost" {
  host = "192.0.2.1"
  structured_data {
    brief = true
  }
}
