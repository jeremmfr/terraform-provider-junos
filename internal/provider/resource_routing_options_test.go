package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRoutingOptions_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.number", "65000"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.asdot_notation", "true"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.loops", "5"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.restart_duration", "120"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.disable", "true"),
					),
				},
				{
					ResourceName:      "junos_routing_options.testacc_routing_options",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
