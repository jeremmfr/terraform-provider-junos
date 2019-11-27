package junos

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccJunosStaticRoute_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosStaticRouteConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"routing_instance", "testacc_staticRoute"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"qualified_next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"qualified_next_hop.0.metric", "101"),
					),
				},
				{
					Config: testAccJunosStaticRouteConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"qualified_next_hop.1.next_hop", "dsc.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"qualified_next_hop.1.preference", "102"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute",
							"qualified_next_hop.1.metric", "102"),
					),
				},
				{
					ResourceName:      "junos_static_route.testacc_staticRoute",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosStaticRouteConfigCreate() string {
	return fmt.Sprintf(`
resource junos_routing_instance testacc_staticRoute {
  name = "testacc_staticRoute"
}
resource junos_static_route testacc_staticRoute {
  destination = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
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
func testAccJunosStaticRouteConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_routing_instance testacc_staticRoute {
  name = "testacc_staticRoute"
}
resource junos_static_route testacc_staticRoute {
  destination = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
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
