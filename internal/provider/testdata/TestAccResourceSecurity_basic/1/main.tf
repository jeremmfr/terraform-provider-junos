import {
  to = junos_system.system
  id = "system"
}

resource "junos_system" "system" {
  tracing_dest_override_syslog_host = "192.0.2.13"
}
