package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVirtualChassis_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" &&
		os.Getenv("TESTACC_VIRTUAL_CHASSIS_SN") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"serial_number": config.StringVariable(os.Getenv("TESTACC_VIRTUAL_CHASSIS_SN")),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ResourceName:      "junos_virtual_chassis.testacc_virtual_chassis",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
