package junos

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccJunosSecurityZone_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityZoneConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_protocols.#", "1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_protocols.0", "bgp"),
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
					),
				},
				{
					Config: testAccJunosSecurityZoneConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_services.#", "1"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"inbound_services.0", "ssh"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.#", "2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book.1.name", "testacc_address2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.#", "2"),
						resource.TestCheckResourceAttr("junos_security_zone.testacc_securityZone",
							"address_book_set.0.address.1", "testacc_address2"),
					),
				},
				{
					ResourceName:      "junos_security_zone.testacc_securityZone",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityZoneConfigCreate() string {
	return fmt.Sprintf(`
resource junos_security_zone "testacc_securityZone" {
  name = "testacc_securityZone"
  inbound_protocols = [ "bgp" ]
  address_book {
    name = "testacc_address1"
    network = "192.0.2.0/25"
  }
  address_book_set {
    name = "testacc_addressSet"
    address = [ "testacc_address1" ]
  }
}
`)
}
func testAccJunosSecurityZoneConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_security_zone "testacc_securityZone" {
  name = "testacc_securityZone"
  inbound_protocols = [ "bgp" ]
  inbound_services = [ "ssh"]
  address_book {
    name = "testacc_address1"
    network = "192.0.2.0/25"
  }
  address_book {
    name = "testacc_address2"
    network = "192.0.2.128/25"
  }
  address_book_set {
    name = "testacc_addressSet"
    address = [ "testacc_address1", "testacc_address2" ]
  }
}
`)
}
