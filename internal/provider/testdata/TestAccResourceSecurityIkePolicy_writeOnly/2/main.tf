resource "junos_security_ike_proposal" "testacc_ikepol_wo" {
  name                     = "testacc_ikepol_wo"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group2"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol_wo" {
  name                           = "testacc_ikepol_wo"
  proposals                      = [junos_security_ike_proposal.testacc_ikepol_wo.name]
  mode                           = "main"
  pre_shared_key_text_wo         = "thePassWord"
  pre_shared_key_text_wo_version = 1
}
# read the configuration to check that the write-only key has been sent to the device,
# the resource is unchanged in this step so the data source is read during the plan
data "junos_config_raw" "testacc_ikepol_wo" {
  format = "set"
}
