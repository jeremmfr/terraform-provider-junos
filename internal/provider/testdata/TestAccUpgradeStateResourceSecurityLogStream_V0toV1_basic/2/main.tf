resource "junos_routing_instance" "testacc_logstream" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacclogstream"
}

resource "junos_security_log_stream" "testacc_logstream" {
  name     = "testacc_logstream"
  category = ["idp"]
  format   = "syslog"
  host {
    ip_address       = "192.0.2.1"
    port             = 514
    routing_instance = junos_routing_instance.testacc_logstream.name
  }
  rate_limit = 50
  severity   = "error"
}
resource "junos_security_log_stream" "testacc_logstream2" {
  name = "testacc_logstream2"
  file {
    name             = "test"
    allow_duplicates = true
    size             = 3
    rotation         = 3
  }
  filter_threat_attack = true
}
