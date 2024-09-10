resource "junos_system_radius_server" "testacc_radiusServer" {
  address                  = "192.0.2.1"
  secret                   = "$9$dZV2aZGi.fzDiORSeXxDikqmT"
  preauthentication_secret = "$9$6cgx/pBIRSeMXhS4ZjqQzhSrlK8"
  port                     = 1645
}
