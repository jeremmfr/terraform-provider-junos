package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityUtmProfileWebFE_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecurityUtmProfileWebFEConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"block_message.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"block_message.0.url", "block.local"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"block_message.0.type_custom_redirect_url", "true"),
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
							"fallback_settings.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.0.default", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.0.server_connectivity", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.0.timeout", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.0.timeout", "log-and-permit"),
					),
				},
				{
					Config: testAccResourceSecurityUtmProfileWebFEConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"block_message.#", "0"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"category.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"fallback_settings.#", "0"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"quarantine_custom_message", "Quarantine by Juniper"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"no_safe_search", "true"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"quarantine_message.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"quarantine_message.0.url", "quarantine.local"),
						resource.TestCheckResourceAttr("junos_security_utm_profile_web_filtering_juniper_enhanced.testacc_ProfileWebFE",
							"quarantine_message.0.type_custom_redirect_url", "true"),
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

func testAccResourceSecurityUtmProfileWebFEConfigCreate() string {
	return `
resource "junos_security_utm_profile_web_filtering_juniper_enhanced" "testacc_ProfileWebFE" {
  name = "testacc ProfileWebFE"
  block_message {
    url                      = "block.local"
    type_custom_redirect_url = true
  }
  category {
    name   = "Enhanced_Network_Errors"
    action = "block"
  }
  category {
    name   = "Enhanced_Suspicious_Content"
    action = "quarantine"
    reputation_action {
      site_reputation = "very-safe"
      action          = "log-and-permit"
    }
    reputation_action {
      site_reputation = "moderately-safe"
      action          = "log-and-permit"
    }
  }
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

func testAccResourceSecurityUtmProfileWebFEConfigUpdate() string {
	return `
resource "junos_security_utm_profile_web_filtering_juniper_enhanced" "testacc_ProfileWebFE" {
  name = "testacc ProfileWebFE"
  category {
    name   = "Enhanced_Network_Errors"
    action = "block"
  }
  custom_block_message      = "Blocked by Juniper"
  default_action            = "log-and-permit"
  no_safe_search            = true
  quarantine_custom_message = "Quarantine by Juniper"
  quarantine_message {
    url                      = "quarantine.local"
    type_custom_redirect_url = true
  }
  site_reputation_action {
    site_reputation = "harmful"
    action          = "block"
  }
  timeout = 3
}
`
}
