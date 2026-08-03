resource "junos_system_radius_server" "testacc_radiusServer_wo" {
  address                             = "192.0.2.11"
  secret_wo                           = "a@Secret"
  secret_wo_version                   = 1
  preauthentication_secret_wo         = "a@PreauthSecret"
  preauthentication_secret_wo_version = 1
}
