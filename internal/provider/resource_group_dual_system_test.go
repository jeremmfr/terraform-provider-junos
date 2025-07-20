package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceGroupDualSystem_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"apply_groups", "true"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"interface_fxp0.family_inet_address.#", "2"),
						resource.TestCheckResourceAttr("junos_group_dual_system.testacc_node0",
							"system.backup_router_destination.#", "2"),
						resource.TestCheckTypeSetElemAttr("junos_group_dual_system.testacc_node0",
							"system.backup_router_destination.*", "192.0.2.0/26"),
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
