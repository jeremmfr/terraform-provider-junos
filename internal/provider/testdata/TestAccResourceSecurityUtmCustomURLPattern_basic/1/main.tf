resource "junos_security_utm_custom_url_pattern" "testacc_UrlPattern" {
  name  = "testacc_UrlPattern"
  value = ["*.google.com"]
}
