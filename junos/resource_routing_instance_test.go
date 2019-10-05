package junos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccJunosRoutingInstance_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
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

func testAccJunosRoutingInstanceConfigCreate() string {
	return fmt.Sprintf(`
resource junos_routing_instance "testacc_routingInst" {
  name = "testacc_routingInst"
  as = "65000"
}
`)
}
func testAccJunosRoutingInstanceConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_routing_instance "testacc_routingInst" {
  name = "testacc_routingInst"
  as = "65001"
}
`)
}
