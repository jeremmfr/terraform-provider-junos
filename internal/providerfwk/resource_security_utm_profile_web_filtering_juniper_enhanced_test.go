package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityUtmProfileWebFilteringJuniperEnhanced_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"block_message.url", "block.local"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"block_message.type_custom_redirect_url", "true"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.#", "2"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.0.name", "Enhanced_Network_Errors"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.0.action", "block"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.1.name", "Enhanced_Suspicious_Content"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.1.reputation_action.#", "2"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.1.reputation_action.0.site_reputation", "very-safe"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.1.reputation_action.0.action", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.1.reputation_action.1.site_reputation", "moderately-safe"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.1.reputation_action.1.action", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"default_action", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.default", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.server_connectivity", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.timeout", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.timeout", "log-and-permit"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"quarantine_custom_message", "Quarantine by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"no_safe_search", "true"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"quarantine_message.url", "quarantine.local"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"quarantine_message.type_custom_redirect_url", "true"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"site_reputation_action.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"site_reputation_action.0.site_reputation", "harmful"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"site_reputation_action.0.action", "block"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"timeout", "3"),
					),
				},
				{
					ResourceName:      "junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
