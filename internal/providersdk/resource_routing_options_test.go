package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRoutingOptions_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceRoutingOptionsConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.#", "1"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.0.number", "65000"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.0.asdot_notation", "true"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.0.loops", "5"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.#", "1"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.0.restart_duration", "120"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.0.disable", "true"),
					),
				},
				{
					ResourceName:      "junos_routing_options.testacc_routing_options",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceRoutingOptionsConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.#", "1"),
					),
				},
			},
		})
	}
}

func testAccResourceRoutingOptionsConfigCreate() string {
	return `
resource "junos_policyoptions_policy_statement" "testacc_routing_options" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_routing_options"
  from {
    protocol = ["bgp"]
    route_filter {
      route  = "192.0.2.0/28"
      option = "orlonger"
    }
  }
  then {
    load_balance = "per-packet"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_routing_options2" {
  name = "testacc_routing_options2"
  from {
    route_filter {
      route  = "192.0.2.0/28"
      option = "orlonger"
    }
  }
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_routing_options3" {
  name = "testacc_routing_options3"
  from {
    route_filter {
      route  = "192.0.2.16/28"
      option = "orlonger"
    }
  }
  then {
    action = "accept"
  }
}

resource "junos_routing_options" "testacc_routing_options" {
  autonomous_system {
    number         = "65000"
    asdot_notation = true
    loops          = 5
  }
  forwarding_table {
    dynamic_list_next_hop                     = true
    ecmp_fast_reroute                         = true
    export                                    = [junos_policyoptions_policy_statement.testacc_routing_options.name]
    indirect_next_hop                         = true
    indirect_next_hop_change_acknowledgements = true
    krt_nexthop_ack_timeout                   = 200
    remnant_holdtime                          = 0
    unicast_reverse_path                      = "active-paths"
  }
  graceful_restart {
    restart_duration = 120
    disable          = true
  }
  instance_export = [junos_policyoptions_policy_statement.testacc_routing_options2.name]
  instance_import = [junos_policyoptions_policy_statement.testacc_routing_options3.name]
  router_id       = "192.0.2.4"
}
`
}

func testAccResourceRoutingOptionsConfigUpdate() string {
	return `
resource "junos_policyoptions_policy_statement" "testacc_routing_options2" {
  name = "testacc_routing_options2"
  from {
    route_filter {
      route  = "192.0.2.0/28"
      option = "orlonger"
    }
  }
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_routing_options3" {
  name = "testacc_routing_options3"
  from {
    route_filter {
      route  = "192.0.2.16/28"
      option = "orlonger"
    }
  }
  then {
    action = "accept"
  }
}

resource "junos_routing_options" "testacc_routing_options" {
  clean_on_destroy                         = true
  forwarding_table_export_configure_singly = true
  forwarding_table {
    no_ecmp_fast_reroute                         = true
    no_indirect_next_hop                         = true
    no_indirect_next_hop_change_acknowledgements = true
    unicast_reverse_path                         = "feasible-paths"
  }
  graceful_restart {}
  instance_export = [
    junos_policyoptions_policy_statement.testacc_routing_options3.name,
    junos_policyoptions_policy_statement.testacc_routing_options2.name,
  ]
  instance_import = [
    junos_policyoptions_policy_statement.testacc_routing_options2.name,
    junos_policyoptions_policy_statement.testacc_routing_options3.name,
  ]
}
resource "junos_policyoptions_policy_statement" "testacc_routing_options" {
  name                              = "testacc_routing_options"
  add_it_to_forwarding_table_export = true
  from {
    route_filter {
      route  = "192.0.2.0/25"
      option = "orlonger"
    }
  }
  then {
    load_balance = "per-packet"
  }
}
`
}
