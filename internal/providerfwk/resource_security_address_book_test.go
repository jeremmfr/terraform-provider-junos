package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosSecurityAddressBook_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityAddressBookConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"name", "global"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"description", "testacc global description"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.#", "2"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.0.name", "testacc_network"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.0.description", "testacc_network description"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.0.value", "10.0.0.0/24"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.1.name", "testacc_network2"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.1.description", "testacc_network description2"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.1.value", "10.1.0.0/24"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"wildcard_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"wildcard_address.0.value", "10.0.0.0/255.255.0.255"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"wildcard_address.0.name", "testacc_wildcard"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"range_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"range_address.0.name", "testacc_range"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"range_address.0.from", "10.1.1.1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"range_address.0.to", "10.1.1.5"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"dns_name.#", "1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"dns_name.0.name", "testacc_dns"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"dns_name.0.value", "google.com"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"address_set.#", "1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"address_set.0.address.#", "3"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"name", "testacc_secAddrBook"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"attach_zone.#", "2"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"attach_zone.0", "testacc_secZoneAddr1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"attach_zone.1", "testacc_secZoneAddr2"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"network_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"network_address.0.name", "testacc_network"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"network_address.0.value", "10.1.2.3/32"),
					),
				},
				{
					Config: testAccJunosSecurityAddressBookConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityGlobalAddressBook",
							"network_address.0.value", "10.1.0.0/24"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"network_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_address_book.testacc_securityNamedAddressBook",
							"network_address.0.value", "10.1.2.4/32"),
					),
				},
				{
					ResourceName:      "junos_security_address_book.testacc_securityGlobalAddressBook",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityAddressBookConfigCreate() string {
	return `
resource "junos_security_zone" "testacc_secZoneAddr1" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_secZoneAddr1"
}
resource "junos_security_zone" "testacc_secZoneAddr2" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_secZoneAddr2"
}

resource "junos_security_address_book" "testacc_securityGlobalAddressBook" {
  description = "testacc global description"
  network_address {
    name        = "testacc_network"
    description = "testacc_network description"
    value       = "10.0.0.0/24"
  }
  wildcard_address {
    name        = "testacc_wildcard"
    description = "testacc_wildcard description"
    value       = "10.0.0.0/255.255.0.255"
  }
  network_address {
    name        = "testacc_network2"
    description = "testacc_network description2"
    value       = "10.1.0.0/24"
  }
  range_address {
    name        = "testacc_range"
    description = "testacc_range description"
    from        = "10.1.1.1"
    to          = "10.1.1.5"
  }
  dns_name {
    name  = "testacc_dns"
    value = "google.com"
  }
  address_set {
    name    = "testacc_addressSet"
    address = ["testacc_network", "testacc_wildcard", "testacc_network2"]
  }
}

resource "junos_security_address_book" "testacc_securityNamedAddressBook" {
  name        = "testacc_secAddrBook"
  attach_zone = [junos_security_zone.testacc_secZoneAddr1.name, junos_security_zone.testacc_secZoneAddr2.name]
  network_address {
    name  = "testacc_network"
    value = "10.1.2.3/32"
  }
}
`
}

func testAccJunosSecurityAddressBookConfigUpdate() string {
	return `
resource "junos_security_address_book" "testacc_securityGlobalAddressBook" {
  description = "testacc global description"
  network_address {
    name        = "testacc_network"
    description = "testacc_network description"
    value       = "10.1.0.0/24"
  }
  dns_name {
    name        = "testacc_dns"
    description = "testacc_dns description"
    value       = "google.com"
    ipv4_only   = true
  }
  dns_name {
    name      = "testacc_dns6"
    value     = "google.com"
    ipv6_only = true
  }
  address_set {
    name        = "testacc_addressSet"
    description = "testacc_addressSet description"
    address     = ["testacc_network", "testacc_dns"]
  }
  address_set {
    name        = "testacc_addressSet2"
    address     = ["testacc_dns"]
    address_set = ["testacc_addressSet"]
  }
}

resource "junos_security_address_book" "testacc_securityNamedAddressBook" {
  name = "testacc_secAddrBook"
  network_address {
    name  = "testacc_network"
    value = "10.1.2.4/32"
  }
}
`
}
