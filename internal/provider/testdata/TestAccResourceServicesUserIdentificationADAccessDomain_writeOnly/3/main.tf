resource "junos_services" "testacc_userID_addomain_wo" {
  clean_on_destroy = true
  user_identification {
    ad_access {}
  }
}
resource "junos_services_user_identification_ad_access_domain" "testacc_userID_addomain_wo" {
  name                     = "testacc_userID_addomain_wo.local"
  user_name                = "user_dom"
  user_password_wo         = "userPassTwo"
  user_password_wo_version = 2
  domain_controller {
    name    = "server1"
    address = "192.0.2.3"
  }
  user_group_mapping_ldap {
    base                     = "CN=xxx"
    user_name                = "user_ldap"
    user_password_wo         = "ldapPassTwo"
    user_password_wo_version = 2
  }
}
