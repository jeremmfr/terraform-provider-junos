package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccResourceLldpMedInterface_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface := junos.DefaultInterfaceSwitchTestAcc
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceLldpMedInterfaceSWConfigCreate(testaccInterface),
				},
				{
					Config: testAccResourceLldpMedInterfaceSWConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_lldpmed_interface.testacc_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		testaccInterface := junos.DefaultInterfaceTestAcc
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceLldpMedInterfaceConfigCreate(testaccInterface),
				},
				{
					Config: testAccResourceLldpMedInterfaceConfigUpdate(testaccInterface),
				},
				{
					Config: testAccResourceLldpMedInterfaceConfigUpdate2(testaccInterface),
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

func testAccResourceLldpMedInterfaceSWConfigCreate(interFace string) string {
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

func testAccResourceLldpMedInterfaceSWConfigUpdate(interFace string) string {
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

func testAccResourceLldpMedInterfaceConfigCreate(interFace string) string {
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

func testAccResourceLldpMedInterfaceConfigUpdate(interFace string) string {
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

func testAccResourceLldpMedInterfaceConfigUpdate2(interFace string) string {
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
