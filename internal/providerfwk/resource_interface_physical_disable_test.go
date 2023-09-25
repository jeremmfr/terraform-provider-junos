package providerfwk_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccResourceInterfacePhysicalDisable_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface = junos.DefaultInterfaceSwitchTestAcc
	}
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
				Destroy: true,
			},
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
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
				},
				ExpectError: regexp.MustCompile("interface \"" + testaccInterface + "\" is configured"),
			},
		},
	})
}
