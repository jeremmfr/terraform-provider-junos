resource "junos_system" "testacc_system" {
  host_name = "testacc-terraform"

  accounting {
    events              = ["login"]
    destination_radius  = true
    destination_tacplus = true
    destination_radius_server {
      address                  = "192.0.2.53"
      secret                   = "wordPass"
      preauthentication_secret = "passWord"
      source_address           = "192.0.2.54"
      port                     = 1645
      accounting_port          = 1646
      dynamic_request_port     = 3799
      preauthentication_port   = 1812
      timeout                  = 11
      accounting_timeout       = 5
      retry                    = 3
      accounting_retry         = 2
      max_outstanding_requests = 1000
      routing_instance         = junos_routing_instance.testacc_system.name
    }
    destination_radius_server {
      address = "192.0.2.43"
      secret  = "aPass"
    }
    destination_tacplus_server {
      address           = "192.0.2.55"
      secret            = "password"
      source_address    = "192.0.2.56"
      port              = 49
      timeout           = 12
      single_connection = true
      routing_instance  = junos_routing_instance.testacc_system.name
    }
    destination_tacplus_server {
      address = "192.0.2.45"
    }
  }
  archival_configuration {
    archive_site {
      url      = "scp://juniper-configs@192.0.2.30:/destination/directory"
      password = "password/&"
    }
    archive_site {
      url = "http://juniper-configs@192.0.2.30:/destination/directory"
    }
    transfer_on_commit = true
  }
  name_server_opts {
    address          = "192.0.2.10"
    routing_instance = junos_routing_instance.testacc_system.name
  }
  internet_options {
    no_gre_path_mtu_discovery     = true
    no_ipip_path_mtu_discovery    = true
    no_ipv6_path_mtu_discovery    = true
    no_ipv6_reject_zero_hop_limit = true
    no_path_mtu_discovery         = true
    no_source_quench              = true
    no_tcp_reset                  = "drop-tcp-with-syn-only"
    no_tcp_rfc1323                = true
    no_tcp_rfc1323_paws           = true
  }
  services {
    netconf_traceoptions {
      file_name              = "testacc_netconf"
      file_no_world_readable = true
      file_size              = 40960
      flag                   = ["incoming", "outgoing"]
    }
    ssh {
      ciphers           = ["aes256-ctr"]
      no_tcp_forwarding = true
    }
    web_management_http {}
    web_management_https {
      system_generated_certificate = true
    }
  }
  syslog {
    archive {
      no_binary_data = true
      files          = 5
      size           = 10000000
      world_readable = true
    }
    console {
      any_severity = "emergency"
    }
  }
  tacplus_options_authorization_time_interval = 900
  tacplus_options_no_cmd_attribute_value      = true
  tacplus_options_no_strict_authorization     = true
  time_zone                                   = "Europe/Paris"
}

resource "junos_routing_instance" "testacc_system" {
  lifecycle {
    create_before_destroy = true
  }

  name = "testacc_system"
}
