package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccResourceSystem_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
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
							"inet6_backup_router.destination.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"inet6_backup_router.destination.*", "::/0"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"inet6_backup_router.address", "fe80::1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.gre_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.icmpv4_rate_limit.bucket_size", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.icmpv4_rate_limit.packet_rate", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.icmpv6_rate_limit.bucket_size", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.icmpv6_rate_limit.packet_rate", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.ipip_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.ipv6_duplicate_addr_detection_transmits", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.ipv6_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.ipv6_path_mtu_discovery_timeout", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.ipv6_reject_zero_hop_limit", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.source_port_upper_limit", "50000"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.source_quench", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.tcp_drop_synfin_set", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.tcp_mss", "1400"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.autoupdate_password", "some_password"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.autoupdate_url", "some_url"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.renew_interval", "24"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"license.renew_before_expiration", "30"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"login.deny_sources_address.#", "1"),
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
							"services.ssh.authentication_order.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.authentication_order.0", "password"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.ciphers.#", "2"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.ssh.ciphers.*", "aes256-ctr"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.ssh.ciphers.*", "aes256-cbc"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.client_alive_count_max", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.client_alive_interval", "30"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.connection_limit", "10"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.fingerprint_hash", "md5"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.hostkey_algorithm.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.ssh.hostkey_algorithm.*", "no-ssh-dss"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.key_exchange.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.ssh.key_exchange.*", "ecdh-sha2-nistp256"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.macs.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.ssh.macs.*", "hmac-sha2-256"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.max_pre_authentication_packets", "10000"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.max_sessions_per_connection", "100"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.port", "22"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.protocol_version.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.ssh.protocol_version.*", "v2"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.rate_limit", "200"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.root_login", "deny"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.tcp_forwarding", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.web_management_http.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.web_management_http.interface.*", "fxp0.0"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.web_management_http.port", "80"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.web_management_https.port", "443"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.web_management_https.system_generated_certificate", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.web_management_https.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_system.testacc_system",
							"services.web_management_https.interface.*", "fxp0.0"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.archive.binary_data", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.archive.files", "5"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.archive.size", "10000000"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.archive.no_world_readable", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.log_rotate_frequency", "30"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.source_address", "192.0.2.1"),
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
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_gre_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_ipip_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_ipv6_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_ipv6_reject_zero_hop_limit", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_path_mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_source_quench", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_tcp_reset", "drop-tcp-with-syn-only"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_tcp_rfc1323", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"internet_options.no_tcp_rfc1323_paws", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.ciphers.#", "1"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"services.ssh.no_tcp_forwarding", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.archive.no_binary_data", "true"),
						resource.TestCheckResourceAttr("junos_system.testacc_system",
							"syslog.archive.world_readable", "true"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
