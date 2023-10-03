package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAggregateRoute_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
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
					ConfigDirectory: config.TestStepDirectory(),
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
