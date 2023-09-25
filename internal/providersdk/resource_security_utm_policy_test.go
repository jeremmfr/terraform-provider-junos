package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityUtmPolicy_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecurityUtmPolicyConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"anti_virus.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"anti_virus.0.http_profile", "junos-sophos-av-defaults"),
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"traffic_sessions_per_client.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"traffic_sessions_per_client.0.over_limit", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"web_filtering_profile", "junos-wf-local-default"),
					),
				},
				{
					Config: testAccResourceSecurityUtmPolicyConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"anti_virus.#", "0"),
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"traffic_sessions_per_client.#", "0"),
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"web_filtering_profile", "junos-wf-enhanced-default"),
					),
				},
				{
					ResourceName:      "junos_security_utm_policy.testacc_Policy",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceSecurityUtmPolicyConfigCreate() string {
	return `
resource "junos_security_utm_policy" "testacc_Policy" {
  name = "testacc Policy"
  anti_virus {
    http_profile = "junos-sophos-av-defaults"
  }
  traffic_sessions_per_client {
    over_limit = "log-and-permit"
  }
  web_filtering_profile = "junos-wf-local-default"
}
`
}

func testAccResourceSecurityUtmPolicyConfigUpdate() string {
	return `
resource "junos_security_utm_policy" "testacc_Policy" {
  name                  = "testacc Policy"
  web_filtering_profile = "junos-wf-enhanced-default"
}
`
}
