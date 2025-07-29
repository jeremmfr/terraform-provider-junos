package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityUtmPolicy_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					// 1
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"anti_virus.http_profile", "junos-sophos-av-defaults"),
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"traffic_sessions_per_client.over_limit", "log-and-permit"),
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"web_filtering_profile", "junos-wf-local-default"),
					),
				},
				{
					// 2
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_policy.testacc_Policy",
							"web_filtering_profile", "junos-wf-enhanced-default"),
					),
				},
				{
					// 3
					ResourceName:      "junos_security_utm_policy.testacc_Policy",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					// 4
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 5
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
