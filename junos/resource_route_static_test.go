package junos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccJunosRouteStaticAndInstance_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosRouteStaticAndInstanceConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_routing_instance.testacc_routeSttIns",
						"as", "65000"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"routing_instance", "testacc_routeSttIns"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"preference", "100"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"metric", "100"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"next_hop.#", "1"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"next_hop.0", "st0.0"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"qualified_next_hop.#", "1"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"qualified_next_hop.0.next_hop", "st0.0"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"qualified_next_hop.0.preference", "101"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"qualified_next_hop.0.metric", "101"),
				),
			},
			{
				Config: testAccJunosRouteStaticAndInstanceConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_routing_instance.testacc_routeSttIns",
						"as", "65001"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"qualified_next_hop.#", "2"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"qualified_next_hop.1.next_hop", "dsc.0"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"qualified_next_hop.1.preference", "102"),
					resource.TestCheckResourceAttr("junos_route_static.testacc_routeStt",
						"qualified_next_hop.1.metric", "102"),
				),
			},
			{
				ResourceName:      "junos_routing_instance.testacc_routeSttIns",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "junos_route_static.testacc_routeStt",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosRouteStaticAndInstanceConfigCreate() string {
	return fmt.Sprintf(`
resource junos_routing_instance testacc_routeSttIns {
  name = "testacc_routeSttIns"
  as = 65000
}
resource junos_route_static testacc_routeStt {
  destination = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_routeSttIns.name
  preference = 100
  metric = 100
  next_hop = [ "st0.0" ]
  qualified_next_hop {
    next_hop = "st0.0"
    preference = 101
    metric = 101
  }
}
`)
}
func testAccJunosRouteStaticAndInstanceConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_routing_instance testacc_routeSttIns {
  name = "testacc_routeSttIns"
  as = 65001
}
resource junos_route_static testacc_routeStt {
  destination = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_routeSttIns.name
  preference = 100
  metric = 100
  next_hop = [ "st0.0" ]
  qualified_next_hop {
    next_hop = "st0.0"
    preference = 101
    metric = 101
  }
  qualified_next_hop {
    next_hop = "dsc.0"
    preference = 102
    metric = 102
  }
}
`)
}
