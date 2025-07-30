resource "junos_services_user_identification_device_identity_profile" "profile" {
  lifecycle {
    create_before_destroy = true
  }
  name   = "testacc_securityPolicy"
  domain = "testacc_securityPolicy"
  attribute {
    name  = "device-identity"
    value = ["testacc_securityPolicy"]
  }
}
resource "junos_security_policy" "testacc_securityPolicy" {
  from_zone = junos_security_zone.testacc_seczonePolicy1.name
  to_zone   = junos_security_zone.testacc_seczonePolicy1.name
  policy {
    name                          = "testacc_Policy_1"
    match_source_address          = ["testacc_address1"]
    match_destination_address     = ["any"]
    match_application             = ["junos-ssh"]
    match_dynamic_application     = ["any"]
    match_source_end_user_profile = junos_services_user_identification_device_identity_profile.profile.name
    log_init                      = true
    log_close                     = true
    count                         = true
  }
}

resource "junos_security_zone" "testacc_seczonePolicy1" {
  name = "testacc_seczonePolicy1"
  address_book {
    name    = "testacc_address1"
    network = "192.0.2.0/25"
  }
}
