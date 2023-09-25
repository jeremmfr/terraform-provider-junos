package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
// export TESTACC_INTERFACE2=<interface> for choose 2nd interface available else it's ge-0/0/4.
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
					Config: testAccResourceChassisClusterConfigCreate(testaccInterface, testaccInterface2),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"fab0.#", "1"),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"fab0.0.member_interfaces.#", "1"),
						resource.TestCheckResourceAttr("junos_chassis_cluster.testacc_cluster",
							"fab0.0.member_interfaces.0", testaccInterface),
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
					Config: testAccResourceChassisClusterConfigUpdate(testaccInterface, testaccInterface2),
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
					ResourceName:      "junos_chassis_cluster.testacc_cluster",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceChassisClusterConfigCreate(interFace, interFace2 string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_cluster_int2" {
  name        = "%s"
  description = "testacc_cluster_int2"
  gigether_opts {
    redundant_parent = "reth0"
  }
}
resource "junos_interface_physical" "testacc_cluster" {
  name = "%s"
}
resource "junos_chassis_cluster" "testacc_cluster" {
  fab0 {
    member_interfaces = [junos_interface_physical.testacc_cluster.name]
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  redundancy_group {
    node0_priority = 98
    node1_priority = 97
    interface_monitor {
      name   = junos_interface_physical.testacc_cluster_int2.name
      weight = 255
    }
    preempt = true
  }
  reth_count = 2
}
`, interFace2, interFace)
}

func testAccResourceChassisClusterConfigUpdate(interFace, interFace2 string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_cluster_int2" {
  name        = "%s"
  description = "testacc_cluster_int2"
  gigether_opts {
    redundant_parent = "reth0"
  }
}
resource "junos_interface_physical" "testacc_cluster" {
  name = "%s"
}
resource "junos_chassis_cluster" "testacc_cluster" {
  fab0 {
    member_interfaces = [junos_interface_physical.testacc_cluster.name]
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
    interface_monitor {
      name   = junos_interface_physical.testacc_cluster_int2.name
      weight = 255
    }
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
    preempt        = true
    preempt_delay  = 2
    preempt_limit  = 3
    preempt_period = 4
  }
  reth_count = 3
}
`, interFace2, interFace)
}
