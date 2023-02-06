package providersdk_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceSecurityZone_basic(t *testing.T) {
	testaccInterface := defaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceSecurityZoneConfigCreate(testaccInterface),
				},
				{
					Config: testAccDataSourceSecurityZoneConfigData(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_security_zone.testacc_dataSecurityZone",
							"id", "testacc_dataSecurityZone"),
						resource.TestCheckResourceAttr("data.junos_security_zone.testacc_dataSecurityZone",
							"address_book.#", "1"),
						resource.TestCheckResourceAttr("data.junos_security_zone.testacc_dataSecurityZone",
							"interface.#", "1"),
						resource.TestCheckResourceAttr("data.junos_security_zone.testacc_dataSecurityZone",
							"interface.0.inbound_services.#", "1"),
					),
				},
				{
					Config:      testAccDataSourceSecurityZoneConfigDataFailed(),
					ExpectError: regexp.MustCompile("routing instance .* doesn't exist"),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}

func testAccDataSourceSecurityZoneConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_dataSecurityZone" {
  name         = "%s"
  description  = "testacc_dataSecurityZone"
  vlan_tagging = true
}
resource "junos_security_zone" "testacc_dataSecurityZone" {
  name                          = "testacc_dataSecurityZone"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_dataSecurityZone" {
  name = "testacc_dataSecurityZone"
  zone = junos_security_zone.testacc_dataSecurityZone.name
  cidr = "192.0.2.0/25"
}
resource "junos_interface_logical" "testacc_dataSecurityZone" {
  name                      = "${junos_interface_physical.testacc_dataSecurityZone.name}.100"
  description               = "testacc_dataSecurityZone"
  security_zone             = junos_security_zone.testacc_dataSecurityZone.name
  security_inbound_services = ["ssh"]
}
`, interFace)
}

func testAccDataSourceSecurityZoneConfigData(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_dataSecurityZone" {
  name         = "%s"
  description  = "testacc_dataSecurityZone"
  vlan_tagging = true
}
resource "junos_security_zone" "testacc_dataSecurityZone" {
  name                          = "testacc_dataSecurityZone"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_dataSecurityZone" {
  name = "testacc_dataSecurityZone"
  zone = junos_security_zone.testacc_dataSecurityZone.name
  cidr = "192.0.2.0/25"
}
resource "junos_interface_logical" "testacc_dataSecurityZone" {
  name                      = "${junos_interface_physical.testacc_dataSecurityZone.name}.100"
  description               = "testacc_dataSecurityZone"
  security_zone             = junos_security_zone.testacc_dataSecurityZone.name
  security_inbound_services = ["ssh"]
}

data "junos_security_zone" "testacc_dataSecurityZone" {
  name = "testacc_dataSecurityZone"
}
`, interFace)
}

func testAccDataSourceSecurityZoneConfigDataFailed() string {
	return `
data "junos_routing_instance" "testacc_dataSecurityZone" {
  name = "testacc"
}
`
}
