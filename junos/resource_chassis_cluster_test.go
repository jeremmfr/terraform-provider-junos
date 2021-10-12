package junos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
// export TESTACC_INTERFACE2=<interface> for choose 2nd interface available else it's ge-0/0/4.
func TestAccJunosCluster_basic(t *testing.T) {
	var testaccInterface string
	var testaccInterface2 string
	if os.Getenv("TESTACC_INTERFACE") != "" {
		testaccInterface = os.Getenv("TESTACC_INTERFACE")
	} else {
		testaccInterface = defaultInterfaceTestAcc
	}
	if os.Getenv("TESTACC_INTERFACE2") != "" {
		testaccInterface2 = os.Getenv("TESTACC_INTERFACE2")
	} else {
		testaccInterface2 = defaultInterfaceTestAcc2
	}
	if os.Getenv("TESTACC_SWITCH") == "" && os.Getenv("TESTACC_ROUTER") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosClusterConfigCreate(testaccInterface, testaccInterface2),
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
					Config: testAccJunosClusterConfigUpdate(testaccInterface, testaccInterface2),
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

func testAccJunosClusterConfigCreate(interFace, interFace2 string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_cluster_int2" {
  name        = "` + interFace2 + `"
  description = "testacc_cluster_int2"
  gigether_opts {
    redundant_parent = "reth0"
  }
}
resource "junos_chassis_cluster" "testacc_cluster" {
  fab0 {
    member_interfaces = ["` + interFace + `"]
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
`)
}

func testAccJunosClusterConfigUpdate(interFace, interFace2 string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_cluster_int2" {
  name        = "` + interFace2 + `"
  description = "testacc_cluster_int2"
  gigether_opts {
    redundant_parent = "reth0"
  }
}
resource "junos_chassis_cluster" "testacc_cluster" {
  fab0 {
    member_interfaces = ["` + interFace + `"]
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
`)
}
