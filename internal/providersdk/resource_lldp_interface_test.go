package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccResourceLldpInterface_basic(t *testing.T) {
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
					Config: testAccResourceLldpInterfaceSWConfigCreate(testaccInterface),
				},
				{
					Config: testAccResourceLldpInterfaceSWConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_lldp_interface.testacc_interface",
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
					Config: testAccResourceLldpInterfaceConfigCreate(testaccInterface),
				},
				{
					Config: testAccResourceLldpInterfaceConfigUpdate(testaccInterface),
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

func testAccResourceLldpInterfaceSWConfigCreate(interFace string) string {
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

func testAccResourceLldpInterfaceSWConfigUpdate(interFace string) string {
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

func testAccResourceLldpInterfaceConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_lldp_interface" "testacc_all" {
  name = "all"
}
resource "junos_lldp_interface" "testacc_interface" {
  name = "%s"
}
`, interFace)
}

func testAccResourceLldpInterfaceConfigUpdate(interFace string) string {
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
