package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityZone_basic(t *testing.T) {
	var testaccInterface string
	if os.Getenv("TESTACC_INTERFACE") != "" {
		testaccInterface = os.Getenv("TESTACC_INTERFACE")
	} else {
		testaccInterface = defaultInterfaceTestAcc
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityZoneConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.#", "1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.0.name", "testacc_address1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.#", "1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.name", "testacc_addressSet"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.0", "testacc_address1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"application_tracking", "true"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_protocols.0", "bgp"),
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
					Config: testAccJunosSecurityZoneConfigUpdate(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.#", "2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.1.name", "testacc_address2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.#", "2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.1", "testacc_address2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_services.#", "1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_services.0", "ssh"),
					),
				},
				{
					ResourceName:      "junos_security_zone.testacc_securityZone",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSecurityZoneConfigUpdate2(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_securityZone",
							"security_zone", "testacc_securityZone"),
					),
				},
			},
		})
	}
}

func testAccJunosSecurityZoneConfigCreate() string {
	return `
resource junos_security_screen "testaccZone" {
  name        = "testaccZone"
  description = "testaccZone"
}
resource junos_security_zone "testacc_securityZone" {
  name = "testacc_securityZone"
  address_book {
    name    = "testacc_address1"
    network = "192.0.2.0/25"
  }
  address_book_set {
    name    = "testacc_addressSet"
    address = ["testacc_address1"]
  }
  application_tracking = true
  inbound_protocols    = ["bgp"]
  description          = "testacc securityZone"
  reverse_reroute      = true
  screen               = junos_security_screen.testaccZone.id
  source_identity_log  = true
  tcp_rst              = true
}
`
}
func testAccJunosSecurityZoneConfigUpdate(interFace string) string {
	return `
resource junos_security_screen "testaccZone" {
  name        = "testaccZone"
  description = "testaccZone"
}
resource junos_security_zone "testacc_securityZone" {
  name = "testacc_securityZone"
  address_book {
    name    = "testacc_address1"
    network = "192.0.2.0/25"
  }
  address_book {
    name    = "testacc_address2"
    network = "192.0.2.128/25"
  }
  address_book_set {
    name    = "testacc_addressSet"
    address = ["testacc_address1", "testacc_address2"]
  }
  inbound_protocols = ["bgp"]
  inbound_services  = ["ssh"]
}
resource junos_interface_logical "testacc_securityZone" {
  name          = "` + interFace + `.0"
  security_zone = junos_security_zone.testacc_securityZone.name
}
`
}
func testAccJunosSecurityZoneConfigUpdate2(interFace string) string {
	return `
resource junos_security_zone "testacc_securityZone" {
  name = "testacc_securityZone"
}
resource junos_interface_logical "testacc_securityZone" {
  name          = "` + interFace + `.0"
  security_zone = junos_security_zone.testacc_securityZone.name
}
`
}
