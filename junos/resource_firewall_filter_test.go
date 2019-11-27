package junos

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.name", "testacc_fwFilter_term1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.address.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.address.0", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.address_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.address_except.0", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.port.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.port.0", "22-23"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.prefix_list.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.prefix_list.0", "testacc_fwFilter"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.prefix_list_except.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.prefix_list_except.0", "testacc_fwFilter2"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.protocol.#", "1"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.protocol.0", "tcp"),
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.0.from.0.tcp_flags", "!0x3"),
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
					),
				},
				{
					Config: testAccJunosFirewallFilterConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_firewall_filter.testacc_fwFilter",
							"term.#", "4"),
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
	return fmt.Sprintf(`
resource junos_firewall_filter "testacc_fwFilter" {
  name = "testacc_fwFilter"
  family = "inet"
  interface_specific = true
  term {
    name = "testacc_fwFilter_term1"
    from {
      address = [ "192.0.2.0/25" ]
      address_except = [ "192.0.2.128/25" ]
      port = [ "22-23" ]
      prefix_list = [ junos_policyoptions_prefix_list.testacc_fwFilter.name ]
      prefix_list_except = [ junos_policyoptions_prefix_list.testacc_fwFilter2.name ]
      protocol = [ "tcp" ]
      tcp_flags = "!0x3"
    }
    then {
      action = "next term"
      syslog = true
      log = true
      port_mirror = true
      service_accounting = true
    }
  }
}
resource junos_policyoptions_prefix_list "testacc_fwFilter" {
  name = "testacc_fwFilter"
  prefix = [ "192.0.2.0/25" ]
}
resource junos_policyoptions_prefix_list "testacc_fwFilter2" {
  name = "testacc_fwFilter2"
  prefix = [ "192.0.2.128/25" ]
}
`)
}
func testAccJunosFirewallFilterConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_firewall_filter "testacc_fwFilter" {
  name = "testacc_fwFilter"
  family = "inet"
  interface_specific = true
  term {
    name = "testacc_fwFilter_term1"
    from {
      address = [ "192.0.2.0/25" ]
      address_except = [ "192.0.2.128/25" ]
      port = [ "22-23" ]
      prefix_list = [ junos_policyoptions_prefix_list.testacc_fwFilter.name ]
      prefix_list_except = [ junos_policyoptions_prefix_list.testacc_fwFilter2.name ]
      protocol = [ "tcp" ]
      tcp_flags = "!0x3"
    }
    then {
      action = "next term"
      syslog = true
      log = true
      port_mirror = true
      service_accounting = true
    }
  }
  term {
    name = "testacc_fwFilter_term2"
    from {
      source_address = [ "192.0.2.0/25" ]
      source_address_except = [ "192.0.2.128/25" ]
      port_except = [ "23" ]
      source_prefix_list = [ junos_policyoptions_prefix_list.testacc_fwFilter.name ]
      source_prefix_list_except = [ junos_policyoptions_prefix_list.testacc_fwFilter2.name ]
      tcp_established = true
      protocol_except = [ "icmp" ]
    }
    then {
      policer = junos_firewall_policer.testacc_fwfilter.name
      action = "accept"
    }
  }
  term {
    name = "testacc_fwFilter_term3"
    from {
      destination_address = [ "192.0.2.0/25" ]
      destination_address_except = [ "192.0.2.128/25" ]
      destination_port = [ "22-23" ]
      source_port_except = [ "23" ]
      destination_prefix_list = [ junos_policyoptions_prefix_list.testacc_fwFilter.name ]
      destination_prefix_list_except = [ junos_policyoptions_prefix_list.testacc_fwFilter2.name ]
      tcp_initial = true
    }
    then {
      action = "discard"
    }
  }
  term {
    name = "testacc_fwFilter_term4"
    from {
      source_port = [ "22-23" ]
      destination_port_except = [ "23" ]
    }
    then {
      action = "reject"
    }
  }
}
resource junos_policyoptions_prefix_list "testacc_fwFilter" {
  name = "testacc_fwFilter"
  prefix = [ "192.0.2.0/25" ]
}
resource junos_policyoptions_prefix_list "testacc_fwFilter2" {
  name = "testacc_fwFilter2"
  prefix = [ "192.0.2.128/25" ]
}
resource junos_firewall_policer testacc_fwfilter {
  name = "testacc_fwfilter"
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit = "50k"
  }
  then {
    discard = true
  }
}
`)
}
