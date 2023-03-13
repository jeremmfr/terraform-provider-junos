package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceInterfacePhysical_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
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
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_datainterfaceP" {
  name         = "%s"
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}
`, interFace)
}

func testAccDataSourceInterfacePhysicalConfigData(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_datainterfaceP" {
  name         = "%s"
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}

data "junos_interface_physical" "testacc_datainterfaceP" {
  config_interface = "%s"
}
`, interFace, interFace)
}
