package providerfwk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccJunosInterfacePhysicalUpgradeStateV0toV1_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	testaccInterface2 := junos.DefaultInterfaceTestAcc2
	testaccInterfaceAE := "ae0"
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface = junos.DefaultInterfaceSwitchTestAcc
	}
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE_AE"); iface != "" {
		testaccInterfaceAE = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccInterface2 = iface
	}
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					ExternalProviders: map[string]resource.ExternalProvider{
						"junos": {
							VersionConstraint: "1.33.0",
							Source:            "jeremmfr/junos",
						},
					},
					Config: testAccJunosInterfacePhysicalSWConfigV0(testaccInterface),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosInterfacePhysicalSWConfigV0(testaccInterface),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	} else {
		if os.Getenv("TESTACC_ROUTER") != "" {
			resource.Test(t, resource.TestCase{
				Steps: []resource.TestStep{
					{
						ExternalProviders: map[string]resource.ExternalProvider{
							"junos": {
								VersionConstraint: "1.33.0",
								Source:            "jeremmfr/junos",
							},
						},
						Config: testAccJunosInterfacePhysicalRouterConfigV0(testaccInterface, testaccInterfaceAE),
					},
					{
						ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
						Config:                   testAccJunosInterfacePhysicalRouterConfigV0(testaccInterface, testaccInterfaceAE),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectEmptyPlan(),
							},
						},
					},
				},
			})
		}
		if os.Getenv("TESTACC_SRX") != "" {
			resource.Test(t, resource.TestCase{
				Steps: []resource.TestStep{
					{
						ExternalProviders: map[string]resource.ExternalProvider{
							"junos": {
								VersionConstraint: "1.33.0",
								Source:            "jeremmfr/junos",
							},
						},
						Config: testAccJunosInterfacePhysicalSRXConfigV0(testaccInterface, testaccInterface2),
					},
					{
						ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
						Config:                   testAccJunosInterfacePhysicalSRXConfigV0(testaccInterface, testaccInterface2),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectEmptyPlan(),
							},
						},
					},
				},
			})
		}
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					ExternalProviders: map[string]resource.ExternalProvider{
						"junos": {
							VersionConstraint: "1.33.0",
							Source:            "jeremmfr/junos",
						},
					},
					Config: testAccJunosInterfacePhysicalConfigV0(testaccInterface, testaccInterfaceAE, testaccInterface2),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config: testAccJunosInterfacePhysicalConfigV0(
						testaccInterface, testaccInterfaceAE, testaccInterface2),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	}
}

func testAccJunosInterfacePhysicalSWConfigV0(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface" {
  name         = "%s"
  description  = "testacc_interface"
  trunk        = true
  vlan_native  = 100
  vlan_members = ["100-110"]
}
`, interFace)
}

func testAccJunosInterfacePhysicalConfigV0(interFace, interfaceAE, interFace2 string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface" {
  name        = "%s"
  description = "testacc_interfaceU"
  gigether_opts {
    ae_8023ad = "%s"
  }
}
resource "junos_interface_physical" "testacc_interface2" {
  name        = "%s"
  description = "testacc_interface2"
  ether_opts {
    flow_control     = true
    loopback         = true
    auto_negotiation = true
  }
  mtu = 9000
}
resource "junos_interface_logical" "testacc_interfaceLO" {
  name = "lo0.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/32"
    }
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  depends_on = [
    junos_interface_physical.testacc_interface,
    junos_interface_logical.testacc_interfaceLO,
  ]
  name        = "%s"
  description = "testacc_interfaceAE"
  parent_ether_opts {
    bfd_liveness_detection {
      local_address                      = "192.0.2.1"
      detection_time_threshold           = 30
      holddown_interval                  = 30
      minimum_interval                   = 30
      minimum_receive_interval           = 10
      multiplier                         = 1
      neighbor                           = "192.0.2.2"
      no_adaptation                      = true
      transmit_interval_minimum_interval = 10
      transmit_interval_threshold        = 30
      version                            = "automatic"
    }
    no_flow_control   = true
    no_loopback       = true
    link_speed        = "1g"
    minimum_bandwidth = "1 gbps"
  }
  vlan_tagging = true
}
`, interFace, interfaceAE, interFace2, interfaceAE)
}

func testAccJunosInterfacePhysicalRouterConfigV0(interFace, interfaceAE string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface" {
  name        = "%s"
  description = "testacc_interface"
  gigether_opts {
    ae_8023ad = "%s"
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  name        = "%s"
  description = "testacc_interfaceAE"
  esi {
    identifier = "00:11:11:11:11:11:11:11:11:11"
    mode       = "all-active"
  }
  vlan_tagging = true
}
`, interFace, interfaceAE, interfaceAE)
}

func testAccJunosInterfacePhysicalSRXConfigV0(interFace, interFace2 string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface" {
  depends_on = [
    junos_chassis_cluster.testacc_interface
  ]
  name        = "%s"
  description = "testacc_interface"
  gigether_opts {
    redundant_parent = "reth0"
  }
}
resource "junos_interface_physical" "testacc_interface2" {
  name = "%s"
}
resource "junos_chassis_cluster" "testacc_interface" {
  fab0 {
    member_interfaces = [junos_interface_physical.testacc_interface2.name]
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  reth_count = 1
}
resource "junos_interface_physical" "testacc_interface_reth" {
  depends_on = [
    junos_interface_physical.testacc_interface
  ]
  name        = "reth0"
  description = "testacc_interface_reth"
  parent_ether_opts {
    redundancy_group = 1
    minimum_links    = 1
  }
}
`, interFace, interFace2)
}
