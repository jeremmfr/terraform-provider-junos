resource "junos_security_utm_profile_web_filtering_juniper_local" "testacc_ProfileWebFL" {
  name = "testacc ProfileWebFL"
  category {
    name           = junos_security_utm_custom_url_category.testacc_ProfileWebFL2.name
    action         = "block"
    custom_message = junos_security_utm_custom_message.testacc_ProfileWebFL.name
  }
  category {
    name   = junos_security_utm_custom_url_category.testacc_ProfileWebFL.name
    action = "permit"
  }
  custom_block_message = "Blocked by Juniper"
  custom_message       = junos_security_utm_custom_message.testacc_ProfileWebFL.name
  default_action       = "log-and-permit"
  timeout              = 3
}

resource "junos_security_utm_custom_message" "testacc_ProfileWebFL" {
  name    = "testacc-profilewebfl"
  type    = "user-message"
  content = "testacc_ProfileWebFL"
}

resource "junos_security_utm_custom_url_pattern" "testacc_ProfileWebFL" {
  name  = "testacc-ProfileWebFL"
  value = ["*.google.com"]
}
resource "junos_security_utm_custom_url_category" "testacc_ProfileWebFL" {
  name = "testacc-ProfileWebFL"
  value = [
    junos_security_utm_custom_url_pattern.testacc_ProfileWebFL.name,
  ]
}


resource "junos_security_utm_custom_url_pattern" "testacc_ProfileWebFL2" {
  name  = "testacc-ProfileWebFL2"
  value = ["api.google.com"]
}
resource "junos_security_utm_custom_url_category" "testacc_ProfileWebFL2" {
  name = "testacc-ProfileWebFL2"
  value = [
    junos_security_utm_custom_url_pattern.testacc_ProfileWebFL2.name,
  ]
}
