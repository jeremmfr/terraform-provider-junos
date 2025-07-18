package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testaccNullCommitFile = "/tmp/testacc_terraform-provider-junos_null-commit-file"

// export TESTACC_INTERFACE=<inteface> to choose interface available else it's ge-0/0/3.
func TestAccResourceNullCommitFile_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				// 1
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
					"file":      config.StringVariable(testaccNullCommitFile),
				},
			},
			{
				// 2
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
					"file":      config.StringVariable(testaccNullCommitFile),
				},
			},
			{
				// 3
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
					"file":      config.StringVariable(testaccNullCommitFile),
				},
			},
			{
				// 4
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
					"file":      config.StringVariable(testaccNullCommitFile),
				},
				ExpectNonEmptyPlan: true,
			},
			{
				// 5
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
				},
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_interface_physical.testacc_nullcommitfile",
						"description", "testacc_nullfile"),
					resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_nullcommitfile",
						"description", "testacc_nullfile"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
