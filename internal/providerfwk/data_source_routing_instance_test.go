package providerfwk_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceRoutingInstance_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
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
						resource.TestCheckResourceAttr("data.junos_routing_instance.testacc_dataRoutingInstance",
							"id", "testacc_dataRoutingInstance"),
						resource.TestCheckResourceAttr("data.junos_routing_instance.testacc_dataRoutingInstance",
							"type", "virtual-router"),
						resource.TestCheckResourceAttr("data.junos_routing_instance.testacc_dataRoutingInstance",
							"interface.#", "1"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ExpectError:     regexp.MustCompile("routing instance .* doesn't exist"),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}
