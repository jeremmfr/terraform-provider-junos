resource "junos_security_utm_profile_web_filtering_websense_redirect" "testacc_ProfileWebFWebS" {
  name = "testacc ProfileWebFWebS"
  category {
    name           = junos_security_utm_custom_url_category.testacc_ProfileWebFWebS2.name
    action         = "block"
    custom_message = junos_security_utm_custom_message.testacc_ProfileWebFWebS.name
  }
  category {
    name   = junos_security_utm_custom_url_category.testacc_ProfileWebFWebS.name
    action = "permit"
  }
  custom_block_message = "Blocked by Juniper"
  custom_message       = junos_security_utm_custom_message.testacc_ProfileWebFWebS.name
  timeout              = 3
  server {
    host             = "10.0.0.1"
    port             = 1024
    routing_instance = junos_routing_instance.testacc_ProfileWebFWebS.name
    source_address   = "10.0.0.2"
  }
  sockets = 16
}
resource "junos_routing_instance" "testacc_ProfileWebFWebS" {
  name = "testacc_ProfileWebFWebS"
}

resource "junos_security_utm_custom_message" "testacc_ProfileWebFWebS" {
  name    = "testacc-profilewebfwebs"
  type    = "user-message"
  content = "testacc_ProfileWebFWebS"
}

resource "junos_security_utm_custom_url_pattern" "testacc_ProfileWebFWebS" {
  name  = "testacc_ProfileWebFWebS"
  value = ["*.google.com"]
}
resource "junos_security_utm_custom_url_category" "testacc_ProfileWebFWebS" {
  name = "testacc_ProfileWebFWebS"
  value = [
    junos_security_utm_custom_url_pattern.testacc_ProfileWebFWebS.name,
  ]
}


resource "junos_security_utm_custom_url_pattern" "testacc_ProfileWebFWebS2" {
  name  = "testacc_ProfileWebFWebS2"
  value = ["api.google.com"]
}
resource "junos_security_utm_custom_url_category" "testacc_ProfileWebFWebS2" {
  name = "testacc_ProfileWebFWebS2"
  value = [
    junos_security_utm_custom_url_pattern.testacc_ProfileWebFWebS2.name,
  ]
}
