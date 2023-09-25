package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityUtmProfileWebFL_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecurityUtmProfileWebFLConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"default_action", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.0.default", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.0.server_connectivity", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.0.timeout", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.0.timeout", "log-and-permit"),
					),
				},
				{
					Config: testAccResourceSecurityUtmProfileWebFLConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_local.testacc_ProfileWebFL",
							"fallback_settings.#", "0"),
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

func testAccResourceSecurityUtmProfileWebFLConfigCreate() string {
	return `
resource "junos_security_utm_profile_web_filtering_juniper_local" "testacc_ProfileWebFL" {
  name                 = "testacc ProfileWebFL"
  custom_block_message = "Blocked by Juniper"
  default_action       = "log-and-permit"
  fallback_settings {
    default             = "log-and-permit"
    server_connectivity = "log-and-permit"
    timeout             = "log-and-permit"
  }
}
`
}

func testAccResourceSecurityUtmProfileWebFLConfigUpdate() string {
	return `
resource "junos_security_utm_profile_web_filtering_juniper_local" "testacc_ProfileWebFL" {
  name                 = "testacc ProfileWebFL"
  custom_block_message = "Blocked by Juniper"
  default_action       = "log-and-permit"
  timeout              = 3
}
`
}
