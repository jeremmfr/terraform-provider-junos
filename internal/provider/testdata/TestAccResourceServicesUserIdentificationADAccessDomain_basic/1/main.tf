resource "junos_services" "testacc_userID_addomain" {
  clean_on_destroy = true
  user_identification {
    ad_access {}
  }
}
resource "junos_services_user_identification_ad_access_domain" "testacc_userID_addomain" {
  name          = "testacc_userID_addomain.local"
  user_name     = "user_dom"
  user_password = "user_pass"
  domain_controller {
    name    = "server1"
    address = "192.0.2.3"
  }
  ip_user_mapping_discovery_wmi {}
  user_group_mapping_ldap {
    base = "CN=xxx"
  }
}
