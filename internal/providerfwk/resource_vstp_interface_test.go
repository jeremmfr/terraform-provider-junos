package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's xe-0/0/3.
func TestAccResourceVstpInterface_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface := junos.DefaultInterfaceSwitchTestAcc
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
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
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface3",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface4",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface5",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface6",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
