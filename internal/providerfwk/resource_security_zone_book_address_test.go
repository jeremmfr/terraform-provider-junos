package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosSecurityZoneBookAddress_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityZoneBookAddressConfigCreate(),
				},
				{
					ResourceName:      "junos_security_zone_book_address.testacc_szone_bookaddress1",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_zone_book_address.testacc_szone_bookaddress2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_zone_book_address.testacc_szone_bookaddress3",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_zone_book_address.testacc_szone_bookaddress5",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_zone_book_address.testacc_szone_bookaddress6",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSecurityZoneBookAddressConfigUpdate(),
				},
			},
		})
	}
}

func testAccJunosSecurityZoneBookAddressConfigCreate() string {
	return `
resource "junos_security_zone" "testacc_szone_bookaddress" {
  name                          = "testacc_szone_bookaddress"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress1" {
  name        = "testacc_szone_bookaddress1"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  cidr        = "192.0.2.0/25"
  description = "testacc szone bookaddress1"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress2" {
  name        = "testacc_szone_bookaddress2"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  dns_name    = "test.com"
  description = "testacc szone bookaddress2"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress3" {
  name          = "testacc_szone_bookaddress3"
  zone          = junos_security_zone.testacc_szone_bookaddress.name
  dns_name      = "test.com"
  description   = "testacc szone bookaddress3"
  dns_ipv4_only = true
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress4" {
  name          = "testacc_szone_bookaddress4"
  zone          = junos_security_zone.testacc_szone_bookaddress.name
  dns_name      = "test.com"
  description   = "testacc szone bookaddress4"
  dns_ipv6_only = true
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress5" {
  name        = "testacc_szone_bookaddress5"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  range_from  = "192.0.2.10"
  range_to    = "192.0.2.12"
  description = "testacc szone bookaddress5"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress6" {
  name        = "testacc_szone_bookaddress6"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  wildcard    = "192.0.2.0/255.0.255.255"
  description = "testacc szone bookaddress6"
}
`
}

func testAccJunosSecurityZoneBookAddressConfigUpdate() string {
	return `
resource "junos_security_zone" "testacc_szone_bookaddress" {
  name                          = "testacc_szone_bookaddress"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress1" {
  name        = "testacc_szone_bookaddress1"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  cidr        = "192.0.2.128/25"
  description = "testacc szone address1"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress2" {
  name        = "testacc_szone_bookaddress2"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  dns_name    = "test.fr"
  description = "testacc szone address2"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress3" {
  name        = "testacc_szone_bookaddress3"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  dns_name    = "test.net"
  description = "testacc szone address3"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress5" {
  name        = "testacc_szone_bookaddress5"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  range_from  = "192.0.2.20"
  range_to    = "192.0.2.22"
  description = "testacc szone ddress5"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress6" {
  name        = "testacc_szone_bookaddress6"
  zone        = junos_security_zone.testacc_szone_bookaddress.name
  wildcard    = "192.0.4.0/255.0.255.255"
  description = "testacc szone address6"
}
`
}
