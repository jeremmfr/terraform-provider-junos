package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSnmpView_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSnmpViewConfigCreate(),
			},
			{
				Config: testAccResourceSnmpViewConfigUpdate(),
			},
			{
				ResourceName:      "junos_snmp_view.testacc_snmpview",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceSnmpViewConfigCreate() string {
	return `
resource "junos_snmp_view" "testacc_snmpview" {
  name        = "testacc_snmpview"
  oid_include = [".1"]
}
`
}

func testAccResourceSnmpViewConfigUpdate() string {
	return `
resource "junos_snmp_view" "testacc_snmpview" {
  name        = "testacc_snmpview"
  oid_include = [".1", ".1.1"]
  oid_exclude = [".1.1.2"]
}
`
}
