resource "junos_system_radius_server" "testacc_radiusServer" {
  address = "192.0.2.1"
  secret  = "password"
}
