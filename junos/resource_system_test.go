package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSystem_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSystemConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"host_name", "testacc-terraform"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"authentication_order.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"authentication_order.0", "password"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"auto_snapshot", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"default_address_selection", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"inet6_backup_router.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"inet6_backup_router.0.destination.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"inet6_backup_router.0.destination.*", "::/0"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"inet6_backup_router.0.address", "fe80::1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.gre_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.icmpv4_rate_limit.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.icmpv4_rate_limit.0.bucket_size", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.icmpv4_rate_limit.0.packet_rate", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.icmpv6_rate_limit.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.icmpv6_rate_limit.0.bucket_size", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.icmpv6_rate_limit.0.packet_rate", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.ipip_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.ipv6_duplicate_addr_detection_transmits", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.ipv6_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.ipv6_path_mtu_discovery_timeout", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.ipv6_reject_zero_hop_limit", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.source_port_upper_limit", "50000"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.source_quench", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.tcp_drop_synfin_set", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.tcp_mss", "1400"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.0.autoupdate_password", "some_password"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.0.autoupdate_url", "some_url"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.0.renew_interval", "24"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.0.renew_before_expiration", "30"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"login.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"login.0.deny_sources_address.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"max_configuration_rollbacks", "49"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"max_configurations_on_flash", "49"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"name_server.#", "2"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"name_server.0", "192.0.2.10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"name_server.1", "192.0.2.11"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"no_multicast_echo", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"no_ping_record_route", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"no_ping_time_stamp", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"no_redirects", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"no_redirects_ipv6", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.authentication_order.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.authentication_order.0", "password"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.ciphers.#", "2"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.0.ssh.0.ciphers.*", "aes256-ctr"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.0.ssh.0.ciphers.*", "aes256-cbc"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.client_alive_count_max", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.client_alive_interval", "30"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.connection_limit", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.fingerprint_hash", "md5"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.hostkey_algorithm.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.0.ssh.0.hostkey_algorithm.*", "no-ssh-dss"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.key_exchange.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.0.ssh.0.key_exchange.*", "ecdh-sha2-nistp256"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.macs.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.0.ssh.0.macs.*", "hmac-sha2-256"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.max_pre_authentication_packets", "10000"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.max_sessions_per_connection", "100"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.port", "22"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.protocol_version.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.0.ssh.0.protocol_version.*", "v2"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.rate_limit", "200"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.root_login", "deny"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.tcp_forwarding", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.web_management_http.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.web_management_http.0.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.0.web_management_http.0.interface.*", "fxp0.0"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.web_management_http.0.port", "80"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.web_management_https.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.web_management_https.0.port", "443"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.web_management_https.0.system_generated_certificate", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.web_management_https.0.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.0.web_management_https.0.interface.*", "fxp0.0"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.archive.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.archive.0.binary_data", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.archive.0.files", "5"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.archive.0.size", "10000000"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.archive.0.no_world_readable", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.log_rotate_frequency", "30"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.source_address", "192.0.2.1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"time_zone", "Europe/Paris"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"tracing_dest_override_syslog_host", "192.0.2.50"),
					),
				},
				{
					ResourceName:      "junos_system.testacc_system",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSystemConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_gre_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_ipip_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_ipv6_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_ipv6_reject_zero_hop_limit", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_source_quench", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_tcp_reset", "drop-tcp-with-syn-only"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_tcp_rfc1323", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.0.no_tcp_rfc1323_paws", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.ciphers.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.no_tcp_forwarding", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.archive.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.archive.0.no_binary_data", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.0.archive.0.world_readable", "true"),
					),
				},
				{
					Config: testAccJunosSystemPostTest(),
				},
			},
		})
	}
}

