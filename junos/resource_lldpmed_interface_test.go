package junos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccJunosLldpMedInterface_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface := defaultInterfaceSwitchTestAcc
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosLldpMedInterfaceSWConfigCreate(testaccInterface),
				},
				{
					Config: testAccJunosLldpMedInterfaceSWConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_lldpmed_interface.testacc_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		testaccInterface := defaultInterfaceTestAcc
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosLldpMedInterfaceConfigCreate(testaccInterface),
				},
				{
					Config: testAccJunosLldpMedInterfaceConfigUpdate(testaccInterface),
				},
				{
					Config: testAccJunosLldpMedInterfaceConfigUpdate2(testaccInterface),
				},
				{
					ResourceName:      "junos_lldpmed_interface.testacc_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosLldpMedInterfaceSWConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldpmed_interface" "testacc_all" {
  name = "all"
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name = "%s"
  location {}
}
`, interFace)
}

func testAccJunosLldpMedInterfaceSWConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldpmed_interface" "testacc_all" {
  name = "all"
  location {
    civic_based_country_code = "FR"
    civic_based_what         = 1
  }
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name = "%s"
  location {
    civic_based_country_code = "FR"
    civic_based_ca_type {
      ca_type  = 10
      ca_value = "testacc"
    }
    civic_based_ca_type {
      ca_type = 0
    }
  }
}
`, interFace)
}

func testAccJunosLldpMedInterfaceConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldpmed_interface" "testacc_all" {
  name    = "all"
  disable = true
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name   = "%s"
  enable = true
  location {}
}
`, interFace)
}

func testAccJunosLldpMedInterfaceConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldpmed_interface" "testacc_all" {
  name = "all"
  location {
    civic_based_country_code = "FR"
    civic_based_what         = 1
  }
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name = "%s"
  location {
    civic_based_country_code = "UK"
    civic_based_ca_type {
      ca_type  = 20
      ca_value = "testacc"
    }
    civic_based_ca_type {
      ca_type = 0
    }
  }
}
`, interFace)
}

func testAccJunosLldpMedInterfaceConfigUpdate2(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldpmed_interface" "testacc_all" {
  name = "all"
  location {
    elin = "testacc_lldpmed"
  }
}
resource "junos_lldpmed_interface" "testacc_interface" {
  name = "%s"
  location {
    co_ordinate_latitude  = 180
    co_ordinate_longitude = 180
  }
}
`, interFace)
}
