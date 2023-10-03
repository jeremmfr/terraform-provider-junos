package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceInterfaceLogical_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
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
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL",
							"id", testaccInterface+".100"),
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL",
							"name", testaccInterface+".100"),
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL",
							"family_inet.address.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL",
							"family_inet.address.0.cidr_ip", "192.0.2.1/25"),
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL2",
							"id", testaccInterface+".100"),
					),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}
