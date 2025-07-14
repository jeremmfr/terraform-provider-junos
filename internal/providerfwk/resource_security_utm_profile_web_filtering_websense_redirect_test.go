package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityUtmProfileWebFilteringWebsenseRedirect_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.default", "log-and-permit"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.server_connectivity", "log-and-permit"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.timeout", "log-and-permit"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.timeout", "log-and-permit"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"timeout", "3"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"server.host", "10.0.0.1"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"server.port", "1024"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"sockets", "16"),
					),
				},
				{
					ResourceName:      "junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
