package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurity_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config:             testAccJunosSecurityConfigPreCreate(),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: testAccJunosSecurityConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
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
							"ike_traceoptions.0.file.0.size", "100000"),
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
							"utm.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"utm.0.feature_profile_web_filtering_type", "juniper-enhanced"),
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
					),
				},
				{
					ResourceName:      "junos_security.testacc_security",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSecurityConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
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
			},
		})
	}
}

func testAccJunosSecurityConfigPreCreate() string {
	return `
resource junos_system "system" {
  tracing_dest_override_syslog_host = "192.0.2.13"
}
`
}

func testAccJunosSecurityConfigCreate() string {
	return `
resource junos_security "testacc_security" {
  ike_traceoptions {
    file {
      name           = "ike.log"
      files          = 5
      match          = "test"
      size           = 100000
      world_readable = true
    }
    flag       = ["all"]
    rate_limit = 100
    no_remote_trace = true
  }
  utm {
    feature_profile_web_filtering_type = "juniper-enhanced"
  }
  alg {
    dns_disable    = true
    ftp_disable    = true
    msrpc_disable  = true
    pptp_disable   = true
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
      block_non_ip_all = true
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
}
`
}
func testAccJunosSecurityConfigUpdate() string {
	return `
resource junos_security "testacc_security" {
  ike_traceoptions {
    file {
      name           = "ike.log"
      files          = 5
      size           = 100000
      no_world_readable = true
    }
    rate_limit = 100
    # no_remote_trace = true
  }
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
}
`
}

func testAccJunosSecurityConfigUpdate2() string {
	return `
resource junos_security "testacc_security" {
  flow {
    tcp_session {
      time_wait_state {
        session_timeout = 90
      }
    }
  }
}
`
}
