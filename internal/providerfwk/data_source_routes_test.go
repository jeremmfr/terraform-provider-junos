package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceRoutes_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" || os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckTypeSetElemNestedAttrs("data.junos_routes.all",
							"table.*", map[string]string{"name": "inet.0"}),
						resource.TestCheckTypeSetElemNestedAttrs("data.junos_routes.all",
							"table.*", map[string]string{"name": "testacc_data_routes.inet.0"}),
						resource.TestCheckTypeSetElemNestedAttrs("data.junos_routes.all",
							"table.*.route.*", map[string]string{"destination": "0.0.0.0/0"}),
						resource.TestCheckTypeSetElemNestedAttrs("data.junos_routes.all",
							"table.*.route.*.entry.*", map[string]string{
								"current_active": "true",
								"protocol":       "Local",
								"next_hop_type":  "Local",
							}),
						resource.TestCheckResourceAttr("data.junos_routes.default",
							"table.#", "1"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.#", "1"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.#", "1"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.0.destination", "192.0.2.1/32"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.0.entry.#", "1"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.0.entry.0.current_active", "true"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.0.entry.0.protocol", "Local"),
					),
				},
			},
		})
	}
}
