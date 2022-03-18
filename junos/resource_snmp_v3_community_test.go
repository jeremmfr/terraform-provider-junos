package junos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSnmpV3Communitry_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosSnmpV3CommunitryConfigCreate(),
			},
			{
				Config: testAccJunosSnmpV3CommunitryConfigUpdate(),
			},
			{
				ResourceName:      "junos_snmp_v3_community.testacc_snmpv3comm",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosSnmpV3CommunitryConfigCreate() string {
	return `
resource "junos_snmp_v3_community" "testacc_snmpv3comm" {
  community_index = "testacc_snmpv3comm#1"
  security_name   = "testacc_snmpv3comm#1_security"
}
`
}

func testAccJunosSnmpV3CommunitryConfigUpdate() string {
	return `
resource "junos_snmp_v3_community" "testacc_snmpv3comm" {
  community_index = "testacc_snmpv3comm#1"
  security_name   = "testacc_snmpv3comm#1_security2"
  community_name  = "testacc_snmpcomm#1"
  context         = "testacc_snmpv3comm#1_ctx"
  tag             = "testacc_snmpv3comm#1_tag"
}
`
}
