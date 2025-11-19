package provider_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const testaccActionCommitFile = "/tmp/testacc_terraform-provider-junos_action-commit-file"

// export TESTACC_INTERFACE=<interface> to choose interface available else it's ge-0/0/3.
func TestAccActionCommitFile_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// tfversion.SkipBelow(tfversion.Version1_14_0),
			tfversion.SkipBelow(version.Must(version.NewVersion("1.14.0"))),
		},
		Steps: []resource.TestStep{
			{
				// 1
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
					"file":      config.StringVariable(testaccActionCommitFile),
				},
			},
			{
				// 2
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
					"file":      config.StringVariable(testaccActionCommitFile),
				},
			},
			{
				// 3
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
					"file":      config.StringVariable(testaccActionCommitFile),
				},
				ExpectNonEmptyPlan: true,
			},
			{
				// 4
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface": config.StringVariable(testaccInterface),
				},
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_interface_physical.testacc_actioncommitfile",
						"description", "testacc_actionfile"),
					resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_actioncommitfile",
						"description", "testacc_actionfile"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
