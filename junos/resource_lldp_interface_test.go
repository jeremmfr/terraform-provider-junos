package junos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccJunosLldpInterface_basic(t *testing.T) {
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
					Config: testAccJunosLldpInterfaceSWConfigCreate(testaccInterface),
				},
				{
					Config: testAccJunosLldpInterfaceSWConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_lldp_interface.testacc_interface",
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
					Config: testAccJunosLldpInterfaceConfigCreate(testaccInterface),
				},
				{
					Config: testAccJunosLldpInterfaceConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_lldp_interface.testacc_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosLldpInterfaceSWConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldp_interface" "testacc_all" {
  name = "all"
}
resource "junos_lldp_interface" "testacc_interface" {
  name = "%s"
  power_negotiation {}
}
`, interFace)
}

func testAccJunosLldpInterfaceSWConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldp_interface" "testacc_all" {
  name   = "all"
  enable = true
  power_negotiation {
    enable = true
  }

}
resource "junos_lldp_interface" "testacc_interface" {
  name    = "%s"
  disable = true
  power_negotiation {
    disable = true
  }
}
`, interFace)
}

func testAccJunosLldpInterfaceConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldp_interface" "testacc_all" {
  name = "all"
}
resource "junos_lldp_interface" "testacc_interface" {
  name = "%s"
}
`, interFace)
}

func testAccJunosLldpInterfaceConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldp_interface" "testacc_all" {
  name                     = "all"
  enable                   = true
  trap_notification_enable = true

}
resource "junos_lldp_interface" "testacc_interface" {
  name                      = "%s"
  disable                   = true
  trap_notification_disable = true
}
`, interFace)
}
