package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityUtmProfileWebFWebS_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" && os.Getenv("TESTACC_ROUTER") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityUtmProfileWebFWebSConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.#", "1"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.0.default", "log-and-permit"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.0.server_connectivity", "log-and-permit"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.0.timeout", "log-and-permit"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.0.timeout", "log-and-permit"),
					),
				},
				{
					Config: testAccJunosSecurityUtmProfileWebFWebSConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"custom_block_message", "Blocked by Juniper"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"fallback_settings.#", "0"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"timeout", "3"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"server.#", "1"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"server.0.host", "10.0.0.1"),
						resource.TestCheckResourceAttr(
							"junos_security_utm_profile_web_filtering_websense_redirect.testacc_ProfileWebFWebS",
							"server.0.port", "1024"),
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

func testAccJunosSecurityUtmProfileWebFWebSConfigCreate() string {
	return `
resource junos_security_utm_profile_web_filtering_websense_redirect "testacc_ProfileWebFWebS" {
  name                 = "testacc ProfileWebFWebS"
  custom_block_message = "Blocked by Juniper"
  fallback_settings {
    default             = "log-and-permit"
    server_connectivity = "log-and-permit"
    timeout             = "log-and-permit"
  }
}
`
}

func testAccJunosSecurityUtmProfileWebFWebSConfigUpdate() string {
	return `
resource junos_security_utm_profile_web_filtering_websense_redirect "testacc_ProfileWebFWebS" {
  name                 = "testacc ProfileWebFWebS"
  custom_block_message = "Blocked by Juniper"
  timeout              = 3
  server {
    host = "10.0.0.1"
    port = 1024
  }
  sockets = 16
}
`
}
