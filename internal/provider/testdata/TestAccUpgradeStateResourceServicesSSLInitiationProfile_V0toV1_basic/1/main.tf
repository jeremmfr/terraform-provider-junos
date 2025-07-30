resource "junos_services_ssl_initiation_profile" "testacc_sslInitProf" {
  name = "testacc_sslInitProf.1"
  actions {
    crl_disable                      = true
    crl_if_not_present               = "allow"
    crl_ignore_hold_instruction_code = true
    ignore_server_auth_failure       = true
  }
  custom_ciphers       = ["tls12-rsa-aes-256-cbc-sha256", "tls12-rsa-aes-128-gcm-sha256"]
  enable_flow_tracing  = true
  enable_session_cache = true
  preferred_ciphers    = "medium"
  protocol_version     = "tls12"
}
