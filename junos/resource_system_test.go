package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSystem_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSystemConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"name_server.#", "2"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"name_server.0", "192.0.2.10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"name_server.1", "192.0.2.11"),
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
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.ciphers.0", "aes256-ctr"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.ciphers.1", "aes256-cbc"),
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
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.hostkey_algorithm.0", "no-ssh-dss"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.key_exchange.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.key_exchange.0", "ecdh-sha2-nistp256"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.macs.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.macs.0", "hmac-sha2-256"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.max_pre_authentication_packets", "10000"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.max_sessions_per_connection", "100"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.port", "22"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.protocol_version.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.protocol_version.0", "v2"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.rate_limit", "200"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.root_login", "deny"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.0.ssh.0.tcp_forwarding", "true"),
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
			},
		})
	}
}

func testAccJunosSystemConfigCreate() string {
	return `
resource junos_system "testacc_system" {
  name_server = ["192.0.2.10","192.0.2.11"]
  services {
    ssh {
	  authentication_order           = ["password"]
	  ciphers                        = ["aes256-ctr","aes256-cbc"]
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
  }
  syslog {
    archive {
      binary_data       = true
	  files             = 5
	  size              = 10000000
      no_world_readable = true
	}
	log_rotate_frequency = 30
	source_address       = "192.0.2.1"
  }
  tracing_dest_override_syslog_host = "192.0.2.50"
}
`
}
func testAccJunosSystemConfigUpdate() string {
	return `
resource junos_system "testacc_system" {
  name_server = ["192.0.2.10"]
  services {
    ssh {
       ciphers                = ["aes256-ctr"]
       no_tcp_forwarding      = true
    }
  }
  syslog {
    archive {
      no_binary_data = true
      files          = 5
      size           = 10000000
      world_readable = true
    }
  }
}
`
}
