package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// the manual block of junos_security_ipsec_vpn is only available on Junos versions before 22,
// so the resource is created with a count in the test configurations.
const testAccSecurityIpsecVpnWriteOnlyName = "junos_security_ipsec_vpn.testacc_ipsecvpn_wo[0]"

// testCheckIfSecurityIpsecVpnManual runs the checks only when the manual block is available:
// on a recent device the count is 0 and there is nothing to check,
// but the checks would fail on the missing resource instead of being skipped.
func testCheckIfSecurityIpsecVpnManual(checks ...resource.TestCheckFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, ok := s.RootModule().Resources[testAccSecurityIpsecVpnWriteOnlyName]; !ok {
			return nil
		}

		return resource.ComposeTestCheckFunc(checks...)(s)
	}
}

// export TESTACC_INTERFACE=<inteface> to choose interface available else it's ge-0/0/3.
func TestAccResourceSecurityIpsecVpn_writeOnly(t *testing.T) {
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
					Check: testCheckIfSecurityIpsecVpnManual(
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_text"),
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_text_wo"),
						resource.TestCheckResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_text_wo_version", "1"),
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.encryption_key_text"),
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.encryption_key_text_wo"),
						resource.TestCheckResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.encryption_key_text_wo_version", "1"),
					),
				},
				{
					// check that the write-only keys have really been sent to the device
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: testCheckIfSecurityIpsecVpnManual(
						resource.TestMatchResourceAttr("data.junos_config_raw.testacc_ipsecvpn_wo",
							"config", regexp.MustCompile(
								`set security ipsec vpn "?testacc_ipsecvpn_wo"? manual authentication key ascii-text `)),
						resource.TestMatchResourceAttr("data.junos_config_raw.testacc_ipsecvpn_wo",
							"config", regexp.MustCompile(
								`set security ipsec vpn "?testacc_ipsecvpn_wo"? manual encryption key ascii-text `)),
					),
				},
				{
					// increment the version of the authentication key only
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: testCheckIfSecurityIpsecVpnManual(
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_text"),
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_text_wo"),
						resource.TestCheckResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_text_wo_version", "2"),
						resource.TestCheckResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.encryption_key_text_wo_version", "1"),
					),
				},
				{
					// switch the authentication key to the hexadecimal format, also write-only,
					// and the encryption key back to the standard argument
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: testCheckIfSecurityIpsecVpnManual(
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_hexa"),
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_hexa_wo"),
						resource.TestCheckResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_hexa_wo_version", "1"),
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.authentication_key_text_wo_version"),
						resource.TestCheckResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.encryption_key_text", "Encryp"),
						resource.TestCheckNoResourceAttr(testAccSecurityIpsecVpnWriteOnlyName,
							"manual.encryption_key_text_wo_version"),
					),
				},
			},
		})
	}
}
