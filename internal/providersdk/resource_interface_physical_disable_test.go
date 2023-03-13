package providersdk_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccJunosInterfacePhysicalDisable_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface = junos.DefaultInterfaceSwitchTestAcc
	}
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosInterfacePhysicalDisablePreConfigCreate(testaccInterface),
			},
			{
				Config:  testAccJunosInterfacePhysicalDisablePreConfigCreate(testaccInterface),
				Destroy: true,
			},
			{
				Config: testAccJunosInterfacePhysicalDisableConfigCreate(testaccInterface),
			},
			{
				Config:             testAccJunosInterfacePhysicalDisableConfigConflict(testaccInterface),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      testAccJunosInterfacePhysicalDisableConfigConflict(testaccInterface),
				ExpectError: regexp.MustCompile("interface " + testaccInterface + " is configured"),
			},
		},
	})
}

func testAccJunosInterfacePhysicalDisablePreConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface_disable" {
  name                  = "%s"
  no_disable_on_destroy = true
}
resource "junos_interface_logical" "testacc_interface_disable" {
  name        = "${junos_interface_physical.testacc_interface_disable.name}.0"
  description = "testacc_interface_disable"
}
`, interFace)
}

func testAccJunosInterfacePhysicalDisableConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical_disable" "testacc_interface_disable" {
  name = "%s"
}
resource "junos_interface_physical_disable" "testacc_interface_disable2" {
  name = "%s"
}
`, interFace, interFace)
}

func testAccJunosInterfacePhysicalDisableConfigConflict(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface_disable" {
  name                  = "%s"
  description           = "testacc_interface_disable"
  no_disable_on_destroy = true
}
resource "junos_interface_physical_disable" "testacc_interface_disable" {
  name = "%s"
}
`, interFace, interFace)
}
