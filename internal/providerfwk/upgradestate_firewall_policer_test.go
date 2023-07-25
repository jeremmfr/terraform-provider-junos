package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccJunosFirewallPolicerUpgradeStateV0toV1_basic(t *testing.T) {
	if os.Getenv("TESTACC_UPGRADE_STATE") == "" {
		return
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					ExternalProviders: map[string]resource.ExternalProvider{
						"junos": {
							VersionConstraint: "1.33.0",
							Source:            "jeremmfr/junos",
						},
					},
					Config: testAccJunosFirewallPolicerConfigV0(),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosFirewallPolicerConfigV0(),
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

func testAccJunosFirewallPolicerConfigV0() string {
	return `
resource "junos_firewall_policer" "testacc_fwPolic" {
  name = "testacc_fwPolic"
  if_exceeding {
    bandwidth_limit  = "32k"
    burst_size_limit = "50k"
  }
  then {
    forwarding_class = "best-effort"
    loss_priority    = "high"
    out_of_profile   = true
  }
}
`
}
