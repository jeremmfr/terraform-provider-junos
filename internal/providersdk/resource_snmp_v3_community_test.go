package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSnmpV3Communitry_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSnmpV3CommunitryConfigCreate(),
			},
			{
				Config: testAccResourceSnmpV3CommunitryConfigUpdate(),
			},
			{
				ResourceName:      "junos_snmp_v3_community.testacc_snmpv3comm",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceSnmpV3CommunitryConfigCreate() string {
	return `
resource "junos_snmp_v3_community" "testacc_snmpv3comm" {
  community_index = "testacc_snmpv3comm#1"
  security_name   = "testacc_snmpv3comm#1_security"
}
`
}

func testAccResourceSnmpV3CommunitryConfigUpdate() string {
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
