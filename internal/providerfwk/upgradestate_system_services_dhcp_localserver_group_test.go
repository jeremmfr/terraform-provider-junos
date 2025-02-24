package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccUpgradeStateResourceSystemServicesDhcpLocalserverGroup_V0toV1_basic(t *testing.T) {
	if os.Getenv("TESTACC_UPGRADE_STATE") == "" {
		return
	}
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
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
