package providerfwk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccJunosSecurityIkeGatewayUpgradeStateV0toV1_basic(t *testing.T) {
	if os.Getenv("TESTACC_UPGRADE_STATE") == "" {
		return
	}
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
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
					Config: testAccJunosSecurityIkeGatewayConfigV0(testaccInterface),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosSecurityIkeGatewayConfigV0(testaccInterface),
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

func testAccJunosSecurityIkeGatewayConfigV0(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_v0to1_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_policy" "testacc_v0to1_ikegateway" {
  name                = "testacc_v0to1_ikegateway"
  proposal_set        = "basic"
  mode                = "aggressive"
  pre_shared_key_text = "thePassWord"
}
resource "junos_security_ike_gateway" "testacc_v0to1_ikegateway" {
  name               = "testacc_v0to1_ikegateway"
  policy             = junos_security_ike_policy.testacc_v0to1_ikegateway.name
  external_interface = junos_interface_logical.testacc_v0to1_ikegateway.name
  dynamic_remote {
    distinguished_name {
      container = "dc=example,dc=com"
    }
    connections_limit = 10
  }
  aaa {
    client_username = "user"
    client_password = "password"
  }
}
`, interFace)
}
