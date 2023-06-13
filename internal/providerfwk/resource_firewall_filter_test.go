package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosFirewallFilter_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosFirewallFilterConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"family", "inet"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"interface_specific", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.#", "2"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.name", "testacc_fwFilter_term1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.address.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.address.*", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.address_except.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.address_except.*", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.port.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.port.*", "22-23"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.prefix_list.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.prefix_list.*", "testacc_fwFilter#1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.prefix_list_except.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.prefix_list_except.*", "testacc_fwFilter#2"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.protocol.*", "tcp"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.tcp_flags", "!0x3"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.is_fragment", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.action", "next term"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.syslog", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.log", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.port_mirror", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.service_accounting", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.icmp_code.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.icmp_code.*", "network-unreachable"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.icmp_type.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.icmp_type.*", "router-advertisement"),
					),
				},
				{
					Config: testAccJunosFirewallFilterConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.#", "5"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.source_address.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.source_address_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.port_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.source_prefix_list.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.source_prefix_list_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.tcp_established", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.protocol_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.then.policer", "testacc_fwfilter#1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.then.action", "accept"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.destination_address.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.destination_address_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.destination_port.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.source_port_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.destination_prefix_list.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.destination_prefix_list_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.tcp_initial", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.then.action", "discard"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.3.from.source_port.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.3.from.destination_port_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.3.then.action", "reject"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.4.from.icmp_code_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.4.from.icmp_type_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"family", "inet6"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.0.from.next_header.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.0.from.next_header.*", "icmp6"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.0.then.action", "discard"),
					),
				},
				{
					ResourceName:      "junos_firewall_filter.testacc_fwFilter",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosFirewallFilterConfigCreate() string {
	return `
resource "junos_firewall_filter" "testacc_fwFilter" {
  name               = "testacc_fwFilter"
  family             = "inet"
  interface_specific = true
  term {
    name = "testacc_fwFilter_term1"
    from {
      address            = ["192.0.2.0/25"]
      address_except     = ["192.0.2.128/25"]
      port               = ["22-23"]
      prefix_list        = [junos_policyoptions_prefix_list.testacc_fwFilter.name]
      prefix_list_except = [junos_policyoptions_prefix_list.testacc_fwFilter2.name]
      protocol           = ["tcp"]
      tcp_flags          = "!0x3"
      is_fragment        = true
    }
    then {
      action             = "next term"
      syslog             = true
      log                = true
      packet_mode        = true
      port_mirror        = true
      service_accounting = true
    }
  }
  term {
    name = "testacc_fwFilter_term2"
    from {
      icmp_code = ["network-unreachable"]
      icmp_type = ["router-advertisement"]
    }
    then {
      action = "accept"
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter_vpls" {
  name   = "testacc_fwFilter vpls"
  family = "vpls"
  term {
    name = "testacc_fwFilter vpls term1"
    from {
      destination_mac_address = [
        "aa:bb:cc:dd:ee:ff/48",
      ]
      destination_mac_address_except = [
        "aa:bb:cc:dd:ee:f0/48",
      ]
      forwarding_class = [
        "best-effort",
      ]
      source_mac_address_except = [
        "aa:bb:cc:dd:ee:01/48",
      ]
      source_mac_address = [
        "aa:bb:cc:dd:ee:02/48",
      ]
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter_any" {
  name   = "testacc_fwFilter any"
  family = "any"
  term {
    name = "testacc_fwFilter any term1"
    from {
      packet_length        = ["1-500"]
      loss_priority_except = ["medium-high"]
    }
  }
}

resource "junos_policyoptions_prefix_list" "testacc_fwFilter" {
  name   = "testacc_fwFilter#1"
  prefix = ["192.0.2.0/25"]
}
resource "junos_policyoptions_prefix_list" "testacc_fwFilter2" {
  name   = "testacc_fwFilter#2"
  prefix = ["192.0.2.128/25"]
}
`
}

func testAccJunosFirewallFilterConfigUpdate() string {
	return `
resource "junos_firewall_filter" "testacc_fwFilter" {
  name               = "testacc_fwFilter"
  family             = "inet"
  interface_specific = true
  term {
    name = "testacc_fwFilter_term1"
    from {
      address            = ["192.0.2.0/25"]
      address_except     = ["192.0.2.128/25"]
      port               = ["22-23"]
      prefix_list        = [junos_policyoptions_prefix_list.testacc_fwFilter.name]
      prefix_list_except = [junos_policyoptions_prefix_list.testacc_fwFilter2.name]
      protocol           = ["tcp"]
      tcp_flags          = "!0x3"
    }
    then {
      action             = "next term"
      syslog             = true
      log                = true
      port_mirror        = true
      service_accounting = true
    }
  }
  term {
    name = "testacc_fwFilter_term2"
    from {
      source_address            = ["192.0.2.0/25"]
      source_address_except     = ["192.0.2.128/25"]
      port_except               = ["23"]
      source_prefix_list        = [junos_policyoptions_prefix_list.testacc_fwFilter.name]
      source_prefix_list_except = [junos_policyoptions_prefix_list.testacc_fwFilter2.name]
      tcp_established           = true
      protocol_except           = ["icmp"]
    }
    then {
      policer = junos_firewall_policer.testacc_fwfilter.name
      action  = "accept"
    }
  }
  term {
    name = "testacc_fwFilter_term3"
    from {
      destination_address            = ["192.0.2.0/25"]
      destination_address_except     = ["192.0.2.128/25"]
      destination_port               = ["22-23"]
      source_port_except             = ["23"]
      destination_prefix_list        = [junos_policyoptions_prefix_list.testacc_fwFilter.name]
      destination_prefix_list_except = [junos_policyoptions_prefix_list.testacc_fwFilter2.name]
      tcp_initial                    = true
    }
    then {
      action = "discard"
    }
  }
  term {
    name = "testacc_fwFilter_term4"
    from {
      source_port             = ["22-23"]
      destination_port_except = ["23"]
    }
    then {
      action = "reject"
    }
  }
  term {
    name = "testacc_fwFilter_term5"
    from {
      icmp_code_except = ["network-unreachable"]
      icmp_type_except = ["router-advertisement"]
    }
    then {
      action = "reject"
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter6" {
  name   = "testacc_fwFilter#6"
  family = "inet6"
  term {
    name = "testacc_fwFilter#6 term1"
    from {
      interface     = ["fe-*"]
      next_header   = ["icmp6"]
      loss_priority = ["low"]
    }
    then {
      action = "discard"
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter62" {
  name   = "testacc_fwFilte #62"
  family = "inet6"
  term {
    name   = "testacc_fwFilter#62 term1"
    filter = junos_firewall_filter.testacc_fwFilter6.name
  }
}
resource "junos_policyoptions_prefix_list" "testacc_fwFilter" {
  name   = "testacc_fwFilter#1"
  prefix = ["192.0.2.0/25"]
}
resource "junos_policyoptions_prefix_list" "testacc_fwFilter2" {
  name   = "testacc_fwFilter#2"
  prefix = ["192.0.2.128/25"]
}
resource "junos_firewall_policer" "testacc_fwfilter" {
  name = "testacc_fwfilter#1"
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit  = "50k"
  }
  then {
    discard = true
  }
}
resource "junos_firewall_filter" "testacc_fwFilter_vpls" {
  name   = "testacc_fwFilter vpls"
  family = "vpls"
  term {
    name = "testacc_fwFilter vpls term1"
    from {
      forwarding_class_except = [
        "network-control",
      ]
    }
    then {
      loss_priority = "high"
    }
  }
}
resource "junos_firewall_filter" "testacc_fwFilter_any" {
  name   = "testacc_fwFilter any"
  family = "any"
  term {
    name = "testacc_fwFilter any term1"
    from {
      packet_length_except = ["1-500"]
    }
    then {
      forwarding_class = "best-effort"
    }
  }
}
`
}
