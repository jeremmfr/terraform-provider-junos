package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceServices_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("junos_services.testacc",
							"security_intelligence.authentication_token"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("junos_services.testacc",
							"security_intelligence.authentication_token"),
						resource.TestCheckResourceAttr("junos_services.testacc",
							"security_intelligence.default_policy.#", "1"),
					),
				},
				{
					ResourceName:      "junos_services_proxy_profile.testacc_services",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_services.testacc",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
