package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccJunosSecurityNatStaticRuleUpgradeStateV0toV1_basic(t *testing.T) {
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
					Config: testAccJunosSecurityNatStaticRuleConfigV0(),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosSecurityNatStaticRuleConfigV0(),
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

func testAccJunosSecurityNatStaticRuleConfigV0() string {
	return `
resource "junos_security_nat_static" "testacc_securityNATSttRule" {
  name = "testacc_secNATSttRule_upgrade"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATSttRule.name]
  }
  configure_rules_singly = true
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule" {
  name                = "testacc_secNATSttRule_upgrade"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "192.0.2.0/25"
  then {
    type   = "prefix"
    prefix = "192.0.2.128/25"
  }
}
resource "junos_security_zone" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule_upgrade"
}
`
}
