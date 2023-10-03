resource "junos_security_zone" "testacc_secglobpolicy1" {
  name = "testacc_secglobpolicy1"
}
resource "junos_security_zone" "testacc_secglobpolicy2" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_secglobpolicy2"
}
resource "junos_security_address_book" "testacc_secglobpolicy" {
  network_address {
    name  = "blue"
    value = "192.0.2.1/32"
  }
  network_address {
    name  = "green"
    value = "192.0.2.2/32"
  }
}
resource "junos_services_user_identification_device_identity_profile" "profile" {
  lifecycle {
    create_before_destroy = true
  }
  name   = "testacc_secglobpolicy"
  domain = "testacc_secglobpolicy"
  attribute {
    name  = "device-identity"
    value = ["testacc_secglobpolicy"]
  }
}
resource "junos_security_global_policy" "testacc_secglobpolicy" {
  depends_on = [
    junos_security_address_book.testacc_secglobpolicy
  ]
  policy {
    name                               = "test"
    match_source_address               = ["blue"]
    match_destination_address          = ["green"]
    match_destination_address_excluded = true
    match_application                  = ["any"]
    match_dynamic_application          = ["any"]
    match_source_end_user_profile      = junos_services_user_identification_device_identity_profile.profile.name
    match_from_zone                    = [junos_security_zone.testacc_secglobpolicy1.name]
    match_to_zone                      = [junos_security_zone.testacc_secglobpolicy2.name]
  }
}
