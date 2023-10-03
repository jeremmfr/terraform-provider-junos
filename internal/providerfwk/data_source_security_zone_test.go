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
func TestAccDataSourceSecurityZone_basic(t *testing.T) {
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
						resource.TestCheckResourceAttr("data.junos_security_zone.testacc_dataSecurityZone",
							"id", "testacc_dataSecurityZone"),
						resource.TestCheckResourceAttr("data.junos_security_zone.testacc_dataSecurityZone",
							"address_book.#", "1"),
						resource.TestCheckResourceAttr("data.junos_security_zone.testacc_dataSecurityZone",
							"interface.#", "1"),
						resource.TestCheckResourceAttr("data.junos_security_zone.testacc_dataSecurityZone",
							"interface.0.inbound_services.#", "1"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ExpectError:     regexp.MustCompile("security zone .* doesn't exist"),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}
