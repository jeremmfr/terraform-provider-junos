package junos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccJunosSecurityZone_basic(t *testing.T) {
	testaccInterface := defaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
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
					Config: testAccJunosSecurityZoneConfigUpdate(testaccInterface),
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
resource "junos_security_screen" "testaccZone" {
  lifecycle {
    create_before_destroy = true
  }
  name        = "testaccZone"
  description = "testaccZone"
}
resource "junos_security_zone" "testacc_securityZone" {
  name = "testacc_securityZone"
  address_book {
    name        = "testacc_zone1"
    description = "testacc_zone 1"
    network     = "192.0.2.0/25"
  }
  address_book_dns {
    name        = "testacc_zone2"
    description = "testacc_zone 2"
    fqdn        = "test.com"
  }
  address_book_dns {
    name        = "testacc_zone2b"
    description = "testacc_zone 2b"
    fqdn        = "test.com"
    ipv4_only   = true
  }
  address_book_dns {
    name        = "testacc_zone2c"
    description = "testacc_zone 2c"
    fqdn        = "test.com"
    ipv6_only   = true
  }
  address_book_range {
    name        = "testacc_zone3"
    description = "testacc_zone 3"
    from        = "192.0.2.10"
    to          = "192.0.2.12"
  }
  address_book_set {
    name        = "testacc_zoneSet"
    description = "testacc_zone Set"
    address     = ["testacc_zone1"]
  }
  address_book_set {
    name        = "testacc_zoneSet2"
    description = "testacc_zone Set2"
    address     = ["testacc_zone2c"]
    address_set = ["testacc_zoneSet"]
  }
  address_book_wildcard {
    name        = "testacc_zone4"
    description = "testacc_zone 4"
    network     = "192.0.2.0/255.0.255.255"
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
	return fmt.Sprintf(`
resource "junos_security_zone" "testacc_securityZone" {
  name = "testacc_securityZone"
  address_book {
    name    = "testacc_zone1"
    network = "192.0.2.0/25"
  }
  address_book {
    name    = "testacc_zone2"
    network = "192.0.2.128/25"
  }
  address_book_set {
    name    = "testacc_zoneSet"
    address = ["testacc_zone1", "testacc_zone2"]
  }
  inbound_protocols = ["bgp"]
  inbound_services  = ["ssh"]
}
resource "junos_interface_logical" "testacc_securityZone" {
  name          = "%s.0"
  security_zone = junos_security_zone.testacc_securityZone.name
}
`, interFace)
}

func testAccJunosSecurityZoneConfigUpdate2(interFace string) string {
	return fmt.Sprintf(`
resource "junos_security_zone" "testacc_securityZone" {
  name = "testacc_securityZone"
}
resource "junos_interface_logical" "testacc_securityZone" {
  name          = "%s.0"
  security_zone = junos_security_zone.testacc_securityZone.name
}
`, interFace)
}
