package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccJunosSecurity_basic(t *testing.T) {
	testaccSecurity := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccSecurity = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:             testAccJunosSecurityConfigPreCreate(),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: testAccJunosSecurityConfigCreate(testaccSecurity),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"alg.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.advanced_options.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.advanced_options.0.drop_matching_reserved_ip_address", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.advanced_options.0.drop_matching_link_local_address", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.advanced_options.0.reverse_route_packet_mode_vr", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.aging.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.aging.0.early_ageout", "10"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.aging.0.high_watermark", "90"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.aging.0.low_watermark", "80"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.allow_dns_reply", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.allow_embedded_icmp", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.allow_reverse_ecmp", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.enable_reroute_uniform_link_check_nat", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ethernet_switching.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ethernet_switching.0.block_non_ip_all", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ethernet_switching.0.bpdu_vlan_flooding", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ethernet_switching.0.no_packet_flooding.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.force_ip_reassembly", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ipsec_performance_acceleration", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.mcast_buffer_enhance", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.pending_sess_queue_length", "normal"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.preserve_incoming_fragment_size", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.route_change_timeout", "10"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.syn_flood_protection_mode", "syn-proxy"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.sync_icmp_session", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.all_tcp_mss", "1499"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.gre_in.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.gre_out.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.ipsec_vpn.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.fin_invalidate_session", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.maximum_window", "512K"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.no_sequence_check", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.rst_invalidate_session", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.rst_sequence_check", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.strict_syn_check", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.tcp_initial_timeout", "10"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.time_wait_state.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"forwarding_options.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"forwarding_options.0.mpls_mode", "flow-based"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"forwarding_options.0.inet6_mode", "flow-based"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"forwarding_options.0.iso_mode_packet_based", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.name", "ike.log"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.files", "5"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.match", "test"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.size", "102400"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.world_readable", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.flag.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.flag.0", "all"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.rate_limit", "100"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.no_remote_trace", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.disable", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.facility_override", "local7"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.file.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.file.0.files", "10"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.file.0.name", "security.log"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.file.0.path", "/"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.file.0.size", "10"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.format", "syslog"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.mode", "event"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.report", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.source_interface", testaccSecurity+".0"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.transport.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.transport.0.protocol", "tcp"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.transport.0.tcp_connections", "5"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.transport.0.tls_profile", "testacc_security"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.utc_timestamp", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"policies.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"policies.0.policy_rematch", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"utm.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"utm.0.feature_profile_web_filtering_type", "juniper-enhanced"),
					),
				},
				{
					ResourceName:      "junos_security.testacc_security",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSecurityConfigUpdate(testaccSecurity),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ethernet_switching.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ethernet_switching.0.bypass_non_ip_unicast", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ethernet_switching.0.no_packet_flooding.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.ethernet_switching.0.no_packet_flooding.0.no_trace_route", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.gre_in.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.gre_in.0.mss", "1399"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.gre_out.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.gre_out.0.mss", "1399"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.ipsec_vpn.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_mss.0.ipsec_vpn.0.mss", "1399"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.no_syn_check", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.no_syn_check_in_tunnel", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.time_wait_state.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.time_wait_state.0.apply_to_half_close_state", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.time_wait_state.0.session_ageout", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.match", ""),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.no_world_readable", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.flag.#", "0"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.event_rate", "100"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.max_database_record", "1000"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.rate_cap", "100"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"log.0.source_address", "192.0.2.1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"policies.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"policies.0.policy_rematch_extensive", "true"),
					),
				},
				{
					Config: testAccJunosSecurityConfigUpdate2(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.time_wait_state.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"flow.0.tcp_session.0.time_wait_state.0.session_timeout", "90"),
					),
				},
				{
					Config: testAccJunosSecurityConfigPostCheck(),
				},
			},
		})
	}
}

func testAccJunosSecurityConfigPreCreate() string {
	return `
resource "junos_system" "system" {
  tracing_dest_override_syslog_host = "192.0.2.13"
}
`
}

func testAccJunosSecurityConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_security" {
  lifecycle {
    create_before_destroy = true
  }
  name        = "%s.0"
  description = "testacc_security"
}
resource "junos_services_proxy_profile" "testacc_security" {
  lifecycle {
    create_before_destroy = true
  }
  name               = "testacc_security"
  protocol_http_host = "192.0.2.11"
  protocol_http_port = 3128
}
resource "junos_services_ssl_initiation_profile" "testacc_security" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_security"
}

