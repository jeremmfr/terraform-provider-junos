package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityUtmPolicy_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityUtmPolicyConfigCreate(),
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
					Config: testAccJunosSecurityUtmPolicyConfigUpdate(),
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

func testAccJunosSecurityUtmPolicyConfigCreate() string {
	return `
resource junos_security_utm_policy "testacc_Policy" {
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

func testAccJunosSecurityUtmPolicyConfigUpdate() string {
	return `
resource junos_security_utm_policy "testacc_Policy" {
  name                  = "testacc Policy"
  web_filtering_profile = "junos-wf-enhanced-default"
}
`
}
