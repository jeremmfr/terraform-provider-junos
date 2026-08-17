resource "junos_services" "testacc_userID_addomain_wo" {
  clean_on_destroy = true
  user_identification {
    ad_access {}
  }
}
resource "junos_services_user_identification_ad_access_domain" "testacc_userID_addomain_wo" {
  name                     = "testacc_userID_addomain_wo.local"
  user_name                = "user_dom"
  user_password_wo         = "userPassOne"
  user_password_wo_version = 1
  domain_controller {
    name    = "server1"
    address = "192.0.2.3"
  }
  user_group_mapping_ldap {
    base                     = "CN=xxx"
    user_name                = "user_ldap"
    user_password_wo         = "ldapPassOne"
    user_password_wo_version = 1
  }
}
# read the configuration to check that the write-only passwords have been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_userID_addomain_wo" {
  format = "set"
}
