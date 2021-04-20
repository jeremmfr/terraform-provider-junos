package junos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSnmpView_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosSnmpViewConfigCreate(),
			},
			{
				Config: testAccJunosSnmpViewConfigUpdate(),
			},
			{
				ResourceName:      "junos_snmp_view.testacc_snmpview",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosSnmpViewConfigCreate() string {
	return `
resource "junos_snmp_view" "testacc_snmpview" {
  name        = "testacc_snmpview"
  oid_include = [".1"]
}
`
}

func testAccJunosSnmpViewConfigUpdate() string {
	return `
resource "junos_snmp_view" "testacc_snmpview" {
  name        = "testacc_snmpview"
  oid_include = [".1", ".1.1"]
  oid_exclude = [".1.1.2"]
}
`
}