resource "junos_security" "testacc_security" {
  alg {
    dns_disable    = true
    ftp_disable    = true
    h323_disable   = true
    mgcp_disable   = true
    msrpc_disable  = true
    pptp_disable   = true
    rsh_disable    = true
    rtsp_disable   = true
    sccp_disable   = true
    sip_disable    = true
    sql_disable    = true
    sunrpc_disable = true
    talk_disable   = true
    tftp_disable   = true
  }
  flow {
    advanced_options {
      drop_matching_reserved_ip_address = true
      drop_matching_link_local_address  = true
      reverse_route_packet_mode_vr      = true
    }
    aging {
      early_ageout   = 10
      high_watermark = 90
      low_watermark  = 80
    }
    allow_dns_reply                       = true
    allow_embedded_icmp                   = true
    allow_reverse_ecmp                    = true
    enable_reroute_uniform_link_check_nat = true
    ethernet_switching {
      block_non_ip_all   = true
      bpdu_vlan_flooding = true
      no_packet_flooding {}
    }
    force_ip_reassembly             = true
    ipsec_performance_acceleration  = true
    mcast_buffer_enhance            = true
    pending_sess_queue_length       = "normal"
    preserve_incoming_fragment_size = true
    route_change_timeout            = 10
    syn_flood_protection_mode       = "syn-proxy"
    sync_icmp_session               = true
    tcp_mss {
      all_tcp_mss = 1499
      gre_in {}
      gre_out {}
      ipsec_vpn {}
    }
    tcp_session {
      fin_invalidate_session = true
      maximum_window         = "512K"
      no_sequence_check      = true
      rst_invalidate_session = true
      rst_sequence_check     = true
      strict_syn_check       = true
      tcp_initial_timeout    = 10
      time_wait_state {}
    }
  }
  forwarding_options {
    inet6_mode            = "flow-based"
    mpls_mode             = "flow-based"
    iso_mode_packet_based = "true"
  }
  forwarding_process {
    enhanced_services_mode = true
  }
  idp_security_package {
    automatic_enable             = true
    automatic_interval           = 24
    automatic_start_time         = "2016-1-1.02:00:00"
    install_ignore_version_check = true
    proxy_profile                = junos_services_proxy_profile.testacc_security.name
    source_address               = "192.0.2.6"
    url                          = "https://signatures.juniper.net/cgi-bin/index.cgi"
  }
  idp_sensor_configuration {
    log_cache_size = 10
    log_suppression {
      disable                        = true
      no_include_destination_address = true
      max_logs_operate               = 1000
      max_time_report                = 30
      start_log                      = 35
    }
    packet_log {
      source_address             = "192.0.2.4"
      host_address               = "192.0.2.5"
      host_port                  = 100
      max_sessions               = 10
      threshold_logging_interval = 20
      total_memory               = 25
    }
    security_configuration_protection_mode = "datacenter"
  }
  ike_traceoptions {
    file {
      name           = "ike.log"
      files          = 5
      match          = "test"
      size           = 102400
      world_readable = true
    }
    flag            = ["all"]
    rate_limit      = 100
    no_remote_trace = true
  }
  log {
    disable           = true
    facility_override = "local7"
    file {
      files = 10
      name  = "security.log"
      path  = "/"
      size  = 10
    }
    format           = "syslog"
    mode             = "event"
    report           = true
    source_interface = junos_interface_logical.testacc_security.name
    transport {
      protocol        = "tcp"
      tcp_connections = 5
      tls_profile     = junos_services_ssl_initiation_profile.testacc_security.name
    }
    utc_timestamp = true
  }
  policies {
    policy_rematch = true
  }
  user_identification_auth_source {
    ad_auth_priority                = 1
    aruba_clearpass_priority        = 2
    firewall_auth_priority          = 3
    local_auth_priority             = 4
    unified_access_control_priority = 0
  }
  utm {
    feature_profile_web_filtering_type = "juniper-enhanced"
    feature_profile_web_filtering_juniper_enhanced_server {
      host = "192.0.2.1"
      port = 1500
    }
  }
}
`, interFace)
}

func testAccJunosSecurityConfigUpdate(interFace string) string {
	return `
resource "junos_security" "testacc_security" {
  flow {
    ethernet_switching {
      bypass_non_ip_unicast = true
      no_packet_flooding {
        no_trace_route = true
      }
    }
    tcp_mss {
      all_tcp_mss = 1499
      gre_in {
        mss = 1399
      }
      gre_out {
        mss = 1399
      }
      ipsec_vpn {
        mss = 1399
      }
    }
    tcp_session {
      no_syn_check           = true
      no_syn_check_in_tunnel = true
      time_wait_state {
        apply_to_half_close_state = true
        session_ageout            = true
      }
    }
  }
  idp_sensor_configuration {
    log_suppression {
      include_destination_address = true
    }
    packet_log {
      source_address = "192.0.2.4"
    }
  }
  ike_traceoptions {
    file {
      name              = "ike.log"
      files             = 5
      size              = 100000
      no_world_readable = true
    }
    rate_limit = 100
    # no_remote_trace = true
  }
  log {
    mode                = "event"
    event_rate          = 100
    max_database_record = 1000
    rate_cap            = 100
    source_address      = "192.0.2.1"
  }
  policies {
    policy_rematch_extensive = true
  }
  utm {
    feature_profile_web_filtering_type = "juniper-enhanced"
    feature_profile_web_filtering_juniper_enhanced_server {
      host             = "192.0.2.1"
      port             = 1500
      routing_instance = junos_routing_instance.testacc_security.name
    }
  }
}
resource "junos_routing_instance" "testacc_security" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_security"
}
`
}

func testAccJunosSecurityConfigUpdate2() string {
	return `
resource "junos_security" "testacc_security" {
  flow {
    tcp_session {
      time_wait_state {
        session_timeout = 90
      }
    }
  }
  idp_sensor_configuration {
    log_suppression {}
  }
}
`
}

func testAccJunosSecurityConfigPostCheck() string {
	return `
resource "junos_security" "testacc_security" {
  alg {
    dns_disable    = true
    ftp_disable    = true
    h323_disable   = true
    mgcp_disable   = true
    msrpc_disable  = true
    pptp_disable   = true
    rsh_disable    = true
    rtsp_disable   = true
    sccp_disable   = true
    sip_disable    = true
    sql_disable    = true
    sunrpc_disable = true
    talk_disable   = true
    tftp_disable   = true
  }
}
`
}
