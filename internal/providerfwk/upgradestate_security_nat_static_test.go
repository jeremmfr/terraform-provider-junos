package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccJunosSecurityNatStaticUpgradeStateV0toV1_basic(t *testing.T) {
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
					Config: testAccJunosSecurityNatStaticConfigV0(),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosSecurityNatStaticConfigV0(),
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

func testAccJunosSecurityNatStaticConfigV0() string {
	return `
resource "junos_security_nat_static" "testacc_securityNATStt" {
  name        = "testacc_secNATStt_upgrade"
  description = "testacc securityNATStt upgrade"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATStt.name]
  }
  rule {
    name                = "testacc_securityNATSttRule"
    destination_address = "192.0.2.0/25"
    then {
      type             = "prefix"
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
      prefix           = "192.0.2.128/25"
    }
  }
  rule {
    name                = "testacc_securityNATSttRule2"
    destination_address = "64:ff9b::/96"
    then {
      type = "inet"
    }
  }
}

resource "junos_security_zone" "testacc_securityNATStt" {
  name = "testacc_securityNATStt_upgrade"
}
resource "junos_routing_instance" "testacc_securityNATStt" {
  name = "testacc_securityNATStt_upgrade"
}
`
}
