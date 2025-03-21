package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccUpgradeStateResourceLayer2Control_V0toV1_basic(t *testing.T) {
	if os.Getenv("TESTACC_UPGRADE_STATE") == "" {
		return
	}
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface := junos.DefaultInterfaceTestAcc
		testaccInterface2 := junos.DefaultInterfaceTestAcc2
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
			testaccInterface2 = iface
		}
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	}
}
