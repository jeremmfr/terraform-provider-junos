
resource "junos_security_zone" "testacc_secglobpolicy1" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_secglobpolicy1"
}
resource "junos_security_idp_policy" "testacc_secglobpolicy" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_secglobpolicy"
}
resource "junos_security_address_book" "testacc_secglobpolicy" {
  lifecycle {
    create_before_destroy = true
  }
  network_address {
    name  = "blue"
    value = "192.0.2.1/32"
  }
  network_address {
    name  = "green"
    value = "192.0.2.2/32"
  }
}
resource "junos_security_global_policy" "testacc_secglobpolicy" {
  depends_on = [
    junos_security_address_book.testacc_secglobpolicy
  ]
  policy {
    name                      = "test"
    match_source_address      = ["blue"]
    match_destination_address = ["any"]
    match_application         = ["any"]
    match_from_zone           = [junos_security_zone.testacc_secglobpolicy1.name]
    match_to_zone             = [junos_security_zone.testacc_secglobpolicy1.name]
    count                     = true
    log_init                  = true
    log_close                 = true
    permit_application_services {
      idp_policy = junos_security_idp_policy.testacc_secglobpolicy.name
    }
  }
}
