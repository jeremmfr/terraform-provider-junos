package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceInterfacePhysical_basic(t *testing.T) {
	var testaccInterface string
	if os.Getenv("TESTACC_INTERFACE") != "" {
		testaccInterface = os.Getenv("TESTACC_INTERFACE")
	} else {
		testaccInterface = defaultInterfaceTestAcc
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceInterfacePhysicalConfigCreate(testaccInterface),
				},
				{
					Config: testAccDataSourceInterfacePhysicalConfigData(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_datainterfaceP",
							"id", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_datainterfaceP",
							"vlan_tagging", "true"),
					),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}

func testAccDataSourceInterfacePhysicalConfigCreate(interFace string) string {
	return `
resource junos_interface_physical testacc_datainterfaceP {
  name         = "` + interFace + `"
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}
`
}

func testAccDataSourceInterfacePhysicalConfigData(interFace string) string {
	return `
resource junos_interface_physical testacc_datainterfaceP {
  name         = "` + interFace + `"
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}

data junos_interface_physical testacc_datainterfaceP {
  config_interface = "` + interFace + `"
}
`
}
