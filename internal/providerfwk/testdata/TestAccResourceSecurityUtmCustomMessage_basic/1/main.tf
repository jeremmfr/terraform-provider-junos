resource "junos_security_utm_custom_message" "testacc_Message" {
  name             = "testacc-custom_message"
  type             = "custom-page"
  custom_page_file = "afile#1.html"
}

resource "junos_security_utm_custom_message" "testacc_Message2" {
  name    = "testacc-custom_message2"
  type    = "redirect-url"
  content = "http://redirect-url.com"
}
