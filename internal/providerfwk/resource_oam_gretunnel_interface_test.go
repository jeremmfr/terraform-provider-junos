package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceOamGretunnelInterface_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
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
