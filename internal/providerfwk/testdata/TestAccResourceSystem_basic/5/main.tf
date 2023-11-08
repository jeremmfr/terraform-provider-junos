resource "junos_system" "testacc_system" {
  host_name = "testacc-terraform"
  services {
    ssh {
      root_login = "allow"
    }
  }
  name_server_opts {
    address = "192.0.2.10"
  }
  name_server_opts {
    address = "192.0.2.11"
  }
  time_zone = "Europe/Paris"
}

resource "junos_routing_instance" "testacc_system" {
  name = "testacc_system"
}
