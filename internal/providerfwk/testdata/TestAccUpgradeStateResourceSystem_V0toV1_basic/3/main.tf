resource "junos_system" "testacc_system" {
  host_name = "testacc-terraform"
  services {
    ssh {
      root_login = "allow"
    }
  }
  time_zone = "Europe/Paris"
}
