resource "junos_system" "testacc_system" {
  host_name = "testacc-terraform"
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
  name_server = ["192.0.2.10"]
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
  time_zone = "Europe/Paris"
}
