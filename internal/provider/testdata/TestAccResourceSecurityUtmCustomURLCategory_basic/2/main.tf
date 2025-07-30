resource "junos_security_utm_custom_url_pattern" "testacc_URLCategory1" {
  name  = "testacc-custom-pattern1"
  value = ["*.google.com"]
}
resource "junos_security_utm_custom_url_pattern" "testacc_URLCategory2" {
  name  = "testacc-custom-pattern2"
  value = ["*.google.fr"]
}
resource "junos_security_utm_custom_url_category" "testacc_URLCategory" {
  name = "testacc_URLCategory"
  value = [
    junos_security_utm_custom_url_pattern.testacc_URLCategory1.name,
    junos_security_utm_custom_url_pattern.testacc_URLCategory2.name,
  ]
}
