package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccResourceSecurityIkePolicy_writeOnly(t *testing.T) {
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
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_text"),
						resource.TestCheckNoResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_text_wo"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_text_wo_version", "1"),
					),
				},
				{
					// check that the write-only key has really been sent to the device
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_ikepol_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`set security ike policy "?testacc_ikepol_wo"? pre-shared-key ascii-text `)),
						),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_text"),
						resource.TestCheckNoResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_text_wo"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_text_wo_version", "2"),
					),
				},
				{
					// switch to the hexadecimal format, also write-only
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_hexa"),
						resource.TestCheckNoResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_hexa_wo"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_hexa_wo_version", "1"),
						resource.TestCheckNoResourceAttr("junos_security_ike_policy.testacc_ikepol_wo",
							"pre_shared_key_text_wo_version"),
					),
				},
			},
		})
	}
}
