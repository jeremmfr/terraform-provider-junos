package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosRoutingInstance_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosRoutingInstanceConfigSRXCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"type", "virtual-router"),
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65000"),
					),
				},
				{
					Config: testAccJunosRoutingInstanceConfigSRXUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65001"),
					),
				},
				{
					ResourceName:      "junos_routing_instance.testacc_routingInst",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosRoutingInstanceConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"type", "virtual-router"),
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65000"),
					),
				},
				{
					Config: testAccJunosRoutingInstanceConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65001"),
					),
				},
				{
					ResourceName:      "junos_routing_instance.testacc_routingInst",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_routing_instance.testacc_routingInst2",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosRoutingInstanceConfigSRXCreate() string {
	return `
resource "junos_routing_instance" "testacc_routingInst" {
  name            = "testacc_routingInst"
  as              = "65000"
  description     = "testacc routingInst"
  instance_export = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  instance_import = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  router_id       = "192.0.2.65"
}
resource "junos_policyoptions_community" "testacc_routingInst2" {
  name    = "testacc_routingInst2"
  members = ["target:65000:100"]
}
resource "junos_policyoptions_policy_statement" "testacc_routingInst2" {
  name = "testacc_routingInst2"
  from {
    bgp_community = [junos_policyoptions_community.testacc_routingInst2.name]
  }
  then {
    action = "accept"
  }
}
`
}

func testAccJunosRoutingInstanceConfigSRXUpdate() string {
	return `
resource "junos_routing_instance" "testacc_routingInst" {
  name = "testacc_routingInst"
  as   = "65001"
  instance_export = [
    junos_policyoptions_policy_statement.testacc_routingInst3.name,
    junos_policyoptions_policy_statement.testacc_routingInst2.name,
  ]
  instance_import = [
    junos_policyoptions_policy_statement.testacc_routingInst3.name,
    junos_policyoptions_policy_statement.testacc_routingInst2.name,
  ]
}
resource "junos_policyoptions_community" "testacc_routingInst2" {
  name    = "testacc_routingInst2"
  members = ["target:65000:100"]
}
resource "junos_policyoptions_community" "testacc_routingInst3" {
  name    = "testacc_routingInst3"
  members = ["target:65000:200"]
}
resource "junos_policyoptions_policy_statement" "testacc_routingInst2" {
  name = "testacc_routingInst2"
  from {
    bgp_community = [junos_policyoptions_community.testacc_routingInst2.name]
  }
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_routingInst3" {
  name = "testacc_routingInst3"
  from {
    bgp_community = [junos_policyoptions_community.testacc_routingInst3.name]
  }
  then {
    action = "accept"
  }
}
`
}

func testAccJunosRoutingInstanceConfigCreate() string {
	return `
resource "junos_routing_instance" "testacc_routingInst" {
  name            = "testacc_routingInst"
  as              = "65000"
  description     = "testacc routingInst"
  instance_export = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  instance_import = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  router_id       = "192.0.2.65"
}
resource "junos_policyoptions_community" "testacc_routingInst2" {
  name    = "testacc_routingInst2"
  members = ["target:65000:100"]
}
resource "junos_policyoptions_policy_statement" "testacc_routingInst2" {
  name = "testacc_routingInst2"
  from {
    bgp_community = [junos_policyoptions_community.testacc_routingInst2.name]
  }
  then {
    action = "accept"
  }
}
resource "junos_interface_logical" "testacc_routingInst2" {
  name = "lo0.1"
  family_inet {
    address {
      cidr_ip = "192.0.2.15/32"
    }
  }
}
resource "junos_routing_instance" "testacc_routingInst2" {
  name                  = "testacc_routingInst2"
  type                  = "l2vpn"
  route_distinguisher   = "1:2"
  vrf_export            = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  vrf_import            = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  vrf_target            = "target:2:3"
  vrf_target_export     = "target:4:5"
  vrf_target_import     = "target:6:7"
  vtep_source_interface = junos_interface_logical.testacc_routingInst2.name
}

resource "junos_routing_instance" "testacc_routingInst3" {
  name                  = "testacc_routingInst3"
  configure_type_singly = true
  type                  = ""
  vtep_source_interface = junos_interface_logical.testacc_routingInst2.name
}
resource "junos_routing_instance" "testacc_routingInst4" {
  name                = "testacc_routingInst4"
  type                = "virtual-switch"
  route_distinguisher = "8:9"
  vrf_export          = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  vrf_target_auto     = true
}
`
}

func testAccJunosRoutingInstanceConfigUpdate() string {
	return `
resource "junos_routing_instance" "testacc_routingInst" {
  name = "testacc_routingInst"
  as   = "65001"
  instance_export = [
    junos_policyoptions_policy_statement.testacc_routingInst3.name,
    junos_policyoptions_policy_statement.testacc_routingInst2.name,
  ]
  instance_import = [
    junos_policyoptions_policy_statement.testacc_routingInst3.name,
    junos_policyoptions_policy_statement.testacc_routingInst2.name,
  ]
}
resource "junos_policyoptions_community" "testacc_routingInst2" {
  name    = "testacc_routingInst2"
  members = ["target:65000:100"]
}
resource "junos_policyoptions_community" "testacc_routingInst3" {
  name    = "testacc_routingInst3"
  members = ["target:65000:200"]
}
resource "junos_policyoptions_policy_statement" "testacc_routingInst2" {
  name = "testacc_routingInst2"
  from {
    bgp_community = [junos_policyoptions_community.testacc_routingInst2.name]
  }
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_routingInst3" {
  name = "testacc_routingInst3"
  from {
    bgp_community = [junos_policyoptions_community.testacc_routingInst3.name]
  }
  then {
    action = "accept"
  }
}
resource "junos_interface_logical" "testacc_routingInst2" {
  name = "lo0.1"
  family_inet {
    address {
      cidr_ip = "192.0.2.15/32"
    }
  }
}
resource "junos_routing_instance" "testacc_routingInst2" {
  name                = "testacc_routingInst2"
  type                = "l2vpn"
  route_distinguisher = "1:2"
  vrf_export = [
    junos_policyoptions_policy_statement.testacc_routingInst3.name,
    junos_policyoptions_policy_statement.testacc_routingInst2.name,
  ]
  vrf_import = [
    junos_policyoptions_policy_statement.testacc_routingInst3.name,
    junos_policyoptions_policy_statement.testacc_routingInst2.name,
  ]
  vrf_target            = "target:2:3"
  vrf_target_export     = "target:4:5"
  vrf_target_import     = "target:8:7"
  vtep_source_interface = junos_interface_logical.testacc_routingInst2.name
}

resource "junos_routing_instance" "testacc_routingInst3" {
  name                  = "testacc_routingInst3"
  configure_type_singly = true
  type                  = ""
  vtep_source_interface = junos_interface_logical.testacc_routingInst2.name
}
resource "junos_routing_instance" "testacc_routingInst4" {
  name                = "testacc_routingInst4"
  type                = "virtual-switch"
  route_distinguisher = "10:11"
  vrf_export          = [junos_policyoptions_policy_statement.testacc_routingInst2.name]
  vrf_target_auto     = true
}
`
}
