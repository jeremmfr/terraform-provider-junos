resource "junos_system_radius_server" "testacc_radiusServer_wo" {
  address                             = "192.0.2.11"
  secret_wo                           = "a@Secret2"
  secret_wo_version                   = 2
  preauthentication_secret_wo         = "a@PreauthSecret2"
  preauthentication_secret_wo_version = 2
}
