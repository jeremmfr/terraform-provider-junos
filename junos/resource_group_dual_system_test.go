package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosGroupDualSystem_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosGroupDualSystemConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"apply_groups", "true"),
					),
				},
				{
					Config: testAccJunosGroupDualSystemConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"interface_fxp0.#", "1"),
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"interface_fxp0.0.family_inet_address.#", "2"),
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"system.#", "1"),
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"system.0.backup_router_destination.#", "2"),
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"system.#", "1"),
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"system.0.backup_router_destination.1", "192.0.2.0/26"),
					),
				},
				{
					ResourceName:      "junos_group_dual_system.testacc_node0",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosGroupDualSystemConfigCreate() string {
	return `
resource "junos_group_dual_system" "testacc_node0" {
  name = "node0"
  interface_fxp0 {
    description = "test_"
    family_inet_address {
      cidr_ip = "192.0.2.193/26"
    }
  }
  routing_options {
    static_route {
      destination = "192.0.2.0/26"
      next_hop    = ["192.0.2.254"]
    }
    static_route {
      destination = "192.0.2.64/26"
      next_hop    = ["192.0.2.254"]
    }
  }
  security {
    log_source_address = "192.0.2.128"
  }
  system {
    host_name             = "test_node"
    backup_router_address = "192.0.2.254"
    backup_router_destination = [
      "192.0.2.0/26",
    ]
  }
}
`
}
func testAccJunosGroupDualSystemConfigUpdate() string {
	return `
resource "junos_group_dual_system" "testacc_node0" {
  name         = "node0"
  apply_groups = false
  interface_fxp0 {
    description = "test_"
    family_inet_address {
      cidr_ip = "192.0.2.193/26"
    }
    family_inet_address {
      cidr_ip     = "192.0.2.194/26"
      master_only = true
    }
  }
  routing_options {
    static_route {
      destination = "192.0.2.0/26"
      next_hop    = ["192.0.2.254"]
    }
  }
  security {
    log_source_address = "192.0.2.128"
  }
  system {
    host_name             = "test_node"
    backup_router_address = "192.0.2.254"
    backup_router_destination = [
      "192.0.2.64/26",
      "192.0.2.0/26",

    ]
  }
}
`
}
