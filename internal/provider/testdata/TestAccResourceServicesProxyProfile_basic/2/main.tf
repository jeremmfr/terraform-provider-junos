resource "junos_services_proxy_profile" "testacc_services_proxy_profile" {
  name               = "testacc_services_proxy_profile"
  protocol_http_host = "192.0.2.1"
  protocol_http_port = 3128
}

resource "junos_services_proxy_profile" "testacc_services_proxy_profile2" {
  name               = "testacc Services pr%xy_profile2!"
  protocol_http_host = "test.url@char"
}

resource "junos_services_proxy_profile" "testacc_services_proxy_profile3" {
  name               = "testacc_services_proxy_profile3"
  protocol_http_host = "fe80::fe80"
}
