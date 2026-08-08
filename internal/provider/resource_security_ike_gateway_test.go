package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// export TESTACC_INTERFACE=<inteface> to choose interface available else it's ge-0/0/3.
func TestAccResourceSecurityIkeGateway_writeOnly(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_11_0),
			},
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_security_ike_gateway.testacc_ikegw_wo",
							"aaa.client_password"),
						resource.TestCheckNoResourceAttr("junos_security_ike_gateway.testacc_ikegw_wo",
							"aaa.client_password_wo"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegw_wo",
							"aaa.client_password_wo_version", "1"),
					),
				},
				{
					// check that the write-only password has really been sent to the device
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_ikegw_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`set security ike gateway "?testacc_ikegw_wo"? aaa client password `)),
						),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_security_ike_gateway.testacc_ikegw_wo",
							"aaa.client_password"),
						resource.TestCheckNoResourceAttr("junos_security_ike_gateway.testacc_ikegw_wo",
							"aaa.client_password_wo"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegw_wo",
							"aaa.client_password_wo_version", "2"),
					),
				},
				{
					// switch back to the standard argument
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegw_wo",
							"aaa.client_password", "thePassWord3"),
						resource.TestCheckNoResourceAttr("junos_security_ike_gateway.testacc_ikegw_wo",
							"aaa.client_password_wo_version"),
					),
				},
			},
		})
	}
}
