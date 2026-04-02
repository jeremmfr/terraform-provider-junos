package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePolicyoptionsASPath_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.junos_policyoptions_as_path.testacc_dataASPath",
						"id", "testacc_dataASPath"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_as_path.testacc_dataASPath",
						"path", "5|12|18"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile("policy-options as-path .* doesn't exist"),
			},
		},
		PreventPostDestroyRefresh: true,
	})
}
