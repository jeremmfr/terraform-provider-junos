package provider_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> to choose interface available else it's ge-0/0/3.
// export TESTACC_INTERFACE2=<interface> to choose 2nd interface available else it's ge-0/0/4.
func TestAccResourceChassisCluster_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	testaccInterface2 := junos.DefaultInterfaceTestAcc2
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccInterface2 = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"fab0.member_interfaces.#", "1"),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"fab0.member_interfaces.0", testaccInterface),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"redundancy_group.#", "2"),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"redundancy_group.0.node0_priority", "100"),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"redundancy_group.1.interface_monitor.#", "1"),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"redundancy_group.1.interface_monitor.0.name", testaccInterface2),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"reth_count", "2"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"redundancy_group.#", "3"),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"redundancy_group.1.node0_priority", "100"),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"reth_count", "3"),
					),
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
					ResourceName:      "junos_chassis_cluster.testacc_cluster",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
