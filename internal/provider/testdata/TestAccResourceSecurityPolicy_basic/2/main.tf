resource "junos_services_advanced_anti_malware_policy" "testacc_securityPolicy" {
  name                     = "testacc_securityPolicy"
  verdict_threshold        = "recommended"
  default_notification_log = true
}
resource "junos_security_idp_policy" "testacc_securityPolicy" {
  name = "testacc_securityPolicy"
}
resource "junos_security_policy" "testacc_securityPolicy" {
  from_zone = junos_security_zone.testacc_seczonePolicy1.name
  to_zone   = junos_security_zone.testacc_seczonePolicy1.name
  policy {
    name                          = "testacc_Policy_1"
    match_source_address          = ["testacc_address1"]
    match_destination_address     = ["any"]
    match_application             = ["junos-ssh"]
    match_source_address_excluded = true
    log_init                      = true
    log_close                     = true
    count                         = true
    permit_application_services {
      advanced_anti_malware_policy = junos_services_advanced_anti_malware_policy.testacc_securityPolicy.name
      idp_policy                   = junos_security_idp_policy.testacc_securityPolicy.name
      redirect_wx                  = true
      ssl_proxy {}
      uac_policy {}
    }
  }
  policy {
    name                               = "testacc_Policy_2"
    match_source_address               = ["testacc_address1"]
    match_destination_address          = ["testacc_address1"]
    match_destination_address_excluded = true
    match_application                  = ["any"]
    then                               = "reject"
  }
}

resource "junos_security_zone" "testacc_seczonePolicy1" {
  name = "testacc_seczonePolicy1"
  address_book {
    name    = "testacc_address1"
    network = "192.0.2.0/25"
  }
}
