package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityUtmProfileWebFilteringJuniperLocal_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"default_action", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.default", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.server_connectivity", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.timeout", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.timeout", "log-and-permit"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"timeout", "3"),
					),
				},
				{
					ResourceName:      "junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
