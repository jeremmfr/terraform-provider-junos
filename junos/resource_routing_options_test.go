package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosRoutingOptions_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosRoutingOptionsConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.#", "1"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.0.number", "65000"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.0.asdot_notation", "true"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"autonomous_system.0.loops", "5"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.#", "1"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.0.restart_duration", "120"),
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.0.disable", "true"),
					),
				},
				{
					ResourceName:      "junos_routing_options.testacc_routing_options",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosRoutingOptionsConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_options.testacc_routing_options",
							"graceful_restart.#", "1"),
					),
				},
			},
		})
	}
}

func testAccJunosRoutingOptionsConfigCreate() string {
	return `
resource junos_routing_options "testacc_routing_options" {
  autonomous_system {
    number         = "65000"
    asdot_notation = true
    loops          = 5
  }
  graceful_restart {
    restart_duration = 120
    disable          = true
  }
}
`
}
func testAccJunosRoutingOptionsConfigUpdate() string {
	return `
resource junos_routing_options "testacc_routing_options" {
  graceful_restart {}
}
`
}
