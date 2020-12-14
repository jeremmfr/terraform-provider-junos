package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"routing_instance", "testacc_staticRoute"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.0.metric", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.next_hop", "192.0.2.250"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.interface", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"qualified_next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"qualified_next_hop.0.metric", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"routing_instance", "testacc_staticRoute"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.0.metric", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.1.next_hop", "2001:db8:85a4::1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.1.interface", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.0.metric", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"community.0", "no-advertise"),
					),
				},
				{
					Config: testAccJunosStaticRouteConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.next_hop", "dsc.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.preference", "102"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.metric", "102"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.1.next_hop", "dsc.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.1.preference", "102"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.1.metric", "102"),
					),
				},
				{
					ResourceName:      "junos_static_route.testacc_staticRoute_instance",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosStaticRouteConfigCreate() string {
	return `
resource junos_routing_instance testacc_staticRoute {
  name = "testacc_staticRoute"
}
resource junos_static_route testacc_staticRoute_instance {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  preference       = 100
  metric           = 100
  next_hop         = ["st0.0"]
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop   = "192.0.2.250"
    interface  = "st0.0"
  }
  community = ["no-advertise"]
}
resource junos_static_route testacc_staticRoute_default {
  destination = "192.0.2.0/24"
  preference  = 100
  metric      = 100
  next_hop    = ["st0.0"]
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  community = ["no-advertise"]
}
resource junos_static_route testacc_staticRoute_ipv6_default {
  destination = "2001:db8:85a3::/48"
  preference  = 100
  metric      = 100
  next_hop    = ["st0.0"]
  qualified_next_hop {
    next_hop = "st0.0"
    preference = 101
    metric = 101
  }
  community = ["no-advertise"]
}
resource junos_static_route testacc_staticRoute_ipv6_instance {
  destination      = "2001:db8:85a3::/48"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  preference       = 100
  metric           = 100
  next_hop         = ["st0.0"]
  qualified_next_hop {
    next_hop = "st0.0"
    preference = 101
    metric = 101
  }
  qualified_next_hop {
    next_hop   = "2001:db8:85a4::1"
    interface  = "st0.0"
  }
  community = ["no-advertise"]
}
`
}
func testAccJunosStaticRouteConfigUpdate() string {
	return `
resource junos_routing_instance testacc_staticRoute {
  name = "testacc_staticRoute"
}
resource junos_static_route testacc_staticRoute_instance {
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
resource junos_static_route testacc_staticRoute_ipv6_default {
  destination = "2001:db8:85a3::/48"
  preference  = 100
  metric      = 100
  next_hop    = ["st0.0"]
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
  community = ["no-advertise"]
}
`
}
