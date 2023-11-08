resource "junos_security" "testacc_security" {
  flow {
    tcp_session {
      time_wait_state {
        session_timeout = 90
      }
    }
  }
  idp_sensor_configuration {
    log_suppression {}
  }
}
