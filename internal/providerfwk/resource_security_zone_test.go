package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccResourceSecurityZone_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.#", "1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.0.name", "testacc_zone1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.#", "2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.name", "testacc_zoneSet"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.*", "testacc_zone1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"application_tracking", "true"),
						resource.TestCheckTypeSetElemAttr("junos_security_zone.testacc_securityZone",
							"inbound_protocols.*", "bgp"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_protocols.#", "1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"reverse_reroute", "true"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"source_identity_log", "true"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"screen", "testaccZone"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"tcp_rst", "true"),
					),
				},
				{
					ResourceName:      "junos_security_zone.testacc_securityZone",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.#", "2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.1.name", "testacc_zone2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.#", "2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.1", "testacc_zone2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_services.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_security_zone.testacc_securityZone",
							"inbound_services.*", "ssh"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_securityZone",
							"security_zone", "testacc_securityZone"),
					),
				},
			},
		})
	}
}
