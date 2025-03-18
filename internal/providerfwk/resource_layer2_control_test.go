package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> to choose interface available else it's xe-0/0/3.
// export TESTACC_INTERFACE2=<interface> to choose 2nd interface available else it's xe-0/0/4.
func TestAccResourceLayer2Control_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface := junos.DefaultInterfaceSwitchTestAcc
		testaccInterface2 := junos.DefaultInterfaceSwitchTestAcc2
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
			testaccInterface2 = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_layer2_control.l2c",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
				},
			},
		})
	}
}
