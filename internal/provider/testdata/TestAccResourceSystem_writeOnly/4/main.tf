resource "junos_system" "testacc_system_wo" {
  host_name = "testacc-terraform-wo"

  name_server_opts {
    address = "192.0.2.10"
  }
  name_server_opts {
    address = "192.0.2.11"
  }
  services {
    ssh {
      root_login = "allow"
    }
  }
  time_zone = "Europe/Paris"
}
