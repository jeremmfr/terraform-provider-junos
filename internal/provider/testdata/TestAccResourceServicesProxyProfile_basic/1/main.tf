resource "junos_services_proxy_profile" "testacc_services_proxy_profile" {
  name               = "testacc_services_proxy_profile"
  protocol_http_host = "192.0.2.2"
  protocol_http_port = 3129
}
