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
func TestAccJunosSecurityIpsecVPNUpgradeStateV0toV1_basic(t *testing.T) {
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
					Config: testAccJunosSecurityIpsecVPNConfigV0(testaccInterface),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosSecurityIpsecVPNConfigV0(testaccInterface),
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

func testAccJunosSecurityIpsecVPNConfigV0(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_v0to1_ipsecvpn" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_policy" "testacc_v0to1_ipsecvpn" {
  name                = "testacc_v0to1_ipsecvpn"
  proposal_set        = "basic"
  mode                = "main"
  pre_shared_key_text = "thePassWord"
}
resource "junos_security_ike_gateway" "testacc_v0to1_ipsecvpn" {
  name               = "testacc_v0to1_ipsecvpn"
  address            = ["192.0.2.3"]
  policy             = junos_security_ike_policy.testacc_v0to1_ipsecvpn.name
  external_interface = junos_interface_logical.testacc_v0to1_ipsecvpn.name
}
resource "junos_security_ipsec_policy" "testacc_v0to1_ipsecvpn" {
  name         = "testacc_ipsecpol"
  proposal_set = "basic"
  pfs_keys     = "group2"
}
resource "junos_interface_st0_unit" "testacc_v0to1_ipsecvpn" {}
resource "junos_security_ipsec_vpn" "testacc_v0to1_ipsecvpn" {
  name           = "testacc_v0to1_ipsecvpn"
  bind_interface = junos_interface_st0_unit.testacc_v0to1_ipsecvpn.id
  ike {
    gateway          = junos_security_ike_gateway.testacc_v0to1_ipsecvpn.name
    policy           = junos_security_ipsec_policy.testacc_v0to1_ipsecvpn.name
    identity_local   = "192.0.2.64/26"
    identity_remote  = "192.0.2.128/26"
    identity_service = "any"
  }
  vpn_monitor {
    destination_ip = "192.0.2.129"
    optimized      = true
  }
  establish_tunnels = "on-traffic"
}
`, interFace)
}
