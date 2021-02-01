package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityScreen_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityScreenConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_screen.testacc_securityScreen",
							"icmp.#", "1"),
						resource.TestCheckResourceAttr("junos_security_screen.testacc_securityScreen",
							"ip.#", "1"),
						resource.TestCheckResourceAttr("junos_security_screen.testacc_securityScreen",
							"limit_session.#", "1"),
						resource.TestCheckResourceAttr("junos_security_screen.testacc_securityScreen",
							"tcp.#", "1"),
						resource.TestCheckResourceAttr("junos_security_screen.testacc_securityScreen",
							"udp.#", "1"),
						resource.TestCheckResourceAttr("junos_security_screen_whitelist.testacc1",
							"address.#", "2"),
					),
				},
				{
					Config: testAccJunosSecurityScreenConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_screen_whitelist.testacc2",
							"address.#", "1"),
					),
				},
				{
					ResourceName:      "junos_security_screen.testacc_securityScreen",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_screen_whitelist.testacc1",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityScreenConfigCreate() string {
	return `
resource junos_security_screen testacc_securityScreen {
  name               = "testacc 1"
  alarm_without_drop = true
  description        = "desc testacc 1"
  icmp {
    flood {}
    fragment         = true
    icmpv6_malformed = true
    large            = true
    ping_death       = true
    sweep {}
  }
  ip {
    bad_option = true
    block_frag = true
    ipv6_extension_header {
      ah_header  = true
      esp_header = true
      hip_header = true
      destination_header {}
      fragment_header = true
      hop_by_hop_header {}
      mobility_header = true
      no_next_header  = true
      routing_header  = true
      shim6_header    = true
      user_defined_header_type = [
        "10 to 20",
        "2 to 5",
        "1",
      ]
    }
    ipv6_extension_header_limit = 32
    ipv6_malformed_header       = true
    loose_source_route_option   = true
    record_route_option         = true
    security_option             = true
    source_route_option         = true
    spoofing                    = true
    stream_option               = true
    strict_source_route_option  = true
    tear_drop                   = true
    timestamp_option            = true
    tunnel {
      bad_inner_header = true
      gre {
        gre_4in4 = true
        gre_4in6 = true
        gre_6in4 = true
        gre_6in6 = true
      }
      ip_in_udp_teredo = true
      ipip {
        ipip_4in4      = true
        ipip_4in6      = true
        ipip_6in4      = true
        ipip_6in6      = true
        ipip_6over4    = true
        ipip_6to4relay = true
        dslite         = true
        isatap         = true
      }
    }
    unknown_protocol = true
  }
  limit_session {
    destination_ip_based = 2000
    source_ip_based      = 3000
  }
  tcp {
    fin_no_ack = true
    land       = true
    no_flag    = true
    port_scan {}
    syn_ack_ack_proxy {}
    syn_fin = true
    syn_flood {
      alarm_threshold       = 10011
      attack_threshold      = 10012
      destination_threshold = 10013
      source_threshold      = 10014
      timeout               = 10
      whitelist {
        name                = "test3"
        source_address      = ["192.0.2.0/26"]
        destination_address = ["192.0.2.64/26"]
      }
    }
    syn_frag = true
    sweep {}
    winnuke = true
  }
  udp {
    flood {}
    port_scan {}
    sweep {}
  }
}
resource "junos_security_screen_whitelist" "testacc1" {
  name = "testacc1"
  address = [
    "192.0.2.128/26",
    "192.0.2.64/26",
  ]
}
`
}
func testAccJunosSecurityScreenConfigUpdate() string {
	return `
resource junos_security_screen testacc_securityScreen {
  name               = "testacc 1"
  alarm_without_drop = true
  description        = "desc testacc 1"
  icmp {
    flood {
      threshold = 10000
    }
    fragment         = true
    icmpv6_malformed = true
    large            = true
    ping_death       = true
    sweep {
      threshold = 10000
    }
  }
  ip {
    bad_option = true
    block_frag = true
    ipv6_extension_header {
      ah_header  = true
      esp_header = true
      hip_header = true
      destination_header {
        ilnp_nonce_option                 = true
        home_address_option               = true
        line_identification_option        = true
        tunnel_encapsulation_limit_option = true
        user_defined_option_type = [
          "10 to 20",
          "2 to 5",
          "1",
        ]
      }
      fragment_header = true
      hop_by_hop_header {
        calipso_option       = true
        rpl_option           = true
        smf_dpd_option       = true
        jumbo_payload_option = true
        quick_start_option   = true
        router_alert_option  = true
        user_defined_option_type = [
          "10 to 20",
          "2 to 5",
          "1",
        ]
      }
      mobility_header = true
      no_next_header  = true
      routing_header  = true
      shim6_header    = true
      user_defined_header_type = [
        "10 to 20",
        "2 to 5",
        "1",
      ]
    }
    ipv6_extension_header_limit = 32
    ipv6_malformed_header       = true
    loose_source_route_option   = true
    record_route_option         = true
    security_option             = true
    source_route_option         = true
    spoofing                    = true
    stream_option               = true
    strict_source_route_option  = true
    tear_drop                   = true
    timestamp_option            = true
    tunnel {
      bad_inner_header = true
      gre {
        gre_4in4 = true
        gre_4in6 = true
        gre_6in4 = true
        gre_6in6 = true
      }
      ip_in_udp_teredo = true
      ipip {
        ipip_4in4      = true
        ipip_4in6      = true
        ipip_6in4      = true
        ipip_6in6      = true
        ipip_6over4    = true
        ipip_6to4relay = true
        dslite         = true
        isatap         = true
      }
    }
    unknown_protocol = true
  }
  limit_session {
    destination_ip_based = 2000
    source_ip_based      = 3000
  }
  tcp {
    fin_no_ack = true
    land       = true
    no_flag    = true
    port_scan {
      threshold = 10000
    }
    syn_ack_ack_proxy {
      threshold = 10001
    }
    syn_fin = true
    syn_flood {
      alarm_threshold       = 10011
      attack_threshold      = 10012
      destination_threshold = 10013
      source_threshold      = 10014
      timeout               = 10
      whitelist {
        name                = "test3"
        source_address      = ["192.0.2.0/26"]
        destination_address = ["192.0.2.64/26"]
      }
      whitelist {
        name                = "test2"
        source_address      = ["192.0.2.128/26"]
        destination_address = ["192.0.2.192/26"]
      }
    }
    syn_frag = true
    sweep {
      threshold = 10002
    }
    winnuke = true
  }
  udp {
    flood {
      threshold = 1000
      whitelist = [
        junos_security_screen_whitelist.testacc2.name,
        junos_security_screen_whitelist.testacc1.name,
      ]
    }
    port_scan {
      threshold = 1000
    }
    sweep {
      threshold = 1000
    }
  }
}
resource "junos_security_screen_whitelist" "testacc1" {
  name = "testacc1"
  address = [
    "192.0.2.128/26",
    "192.0.2.64/26",
  ]
}
resource "junos_security_screen_whitelist" "testacc2" {
  name = "testacc2"
  address = [
    "192.0.2.0/26",
  ]
}
`
}
