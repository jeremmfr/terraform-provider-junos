package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccJunosSecurityNatSourceUpgradeStateV0toV1_basic(t *testing.T) {
	if os.Getenv("TESTACC_UPGRADE_STATE") == "" {
		return
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					ExternalProviders: map[string]resource.ExternalProvider{
						"junos": {
							VersionConstraint: "1.33.0",
							Source:            "jeremmfr/junos",
						},
					},
					Config: testAccJunosSecurityNatSourceConfigV0(),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosSecurityNatSourceConfigV0(),
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

func testAccJunosSecurityNatSourceConfigV0() string {
	return `
resource "junos_security_nat_source" "testacc_securitySNAT" {
  name        = "testacc_securitySNAT_upgrade"
  description = "testacc securitySNAT upgrade"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT_upgrade.name]
  }
  to {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT_upgrade.name]
  }
  rule {
    name = "testacc_securitySNATRule"
    match {
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      protocol            = ["tcp"]
    }
    then {
      type = "off"
    }
  }
}
resource "junos_security_zone" "testacc_securitySNAT_upgrade" {
  name = "testacc_securitySNAT_upgrade"
}
`
}
