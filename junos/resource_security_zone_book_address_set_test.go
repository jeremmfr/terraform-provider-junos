package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityZoneBookAddressSet_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityZoneBookAddressSetConfigCreate(),
				},
				{
					ResourceName:      "junos_security_zone_book_address_set.testacc_szone_bookaddress_set",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSecurityZoneBookAddressSetConfigUpdate(),
				},
			},
		})
	}
}

func testAccJunosSecurityZoneBookAddressSetConfigCreate() string {
	return `
resource "junos_security_zone" "testacc_szone_bookaddressset" {
  name                          = "testacc_szone_bookaddressset"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress_set" {
  name = "testacc_szone_bookaddress_set1"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  cidr = "192.0.2.0/25"
}
resource "junos_security_zone_book_address_set" "testacc_szone_bookaddress_set" {
  name = "testacc_szone_bookaddress_set"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  address = [
    junos_security_zone_book_address.testacc_szone_bookaddress_set.name,
  ]
  description = "testacc szone bookaddress set"
}
`
}

func testAccJunosSecurityZoneBookAddressSetConfigUpdate() string {
	return `
resource "junos_security_zone" "testacc_szone_bookaddressset" {
  name                          = "testacc_szone_bookaddressset"
  address_book_configure_singly = true
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress_set" {
  name = "testacc_szone_bookaddress_set1"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  cidr = "192.0.2.0/25"
}
resource "junos_security_zone_book_address" "testacc_szone_bookaddress_set2" {
  name = "testacc_szone_bookaddress_set2"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  cidr = "192.0.2.128/25"
}
resource "junos_security_zone_book_address_set" "testacc_szone_bookaddress_set" {
  name = "testacc_szone_bookaddress_set"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  address = [
    junos_security_zone_book_address.testacc_szone_bookaddress_set.name,
    junos_security_zone_book_address.testacc_szone_bookaddress_set2.name,
  ]
}
resource "junos_security_zone_book_address_set" "testacc_szone_bookaddress_setS2" {
  name = "testacc_szone_bookaddress_setS2"
  zone = junos_security_zone.testacc_szone_bookaddressset.name
  address = [
    junos_security_zone_book_address.testacc_szone_bookaddress_set2.name,
  ]
  address_set = [
    junos_security_zone_book_address_set.testacc_szone_bookaddress_set.name,
  ]
}
`
}