// nolint: lll, nolintlint
func testAccJunosSystemConfigCreate() string {
	return `
data "junos_system_information" "srx" {}
locals {
  netconfSSHCltAliveCountMax    = "${tonumber(replace(data.junos_system_information.srx.os_version, "/\\..*$/", "")) >= 21 ? 100 : null}"
  netconfSSHClientAliveInterval = "${tonumber(replace(data.junos_system_information.srx.os_version, "/\\..*$/", "")) >= 21 ? 1000 : null}"
}
resource "junos_system" "testacc_system" {
  host_name = "testacc-terraform"
  archival_configuration {
    archive_site {
      url      = "scp://juniper-configs@192.0.2.30:/destination/directory"
      password = "password/&"
    }
    transfer_interval = 1440
  }
  authentication_order      = ["password"]
  auto_snapshot             = true
  default_address_selection = true
  domain_name               = "domain.local"
  inet6_backup_router {
    destination = ["::/0"]
    address     = "fe80::1"
  }
  internet_options {
    gre_path_mtu_discovery = true
    icmpv4_rate_limit {
      bucket_size = 10
      packet_rate = 10
    }
    icmpv6_rate_limit {
      bucket_size = 10
      packet_rate = 10
    }
    ipip_path_mtu_discovery                 = true
    ipv6_duplicate_addr_detection_transmits = 10
    ipv6_path_mtu_discovery                 = true
    ipv6_path_mtu_discovery_timeout         = 10
    ipv6_reject_zero_hop_limit              = true
    path_mtu_discovery                      = true
    source_port_upper_limit                 = 50000
    source_quench                           = true
    tcp_drop_synfin_set                     = true
    tcp_mss                                 = 1400
  }
  license {
    autoupdate              = true
    autoupdate_password     = "some_password"
    autoupdate_url          = "some_url"
    renew_interval          = 24
    renew_before_expiration = 30
  }
  login {
    announcement         = "test announce"
    deny_sources_address = ["127.0.0.1"]
    idle_timeout         = 60
    message              = "test message"
    password {
      change_type               = "character-sets"
      format                    = "sha512"
      maximum_length            = 128
      minimum_changes           = 1
      minimum_character_changes = 4
      minimum_length            = 6
      minimum_lower_cases       = 1
      minimum_numerics          = 1
      minimum_punctuations      = 1
      minimum_reuse             = 1
      minimum_upper_cases       = 1
    }
    retry_options {
      backoff_factor          = 5
      backoff_threshold       = 1
      lockout_period          = 1
      maximum_time            = 300
      minimum_time            = 20
      tries_before_disconnect = 10
    }
  }
  max_configuration_rollbacks = 49
  max_configurations_on_flash = 49
  name_server                 = ["192.0.2.10", "192.0.2.11"]
  no_multicast_echo           = true
  no_ping_record_route        = true
  no_ping_time_stamp          = true
  no_redirects                = true
  no_redirects_ipv6           = true
  ntp {
    boot_server              = "192.0.2.13"
    broadcast_client         = true
    interval_range           = 2
    multicast_client         = true
    multicast_client_address = "224.0.0.3"
    threshold_action         = "accept"
    threshold_value          = 30
  }
  ports {
    auxiliary_authentication_order = ["password", "radius"]
    auxiliary_disable              = true
    auxiliary_insecure             = true
    auxiliary_logout_on_disconnect = true
    auxiliary_type                 = "vt100"
    console_authentication_order   = ["radius", "password"]
    console_disable                = true
    console_insecure               = true
    console_logout_on_disconnect   = true
    console_type                   = "vt100"
  }
  radius_options_attributes_nas_ipaddress   = "192.0.2.12"
  radius_options_enhanced_accounting        = true
  radius_options_password_protocol_mschapv2 = true
  services {
    netconf_ssh {
      client_alive_count_max = local.netconfSSHCltAliveCountMax
      client_alive_interval  = local.netconfSSHClientAliveInterval
      connection_limit       = 200
      rate_limit             = 200
    }
    netconf_traceoptions {
      file_name           = "testacc_netconf"
      file_match          = "test"
      file_world_readable = true
      file_size           = 20480
      flag                = ["all"]
      no_remote_trace     = true
      on_demand           = true
    }
    ssh {
      authentication_order           = ["password"]
      ciphers                        = ["aes256-ctr", "aes256-cbc"]
      client_alive_count_max         = 10
      client_alive_interval          = 30
      connection_limit               = 10
      fingerprint_hash               = "md5"
      hostkey_algorithm              = ["no-ssh-dss"]
      key_exchange                   = ["ecdh-sha2-nistp256"]
      macs                           = ["hmac-sha2-256"]
      max_pre_authentication_packets = 10000
      max_sessions_per_connection    = 100
      port                           = 22
      protocol_version               = ["v2"]
      rate_limit                     = 200
      root_login                     = "deny"
      tcp_forwarding                 = true
    }
    web_management_http {
      interface = ["fxp0.0"]
      port      = 80
    }
    web_management_https {
      interface                    = ["fxp0.0"]
      system_generated_certificate = true
      port                         = 443
    }
  }
  syslog {
    archive {
      binary_data       = true
      files             = 5
      size              = 10000000
      no_world_readable = true
    }
    console {
      any_severity                 = "emergency"
      authorization_severity       = "none"
      changelog_severity           = "emergency"
      conflictlog_severity         = "error"
      daemon_severity              = "none"
      dfc_severity                 = "alert"
      external_severity            = "any"
      firewall_severity            = "info"
      ftp_severity                 = "none"
      interactivecommands_severity = "critical"
      kernel_severity              = "emergency"
      ntp_severity                 = "emergency"
      pfe_severity                 = "emergency"
      security_severity            = "emergency"
      user_severity                = "emergency"
    }
    log_rotate_frequency    = 30
    source_address          = "192.0.2.1"
    time_format_millisecond = true
    time_format_year        = true
  }
  time_zone                         = "Europe/Paris"
  tracing_dest_override_syslog_host = "192.0.2.50"
}
`
}

func testAccJunosSystemConfigUpdate() string {
	return `
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
`
}

func testAccJunosSystemPostTest() string {
	return `
resource "junos_system" "testacc_system" {
  host_name = "testacc-terraform"
  services {
    ssh {
      root_login = "allow"
    }
  }
  time_zone = "Europe/Paris"
}
`
}
