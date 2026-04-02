package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePolicyoptionsASPathGroup_basic(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.junos_policyoptions_as_path_group.testacc_dataASPathGroup",
						"id", "testacc_dataASPathGroup"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_as_path_group.testacc_dataASPathGroup",
						"as_path.#", "1"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_as_path_group.testacc_dataASPathGroup",
						"as_path.0.name", "testacc_dataASPathGroup"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_as_path_group.testacc_dataASPathGroup",
						"as_path.0.path", "5|12|18"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile("policy-options as-path-group .* doesn't exist"),
			},
		},
		PreventPostDestroyRefresh: true,
	})
}
