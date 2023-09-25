package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceGenerateRoute_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
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
					ConfigDirectory: config.TestStepDirectory(),
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
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
