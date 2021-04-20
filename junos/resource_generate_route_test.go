package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosGenerateRoute_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosGenerateRouteConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"routing_instance", "testacc_generateRoute"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"active", "true"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"full", "true"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"discard", "true"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"policy.#", "1"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"policy.0", "testacc_generateRoute"),
					),
				},
				{
					ResourceName:      "junos_generate_route.testacc_generateRoute",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_generate_route.testacc_generateRoute6",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosGenerateRouteConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"passive", "true"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"brief", "true"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"community.#", "0"),
						resource.TestCheckResourceAttr("junos_generate_route.testacc_generateRoute",
							"policy.#", "0"),
					),
				},
				{
					Config: testAccJunosGenerateRouteConfigUpdate2(),
				},
			},
		})
	}
}

func testAccJunosGenerateRouteConfigCreate() string {
	return `
resource junos_routing_instance testacc_generateRoute {
  name = "testacc_generateRoute"
}
resource junos_policyoptions_policy_statement "testacc_generateRoute" {
  name = "testacc_generateRoute"
  then {
    action = "accept"
  }
}

resource junos_generate_route testacc_generateRoute {
  destination                  = "192.0.2.0/24"
  routing_instance             = junos_routing_instance.testacc_generateRoute.name
  preference                   = 100
  metric                       = 100
  active                       = true
  full                         = true
  discard                      = true
  community                    = ["no-advertise"]
  policy                       = [junos_policyoptions_policy_statement.testacc_generateRoute.name]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
resource junos_generate_route testacc_generateRoute6 {
  destination                  = "2001:db8:85a3::/48"
  routing_instance             = junos_routing_instance.testacc_generateRoute.name
  preference                   = 100
  metric                       = 100
  active                       = true
  full                         = true
  discard                      = true
  community                    = ["no-advertise"]
  policy                       = [junos_policyoptions_policy_statement.testacc_generateRoute.name]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
`
}

func testAccJunosGenerateRouteConfigUpdate() string {
	return `
resource junos_routing_instance testacc_generateRoute {
  name = "testacc_generateRoute"
}
resource junos_policyoptions_policy_statement "testacc_generateRoute" {
  name = "testacc_generateRoute"
  then {
    action = "accept"
  }
}

resource junos_generate_route testacc_generateRoute {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_generateRoute.name
  passive          = true
  brief            = true
}
resource junos_generate_route testacc_generateRoute6 {
  destination      = "2001:db8:85a3::/48"
  routing_instance = junos_routing_instance.testacc_generateRoute.name
  passive          = true
  brief            = true
}
`
}

func testAccJunosGenerateRouteConfigUpdate2() string {
	return `
resource junos_routing_instance testacc_generateRoute {
  name = "testacc_generateRoute"
}
resource junos_routing_instance testacc_generateRoute2 {
  name = "testacc_generateRoute2"
}

resource junos_generate_route testacc_generateRoute {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_generateRoute.name
  next_table       = "${junos_routing_instance.testacc_generateRoute2.name}.inet.0"
}
`
}
