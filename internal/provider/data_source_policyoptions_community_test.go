package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePolicyoptionsCommunity_basic(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.junos_policyoptions_community.testacc_dataCommunity",
						"id", "testacc_dataCommunity"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_community.testacc_dataCommunity",
						"members.#", "1"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_community.testacc_dataCommunity",
						"members.0", "65000:100"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile("policy-options community .* doesn't exist"),
			},
		},
		PreventPostDestroyRefresh: true,
	})
}
