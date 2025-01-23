package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityDynamicAddressFeedServer_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_dynamic_address_feed_server.testacc_dyn_add_feed_srv",
							"feed_name.#", "2"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_dynamic_address_feed_server.testacc_dyn_add_feed_srv",
							"feed_name.#", "3"),
					),
				},
				{
					ResourceName:      "junos_security_dynamic_address_feed_server.testacc_dyn_add_feed_srv",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
