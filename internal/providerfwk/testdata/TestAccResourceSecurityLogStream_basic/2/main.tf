resource "junos_routing_instance" "testacc_logstream" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacclogstream"
}
resource "junos_services_ssl_initiation_profile" "testacc_logstream" {
  lifecycle {
    create_before_destroy = true
  }
  name = "test@cc logstream"
}

resource "junos_security_log_stream" "testacc_logstream" {
  name     = "testacc_logstream"
  category = ["idp"]
  format   = "syslog"
  host {
    ip_address       = "AHostN@me"
    port             = 514
    routing_instance = junos_routing_instance.testacc_logstream.name
  }
  rate_limit = 50
  severity   = "error"
  transport {
    protocol        = "tls"
    tcp_connections = 3
    tls_profile     = junos_services_ssl_initiation_profile.testacc_logstream.name
  }
}


resource "junos_security_log_stream" "testacc_logstream2" {
  name   = "testacc_logstream2"
  format = "syslog"
  host {
    ip_address = "192.0.2.1"
  }
  transport {
    protocol = "udp"
  }
}
