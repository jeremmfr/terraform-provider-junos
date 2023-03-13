package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSnmpClientlist_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosSnmpClientlistConfigCreate(),
			},
			{
				Config: testAccJunosSnmpClientlistConfigUpdate(),
			},
			{
				ResourceName:      "junos_snmp_clientlist.testacc_snmpclientlist",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosSnmpClientlistConfigCreate() string {
	return `
resource "junos_snmp_clientlist" "testacc_snmpclientlist" {
  name = "testacc@snmpclientlist"
}
`
}

func testAccJunosSnmpClientlistConfigUpdate() string {
	return `
resource "junos_snmp_clientlist" "testacc_snmpclientlist" {
  name   = "testacc@snmpclientlist"
  prefix = ["192.0.2.1/32", "192.0.2.2/32"]
}
`
}
