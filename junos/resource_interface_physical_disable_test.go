package junos_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccJunosInterfacePhysicalDisable_basic(t *testing.T) {
	var testaccInterface string
	if os.Getenv("TESTACC_INTERFACE") != "" {
		testaccInterface = os.Getenv("TESTACC_INTERFACE")
	} else {
		testaccInterface = defaultInterfaceTestAcc
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
	return `
resource junos_interface_physical testacc_interface_disable {
  name                  = "` + interFace + `"
  no_disable_on_destroy = true
}
`
}

func testAccJunosInterfacePhysicalDisableConfigCreate(interFace string) string {
	return `
resource junos_interface_physical_disable testacc_interface_disable {
  name = "` + interFace + `"
}
resource junos_interface_physical_disable testacc_interface_disable2 {
  name = "` + interFace + `"
}
`
}

func testAccJunosInterfacePhysicalDisableConfigConflict(interFace string) string {
	return `
resource junos_interface_physical testacc_interface_disable {
  name                  = "` + interFace + `"
  description           = "testacc_interface_disable"
  no_disable_on_destroy = true
}
resource junos_interface_physical_disable testacc_interface_disable {
  name = "` + interFace + `"
}
`
}
