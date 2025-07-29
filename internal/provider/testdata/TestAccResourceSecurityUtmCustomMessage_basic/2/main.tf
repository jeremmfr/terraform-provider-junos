resource "junos_security_utm_custom_message" "testacc_Message" {
  name    = "testacc-custom_message"
  type    = "user-message"
  content = "a mess@ge for user"
}

