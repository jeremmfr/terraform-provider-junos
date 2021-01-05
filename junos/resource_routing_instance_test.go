package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosRoutingInstance_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosRoutingInstanceConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"type", "virtual-router"),
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65000"),
					),
				},
				{
					Config: testAccJunosRoutingInstanceConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65001"),
					),
				},
				{
					ResourceName:      "junos_routing_instance.testacc_routingInst",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosRoutingInstanceConfigCreate() string {
	return `
resource junos_routing_instance "testacc_routingInst" {
  name = "testacc_routingInst"
  as   = "65000"
}
`
}
func testAccJunosRoutingInstanceConfigUpdate() string {
	return `
resource junos_routing_instance "testacc_routingInst" {
  name = "testacc_routingInst"
  as   = "65001"
}
`
}
