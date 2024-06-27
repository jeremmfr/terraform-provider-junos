import {
  to = junos_security.security
  id = "security"
}

resource "junos_security" "security" {
  log {
    source_address = "192.0.2.2"
  }
}
