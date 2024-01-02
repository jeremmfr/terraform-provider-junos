package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceChassisInventory_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" || os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_chassis_inventory.testacc",
							"chassis.#", "1"),
						resource.TestCheckResourceAttrSet("data.junos_chassis_inventory.testacc",
							"chassis.0.serial_number"),
						resource.TestCheckResourceAttrSet("data.junos_chassis_inventory.testacc",
							"chassis.0.module.#"),
						resource.TestCheckResourceAttrSet("data.junos_chassis_inventory.testacc",
							"chassis.0.module.0.name"),
					),
				},
			},
		})
	}
}
