package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccJunosSecurityNatDestinationUpgradeStateV0toV1_basic(t *testing.T) {
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
					Config: testAccJunosSecurityNatDestinationConfigV0(),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosSecurityNatDestinationConfigV0(),
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

func testAccJunosSecurityNatDestinationConfigV0() string {
	return `
resource "junos_security_nat_destination" "testacc_securityDNAT" {
  name        = "testacc_securityDNAT_upgrade"
  description = "testacc securityDNAT upgrade"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityDNAT_upgrade.name]
  }
  rule {
    name                = "testacc_securityDNATRule"
    destination_address = "192.0.2.1/32"
    then {
      type = "off"
    }
  }
}
resource "junos_security_zone" "testacc_securityDNAT_upgrade" {
  name = "testacc_securityDNAT_upgrade"
}
`
}
