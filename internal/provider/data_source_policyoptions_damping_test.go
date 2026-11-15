package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePolicyoptionsDamping_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.junos_policyoptions_damping.testacc_dataDamping",
						"id", "testacc_dataDamping"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_damping.testacc_dataDamping",
						"half_life", "15"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_damping.testacc_dataDamping",
						"max_suppress", "60"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_damping.testacc_dataDamping",
						"reuse", "750"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_damping.testacc_dataDamping",
						"suppress", "3000"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile("policy-options damping .* doesn't exist"),
			},
		},
		PreventPostDestroyRefresh: true,
	})
}
