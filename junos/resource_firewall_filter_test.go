package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosFirewallFilter_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
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
							"term.0.from.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.address.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.address.*", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.address_except.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.address_except.*", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.port.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.port.*", "22-23"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.prefix_list.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.prefix_list.*", "testacc_fwFilter"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.prefix_list_except.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.prefix_list_except.*", "testacc_fwFilter2"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.protocol.*", "tcp"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.tcp_flags", "!0x3"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.is_fragment", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.0.action", "next term"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.0.syslog", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.0.log", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.0.port_mirror", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.then.0.service_accounting", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.icmp_code.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.icmp_code.*", "network-unreachable"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.icmp_type.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.icmp_type.*", "router-advertisement"),
					),
				},
				{
					Config: testAccJunosFirewallFilterConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.#", "5"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.source_address.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.source_address_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.port_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.source_prefix_list.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.source_prefix_list_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.tcp_established", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.from.0.protocol_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.then.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.then.0.policer", "testacc_fwfilter"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.1.then.0.action", "accept"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.0.destination_address.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.0.destination_address_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.0.destination_port.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.0.source_port_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.0.destination_prefix_list.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.0.destination_prefix_list_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.from.0.tcp_initial", "true"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.then.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.2.then.0.action", "discard"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.3.from.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.3.from.0.source_port.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.3.from.0.destination_port_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.3.then.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.3.then.0.action", "reject"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.4.from.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.4.from.0.icmp_code_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.4.from.0.icmp_type_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"family", "inet6"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.0.from.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.0.from.0.next_header.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.0.from.0.next_header.*", "icmp6"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter6",
							"term.0.then.0.action", "discard"),
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
resource junos_firewall_filter "testacc_fwFilter" {
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
resource junos_policyoptions_prefix_list "testacc_fwFilter" {
  name   = "testacc_fwFilter"
  prefix = ["192.0.2.0/25"]
}
resource junos_policyoptions_prefix_list "testacc_fwFilter2" {
  name   = "testacc_fwFilter2"
  prefix = ["192.0.2.128/25"]
}
`
}

func testAccJunosFirewallFilterConfigUpdate() string {
	return `
resource junos_firewall_filter "testacc_fwFilter" {
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
resource junos_firewall_filter "testacc_fwFilter6" {
  name   = "testacc_fwFilter6"
  family = "inet6"
  term {
    name = "testacc_fwFilter6_term1"
    from {
      next_header = ["icmp6"]
    }
    then {
      action = "discard"
    }
  }
}
resource junos_policyoptions_prefix_list "testacc_fwFilter" {
  name   = "testacc_fwFilter"
  prefix = ["192.0.2.0/25"]
}
resource junos_policyoptions_prefix_list "testacc_fwFilter2" {
  name   = "testacc_fwFilter2"
  prefix = ["192.0.2.128/25"]
}
resource junos_firewall_policer testacc_fwfilter {
  name = "testacc_fwfilter"
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit  = "50k"
  }
  then {
    discard = true
  }
}
`
}
