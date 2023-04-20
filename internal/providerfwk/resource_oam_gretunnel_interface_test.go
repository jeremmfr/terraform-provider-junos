package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosOamGretunnelInterface_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosOamGretunnelInterfaceConfigCreate(),
				},
				{
					Config: testAccJunosOamGretunnelInterfaceConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_oam_gretunnel_interface.testacc_oam_gretunnel_interface",
							"hold_time", "10"),
						resource.TestCheckResourceAttr("junos_oam_gretunnel_interface.testacc_oam_gretunnel_interface",
							"keepalive_time", "5"),
					),
				},
				{
					ResourceName:      "junos_oam_gretunnel_interface.testacc_oam_gretunnel_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosOamGretunnelInterfaceConfigCreate() string {
	return `
resource "junos_oam_gretunnel_interface" "testacc_oam_gretunnel_interface" {
  name = "gr-3/3/0.3"
}
`
}

func testAccJunosOamGretunnelInterfaceConfigUpdate() string {
	return `
resource "junos_oam_gretunnel_interface" "testacc_oam_gretunnel_interface" {
  name           = "gr-3/3/0.3"
  hold_time      = 10
  keepalive_time = 5
}
resource "junos_oam_gretunnel_interface" "testacc_oam_gretunnel_interface2" {
  name      = "gr-3/3/0.2"
  hold_time = 11
}
resource "junos_oam_gretunnel_interface" "testacc_oam_gretunnel_interface3" {
  name           = "gr-3/3/0.1"
  keepalive_time = 2
}
`
}
