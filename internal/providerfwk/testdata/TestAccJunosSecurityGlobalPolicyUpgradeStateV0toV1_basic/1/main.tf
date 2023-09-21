resource "junos_security_zone" "testacc_v0to1_globalpolicy" {
  name = "testacc_v0to1_globalpolicy"
}

resource "junos_services_advanced_anti_malware_policy" "testacc_v0to1_globalpolicy" {
  name                     = "testacc_v0to1_globalpolicy"
  verdict_threshold        = "recommended"
  default_notification_log = true
}
resource "junos_security_global_policy" "testacc_v0to1_globalpolicy" {
  policy {
    name                      = "test"
    match_source_address      = ["any"]
    match_destination_address = ["any"]
    match_application         = ["any"]
    match_from_zone           = [junos_security_zone.testacc_v0to1_globalpolicy.name]
    match_to_zone             = [junos_security_zone.testacc_v0to1_globalpolicy.name]
    count                     = true
    log_init                  = true
    log_close                 = true
    permit_application_services {
      advanced_anti_malware_policy = junos_services_advanced_anti_malware_policy.testacc_v0to1_globalpolicy.name
      idp                          = true
      redirect_wx                  = true
      ssl_proxy {}
      uac_policy {}
    }
  }
  policy {
    name                      = "test2"
    match_source_address      = ["any"]
    match_destination_address = ["any"]
    match_application         = ["any"]
    match_from_zone           = ["any"]
    match_to_zone             = ["any"]
    then                      = "deny"
  }
}
