resource "junos_system_tacplus_server" "testacc_tacplusServer_wo" {
  address           = "192.0.2.11"
  secret_wo         = "a@Secret"
  secret_wo_version = 1
}
