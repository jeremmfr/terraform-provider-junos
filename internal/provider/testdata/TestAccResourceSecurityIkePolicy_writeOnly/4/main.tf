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
  pre_shared_key_hexa_wo         = "0123456789abcdef"
  pre_shared_key_hexa_wo_version = 1
}
