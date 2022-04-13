package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosAggregateRoute_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosAggregateRouteConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"routing_instance", "testacc_aggregateRoute"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"active", "true"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"full", "true"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"discard", "true"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"policy.#", "1"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"policy.0", "testacc_aggregateRoute"),
					),
				},
				{
					ResourceName:      "junos_aggregate_route.testacc_aggregateRoute",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_aggregate_route.testacc_aggregateRoute6",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosAggregateRouteConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"passive", "true"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"brief", "true"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"community.#", "0"),
						resource.TestCheckResourceAttr("junos_aggregate_route.testacc_aggregateRoute",
							"policy.#", "0"),
					),
				},
			},
		})
	}
}

func testAccJunosAggregateRouteConfigCreate() string {
	return `
resource "junos_routing_instance" "testacc_aggregateRoute" {
  name = "testacc_aggregateRoute"
}
resource "junos_policyoptions_policy_statement" "testacc_aggregateRoute" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_aggregateRoute"
  then {
    action = "accept"
  }
}

resource "junos_aggregate_route" "testacc_aggregateRoute" {
  destination                  = "192.0.2.0/24"
  routing_instance             = junos_routing_instance.testacc_aggregateRoute.name
  preference                   = 100
  metric                       = 100
  active                       = true
  full                         = true
  discard                      = true
  community                    = ["no-advertise"]
  policy                       = [junos_policyoptions_policy_statement.testacc_aggregateRoute.name]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
resource "junos_aggregate_route" "testacc_aggregateRoute6" {
  destination                  = "2001:db8:85a3::/48"
  routing_instance             = junos_routing_instance.testacc_aggregateRoute.name
  preference                   = 100
  metric                       = 100
  active                       = true
  full                         = true
  discard                      = true
  community                    = ["no-advertise"]
  policy                       = [junos_policyoptions_policy_statement.testacc_aggregateRoute.name]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
`
}

func testAccJunosAggregateRouteConfigUpdate() string {
	return `
resource "junos_routing_instance" "testacc_aggregateRoute" {
  name = "testacc_aggregateRoute"
}

resource "junos_aggregate_route" "testacc_aggregateRoute" {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_aggregateRoute.name
  passive          = true
  brief            = true
}
resource "junos_aggregate_route" "testacc_aggregateRoute6" {
  destination      = "2001:db8:85a3::/48"
  routing_instance = junos_routing_instance.testacc_aggregateRoute.name
  passive          = true
  brief            = true
}
`
}
